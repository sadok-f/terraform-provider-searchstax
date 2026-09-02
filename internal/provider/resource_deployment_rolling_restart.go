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

var _ resource.ResourceWithImportState = &deploymentRollingRestartResource{}

func NewDeploymentRollingRestartResource() resource.Resource {
	return &deploymentRollingRestartResource{}
}

type deploymentRollingRestartResource struct{ client *searchstaxClient.Client }

func (r *deploymentRollingRestartResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_deployment_rolling_restart"
}

func (r *deploymentRollingRestartResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{Attributes: map[string]schema.Attribute{
		"id":             schema.StringAttribute{Computed: true},
		"account_name":   schema.StringAttribute{Required: true},
		"deployment_uid": schema.StringAttribute{Required: true},
		"solr":           schema.BoolAttribute{Optional: true},
		"zookeeper":      schema.BoolAttribute{Optional: true},
		"triggers": schema.MapAttribute{
			ElementType: types.StringType,
			Optional:    true,
			Description: "Arbitrary map of values that, when changed, initiates a rolling restart. Use it to trigger a single restart when the custom jar list changes.",
		},
		"message": schema.StringAttribute{Computed: true},
	}}
}

func (r *deploymentRollingRestartResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *deploymentRollingRestartResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan deploymentRollingRestartResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	solr := true
	zookeeper := false
	if !plan.Solr.IsNull() {
		solr = plan.Solr.ValueBool()
	}
	if !plan.Zookeeper.IsNull() {
		zookeeper = plan.Zookeeper.ValueBool()
	}
	out, err := r.client.RollingRestart(plan.AccountName.ValueString(), plan.DeploymentUID.ValueString(), searchstaxClient.RollingRestartRequest{
		Solr:      solr,
		Zookeeper: zookeeper,
	})
	if err != nil {
		resp.Diagnostics.AddError("Error initiating rolling restart", err.Error())
		return
	}
	msg := out.Message
	if msg == "" {
		msg = out.Detail
	}
	plan.Message = types.StringValue(msg)
	plan.ID = types.StringValue(plan.AccountName.ValueString() + "/" + plan.DeploymentUID.ValueString())
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *deploymentRollingRestartResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state deploymentRollingRestartResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *deploymentRollingRestartResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: account_name/deployment_uid. Got: %q", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("account_name"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("deployment_uid"), parts[1])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
}

func (r *deploymentRollingRestartResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan deploymentRollingRestartResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state deploymentRollingRestartResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.Message.IsNull() {
		plan.ID = types.StringValue(plan.AccountName.ValueString() + "/" + plan.DeploymentUID.ValueString())
		plan.Message = types.StringValue("imported")
		resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
		return
	}

	solr := true
	zookeeper := false
	if !plan.Solr.IsNull() {
		solr = plan.Solr.ValueBool()
	}
	if !plan.Zookeeper.IsNull() {
		zookeeper = plan.Zookeeper.ValueBool()
	}
	out, err := r.client.RollingRestart(plan.AccountName.ValueString(), plan.DeploymentUID.ValueString(), searchstaxClient.RollingRestartRequest{
		Solr:      solr,
		Zookeeper: zookeeper,
	})
	if err != nil {
		resp.Diagnostics.AddError("Error initiating rolling restart", err.Error())
		return
	}
	msg := out.Message
	if msg == "" {
		msg = out.Detail
	}
	plan.Message = types.StringValue(msg)
	plan.ID = types.StringValue(plan.AccountName.ValueString() + "/" + plan.DeploymentUID.ValueString())
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *deploymentRollingRestartResource) Delete(context.Context, resource.DeleteRequest, *resource.DeleteResponse) {
}

type deploymentRollingRestartResourceModel struct {
	ID            types.String `tfsdk:"id"`
	AccountName   types.String `tfsdk:"account_name"`
	DeploymentUID types.String `tfsdk:"deployment_uid"`
	Solr          types.Bool   `tfsdk:"solr"`
	Zookeeper     types.Bool   `tfsdk:"zookeeper"`
	Triggers      types.Map    `tfsdk:"triggers"`
	Message       types.String `tfsdk:"message"`
}
