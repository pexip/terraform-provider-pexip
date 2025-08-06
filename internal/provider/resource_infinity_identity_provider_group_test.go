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

func TestInfinityIdentityProviderGroup(t *testing.T) {
	t.Parallel()
	_ = os.Setenv("TF_ACC", "1")

	// Create a mock client and set up expectations
	client := infinity.NewClientMock()

	// Mock the CreateIdentityprovidergroup API call
	createResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/identity_provider_group/123/",
	}
	client.On("PostWithResponse", mock.Anything, "configuration/v1/identity_provider_group/", mock.Anything, mock.Anything).Return(createResponse, nil)

	// Shared state for mocking
	mockState := &config.IdentityProviderGroup{
		ID:          123,
		ResourceURI: "/api/admin/configuration/v1/identity_provider_group/123/",
		Name:        "identity_provider_group-test",
		Description: "Test IdentityProviderGroup",
	}

	// Mock the GetIdentityprovidergroup API call for Read operations
	client.On("GetJSON", mock.Anything, "configuration/v1/identity_provider_group/123/", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		identity_provider_group := args.Get(2).(*config.IdentityProviderGroup)
		*identity_provider_group = *mockState
	}).Maybe()

	// Mock the UpdateIdentityprovidergroup API call
	client.On("PutJSON", mock.Anything, "configuration/v1/identity_provider_group/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		updateRequest := args.Get(2).(*config.IdentityProviderGroupUpdateRequest)
		identity_provider_group := args.Get(3).(*config.IdentityProviderGroup)

		// Update mock state based on request
		if updateRequest.Name != "" {
			mockState.Name = updateRequest.Name
		}
		if updateRequest.Description != "" {
			mockState.Description = updateRequest.Description
		}

		// Return updated state
		*identity_provider_group = *mockState
	}).Maybe()

	// Mock the DeleteIdentityprovidergroup API call
	client.On("DeleteJSON", mock.Anything, mock.MatchedBy(func(path string) bool {
		return path == "configuration/v1/identity_provider_group/123/"
	}), mock.Anything).Return(nil)

	testInfinityIdentityProviderGroup(t, client)
}

func testInfinityIdentityProviderGroup(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: getTestProtoV5ProviderFactories(client),
		Steps: []resource.TestStep{
			{
				Config: test.LoadTestFolder(t, "resource_infinity_identity_provider_group_basic"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_identity_provider_group.identity_provider_group-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_identity_provider_group.identity_provider_group-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider_group.identity_provider_group-test", "name", "identity_provider_group-test"),
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider_group.identity_provider_group-test", "description", "Test IdentityProviderGroup"),
				),
			},
			{
				Config: test.LoadTestFolder(t, "resource_infinity_identity_provider_group_basic_updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_identity_provider_group.identity_provider_group-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_identity_provider_group.identity_provider_group-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider_group.identity_provider_group-test", "name", "identity_provider_group-test"),
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider_group.identity_provider_group-test", "description", "Updated Test IdentityProviderGroup"),
				),
			},
		},
	})
}
