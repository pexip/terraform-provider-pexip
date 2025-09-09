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

func TestInfinitySystemTuneable(t *testing.T) {
	t.Parallel()
	_ = os.Setenv("TF_ACC", "1")

	// Create a mock client and set up expectations
	client := infinity.NewClientMock()

	// Mock the CreateSystemtuneable API call
	createResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/system_tuneable/123/",
	}
	client.On("PostWithResponse", mock.Anything, "configuration/v1/system_tuneable/", mock.Anything, mock.Anything).Return(createResponse, nil)

	// Shared state for mocking
	mockState := &config.SystemTuneable{
		ID:          123,
		ResourceURI: "/api/admin/configuration/v1/system_tuneable/123/",
		Name:        "system_tuneable-test",
		Setting:     "test-value",
	}

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
			{
				Config: test.LoadTestFolder(t, "resource_infinity_system_tuneable_basic"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_system_tuneable.system_tuneable-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_system_tuneable.system_tuneable-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_system_tuneable.system_tuneable-test", "name", "system_tuneable-test"),
					resource.TestCheckResourceAttr("pexip_infinity_system_tuneable.system_tuneable-test", "setting", "test-value"),
				),
			},
			{
				Config: test.LoadTestFolder(t, "resource_infinity_system_tuneable_basic_updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_system_tuneable.system_tuneable-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_system_tuneable.system_tuneable-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_system_tuneable.system_tuneable-test", "name", "system_tuneable-test"),
					resource.TestCheckResourceAttr("pexip_infinity_system_tuneable.system_tuneable-test", "setting", "updated-value"),
				),
			},
		},
	})
}
