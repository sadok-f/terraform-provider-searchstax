package provider

import (
	"context"
	"fmt"

	searchstaxClient "terraform-provider-searchstax/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func NewWebhooksDataSource() datasource.DataSource { return &webhooksDataSource{} }

type webhooksDataSource struct{ client *searchstaxClient.Client }

func (d *webhooksDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_webhooks"
}
func (d *webhooksDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{Attributes: map[string]schema.Attribute{
		"id":           schema.StringAttribute{Computed: true},
		"account_name": schema.StringAttribute{Required: true},
		"webhooks": schema.ListNestedAttribute{Computed: true, NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
			"id":     schema.Int64Attribute{Computed: true},
			"name":   schema.StringAttribute{Computed: true},
			"url":    schema.StringAttribute{Computed: true},
			"paused": schema.BoolAttribute{Computed: true},
		}}},
	}}
}
func (d *webhooksDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *webhooksDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state webhooksDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	list, err := d.client.GetWebhooks(state.AccountName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Unable to Read SearchStax webhooks", err.Error())
		return
	}
	state.ID = types.StringValue("placeholder")
	state.Webhooks = nil
	for _, w := range list.Results {
		state.Webhooks = append(state.Webhooks, webhookModel{
			WebhookID: types.Int64Value(w.ID),
			Name:      types.StringValue(w.Name),
			URL:       types.StringValue(w.URL),
			Paused:    types.BoolValue(w.Paused),
		})
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

type webhooksDataSourceModel struct {
	ID          types.String   `tfsdk:"id"`
	AccountName types.String   `tfsdk:"account_name"`
	Webhooks    []webhookModel `tfsdk:"webhooks"`
}
type webhookModel struct {
	WebhookID types.Int64  `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	URL       types.String `tfsdk:"url"`
	Paused    types.Bool   `tfsdk:"paused"`
}
