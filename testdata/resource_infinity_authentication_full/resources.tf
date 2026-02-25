/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

# Create TLS private key for OIDC authentication
resource "tls_private_key" "oidc" {
  algorithm = "RSA"
  rsa_bits  = 4096
}

resource "pexip_infinity_authentication" "authentication-test" {
  # For testing purposes, we set all fields to non-default values but leave the source as LOCAL, otherwise all the other settings would be validated.
  #source = "LOCAL"
  # client_certificate and api_oauth2_disable_basic should not be set or Terraform will not be able to access the mgmt API

  api_oauth2_allow_all_perms   = true
  api_oauth2_expiration        = 7200
  ldap_base_dn                 = "dc=example,dc=com"
  ldap_bind_password           = "SuperSecretLdapPassword123!"
  ldap_bind_username           = "CN=Service Account,OU=Service Accounts,DC=example,DC=com"
  ldap_group_filter            = "(|(objectclass=group)(objectclass=groupOfNames)(objectclass=groupOfUniqueNames)(objectclass=posixGroup)(objectclass=testGroup))"
  ldap_group_membership_filter = "(|(member={userdn})(uniquemember={userdn}))"
  ldap_group_search_dn         = "OU=Groups,DC=example,DC=com"
  ldap_permit_no_tls           = true
  ldap_server                  = "ldap.example.com"
  ldap_use_global_catalog      = true
  ldap_user_filter             = "(&(objectclass=person))"
  ldap_user_group_attributes   = "memberOftest"
  ldap_user_search_dn          = "OU=Users,DC=example,DC=com"
  ldap_user_search_filter      = "(|(uid={username}))"
  oidc_auth_method             = "private_key"
  oidc_authorize_url           = "https://auth.example.com/oauth2/authorize"
  oidc_client_id               = "pexip-infinity-client-id"
  oidc_client_secret           = "SuperSecretOidcClientSecret456!"
  oidc_domain_hint             = "example.com"
  oidc_groups_field            = "testgroups"
  oidc_login_button            = "Sign in with Corporate SSO"
  #oidc_metadata = ""
  oidc_metadata_url       = "https://auth.example.com/.well-known/openid-configuration"
  oidc_private_key        = tls_private_key.oidc.private_key_pem
  oidc_required_key       = "department"
  oidc_required_value     = "IT"
  oidc_scope              = "openid profile email groups"
  oidc_token_endpoint_url = "https://auth.example.com/oauth2/token"
  oidc_username_field     = "preferred_username_test"
}