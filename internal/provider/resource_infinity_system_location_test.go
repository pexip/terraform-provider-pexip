package provider

import (
	"crypto/tls"
	"github.com/stretchr/testify/require"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/pexip/go-infinity-sdk/v38"
	"github.com/pexip/terraform-provider-pexip/internal/test"
)

func TestInfinitySystemLocation(t *testing.T) {
	_ = os.Setenv("TF_ACC", "1")

	// Create a mock client and set up expectations
	//client := infinity.NewClientMock()
	client, err := infinity.New(
		infinity.WithBaseURL("https://dev-manager.dev.pexip.network"),
		infinity.WithBasicAuth("admin", "admin"),
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
	require.NoError(t, err)

	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: getTestProtoV5ProviderFactories(client),
		Steps: []resource.TestStep{
			{
				Config: test.LoadTestData(t, "resource_infinity_system_location_basic.tf"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_system_location.main-location", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_system_location.main-location", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_system_location.main-location", "name", "main"),
					resource.TestCheckResourceAttr("pexip_infinity_system_location.main-location", "description", "Main location for Pexip Infinity System"),
					resource.TestCheckResourceAttr("pexip_infinity_system_location.main-location", "mtu", "1460"),
					resource.TestCheckResourceAttr("pexip_infinity_system_location.main-location", "dns_servers.#", "2"),
					resource.TestCheckResourceAttr("pexip_infinity_system_location.main-location", "dns_servers.0", "/api/admin/configuration/v1/dns_server/1/"),
					resource.TestCheckResourceAttr("pexip_infinity_system_location.main-location", "dns_servers.1", "/api/admin/configuration/v1/dns_server/2/")),
			},
			{
				Config: test.LoadTestData(t, "resource_infinity_system_location_basic_updated.tf"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_system_location.main-location", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_system_location.main-location", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_system_location.main-location", "name", "main"),
					resource.TestCheckResourceAttr("pexip_infinity_system_location.main-location", "description", "Main location for Pexip Infinity System - updated"),
					resource.TestCheckResourceAttr("pexip_infinity_system_location.main-location", "mtu", "1460"),
					resource.TestCheckResourceAttr("pexip_infinity_system_location.main-location", "dns_servers.#", "1"),
					resource.TestCheckResourceAttr("pexip_infinity_system_location.main-location", "dns_servers.0", "/api/admin/configuration/v1/dns_server/1/"),
				),
			},
		},
	})
}
