package provider

import (
	"os"
	"testing"

	"github.com/pexip/go-infinity-sdk/v38/config"
	"github.com/pexip/go-infinity-sdk/v38/types"
	"github.com/stretchr/testify/mock"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/pexip/go-infinity-sdk/v38"

	"github.com/pexip/terraform-provider-pexip/internal/test"
)

func TestInfinityManagementVM(t *testing.T) {
	t.Parallel()
	_ = os.Setenv("TF_ACC", "1")

	// Create a mock client and set up expectations
	client := infinity.NewClientMock()

	// Mock the CreateManagementvm API call
	createResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/management_vm/123/",
	}
	client.On("PostWithResponse", mock.Anything, "configuration/v1/management_vm/", mock.Anything, mock.Anything).Return(createResponse, nil)

	// Shared state for mocking
	mockState := &config.ManagementVM{
		ID:                          123,
		ResourceURI:                 "/api/admin/configuration/v1/management_vm/123/",
		Name:                        "management_vm-test",
		Description:                 "Test ManagementVM",
		Address:                     "192.168.1.100",
		Netmask:                     "255.255.255.0",
		Gateway:                     "192.168.1.1",
		MTU:                         1500,
		Hostname:                    "management_vm-test",
		Domain:                      "example.com",
		AlternativeFQDN:             "alt.example.com",
		IPV6Address:                 test.StringPtr("2001:db8::1"),
		IPV6Gateway:                 test.StringPtr("2001:db8::1"),
		StaticNATAddress:            test.StringPtr("203.0.113.1"),
		HTTPProxy:                   test.StringPtr("http://proxy.example.com:8080"),
		TLSCertificate:              test.StringPtr("test-certificate"),
		EnableSSH:                   "yes",
		SSHAuthorizedKeysUseCloud:   true,
		SecondaryConfigPassphrase:   "test-passphrase",
		SNMPMode:                    "v1v2c",
		SNMPCommunity:               "public",
		SNMPUsername:                "management_vm-test",
		SNMPAuthenticationPassword:  "test-auth-pass",
		SNMPPrivacyPassword:         "test-priv-pass",
		SNMPSystemContact:           "admin@example.com",
		SNMPSystemLocation:          "datacenter",
		SNMPNetworkManagementSystem: test.StringPtr("192.168.1.200"),
		Initializing:                true,
		Primary:                     true,
	}

	// Mock the GetManagementvm API call for Read operations
	client.On("GetJSON", mock.Anything, "configuration/v1/management_vm/123/", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		management_vm := args.Get(2).(*config.ManagementVM)
		*management_vm = *mockState
	}).Maybe()

	// ManagementVM doesn't support update operations

	// Mock the DeleteManagementvm API call
	client.On("DeleteJSON", mock.Anything, mock.MatchedBy(func(path string) bool {
		return path == "configuration/v1/management_vm/123/"
	}), mock.Anything).Return(nil)

	testInfinityManagementVM(t, client)
}

func testInfinityManagementVM(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: getTestProtoV5ProviderFactories(client),
		Steps: []resource.TestStep{
			{
				Config: test.LoadTestFolder(t, "resource_infinity_management_vm_basic"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_management_vm.management_vm-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_management_vm.management_vm-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_management_vm.management_vm-test", "name", "management_vm-test"),
					resource.TestCheckResourceAttr("pexip_infinity_management_vm.management_vm-test", "description", "Test ManagementVM"),
					resource.TestCheckResourceAttr("pexip_infinity_management_vm.management_vm-test", "hostname", "management_vm-test"),
					resource.TestCheckResourceAttr("pexip_infinity_management_vm.management_vm-test", "ssh_authorized_keys_use_cloud", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_management_vm.management_vm-test", "snmp_username", "management_vm-test"),
					resource.TestCheckResourceAttr("pexip_infinity_management_vm.management_vm-test", "initializing", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_management_vm.management_vm-test", "primary", "true"),
				),
			},
			// ManagementVM doesn't support updates, so only test creation/read
		},
	})
}
