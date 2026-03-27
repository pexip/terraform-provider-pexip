resource "pexip_infinity_mjx_exchange_deployment" "test_min" {
  name                     = "tf-test mjx-exchange-deployment min"
  service_account_username = "exchange-service@example.com"
  service_account_password = "test-password"
}

resource "pexip_infinity_mjx_exchange_deployment" "test_full" {
  name                     = "tf-test mjx-exchange-deployment full"
  service_account_username = "exchange-service@example.com"
  service_account_password = "test-password"
}

resource "pexip_infinity_mjx_exchange_autodiscover_url" "test_full" {
  name                = "tf-test mjx-exchange-autodiscover-url full"
  description         = "Test Exchange Autodiscover URL description"
  url                 = "https://autodiscover-full.example.com/autodiscover/autodiscover.xml"
  exchange_deployment = pexip_infinity_mjx_exchange_deployment.test_full.id
}
