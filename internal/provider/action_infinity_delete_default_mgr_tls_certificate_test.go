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

func TestInfinityDeleteDefaultMgrTLSCertificateAction(t *testing.T) {
	t.Parallel()
	_ = os.Setenv("TF_ACC", "1")

	// Create a mock client and set up expectations
	client := infinity.NewClientMock()

	// Mock the DeleteTLSCertificate API call (certificate ID 1 is the default)
	client.On("DeleteJSON", mock.Anything, "configuration/v1/tls_certificate/1/", mock.Anything).Return(nil)

	testInfinityDeleteDefaultMgrTLSCertificateAction(t, client)
}

func testInfinityDeleteDefaultMgrTLSCertificateAction(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: getTestProtoV5ProviderFactories(client),
		Steps: []resource.TestStep{
			{
				Config: test.LoadTestFolder(t, "action_infinity_delete_default_mgr_tls_certificate_basic"),
				// Actions don't have persistent state, so we just verify the configuration is valid
				// The mock client will verify that DeleteTLSCertificate was called
			},
		},
	})
}
