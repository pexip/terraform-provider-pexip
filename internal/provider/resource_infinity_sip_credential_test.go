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

func TestInfinitySIPCredential(t *testing.T) {
	t.Parallel()
	_ = os.Setenv("TF_ACC", "1")

	// Create a mock client and set up expectations
	client := infinity.NewClientMock()

	// Mock the CreateSipcredential API call
	createResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/sip_credential/123/",
	}
	client.On("PostWithResponse", mock.Anything, "configuration/v1/sip_credential/", mock.Anything, mock.Anything).Return(createResponse, nil)

	// Shared state for mocking
	mockState := &config.SIPCredential{
		ID:          123,
		ResourceURI: "/api/admin/configuration/v1/sip_credential/123/",
		Username:    "sip_credential-test",
		Realm:       "test-value",
		Password:    "test-value",
	}

	// Mock the GetSipcredential API call for Read operations
	client.On("GetJSON", mock.Anything, "configuration/v1/sip_credential/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		sip_credential := args.Get(3).(*config.SIPCredential)
		*sip_credential = *mockState
	}).Maybe()

	// Mock the UpdateSipcredential API call
	client.On("PutJSON", mock.Anything, "configuration/v1/sip_credential/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		updateReq := args.Get(2).(*config.SIPCredentialUpdateRequest)
		sip_credential := args.Get(3).(*config.SIPCredential)

		// Update mock state based on request
		if updateReq.Realm != "" {
			mockState.Realm = updateReq.Realm
		}
		if updateReq.Username != "" {
			mockState.Username = updateReq.Username
		}
		if updateReq.Password != "" {
			mockState.Password = updateReq.Password
		}

		// Return updated state
		*sip_credential = *mockState
	}).Maybe()

	// Mock the DeleteSipcredential API call
	client.On("DeleteJSON", mock.Anything, mock.MatchedBy(func(path string) bool {
		return path == "configuration/v1/sip_credential/123/"
	}), mock.Anything).Return(nil)

	testInfinitySIPCredential(t, client)
}

func testInfinitySIPCredential(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: getTestProtoV5ProviderFactories(client),
		Steps: []resource.TestStep{
			{
				Config: test.LoadTestFolder(t, "resource_infinity_sip_credential_basic"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_sip_credential.sip_credential-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_sip_credential.sip_credential-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_sip_credential.sip_credential-test", "username", "sip_credential-test"),
				),
			},
			{
				Config: test.LoadTestFolder(t, "resource_infinity_sip_credential_basic_updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_sip_credential.sip_credential-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_sip_credential.sip_credential-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_sip_credential.sip_credential-test", "username", "sip_credential-test"),
					resource.TestCheckResourceAttr("pexip_infinity_sip_credential.sip_credential-test", "realm", "updated-value"),
					resource.TestCheckResourceAttr("pexip_infinity_sip_credential.sip_credential-test", "password", "updated-value"),
				),
			},
		},
	})
}
