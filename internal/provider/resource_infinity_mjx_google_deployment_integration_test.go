//go:build integration

/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package provider

import (
	"crypto/tls"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"

	"github.com/pexip/go-infinity-sdk/v38"

	"github.com/pexip/terraform-provider-pexip/internal/test"
)

func TestInfinityMjxGoogleDeploymentIntegration(t *testing.T) {
	_ = os.Setenv("TF_ACC", "1")

	client, err := infinity.New(
		infinity.WithBaseURL(test.INFINITY_BASE_URL),
		infinity.WithBasicAuth(test.INFINITY_USERNAME, test.INFINITY_PASSWORD),
		infinity.WithMaxRetries(2),
		infinity.WithTransport(&http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, // We need this because default certificate is not trusted
				MinVersion:         tls.VersionTLS12,
			},
			MaxIdleConns:        30,
			MaxIdleConnsPerHost: 5,
			IdleConnTimeout:     60 * time.Second,
		}),
	)
	require.NoError(t, err)

	testInfinityMjxGoogleDeploymentIntegration(t, client)
}

func testInfinityMjxGoogleDeploymentIntegration(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: getTestProtoV6ProviderFactories(client),
		ExternalProviders: map[string]resource.ExternalProvider{
			"tls": {
				Source: "hashicorp/tls",
			},
		},
		Steps: []resource.TestStep{
			// Step 1: Create with full config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_mjx_google_deployment_full_integration"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_google_deployment.test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_google_deployment.test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_google_deployment.test", "name", "tf-test mjx-google-deployment full"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_google_deployment.test", "description", "Test MJX Google deployment description"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_google_deployment.test", "client_email", "test-service@my-project.iam.gserviceaccount.com"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_google_deployment.test", "client_id", "123456789012345678901"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_google_deployment.test", "use_user_consent", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_google_deployment.test", "auth_endpoint", "https://accounts.google.com/o/oauth2/v2/auth"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_google_deployment.test", "token_endpoint", "https://oauth2.googleapis.com/token"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_google_deployment.test", "redirect_uri", "https://pexip.example.com/admin/platform/mjxgoogledeployment/oauth_redirect/"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_google_deployment.test", "maximum_number_of_api_requests", "500000"),
				),
			},
			// Step 2: Update with min config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_mjx_google_deployment_min_integration"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_google_deployment.test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_google_deployment.test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_google_deployment.test", "name", "tf-test mjx-google-deployment min"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_google_deployment.test", "description", ""),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_google_deployment.test", "client_email", "test-service@my-project.iam.gserviceaccount.com"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_google_deployment.test", "use_user_consent", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_google_deployment.test", "auth_endpoint", "https://accounts.google.com/o/oauth2/v2/auth"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_google_deployment.test", "token_endpoint", "https://oauth2.googleapis.com/token"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_google_deployment.test", "redirect_uri", ""),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_google_deployment.test", "maximum_number_of_api_requests", "900000"),
				),
			},
			// Step 3: Destroy
			{
				Config:  test.LoadTestFolder(t, "resource_infinity_mjx_google_deployment_min_integration"),
				Destroy: true,
			},
			// Step 4: Recreate with min config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_mjx_google_deployment_min_integration"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_google_deployment.test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_google_deployment.test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_google_deployment.test", "name", "tf-test mjx-google-deployment min"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_google_deployment.test", "description", ""),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_google_deployment.test", "client_email", "test-service@my-project.iam.gserviceaccount.com"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_google_deployment.test", "use_user_consent", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_google_deployment.test", "auth_endpoint", "https://accounts.google.com/o/oauth2/v2/auth"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_google_deployment.test", "token_endpoint", "https://oauth2.googleapis.com/token"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_google_deployment.test", "redirect_uri", ""),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_google_deployment.test", "maximum_number_of_api_requests", "900000"),
				),
			},
			// Step 5: Update to full config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_mjx_google_deployment_full_integration"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_google_deployment.test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_google_deployment.test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_google_deployment.test", "name", "tf-test mjx-google-deployment full"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_google_deployment.test", "description", "Test MJX Google deployment description"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_google_deployment.test", "client_email", "test-service@my-project.iam.gserviceaccount.com"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_google_deployment.test", "client_id", "123456789012345678901"),
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
