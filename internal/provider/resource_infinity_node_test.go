//go:build integration

package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/pexip/terraform-provider-pexip/internal/test"
	"testing"
)

func TestInfinityNode(t *testing.T) {
	config := test.LoadTestData(t, "infinity_node_basic.tf")
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: testProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_node.node", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_node.node", "name"),
					resource.TestCheckResourceAttrSet("pexip_infinity_node.node", "config"),
					resource.TestCheckResourceAttr("pexip_infinity_node.node", "name", "test-node-1"),
					resource.TestCheckResourceAttr("pexip_infinity_node.node", "hostname", "test-node-1"),
					resource.TestCheckResourceAttr("pexip_infinity_node.node", "address", "192.168.1.100"),
					resource.TestCheckResourceAttr("pexip_infinity_node.node", "node_type", "CONFERENCING"),
					resource.TestCheckResourceAttr("pexip_infinity_node.node", "system_location", "Test Location"),
					resource.TestCheckResourceAttr("pexip_infinity_node.node", "transcoding", "true"),
				),
			},
		},
	})
}
