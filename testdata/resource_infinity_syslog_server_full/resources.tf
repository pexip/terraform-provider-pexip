/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */
resource "pexip_infinity_syslog_server" "tf-test-syslog-server" {
  address     = "syslog-full.example.com"
  description = "tf-test syslog server description"
  port        = 1514
  transport   = "tls"
  audit_log   = true
  support_log = true
  web_log     = true
}
