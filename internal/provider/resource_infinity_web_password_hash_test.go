/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/pexip/go-infinity-sdk/v38"

	"github.com/pexip/terraform-provider-pexip/internal/test"
)

func TestInfinityWebPasswordHash(t *testing.T) {
	t.Parallel()
	_ = os.Setenv("TF_ACC", "1")

	// Create a mock client and set up expectations
	client := infinity.NewClientMock()

	testInfinityWebPasswordHash(t, client)
}

func testInfinityWebPasswordHash(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: getTestProtoV5ProviderFactories(client),
		Steps: []resource.TestStep{
			{
				Config: test.LoadTestFolder(t, "resource_infinity_web_password_hash_basic"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_web_password_hash.web_password_hash-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_web_password_hash.web_password_hash-test", "hash"),
					resource.TestCheckResourceAttr("pexip_infinity_web_password_hash.web_password_hash-test", "rounds", "5000"),
					resource.TestCheckResourceAttr("pexip_infinity_web_password_hash.web_password_hash-test", "salt", "abcdefghijkl"),
				),
			},
			{
				Config: test.LoadTestFolder(t, "resource_infinity_web_password_hash_basic_updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_web_password_hash.web_password_hash-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_web_password_hash.web_password_hash-test", "hash"),
					resource.TestCheckResourceAttr("pexip_infinity_web_password_hash.web_password_hash-test", "rounds", "6000"),
					resource.TestCheckResourceAttr("pexip_infinity_web_password_hash.web_password_hash-test", "salt", "mnopqrstuvwx"),
				),
			},
		},
	})
}
