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

func TestInfinityMjxExchangeDeployment(t *testing.T) {
	t.Parallel()
	_ = os.Setenv("TF_ACC", "1")

	client := infinity.NewClientMock()

	mockState := &config.MjxExchangeDeployment{}

	createResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/mjx_exchange_deployment/123/",
	}
	client.On("PostWithResponse", mock.Anything, "configuration/v1/mjx_exchange_deployment/", mock.Anything, mock.Anything).Return(createResponse, nil).Run(func(args mock.Arguments) {
		createReq := args.Get(2).(*config.MjxExchangeDeploymentCreateRequest)
		*mockState = config.MjxExchangeDeployment{
			ID:                             123,
			ResourceURI:                    "/api/admin/configuration/v1/mjx_exchange_deployment/123/",
			Name:                           createReq.Name,
			Description:                    createReq.Description,
			ServiceAccountUsername:         createReq.ServiceAccountUsername,
			AuthenticationMethod:           createReq.AuthenticationMethod,
			EWSURL:                         createReq.EWSURL,
			DisableProxy:                   createReq.DisableProxy,
			FindItemsRequestQuota:          createReq.FindItemsRequestQuota,
			KerberosRealm:                  createReq.KerberosRealm,
			KerberosKDC:                    createReq.KerberosKDC,
			KerberosExchangeSPN:            createReq.KerberosExchangeSPN,
			KerberosAuthEveryRequest:       createReq.KerberosAuthEveryRequest,
			KerberosEnableTLS:              createReq.KerberosEnableTLS,
			KerberosKDCHTTPSProxy:          createReq.KerberosKDCHTTPSProxy,
			KerberosVerifyTLSUsingCustomCA: createReq.KerberosVerifyTLSUsingCustomCA,
			OAuthClientID:                  createReq.OAuthClientID,
			OAuthAuthEndpoint:              createReq.OAuthAuthEndpoint,
			OAuthTokenEndpoint:             createReq.OAuthTokenEndpoint,
			OAuthRedirectURI:               createReq.OAuthRedirectURI,
			OAuthState:                     createReq.OAuthState,
			MjxIntegrations:                createReq.MjxIntegrations,
		}
	})

	client.On("GetJSON", mock.Anything, "configuration/v1/mjx_exchange_deployment/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		deployment := args.Get(3).(*config.MjxExchangeDeployment)
		*deployment = *mockState
	}).Maybe()

	client.On("PutJSON", mock.Anything, "configuration/v1/mjx_exchange_deployment/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		updateReq := args.Get(2).(*config.MjxExchangeDeploymentUpdateRequest)

		mockState.Name = updateReq.Name
		mockState.Description = updateReq.Description
		mockState.ServiceAccountUsername = updateReq.ServiceAccountUsername
		mockState.AuthenticationMethod = updateReq.AuthenticationMethod
		mockState.EWSURL = updateReq.EWSURL
		mockState.KerberosRealm = updateReq.KerberosRealm
		mockState.KerberosKDC = updateReq.KerberosKDC
		mockState.KerberosExchangeSPN = updateReq.KerberosExchangeSPN
		mockState.KerberosKDCHTTPSProxy = updateReq.KerberosKDCHTTPSProxy
		mockState.OAuthAuthEndpoint = updateReq.OAuthAuthEndpoint
		mockState.OAuthTokenEndpoint = updateReq.OAuthTokenEndpoint
		mockState.OAuthRedirectURI = updateReq.OAuthRedirectURI
		mockState.OAuthClientID = updateReq.OAuthClientID
		mockState.OAuthState = updateReq.OAuthState
		mockState.MjxIntegrations = updateReq.MjxIntegrations
		if updateReq.FindItemsRequestQuota != 0 {
			mockState.FindItemsRequestQuota = updateReq.FindItemsRequestQuota
		}
		if updateReq.DisableProxy != nil {
			mockState.DisableProxy = *updateReq.DisableProxy
		}
		if updateReq.KerberosAuthEveryRequest != nil {
			mockState.KerberosAuthEveryRequest = *updateReq.KerberosAuthEveryRequest
		}
		if updateReq.KerberosEnableTLS != nil {
			mockState.KerberosEnableTLS = *updateReq.KerberosEnableTLS
		}
		if updateReq.KerberosVerifyTLSUsingCustomCA != nil {
			mockState.KerberosVerifyTLSUsingCustomCA = *updateReq.KerberosVerifyTLSUsingCustomCA
		}
	}).Maybe()

	client.On("DeleteJSON", mock.Anything, "configuration/v1/mjx_exchange_deployment/123/", mock.Anything).Return(nil)

	testInfinityMjxExchangeDeployment(t, client)
}

func testInfinityMjxExchangeDeployment(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: getTestProtoV6ProviderFactories(client),
		Steps: []resource.TestStep{
			// Step 1: Create with full config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_mjx_exchange_deployment_full"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_exchange_deployment.test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_exchange_deployment.test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_exchange_deployment.test", "name", "tf-test mjx-exchange-deployment full"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_exchange_deployment.test", "description", "Test MJX Exchange deployment description"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_exchange_deployment.test", "service_account_username", "exchange-service-full@example.com"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_exchange_deployment.test", "service_account_password", "test-password-full"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_exchange_deployment.test", "authentication_method", "NTLM"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_exchange_deployment.test", "ews_url", "https://exchange.example.com/EWS/Exchange.asmx"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_exchange_deployment.test", "disable_proxy", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_exchange_deployment.test", "find_items_request_quota", "500000"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_exchange_deployment.test", "kerberos_realm", "EXAMPLE.COM"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_exchange_deployment.test", "kerberos_kdc", "kdc.example.com"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_exchange_deployment.test", "kerberos_exchange_spn", "exchangeMDB/exchange.example.com"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_exchange_deployment.test", "kerberos_auth_every_request", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_exchange_deployment.test", "kerberos_enable_tls", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_exchange_deployment.test", "kerberos_kdc_https_proxy", "https://kdc-proxy.example.com"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_exchange_deployment.test", "kerberos_verify_tls_using_custom_ca", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_exchange_deployment.test", "oauth_client_id", "a1b2c3d4-e5f6-7890-abcd-ef1234567890"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_exchange_deployment.test", "oauth_auth_endpoint", "https://login.microsoftonline.com/tenant/oauth2/v2.0/authorize"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_exchange_deployment.test", "oauth_token_endpoint", "https://login.microsoftonline.com/tenant/oauth2/v2.0/token"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_exchange_deployment.test", "oauth_redirect_uri", "https://pexip.example.com/admin/platform/mjxexchangedeployment/oauth_redirect/"),
				),
			},
			// Step 2: Update with min config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_mjx_exchange_deployment_min"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_exchange_deployment.test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_exchange_deployment.test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_exchange_deployment.test", "name", "tf-test mjx-exchange-deployment min"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_exchange_deployment.test", "description", ""),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_exchange_deployment.test", "service_account_username", "exchange-service@example.com"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_exchange_deployment.test", "service_account_password", "test-password"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_exchange_deployment.test", "authentication_method", "BASIC"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_exchange_deployment.test", "ews_url", ""),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_exchange_deployment.test", "disable_proxy", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_exchange_deployment.test", "find_items_request_quota", "1000000"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_exchange_deployment.test", "kerberos_realm", ""),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_exchange_deployment.test", "kerberos_kdc", ""),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_exchange_deployment.test", "kerberos_exchange_spn", ""),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_exchange_deployment.test", "kerberos_auth_every_request", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_exchange_deployment.test", "kerberos_enable_tls", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_exchange_deployment.test", "kerberos_kdc_https_proxy", ""),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_exchange_deployment.test", "kerberos_verify_tls_using_custom_ca", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_exchange_deployment.test", "oauth_client_id", "12345678-1234-1234-1234-123456789012"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_exchange_deployment.test", "oauth_auth_endpoint", ""),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_exchange_deployment.test", "oauth_token_endpoint", ""),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_exchange_deployment.test", "oauth_redirect_uri", ""),
				),
			},
			// Step 3: Destroy
			{
				Config:  test.LoadTestFolder(t, "resource_infinity_mjx_exchange_deployment_min"),
				Destroy: true,
			},
			// Step 4: Recreate with min config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_mjx_exchange_deployment_min"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_exchange_deployment.test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_exchange_deployment.test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_exchange_deployment.test", "name", "tf-test mjx-exchange-deployment min"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_exchange_deployment.test", "description", ""),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_exchange_deployment.test", "service_account_username", "exchange-service@example.com"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_exchange_deployment.test", "service_account_password", "test-password"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_exchange_deployment.test", "authentication_method", "BASIC"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_exchange_deployment.test", "disable_proxy", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_exchange_deployment.test", "find_items_request_quota", "1000000"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_exchange_deployment.test", "kerberos_enable_tls", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_exchange_deployment.test", "oauth_client_id", "12345678-1234-1234-1234-123456789012"),
				),
			},
			// Step 5: Update to full config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_mjx_exchange_deployment_full"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_exchange_deployment.test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_exchange_deployment.test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_exchange_deployment.test", "name", "tf-test mjx-exchange-deployment full"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_exchange_deployment.test", "description", "Test MJX Exchange deployment description"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_exchange_deployment.test", "service_account_username", "exchange-service-full@example.com"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_exchange_deployment.test", "service_account_password", "test-password-full"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_exchange_deployment.test", "authentication_method", "NTLM"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_exchange_deployment.test", "ews_url", "https://exchange.example.com/EWS/Exchange.asmx"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_exchange_deployment.test", "disable_proxy", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_exchange_deployment.test", "find_items_request_quota", "500000"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_exchange_deployment.test", "kerberos_realm", "EXAMPLE.COM"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_exchange_deployment.test", "kerberos_kdc", "kdc.example.com"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_exchange_deployment.test", "kerberos_exchange_spn", "exchangeMDB/exchange.example.com"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_exchange_deployment.test", "kerberos_auth_every_request", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_exchange_deployment.test", "kerberos_enable_tls", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_exchange_deployment.test", "kerberos_kdc_https_proxy", "https://kdc-proxy.example.com"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_exchange_deployment.test", "kerberos_verify_tls_using_custom_ca", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_exchange_deployment.test", "oauth_client_id", "a1b2c3d4-e5f6-7890-abcd-ef1234567890"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_exchange_deployment.test", "oauth_auth_endpoint", "https://login.microsoftonline.com/tenant/oauth2/v2.0/authorize"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_exchange_deployment.test", "oauth_token_endpoint", "https://login.microsoftonline.com/tenant/oauth2/v2.0/token"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_exchange_deployment.test", "oauth_redirect_uri", "https://pexip.example.com/admin/platform/mjxexchangedeployment/oauth_redirect/"),
				),
			},
		},
	})
}
