package provider

import (
	"context"
	"fmt"

	searchstaxClient "terraform-provider-searchstax/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func NewDeploymentServersDataSource() datasource.DataSource { return &deploymentServersDataSource{} }

type deploymentServersDataSource struct{ client *searchstaxClient.Client }

func (d *deploymentServersDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_deployment_servers"
}

func (d *deploymentServersDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{Attributes: map[string]schema.Attribute{
		"id":             schema.StringAttribute{Computed: true},
		"account_name":   schema.StringAttribute{Required: true},
		"deployment_uid": schema.StringAttribute{Required: true},
		"servers": schema.ListNestedAttribute{Computed: true, NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
			"sn":              schema.Int64Attribute{Computed: true},
			"node":            schema.StringAttribute{Computed: true},
			"private_address": schema.StringAttribute{Computed: true},
			"dns_address":     schema.StringAttribute{Computed: true},
			"status":          schema.StringAttribute{Computed: true},
			"status_details":  schema.StringAttribute{Computed: true},
			"role":            schema.StringAttribute{Computed: true},
			"solr":            schema.BoolAttribute{Computed: true},
			"zookeeper":       schema.BoolAttribute{Computed: true},
		}}},
	}}
}

func (d *deploymentServersDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *deploymentServersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state deploymentServersDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	out, err := d.client.GetDeploymentServers(state.AccountName.ValueString(), state.DeploymentUID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Unable to read deployment servers", err.Error())
		return
	}
	state.ID = types.StringValue("placeholder")
	state.Servers = nil
	for _, s := range out.Results {
		state.Servers = append(state.Servers, deploymentServerModel{
			SN:             types.Int64Value(s.SN),
			Node:           types.StringValue(s.Node),
			PrivateAddress: types.StringValue(s.PrivateAddress),
			DNSAddress:     types.StringValue(s.DNSAddress),
			Status:         types.StringValue(s.Status),
			StatusDetails:  types.StringValue(s.StatusDetails),
			Role:           types.StringValue(s.Role),
			Solr:           types.BoolValue(s.Solr),
			Zookeeper:      types.BoolValue(s.Zookeeper),
		})
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

type deploymentServersDataSourceModel struct {
	ID            types.String            `tfsdk:"id"`
	AccountName   types.String            `tfsdk:"account_name"`
	DeploymentUID types.String            `tfsdk:"deployment_uid"`
	Servers       []deploymentServerModel `tfsdk:"servers"`
}

type deploymentServerModel struct {
	SN             types.Int64  `tfsdk:"sn"`
	Node           types.String `tfsdk:"node"`
	PrivateAddress types.String `tfsdk:"private_address"`
	DNSAddress     types.String `tfsdk:"dns_address"`
	Status         types.String `tfsdk:"status"`
	StatusDetails  types.String `tfsdk:"status_details"`
	Role           types.String `tfsdk:"role"`
	Solr           types.Bool   `tfsdk:"solr"`
	Zookeeper      types.Bool   `tfsdk:"zookeeper"`
}
