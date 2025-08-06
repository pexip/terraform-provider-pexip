/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_conference_alias" "conference_alias-test" {
  alias       = "updated-value"                // Updated value
  description = "Updated Test ConferenceAlias" // Updated description
  conference  = "updated-value"                // Updated value
}