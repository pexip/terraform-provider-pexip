package provider

import (
	"github.com/pexip/go-infinity-sdk/v38/config"
	"github.com/pexip/go-infinity-sdk/v38/types"
	"github.com/stretchr/testify/mock"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/pexip/go-infinity-sdk/v38"
	"github.com/pexip/terraform-provider-pexip/internal/test"
)

func TestInfinitySMTPServer(t *testing.T) {
	_ = os.Setenv("TF_ACC", "1")

	// Create a mock client and set up expectations
	client := infinity.NewClientMock()

	// Mock the CreateSmtpserver API call
	createResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/smtp_server/123/",
	}
	client.On("PostWithResponse", mock.Anything, "configuration/v1/smtp_server/", mock.Anything, mock.Anything).Return(createResponse, nil)

	// Shared state for mocking
	mockState := &config.SMTPServer{
		ID:                 123,
		ResourceURI:        "/api/admin/configuration/v1/smtp_server/123/",
		Name:               "smtp_server-test",
		Description:        "Test SMTPServer",
		Address:            "test-server.example.com",
		Port:               587,
		Username:           "smtp_server-test",
		Password:           "test-value",
		FromEmailAddress:   "test@example.com",
		ConnectionSecurity: "starttls",
	}

	// Mock the GetSmtpserver API call for Read operations
	client.On("GetJSON", mock.Anything, "configuration/v1/smtp_server/123/", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		smtp_server := args.Get(2).(*config.SMTPServer)
		*smtp_server = *mockState
	}).Maybe()

	// Mock the UpdateSmtpserver API call
	client.On("PutJSON", mock.Anything, "configuration/v1/smtp_server/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		updateReq := args.Get(2).(*config.SMTPServerUpdateRequest)
		smtp_server := args.Get(3).(*config.SMTPServer)

		// Update mock state based on request
		if updateReq.Description != "" {
			mockState.Description = updateReq.Description
		}
		if updateReq.Address != "" {
			mockState.Address = updateReq.Address
		}
		if updateReq.Port != nil {
			mockState.Port = *updateReq.Port
		}
		if updateReq.Username != "" {
			mockState.Username = updateReq.Username
		}
		if updateReq.Password != "" {
			mockState.Password = updateReq.Password
		}
		if updateReq.FromEmailAddress != "" {
			mockState.FromEmailAddress = updateReq.FromEmailAddress
		}
		if updateReq.ConnectionSecurity != "" {
			mockState.ConnectionSecurity = updateReq.ConnectionSecurity
		}

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
		ProtoV5ProviderFactories: getTestProtoV5ProviderFactories(client),
		Steps: []resource.TestStep{
			{
				Config: test.LoadTestFolder(t, "resource_infinity_smtp_server_basic"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_smtp_server.smtp_server-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_smtp_server.smtp_server-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_smtp_server.smtp_server-test", "name", "smtp_server-test"),
					resource.TestCheckResourceAttr("pexip_infinity_smtp_server.smtp_server-test", "description", "Test SMTPServer"),
					resource.TestCheckResourceAttr("pexip_infinity_smtp_server.smtp_server-test", "username", "smtp_server-test"),
				),
			},
			{
				Config: test.LoadTestFolder(t, "resource_infinity_smtp_server_basic_updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_smtp_server.smtp_server-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_smtp_server.smtp_server-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_smtp_server.smtp_server-test", "name", "smtp_server-test"),
					resource.TestCheckResourceAttr("pexip_infinity_smtp_server.smtp_server-test", "description", "Updated Test SMTPServer"),
					resource.TestCheckResourceAttr("pexip_infinity_smtp_server.smtp_server-test", "address", "updated-server.example.com"),
					resource.TestCheckResourceAttr("pexip_infinity_smtp_server.smtp_server-test", "port", "465"),
					resource.TestCheckResourceAttr("pexip_infinity_smtp_server.smtp_server-test", "username", "smtp_server-test"),
					resource.TestCheckResourceAttr("pexip_infinity_smtp_server.smtp_server-test", "password", "updated-value"),
					resource.TestCheckResourceAttr("pexip_infinity_smtp_server.smtp_server-test", "from_email_address", "updated@example.com"),
					resource.TestCheckResourceAttr("pexip_infinity_smtp_server.smtp_server-test", "connection_security", "ssl_tls"),
				),
			},
		},
	})
}
