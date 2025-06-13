package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/pexip/go-infinity-sdk/v38"
	"github.com/pexip/terraform-provider-pexip/internal/provider/validators"
	"sync"
)

var (
	_ provider.Provider = (*pexipProvider)(nil)
)

type pexipProvider struct {
	Address        string
	Mutex          *sync.Mutex
	InfinityClient *infinity.Client
}

func New() provider.Provider {
	return &pexipProvider{}
}

func (p *pexipProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "pexip"
}

func (p *pexipProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "TThe Pexip Terraform provider exposes data sources and resources to deploy Pexip Infinity.",
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

func (p *pexipProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
}

func (p *pexipProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		//func() resource.Resource {
		//	return &infinityNodeResource{}
		//},
	}
}

func (p *pexipProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		func() datasource.DataSource {
			return &infinityManagerConfigDataSource{}
		},
	}
}

//func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
//	// Warning or errors can be collected in a slice type
//	var diags diag.Diagnostics
//
//	// load provider config vars
//	address := helpers.ResourceToString(d, "address")
//	username := helpers.ResourceToString(d, "username")
//	password := helpers.ResourceToString(d, "password")
//
//	// Initialize the InfinityClient SDK client with the base URL and authentication
//	client, err := infinity.New(
//		infinity.WithBaseURL(address),
//		infinity.WithBasicAuth(username, password),
//		infinity.WithMaxRetries(2),
//	)
//	if err != nil {
//		log.Error().Err(err).Msg("failed to create Infinity SDK client")
//	}
//
//	var mut sync.Mutex
//	conf := providerConfiguration{
//		Address:        address,
//		Mutex:          &mut,
//		InfinityClient: client,
//	}
//	return conf, diags
//}
