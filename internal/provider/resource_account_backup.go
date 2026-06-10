package provider

import (
	"context"
	"fmt"
	"strings"

	searchstaxClient "terraform-provider-searchstax/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func NewAccountBackupResource() resource.Resource { return &accountBackupResource{} }

type accountBackupResource struct{ client *searchstaxClient.Client }

func (r *accountBackupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_account_backup"
}

func (r *accountBackupResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{Attributes: map[string]schema.Attribute{
		"id":           schema.StringAttribute{Computed: true},
		"account_name": schema.StringAttribute{Required: true},
		"backup_id":    schema.StringAttribute{Required: true},
	}}
}

func (r *accountBackupResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	c, ok := req.ProviderData.(*searchstaxClient.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Resource Configure Type", fmt.Sprintf("Expected *searchstaxClient.Client, got: %T.", req.ProviderData))
		return
	}
	r.client = c
}

func (r *accountBackupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan accountBackupResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.ID = types.StringValue(plan.AccountName.ValueString() + "/" + plan.BackupID.ValueString())
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *accountBackupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state accountBackupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	list, err := r.client.GetAccountBackups(state.AccountName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading account backups", err.Error())
		return
	}
	for _, b := range list.Results {
		if b.ID == state.BackupID.ValueString() {
			state.ID = types.StringValue(state.AccountName.ValueString() + "/" + state.BackupID.ValueString())
			resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
			return
		}
	}
	state.ID = types.StringValue(state.AccountName.ValueString() + "/" + state.BackupID.ValueString())
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *accountBackupResource) Update(context.Context, resource.UpdateRequest, *resource.UpdateResponse) {
}

func (r *accountBackupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state accountBackupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.DeleteAccountBackup(state.AccountName.ValueString(), state.BackupID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting account backup", err.Error())
	}
}

func (r *accountBackupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, "/")
	if len(parts) != 2 {
		resp.Diagnostics.AddError("Unexpected Import Identifier", "Expected account_name/backup_id")
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("account_name"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("backup_id"), parts[1])...)
}

type accountBackupResourceModel struct {
	ID          types.String `tfsdk:"id"`
	AccountName types.String `tfsdk:"account_name"`
	BackupID    types.String `tfsdk:"backup_id"`
}
