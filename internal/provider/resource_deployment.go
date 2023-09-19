package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strings"
	searchstaxClient "terraform-provider-searchstax/internal/client"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.ResourceWithImportState = &deploymentResource{}
)

// NewDeploymentResource is a helper function to simplify the provider implementation.
func NewDeploymentResource() resource.Resource {
	return &deploymentResource{}
}

// deploymentResource is the resource implementation.
type deploymentResource struct {
	client *searchstaxClient.Client
}

// Configure adds the provider configured client to the resource.
func (d *deploymentResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (d *deploymentResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_deployment"
}

// Schema defines the schema for the resource.
func (d *deploymentResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			"name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"application": schema.StringAttribute{
				Required: true,
			},
			"application_version": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"plan_type": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"region_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"plan": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"cloud_provider_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"termination_lock": schema.BoolAttribute{
				Required: true,
			},
			"private_vpc": schema.Int64Attribute{
				Optional: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"uid": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"tier": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"http_endpoint": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"status": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"provision_state": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"is_master_slave": schema.BoolAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"vpc_type": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"vpc_name": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"cloud_provider": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"num_additional_app_nodes": schema.Int64Attribute{
				Optional: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"deployment_type": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"num_nodes_default": schema.Int64Attribute{
				Computed: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			//"servers": schema.MapAttribute{
			//	ElementType: types.StringType,
			//	Computed:    true,
			//},

			"date_created": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			//TODO to list all missing attributes
		},
	}
}

// Create a new resource.
func (d *deploymentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan deploymentModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	var item = searchstaxClient.Deployment{
		Name:                  plan.Name.ValueString(),
		Application:           plan.Application.ValueString(),
		ApplicationVersion:    plan.ApplicationVersion.ValueString(),
		TerminationLock:       plan.TerminationLock.ValueBool(),
		PlanType:              plan.PlanType.ValueString(),
		Plan:                  plan.Plan.ValueString(),
		RegionId:              plan.RegionId.ValueString(),
		CloudProviderId:       plan.CloudProviderId.ValueString(),
		NumAdditionalAppNodes: plan.NumAdditionalAppNodes.ValueInt64(),
		PrivateVpc:            plan.PrivateVpc.ValueInt64(),
	}

	// Create new deployment
	var deployment, err = d.client.CreateDeployment(item, plan.AccountName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating deployment",
			"Could not create deployment, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan.ID = types.StringValue("placeholder")
	plan.UID = types.StringValue(deployment.UID)
	plan.CloudProvider = types.StringValue(deployment.CloudProvider)
	plan.DateCreated = types.StringValue(deployment.DateCreated)
	plan.DeploymentType = types.StringValue(deployment.DeploymentType)
	plan.HttpEndpoint = types.StringValue(deployment.HttpEndpoint)
	plan.IsMasterSlave = types.BoolValue(deployment.IsMasterSlave)
	plan.NumNodesDefault = types.Int64Value(deployment.NumNodesDefault)
	plan.ProvisionState = types.StringValue(deployment.ProvisionState)
	plan.Status = types.StringValue(deployment.Status)
	plan.Tier = types.StringValue(deployment.Tier)
	plan.VpcName = types.StringValue(deployment.VpcName)
	plan.VpcType = types.StringValue(deployment.VpcType)
	plan.RegionId = types.StringValue(deployment.RegionId)
	plan.NumAdditionalAppNodes = types.Int64Value(deployment.NumAdditionalAppNodes)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information.
func (d *deploymentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state deploymentModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed deployment value from SearchStax
	var deployment, err = d.client.GetDeployment(state.AccountName.ValueString(), state.UID.ValueString())
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
	state.Name = types.StringValue(deployment.Name)
	state.Application = types.StringValue(deployment.Application)
	state.ApplicationVersion = types.StringValue(deployment.ApplicationVersion)
	state.Tier = types.StringValue(deployment.Tier)
	state.HttpEndpoint = types.StringValue(deployment.HttpEndpoint)
	state.Status = types.StringValue(deployment.Status)
	state.ProvisionState = types.StringValue(deployment.ProvisionState)
	state.PlanType = types.StringValue(deployment.PlanType)
	state.Plan = types.StringValue(deployment.Plan)
	state.CloudProvider = types.StringValue(deployment.CloudProvider)
	state.CloudProviderId = types.StringValue(deployment.CloudProviderId)
	state.NumAdditionalAppNodes = types.Int64Value(deployment.NumAdditionalAppNodes)
	state.NumNodesDefault = types.Int64Value(deployment.NumNodesDefault)
	state.IsMasterSlave = types.BoolValue(deployment.IsMasterSlave)
	state.VpcName = types.StringValue(deployment.VpcName)
	state.VpcType = types.StringValue(deployment.VpcType)
	state.RegionId = types.StringValue(deployment.RegionId)
	state.TerminationLock = types.BoolValue(deployment.TerminationLock)
	state.DateCreated = types.StringValue(deployment.DateCreated)
	state.DeploymentType = types.StringValue(deployment.DeploymentType)
	//TODO list the rest of attributes

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (d *deploymentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
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

func (d *deploymentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
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
func (d *deploymentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, "/")

	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: account_name/uid. Got: %q", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("account_name"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("uid"), idParts[1])...)
}

// deploymentModel maps deployment schema data.
type deploymentModel struct {
	ID                    types.String `tfsdk:"id"`
	AccountName           types.String `tfsdk:"account_name"`
	UID                   types.String `tfsdk:"uid"`
	Name                  types.String `tfsdk:"name"`
	Application           types.String `tfsdk:"application"`
	ApplicationVersion    types.String `tfsdk:"application_version"`
	TerminationLock       types.Bool   `tfsdk:"termination_lock"`
	PlanType              types.String `tfsdk:"plan_type"`
	Plan                  types.String `tfsdk:"plan"`
	RegionId              types.String `tfsdk:"region_id"`
	CloudProvider         types.String `tfsdk:"cloud_provider"`
	CloudProviderId       types.String `tfsdk:"cloud_provider_id"`
	NumAdditionalAppNodes types.Int64  `tfsdk:"num_additional_app_nodes"`
	PrivateVpc            types.Int64  `tfsdk:"private_vpc"`
	Tier                  types.String `tfsdk:"tier"`
	HttpEndpoint          types.String `tfsdk:"http_endpoint"`
	ProvisionState        types.String `tfsdk:"provision_state"`
	Status                types.String `tfsdk:"status"`
	DateCreated           types.String `tfsdk:"date_created"`
	IsMasterSlave         types.Bool   `tfsdk:"is_master_slave"`
	VpcType               types.String `tfsdk:"vpc_type"`
	VpcName               types.String `tfsdk:"vpc_name"`
	DeploymentType        types.String `tfsdk:"deployment_type"`
	NumNodesDefault       types.Int64  `tfsdk:"num_nodes_default"`
}
