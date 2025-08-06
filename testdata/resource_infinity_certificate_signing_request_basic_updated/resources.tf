/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_certificate_signing_request" "certificate_signing_request-test" {
  subject_name                 = "certificate_signing_request-test"
  dn                           = "updated-value" // Updated value
  additional_subject_alt_names = "certificate_signing_request-test"
  private_key_type             = "ecdsa256"      // Updated value
  private_key                  = "updated-value" // Updated value
  private_key_passphrase       = "updated-value" // Updated value
  ad_compatible                = false           // Updated to false
  tls_certificate              = "updated-value" // Updated value
}