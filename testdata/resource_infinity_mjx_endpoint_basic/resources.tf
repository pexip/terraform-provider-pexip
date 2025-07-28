resource "pexip_infinity_mjx_endpoint" "mjx_endpoint-test" {
  name = "mjx_endpoint-test"
  description = "Test MjxEndpoint"
  endpoint_type = "polycom"
  room_resource_email = "test@example.com"
  mjx_endpoint_group = "test-value"
  api_address = "test-server.example.com"
  api_username = "mjx_endpoint-test"
  api_password = "test-value"
  use_https = "yes"
  verify_cert = "yes"
  poly_username = "mjx_endpoint-test"
  poly_password = "test-value"
  poly_raise_alarms_for_this_endpoint = true
  webex_device_id = "test-value"
}