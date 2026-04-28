package provider

import (
	"context"
	"fmt"
	"strings"

	searchstaxClient "terraform-provider-searchstax/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.ResourceWithImportState = &dnsRecordResource{}

func NewDNSRecordResource() resource.Resource { return &dnsRecordResource{} }

type dnsRecordResource struct{ client *searchstaxClient.Client }

func (r *dnsRecordResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dns_record"
}

func (r *dnsRecordResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{Attributes: map[string]schema.Attribute{
		"id":           schema.StringAttribute{Computed: true},
		"account_name": schema.StringAttribute{Required: true, PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
		"name":         schema.StringAttribute{Required: true, PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
		"deployment":   schema.StringAttribute{Required: true},
		"ttl":          schema.StringAttribute{Optional: true, Computed: true},
	}}
}

func (r *dnsRecordResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*searchstaxClient.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Resource Configure Type", fmt.Sprintf("Expected *searchstaxClient.Client, got: %T.", req.ProviderData))
		return
	}
	r.client = client
}

func (r *dnsRecordResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan dnsRecordResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	record, err := r.client.AssociateDNSRecord(plan.AccountName.ValueString(), plan.Name.ValueString(), searchstaxClient.AssociateDNSRecordRequest{
		Deployment: plan.Deployment.ValueString(),
		TTL:        plan.TTL.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error associating DNS record", err.Error())
		return
	}
	plan.ID = types.StringValue(plan.AccountName.ValueString() + "/" + plan.Name.ValueString())
	plan.Deployment = types.StringValue(record.Deployment)
	plan.TTL = types.StringValue(record.TTL)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *dnsRecordResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state dnsRecordResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	record, err := r.client.GetDNSRecord(state.AccountName.ValueString(), state.Name.ValueString())
	if err != nil {
		resp.State.RemoveResource(ctx)
		return
	}
	state.Deployment = types.StringValue(record.Deployment)
	state.TTL = types.StringValue(record.TTL)
	state.ID = types.StringValue(state.AccountName.ValueString() + "/" + state.Name.ValueString())
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *dnsRecordResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan dnsRecordResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	record, err := r.client.AssociateDNSRecord(plan.AccountName.ValueString(), plan.Name.ValueString(), searchstaxClient.AssociateDNSRecordRequest{
		Deployment: plan.Deployment.ValueString(),
		TTL:        plan.TTL.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error updating DNS record", err.Error())
		return
	}
	plan.ID = types.StringValue(plan.AccountName.ValueString() + "/" + plan.Name.ValueString())
	plan.Deployment = types.StringValue(record.Deployment)
	plan.TTL = types.StringValue(record.TTL)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *dnsRecordResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// no delete endpoint exposed by API for DNS alias; disassociate by setting blank deployment
	var state dnsRecordResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	_, _ = r.client.AssociateDNSRecord(state.AccountName.ValueString(), state.Name.ValueString(), searchstaxClient.AssociateDNSRecordRequest{Deployment: ""})
}

func (r *dnsRecordResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, "/")
	if len(parts) != 2 {
		resp.Diagnostics.AddError("Unexpected Import Identifier", "Expected account_name/name")
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("account_name"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), parts[1])...)
}

type dnsRecordResourceModel struct {
	ID          types.String `tfsdk:"id"`
	AccountName types.String `tfsdk:"account_name"`
	Name        types.String `tfsdk:"name"`
	Deployment  types.String `tfsdk:"deployment"`
	TTL         types.String `tfsdk:"ttl"`
}
