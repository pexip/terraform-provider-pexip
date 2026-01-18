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

	// Mock the CreateSyslogserver API call
	createResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/syslog_server/123/",
	}
	client.On("PostWithResponse", mock.Anything, "configuration/v1/syslog_server/", mock.Anything, mock.Anything).Return(createResponse, nil)

	// Shared state for mocking
	mockState := &config.SyslogServer{
		ID:          123,
		ResourceURI: "/api/admin/configuration/v1/syslog_server/123/",
		Address:     "192.168.1.50",
		Port:        514,
		Description: "Test SyslogServer",
		Transport:   "udp",
		AuditLog:    true,
		SupportLog:  true,
		WebLog:      true,
	}

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
		if updateReq.Description != "" {
			mockState.Description = updateReq.Description
		}
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
			{
				Config: test.LoadTestFolder(t, "resource_infinity_syslog_server_basic"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_syslog_server.syslog_server-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_syslog_server.syslog_server-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_syslog_server.syslog_server-test", "description", "Test SyslogServer"),
					resource.TestCheckResourceAttr("pexip_infinity_syslog_server.syslog_server-test", "audit_log", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_syslog_server.syslog_server-test", "support_log", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_syslog_server.syslog_server-test", "web_log", "true"),
				),
			},
			{
				Config: test.LoadTestFolder(t, "resource_infinity_syslog_server_basic_updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_syslog_server.syslog_server-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_syslog_server.syslog_server-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_syslog_server.syslog_server-test", "address", "10.1.1.50"),
					resource.TestCheckResourceAttr("pexip_infinity_syslog_server.syslog_server-test", "port", "1514"),
					resource.TestCheckResourceAttr("pexip_infinity_syslog_server.syslog_server-test", "description", "Updated Test SyslogServer"),
					resource.TestCheckResourceAttr("pexip_infinity_syslog_server.syslog_server-test", "transport", "tcp"),
					resource.TestCheckResourceAttr("pexip_infinity_syslog_server.syslog_server-test", "audit_log", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_syslog_server.syslog_server-test", "support_log", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_syslog_server.syslog_server-test", "web_log", "false"),
				),
			},
		},
	})
}
