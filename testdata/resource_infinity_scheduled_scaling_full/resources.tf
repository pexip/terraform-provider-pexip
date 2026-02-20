/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_scheduled_scaling" "test" {
  policy_name         = "tf-test scheduled scaling full"
  policy_type         = "TeamsConnectorScaling"
  resource_identifier = "tf-test-resource-group-full"
  enabled             = true
  local_timezone      = "America/New_York"
  start_date          = "2024-06-15"
  time_from           = "08:30"
  time_to             = "18:30"
  instances_to_add    = 5
  minutes_in_advance  = 30
  mon                 = true
  tue                 = true
  wed                 = true
  thu                 = true
  fri                 = true
  sat                 = false
  sun                 = false
}
