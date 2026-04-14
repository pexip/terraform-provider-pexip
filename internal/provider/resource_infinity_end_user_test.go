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

func TestInfinityEndUser(t *testing.T) {
	t.Parallel()
	_ = os.Setenv("TF_ACC", "1")

	// Create a mock client and set up expectations
	client := infinity.NewClientMock()

	// Shared state for mocking
	mockState := &config.EndUser{
		ID:                  1,
		ResourceURI:         "/api/admin/configuration/v1/end_user/1/",
		PrimaryEmailAddress: "tf-test-user@example.com",
		FirstName:           "",
		LastName:            "",
		DisplayName:         "",
		TelephoneNumber:     "",
		MobileNumber:        "",
		Title:               "",
		Department:          "",
		AvatarURL:           "",
		UserGroups:          []string{},
		UserOID:             test.StringPtr("user-oid-123"),
		ExchangeUserID:      test.StringPtr("exchange-123"),
		MSExchangeGUID:      nil,
		SyncTag:             "",
	}

	// Step 1: Create with full config
	client.On("PostWithResponse", mock.Anything, "configuration/v1/end_user/", mock.Anything, mock.Anything).Return(&types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/end_user/1/",
	}, nil).Run(func(args mock.Arguments) {
		req := args.Get(2).(*config.EndUserCreateRequest)
		mockState.PrimaryEmailAddress = req.PrimaryEmailAddress
		mockState.FirstName = req.FirstName
		mockState.LastName = req.LastName
		mockState.DisplayName = req.DisplayName
		mockState.TelephoneNumber = req.TelephoneNumber
		mockState.MobileNumber = req.MobileNumber
		mockState.Title = req.Title
		mockState.Department = req.Department
		mockState.AvatarURL = req.AvatarURL
		mockState.MSExchangeGUID = req.MSExchangeGUID
		mockState.SyncTag = req.SyncTag
		mockState.UserGroups = req.UserGroups
	}).Once()

	// Step 2: Update to min config (clear all optional fields)
	client.On("PutJSON", mock.Anything, "configuration/v1/end_user/1/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		req := args.Get(2).(*config.EndUserUpdateRequest)
		mockState.PrimaryEmailAddress = req.PrimaryEmailAddress
		mockState.FirstName = req.FirstName
		mockState.LastName = req.LastName
		mockState.DisplayName = req.DisplayName
		mockState.TelephoneNumber = req.TelephoneNumber
		mockState.MobileNumber = req.MobileNumber
		mockState.Title = req.Title
		mockState.Department = req.Department
		mockState.AvatarURL = req.AvatarURL
		mockState.MSExchangeGUID = req.MSExchangeGUID
		mockState.SyncTag = req.SyncTag
		if req.UserGroups != nil {
			mockState.UserGroups = req.UserGroups
		}
		if args.Get(3) != nil {
			endUser := args.Get(3).(*config.EndUser)
			*endUser = *mockState
		}
	}).Once()

	// Step 3: Delete
	client.On("DeleteJSON", mock.Anything, "configuration/v1/end_user/1/", mock.Anything).Return(nil).Maybe()

	// Step 4: Recreate with min config
	client.On("PostWithResponse", mock.Anything, "configuration/v1/end_user/", mock.Anything, mock.Anything).Return(&types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/end_user/1/",
	}, nil).Run(func(args mock.Arguments) {
		req := args.Get(2).(*config.EndUserCreateRequest)
		mockState.PrimaryEmailAddress = req.PrimaryEmailAddress
		mockState.FirstName = req.FirstName
		mockState.LastName = req.LastName
		mockState.DisplayName = req.DisplayName
		mockState.TelephoneNumber = req.TelephoneNumber
		mockState.MobileNumber = req.MobileNumber
		mockState.Title = req.Title
		mockState.Department = req.Department
		mockState.AvatarURL = req.AvatarURL
		mockState.MSExchangeGUID = req.MSExchangeGUID
		mockState.SyncTag = req.SyncTag
		mockState.UserGroups = req.UserGroups
	}).Once()

	// Step 5: Update to full config
	client.On("PutJSON", mock.Anything, "configuration/v1/end_user/1/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		req := args.Get(2).(*config.EndUserUpdateRequest)
		mockState.PrimaryEmailAddress = req.PrimaryEmailAddress
		mockState.FirstName = req.FirstName
		mockState.LastName = req.LastName
		mockState.DisplayName = req.DisplayName
		mockState.TelephoneNumber = req.TelephoneNumber
		mockState.MobileNumber = req.MobileNumber
		mockState.Title = req.Title
		mockState.Department = req.Department
		mockState.AvatarURL = req.AvatarURL
		mockState.MSExchangeGUID = req.MSExchangeGUID
		mockState.SyncTag = req.SyncTag
		if req.UserGroups != nil {
			mockState.UserGroups = req.UserGroups
		}
		if args.Get(3) != nil {
			endUser := args.Get(3).(*config.EndUser)
			*endUser = *mockState
		}
	}).Once()

	// Mock Read operations (GetJSON) - used throughout all steps
	client.On("GetJSON", mock.Anything, "configuration/v1/end_user/1/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		endUser := args.Get(3).(*config.EndUser)
		*endUser = *mockState
	}).Maybe()

	testInfinityEndUser(t, client)
}

func testInfinityEndUser(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: getTestProtoV6ProviderFactories(client),
		Steps: []resource.TestStep{
			{
				// Step 1: Create with full config
				Config: test.LoadTestFolder(t, "resource_infinity_end_user_full"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_end_user.tf-test-end-user", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_end_user.tf-test-end-user", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_end_user.tf-test-end-user", "primary_email_address", "tf-test-user@example.com"),
					resource.TestCheckResourceAttr("pexip_infinity_end_user.tf-test-end-user", "first_name", "tf-test John"),
					resource.TestCheckResourceAttr("pexip_infinity_end_user.tf-test-end-user", "last_name", "tf-test Doe"),
					resource.TestCheckResourceAttr("pexip_infinity_end_user.tf-test-end-user", "display_name", "tf-test John Doe"),
					resource.TestCheckResourceAttr("pexip_infinity_end_user.tf-test-end-user", "telephone_number", "+1234567890"),
					resource.TestCheckResourceAttr("pexip_infinity_end_user.tf-test-end-user", "mobile_number", "+0987654321"),
					resource.TestCheckResourceAttr("pexip_infinity_end_user.tf-test-end-user", "title", "tf-test Software Engineer"),
					resource.TestCheckResourceAttr("pexip_infinity_end_user.tf-test-end-user", "department", "tf-test Engineering"),
					resource.TestCheckResourceAttr("pexip_infinity_end_user.tf-test-end-user", "avatar_url", "https://example.com/avatar.jpg"),
					resource.TestCheckResourceAttr("pexip_infinity_end_user.tf-test-end-user", "ms_exchange_guid", "11111111-2222-3333-4444-555555555555"),
					resource.TestCheckResourceAttr("pexip_infinity_end_user.tf-test-end-user", "sync_tag", "tf-test-sync-tag"),
				),
			},
			{
				// Step 2: Update to min config (clear all optional fields)
				Config: test.LoadTestFolder(t, "resource_infinity_end_user_min"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_end_user.tf-test-end-user", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_end_user.tf-test-end-user", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_end_user.tf-test-end-user", "primary_email_address", "tf-test-user@example.com"),
					resource.TestCheckResourceAttr("pexip_infinity_end_user.tf-test-end-user", "first_name", ""),
					resource.TestCheckResourceAttr("pexip_infinity_end_user.tf-test-end-user", "last_name", ""),
					resource.TestCheckResourceAttr("pexip_infinity_end_user.tf-test-end-user", "display_name", ""),
					resource.TestCheckResourceAttr("pexip_infinity_end_user.tf-test-end-user", "telephone_number", ""),
					resource.TestCheckResourceAttr("pexip_infinity_end_user.tf-test-end-user", "mobile_number", ""),
					resource.TestCheckResourceAttr("pexip_infinity_end_user.tf-test-end-user", "title", ""),
					resource.TestCheckResourceAttr("pexip_infinity_end_user.tf-test-end-user", "department", ""),
					resource.TestCheckResourceAttr("pexip_infinity_end_user.tf-test-end-user", "avatar_url", ""),
					resource.TestCheckResourceAttr("pexip_infinity_end_user.tf-test-end-user", "sync_tag", ""),
				),
			},
			{
				// Step 3: Destroy
				Config:  test.LoadTestFolder(t, "resource_infinity_end_user_min"),
				Destroy: true,
			},
			{
				// Step 4: Create with min config (after destroy)
				Config: test.LoadTestFolder(t, "resource_infinity_end_user_min"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_end_user.tf-test-end-user", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_end_user.tf-test-end-user", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_end_user.tf-test-end-user", "primary_email_address", "tf-test-user@example.com"),
					resource.TestCheckResourceAttr("pexip_infinity_end_user.tf-test-end-user", "first_name", ""),
					resource.TestCheckResourceAttr("pexip_infinity_end_user.tf-test-end-user", "last_name", ""),
					resource.TestCheckResourceAttr("pexip_infinity_end_user.tf-test-end-user", "display_name", ""),
					resource.TestCheckResourceAttr("pexip_infinity_end_user.tf-test-end-user", "telephone_number", ""),
					resource.TestCheckResourceAttr("pexip_infinity_end_user.tf-test-end-user", "mobile_number", ""),
					resource.TestCheckResourceAttr("pexip_infinity_end_user.tf-test-end-user", "title", ""),
					resource.TestCheckResourceAttr("pexip_infinity_end_user.tf-test-end-user", "department", ""),
					resource.TestCheckResourceAttr("pexip_infinity_end_user.tf-test-end-user", "avatar_url", ""),
					resource.TestCheckResourceAttr("pexip_infinity_end_user.tf-test-end-user", "sync_tag", ""),
				),
			},
			{
				// Step 5: Update to full config
				Config: test.LoadTestFolder(t, "resource_infinity_end_user_full"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_end_user.tf-test-end-user", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_end_user.tf-test-end-user", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_end_user.tf-test-end-user", "primary_email_address", "tf-test-user@example.com"),
					resource.TestCheckResourceAttr("pexip_infinity_end_user.tf-test-end-user", "first_name", "tf-test John"),
					resource.TestCheckResourceAttr("pexip_infinity_end_user.tf-test-end-user", "last_name", "tf-test Doe"),
					resource.TestCheckResourceAttr("pexip_infinity_end_user.tf-test-end-user", "display_name", "tf-test John Doe"),
					resource.TestCheckResourceAttr("pexip_infinity_end_user.tf-test-end-user", "telephone_number", "+1234567890"),
					resource.TestCheckResourceAttr("pexip_infinity_end_user.tf-test-end-user", "mobile_number", "+0987654321"),
					resource.TestCheckResourceAttr("pexip_infinity_end_user.tf-test-end-user", "title", "tf-test Software Engineer"),
					resource.TestCheckResourceAttr("pexip_infinity_end_user.tf-test-end-user", "department", "tf-test Engineering"),
					resource.TestCheckResourceAttr("pexip_infinity_end_user.tf-test-end-user", "avatar_url", "https://example.com/avatar.jpg"),
					resource.TestCheckResourceAttr("pexip_infinity_end_user.tf-test-end-user", "ms_exchange_guid", "11111111-2222-3333-4444-555555555555"),
					resource.TestCheckResourceAttr("pexip_infinity_end_user.tf-test-end-user", "sync_tag", "tf-test-sync-tag"),
				),
			},
		},
	})
}
