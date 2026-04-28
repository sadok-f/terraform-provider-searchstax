package provider

import (
	"context"
	"fmt"

	searchstaxClient "terraform-provider-searchstax/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Note: SearchStax mock API currently exposes webhook list/read semantics only.
// This resource models a managed reference to an existing webhook by id.
func NewWebhookResource() resource.Resource { return &webhookResource{} }

type webhookResource struct{ client *searchstaxClient.Client }

func (r *webhookResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_webhook"
}
func (r *webhookResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{Attributes: map[string]schema.Attribute{
		"id":           schema.StringAttribute{Computed: true},
		"account_name": schema.StringAttribute{Required: true},
		"webhook_id":   schema.Int64Attribute{Required: true},
		"name":         schema.StringAttribute{Computed: true},
		"url":          schema.StringAttribute{Computed: true},
		"paused":       schema.BoolAttribute{Computed: true},
	}}
}
func (r *webhookResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *webhookResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan webhookResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	list, err := r.client.GetWebhooks(plan.AccountName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading SearchStax webhooks", err.Error())
		return
	}
	for _, w := range list.Results {
		if w.ID == plan.WebhookID.ValueInt64() {
			plan.ID = types.StringValue(fmt.Sprintf("%s/%d", plan.AccountName.ValueString(), w.ID))
			plan.Name = types.StringValue(w.Name)
			plan.URL = types.StringValue(w.URL)
			plan.Paused = types.BoolValue(w.Paused)
			resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
			return
		}
	}
	resp.Diagnostics.AddError("Webhook not found", fmt.Sprintf("Webhook id %d not found for account %s", plan.WebhookID.ValueInt64(), plan.AccountName.ValueString()))
}
func (r *webhookResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state webhookResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	list, err := r.client.GetWebhooks(state.AccountName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading SearchStax webhooks", err.Error())
		return
	}
	for _, w := range list.Results {
		if w.ID == state.WebhookID.ValueInt64() {
			state.ID = types.StringValue(fmt.Sprintf("%s/%d", state.AccountName.ValueString(), w.ID))
			state.Name = types.StringValue(w.Name)
			state.URL = types.StringValue(w.URL)
			state.Paused = types.BoolValue(w.Paused)
			resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
			return
		}
	}
	resp.State.RemoveResource(ctx)
}
func (r *webhookResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan webhookResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	list, err := r.client.GetWebhooks(plan.AccountName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading SearchStax webhooks", err.Error())
		return
	}
	for _, w := range list.Results {
		if w.ID == plan.WebhookID.ValueInt64() {
			plan.ID = types.StringValue(fmt.Sprintf("%s/%d", plan.AccountName.ValueString(), w.ID))
			plan.Name = types.StringValue(w.Name)
			plan.URL = types.StringValue(w.URL)
			plan.Paused = types.BoolValue(w.Paused)
			resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
			return
		}
	}
	resp.Diagnostics.AddError("Webhook not found", fmt.Sprintf("Webhook id %d not found for account %s", plan.WebhookID.ValueInt64(), plan.AccountName.ValueString()))
}
func (r *webhookResource) Delete(context.Context, resource.DeleteRequest, *resource.DeleteResponse) {}

type webhookResourceModel struct {
	ID          types.String `tfsdk:"id"`
	AccountName types.String `tfsdk:"account_name"`
	WebhookID   types.Int64  `tfsdk:"webhook_id"`
	Name        types.String `tfsdk:"name"`
	URL         types.String `tfsdk:"url"`
	Paused      types.Bool   `tfsdk:"paused"`
}
