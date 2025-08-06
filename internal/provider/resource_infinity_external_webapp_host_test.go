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

func TestInfinityExternalWebappHost(t *testing.T) {
	_ = os.Setenv("TF_ACC", "1")

	// Create a mock client and set up expectations
	client := infinity.NewClientMock()

	// Mock the CreateExternalwebapphost API call
	createResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/external_webapp_host/123/",
	}
	client.On("PostWithResponse", mock.Anything, "configuration/v1/external_webapp_host/", mock.Anything, mock.Anything).Return(createResponse, nil)

	// Shared state for mocking
	mockState := &config.ExternalWebappHost{
		ID:          123,
		ResourceURI: "/api/admin/configuration/v1/external_webapp_host/123/",
		Address:     "test-server.example.com",
	}

	// Mock the GetExternalwebapphost API call for Read operations
	client.On("GetJSON", mock.Anything, "configuration/v1/external_webapp_host/123/", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		external_webapp_host := args.Get(2).(*config.ExternalWebappHost)
		*external_webapp_host = *mockState
	}).Maybe()

	// Mock the UpdateExternalwebapphost API call
	client.On("PutJSON", mock.Anything, "configuration/v1/external_webapp_host/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		updateRequest := args.Get(2).(*config.ExternalWebappHostUpdateRequest)
		external_webapp_host := args.Get(3).(*config.ExternalWebappHost)

		// Update mock state based on request
		mockState.Address = updateRequest.Address

		// Return updated state
		*external_webapp_host = *mockState
	}).Maybe()

	// Mock the DeleteExternalwebapphost API call
	client.On("DeleteJSON", mock.Anything, mock.MatchedBy(func(path string) bool {
		return path == "configuration/v1/external_webapp_host/123/"
	}), mock.Anything).Return(nil)

	testInfinityExternalWebappHost(t, client)
}

func testInfinityExternalWebappHost(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: getTestProtoV5ProviderFactories(client),
		Steps: []resource.TestStep{
			{
				Config: test.LoadTestFolder(t, "resource_infinity_external_webapp_host_basic"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_external_webapp_host.external_webapp_host-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_external_webapp_host.external_webapp_host-test", "resource_id"),
				),
			},
			{
				Config: test.LoadTestFolder(t, "resource_infinity_external_webapp_host_basic_updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_external_webapp_host.external_webapp_host-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_external_webapp_host.external_webapp_host-test", "resource_id"),
				),
			},
		},
	})
}
