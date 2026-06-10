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

func NewAPIKeyAssociationResource() resource.Resource { return &apiKeyAssociationResource{} }

type apiKeyAssociationResource struct{ client *searchstaxClient.Client }

func (r *apiKeyAssociationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_api_key_association"
}

func (r *apiKeyAssociationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{Attributes: map[string]schema.Attribute{
		"id":             schema.StringAttribute{Computed: true},
		"account_name":   schema.StringAttribute{Required: true},
		"api_key":        schema.StringAttribute{Required: true, Sensitive: true},
		"deployment_uid": schema.StringAttribute{Required: true},
	}}
}

func (r *apiKeyAssociationResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *apiKeyAssociationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan apiKeyAssociationResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	_, err := r.client.AssociateAPIKey(plan.AccountName.ValueString(), searchstaxClient.AssociateAPIKeyRequest{
		APIKey:     plan.APIKey.ValueString(),
		Deployment: plan.DeploymentUID.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error associating API key", err.Error())
		return
	}
	plan.ID = types.StringValue(plan.AccountName.ValueString() + "/" + plan.DeploymentUID.ValueString())
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *apiKeyAssociationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state apiKeyAssociationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	out, err := r.client.GetDeploymentAPIKeys(state.AccountName.ValueString(), searchstaxClient.DeploymentAPIKeysRequest{
		Deployment: state.DeploymentUID.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error reading API key association", err.Error())
		return
	}
	for _, key := range out.APIKey {
		if key == state.APIKey.ValueString() {
			state.ID = types.StringValue(state.AccountName.ValueString() + "/" + state.DeploymentUID.ValueString())
			resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
			return
		}
	}
	resp.State.RemoveResource(ctx)
}

func (r *apiKeyAssociationResource) Update(context.Context, resource.UpdateRequest, *resource.UpdateResponse) {
}

func (r *apiKeyAssociationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state apiKeyAssociationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	_, err := r.client.DisassociateAPIKey(state.AccountName.ValueString(), searchstaxClient.AssociateAPIKeyRequest{
		APIKey:     state.APIKey.ValueString(),
		Deployment: state.DeploymentUID.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error disassociating API key", err.Error())
	}
}

func (r *apiKeyAssociationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, "/")
	if len(parts) != 2 {
		resp.Diagnostics.AddError("Unexpected Import Identifier", "Expected account_name/deployment_uid")
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("account_name"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("deployment_uid"), parts[1])...)
}

type apiKeyAssociationResourceModel struct {
	ID            types.String `tfsdk:"id"`
	AccountName   types.String `tfsdk:"account_name"`
	APIKey        types.String `tfsdk:"api_key"`
	DeploymentUID types.String `tfsdk:"deployment_uid"`
}
