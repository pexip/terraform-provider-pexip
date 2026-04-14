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

func TestInfinitySMTPServer(t *testing.T) {
	t.Parallel()
	_ = os.Setenv("TF_ACC", "1")

	// Create a mock client and set up expectations
	client := infinity.NewClientMock()

	// Shared state for mocking
	mockState := &config.SMTPServer{
		ID:                 123,
		ResourceURI:        "/api/admin/configuration/v1/smtp_server/123/",
		Name:               "smtp_server-test",
		Description:        "Test SMTPServer",
		Address:            "test-server.example.com",
		Port:               587,
		Username:           "smtp_server-test",
		Password:           "", // API returns empty string for hashed passwords
		FromEmailAddress:   "test@example.com",
		ConnectionSecurity: "NONE",
	}

	// Mock the CreateSmtpserver API call
	createResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/smtp_server/123/",
	}
	client.On("PostWithResponse", mock.Anything, "configuration/v1/smtp_server/", mock.Anything, mock.Anything).Return(createResponse, nil).Run(func(args mock.Arguments) {
		createReq := args.Get(2).(*config.SMTPServerCreateRequest)

		// Update mockState to reflect the created resource
		mockState.Name = createReq.Name
		mockState.Description = createReq.Description
		mockState.Address = createReq.Address
		mockState.Port = createReq.Port
		mockState.Username = createReq.Username
		// Don't update password - keep empty to simulate hashing
		mockState.FromEmailAddress = createReq.FromEmailAddress
		mockState.ConnectionSecurity = createReq.ConnectionSecurity
	})

	// Mock the GetSmtpserver API call for Read operations
	client.On("GetJSON", mock.Anything, "configuration/v1/smtp_server/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		smtp_server := args.Get(3).(*config.SMTPServer)
		*smtp_server = *mockState
	}).Maybe()

	// Mock the UpdateSmtpserver API call
	client.On("PutJSON", mock.Anything, "configuration/v1/smtp_server/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		updateReq := args.Get(2).(*config.SMTPServerUpdateRequest)
		smtp_server := args.Get(3).(*config.SMTPServer)

		// Update mock state based on request
		mockState.Name = updateReq.Name
		mockState.Description = updateReq.Description
		mockState.Address = updateReq.Address
		if updateReq.Port != nil {
			mockState.Port = *updateReq.Port
		}
		mockState.Username = updateReq.Username
		// Don't update password from request - keep empty to simulate hashing
		mockState.FromEmailAddress = updateReq.FromEmailAddress
		mockState.ConnectionSecurity = updateReq.ConnectionSecurity

		// Return updated state
		*smtp_server = *mockState
	}).Maybe()

	// Mock the DeleteSmtpserver API call
	client.On("DeleteJSON", mock.Anything, mock.MatchedBy(func(path string) bool {
		return path == "configuration/v1/smtp_server/123/"
	}), mock.Anything).Return(nil)

	testInfinitySMTPServer(t, client)
}

func testInfinitySMTPServer(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: getTestProtoV6ProviderFactories(client),
		Steps: []resource.TestStep{
			{
				// Step 1: Create with full config
				Config: test.LoadTestFolder(t, "resource_infinity_smtp_server_full"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_smtp_server.smtp_server-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_smtp_server.smtp_server-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_smtp_server.smtp_server-test", "name", "tf-test SMTP Server full"),
					resource.TestCheckResourceAttr("pexip_infinity_smtp_server.smtp_server-test", "description", "full Test SMTPServer"),
					resource.TestCheckResourceAttr("pexip_infinity_smtp_server.smtp_server-test", "address", "updated-server.example.com"),
					resource.TestCheckResourceAttr("pexip_infinity_smtp_server.smtp_server-test", "port", "465"),
					resource.TestCheckResourceAttr("pexip_infinity_smtp_server.smtp_server-test", "username", "smtp_server-test"),
					resource.TestCheckResourceAttrSet("pexip_infinity_smtp_server.smtp_server-test", "password"),
					resource.TestCheckResourceAttr("pexip_infinity_smtp_server.smtp_server-test", "from_email_address", "updated@example.com"),
					resource.TestCheckResourceAttr("pexip_infinity_smtp_server.smtp_server-test", "connection_security", "STARTTLS"),
				),
			},
			{
				// Step 2: Update to min config and delete
				Config: test.LoadTestFolder(t, "resource_infinity_smtp_server_min"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_smtp_server.smtp_server-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_smtp_server.smtp_server-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_smtp_server.smtp_server-test", "name", "tf-test SMTP Server min"),
					resource.TestCheckResourceAttr("pexip_infinity_smtp_server.smtp_server-test", "address", "test-server.example.com"),
					resource.TestCheckResourceAttr("pexip_infinity_smtp_server.smtp_server-test", "from_email_address", "test@example.com"),
				),
			},
			{
				// Step 3: Destroy
				Config:  test.LoadTestFolder(t, "resource_infinity_smtp_server_min"),
				Destroy: true,
			},
			{
				// Step 4: Recreate with min config
				Config: test.LoadTestFolder(t, "resource_infinity_smtp_server_min"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_smtp_server.smtp_server-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_smtp_server.smtp_server-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_smtp_server.smtp_server-test", "name", "tf-test SMTP Server min"),
					resource.TestCheckResourceAttr("pexip_infinity_smtp_server.smtp_server-test", "address", "test-server.example.com"),
					resource.TestCheckResourceAttr("pexip_infinity_smtp_server.smtp_server-test", "from_email_address", "test@example.com"),
				),
			},
			{
				// Step 5: Update to full config
				Config: test.LoadTestFolder(t, "resource_infinity_smtp_server_full"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_smtp_server.smtp_server-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_smtp_server.smtp_server-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_smtp_server.smtp_server-test", "name", "tf-test SMTP Server full"),
					resource.TestCheckResourceAttr("pexip_infinity_smtp_server.smtp_server-test", "description", "full Test SMTPServer"),
					resource.TestCheckResourceAttr("pexip_infinity_smtp_server.smtp_server-test", "address", "updated-server.example.com"),
					resource.TestCheckResourceAttr("pexip_infinity_smtp_server.smtp_server-test", "port", "465"),
					resource.TestCheckResourceAttr("pexip_infinity_smtp_server.smtp_server-test", "username", "smtp_server-test"),
					resource.TestCheckResourceAttrSet("pexip_infinity_smtp_server.smtp_server-test", "password"),
					resource.TestCheckResourceAttr("pexip_infinity_smtp_server.smtp_server-test", "from_email_address", "updated@example.com"),
					resource.TestCheckResourceAttr("pexip_infinity_smtp_server.smtp_server-test", "connection_security", "STARTTLS"),
				),
			},
		},
	})
}
