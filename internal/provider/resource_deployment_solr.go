package provider

import (
	"context"
	"fmt"

	searchstaxClient "terraform-provider-searchstax/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func NewDeploymentSolrResource() resource.Resource { return &deploymentSolrResource{} }

type deploymentSolrResource struct{ client *searchstaxClient.Client }

func (r *deploymentSolrResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_deployment_solr"
}

func (r *deploymentSolrResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{Attributes: map[string]schema.Attribute{
		"id":             schema.StringAttribute{Computed: true},
		"account_name":   schema.StringAttribute{Required: true},
		"deployment_uid": schema.StringAttribute{Required: true},
		"node":           schema.StringAttribute{Required: true},
		"action": schema.StringAttribute{
			Required: true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
	}}
}

func (r *deploymentSolrResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *deploymentSolrResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan deploymentSolrResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var err error
	switch plan.Action.ValueString() {
	case "start":
		err = r.client.StartSolr(plan.AccountName.ValueString(), plan.DeploymentUID.ValueString(), plan.Node.ValueString())
	case "stop":
		err = r.client.StopSolr(plan.AccountName.ValueString(), plan.DeploymentUID.ValueString(), plan.Node.ValueString())
	default:
		resp.Diagnostics.AddError("Invalid action", "action must be start or stop")
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Error performing solr action", err.Error())
		return
	}
	plan.ID = types.StringValue(plan.AccountName.ValueString() + "/" + plan.DeploymentUID.ValueString() + "/" + plan.Node.ValueString() + "/" + plan.Action.ValueString())
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *deploymentSolrResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state deploymentSolrResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *deploymentSolrResource) Update(context.Context, resource.UpdateRequest, *resource.UpdateResponse) {}

func (r *deploymentSolrResource) Delete(context.Context, resource.DeleteRequest, *resource.DeleteResponse) {}

type deploymentSolrResourceModel struct {
	ID            types.String `tfsdk:"id"`
	AccountName   types.String `tfsdk:"account_name"`
	DeploymentUID types.String `tfsdk:"deployment_uid"`
	Node          types.String `tfsdk:"node"`
	Action        types.String `tfsdk:"action"`
}
