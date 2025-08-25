/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package provider

import (
	"testing"
	"time"

	"github.com/pexip/go-infinity-sdk/v38/config"
	"github.com/pexip/go-infinity-sdk/v38/types"
	"github.com/pexip/go-infinity-sdk/v38/util"
	"github.com/stretchr/testify/mock"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/pexip/go-infinity-sdk/v38"

	"github.com/pexip/terraform-provider-pexip/internal/test"
)

func TestInfinitySystemSyncpoint(t *testing.T) {
	t.Parallel()

	// Create a mock client and set up expectations
	client := infinity.NewClientMock()

	// Mock the CreateSystemsyncpoint API call
	createResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/system_syncpoint/1/",
	}
	client.On("PostWithResponse", mock.Anything, "configuration/v1/system_syncpoint/", mock.Anything, mock.Anything).Return(createResponse, nil)

	// Shared state for mocking
	mockState := &config.SystemSyncpoint{
		ID:           1,
		ResourceURI:  "/api/admin/configuration/v1/system_syncpoint/1/",
		CreationTime: util.InfinityTime{Time: time.Now()},
	}

	// Mock the GetSystemsyncpoint API call for Read operations
	client.On("GetJSON", mock.Anything, "configuration/v1/system_syncpoint/1/", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		system_syncpoint := args.Get(2).(*config.SystemSyncpoint)
		*system_syncpoint = *mockState
	}).Maybe()

	testInfinitySystemSyncpoint(t, client)
}

func testInfinitySystemSyncpoint(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: getTestProtoV5ProviderFactories(client),
		Steps: []resource.TestStep{
			{
				Config: test.LoadTestFolder(t, "resource_infinity_system_syncpoint_basic"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_system_syncpoint.system_syncpoint-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_system_syncpoint.system_syncpoint-test", "resource_id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_system_syncpoint.system_syncpoint-test", "creation_time"),
				),
				// Use PlanOnly to avoid attempting to destroy the resource after the test
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}
