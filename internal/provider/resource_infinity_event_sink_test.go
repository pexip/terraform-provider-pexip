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

	// Mock the CreateEventSink API call
	createResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/event_sink/123/",
	}
	client.On("PostWithResponse", mock.Anything, "configuration/v1/event_sink/", mock.Anything, mock.Anything).Return(createResponse, nil)

	// Shared state for mocking
	description := "Test Event Sink"
	username := "testuser"
	password := "testpassword"
	mockState := &config.EventSink{
		ID:                   123,
		ResourceURI:          "/api/admin/configuration/v1/event_sink/123/",
		Name:                 "test-event-sink",
		Description:          &description,
		URL:                  "https://test-event-sink.dev.pexip.network",
		Username:             &username,
		Password:             &password,
		BulkSupport:          false,
		VerifyTLSCertificate: false,
		Version:              1,
	}

	// Mock the GetEventSink API call for Read operations
	client.On("GetJSON", mock.Anything, "configuration/v1/event_sink/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		eventSink := args.Get(3).(*config.EventSink)
		*eventSink = *mockState
	}).Maybe()

	// Mock the UpdateEventSink API call
	client.On("PutJSON", mock.Anything, "configuration/v1/event_sink/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		updateRequest := args.Get(2).(*config.EventSinkUpdateRequest)
		eventSink := args.Get(3).(*config.EventSink)

		// Update mock state
		mockState.Name = updateRequest.Name
		mockState.URL = updateRequest.URL
		if updateRequest.Description != nil {
			mockState.Description = updateRequest.Description
		}
		if updateRequest.Username != nil {
			mockState.Username = updateRequest.Username
		}
		if updateRequest.Password != nil {
			mockState.Password = updateRequest.Password
		}
		if updateRequest.BulkSupport != nil {
			mockState.BulkSupport = *updateRequest.BulkSupport
		}
		if updateRequest.VerifyTLSCertificate != nil {
			mockState.VerifyTLSCertificate = *updateRequest.VerifyTLSCertificate
		}
		if updateRequest.Version != nil {
			mockState.Version = *updateRequest.Version
		}

		// Return updated state
		*eventSink = *mockState
	}).Maybe()

	// Mock the DeleteEventSink API call
	client.On("DeleteJSON", mock.Anything, mock.MatchedBy(func(path string) bool {
		return path == "configuration/v1/event_sink/123/"
	}), mock.Anything).Return(nil)

	testInfinityEventSink(t, client)
}

func testInfinityEventSink(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: getTestProtoV5ProviderFactories(client),
		Steps: []resource.TestStep{
			{
				Config: test.LoadTestFolder(t, "resource_infinity_event_sink_basic"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_event_sink.event-sink-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_event_sink.event-sink-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_event_sink.event-sink-test", "name", "test-event-sink"),
					resource.TestCheckResourceAttr("pexip_infinity_event_sink.event-sink-test", "description", "Test Event Sink"),
					resource.TestCheckResourceAttr("pexip_infinity_event_sink.event-sink-test", "url", "https://test-event-sink.dev.pexip.network"),
					resource.TestCheckResourceAttr("pexip_infinity_event_sink.event-sink-test", "username", "testuser"),
					resource.TestCheckResourceAttr("pexip_infinity_event_sink.event-sink-test", "password", "testpassword"),
				),
			},
			{
				Config: test.LoadTestFolder(t, "resource_infinity_event_sink_basic_updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_event_sink.event-sink-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_event_sink.event-sink-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_event_sink.event-sink-test", "name", "test-event-sink"),
					resource.TestCheckResourceAttr("pexip_infinity_event_sink.event-sink-test", "description", "Test Event Sink"),
					resource.TestCheckResourceAttr("pexip_infinity_event_sink.event-sink-test", "url", "https://test-event-sink.dev.pexip.network"),
					resource.TestCheckResourceAttr("pexip_infinity_event_sink.event-sink-test", "username", "testuser"),
					resource.TestCheckResourceAttr("pexip_infinity_event_sink.event-sink-test", "password", "updatedpassword"),
				),
			},
		},
	})
}
