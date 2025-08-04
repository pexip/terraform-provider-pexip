package provider

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/pexip/go-infinity-sdk/v38"
	"github.com/pexip/go-infinity-sdk/v38/command"
	"github.com/pexip/go-infinity-sdk/v38/config"
	"github.com/pexip/go-infinity-sdk/v38/history"
	"github.com/pexip/go-infinity-sdk/v38/interfaces"
	"github.com/pexip/go-infinity-sdk/v38/status"
	"github.com/pexip/terraform-provider-pexip/internal/provider/validators"
	"github.com/pexip/terraform-provider-pexip/internal/version"
)

var (
	_ provider.Provider = (*PexipProvider)(nil)
)

type PexipProviderModel struct {
	Address  types.String `tfsdk:"address"`
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
	Insecure types.Bool   `tfsdk:"insecure"`
}

type PexipProvider struct {
	Address string
	Mutex   *sync.Mutex
	client  InfinityClient
}

type InfinityClient interface {
	interfaces.HTTPClient
	Config() *config.Service
	Status() *status.Service
	History() *history.Service
	Command() *command.Service
}

func New() provider.Provider {
	return &PexipProvider{
		Mutex: &sync.Mutex{},
	}
}

func newTestProvider(client InfinityClient) provider.Provider {
	return &PexipProvider{
		Mutex:  &sync.Mutex{},
		client: client,
	}
}

func (p *PexipProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "pexip"
	resp.Version = version.Version().String()
}

func (p *PexipProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "The Pexip Terraform provider exposes data sources and resources to deploy Pexip Infinity.",
		Attributes: map[string]schema.Attribute{
			"address": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					validators.URL(true),
				},
				MarkdownDescription: "URL of the Infinity Manager API, e.g. https://infinity.example.com",
			},
			"username": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(4),
				},
				MarkdownDescription: "Pexip Infinity Manager username for authentication.",
			},
			"password": schema.StringAttribute{
				Required:  true,
				Sensitive: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(4),
				},
				MarkdownDescription: "Pexip Infinity Manager password for authentication.",
			},
			"insecure": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Trust self-signed or otherwise invalid certificates. Defaults to `false`.",
			},
		},
	}
}

func (p *PexipProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data PexipProviderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if p.client == nil {
		var err error
		userAgent := fmt.Sprintf("terraform-provider-pexip/%s", version.Version().String())

		p.client, err = infinity.New(
			infinity.WithBaseURL(data.Address.ValueString()),
			infinity.WithBasicAuth(data.Username.ValueString(), data.Password.ValueString()),
			infinity.WithUserAgent(userAgent),
			infinity.WithMaxRetries(2),
			infinity.WithTransport(&http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: data.Insecure.ValueBool(),
					MinVersion:         tls.VersionTLS12,
				},
				MaxIdleConns:        30,
				MaxIdleConnsPerHost: 5,
				IdleConnTimeout:     60 * time.Second,
			}),
		)
		if err != nil {
			resp.Diagnostics.AddError(
				"Failed to create Infinity SDK client",
				fmt.Sprintf("Could not create Infinity SDK client: %s", err),
			)
			return
		}
	}

	// Pass the configured provider to resources and data sources.
	resp.DataSourceData = p
	resp.ResourceData = p
}

func (p *PexipProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		func() resource.Resource { return &InfinitySSHPasswordHashResource{} },
		func() resource.Resource { return &InfinityWebPasswordHashResource{} },
		func() resource.Resource { return &InfinityDnsServerResource{} },
		func() resource.Resource { return &InfinityNtpServerResource{} },
		func() resource.Resource { return &InfinitySystemLocationResource{} },
		func() resource.Resource { return &InfinityTeamsProxyResource{} },
		func() resource.Resource { return &InfinityEventSinkResource{} },
		func() resource.Resource { return &InfinityPolicyServerResource{} },
		func() resource.Resource { return &InfinitySIPProxyResource{} },
		func() resource.Resource { return &InfinitySTUNServerResource{} },
		func() resource.Resource { return &InfinityTURNServerResource{} },
		func() resource.Resource { return &InfinityHTTPProxyResource{} },
		func() resource.Resource { return &InfinityMSSIPProxyResource{} },
		func() resource.Resource { return &InfinityRoleResource{} },
		func() resource.Resource { return &InfinityUserGroupResource{} },
		func() resource.Resource { return &InfinityEndUserResource{} },
		func() resource.Resource { return &InfinityLdapSyncSourceResource{} },
		func() resource.Resource { return &InfinityConferenceResource{} },
		func() resource.Resource { return &InfinityConferenceAliasResource{} },
		func() resource.Resource { return &InfinityLocationResource{} },
		func() resource.Resource { return &InfinityDeviceResource{} },
		func() resource.Resource { return &InfinityIvrThemeResource{} },
		func() resource.Resource { return &InfinityGatewayRoutingRuleResource{} },
		func() resource.Resource { return &InfinityADFSAuthServerResource{} },
		func() resource.Resource { return &InfinityIdentityProviderResource{} },
		func() resource.Resource { return &InfinityAutomaticParticipantResource{} },
		func() resource.Resource { return &InfinityCACertificateResource{} },
		func() resource.Resource { return &InfinitySSHAuthorizedKeyResource{} },
		func() resource.Resource { return &InfinityScheduledConferenceResource{} },
		func() resource.Resource { return &InfinitySystemTuneableResource{} },
		func() resource.Resource { return &InfinityCertificateSigningRequestResource{} },
		func() resource.Resource { return &InfinityWorkerVMResource{} },
		func() resource.Resource { return &InfinityTLSCertificateResource{} },
		func() resource.Resource { return &InfinityLicenceResource{} },
		func() resource.Resource { return &InfinityStaticRouteResource{} },
		func() resource.Resource { return &InfinityRegistrationResource{} },
		func() resource.Resource { return &InfinityRecurringConferenceResource{} },
		func() resource.Resource { return &InfinityGlobalConfigurationResource{} },
		func() resource.Resource { return &InfinityGMSAccessTokenResource{} },
		func() resource.Resource { return &InfinityMediaLibraryEntryResource{} },
		func() resource.Resource { return &InfinityAuthenticationResource{} },
		func() resource.Resource { return &InfinityUserGroupEntityMappingResource{} },
		func() resource.Resource { return &InfinityLogLevelResource{} },
		func() resource.Resource { return &InfinityMediaLibraryPlaylistResource{} },
		func() resource.Resource { return &InfinityManagementVMResource{} },
		func() resource.Resource { return &InfinityMsExchangeConnectorResource{} },
		func() resource.Resource { return &InfinitySnmpNetworkManagementSystemResource{} },
		func() resource.Resource { return &InfinitySMTPServerResource{} },
		func() resource.Resource { return &InfinityIdentityProviderGroupResource{} },
		func() resource.Resource { return &InfinityIdentityProviderAttributeResource{} },
		func() resource.Resource { return &InfinitySIPCredentialResource{} },
		func() resource.Resource { return &InfinityH323GatekeeperResource{} },
		func() resource.Resource { return &InfinityLdapRoleResource{} },
		func() resource.Resource { return &InfinitySyslogServerResource{} },
		func() resource.Resource { return &InfinityRoleMappingResource{} },
		func() resource.Resource { return &InfinityOAuth2ClientResource{} },
		func() resource.Resource { return &InfinityDiagnosticGraphResource{} },
		func() resource.Resource { return &InfinitySystemSyncpointResource{} },
		func() resource.Resource { return &InfinityUpgradeResource{} },
		func() resource.Resource { return &InfinityAzureTenantResource{} },
		func() resource.Resource { return &InfinityWebappBrandingResource{} },
		func() resource.Resource { return &InfinityBreakInAllowListAddressResource{} },
		func() resource.Resource { return &InfinityMjxEndpointResource{} },
		func() resource.Resource { return &InfinityLdapSyncFieldResource{} },
		func() resource.Resource { return &InfinityMediaProcessingServerResource{} },
		func() resource.Resource { return &InfinityPexipStreamingCredentialResource{} },
		func() resource.Resource { return &InfinityScheduledAliasResource{} },
		func() resource.Resource { return &InfinityScheduledScalingResource{} },
		func() resource.Resource { return &InfinityWebappAliasResource{} },
		func() resource.Resource { return &InfinityLicenceRequestResource{} },
		func() resource.Resource { return &InfinityMjxIntegrationResource{} },
		func() resource.Resource { return &InfinityExternalWebappHostResource{} },
		func() resource.Resource { return &InfinityGoogleAuthServerResource{} },
	}
}

func (p *PexipProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		func() datasource.DataSource {
			return &InfinityManagerConfigDataSource{}
		},
	}
}
