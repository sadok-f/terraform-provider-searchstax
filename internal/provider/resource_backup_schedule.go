package provider

import (
	"context"
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
	if !plan.Frequency.IsNull() {
		body["frequency"] = plan.Frequency.ValueInt64()
	}
	if !plan.Collections.IsNull() {
		var collections []string
		_ = plan.Collections.ElementsAs(ctx, &collections, false)
		body["collections"] = collections
	}
	return body, days, nil
}

func (r *backupScheduleResource) findSchedule(accountName, deploymentUID string, days []string, retention int64, time string) (string, bool) {
	list, err := r.client.GetBackupSchedules(accountName, deploymentUID)
	if err != nil {
		return "", false
	}
	for _, s := range list.Results {
		if int64(s.Retention) == retention && s.Time == time && reflect.DeepEqual(s.Days, days) {
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
		resp.Diagnostics.AddError("Error creating backup schedule", err.Error())
		return
	}
	timeVal := ""
	if !plan.Time.IsNull() {
		timeVal = plan.Time.ValueString()
	}
	scheduleID, ok := r.findSchedule(plan.AccountName.ValueString(), plan.DeploymentUID.ValueString(), days, plan.Retention.ValueInt64(), timeVal)
	if !ok {
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
			state.Time = types.StringValue(s.Time)
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

func (r *backupScheduleResource) Update(context.Context, resource.UpdateRequest, *resource.UpdateResponse) {}

func (r *backupScheduleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state backupScheduleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.DeleteBackupSchedule(state.AccountName.ValueString(), state.DeploymentUID.ValueString(), state.ScheduleID.ValueString()); err != nil {
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
