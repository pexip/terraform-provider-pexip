/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_ssh_authorized_key" "tf-test-ssh-key" {
  key     = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIKTj7PIu5ycIpVVxMYlnHmVKlhG4ALxqryNSfy59XIGf tf-test"
  keytype = "ssh-ed25519"
  comment = "tf-test"
}
