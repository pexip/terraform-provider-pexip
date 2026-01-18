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
MIIDeTCCAmGgAwIBAgIUBM8euKK5qdSr9d5bFtCFOGk/GnQwDQYJKoZIhvcNAQEL
BQAwTDELMAkGA1UEBhMCVVMxDTALBgNVBAgMBFRlc3QxDTALBgNVBAcMBFRlc3Qx
DTALBgNVBAoMBFRlc3QxEDAOBgNVBAMMB1Rlc3QgQ0EwHhcNMjUxMjEyMTk0MjIw
WhcNMjYxMjEyMTk0MjIwWjBMMQswCQYDVQQGEwJVUzENMAsGA1UECAwEVGVzdDEN
MAsGA1UEBwwEVGVzdDENMAsGA1UECgwEVGVzdDEQMA4GA1UEAwwHVGVzdCBDQTCC
ASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAMu69O8MRO1gNjTgmkRIyNJV
ASmVqT1GhM7R+epArMfRnPYkMCsv2wPrWzqnwvzMY/UUIrwplon0gdb6fopJR40y
bmEY26sYfgLK9n7Nij0QGRMBvlz+fvXabJZE+pkVu5u9iompUuqGgSo0R/jqDi6+
gVupTtDGQVL2olj5jaeBKg/WMcxSNHTIkCzO7402nYyULyyG4n7a8KHWhR7SGaoM
ssl5MYRiUZ1lt6Rt31oIFwCmWks0LxrFcHiT7YiVbgSh5g9G942E9MMcSf2H/D2J
JwzQUFk7CSQRmFZa1H9WkGsAMSuJ8zKR81zIseowBJDegX7+mKcqePKyQkkLyO8C
AwEAAaNTMFEwHQYDVR0OBBYEFBJp92iZ2kfjm7sjet7PAOgxEt6TMB8GA1UdIwQY
MBaAFBJp92iZ2kfjm7sjet7PAOgxEt6TMA8GA1UdEwEB/wQFMAMBAf8wDQYJKoZI
hvcNAQELBQADggEBALqCSgcWwSxia54EKv5dWqEUmNWbchHC+/08TyjiFhwN+q6f
TAZ0ZzKueB/FguaPsyCRDkbc3Y2eSDNwMf4wJsf913bw52bTQ6DwtauerHsF9Ywp
H/M2VwkoN31qjaV1BRJmyvSpgMWyoBuCY7j24fSNACRrVCqdclW3wer9fcNY0O+8
pUt7qzJf9L2tYEpehIoI69tX2UWfOjkSp7kWmFTIT3UfF+UWii3RXab5WaVF0DV6
uD8tWBa2wjkjNBGri2IIRYdgDNy9YcnExLeu6p2XrlxgHiyclA52AspuSwKcwnAS
g2OX6o8rGmwv2UgI1X+x9kLrdj0OFenKrwaBiEI=
-----END CERTIFICATE-----
`,
		TrustedIntermediate: false, // Root CAs (Subject == Issuer) have trusted_intermediate = false
		SubjectName:         "Test CA",
		IssuerName:          "Test CA",
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
		if updateRequest.Certificate != "" && updateRequest.Certificate != mockState.Certificate {
			// Certificate is being updated - update subject and issuer
			mockState.Certificate = updateRequest.Certificate
			mockState.SubjectName = "Updated CA"
			mockState.IssuerName = "Updated CA"
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
					resource.TestCheckResourceAttr("pexip_infinity_ca_certificate.ca_certificate-test", "trusted_intermediate", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_ca_certificate.ca_certificate-test", "subject_name", "Test CA"),
					resource.TestCheckResourceAttr("pexip_infinity_ca_certificate.ca_certificate-test", "issuer_name", "Test CA"),
				),
			},
			{
				Config: test.LoadTestFolder(t, "resource_infinity_ca_certificate_basic_updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_ca_certificate.ca_certificate-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_ca_certificate.ca_certificate-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_ca_certificate.ca_certificate-test", "trusted_intermediate", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_ca_certificate.ca_certificate-test", "subject_name", "Updated CA"),
					resource.TestCheckResourceAttr("pexip_infinity_ca_certificate.ca_certificate-test", "issuer_name", "Updated CA"),
				),
			},
		},
	})
}
