/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package provider

import (
	"os"
	"regexp"
	"testing"

	"github.com/pexip/go-infinity-sdk/v38/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/pexip/go-infinity-sdk/v38"

	"github.com/pexip/terraform-provider-pexip/internal/test"
)

func TestInfinityRegistration(t *testing.T) {
	t.Parallel()
	_ = os.Setenv("TF_ACC", "1")

	client := infinity.NewClientMock()

	mockState := &config.Registration{
		ID:                         1,
		ResourceURI:                "/api/admin/configuration/v1/registration/1/",
		Enable:                     true,       // default
		RefreshStrategy:            "adaptive", // default
		AdaptiveMinRefresh:         60,         // default
		AdaptiveMaxRefresh:         3600,       // default
		MaximumMinRefresh:          60,         // default
		MaximumMaxRefresh:          300,        // default
		NattedMinRefresh:           60,         // default
		NattedMaxRefresh:           90,         // default
		RouteViaRegistrar:          true,       // default
		EnablePushNotifications:    false,      // default
		EnableGoogleCloudMessaging: true,       // default
	}

	// Delete mock — registered first so it takes priority over the general mock.
	// Fingerprinted by both AdaptiveMinRefresh == 60 AND MaximumMinRefresh == 60,
	// which only Delete sends (Create/Update only sends fields matching the active strategy).
	client.On("PatchJSON", mock.Anything, "configuration/v1/registration/1/",
		mock.MatchedBy(func(req *config.RegistrationUpdateRequest) bool {
			return req.AdaptiveMinRefresh != nil && *req.AdaptiveMinRefresh == 60 &&
				req.MaximumMinRefresh != nil && *req.MaximumMinRefresh == 60
		}), mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		req := args.Get(2).(*config.RegistrationUpdateRequest)

		assert.NotNil(t, req.Enable)
		assert.True(t, *req.Enable)
		assert.Equal(t, "adaptive", req.RefreshStrategy)
		assert.NotNil(t, req.AdaptiveMinRefresh)
		assert.Equal(t, 60, *req.AdaptiveMinRefresh)
		assert.NotNil(t, req.AdaptiveMaxRefresh)
		assert.Equal(t, 3600, *req.AdaptiveMaxRefresh)
		assert.NotNil(t, req.MaximumMinRefresh)
		assert.Equal(t, 60, *req.MaximumMinRefresh)
		assert.NotNil(t, req.MaximumMaxRefresh)
		assert.Equal(t, 300, *req.MaximumMaxRefresh)
		assert.NotNil(t, req.NattedMinRefresh)
		assert.Equal(t, 60, *req.NattedMinRefresh)
		assert.NotNil(t, req.NattedMaxRefresh)
		assert.Equal(t, 90, *req.NattedMaxRefresh)
		assert.NotNil(t, req.RouteViaRegistrar)
		assert.True(t, *req.RouteViaRegistrar)
		assert.NotNil(t, req.EnablePushNotifications)
		assert.False(t, *req.EnablePushNotifications)
		assert.NotNil(t, req.EnableGoogleCloudMessaging)
		assert.True(t, *req.EnableGoogleCloudMessaging)

		// Reset mockState to defaults
		mockState.Enable = true
		mockState.RefreshStrategy = "adaptive"
		mockState.AdaptiveMinRefresh = 60
		mockState.AdaptiveMaxRefresh = 3600
		mockState.MaximumMinRefresh = 60
		mockState.MaximumMaxRefresh = 300
		mockState.NattedMinRefresh = 60
		mockState.NattedMaxRefresh = 90
		mockState.RouteViaRegistrar = true
		mockState.EnablePushNotifications = false
		mockState.EnableGoogleCloudMessaging = true
		mockState.PushToken = ""
	}).Once()

	// General PatchJSON mock — handles all create and update calls.
	client.On("PatchJSON", mock.Anything, "configuration/v1/registration/1/",
		mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		req := args.Get(2).(*config.RegistrationUpdateRequest)
		result := args.Get(3).(*config.Registration)

		if req.Enable != nil {
			mockState.Enable = *req.Enable
		}
		if req.RefreshStrategy != "" {
			mockState.RefreshStrategy = req.RefreshStrategy
		}
		if req.AdaptiveMinRefresh != nil {
			mockState.AdaptiveMinRefresh = *req.AdaptiveMinRefresh
		}
		if req.AdaptiveMaxRefresh != nil {
			mockState.AdaptiveMaxRefresh = *req.AdaptiveMaxRefresh
		}
		if req.MaximumMinRefresh != nil {
			mockState.MaximumMinRefresh = *req.MaximumMinRefresh
		}
		if req.MaximumMaxRefresh != nil {
			mockState.MaximumMaxRefresh = *req.MaximumMaxRefresh
		}
		if req.NattedMinRefresh != nil {
			mockState.NattedMinRefresh = *req.NattedMinRefresh
		}
		if req.NattedMaxRefresh != nil {
			mockState.NattedMaxRefresh = *req.NattedMaxRefresh
		}
		if req.RouteViaRegistrar != nil {
			mockState.RouteViaRegistrar = *req.RouteViaRegistrar
		}
		if req.EnablePushNotifications != nil {
			mockState.EnablePushNotifications = *req.EnablePushNotifications
		}
		if req.EnableGoogleCloudMessaging != nil {
			mockState.EnableGoogleCloudMessaging = *req.EnableGoogleCloudMessaging
		}
		if req.PushToken != "" {
			mockState.PushToken = req.PushToken
		}

		*result = *mockState
	}).Maybe()

	// GetJSON mock — returns current mockState for all Read operations.
	client.On(
		"GetJSON",
		mock.Anything,
		"configuration/v1/registration/1/",
		mock.Anything,
		mock.AnythingOfType("*config.Registration"),
	).Return(nil).Run(func(args mock.Arguments) {
		reg := args.Get(3).(*config.Registration)
		*reg = *mockState
	}).Maybe()

	testInfinityRegistration(t, client)
}

func testInfinityRegistration(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: getTestProtoV6ProviderFactories(client),
		Steps: []resource.TestStep{
			// Step 1: Apply full configuration with all fields set to non-default values.
			{
				Config: test.LoadTestFolder(t, "resource_infinity_registration_full"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_registration.registration-test", "id"),
					resource.TestCheckResourceAttr("pexip_infinity_registration.registration-test", "enable", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_registration.registration-test", "refresh_strategy", "maximum"),
					resource.TestCheckResourceAttr("pexip_infinity_registration.registration-test", "adaptive_min_refresh", "120"),
					resource.TestCheckResourceAttr("pexip_infinity_registration.registration-test", "adaptive_max_refresh", "600"),
					resource.TestCheckResourceAttr("pexip_infinity_registration.registration-test", "maximum_min_refresh", "120"),
					resource.TestCheckResourceAttr("pexip_infinity_registration.registration-test", "maximum_max_refresh", "600"),
					resource.TestCheckResourceAttr("pexip_infinity_registration.registration-test", "natted_min_refresh", "120"),
					resource.TestCheckResourceAttr("pexip_infinity_registration.registration-test", "natted_max_refresh", "180"),
					resource.TestCheckResourceAttr("pexip_infinity_registration.registration-test", "route_via_registrar", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_registration.registration-test", "enable_push_notifications", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_registration.registration-test", "enable_google_cloud_messaging", "false"),
					resource.TestCheckResourceAttrSet("pexip_infinity_registration.registration-test", "push_token"),
				),
			},
			// Step 2: Destroy — triggers Delete which must reset all fields to API defaults.
			{
				Config:  test.LoadTestFolder(t, "resource_infinity_registration_full"),
				Destroy: true,
			},
			// Step 3: Re-apply min config and verify the API returned all fields to defaults.
			{
				Config: test.LoadTestFolder(t, "resource_infinity_registration_min"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_registration.registration-test", "id"),
					resource.TestCheckResourceAttr("pexip_infinity_registration.registration-test", "enable", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_registration.registration-test", "refresh_strategy", "adaptive"),
					resource.TestCheckResourceAttr("pexip_infinity_registration.registration-test", "adaptive_min_refresh", "60"),
					resource.TestCheckResourceAttr("pexip_infinity_registration.registration-test", "adaptive_max_refresh", "3600"),
					resource.TestCheckResourceAttr("pexip_infinity_registration.registration-test", "maximum_min_refresh", "60"),
					resource.TestCheckResourceAttr("pexip_infinity_registration.registration-test", "maximum_max_refresh", "300"),
					resource.TestCheckResourceAttr("pexip_infinity_registration.registration-test", "natted_min_refresh", "60"),
					resource.TestCheckResourceAttr("pexip_infinity_registration.registration-test", "natted_max_refresh", "90"),
					resource.TestCheckResourceAttr("pexip_infinity_registration.registration-test", "route_via_registrar", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_registration.registration-test", "enable_push_notifications", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_registration.registration-test", "enable_google_cloud_messaging", "true"),
					resource.TestCheckNoResourceAttr("pexip_infinity_registration.registration-test", "push_token"),
				),
			},
		},
	})
}

func TestInfinityRegistrationValidation(t *testing.T) {
	t.Parallel()
	_ = os.Setenv("TF_ACC", "1")

	client := infinity.NewClientMock()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: getTestProtoV6ProviderFactories(client),
		Steps: []resource.TestStep{
			{
				Config: `
resource "pexip_infinity_registration" "registration-test" {
  refresh_strategy = "INVALID"
}
`,
				ExpectError: regexp.MustCompile(`value must be one of`),
			},
			{
				Config: `
resource "pexip_infinity_registration" "registration-test" {
  adaptive_min_refresh = 59
}
`,
				ExpectError: regexp.MustCompile(`value must be between 60 and 3600`),
			},
			{
				Config: `
resource "pexip_infinity_registration" "registration-test" {
  adaptive_min_refresh = 3601
}
`,
				ExpectError: regexp.MustCompile(`value must be between 60 and 3600`),
			},
			{
				Config: `
resource "pexip_infinity_registration" "registration-test" {
  adaptive_max_refresh = 59
}
`,
				ExpectError: regexp.MustCompile(`value must be between 60 and 7200`),
			},
			{
				Config: `
resource "pexip_infinity_registration" "registration-test" {
  adaptive_max_refresh = 7201
}
`,
				ExpectError: regexp.MustCompile(`value must be between 60 and 7200`),
			},
			{
				Config: `
resource "pexip_infinity_registration" "registration-test" {
  maximum_min_refresh = 59
}
`,
				ExpectError: regexp.MustCompile(`value must be between 60 and 3600`),
			},
			{
				Config: `
resource "pexip_infinity_registration" "registration-test" {
  maximum_min_refresh = 3601
}
`,
				ExpectError: regexp.MustCompile(`value must be between 60 and 3600`),
			},
			{
				Config: `
resource "pexip_infinity_registration" "registration-test" {
  maximum_max_refresh = 59
}
`,
				ExpectError: regexp.MustCompile(`value must be between 60 and 7200`),
			},
			{
				Config: `
resource "pexip_infinity_registration" "registration-test" {
  maximum_max_refresh = 7201
}
`,
				ExpectError: regexp.MustCompile(`value must be between 60 and 7200`),
			},
			{
				Config: `
resource "pexip_infinity_registration" "registration-test" {
  natted_min_refresh = 59
}
`,
				ExpectError: regexp.MustCompile(`value must be between 60 and 3600`),
			},
			{
				Config: `
resource "pexip_infinity_registration" "registration-test" {
  natted_min_refresh = 3601
}
`,
				ExpectError: regexp.MustCompile(`value must be between 60 and 3600`),
			},
			{
				Config: `
resource "pexip_infinity_registration" "registration-test" {
  natted_max_refresh = 59
}
`,
				ExpectError: regexp.MustCompile(`value must be between 60 and 3600`),
			},
			{
				Config: `
resource "pexip_infinity_registration" "registration-test" {
  natted_max_refresh = 3601
}
`,
				ExpectError: regexp.MustCompile(`value must be between 60 and 3600`),
			},
		},
	})
}
