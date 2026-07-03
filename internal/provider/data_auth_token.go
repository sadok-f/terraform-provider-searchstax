package provider

import (
	"context"
	"fmt"

	searchstaxClient "terraform-provider-searchstax/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func NewAuthTokenDataSource() datasource.DataSource { return &authTokenDataSource{} }

type authTokenDataSource struct{ client *searchstaxClient.Client }

func (d *authTokenDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_auth_token"
}

func (d *authTokenDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{Attributes: map[string]schema.Attribute{
		"id":                 schema.StringAttribute{Computed: true},
		"token":              schema.StringAttribute{Computed: true, Sensitive: true},
		"valid":              schema.BoolAttribute{Computed: true},
		"expires_in_seconds": schema.Int64Attribute{Computed: true},
		"token_expires_at":   schema.StringAttribute{Computed: true},
	}}
}

func (d *authTokenDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *authTokenDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	token, err := d.client.SignIn()
	if err != nil {
		resp.Diagnostics.AddError("Unable to obtain SearchStax auth token", err.Error())
		return
	}
	state := authTokenDataSourceModel{
		ID:    types.StringValue("placeholder"),
		Token: types.StringValue(token.Token),
	}
	verify, err := d.client.VerifyAuthToken()
	if err == nil {
		state.Valid = types.BoolValue(verify.Valid)
		if verify.TokenExpiresInSeconds > 0 {
			state.ExpiresInSeconds = types.Int64Value(verify.TokenExpiresInSeconds)
		} else if verify.Token.DurationDays > 0 {
			// The real API reports the token lifetime as duration_days rather
			// than an explicit seconds value.
			state.ExpiresInSeconds = types.Int64Value(int64(verify.Token.DurationDays) * 86400)
		}
		if verify.Token.Expires != "" {
			state.TokenExpiresAt = types.StringValue(verify.Token.Expires)
		}
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

type authTokenDataSourceModel struct {
	ID               types.String `tfsdk:"id"`
	Token            types.String `tfsdk:"token"`
	Valid            types.Bool   `tfsdk:"valid"`
	ExpiresInSeconds types.Int64  `tfsdk:"expires_in_seconds"`
	TokenExpiresAt   types.String `tfsdk:"token_expires_at"`
}
