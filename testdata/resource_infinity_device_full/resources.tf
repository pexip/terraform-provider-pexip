/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_identity_provider_group" "tf-test-identity-provider-group" {
  name        = "tf-test-identity-provider-group"
  description = "Test Identity Provider Group for Device"
}

resource "pexip_infinity_device" "tf-test-device" {
  alias                           = "tf-test-device"
  description                     = "Test Device Description"
  username                        = "tf-test-user"
  password                        = "tf-test-pass"
  primary_owner_email_address     = "tf-test@example.com"
  enable_sip                      = true
  enable_h323                     = true
  enable_infinity_connect_non_sso = true
  enable_infinity_connect_sso     = true
  enable_standard_sso             = true
  sso_identity_provider_group     = pexip_infinity_identity_provider_group.tf-test-identity-provider-group.id
  tag                             = "tf-test-tag"
  sync_tag                        = "tf-test-sync-tag"
}
