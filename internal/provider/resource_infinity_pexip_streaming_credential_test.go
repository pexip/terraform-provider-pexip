package provider

import (
	"github.com/pexip/go-infinity-sdk/v38/config"
	"github.com/pexip/go-infinity-sdk/v38/types"
	"github.com/stretchr/testify/mock"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/pexip/go-infinity-sdk/v38"
	"github.com/pexip/terraform-provider-pexip/internal/test"
)

func TestInfinityPexipStreamingCredential(t *testing.T) {
	_ = os.Setenv("TF_ACC", "1")

	// Create a mock client and set up expectations
	client := infinity.NewClientMock()

	// Mock the CreatePexipstreamingcredential API call
	createResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/pexip_streaming_credential/123/",
	}
	client.On("PostWithResponse", mock.Anything, "configuration/v1/pexip_streaming_credential/", mock.Anything, mock.Anything).Return(createResponse, nil)

	// Shared state for mocking
	mockState := &config.PexipStreamingCredential{
		ID:          123,
		ResourceURI: "/api/admin/configuration/v1/pexip_streaming_credential/123/",
		Kid:         "test-value",
		PublicKey:   "test-value",
	}

	// Mock the GetPexipstreamingcredential API call for Read operations
	client.On("GetJSON", mock.Anything, "configuration/v1/pexip_streaming_credential/123/", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		pexip_streaming_credential := args.Get(2).(*config.PexipStreamingCredential)
		*pexip_streaming_credential = *mockState
	}).Maybe()

	// Mock the UpdatePexipstreamingcredential API call
	client.On("PutJSON", mock.Anything, "configuration/v1/pexip_streaming_credential/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		updateRequest := args.Get(2).(*config.PexipStreamingCredentialUpdateRequest)
		pexip_streaming_credential := args.Get(3).(*config.PexipStreamingCredential)

		// Update mock state based on request
		if updateRequest.Kid != "" {
			mockState.Kid = updateRequest.Kid
		}
		if updateRequest.PublicKey != "" {
			mockState.PublicKey = updateRequest.PublicKey
		}

		// Return updated state
		*pexip_streaming_credential = *mockState
	}).Maybe()

	// Mock the DeletePexipstreamingcredential API call
	client.On("DeleteJSON", mock.Anything, mock.MatchedBy(func(path string) bool {
		return path == "configuration/v1/pexip_streaming_credential/123/"
	}), mock.Anything).Return(nil)

	testInfinityPexipStreamingCredential(t, client)
}

func testInfinityPexipStreamingCredential(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: getTestProtoV5ProviderFactories(client),
		Steps: []resource.TestStep{
			{
				Config: test.LoadTestFolder(t, "resource_infinity_pexip_streaming_credential_basic"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_pexip_streaming_credential.pexip_streaming_credential-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_pexip_streaming_credential.pexip_streaming_credential-test", "resource_id"),
				),
			},
			{
				Config: test.LoadTestFolder(t, "resource_infinity_pexip_streaming_credential_basic_updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_pexip_streaming_credential.pexip_streaming_credential-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_pexip_streaming_credential.pexip_streaming_credential-test", "resource_id"),
				),
			},
		},
	})
}
