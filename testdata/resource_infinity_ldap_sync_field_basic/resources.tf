resource "pexip_infinity_ldap_sync_field" "ldap_sync_field-test" {
  name                   = "ldap_sync_field-test"
  description            = "Test LdapSyncField"
  template_variable_name = "ldap_sync_field-test"
  is_binary              = true
}