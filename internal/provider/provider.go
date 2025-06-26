package provider

import (
	"context"
	"crypto/tls"
	"fmt"
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
	"net/http"
	"sync"
	"time"
)

var (
	_ provider.Provider = (*PexipProvider)(nil)
)

type PexipProviderModel struct {
	Address  types.String `tfsdk:"address"`
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
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
				MarkdownDescription: "Pexip Infinity Manager username for authentication",
			},
			"password": schema.StringAttribute{
				Required:  true,
				Sensitive: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(4),
				},
				MarkdownDescription: "Pexip Infinity Manager password for authentication",
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
		p.client, err = infinity.New(
			infinity.WithBaseURL(data.Address.ValueString()),
			infinity.WithBasicAuth(data.Username.ValueString(), data.Password.ValueString()),
			infinity.WithMaxRetries(2),
			infinity.WithTransport(&http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true, // We need this because default certificate is not trusted
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

	p.Address = data.Address.ValueString()
	p.Mutex = &sync.Mutex{}

	resp.DataSourceData = p //TODO check if this is correct or if we should use a dedicated structure
	resp.ResourceData = p   //TODO check if this is correct or if we should use a dedicated structure
}

func (p *PexipProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		func() resource.Resource {
			return &InfinityNodeResource{}
		},
		func() resource.Resource { return &InfinitySSHPasswordHashResource{} },
		func() resource.Resource { return &InfinityWebPasswordHashResource{} },
		func() resource.Resource { return &InfinityDnsServerResource{} },
	}
}

func (p *PexipProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		func() datasource.DataSource {
			return &InfinityManagerConfigDataSource{}
		},
	}
}
