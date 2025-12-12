/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_scheduled_scaling" "scheduled_scaling-test" {
  policy_name         = "scheduled_scaling-test"
  policy_type         = "TeamsConnectorScaling"
  resource_identifier = "updated-value"   // Updated value
  enabled             = false              // Updated to false
  local_timezone      = "America/New_York" // Updated value
  start_date          = "2024-02-01"       // Updated value
  time_from           = "08:00"            // Updated value
  time_to             = "18:00"            // Updated value
  instances_to_add    = 3                  // Updated value
  minutes_in_advance  = 30                 // Updated value
  mon                 = false              // Updated to false
  tue                 = false              // Updated to false
  wed                 = false              // Updated to false
  thu                 = false              // Updated to false
  fri                 = false              // Updated to false
  sat                 = false              // Updated to false
  sun                 = false              // Updated to false
}