package provider

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"terraform-provider-zosmf/zosmf"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ provider.Provider = &zosmfProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &zosmfProvider{
			version: version,
		}
	}
}

type zosmfProviderModel struct {
	Host     types.String `tfsdk:"host"`
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
	Insecure types.Bool   `tfsdk:"insecure"`
}

// zosmfProvider is the provider implementation.
type zosmfProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// Metadata returns the provider type name.
func (p *zosmfProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "zosmf"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *zosmfProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				Optional: false,
				Required: true,
			},
			"username": schema.StringAttribute{
				Optional: false,
				Required: true,
			},
			"password": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
			},
			"insecure": schema.StringAttribute{
				Optional: true,
			},
		},
	}
}

// Configure prepares a zosmf API client for data sources and resources.
func (p *zosmfProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config zosmfProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	zos_password := os.Getenv("ZOS_PASSWORD")

	if !config.Password.IsNull() {
		zos_password = config.Password.ValueString()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if zos_password == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Missing Zosmf password",
			"Password is needed",
		)
	}

	client, err := zosmf.NewClient(config.Host.ValueString(), config.Username.ValueString(), config.Password.ValueString(), config.Insecure.ValueBool())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create zosmf API Client",
			"An unexpected error occurred when creating the zosmf API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"Zosmf Client Error: "+err.Error(),
		)
		return
	}
	tflog.Warn(ctx, fmt.Sprintf("the username %s", client.Username))
	resp.DataSourceData = client
	resp.ResourceData = client
}

// DataSources defines the data sources implemented in the provider.
func (p *zosmfProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewDatasetDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *zosmfProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewDatasetResource,
	}
}
