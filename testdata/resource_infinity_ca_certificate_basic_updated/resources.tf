/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_ca_certificate" "ca_certificate-test" {
  certificate          = "updated-value" // Updated value
  trusted_intermediate = false           // Updated to false
}