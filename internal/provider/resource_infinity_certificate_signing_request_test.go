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

func TestInfinityCertificateSigningRequest(t *testing.T) {
	t.Parallel()
	_ = os.Setenv("TF_ACC", "1")

	// Create a mock client and set up expectations
	client := infinity.NewClientMock()

	// Mock the CreateCertificatesigningrequest API call
	createResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/certificate_signing_request/123/",
	}
	client.On("PostWithResponse", mock.Anything, "configuration/v1/certificate_signing_request/", mock.Anything, mock.Anything).Return(createResponse, nil)

	// Shared state for mocking
	mockState := &config.CertificateSigningRequest{
		ID:                        123,
		ResourceURI:               "/api/admin/configuration/v1/certificate_signing_request/123/",
		SubjectName:               "certificate_signing_request-test",
		DN:                        "test-value",
		AdditionalSubjectAltNames: "certificate_signing_request-test",
		PrivateKeyType:            "rsa2048",
		PrivateKey:                test.StringPtr("test-value"),
		PrivateKeyPassphrase:      "test-value",
		AdCompatible:              true,
		TLSCertificate:            test.StringPtr("test-value"),
	}

	// Mock the GetCertificatesigningrequest API call for Read operations
	client.On("GetJSON", mock.Anything, "configuration/v1/certificate_signing_request/123/", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		certificate_signing_request := args.Get(2).(*config.CertificateSigningRequest)
		*certificate_signing_request = *mockState
	}).Maybe()

	// Mock the UpdateCertificatesigningrequest API call
	client.On("PutJSON", mock.Anything, "configuration/v1/certificate_signing_request/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		updateRequest := args.Get(2).(*config.CertificateSigningRequestUpdateRequest)
		certificate_signing_request := args.Get(3).(*config.CertificateSigningRequest)

		// Update mock state based on request
		mockState.SubjectName = updateRequest.SubjectName
		mockState.DN = updateRequest.DN
		mockState.AdditionalSubjectAltNames = updateRequest.AdditionalSubjectAltNames
		mockState.PrivateKeyType = updateRequest.PrivateKeyType
		mockState.PrivateKey = updateRequest.PrivateKey
		mockState.PrivateKeyPassphrase = updateRequest.PrivateKeyPassphrase
		if updateRequest.AdCompatible != nil {
			mockState.AdCompatible = *updateRequest.AdCompatible
		}
		mockState.TLSCertificate = updateRequest.TLSCertificate

		// Return updated state
		*certificate_signing_request = *mockState
	}).Maybe()

	// Mock the DeleteCertificatesigningrequest API call
	client.On("DeleteJSON", mock.Anything, mock.MatchedBy(func(path string) bool {
		return path == "configuration/v1/certificate_signing_request/123/"
	}), mock.Anything).Return(nil)

	testInfinityCertificateSigningRequest(t, client)
}

func testInfinityCertificateSigningRequest(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: getTestProtoV5ProviderFactories(client),
		Steps: []resource.TestStep{
			{
				Config: test.LoadTestFolder(t, "resource_infinity_certificate_signing_request_basic"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_certificate_signing_request.certificate_signing_request-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_certificate_signing_request.certificate_signing_request-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_certificate_signing_request.certificate_signing_request-test", "subject_name", "certificate_signing_request-test"),
					resource.TestCheckResourceAttr("pexip_infinity_certificate_signing_request.certificate_signing_request-test", "additional_subject_alt_names", "certificate_signing_request-test"),
					resource.TestCheckResourceAttr("pexip_infinity_certificate_signing_request.certificate_signing_request-test", "ad_compatible", "true"),
				),
			},
			{
				Config: test.LoadTestFolder(t, "resource_infinity_certificate_signing_request_basic_updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_certificate_signing_request.certificate_signing_request-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_certificate_signing_request.certificate_signing_request-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_certificate_signing_request.certificate_signing_request-test", "subject_name", "certificate_signing_request-test"),
					resource.TestCheckResourceAttr("pexip_infinity_certificate_signing_request.certificate_signing_request-test", "additional_subject_alt_names", "certificate_signing_request-test"),
					resource.TestCheckResourceAttr("pexip_infinity_certificate_signing_request.certificate_signing_request-test", "ad_compatible", "false"),
				),
			},
		},
	})
}
