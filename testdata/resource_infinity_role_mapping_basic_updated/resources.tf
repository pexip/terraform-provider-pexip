resource "pexip_infinity_role_mapping" "role_mapping-test" {
  name   = "role_mapping-test"
  source = "ldap_attribute" // Updated value
  value  = "updated-value"  // Updated value
}