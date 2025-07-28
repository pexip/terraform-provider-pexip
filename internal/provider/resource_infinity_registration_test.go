package provider

import (
	"github.com/pexip/go-infinity-sdk/v38/config"
	"github.com/stretchr/testify/mock"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/pexip/go-infinity-sdk/v38"
	"github.com/pexip/terraform-provider-pexip/internal/test"
)

func TestInfinityRegistration(t *testing.T) {
	_ = os.Setenv("TF_ACC", "1")

	// Create a mock client and set up expectations
	client := infinity.NewClientMock()

	// Registration is a singleton resource with ID 1
	// Mock the UpdateRegistration API call (Registration uses PUT for create/update)
	mockState := &config.Registration{
		ID:                         1,
		ResourceURI:                "/api/admin/configuration/v1/registration/1/",
		Enable:                     true,
		RefreshStrategy:            "adaptive",
		AdaptiveMinRefresh:         300,
		AdaptiveMaxRefresh:         600,
		MaximumMinRefresh:          0,
		MaximumMaxRefresh:          0,
		NattedMinRefresh:           0,
		NattedMaxRefresh:           0,
		RouteViaRegistrar:          true,
		EnablePushNotifications:    true,
		EnableGoogleCloudMessaging: true,
		PushToken:                  "test-value",
	}

	client.On("PutJSON", mock.Anything, "configuration/v1/registration/1/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		updateRequest := args.Get(2).(*config.RegistrationUpdateRequest)
		registration := args.Get(3).(*config.Registration)

		// Update mock state based on request
		if updateRequest.Enable != nil {
			mockState.Enable = *updateRequest.Enable
		}
		if updateRequest.RefreshStrategy != "" {
			mockState.RefreshStrategy = updateRequest.RefreshStrategy
			// Update refresh-related fields based on strategy
			if updateRequest.RefreshStrategy == "maximum" {
				if updateRequest.MaximumMinRefresh != nil {
					mockState.MaximumMinRefresh = *updateRequest.MaximumMinRefresh
				}
				if updateRequest.MaximumMaxRefresh != nil {
					mockState.MaximumMaxRefresh = *updateRequest.MaximumMaxRefresh
				}
				// Reset adaptive fields when switching to maximum
				mockState.AdaptiveMinRefresh = 0
				mockState.AdaptiveMaxRefresh = 0
			} else if updateRequest.RefreshStrategy == "adaptive" {
				if updateRequest.AdaptiveMinRefresh != nil {
					mockState.AdaptiveMinRefresh = *updateRequest.AdaptiveMinRefresh
				}
				if updateRequest.AdaptiveMaxRefresh != nil {
					mockState.AdaptiveMaxRefresh = *updateRequest.AdaptiveMaxRefresh
				}
				// Reset maximum fields when switching to adaptive
				mockState.MaximumMinRefresh = 0
				mockState.MaximumMaxRefresh = 0
			}
		}
		if updateRequest.NattedMinRefresh != nil {
			mockState.NattedMinRefresh = *updateRequest.NattedMinRefresh
		}
		if updateRequest.NattedMaxRefresh != nil {
			mockState.NattedMaxRefresh = *updateRequest.NattedMaxRefresh
		}
		if updateRequest.RouteViaRegistrar != nil {
			mockState.RouteViaRegistrar = *updateRequest.RouteViaRegistrar
		}
		if updateRequest.EnablePushNotifications != nil {
			mockState.EnablePushNotifications = *updateRequest.EnablePushNotifications
		}
		if updateRequest.EnableGoogleCloudMessaging != nil {
			mockState.EnableGoogleCloudMessaging = *updateRequest.EnableGoogleCloudMessaging
		}
		if updateRequest.PushToken != "" {
			mockState.PushToken = updateRequest.PushToken
		}

		// Return updated state
		*registration = *mockState
	}).Maybe()

	// Mock the GetRegistration API call for Read operations
	client.On("GetJSON", mock.Anything, "configuration/v1/registration/1/", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		registration := args.Get(2).(*config.Registration)
		*registration = *mockState
	}).Maybe()

	// Mock the DeleteRegistration API call
	client.On("DeleteJSON", mock.Anything, mock.MatchedBy(func(path string) bool {
		return path == "configuration/v1/registration/1/"
	}), mock.Anything).Return(nil)

	testInfinityRegistration(t, client)
}

func testInfinityRegistration(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: getTestProtoV5ProviderFactories(client),
		Steps: []resource.TestStep{
			{
				Config: test.LoadTestFolder(t, "resource_infinity_registration_basic"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_registration.registration-test", "id"),
					resource.TestCheckResourceAttr("pexip_infinity_registration.registration-test", "enable", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_registration.registration-test", "refresh_strategy", "adaptive"),
					resource.TestCheckResourceAttr("pexip_infinity_registration.registration-test", "route_via_registrar", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_registration.registration-test", "enable_push_notifications", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_registration.registration-test", "enable_google_cloud_messaging", "true"),
				),
			},
			{
				Config: test.LoadTestFolder(t, "resource_infinity_registration_basic_updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_registration.registration-test", "id"),
					resource.TestCheckResourceAttr("pexip_infinity_registration.registration-test", "enable", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_registration.registration-test", "refresh_strategy", "maximum"),
					resource.TestCheckResourceAttr("pexip_infinity_registration.registration-test", "route_via_registrar", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_registration.registration-test", "enable_push_notifications", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_registration.registration-test", "enable_google_cloud_messaging", "false"),
				),
			},
		},
	})
}
