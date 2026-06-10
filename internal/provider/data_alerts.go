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

func NewAlertsDataSource() datasource.DataSource { return &alertsDataSource{} }

type alertsDataSource struct{ client *searchstaxClient.Client }

func (d *alertsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_alerts"
}

func (d *alertsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{Attributes: map[string]schema.Attribute{
		"id":             schema.StringAttribute{Computed: true},
		"account_name":   schema.StringAttribute{Required: true},
		"deployment_uid": schema.StringAttribute{Required: true},
		"alerts": schema.ListNestedAttribute{Computed: true, NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Computed: true},
		}}},
	}}
}

func (d *alertsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *alertsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state alertsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	out, err := d.client.GetAlerts(state.AccountName.ValueString(), state.DeploymentUID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Unable to read alerts", err.Error())
		return
	}
	state.ID = types.StringValue("placeholder")
	state.Alerts = nil
	for _, a := range out.Results {
		state.Alerts = append(state.Alerts, alertListModel{ID: types.StringValue(strconv.FormatInt(a.ID, 10))})
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

type alertsDataSourceModel struct {
	ID            types.String     `tfsdk:"id"`
	AccountName   types.String     `tfsdk:"account_name"`
	DeploymentUID types.String     `tfsdk:"deployment_uid"`
	Alerts        []alertListModel `tfsdk:"alerts"`
}

type alertListModel struct {
	ID types.String `tfsdk:"id"`
}
