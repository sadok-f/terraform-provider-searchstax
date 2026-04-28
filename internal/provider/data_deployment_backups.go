package provider

import (
	"context"
	"fmt"

	searchstaxClient "terraform-provider-searchstax/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func NewDeploymentBackupsDataSource() datasource.DataSource { return &deploymentBackupsDataSource{} }

type deploymentBackupsDataSource struct{ client *searchstaxClient.Client }

func (d *deploymentBackupsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_deployment_backups"
}
func (d *deploymentBackupsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{Attributes: map[string]schema.Attribute{
		"id":             schema.StringAttribute{Computed: true},
		"account_name":   schema.StringAttribute{Required: true},
		"deployment_uid": schema.StringAttribute{Required: true},
		"backups": schema.ListNestedAttribute{Computed: true, NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Computed: true},
		}}},
	}}
}
func (d *deploymentBackupsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *deploymentBackupsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state deploymentBackupsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	out, err := d.client.GetDeploymentBackups(state.AccountName.ValueString(), state.DeploymentUID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Unable to Read deployment backups", err.Error())
		return
	}
	state.ID = types.StringValue("placeholder")
	state.Backups = nil
	for _, b := range out.Results {
		state.Backups = append(state.Backups, backupModel{ID: types.StringValue(b.ID)})
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

type deploymentBackupsDataSourceModel struct {
	ID            types.String  `tfsdk:"id"`
	AccountName   types.String  `tfsdk:"account_name"`
	DeploymentUID types.String  `tfsdk:"deployment_uid"`
	Backups       []backupModel `tfsdk:"backups"`
}
type backupModel struct {
	ID types.String `tfsdk:"id"`
}
