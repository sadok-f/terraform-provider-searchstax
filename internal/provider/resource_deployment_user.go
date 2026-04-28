package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strings"
	searchstaxClient "terraform-provider-searchstax/internal/client"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.ResourceWithImportState = &deploymentUserResource{}
)

// NewDeploymentUserResource is a helper function to simplify the provider implementation.
func NewDeploymentUserResource() resource.Resource {
	return &deploymentUserResource{}
}

// deploymentUserResource is the resource implementation.
type deploymentUserResource struct {
	client *searchstaxClient.Client
}

// Configure adds the provider configured client to the resource.
func (d *deploymentUserResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*searchstaxClient.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *searchstaxClient.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

// Metadata returns the resource type name.
func (d *deploymentUserResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_deployment_user"
}

// Schema defines the schema for the resource.
func (d *deploymentUserResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			// id is required by the testing framework
			"id": schema.StringAttribute{
				Computed: true,
			},
			"account_name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"deployment_uid": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"username": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"password": schema.StringAttribute{
				Required:      true,
				Sensitive:     true,
				PlanModifiers: []planmodifier.String{
					// password rotation is modeled as update via delete+add in the client.
					// keep it updatable without forcing a replace in state.
				},
			},
			"role": schema.StringAttribute{
				Required: true,
			},
		},
	}
}

// Create a new resource.
func (d *deploymentUserResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan deploymentUserModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	var item = searchstaxClient.DeploymentUser{
		Username: plan.Username.ValueString(),
		Password: plan.Password.ValueString(),
		Role:     plan.Role.ValueString(),
	}

	// Create new basic-auth user
	var _, err = d.client.CreateDeploymentUser(item, plan.AccountName.ValueString(), plan.DeploymentUID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating deployment user",
			"Could not create deployment user, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan.ID = types.StringValue(fmt.Sprintf("%s/%s/%s", plan.AccountName.ValueString(), plan.DeploymentUID.ValueString(), plan.Username.ValueString()))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information.
func (d *deploymentUserResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state deploymentUserModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed user value from SearchStax
	var user, err = d.client.GetDeploymentUser(state.AccountName.ValueString(), state.DeploymentUID.ValueString(), state.Username.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading SearchStax deployment user",
			fmt.Sprintf("Could not read SearchStax deployment user %q (deployment %s, account %s): %s", state.Username.ValueString(), state.DeploymentUID.ValueString(), state.AccountName.ValueString(), err.Error()),
		)
		return
	}

	// Overwrite items with refreshed state
	state.ID = types.StringValue(fmt.Sprintf("%s/%s/%s", state.AccountName.ValueString(), state.DeploymentUID.ValueString(), state.Username.ValueString()))
	state.Username = types.StringValue(user.Username)
	state.Role = types.StringValue(user.Role)
	// keep password from state (API does not return it)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (d *deploymentUserResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan deploymentUserModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	var item = searchstaxClient.DeploymentUser{
		Username: plan.Username.ValueString(),
		Password: plan.Password.ValueString(),
		Role:     plan.Role.ValueString(),
	}

	_, err := d.client.UpdateDeploymentUser(plan.AccountName.ValueString(), plan.DeploymentUID.ValueString(), item)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating SearchStax Deployment User",
			"Could not update deployment user, unexpected error: "+err.Error(),
		)
		return
	}

	plan.ID = types.StringValue(fmt.Sprintf("%s/%s/%s", plan.AccountName.ValueString(), plan.DeploymentUID.ValueString(), plan.Username.ValueString()))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (d *deploymentUserResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state deploymentUserModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete user
	err := d.client.DeleteDeploymentUser(state.AccountName.ValueString(), state.DeploymentUID.ValueString(), state.Username.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting SearchStax Deployment User",
			"Could not delete deployment user, unexpected error: "+err.Error(),
		)
		return
	}
}

// ImportState - Import existing deployment cluster into terraform state.
func (d *deploymentUserResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, "/")

	if len(idParts) != 3 || idParts[0] == "" || idParts[1] == "" || idParts[2] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: account_name/deployment_uid/username. Got: %q", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("account_name"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("deployment_uid"), idParts[1])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("username"), idParts[2])...)
}

// deploymentUserModel maps deployment schema data.
type deploymentUserModel struct {
	ID            types.String `tfsdk:"id"`
	AccountName   types.String `tfsdk:"account_name"`
	DeploymentUID types.String `tfsdk:"deployment_uid"`
	Username      types.String `tfsdk:"username"`
	Password      types.String `tfsdk:"password"`
	Role          types.String `tfsdk:"role"`
}
