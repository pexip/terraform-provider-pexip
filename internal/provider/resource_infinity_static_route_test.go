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

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/pexip/go-infinity-sdk/v38"

	"github.com/pexip/terraform-provider-pexip/internal/test"
)

func TestInfinityStaticRoute(t *testing.T) {
	t.Parallel()
	_ = os.Setenv("TF_ACC", "1")

	// Create a mock client and set up expectations
	client := infinity.NewClientMock()

	// Mock the CreateStaticroute API call
	createResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/static_route/123/",
	}
	client.On("PostWithResponse", mock.Anything, "configuration/v1/static_route/", mock.Anything, mock.Anything).Return(createResponse, nil)

	// Shared state for mocking
	mockState := &config.StaticRoute{
		ID:          123,
		ResourceURI: "/api/admin/configuration/v1/static_route/123/",
		Name:        "static_route-test",
		Address:     "192.168.1.0",
		Prefix:      24,
		Gateway:     "192.168.1.1",
	}

	// Mock the GetStaticroute API call for Read operations
	client.On("GetJSON", mock.Anything, "configuration/v1/static_route/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		static_route := args.Get(3).(*config.StaticRoute)
		*static_route = *mockState
	}).Maybe()

	// Mock the UpdateStaticroute API call
	client.On("PutJSON", mock.Anything, "configuration/v1/static_route/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		updateReq := args.Get(2).(*config.StaticRouteUpdateRequest)
		static_route := args.Get(3).(*config.StaticRoute)

		// Update mock state based on request
		if updateReq.Address != "" {
			mockState.Address = updateReq.Address
		}
		if updateReq.Prefix != 0 {
			mockState.Prefix = updateReq.Prefix
		}
		if updateReq.Gateway != "" {
			mockState.Gateway = updateReq.Gateway
		}

		// Return updated state
		*static_route = *mockState
	}).Maybe()

	// Mock the DeleteStaticroute API call
	client.On("DeleteJSON", mock.Anything, mock.MatchedBy(func(path string) bool {
		return path == "configuration/v1/static_route/123/"
	}), mock.Anything).Return(nil)

	testInfinityStaticRoute(t, client)
}

func testInfinityStaticRoute(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: getTestProtoV5ProviderFactories(client),
		Steps: []resource.TestStep{
			{
				Config: test.LoadTestFolder(t, "resource_infinity_static_route_basic"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_static_route.static_route-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_static_route.static_route-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_static_route.static_route-test", "name", "static_route-test"),
				),
			},
			{
				Config: test.LoadTestFolder(t, "resource_infinity_static_route_basic_updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_static_route.static_route-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_static_route.static_route-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_static_route.static_route-test", "name", "static_route-test"),
				),
			},
		},
	})
}
