resource "pexip_infinity_scheduled_scaling" "scheduled_scaling-test" {
  policy_name         = "scheduled_scaling-test"
  policy_type         = "management_vm" // Updated value
  resource_identifier = "updated-value" // Updated value
  enabled             = false           // Updated to false
  local_timezone      = "updated-value" // Updated value
  start_date          = "updated-value" // Updated value
  time_from           = "updated-value" // Updated value
  time_to             = "updated-value" // Updated value
  instances_to_add    = 3               // Updated value
  minutes_in_advance  = 30              // Updated value
  mon                 = false           // Updated to false
  tue                 = false           // Updated to false
  wed                 = false           // Updated to false
  thu                 = false           // Updated to false
  fri                 = false           // Updated to false
  sat                 = false           // Updated to false
  sun                 = false           // Updated to false
}