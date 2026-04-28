package provider

import (
	"context"
	"fmt"

	searchstaxClient "terraform-provider-searchstax/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func NewDeploymentUsersDataSource() datasource.DataSource { return &deploymentUsersDataSource{} }

type deploymentUsersDataSource struct{ client *searchstaxClient.Client }

func (d *deploymentUsersDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_deployment_users"
}

func (d *deploymentUsersDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{Attributes: map[string]schema.Attribute{
		"id":             schema.StringAttribute{Computed: true},
		"account_name":   schema.StringAttribute{Required: true},
		"deployment_uid": schema.StringAttribute{Required: true},
		"users": schema.ListNestedAttribute{Computed: true, NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
			"username": schema.StringAttribute{Computed: true},
			"role":     schema.StringAttribute{Computed: true},
		}}},
	}}
}

func (d *deploymentUsersDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *deploymentUsersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state deploymentUsersDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	list, err := d.client.GetDeploymentUsers(state.AccountName.ValueString(), state.DeploymentUID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Unable to Read SearchStax deployment users", err.Error())
		return
	}
	state.ID = types.StringValue("placeholder")
	state.Users = nil
	for _, u := range list.Results {
		state.Users = append(state.Users, deploymentUserDataModel{
			Username: types.StringValue(u.Username),
			Role:     types.StringValue(u.Role),
		})
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

type deploymentUsersDataSourceModel struct {
	ID            types.String              `tfsdk:"id"`
	AccountName   types.String              `tfsdk:"account_name"`
	DeploymentUID types.String              `tfsdk:"deployment_uid"`
	Users         []deploymentUserDataModel `tfsdk:"users"`
}
type deploymentUserDataModel struct {
	Username types.String `tfsdk:"username"`
	Role     types.String `tfsdk:"role"`
}
