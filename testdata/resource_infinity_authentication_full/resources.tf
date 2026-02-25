/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_authentication" "authentication-test" {
  # For testing purposes, we set all fields to non-default values but leave the source as LOCAL, otherwise all the other settings would be validated.
  #source = "LOCAL"
  # client_certificate and api_oauth2_disable_basic should not be set or Terraform will not be able to access the mgmt API

  api_oauth2_allow_all_perms = true
  api_oauth2_expiration      = 7200
  ldap_base_dn      = "dc=example,dc=com"
  ldap_bind_password = "SuperSecretLdapPassword123!"
  ldap_bind_username = "CN=Service Account,OU=Service Accounts,DC=example,DC=com"
  ldap_group_filter  = "(|(objectclass=group)(objectclass=groupOfNames)(objectclass=groupOfUniqueNames)(objectclass=posixGroup)(objectclass=testGroup))"
  ldap_group_membership_filter = "(|(member={userdn})(uniquemember={userdn}))"
  ldap_group_search_dn        = "OU=Groups,DC=example,DC=com"
  ldap_permit_no_tls      = true
  ldap_server       = "ldap.example.com"
  ldap_use_global_catalog = true
  ldap_user_filter        = "(&(objectclass=person))"
  ldap_user_group_attributes = "memberOftest"
  ldap_user_search_dn     = "OU=Users,DC=example,DC=com"
  ldap_user_search_filter = "(|(uid={username}))"
  oidc_auth_method = "private_key"
  oidc_authorize_url       = "https://auth.example.com/oauth2/authorize"
  oidc_client_id         = "pexip-infinity-client-id"
  oidc_client_secret     = "SuperSecretOidcClientSecret456!"
  oidc_domain_hint  = "example.com"
  oidc_groups_field   = "testgroups"
  oidc_login_button = "Sign in with Corporate SSO"
  #oidc_metadata = ""
  oidc_metadata_url      = "https://auth.example.com/.well-known/openid-configuration"
  oidc_private_key       = "-----BEGIN PRIVATE KEY-----\n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\nWOMLmDqWZ7RW2tX5wIKroWeAAvsVjg4=\n-----END PRIVATE KEY-----"
  oidc_required_key   = "department"
  oidc_required_value = "IT"
  oidc_scope       = "openid profile email groups"
  oidc_token_endpoint_url  = "https://auth.example.com/oauth2/token"
  oidc_username_field = "preferred_username_test"
}