package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/pexip/go-infinity-sdk/v38"
	"github.com/pexip/terraform-provider-pexip/internal/helpers"
	"github.com/rs/zerolog/log"
	"sync"
)

type providerConfiguration struct {
	Address        string
	Mutex          *sync.Mutex
	InfinityClient *infinity.Client
}

// Provider represents a terraform provider definition
func Provider() *schema.Provider {
	return New()
}

// New represents a terraform provider definition
func New() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"address": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("PEXIP_INFINITY__MANAGER_ADDRESS", nil),
				ValidateFunc: validation.All(
					validation.NoZeroValues,
					validation.IsURLWithHTTPS,
				),
				Description: "URL of the Infinity Manager API, e.g. https://infinity.example.com",
			},
			"username": {
				Type:         schema.TypeString,
				Required:     true,
				DefaultFunc:  schema.EnvDefaultFunc("PEXIP_INFINITY__MANAGER_USERNAME", nil),
				ValidateFunc: validation.NoZeroValues,
				Description:  "Pexip Infinity Manager username for authentication",
			},
			"password": {
				Type:         schema.TypeString,
				Required:     true,
				DefaultFunc:  schema.EnvDefaultFunc("PEXIP_INFINITY__MANAGER_PASSWORD", nil),
				ValidateFunc: validation.NoZeroValues,
				Description:  "Pexip Infinity Manager password for authentication",
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"infinity_manager_config": dataSourceInfinityManagerConfig(),
			"infinity_manager":        dataSourceInfinityManager(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"infinity_node": infinityNodeResourceQuery(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	// load provider config vars
	address := helpers.ResourceToString(d, "address")
	username := helpers.ResourceToString(d, "username")
	password := helpers.ResourceToString(d, "password")

	// Initialize the InfinityClient SDK client with the base URL and authentication
	client, err := infinity.New(
		infinity.WithBaseURL(address),
		infinity.WithBasicAuth(username, password),
		infinity.WithMaxRetries(2),
	)
	if err != nil {
		log.Error().Err(err).Msg("failed to create Infinity SDK client")
	}

	var mut sync.Mutex
	conf := providerConfiguration{
		Address:        address,
		Mutex:          &mut,
		InfinityClient: client,
	}
	return conf, diags
}
