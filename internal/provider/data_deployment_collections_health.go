package provider

import (
	"context"
	"fmt"

	searchstaxClient "terraform-provider-searchstax/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func NewDeploymentCollectionsHealthDataSource() datasource.DataSource {
	return &deploymentCollectionsHealthDataSource{}
}

type deploymentCollectionsHealthDataSource struct{ client *searchstaxClient.Client }

func (d *deploymentCollectionsHealthDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_deployment_collections_health"
}

func (d *deploymentCollectionsHealthDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{Attributes: map[string]schema.Attribute{
		"id":             schema.StringAttribute{Computed: true},
		"account_name":   schema.StringAttribute{Required: true},
		"deployment_uid": schema.StringAttribute{Required: true},
		"healthy":        schema.BoolAttribute{Computed: true},
		"success":        schema.BoolAttribute{Computed: true},
		"error":          schema.StringAttribute{Computed: true},
		"collections":    schema.ListAttribute{Computed: true, ElementType: types.StringType},
	}}
}

func (d *deploymentCollectionsHealthDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *deploymentCollectionsHealthDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state deploymentCollectionsHealthDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	out, err := d.client.GetCollectionsHealth(state.AccountName.ValueString(), state.DeploymentUID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Unable to read collections health", err.Error())
		return
	}
	collections, diags := types.ListValueFrom(ctx, types.StringType, out.Collections)
	resp.Diagnostics.Append(diags...)
	state.ID = types.StringValue("placeholder")
	state.Healthy = types.BoolValue(out.Healthy)
	state.Success = types.BoolValue(out.Success)
	state.Error = types.StringValue(out.Error)
	state.Collections = collections
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

type deploymentCollectionsHealthDataSourceModel struct {
	ID            types.String `tfsdk:"id"`
	AccountName   types.String `tfsdk:"account_name"`
	DeploymentUID types.String `tfsdk:"deployment_uid"`
	Healthy       types.Bool   `tfsdk:"healthy"`
	Success       types.Bool   `tfsdk:"success"`
	Error         types.String `tfsdk:"error"`
	Collections   types.List   `tfsdk:"collections"`
}
