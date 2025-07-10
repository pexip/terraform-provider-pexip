resource "pexip_infinity_policy_server" "policy-server-test" {
  name        = "test-policy-server"
  description = "Test Policy Server"
  url         = "https://test-policy-server.dev.pexip.network"
  username    = "testuser"
  password    = "updatedpassword"
}