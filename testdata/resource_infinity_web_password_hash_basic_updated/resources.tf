/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_web_password_hash" "web_password_hash-test" {
  password = "updated-value" // Updated value
  salt     = "mnopqrstuvwx"  // Updated value, exactly 12 characters
  rounds   = 6000            // Updated value, must be >= 5000
}