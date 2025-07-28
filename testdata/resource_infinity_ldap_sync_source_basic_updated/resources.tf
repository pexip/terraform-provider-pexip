resource "pexip_infinity_ldap_sync_source" "ldap_sync_source-test" {
  name = "ldap_sync_source-test"
  description = "Updated Test LdapSyncSource"  // Updated description
  ldap_server = "updated-value"  // Updated value
  ldap_base_dn = "updated-value"  // Updated value
  ldap_bind_username = "ldap_sync_source-test"
  ldap_bind_password = "updated-value"  // Updated value
  ldap_use_global_catalog = false  // Updated to false
  ldap_permit_no_tls = false  // Updated to false
}