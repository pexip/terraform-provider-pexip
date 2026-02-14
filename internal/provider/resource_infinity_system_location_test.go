/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package provider

import (
	"os"
	"testing"

	"github.com/pexip/go-infinity-sdk/v38"
	"github.com/pexip/go-infinity-sdk/v38/config"
	"github.com/pexip/go-infinity-sdk/v38/types"
	"github.com/stretchr/testify/mock"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/pexip/terraform-provider-pexip/internal/test"
)

// Helper function to convert URI slices to resource slices
func convertURIsToResources[T any](uris []string) []T {
	if len(uris) == 0 {
		return []T{}
	}

	resources := make([]T, len(uris))
	for i, uri := range uris {
		var zero T
		// Extract ID from URI (e.g., "/api/admin/configuration/v1/dns_server/1/" -> 1)
		// Simple extraction: find the last number before the trailing slash
		var id int
		if len(uri) > 0 {
			// Remove trailing slash
			trimmed := uri
			if trimmed[len(trimmed)-1] == '/' {
				trimmed = trimmed[:len(trimmed)-1]
			}
			// Find last slash
			lastSlash := -1
			for j := len(trimmed) - 1; j >= 0; j-- {
				if trimmed[j] == '/' {
					lastSlash = j
					break
				}
			}
			if lastSlash >= 0 && lastSlash < len(trimmed)-1 {
				idStr := trimmed[lastSlash+1:]
				// Parse the ID
				for _, c := range idStr {
					if c >= '0' && c <= '9' {
						id = id*10 + int(c-'0')
					}
				}
			}
		}

		switch any(zero).(type) {
		case config.DNSServer:
			resources[i] = any(config.DNSServer{ID: id, ResourceURI: uri}).(T)
		case config.NTPServer:
			resources[i] = any(config.NTPServer{ID: id, ResourceURI: uri}).(T)
		case config.SyslogServer:
			resources[i] = any(config.SyslogServer{ID: id, ResourceURI: uri}).(T)
		case config.EventSink:
			resources[i] = any(config.EventSink{ID: id, ResourceURI: uri}).(T)
		}
	}
	return resources
}

func TestInfinitySystemLocation(t *testing.T) {
	t.Parallel()
	_ = os.Setenv("TF_ACC", "1")

	client := infinity.NewClientMock()

	// Mock DNS Servers (created twice: step 1 and step 4)
	mockDNS1 := &config.DNSServer{}
	mockDNS2 := &config.DNSServer{}

	// Step 1: Create DNS1 (68.94.156.1)
	client.On("PostWithResponse", mock.Anything, "configuration/v1/dns_server/", mock.MatchedBy(func(req *config.DNSServerCreateRequest) bool {
		return req.Address == "68.94.156.1"
	}), mock.Anything).Return(&types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/dns_server/1/",
	}, nil).Run(func(args mock.Arguments) {
		req := args.Get(2).(*config.DNSServerCreateRequest)
		*mockDNS1 = config.DNSServer{
			ID:          1,
			ResourceURI: "/api/admin/configuration/v1/dns_server/1/",
			Address:     req.Address,
			Description: req.Description,
		}
	}).Once()

	// Step 1: Create DNS2 (68.94.157.1)
	client.On("PostWithResponse", mock.Anything, "configuration/v1/dns_server/", mock.MatchedBy(func(req *config.DNSServerCreateRequest) bool {
		return req.Address == "68.94.157.1"
	}), mock.Anything).Return(&types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/dns_server/2/",
	}, nil).Run(func(args mock.Arguments) {
		req := args.Get(2).(*config.DNSServerCreateRequest)
		*mockDNS2 = config.DNSServer{
			ID:          2,
			ResourceURI: "/api/admin/configuration/v1/dns_server/2/",
			Address:     req.Address,
			Description: req.Description,
		}
	}).Once()

	// Step 4: Recreate DNS1 (68.94.156.1)
	client.On("PostWithResponse", mock.Anything, "configuration/v1/dns_server/", mock.MatchedBy(func(req *config.DNSServerCreateRequest) bool {
		return req.Address == "68.94.156.1"
	}), mock.Anything).Return(&types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/dns_server/1/",
	}, nil).Run(func(args mock.Arguments) {
		req := args.Get(2).(*config.DNSServerCreateRequest)
		*mockDNS1 = config.DNSServer{
			ID:          1,
			ResourceURI: "/api/admin/configuration/v1/dns_server/1/",
			Address:     req.Address,
			Description: req.Description,
		}
	}).Once()

	// Step 4: Recreate DNS2 (68.94.157.1)
	client.On("PostWithResponse", mock.Anything, "configuration/v1/dns_server/", mock.MatchedBy(func(req *config.DNSServerCreateRequest) bool {
		return req.Address == "68.94.157.1"
	}), mock.Anything).Return(&types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/dns_server/2/",
	}, nil).Run(func(args mock.Arguments) {
		req := args.Get(2).(*config.DNSServerCreateRequest)
		*mockDNS2 = config.DNSServer{
			ID:          2,
			ResourceURI: "/api/admin/configuration/v1/dns_server/2/",
			Address:     req.Address,
			Description: req.Description,
		}
	}).Once()
	client.On("GetJSON", mock.Anything, "configuration/v1/dns_server/1/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		dns := args.Get(3).(*config.DNSServer)
		*dns = *mockDNS1
	}).Maybe()

	client.On("DeleteJSON", mock.Anything, "configuration/v1/dns_server/1/", mock.Anything).Return(nil).Maybe()

	client.On("GetJSON", mock.Anything, "configuration/v1/dns_server/2/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		dns := args.Get(3).(*config.DNSServer)
		*dns = *mockDNS2
	}).Maybe()

	client.On("DeleteJSON", mock.Anything, "configuration/v1/dns_server/2/", mock.Anything).Return(nil).Maybe()

	client.On("GetJSON", mock.Anything, "configuration/v1/dns_server/2/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		dns := args.Get(3).(*config.DNSServer)
		*dns = *mockDNS2
	}).Maybe()

	client.On("DeleteJSON", mock.Anything, "configuration/v1/dns_server/2/", mock.Anything).Return(nil).Maybe()

	// Mock NTP Server
	mockNTP := &config.NTPServer{}

	// Step 1: Create NTP (time.google.com)
	client.On("PostWithResponse", mock.Anything, "configuration/v1/ntp_server/", mock.MatchedBy(func(req *config.NTPServerCreateRequest) bool {
		return req.Address == "time.google.com"
	}), mock.Anything).Return(&types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/ntp_server/3/",
	}, nil).Run(func(args mock.Arguments) {
		req := args.Get(2).(*config.NTPServerCreateRequest)
		*mockNTP = config.NTPServer{
			ID:          3,
			ResourceURI: "/api/admin/configuration/v1/ntp_server/3/",
			Address:     req.Address,
			Description: req.Description,
			Key:         req.Key,
			KeyID:       req.KeyID,
		}
	}).Once()

	// Step 4: Recreate NTP (time.google.com)
	client.On("PostWithResponse", mock.Anything, "configuration/v1/ntp_server/", mock.MatchedBy(func(req *config.NTPServerCreateRequest) bool {
		return req.Address == "time.google.com"
	}), mock.Anything).Return(&types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/ntp_server/3/",
	}, nil).Run(func(args mock.Arguments) {
		req := args.Get(2).(*config.NTPServerCreateRequest)
		*mockNTP = config.NTPServer{
			ID:          3,
			ResourceURI: "/api/admin/configuration/v1/ntp_server/3/",
			Address:     req.Address,
			Description: req.Description,
			Key:         req.Key,
			KeyID:       req.KeyID,
		}
	}).Once()

	client.On("GetJSON", mock.Anything, "configuration/v1/ntp_server/3/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		ntp := args.Get(3).(*config.NTPServer)
		*ntp = *mockNTP
	}).Maybe()

	client.On("DeleteJSON", mock.Anything, "configuration/v1/ntp_server/3/", mock.Anything).Return(nil).Maybe()

	// Mock STUN Server 1
	mockSTUN1 := &config.STUNServer{}

	// Step 1: Create STUN1 (tf-test-stun1)
	client.On("PostWithResponse", mock.Anything, "configuration/v1/stun_server/", mock.MatchedBy(func(req *config.STUNServerCreateRequest) bool {
		return req.Name == "tf-test-stun1"
	}), mock.Anything).Return(&types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/stun_server/4/",
	}, nil).Run(func(args mock.Arguments) {
		req := args.Get(2).(*config.STUNServerCreateRequest)
		*mockSTUN1 = config.STUNServer{
			ID:          4,
			ResourceURI: "/api/admin/configuration/v1/stun_server/4/",
			Name:        req.Name,
			Address:     req.Address,
			Port:        req.Port,
			Description: req.Description,
		}
	}).Once()

	// Step 4: Recreate STUN1 (tf-test-stun1)
	client.On("PostWithResponse", mock.Anything, "configuration/v1/stun_server/", mock.MatchedBy(func(req *config.STUNServerCreateRequest) bool {
		return req.Name == "tf-test-stun1"
	}), mock.Anything).Return(&types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/stun_server/4/",
	}, nil).Run(func(args mock.Arguments) {
		req := args.Get(2).(*config.STUNServerCreateRequest)
		*mockSTUN1 = config.STUNServer{
			ID:          4,
			ResourceURI: "/api/admin/configuration/v1/stun_server/4/",
			Name:        req.Name,
			Address:     req.Address,
			Port:        req.Port,
			Description: req.Description,
		}
	}).Once()
	client.On("GetJSON", mock.Anything, "configuration/v1/stun_server/4/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		stun := args.Get(3).(*config.STUNServer)
		*stun = *mockSTUN1
	}).Maybe()

	client.On("DeleteJSON", mock.Anything, "configuration/v1/stun_server/4/", mock.Anything).Return(nil).Maybe()

	// Mock STUN Server 2
	mockSTUN2 := &config.STUNServer{}

	// Step 1: Create STUN2 (tf-test-stun2)
	client.On("PostWithResponse", mock.Anything, "configuration/v1/stun_server/", mock.MatchedBy(func(req *config.STUNServerCreateRequest) bool {
		return req.Name == "tf-test-stun2"
	}), mock.Anything).Return(&types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/stun_server/5/",
	}, nil).Run(func(args mock.Arguments) {
		req := args.Get(2).(*config.STUNServerCreateRequest)
		*mockSTUN2 = config.STUNServer{
			ID:          5,
			ResourceURI: "/api/admin/configuration/v1/stun_server/5/",
			Name:        req.Name,
			Address:     req.Address,
			Port:        req.Port,
			Description: req.Description,
		}
	}).Once()

	// Step 4: Recreate STUN2 (tf-test-stun2)
	client.On("PostWithResponse", mock.Anything, "configuration/v1/stun_server/", mock.MatchedBy(func(req *config.STUNServerCreateRequest) bool {
		return req.Name == "tf-test-stun2"
	}), mock.Anything).Return(&types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/stun_server/5/",
	}, nil).Run(func(args mock.Arguments) {
		req := args.Get(2).(*config.STUNServerCreateRequest)
		*mockSTUN2 = config.STUNServer{
			ID:          5,
			ResourceURI: "/api/admin/configuration/v1/stun_server/5/",
			Name:        req.Name,
			Address:     req.Address,
			Port:        req.Port,
			Description: req.Description,
		}
	}).Once()
	client.On("GetJSON", mock.Anything, "configuration/v1/stun_server/5/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		stun := args.Get(3).(*config.STUNServer)
		*stun = *mockSTUN2
	}).Maybe()

	client.On("DeleteJSON", mock.Anything, "configuration/v1/stun_server/5/", mock.Anything).Return(nil).Maybe()

	// Mock TURN Server 1
	mockTURN1 := &config.TURNServer{}

	// Step 1: Create TURN1 (tf-test-turn1)
	client.On("PostWithResponse", mock.Anything, "configuration/v1/turn_server/", mock.MatchedBy(func(req *config.TURNServerCreateRequest) bool {
		return req.Name == "tf-test-turn1"
	}), mock.Anything).Return(&types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/turn_server/6/",
	}, nil).Run(func(args mock.Arguments) {
		req := args.Get(2).(*config.TURNServerCreateRequest)
		*mockTURN1 = config.TURNServer{
			ID:            6,
			ResourceURI:   "/api/admin/configuration/v1/turn_server/6/",
			Name:          req.Name,
			Address:       req.Address,
			Port:          req.Port,
			ServerType:    req.ServerType,
			TransportType: req.TransportType,
			Description:   req.Description,
			Username:      req.Username,
			Password:      req.Password,
		}
	}).Once()

	// Step 4: Recreate TURN1 (tf-test-turn1)
	client.On("PostWithResponse", mock.Anything, "configuration/v1/turn_server/", mock.MatchedBy(func(req *config.TURNServerCreateRequest) bool {
		return req.Name == "tf-test-turn1"
	}), mock.Anything).Return(&types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/turn_server/6/",
	}, nil).Run(func(args mock.Arguments) {
		req := args.Get(2).(*config.TURNServerCreateRequest)
		*mockTURN1 = config.TURNServer{
			ID:            6,
			ResourceURI:   "/api/admin/configuration/v1/turn_server/6/",
			Name:          req.Name,
			Address:       req.Address,
			Port:          req.Port,
			ServerType:    req.ServerType,
			TransportType: req.TransportType,
			Description:   req.Description,
			Username:      req.Username,
			Password:      req.Password,
		}
	}).Once()
	client.On("GetJSON", mock.Anything, "configuration/v1/turn_server/6/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		turn := args.Get(3).(*config.TURNServer)
		*turn = *mockTURN1
	}).Maybe()

	client.On("DeleteJSON", mock.Anything, "configuration/v1/turn_server/6/", mock.Anything).Return(nil).Maybe()

	// Mock TURN Server 2
	mockTURN2 := &config.TURNServer{}

	// Step 1: Create TURN2 (tf-test-turn2)
	client.On("PostWithResponse", mock.Anything, "configuration/v1/turn_server/", mock.MatchedBy(func(req *config.TURNServerCreateRequest) bool {
		return req.Name == "tf-test-turn2"
	}), mock.Anything).Return(&types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/turn_server/7/",
	}, nil).Run(func(args mock.Arguments) {
		req := args.Get(2).(*config.TURNServerCreateRequest)
		*mockTURN2 = config.TURNServer{
			ID:            7,
			ResourceURI:   "/api/admin/configuration/v1/turn_server/7/",
			Name:          req.Name,
			Address:       req.Address,
			Port:          req.Port,
			ServerType:    req.ServerType,
			TransportType: req.TransportType,
			Description:   req.Description,
			Username:      req.Username,
			Password:      req.Password,
		}
	}).Once()

	// Step 4: Recreate TURN2 (tf-test-turn2)
	client.On("PostWithResponse", mock.Anything, "configuration/v1/turn_server/", mock.MatchedBy(func(req *config.TURNServerCreateRequest) bool {
		return req.Name == "tf-test-turn2"
	}), mock.Anything).Return(&types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/turn_server/7/",
	}, nil).Run(func(args mock.Arguments) {
		req := args.Get(2).(*config.TURNServerCreateRequest)
		*mockTURN2 = config.TURNServer{
			ID:            7,
			ResourceURI:   "/api/admin/configuration/v1/turn_server/7/",
			Name:          req.Name,
			Address:       req.Address,
			Port:          req.Port,
			ServerType:    req.ServerType,
			TransportType: req.TransportType,
			Description:   req.Description,
			Username:      req.Username,
			Password:      req.Password,
		}
	}).Once()
	client.On("GetJSON", mock.Anything, "configuration/v1/turn_server/7/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		turn := args.Get(3).(*config.TURNServer)
		*turn = *mockTURN2
	}).Maybe()

	client.On("DeleteJSON", mock.Anything, "configuration/v1/turn_server/7/", mock.Anything).Return(nil).Maybe()

	// Mock Event Sink
	mockEventSink := &config.EventSink{}

	// Step 1: Create Event Sink (tf-test-event-sink)
	client.On("PostWithResponse", mock.Anything, "configuration/v1/event_sink/", mock.MatchedBy(func(req *config.EventSinkCreateRequest) bool {
		return req.Name == "tf-test-event-sink"
	}), mock.Anything).Return(&types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/event_sink/8/",
	}, nil).Run(func(args mock.Arguments) {
		req := args.Get(2).(*config.EventSinkCreateRequest)
		*mockEventSink = config.EventSink{
			ID:                   8,
			ResourceURI:          "/api/admin/configuration/v1/event_sink/8/",
			Name:                 req.Name,
			URL:                  req.URL,
			Description:          req.Description,
			BulkSupport:          req.BulkSupport,
			VerifyTLSCertificate: req.VerifyTLSCertificate,
			Version:              req.Version,
		}
	}).Once()

	// Step 4: Recreate Event Sink (tf-test-event-sink)
	client.On("PostWithResponse", mock.Anything, "configuration/v1/event_sink/", mock.MatchedBy(func(req *config.EventSinkCreateRequest) bool {
		return req.Name == "tf-test-event-sink"
	}), mock.Anything).Return(&types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/event_sink/8/",
	}, nil).Run(func(args mock.Arguments) {
		req := args.Get(2).(*config.EventSinkCreateRequest)
		*mockEventSink = config.EventSink{
			ID:                   8,
			ResourceURI:          "/api/admin/configuration/v1/event_sink/8/",
			Name:                 req.Name,
			URL:                  req.URL,
			Description:          req.Description,
			BulkSupport:          req.BulkSupport,
			VerifyTLSCertificate: req.VerifyTLSCertificate,
			Version:              req.Version,
		}
	}).Once()

	client.On("GetJSON", mock.Anything, "configuration/v1/event_sink/8/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		event := args.Get(3).(*config.EventSink)
		*event = *mockEventSink
	}).Maybe()

	client.On("DeleteJSON", mock.Anything, "configuration/v1/event_sink/8/", mock.Anything).Return(nil).Maybe()

	// Mock H323 Gatekeeper
	mockH323 := &config.H323Gatekeeper{}

	// Step 1: Create H323 Gatekeeper (tf-test-h323-gk)
	client.On("PostWithResponse", mock.Anything, "configuration/v1/h323_gatekeeper/", mock.MatchedBy(func(req *config.H323GatekeeperCreateRequest) bool {
		return req.Name == "tf-test-h323-gk"
	}), mock.Anything).Return(&types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/h323_gatekeeper/9/",
	}, nil).Run(func(args mock.Arguments) {
		req := args.Get(2).(*config.H323GatekeeperCreateRequest)
		*mockH323 = config.H323Gatekeeper{
			ID:          9,
			ResourceURI: "/api/admin/configuration/v1/h323_gatekeeper/9/",
			Name:        req.Name,
			Address:     req.Address,
			Port:        req.Port,
			Description: req.Description,
		}
	}).Once()

	// Step 4: Recreate H323 Gatekeeper (tf-test-h323-gk)
	client.On("PostWithResponse", mock.Anything, "configuration/v1/h323_gatekeeper/", mock.MatchedBy(func(req *config.H323GatekeeperCreateRequest) bool {
		return req.Name == "tf-test-h323-gk"
	}), mock.Anything).Return(&types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/h323_gatekeeper/9/",
	}, nil).Run(func(args mock.Arguments) {
		req := args.Get(2).(*config.H323GatekeeperCreateRequest)
		*mockH323 = config.H323Gatekeeper{
			ID:          9,
			ResourceURI: "/api/admin/configuration/v1/h323_gatekeeper/9/",
			Name:        req.Name,
			Address:     req.Address,
			Port:        req.Port,
			Description: req.Description,
		}
	}).Once()

	client.On("GetJSON", mock.Anything, "configuration/v1/h323_gatekeeper/9/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		h323 := args.Get(3).(*config.H323Gatekeeper)
		*h323 = *mockH323
	}).Maybe()

	client.On("DeleteJSON", mock.Anything, "configuration/v1/h323_gatekeeper/9/", mock.Anything).Return(nil).Maybe()

	// Mock HTTP Proxy
	mockHTTP := &config.HTTPProxy{}

	// Step 1: Create HTTP Proxy (tf-test-http-proxy)
	client.On("PostWithResponse", mock.Anything, "configuration/v1/http_proxy/", mock.MatchedBy(func(req *config.HTTPProxyCreateRequest) bool {
		return req.Name == "tf-test-http-proxy"
	}), mock.Anything).Return(&types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/http_proxy/10/",
	}, nil).Run(func(args mock.Arguments) {
		req := args.Get(2).(*config.HTTPProxyCreateRequest)
		*mockHTTP = config.HTTPProxy{
			ID:          10,
			ResourceURI: "/api/admin/configuration/v1/http_proxy/10/",
			Name:        req.Name,
			Address:     req.Address,
			Port:        req.Port,
			Protocol:    req.Protocol,
			Username:    req.Username,
			Password:    req.Password,
		}
	}).Once()

	// Step 4: Recreate HTTP Proxy (tf-test-http-proxy)
	client.On("PostWithResponse", mock.Anything, "configuration/v1/http_proxy/", mock.MatchedBy(func(req *config.HTTPProxyCreateRequest) bool {
		return req.Name == "tf-test-http-proxy"
	}), mock.Anything).Return(&types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/http_proxy/10/",
	}, nil).Run(func(args mock.Arguments) {
		req := args.Get(2).(*config.HTTPProxyCreateRequest)
		*mockHTTP = config.HTTPProxy{
			ID:          10,
			ResourceURI: "/api/admin/configuration/v1/http_proxy/10/",
			Name:        req.Name,
			Address:     req.Address,
			Port:        req.Port,
			Protocol:    req.Protocol,
			Username:    req.Username,
			Password:    req.Password,
		}
	}).Once()

	client.On("GetJSON", mock.Anything, "configuration/v1/http_proxy/10/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		http := args.Get(3).(*config.HTTPProxy)
		*http = *mockHTTP
	}).Maybe()

	client.On("DeleteJSON", mock.Anything, "configuration/v1/http_proxy/10/", mock.Anything).Return(nil).Maybe()

	// Mock MSSIP Proxy
	mockMSSIP := &config.MSSIPProxy{}

	// Step 1: Create MSSIP Proxy (tf-test-mssip)
	client.On("PostWithResponse", mock.Anything, "configuration/v1/mssip_proxy/", mock.MatchedBy(func(req *config.MSSIPProxyCreateRequest) bool {
		return req.Name == "tf-test-mssip"
	}), mock.Anything).Return(&types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/mssip_proxy/11/",
	}, nil).Run(func(args mock.Arguments) {
		req := args.Get(2).(*config.MSSIPProxyCreateRequest)
		*mockMSSIP = config.MSSIPProxy{
			ID:          11,
			ResourceURI: "/api/admin/configuration/v1/mssip_proxy/11/",
			Name:        req.Name,
			Address:     req.Address,
			Port:        req.Port,
			Transport:   req.Transport,
			Description: req.Description,
		}
	}).Once()

	// Step 4: Recreate MSSIP Proxy (tf-test-mssip)
	client.On("PostWithResponse", mock.Anything, "configuration/v1/mssip_proxy/", mock.MatchedBy(func(req *config.MSSIPProxyCreateRequest) bool {
		return req.Name == "tf-test-mssip"
	}), mock.Anything).Return(&types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/mssip_proxy/11/",
	}, nil).Run(func(args mock.Arguments) {
		req := args.Get(2).(*config.MSSIPProxyCreateRequest)
		*mockMSSIP = config.MSSIPProxy{
			ID:          11,
			ResourceURI: "/api/admin/configuration/v1/mssip_proxy/11/",
			Name:        req.Name,
			Address:     req.Address,
			Port:        req.Port,
			Transport:   req.Transport,
			Description: req.Description,
		}
	}).Once()

	client.On("GetJSON", mock.Anything, "configuration/v1/mssip_proxy/11/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		mssip := args.Get(3).(*config.MSSIPProxy)
		*mssip = *mockMSSIP
	}).Maybe()

	client.On("DeleteJSON", mock.Anything, "configuration/v1/mssip_proxy/11/", mock.Anything).Return(nil).Maybe()

	// Mock Policy Server
	mockPolicy := &config.PolicyServer{}

	// Step 1: Create Policy Server (tf-test-policy)
	client.On("PostWithResponse", mock.Anything, "configuration/v1/policy_server/", mock.MatchedBy(func(req *config.PolicyServerCreateRequest) bool {
		return req.Name == "tf-test-policy"
	}), mock.Anything).Return(&types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/policy_server/12/",
	}, nil).Run(func(args mock.Arguments) {
		req := args.Get(2).(*config.PolicyServerCreateRequest)
		*mockPolicy = config.PolicyServer{
			ID:                                 12,
			ResourceURI:                        "/api/admin/configuration/v1/policy_server/12/",
			Name:                               req.Name,
			URL:                                req.URL,
			EnableAvatarLookup:                 req.EnableAvatarLookup,
			EnableDirectoryLookup:              req.EnableDirectoryLookup,
			EnableInternalMediaLocationPolicy:  req.EnableInternalMediaLocationPolicy,
			EnableInternalParticipantPolicy:    req.EnableInternalParticipantPolicy,
			EnableInternalServicePolicy:        req.EnableInternalServicePolicy,
			Description:                        req.Description,
		}
	}).Once()

	// Step 4: Recreate Policy Server (tf-test-policy)
	client.On("PostWithResponse", mock.Anything, "configuration/v1/policy_server/", mock.MatchedBy(func(req *config.PolicyServerCreateRequest) bool {
		return req.Name == "tf-test-policy"
	}), mock.Anything).Return(&types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/policy_server/12/",
	}, nil).Run(func(args mock.Arguments) {
		req := args.Get(2).(*config.PolicyServerCreateRequest)
		*mockPolicy = config.PolicyServer{
			ID:                                 12,
			ResourceURI:                        "/api/admin/configuration/v1/policy_server/12/",
			Name:                               req.Name,
			URL:                                req.URL,
			EnableAvatarLookup:                 req.EnableAvatarLookup,
			EnableDirectoryLookup:              req.EnableDirectoryLookup,
			EnableInternalMediaLocationPolicy:  req.EnableInternalMediaLocationPolicy,
			EnableInternalParticipantPolicy:    req.EnableInternalParticipantPolicy,
			EnableInternalServicePolicy:        req.EnableInternalServicePolicy,
			Description:                        req.Description,
		}
	}).Once()

	client.On("GetJSON", mock.Anything, "configuration/v1/policy_server/12/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		policy := args.Get(3).(*config.PolicyServer)
		*policy = *mockPolicy
	}).Maybe()

	client.On("DeleteJSON", mock.Anything, "configuration/v1/policy_server/12/", mock.Anything).Return(nil).Maybe()

	// Mock SIP Proxy
	mockSIP := &config.SIPProxy{}

	// Step 1: Create SIP Proxy (tf-test-sip)
	client.On("PostWithResponse", mock.Anything, "configuration/v1/sip_proxy/", mock.MatchedBy(func(req *config.SIPProxyCreateRequest) bool {
		return req.Name == "tf-test-sip"
	}), mock.Anything).Return(&types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/sip_proxy/13/",
	}, nil).Run(func(args mock.Arguments) {
		req := args.Get(2).(*config.SIPProxyCreateRequest)
		*mockSIP = config.SIPProxy{
			ID:          13,
			ResourceURI: "/api/admin/configuration/v1/sip_proxy/13/",
			Name:        req.Name,
			Address:     req.Address,
			Port:        req.Port,
			Transport:   req.Transport,
			Description: req.Description,
		}
	}).Once()

	// Step 4: Recreate SIP Proxy (tf-test-sip)
	client.On("PostWithResponse", mock.Anything, "configuration/v1/sip_proxy/", mock.MatchedBy(func(req *config.SIPProxyCreateRequest) bool {
		return req.Name == "tf-test-sip"
	}), mock.Anything).Return(&types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/sip_proxy/13/",
	}, nil).Run(func(args mock.Arguments) {
		req := args.Get(2).(*config.SIPProxyCreateRequest)
		*mockSIP = config.SIPProxy{
			ID:          13,
			ResourceURI: "/api/admin/configuration/v1/sip_proxy/13/",
			Name:        req.Name,
			Address:     req.Address,
			Port:        req.Port,
			Transport:   req.Transport,
			Description: req.Description,
		}
	}).Once()

	client.On("GetJSON", mock.Anything, "configuration/v1/sip_proxy/13/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		sip := args.Get(3).(*config.SIPProxy)
		*sip = *mockSIP
	}).Maybe()

	client.On("DeleteJSON", mock.Anything, "configuration/v1/sip_proxy/13/", mock.Anything).Return(nil).Maybe()

	// Mock SNMP Network Management System
	mockSNMP := &config.SnmpNetworkManagementSystem{}

	// Step 1: Create SNMP NMS (tf-test-snmp)
	client.On("PostWithResponse", mock.Anything, "configuration/v1/snmp_network_management_system/", mock.MatchedBy(func(req *config.SnmpNetworkManagementSystemCreateRequest) bool {
		return req.Name == "tf-test-snmp"
	}), mock.Anything).Return(&types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/snmp_network_management_system/14/",
	}, nil).Run(func(args mock.Arguments) {
		req := args.Get(2).(*config.SnmpNetworkManagementSystemCreateRequest)
		*mockSNMP = config.SnmpNetworkManagementSystem{
			ID:                14,
			ResourceURI:       "/api/admin/configuration/v1/snmp_network_management_system/14/",
			Name:              req.Name,
			Address:           req.Address,
			Port:              req.Port,
			SnmpTrapCommunity: req.SnmpTrapCommunity,
			Description:       req.Description,
		}
	}).Once()

	// Step 4: Recreate SNMP NMS (tf-test-snmp)
	client.On("PostWithResponse", mock.Anything, "configuration/v1/snmp_network_management_system/", mock.MatchedBy(func(req *config.SnmpNetworkManagementSystemCreateRequest) bool {
		return req.Name == "tf-test-snmp"
	}), mock.Anything).Return(&types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/snmp_network_management_system/14/",
	}, nil).Run(func(args mock.Arguments) {
		req := args.Get(2).(*config.SnmpNetworkManagementSystemCreateRequest)
		*mockSNMP = config.SnmpNetworkManagementSystem{
			ID:                14,
			ResourceURI:       "/api/admin/configuration/v1/snmp_network_management_system/14/",
			Name:              req.Name,
			Address:           req.Address,
			Port:              req.Port,
			SnmpTrapCommunity: req.SnmpTrapCommunity,
			Description:       req.Description,
		}
	}).Once()

	client.On("GetJSON", mock.Anything, "configuration/v1/snmp_network_management_system/14/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		snmp := args.Get(3).(*config.SnmpNetworkManagementSystem)
		*snmp = *mockSNMP
	}).Maybe()

	client.On("DeleteJSON", mock.Anything, "configuration/v1/snmp_network_management_system/14/", mock.Anything).Return(nil).Maybe()

	// Mock Azure Tenant
	mockAzure := &config.AzureTenant{}

	// Step 1: Create Azure Tenant (tf-test-azure-tenant-for-location)
	client.On("PostWithResponse", mock.Anything, "configuration/v1/azure_tenant/", mock.MatchedBy(func(req *config.AzureTenantCreateRequest) bool {
		return req.Name == "tf-test-azure-tenant-for-location"
	}), mock.Anything).Return(&types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/azure_tenant/15/",
	}, nil).Run(func(args mock.Arguments) {
		req := args.Get(2).(*config.AzureTenantCreateRequest)
		*mockAzure = config.AzureTenant{
			ID:          15,
			ResourceURI: "/api/admin/configuration/v1/azure_tenant/15/",
			Name:        req.Name,
			TenantID:    req.TenantID,
			Description: req.Description,
		}
	}).Once()

	// Step 4: Recreate Azure Tenant (tf-test-azure-tenant-for-location)
	client.On("PostWithResponse", mock.Anything, "configuration/v1/azure_tenant/", mock.MatchedBy(func(req *config.AzureTenantCreateRequest) bool {
		return req.Name == "tf-test-azure-tenant-for-location"
	}), mock.Anything).Return(&types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/azure_tenant/15/",
	}, nil).Run(func(args mock.Arguments) {
		req := args.Get(2).(*config.AzureTenantCreateRequest)
		*mockAzure = config.AzureTenant{
			ID:          15,
			ResourceURI: "/api/admin/configuration/v1/azure_tenant/15/",
			Name:        req.Name,
			TenantID:    req.TenantID,
			Description: req.Description,
		}
	}).Once()

	client.On("GetJSON", mock.Anything, "configuration/v1/azure_tenant/15/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		azure := args.Get(3).(*config.AzureTenant)
		*azure = *mockAzure
	}).Maybe()

	client.On("DeleteJSON", mock.Anything, "configuration/v1/azure_tenant/15/", mock.Anything).Return(nil).Maybe()

	// Mock Teams Proxy
	mockTeams := &config.TeamsProxy{}

	// Step 1: Create Teams Proxy (tf-test-teams-proxy-for-location)
	client.On("PostWithResponse", mock.Anything, "configuration/v1/teams_proxy/", mock.MatchedBy(func(req *config.TeamsProxyCreateRequest) bool {
		return req.Name == "tf-test-teams-proxy-for-location"
	}), mock.Anything).Return(&types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/teams_proxy/16/",
	}, nil).Run(func(args mock.Arguments) {
		req := args.Get(2).(*config.TeamsProxyCreateRequest)
		*mockTeams = config.TeamsProxy{
			ID:                    16,
			ResourceURI:           "/api/admin/configuration/v1/teams_proxy/16/",
			Name:                  req.Name,
			Address:               req.Address,
			Port:                  req.Port,
			AzureTenant:           req.AzureTenant,
			MinNumberOfInstances:  req.MinNumberOfInstances,
			NotificationsEnabled:  req.NotificationsEnabled,
			NotificationsQueue:    req.NotificationsQueue,
			Description:           req.Description,
		}
	}).Once()

	// Step 4: Recreate Teams Proxy (tf-test-teams-proxy-for-location)
	client.On("PostWithResponse", mock.Anything, "configuration/v1/teams_proxy/", mock.MatchedBy(func(req *config.TeamsProxyCreateRequest) bool {
		return req.Name == "tf-test-teams-proxy-for-location"
	}), mock.Anything).Return(&types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/teams_proxy/16/",
	}, nil).Run(func(args mock.Arguments) {
		req := args.Get(2).(*config.TeamsProxyCreateRequest)
		*mockTeams = config.TeamsProxy{
			ID:                    16,
			ResourceURI:           "/api/admin/configuration/v1/teams_proxy/16/",
			Name:                  req.Name,
			Address:               req.Address,
			Port:                  req.Port,
			AzureTenant:           req.AzureTenant,
			MinNumberOfInstances:  req.MinNumberOfInstances,
			NotificationsEnabled:  req.NotificationsEnabled,
			NotificationsQueue:    req.NotificationsQueue,
			Description:           req.Description,
		}
	}).Once()

	client.On("GetJSON", mock.Anything, "configuration/v1/teams_proxy/16/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		teams := args.Get(3).(*config.TeamsProxy)
		*teams = *mockTeams
	}).Maybe()

	client.On("DeleteJSON", mock.Anything, "configuration/v1/teams_proxy/16/", mock.Anything).Return(nil).Maybe()



	// Mock System Locations
	// Mock System Locations (for circular references and main location)
	// Register in order: test1 (ID 17), test2 (ID 18), test3 (ID 19), then main-location (ID 123)
	mockSysLoc1 := &config.SystemLocation{}
	mockSysLoc2 := &config.SystemLocation{}
	mockSysLoc3 := &config.SystemLocation{}
	mockState := &config.SystemLocation{}

	// Step 1: Create test1, test2, test3
	// System Location 1: tf-test 1 (step 1)
	client.On("PostWithResponse", mock.Anything, "configuration/v1/system_location/", mock.MatchedBy(func(req *config.SystemLocationCreateRequest) bool {
		return req.Name == "tf-test 1"
	}), mock.Anything).Return(&types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/system_location/17/",
	}, nil).Run(func(args mock.Arguments) {
		req := args.Get(2).(*config.SystemLocationCreateRequest)
		*mockSysLoc1 = config.SystemLocation{
			ID:                         17,
			ResourceURI:                "/api/admin/configuration/v1/system_location/17/",
			Name:                       req.Name,
			Description:                req.Description,
			MTU:                        1500,
			MediaQoS:                   test.IntPtr(0),
			SignallingQoS:              test.IntPtr(0),
			BDPMPinChecksEnabled:       "GLOBAL",
			BDPMScanQuarantineEnabled:  "GLOBAL",
		}
	}).Once()
	// System Location 2: tf-test 2 (step 1)
	client.On("PostWithResponse", mock.Anything, "configuration/v1/system_location/", mock.MatchedBy(func(req *config.SystemLocationCreateRequest) bool {
		return req.Name == "tf-test 2"
	}), mock.Anything).Return(&types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/system_location/18/",
	}, nil).Run(func(args mock.Arguments) {
		req := args.Get(2).(*config.SystemLocationCreateRequest)
		*mockSysLoc2 = config.SystemLocation{
			ID:                         18,
			ResourceURI:                "/api/admin/configuration/v1/system_location/18/",
			Name:                       req.Name,
			Description:                req.Description,
			MTU:                        1500,
			MediaQoS:                   test.IntPtr(0),
			SignallingQoS:              test.IntPtr(0),
			BDPMPinChecksEnabled:       "GLOBAL",
			BDPMScanQuarantineEnabled:  "GLOBAL",
		}
	}).Once()
	// System Location 3: tf-test 3 (step 1)
	client.On("PostWithResponse", mock.Anything, "configuration/v1/system_location/", mock.MatchedBy(func(req *config.SystemLocationCreateRequest) bool {
		return req.Name == "tf-test 3"
	}), mock.Anything).Return(&types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/system_location/19/",
	}, nil).Run(func(args mock.Arguments) {
		req := args.Get(2).(*config.SystemLocationCreateRequest)
		*mockSysLoc3 = config.SystemLocation{
			ID:                         19,
			ResourceURI:                "/api/admin/configuration/v1/system_location/19/",
			Name:                       req.Name,
			Description:                req.Description,
			MTU:                        1500,
			MediaQoS:                   test.IntPtr(0),
			SignallingQoS:              test.IntPtr(0),
			BDPMPinChecksEnabled:       "GLOBAL",
			BDPMScanQuarantineEnabled:  "GLOBAL",
		}
	}).Once()

	// Step 4: Recreate test1, test2, test3 after destroy
	// System Location 1: tf-test 1 (step 4)
	client.On("PostWithResponse", mock.Anything, "configuration/v1/system_location/", mock.MatchedBy(func(req *config.SystemLocationCreateRequest) bool {
		return req.Name == "tf-test 1"
	}), mock.Anything).Return(&types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/system_location/17/",
	}, nil).Run(func(args mock.Arguments) {
		req := args.Get(2).(*config.SystemLocationCreateRequest)
		*mockSysLoc1 = config.SystemLocation{
			ID:                         17,
			ResourceURI:                "/api/admin/configuration/v1/system_location/17/",
			Name:                       req.Name,
			Description:                req.Description,
			MTU:                        1500,
			MediaQoS:                   test.IntPtr(0),
			SignallingQoS:              test.IntPtr(0),
			BDPMPinChecksEnabled:       "GLOBAL",
			BDPMScanQuarantineEnabled:  "GLOBAL",
		}
	}).Once()
	// System Location 2: tf-test 2 (step 4)
	client.On("PostWithResponse", mock.Anything, "configuration/v1/system_location/", mock.MatchedBy(func(req *config.SystemLocationCreateRequest) bool {
		return req.Name == "tf-test 2"
	}), mock.Anything).Return(&types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/system_location/18/",
	}, nil).Run(func(args mock.Arguments) {
		req := args.Get(2).(*config.SystemLocationCreateRequest)
		*mockSysLoc2 = config.SystemLocation{
			ID:                         18,
			ResourceURI:                "/api/admin/configuration/v1/system_location/18/",
			Name:                       req.Name,
			Description:                req.Description,
			MTU:                        1500,
			MediaQoS:                   test.IntPtr(0),
			SignallingQoS:              test.IntPtr(0),
			BDPMPinChecksEnabled:       "GLOBAL",
			BDPMScanQuarantineEnabled:  "GLOBAL",
		}
	}).Once()
	// System Location 3: tf-test 3 (step 4)
	client.On("PostWithResponse", mock.Anything, "configuration/v1/system_location/", mock.MatchedBy(func(req *config.SystemLocationCreateRequest) bool {
		return req.Name == "tf-test 3"
	}), mock.Anything).Return(&types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/system_location/19/",
	}, nil).Run(func(args mock.Arguments) {
		req := args.Get(2).(*config.SystemLocationCreateRequest)
		*mockSysLoc3 = config.SystemLocation{
			ID:                         19,
			ResourceURI:                "/api/admin/configuration/v1/system_location/19/",
			Name:                       req.Name,
			Description:                req.Description,
			MTU:                        1500,
			MediaQoS:                   test.IntPtr(0),
			SignallingQoS:              test.IntPtr(0),
			BDPMPinChecksEnabled:       "GLOBAL",
			BDPMScanQuarantineEnabled:  "GLOBAL",
		}
	}).Once()

	// Add GetJSON and DeleteJSON mocks for test1, test2, test3
	client.On("GetJSON", mock.Anything, "configuration/v1/system_location/17/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		sysLoc := args.Get(3).(*config.SystemLocation)
		*sysLoc = *mockSysLoc1
	}).Maybe()
	client.On("DeleteJSON", mock.Anything, "configuration/v1/system_location/17/", mock.Anything).Return(nil).Maybe()

	client.On("GetJSON", mock.Anything, "configuration/v1/system_location/18/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		sysLoc := args.Get(3).(*config.SystemLocation)
		*sysLoc = *mockSysLoc2
	}).Maybe()
	client.On("DeleteJSON", mock.Anything, "configuration/v1/system_location/18/", mock.Anything).Return(nil).Maybe()

	client.On("GetJSON", mock.Anything, "configuration/v1/system_location/19/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		sysLoc := args.Get(3).(*config.SystemLocation)
		*sysLoc = *mockSysLoc3
	}).Maybe()
	client.On("DeleteJSON", mock.Anything, "configuration/v1/system_location/19/", mock.Anything).Return(nil).Maybe()

	// Main system location - Step 1: Create with full config
	client.On("PostWithResponse", mock.Anything, "configuration/v1/system_location/", mock.MatchedBy(func(req *config.SystemLocationCreateRequest) bool {
		return req.Name == "tf-test-system-location-full"
	}), mock.Anything).Return(&types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/system_location/123/",
	}, nil).Run(func(args mock.Arguments) {
		req := args.Get(2).(*config.SystemLocationCreateRequest)
		*mockState = config.SystemLocation{
			ID:                          123,
			ResourceURI:                 "/api/admin/configuration/v1/system_location/123/",
			Name:                        req.Name,
			Description:                 req.Description,
			MTU:                         req.MTU,
			MediaQoS:                    req.MediaQoS,
			SignallingQoS:               req.SignallingQoS,
			LocalMSSIPDomain:            req.LocalMSSIPDomain,
			BDPMPinChecksEnabled:        req.BDPMPinChecksEnabled,
			BDPMScanQuarantineEnabled:   req.BDPMScanQuarantineEnabled,
			UseRelayCandidatesOnly:      req.UseRelayCandidatesOnly,
			H323Gatekeeper:              req.H323Gatekeeper,
			HTTPProxy:                   req.HTTPProxy,
			MSSIPProxy:                  req.MSSIPProxy,
			PolicyServer:                req.PolicyServer,
			SIPProxy:                    req.SIPProxy,
			SNMPNetworkManagementSystem: req.SNMPNetworkManagementSystem,
			STUNServer:                  req.STUNServer,
			TeamsProxy:                  req.TeamsProxy,
			TURNServer:                  req.TURNServer,
			OverflowLocation1:           req.OverflowLocation1,
			OverflowLocation2:           req.OverflowLocation2,
			TranscodingLocation:         req.TranscodingLocation,
			LiveCaptionsDialOut1:        req.LiveCaptionsDialOut1,
			LiveCaptionsDialOut2:        req.LiveCaptionsDialOut2,
			LiveCaptionsDialOut3:        req.LiveCaptionsDialOut3,
			ClientSTUNServers:           req.ClientSTUNServers,
			ClientTURNServers:           req.ClientTURNServers,
		}
		mockState.DNSServers = convertURIsToResources[config.DNSServer](req.DNSServers)
		mockState.NTPServers = convertURIsToResources[config.NTPServer](req.NTPServers)
		mockState.SyslogServers = convertURIsToResources[config.SyslogServer](req.SyslogServers)
		mockState.EventSinks = convertURIsToResources[config.EventSink](req.EventSinks)
	}).Once()

	// Main system location - Step 4: Create with min config (after destroy)
	client.On("PostWithResponse", mock.Anything, "configuration/v1/system_location/", mock.MatchedBy(func(req *config.SystemLocationCreateRequest) bool {
		return req.Name == "tf-test-system-location-min"
	}), mock.Anything).Return(&types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/system_location/123/",
	}, nil).Run(func(args mock.Arguments) {
		req := args.Get(2).(*config.SystemLocationCreateRequest)
		*mockState = config.SystemLocation{
			ID:                          123,
			ResourceURI:                 "/api/admin/configuration/v1/system_location/123/",
			Name:                        req.Name,
			Description:                 req.Description,
			MTU:                         req.MTU,
			MediaQoS:                    req.MediaQoS,
			SignallingQoS:               req.SignallingQoS,
			LocalMSSIPDomain:            req.LocalMSSIPDomain,
			BDPMPinChecksEnabled:        req.BDPMPinChecksEnabled,
			BDPMScanQuarantineEnabled:   req.BDPMScanQuarantineEnabled,
			UseRelayCandidatesOnly:      req.UseRelayCandidatesOnly,
			H323Gatekeeper:              req.H323Gatekeeper,
			HTTPProxy:                   req.HTTPProxy,
			MSSIPProxy:                  req.MSSIPProxy,
			PolicyServer:                req.PolicyServer,
			SIPProxy:                    req.SIPProxy,
			SNMPNetworkManagementSystem: req.SNMPNetworkManagementSystem,
			STUNServer:                  req.STUNServer,
			TeamsProxy:                  req.TeamsProxy,
			TURNServer:                  req.TURNServer,
			OverflowLocation1:           req.OverflowLocation1,
			OverflowLocation2:           req.OverflowLocation2,
			TranscodingLocation:         req.TranscodingLocation,
			LiveCaptionsDialOut1:        req.LiveCaptionsDialOut1,
			LiveCaptionsDialOut2:        req.LiveCaptionsDialOut2,
			LiveCaptionsDialOut3:        req.LiveCaptionsDialOut3,
			ClientSTUNServers:           req.ClientSTUNServers,
			ClientTURNServers:           req.ClientTURNServers,
		}
		mockState.DNSServers = convertURIsToResources[config.DNSServer](req.DNSServers)
		mockState.NTPServers = convertURIsToResources[config.NTPServer](req.NTPServers)
		mockState.SyslogServers = convertURIsToResources[config.SyslogServer](req.SyslogServers)
		mockState.EventSinks = convertURIsToResources[config.EventSink](req.EventSinks)
	}).Once()

	// Mock the system_location update for step 2 (update to min config) and step 5 (update to full config)
	client.On("PutJSON", mock.Anything, "configuration/v1/system_location/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		updateReq := args.Get(2).(*config.SystemLocationUpdateRequest)
		sysLoc := args.Get(3).(*config.SystemLocation)

		// Update the mockState with new values
		mockState.Name = updateReq.Name
		mockState.Description = updateReq.Description
		mockState.MTU = updateReq.MTU
		mockState.MediaQoS = updateReq.MediaQoS
		mockState.SignallingQoS = updateReq.SignallingQoS
		mockState.LocalMSSIPDomain = updateReq.LocalMSSIPDomain
		mockState.BDPMPinChecksEnabled = updateReq.BDPMPinChecksEnabled
		mockState.BDPMScanQuarantineEnabled = updateReq.BDPMScanQuarantineEnabled
		mockState.UseRelayCandidatesOnly = updateReq.UseRelayCandidatesOnly
		mockState.H323Gatekeeper = updateReq.H323Gatekeeper
		mockState.SNMPNetworkManagementSystem = updateReq.SNMPNetworkManagementSystem
		mockState.SIPProxy = updateReq.SIPProxy
		mockState.HTTPProxy = updateReq.HTTPProxy
		mockState.MSSIPProxy = updateReq.MSSIPProxy
		mockState.TeamsProxy = updateReq.TeamsProxy
		mockState.STUNServer = updateReq.STUNServer
		mockState.TURNServer = updateReq.TURNServer
		mockState.PolicyServer = updateReq.PolicyServer
		mockState.OverflowLocation1 = updateReq.OverflowLocation1
		mockState.OverflowLocation2 = updateReq.OverflowLocation2
		mockState.TranscodingLocation = updateReq.TranscodingLocation
		mockState.LiveCaptionsDialOut1 = updateReq.LiveCaptionsDialOut1
		mockState.LiveCaptionsDialOut2 = updateReq.LiveCaptionsDialOut2
		mockState.LiveCaptionsDialOut3 = updateReq.LiveCaptionsDialOut3
		mockState.ClientSTUNServers = updateReq.ClientSTUNServers
		mockState.ClientTURNServers = updateReq.ClientTURNServers

		// Update collections
		mockState.DNSServers = convertURIsToResources[config.DNSServer](updateReq.DNSServers)
		mockState.NTPServers = convertURIsToResources[config.NTPServer](updateReq.NTPServers)
		mockState.SyslogServers = convertURIsToResources[config.SyslogServer](updateReq.SyslogServers)
		mockState.EventSinks = convertURIsToResources[config.EventSink](updateReq.EventSinks)

		// Return updated state
		*sysLoc = *mockState
	}).Twice()

	// Mock the GetJSON calls for reading the updated system_location after each update (steps 2 and 5)
	client.On("GetJSON", mock.Anything, "configuration/v1/system_location/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		sysLoc := args.Get(3).(*config.SystemLocation)
		*sysLoc = *mockState
	}).Maybe()

	// Mock the DeleteSystemLocation API call
	client.On("DeleteJSON", mock.Anything, "configuration/v1/system_location/123/", mock.Anything).Return(nil)

	testInfinitySystemLocation(t, client)
}

func testInfinitySystemLocation(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: getTestProtoV5ProviderFactories(client),
		Steps: []resource.TestStep{
			{
				Config: test.LoadTestFolder(t, "resource_infinity_system_location_full"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_system_location.main-location", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_system_location.main-location", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_system_location.main-location", "name", "tf-test-system-location-full"),
					resource.TestCheckResourceAttr("pexip_infinity_system_location.main-location", "description", "Full configuration test location"),
					resource.TestCheckResourceAttr("pexip_infinity_system_location.main-location", "mtu", "1460"),
					resource.TestCheckResourceAttr("pexip_infinity_system_location.main-location", "media_qos", "46"),
					resource.TestCheckResourceAttr("pexip_infinity_system_location.main-location", "signalling_qos", "24"),
					resource.TestCheckResourceAttr("pexip_infinity_system_location.main-location", "local_mssip_domain", "test-mssip.pexvclab.com"),
					resource.TestCheckResourceAttr("pexip_infinity_system_location.main-location", "bdpm_pin_checks_enabled", "ON"),
					resource.TestCheckResourceAttr("pexip_infinity_system_location.main-location", "bdpm_scan_quarantine_enabled", "ON"),
					resource.TestCheckResourceAttr("pexip_infinity_system_location.main-location", "use_relay_candidates_only", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_system_location.main-location", "dns_servers.#", "2"),
					resource.TestCheckResourceAttr("pexip_infinity_system_location.main-location", "ntp_servers.#", "1"),
					resource.TestCheckResourceAttr("pexip_infinity_system_location.main-location", "client_stun_servers.#", "2"),
					resource.TestCheckResourceAttr("pexip_infinity_system_location.main-location", "client_turn_servers.#", "2"),
					resource.TestCheckResourceAttr("pexip_infinity_system_location.main-location", "event_sinks.#", "1"),
					resource.TestCheckResourceAttr("pexip_infinity_system_location.main-location", "syslog_servers.#", "0"),
					resource.TestCheckResourceAttrSet("pexip_infinity_system_location.main-location", "h323_gatekeeper"),
					resource.TestCheckResourceAttrSet("pexip_infinity_system_location.main-location", "http_proxy"),
					resource.TestCheckResourceAttrSet("pexip_infinity_system_location.main-location", "mssip_proxy"),
					resource.TestCheckResourceAttrSet("pexip_infinity_system_location.main-location", "policy_server"),
					resource.TestCheckResourceAttrSet("pexip_infinity_system_location.main-location", "sip_proxy"),
					resource.TestCheckResourceAttrSet("pexip_infinity_system_location.main-location", "snmp_network_management_system"),
					resource.TestCheckResourceAttrSet("pexip_infinity_system_location.main-location", "stun_server"),
					resource.TestCheckResourceAttrSet("pexip_infinity_system_location.main-location", "teams_proxy"),
					resource.TestCheckResourceAttrSet("pexip_infinity_system_location.main-location", "turn_server"),
					resource.TestCheckResourceAttrSet("pexip_infinity_system_location.main-location", "overflow_location1"),
					resource.TestCheckResourceAttrSet("pexip_infinity_system_location.main-location", "overflow_location2"),
					resource.TestCheckResourceAttrSet("pexip_infinity_system_location.main-location", "transcoding_location"),
					resource.TestCheckResourceAttrSet("pexip_infinity_system_location.main-location", "live_captions_dial_out_1"),
					resource.TestCheckResourceAttrSet("pexip_infinity_system_location.main-location", "live_captions_dial_out_2"),
					resource.TestCheckResourceAttrSet("pexip_infinity_system_location.main-location", "live_captions_dial_out_3"),
				),
			},
			{
				Config: test.LoadTestFolder(t, "resource_infinity_system_location_min"),
				Check: resource.ComposeTestCheckFunc(
					// IDs and required fields
					resource.TestCheckResourceAttrSet("pexip_infinity_system_location.main-location", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_system_location.main-location", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_system_location.main-location", "name", "tf-test-system-location-min"),

					// Optional fields cleared - verify defaults
					resource.TestCheckResourceAttr("pexip_infinity_system_location.main-location", "description", ""),
					resource.TestCheckResourceAttr("pexip_infinity_system_location.main-location", "mtu", "1500"),
					resource.TestCheckResourceAttr("pexip_infinity_system_location.main-location", "media_qos", "0"),
					resource.TestCheckResourceAttr("pexip_infinity_system_location.main-location", "signalling_qos", "0"),
					resource.TestCheckResourceAttr("pexip_infinity_system_location.main-location", "local_mssip_domain", ""),
					resource.TestCheckResourceAttr("pexip_infinity_system_location.main-location", "bdpm_pin_checks_enabled", "GLOBAL"),
					resource.TestCheckResourceAttr("pexip_infinity_system_location.main-location", "bdpm_scan_quarantine_enabled", "GLOBAL"),
					resource.TestCheckResourceAttr("pexip_infinity_system_location.main-location", "use_relay_candidates_only", "false"),

					// Collections cleared
					resource.TestCheckResourceAttr("pexip_infinity_system_location.main-location", "dns_servers.#", "0"),
					resource.TestCheckResourceAttr("pexip_infinity_system_location.main-location", "ntp_servers.#", "0"),
					resource.TestCheckResourceAttr("pexip_infinity_system_location.main-location", "client_stun_servers.#", "0"),
					resource.TestCheckResourceAttr("pexip_infinity_system_location.main-location", "client_turn_servers.#", "0"),
					resource.TestCheckResourceAttr("pexip_infinity_system_location.main-location", "event_sinks.#", "0"),
					resource.TestCheckResourceAttr("pexip_infinity_system_location.main-location", "syslog_servers.#", "0"),

					// Nullable fields cleared via UPDATE - should be null
					resource.TestCheckNoResourceAttr("pexip_infinity_system_location.main-location", "h323_gatekeeper"),
					resource.TestCheckNoResourceAttr("pexip_infinity_system_location.main-location", "http_proxy"),
					resource.TestCheckNoResourceAttr("pexip_infinity_system_location.main-location", "mssip_proxy"),
					resource.TestCheckNoResourceAttr("pexip_infinity_system_location.main-location", "policy_server"),
					resource.TestCheckNoResourceAttr("pexip_infinity_system_location.main-location", "sip_proxy"),
					resource.TestCheckNoResourceAttr("pexip_infinity_system_location.main-location", "snmp_network_management_system"),
					resource.TestCheckNoResourceAttr("pexip_infinity_system_location.main-location", "stun_server"),
					resource.TestCheckNoResourceAttr("pexip_infinity_system_location.main-location", "teams_proxy"),
					resource.TestCheckNoResourceAttr("pexip_infinity_system_location.main-location", "turn_server"),
					resource.TestCheckNoResourceAttr("pexip_infinity_system_location.main-location", "overflow_location1"),
					resource.TestCheckNoResourceAttr("pexip_infinity_system_location.main-location", "overflow_location2"),
					resource.TestCheckNoResourceAttr("pexip_infinity_system_location.main-location", "transcoding_location"),
					resource.TestCheckNoResourceAttr("pexip_infinity_system_location.main-location", "live_captions_dial_out_1"),
					resource.TestCheckNoResourceAttr("pexip_infinity_system_location.main-location", "live_captions_dial_out_2"),
					resource.TestCheckNoResourceAttr("pexip_infinity_system_location.main-location", "live_captions_dial_out_3"),
				),
			},
			{
				// Step 3: Destroy and recreate with minimal config
				Config:  test.LoadTestFolder(t, "resource_infinity_system_location_min"),
				Destroy: true,
			},
			{
				// Step 4: Recreate with minimal config (after destroy)
				Config: test.LoadTestFolder(t, "resource_infinity_system_location_min"),
				Check: resource.ComposeTestCheckFunc(
					// IDs and required fields
					resource.TestCheckResourceAttrSet("pexip_infinity_system_location.main-location", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_system_location.main-location", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_system_location.main-location", "name", "tf-test-system-location-min"),

					// Optional fields - verify defaults
					resource.TestCheckResourceAttr("pexip_infinity_system_location.main-location", "description", ""),
					resource.TestCheckResourceAttr("pexip_infinity_system_location.main-location", "mtu", "1500"),
					resource.TestCheckResourceAttr("pexip_infinity_system_location.main-location", "media_qos", "0"),
					resource.TestCheckResourceAttr("pexip_infinity_system_location.main-location", "signalling_qos", "0"),
					resource.TestCheckResourceAttr("pexip_infinity_system_location.main-location", "local_mssip_domain", ""),
					resource.TestCheckResourceAttr("pexip_infinity_system_location.main-location", "bdpm_pin_checks_enabled", "GLOBAL"),
					resource.TestCheckResourceAttr("pexip_infinity_system_location.main-location", "bdpm_scan_quarantine_enabled", "GLOBAL"),
					resource.TestCheckResourceAttr("pexip_infinity_system_location.main-location", "use_relay_candidates_only", "false"),

					// Collections
					resource.TestCheckResourceAttr("pexip_infinity_system_location.main-location", "dns_servers.#", "0"),
					resource.TestCheckResourceAttr("pexip_infinity_system_location.main-location", "ntp_servers.#", "0"),
					resource.TestCheckResourceAttr("pexip_infinity_system_location.main-location", "client_stun_servers.#", "0"),
					resource.TestCheckResourceAttr("pexip_infinity_system_location.main-location", "client_turn_servers.#", "0"),
					resource.TestCheckResourceAttr("pexip_infinity_system_location.main-location", "event_sinks.#", "0"),
					resource.TestCheckResourceAttr("pexip_infinity_system_location.main-location", "syslog_servers.#", "0"),
				),
			},
			{
				// Step 5: Update to full config
				Config: test.LoadTestFolder(t, "resource_infinity_system_location_full"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_system_location.main-location", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_system_location.main-location", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_system_location.main-location", "name", "tf-test-system-location-full"),
					resource.TestCheckResourceAttr("pexip_infinity_system_location.main-location", "description", "Full configuration test location"),
					resource.TestCheckResourceAttr("pexip_infinity_system_location.main-location", "mtu", "1460"),
					resource.TestCheckResourceAttr("pexip_infinity_system_location.main-location", "media_qos", "46"),
					resource.TestCheckResourceAttr("pexip_infinity_system_location.main-location", "signalling_qos", "24"),
					resource.TestCheckResourceAttr("pexip_infinity_system_location.main-location", "local_mssip_domain", "test-mssip.pexvclab.com"),
					resource.TestCheckResourceAttr("pexip_infinity_system_location.main-location", "bdpm_pin_checks_enabled", "ON"),
					resource.TestCheckResourceAttr("pexip_infinity_system_location.main-location", "bdpm_scan_quarantine_enabled", "ON"),
					resource.TestCheckResourceAttr("pexip_infinity_system_location.main-location", "use_relay_candidates_only", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_system_location.main-location", "dns_servers.#", "2"),
					resource.TestCheckResourceAttr("pexip_infinity_system_location.main-location", "ntp_servers.#", "1"),
					resource.TestCheckResourceAttr("pexip_infinity_system_location.main-location", "client_stun_servers.#", "2"),
					resource.TestCheckResourceAttr("pexip_infinity_system_location.main-location", "client_turn_servers.#", "2"),
					resource.TestCheckResourceAttr("pexip_infinity_system_location.main-location", "event_sinks.#", "1"),
					resource.TestCheckResourceAttr("pexip_infinity_system_location.main-location", "syslog_servers.#", "0"),
					resource.TestCheckResourceAttrSet("pexip_infinity_system_location.main-location", "h323_gatekeeper"),
					resource.TestCheckResourceAttrSet("pexip_infinity_system_location.main-location", "http_proxy"),
					resource.TestCheckResourceAttrSet("pexip_infinity_system_location.main-location", "mssip_proxy"),
					resource.TestCheckResourceAttrSet("pexip_infinity_system_location.main-location", "policy_server"),
					resource.TestCheckResourceAttrSet("pexip_infinity_system_location.main-location", "sip_proxy"),
					resource.TestCheckResourceAttrSet("pexip_infinity_system_location.main-location", "snmp_network_management_system"),
					resource.TestCheckResourceAttrSet("pexip_infinity_system_location.main-location", "stun_server"),
					resource.TestCheckResourceAttrSet("pexip_infinity_system_location.main-location", "teams_proxy"),
					resource.TestCheckResourceAttrSet("pexip_infinity_system_location.main-location", "turn_server"),
					resource.TestCheckResourceAttrSet("pexip_infinity_system_location.main-location", "overflow_location1"),
					resource.TestCheckResourceAttrSet("pexip_infinity_system_location.main-location", "overflow_location2"),
					resource.TestCheckResourceAttrSet("pexip_infinity_system_location.main-location", "transcoding_location"),
					resource.TestCheckResourceAttrSet("pexip_infinity_system_location.main-location", "live_captions_dial_out_1"),
					resource.TestCheckResourceAttrSet("pexip_infinity_system_location.main-location", "live_captions_dial_out_2"),
					resource.TestCheckResourceAttrSet("pexip_infinity_system_location.main-location", "live_captions_dial_out_3"),
				),
			},
		},
	})
}
