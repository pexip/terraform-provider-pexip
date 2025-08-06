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

func TestInfinityGatewayRoutingRule(t *testing.T) {
	_ = os.Setenv("TF_ACC", "1")

	// Create a mock client and set up expectations
	client := infinity.NewClientMock()

	// Mock the CreateGatewayroutingrule API call
	createResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/gateway_routing_rule/123/",
	}
	client.On("PostWithResponse", mock.Anything, "configuration/v1/gateway_routing_rule/", mock.Anything, mock.Anything).Return(createResponse, nil)

	// Shared state for mocking
	mockState := &config.GatewayRoutingRule{
		ID:               123,
		ResourceURI:      "/api/admin/configuration/v1/gateway_routing_rule/123/",
		Name:             "gateway_routing_rule-test",
		Description:      "Test GatewayRoutingRule",
		Priority:         123,
		Enable:           true,
		MatchString:      "test-value",
		ReplaceString:    "test-value",
		CalledDeviceType: "unknown",
		OutgoingProtocol: "sip",
		CallType:         "video",
		IVRTheme:         test.StringPtr("test-value"),
	}

	// Mock the GetGatewayroutingrule API call for Read operations
	client.On("GetJSON", mock.Anything, "configuration/v1/gateway_routing_rule/123/", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		gateway_routing_rule := args.Get(2).(*config.GatewayRoutingRule)
		*gateway_routing_rule = *mockState
	}).Maybe()

	// Mock the UpdateGatewayroutingrule API call
	client.On("PutJSON", mock.Anything, "configuration/v1/gateway_routing_rule/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		updateRequest := args.Get(2).(*config.GatewayRoutingRuleUpdateRequest)
		gateway_routing_rule := args.Get(3).(*config.GatewayRoutingRule)

		// Update mock state based on request
		mockState.Name = updateRequest.Name
		mockState.Description = updateRequest.Description
		if updateRequest.Priority != nil {
			mockState.Priority = *updateRequest.Priority
		}
		if updateRequest.Enable != nil {
			mockState.Enable = *updateRequest.Enable
		}
		mockState.MatchString = updateRequest.MatchString
		mockState.ReplaceString = updateRequest.ReplaceString
		mockState.CalledDeviceType = updateRequest.CalledDeviceType
		mockState.OutgoingProtocol = updateRequest.OutgoingProtocol
		mockState.CallType = updateRequest.CallType
		mockState.IVRTheme = updateRequest.IVRTheme

		// Return updated state
		*gateway_routing_rule = *mockState
	}).Maybe()

	// Mock the DeleteGatewayroutingrule API call
	client.On("DeleteJSON", mock.Anything, mock.MatchedBy(func(path string) bool {
		return path == "configuration/v1/gateway_routing_rule/123/"
	}), mock.Anything).Return(nil)

	testInfinityGatewayRoutingRule(t, client)
}

func testInfinityGatewayRoutingRule(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: getTestProtoV5ProviderFactories(client),
		Steps: []resource.TestStep{
			{
				Config: test.LoadTestFolder(t, "resource_infinity_gateway_routing_rule_basic"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_gateway_routing_rule.gateway_routing_rule-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_gateway_routing_rule.gateway_routing_rule-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_gateway_routing_rule.gateway_routing_rule-test", "name", "gateway_routing_rule-test"),
					resource.TestCheckResourceAttr("pexip_infinity_gateway_routing_rule.gateway_routing_rule-test", "description", "Test GatewayRoutingRule"),
					resource.TestCheckResourceAttr("pexip_infinity_gateway_routing_rule.gateway_routing_rule-test", "priority", "123"),
					resource.TestCheckResourceAttr("pexip_infinity_gateway_routing_rule.gateway_routing_rule-test", "enable", "true"),
				),
			},
			{
				Config: test.LoadTestFolder(t, "resource_infinity_gateway_routing_rule_basic_updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_gateway_routing_rule.gateway_routing_rule-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_gateway_routing_rule.gateway_routing_rule-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_gateway_routing_rule.gateway_routing_rule-test", "name", "gateway_routing_rule-test"),
					resource.TestCheckResourceAttr("pexip_infinity_gateway_routing_rule.gateway_routing_rule-test", "description", "Updated Test GatewayRoutingRule"),
					resource.TestCheckResourceAttr("pexip_infinity_gateway_routing_rule.gateway_routing_rule-test", "priority", "456"),
					resource.TestCheckResourceAttr("pexip_infinity_gateway_routing_rule.gateway_routing_rule-test", "enable", "false"),
				),
			},
		},
	})
}
