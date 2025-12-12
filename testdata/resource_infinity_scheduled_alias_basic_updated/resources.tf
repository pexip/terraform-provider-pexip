/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_ms_exchange_connector" "test-connector" {
  name                                             = "test-exchange-connector"
  description                                      = "Test Exchange Connector"
  meeting_buffer_before                            = 300
  meeting_buffer_after                             = 300
  scheduled_alias_suffix_length                    = 6
  room_mailbox_email_address                       = "test@example.com"
  room_mailbox_name                                = "test-exchange-connector"
  url                                              = "https://exchange.test.local"
  username                                         = "testuser"
  password                                         = "testpass"
  authentication_method                            = "oauth2"
  auth_provider                                    = "azure"
  uuid                                             = "test-uuid-value"
  scheduled_alias_prefix                           = "test"
  scheduled_alias_domain                           = "example.com"
  enable_dynamic_vmrs                              = true
  enable_personal_vmrs                             = true
  allow_new_users                                  = true
  disable_proxy                                    = true
  use_custom_add_in_sources                        = true
  enable_addin_debug_logs                          = true
  oauth_client_id                                  = "test-client-id"
  oauth_client_secret                              = "test-secret"
  oauth_auth_endpoint                              = "test-auth-endpoint"
  oauth_token_endpoint                             = "test-token-endpoint"
  oauth_redirect_uri                               = "test-redirect-uri"
  oauth_refresh_token                              = "test-refresh-token"
  oauth_state                                      = "test-state"
  kerberos_realm                                   = "test-realm"
  kerberos_kdc                                     = "test-kdc"
  kerberos_kdc_https_proxy                         = "test-proxy"
  kerberos_exchange_spn                            = "test-spn"
  kerberos_enable_tls                              = true
  kerberos_auth_every_request                      = true
  kerberos_verify_tls_using_custom_ca              = true
  addin_server_domain                              = "test-domain"
  addin_display_name                               = "test-exchange-connector"
  addin_description                                = "Test Exchange Connector"
  addin_provider_name                              = "test-exchange-connector"
  addin_button_label                               = "test-button"
  addin_group_label                                = "test-group"
  addin_supertip_title                             = "test-title"
  addin_supertip_description                       = "Test Exchange Connector"
  addin_application_id                             = "test-app-id"
  addin_authority_url                              = "https://example.com"
  addin_oidc_metadata_url                          = "https://example.com"
  addin_authentication_method                      = "web_api"
  addin_naa_web_api_application_id                 = "test-naa-app-id"
  personal_vmr_oauth_client_id                     = "test-vmr-client-id"
  personal_vmr_oauth_client_secret                 = "test-vmr-secret"
  personal_vmr_oauth_auth_endpoint                 = "test-vmr-auth"
  personal_vmr_oauth_token_endpoint                = "test-vmr-token"
  personal_vmr_adfs_relying_party_trust_identifier = "test-adfs"
  office_js_url                                    = "https://example.com"
  microsoft_fabric_url                             = "https://example.com"
  microsoft_fabric_components_url                  = "https://example.com"
  additional_add_in_script_sources                 = "test-sources"
  domains                                          = "test-domain"
  non_idp_participants                             = "test-participants"
}

resource "pexip_infinity_scheduled_alias" "scheduled_alias-test" {
  alias              = "updated-scheduled-alias"                   // Updated value
  alias_number       = 9876543210                                  // Updated value
  numeric_alias      = "987654"                                    // Updated value
  uuid               = "22222222-2222-2222-2222-222222222222"      // Updated value
  exchange_connector = pexip_infinity_ms_exchange_connector.test-connector.id
  is_used            = false                                       // Updated to false
  ews_item_uid       = "updated-ews-uid"                           // Updated value

  depends_on = [
    pexip_infinity_ms_exchange_connector.test-connector
  ]
}