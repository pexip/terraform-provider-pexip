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

resource "pexip_infinity_mjx_exchange_autodiscover_url" "test" {
  name                = "tf-test mjx-exchange-autodiscover-url min"
  url                 = "https://autodiscover.example.com/autodiscover/autodiscover.xml"
  exchange_deployment = pexip_infinity_mjx_exchange_deployment.test_min.id
}
