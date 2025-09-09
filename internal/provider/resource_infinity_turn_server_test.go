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

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/pexip/go-infinity-sdk/v38"

	"github.com/pexip/terraform-provider-pexip/internal/test"
)

func TestInfinityTurnServer(t *testing.T) {
	t.Parallel()
	_ = os.Setenv("TF_ACC", "1")

	// Create a mock client and set up expectations
	client := infinity.NewClientMock()

	// Mock the CreateTURNServer API call
	createResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/turn_server/123/",
	}
	client.On("PostWithResponse", mock.Anything, "configuration/v1/turn_server/", mock.Anything, mock.Anything).Return(createResponse, nil)

	// Shared state for mocking
	mockState := &config.TURNServer{
		ID:            123,
		ResourceURI:   "/api/admin/configuration/v1/turn_server/123/",
		Name:          "turn-server-test",
		Description:   "Test TURN server",
		Address:       "test-turn-server.dev.pexip.network",
		Port:          func() *int { port := 8080; return &port }(),
		ServerType:    "namepsw",
		TransportType: "udp",
		Username:      "turnuser",
		SecretKey:     "turnsecretkey",
	}

	// Mock the GetTURNServer API call for Read operations
	client.On("GetJSON", mock.Anything, "configuration/v1/turn_server/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		turnServer := args.Get(3).(*config.TURNServer)
		*turnServer = *mockState
	}).Maybe()

	// Mock the UpdateTURNServer API call
	client.On("PutJSON", mock.Anything, "configuration/v1/turn_server/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		updateRequest := args.Get(2).(*config.TURNServerUpdateRequest)
		turnServer := args.Get(3).(*config.TURNServer)

		// Update mock state
		mockState.Name = updateRequest.Name
		if updateRequest.Description != "" {
			mockState.Description = updateRequest.Description
		}
		if updateRequest.Address != "" {
			mockState.Address = updateRequest.Address
		}
		if updateRequest.Port != nil {
			mockState.Port = updateRequest.Port
		}
		if updateRequest.ServerType != "" {
			mockState.ServerType = updateRequest.ServerType
		}
		if updateRequest.TransportType != "" {
			mockState.TransportType = updateRequest.TransportType
		}
		if updateRequest.Username != "" {
			mockState.Username = updateRequest.Username
		}
		if updateRequest.SecretKey != "" {
			mockState.SecretKey = updateRequest.SecretKey
		}

		// Return updated state
		*turnServer = *mockState
	}).Maybe()

	// Mock the DeleteTURNServer API call
	client.On("DeleteJSON", mock.Anything, mock.MatchedBy(func(path string) bool {
		return path == "configuration/v1/turn_server/123/"
	}), mock.Anything).Return(nil)

	testInfinityTurnServer(t, client)
}

func testInfinityTurnServer(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: getTestProtoV5ProviderFactories(client),
		Steps: []resource.TestStep{
			{
				Config: test.LoadTestFolder(t, "resource_infinity_turn_server_basic"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_turn_server.turn-server-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_turn_server.turn-server-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_turn_server.turn-server-test", "name", "turn-server-test"),
					resource.TestCheckResourceAttr("pexip_infinity_turn_server.turn-server-test", "description", "Test TURN server"),
					resource.TestCheckResourceAttr("pexip_infinity_turn_server.turn-server-test", "address", "test-turn-server.dev.pexip.network"),
					resource.TestCheckResourceAttr("pexip_infinity_turn_server.turn-server-test", "port", "8080"),
					resource.TestCheckResourceAttr("pexip_infinity_turn_server.turn-server-test", "server_type", "namepsw"),
					resource.TestCheckResourceAttr("pexip_infinity_turn_server.turn-server-test", "transport_type", "udp"),
					resource.TestCheckResourceAttr("pexip_infinity_turn_server.turn-server-test", "username", "turnuser"),
					resource.TestCheckResourceAttr("pexip_infinity_turn_server.turn-server-test", "password", "turnpassword"),
					resource.TestCheckResourceAttr("pexip_infinity_turn_server.turn-server-test", "secret_key", "turnsecretkey"),
				),
			},
			{
				Config: test.LoadTestFolder(t, "resource_infinity_turn_server_basic_updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_turn_server.turn-server-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_turn_server.turn-server-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_turn_server.turn-server-test", "name", "turn-server-test"),
					resource.TestCheckResourceAttr("pexip_infinity_turn_server.turn-server-test", "description", "Test TURN server"),
					resource.TestCheckResourceAttr("pexip_infinity_turn_server.turn-server-test", "address", "test-turn-server.dev.pexip.network"),
					resource.TestCheckResourceAttr("pexip_infinity_turn_server.turn-server-test", "port", "8081"),
					resource.TestCheckResourceAttr("pexip_infinity_turn_server.turn-server-test", "server_type", "namepsw"),
					resource.TestCheckResourceAttr("pexip_infinity_turn_server.turn-server-test", "transport_type", "udp"),
					resource.TestCheckResourceAttr("pexip_infinity_turn_server.turn-server-test", "username", "turnuser"),
					resource.TestCheckResourceAttr("pexip_infinity_turn_server.turn-server-test", "password", "updatedturnpassword"),
					resource.TestCheckResourceAttr("pexip_infinity_turn_server.turn-server-test", "secret_key", "updatedturnsecretkey"),
				),
			},
		},
	})
}
