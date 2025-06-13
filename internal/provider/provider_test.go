package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"os"
	"testing"
)

var testProtoV5ProviderFactories = map[string]func() (tfprotov5.ProviderServer, error){
	"pexip": providerserver.NewProtocol5WithError(New()),
}

func TestMain(m *testing.M) {
	//if os.Getenv("TF_ACC") == "" {
	//	os.Exit(m.Run())
	//}
	os.Setenv("TF_ACC", "true") // Enable acceptance tests
	resource.TestMain(m)
}
