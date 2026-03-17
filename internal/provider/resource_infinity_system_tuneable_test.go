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

func TestInfinitySystemTuneable(t *testing.T) {
	t.Parallel()
	_ = os.Setenv("TF_ACC", "1")

	// Create a mock client and set up expectations
	client := infinity.NewClientMock()

	// Shared state for mocking - starts with full config
	mockState := &config.SystemTuneable{
		ID:          123,
		ResourceURI: "/api/admin/configuration/v1/system_tuneable/123/",
		Name:        "tf-test-system-tuneable-full",
		Setting:     "full-test-value",
	}

	// Mock the CreateSystemtuneable API call
	createResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/system_tuneable/123/",
	}
	client.On("PostWithResponse", mock.Anything, "configuration/v1/system_tuneable/", mock.Anything, mock.Anything).Return(createResponse, nil).Run(func(args mock.Arguments) {
		createReq := args.Get(2).(*config.SystemTuneableCreateRequest)
		// Update mock state based on create request
		mockState.Name = createReq.Name
		mockState.Setting = createReq.Setting
	}).Maybe()

	// Mock the GetSystemtuneable API call for Read operations
	client.On("GetJSON", mock.Anything, "configuration/v1/system_tuneable/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		system_tuneable := args.Get(3).(*config.SystemTuneable)
		*system_tuneable = *mockState
	}).Maybe()

	// Mock the UpdateSystemtuneable API call
	client.On("PutJSON", mock.Anything, "configuration/v1/system_tuneable/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		updateReq := args.Get(2).(*config.SystemTuneableUpdateRequest)
		system_tuneable := args.Get(3).(*config.SystemTuneable)

		// Update mock state based on request
		if updateReq.Name != "" {
			mockState.Name = updateReq.Name
		}
		if updateReq.Setting != "" {
			mockState.Setting = updateReq.Setting
		}

		// Return updated state
		*system_tuneable = *mockState
	}).Maybe()

	// Mock the DeleteSystemtuneable API call
	client.On("DeleteJSON", mock.Anything, mock.MatchedBy(func(path string) bool {
		return path == "configuration/v1/system_tuneable/123/"
	}), mock.Anything).Return(nil)

	testInfinitySystemTuneable(t, client)
}

func testInfinitySystemTuneable(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: getTestProtoV5ProviderFactories(client),
		Steps: []resource.TestStep{
			// Step 1: Create with full config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_system_tuneable_full"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_system_tuneable.system_tuneable-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_system_tuneable.system_tuneable-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_system_tuneable.system_tuneable-test", "name", "tf-test-system-tuneable-full"),
					resource.TestCheckResourceAttr("pexip_infinity_system_tuneable.system_tuneable-test", "setting", "full-test-value"),
				),
			},
			// Step 2: Update to min config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_system_tuneable_min"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_system_tuneable.system_tuneable-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_system_tuneable.system_tuneable-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_system_tuneable.system_tuneable-test", "name", "tf-test-system-tuneable"),
					resource.TestCheckResourceAttr("pexip_infinity_system_tuneable.system_tuneable-test", "setting", "test-value"),
				),
			},
			// Step 3: Destroy
			{
				Config:  test.LoadTestFolder(t, "resource_infinity_system_tuneable_min"),
				Destroy: true,
			},
			// Step 4: Create with min config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_system_tuneable_min"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_system_tuneable.system_tuneable-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_system_tuneable.system_tuneable-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_system_tuneable.system_tuneable-test", "name", "tf-test-system-tuneable"),
					resource.TestCheckResourceAttr("pexip_infinity_system_tuneable.system_tuneable-test", "setting", "test-value"),
				),
			},
			// Step 5: Update to full config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_system_tuneable_full"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_system_tuneable.system_tuneable-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_system_tuneable.system_tuneable-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_system_tuneable.system_tuneable-test", "name", "tf-test-system-tuneable-full"),
					resource.TestCheckResourceAttr("pexip_infinity_system_tuneable.system_tuneable-test", "setting", "full-test-value"),
				),
			},
		},
	})
}
