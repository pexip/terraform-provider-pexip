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
	"github.com/stretchr/testify/mock"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/pexip/go-infinity-sdk/v38"

	"github.com/pexip/terraform-provider-pexip/internal/test"
)

func TestInfinityGMSGatewayToken(t *testing.T) {
	t.Parallel()
	_ = os.Setenv("TF_ACC", "1")

	// Create a mock client and set up expectations
	client := infinity.NewClientMock()

	// Shared state for mocking
	supportsDirectGuestJoin := false
	mockState := &config.GMSGatewayToken{
		ID:                      1,
		ResourceURI:             "/api/admin/configuration/v1/gms_gateway_token/1/",
		Certificate:             "",                                                                                          // API clears this after update
		IntermediateCertificate: test.StringPtr("-----BEGIN CERTIFICATE-----\nintermediate-cert\n-----END CERTIFICATE-----"), // Extracted from chain
		LeafCertificate:         test.StringPtr("-----BEGIN CERTIFICATE-----\nserver-cert\n-----END CERTIFICATE-----"),       // Extracted from chain
		SupportsDirectGuestJoin: &supportsDirectGuestJoin,
	}

	// Mock the GetGMSGatewayToken API call for Read operations
	client.On(
		"GetJSON",
		mock.Anything,                           // context.Context
		"configuration/v1/gms_gateway_token/1/", // string
		mock.Anything,                           // *url.Values or nil
		mock.AnythingOfType("*config.GMSGatewayToken"), // pointer to config.GMSGatewayToken
	).Return(nil).Run(func(args mock.Arguments) {
		gmsGatewayToken := args.Get(3).(*config.GMSGatewayToken)
		*gmsGatewayToken = *mockState
	}).Maybe()

	// Mock the UpdateGMSGatewayToken API call
	client.On("PatchJSON", mock.Anything, "configuration/v1/gms_gateway_token/1/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		gmsGatewayToken := args.Get(3).(*config.GMSGatewayToken)

		// Simulate API behavior:
		// 1. Certificate is cleared after update (API processes the chain)
		// 2. API extracts leaf (server) cert and intermediate cert from the chain
		// 3. PrivateKey is not returned by API
		mockState.Certificate = ""
		mockState.IntermediateCertificate = test.StringPtr("-----BEGIN CERTIFICATE-----\nintermediate-cert\n-----END CERTIFICATE-----")
		mockState.LeafCertificate = test.StringPtr("-----BEGIN CERTIFICATE-----\nserver-cert\n-----END CERTIFICATE-----")

		// Return updated state
		*gmsGatewayToken = *mockState
	}).Maybe()

	testInfinityGMSGatewayToken(t, client)
}

func testInfinityGMSGatewayToken(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: getTestProtoV5ProviderFactories(client),
		Steps: []resource.TestStep{
			{
				Config: test.LoadTestFolder(t, "resource_infinity_gms_gateway_token_basic"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_gms_gateway_token.gms_gateway_token_test", "id"),
					// Certificate from plan is preserved in state
					resource.TestCheckResourceAttr("pexip_infinity_gms_gateway_token.gms_gateway_token_test", "certificate", "-----BEGIN CERTIFICATE-----\nserver-cert\n-----END CERTIFICATE-----\n-----BEGIN CERTIFICATE-----\nintermediate-cert\n-----END CERTIFICATE-----\n"),
					resource.TestCheckResourceAttr("pexip_infinity_gms_gateway_token.gms_gateway_token_test", "private_key", "test-private-key-data"),
					// API extracts and returns individual certs
					resource.TestCheckResourceAttr("pexip_infinity_gms_gateway_token.gms_gateway_token_test", "intermediate_certificate", "-----BEGIN CERTIFICATE-----\nintermediate-cert\n-----END CERTIFICATE-----"),
					resource.TestCheckResourceAttr("pexip_infinity_gms_gateway_token.gms_gateway_token_test", "leaf_certificate", "-----BEGIN CERTIFICATE-----\nserver-cert\n-----END CERTIFICATE-----"),
					resource.TestCheckResourceAttr("pexip_infinity_gms_gateway_token.gms_gateway_token_test", "supports_direct_guest_join", "false"),
				),
			},
		},
	})
}
