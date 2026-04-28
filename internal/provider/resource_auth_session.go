package provider

import (
	"context"
	"fmt"

	searchstaxClient "terraform-provider-searchstax/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func NewAuthSessionResource() resource.Resource { return &authSessionResource{} }

type authSessionResource struct{ client *searchstaxClient.Client }

func (r *authSessionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_auth_session"
}
func (r *authSessionResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{Attributes: map[string]schema.Attribute{
		"id":    schema.StringAttribute{Computed: true},
		"token": schema.StringAttribute{Computed: true, Sensitive: true},
	}}
}
func (r *authSessionResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *authSessionResource) Create(ctx context.Context, _ resource.CreateRequest, resp *resource.CreateResponse) {
	token, err := r.client.SignIn()
	if err != nil {
		resp.Diagnostics.AddError("Error creating auth session", err.Error())
		return
	}
	state := authSessionResourceModel{ID: types.StringValue("placeholder"), Token: types.StringValue(token.Token)}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
func (r *authSessionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state authSessionResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if state.Token.IsNull() || state.Token.ValueString() == "" {
		resp.State.RemoveResource(ctx)
		return
	}
}
func (r *authSessionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state authSessionResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
func (r *authSessionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state authSessionResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if state.Token.IsNull() || state.Token.ValueString() == "" {
		return
	}
	// best-effort sign out (mock may not expose this endpoint)
	token := state.Token.ValueString()
	_ = r.client.SignOut(&token)
}

type authSessionResourceModel struct {
	ID    types.String `tfsdk:"id"`
	Token types.String `tfsdk:"token"`
}
