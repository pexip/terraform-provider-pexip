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

func TestInfinitySSHAuthorizedKey(t *testing.T) {
	t.Parallel()
	_ = os.Setenv("TF_ACC", "1")

	// Create a mock client and set up expectations
	client := infinity.NewClientMock()

	// Mock the CreateSshauthorizedkey API call
	createResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/ssh_authorized_key/123/",
	}
	client.On("PostWithResponse", mock.Anything, "configuration/v1/ssh_authorized_key/", mock.Anything, mock.Anything).Return(createResponse, nil).Maybe()

	// Shared state for mocking - starts with min config
	mockState := &config.SSHAuthorizedKey{
		ID:          123,
		ResourceURI: "/api/admin/configuration/v1/ssh_authorized_key/123/",
		Keytype:     "ssh-ed25519",
		Key:         "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIKTj7PIu5ycIpVVxMYlnHmVKlhG4ALxqryNSfy59XIGf tf-test",
		Comment:     "tf-test",
		Nodes:       []string{},
	}

	// Mock the GetSshauthorizedkey API call for Read operations
	client.On("GetJSON", mock.Anything, "configuration/v1/ssh_authorized_key/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		ssh_authorized_key := args.Get(3).(*config.SSHAuthorizedKey)
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
		// Comment and Nodes always sent now (without omitempty)
		mockState.Comment = updateReq.Comment
		mockState.Nodes = updateReq.Nodes

		// Return updated state - note the API returns the full key with keytype prefix
		*ssh_authorized_key = *mockState
	}).Maybe()

	// Mock the DeleteSshauthorizedkey API call
	client.On("DeleteJSON", mock.Anything, mock.MatchedBy(func(path string) bool {
		return path == "configuration/v1/ssh_authorized_key/123/"
	}), mock.Anything).Return(nil).Maybe()

	testInfinitySSHAuthorizedKey(t, client)
}

func testInfinitySSHAuthorizedKey(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: getTestProtoV6ProviderFactories(client),
		Steps: []resource.TestStep{
			// Step 1: Create with min config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_ssh_authorized_key_min"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_ssh_authorized_key.tf-test-ssh-key", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_ssh_authorized_key.tf-test-ssh-key", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_ssh_authorized_key.tf-test-ssh-key", "keytype", "ssh-ed25519"),
					resource.TestCheckResourceAttr("pexip_infinity_ssh_authorized_key.tf-test-ssh-key", "key", "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIKTj7PIu5ycIpVVxMYlnHmVKlhG4ALxqryNSfy59XIGf tf-test"),
					resource.TestCheckResourceAttr("pexip_infinity_ssh_authorized_key.tf-test-ssh-key", "comment", "tf-test"),
					resource.TestCheckResourceAttr("pexip_infinity_ssh_authorized_key.tf-test-ssh-key", "nodes.#", "0"),
				),
			},
			// Step 2: Update to full config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_ssh_authorized_key_full"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_ssh_authorized_key.tf-test-ssh-key", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_ssh_authorized_key.tf-test-ssh-key", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_ssh_authorized_key.tf-test-ssh-key", "keytype", "ssh-rsa"),
					resource.TestCheckResourceAttr("pexip_infinity_ssh_authorized_key.tf-test-ssh-key", "key", "AAAAB3NzaC1yc2EAAAADAQABAAACAQDGWIGMczXIMRNassH/IuFPSoyryEyn3uhUqn1s3tSSDOV0b3xogwejZJZKZfUo+oFoYKLbeD70CuZCSIHOx5uZmTYk04vN8r4fX0nzEfHYSCty5ZSvPXevdxyZD+CLnTEtYxbBq4k3xIsmprRKWz70MoVXQqM9jZpR5sOc1LarW24HJhM22iVVghrDX6tsI13Kvld3QRg6Y+jh6rZnH8k3EBwqP+BndSp4ECUM+XA5OEFN4ylZSlk/VS6V9XcVnERFbA3m+qkIhx/K8dc5XmGDGO1Aayn78z2lBtdUul4YdQnUYczu6hpJa2Swasatip0CL6o3vJX344MwkU3MMzJ+ynPdOMOLqQjFgX1gNboWa5udNNdzKdLmRYd3//Fwx9ZE6lPlPrApb6C1VZNgqvFl7yz0F0eSVOJZ7iEL6WzYybbtPbrWi0kO5bYpB/muP2jficXwCqaVxG9Qj/at6ALGPAgkZWbLh0MZFlH0fQzQYxnq2aLRe0KPdgoWXOW1gU7fycR/0j28yBYX5XAI1DMvwB+6vONuEo27Ty6etwHHJWYpVzmzwoElBcqfeBRxtdAgB4Rbq+SX3kNE4J5bsWxY0D6UkUuZ0xdRAgjcRwWxcJsTwIMKSyjUWzoihtIaANQE2sX/6LuR8xI3tI7ckSpY9QZzci6W/o6PuOWeP/Njkw=="),
					resource.TestCheckResourceAttr("pexip_infinity_ssh_authorized_key.tf-test-ssh-key", "comment", "tf-test SSH Key"),
					resource.TestCheckResourceAttr("pexip_infinity_ssh_authorized_key.tf-test-ssh-key", "nodes.#", "0"),
				),
			},
		},
	})
}
