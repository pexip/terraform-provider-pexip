/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_tls_certificate" "tls-cert-test" {
  certificate = tls_self_signed_cert.manager_cert.cert_pem
  private_key = tls_private_key.manager_private_key.private_key_pem
  nodes       = ["${local.hostname}.${local.domain}"]

  depends_on = [
    google_compute_instance.infinity_manager,
    null_resource.wait_for_infinity_manager_http
  ]
}

// specify a TLS cert and private key
/*
resource "pexip_infinity_tls_certificate" "tls-cert-test2" {
  certificate = file("gcp-infinity-manager/test-cert.pem")
  private_key = file("gcp-infinity-manager/test-key.key")

  depends_on = [
    google_compute_instance.infinity_manager,
    null_resource.wait_for_infinity_manager_http
  ]
}
*/

resource "pexip_infinity_licence" "license" {
  entitlement_id = var.license_key

  depends_on = [
    google_compute_instance.infinity_manager,
    null_resource.wait_for_infinity_manager_http
  ]
}

resource "pexip_infinity_snmp_network_management_system" "snmp-nms-test" {
  address             = var.snmp_trap_host
  port                = 161
  name                = "SNMP NMS Test"
  description         = "Test SNMP NMS"
  snmp_trap_community = "public-test-trap"

  depends_on = [
    google_compute_instance.infinity_manager,
    null_resource.wait_for_infinity_manager_http
  ]
}