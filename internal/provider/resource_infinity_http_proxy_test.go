/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package provider

import (
	"os"
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
	_ = os.Setenv("TF_ACC", "1")

	// Create a mock client and set up expectations
	client := infinity.NewClientMock()

	// Shared state for mocking - initialize with defaults
	mockState := &config.HTTPProxy{
		ID:          1,
		ResourceURI: "/api/admin/configuration/v1/http_proxy/1/",
		Name:        "tf-test-http-proxy",
		Address:     "proxy-min.example.com",
		Port:        intPtr(8080),
		Protocol:    "http",
		Username:    "",
		Password:    "",
	}

	// Step 1: Create with full config
	client.On("PostWithResponse", mock.Anything, "configuration/v1/http_proxy/", mock.Anything, mock.Anything).Return(&types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/http_proxy/1/",
	}, nil).Run(func(args mock.Arguments) {
		req := args.Get(2).(*config.HTTPProxyCreateRequest)
		mockState.Name = req.Name
		mockState.Address = req.Address
		mockState.Port = req.Port
		mockState.Protocol = req.Protocol
		mockState.Username = req.Username
		mockState.Password = req.Password
	}).Once()

	// Step 2: Update to min config
	client.On("PutJSON", mock.Anything, "configuration/v1/http_proxy/1/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		req := args.Get(2).(*config.HTTPProxyUpdateRequest)
		mockState.Name = req.Name
		mockState.Address = req.Address
		mockState.Port = req.Port
		mockState.Protocol = req.Protocol
		mockState.Username = req.Username
		mockState.Password = req.Password
		if args.Get(3) != nil {
			proxy := args.Get(3).(*config.HTTPProxy)
			*proxy = *mockState
		}
	}).Once()

	// Step 3: Delete
	client.On("DeleteJSON", mock.Anything, "configuration/v1/http_proxy/1/", mock.Anything).Return(nil).Maybe()

	// Step 4: Recreate with min config
	client.On("PostWithResponse", mock.Anything, "configuration/v1/http_proxy/", mock.Anything, mock.Anything).Return(&types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/http_proxy/1/",
	}, nil).Run(func(args mock.Arguments) {
		req := args.Get(2).(*config.HTTPProxyCreateRequest)
		mockState.Name = req.Name
		mockState.Address = req.Address
		mockState.Port = req.Port
		mockState.Protocol = req.Protocol
		mockState.Username = req.Username
		mockState.Password = req.Password
	}).Once()

	// Step 5: Update to full config
	client.On("PutJSON", mock.Anything, "configuration/v1/http_proxy/1/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		req := args.Get(2).(*config.HTTPProxyUpdateRequest)
		mockState.Name = req.Name
		mockState.Address = req.Address
		mockState.Port = req.Port
		mockState.Protocol = req.Protocol
		mockState.Username = req.Username
		mockState.Password = req.Password
		if args.Get(3) != nil {
			proxy := args.Get(3).(*config.HTTPProxy)
			*proxy = *mockState
		}
	}).Once()

	// Mock Read operations (GetJSON) - used throughout all steps
	client.On("GetJSON", mock.Anything, "configuration/v1/http_proxy/1/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		proxy := args.Get(3).(*config.HTTPProxy)
		*proxy = *mockState
	}).Maybe()

	testInfinityHTTPProxy(t, client)
}

func testInfinityHTTPProxy(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: getTestProtoV5ProviderFactories(client),
		Steps: []resource.TestStep{
			{
				// Step 1: Create with full config
				Config: test.LoadTestFolder(t, "resource_infinity_http_proxy_full"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_http_proxy.tf-test-http-proxy", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_http_proxy.tf-test-http-proxy", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_http_proxy.tf-test-http-proxy", "name", "tf-test-http-proxy"),
					resource.TestCheckResourceAttr("pexip_infinity_http_proxy.tf-test-http-proxy", "address", "proxy.example.com"),
					resource.TestCheckResourceAttr("pexip_infinity_http_proxy.tf-test-http-proxy", "port", "8081"),
					resource.TestCheckResourceAttr("pexip_infinity_http_proxy.tf-test-http-proxy", "protocol", "http"),
					resource.TestCheckResourceAttr("pexip_infinity_http_proxy.tf-test-http-proxy", "username", "tf-test-user"),
					resource.TestCheckResourceAttr("pexip_infinity_http_proxy.tf-test-http-proxy", "password", "tf-test-password"),
				),
			},
			{
				// Step 2: Update to min config (clear optional fields, reset to defaults)
				Config: test.LoadTestFolder(t, "resource_infinity_http_proxy_min"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_http_proxy.tf-test-http-proxy", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_http_proxy.tf-test-http-proxy", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_http_proxy.tf-test-http-proxy", "name", "tf-test-http-proxy"),
					resource.TestCheckResourceAttr("pexip_infinity_http_proxy.tf-test-http-proxy", "address", "proxy-min.example.com"),
					resource.TestCheckResourceAttr("pexip_infinity_http_proxy.tf-test-http-proxy", "port", "8080"),
					resource.TestCheckResourceAttr("pexip_infinity_http_proxy.tf-test-http-proxy", "protocol", "http"),
					resource.TestCheckResourceAttr("pexip_infinity_http_proxy.tf-test-http-proxy", "username", ""),
					resource.TestCheckResourceAttr("pexip_infinity_http_proxy.tf-test-http-proxy", "password", ""),
				),
			},
			{
				// Step 3: Destroy
				Config:  test.LoadTestFolder(t, "resource_infinity_http_proxy_min"),
				Destroy: true,
			},
			{
				// Step 4: Create with min config (after destroy)
				Config: test.LoadTestFolder(t, "resource_infinity_http_proxy_min"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_http_proxy.tf-test-http-proxy", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_http_proxy.tf-test-http-proxy", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_http_proxy.tf-test-http-proxy", "name", "tf-test-http-proxy"),
					resource.TestCheckResourceAttr("pexip_infinity_http_proxy.tf-test-http-proxy", "address", "proxy-min.example.com"),
					resource.TestCheckResourceAttr("pexip_infinity_http_proxy.tf-test-http-proxy", "port", "8080"),
					resource.TestCheckResourceAttr("pexip_infinity_http_proxy.tf-test-http-proxy", "protocol", "http"),
					resource.TestCheckResourceAttr("pexip_infinity_http_proxy.tf-test-http-proxy", "username", ""),
					resource.TestCheckResourceAttr("pexip_infinity_http_proxy.tf-test-http-proxy", "password", ""),
				),
			},
			{
				// Step 5: Update to full config
				Config: test.LoadTestFolder(t, "resource_infinity_http_proxy_full"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_http_proxy.tf-test-http-proxy", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_http_proxy.tf-test-http-proxy", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_http_proxy.tf-test-http-proxy", "name", "tf-test-http-proxy"),
					resource.TestCheckResourceAttr("pexip_infinity_http_proxy.tf-test-http-proxy", "address", "proxy.example.com"),
					resource.TestCheckResourceAttr("pexip_infinity_http_proxy.tf-test-http-proxy", "port", "8081"),
					resource.TestCheckResourceAttr("pexip_infinity_http_proxy.tf-test-http-proxy", "protocol", "http"),
					resource.TestCheckResourceAttr("pexip_infinity_http_proxy.tf-test-http-proxy", "username", "tf-test-user"),
					resource.TestCheckResourceAttr("pexip_infinity_http_proxy.tf-test-http-proxy", "password", "tf-test-password"),
				),
			},
		},
	})
}
