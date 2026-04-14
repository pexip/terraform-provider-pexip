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

func TestInfinityPolicyServer(t *testing.T) {
	t.Parallel()
	_ = os.Setenv("TF_ACC", "1")

	// Create a mock client and set up expectations
	client := infinity.NewClientMock()

	// Shared state for mocking - initialize with defaults
	mockState := &config.PolicyServer{
		ID:                                  1,
		ResourceURI:                         "/api/admin/configuration/v1/policy_server/1/",
		Name:                                "tf-test-policy-server",
		Description:                         "",
		URL:                                 "",
		Username:                            "",
		Password:                            "",
		EnableServiceLookup:                 false,
		EnableParticipantLookup:             false,
		EnableRegistrationLookup:            false,
		EnableDirectoryLookup:               false,
		EnableAvatarLookup:                  false,
		EnableMediaLocationLookup:           false,
		EnableInternalServicePolicy:         false,
		EnableInternalParticipantPolicy:     false,
		EnableInternalMediaLocationPolicy:   false,
		InternalServicePolicyTemplate:       "",
		InternalParticipantPolicyTemplate:   "",
		InternalMediaLocationPolicyTemplate: "",
		PreferLocalAvatarConfiguration:      false,
	}

	// Step 1: Create with full config
	client.On("PostWithResponse", mock.Anything, "configuration/v1/policy_server/", mock.Anything, mock.Anything).Return(&types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/policy_server/1/",
	}, nil).Run(func(args mock.Arguments) {
		req := args.Get(2).(*config.PolicyServerCreateRequest)
		mockState.Name = req.Name
		mockState.Description = req.Description
		mockState.URL = req.URL
		mockState.Username = req.Username
		mockState.Password = req.Password
		mockState.EnableServiceLookup = req.EnableServiceLookup
		mockState.EnableParticipantLookup = req.EnableParticipantLookup
		mockState.EnableRegistrationLookup = req.EnableRegistrationLookup
		mockState.EnableDirectoryLookup = req.EnableDirectoryLookup
		mockState.EnableAvatarLookup = req.EnableAvatarLookup
		mockState.EnableMediaLocationLookup = req.EnableMediaLocationLookup
		mockState.EnableInternalServicePolicy = req.EnableInternalServicePolicy
		mockState.EnableInternalParticipantPolicy = req.EnableInternalParticipantPolicy
		mockState.EnableInternalMediaLocationPolicy = req.EnableInternalMediaLocationPolicy
		mockState.InternalServicePolicyTemplate = req.InternalServicePolicyTemplate
		mockState.InternalParticipantPolicyTemplate = req.InternalParticipantPolicyTemplate
		mockState.InternalMediaLocationPolicyTemplate = req.InternalMediaLocationPolicyTemplate
		mockState.PreferLocalAvatarConfiguration = req.PreferLocalAvatarConfiguration
	}).Once()

	// Step 2: Update to min config
	client.On("PutJSON", mock.Anything, "configuration/v1/policy_server/1/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		req := args.Get(2).(*config.PolicyServerUpdateRequest)
		mockState.Name = req.Name
		mockState.Description = req.Description
		mockState.URL = req.URL
		mockState.Username = req.Username
		mockState.Password = req.Password
		if req.EnableServiceLookup != nil {
			mockState.EnableServiceLookup = *req.EnableServiceLookup
		}
		if req.EnableParticipantLookup != nil {
			mockState.EnableParticipantLookup = *req.EnableParticipantLookup
		}
		if req.EnableRegistrationLookup != nil {
			mockState.EnableRegistrationLookup = *req.EnableRegistrationLookup
		}
		if req.EnableDirectoryLookup != nil {
			mockState.EnableDirectoryLookup = *req.EnableDirectoryLookup
		}
		if req.EnableAvatarLookup != nil {
			mockState.EnableAvatarLookup = *req.EnableAvatarLookup
		}
		if req.EnableMediaLocationLookup != nil {
			mockState.EnableMediaLocationLookup = *req.EnableMediaLocationLookup
		}
		if req.EnableInternalServicePolicy != nil {
			mockState.EnableInternalServicePolicy = *req.EnableInternalServicePolicy
		}
		if req.EnableInternalParticipantPolicy != nil {
			mockState.EnableInternalParticipantPolicy = *req.EnableInternalParticipantPolicy
		}
		if req.EnableInternalMediaLocationPolicy != nil {
			mockState.EnableInternalMediaLocationPolicy = *req.EnableInternalMediaLocationPolicy
		}
		mockState.InternalServicePolicyTemplate = req.InternalServicePolicyTemplate
		mockState.InternalParticipantPolicyTemplate = req.InternalParticipantPolicyTemplate
		mockState.InternalMediaLocationPolicyTemplate = req.InternalMediaLocationPolicyTemplate
		if req.PreferLocalAvatarConfiguration != nil {
			mockState.PreferLocalAvatarConfiguration = *req.PreferLocalAvatarConfiguration
		}
		if args.Get(3) != nil {
			policyServer := args.Get(3).(*config.PolicyServer)
			*policyServer = *mockState
		}
	}).Once()

	// Step 3: Delete
	client.On("DeleteJSON", mock.Anything, "configuration/v1/policy_server/1/", mock.Anything).Return(nil).Maybe()

	// Step 4: Recreate with min config
	client.On("PostWithResponse", mock.Anything, "configuration/v1/policy_server/", mock.Anything, mock.Anything).Return(&types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/policy_server/1/",
	}, nil).Run(func(args mock.Arguments) {
		req := args.Get(2).(*config.PolicyServerCreateRequest)
		mockState.Name = req.Name
		mockState.Description = req.Description
		mockState.URL = req.URL
		mockState.Username = req.Username
		mockState.Password = req.Password
		mockState.EnableServiceLookup = req.EnableServiceLookup
		mockState.EnableParticipantLookup = req.EnableParticipantLookup
		mockState.EnableRegistrationLookup = req.EnableRegistrationLookup
		mockState.EnableDirectoryLookup = req.EnableDirectoryLookup
		mockState.EnableAvatarLookup = req.EnableAvatarLookup
		mockState.EnableMediaLocationLookup = req.EnableMediaLocationLookup
		mockState.EnableInternalServicePolicy = req.EnableInternalServicePolicy
		mockState.EnableInternalParticipantPolicy = req.EnableInternalParticipantPolicy
		mockState.EnableInternalMediaLocationPolicy = req.EnableInternalMediaLocationPolicy
		mockState.InternalServicePolicyTemplate = req.InternalServicePolicyTemplate
		mockState.InternalParticipantPolicyTemplate = req.InternalParticipantPolicyTemplate
		mockState.InternalMediaLocationPolicyTemplate = req.InternalMediaLocationPolicyTemplate
		mockState.PreferLocalAvatarConfiguration = req.PreferLocalAvatarConfiguration
	}).Once()

	// Step 5: Update to full config
	client.On("PutJSON", mock.Anything, "configuration/v1/policy_server/1/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		req := args.Get(2).(*config.PolicyServerUpdateRequest)
		mockState.Name = req.Name
		mockState.Description = req.Description
		mockState.URL = req.URL
		mockState.Username = req.Username
		mockState.Password = req.Password
		if req.EnableServiceLookup != nil {
			mockState.EnableServiceLookup = *req.EnableServiceLookup
		}
		if req.EnableParticipantLookup != nil {
			mockState.EnableParticipantLookup = *req.EnableParticipantLookup
		}
		if req.EnableRegistrationLookup != nil {
			mockState.EnableRegistrationLookup = *req.EnableRegistrationLookup
		}
		if req.EnableDirectoryLookup != nil {
			mockState.EnableDirectoryLookup = *req.EnableDirectoryLookup
		}
		if req.EnableAvatarLookup != nil {
			mockState.EnableAvatarLookup = *req.EnableAvatarLookup
		}
		if req.EnableMediaLocationLookup != nil {
			mockState.EnableMediaLocationLookup = *req.EnableMediaLocationLookup
		}
		if req.EnableInternalServicePolicy != nil {
			mockState.EnableInternalServicePolicy = *req.EnableInternalServicePolicy
		}
		if req.EnableInternalParticipantPolicy != nil {
			mockState.EnableInternalParticipantPolicy = *req.EnableInternalParticipantPolicy
		}
		if req.EnableInternalMediaLocationPolicy != nil {
			mockState.EnableInternalMediaLocationPolicy = *req.EnableInternalMediaLocationPolicy
		}
		mockState.InternalServicePolicyTemplate = req.InternalServicePolicyTemplate
		mockState.InternalParticipantPolicyTemplate = req.InternalParticipantPolicyTemplate
		mockState.InternalMediaLocationPolicyTemplate = req.InternalMediaLocationPolicyTemplate
		if req.PreferLocalAvatarConfiguration != nil {
			mockState.PreferLocalAvatarConfiguration = *req.PreferLocalAvatarConfiguration
		}
		if args.Get(3) != nil {
			policyServer := args.Get(3).(*config.PolicyServer)
			*policyServer = *mockState
		}
	}).Once()

	// Mock Read operations (GetJSON) - used throughout all steps
	client.On("GetJSON", mock.Anything, "configuration/v1/policy_server/1/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		policyServer := args.Get(3).(*config.PolicyServer)
		*policyServer = *mockState
	}).Maybe()

	testInfinityPolicyServer(t, client)
}

func testInfinityPolicyServer(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: getTestProtoV6ProviderFactories(client),
		Steps: []resource.TestStep{
			{
				// Step 1: Create with full config
				Config: test.LoadTestFolder(t, "resource_infinity_policy_server_full"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_policy_server.tf-test-policy-server", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_policy_server.tf-test-policy-server", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_policy_server.tf-test-policy-server", "name", "tf-test-policy-server"),
					resource.TestCheckResourceAttr("pexip_infinity_policy_server.tf-test-policy-server", "description", "tf-test Policy Server Description"),
					resource.TestCheckResourceAttr("pexip_infinity_policy_server.tf-test-policy-server", "url", "https://policy.example.com"),
					resource.TestCheckResourceAttr("pexip_infinity_policy_server.tf-test-policy-server", "username", "tf-test-user"),
					resource.TestCheckResourceAttr("pexip_infinity_policy_server.tf-test-policy-server", "password", "tf-test-password"),
					resource.TestCheckResourceAttr("pexip_infinity_policy_server.tf-test-policy-server", "enable_service_lookup", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_policy_server.tf-test-policy-server", "enable_participant_lookup", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_policy_server.tf-test-policy-server", "enable_registration_lookup", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_policy_server.tf-test-policy-server", "enable_directory_lookup", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_policy_server.tf-test-policy-server", "enable_avatar_lookup", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_policy_server.tf-test-policy-server", "enable_media_location_lookup", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_policy_server.tf-test-policy-server", "enable_internal_service_policy", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_policy_server.tf-test-policy-server", "enable_internal_participant_policy", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_policy_server.tf-test-policy-server", "enable_internal_media_location_policy", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_policy_server.tf-test-policy-server", "prefer_local_avatar_configuration", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_policy_server.tf-test-policy-server", "internal_service_policy_template", "tf-test service template"),
					resource.TestCheckResourceAttr("pexip_infinity_policy_server.tf-test-policy-server", "internal_participant_policy_template", "tf-test participant template"),
					resource.TestCheckResourceAttr("pexip_infinity_policy_server.tf-test-policy-server", "internal_media_location_policy_template", "tf-test media location template"),
				),
			},
			{
				// Step 2: Update to min config (clear optional fields, reset to defaults)
				Config: test.LoadTestFolder(t, "resource_infinity_policy_server_min"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_policy_server.tf-test-policy-server", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_policy_server.tf-test-policy-server", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_policy_server.tf-test-policy-server", "name", "tf-test-policy-server"),
					resource.TestCheckResourceAttr("pexip_infinity_policy_server.tf-test-policy-server", "description", ""),
					resource.TestCheckResourceAttr("pexip_infinity_policy_server.tf-test-policy-server", "url", ""),
					resource.TestCheckResourceAttr("pexip_infinity_policy_server.tf-test-policy-server", "username", ""),
					resource.TestCheckResourceAttr("pexip_infinity_policy_server.tf-test-policy-server", "password", ""),
					resource.TestCheckResourceAttr("pexip_infinity_policy_server.tf-test-policy-server", "enable_service_lookup", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_policy_server.tf-test-policy-server", "enable_participant_lookup", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_policy_server.tf-test-policy-server", "enable_registration_lookup", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_policy_server.tf-test-policy-server", "enable_directory_lookup", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_policy_server.tf-test-policy-server", "enable_avatar_lookup", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_policy_server.tf-test-policy-server", "enable_media_location_lookup", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_policy_server.tf-test-policy-server", "enable_internal_service_policy", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_policy_server.tf-test-policy-server", "enable_internal_participant_policy", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_policy_server.tf-test-policy-server", "enable_internal_media_location_policy", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_policy_server.tf-test-policy-server", "prefer_local_avatar_configuration", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_policy_server.tf-test-policy-server", "internal_service_policy_template", ""),
					resource.TestCheckResourceAttr("pexip_infinity_policy_server.tf-test-policy-server", "internal_participant_policy_template", ""),
					resource.TestCheckResourceAttr("pexip_infinity_policy_server.tf-test-policy-server", "internal_media_location_policy_template", ""),
				),
			},
			{
				// Step 3: Destroy
				Config:  test.LoadTestFolder(t, "resource_infinity_policy_server_min"),
				Destroy: true,
			},
			{
				// Step 4: Create with min config (after destroy)
				Config: test.LoadTestFolder(t, "resource_infinity_policy_server_min"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_policy_server.tf-test-policy-server", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_policy_server.tf-test-policy-server", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_policy_server.tf-test-policy-server", "name", "tf-test-policy-server"),
					resource.TestCheckResourceAttr("pexip_infinity_policy_server.tf-test-policy-server", "description", ""),
					resource.TestCheckResourceAttr("pexip_infinity_policy_server.tf-test-policy-server", "url", ""),
					resource.TestCheckResourceAttr("pexip_infinity_policy_server.tf-test-policy-server", "username", ""),
					resource.TestCheckResourceAttr("pexip_infinity_policy_server.tf-test-policy-server", "password", ""),
					resource.TestCheckResourceAttr("pexip_infinity_policy_server.tf-test-policy-server", "enable_service_lookup", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_policy_server.tf-test-policy-server", "enable_participant_lookup", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_policy_server.tf-test-policy-server", "enable_registration_lookup", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_policy_server.tf-test-policy-server", "enable_directory_lookup", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_policy_server.tf-test-policy-server", "enable_avatar_lookup", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_policy_server.tf-test-policy-server", "enable_media_location_lookup", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_policy_server.tf-test-policy-server", "enable_internal_service_policy", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_policy_server.tf-test-policy-server", "enable_internal_participant_policy", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_policy_server.tf-test-policy-server", "enable_internal_media_location_policy", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_policy_server.tf-test-policy-server", "prefer_local_avatar_configuration", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_policy_server.tf-test-policy-server", "internal_service_policy_template", ""),
					resource.TestCheckResourceAttr("pexip_infinity_policy_server.tf-test-policy-server", "internal_participant_policy_template", ""),
					resource.TestCheckResourceAttr("pexip_infinity_policy_server.tf-test-policy-server", "internal_media_location_policy_template", ""),
				),
			},
			{
				// Step 5: Update to full config
				Config: test.LoadTestFolder(t, "resource_infinity_policy_server_full"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_policy_server.tf-test-policy-server", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_policy_server.tf-test-policy-server", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_policy_server.tf-test-policy-server", "name", "tf-test-policy-server"),
					resource.TestCheckResourceAttr("pexip_infinity_policy_server.tf-test-policy-server", "description", "tf-test Policy Server Description"),
					resource.TestCheckResourceAttr("pexip_infinity_policy_server.tf-test-policy-server", "url", "https://policy.example.com"),
					resource.TestCheckResourceAttr("pexip_infinity_policy_server.tf-test-policy-server", "username", "tf-test-user"),
					resource.TestCheckResourceAttr("pexip_infinity_policy_server.tf-test-policy-server", "password", "tf-test-password"),
					resource.TestCheckResourceAttr("pexip_infinity_policy_server.tf-test-policy-server", "enable_service_lookup", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_policy_server.tf-test-policy-server", "enable_participant_lookup", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_policy_server.tf-test-policy-server", "enable_registration_lookup", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_policy_server.tf-test-policy-server", "enable_directory_lookup", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_policy_server.tf-test-policy-server", "enable_avatar_lookup", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_policy_server.tf-test-policy-server", "enable_media_location_lookup", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_policy_server.tf-test-policy-server", "enable_internal_service_policy", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_policy_server.tf-test-policy-server", "enable_internal_participant_policy", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_policy_server.tf-test-policy-server", "enable_internal_media_location_policy", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_policy_server.tf-test-policy-server", "prefer_local_avatar_configuration", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_policy_server.tf-test-policy-server", "internal_service_policy_template", "tf-test service template"),
					resource.TestCheckResourceAttr("pexip_infinity_policy_server.tf-test-policy-server", "internal_participant_policy_template", "tf-test participant template"),
					resource.TestCheckResourceAttr("pexip_infinity_policy_server.tf-test-policy-server", "internal_media_location_policy_template", "tf-test media location template"),
				),
			},
		},
	})
}
