/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package provider

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/pexip/go-infinity-sdk/v38/config"
	"github.com/pexip/go-infinity-sdk/v38/types"
	"github.com/stretchr/testify/mock"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/pexip/go-infinity-sdk/v38"

	"github.com/pexip/terraform-provider-pexip/internal/test"
)

func TestInfinityWebappBranding(t *testing.T) {
	t.Parallel()
	_ = os.Setenv("TF_ACC", "1")

	// Create a mock client and set up expectations
	client := infinity.NewClientMock()

	// Shared state for mocking
	mockState := &config.WebappBranding{
		ResourceURI: "/api/admin/configuration/v1/webapp_branding/123/",
		Name:        "webapp_branding-test",
		Description: "Test WebappBranding",
		WebappType:  "webapp1",
	}

	// Mock the CreateWebappBranding API call (multipart form)
	createResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/webapp_branding/123/",
	}
	client.On("PostMultipartFormWithFieldsAndResponse", mock.Anything, "configuration/v1/webapp_branding/", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(createResponse, nil).Run(func(args mock.Arguments) {
		fields := args.Get(2).(map[string]string)
		// Update mock state based on create request fields
		mockState.Name = fields["name"]
		mockState.Description = fields["description"]
		mockState.UUID = fields["uuid"]
		mockState.WebappType = fields["webapp_type"]
		mockState.IsDefault = fields["is_default"] == "True"
	})

	// Mock the GetWebappBranding API call for Read operations (uses UUID)
	client.On("GetJSON", mock.Anything, "configuration/v1/webapp_branding/12345678-1234-1234-1234-123456789012/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		webapp_branding := args.Get(3).(*config.WebappBranding)
		*webapp_branding = *mockState
	}).Maybe()
	// Also mock for updated UUID
	client.On("GetJSON", mock.Anything, "configuration/v1/webapp_branding/87654321-4321-4321-4321-210987654321/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		webapp_branding := args.Get(3).(*config.WebappBranding)
		*webapp_branding = *mockState
	}).Maybe()

	// Mock the UpdateWebappBranding API call (multipart form)
	updateResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/webapp_branding/123/",
	}
	client.On("PutMultipartFormWithFieldsAndResponse", mock.Anything, mock.MatchedBy(func(path string) bool {
		return path == "configuration/v1/webapp_branding/webapp_branding-test/" || path == "configuration/v1/webapp_branding/12345678-1234-1234-1234-123456789012/" || path == "configuration/v1/webapp_branding/87654321-4321-4321-4321-210987654321/"
	}), mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(updateResponse, nil).Run(func(args mock.Arguments) {
		fields := args.Get(2).(map[string]string)
		// Update mock state based on update request fields
		if desc, ok := fields["description"]; ok && desc != "" {
			mockState.Description = desc
		}
		if uuid, ok := fields["uuid"]; ok && uuid != "" {
			mockState.UUID = uuid
		}
		if webappType, ok := fields["webapp_type"]; ok && webappType != "" {
			mockState.WebappType = webappType
		}
		if isDefault, ok := fields["is_default"]; ok {
			mockState.IsDefault = isDefault == "True"
		}
	}).Maybe()

	// Mock the DeleteWebappBranding API call (uses UUID)
	client.On("DeleteJSON", mock.Anything, mock.MatchedBy(func(path string) bool {
		return path == "configuration/v1/webapp_branding/12345678-1234-1234-1234-123456789012/" ||
			path == "configuration/v1/webapp_branding/87654321-4321-4321-4321-210987654321/"
	}), mock.Anything).Return(nil)

	testInfinityWebappBranding(t, client)
}

func testInfinityWebappBranding(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: getTestProtoV5ProviderFactories(client),
		Steps: []resource.TestStep{
			{
				Config: test.LoadTestFolder(t, "resource_infinity_webapp_branding_basic"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_webapp_branding.webapp_branding-test", "id"),
					resource.TestCheckResourceAttr("pexip_infinity_webapp_branding.webapp_branding-test", "name", "webapp_branding-test"),
					resource.TestCheckResourceAttr("pexip_infinity_webapp_branding.webapp_branding-test", "description", "Test WebappBranding"),
					resource.TestCheckResourceAttr("pexip_infinity_webapp_branding.webapp_branding-test", "is_default", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_webapp_branding.webapp_branding-test", "branding_file", "./webapp2-brand.zip"),
				),
			},
			{
				Config: test.LoadTestFolder(t, "resource_infinity_webapp_branding_basic_updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_webapp_branding.webapp_branding-test", "id"),
					resource.TestCheckResourceAttr("pexip_infinity_webapp_branding.webapp_branding-test", "name", "webapp_branding-test"),
					resource.TestCheckResourceAttr("pexip_infinity_webapp_branding.webapp_branding-test", "description", "Updated Test WebappBranding"),
					resource.TestCheckResourceAttr("pexip_infinity_webapp_branding.webapp_branding-test", "uuid", "87654321-4321-4321-4321-210987654321"),
					resource.TestCheckResourceAttr("pexip_infinity_webapp_branding.webapp_branding-test", "webapp_type", "webapp3"),
					resource.TestCheckResourceAttr("pexip_infinity_webapp_branding.webapp_branding-test", "is_default", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_webapp_branding.webapp_branding-test", "branding_file", "./webapp3-brand.zip"),
				),
			},
		},
	})
}

func TestInfinityWebappBrandingAutoGeneratedUUID(t *testing.T) {
	t.Parallel()
	_ = os.Setenv("TF_ACC", "1")

	// Create a mock client
	client := infinity.NewClientMock()

	// Shared state for mocking
	mockState := &config.WebappBranding{
		ResourceURI: "/api/admin/configuration/v1/webapp_branding/456/",
		Name:        "webapp_branding-autogen",
		Description: "Test Auto-generated UUID",
		UUID:        "", // Will be set by the create handler
		WebappType:  "webapp2",
		IsDefault:   false,
	}

	// Mock the CreateWebappBranding API call (multipart form)
	createResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/webapp_branding/456/",
	}
	client.On("PostMultipartFormWithFieldsAndResponse", mock.Anything, "configuration/v1/webapp_branding/", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(createResponse, nil).Run(func(args mock.Arguments) {
		fields := args.Get(2).(map[string]string)
		// Capture the auto-generated UUID
		mockState.UUID = fields["uuid"]
		mockState.Name = fields["name"]
		mockState.Description = fields["description"]
		mockState.WebappType = fields["webapp_type"]
		mockState.IsDefault = fields["is_default"] == "True"
	})

	// Mock the GetWebappBranding API call for Read operations (uses UUID, which is auto-generated)
	client.On("GetJSON", mock.Anything, mock.MatchedBy(func(path string) bool {
		// Accept any UUID path for this test since UUID is auto-generated
		return regexp.MustCompile(`configuration/v1/webapp_branding/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}/`).MatchString(path)
	}), mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		webapp_branding := args.Get(3).(*config.WebappBranding)
		*webapp_branding = *mockState
	}).Maybe()

	// Mock the DeleteWebappBranding API call (uses UUID)
	client.On("DeleteJSON", mock.Anything, mock.MatchedBy(func(path string) bool {
		// Accept any UUID path for delete
		return regexp.MustCompile(`configuration/v1/webapp_branding/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}/`).MatchString(path)
	}), mock.Anything).Return(nil)

	// RFC 4122 UUID regex pattern
	uuidPattern := regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)

	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: getTestProtoV5ProviderFactories(client),
		Steps: []resource.TestStep{
			{
				Config: test.LoadTestFolder(t, "resource_infinity_webapp_branding_auto_generated_uuid"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_webapp_branding.webapp_branding-autogen", "id"),
					resource.TestCheckResourceAttr("pexip_infinity_webapp_branding.webapp_branding-autogen", "name", "webapp_branding-autogen"),
					resource.TestCheckResourceAttr("pexip_infinity_webapp_branding.webapp_branding-autogen", "description", "Test Auto-generated UUID"),
					// Check that UUID was generated and matches RFC 4122 format
					resource.TestCheckResourceAttrSet("pexip_infinity_webapp_branding.webapp_branding-autogen", "uuid"),
					resource.TestMatchResourceAttr("pexip_infinity_webapp_branding.webapp_branding-autogen", "uuid", uuidPattern),
					resource.TestCheckResourceAttr("pexip_infinity_webapp_branding.webapp_branding-autogen", "webapp_type", "webapp2"),
					resource.TestCheckResourceAttr("pexip_infinity_webapp_branding.webapp_branding-autogen", "is_default", "false"),
				),
			},
		},
	})
}

func TestInfinityWebappBrandingUserProvidedValidUUID(t *testing.T) {
	t.Parallel()
	_ = os.Setenv("TF_ACC", "1")

	// Create a mock client
	client := infinity.NewClientMock()

	testUUID := "550e8400-e29b-41d4-a716-446655440000"

	// Shared state for mocking
	mockState := &config.WebappBranding{
		ResourceURI: "/api/admin/configuration/v1/webapp_branding/789/",
		Name:        "webapp_branding-valid-uuid",
		Description: "Test Valid UUID",
		UUID:        testUUID,
		WebappType:  "webapp3",
		IsDefault:   true,
	}

	// Mock the CreateWebappBranding API call (multipart form)
	createResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/webapp_branding/789/",
	}
	client.On("PostMultipartFormWithFieldsAndResponse", mock.Anything, "configuration/v1/webapp_branding/", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(createResponse, nil).Run(func(args mock.Arguments) {
		fields := args.Get(2).(map[string]string)
		mockState.UUID = fields["uuid"]
		mockState.Name = fields["name"]
		mockState.Description = fields["description"]
		mockState.WebappType = fields["webapp_type"]
		mockState.IsDefault = fields["is_default"] == "True"
	})

	// Mock the GetWebappBranding API call for Read operations (uses UUID)
	client.On("GetJSON", mock.Anything, fmt.Sprintf("configuration/v1/webapp_branding/%s/", testUUID), mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		webapp_branding := args.Get(3).(*config.WebappBranding)
		*webapp_branding = *mockState
	}).Maybe()

	// Mock the DeleteWebappBranding API call (uses UUID)
	client.On("DeleteJSON", mock.Anything, mock.MatchedBy(func(path string) bool {
		return path == fmt.Sprintf("configuration/v1/webapp_branding/%s/", testUUID)
	}), mock.Anything).Return(nil)

	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: getTestProtoV5ProviderFactories(client),
		Steps: []resource.TestStep{
			{
				Config: test.LoadTestFolder(t, "resource_infinity_webapp_branding_user_provided_valid_uuid"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_webapp_branding.webapp_branding-valid-uuid", "id"),
					resource.TestCheckResourceAttr("pexip_infinity_webapp_branding.webapp_branding-valid-uuid", "name", "webapp_branding-valid-uuid"),
					resource.TestCheckResourceAttr("pexip_infinity_webapp_branding.webapp_branding-valid-uuid", "uuid", testUUID),
					resource.TestCheckResourceAttr("pexip_infinity_webapp_branding.webapp_branding-valid-uuid", "webapp_type", "webapp3"),
				),
			},
		},
	})
}
