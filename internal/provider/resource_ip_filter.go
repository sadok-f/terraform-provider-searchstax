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

func NewIPFilterResource() resource.Resource { return &ipFilterResource{} }

type ipFilterResource struct{ client *searchstaxClient.Client }

func (r *ipFilterResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ip_filter"
}
func (r *ipFilterResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{Attributes: map[string]schema.Attribute{
		"id":             schema.StringAttribute{Computed: true},
		"account_name":   schema.StringAttribute{Required: true},
		"deployment_uid": schema.StringAttribute{Required: true},
		"cidr_ip":        schema.StringAttribute{Required: true},
		"description":    schema.StringAttribute{Optional: true},
		"services":       schema.ListAttribute{Required: true, ElementType: types.StringType},
	}}
}
func (r *ipFilterResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ipFilterResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ipFilterResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var services []string
	resp.Diagnostics.Append(plan.Services.ElementsAs(ctx, &services, false)...)
	if resp.Diagnostics.HasError() {
		return
	}
	err := r.client.AddIPFilter(plan.AccountName.ValueString(), plan.DeploymentUID.ValueString(), searchstaxClient.IPFilterUpsertRequest{
		CIDRIP:      plan.CIDRIP.ValueString(),
		Description: plan.Description.ValueString(),
		Services:    services,
	})
	if err != nil {
		resp.Diagnostics.AddError("Error creating IP filter", err.Error())
		return
	}
	plan.ID = types.StringValue(plan.AccountName.ValueString() + "/" + plan.DeploymentUID.ValueString() + "/" + plan.CIDRIP.ValueString())
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}
func (r *ipFilterResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ipFilterResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	list, err := r.client.GetIPFilters(state.AccountName.ValueString(), state.DeploymentUID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading IP filters", err.Error())
		return
	}
	for _, f := range list.Results {
		if f.CIDRIP == state.CIDRIP.ValueString() {
			services, diags := types.ListValueFrom(ctx, types.StringType, f.Services)
			resp.Diagnostics.Append(diags...)
			state.Services = services
			state.Description = types.StringValue(f.Description)
			state.ID = types.StringValue(state.AccountName.ValueString() + "/" + state.DeploymentUID.ValueString() + "/" + state.CIDRIP.ValueString())
			resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
			return
		}
	}
	resp.State.RemoveResource(ctx)
}
func (r *ipFilterResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ipFilterResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	var services []string
	resp.Diagnostics.Append(plan.Services.ElementsAs(ctx, &services, false)...)
	if resp.Diagnostics.HasError() {
		return
	}
	err := r.client.UpdateIPFilter(plan.AccountName.ValueString(), plan.DeploymentUID.ValueString(), searchstaxClient.IPFilterUpsertRequest{
		CIDRIP:      plan.CIDRIP.ValueString(),
		Description: plan.Description.ValueString(),
		Services:    services,
	})
	if err != nil {
		resp.Diagnostics.AddError("Error updating IP filter", err.Error())
		return
	}
	plan.ID = types.StringValue(plan.AccountName.ValueString() + "/" + plan.DeploymentUID.ValueString() + "/" + plan.CIDRIP.ValueString())
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}
func (r *ipFilterResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ipFilterResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	err := r.client.DeleteIPFilter(state.AccountName.ValueString(), state.DeploymentUID.ValueString(), searchstaxClient.IPFilterDeleteRequest{
		CIDRIP: state.CIDRIP.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error deleting IP filter", err.Error())
	}
}
func (r *ipFilterResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, "/")
	if len(parts) != 3 {
		resp.Diagnostics.AddError("Unexpected Import Identifier", "Expected account_name/deployment_uid/cidr_ip")
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("account_name"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("deployment_uid"), parts[1])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("cidr_ip"), parts[2])...)
}

type ipFilterResourceModel struct {
	ID            types.String `tfsdk:"id"`
	AccountName   types.String `tfsdk:"account_name"`
	DeploymentUID types.String `tfsdk:"deployment_uid"`
	CIDRIP        types.String `tfsdk:"cidr_ip"`
	Description   types.String `tfsdk:"description"`
	Services      types.List   `tfsdk:"services"`
}
