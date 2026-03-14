/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */
resource "pexip_infinity_turn_server" "tf-test-turn-server" {
  name           = "tf-test-turn-server-full"
  description    = "tf-test TURN server description"
  address        = "turn-full.example.com"
  port           = 5349
  server_type    = "coturn_shared"
  transport_type = "tls"
  username       = "tf-test-username"
  password       = "tf-test-password"
  secret_key     = "tf-test-secret-key"
}
