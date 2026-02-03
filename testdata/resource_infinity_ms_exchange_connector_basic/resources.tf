/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_ms_exchange_connector" "ms_exchange_connector-test" {
  name                                             = "ms_exchange_connector-test"
  description                                      = "Test MsExchangeConnector"
  meeting_buffer_before                            = 300
  meeting_buffer_after                             = 300
  scheduled_alias_suffix_length                    = 6
  room_mailbox_email_address                       = "test@example.com"
  room_mailbox_name                                = "ms_exchange_connector-test"
  url                                              = "https://example.com"
  username                                         = "ms_exchange_connector-test"
  password                                         = "test-value"
  authentication_method                            = "OAUTH"
  auth_provider                                    = "AZURE"
  scheduled_alias_prefix                           = "test-value"
  scheduled_alias_domain                           = "example.com"
  enable_dynamic_vmrs                              = true
  enable_personal_vmrs                             = true
  allow_new_users                                  = true
  disable_proxy                                    = true
  use_custom_add_in_sources                        = true
  enable_addin_debug_logs                          = true
  oauth_client_id                                  = "test-value"
  oauth_client_secret                              = "test-value"
  oauth_auth_endpoint                              = "test-value"
  oauth_token_endpoint                             = "test-value"
  oauth_redirect_uri                               = "test-value"
  oauth_state                                      = "test-value"
  kerberos_realm                                   = "test-value"
  kerberos_kdc                                     = "test-value"
  kerberos_kdc_https_proxy                         = "test-value"
  kerberos_exchange_spn                            = "test-value"
  kerberos_enable_tls                              = true
  kerberos_auth_every_request                      = true
  kerberos_verify_tls_using_custom_ca              = true
  addin_server_domain                              = "test-value"
  addin_display_name                               = "ms_exchange_connector-test"
  addin_description                                = "Test MsExchangeConnector"
  addin_provider_name                              = "ms_exchange_connector-test"
  addin_button_label                               = "test-value"
  addin_group_label                                = "test-value"
  addin_supertip_title                             = "test-value"
  addin_supertip_description                       = "Test MsExchangeConnector"
  addin_application_id                             = "test-value"
  addin_authority_url                              = "https://example.com"
  addin_oidc_metadata_url                          = "https://example.com"
  addin_authentication_method                      = "EXCHANGE_USER_ID_TOKEN"
  addin_naa_web_api_application_id                 = "test-value"
  personal_vmr_oauth_client_id                     = "test-value"
  personal_vmr_oauth_client_secret                 = "test-value"
  personal_vmr_oauth_auth_endpoint                 = "test-value"
  personal_vmr_oauth_token_endpoint                = "test-value"
  personal_vmr_adfs_relying_party_trust_identifier = "test-value"
  office_js_url                                    = "https://example.com"
  microsoft_fabric_url                             = "https://example.com"
  microsoft_fabric_components_url                  = "https://example.com"
  additional_add_in_script_sources                 = "test-value"
  host_identity_provider_group                     = "test-server.example.com"
  ivr_theme                                        = "test-value"
  non_idp_participants                             = "disallow_all"
}