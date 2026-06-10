package provider

import (
	"context"
	"fmt"

	searchstaxClient "terraform-provider-searchstax/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func NewBackupSchedulesDataSource() datasource.DataSource { return &backupSchedulesDataSource{} }

type backupSchedulesDataSource struct{ client *searchstaxClient.Client }

func (d *backupSchedulesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_backup_schedules"
}

func (d *backupSchedulesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{Attributes: map[string]schema.Attribute{
		"id":             schema.StringAttribute{Computed: true},
		"account_name":   schema.StringAttribute{Required: true},
		"deployment_uid": schema.StringAttribute{Required: true},
		"schedules": schema.ListNestedAttribute{Computed: true, NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
			"id":          schema.StringAttribute{Computed: true},
			"time":        schema.StringAttribute{Computed: true},
			"retention":   schema.Int64Attribute{Computed: true},
			"frequency":   schema.Int64Attribute{Computed: true},
			"region_id":   schema.StringAttribute{Computed: true},
			"days":        schema.ListAttribute{Computed: true, ElementType: types.StringType},
			"collections": schema.ListAttribute{Computed: true, ElementType: types.StringType},
		}}},
	}}
}

func (d *backupSchedulesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *backupSchedulesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state backupSchedulesDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	out, err := d.client.GetBackupSchedules(state.AccountName.ValueString(), state.DeploymentUID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Unable to read backup schedules", err.Error())
		return
	}
	state.ID = types.StringValue("placeholder")
	state.Schedules = nil
	for _, s := range out.Results {
		days, diags := types.ListValueFrom(ctx, types.StringType, s.Days)
		resp.Diagnostics.Append(diags...)
		collections, diags := types.ListValueFrom(ctx, types.StringType, s.Collections)
		resp.Diagnostics.Append(diags...)
		state.Schedules = append(state.Schedules, backupScheduleModel{
			ID:          types.StringValue(s.ID),
			Time:        types.StringValue(s.Time),
			Retention:   types.Int64Value(int64(s.Retention)),
			Frequency:   types.Int64Value(int64(s.Frequency)),
			RegionID:    types.StringValue(s.RegionID),
			Days:        days,
			Collections: collections,
		})
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

type backupSchedulesDataSourceModel struct {
	ID            types.String          `tfsdk:"id"`
	AccountName   types.String          `tfsdk:"account_name"`
	DeploymentUID types.String          `tfsdk:"deployment_uid"`
	Schedules     []backupScheduleModel `tfsdk:"schedules"`
}

type backupScheduleModel struct {
	ID          types.String `tfsdk:"id"`
	Time        types.String `tfsdk:"time"`
	Retention   types.Int64  `tfsdk:"retention"`
	Frequency   types.Int64  `tfsdk:"frequency"`
	RegionID    types.String `tfsdk:"region_id"`
	Days        types.List   `tfsdk:"days"`
	Collections types.List   `tfsdk:"collections"`
}
