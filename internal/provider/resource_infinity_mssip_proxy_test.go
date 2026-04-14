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

func TestInfinityMSSIPProxy(t *testing.T) {
	t.Parallel()
	_ = os.Setenv("TF_ACC", "1")

	// Create a mock client and set up expectations
	client := infinity.NewClientMock()

	// Shared state for mocking - initialize with defaults
	defaultPort := 5061
	mockState := &config.MSSIPProxy{
		ID:          1,
		ResourceURI: "/api/admin/configuration/v1/mssip_proxy/1/",
		Name:        "tf-test-mssip-proxy",
		Description: "",
		Address:     "mssip-proxy-min.example.com",
		Port:        &defaultPort,
		Transport:   "tls",
	}

	// Step 1: Create with full config
	client.On("PostWithResponse", mock.Anything, "configuration/v1/mssip_proxy/", mock.Anything, mock.Anything).Return(&types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/mssip_proxy/1/",
	}, nil).Run(func(args mock.Arguments) {
		req := args.Get(2).(*config.MSSIPProxyCreateRequest)
		mockState.Name = req.Name
		mockState.Description = req.Description
		mockState.Address = req.Address
		mockState.Port = req.Port
		mockState.Transport = req.Transport
	}).Once()

	// Step 2: Update to min config
	client.On("PutJSON", mock.Anything, "configuration/v1/mssip_proxy/1/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		req := args.Get(2).(*config.MSSIPProxyUpdateRequest)
		mockState.Name = req.Name
		mockState.Description = req.Description
		mockState.Address = req.Address
		mockState.Port = req.Port
		mockState.Transport = req.Transport
		if args.Get(3) != nil {
			proxy := args.Get(3).(*config.MSSIPProxy)
			*proxy = *mockState
		}
	}).Once()

	// Step 3: Delete
	client.On("DeleteJSON", mock.Anything, "configuration/v1/mssip_proxy/1/", mock.Anything).Return(nil).Maybe()

	// Step 4: Recreate with min config
	client.On("PostWithResponse", mock.Anything, "configuration/v1/mssip_proxy/", mock.Anything, mock.Anything).Return(&types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/mssip_proxy/1/",
	}, nil).Run(func(args mock.Arguments) {
		req := args.Get(2).(*config.MSSIPProxyCreateRequest)
		mockState.Name = req.Name
		mockState.Description = req.Description
		mockState.Address = req.Address
		mockState.Port = req.Port
		mockState.Transport = req.Transport
	}).Once()

	// Step 5: Update to full config
	client.On("PutJSON", mock.Anything, "configuration/v1/mssip_proxy/1/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		req := args.Get(2).(*config.MSSIPProxyUpdateRequest)
		mockState.Name = req.Name
		mockState.Description = req.Description
		mockState.Address = req.Address
		mockState.Port = req.Port
		mockState.Transport = req.Transport
		if args.Get(3) != nil {
			proxy := args.Get(3).(*config.MSSIPProxy)
			*proxy = *mockState
		}
	}).Once()

	// Mock Read operations (GetJSON) - used throughout all steps
	client.On("GetJSON", mock.Anything, "configuration/v1/mssip_proxy/1/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		proxy := args.Get(3).(*config.MSSIPProxy)
		*proxy = *mockState
	}).Maybe()

	testInfinityMSSIPProxy(t, client)
}

func testInfinityMSSIPProxy(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: getTestProtoV6ProviderFactories(client),
		Steps: []resource.TestStep{
			{
				// Step 1: Create with full config
				Config: test.LoadTestFolder(t, "resource_infinity_mssip_proxy_full"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_mssip_proxy.tf-test-mssip-proxy", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_mssip_proxy.tf-test-mssip-proxy", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_mssip_proxy.tf-test-mssip-proxy", "name", "tf-test-mssip-proxy"),
					resource.TestCheckResourceAttr("pexip_infinity_mssip_proxy.tf-test-mssip-proxy", "description", "tf-test MSSIP Proxy Description"),
					resource.TestCheckResourceAttr("pexip_infinity_mssip_proxy.tf-test-mssip-proxy", "address", "mssip-proxy.example.com"),
					resource.TestCheckResourceAttr("pexip_infinity_mssip_proxy.tf-test-mssip-proxy", "port", "5060"),
					resource.TestCheckResourceAttr("pexip_infinity_mssip_proxy.tf-test-mssip-proxy", "transport", "tcp"),
				),
			},
			{
				// Step 2: Update to min config (clear optional fields, reset to defaults)
				Config: test.LoadTestFolder(t, "resource_infinity_mssip_proxy_min"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_mssip_proxy.tf-test-mssip-proxy", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_mssip_proxy.tf-test-mssip-proxy", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_mssip_proxy.tf-test-mssip-proxy", "name", "tf-test-mssip-proxy"),
					resource.TestCheckResourceAttr("pexip_infinity_mssip_proxy.tf-test-mssip-proxy", "description", ""),
					resource.TestCheckResourceAttr("pexip_infinity_mssip_proxy.tf-test-mssip-proxy", "address", "mssip-proxy-min.example.com"),
					resource.TestCheckResourceAttr("pexip_infinity_mssip_proxy.tf-test-mssip-proxy", "port", "5061"),
					resource.TestCheckResourceAttr("pexip_infinity_mssip_proxy.tf-test-mssip-proxy", "transport", "tls"),
				),
			},
			{
				// Step 3: Destroy
				Config:  test.LoadTestFolder(t, "resource_infinity_mssip_proxy_min"),
				Destroy: true,
			},
			{
				// Step 4: Create with min config (after destroy)
				Config: test.LoadTestFolder(t, "resource_infinity_mssip_proxy_min"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_mssip_proxy.tf-test-mssip-proxy", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_mssip_proxy.tf-test-mssip-proxy", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_mssip_proxy.tf-test-mssip-proxy", "name", "tf-test-mssip-proxy"),
					resource.TestCheckResourceAttr("pexip_infinity_mssip_proxy.tf-test-mssip-proxy", "description", ""),
					resource.TestCheckResourceAttr("pexip_infinity_mssip_proxy.tf-test-mssip-proxy", "address", "mssip-proxy-min.example.com"),
					resource.TestCheckResourceAttr("pexip_infinity_mssip_proxy.tf-test-mssip-proxy", "port", "5061"),
					resource.TestCheckResourceAttr("pexip_infinity_mssip_proxy.tf-test-mssip-proxy", "transport", "tls"),
				),
			},
			{
				// Step 5: Update to full config
				Config: test.LoadTestFolder(t, "resource_infinity_mssip_proxy_full"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_mssip_proxy.tf-test-mssip-proxy", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_mssip_proxy.tf-test-mssip-proxy", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_mssip_proxy.tf-test-mssip-proxy", "name", "tf-test-mssip-proxy"),
					resource.TestCheckResourceAttr("pexip_infinity_mssip_proxy.tf-test-mssip-proxy", "description", "tf-test MSSIP Proxy Description"),
					resource.TestCheckResourceAttr("pexip_infinity_mssip_proxy.tf-test-mssip-proxy", "address", "mssip-proxy.example.com"),
					resource.TestCheckResourceAttr("pexip_infinity_mssip_proxy.tf-test-mssip-proxy", "port", "5060"),
					resource.TestCheckResourceAttr("pexip_infinity_mssip_proxy.tf-test-mssip-proxy", "transport", "tcp"),
				),
			},
		},
	})
}
