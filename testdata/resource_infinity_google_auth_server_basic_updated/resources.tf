/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_google_auth_server" "google_auth_server-test" {
  name             = "google_auth_server-test"
  description      = "Updated Test GoogleAuthServer" // Updated description
  application_type = "installed"                     // Updated value
  client_id        = "updated-value"                 // Updated value
  client_secret    = "updated-value"                 // Updated value
}