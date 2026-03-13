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

func TestInfinityExternalWebappHost(t *testing.T) {
	t.Parallel()
	_ = os.Setenv("TF_ACC", "1")

	// Create a mock client and set up expectations
	client := infinity.NewClientMock()

	// Shared state for mocking
	mockState := &config.ExternalWebappHost{
		ID:          1,
		ResourceURI: "/api/admin/configuration/v1/external_webapp_host/1/",
		Address:     "tf-test-webapp.example.com",
	}

	// Step 1: Create
	client.On("PostWithResponse", mock.Anything, "configuration/v1/external_webapp_host/", mock.Anything, mock.Anything).Return(&types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/external_webapp_host/1/",
	}, nil).Run(func(args mock.Arguments) {
		req := args.Get(2).(*config.ExternalWebappHostCreateRequest)
		mockState.Address = req.Address
	}).Once()

	// Step 2: Update
	client.On("PutJSON", mock.Anything, "configuration/v1/external_webapp_host/1/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		req := args.Get(2).(*config.ExternalWebappHostUpdateRequest)
		mockState.Address = req.Address
		if args.Get(3) != nil {
			host := args.Get(3).(*config.ExternalWebappHost)
			*host = *mockState
		}
	}).Once()

	// Mock Delete for cleanup
	client.On("DeleteJSON", mock.Anything, "configuration/v1/external_webapp_host/1/", mock.Anything).Return(nil).Maybe()

	// Mock Read operations (GetJSON) - used throughout all steps
	client.On("GetJSON", mock.Anything, "configuration/v1/external_webapp_host/1/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		host := args.Get(3).(*config.ExternalWebappHost)
		*host = *mockState
	}).Maybe()

	testInfinityExternalWebappHost(t, client)
}

func testInfinityExternalWebappHost(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: getTestProtoV5ProviderFactories(client),
		Steps: []resource.TestStep{
			{
				// Step 1: Create
				Config: test.LoadTestFolder(t, "resource_infinity_external_webapp_host_full"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_external_webapp_host.tf-test-external-webapp-host", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_external_webapp_host.tf-test-external-webapp-host", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_external_webapp_host.tf-test-external-webapp-host", "address", "tf-test-webapp.example.com"),
				),
			},
			{
				// Step 2: Update
				Config: test.LoadTestFolder(t, "resource_infinity_external_webapp_host_min"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_external_webapp_host.tf-test-external-webapp-host", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_external_webapp_host.tf-test-external-webapp-host", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_external_webapp_host.tf-test-external-webapp-host", "address", "tf-test-webapp-min.example.com"),
				),
			},
		},
	})
}
