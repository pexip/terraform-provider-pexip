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

func TestInfinityMjxGraphDeployment(t *testing.T) {
	t.Parallel()
	_ = os.Setenv("TF_ACC", "1")

	client := infinity.NewClientMock()

	mockState := &config.MjxGraphDeployment{}

	createResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/mjx_graph_deployment/123/",
	}
	client.On("PostWithResponse", mock.Anything, "configuration/v1/mjx_graph_deployment/", mock.Anything, mock.Anything).Return(createResponse, nil).Run(func(args mock.Arguments) {
		createReq := args.Get(2).(*config.MjxGraphDeploymentCreateRequest)
		*mockState = config.MjxGraphDeployment{
			ID:              123,
			ResourceURI:     "/api/admin/configuration/v1/mjx_graph_deployment/123/",
			Name:            createReq.Name,
			Description:     createReq.Description,
			ClientID:        createReq.ClientID,
			OAuthTokenURL:   createReq.OAuthTokenURL,
			GraphAPIDomain:  createReq.GraphAPIDomain,
			RequestQuota:    createReq.RequestQuota,
			DisableProxy:    createReq.DisableProxy,
			MjxIntegrations: createReq.MjxIntegrations,
		}
	})

	client.On("GetJSON", mock.Anything, "configuration/v1/mjx_graph_deployment/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		deployment := args.Get(3).(*config.MjxGraphDeployment)
		*deployment = *mockState
	}).Maybe()

	client.On("PutJSON", mock.Anything, "configuration/v1/mjx_graph_deployment/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		updateReq := args.Get(2).(*config.MjxGraphDeploymentUpdateRequest)

		mockState.Name = updateReq.Name
		mockState.Description = updateReq.Description
		mockState.ClientID = updateReq.ClientID
		mockState.OAuthTokenURL = updateReq.OAuthTokenURL
		mockState.GraphAPIDomain = updateReq.GraphAPIDomain
		if updateReq.RequestQuota != nil {
			mockState.RequestQuota = *updateReq.RequestQuota
		}
		if updateReq.DisableProxy != nil {
			mockState.DisableProxy = *updateReq.DisableProxy
		}
		mockState.MjxIntegrations = updateReq.MjxIntegrations
	}).Maybe()

	client.On("DeleteJSON", mock.Anything, "configuration/v1/mjx_graph_deployment/123/", mock.Anything).Return(nil)

	testInfinityMjxGraphDeployment(t, client)
}

func testInfinityMjxGraphDeployment(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: getTestProtoV5ProviderFactories(client),
		Steps: []resource.TestStep{
			// Step 1: Create with full config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_mjx_graph_deployment_full"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_graph_deployment.test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_graph_deployment.test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_graph_deployment.test", "name", "tf-test mjx-graph-deployment full"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_graph_deployment.test", "description", "Test MJX Graph deployment description"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_graph_deployment.test", "client_id", "a1b2c3d4-e5f6-7890-abcd-ef1234567890"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_graph_deployment.test", "client_secret", "test-client-secret"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_graph_deployment.test", "oauth_token_url", "https://login.microsoftonline.com/updated-tenant/oauth2/v2.0/token"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_graph_deployment.test", "graph_api_domain", "graph.microsoft.com"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_graph_deployment.test", "request_quota", "500000"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_graph_deployment.test", "disable_proxy", "true"),
				),
			},
			// Step 2: Update with min config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_mjx_graph_deployment_min"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_graph_deployment.test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_graph_deployment.test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_graph_deployment.test", "name", "tf-test mjx-graph-deployment min"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_graph_deployment.test", "description", ""),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_graph_deployment.test", "client_id", "12345678-1234-1234-1234-123456789012"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_graph_deployment.test", "oauth_token_url", "https://login.microsoftonline.com/test-tenant/oauth2/v2.0/token"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_graph_deployment.test", "graph_api_domain", "graph.microsoft.com"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_graph_deployment.test", "request_quota", "1000000"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_graph_deployment.test", "disable_proxy", "false"),
				),
			},
			// Step 3: Destroy
			{
				Config:  test.LoadTestFolder(t, "resource_infinity_mjx_graph_deployment_min"),
				Destroy: true,
			},
			// Step 4: Recreate with min config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_mjx_graph_deployment_min"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_graph_deployment.test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_graph_deployment.test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_graph_deployment.test", "name", "tf-test mjx-graph-deployment min"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_graph_deployment.test", "description", ""),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_graph_deployment.test", "client_id", "12345678-1234-1234-1234-123456789012"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_graph_deployment.test", "oauth_token_url", "https://login.microsoftonline.com/test-tenant/oauth2/v2.0/token"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_graph_deployment.test", "graph_api_domain", "graph.microsoft.com"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_graph_deployment.test", "request_quota", "1000000"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_graph_deployment.test", "disable_proxy", "false"),
				),
			},
			// Step 5: Update to full config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_mjx_graph_deployment_full"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_graph_deployment.test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_graph_deployment.test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_graph_deployment.test", "name", "tf-test mjx-graph-deployment full"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_graph_deployment.test", "description", "Test MJX Graph deployment description"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_graph_deployment.test", "client_id", "a1b2c3d4-e5f6-7890-abcd-ef1234567890"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_graph_deployment.test", "client_secret", "test-client-secret"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_graph_deployment.test", "oauth_token_url", "https://login.microsoftonline.com/updated-tenant/oauth2/v2.0/token"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_graph_deployment.test", "graph_api_domain", "graph.microsoft.com"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_graph_deployment.test", "request_quota", "500000"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_graph_deployment.test", "disable_proxy", "true"),
				),
			},
		},
	})
}
