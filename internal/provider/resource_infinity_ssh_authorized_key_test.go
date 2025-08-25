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

func TestInfinitySSHAuthorizedKey(t *testing.T) {
	t.Parallel()

	// Create a mock client and set up expectations
	client := infinity.NewClientMock()

	// Mock the CreateSshauthorizedkey API call
	createResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/ssh_authorized_key/123/",
	}
	client.On("PostWithResponse", mock.Anything, "configuration/v1/ssh_authorized_key/", mock.Anything, mock.Anything).Return(createResponse, nil)

	// Shared state for mocking
	mockState := &config.SSHAuthorizedKey{
		ID:          123,
		ResourceURI: "/api/admin/configuration/v1/ssh_authorized_key/123/",
		Keytype:     "ssh-rsa",
		Key:         "test-value",
		Comment:     "test-value",
		Nodes:       []string{},
	}

	// Mock the GetSshauthorizedkey API call for Read operations
	client.On("GetJSON", mock.Anything, "configuration/v1/ssh_authorized_key/123/", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		ssh_authorized_key := args.Get(2).(*config.SSHAuthorizedKey)
		*ssh_authorized_key = *mockState
	}).Maybe()

	// Mock the UpdateSshauthorizedkey API call
	client.On("PutJSON", mock.Anything, "configuration/v1/ssh_authorized_key/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		updateReq := args.Get(2).(*config.SSHAuthorizedKeyUpdateRequest)
		ssh_authorized_key := args.Get(3).(*config.SSHAuthorizedKey)

		// Update mock state based on request
		if updateReq.Keytype != "" {
			mockState.Keytype = updateReq.Keytype
		}
		if updateReq.Key != "" {
			mockState.Key = updateReq.Key
		}
		if updateReq.Comment != "" {
			mockState.Comment = updateReq.Comment
		}
		if updateReq.Nodes != nil {
			mockState.Nodes = updateReq.Nodes
		}

		// Return updated state
		*ssh_authorized_key = *mockState
	}).Maybe()

	// Mock the DeleteSshauthorizedkey API call
	client.On("DeleteJSON", mock.Anything, mock.MatchedBy(func(path string) bool {
		return path == "configuration/v1/ssh_authorized_key/123/"
	}), mock.Anything).Return(nil)

	testInfinitySSHAuthorizedKey(t, client)
}

func testInfinitySSHAuthorizedKey(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: getTestProtoV5ProviderFactories(client),
		Steps: []resource.TestStep{
			{
				Config: test.LoadTestFolder(t, "resource_infinity_ssh_authorized_key_basic"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_ssh_authorized_key.ssh_authorized_key-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_ssh_authorized_key.ssh_authorized_key-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_ssh_authorized_key.ssh_authorized_key-test", "keytype", "ssh-rsa"),
					resource.TestCheckResourceAttr("pexip_infinity_ssh_authorized_key.ssh_authorized_key-test", "key", "test-value"),
					resource.TestCheckResourceAttr("pexip_infinity_ssh_authorized_key.ssh_authorized_key-test", "comment", "test-value"),
				),
			},
			{
				Config: test.LoadTestFolder(t, "resource_infinity_ssh_authorized_key_basic_updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_ssh_authorized_key.ssh_authorized_key-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_ssh_authorized_key.ssh_authorized_key-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_ssh_authorized_key.ssh_authorized_key-test", "keytype", "ssh-ed25519"),
					resource.TestCheckResourceAttr("pexip_infinity_ssh_authorized_key.ssh_authorized_key-test", "key", "updated-value"),
					resource.TestCheckResourceAttr("pexip_infinity_ssh_authorized_key.ssh_authorized_key-test", "comment", "updated-value"),
				),
			},
		},
	})
}
