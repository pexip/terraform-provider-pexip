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

func TestInfinityMjxExchangeAutodiscoverURL(t *testing.T) {
	t.Parallel()
	_ = os.Setenv("TF_ACC", "1")

	client := infinity.NewClientMock()

	mockState := &config.MjxExchangeAutodiscoverURL{}

	createResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/mjx_exchange_autodiscover_url/123/",
	}
	client.On("PostWithResponse", mock.Anything, "configuration/v1/mjx_exchange_autodiscover_url/", mock.Anything, mock.Anything).Return(createResponse, nil).Run(func(args mock.Arguments) {
		createReq := args.Get(2).(*config.MjxExchangeAutodiscoverURLCreateRequest)
		*mockState = config.MjxExchangeAutodiscoverURL{
			ID:                 123,
			ResourceURI:        "/api/admin/configuration/v1/mjx_exchange_autodiscover_url/123/",
			Name:               createReq.Name,
			Description:        createReq.Description,
			URL:                createReq.URL,
			ExchangeDeployment: createReq.ExchangeDeployment,
		}
	})

	client.On("GetJSON", mock.Anything, "configuration/v1/mjx_exchange_autodiscover_url/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		autodiscoverURL := args.Get(3).(*config.MjxExchangeAutodiscoverURL)
		*autodiscoverURL = *mockState
	}).Maybe()

	client.On("PutJSON", mock.Anything, "configuration/v1/mjx_exchange_autodiscover_url/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		updateReq := args.Get(2).(*config.MjxExchangeAutodiscoverURLUpdateRequest)
		mockState.Name = updateReq.Name
		mockState.Description = updateReq.Description
		mockState.URL = updateReq.URL
		mockState.ExchangeDeployment = updateReq.ExchangeDeployment
	}).Maybe()

	client.On("DeleteJSON", mock.Anything, "configuration/v1/mjx_exchange_autodiscover_url/123/", mock.Anything).Return(nil)

	testInfinityMjxExchangeAutodiscoverURL(t, client)
}

func testInfinityMjxExchangeAutodiscoverURL(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: getTestProtoV6ProviderFactories(client),
		Steps: []resource.TestStep{
			// Step 1: Create with full config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_mjx_exchange_autodiscover_url_full"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_exchange_autodiscover_url.test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_exchange_autodiscover_url.test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_exchange_autodiscover_url.test", "name", "tf-test mjx-exchange-autodiscover-url full"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_exchange_autodiscover_url.test", "description", "Test Exchange Autodiscover URL description"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_exchange_autodiscover_url.test", "url", "https://autodiscover-full.example.com/autodiscover/autodiscover.xml"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_exchange_autodiscover_url.test", "exchange_deployment", "/api/admin/configuration/v1/mjx_exchange_deployment/2/"),
				),
			},
			// Step 2: Update with min config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_mjx_exchange_autodiscover_url_min"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_exchange_autodiscover_url.test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_exchange_autodiscover_url.test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_exchange_autodiscover_url.test", "name", "tf-test mjx-exchange-autodiscover-url min"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_exchange_autodiscover_url.test", "description", ""),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_exchange_autodiscover_url.test", "url", "https://autodiscover.example.com/autodiscover/autodiscover.xml"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_exchange_autodiscover_url.test", "exchange_deployment", "/api/admin/configuration/v1/mjx_exchange_deployment/1/"),
				),
			},
			// Step 3: Destroy
			{
				Config:  test.LoadTestFolder(t, "resource_infinity_mjx_exchange_autodiscover_url_min"),
				Destroy: true,
			},
			// Step 4: Recreate with min config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_mjx_exchange_autodiscover_url_min"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_exchange_autodiscover_url.test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_exchange_autodiscover_url.test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_exchange_autodiscover_url.test", "name", "tf-test mjx-exchange-autodiscover-url min"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_exchange_autodiscover_url.test", "description", ""),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_exchange_autodiscover_url.test", "url", "https://autodiscover.example.com/autodiscover/autodiscover.xml"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_exchange_autodiscover_url.test", "exchange_deployment", "/api/admin/configuration/v1/mjx_exchange_deployment/1/"),
				),
			},
			// Step 5: Update to full config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_mjx_exchange_autodiscover_url_full"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_exchange_autodiscover_url.test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_exchange_autodiscover_url.test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_exchange_autodiscover_url.test", "name", "tf-test mjx-exchange-autodiscover-url full"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_exchange_autodiscover_url.test", "description", "Test Exchange Autodiscover URL description"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_exchange_autodiscover_url.test", "url", "https://autodiscover-full.example.com/autodiscover/autodiscover.xml"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_exchange_autodiscover_url.test", "exchange_deployment", "/api/admin/configuration/v1/mjx_exchange_deployment/2/"),
				),
			},
		},
	})
}
