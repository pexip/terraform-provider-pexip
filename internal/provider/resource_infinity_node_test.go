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
					resource.TestCheckResourceAttrSet("pexip_infinity_node.node", "config"),
				),
			},
		},
	})
}
