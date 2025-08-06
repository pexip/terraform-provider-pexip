/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_ssh_password_hash" "ssh_password_hash-test" {
  password = "test-value"
  salt     = "abcdefghijklmnop"
  rounds   = 5000
}