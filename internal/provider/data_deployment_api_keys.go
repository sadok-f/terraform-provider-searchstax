package provider

import (
	"context"
	"fmt"

	searchstaxClient "terraform-provider-searchstax/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func NewDeploymentAPIKeysDataSource() datasource.DataSource { return &deploymentAPIKeysDataSource{} }

type deploymentAPIKeysDataSource struct{ client *searchstaxClient.Client }

func (d *deploymentAPIKeysDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_deployment_api_keys"
}

func (d *deploymentAPIKeysDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{Attributes: map[string]schema.Attribute{
		"id":             schema.StringAttribute{Computed: true},
		"account_name":   schema.StringAttribute{Required: true},
		"deployment_uid": schema.StringAttribute{Required: true},
		"api_keys":       schema.ListAttribute{Computed: true, ElementType: types.StringType, Sensitive: true},
	}}
}

func (d *deploymentAPIKeysDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *deploymentAPIKeysDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state deploymentAPIKeysDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	out, err := d.client.GetDeploymentAPIKeys(state.AccountName.ValueString(), searchstaxClient.DeploymentAPIKeysRequest{
		Deployment: state.DeploymentUID.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Unable to read deployment API keys", err.Error())
		return
	}
	keys, diags := types.ListValueFrom(ctx, types.StringType, out.APIKey)
	resp.Diagnostics.Append(diags...)
	state.APIKeys = keys
	state.ID = types.StringValue("placeholder")
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

type deploymentAPIKeysDataSourceModel struct {
	ID            types.String `tfsdk:"id"`
	AccountName   types.String `tfsdk:"account_name"`
	DeploymentUID types.String `tfsdk:"deployment_uid"`
	APIKeys       types.List   `tfsdk:"api_keys"`
}
