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

func TestInfinityH323Gatekeeper(t *testing.T) {
	t.Parallel()
	_ = os.Setenv("TF_ACC", "1")

	// Create a mock client and set up expectations
	client := infinity.NewClientMock()

	// Mock the CreateH323Gatekeeper API call
	createResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/h323_gatekeeper/123/",
	}
	client.On("PostWithResponse", mock.Anything, "configuration/v1/h323_gatekeeper/", mock.Anything, mock.Anything).Return(createResponse, nil)

	// Shared state for mocking
	mockState := &config.H323Gatekeeper{
		ID:          123,
		ResourceURI: "/api/admin/configuration/v1/h323_gatekeeper/123/",
		Name:        "h323_gatekeeper-test",
		Description: "Test H323Gatekeeper",
		Address:     "192.168.1.100",
	}

	// Mock the GetH323Gatekeeper API call for Read operations
	client.On("GetJSON", mock.Anything, "configuration/v1/h323_gatekeeper/123/", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		h323_gatekeeper := args.Get(2).(*config.H323Gatekeeper)
		*h323_gatekeeper = *mockState
	}).Maybe()

	// Mock the UpdateH323Gatekeeper API call
	client.On("PutJSON", mock.Anything, "configuration/v1/h323_gatekeeper/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		updateRequest := args.Get(2).(*config.H323GatekeeperUpdateRequest)
		h323_gatekeeper := args.Get(3).(*config.H323Gatekeeper)

		// Update mock state based on request
		if updateRequest.Name != "" {
			mockState.Name = updateRequest.Name
		}
		if updateRequest.Description != "" {
			mockState.Description = updateRequest.Description
		}
		if updateRequest.Address != "" {
			mockState.Address = updateRequest.Address
		}

		// Return updated state
		*h323_gatekeeper = *mockState
	}).Maybe()

	// Mock the DeleteH323Gatekeeper API call
	client.On("DeleteJSON", mock.Anything, mock.MatchedBy(func(path string) bool {
		return path == "configuration/v1/h323_gatekeeper/123/"
	}), mock.Anything).Return(nil)

	testInfinityH323Gatekeeper(t, client)
}

func testInfinityH323Gatekeeper(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: getTestProtoV5ProviderFactories(client),
		Steps: []resource.TestStep{
			{
				Config: test.LoadTestFolder(t, "resource_infinity_h323_gatekeeper_basic"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_h323_gatekeeper.h323_gatekeeper-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_h323_gatekeeper.h323_gatekeeper-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_h323_gatekeeper.h323_gatekeeper-test", "name", "h323_gatekeeper-test"),
					resource.TestCheckResourceAttr("pexip_infinity_h323_gatekeeper.h323_gatekeeper-test", "description", "Test H323Gatekeeper"),
				),
			},
			{
				Config: test.LoadTestFolder(t, "resource_infinity_h323_gatekeeper_basic_updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_h323_gatekeeper.h323_gatekeeper-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_h323_gatekeeper.h323_gatekeeper-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_h323_gatekeeper.h323_gatekeeper-test", "name", "h323_gatekeeper-test"),
					resource.TestCheckResourceAttr("pexip_infinity_h323_gatekeeper.h323_gatekeeper-test", "description", "Updated Test H323Gatekeeper"),
				),
			},
		},
	})
}
