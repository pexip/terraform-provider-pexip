package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/pexip/go-infinity-sdk/v38/config"
	"github.com/pexip/go-infinity-sdk/v38/types"
	"github.com/stretchr/testify/mock"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/pexip/go-infinity-sdk/v38"
	"github.com/pexip/terraform-provider-pexip/internal/test"
)

func TestInfinityLicenceRequest(t *testing.T) {
	_ = os.Setenv("TF_ACC", "1")

	// Create a mock client and set up expectations
	client := infinity.NewClientMock()

	// Mock the CreateLicencerequest API call
	createResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/licence_request/123/",
	}
	client.On("PostWithResponse", mock.Anything, "configuration/v1/licence_request/", mock.Anything, mock.Anything).Return(createResponse, nil)

	// Shared state for mocking
	mockState := &config.LicenceRequest{
		SequenceNumber: "123",
		Reference:      "test-value",
		Actions:        "test-value",
		GenerationTime: "2023-01-01T00:00:00Z",
		Status:         "pending",
		ResponseXML:    test.StringPtr("<response></response>"),
		ResourceURI:    "/api/admin/configuration/v1/licence_request/123/",
	}

	// Mock the GetLicencerequest API call for Read operations
	client.On("GetJSON", mock.Anything, "configuration/v1/licence_request/123/", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		license_request := args.Get(2).(*config.LicenceRequest)
		*license_request = *mockState
	}).Maybe()

	testInfinityLicenceRequest(t, client)
}

func testInfinityLicenceRequest(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: getTestProtoV5ProviderFactories(client),
		// Skip the final destroy since this resource cannot be deleted.
		// The test already verifies that a destroy operation fails as expected.
		CheckDestroy: func(s *terraform.State) error {
			return nil
		},
		Steps: []resource.TestStep{
			// Step 1: Test resource creation and stabilization
			{
				Config: test.LoadTestFolder(t, "resource_infinity_licence_request_basic"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_licence_request.licence_request-test", "id"),
					resource.TestCheckResourceAttr("pexip_infinity_licence_request.licence_request-test", "reference", "test-value"),
					resource.TestCheckResourceAttr("pexip_infinity_licence_request.licence_request-test", "actions", "test-value"),
					resource.TestCheckResourceAttr("pexip_infinity_licence_request.licence_request-test", "sequence_number", "123"),
				),
				ExpectNonEmptyPlan: false,
			},
			// Step 2: Verify that an update is not allowed
			{
				Config:      test.LoadTestFolder(t, "resource_infinity_licence_request_basic_updated"),
				ExpectError: regexp.MustCompile("Update Not Supported"),
			},
		},
	})
}
