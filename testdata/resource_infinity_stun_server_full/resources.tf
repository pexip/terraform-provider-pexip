/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */
resource "pexip_infinity_stun_server" "tf-test-stun-server" {
  name        = "tf-test-stun-server-full"
  description = "tf-test STUN server description"
  address     = "stun-full.example.com"
  port        = 5349
}
