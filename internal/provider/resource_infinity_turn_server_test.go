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

func TestInfinityTurnServer(t *testing.T) {
	t.Parallel()
	_ = os.Setenv("TF_ACC", "1")

	// Create a mock client and set up expectations
	client := infinity.NewClientMock()

	// Shared state for mocking - starts with full config
	mockState := &config.TURNServer{
		ID:            123,
		ResourceURI:   "/api/admin/configuration/v1/turn_server/123/",
		Name:          "tf-test-turn-server-full",
		Description:   "tf-test TURN server description",
		Address:       "turn-full.example.com",
		Port:          func() *int { port := 5349; return &port }(),
		ServerType:    "coturn_shared",
		TransportType: "tls",
		Username:      "tf-test-username",
		SecretKey:     "tf-test-secret-key",
	}

	// Mock the CreateTURNServer API call
	createResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/turn_server/123/",
	}
	client.On("PostWithResponse", mock.Anything, "configuration/v1/turn_server/", mock.Anything, mock.Anything).Return(createResponse, nil).Run(func(args mock.Arguments) {
		createReq := args.Get(2).(*config.TURNServerCreateRequest)
		// Update mock state based on create request
		mockState.Name = createReq.Name
		mockState.Description = createReq.Description
		mockState.Address = createReq.Address
		mockState.Port = createReq.Port
		mockState.ServerType = createReq.ServerType
		mockState.TransportType = createReq.TransportType
		mockState.Username = createReq.Username
		mockState.SecretKey = createReq.SecretKey
	}).Maybe()

	// Mock the GetTURNServer API call for Read operations
	client.On("GetJSON", mock.Anything, "configuration/v1/turn_server/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		turnServer := args.Get(3).(*config.TURNServer)
		*turnServer = *mockState
	}).Maybe()

	// Mock the UpdateTURNServer API call
	client.On("PutJSON", mock.Anything, "configuration/v1/turn_server/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		updateRequest := args.Get(2).(*config.TURNServerUpdateRequest)
		turnServer := args.Get(3).(*config.TURNServer)

		// Update mock state based on request
		if updateRequest.Name != "" {
			mockState.Name = updateRequest.Name
		}
		// Description can be empty string, so always update
		mockState.Description = updateRequest.Description
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
		// Username, password, and secret_key can be empty strings, so always update
		mockState.Username = updateRequest.Username
		mockState.SecretKey = updateRequest.SecretKey

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
		ProtoV6ProviderFactories: getTestProtoV6ProviderFactories(client),
		Steps: []resource.TestStep{
			// Step 1: Create with full config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_turn_server_full"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_turn_server.tf-test-turn-server", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_turn_server.tf-test-turn-server", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_turn_server.tf-test-turn-server", "name", "tf-test-turn-server-full"),
					resource.TestCheckResourceAttr("pexip_infinity_turn_server.tf-test-turn-server", "description", "tf-test TURN server description"),
					resource.TestCheckResourceAttr("pexip_infinity_turn_server.tf-test-turn-server", "address", "turn-full.example.com"),
					resource.TestCheckResourceAttr("pexip_infinity_turn_server.tf-test-turn-server", "port", "5349"),
					resource.TestCheckResourceAttr("pexip_infinity_turn_server.tf-test-turn-server", "server_type", "coturn_shared"),
					resource.TestCheckResourceAttr("pexip_infinity_turn_server.tf-test-turn-server", "transport_type", "tls"),
					resource.TestCheckResourceAttr("pexip_infinity_turn_server.tf-test-turn-server", "username", "tf-test-username"),
					resource.TestCheckResourceAttr("pexip_infinity_turn_server.tf-test-turn-server", "password", "tf-test-password"),
					resource.TestCheckResourceAttr("pexip_infinity_turn_server.tf-test-turn-server", "secret_key", "tf-test-secret-key"),
				),
			},
			// Step 2: Update to min config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_turn_server_min"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_turn_server.tf-test-turn-server", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_turn_server.tf-test-turn-server", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_turn_server.tf-test-turn-server", "name", "tf-test-turn-server"),
					resource.TestCheckResourceAttr("pexip_infinity_turn_server.tf-test-turn-server", "description", ""),
					resource.TestCheckResourceAttr("pexip_infinity_turn_server.tf-test-turn-server", "address", "turn.example.com"),
					resource.TestCheckResourceAttr("pexip_infinity_turn_server.tf-test-turn-server", "port", "3478"),
					resource.TestCheckResourceAttr("pexip_infinity_turn_server.tf-test-turn-server", "server_type", "namepsw"),
					resource.TestCheckResourceAttr("pexip_infinity_turn_server.tf-test-turn-server", "transport_type", "udp"),
					resource.TestCheckResourceAttr("pexip_infinity_turn_server.tf-test-turn-server", "username", ""),
					resource.TestCheckResourceAttr("pexip_infinity_turn_server.tf-test-turn-server", "password", ""),
					resource.TestCheckResourceAttr("pexip_infinity_turn_server.tf-test-turn-server", "secret_key", ""),
				),
			},
			// Step 3: Destroy
			{
				Config:  test.LoadTestFolder(t, "resource_infinity_turn_server_min"),
				Destroy: true,
			},
			// Step 4: Create with min config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_turn_server_min"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_turn_server.tf-test-turn-server", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_turn_server.tf-test-turn-server", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_turn_server.tf-test-turn-server", "name", "tf-test-turn-server"),
					resource.TestCheckResourceAttr("pexip_infinity_turn_server.tf-test-turn-server", "description", ""),
					resource.TestCheckResourceAttr("pexip_infinity_turn_server.tf-test-turn-server", "address", "turn.example.com"),
					resource.TestCheckResourceAttr("pexip_infinity_turn_server.tf-test-turn-server", "port", "3478"),
					resource.TestCheckResourceAttr("pexip_infinity_turn_server.tf-test-turn-server", "server_type", "namepsw"),
					resource.TestCheckResourceAttr("pexip_infinity_turn_server.tf-test-turn-server", "transport_type", "udp"),
					resource.TestCheckResourceAttr("pexip_infinity_turn_server.tf-test-turn-server", "username", ""),
					resource.TestCheckResourceAttr("pexip_infinity_turn_server.tf-test-turn-server", "password", ""),
					resource.TestCheckResourceAttr("pexip_infinity_turn_server.tf-test-turn-server", "secret_key", ""),
				),
			},
			// Step 5: Update to full config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_turn_server_full"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_turn_server.tf-test-turn-server", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_turn_server.tf-test-turn-server", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_turn_server.tf-test-turn-server", "name", "tf-test-turn-server-full"),
					resource.TestCheckResourceAttr("pexip_infinity_turn_server.tf-test-turn-server", "description", "tf-test TURN server description"),
					resource.TestCheckResourceAttr("pexip_infinity_turn_server.tf-test-turn-server", "address", "turn-full.example.com"),
					resource.TestCheckResourceAttr("pexip_infinity_turn_server.tf-test-turn-server", "port", "5349"),
					resource.TestCheckResourceAttr("pexip_infinity_turn_server.tf-test-turn-server", "server_type", "coturn_shared"),
					resource.TestCheckResourceAttr("pexip_infinity_turn_server.tf-test-turn-server", "transport_type", "tls"),
					resource.TestCheckResourceAttr("pexip_infinity_turn_server.tf-test-turn-server", "username", "tf-test-username"),
					resource.TestCheckResourceAttr("pexip_infinity_turn_server.tf-test-turn-server", "password", "tf-test-password"),
					resource.TestCheckResourceAttr("pexip_infinity_turn_server.tf-test-turn-server", "secret_key", "tf-test-secret-key"),
				),
			},
		},
	})
}
