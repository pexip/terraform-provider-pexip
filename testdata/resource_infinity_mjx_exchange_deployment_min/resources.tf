resource "pexip_infinity_mjx_exchange_deployment" "test" {
  name                     = "tf-test mjx-exchange-deployment min"
  service_account_username = "exchange-service@example.com"
  service_account_password = "test-password"
  oauth_client_id          = "12345678-1234-1234-1234-123456789012"
}
