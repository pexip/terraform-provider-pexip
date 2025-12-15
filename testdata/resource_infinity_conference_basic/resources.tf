/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_conference" "conference-test" {
  name         = "conference-test"
  description  = "Test Conference"
  service_type = "conference"
  pin          = "1234"
  tag          = "test-tag"
}