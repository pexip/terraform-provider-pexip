/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package provider

import (
	"testing"

	"github.com/pexip/go-infinity-sdk/v38/config"
	"github.com/pexip/go-infinity-sdk/v38/types"
	"github.com/stretchr/testify/mock"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/pexip/go-infinity-sdk/v38"

	"github.com/pexip/terraform-provider-pexip/internal/test"
)

func TestInfinityDiagnosticGraph(t *testing.T) {
	t.Parallel()

	// Create a mock client and set up expectations
	client := infinity.NewClientMock()

	// Mock the CreateDiagnosticgraph API call
	createResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/diagnostic_graphs/123/",
	}
	client.On("PostWithResponse", mock.Anything, "configuration/v1/diagnostic_graphs/", mock.Anything, mock.Anything).Return(createResponse, nil)

	// Shared state for mocking
	mockState := &config.DiagnosticGraph{
		ID:          123,
		ResourceURI: "/api/admin/configuration/v1/diagnostic_graphs/123/",
		Title:       "Test Diagnostic Graph",
		Order:       1,
	}

	// Mock the GetDiagnosticgraph API call for Read operations
	client.On("GetJSON", mock.Anything, "configuration/v1/diagnostic_graphs/123/", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		diagnostic_graph := args.Get(2).(*config.DiagnosticGraph)
		*diagnostic_graph = *mockState
	}).Maybe()

	// Mock the UpdateDiagnosticgraph API call
	client.On("PutJSON", mock.Anything, "configuration/v1/diagnostic_graphs/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		updateRequest := args.Get(2).(*config.DiagnosticGraphUpdateRequest)
		diagnostic_graph := args.Get(3).(*config.DiagnosticGraph)

		// Update mock state based on request
		if updateRequest.Title != "" {
			mockState.Title = updateRequest.Title
		}
		if updateRequest.Order != nil {
			mockState.Order = *updateRequest.Order
		}

		// Return updated state
		*diagnostic_graph = *mockState
	}).Maybe()

	// Mock the DeleteDiagnosticgraph API call
	client.On("DeleteJSON", mock.Anything, mock.MatchedBy(func(path string) bool {
		return path == "configuration/v1/diagnostic_graphs/123/"
	}), mock.Anything).Return(nil)

	testInfinityDiagnosticGraph(t, client)
}

func testInfinityDiagnosticGraph(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: getTestProtoV5ProviderFactories(client),
		Steps: []resource.TestStep{
			{
				Config: test.LoadTestFolder(t, "resource_infinity_diagnostic_graph_basic"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_diagnostic_graph.diagnostic_graph-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_diagnostic_graph.diagnostic_graph-test", "resource_id"),
				),
			},
			{
				Config: test.LoadTestFolder(t, "resource_infinity_diagnostic_graph_basic_updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_diagnostic_graph.diagnostic_graph-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_diagnostic_graph.diagnostic_graph-test", "resource_id"),
				),
			},
		},
	})
}
