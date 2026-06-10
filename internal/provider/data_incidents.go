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

func NewIncidentsDataSource() datasource.DataSource { return &incidentsDataSource{} }

type incidentsDataSource struct{ client *searchstaxClient.Client }

func (d *incidentsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_incidents"
}

func (d *incidentsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{Attributes: map[string]schema.Attribute{
		"id":             schema.StringAttribute{Computed: true},
		"account_name":   schema.StringAttribute{Required: true},
		"deployment_uid": schema.StringAttribute{Required: true},
		"incidents": schema.ListNestedAttribute{Computed: true, NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Computed: true},
		}}},
	}}
}

func (d *incidentsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *incidentsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state incidentsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	out, err := d.client.GetIncidents(state.AccountName.ValueString(), state.DeploymentUID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Unable to read incidents", err.Error())
		return
	}
	state.ID = types.StringValue("placeholder")
	state.Incidents = nil
	for _, i := range out.Results {
		state.Incidents = append(state.Incidents, incidentModel{ID: types.StringValue(strconv.FormatInt(i.ID, 10))})
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

type incidentsDataSourceModel struct {
	ID            types.String    `tfsdk:"id"`
	AccountName   types.String    `tfsdk:"account_name"`
	DeploymentUID types.String    `tfsdk:"deployment_uid"`
	Incidents     []incidentModel `tfsdk:"incidents"`
}

type incidentModel struct {
	ID types.String `tfsdk:"id"`
}
