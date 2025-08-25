/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/pexip/go-infinity-sdk/v38"

	"github.com/pexip/terraform-provider-pexip/internal/test"
)

func TestInfinityManagerConfig(t *testing.T) {
	t.Parallel()

	// Create a mock client and set up expectations
	client := infinity.NewClientMock()

	testInfinityManagerConfig(t, client)
}

func testInfinityManagerConfig(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: getTestProtoV5ProviderFactories(client),
		Steps: []resource.TestStep{
			{
				Config: test.LoadTestFolder(t, "data_infinity_manager_config_basic"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.pexip_infinity_manager_config.master", "hostname"),
					resource.TestCheckResourceAttrSet("data.pexip_infinity_manager_config.master", "domain"),
					resource.TestCheckResourceAttrSet("data.pexip_infinity_manager_config.master", "ip"),
					resource.TestCheckResourceAttrSet("data.pexip_infinity_manager_config.master", "mask"),
					resource.TestCheckResourceAttrSet("data.pexip_infinity_manager_config.master", "gw"),
					resource.TestCheckResourceAttrSet("data.pexip_infinity_manager_config.master", "dns"),
					resource.TestCheckResourceAttrSet("data.pexip_infinity_manager_config.master", "ntp"),
					resource.TestCheckResourceAttrSet("data.pexip_infinity_manager_config.master", "user"),
					resource.TestCheckResourceAttrSet("data.pexip_infinity_manager_config.master", "pass"),
					resource.TestCheckResourceAttrSet("data.pexip_infinity_manager_config.master", "admin_password"),
					resource.TestCheckResourceAttrSet("data.pexip_infinity_manager_config.master", "error_reports"),
					resource.TestCheckResourceAttrSet("data.pexip_infinity_manager_config.master", "enable_analytics"),
					resource.TestCheckResourceAttrSet("data.pexip_infinity_manager_config.master", "contact_email_address"),
				),
			},
		},
	})
}
