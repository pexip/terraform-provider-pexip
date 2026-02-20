/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_scheduled_scaling" "test" {
  policy_name         = "tf-test scheduled scaling min"
  policy_type         = "TeamsConnectorScaling"
  resource_identifier = "tf-test-resource-group-min"
  enabled             = true
  local_timezone      = "UTC"
  start_date          = "2024-01-01"
  time_from           = "09:00:00"
  time_to             = "17:00:00"
  instances_to_add    = 1
  minutes_in_advance  = 15
  mon                 = false
  tue                 = false
  wed                 = false
  thu                 = false
  fri                 = false
  sat                 = false
  sun                 = false
}
