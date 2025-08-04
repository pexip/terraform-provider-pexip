# infinity_ldap_sync_source

Manages an LDAP synchronization source configuration with the Infinity service. LDAP sync sources enable automatic synchronization of user and group information from Active Directory or other LDAP directories.

## Example Usage

```hcl
resource "pexip_infinity_ldap_sync_source" "example" {
  name                     = "Corporate AD"
  description              = "Corporate Active Directory synchronization"
  ldap_server              = "corp.example.com"
  ldap_base_dn             = "dc=corp,dc=example,dc=com"
  ldap_bind_username       = "pexip-service@corp.example.com"
  ldap_bind_password       = "secure-password"
  ldap_use_global_catalog  = true
  ldap_permit_no_tls       = false
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The unique name of the LDAP synchronization source. Maximum length: 250 characters.
* `ldap_server` - (Required) The hostname of the LDAP server. Enter a domain name for DNS SRV lookup or an FQDN for DNS A/AAAA lookup. Maximum length: 255 characters.
* `ldap_base_dn` - (Required) The base DN of the LDAP forest to query (e.g. dc=example,dc=com). Maximum length: 255 characters.
* `ldap_bind_username` - (Required) The username used to bind to the LDAP server. This should be a domain user service account. Maximum length: 255 characters.
* `ldap_bind_password` - (Required) The password used to bind to the LDAP server. Maximum length: 100 characters. This field is sensitive.
* `description` - (Optional) A description of the LDAP synchronization source. Maximum length: 250 characters.
* `ldap_use_global_catalog` - (Optional) Search the Active Directory Global Catalog instead of traditional LDAP. Defaults to false.
* `ldap_permit_no_tls` - (Optional) Permit LDAP queries to be sent over an insecure connection. Defaults to false.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Resource URI for the LDAP sync source in Infinity.
* `resource_id` - The resource integer identifier for the LDAP sync source in Infinity.

## Import

LDAP sync sources can be imported using their resource ID:

```bash
terraform import pexip_infinity_ldap_sync_source.example 123
```

## Security Notes

- The `ldap_bind_password` field is marked as sensitive and will not be displayed in Terraform output.
- Use a dedicated service account with minimal required permissions for LDAP binding.
- Set `ldap_permit_no_tls` to `false` (default) to ensure encrypted LDAP communications.
- Consider using `ldap_use_global_catalog` for better performance in multi-domain Active Directory environments.