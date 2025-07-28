resource "pexip_infinity_break_in_allow_list_address" "break_in_allow_list_address-test" {
  name = "break_in_allow_list_address-test"
  description = "Test BreakInAllowListAddress"
  address = "192.168.1.0"
  prefix = 24
  allowlist_entry_type = "temporary"
  ignore_incorrect_aliases = true
  ignore_incorrect_pins = true
}