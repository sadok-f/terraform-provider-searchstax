package provider

import (
	"context"
	"fmt"

	searchstaxClient "terraform-provider-searchstax/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func NewIPFiltersDataSource() datasource.DataSource { return &ipFiltersDataSource{} }

type ipFiltersDataSource struct{ client *searchstaxClient.Client }

func (d *ipFiltersDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ip_filters"
}

func (d *ipFiltersDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{Attributes: map[string]schema.Attribute{
		"id":             schema.StringAttribute{Computed: true},
		"account_name":   schema.StringAttribute{Required: true},
		"deployment_uid": schema.StringAttribute{Required: true},
		"filters": schema.ListNestedAttribute{Computed: true, NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
			"cidr_ip":     schema.StringAttribute{Computed: true},
			"description": schema.StringAttribute{Computed: true},
			"services":    schema.ListAttribute{Computed: true, ElementType: types.StringType},
		}}},
	}}
}

func (d *ipFiltersDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*searchstaxClient.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Data Source Configure Type", fmt.Sprintf("Expected *searchstaxClient.Client, got: %T.", req.ProviderData))
		return
	}
	d.client = client
}

func (d *ipFiltersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state ipFiltersDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	list, err := d.client.GetIPFilters(state.AccountName.ValueString(), state.DeploymentUID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Unable to Read SearchStax IP filters", err.Error())
		return
	}
	state.ID = types.StringValue("placeholder")
	state.Filters = nil
	for _, f := range list.Results {
		services, diags := types.ListValueFrom(ctx, types.StringType, f.Services)
		resp.Diagnostics.Append(diags...)
		state.Filters = append(state.Filters, ipFilterModel{
			CIDRIP:      types.StringValue(f.CIDRIP),
			Description: types.StringValue(f.Description),
			Services:    services,
		})
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

type ipFiltersDataSourceModel struct {
	ID            types.String    `tfsdk:"id"`
	AccountName   types.String    `tfsdk:"account_name"`
	DeploymentUID types.String    `tfsdk:"deployment_uid"`
	Filters       []ipFilterModel `tfsdk:"filters"`
}

type ipFilterModel struct {
	CIDRIP      types.String `tfsdk:"cidr_ip"`
	Description types.String `tfsdk:"description"`
	Services    types.List   `tfsdk:"services"`
}
