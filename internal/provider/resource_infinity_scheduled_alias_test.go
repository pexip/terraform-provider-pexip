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

func TestInfinityScheduledAlias(t *testing.T) {
	_ = os.Setenv("TF_ACC", "1")

	// Create a mock client and set up expectations
	client := infinity.NewClientMock()

	// Mock the CreateScheduledalias API call
	createResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/scheduled_alias/123/",
	}
	client.On("PostWithResponse", mock.Anything, "configuration/v1/scheduled_alias/", mock.Anything, mock.Anything).Return(createResponse, nil)

	// Shared state for mocking
	mockState := &config.ScheduledAlias{
		ID:                123,
		ResourceURI:       "/api/admin/configuration/v1/scheduled_alias/123/",
		Alias:             "test-value",
		AliasNumber:       1234567890,
		NumericAlias:      "test-value",
		UUID:              "test-value",
		ExchangeConnector: "test-value",
		IsUsed:            true,
		EWSItemUID:        test.StringPtr("test-value"),
	}

	// Mock the GetScheduledalias API call for Read operations
	client.On("GetJSON", mock.Anything, "configuration/v1/scheduled_alias/123/", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		scheduled_alias := args.Get(2).(*config.ScheduledAlias)
		*scheduled_alias = *mockState
	}).Maybe()

	// Mock the UpdateScheduledalias API call
	client.On("PutJSON", mock.Anything, "configuration/v1/scheduled_alias/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		updateRequest := args.Get(2).(*config.ScheduledAliasUpdateRequest)
		scheduled_alias := args.Get(3).(*config.ScheduledAlias)

		// Update mock state based on request
		if updateRequest.Alias != "" {
			mockState.Alias = updateRequest.Alias
		}
		if updateRequest.AliasNumber != nil {
			mockState.AliasNumber = *updateRequest.AliasNumber
		}
		if updateRequest.NumericAlias != "" {
			mockState.NumericAlias = updateRequest.NumericAlias
		}
		if updateRequest.UUID != "" {
			mockState.UUID = updateRequest.UUID
		}
		if updateRequest.ExchangeConnector != "" {
			mockState.ExchangeConnector = updateRequest.ExchangeConnector
		}
		if updateRequest.IsUsed != nil {
			mockState.IsUsed = *updateRequest.IsUsed
		}
		if updateRequest.EWSItemUID != nil {
			mockState.EWSItemUID = updateRequest.EWSItemUID
		}

		// Return updated state
		*scheduled_alias = *mockState
	}).Maybe()

	// Mock the DeleteScheduledalias API call
	client.On("DeleteJSON", mock.Anything, mock.MatchedBy(func(path string) bool {
		return path == "configuration/v1/scheduled_alias/123/"
	}), mock.Anything).Return(nil)

	testInfinityScheduledAlias(t, client)
}

func testInfinityScheduledAlias(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: getTestProtoV5ProviderFactories(client),
		Steps: []resource.TestStep{
			{
				Config: test.LoadTestFolder(t, "resource_infinity_scheduled_alias_basic"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_scheduled_alias.scheduled_alias-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_scheduled_alias.scheduled_alias-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_scheduled_alias.scheduled_alias-test", "is_used", "true"),
				),
			},
			{
				Config: test.LoadTestFolder(t, "resource_infinity_scheduled_alias_basic_updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_scheduled_alias.scheduled_alias-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_scheduled_alias.scheduled_alias-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_scheduled_alias.scheduled_alias-test", "is_used", "false"),
				),
			},
		},
	})
}
