---
page_title: "pexip_infinity_authentication Resource - terraform-provider-pexip"
subcategory: ""
description: |-
  Manages the Pexip Infinity authentication configuration.
---

# pexip_infinity_authentication (Resource)

Manages the authentication configuration with the Infinity service. This is a singleton resource - only one authentication configuration exists per system. The resource configures how users authenticate to Pexip Infinity, including local authentication, LDAP integration, and OpenID Connect (OIDC) single sign-on.

## Example Usage

### Local Authentication (Default)

```terraform
resource "pexip_infinity_authentication" "local_auth" {
  source = "local"
}
```

### LDAP Authentication

```terraform
resource "pexip_infinity_authentication" "ldap_auth" {
  source = "ldap"
  
  # LDAP server configuration
  ldap_server                = "ldaps://ldap.company.com:636"
  ldap_base_dn              = "dc=company,dc=com"
  ldap_bind_username        = "cn=pexip,ou=service,dc=company,dc=com"
  ldap_bind_password        = var.ldap_bind_password
  
  # User search configuration
  ldap_user_search_dn       = "ou=users,dc=company,dc=com"
  ldap_user_filter          = "(&(objectclass=person)(!(objectclass=computer)))"
  ldap_user_search_filter   = "(|(uid={username})(sAMAccountName={username})(mail={username}))"
  ldap_user_group_attributes = "memberOf"
  
  # Group search configuration
  ldap_group_search_dn         = "ou=groups,dc=company,dc=com"
  ldap_group_filter            = "(objectclass=group)"
  ldap_group_membership_filter = "(member={userdn})"
  
  # Security settings
  ldap_use_global_catalog = false
  ldap_permit_no_tls      = false
}
```

### OpenID Connect (OIDC) Authentication

```terraform
resource "pexip_infinity_authentication" "oidc_auth" {
  source = "oidc"
  
  # OIDC provider configuration
  oidc_metadata_url     = "https://login.microsoftonline.com/tenant-id/.well-known/openid_configuration"
  oidc_client_id        = var.oidc_client_id
  oidc_client_secret    = var.oidc_client_secret
  oidc_auth_method      = "client_secret_basic"
  oidc_scope           = "openid profile email"
  
  # Claim mapping
  oidc_username_field   = "preferred_username"
  oidc_groups_field     = "groups"
  
  # Access control
  oidc_required_key     = "department"
  oidc_required_value   = "IT"
  
  # UI customization
  oidc_domain_hint      = "company.com"
  oidc_login_button     = "Sign in with Corporate SSO"
}
```

### OIDC with JWT Private Key Authentication

```terraform
resource "pexip_infinity_authentication" "oidc_jwt_auth" {
  source = "oidc"
  
  # OIDC configuration with JWT
  oidc_client_id        = var.oidc_client_id
  oidc_private_key      = var.oidc_private_key
  oidc_auth_method      = "private_key_jwt"
  oidc_authorize_url    = "https://auth.company.com/oauth2/authorize"
  oidc_token_endpoint_url = "https://auth.company.com/oauth2/token"
  oidc_scope           = "openid profile email groups"
  
  # Claim configuration
  oidc_username_field   = "sub"
  oidc_groups_field     = "department_groups"
}
```

### OAuth2 API Configuration

```terraform
resource "pexip_infinity_authentication" "api_config" {
  source = "local"
  
  # OAuth2 API settings
  api_oauth2_disable_basic    = true
  api_oauth2_allow_all_perms  = false
  api_oauth2_expiration       = 7200  # 2 hours
}
```

### Client Certificate Authentication

```terraform
resource "pexip_infinity_authentication" "cert_auth" {
  source             = "local"
  client_certificate = "required"
  
  # Combined with OAuth2 API settings
  api_oauth2_disable_basic = true
  api_oauth2_expiration    = 3600
}
```

## Schema

### Optional

- `source` (String) - Authentication source. Valid values: local, ldap, oidc.
- `client_certificate` (String) - Client certificate requirement. Valid values: disabled, optional, required.
- `api_oauth2_disable_basic` (Boolean) - Whether to disable basic authentication for OAuth2 API access. Defaults to `false`.
- `api_oauth2_allow_all_perms` (Boolean) - Whether to allow all permissions for OAuth2 API access. Defaults to `false`.
- `api_oauth2_expiration` (Number) - OAuth2 token expiration time in seconds. Minimum: 60 seconds. Defaults to `3600`.

#### LDAP Configuration

- `ldap_server` (String) - LDAP server URL (ldap:// or ldaps://).
- `ldap_base_dn` (String) - LDAP base distinguished name for searches.
- `ldap_bind_username` (String) - LDAP bind username for authentication.
- `ldap_bind_password` (String, Sensitive) - LDAP bind password for authentication.
- `ldap_user_search_dn` (String) - LDAP distinguished name for user searches.
- `ldap_user_filter` (String) - LDAP filter for user searches. Defaults to `"(&(objectclass=person)(!(objectclass=computer)))"`.
- `ldap_user_search_filter` (String) - LDAP search filter for users. Defaults to `"(|(uid={username})(sAMAccountName={username}))"`.
- `ldap_user_group_attributes` (String) - LDAP attributes to use for user group membership. Defaults to `"memberOf"`.
- `ldap_group_search_dn` (String) - LDAP distinguished name for group searches.
- `ldap_group_filter` (String) - LDAP filter for group searches. Defaults to `"(|(objectclass=group)(objectclass=groupOfNames)(objectclass=groupOfUniqueNames)(objectclass=posixGroup))"`.
- `ldap_group_membership_filter` (String) - LDAP filter for group membership queries. Defaults to `"(|(member={userdn})(uniquemember={userdn})(memberuid={useruid}))"`.
- `ldap_use_global_catalog` (Boolean) - Whether to use LDAP global catalog. Defaults to `false`.
- `ldap_permit_no_tls` (Boolean) - Whether to permit LDAP connections without TLS. Defaults to `false`.

#### OIDC Configuration

- `oidc_metadata_url` (String) - OpenID Connect metadata URL.
- `oidc_metadata` (String) - OpenID Connect metadata as JSON string. Defaults to `"{}"`.
- `oidc_client_id` (String) - OpenID Connect client ID.
- `oidc_client_secret` (String, Sensitive) - OpenID Connect client secret.
- `oidc_private_key` (String, Sensitive) - OpenID Connect private key for JWT signing.
- `oidc_auth_method` (String) - OpenID Connect authentication method. Valid values: client_secret_basic, client_secret_post, private_key_jwt. Defaults to `"client_secret"`.
- `oidc_scope` (String) - OpenID Connect scope for authentication requests. Defaults to `"openid profile email"`.
- `oidc_authorize_url` (String) - OpenID Connect authorization URL.
- `oidc_token_endpoint_url` (String) - OpenID Connect token endpoint URL.
- `oidc_username_field` (String) - OpenID Connect claim field for username. Defaults to `"preferred_username"`.
- `oidc_groups_field` (String) - OpenID Connect claim field for groups. Defaults to `"groups"`.
- `oidc_required_key` (String) - OpenID Connect required claim key for access control.
- `oidc_required_value` (String) - OpenID Connect required claim value for access control.
- `oidc_domain_hint` (String) - OpenID Connect domain hint for login.
- `oidc_login_button` (String) - Text for the OpenID Connect login button. Maximum length: 128 characters.

### Read-Only

- `id` (String) - Resource URI for the authentication configuration in Infinity.

## Import

Import is supported using the following syntax:

```shell
terraform import pexip_infinity_authentication.example auth
```

For singleton resources, the import ID doesn't matter since there's only one instance.

## Usage Notes

### Authentication Sources
- **local**: Built-in Pexip user database
- **ldap**: Active Directory or LDAP integration
- **oidc**: OpenID Connect single sign-on

### Singleton Resource Behavior
- Only one authentication configuration exists per Pexip Infinity system
- Creating this resource updates the existing configuration
- Deleting resets authentication to local mode
- Import ID is ignored since there's only one instance

### LDAP Integration
- Use secure LDAP (ldaps://) for production environments
- Test LDAP connectivity before applying configuration
- Ensure bind user has appropriate read permissions
- Configure proper group membership attributes

### OIDC Integration
- Requires registered application in identity provider
- Use appropriate authentication method for security requirements
- Configure proper scopes and claim mappings
- Test token validation and user attribute mapping

### Client Certificate Authentication
- Can be combined with other authentication methods
- Requires proper PKI infrastructure
- Configure certificate validation rules separately
- Monitor certificate expiration and renewal

### OAuth2 API Configuration
- Controls API authentication behavior
- Longer expiration times reduce re-authentication
- Disable basic auth for enhanced security
- Monitor API token usage and revocation

## Troubleshooting

### Common Issues

**LDAP Authentication Failures**
- Verify LDAP server connectivity and credentials
- Check bind user permissions and DN format
- Test LDAP queries using tools like ldapsearch
- Verify SSL/TLS certificate trust for ldaps://

**OIDC Authentication Problems**
- Verify client ID and secret are correct
- Check OIDC provider metadata URL accessibility
- Ensure proper redirect URI configuration
- Verify claim mappings match provider attributes

**User Not Found Errors**
- Check LDAP user search DN and filter configuration
- Verify username format matches search filter
- Ensure user exists in specified search base
- Test user search filters manually

**Group Membership Issues**
- Verify group search DN and filter configuration
- Check group membership attribute mapping
- Ensure users have proper group memberships
- Test group queries independently

**Certificate Authentication Problems**
- Verify client certificate trust chain
- Check certificate validity and expiration
- Ensure proper certificate subject validation
- Verify certificate revocation status

**API Authentication Issues**
- Check OAuth2 token expiration settings
- Verify API permissions and scopes
- Ensure proper token format and validation
- Monitor token usage and revocation

**Import Fails**
- Authentication configuration always exists as singleton
- Import should work with any ID value
- Verify provider authentication and permissions
- Check system connectivity to Infinity Manager

**Performance Problems**
- Monitor LDAP query response times
- Check OIDC token validation performance
- Verify network connectivity to authentication sources
- Consider caching strategies for external authentication

**Security Configuration Issues**
- Ensure TLS is properly configured for external connections
- Verify certificate validation settings
- Check authentication timeout settings
- Monitor for authentication bypass attempts

**Migration Between Authentication Sources**
- Plan user migration strategy carefully
- Test authentication with pilot users
- Ensure user attribute mapping is correct
- Have rollback plan ready for issues