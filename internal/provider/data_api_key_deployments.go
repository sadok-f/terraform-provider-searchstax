package provider

import (
	"context"
	"fmt"

	searchstaxClient "terraform-provider-searchstax/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func NewAPIKeyDeploymentsDataSource() datasource.DataSource { return &apiKeyDeploymentsDataSource{} }

type apiKeyDeploymentsDataSource struct{ client *searchstaxClient.Client }

func (d *apiKeyDeploymentsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_api_key_deployments"
}

func (d *apiKeyDeploymentsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{Attributes: map[string]schema.Attribute{
		"id":           schema.StringAttribute{Computed: true},
		"account_name": schema.StringAttribute{Required: true},
		"api_key":      schema.StringAttribute{Required: true, Sensitive: true},
		"deployments":  schema.ListAttribute{Computed: true, ElementType: types.StringType},
	}}
}

func (d *apiKeyDeploymentsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	c, ok := req.ProviderData.(*searchstaxClient.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Data Source Configure Type", fmt.Sprintf("Expected *searchstaxClient.Client, got: %T.", req.ProviderData))
		return
	}
	d.client = c
}

func (d *apiKeyDeploymentsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state apiKeyDeploymentsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	out, err := d.client.GetAPIKeyDeployments(state.AccountName.ValueString(), searchstaxClient.APIKeyDeploymentsRequest{
		APIKey: state.APIKey.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Unable to read API key deployments", err.Error())
		return
	}
	deployments, diags := types.ListValueFrom(ctx, types.StringType, out.Deployments)
	resp.Diagnostics.Append(diags...)
	state.Deployments = deployments
	state.ID = types.StringValue("placeholder")
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

type apiKeyDeploymentsDataSourceModel struct {
	ID          types.String `tfsdk:"id"`
	AccountName types.String `tfsdk:"account_name"`
	APIKey      types.String `tfsdk:"api_key"`
	Deployments types.List   `tfsdk:"deployments"`
}
