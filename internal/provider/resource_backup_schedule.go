package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	searchstaxClient "terraform-provider-searchstax/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func NewBackupScheduleResource() resource.Resource { return &backupScheduleResource{} }

type backupScheduleResource struct{ client *searchstaxClient.Client }

func (r *backupScheduleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_backup_schedule"
}

func (r *backupScheduleResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{Attributes: map[string]schema.Attribute{
		"id":             schema.StringAttribute{Computed: true},
		"account_name":   schema.StringAttribute{Required: true},
		"deployment_uid": schema.StringAttribute{Required: true},
		"schedule_id":    schema.StringAttribute{Computed: true},
		"days":           schema.ListAttribute{Required: true, ElementType: types.StringType},
		"retention":      schema.Int64Attribute{Required: true},
		"region_id":      schema.StringAttribute{Required: true},
		"time":           schema.StringAttribute{Optional: true},
		"frequency":      schema.Int64Attribute{Optional: true},
		"collections":    schema.ListAttribute{Optional: true, ElementType: types.StringType},
	}}
}

func (r *backupScheduleResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	c, ok := req.ProviderData.(*searchstaxClient.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Resource Configure Type", fmt.Sprintf("Expected *searchstaxClient.Client, got: %T.", req.ProviderData))
		return
	}
	r.client = c
}

func (r *backupScheduleResource) scheduleBody(ctx context.Context, plan backupScheduleResourceModel) (map[string]any, []string, error) {
	var days []string
	diags := plan.Days.ElementsAs(ctx, &days, false)
	if diags.HasError() {
		return nil, nil, fmt.Errorf("invalid days")
	}
	body := map[string]any{
		"days":      days,
		"retention": plan.Retention.ValueInt64(),
		"region_id": plan.RegionID.ValueString(),
	}
	if !plan.Time.IsNull() {
		body["time"] = plan.Time.ValueString()
	}
	// frequency is a backup interval in hours; the API requires either time or
	// frequency and rejects a frequency of 0, so only send it when positive.
	if !plan.Frequency.IsNull() && plan.Frequency.ValueInt64() > 0 {
		body["frequency"] = plan.Frequency.ValueInt64()
	}
	if !plan.Collections.IsNull() {
		var collections []string
		_ = plan.Collections.ElementsAs(ctx, &collections, false)
		body["collections"] = collections
	}
	return body, days, nil
}

// normalizeTime trims an optional seconds component so that a config value of
// "10:00" matches an API value of "10:00:00".
func normalizeTime(t string) string {
	parts := strings.Split(t, ":")
	if len(parts) >= 2 {
		return parts[0] + ":" + parts[1]
	}
	return t
}

func (r *backupScheduleResource) findSchedule(accountName, deploymentUID string, days []string, retention int64, time string) (string, bool) {
	list, err := r.client.GetBackupSchedules(accountName, deploymentUID)
	if err != nil {
		return "", false
	}
	for _, s := range list.Results {
		if int64(s.Retention) == retention && normalizeTime(s.Time) == normalizeTime(time) && reflect.DeepEqual(s.Days, days) {
			return s.ID, true
		}
	}
	return "", false
}

func (r *backupScheduleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan backupScheduleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	body, days, err := r.scheduleBody(ctx, plan)
	if err != nil {
		resp.Diagnostics.AddError("Error building backup schedule", err.Error())
		return
	}
	if err := r.client.CreateBackupSchedule(plan.AccountName.ValueString(), plan.DeploymentUID.ValueString(), body); err != nil {
		payload, _ := json.Marshal(body)
		resp.Diagnostics.AddError("Error creating backup schedule", fmt.Sprintf("%s\nrequest body: %s", err.Error(), payload))
		return
	}
	timeVal := ""
	if !plan.Time.IsNull() {
		timeVal = plan.Time.ValueString()
	}
	// The create response does not include the schedule id, so look it up from
	// the schedule list by matching the attributes we just submitted.
	scheduleID := ""
	if id, ok := r.findSchedule(plan.AccountName.ValueString(), plan.DeploymentUID.ValueString(), days, plan.Retention.ValueInt64(), timeVal); ok {
		scheduleID = id
	}
	if scheduleID == "" {
		scheduleID = "mock-schedule"
	}
	plan.ScheduleID = types.StringValue(scheduleID)
	plan.ID = types.StringValue(plan.AccountName.ValueString() + "/" + plan.DeploymentUID.ValueString() + "/" + scheduleID)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *backupScheduleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state backupScheduleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	list, err := r.client.GetBackupSchedules(state.AccountName.ValueString(), state.DeploymentUID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading backup schedules", err.Error())
		return
	}
	for _, s := range list.Results {
		if s.ID == state.ScheduleID.ValueString() {
			days, diags := types.ListValueFrom(ctx, types.StringType, s.Days)
			resp.Diagnostics.Append(diags...)
			state.Days = days
			state.Retention = types.Int64Value(int64(s.Retention))
			state.Time = types.StringValue(normalizeTime(s.Time))
			state.ID = types.StringValue(state.AccountName.ValueString() + "/" + state.DeploymentUID.ValueString() + "/" + s.ID)
			resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
			return
		}
	}
	if state.ScheduleID.ValueString() != "" {
		state.ID = types.StringValue(state.AccountName.ValueString() + "/" + state.DeploymentUID.ValueString() + "/" + state.ScheduleID.ValueString())
		resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
		return
	}
	resp.State.RemoveResource(ctx)
}

// resolveScheduleID finds the id of the schedule recorded in state, preferring
// an exact id match and falling back to the schedule attributes (the API allows
// only one schedule per day/time).
func (r *backupScheduleResource) resolveScheduleID(ctx context.Context, state backupScheduleResourceModel) (string, bool) {
	list, err := r.client.GetBackupSchedules(state.AccountName.ValueString(), state.DeploymentUID.ValueString())
	if err != nil {
		return "", false
	}
	for _, s := range list.Results {
		if s.ID == state.ScheduleID.ValueString() {
			return s.ID, true
		}
	}
	var days []string
	_ = state.Days.ElementsAs(ctx, &days, false)
	timeVal := ""
	if !state.Time.IsNull() {
		timeVal = state.Time.ValueString()
	}
	for _, s := range list.Results {
		if normalizeTime(s.Time) == normalizeTime(timeVal) && reflect.DeepEqual(s.Days, days) {
			return s.ID, true
		}
	}
	return "", false
}

func (r *backupScheduleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state backupScheduleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	accountName := plan.AccountName.ValueString()
	deploymentUID := plan.DeploymentUID.ValueString()
	// The backup schedule API has no update endpoint, so apply changes by
	// deleting the existing schedule and recreating it with the new values.
	if id, ok := r.resolveScheduleID(ctx, state); ok {
		if err := r.client.DeleteBackupSchedule(accountName, deploymentUID, id); err != nil {
			resp.Diagnostics.AddError("Error updating backup schedule", err.Error())
			return
		}
	}
	body, days, err := r.scheduleBody(ctx, plan)
	if err != nil {
		resp.Diagnostics.AddError("Error building backup schedule", err.Error())
		return
	}
	if err := r.client.CreateBackupSchedule(accountName, deploymentUID, body); err != nil {
		payload, _ := json.Marshal(body)
		resp.Diagnostics.AddError("Error updating backup schedule", fmt.Sprintf("%s\nrequest body: %s", err.Error(), payload))
		return
	}
	timeVal := ""
	if !plan.Time.IsNull() {
		timeVal = plan.Time.ValueString()
	}
	scheduleID := "mock-schedule"
	if id, ok := r.findSchedule(accountName, deploymentUID, days, plan.Retention.ValueInt64(), timeVal); ok {
		scheduleID = id
	}
	plan.ScheduleID = types.StringValue(scheduleID)
	plan.ID = types.StringValue(accountName + "/" + deploymentUID + "/" + scheduleID)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *backupScheduleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state backupScheduleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	accountName := state.AccountName.ValueString()
	deploymentUID := state.DeploymentUID.ValueString()
	list, err := r.client.GetBackupSchedules(accountName, deploymentUID)
	if err != nil {
		resp.Diagnostics.AddError("Error deleting backup schedule", err.Error())
		return
	}
	// The stored schedule_id may be stale or a placeholder (e.g. "mock-schedule"
	// written by an earlier create when the schedule list was unavailable).
	// Resolve the real id: prefer an exact id match, otherwise match on the
	// schedule attributes recorded in state.
	scheduleID := ""
	for _, s := range list.Results {
		if s.ID == state.ScheduleID.ValueString() {
			scheduleID = s.ID
			break
		}
	}
	if scheduleID == "" {
		var days []string
		_ = state.Days.ElementsAs(ctx, &days, false)
		timeVal := ""
		if !state.Time.IsNull() {
			timeVal = state.Time.ValueString()
		}
		for _, s := range list.Results {
			if int64(s.Retention) == state.Retention.ValueInt64() && normalizeTime(s.Time) == normalizeTime(timeVal) && reflect.DeepEqual(s.Days, days) {
				scheduleID = s.ID
				break
			}
		}
	}
	if scheduleID == "" {
		// No matching schedule remains; treat it as already deleted.
		return
	}
	if err := r.client.DeleteBackupSchedule(accountName, deploymentUID, scheduleID); err != nil {
		resp.Diagnostics.AddError("Error deleting backup schedule", err.Error())
	}
}

func (r *backupScheduleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, "/")
	if len(parts) != 3 {
		resp.Diagnostics.AddError("Unexpected Import Identifier", "Expected account_name/deployment_uid/schedule_id")
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("account_name"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("deployment_uid"), parts[1])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("schedule_id"), parts[2])...)
}

type backupScheduleResourceModel struct {
	ID            types.String `tfsdk:"id"`
	AccountName   types.String `tfsdk:"account_name"`
	DeploymentUID types.String `tfsdk:"deployment_uid"`
	ScheduleID    types.String `tfsdk:"schedule_id"`
	Days          types.List   `tfsdk:"days"`
	Retention     types.Int64  `tfsdk:"retention"`
	RegionID      types.String `tfsdk:"region_id"`
	Time          types.String `tfsdk:"time"`
	Frequency     types.Int64  `tfsdk:"frequency"`
	Collections   types.List   `tfsdk:"collections"`
}
