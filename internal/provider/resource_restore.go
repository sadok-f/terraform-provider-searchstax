package provider

import (
	"context"
	"fmt"
	"strings"

	searchstaxClient "terraform-provider-searchstax/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func NewRestoreResource() resource.Resource { return &restoreResource{} }

type restoreResource struct{ client *searchstaxClient.Client }

func (r *restoreResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_restore"
}

func (r *restoreResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{Computed: true},
		"account_name": schema.StringAttribute{Required: true},
		"backup_id": schema.StringAttribute{
			Required: true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
		"deployment_uid": schema.StringAttribute{
			Optional: true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
		"restore_id": schema.StringAttribute{Computed: true},
		"status":     schema.StringAttribute{Computed: true},
	}}
}

func (r *restoreResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *restoreResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan restoreResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	reqBody := searchstaxClient.RestoreRequest{BackupID: plan.BackupID.ValueString()}
	var out *searchstaxClient.RestoreResponse
	var err error
	if plan.DeploymentUID.IsNull() {
		out, err = r.client.CreateAccountRestore(plan.AccountName.ValueString(), reqBody)
	} else {
		out, err = r.client.CreateDeploymentRestore(plan.AccountName.ValueString(), plan.DeploymentUID.ValueString(), reqBody)
	}
	if err != nil {
		resp.Diagnostics.AddError("Error creating restore", err.Error())
		return
	}
	plan.RestoreID = types.StringValue(out.RestoreID)
	if out.Status != "" {
		plan.Status = types.StringValue(out.Status)
	} else if !plan.DeploymentUID.IsNull() {
		status, err := r.client.GetDeploymentRestoreStatus(plan.AccountName.ValueString(), plan.DeploymentUID.ValueString(), reqBody)
		if err == nil && status.Status != "" {
			plan.Status = types.StringValue(status.Status)
		}
	}
	if plan.DeploymentUID.IsNull() {
		plan.ID = types.StringValue(plan.AccountName.ValueString() + "/" + out.RestoreID)
	} else {
		plan.ID = types.StringValue(plan.AccountName.ValueString() + "/" + plan.DeploymentUID.ValueString() + "/" + out.RestoreID)
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *restoreResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state restoreResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	reqBody := searchstaxClient.RestoreRequest{BackupID: state.BackupID.ValueString()}
	var out *searchstaxClient.RestoreResponse
	var err error
	if state.DeploymentUID.IsNull() {
		resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
		return
	}
	out, err = r.client.GetDeploymentRestoreStatus(state.AccountName.ValueString(), state.DeploymentUID.ValueString(), reqBody)
	if err != nil {
		resp.Diagnostics.AddError("Error reading restore status", err.Error())
		return
	}
	state.Status = types.StringValue(out.Status)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *restoreResource) Update(context.Context, resource.UpdateRequest, *resource.UpdateResponse) {}

func (r *restoreResource) Delete(context.Context, resource.DeleteRequest, *resource.DeleteResponse) {}

func (r *restoreResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, "/")
	if len(parts) == 2 {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("account_name"), parts[0])...)
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("restore_id"), parts[1])...)
		return
	}
	if len(parts) == 3 {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("account_name"), parts[0])...)
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("deployment_uid"), parts[1])...)
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("restore_id"), parts[2])...)
		return
	}
	resp.Diagnostics.AddError("Unexpected Import Identifier", "Expected account_name/restore_id or account_name/deployment_uid/restore_id")
}

type restoreResourceModel struct {
	ID            types.String `tfsdk:"id"`
	AccountName   types.String `tfsdk:"account_name"`
	DeploymentUID types.String `tfsdk:"deployment_uid"`
	BackupID      types.String `tfsdk:"backup_id"`
	RestoreID     types.String `tfsdk:"restore_id"`
	Status        types.String `tfsdk:"status"`
}
