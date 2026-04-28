package provider

import (
	"context"
	"fmt"

	searchstaxClient "terraform-provider-searchstax/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSourceWithConfigure = &dnsRecordsDataSource{}

func NewDNSRecordsDataSource() datasource.DataSource { return &dnsRecordsDataSource{} }

type dnsRecordsDataSource struct{ client *searchstaxClient.Client }

func (d *dnsRecordsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dns_records"
}

func (d *dnsRecordsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{Attributes: map[string]schema.Attribute{
		"id":           schema.StringAttribute{Computed: true},
		"account_name": schema.StringAttribute{Required: true},
		"records": schema.ListNestedAttribute{Computed: true, NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
			"name":       schema.StringAttribute{Computed: true},
			"deployment": schema.StringAttribute{Computed: true},
			"ttl":        schema.StringAttribute{Computed: true},
		}}},
	}}
}

func (d *dnsRecordsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *dnsRecordsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state dnsRecordsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	list, err := d.client.GetDNSRecords(state.AccountName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Unable to Read SearchStax DNS Records", err.Error())
		return
	}
	state.ID = types.StringValue("placeholder")
	state.Records = nil
	for _, r := range list.Results {
		state.Records = append(state.Records, dnsRecordModel{
			Name:       types.StringValue(r.Name),
			Deployment: types.StringValue(r.Deployment),
			TTL:        types.StringValue(r.TTL),
		})
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

type dnsRecordsDataSourceModel struct {
	ID          types.String     `tfsdk:"id"`
	AccountName types.String     `tfsdk:"account_name"`
	Records     []dnsRecordModel `tfsdk:"records"`
}

type dnsRecordModel struct {
	Name       types.String `tfsdk:"name"`
	Deployment types.String `tfsdk:"deployment"`
	TTL        types.String `tfsdk:"ttl"`
}
