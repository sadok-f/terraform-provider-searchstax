package provider

import (
	"context"
	"fmt"

	searchstaxClient "terraform-provider-searchstax/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func NewZookeeperConfigDataSource() datasource.DataSource { return &zookeeperConfigDataSource{} }

type zookeeperConfigDataSource struct{ client *searchstaxClient.Client }

func (d *zookeeperConfigDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_zookeeper_config"
}

func (d *zookeeperConfigDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{Attributes: map[string]schema.Attribute{
		"id":             schema.StringAttribute{Computed: true},
		"account_name":   schema.StringAttribute{Required: true},
		"deployment_uid": schema.StringAttribute{Required: true},
		"name":           schema.StringAttribute{Required: true},
		"files":          schema.ListAttribute{Computed: true, ElementType: types.StringType},
	}}
}

func (d *zookeeperConfigDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *zookeeperConfigDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state zookeeperConfigDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	out, err := d.client.GetZookeeperConfig(state.AccountName.ValueString(), state.DeploymentUID.ValueString(), state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Unable to read zookeeper config", err.Error())
		return
	}
	files, diags := types.ListValueFrom(ctx, types.StringType, out.Files)
	resp.Diagnostics.Append(diags...)
	state.ID = types.StringValue(state.AccountName.ValueString() + "/" + state.DeploymentUID.ValueString() + "/" + out.Name)
	if out.Name != "" {
		state.Name = types.StringValue(out.Name)
	}
	state.Files = files
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

type zookeeperConfigDataSourceModel struct {
	ID            types.String `tfsdk:"id"`
	AccountName   types.String `tfsdk:"account_name"`
	DeploymentUID types.String `tfsdk:"deployment_uid"`
	Name          types.String `tfsdk:"name"`
	Files         types.List   `tfsdk:"files"`
}
