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

func TestInfinityTLSCertificate(t *testing.T) {
	t.Parallel()
	_ = os.Setenv("TF_ACC", "1")

	// Create a mock client and set up expectations
	client := infinity.NewClientMock()

	// Mock the CreateTLSCertificate API call
	createResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/tls_certificate/123/",
	}
	client.On("PostWithResponse", mock.Anything, "configuration/v1/tls_certificate/", mock.Anything, mock.Anything).Return(createResponse, nil)

	tlsCertificate := test.LoadTestFile(t, "tls_certificate.pem")
	tlsPrivateKey := test.LoadTestFile(t, "tls_private_key.pem")

	// Shared state for mocking
	mockState := &config.TLSCertificate{
		ID:          123,
		ResourceURI: "/api/admin/configuration/v1/tls_certificate/123/",
		Certificate: tlsCertificate,
		PrivateKey:  tlsPrivateKey,
		Nodes:       []string{}, // Start with empty slice (empty list in Terraform)
	}

	// Mock the GetTLSCertificate API call for Read operations
	client.On("GetJSON", mock.Anything, "configuration/v1/tls_certificate/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		tlsCert := args.Get(3).(*config.TLSCertificate)
		*tlsCert = *mockState
	}).Maybe()

	// Mock the UpdateTLSCertificate API call
	client.On("PatchJSON", mock.Anything, "configuration/v1/tls_certificate/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		updateRequest := args.Get(2).(*config.TLSCertificateUpdateRequest)
		tlsCert := args.Get(3).(*config.TLSCertificate)

		// Update mock state
		if updateRequest.Certificate != "" {
			mockState.Certificate = updateRequest.Certificate
		}
		if updateRequest.PrivateKey != "" {
			mockState.PrivateKey = updateRequest.PrivateKey
		}
		if updateRequest.Nodes != nil {
			mockState.Nodes = updateRequest.Nodes
		}

		// Return updated state
		*tlsCert = *mockState
	}).Maybe()

	// Mock the DeleteTLSCertificate API call
	client.On("DeleteJSON", mock.Anything, mock.MatchedBy(func(path string) bool {
		return path == "configuration/v1/tls_certificate/123/"
	}), mock.Anything).Return(nil)

	testInfinityTLSCertificate(t, client, tlsPrivateKey, tlsCertificate)
}

func testInfinityTLSCertificate(t *testing.T, client InfinityClient, privateKey, certificate string) {
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: getTestProtoV5ProviderFactories(client),
		Steps: []resource.TestStep{
			{
				Config: test.LoadTestFolder(t, "resource_infinity_tls_certificate_basic"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_tls_certificate.tls-cert-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_tls_certificate.tls-cert-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_tls_certificate.tls-cert-test", "certificate", certificate),
					resource.TestCheckResourceAttr("pexip_infinity_tls_certificate.tls-cert-test", "private_key", privateKey),
				),
			},
			{
				Config: test.LoadTestFolder(t, "resource_infinity_tls_certificate_basic_updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_tls_certificate.tls-cert-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_tls_certificate.tls-cert-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_tls_certificate.tls-cert-test", "certificate", certificate),
					resource.TestCheckResourceAttr("pexip_infinity_tls_certificate.tls-cert-test", "private_key", privateKey),
				),
			},
		},
	})
}
