/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_google_auth_server" "google_auth_server-test" {
  name             = "google_auth_server-test"
  description      = "Test GoogleAuthServer"
  application_type = "installed"
  client_id        = "123456789012-abcdefghijklmnopqrstuvwxyz123456.apps.googleusercontent.com"
  client_secret    = "GOCSPX-abcdefghijklmnopqrstuvwxyz"
}