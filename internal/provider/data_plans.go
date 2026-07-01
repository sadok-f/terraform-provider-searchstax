package provider

import (
	"context"
	"fmt"

	searchstaxClient "terraform-provider-searchstax/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func NewPlansDataSource() datasource.DataSource { return &plansDataSource{} }

type plansDataSource struct{ client *searchstaxClient.Client }

func (d *plansDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_plans"
}

func (d *plansDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{Attributes: map[string]schema.Attribute{
		"id":           schema.StringAttribute{Computed: true},
		"account_name": schema.StringAttribute{Required: true},
		"application":  schema.StringAttribute{Optional: true},
		"plan_type":    schema.StringAttribute{Optional: true},
		"page":         schema.Int64Attribute{Optional: true},
		"total_count":  schema.Int64Attribute{Computed: true},
		"plans": schema.ListNestedAttribute{Computed: true, NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
			"name":                 schema.StringAttribute{Computed: true},
			"description":          schema.StringAttribute{Computed: true},
			"plan_type":            schema.StringAttribute{Computed: true},
			"application":          schema.StringAttribute{Computed: true},
			"trial_available":      schema.BoolAttribute{Computed: true},
			"application_versions": schema.ListAttribute{Computed: true, ElementType: types.StringType},
			"regions": schema.ListNestedAttribute{Computed: true, NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
				"region_id":                         schema.StringAttribute{Computed: true},
				"cloud_provider":                    schema.StringAttribute{Computed: true},
				"cloud_provider_id":                 schema.StringAttribute{Computed: true},
				"price":                             schema.Float64Attribute{Computed: true},
				"additional_application_node_price": schema.Float64Attribute{Computed: true},
				"additional_zookeeper_node_price":   schema.Float64Attribute{Computed: true},
			}}},
		}}},
	}}
}

func (d *plansDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *plansDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state plansDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	page := int(state.Page.ValueInt64())
	var (
		out *searchstaxClient.PlansList
		err error
	)
	if page > 0 {
		out, err = d.client.GetPlans(state.AccountName.ValueString(), state.Application.ValueString(), state.PlanType.ValueString(), page)
	} else {
		out, err = d.client.GetAllPlans(state.AccountName.ValueString(), state.Application.ValueString(), state.PlanType.ValueString())
	}
	if err != nil {
		resp.Diagnostics.AddError("Unable to read plans", err.Error())
		return
	}
	state.ID = types.StringValue("placeholder")
	state.TotalCount = types.Int64Value(int64(out.Count))
	state.Plans = nil
	for _, p := range out.Results {
		versions, diags := types.ListValueFrom(ctx, types.StringType, p.ApplicationVersions)
		resp.Diagnostics.Append(diags...)
		regions := make([]planRegionModel, 0, len(p.PlanRegions))
		for _, r := range p.PlanRegions {
			regions = append(regions, planRegionModel{
				RegionID:                       types.StringValue(r.RegionID),
				CloudProvider:                  types.StringValue(r.CloudProvider),
				CloudProviderID:                types.StringValue(r.CloudProviderID),
				Price:                          types.Float64Value(r.Price),
				AdditionalApplicationNodePrice: types.Float64Value(r.AdditionalApplicationNodePrice),
				AdditionalZookeeperNodePrice:   types.Float64Value(r.AdditionalZookeeperNodePrice),
			})
		}
		state.Plans = append(state.Plans, planModel{
			Name:                types.StringValue(p.Name),
			Description:         types.StringValue(p.Description),
			PlanType:            types.StringValue(p.PlanType),
			Application:         types.StringValue(p.Application),
			TrialAvailable:      types.BoolValue(p.TrialAvailable),
			ApplicationVersions: versions,
			Regions:             regions,
		})
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

type plansDataSourceModel struct {
	ID          types.String `tfsdk:"id"`
	AccountName types.String `tfsdk:"account_name"`
	Application types.String `tfsdk:"application"`
	PlanType    types.String `tfsdk:"plan_type"`
	Page        types.Int64  `tfsdk:"page"`
	TotalCount  types.Int64  `tfsdk:"total_count"`
	Plans       []planModel  `tfsdk:"plans"`
}

type planModel struct {
	Name                types.String      `tfsdk:"name"`
	Description         types.String      `tfsdk:"description"`
	PlanType            types.String      `tfsdk:"plan_type"`
	Application         types.String      `tfsdk:"application"`
	TrialAvailable      types.Bool        `tfsdk:"trial_available"`
	ApplicationVersions types.List        `tfsdk:"application_versions"`
	Regions             []planRegionModel `tfsdk:"regions"`
}

type planRegionModel struct {
	RegionID                       types.String  `tfsdk:"region_id"`
	CloudProvider                  types.String  `tfsdk:"cloud_provider"`
	CloudProviderID                types.String  `tfsdk:"cloud_provider_id"`
	Price                          types.Float64 `tfsdk:"price"`
	AdditionalApplicationNodePrice types.Float64 `tfsdk:"additional_application_node_price"`
	AdditionalZookeeperNodePrice   types.Float64 `tfsdk:"additional_zookeeper_node_price"`
}
