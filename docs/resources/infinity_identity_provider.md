# infinity_identity_provider

Manages an identity provider configuration with the Infinity service. Identity providers enable single sign-on (SSO) integration using SAML or OIDC protocols.

## Example Usage

### SAML Identity Provider

```hcl
resource "pexip_infinity_identity_provider" "saml_example" {
  name                            = "Corporate SAML IdP"
  description                     = "Corporate SAML identity provider"
  idp_type                        = "saml"
  signature_algorithm             = "rsa-sha256"
  digest_algorithm                = "sha256"
  assertion_consumer_service_url  = "https://pexip.example.com/saml/acs"
  sso_url                         = "https://idp.example.com/sso"
  idp_entity_id                   = "https://idp.example.com/entity"
  service_entity_id               = "https://pexip.example.com/entity"
  display_name_attribute_name     = "displayName"
  registration_alias_attribute_name = "mail"
  worker_fqdn_acs_urls            = false
}
```

### OIDC Identity Provider

```hcl
resource "pexip_infinity_identity_provider" "oidc_example" {
  name                               = "Azure AD OIDC"
  description                        = "Azure Active Directory OIDC provider"
  idp_type                           = "oidc"
  signature_algorithm                = "rsa-sha256"
  digest_algorithm                   = "sha256"
  assertion_consumer_service_url     = "https://pexip.example.com/oidc/callback"
  oidc_flow                          = "authorization_code"
  oidc_client_id                     = "12345678-1234-1234-1234-123456789012"
  oidc_client_secret                 = "secret-value"
  oidc_token_url                     = "https://login.microsoftonline.com/tenant-id/oauth2/v2.0/token"
  oidc_user_info_url                 = "https://graph.microsoft.com/oidc/userinfo"
  oidc_jwks_url                      = "https://login.microsoftonline.com/tenant-id/discovery/v2.0/keys"
  oidc_token_endpoint_auth_scheme    = "client_secret_post"
  oidc_token_signature_scheme        = "rs256"
  oidc_display_name_claim_name       = "name"
  oidc_registration_alias_claim_name = "preferred_username"
  oidc_additional_scopes             = "profile email"
  disable_popup_flow                 = false
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The unique name of the identity provider. Maximum length: 250 characters.
* `idp_type` - (Required) The identity provider type. Valid choices: `saml`, `oidc`.
* `signature_algorithm` - (Required) The signature algorithm. Valid choices: `rsa-sha256`, `rsa-sha1`.
* `digest_algorithm` - (Required) The digest algorithm. Valid choices: `sha256`, `sha1`.
* `assertion_consumer_service_url` - (Required) The assertion consumer service URL. Maximum length: 250 characters.
* `description` - (Optional) A description of the identity provider. Maximum length: 250 characters.

### SAML-specific Arguments

* `sso_url` - (Optional) The SSO URL for SAML identity providers. Maximum length: 250 characters.
* `idp_entity_id` - (Optional) The identity provider entity ID for SAML. Maximum length: 250 characters.
* `idp_public_key` - (Optional) The identity provider public key for SAML.
* `service_entity_id` - (Optional) The service entity ID for SAML. Maximum length: 250 characters.
* `service_public_key` - (Optional) The service public key for SAML.
* `service_private_key` - (Optional) The service private key for SAML. This field is sensitive.
* `display_name_attribute_name` - (Optional) The display name attribute name. Maximum length: 250 characters.
* `registration_alias_attribute_name` - (Optional) The registration alias attribute name. Maximum length: 250 characters.
* `worker_fqdn_acs_urls` - (Optional) Whether to use worker FQDN in ACS URLs. Defaults to false.

### OIDC-specific Arguments

* `disable_popup_flow` - (Optional) Whether to disable popup flow for OIDC.
* `oidc_flow` - (Optional) The OIDC flow type.
* `oidc_client_id` - (Optional) The OIDC client ID.
* `oidc_client_secret` - (Optional) The OIDC client secret. This field is sensitive.
* `oidc_token_url` - (Optional) The OIDC token URL.
* `oidc_user_info_url` - (Optional) The OIDC user info URL.
* `oidc_jwks_url` - (Optional) The OIDC JWKS URL.
* `oidc_token_endpoint_auth_scheme` - (Optional) The OIDC token endpoint authentication scheme.
* `oidc_token_signature_scheme` - (Optional) The OIDC token signature scheme.
* `oidc_display_name_claim_name` - (Optional) The OIDC display name claim name.
* `oidc_registration_alias_claim_name` - (Optional) The OIDC registration alias claim name.
* `oidc_additional_scopes` - (Optional) Additional OIDC scopes.
* `oidc_france_connect_required_eidas_level` - (Optional) Required eIDAS level for France Connect.

### Other Arguments

* `attributes` - (Optional) Additional attributes configuration.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Resource URI for the identity provider in Infinity.
* `resource_id` - The resource integer identifier for the identity provider in Infinity.
* `uuid` - The UUID of the identity provider.

## Import

Identity providers can be imported using their resource ID:

```bash
terraform import pexip_infinity_identity_provider.example 123
```