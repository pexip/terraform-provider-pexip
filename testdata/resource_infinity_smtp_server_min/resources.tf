/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_smtp_server" "smtp_server-test" {
  name                = "tf-test SMTP Server min"
  address             = "test-server.example.com"
}