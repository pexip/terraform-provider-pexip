resource "pexip_infinity_break_in_allow_list_address" "break_in_allow_list_address-test" {
  name = "break_in_allow_list_address-test"
  description = "Updated Test BreakInAllowListAddress"
  address = "10.0.0.0"
  prefix = 16
  allowlist_entry_type = "permanent"
  ignore_incorrect_aliases = false
  ignore_incorrect_pins = false
}