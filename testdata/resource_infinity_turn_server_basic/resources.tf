/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_turn_server" "turn-server-test" {
  name           = "turn-server-test"
  description    = "Test TURN server"
  address        = "test-turn-server.dev.pexip.network"
  port           = 8080
  server_type    = "namepsw"
  transport_type = "udp"
  username       = "turnuser"
  password       = "turnpassword"
  secret_key     = "turnsecretkey"
}