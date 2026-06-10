package provider

import (
	"context"
	"fmt"

	searchstaxClient "terraform-provider-searchstax/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func NewBasicAuthResource() resource.Resource { return &basicAuthResource{} }

type basicAuthResource struct{ client *searchstaxClient.Client }

func (r *basicAuthResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_basic_auth"
}

func (r *basicAuthResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{Attributes: map[string]schema.Attribute{
		"id":             schema.StringAttribute{Computed: true},
		"account_name":   schema.StringAttribute{Required: true},
		"deployment_uid": schema.StringAttribute{Required: true},
		"enabled": schema.BoolAttribute{
			Required: true,
			PlanModifiers: []planmodifier.Bool{
				boolplanmodifier.RequiresReplace(),
			},
		},
	}}
}

func (r *basicAuthResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *basicAuthResource) setEnabled(accountName, deploymentUID string, enabled bool) error {
	if enabled {
		_, err := r.client.EnableBasicAuth(accountName, deploymentUID)
		return err
	}
	_, err := r.client.DisableBasicAuth(accountName, deploymentUID)
	return err
}

func (r *basicAuthResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan basicAuthResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.setEnabled(plan.AccountName.ValueString(), plan.DeploymentUID.ValueString(), plan.Enabled.ValueBool()); err != nil {
		resp.Diagnostics.AddError("Error configuring basic auth", err.Error())
		return
	}
	plan.ID = types.StringValue(plan.AccountName.ValueString() + "/" + plan.DeploymentUID.ValueString())
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *basicAuthResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state basicAuthResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	enabled, err := r.client.IsBasicAuthEnabled(state.AccountName.ValueString(), state.DeploymentUID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading basic auth status", err.Error())
		return
	}
	state.Enabled = types.BoolValue(enabled)
	state.ID = types.StringValue(state.AccountName.ValueString() + "/" + state.DeploymentUID.ValueString())
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *basicAuthResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan basicAuthResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.setEnabled(plan.AccountName.ValueString(), plan.DeploymentUID.ValueString(), plan.Enabled.ValueBool()); err != nil {
		resp.Diagnostics.AddError("Error updating basic auth", err.Error())
		return
	}
	plan.ID = types.StringValue(plan.AccountName.ValueString() + "/" + plan.DeploymentUID.ValueString())
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *basicAuthResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state basicAuthResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if state.Enabled.ValueBool() {
		if _, err := r.client.DisableBasicAuth(state.AccountName.ValueString(), state.DeploymentUID.ValueString()); err != nil {
			resp.Diagnostics.AddError("Error disabling basic auth", err.Error())
		}
	}
}

type basicAuthResourceModel struct {
	ID            types.String `tfsdk:"id"`
	AccountName   types.String `tfsdk:"account_name"`
	DeploymentUID types.String `tfsdk:"deployment_uid"`
	Enabled       types.Bool   `tfsdk:"enabled"`
}
