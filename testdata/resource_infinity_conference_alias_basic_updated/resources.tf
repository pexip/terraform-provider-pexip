/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_conference" "test-conference" {
  name         = "test-conference"
  description  = "Test Conference"
  service_type = "conference"
}

resource "pexip_infinity_conference_alias" "conference_alias-test" {
  alias       = "test-alias-updated"           // Updated value
  description = "Updated Test ConferenceAlias" // Updated description
  conference  = pexip_infinity_conference.test-conference.id

  depends_on = [
    pexip_infinity_conference.test-conference
  ]
}