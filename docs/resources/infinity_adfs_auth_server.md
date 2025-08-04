# infinity_adfs_auth_server

Manages an ADFS OAuth 2.0 auth server configuration with the Infinity service.

## Example Usage

```hcl
resource "pexip_infinity_adfs_auth_server" "example" {
  name                               = "ADFS Server"
  description                        = "ADFS authentication server for corporate users"
  client_id                          = "12345678-1234-1234-1234-123456789012"
  federation_service_name            = "adfs.example.com"
  federation_service_identifier      = "http://adfs.example.com/adfs/services/trust"
  relying_party_trust_identifier_url = "https://pexip.example.com"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The unique name of the ADFS auth server. Maximum length: 250 characters.
* `description` - (Optional) A description of the ADFS auth server. Maximum length: 250 characters.
* `client_id` - (Required) The client ID for the ADFS OAuth 2.0 client. Maximum length: 250 characters.
* `federation_service_name` - (Required) The federation service name. Maximum length: 250 characters.
* `federation_service_identifier` - (Required) The federation service identifier. Maximum length: 250 characters.
* `relying_party_trust_identifier_url` - (Required) The relying party trust identifier URL. Maximum length: 250 characters.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Resource URI for the ADFS auth server in Infinity.
* `resource_id` - The resource integer identifier for the ADFS auth server in Infinity.

## Import

ADFS auth servers can be imported using their resource ID:

```bash
terraform import pexip_infinity_adfs_auth_server.example 123
```