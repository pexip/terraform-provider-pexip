package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/pexip/terraform-provider-pexip/internal/test"
	"os"
	"testing"
)

func TestInfinityManagerConfig(t *testing.T) {
	os.Setenv("TF_ACC", "1")
	expected := test.LoadTestData(t, "data_infinity_manager_config_basic_rendered.txt")
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
					resource.TestCheckResourceAttr("data.pexip_infinity_manager_config.master", "rendered", expected),
				),
			},
		},
	})
}
