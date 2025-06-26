package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/pexip/go-infinity-sdk/v38"
	"os"
	"testing"
)

func getTestProtoV5ProviderFactories(client *infinity.ClientMock) map[string]func() (tfprotov5.ProviderServer, error) {
	return map[string]func() (tfprotov5.ProviderServer, error){
		"pexip": providerserver.NewProtocol5WithError(newTestProvider(client)),
	}
}

func TestMain(m *testing.M) {
	if os.Getenv("TF_ACC") == "" {
		os.Exit(m.Run())
	}
	resource.TestMain(m)
}
