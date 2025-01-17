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
	resp.TypeName = req.ProviderTypeName + "_deployment"
}

// Schema defines the schema for the resource.
func (d *deploymentUserResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			// id is required by the testing framework
			"id": schema.StringAttribute{
				Computed: true,
			},
			"uid": schema.StringAttribute{
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
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"role": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
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

	// Create new deployment
	var deploymentUser, err = d.client.CreateDeploymentUser(item, plan.AccountName.ValueString(), plan.UID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating deployment",
			"Could not create deployment, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan.ID = types.StringValue("placeholder")
	plan.UID = types.StringValue(deploymentUser.UID)
	plan.Username = types.StringValue(deploymentUser.Username)
	plan.Password = types.StringValue(deploymentUser.Password)
	plan.Role = types.StringValue(deploymentUser.Role)

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

	// Get refreshed deployment value from SearchStax
	var deployment, err = d.client.GetDeploymentUser(state.AccountName.ValueString(), state.UID.ValueString(), state.Username.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading SearchStax deployment",
			fmt.Sprintf("Could not read SearchStax deployment UID: %s, Account: %s, Error: %s ", state.UID.ValueString(), state.AccountName.ValueString(), err.Error()),
		)
		return
	}

	// Overwrite items with refreshed state
	state.ID = types.StringValue("placeholder")
	state.UID = types.StringValue(deployment.UID)
	state.Username = types.StringValue(deployment.Username)
	state.Password = types.StringValue(deployment.Password)
	state.Role = types.StringValue(deployment.Role)

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
	var plan deploymentModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// if we're changing only the PrivateVPC, then update the state and return
	// this is a WORKAROUND until the API returns a private_vpc id as well
	if plan.PrivateVpc.ValueInt64() != 0 {
		plan.ID = types.StringValue("placeholder")
		diags = resp.State.Set(ctx, plan)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		return
	}

	// Generate API request body from plan
	var item = searchstaxClient.Deployment{
		UID:                   plan.UID.ValueString(),
		Name:                  plan.Name.ValueString(),
		Application:           plan.Application.ValueString(),
		ApplicationVersion:    plan.ApplicationVersion.ValueString(),
		TerminationLock:       plan.TerminationLock.ValueBool(),
		PlanType:              plan.PlanType.ValueString(),
		Plan:                  plan.Plan.ValueString(),
		RegionId:              plan.RegionId.ValueString(),
		CloudProviderId:       plan.CloudProviderId.ValueString(),
		NumAdditionalAppNodes: plan.NumAdditionalAppNodes.ValueInt64(),
		NumNodesDefault:       plan.NumNodesDefault.ValueInt64(),
		PrivateVpc:            plan.PrivateVpc.ValueInt64(),
		HttpEndpoint:          plan.HttpEndpoint.ValueString(),
		Status:                plan.Status.ValueString(),
		ProvisionState:        plan.ProvisionState.ValueString(),
		Tier:                  plan.Tier.ValueString(),
		DateCreated:           plan.DateCreated.ValueString(),
		DeploymentType:        plan.DeploymentType.ValueString(),
	}

	// Update existing deployment (recreate the cluster)
	uid := ""
	req.State.GetAttribute(ctx, path.Root("uid"), &uid)
	deployment, err := d.client.UpdateDeployment(plan.AccountName.ValueString(), uid, item)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating SearchStax Deployment",
			"Could not update Deployment, unexpected error: "+err.Error(),
		)
		return
	}

	// Overwrite items with refreshed state
	plan.ID = types.StringValue("placeholder")
	plan.UID = types.StringValue(deployment.UID)
	plan.CloudProvider = types.StringValue(deployment.CloudProvider)
	plan.DateCreated = types.StringValue(deployment.DateCreated)
	plan.DeploymentType = types.StringValue(deployment.DeploymentType)
	plan.HttpEndpoint = types.StringValue(deployment.HttpEndpoint)
	plan.IsMasterSlave = types.BoolValue(deployment.IsMasterSlave)
	plan.NumNodesDefault = types.Int64Value(deployment.NumNodesDefault)
	plan.NumAdditionalAppNodes = types.Int64Value(deployment.NumAdditionalAppNodes)
	plan.ProvisionState = types.StringValue(deployment.ProvisionState)
	plan.Status = types.StringValue(deployment.Status)
	plan.Tier = types.StringValue(deployment.Tier)
	plan.VpcName = types.StringValue(deployment.VpcName)
	plan.VpcType = types.StringValue(deployment.VpcType)
	plan.VpcType = types.StringValue(deployment.VpcType)
	//TODO list the rest of attributes

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (d *deploymentUserResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state deploymentModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing deployment
	err := d.client.DeleteDeployment(state.AccountName.ValueString(), state.UID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting SearchStax Deployment",
			"Could not delete deployment, unexpected error: "+err.Error(),
		)
		return
	}
}

// ImportState - Import existing deployment cluster into terraform state.
func (d *deploymentUserResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	//idParts := strings.Split(req.ID, "/")
	//
	//if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
	//	resp.Diagnostics.AddError(
	//		"Unexpected Import Identifier",
	//		fmt.Sprintf("Expected import identifier with format: account_name/uid. Got: %q", req.ID),
	//	)
	//	return
	//}
	//
	//resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("account_name"), idParts[0])...)
	//resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("uid"), idParts[1])...)
}

// deploymentUserModel maps deployment schema data.
type deploymentUserModel struct {
	ID          types.String `tfsdk:"id"`
	AccountName types.String `tfsdk:"account_name"`
	UID         types.String `tfsdk:"uid"`
	Username    types.String `tfsdk:"username"`
	Password    types.String `tfsdk:"password"`
	Role        types.String `tfsdk:"role"`
}
