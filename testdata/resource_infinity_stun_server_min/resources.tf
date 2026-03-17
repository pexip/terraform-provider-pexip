/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */
resource "pexip_infinity_stun_server" "tf-test-stun-server" {
  name    = "tf-test-stun-server"
  address = "stun.example.com"
}
