/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package provider

import (
	"testing"

	"github.com/pexip/go-infinity-sdk/v38/config"
	"github.com/pexip/go-infinity-sdk/v38/types"
	"github.com/stretchr/testify/mock"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/pexip/go-infinity-sdk/v38"

	"github.com/pexip/terraform-provider-pexip/internal/test"
)

// Helper function to convert int to *int
func intPtr(i int) *int {
	return &i
}

func TestInfinityHTTPProxy(t *testing.T) {
	t.Parallel()

	// Create a mock client and set up expectations
	client := infinity.NewClientMock()

	// Mock the CreateHttpproxy API call
	createResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/http_proxy/123/",
	}
	client.On("PostWithResponse", mock.Anything, "configuration/v1/http_proxy/", mock.Anything, mock.Anything).Return(createResponse, nil)

	// Shared state for mocking
	mockState := &config.HTTPProxy{
		ID:          123,
		ResourceURI: "/api/admin/configuration/v1/http_proxy/123/",
		Name:        "http_proxy-test",
		Address:     "test-server.example.com",
		Port:        intPtr(8080),
		Protocol:    "https",
		Username:    "http_proxy-test",
		Password:    "test-value",
	}

	// Mock the GetHttpproxy API call for Read operations
	client.On("GetJSON", mock.Anything, "configuration/v1/http_proxy/123/", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		http_proxy := args.Get(2).(*config.HTTPProxy)
		*http_proxy = *mockState
	}).Maybe()

	// Mock the UpdateHttpproxy API call
	client.On("PutJSON", mock.Anything, "configuration/v1/http_proxy/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		updateRequest := args.Get(2).(*config.HTTPProxyUpdateRequest)
		http_proxy := args.Get(3).(*config.HTTPProxy)

		// Update mock state based on request
		if updateRequest.Name != "" {
			mockState.Name = updateRequest.Name
		}
		if updateRequest.Address != "" {
			mockState.Address = updateRequest.Address
		}
		if updateRequest.Port != nil {
			mockState.Port = updateRequest.Port
		}
		if updateRequest.Protocol != "" {
			mockState.Protocol = updateRequest.Protocol
		}
		if updateRequest.Username != "" {
			mockState.Username = updateRequest.Username
		}
		if updateRequest.Password != "" {
			mockState.Password = updateRequest.Password
		}

		// Return updated state
		*http_proxy = *mockState
	}).Maybe()

	// Mock the DeleteHttpproxy API call
	client.On("DeleteJSON", mock.Anything, mock.MatchedBy(func(path string) bool {
		return path == "configuration/v1/http_proxy/123/"
	}), mock.Anything).Return(nil)

	testInfinityHTTPProxy(t, client)
}

func testInfinityHTTPProxy(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: getTestProtoV5ProviderFactories(client),
		Steps: []resource.TestStep{
			{
				Config: test.LoadTestFolder(t, "resource_infinity_http_proxy_basic"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_http_proxy.http_proxy-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_http_proxy.http_proxy-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_http_proxy.http_proxy-test", "name", "http_proxy-test"),
					resource.TestCheckResourceAttr("pexip_infinity_http_proxy.http_proxy-test", "port", "8080"),
					resource.TestCheckResourceAttr("pexip_infinity_http_proxy.http_proxy-test", "username", "http_proxy-test"),
				),
			},
			{
				Config: test.LoadTestFolder(t, "resource_infinity_http_proxy_basic_updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_http_proxy.http_proxy-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_http_proxy.http_proxy-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_http_proxy.http_proxy-test", "name", "http_proxy-test"),
					resource.TestCheckResourceAttr("pexip_infinity_http_proxy.http_proxy-test", "port", "8081"),
					resource.TestCheckResourceAttr("pexip_infinity_http_proxy.http_proxy-test", "username", "http_proxy-test"),
				),
			},
		},
	})
}
