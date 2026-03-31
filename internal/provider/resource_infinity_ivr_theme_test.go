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

func TestInfinityIvrTheme(t *testing.T) {
	t.Parallel()
	_ = os.Setenv("TF_ACC", "1")

	// Create a mock client and set up expectations
	client := infinity.NewClientMock()

	// Shared state for mocking - starts with min config
	mockState := &config.IVRTheme{
		ID:          123,
		ResourceURI: "/api/admin/configuration/v1/ivr_theme/123/",
		Name:        "tf-test-ivr-theme",
	}

	// Mock the CreateIVRTheme API call (single multipart POST with file)
	createResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/ivr_theme/123/",
	}
	client.On("PostMultipartFormWithFieldsAndResponse", mock.Anything, "configuration/v1/ivr_theme/",
		mock.Anything, "package", mock.Anything, mock.Anything, mock.Anything).Return(createResponse, nil).Run(func(args mock.Arguments) {
		fields := args.Get(2).(map[string]string)
		if name, ok := fields["name"]; ok {
			mockState.Name = name
		}
	}).Maybe()

	// Mock the GetIVRTheme API call for Read operations
	client.On("GetJSON", mock.Anything, "configuration/v1/ivr_theme/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		ivr_theme := args.Get(3).(*config.IVRTheme)
		*ivr_theme = *mockState
	}).Maybe()

	// Mock the UpdateIVRTheme API call (PATCH with file)
	client.On("PatchMultipartFormWithFieldsAndResponse", mock.Anything, "configuration/v1/ivr_theme/123/",
		mock.Anything, "package", mock.Anything, mock.Anything, mock.Anything).Return(nil, nil).Run(func(args mock.Arguments) {
		fields := args.Get(2).(map[string]string)
		ivr_theme := args.Get(6).(*config.IVRTheme)
		if name, ok := fields["name"]; ok && name != "" {
			mockState.Name = name
		}
		*ivr_theme = *mockState
	}).Maybe()

	// Mock the DeleteIVRTheme API call
	client.On("DeleteJSON", mock.Anything, mock.MatchedBy(func(path string) bool {
		return path == "configuration/v1/ivr_theme/123/"
	}), mock.Anything).Return(nil)

	testInfinityIvrTheme(t, client)
}

func testInfinityIvrTheme(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: getTestProtoV5ProviderFactories(client),
		Steps: []resource.TestStep{
			// Step 1: Create with min config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_ivr_theme_min"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_ivr_theme.ivr_theme-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_ivr_theme.ivr_theme-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_ivr_theme.ivr_theme-test", "name", "tf-test-ivr-theme"),
				),
			},
			// Step 2: Update to full config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_ivr_theme_full"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_ivr_theme.ivr_theme-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_ivr_theme.ivr_theme-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_ivr_theme.ivr_theme-test", "name", "tf-test-ivr-theme-full"),
				),
			},
		},
	})
}
