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

func TestInfinitySyslogServer(t *testing.T) {
	t.Parallel()
	_ = os.Setenv("TF_ACC", "1")

	// Create a mock client and set up expectations
	client := infinity.NewClientMock()

	// Shared state for mocking - starts with full config
	mockState := &config.SyslogServer{
		ID:          123,
		ResourceURI: "/api/admin/configuration/v1/syslog_server/123/",
		Address:     "syslog.example.com",
		Port:        1514,
		Description: "tf-test syslog server description",
		Transport:   "tls",
		AuditLog:    true,
		SupportLog:  true,
		WebLog:      true,
	}

	// Mock the CreateSyslogserver API call
	createResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/syslog_server/123/",
	}
	client.On("PostWithResponse", mock.Anything, "configuration/v1/syslog_server/", mock.Anything, mock.Anything).Return(createResponse, nil).Run(func(args mock.Arguments) {
		createReq := args.Get(2).(*config.SyslogServerCreateRequest)
		// Update mock state based on create request
		mockState.Address = createReq.Address
		mockState.Description = createReq.Description
		mockState.Port = createReq.Port
		mockState.Transport = createReq.Transport
		mockState.AuditLog = createReq.AuditLog
		mockState.SupportLog = createReq.SupportLog
		mockState.WebLog = createReq.WebLog
	}).Maybe()

	// Mock the GetSyslogserver API call for Read operations
	client.On("GetJSON", mock.Anything, "configuration/v1/syslog_server/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		syslog_server := args.Get(3).(*config.SyslogServer)
		*syslog_server = *mockState
	}).Maybe()

	// Mock the UpdateSyslogserver API call
	client.On("PutJSON", mock.Anything, "configuration/v1/syslog_server/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		updateReq := args.Get(2).(*config.SyslogServerUpdateRequest)
		syslog_server := args.Get(3).(*config.SyslogServer)

		// Update mock state based on request
		if updateReq.Address != "" {
			mockState.Address = updateReq.Address
		}
		if updateReq.Port != 0 {
			mockState.Port = updateReq.Port
		}
		// Description can be empty string, so always update
		mockState.Description = updateReq.Description
		if updateReq.Transport != "" {
			mockState.Transport = updateReq.Transport
		}
		if updateReq.AuditLog != nil {
			mockState.AuditLog = *updateReq.AuditLog
		}
		if updateReq.SupportLog != nil {
			mockState.SupportLog = *updateReq.SupportLog
		}
		if updateReq.WebLog != nil {
			mockState.WebLog = *updateReq.WebLog
		}

		// Return updated state
		*syslog_server = *mockState
	}).Maybe()

	// Mock the DeleteSyslogserver API call
	client.On("DeleteJSON", mock.Anything, mock.MatchedBy(func(path string) bool {
		return path == "configuration/v1/syslog_server/123/"
	}), mock.Anything).Return(nil)

	testInfinitySyslogServer(t, client)
}

func testInfinitySyslogServer(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: getTestProtoV5ProviderFactories(client),
		Steps: []resource.TestStep{
			// Step 1: Create with full config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_syslog_server_full"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_syslog_server.tf-test-syslog-server", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_syslog_server.tf-test-syslog-server", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_syslog_server.tf-test-syslog-server", "address", "syslog.example.com"),
					resource.TestCheckResourceAttr("pexip_infinity_syslog_server.tf-test-syslog-server", "description", "tf-test syslog server description"),
					resource.TestCheckResourceAttr("pexip_infinity_syslog_server.tf-test-syslog-server", "port", "1514"),
					resource.TestCheckResourceAttr("pexip_infinity_syslog_server.tf-test-syslog-server", "transport", "tls"),
					resource.TestCheckResourceAttr("pexip_infinity_syslog_server.tf-test-syslog-server", "audit_log", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_syslog_server.tf-test-syslog-server", "support_log", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_syslog_server.tf-test-syslog-server", "web_log", "true"),
				),
			},
			// Step 2: Update to min config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_syslog_server_min"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_syslog_server.tf-test-syslog-server", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_syslog_server.tf-test-syslog-server", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_syslog_server.tf-test-syslog-server", "address", "syslog.example.com"),
					resource.TestCheckResourceAttr("pexip_infinity_syslog_server.tf-test-syslog-server", "description", ""),
					resource.TestCheckResourceAttr("pexip_infinity_syslog_server.tf-test-syslog-server", "port", "514"),
					resource.TestCheckResourceAttr("pexip_infinity_syslog_server.tf-test-syslog-server", "transport", "udp"),
					resource.TestCheckResourceAttr("pexip_infinity_syslog_server.tf-test-syslog-server", "audit_log", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_syslog_server.tf-test-syslog-server", "support_log", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_syslog_server.tf-test-syslog-server", "web_log", "false"),
				),
			},
			// Step 3: Destroy
			{
				Config:  test.LoadTestFolder(t, "resource_infinity_syslog_server_min"),
				Destroy: true,
			},
			// Step 4: Create with min config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_syslog_server_min"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_syslog_server.tf-test-syslog-server", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_syslog_server.tf-test-syslog-server", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_syslog_server.tf-test-syslog-server", "address", "syslog.example.com"),
					resource.TestCheckResourceAttr("pexip_infinity_syslog_server.tf-test-syslog-server", "description", ""),
					resource.TestCheckResourceAttr("pexip_infinity_syslog_server.tf-test-syslog-server", "port", "514"),
					resource.TestCheckResourceAttr("pexip_infinity_syslog_server.tf-test-syslog-server", "transport", "udp"),
					resource.TestCheckResourceAttr("pexip_infinity_syslog_server.tf-test-syslog-server", "audit_log", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_syslog_server.tf-test-syslog-server", "support_log", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_syslog_server.tf-test-syslog-server", "web_log", "false"),
				),
			},
			// Step 5: Update to full config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_syslog_server_full"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_syslog_server.tf-test-syslog-server", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_syslog_server.tf-test-syslog-server", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_syslog_server.tf-test-syslog-server", "address", "syslog.example.com"),
					resource.TestCheckResourceAttr("pexip_infinity_syslog_server.tf-test-syslog-server", "description", "tf-test syslog server description"),
					resource.TestCheckResourceAttr("pexip_infinity_syslog_server.tf-test-syslog-server", "port", "1514"),
					resource.TestCheckResourceAttr("pexip_infinity_syslog_server.tf-test-syslog-server", "transport", "tls"),
					resource.TestCheckResourceAttr("pexip_infinity_syslog_server.tf-test-syslog-server", "audit_log", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_syslog_server.tf-test-syslog-server", "support_log", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_syslog_server.tf-test-syslog-server", "web_log", "true"),
				),
			},
		},
	})
}
