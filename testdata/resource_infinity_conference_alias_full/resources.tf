/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_conference" "tf-test-conference" {
  name         = "tf-test-conference"
  description  = "Test Conference"
  service_type = "conference"
}

resource "pexip_infinity_conference_alias" "tf-test-conference-alias" {
  alias       = "tf-test-alias"
  description = "Test Conference Alias Description"
  conference  = pexip_infinity_conference.tf-test-conference.id
}
