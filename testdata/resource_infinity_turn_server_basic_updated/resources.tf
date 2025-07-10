resource "pexip_infinity_turn_server" "turn-server-test" {
  name                    = "turn-server-test"
  description             = "Test TURN server"
  address                 = "test-turn-server.dev.pexip.network"
  port                    = 8081 // Updated port
  server_type             = "namepsw"
  transport_type          = "udp"
  username                = "turnuser"
  password                = "updatedturnpassword" // Updated password
  secret_key              = "updatedturnsecretkey" // Updated secret key
}