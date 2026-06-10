package provider

import (
	"context"
	"fmt"

	searchstaxClient "terraform-provider-searchstax/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func NewRestoreStatusDataSource() datasource.DataSource { return &restoreStatusDataSource{} }

type restoreStatusDataSource struct{ client *searchstaxClient.Client }

func (d *restoreStatusDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_restore_status"
}

func (d *restoreStatusDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{Attributes: map[string]schema.Attribute{
		"id":             schema.StringAttribute{Computed: true},
		"account_name":   schema.StringAttribute{Required: true},
		"deployment_uid": schema.StringAttribute{Required: true},
		"backup_id":      schema.StringAttribute{Required: true},
		"restore_id":     schema.StringAttribute{Computed: true},
		"status":         schema.StringAttribute{Computed: true},
	}}
}

func (d *restoreStatusDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *restoreStatusDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state restoreStatusDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	reqBody := searchstaxClient.RestoreRequest{BackupID: state.BackupID.ValueString()}
	out, err := d.client.GetDeploymentRestoreStatus(state.AccountName.ValueString(), state.DeploymentUID.ValueString(), reqBody)
	if err != nil {
		resp.Diagnostics.AddError("Unable to read restore status", err.Error())
		return
	}
	state.ID = types.StringValue(state.AccountName.ValueString() + "/" + state.DeploymentUID.ValueString() + "/" + state.BackupID.ValueString())
	state.RestoreID = types.StringValue(out.RestoreID)
	state.Status = types.StringValue(out.Status)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

type restoreStatusDataSourceModel struct {
	ID            types.String `tfsdk:"id"`
	AccountName   types.String `tfsdk:"account_name"`
	DeploymentUID types.String `tfsdk:"deployment_uid"`
	BackupID      types.String `tfsdk:"backup_id"`
	RestoreID     types.String `tfsdk:"restore_id"`
	Status        types.String `tfsdk:"status"`
}
