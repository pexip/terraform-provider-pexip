/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_conference" "test-adp-conference" {
  name         = "tf-test-adp-conference"
  service_type = "conference"
}

resource "pexip_infinity_system_location" "test-adp-location" {
  name = "tf-test-adp-location"
}

resource "pexip_infinity_automatic_participant" "automatic-participant-test" {
  alias = "tf-test-adp" 
}