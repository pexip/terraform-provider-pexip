/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_google_auth_server" "google_auth_server-test" {
  name             = "google_auth_server-test"
  description      = "Updated Test GoogleAuthServer"                                            // Updated description
  application_type = "installed"                                                                // Updated value
  client_id        = "987654321098-zyxwvutsrqponmlkjihgfedcba987654.apps.googleusercontent.com" // Updated value
  client_secret    = "GOCSPX-zyxwvutsrqponmlkjihgfedcba"                                        // Updated value
}