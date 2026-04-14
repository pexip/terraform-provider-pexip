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

	// Set test environment variables for Terraform variables
	_ = os.Setenv("TF_VAR_infinity_gms_gw_token_cert", "-----BEGIN CERTIFICATE-----\nserver-cert\n-----END CERTIFICATE-----\n-----BEGIN CERTIFICATE-----\nintermediate-cert\n-----END CERTIFICATE-----\n")
	_ = os.Setenv("TF_VAR_infinity_gms_gw_token_key", "test-private-key-data")
	_ = os.Setenv("TF_VAR_infinity_gms_gw_token_cert2", "-----BEGIN CERTIFICATE-----\nserver-cert2\n-----END CERTIFICATE-----\n-----BEGIN CERTIFICATE-----\nintermediate-cert2\n-----END CERTIFICATE-----\n")
	_ = os.Setenv("TF_VAR_infinity_gms_gw_token_key2", "test-private-key-data2")

	testInfinityGMSGatewayToken(t, client)
}

func testInfinityGMSGatewayToken(t *testing.T, client InfinityClient) {
	// Verify required environment variables are set
	requiredEnvVars := map[string]string{
		"TF_VAR_infinity_gms_gw_token_cert":  os.Getenv("TF_VAR_infinity_gms_gw_token_cert"),
		"TF_VAR_infinity_gms_gw_token_key":   os.Getenv("TF_VAR_infinity_gms_gw_token_key"),
		"TF_VAR_infinity_gms_gw_token_cert2": os.Getenv("TF_VAR_infinity_gms_gw_token_cert2"),
		"TF_VAR_infinity_gms_gw_token_key2":  os.Getenv("TF_VAR_infinity_gms_gw_token_key2"),
	}

	for name, value := range requiredEnvVars {
		if value == "" {
			t.Fatalf("Required environment variable %s is not set", name)
		}
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: getTestProtoV6ProviderFactories(client),
		Steps: []resource.TestStep{
			// Step 1: Create with min config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_gms_gateway_token_min"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_gms_gateway_token.test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_gms_gateway_token.test", "certificate"),
					resource.TestCheckResourceAttrSet("pexip_infinity_gms_gateway_token.test", "private_key"),
					resource.TestCheckResourceAttrSet("pexip_infinity_gms_gateway_token.test", "intermediate_certificate"),
					resource.TestCheckResourceAttrSet("pexip_infinity_gms_gateway_token.test", "leaf_certificate"),
					resource.TestCheckResourceAttrSet("pexip_infinity_gms_gateway_token.test", "supports_direct_guest_join"),
				),
			},
			// Step 2: Update with full config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_gms_gateway_token_full"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_gms_gateway_token.test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_gms_gateway_token.test", "certificate"),
					resource.TestCheckResourceAttrSet("pexip_infinity_gms_gateway_token.test", "private_key"),
					resource.TestCheckResourceAttrSet("pexip_infinity_gms_gateway_token.test", "intermediate_certificate"),
					resource.TestCheckResourceAttrSet("pexip_infinity_gms_gateway_token.test", "leaf_certificate"),
					resource.TestCheckResourceAttrSet("pexip_infinity_gms_gateway_token.test", "supports_direct_guest_join"),
				),
			},
		},
	})
}
