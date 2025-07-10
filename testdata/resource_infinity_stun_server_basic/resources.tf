resource "pexip_infinity_stun_server" "stun-server-test" {
  name        = "stun-server-test"
  description = "Test STUN server"
  address     = "test-stun-server.dev.pexip.network"
  port        = 8080
}