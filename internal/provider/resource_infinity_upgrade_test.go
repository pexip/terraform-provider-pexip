/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package provider

import (
	"testing"

	"github.com/pexip/go-infinity-sdk/v38/types"
	"github.com/stretchr/testify/mock"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/pexip/go-infinity-sdk/v38"

	"github.com/pexip/terraform-provider-pexip/internal/test"
)

func TestInfinityUpgrade(t *testing.T) {
	t.Parallel()

	// Create a mock client and set up expectations
	client := infinity.NewClientMock()

	// Mock the CreateUpgrade API call
	createResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/upgrade/123/",
	}
	client.On("PostWithResponse", mock.Anything, "configuration/v1/upgrade/", mock.Anything, mock.Anything).Return(createResponse, nil)

	// Note: Upgrade is an action resource that doesn't have persistent state to read

	testInfinityUpgrade(t, client)
}

func testInfinityUpgrade(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: getTestProtoV5ProviderFactories(client),
		Steps: []resource.TestStep{
			{
				Config: test.LoadTestFolder(t, "resource_infinity_upgrade_basic"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_upgrade.upgrade-test", "id"),
					resource.TestCheckResourceAttr("pexip_infinity_upgrade.upgrade-test", "package", "test-value"),
					resource.TestCheckResourceAttrSet("pexip_infinity_upgrade.upgrade-test", "timestamp"),
				),
			},
		},
	})
}
