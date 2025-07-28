resource "pexip_infinity_scheduled_scaling" "scheduled_scaling-test" {
  policy_name = "scheduled_scaling-test"
  policy_type = "worker_vm"
  resource_identifier = "test-value"
  enabled = true
  local_timezone = "test-value"
  start_date = "test-value"
  time_from = "test-value"
  time_to = "test-value"
  instances_to_add = 2
  minutes_in_advance = 15
  mon = true
  tue = true
  wed = true
  thu = true
  fri = true
  sat = true
  sun = true
}