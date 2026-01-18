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

	// Mock the CreateLoglevel API call
	createResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/log_level/123/",
	}
	client.On("PostWithResponse", mock.Anything, "configuration/v1/log_level/", mock.Anything, mock.Anything).Return(createResponse, nil)

	// Shared state for mocking
	mockState := &config.LogLevel{
		ID:          123,
		ResourceURI: "/api/admin/configuration/v1/log_level/123/",
		Name:        "log_level-test",
		Level:       "INFO",
	}

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
			{
				Config: test.LoadTestFolder(t, "resource_infinity_log_level_basic"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_log_level.log_level-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_log_level.log_level-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_log_level.log_level-test", "name", "log_level-test"),
					resource.TestCheckResourceAttr("pexip_infinity_log_level.log_level-test", "level", "INFO"),
				),
			},
			{
				Config: test.LoadTestFolder(t, "resource_infinity_log_level_basic_updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_log_level.log_level-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_log_level.log_level-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_log_level.log_level-test", "name", "log_level-test"),
					resource.TestCheckResourceAttr("pexip_infinity_log_level.log_level-test", "level", "ERROR"),
				),
			},
		},
	})
}
