package provider

import (
	"context"
	"os"
	searchstaxClient "terraform-provider-searchstax/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ provider.Provider = &searchstaxProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &searchstaxProvider{
			version: version,
		}
	}
}

// searchstaxProvider is the provider implementation.
type searchstaxProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// searchstaxProviderModel maps provider schema data to a Go type.
type searchstaxProviderModel struct {
	Host     types.String `tfsdk:"host"`
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
}

// Metadata returns the provider type name.
func (p *searchstaxProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "searchstax"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *searchstaxProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				Optional: true,
			},
			"username": schema.StringAttribute{
				Optional: true,
			},
			"password": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
			},
		},
	}
}

// Configure prepares a SearchStax API client for data sources and resources.
func (p *searchstaxProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {

	tflog.Info(ctx, "Configuring SearchStax client")

	// Retrieve provider data from configuration
	var config searchstaxProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.

	if config.Host.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Unknown SearchStax API Host",
			"The provider cannot create the SearchStax API client as there is an unknown configuration value for the SearchStax API host. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the SEARCHSTAX_HOST environment variable.",
		)
	}

	if config.Username.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Unknown SearchStax API Username",
			"The provider cannot create the SearchStax API client as there is an unknown configuration value for the SearchStax API username. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the SEARCHSTAX_USERNAME environment variable.",
		)
	}

	if config.Password.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Unknown SearchStax API Password",
			"The provider cannot create the SearchStax API client as there is an unknown configuration value for the SearchStax API password. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the SEARCHSTAX_PASSWORD environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	host := os.Getenv("SEARCHSTAX_HOST")
	username := os.Getenv("SEARCHSTAX_USERNAME")
	password := os.Getenv("SEARCHSTAX_PASSWORD")

	if !config.Host.IsNull() {
		host = config.Host.ValueString()
	}

	if !config.Username.IsNull() {
		username = config.Username.ValueString()
	}

	if !config.Password.IsNull() {
		password = config.Password.ValueString()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if username == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Missing SearchStax API Username",
			"The provider cannot create the SearchStax API client as there is a missing or empty value for the SearchStax API username. "+
				"Set the username value in the configuration or use the SEARCHSTAX_USERNAME environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if password == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Missing SearchStax API Password",
			"The provider cannot create the SearchStax API client as there is a missing or empty value for the SearchStax API password. "+
				"Set the password value in the configuration or use the SEARCHSTAX_PASSWORD environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Create a new SearchStax client using the configuration values
	client, err := searchstaxClient.NewClient(&host, &username, &password)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create SearchStax API Client",
			"An unexpected error occurred when creating the SearchStax API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"SearchStax Client Error: "+err.Error(),
		)
		return
	}

	// Make the SearchStax client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client
}

// DataSources defines the data sources implemented in the provider.
func (p *searchstaxProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewDeploymentsDataSource,
		NewPrivateVpcDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *searchstaxProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewDeploymentResource,
	}
}
