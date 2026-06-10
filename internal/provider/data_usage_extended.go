package provider

import (
	"context"
	"fmt"
	"strconv"

	searchstaxClient "terraform-provider-searchstax/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func NewUsageExtendedDataSource() datasource.DataSource { return &usageExtendedDataSource{} }

type usageExtendedDataSource struct{ client *searchstaxClient.Client }

func (d *usageExtendedDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_usage_extended"
}

func (d *usageExtendedDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{Attributes: map[string]schema.Attribute{
		"id":           schema.StringAttribute{Computed: true},
		"account_name": schema.StringAttribute{Required: true},
		"year":         schema.Int64Attribute{Required: true},
		"month":        schema.Int64Attribute{Required: true},
		"events": schema.ListNestedAttribute{Computed: true, NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
			"object_id":             schema.StringAttribute{Computed: true},
			"sku":                   schema.StringAttribute{Computed: true},
			"usage":                 schema.Int64Attribute{Computed: true},
			"start_date":            schema.StringAttribute{Computed: true},
			"end_date":              schema.StringAttribute{Computed: true},
			"currency":              schema.StringAttribute{Computed: true},
			"amount":                schema.StringAttribute{Computed: true},
			"tag_collection":        schema.ListAttribute{Computed: true, ElementType: types.StringType},
			"extended_attributes":   schema.ListAttribute{Computed: true, ElementType: types.StringType},
		}}},
	}}
}

func (d *usageExtendedDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func formatUsageAmount(amount any) string {
	switch v := amount.(type) {
	case nil:
		return ""
	case string:
		return v
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	default:
		return fmt.Sprint(v)
	}
}

func (d *usageExtendedDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state usageExtendedDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	out, err := d.client.GetUsageExtended(state.AccountName.ValueString(), int(state.Year.ValueInt64()), int(state.Month.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError("Unable to read extended usage", err.Error())
		return
	}
	state.ID = types.StringValue("placeholder")
	state.Events = nil
	for _, e := range out.Results {
		tags, diags := types.ListValueFrom(ctx, types.StringType, e.TagCollection)
		resp.Diagnostics.Append(diags...)
		attrs, diags := types.ListValueFrom(ctx, types.StringType, e.ExtendedAttributes)
		resp.Diagnostics.Append(diags...)
		state.Events = append(state.Events, usageExtendedEventModel{
			ObjectID:            types.StringValue(e.ObjectID),
			SKU:                 types.StringValue(e.SKU),
			Usage:               types.Int64Value(int64(e.Usage)),
			StartDate:           types.StringValue(e.StartDate),
			EndDate:             types.StringValue(e.EndDate),
			Currency:            types.StringValue(e.Currency),
			Amount:              types.StringValue(formatUsageAmount(e.Amount)),
			TagCollection:       tags,
			ExtendedAttributes:  attrs,
		})
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

type usageExtendedDataSourceModel struct {
	ID          types.String              `tfsdk:"id"`
	AccountName types.String              `tfsdk:"account_name"`
	Year        types.Int64               `tfsdk:"year"`
	Month       types.Int64               `tfsdk:"month"`
	Events      []usageExtendedEventModel `tfsdk:"events"`
}

type usageExtendedEventModel struct {
	ObjectID           types.String `tfsdk:"object_id"`
	SKU                types.String `tfsdk:"sku"`
	Usage              types.Int64  `tfsdk:"usage"`
	StartDate          types.String `tfsdk:"start_date"`
	EndDate            types.String `tfsdk:"end_date"`
	Currency           types.String `tfsdk:"currency"`
	Amount             types.String `tfsdk:"amount"`
	TagCollection      types.List   `tfsdk:"tag_collection"`
	ExtendedAttributes types.List   `tfsdk:"extended_attributes"`
}
