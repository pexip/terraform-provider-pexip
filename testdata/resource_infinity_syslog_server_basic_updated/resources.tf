/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_syslog_server" "syslog_server-test" {
  address     = "10.1.1.50"                 // Updated address
  port        = 1514                        // Updated port
  description = "Updated Test SyslogServer" // Updated description
  transport   = "tcp"                       // Updated value
  audit_log   = false                       // Updated to false
  support_log = false                       // Updated to false
  web_log     = false                       // Updated to false
}