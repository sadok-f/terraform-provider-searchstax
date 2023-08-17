package provider

import (
	"context"
	"fmt"

	searchstaxClient "terraform-provider-searchstax/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &deploymentsDataSource{}
	_ datasource.DataSourceWithConfigure = &deploymentsDataSource{}
)

// NewDeploymentsDataSource is a helper function to simplify the provider implementation.
func NewDeploymentsDataSource() datasource.DataSource {
	return &deploymentsDataSource{}
}

// deploymentsDataSource is the data source implementation.
type deploymentsDataSource struct {
	client *searchstaxClient.Client
}

// Metadata returns the data source type name.
func (d *deploymentsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_deployments"
}

// Schema - defines the schema for the data source.go.
func (d *deploymentsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"deployments_list": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"uid": schema.StringAttribute{
							Computed: true,
						},
						"name": schema.StringAttribute{
							Computed: true,
						},
						"application": schema.StringAttribute{
							Computed: true,
						},
						"application_version": schema.StringAttribute{
							Computed: true,
						},
						"tier": schema.StringAttribute{
							Computed: true,
						},
						"http_endpoint": schema.StringAttribute{
							Computed: true,
						},
						"status": schema.StringAttribute{
							Computed: true,
						},
						"provision_state": schema.StringAttribute{
							Computed: true,
						},
						"termination_lock": schema.BoolAttribute{
							Computed: true,
						},
						"plan_type": schema.StringAttribute{
							Computed: true,
						},
						"plan": schema.StringAttribute{
							Computed: true,
						},
						"is_master_slave": schema.BoolAttribute{
							Computed: true,
						},
						"vpc_type": schema.StringAttribute{
							Computed: true,
						},
						"vpc_name": schema.StringAttribute{
							Computed: true,
						},
						"region_id": schema.StringAttribute{
							Computed: true,
						},
						"cloud_provider": schema.StringAttribute{
							Computed: true,
						},
						"cloud_provider_id": schema.StringAttribute{
							Computed: true,
						},
						"num_additional_app_nodes": schema.NumberAttribute{
							Computed: true,
						},
						"deployment_type": schema.StringAttribute{
							Computed: true,
						},
						"num_nodes_default": schema.NumberAttribute{
							Computed: true,
						},
						//"servers": schema.MapAttribute{
						//	ElementType: types.StringType,
						//	Computed:    true,
						//},

						"date_created": schema.StringAttribute{
							Computed: true,
						},
						//TODO to list all missing attributes
					},
				},
			},
			"account_name": schema.StringAttribute{
				Required: true,
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *deploymentsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config deploymentsDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	deployments, err := d.client.GetDeployments(config.AccountName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read SearchStax Deployments",
			err.Error(),
		)
		return
	}
	config.ID = types.StringValue("placeholder")

	// Map response body to model
	for _, deployment := range deployments.Results {
		deploymentState := deploymentsModel{
			UID:                   types.StringValue(deployment.UID),
			Name:                  types.StringValue(deployment.Name),
			Application:           types.StringValue(deployment.Application),
			ApplicationVersion:    types.StringValue(deployment.ApplicationVersion),
			Tier:                  types.StringValue(deployment.Tier),
			Plan:                  types.StringValue(deployment.Plan),
			PlanType:              types.StringValue(deployment.PlanType),
			HttpEndpoint:          types.StringValue(deployment.HttpEndpoint),
			Status:                types.StringValue(deployment.Status),
			DateCreated:           types.StringValue(deployment.DateCreated),
			CloudProvider:         types.StringValue(deployment.CloudProvider),
			CloudProviderId:       types.StringValue(deployment.CloudProviderId),
			DeploymentType:        types.StringValue(deployment.DeploymentType),
			ProvisionState:        types.StringValue(deployment.ProvisionState),
			RegionId:              types.StringValue(deployment.RegionId),
			VpcType:               types.StringValue(deployment.VpcType),
			VpcName:               types.StringValue(deployment.VpcName),
			TerminationLock:       types.BoolValue(deployment.TerminationLock),
			IsMasterSlave:         types.BoolValue(deployment.IsMasterSlave),
			NumNodesDefault:       types.Int64Value(deployment.NumNodesDefault),
			NumAdditionalAppNodes: types.Int64Value(deployment.NumAdditionalAppNodes),
			//TODO to list all missing attributes
		}

		config.Deployments = append(config.Deployments, deploymentState)
	}

	// Set state
	diags := resp.State.Set(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the data source.
func (d *deploymentsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// deploymentsDataSourceModel  maps the data source schema data.
type deploymentsDataSourceModel struct {
	ID          types.String       `tfsdk:"id"`
	Deployments []deploymentsModel `tfsdk:"deployments_list"`
	AccountName types.String       `tfsdk:"account_name"`
}

// deploymentsModel maps deployments_list schema data.
type deploymentsModel struct {
	UID                   types.String `tfsdk:"uid"`
	Name                  types.String `tfsdk:"name"`
	Application           types.String `tfsdk:"application"`
	ApplicationVersion    types.String `tfsdk:"application_version"`
	Tier                  types.String `tfsdk:"tier"`
	HttpEndpoint          types.String `tfsdk:"http_endpoint"`
	Status                types.String `tfsdk:"status"`
	ProvisionState        types.String `tfsdk:"provision_state"`
	TerminationLock       types.Bool   `tfsdk:"termination_lock"`
	Plan                  types.String `tfsdk:"plan"`
	PlanType              types.String `tfsdk:"plan_type"`
	IsMasterSlave         types.Bool   `tfsdk:"is_master_slave"`
	VpcType               types.String `tfsdk:"vpc_type"`
	VpcName               types.String `tfsdk:"vpc_name"`
	RegionId              types.String `tfsdk:"region_id"`
	CloudProvider         types.String `tfsdk:"cloud_provider"`
	CloudProviderId       types.String `tfsdk:"cloud_provider_id"`
	NumAdditionalAppNodes types.Int64  `tfsdk:"num_additional_app_nodes"`
	DeploymentType        types.String `tfsdk:"deployment_type"`
	NumNodesDefault       types.Int64  `tfsdk:"num_nodes_default"`
	//Servers               types.List   `tfsdk:"servers"`
	DateCreated types.String `tfsdk:"date_created"`
	//TODO to list all missing attributes
}
