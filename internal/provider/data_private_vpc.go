package provider

import (
	"context"
	"fmt"

	searchstaxClient "terraform-provider-searchstax/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &privateVpcDataSource{}
	_ datasource.DataSourceWithConfigure = &privateVpcDataSource{}
)

// NewPrivateVpcDataSource is a helper function to simplify the provider implementation.
func NewPrivateVpcDataSource() datasource.DataSource {
	return &privateVpcDataSource{}
}

// privateVpcDataSource is the data source implementation.
type privateVpcDataSource struct {
	client *searchstaxClient.Client
}

// Metadata returns the data source type name.
func (d *privateVpcDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_private_vpc"
}

// Schema - defines the schema for the data source.go.
func (d *privateVpcDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"private_vpc_list": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Computed: true,
						},
						"account": schema.StringAttribute{
							Computed: true,
						},
						"region": schema.StringAttribute{
							Computed: true,
						},
						"status": schema.StringAttribute{
							Computed: true,
						},
						"address_space": schema.StringAttribute{
							Computed: true,
						},
						"name": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
			"account_name": schema.StringAttribute{
				Required: true,
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *privateVpcDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config privateVpcDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	privateVpc, err := d.client.GetPrivateVpc(config.AccountName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read SearchStax PrivateVpc",
			err.Error(),
		)
		return
	}
	config.ID = types.StringValue("placeholder")

	// Map response body to model
	for _, privateVpc := range privateVpc.Results {
		privateVpcState := privateVpcModel{
			ID:           types.Int64Value(privateVpc.ID),
			Name:         types.StringValue(privateVpc.Name),
			Account:      types.StringValue(privateVpc.Account),
			Region:       types.StringValue(privateVpc.Region),
			AddressSpace: types.StringValue(privateVpc.AddressSpace),
			Status:       types.StringValue(privateVpc.Status),
		}

		config.PrivateVpcList = append(config.PrivateVpcList, privateVpcState)
	}

	// Set state
	diags := resp.State.Set(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the data source.
func (d *privateVpcDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*searchstaxClient.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *searchstaxClient.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

// privateVpcDataSourceModel  maps the data source schema data.
type privateVpcDataSourceModel struct {
	ID             types.String      `tfsdk:"id"`
	PrivateVpcList []privateVpcModel `tfsdk:"private_vpc_list"`
	AccountName    types.String      `tfsdk:"account_name"`
}

// privateVpcModel maps datasource_privateVpc schema data.
type privateVpcModel struct {
	ID           types.Int64  `tfsdk:"id"`
	Account      types.String `tfsdk:"account"`
	Region       types.String `tfsdk:"region"`
	AddressSpace types.String `tfsdk:"address_space"`
	Status       types.String `tfsdk:"status"`
	Name         types.String `tfsdk:"name"`
}
