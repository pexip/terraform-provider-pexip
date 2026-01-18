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

	// Mock the CreateMSSIPProxy API call
	createResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/mssip_proxy/123/",
	}
	client.On("PostWithResponse", mock.Anything, "configuration/v1/mssip_proxy/", mock.Anything, mock.Anything).Return(createResponse, nil)

	// Shared state for mocking
	mockState := &config.MSSIPProxy{
		ID:          123,
		ResourceURI: "/api/admin/configuration/v1/mssip_proxy/123/",
		Name:        "mssip-proxy-test",
		Description: "Test MSSIP proxy",
		Address:     "test-mssip-proxy.dev.pexip.network",
		Port:        func() *int { port := 5060; return &port }(),
		Transport:   "tcp",
	}

	// Mock the GetMSSIPProxy API call for Read operations
	client.On("GetJSON", mock.Anything, "configuration/v1/mssip_proxy/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		mssipProxy := args.Get(3).(*config.MSSIPProxy)
		*mssipProxy = *mockState
	}).Maybe()

	// Mock the UpdateMSSIPProxy API call
	client.On("PutJSON", mock.Anything, "configuration/v1/mssip_proxy/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		updateRequest := args.Get(2).(*config.MSSIPProxyUpdateRequest)
		mssipProxy := args.Get(3).(*config.MSSIPProxy)

		// Update mock state
		mockState.Name = updateRequest.Name
		if updateRequest.Description != "" {
			mockState.Description = updateRequest.Description
		}
		if updateRequest.Address != "" {
			mockState.Address = updateRequest.Address
		}
		if updateRequest.Port != nil {
			mockState.Port = updateRequest.Port
		}
		if updateRequest.Transport != "" {
			mockState.Transport = updateRequest.Transport
		}

		// Return updated state
		*mssipProxy = *mockState
	}).Maybe()

	// Mock the DeleteMSSIPProxy API call
	client.On("DeleteJSON", mock.Anything, mock.MatchedBy(func(path string) bool {
		return path == "configuration/v1/mssip_proxy/123/"
	}), mock.Anything).Return(nil)

	testInfinityMSSIPProxy(t, client)
}

func testInfinityMSSIPProxy(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: getTestProtoV5ProviderFactories(client),
		Steps: []resource.TestStep{
			{
				Config: test.LoadTestFolder(t, "resource_infinity_mssip_proxy_basic"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_mssip_proxy.mssip-proxy-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_mssip_proxy.mssip-proxy-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_mssip_proxy.mssip-proxy-test", "name", "mssip-proxy-test"),
					resource.TestCheckResourceAttr("pexip_infinity_mssip_proxy.mssip-proxy-test", "description", "Test MSSIP proxy"),
					resource.TestCheckResourceAttr("pexip_infinity_mssip_proxy.mssip-proxy-test", "address", "test-mssip-proxy.dev.pexip.network"),
					resource.TestCheckResourceAttr("pexip_infinity_mssip_proxy.mssip-proxy-test", "port", "5060"),
					resource.TestCheckResourceAttr("pexip_infinity_mssip_proxy.mssip-proxy-test", "transport", "tcp"),
				),
			},
			{
				Config: test.LoadTestFolder(t, "resource_infinity_mssip_proxy_basic_updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_mssip_proxy.mssip-proxy-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_mssip_proxy.mssip-proxy-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_mssip_proxy.mssip-proxy-test", "name", "mssip-proxy-test"),
					resource.TestCheckResourceAttr("pexip_infinity_mssip_proxy.mssip-proxy-test", "description", "Updated Test MSSIP proxy"),
					resource.TestCheckResourceAttr("pexip_infinity_mssip_proxy.mssip-proxy-test", "address", "updated-mssip-proxy.dev.pexip.network"),
					resource.TestCheckResourceAttr("pexip_infinity_mssip_proxy.mssip-proxy-test", "port", "5061"),
					resource.TestCheckResourceAttr("pexip_infinity_mssip_proxy.mssip-proxy-test", "transport", "tls"),
				),
			},
		},
	})
}
