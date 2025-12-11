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

func TestInfinityEndUser(t *testing.T) {
	t.Parallel()
	_ = os.Setenv("TF_ACC", "1")

	// Create a mock client and set up expectations
	client := infinity.NewClientMock()

	// Mock the CreateEndUser API call
	createResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/end_user/123/",
	}
	client.On("PostWithResponse", mock.Anything, "configuration/v1/end_user/", mock.Anything, mock.Anything).Return(createResponse, nil)

	// Shared state for mocking
	mockState := &config.EndUser{
		ID:                  123,
		ResourceURI:         "/api/admin/configuration/v1/end_user/123/",
		PrimaryEmailAddress: "user@example.com",
		FirstName:           "John",
		LastName:            "Doe",
		DisplayName:         "John Doe",
		TelephoneNumber:     "+1234567890",
		MobileNumber:        "+0987654321",
		Title:               "Software Engineer",
		Department:          "Engineering",
		AvatarURL:           "https://example.com/avatar.jpg",
		UserGroups:          []string{},
		UserOID:             test.StringPtr("user-oid-123"),
		ExchangeUserID:      test.StringPtr("exchange-123"),
		MSExchangeGUID:      test.StringPtr("11111111-2222-3333-4444-555555555555"),
		SyncTag:             "sync-tag",
	}

	// Mock the GetEndUser API call for Read operations
	client.On("GetJSON", mock.Anything, "configuration/v1/end_user/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		endUser := args.Get(3).(*config.EndUser)
		*endUser = *mockState
	}).Maybe()

	// Mock the UpdateEndUser API call
	client.On("PutJSON", mock.Anything, "configuration/v1/end_user/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		updateRequest := args.Get(2).(*config.EndUserUpdateRequest)
		endUser := args.Get(3).(*config.EndUser)

		// Update mock state
		if updateRequest.PrimaryEmailAddress != "" {
			mockState.PrimaryEmailAddress = updateRequest.PrimaryEmailAddress
		}
		if updateRequest.FirstName != "" {
			mockState.FirstName = updateRequest.FirstName
		}
		if updateRequest.LastName != "" {
			mockState.LastName = updateRequest.LastName
		}
		if updateRequest.DisplayName != "" {
			mockState.DisplayName = updateRequest.DisplayName
		}
		if updateRequest.TelephoneNumber != "" {
			mockState.TelephoneNumber = updateRequest.TelephoneNumber
		}
		if updateRequest.MobileNumber != "" {
			mockState.MobileNumber = updateRequest.MobileNumber
		}
		if updateRequest.Title != "" {
			mockState.Title = updateRequest.Title
		}
		if updateRequest.Department != "" {
			mockState.Department = updateRequest.Department
		}
		if updateRequest.AvatarURL != "" {
			mockState.AvatarURL = updateRequest.AvatarURL
		}
		if updateRequest.UserOID != nil && *updateRequest.UserOID != "" {
			mockState.UserOID = updateRequest.UserOID
		}
		if updateRequest.ExchangeUserID != nil && *updateRequest.ExchangeUserID != "" {
			mockState.ExchangeUserID = updateRequest.ExchangeUserID
		}
		if updateRequest.MSExchangeGUID != nil && *updateRequest.MSExchangeGUID != "" {
			mockState.MSExchangeGUID = updateRequest.MSExchangeGUID
		}
		if updateRequest.SyncTag != "" {
			mockState.SyncTag = updateRequest.SyncTag
		}

		// Return updated state
		*endUser = *mockState
	}).Maybe()

	// Mock the DeleteEndUser API call
	client.On("DeleteJSON", mock.Anything, mock.MatchedBy(func(path string) bool {
		return path == "configuration/v1/end_user/123/"
	}), mock.Anything).Return(nil)

	testInfinityEndUser(t, client)
}

func testInfinityEndUser(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: getTestProtoV5ProviderFactories(client),
		Steps: []resource.TestStep{
			{
				Config: test.LoadTestFolder(t, "resource_infinity_end_user_basic"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_end_user.end-user-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_end_user.end-user-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_end_user.end-user-test", "primary_email_address", "user@example.com"),
					resource.TestCheckResourceAttr("pexip_infinity_end_user.end-user-test", "first_name", "John"),
					resource.TestCheckResourceAttr("pexip_infinity_end_user.end-user-test", "last_name", "Doe"),
					resource.TestCheckResourceAttr("pexip_infinity_end_user.end-user-test", "display_name", "John Doe"),
					resource.TestCheckResourceAttr("pexip_infinity_end_user.end-user-test", "telephone_number", "+1234567890"),
					resource.TestCheckResourceAttr("pexip_infinity_end_user.end-user-test", "mobile_number", "+0987654321"),
					resource.TestCheckResourceAttr("pexip_infinity_end_user.end-user-test", "title", "Software Engineer"),
					resource.TestCheckResourceAttr("pexip_infinity_end_user.end-user-test", "department", "Engineering"),
					resource.TestCheckResourceAttr("pexip_infinity_end_user.end-user-test", "avatar_url", "https://example.com/avatar.jpg"),
					resource.TestCheckResourceAttr("pexip_infinity_end_user.end-user-test", "user_oid", "user-oid-123"),
					resource.TestCheckResourceAttr("pexip_infinity_end_user.end-user-test", "exchange_user_id", "exchange-123"),
					resource.TestCheckResourceAttr("pexip_infinity_end_user.end-user-test", "ms_exchange_guid", "11111111-2222-3333-4444-555555555555"),
					resource.TestCheckResourceAttr("pexip_infinity_end_user.end-user-test", "sync_tag", "sync-tag"),
				),
			},
			{
				Config: test.LoadTestFolder(t, "resource_infinity_end_user_basic_updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_end_user.end-user-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_end_user.end-user-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_end_user.end-user-test", "primary_email_address", "updated@example.com"),
					resource.TestCheckResourceAttr("pexip_infinity_end_user.end-user-test", "first_name", "Jane"),
					resource.TestCheckResourceAttr("pexip_infinity_end_user.end-user-test", "last_name", "Smith"),
					resource.TestCheckResourceAttr("pexip_infinity_end_user.end-user-test", "display_name", "Jane Smith"),
					resource.TestCheckResourceAttr("pexip_infinity_end_user.end-user-test", "telephone_number", "+1111111111"),
					resource.TestCheckResourceAttr("pexip_infinity_end_user.end-user-test", "mobile_number", "+2222222222"),
					resource.TestCheckResourceAttr("pexip_infinity_end_user.end-user-test", "title", "Senior Engineer"),
					resource.TestCheckResourceAttr("pexip_infinity_end_user.end-user-test", "department", "Product"),
					resource.TestCheckResourceAttr("pexip_infinity_end_user.end-user-test", "avatar_url", "https://example.com/updated-avatar.jpg"),
					resource.TestCheckResourceAttr("pexip_infinity_end_user.end-user-test", "user_oid", "user-oid-123"),
					resource.TestCheckResourceAttr("pexip_infinity_end_user.end-user-test", "exchange_user_id", "exchange-123"),
					resource.TestCheckResourceAttr("pexip_infinity_end_user.end-user-test", "ms_exchange_guid", "updated-guid-456"),
					resource.TestCheckResourceAttr("pexip_infinity_end_user.end-user-test", "sync_tag", "updated-sync-tag"),
				),
			},
		},
	})
}
