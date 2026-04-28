package provider

import (
	"context"
	"fmt"

	searchstaxClient "terraform-provider-searchstax/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &usersDataSource{}
	_ datasource.DataSourceWithConfigure = &usersDataSource{}
)

func NewUsersDataSource() datasource.DataSource {
	return &usersDataSource{}
}

type usersDataSource struct {
	client *searchstaxClient.Client
}

func (d *usersDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_users"
}

func (d *usersDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Computed: true},
			"users": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":         schema.Int64Attribute{Computed: true},
						"email":      schema.StringAttribute{Computed: true},
						"role":       schema.StringAttribute{Computed: true},
						"first_name": schema.StringAttribute{Computed: true},
						"last_name":  schema.StringAttribute{Computed: true},
					},
				},
			},
		},
	}
}

func (d *usersDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*searchstaxClient.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *searchstaxClient.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	d.client = client
}

func (d *usersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state usersDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	users, err := d.client.GetUsers()
	if err != nil {
		resp.Diagnostics.AddError("Unable to Read SearchStax Users", err.Error())
		return
	}

	state.ID = types.StringValue("placeholder")
	state.Users = nil
	for _, u := range users.Results {
		state.Users = append(state.Users, userModel{
			ID:        types.Int64Value(u.ID),
			Email:     types.StringValue(u.Email),
			Role:      types.StringValue(u.Role),
			FirstName: types.StringValue(u.FirstName),
			LastName:  types.StringValue(u.LastName),
		})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

type usersDataSourceModel struct {
	ID    types.String `tfsdk:"id"`
	Users []userModel  `tfsdk:"users"`
}

type userModel struct {
	ID        types.Int64  `tfsdk:"id"`
	Email     types.String `tfsdk:"email"`
	Role      types.String `tfsdk:"role"`
	FirstName types.String `tfsdk:"first_name"`
	LastName  types.String `tfsdk:"last_name"`
}
