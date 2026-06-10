package provider

import (
	"context"
	"fmt"

	searchstaxClient "terraform-provider-searchstax/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func NewUsageDataSource() datasource.DataSource { return &usageDataSource{} }

type usageDataSource struct{ client *searchstaxClient.Client }

func (d *usageDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_usage"
}

func (d *usageDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{Attributes: map[string]schema.Attribute{
		"id":           schema.StringAttribute{Computed: true},
		"account_name": schema.StringAttribute{Required: true},
		"year":         schema.Int64Attribute{Required: true},
		"month":        schema.Int64Attribute{Required: true},
		"events": schema.ListNestedAttribute{Computed: true, NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
			"object_id":      schema.StringAttribute{Computed: true},
			"sku":            schema.StringAttribute{Computed: true},
			"usage":          schema.Int64Attribute{Computed: true},
			"start_date":     schema.StringAttribute{Computed: true},
			"end_date":       schema.StringAttribute{Computed: true},
			"currency":       schema.StringAttribute{Computed: true},
			"amount":         schema.StringAttribute{Computed: true},
			"tag_collection": schema.ListAttribute{Computed: true, ElementType: types.StringType},
		}}},
	}}
}

func (d *usageDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *usageDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state usageDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	out, err := d.client.GetUsage(state.AccountName.ValueString(), int(state.Year.ValueInt64()), int(state.Month.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError("Unable to read usage", err.Error())
		return
	}
	state.ID = types.StringValue("placeholder")
	state.Events = nil
	for _, e := range out.Results {
		tags, diags := types.ListValueFrom(ctx, types.StringType, e.TagCollection)
		resp.Diagnostics.Append(diags...)
		state.Events = append(state.Events, usageEventModel{
			ObjectID:      types.StringValue(e.ObjectID),
			SKU:           types.StringValue(e.SKU),
			Usage:         types.Int64Value(int64(e.Usage)),
			StartDate:     types.StringValue(e.StartDate),
			EndDate:       types.StringValue(e.EndDate),
			Currency:      types.StringValue(e.Currency),
			Amount:        types.StringValue(e.Amount),
			TagCollection: tags,
		})
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

type usageDataSourceModel struct {
	ID          types.String      `tfsdk:"id"`
	AccountName types.String      `tfsdk:"account_name"`
	Year        types.Int64       `tfsdk:"year"`
	Month       types.Int64       `tfsdk:"month"`
	Events      []usageEventModel `tfsdk:"events"`
}

type usageEventModel struct {
	ObjectID      types.String `tfsdk:"object_id"`
	SKU           types.String `tfsdk:"sku"`
	Usage         types.Int64  `tfsdk:"usage"`
	StartDate     types.String `tfsdk:"start_date"`
	EndDate       types.String `tfsdk:"end_date"`
	Currency      types.String `tfsdk:"currency"`
	Amount        types.String `tfsdk:"amount"`
	TagCollection types.List   `tfsdk:"tag_collection"`
}
