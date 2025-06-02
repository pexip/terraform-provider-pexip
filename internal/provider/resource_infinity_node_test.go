package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/pexip/terraform-provider-pexip/internal/test"
	"testing"
)

func TestInfinityNode_Basic(t *testing.T) {
	resourceName := "infinity_manager.manager"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAnsiblePreCheck(t, resourceName) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testInfinityNodeDestroy,
		Steps: []resource.TestStep{
			{
				Config: test.LoadTestConfig(t, "infinity_node_basic.tf"),
				Check: resource.ComposeTestCheckFunc(
					testInfinityNodeExists("infinity_manager.manager"),
					resource.TestCheckResourceAttr("infinity_manager.manager", "name", "manager"),
					resource.TestCheckResourceAttrSet("infinity_manager.manager", "group"),
					resource.TestCheckResourceAttrSet("infinity_manager.manager", "inventory"),
					resource.TestCheckResourceAttr("infinity_manager.manager", "variables.name", "manager"),
					resource.TestCheckResourceAttr("infinity_manager.manager", "variables.role", "master"),
				),
			},
		},
	})
}

func TestInfinityNode_Update(t *testing.T) {
	resourceName := "infinity_manager.manager"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAnsiblePreCheck(t, resourceName) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testInfinityNodeDestroy,
		Steps: []resource.TestStep{
			{
				Config: test.LoadTestConfig(t, "infinity_node_basic.tf"),
				Check: resource.ComposeTestCheckFunc(
					testInfinityNodeExists("infinity_manager.manager"),
					resource.TestCheckResourceAttr("infinity_manager.manager", "name", "manager"),
					resource.TestCheckResourceAttrSet("infinity_manager.manager", "group"),
					resource.TestCheckResourceAttrSet("infinity_manager.manager", "inventory"),
					resource.TestCheckResourceAttr("infinity_manager.manager", "variables.name", "manager"),
					resource.TestCheckResourceAttr("infinity_manager.manager", "variables.role", "master"),
				),
			},
			{
				Config: test.LoadTestConfig(t, "infinity_node_basic_updated.tf"),
				Check: resource.ComposeTestCheckFunc(
					testInfinityNodeExists("infinity_manager.manager"),
					resource.TestCheckResourceAttr("infinity_manager.manager", "name", "manager-osl"),
					resource.TestCheckResourceAttrSet("infinity_manager.manager", "group"),
					resource.TestCheckResourceAttrSet("infinity_manager.manager", "inventory"),
					resource.TestCheckResourceAttr("infinity_manager.manager", "variables.name", "manager-osl"),
					resource.TestCheckResourceAttr("infinity_manager.manager", "variables.role", "master"),
				),
			},
		},
	})
}

func nodeExists(hostID string, nodeRef string) bool {
	// This function should implement the logic to check if a node exists in the inventory.
	return false
}

func testInfinityNodeDestroy(s *terraform.State) error {
	var id *string
	var nodeRef string
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "infinity_node" && rs.Primary.Attributes["name"] == "node1" || rs.Primary.Attributes["name"] == "node2" {
			id = &rs.Primary.ID
			nodeRef = rs.Primary.Attributes["ref"]
		}
	}

	if id == nil {
		return fmt.Errorf("unable to find node(s) in state'")
	}

	if nodeExists(*id, nodeRef) {
		return fmt.Errorf("node '%s' still exists", *id)
	}

	return nil
}

func testInfinityNodeExists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("not found: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no resource ID is set")
		}

		if !nodeExists(rs.Primary.ID, rs.Primary.Attributes["ref"]) {
			return fmt.Errorf("node '%s' does not exist", rs.Primary.ID)
		}
		return nil
	}
}
