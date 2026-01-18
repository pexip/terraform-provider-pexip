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
		Name:                   "break_in_allow_list_address-test",
		Description:            "Test BreakInAllowListAddress",
		Address:                "192.168.1.0",
		Prefix:                 24,
		AllowlistEntryType:     "user",
		IgnoreIncorrectAliases: true,
		IgnoreIncorrectPins:    true,
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
		mockState.Name = updateRequest.Name
		if updateRequest.Description != "" {
			mockState.Description = updateRequest.Description
		}
		if updateRequest.Address != "" {
			mockState.Address = updateRequest.Address
		}
		if updateRequest.AllowlistEntryType != "" {
			mockState.AllowlistEntryType = updateRequest.AllowlistEntryType
		}
		if updateRequest.Prefix != nil {
			mockState.Prefix = *updateRequest.Prefix
		}
		if updateRequest.IgnoreIncorrectAliases != nil {
			mockState.IgnoreIncorrectAliases = *updateRequest.IgnoreIncorrectAliases
		}
		if updateRequest.IgnoreIncorrectPins != nil {
			mockState.IgnoreIncorrectPins = *updateRequest.IgnoreIncorrectPins
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
			{
				Config: test.LoadTestFolder(t, "resource_infinity_break_in_allow_list_address_basic"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_break_in_allow_list_address.break_in_allow_list_address-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_break_in_allow_list_address.break_in_allow_list_address-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_break_in_allow_list_address.break_in_allow_list_address-test", "name", "break_in_allow_list_address-test"),
					resource.TestCheckResourceAttr("pexip_infinity_break_in_allow_list_address.break_in_allow_list_address-test", "description", "Test BreakInAllowListAddress"),
					resource.TestCheckResourceAttr("pexip_infinity_break_in_allow_list_address.break_in_allow_list_address-test", "address", "192.168.1.0"),
					resource.TestCheckResourceAttr("pexip_infinity_break_in_allow_list_address.break_in_allow_list_address-test", "prefix", "24"),
					resource.TestCheckResourceAttr("pexip_infinity_break_in_allow_list_address.break_in_allow_list_address-test", "allowlist_entry_type", "user"),
					resource.TestCheckResourceAttr("pexip_infinity_break_in_allow_list_address.break_in_allow_list_address-test", "ignore_incorrect_aliases", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_break_in_allow_list_address.break_in_allow_list_address-test", "ignore_incorrect_pins", "true"),
				),
			},
			{
				Config: test.LoadTestFolder(t, "resource_infinity_break_in_allow_list_address_basic_updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_break_in_allow_list_address.break_in_allow_list_address-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_break_in_allow_list_address.break_in_allow_list_address-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_break_in_allow_list_address.break_in_allow_list_address-test", "name", "break_in_allow_list_address-test"),
					resource.TestCheckResourceAttr("pexip_infinity_break_in_allow_list_address.break_in_allow_list_address-test", "description", "Updated Test BreakInAllowListAddress"),
					resource.TestCheckResourceAttr("pexip_infinity_break_in_allow_list_address.break_in_allow_list_address-test", "address", "10.0.0.0"),
					resource.TestCheckResourceAttr("pexip_infinity_break_in_allow_list_address.break_in_allow_list_address-test", "prefix", "16"),
					resource.TestCheckResourceAttr("pexip_infinity_break_in_allow_list_address.break_in_allow_list_address-test", "allowlist_entry_type", "proxy"),
					resource.TestCheckResourceAttr("pexip_infinity_break_in_allow_list_address.break_in_allow_list_address-test", "ignore_incorrect_aliases", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_break_in_allow_list_address.break_in_allow_list_address-test", "ignore_incorrect_pins", "false"),
				),
			},
		},
	})
}
