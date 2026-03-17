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

func TestInfinityLogLevel(t *testing.T) {
	t.Parallel()
	_ = os.Setenv("TF_ACC", "1")

	// Create a mock client and set up expectations
	client := infinity.NewClientMock()

	// Shared state for mocking - starts with full config
	mockState := &config.LogLevel{
		ID:          123,
		ResourceURI: "/api/admin/configuration/v1/log_level/123/",
		Name:        "tf-test-log-level-full",
		Level:       "CRITICAL",
	}

	// Mock the CreateLoglevel API call
	createResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/log_level/123/",
	}
	client.On("PostWithResponse", mock.Anything, "configuration/v1/log_level/", mock.Anything, mock.Anything).Return(createResponse, nil).Run(func(args mock.Arguments) {
		createReq := args.Get(2).(*config.LogLevelCreateRequest)
		// Update mock state based on create request
		mockState.Name = createReq.Name
		mockState.Level = createReq.Level
	}).Maybe()

	// Mock the GetLoglevel API call for Read operations
	client.On("GetJSON", mock.Anything, "configuration/v1/log_level/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		log_level := args.Get(3).(*config.LogLevel)
		*log_level = *mockState
	}).Maybe()

	// Mock the UpdateLoglevel API call
	client.On("PutJSON", mock.Anything, "configuration/v1/log_level/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		updateRequest := args.Get(2).(*config.LogLevelUpdateRequest)
		log_level := args.Get(3).(*config.LogLevel)

		// Update mock state based on request
		if updateRequest.Name != "" {
			mockState.Name = updateRequest.Name
		}
		if updateRequest.Level != "" {
			mockState.Level = updateRequest.Level
		}

		// Return updated state
		*log_level = *mockState
	}).Maybe()

	// Mock the DeleteLoglevel API call
	client.On("DeleteJSON", mock.Anything, mock.MatchedBy(func(path string) bool {
		return path == "configuration/v1/log_level/123/"
	}), mock.Anything).Return(nil)

	testInfinityLogLevel(t, client)
}

func testInfinityLogLevel(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: getTestProtoV5ProviderFactories(client),
		Steps: []resource.TestStep{
			// Step 1: Create with full config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_log_level_full"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_log_level.log_level-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_log_level.log_level-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_log_level.log_level-test", "name", "tf-test-log-level-full"),
					resource.TestCheckResourceAttr("pexip_infinity_log_level.log_level-test", "level", "CRITICAL"),
				),
			},
			// Step 2: Update to min config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_log_level_min"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_log_level.log_level-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_log_level.log_level-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_log_level.log_level-test", "name", "tf-test-log-level"),
					resource.TestCheckResourceAttr("pexip_infinity_log_level.log_level-test", "level", "NOTSET"),
				),
			},
			// Step 3: Destroy
			{
				Config:  test.LoadTestFolder(t, "resource_infinity_log_level_min"),
				Destroy: true,
			},
			// Step 4: Create with min config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_log_level_min"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_log_level.log_level-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_log_level.log_level-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_log_level.log_level-test", "name", "tf-test-log-level"),
					resource.TestCheckResourceAttr("pexip_infinity_log_level.log_level-test", "level", "NOTSET"),
				),
			},
			// Step 5: Update to full config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_log_level_full"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_log_level.log_level-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_log_level.log_level-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_log_level.log_level-test", "name", "tf-test-log-level-full"),
					resource.TestCheckResourceAttr("pexip_infinity_log_level.log_level-test", "level", "CRITICAL"),
				),
			},
		},
	})
}
