package provider

import (
	"context"
	"fmt"

	searchstaxClient "terraform-provider-searchstax/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func NewZookeeperConfigsDataSource() datasource.DataSource { return &zookeeperConfigsDataSource{} }

type zookeeperConfigsDataSource struct{ client *searchstaxClient.Client }

func (d *zookeeperConfigsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_zookeeper_configs"
}
func (d *zookeeperConfigsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{Attributes: map[string]schema.Attribute{
		"id":             schema.StringAttribute{Computed: true},
		"account_name":   schema.StringAttribute{Required: true},
		"deployment_uid": schema.StringAttribute{Required: true},
		"configs": schema.ListNestedAttribute{Computed: true, NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
			"name":    schema.StringAttribute{Computed: true},
			"created": schema.StringAttribute{Computed: true},
		}}},
	}}
}
func (d *zookeeperConfigsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *zookeeperConfigsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state zookeeperConfigsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	out, err := d.client.GetZookeeperConfigs(state.AccountName.ValueString(), state.DeploymentUID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Unable to Read zookeeper configs", err.Error())
		return
	}
	state.ID = types.StringValue("placeholder")
	state.Configs = nil
	for _, c := range out.Results {
		state.Configs = append(state.Configs, zookeeperConfigModel{Name: types.StringValue(c.Name), Created: types.StringValue(c.Created)})
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

type zookeeperConfigsDataSourceModel struct {
	ID            types.String           `tfsdk:"id"`
	AccountName   types.String           `tfsdk:"account_name"`
	DeploymentUID types.String           `tfsdk:"deployment_uid"`
	Configs       []zookeeperConfigModel `tfsdk:"configs"`
}
type zookeeperConfigModel struct {
	Name    types.String `tfsdk:"name"`
	Created types.String `tfsdk:"created"`
}
