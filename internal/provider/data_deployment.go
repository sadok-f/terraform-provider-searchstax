package provider

import (
	"context"
	"fmt"

	searchstaxClient "terraform-provider-searchstax/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func NewDeploymentDataSource() datasource.DataSource { return &deploymentDataSource{} }

type deploymentDataSource struct{ client *searchstaxClient.Client }

func (d *deploymentDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_deployment"
}

func (d *deploymentDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	attrs := deploymentCommonSchemaAttributes()
	attrs["id"] = schema.StringAttribute{Computed: true}
	attrs["account_name"] = schema.StringAttribute{Required: true}
	attrs["deployment_uid"] = schema.StringAttribute{Required: true}
	resp.Schema = schema.Schema{Attributes: attrs}
}

func (d *deploymentDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *deploymentDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state deploymentDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	dep, err := d.client.GetDeployment(state.AccountName.ValueString(), state.DeploymentUID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Unable to read deployment", err.Error())
		return
	}
	mapped, diags := mapDeploymentModel(ctx, *dep)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	state.deploymentsModel = mapped
	state.ID = types.StringValue(state.AccountName.ValueString() + "/" + dep.UID)
	state.DeploymentUID = types.StringValue(dep.UID)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

type deploymentDataSourceModel struct {
	deploymentsModel
	ID            types.String `tfsdk:"id"`
	AccountName   types.String `tfsdk:"account_name"`
	DeploymentUID types.String `tfsdk:"deployment_uid"`
}
