/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/pexip/go-infinity-sdk/v38"
	"github.com/pexip/go-infinity-sdk/v38/config"
	"github.com/pexip/go-infinity-sdk/v38/types"
	"github.com/stretchr/testify/mock"

	"github.com/pexip/terraform-provider-pexip/internal/test"
)

func TestInfinitySIPProxy(t *testing.T) {
	t.Parallel()
	_ = os.Setenv("TF_ACC", "1")

	client := infinity.NewClientMock()

	// Mock state for the SIP proxy
	mockState := &config.SIPProxy{}

	// Mock the CreateSIPProxy API call (called in step 1 and step 4)
	client.On("PostWithResponse", mock.Anything, "configuration/v1/sip_proxy/", mock.Anything, mock.Anything).Return(&types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/sip_proxy/123/",
	}, nil).Run(func(args mock.Arguments) {
		req := args.Get(2).(*config.SIPProxyCreateRequest)
		*mockState = config.SIPProxy{
			ID:          123,
			ResourceURI: "/api/admin/configuration/v1/sip_proxy/123/",
			Name:        req.Name,
			Description: req.Description,
			Address:     req.Address,
			Port:        req.Port,
			Transport:   req.Transport,
		}
	}).Twice()

	// Mock the UpdateSIPProxy API call (called in step 2 and step 5)
	client.On("PutJSON", mock.Anything, "configuration/v1/sip_proxy/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		updateReq := args.Get(2).(*config.SIPProxyUpdateRequest)
		sipProxy := args.Get(3).(*config.SIPProxy)

		// Update mock state
		mockState.Name = updateReq.Name
		mockState.Description = updateReq.Description
		mockState.Address = updateReq.Address
		mockState.Port = updateReq.Port
		mockState.Transport = updateReq.Transport

		// Return updated state
		*sipProxy = *mockState
	}).Twice()

	// Mock the GetSIPProxy API call for Read operations
	client.On("GetJSON", mock.Anything, "configuration/v1/sip_proxy/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		sipProxy := args.Get(3).(*config.SIPProxy)
		*sipProxy = *mockState
	}).Maybe()

	// Mock the DeleteSIPProxy API call
	client.On("DeleteJSON", mock.Anything, "configuration/v1/sip_proxy/123/", mock.Anything).Return(nil)

	testInfinitySIPProxy(t, client)
}

func testInfinitySIPProxy(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: getTestProtoV5ProviderFactories(client),
		Steps: []resource.TestStep{
			{
				// Step 1: Create with full config
				Config: test.LoadTestFolder(t, "resource_infinity_sip_proxy_full"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_sip_proxy.sip-proxy-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_sip_proxy.sip-proxy-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_sip_proxy.sip-proxy-test", "name", "tf-test-sip-proxy-full"),
					resource.TestCheckResourceAttr("pexip_infinity_sip_proxy.sip-proxy-test", "description", "Full configuration test SIP proxy"),
					resource.TestCheckResourceAttr("pexip_infinity_sip_proxy.sip-proxy-test", "address", "sip.pexvclab.com"),
					resource.TestCheckResourceAttr("pexip_infinity_sip_proxy.sip-proxy-test", "port", "5061"),
					resource.TestCheckResourceAttr("pexip_infinity_sip_proxy.sip-proxy-test", "transport", "tls"),
				),
			},
			{
				// Step 2: Update to min config
				Config: test.LoadTestFolder(t, "resource_infinity_sip_proxy_min"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_sip_proxy.sip-proxy-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_sip_proxy.sip-proxy-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_sip_proxy.sip-proxy-test", "name", "tf-test-sip-proxy-min"),
					resource.TestCheckResourceAttr("pexip_infinity_sip_proxy.sip-proxy-test", "description", ""),
					resource.TestCheckResourceAttr("pexip_infinity_sip_proxy.sip-proxy-test", "address", "sip.pexvclab.com"),
					resource.TestCheckResourceAttr("pexip_infinity_sip_proxy.sip-proxy-test", "port", "5061"),
					resource.TestCheckResourceAttr("pexip_infinity_sip_proxy.sip-proxy-test", "transport", "tls"),
				),
			},
			{
				// Step 3: Destroy (implicit)
				Config:  test.LoadTestFolder(t, "resource_infinity_sip_proxy_min"),
				Destroy: true,
			},
			{
				// Step 4: Create with min config
				Config: test.LoadTestFolder(t, "resource_infinity_sip_proxy_min"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_sip_proxy.sip-proxy-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_sip_proxy.sip-proxy-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_sip_proxy.sip-proxy-test", "name", "tf-test-sip-proxy-min"),
					resource.TestCheckResourceAttr("pexip_infinity_sip_proxy.sip-proxy-test", "description", ""),
					resource.TestCheckResourceAttr("pexip_infinity_sip_proxy.sip-proxy-test", "address", "sip.pexvclab.com"),
					resource.TestCheckResourceAttr("pexip_infinity_sip_proxy.sip-proxy-test", "port", "5061"),
					resource.TestCheckResourceAttr("pexip_infinity_sip_proxy.sip-proxy-test", "transport", "tls"),
				),
			},
			{
				// Step 5: Update to full config
				Config: test.LoadTestFolder(t, "resource_infinity_sip_proxy_full"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_sip_proxy.sip-proxy-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_sip_proxy.sip-proxy-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_sip_proxy.sip-proxy-test", "name", "tf-test-sip-proxy-full"),
					resource.TestCheckResourceAttr("pexip_infinity_sip_proxy.sip-proxy-test", "description", "Full configuration test SIP proxy"),
					resource.TestCheckResourceAttr("pexip_infinity_sip_proxy.sip-proxy-test", "address", "sip.pexvclab.com"),
					resource.TestCheckResourceAttr("pexip_infinity_sip_proxy.sip-proxy-test", "port", "5061"),
					resource.TestCheckResourceAttr("pexip_infinity_sip_proxy.sip-proxy-test", "transport", "tls"),
				),
			},
		},
	})
}
