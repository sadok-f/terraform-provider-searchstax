package provider

import (
	"context"
	"fmt"

	searchstaxClient "terraform-provider-searchstax/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &deploymentsDataSource{}
	_ datasource.DataSourceWithConfigure = &deploymentsDataSource{}
)

func NewDeploymentsDataSource() datasource.DataSource {
	return &deploymentsDataSource{}
}

type deploymentsDataSource struct {
	client *searchstaxClient.Client
}

func (d *deploymentsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_deployments"
}

func (d *deploymentsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	attrs := deploymentCommonSchemaAttributes()
	attrs["id"] = schema.StringAttribute{Computed: true}
	attrs["account_name"] = schema.StringAttribute{Required: true}
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":           attrs["id"],
			"account_name": attrs["account_name"],
			"deployments_list": schema.ListNestedAttribute{
				Computed:     true,
				NestedObject: schema.NestedAttributeObject{Attributes: deploymentCommonSchemaAttributes()},
			},
		},
	}
}

func (d *deploymentsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config deploymentsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	deployments, err := d.client.GetDeployments(config.AccountName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Unable to Read SearchStax Deployments", err.Error())
		return
	}
	config.ID = types.StringValue("placeholder")
	config.Deployments = nil

	for _, deployment := range deployments.Results {
		item, diags := mapDeploymentModel(ctx, deployment)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		config.Deployments = append(config.Deployments, item)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

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

type deploymentsDataSourceModel struct {
	ID          types.String       `tfsdk:"id"`
	Deployments []deploymentsModel `tfsdk:"deployments_list"`
	AccountName types.String       `tfsdk:"account_name"`
}
