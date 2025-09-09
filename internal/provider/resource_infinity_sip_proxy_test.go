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

func TestInfinitySIPProxy(t *testing.T) {
	t.Parallel()
	_ = os.Setenv("TF_ACC", "1")

	// Create a mock client and set up expectations
	client := infinity.NewClientMock()

	// Mock the CreateSIPProxy API call
	createResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/sip_proxy/123/",
	}
	client.On("PostWithResponse", mock.Anything, "configuration/v1/sip_proxy/", mock.Anything, mock.Anything).Return(createResponse, nil)

	// Shared state for mocking
	port := 8080
	mockState := &config.SIPProxy{
		ID:          123,
		ResourceURI: "/api/admin/configuration/v1/sip_proxy/123/",
		Name:        "test-sip-proxy",
		Description: "Test SIP Proxy",
		Address:     "test-sip-proxy.dev.pexip.network",
		Port:        &port,
		Transport:   "tcp",
	}

	// Mock the GetSIPProxy API call for Read operations
	client.On("GetJSON", mock.Anything, "configuration/v1/sip_proxy/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		sipProxy := args.Get(3).(*config.SIPProxy)
		*sipProxy = *mockState
	}).Maybe()

	// Mock the UpdateSIPProxy API call
	client.On("PutJSON", mock.Anything, "configuration/v1/sip_proxy/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		updateRequest := args.Get(2).(*config.SIPProxyUpdateRequest)
		sipProxy := args.Get(3).(*config.SIPProxy)

		// Update mock state
		mockState.Name = updateRequest.Name
		mockState.Description = updateRequest.Description
		mockState.Address = updateRequest.Address
		mockState.Transport = updateRequest.Transport
		if updateRequest.Port != nil {
			mockState.Port = updateRequest.Port
		}

		// Return updated state
		*sipProxy = *mockState
	}).Maybe()

	// Mock the DeleteSIPProxy API call
	client.On("DeleteJSON", mock.Anything, mock.MatchedBy(func(path string) bool {
		return path == "configuration/v1/sip_proxy/123/"
	}), mock.Anything).Return(nil)

	testInfinitySIPProxy(t, client)
}

func testInfinitySIPProxy(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: getTestProtoV5ProviderFactories(client),
		Steps: []resource.TestStep{
			{
				Config: test.LoadTestFolder(t, "resource_infinity_sip_proxy_basic"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_sip_proxy.sip-proxy-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_sip_proxy.sip-proxy-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_sip_proxy.sip-proxy-test", "name", "test-sip-proxy"),
					resource.TestCheckResourceAttr("pexip_infinity_sip_proxy.sip-proxy-test", "description", "Test SIP Proxy"),
					resource.TestCheckResourceAttr("pexip_infinity_sip_proxy.sip-proxy-test", "address", "test-sip-proxy.dev.pexip.network"),
					resource.TestCheckResourceAttr("pexip_infinity_sip_proxy.sip-proxy-test", "port", "8080"),
					resource.TestCheckResourceAttr("pexip_infinity_sip_proxy.sip-proxy-test", "transport", "tcp"),
				),
			},
			{
				Config: test.LoadTestFolder(t, "resource_infinity_sip_proxy_basic_updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_sip_proxy.sip-proxy-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_sip_proxy.sip-proxy-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_sip_proxy.sip-proxy-test", "name", "test-sip-proxy"),
					resource.TestCheckResourceAttr("pexip_infinity_sip_proxy.sip-proxy-test", "description", "Test SIP Proxy"),
					resource.TestCheckResourceAttr("pexip_infinity_sip_proxy.sip-proxy-test", "address", "test-sip-proxy.dev.pexip.network"),
					resource.TestCheckResourceAttr("pexip_infinity_sip_proxy.sip-proxy-test", "port", "8081"),
					resource.TestCheckResourceAttr("pexip_infinity_sip_proxy.sip-proxy-test", "transport", "tls"),
				),
			},
		},
	})
}
