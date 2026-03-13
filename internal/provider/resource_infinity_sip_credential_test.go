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
	client.On("PostWithResponse", mock.Anything, "configuration/v1/sip_credential/", mock.Anything, mock.Anything).Return(createResponse, nil).Maybe()

	// Shared state for mocking - starts with full config
	mockState := &config.SIPCredential{
		ID:          123,
		ResourceURI: "/api/admin/configuration/v1/sip_credential/123/",
		Username:    "tf-test-sip-credential",
		Realm:       "tf-test-realm",
		Password:    "tf-test-password",
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
		// Password field always sent now (without omitempty)
		mockState.Password = updateReq.Password

		// Return updated state
		*sip_credential = *mockState
	}).Maybe()

	// Mock the DeleteSipcredential API call
	client.On("DeleteJSON", mock.Anything, mock.MatchedBy(func(path string) bool {
		return path == "configuration/v1/sip_credential/123/"
	}), mock.Anything).Return(nil).Maybe()

	testInfinitySIPCredential(t, client)
}

func testInfinitySIPCredential(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: getTestProtoV5ProviderFactories(client),
		Steps: []resource.TestStep{
			// Step 1: Create with full config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_sip_credential_full"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_sip_credential.tf-test-sip-credential", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_sip_credential.tf-test-sip-credential", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_sip_credential.tf-test-sip-credential", "realm", "tf-test-realm"),
					resource.TestCheckResourceAttr("pexip_infinity_sip_credential.tf-test-sip-credential", "username", "tf-test-sip-credential"),
					resource.TestCheckResourceAttr("pexip_infinity_sip_credential.tf-test-sip-credential", "password", "tf-test-password"),
				),
			},
			// Step 2: Update to min config (clearing password)
			{
				Config: test.LoadTestFolder(t, "resource_infinity_sip_credential_min"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_sip_credential.tf-test-sip-credential", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_sip_credential.tf-test-sip-credential", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_sip_credential.tf-test-sip-credential", "realm", "tf-test-realm"),
					resource.TestCheckResourceAttr("pexip_infinity_sip_credential.tf-test-sip-credential", "username", "tf-test-sip-credential"),
					resource.TestCheckResourceAttr("pexip_infinity_sip_credential.tf-test-sip-credential", "password", ""),
				),
			},
			// Step 3: Destroy
			{
				Config:  test.LoadTestFolder(t, "resource_infinity_sip_credential_min"),
				Destroy: true,
			},
			// Step 4: Create with min config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_sip_credential_min"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_sip_credential.tf-test-sip-credential", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_sip_credential.tf-test-sip-credential", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_sip_credential.tf-test-sip-credential", "realm", "tf-test-realm"),
					resource.TestCheckResourceAttr("pexip_infinity_sip_credential.tf-test-sip-credential", "username", "tf-test-sip-credential"),
					resource.TestCheckResourceAttr("pexip_infinity_sip_credential.tf-test-sip-credential", "password", ""),
				),
			},
			// Step 5: Update to full config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_sip_credential_full"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_sip_credential.tf-test-sip-credential", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_sip_credential.tf-test-sip-credential", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_sip_credential.tf-test-sip-credential", "realm", "tf-test-realm"),
					resource.TestCheckResourceAttr("pexip_infinity_sip_credential.tf-test-sip-credential", "username", "tf-test-sip-credential"),
					resource.TestCheckResourceAttr("pexip_infinity_sip_credential.tf-test-sip-credential", "password", "tf-test-password"),
				),
			},
		},
	})
}
