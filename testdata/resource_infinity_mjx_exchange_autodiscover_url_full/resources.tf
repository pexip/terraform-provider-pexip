resource "pexip_infinity_mjx_exchange_autodiscover_url" "test" {
  name                = "tf-test mjx-exchange-autodiscover-url full"
  description         = "Test Exchange Autodiscover URL description"
  url                 = "https://autodiscover-full.example.com/autodiscover/autodiscover.xml"
  exchange_deployment = "/api/admin/configuration/v1/mjx_exchange_deployment/2/"
}
