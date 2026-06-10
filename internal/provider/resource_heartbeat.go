package provider

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	searchstaxClient "terraform-provider-searchstax/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func NewHeartbeatResource() resource.Resource { return &heartbeatResource{} }

type heartbeatResource struct{ client *searchstaxClient.Client }

func (r *heartbeatResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_heartbeat"
}

func (r *heartbeatResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{Attributes: map[string]schema.Attribute{
		"id":               schema.StringAttribute{Computed: true},
		"account_name":     schema.StringAttribute{Required: true},
		"deployment_uid":   schema.StringAttribute{Required: true},
		"heartbeat_id":     schema.Int64Attribute{Computed: true},
		"name":             schema.StringAttribute{Required: true},
		"host":             schema.StringAttribute{Required: true},
		"interval":         schema.StringAttribute{Required: true},
		"max_alerts":       schema.StringAttribute{Required: true},
		"email":            schema.ListAttribute{Optional: true, ElementType: types.StringType},
		"webhook_trigger":  schema.Int64Attribute{Optional: true},
		"webhook_resolve":  schema.Int64Attribute{Optional: true},
	}}
}

func (r *heartbeatResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *heartbeatResource) heartbeatFromPlan(ctx context.Context, plan heartbeatResourceModel) searchstaxClient.Heartbeat {
	hb := searchstaxClient.Heartbeat{
		Name:      plan.Name.ValueString(),
		Host:      plan.Host.ValueString(),
		Interval:  plan.Interval.ValueString(),
		MaxAlerts: plan.MaxAlerts.ValueString(),
	}
	if !plan.Email.IsNull() {
		var emails []string
		_ = plan.Email.ElementsAs(ctx, &emails, false)
		hb.Email = emails
	}
	if !plan.WebhookTrigger.IsNull() && plan.WebhookTrigger.ValueInt64() != 0 {
		hb.WebhookTrigger = plan.WebhookTrigger.ValueInt64()
	}
	if !plan.WebhookResolve.IsNull() && plan.WebhookResolve.ValueInt64() != 0 {
		hb.WebhookResolve = plan.WebhookResolve.ValueInt64()
	}
	return hb
}

func (r *heartbeatResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan heartbeatResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	id, err := r.client.CreateHeartbeat(plan.AccountName.ValueString(), plan.DeploymentUID.ValueString(), r.heartbeatFromPlan(ctx, plan))
	if err != nil {
		resp.Diagnostics.AddError("Error creating heartbeat", err.Error())
		return
	}
	plan.HeartbeatID = types.Int64Value(id)
	plan.ID = types.StringValue(plan.AccountName.ValueString() + "/" + plan.DeploymentUID.ValueString() + "/" + strconv.FormatInt(id, 10))
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *heartbeatResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state heartbeatResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	hb, err := r.client.GetHeartbeat(state.AccountName.ValueString(), state.DeploymentUID.ValueString(), state.HeartbeatID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("Error reading heartbeat", err.Error())
		return
	}
	if hb.Name != "" {
		state.Name = types.StringValue(hb.Name)
	}
	if hb.Host != "" {
		state.Host = types.StringValue(hb.Host)
	}
	if hb.Interval != "" {
		state.Interval = types.StringValue(hb.Interval)
	}
	if hb.MaxAlerts != "" {
		state.MaxAlerts = types.StringValue(hb.MaxAlerts)
	}
	if len(hb.Email) > 0 {
		emails, diags := types.ListValueFrom(ctx, types.StringType, hb.Email)
		resp.Diagnostics.Append(diags...)
		state.Email = emails
	}
	if hb.WebhookTrigger != 0 {
		state.WebhookTrigger = types.Int64Value(hb.WebhookTrigger)
	}
	if hb.WebhookResolve != 0 {
		state.WebhookResolve = types.Int64Value(hb.WebhookResolve)
	}
	state.ID = types.StringValue(state.AccountName.ValueString() + "/" + state.DeploymentUID.ValueString() + "/" + strconv.FormatInt(hb.ID, 10))
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *heartbeatResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan heartbeatResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.UpdateHeartbeat(plan.AccountName.ValueString(), plan.DeploymentUID.ValueString(), plan.HeartbeatID.ValueInt64(), r.heartbeatFromPlan(ctx, plan)); err != nil {
		resp.Diagnostics.AddError("Error updating heartbeat", err.Error())
		return
	}
	plan.ID = types.StringValue(plan.AccountName.ValueString() + "/" + plan.DeploymentUID.ValueString() + "/" + strconv.FormatInt(plan.HeartbeatID.ValueInt64(), 10))
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *heartbeatResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state heartbeatResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.DeleteHeartbeat(state.AccountName.ValueString(), state.DeploymentUID.ValueString(), state.HeartbeatID.ValueInt64()); err != nil {
		resp.Diagnostics.AddError("Error deleting heartbeat", err.Error())
	}
}

func (r *heartbeatResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, "/")
	if len(parts) != 3 {
		resp.Diagnostics.AddError("Unexpected Import Identifier", "Expected account_name/deployment_uid/heartbeat_id")
		return
	}
	id, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Unexpected Import Identifier", "Invalid heartbeat_id")
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("account_name"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("deployment_uid"), parts[1])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("heartbeat_id"), id)...)
}

type heartbeatResourceModel struct {
	ID             types.String `tfsdk:"id"`
	AccountName    types.String `tfsdk:"account_name"`
	DeploymentUID  types.String `tfsdk:"deployment_uid"`
	HeartbeatID    types.Int64  `tfsdk:"heartbeat_id"`
	Name           types.String `tfsdk:"name"`
	Host           types.String `tfsdk:"host"`
	Interval       types.String `tfsdk:"interval"`
	MaxAlerts      types.String `tfsdk:"max_alerts"`
	Email          types.List   `tfsdk:"email"`
	WebhookTrigger types.Int64  `tfsdk:"webhook_trigger"`
	WebhookResolve types.Int64  `tfsdk:"webhook_resolve"`
}
