package provider

import (
	"context"
	"fmt"

	searchstaxClient "terraform-provider-searchstax/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func NewAlertMetricsDataSource() datasource.DataSource { return &alertMetricsDataSource{} }

type alertMetricsDataSource struct{ client *searchstaxClient.Client }

func (d *alertMetricsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_alert_metrics"
}
func (d *alertMetricsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{Attributes: map[string]schema.Attribute{
		"id":           schema.StringAttribute{Computed: true},
		"account_name": schema.StringAttribute{Required: true},
		"metrics": schema.ListNestedAttribute{Computed: true, NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
			"metric": schema.StringAttribute{Computed: true},
			"unit":   schema.StringAttribute{Computed: true},
		}}},
	}}
}
func (d *alertMetricsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *alertMetricsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state alertMetricsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	out, err := d.client.GetAlertMetrics(state.AccountName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Unable to Read alert metrics", err.Error())
		return
	}
	state.ID = types.StringValue("placeholder")
	state.Metrics = nil
	for _, m := range out.Results {
		state.Metrics = append(state.Metrics, alertMetricModel{Metric: types.StringValue(m.Metric), Unit: types.StringValue(m.Unit)})
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

type alertMetricsDataSourceModel struct {
	ID          types.String       `tfsdk:"id"`
	AccountName types.String       `tfsdk:"account_name"`
	Metrics     []alertMetricModel `tfsdk:"metrics"`
}
type alertMetricModel struct {
	Metric types.String `tfsdk:"metric"`
	Unit   types.String `tfsdk:"unit"`
}
