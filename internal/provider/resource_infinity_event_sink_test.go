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

func TestInfinityEventSink(t *testing.T) {
	t.Parallel()
	_ = os.Setenv("TF_ACC", "1")

	// Create a mock client and set up expectations
	client := infinity.NewClientMock()

	// Shared state for mocking
	mockState := &config.EventSink{
		ID:                   1,
		ResourceURI:          "/api/admin/configuration/v1/event_sink/1/",
		Name:                 "tf-test-event-sink",
		Description:          nil,
		URL:                  "https://tf-test-webhook.example.com/events",
		Username:             nil,
		Password:             nil,
		BulkSupport:          false,
		VerifyTLSCertificate: false,
		Version:              1,
	}

	// Step 1: Create with full config
	client.On("PostWithResponse", mock.Anything, "configuration/v1/event_sink/", mock.Anything, mock.Anything).Return(&types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/event_sink/1/",
	}, nil).Run(func(args mock.Arguments) {
		req := args.Get(2).(*config.EventSinkCreateRequest)
		mockState.Name = req.Name
		mockState.Description = req.Description
		mockState.URL = req.URL
		mockState.Username = req.Username
		mockState.Password = req.Password
		mockState.BulkSupport = req.BulkSupport
		mockState.VerifyTLSCertificate = req.VerifyTLSCertificate
		mockState.Version = req.Version
	}).Once()

	// Step 2: Update to min config (clear all optional fields)
	client.On("PutJSON", mock.Anything, "configuration/v1/event_sink/1/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		req := args.Get(2).(*config.EventSinkUpdateRequest)
		mockState.Name = req.Name
		mockState.Description = req.Description
		mockState.URL = req.URL
		mockState.Username = req.Username
		mockState.Password = req.Password
		if req.BulkSupport != nil {
			mockState.BulkSupport = *req.BulkSupport
		}
		if req.VerifyTLSCertificate != nil {
			mockState.VerifyTLSCertificate = *req.VerifyTLSCertificate
		}
		if req.Version != nil {
			mockState.Version = *req.Version
		}
		if args.Get(3) != nil {
			eventSink := args.Get(3).(*config.EventSink)
			*eventSink = *mockState
		}
	}).Once()

	// Step 3: Delete
	client.On("DeleteJSON", mock.Anything, "configuration/v1/event_sink/1/", mock.Anything).Return(nil).Maybe()

	// Step 4: Recreate with min config
	client.On("PostWithResponse", mock.Anything, "configuration/v1/event_sink/", mock.Anything, mock.Anything).Return(&types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/event_sink/1/",
	}, nil).Run(func(args mock.Arguments) {
		req := args.Get(2).(*config.EventSinkCreateRequest)
		mockState.Name = req.Name
		mockState.Description = req.Description
		mockState.URL = req.URL
		mockState.Username = req.Username
		mockState.Password = req.Password
		mockState.BulkSupport = req.BulkSupport
		mockState.VerifyTLSCertificate = req.VerifyTLSCertificate
		mockState.Version = req.Version
	}).Once()

	// Step 5: Update to full config
	client.On("PutJSON", mock.Anything, "configuration/v1/event_sink/1/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		req := args.Get(2).(*config.EventSinkUpdateRequest)
		mockState.Name = req.Name
		mockState.Description = req.Description
		mockState.URL = req.URL
		mockState.Username = req.Username
		mockState.Password = req.Password
		if req.BulkSupport != nil {
			mockState.BulkSupport = *req.BulkSupport
		}
		if req.VerifyTLSCertificate != nil {
			mockState.VerifyTLSCertificate = *req.VerifyTLSCertificate
		}
		if req.Version != nil {
			mockState.Version = *req.Version
		}
		if args.Get(3) != nil {
			eventSink := args.Get(3).(*config.EventSink)
			*eventSink = *mockState
		}
	}).Once()

	// Mock Read operations (GetJSON) - used throughout all steps
	client.On("GetJSON", mock.Anything, "configuration/v1/event_sink/1/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		eventSink := args.Get(3).(*config.EventSink)
		*eventSink = *mockState
	}).Maybe()

	testInfinityEventSink(t, client)
}

func testInfinityEventSink(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: getTestProtoV6ProviderFactories(client),
		Steps: []resource.TestStep{
			{
				// Step 1: Create with full config
				Config: test.LoadTestFolder(t, "resource_infinity_event_sink_full"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_event_sink.tf-test-event-sink", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_event_sink.tf-test-event-sink", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_event_sink.tf-test-event-sink", "name", "tf-test-event-sink"),
					resource.TestCheckResourceAttr("pexip_infinity_event_sink.tf-test-event-sink", "description", "tf-test Event Sink Description"),
					resource.TestCheckResourceAttr("pexip_infinity_event_sink.tf-test-event-sink", "url", "https://tf-test-webhook.example.com/events"),
					resource.TestCheckResourceAttr("pexip_infinity_event_sink.tf-test-event-sink", "username", "tf-test-user"),
					resource.TestCheckResourceAttr("pexip_infinity_event_sink.tf-test-event-sink", "password", "tf-test-password"),
					resource.TestCheckResourceAttr("pexip_infinity_event_sink.tf-test-event-sink", "bulk_support", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_event_sink.tf-test-event-sink", "verify_tls_certificate", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_event_sink.tf-test-event-sink", "version", "2"),
				),
			},
			{
				// Step 2: Update to min config (clear all optional fields)
				Config: test.LoadTestFolder(t, "resource_infinity_event_sink_min"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_event_sink.tf-test-event-sink", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_event_sink.tf-test-event-sink", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_event_sink.tf-test-event-sink", "name", "tf-test-event-sink"),
					resource.TestCheckResourceAttr("pexip_infinity_event_sink.tf-test-event-sink", "description", ""),
					resource.TestCheckResourceAttr("pexip_infinity_event_sink.tf-test-event-sink", "url", "https://tf-test-webhook.example.com/events"),
					resource.TestCheckResourceAttr("pexip_infinity_event_sink.tf-test-event-sink", "username", ""),
					resource.TestCheckResourceAttr("pexip_infinity_event_sink.tf-test-event-sink", "password", ""),
					resource.TestCheckResourceAttr("pexip_infinity_event_sink.tf-test-event-sink", "bulk_support", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_event_sink.tf-test-event-sink", "verify_tls_certificate", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_event_sink.tf-test-event-sink", "version", "1"),
				),
			},
			{
				// Step 3: Destroy
				Config:  test.LoadTestFolder(t, "resource_infinity_event_sink_min"),
				Destroy: true,
			},
			{
				// Step 4: Create with min config (after destroy)
				Config: test.LoadTestFolder(t, "resource_infinity_event_sink_min"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_event_sink.tf-test-event-sink", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_event_sink.tf-test-event-sink", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_event_sink.tf-test-event-sink", "name", "tf-test-event-sink"),
					resource.TestCheckResourceAttr("pexip_infinity_event_sink.tf-test-event-sink", "description", ""),
					resource.TestCheckResourceAttr("pexip_infinity_event_sink.tf-test-event-sink", "url", "https://tf-test-webhook.example.com/events"),
					resource.TestCheckResourceAttr("pexip_infinity_event_sink.tf-test-event-sink", "username", ""),
					resource.TestCheckResourceAttr("pexip_infinity_event_sink.tf-test-event-sink", "password", ""),
					resource.TestCheckResourceAttr("pexip_infinity_event_sink.tf-test-event-sink", "bulk_support", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_event_sink.tf-test-event-sink", "verify_tls_certificate", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_event_sink.tf-test-event-sink", "version", "1"),
				),
			},
			{
				// Step 5: Update to full config
				Config: test.LoadTestFolder(t, "resource_infinity_event_sink_full"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_event_sink.tf-test-event-sink", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_event_sink.tf-test-event-sink", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_event_sink.tf-test-event-sink", "name", "tf-test-event-sink"),
					resource.TestCheckResourceAttr("pexip_infinity_event_sink.tf-test-event-sink", "description", "tf-test Event Sink Description"),
					resource.TestCheckResourceAttr("pexip_infinity_event_sink.tf-test-event-sink", "url", "https://tf-test-webhook.example.com/events"),
					resource.TestCheckResourceAttr("pexip_infinity_event_sink.tf-test-event-sink", "username", "tf-test-user"),
					resource.TestCheckResourceAttr("pexip_infinity_event_sink.tf-test-event-sink", "password", "tf-test-password"),
					resource.TestCheckResourceAttr("pexip_infinity_event_sink.tf-test-event-sink", "bulk_support", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_event_sink.tf-test-event-sink", "verify_tls_certificate", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_event_sink.tf-test-event-sink", "version", "2"),
				),
			},
		},
	})
}
