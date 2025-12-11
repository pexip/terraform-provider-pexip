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

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/pexip/go-infinity-sdk/v38"

	"github.com/pexip/terraform-provider-pexip/internal/test"
)

func TestInfinityCACertificate(t *testing.T) {
	t.Parallel()
	_ = os.Setenv("TF_ACC", "1")

	// Create a mock client and set up expectations
	client := infinity.NewClientMock()

	// Mock the CreateCACertificate API call
	createResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/ca_certificate/123/",
	}
	client.On("PostWithResponse", mock.Anything, "configuration/v1/ca_certificate/", mock.Anything, mock.Anything).Return(createResponse, nil)

	// Shared state for mocking
	mockState := &config.CACertificate{
		ID:          123,
		ResourceURI: "/api/admin/configuration/v1/ca_certificate/123/",
		Certificate: `-----BEGIN CERTIFICATE-----
MIIDXTCCAkWgAwIBAgIJAKJ5JqfJzhqBMA0GCSqGSIb3DQEBCwUAMEUxCzAJBgNV
BAYTAkFVMRMwEQYDVQQIDApTb21lLVN0YXRlMSEwHwYDVQQKDBhJbnRlcm5ldCBX
aWRnaXRzIFB0eSBMdGQwHhcNMjMwMTAxMDAwMDAwWhcNMjQwMTAxMDAwMDAwWjBF
MQswCQYDVQQGEwJBVTETMBEGA1UECAwKU29tZS1TdGF0ZTEhMB8GA1UECgwYSW50
ZXJuZXQgV2lkZ2l0cyBQdHkgTHRkMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIB
CgKCAQEAtxvDqPhIZgQEo5bN4qhKqAJhvfxPzMEBjQVf9yQcNwVzXrKKQiJcPPjD
xaZRNGxN7H5qr+xVqJhYnXKjXzGRzCbNn6KnKEWnYJC7qR9pYLqAHWHQpzPqLa2V
lZqsL3cXZ3uOr5qnLsE8nNfnRqV5j8qPqCvKgLtNRhHnNQ2VWKYqNnGzPqLj8LWY
xLCnW3xNcBqH5qRHYPnQpRdLqNQzKVzQqPqLj8LWYxLCnW3xNcBqH5qRHYPnQpRd
LqNQzKVzQqPqLjcXZ3uOr5qnLsE8nNfnRqV5j8qPqCvKgLtNRhHnNQ2VWKYqNnGz
PqLj8LWYxLCnW3xNcBqH5qRHYPnQpRdLqNQzKVzQqPqLjQIDAQABo1AwTjAdBgNV
HQ4EFgQUX8mfqCJhJ8HQNGKjGN6wQVqPXdYwHwYDVR0jBBgwFoAUX8mfqCJhJ8HQ
NGKjGN6wQVqPXdYwDAYDVR0TBAUwAwEB/zANBgkqhkiG9w0BAQsFAAOCAQEAjKpO
nZqNz0bHvYkBvPqLj8LWYxLCnW3xNcBqH5qRHYPnQpRdLqNQzKVzQqPqLjcXZ3uO
r5qnLsE8nNfnRqV5j8qPqCvKgLtNRhHnNQ2VWKYqNnGzPqLj8LWYxLCnW3xNcBqH
-----END CERTIFICATE-----
`,
		TrustedIntermediate: true,
		SubjectName:         "ca_certificate-test",
		IssuerName:          "ca_certificate-test",
	}

	// Mock the GetCACertificate API call for Read operations
	client.On("GetJSON", mock.Anything, "configuration/v1/ca_certificate/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		caCertificate := args.Get(3).(*config.CACertificate)
		*caCertificate = *mockState
	}).Maybe()

	// Mock the UpdateCACertificate API call
	client.On("PatchJSON", mock.Anything, "configuration/v1/ca_certificate/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		updateRequest := args.Get(2).(*config.CACertificateUpdateRequest)
		caCertificate := args.Get(3).(*config.CACertificate)

		// Update mock state
		if updateRequest.Certificate != "" {
			mockState.Certificate = updateRequest.Certificate
		}
		mockState.TrustedIntermediate = updateRequest.TrustedIntermediate

		// Return updated state
		*caCertificate = *mockState
	}).Maybe()

	// Mock the DeleteCACertificate API call
	client.On("DeleteJSON", mock.Anything, mock.MatchedBy(func(path string) bool {
		return path == "configuration/v1/ca_certificate/123/"
	}), mock.Anything).Return(nil)

	testInfinityCACertificate(t, client)
}

func testInfinityCACertificate(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: getTestProtoV5ProviderFactories(client),
		Steps: []resource.TestStep{
			{
				Config: test.LoadTestFolder(t, "resource_infinity_ca_certificate_basic"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_ca_certificate.ca_certificate-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_ca_certificate.ca_certificate-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_ca_certificate.ca_certificate-test", "trusted_intermediate", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_ca_certificate.ca_certificate-test", "subject_name", "ca_certificate-test"),
					resource.TestCheckResourceAttr("pexip_infinity_ca_certificate.ca_certificate-test", "issuer_name", "ca_certificate-test"),
				),
			},
			{
				Config: test.LoadTestFolder(t, "resource_infinity_ca_certificate_basic_updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_ca_certificate.ca_certificate-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_ca_certificate.ca_certificate-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_ca_certificate.ca_certificate-test", "trusted_intermediate", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_ca_certificate.ca_certificate-test", "subject_name", "ca_certificate-test"),
					resource.TestCheckResourceAttr("pexip_infinity_ca_certificate.ca_certificate-test", "issuer_name", "ca_certificate-test"),
				),
			},
		},
	})
}
