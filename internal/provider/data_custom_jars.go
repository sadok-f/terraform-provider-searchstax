package provider

import (
	"context"
	"fmt"

	searchstaxClient "terraform-provider-searchstax/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func NewCustomJarsDataSource() datasource.DataSource { return &customJarsDataSource{} }

type customJarsDataSource struct{ client *searchstaxClient.Client }

func (d *customJarsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_custom_jars"
}
func (d *customJarsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{Attributes: map[string]schema.Attribute{
		"id":             schema.StringAttribute{Computed: true},
		"account_name":   schema.StringAttribute{Required: true},
		"deployment_uid": schema.StringAttribute{Required: true},
		"jars": schema.ListNestedAttribute{Computed: true, NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{Computed: true},
		}}},
	}}
}
func (d *customJarsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *customJarsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state customJarsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	out, err := d.client.GetCustomJars(state.AccountName.ValueString(), state.DeploymentUID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Unable to Read custom jars", err.Error())
		return
	}
	state.ID = types.StringValue("placeholder")
	state.Jars = nil
	for _, j := range out.Results {
		state.Jars = append(state.Jars, customJarModel{Name: types.StringValue(j.Name)})
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

type customJarsDataSourceModel struct {
	ID            types.String     `tfsdk:"id"`
	AccountName   types.String     `tfsdk:"account_name"`
	DeploymentUID types.String     `tfsdk:"deployment_uid"`
	Jars          []customJarModel `tfsdk:"jars"`
}
type customJarModel struct {
	Name types.String `tfsdk:"name"`
}
