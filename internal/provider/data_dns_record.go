package provider

import (
	"context"
	"fmt"

	searchstaxClient "terraform-provider-searchstax/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func NewDNSRecordDataSource() datasource.DataSource { return &dnsRecordDataSource{} }

type dnsRecordDataSource struct{ client *searchstaxClient.Client }

func (d *dnsRecordDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dns_record"
}

func (d *dnsRecordDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{Attributes: map[string]schema.Attribute{
		"id":           schema.StringAttribute{Computed: true},
		"account_name": schema.StringAttribute{Required: true},
		"name":         schema.StringAttribute{Required: true},
		"deployment":   schema.StringAttribute{Computed: true},
		"ttl":          schema.StringAttribute{Computed: true},
	}}
}

func (d *dnsRecordDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *dnsRecordDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state dnsRecordDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	record, err := d.client.GetDNSRecord(state.AccountName.ValueString(), state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Unable to read DNS record", err.Error())
		return
	}
	state.ID = types.StringValue(state.AccountName.ValueString() + "/" + record.Name)
	state.Deployment = types.StringValue(record.Deployment)
	state.TTL = types.StringValue(record.TTL)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

type dnsRecordDataSourceModel struct {
	ID          types.String `tfsdk:"id"`
	AccountName types.String `tfsdk:"account_name"`
	Name        types.String `tfsdk:"name"`
	Deployment  types.String `tfsdk:"deployment"`
	TTL         types.String `tfsdk:"ttl"`
}
