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

func TestInfinityBreakInAllowListAddress(t *testing.T) {
	t.Parallel()
	_ = os.Setenv("TF_ACC", "1")

	// Create a mock client and set up expectations
	client := infinity.NewClientMock()

	// Mock the CreateBreakInAllowListAddress API call
	createResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/break_in_allow_list_address/123/",
	}
	client.On("PostWithResponse", mock.Anything, "configuration/v1/break_in_allow_list_address/", mock.Anything, mock.Anything).Return(createResponse, nil)

	// Shared state for mocking
	mockState := &config.BreakInAllowListAddress{
		ID:                     123,
		ResourceURI:            "/api/admin/configuration/v1/break_in_allow_list_address/123/",
		Name:                   "tf-test-break-in-allow-list-address",
		Description:            "Full test configuration for break-in allow list address",
		Address:                "10.0.0.0",
		Prefix:                 16,
		AllowlistEntryType:     "proxy",
		IgnoreIncorrectAliases: false,
		IgnoreIncorrectPins:    false,
	}

	// Mock the GetBreakInAllowListAddress API call for Read operations
	client.On("GetJSON", mock.Anything, "configuration/v1/break_in_allow_list_address/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		breakInAllowListAddress := args.Get(3).(*config.BreakInAllowListAddress)
		*breakInAllowListAddress = *mockState
	}).Maybe()

	// Mock the UpdateBreakInAllowListAddress API call
	client.On("PutJSON", mock.Anything, "configuration/v1/break_in_allow_list_address/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		updateRequest := args.Get(2).(*config.BreakInAllowListAddressUpdateRequest)
		breakInAllowListAddress := args.Get(3).(*config.BreakInAllowListAddress)

		// Update mock state
		if updateRequest.Name != "" {
			mockState.Name = updateRequest.Name
		}
		// Update description only if provided
		if updateRequest.Description != "" {
			mockState.Description = updateRequest.Description
		}
		if updateRequest.Address != "" {
			mockState.Address = updateRequest.Address
		}
		// Apply default if AllowlistEntryType not provided
		if updateRequest.AllowlistEntryType != "" {
			mockState.AllowlistEntryType = updateRequest.AllowlistEntryType
		} else {
			mockState.AllowlistEntryType = "user"
		}
		if updateRequest.Prefix != nil {
			mockState.Prefix = *updateRequest.Prefix
		}
		// Apply default if IgnoreIncorrectAliases not provided
		if updateRequest.IgnoreIncorrectAliases != nil {
			mockState.IgnoreIncorrectAliases = *updateRequest.IgnoreIncorrectAliases
		} else {
			mockState.IgnoreIncorrectAliases = false
		}
		// Apply default if IgnoreIncorrectPins not provided
		if updateRequest.IgnoreIncorrectPins != nil {
			mockState.IgnoreIncorrectPins = *updateRequest.IgnoreIncorrectPins
		} else {
			mockState.IgnoreIncorrectPins = false
		}

		// Return updated state
		*breakInAllowListAddress = *mockState
	}).Maybe()

	// Mock the DeleteBreakInAllowListAddress API call
	client.On("DeleteJSON", mock.Anything, mock.MatchedBy(func(path string) bool {
		return path == "configuration/v1/break_in_allow_list_address/123/"
	}), mock.Anything).Return(nil)

	testInfinityBreakInAllowListAddress(t, client)
}

func testInfinityBreakInAllowListAddress(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: getTestProtoV5ProviderFactories(client),
		Steps: []resource.TestStep{
			// Test 1: Create with full configuration
			{
				Config: test.LoadTestFolder(t, "resource_infinity_break_in_allow_list_address_full"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_break_in_allow_list_address.tf-test-break-in-allow-list-address", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_break_in_allow_list_address.tf-test-break-in-allow-list-address", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_break_in_allow_list_address.tf-test-break-in-allow-list-address", "name", "tf-test-break-in-allow-list-address"),
					resource.TestCheckResourceAttr("pexip_infinity_break_in_allow_list_address.tf-test-break-in-allow-list-address", "description", "Full test configuration for break-in allow list address"),
					resource.TestCheckResourceAttr("pexip_infinity_break_in_allow_list_address.tf-test-break-in-allow-list-address", "address", "10.0.0.0"),
					resource.TestCheckResourceAttr("pexip_infinity_break_in_allow_list_address.tf-test-break-in-allow-list-address", "prefix", "16"),
					resource.TestCheckResourceAttr("pexip_infinity_break_in_allow_list_address.tf-test-break-in-allow-list-address", "allowlist_entry_type", "proxy"),
					resource.TestCheckResourceAttr("pexip_infinity_break_in_allow_list_address.tf-test-break-in-allow-list-address", "ignore_incorrect_aliases", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_break_in_allow_list_address.tf-test-break-in-allow-list-address", "ignore_incorrect_pins", "false"),
				),
			},
			// Test 2: Update to min configuration, then delete
			{
				Config: test.LoadTestFolder(t, "resource_infinity_break_in_allow_list_address_min"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_break_in_allow_list_address.tf-test-break-in-allow-list-address", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_break_in_allow_list_address.tf-test-break-in-allow-list-address", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_break_in_allow_list_address.tf-test-break-in-allow-list-address", "name", "tf-test-break-in-allow-list-address"),
					resource.TestCheckResourceAttr("pexip_infinity_break_in_allow_list_address.tf-test-break-in-allow-list-address", "description", "Full test configuration for break-in allow list address"),
					resource.TestCheckResourceAttr("pexip_infinity_break_in_allow_list_address.tf-test-break-in-allow-list-address", "address", "192.168.1.0"),
					resource.TestCheckResourceAttr("pexip_infinity_break_in_allow_list_address.tf-test-break-in-allow-list-address", "prefix", "24"),
					resource.TestCheckResourceAttr("pexip_infinity_break_in_allow_list_address.tf-test-break-in-allow-list-address", "allowlist_entry_type", "user"),
					resource.TestCheckResourceAttr("pexip_infinity_break_in_allow_list_address.tf-test-break-in-allow-list-address", "ignore_incorrect_aliases", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_break_in_allow_list_address.tf-test-break-in-allow-list-address", "ignore_incorrect_pins", "false"),
				),
			},
			// Test 3: Create with min configuration
			{
				Config: test.LoadTestFolder(t, "resource_infinity_break_in_allow_list_address_min"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_break_in_allow_list_address.tf-test-break-in-allow-list-address", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_break_in_allow_list_address.tf-test-break-in-allow-list-address", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_break_in_allow_list_address.tf-test-break-in-allow-list-address", "name", "tf-test-break-in-allow-list-address"),
					resource.TestCheckResourceAttr("pexip_infinity_break_in_allow_list_address.tf-test-break-in-allow-list-address", "address", "192.168.1.0"),
					resource.TestCheckResourceAttr("pexip_infinity_break_in_allow_list_address.tf-test-break-in-allow-list-address", "prefix", "24"),
					resource.TestCheckResourceAttr("pexip_infinity_break_in_allow_list_address.tf-test-break-in-allow-list-address", "allowlist_entry_type", "user"),
					resource.TestCheckResourceAttr("pexip_infinity_break_in_allow_list_address.tf-test-break-in-allow-list-address", "ignore_incorrect_aliases", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_break_in_allow_list_address.tf-test-break-in-allow-list-address", "ignore_incorrect_pins", "false"),
				),
			},
			// Test 4: Update to full configuration
			{
				Config: test.LoadTestFolder(t, "resource_infinity_break_in_allow_list_address_full"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_break_in_allow_list_address.tf-test-break-in-allow-list-address", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_break_in_allow_list_address.tf-test-break-in-allow-list-address", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_break_in_allow_list_address.tf-test-break-in-allow-list-address", "name", "tf-test-break-in-allow-list-address"),
					resource.TestCheckResourceAttr("pexip_infinity_break_in_allow_list_address.tf-test-break-in-allow-list-address", "description", "Full test configuration for break-in allow list address"),
					resource.TestCheckResourceAttr("pexip_infinity_break_in_allow_list_address.tf-test-break-in-allow-list-address", "address", "10.0.0.0"),
					resource.TestCheckResourceAttr("pexip_infinity_break_in_allow_list_address.tf-test-break-in-allow-list-address", "prefix", "16"),
					resource.TestCheckResourceAttr("pexip_infinity_break_in_allow_list_address.tf-test-break-in-allow-list-address", "allowlist_entry_type", "proxy"),
					resource.TestCheckResourceAttr("pexip_infinity_break_in_allow_list_address.tf-test-break-in-allow-list-address", "ignore_incorrect_aliases", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_break_in_allow_list_address.tf-test-break-in-allow-list-address", "ignore_incorrect_pins", "false"),
				),
			},
		},
	})
}
