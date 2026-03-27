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

func TestInfinityMjxGoogleDeployment(t *testing.T) {
	t.Parallel()
	_ = os.Setenv("TF_ACC", "1")

	client := infinity.NewClientMock()

	mockState := &config.MjxGoogleDeployment{}

	createResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/mjx_google_deployment/123/",
	}
	client.On("PostWithResponse", mock.Anything, "configuration/v1/mjx_google_deployment/", mock.Anything, mock.Anything).Return(createResponse, nil).Run(func(args mock.Arguments) {
		createReq := args.Get(2).(*config.MjxGoogleDeploymentCreateRequest)
		*mockState = config.MjxGoogleDeployment{
			ID:                         123,
			ResourceURI:                "/api/admin/configuration/v1/mjx_google_deployment/123/",
			Name:                       createReq.Name,
			Description:                createReq.Description,
			ClientEmail:                createReq.ClientEmail,
			ClientID:                   createReq.ClientID,
			UseUserConsent:             createReq.UseUserConsent,
			AuthEndpoint:               createReq.AuthEndpoint,
			TokenEndpoint:              createReq.TokenEndpoint,
			RedirectURI:                createReq.RedirectURI,
			RefreshToken:               createReq.RefreshToken,
			OAuthState:                 createReq.OAuthState,
			MaximumNumberOfAPIRequests: createReq.MaximumNumberOfAPIRequests,
			MjxIntegrations:            createReq.MjxIntegrations,
		}
	})

	client.On("GetJSON", mock.Anything, "configuration/v1/mjx_google_deployment/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		deployment := args.Get(3).(*config.MjxGoogleDeployment)
		*deployment = *mockState
	}).Maybe()

	client.On("PutJSON", mock.Anything, "configuration/v1/mjx_google_deployment/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		updateReq := args.Get(2).(*config.MjxGoogleDeploymentUpdateRequest)

		mockState.Name = updateReq.Name
		mockState.Description = updateReq.Description
		mockState.ClientEmail = updateReq.ClientEmail
		mockState.ClientID = updateReq.ClientID
		mockState.AuthEndpoint = updateReq.AuthEndpoint
		mockState.TokenEndpoint = updateReq.TokenEndpoint
		mockState.RedirectURI = updateReq.RedirectURI
		mockState.OAuthState = updateReq.OAuthState
		mockState.MaximumNumberOfAPIRequests = updateReq.MaximumNumberOfAPIRequests
		mockState.MjxIntegrations = updateReq.MjxIntegrations
		if updateReq.UseUserConsent != nil {
			mockState.UseUserConsent = *updateReq.UseUserConsent
		}
	}).Maybe()

	client.On("DeleteJSON", mock.Anything, "configuration/v1/mjx_google_deployment/123/", mock.Anything).Return(nil)

	testInfinityMjxGoogleDeployment(t, client)
}

func testInfinityMjxGoogleDeployment(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: getTestProtoV5ProviderFactories(client),
		Steps: []resource.TestStep{
			// Step 1: Create with full config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_mjx_google_deployment_full"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_google_deployment.test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_google_deployment.test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_google_deployment.test", "name", "tf-test mjx-google-deployment full"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_google_deployment.test", "description", "Test MJX Google deployment description"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_google_deployment.test", "client_email", "test-service@my-project.iam.gserviceaccount.com"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_google_deployment.test", "client_id", "123456789012345678901"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_google_deployment.test", "client_secret", "test-client-secret"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_google_deployment.test", "private_key", "test-private-key"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_google_deployment.test", "use_user_consent", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_google_deployment.test", "auth_endpoint", "https://accounts.google.com/o/oauth2/v2/auth"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_google_deployment.test", "token_endpoint", "https://oauth2.googleapis.com/token"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_google_deployment.test", "redirect_uri", "https://pexip.example.com/admin/platform/mjxgoogledeployment/oauth_redirect/"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_google_deployment.test", "maximum_number_of_api_requests", "500000"),
				),
			},
			// Step 2: Update with min config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_mjx_google_deployment_min"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_google_deployment.test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_google_deployment.test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_google_deployment.test", "name", "tf-test mjx-google-deployment min"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_google_deployment.test", "description", ""),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_google_deployment.test", "client_email", "test-service@my-project.iam.gserviceaccount.com"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_google_deployment.test", "private_key", "test-private-key"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_google_deployment.test", "use_user_consent", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_google_deployment.test", "auth_endpoint", "https://accounts.google.com/o/oauth2/v2/auth"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_google_deployment.test", "token_endpoint", "https://oauth2.googleapis.com/token"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_google_deployment.test", "redirect_uri", ""),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_google_deployment.test", "maximum_number_of_api_requests", "900000"),
				),
			},
			// Step 3: Destroy
			{
				Config:  test.LoadTestFolder(t, "resource_infinity_mjx_google_deployment_min"),
				Destroy: true,
			},
			// Step 4: Recreate with min config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_mjx_google_deployment_min"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_google_deployment.test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_google_deployment.test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_google_deployment.test", "name", "tf-test mjx-google-deployment min"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_google_deployment.test", "description", ""),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_google_deployment.test", "client_email", "test-service@my-project.iam.gserviceaccount.com"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_google_deployment.test", "private_key", "test-private-key"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_google_deployment.test", "use_user_consent", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_google_deployment.test", "auth_endpoint", "https://accounts.google.com/o/oauth2/v2/auth"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_google_deployment.test", "token_endpoint", "https://oauth2.googleapis.com/token"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_google_deployment.test", "redirect_uri", ""),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_google_deployment.test", "maximum_number_of_api_requests", "900000"),
				),
			},
			// Step 5: Update to full config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_mjx_google_deployment_full"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_google_deployment.test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_google_deployment.test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_google_deployment.test", "name", "tf-test mjx-google-deployment full"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_google_deployment.test", "description", "Test MJX Google deployment description"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_google_deployment.test", "client_email", "test-service@my-project.iam.gserviceaccount.com"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_google_deployment.test", "client_id", "123456789012345678901"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_google_deployment.test", "client_secret", "test-client-secret"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_google_deployment.test", "private_key", "test-private-key"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_google_deployment.test", "use_user_consent", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_google_deployment.test", "auth_endpoint", "https://accounts.google.com/o/oauth2/v2/auth"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_google_deployment.test", "token_endpoint", "https://oauth2.googleapis.com/token"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_google_deployment.test", "redirect_uri", "https://pexip.example.com/admin/platform/mjxgoogledeployment/oauth_redirect/"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_google_deployment.test", "maximum_number_of_api_requests", "500000"),
				),
			},
		},
	})
}
