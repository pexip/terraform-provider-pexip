/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_ssh_authorized_key" "ssh_authorized_key-test" {
  keytype = "ssh-rsa"
  key     = "test-value"
  comment = "test-value"
}