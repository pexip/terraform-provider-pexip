/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package provider

import (
	"os"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/pexip/go-infinity-sdk/v38"

	"github.com/pexip/terraform-provider-pexip/internal/test"
)

func TestInfinityDeleteDefaultTLSCertificateAction(t *testing.T) {
	t.Parallel()
	_ = os.Setenv("TF_ACC", "1")

	// Create a mock client and set up expectations
	client := infinity.NewClientMock()

	// Mock the DeleteTLSCertificate API call (certificate ID 1 is the default)
	client.On("DeleteJSON", mock.Anything, "configuration/v1/tls_certificate/1/", mock.Anything).Return(nil)

	// Note: Delete default TLS certificate is an action resource that doesn't have persistent state to read

	testInfinityDeleteDefaultTLSCertificateAction(t, client)
}

func testInfinityDeleteDefaultTLSCertificateAction(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: getTestProtoV5ProviderFactories(client),
		Steps: []resource.TestStep{
			{
				Config: test.LoadTestFolder(t, "resource_infinity_delete_default_tls_certificate_action_basic"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_delete_default_tls_certificate_action.delete-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_delete_default_tls_certificate_action.delete-test", "timestamp"),
				),
			},
		},
	})
}
