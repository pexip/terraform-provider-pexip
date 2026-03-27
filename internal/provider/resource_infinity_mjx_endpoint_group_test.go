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

func TestInfinityMjxEndpointGroup(t *testing.T) {
	t.Parallel()
	_ = os.Setenv("TF_ACC", "1")

	client := infinity.NewClientMock()

	mockState := &config.MjxEndpointGroup{}

	createResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/mjx_endpoint_group/123/",
	}
	client.On("PostWithResponse", mock.Anything, "configuration/v1/mjx_endpoint_group/", mock.Anything, mock.Anything).Return(createResponse, nil).Run(func(args mock.Arguments) {
		createReq := args.Get(2).(*config.MjxEndpointGroupCreateRequest)
		*mockState = config.MjxEndpointGroup{
			ID:             123,
			ResourceURI:    "/api/admin/configuration/v1/mjx_endpoint_group/123/",
			Name:           createReq.Name,
			Description:    createReq.Description,
			DisableProxy:   createReq.DisableProxy,
			SystemLocation: createReq.SystemLocation,
			MjxIntegration: createReq.MjxIntegration,
		}
	})

	client.On("GetJSON", mock.Anything, "configuration/v1/mjx_endpoint_group/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		group := args.Get(3).(*config.MjxEndpointGroup)
		*group = *mockState
	}).Maybe()

	client.On("PutJSON", mock.Anything, "configuration/v1/mjx_endpoint_group/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		updateReq := args.Get(2).(*config.MjxEndpointGroupUpdateRequest)
		mockState.Name = updateReq.Name
		mockState.Description = updateReq.Description
		mockState.DisableProxy = updateReq.DisableProxy
		mockState.SystemLocation = updateReq.SystemLocation
		mockState.MjxIntegration = updateReq.MjxIntegration
	}).Maybe()

	client.On("DeleteJSON", mock.Anything, "configuration/v1/mjx_endpoint_group/123/", mock.Anything).Return(nil)

	testInfinityMjxEndpointGroup(t, client)
}

func testInfinityMjxEndpointGroup(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: getTestProtoV5ProviderFactories(client),
		Steps: []resource.TestStep{
			// Step 1: Create with full config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_mjx_endpoint_group_full"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_endpoint_group.test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_endpoint_group.test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_endpoint_group.test", "name", "tf-test mjx-endpoint-integration-group full"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_endpoint_group.test", "description", "Test MJX endpoint integration group"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_endpoint_group.test", "system_location", "/api/admin/configuration/v1/system_location/2/"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_endpoint_group.test", "mjx_integration", "/api/admin/configuration/v1/mjx_integration/1/"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_endpoint_group.test", "disable_proxy", "true"),
				),
			},
			// Step 2: Update with min config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_mjx_endpoint_group_min"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_endpoint_group.test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_endpoint_group.test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_endpoint_group.test", "name", "tf-test mjx-endpoint-integration-group min"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_endpoint_group.test", "description", ""),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_endpoint_group.test", "system_location", "/api/admin/configuration/v1/system_location/1/"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_endpoint_group.test", "disable_proxy", "false"),
				),
			},
			// Step 3: Destroy
			{
				Config:  test.LoadTestFolder(t, "resource_infinity_mjx_endpoint_group_min"),
				Destroy: true,
			},
			// Step 4: Recreate with min config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_mjx_endpoint_group_min"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_endpoint_group.test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_endpoint_group.test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_endpoint_group.test", "name", "tf-test mjx-endpoint-integration-group min"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_endpoint_group.test", "system_location", "/api/admin/configuration/v1/system_location/1/"),
				),
			},
			// Step 5: Update to full config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_mjx_endpoint_group_full"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_endpoint_group.test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_endpoint_group.test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_endpoint_group.test", "name", "tf-test mjx-endpoint-integration-group full"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_endpoint_group.test", "description", "Test MJX endpoint integration group"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_endpoint_group.test", "system_location", "/api/admin/configuration/v1/system_location/2/"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_endpoint_group.test", "mjx_integration", "/api/admin/configuration/v1/mjx_integration/1/"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_endpoint_group.test", "disable_proxy", "true"),
				),
			},
		},
	})
}
