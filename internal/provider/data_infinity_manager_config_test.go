package provider

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/pexip/terraform-provider-pexip/internal/test"
)

func TestInfinityManagerConfig(t *testing.T) {
	os.Setenv("TF_ACC", "1")
	config := test.LoadTestData(t, "data_infinity_manager_config_basic.tf")
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: testProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.pexip_infinity_manager_config.master", "hostname"),
					resource.TestCheckResourceAttrSet("data.pexip_infinity_manager_config.master", "domain"),
					resource.TestCheckResourceAttrSet("data.pexip_infinity_manager_config.master", "ip"),
					resource.TestCheckResourceAttrSet("data.pexip_infinity_manager_config.master", "mask"),
					resource.TestCheckResourceAttrSet("data.pexip_infinity_manager_config.master", "gw"),
					resource.TestCheckResourceAttrSet("data.pexip_infinity_manager_config.master", "dns"),
					resource.TestCheckResourceAttrSet("data.pexip_infinity_manager_config.master", "ntp"),
					resource.TestCheckResourceAttrSet("data.pexip_infinity_manager_config.master", "user"),
					resource.TestCheckResourceAttrSet("data.pexip_infinity_manager_config.master", "pass"),
					resource.TestCheckResourceAttrSet("data.pexip_infinity_manager_config.master", "admin_password"),
					resource.TestCheckResourceAttrSet("data.pexip_infinity_manager_config.master", "error_reports"),
					resource.TestCheckResourceAttrSet("data.pexip_infinity_manager_config.master", "enable_analytics"),
					resource.TestCheckResourceAttrSet("data.pexip_infinity_manager_config.master", "contact_email_address"),
					testCheckPasswordHashed("data.pexip_infinity_manager_config.master"),
				),
			},
		},
	})
}

// testCheckPasswordHashed verifies that the admin_password in the rendered JSON is properly hashed
func testCheckPasswordHashed(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		rendered := rs.Primary.Attributes["rendered"]
		if rendered == "" {
			return fmt.Errorf("rendered attribute is empty")
		}

		// Parse the JSON to extract the admin_password
		var config map[string]interface{}
		if err := json.Unmarshal([]byte(rendered), &config); err != nil {
			return fmt.Errorf("failed to parse rendered JSON: %v", err)
		}

		mgmtNode, ok := config["management_node_config"].(map[string]interface{})
		if !ok {
			return fmt.Errorf("management_node_config not found in rendered JSON")
		}

		password, ok := mgmtNode["pass"].(string)
		if !ok {
			return fmt.Errorf("admin_password not found in management_node_config")
		}

		// Check if the password is hashed (should start with pbkdf2_sha256$36000$)
		if !strings.HasPrefix(password, "pbkdf2_sha256$36000$") {
			return fmt.Errorf("admin_password is not properly hashed, got: %s", password)
		}

		// Verify the hash format: pbkdf2_sha256$36000$salt$hash
		parts := strings.Split(password, "$")
		if len(parts) != 4 {
			return fmt.Errorf("admin_password hash format is incorrect, expected 4 parts, got %d", len(parts))
		}

		if parts[0] != "pbkdf2_sha256" {
			return fmt.Errorf("admin_password hash algorithm is incorrect, expected pbkdf2_sha256, got %s", parts[0])
		}

		if parts[1] != "36000" {
			return fmt.Errorf("admin_password hash rounds is incorrect, expected 36000, got %s", parts[1])
		}

		if len(parts[2]) != 12 {
			return fmt.Errorf("admin_password salt length is incorrect, expected 12, got %d", len(parts[2]))
		}

		if len(parts[3]) == 0 {
			return fmt.Errorf("admin_password hash part is empty")
		}

		return nil
	}
}
