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

var (
	_ resource.ResourceWithImportState = &userResource{}
)

func NewUserResource() resource.Resource {
	return &userResource{}
}

type userResource struct {
	client *searchstaxClient.Client
}

func (r *userResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

func (r *userResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Computed: true},
			"email": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"role": schema.StringAttribute{
				Required: true,
			},
			"first_name": schema.StringAttribute{
				Optional: true,
			},
			"last_name": schema.StringAttribute{
				Optional: true,
			},
			// optional admin action
			"new_password": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
			},
		},
	}
}

func (r *userResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*searchstaxClient.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *searchstaxClient.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	r.client = client
}

func (r *userResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan userResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.InviteUser(searchstaxClient.InviteUserRequest{
		Email:     plan.Email.ValueString(),
		Role:      plan.Role.ValueString(),
		FirstName: plan.FirstName.ValueString(),
		LastName:  plan.LastName.ValueString(),
	}); err != nil {
		resp.Diagnostics.AddError("Error inviting SearchStax user", err.Error())
		return
	}

	if !plan.NewPassword.IsNull() && plan.NewPassword.ValueString() != "" {
		if err := r.client.ChangeUserPassword(searchstaxClient.ChangeUserPasswordRequest{
			Email:       plan.Email.ValueString(),
			NewPassword: plan.NewPassword.ValueString(),
		}); err != nil {
			resp.Diagnostics.AddError("Error setting SearchStax user password", err.Error())
			return
		}
	}

	plan.ID = types.StringValue(plan.Email.ValueString())
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *userResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state userResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	users, err := r.client.GetUsers()
	if err != nil {
		resp.Diagnostics.AddError("Error reading SearchStax users", err.Error())
		return
	}

	found := false
	for _, u := range users.Results {
		if strings.EqualFold(u.Email, state.Email.ValueString()) {
			state.Role = types.StringValue(u.Role)
			state.FirstName = types.StringValue(u.FirstName)
			state.LastName = types.StringValue(u.LastName)
			found = true
			break
		}
	}

	if !found {
		resp.State.RemoveResource(ctx)
		return
	}

	state.ID = types.StringValue(state.Email.ValueString())
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *userResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan userResourceModel
	var state userResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.Role.ValueString() != state.Role.ValueString() {
		if err := r.client.SetUserRole(searchstaxClient.SetUserRoleRequest{
			Email: plan.Email.ValueString(),
			Role:  plan.Role.ValueString(),
		}); err != nil {
			resp.Diagnostics.AddError("Error updating SearchStax user role", err.Error())
			return
		}
	}

	if !plan.NewPassword.IsNull() && plan.NewPassword.ValueString() != "" && plan.NewPassword.ValueString() != state.NewPassword.ValueString() {
		if err := r.client.ChangeUserPassword(searchstaxClient.ChangeUserPasswordRequest{
			Email:       plan.Email.ValueString(),
			NewPassword: plan.NewPassword.ValueString(),
		}); err != nil {
			resp.Diagnostics.AddError("Error updating SearchStax user password", err.Error())
			return
		}
	}

	plan.ID = types.StringValue(plan.Email.ValueString())
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *userResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state userResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.DeleteUser(searchstaxClient.DeleteUserRequest{Email: state.Email.ValueString()}); err != nil {
		resp.Diagnostics.AddError("Error deleting SearchStax user", err.Error())
		return
	}
}

func (r *userResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("email"), req.ID)...)
}

type userResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Email       types.String `tfsdk:"email"`
	Role        types.String `tfsdk:"role"`
	FirstName   types.String `tfsdk:"first_name"`
	LastName    types.String `tfsdk:"last_name"`
	NewPassword types.String `tfsdk:"new_password"`
}
