/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */
resource "pexip_infinity_syslog_server" "tf-test-syslog-server" {
  address = "syslog.example.com"
}
