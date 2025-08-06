/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_syslog_server" "syslog_server-test" {
  address     = "192.168.1.50"
  port        = 514
  description = "Test SyslogServer"
  transport   = "udp"
  audit_log   = true
  support_log = true
  web_log     = true
}