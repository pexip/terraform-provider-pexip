/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

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

	// config for update that clears all fields 
	

	// Mock the GetSystemLocation API call for Read operations
	client.On("GetJSON", mock.Anything, "configuration/v1/system_location/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		sysLoc := args.Get(3).(*config.SystemLocation)
		if updated {
			//stunServer := "/api/admin/configuration/v1/stun_server/2/"

			*sysLoc = config.SystemLocation{
				ID:          123,
				Name:        "main",
				ResourceURI:               "/api/admin/configuration/v1/system_location/123/",
				//Description: "Main location for Pexip Infinity System - updated",
				MTU:         1500,
				MediaQoS:      test.IntPtr(0),
				SignallingQoS: test.IntPtr(0),
				//DNSServers: []config.DNSServer{
				//	{ID: 1, ResourceURI: "/api/admin/configuration/v1/dns_server/1/"},
				//},
				//NTPServers: []config.NTPServer{
				//	{ID: 1, ResourceURI: "/api/admin/configuration/v1/ntp_server/1/"},
				//},
				//SyslogServers:             []config.SyslogServer{},
				BDPMPinChecksEnabled:      "GLOBAL",
				BDPMScanQuarantineEnabled: "GLOBAL",
				//UseRelayCandidatesOnly:    false,
				//LocalMSSIPDomain:          "",
				//STUNServer:                &stunServer,
				//ClientTURNServers: []string{
				//	"/api/admin/configuration/v1/turn_server/2/",
				//},
				//ClientSTUNServers: []string{
				//	"/api/admin/configuration/v1/stun_server/2/",
				//},
				//EventSinks: []config.EventSink{
				//	{ID: 1, ResourceURI: "/api/admin/configuration/v1/event_sink/1/"},
				//},
			}
		} else {
			mediaQoS := 46
			signallingQoS := 24
			h323Gatekeeper := "/api/admin/configuration/v1/h323_gatekeeper/1/"
			httpProxy := "/api/admin/configuration/v1/http_proxy/1/"
			liveCaptionsDialOut1 := "/api/admin/configuration/v1/system_location/1/"
			liveCaptionsDialOut2 := "/api/admin/configuration/v1/system_location/2/"
			liveCaptionsDialOut3 := "/api/admin/configuration/v1/system_location/3/"
			localMSSIPDomain := "test-mssip-domain.local"
			mssipProxy := "/api/admin/configuration/v1/mssip_proxy/1/"
			overflowLocation1 := "/api/admin/configuration/v1/system_location/1/"
			overflowLocation2 := "/api/admin/configuration/v1/system_location/2/"
			policyServer := "/api/admin/configuration/v1/policy_server/1/"
			sipProxy := "/api/admin/configuration/v1/sip_proxy/1/"
			snmpNMS := "/api/admin/configuration/v1/snmp_network_management_system/2/"
			stunServer := "/api/admin/configuration/v1/stun_server/1/"
			teamsProxy := "/api/admin/configuration/v1/teams_proxy/1/"
			transcodingLocation := "/api/admin/configuration/v1/system_location/3/"
			turnServer := "/api/admin/configuration/v1/turn_server/3/"
			useRelayCandidatesOnly := true

			*sysLoc = config.SystemLocation{
				ID:               123,
				Name:             "main",
				Description:      "Main location for Pexip Infinity System",
				MTU:              1460,
				MediaQoS:         &mediaQoS,
				SignallingQoS:    &signallingQoS,
				LocalMSSIPDomain: localMSSIPDomain,
				DNSServers: []config.DNSServer{
					{ID: 1, ResourceURI: "/api/admin/configuration/v1/dns_server/1/"},
					{ID: 2, ResourceURI: "/api/admin/configuration/v1/dns_server/2/"},
				},
				NTPServers: []config.NTPServer{
					{ID: 1, ResourceURI: "/api/admin/configuration/v1/ntp_server/1/"},
				},
				SyslogServers:               []config.SyslogServer{},
				H323Gatekeeper:              &h323Gatekeeper,
				SIPProxy:                    &sipProxy,
				MSSIPProxy:                  &mssipProxy,
				TeamsProxy:                  &teamsProxy,
				OverflowLocation1:           &overflowLocation1,
				OverflowLocation2:           &overflowLocation2,
				TranscodingLocation:         &transcodingLocation,
				BDPMPinChecksEnabled:        "ON",
				BDPMScanQuarantineEnabled:   "ON",
				UseRelayCandidatesOnly:      useRelayCandidatesOnly,
				ResourceURI:                 "/api/admin/configuration/v1/system_location/123/",
				SNMPNetworkManagementSystem: &snmpNMS,
				HTTPProxy:                   &httpProxy,
				TURNServer:                  &turnServer,
				STUNServer:                  &stunServer,
				ClientTURNServers: []string{
					"/api/admin/configuration/v1/turn_server/1/",
					"/api/admin/configuration/v1/turn_server/2/",
				},
				ClientSTUNServers: []string{
					"/api/admin/configuration/v1/stun_server/1/",
					"/api/admin/configuration/v1/stun_server/2/",
				},
				EventSinks: []config.EventSink{
					{ID: 1, ResourceURI: "/api/admin/configuration/v1/event_sink/1/"},
				},
				PolicyServer:         &policyServer,
				LiveCaptionsDialOut1: &liveCaptionsDialOut1,
				LiveCaptionsDialOut2: &liveCaptionsDialOut2,
				LiveCaptionsDialOut3: &liveCaptionsDialOut3,
			}
		}
	}).Maybe() // Called multiple times for reads

	// Mock the UpdateSystemLocation API call
	client.On("PutJSON", mock.Anything, mock.MatchedBy(func(path string) bool {
		return path == "configuration/v1/system_location/123/"
	}), mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		updated = true // Mark as updated for subsequent reads
		sysLoc := args.Get(3).(*config.SystemLocation)
		//mediaQoS := 46
		//signallingQoS := 24
		//stunServer := "/api/admin/configuration/v1/stun_server/2/"
		*sysLoc = config.SystemLocation{
			ID:            123,
			Name:          "main",
			ResourceURI: "/api/admin/configuration/v1/system_location/123/",
			//Description:   "Main location for Pexip Infinity System - updated",
			BDPMPinChecksEnabled:      "GLOBAL",
			BDPMScanQuarantineEnabled: "GLOBAL",
			MTU:           1500,
			MediaQoS:      test.IntPtr(0),
			SignallingQoS: test.IntPtr(0),
			//DNSServers: []config.DNSServer{
			//	{ID: 1, ResourceURI: "/api/admin/configuration/v1/dns_server/1/"},
			//},
			//NTPServers: []config.NTPServer{
			//	{ID: 1, ResourceURI: "/api/admin/configuration/v1/ntp_server/1/"},
			//},
			//ClientSTUNServers: []string{
			//	"/api/admin/configuration/v1/stun_server/2/",
			//},
			//ClientTURNServers: []string{
			//	"/api/admin/configuration/v1/turn_server/2/",
			//},
			//EventSinks: []config.EventSink{
			//	{ID: 3, ResourceURI: "/api/admin/configuration/v1/event_sink/3/"},
			//},
			//STUNServer:  &stunServer,
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
					resource.TestCheckResourceAttr("pexip_infinity_system_location.main-location", "ntp_servers.#", "1"),
					resource.TestCheckTypeSetElemAttr("pexip_infinity_system_location.main-location", "ntp_servers.*", "/api/admin/configuration/v1/ntp_server/1/"),
					resource.TestCheckResourceAttr("pexip_infinity_system_location.main-location", "client_stun_servers.#", "2"),
					resource.TestCheckTypeSetElemAttr("pexip_infinity_system_location.main-location", "client_stun_servers.*", "/api/admin/configuration/v1/stun_server/1/"),
					resource.TestCheckTypeSetElemAttr("pexip_infinity_system_location.main-location", "client_stun_servers.*", "/api/admin/configuration/v1/stun_server/2/"),
					resource.TestCheckResourceAttr("pexip_infinity_system_location.main-location", "client_turn_servers.#", "2"),
					resource.TestCheckTypeSetElemAttr("pexip_infinity_system_location.main-location", "client_turn_servers.*", "/api/admin/configuration/v1/turn_server/1/"),
					resource.TestCheckTypeSetElemAttr("pexip_infinity_system_location.main-location", "client_turn_servers.*", "/api/admin/configuration/v1/turn_server/2/"),
					resource.TestCheckResourceAttr("pexip_infinity_system_location.main-location", "event_sinks.#", "1"),
					resource.TestCheckTypeSetElemAttr("pexip_infinity_system_location.main-location", "event_sinks.*", "/api/admin/configuration/v1/event_sink/1/"),
					resource.TestCheckResourceAttr("pexip_infinity_system_location.main-location", "bdpm_pin_checks_enabled", "ON"),
					resource.TestCheckResourceAttr("pexip_infinity_system_location.main-location", "bdpm_scan_quarantine_enabled", "ON"),
					resource.TestCheckResourceAttr("pexip_infinity_system_location.main-location", "h323_gatekeeper", "/api/admin/configuration/v1/h323_gatekeeper/1/"),
					resource.TestCheckResourceAttr("pexip_infinity_system_location.main-location", "http_proxy", "/api/admin/configuration/v1/http_proxy/1/"),
					resource.TestCheckResourceAttr("pexip_infinity_system_location.main-location", "live_captions_dial_out_1", "/api/admin/configuration/v1/system_location/1/"),
					resource.TestCheckResourceAttr("pexip_infinity_system_location.main-location", "live_captions_dial_out_2", "/api/admin/configuration/v1/system_location/2/"),
					resource.TestCheckResourceAttr("pexip_infinity_system_location.main-location", "live_captions_dial_out_3", "/api/admin/configuration/v1/system_location/3/"),
					resource.TestCheckResourceAttr("pexip_infinity_system_location.main-location", "local_mssip_domain", "test-mssip-domain.local"),
					resource.TestCheckResourceAttr("pexip_infinity_system_location.main-location", "media_qos", "46"),
					resource.TestCheckResourceAttr("pexip_infinity_system_location.main-location", "mssip_proxy", "/api/admin/configuration/v1/mssip_proxy/1/"),
					resource.TestCheckResourceAttr("pexip_infinity_system_location.main-location", "overflow_location1", "/api/admin/configuration/v1/system_location/1/"),
					resource.TestCheckResourceAttr("pexip_infinity_system_location.main-location", "overflow_location2", "/api/admin/configuration/v1/system_location/2/"),
					resource.TestCheckResourceAttr("pexip_infinity_system_location.main-location", "policy_server", "/api/admin/configuration/v1/policy_server/1/"),
					resource.TestCheckResourceAttr("pexip_infinity_system_location.main-location", "signalling_qos", "24"),
					resource.TestCheckResourceAttr("pexip_infinity_system_location.main-location", "sip_proxy", "/api/admin/configuration/v1/sip_proxy/1/"),
					resource.TestCheckResourceAttr("pexip_infinity_system_location.main-location", "snmp_network_management_system", "/api/admin/configuration/v1/snmp_network_management_system/2/"),
					resource.TestCheckResourceAttr("pexip_infinity_system_location.main-location", "stun_server", "/api/admin/configuration/v1/stun_server/1/"),
					resource.TestCheckResourceAttr("pexip_infinity_system_location.main-location", "syslog_servers.#", "0"),
					resource.TestCheckResourceAttr("pexip_infinity_system_location.main-location", "teams_proxy", "/api/admin/configuration/v1/teams_proxy/1/"),
					resource.TestCheckResourceAttr("pexip_infinity_system_location.main-location", "transcoding_location", "/api/admin/configuration/v1/system_location/3/"),
					resource.TestCheckResourceAttr("pexip_infinity_system_location.main-location", "turn_server", "/api/admin/configuration/v1/turn_server/3/"),
					resource.TestCheckResourceAttr("pexip_infinity_system_location.main-location", "use_relay_candidates_only", "true"),
				),
			},
			{
				Config: test.LoadTestFolder(t, "resource_infinity_system_location_basic_updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_system_location.main-location", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_system_location.main-location", "resource_id"),
					//resource.TestCheckResourceAttr("pexip_infinity_system_location.main-location", "name", "main"),
					//resource.TestCheckResourceAttr("pexip_infinity_system_location.main-location", "description", "Main location for Pexip Infinity System - updated"),
					resource.TestCheckResourceAttr("pexip_infinity_system_location.main-location", "mtu", "1500"),
					//resource.TestCheckResourceAttr("pexip_infinity_system_location.main-location", "dns_servers.#", "1"),
					//resource.TestCheckTypeSetElemAttr("pexip_infinity_system_location.main-location", "dns_servers.*", "/api/admin/configuration/v1/dns_server/1/"),
					//resource.TestCheckResourceAttr("pexip_infinity_system_location.main-location", "client_stun_servers.#", "1"),
					//resource.TestCheckTypeSetElemAttr("pexip_infinity_system_location.main-location", "client_stun_servers.*", "/api/admin/configuration/v1/stun_server/2/"),
					//resource.TestCheckResourceAttr("pexip_infinity_system_location.main-location", "client_turn_servers.#", "1"),
					//resource.TestCheckTypeSetElemAttr("pexip_infinity_system_location.main-location", "client_turn_servers.*", "/api/admin/configuration/v1/turn_server/2/"),
					//resource.TestCheckResourceAttr("pexip_infinity_system_location.main-location", "event_sinks.#", "1"),
					//resource.TestCheckTypeSetElemAttr("pexip_infinity_system_location.main-location", "event_sinks.*", "/api/admin/configuration/v1/event_sink/1/"),
					//resource.TestCheckResourceAttr("pexip_infinity_system_location.main-location", "stun_server", "/api/admin/configuration/v1/stun_server/2/"),
					resource.TestCheckResourceAttr("pexip_infinity_system_location.main-location", "bdpm_pin_checks_enabled", "GLOBAL"),
					resource.TestCheckResourceAttr("pexip_infinity_system_location.main-location", "bdpm_scan_quarantine_enabled", "GLOBAL"),
					//resource.TestCheckResourceAttr("pexip_infinity_system_location.main-location", "use_relay_candidates_only", "false"),
					//resource.TestCheckResourceAttr("pexip_infinity_system_location.main-location", "local_mssip_domain", ""),
					resource.TestCheckResourceAttr("pexip_infinity_system_location.main-location", "media_qos", "0"),
					resource.TestCheckResourceAttr("pexip_infinity_system_location.main-location", "signalling_qos", "0"),
				),
			},
		},
	})
}
