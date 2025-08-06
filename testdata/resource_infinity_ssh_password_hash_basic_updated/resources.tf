/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_ssh_password_hash" "ssh_password_hash-test" {
  password = "updated-value"    // Updated value
  salt     = "qrstuvwxyzabcdef" // Updated value
  rounds   = 6000               // Updated value
}