package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/pexip/go-infinity-sdk/v38"
	"github.com/pexip/go-infinity-sdk/v38/config"
	"github.com/pexip/go-infinity-sdk/v38/types"
	"github.com/stretchr/testify/mock"

	"github.com/pexip/terraform-provider-pexip/internal/test"
)

func TestInfinitySystemLocation(t *testing.T) {
	t.Parallel()
	_ = os.Setenv("TF_ACC", "1")

	// Create a mock client and set up expectations
	client := infinity.NewClientMock()

	// Mock the CreateSystemLocation API call
	createResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/system_location/123/",
	}
	client.On("PostWithResponse", mock.Anything, "configuration/v1/system_location/", mock.Anything, mock.Anything).Return(createResponse, nil)

	// Track state to return different values before and after update
	updated := false

	// Mock the GetSystemLocation API call for Read operations
	client.On("GetJSON", mock.Anything, "configuration/v1/system_location/123/", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		sysLoc := args.Get(2).(*config.SystemLocation)
		if updated {
			mediaQoS := 0
			signallingQoS := 0
			*sysLoc = config.SystemLocation{
				ID:            123,
				Name:          "main",
				Description:   "Main location for Pexip Infinity System - updated",
				MTU:           1460,
				MediaQoS:      &mediaQoS,
				SignallingQoS: &signallingQoS,
				DNSServers: []config.DNSServer{
					{ID: 1, ResourceURI: "/api/admin/configuration/v1/dns_server/1/"},
				},
				NTPServers: []config.NTPServer{
					{ID: 1, ResourceURI: "/api/admin/configuration/v1/ntp_server/1/"},
				},
				ClientSTUNServers: []config.STUNServer{
					{ID: 2},
				},
				ResourceURI: "/api/admin/configuration/v1/system_location/123/",
			}
		} else {
			mediaQoS := 0
			signallingQoS := 0
			*sysLoc = config.SystemLocation{
				ID:            123,
				Name:          "main",
				Description:   "Main location for Pexip Infinity System",
				MTU:           1460,
				MediaQoS:      &mediaQoS,
				SignallingQoS: &signallingQoS,
				DNSServers: []config.DNSServer{
					{ID: 1, ResourceURI: "/api/admin/configuration/v1/dns_server/1/"},
					{ID: 2, ResourceURI: "/api/admin/configuration/v1/dns_server/2/"},
				},
				NTPServers: []config.NTPServer{
					{ID: 1, ResourceURI: "/api/admin/configuration/v1/ntp_server/1/"},
				},
				ClientSTUNServers: []config.STUNServer{
					{ID: 1},
				},
				ResourceURI: "/api/admin/configuration/v1/system_location/123/",
			}
		}
	}).Maybe() // Called multiple times for reads

	// Mock the UpdateSystemLocation API call
	client.On("PutJSON", mock.Anything, mock.MatchedBy(func(path string) bool {
		return path == "configuration/v1/system_location/123/"
	}), mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		updated = true // Mark as updated for subsequent reads
		sysLoc := args.Get(3).(*config.SystemLocation)
		mediaQoS := 0
		signallingQoS := 0
		*sysLoc = config.SystemLocation{
			ID:            123,
			Name:          "main",
			Description:   "Main location for Pexip Infinity System - updated",
			MTU:           1460,
			MediaQoS:      &mediaQoS,
			SignallingQoS: &signallingQoS,
			DNSServers: []config.DNSServer{
				{ID: 1, ResourceURI: "/api/admin/configuration/v1/dns_server/1/"},
			},
			NTPServers: []config.NTPServer{
				{ID: 1, ResourceURI: "/api/admin/configuration/v1/ntp_server/1/"},
			},
			ClientSTUNServers: []config.STUNServer{
				{ID: 2},
			},
			ResourceURI: "/api/admin/configuration/v1/system_location/123/",
		}
	}).Maybe()

	// Mock the DeleteSystemLocation API call
	client.On("DeleteJSON", mock.Anything, mock.MatchedBy(func(path string) bool {
		return path == "configuration/v1/system_location/123/"
	}), mock.Anything).Return(nil)

	testInfinitySystemLocation(t, client)
}

func testInfinitySystemLocation(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: getTestProtoV5ProviderFactories(client),
		Steps: []resource.TestStep{
			{
				Config: test.LoadTestFolder(t, "resource_infinity_system_location_basic"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_system_location.main-location", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_system_location.main-location", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_system_location.main-location", "name", "main"),
					resource.TestCheckResourceAttr("pexip_infinity_system_location.main-location", "description", "Main location for Pexip Infinity System"),
					resource.TestCheckResourceAttr("pexip_infinity_system_location.main-location", "mtu", "1460"),
					resource.TestCheckResourceAttr("pexip_infinity_system_location.main-location", "dns_servers.#", "2"),
					resource.TestCheckTypeSetElemAttr("pexip_infinity_system_location.main-location", "dns_servers.*", "/api/admin/configuration/v1/dns_server/1/"),
					resource.TestCheckTypeSetElemAttr("pexip_infinity_system_location.main-location", "dns_servers.*", "/api/admin/configuration/v1/dns_server/2/"),
					resource.TestCheckResourceAttr("pexip_infinity_system_location.main-location", "client_stun_servers.#", "1"),
					resource.TestCheckTypeSetElemAttr("pexip_infinity_system_location.main-location", "client_stun_servers.*", "/api/admin/configuration/v1/stun_server/1/")),
			},
			{
				Config: test.LoadTestFolder(t, "resource_infinity_system_location_basic_updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_system_location.main-location", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_system_location.main-location", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_system_location.main-location", "name", "main"),
					resource.TestCheckResourceAttr("pexip_infinity_system_location.main-location", "description", "Main location for Pexip Infinity System - updated"),
					resource.TestCheckResourceAttr("pexip_infinity_system_location.main-location", "mtu", "1460"),
					resource.TestCheckResourceAttr("pexip_infinity_system_location.main-location", "dns_servers.#", "1"),
					resource.TestCheckTypeSetElemAttr("pexip_infinity_system_location.main-location", "dns_servers.*", "/api/admin/configuration/v1/dns_server/1/"),
					resource.TestCheckResourceAttr("pexip_infinity_system_location.main-location", "client_stun_servers.#", "1"),
					resource.TestCheckTypeSetElemAttr("pexip_infinity_system_location.main-location", "client_stun_servers.*", "/api/admin/configuration/v1/stun_server/2/"),
				),
			},
		},
	})
}
