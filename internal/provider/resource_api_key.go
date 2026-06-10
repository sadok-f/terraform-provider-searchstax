package provider

import (
	"context"
	"fmt"

	searchstaxClient "terraform-provider-searchstax/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func NewAPIKeyResource() resource.Resource { return &apiKeyResource{} }

type apiKeyResource struct{ client *searchstaxClient.Client }

func (r *apiKeyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_api_key"
}

func (r *apiKeyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{Attributes: map[string]schema.Attribute{
		"id":           schema.StringAttribute{Computed: true},
		"account_name": schema.StringAttribute{Required: true},
		"api_key":      schema.StringAttribute{Computed: true, Sensitive: true},
		"scope": schema.ListAttribute{Optional: true, ElementType: types.StringType},
	}}
}

func (r *apiKeyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *apiKeyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan apiKeyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var scope []string
	resp.Diagnostics.Append(plan.Scope.ElementsAs(ctx, &scope, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	out, err := r.client.CreateAPIKey(plan.AccountName.ValueString(), searchstaxClient.CreateAPIKeyRequest{Scope: scope})
	if err != nil {
		resp.Diagnostics.AddError("Error creating API key", err.Error())
		return
	}
	plan.APIKey = types.StringValue(out.APIKey)
	plan.ID = types.StringValue(plan.AccountName.ValueString() + "/" + out.APIKey)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *apiKeyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state apiKeyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *apiKeyResource) Update(context.Context, resource.UpdateRequest, *resource.UpdateResponse) {}

func (r *apiKeyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state apiKeyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.RevokeAPIKey(state.AccountName.ValueString(), searchstaxClient.RevokeAPIKeyRequest{APIKey: state.APIKey.ValueString()}); err != nil {
		resp.Diagnostics.AddError("Error revoking API key", err.Error())
	}
}

type apiKeyResourceModel struct {
	ID          types.String `tfsdk:"id"`
	AccountName types.String `tfsdk:"account_name"`
	APIKey      types.String `tfsdk:"api_key"`
	Scope       types.List   `tfsdk:"scope"`
}
