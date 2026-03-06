/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_smtp_server" "smtp_server-test" {
  name                = "tf-test SMTP Server full"
  description         = "full Test SMTPServer"    // Updated description
  address             = "updated-server.example.com" // Updated address
  port                = 465                          // Updated port
  username            = "smtp_server-test"
  password            = "updated-value"       // Updated value
  from_email_address  = "updated@example.com" // Updated email
  connection_security = "STARTTLS"            // Updated value
}