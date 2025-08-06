package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/pexip/go-infinity-sdk/v38"

	"github.com/pexip/terraform-provider-pexip/internal/test"
)

func TestInfinitySSHPasswordHash(t *testing.T) {
	t.Parallel()
	_ = os.Setenv("TF_ACC", "1")

	// Create a mock client and set up expectations
	client := infinity.NewClientMock()

	testInfinitySSHPasswordHash(t, client)
}

func testInfinitySSHPasswordHash(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: getTestProtoV5ProviderFactories(client),
		Steps: []resource.TestStep{
			{
				Config: test.LoadTestFolder(t, "resource_infinity_ssh_password_hash_basic"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_ssh_password_hash.ssh_password_hash-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_ssh_password_hash.ssh_password_hash-test", "hash"),
					resource.TestCheckResourceAttr("pexip_infinity_ssh_password_hash.ssh_password_hash-test", "rounds", "5000"),
					resource.TestCheckResourceAttr("pexip_infinity_ssh_password_hash.ssh_password_hash-test", "salt", "abcdefghijklmnop"),
				),
			},
			{
				Config: test.LoadTestFolder(t, "resource_infinity_ssh_password_hash_basic_updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_ssh_password_hash.ssh_password_hash-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_ssh_password_hash.ssh_password_hash-test", "hash"),
					resource.TestCheckResourceAttr("pexip_infinity_ssh_password_hash.ssh_password_hash-test", "rounds", "6000"),
					resource.TestCheckResourceAttr("pexip_infinity_ssh_password_hash.ssh_password_hash-test", "salt", "qrstuvwxyzabcdef"),
				),
			},
		},
	})
}
