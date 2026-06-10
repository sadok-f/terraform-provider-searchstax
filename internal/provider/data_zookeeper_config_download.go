package provider

import (
	"context"
	"fmt"

	searchstaxClient "terraform-provider-searchstax/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func NewZookeeperConfigDownloadDataSource() datasource.DataSource {
	return &zookeeperConfigDownloadDataSource{}
}

type zookeeperConfigDownloadDataSource struct{ client *searchstaxClient.Client }

func (d *zookeeperConfigDownloadDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_zookeeper_config_download"
}

func (d *zookeeperConfigDownloadDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{Attributes: map[string]schema.Attribute{
		"id":             schema.StringAttribute{Computed: true},
		"account_name":   schema.StringAttribute{Required: true},
		"deployment_uid": schema.StringAttribute{Required: true},
		"name":           schema.StringAttribute{Required: true},
		"download":       schema.StringAttribute{Computed: true},
		"note":           schema.StringAttribute{Computed: true},
	}}
}

func (d *zookeeperConfigDownloadDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *zookeeperConfigDownloadDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state zookeeperConfigDownloadDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	out, err := d.client.DownloadZookeeperConfig(state.AccountName.ValueString(), state.DeploymentUID.ValueString(), state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Unable to download zookeeper config", err.Error())
		return
	}
	state.ID = types.StringValue("placeholder")
	state.Download = types.StringValue(out.Download)
	state.Note = types.StringValue(out.Note)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

type zookeeperConfigDownloadDataSourceModel struct {
	ID            types.String `tfsdk:"id"`
	AccountName   types.String `tfsdk:"account_name"`
	DeploymentUID types.String `tfsdk:"deployment_uid"`
	Name          types.String `tfsdk:"name"`
	Download      types.String `tfsdk:"download"`
	Note          types.String `tfsdk:"note"`
}
