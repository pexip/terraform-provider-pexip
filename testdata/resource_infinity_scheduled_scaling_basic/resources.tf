/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_scheduled_scaling" "scheduled_scaling-test" {
  policy_name         = "scheduled_scaling-test"
  policy_type         = "TeamsConnectorScaling"
  resource_identifier = "test-resource-group"
  enabled             = true
  local_timezone      = "UTC"
  start_date          = "2024-01-01"
  time_from           = "09:00"
  time_to             = "17:00"
  instances_to_add    = 2
  minutes_in_advance  = 15
  mon                 = true
  tue                 = true
  wed                 = true
  thu                 = true
  fri                 = true
  sat                 = true
  sun                 = true
}