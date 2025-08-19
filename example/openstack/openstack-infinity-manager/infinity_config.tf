/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_tls_certificate" "tls-cert-test" {
  certificate = tls_self_signed_cert.manager_cert.cert_pem
  private_key = tls_private_key.manager_private_key.private_key_pem
  nodes       = ["${local.hostname}.${var.domain}"]

  depends_on = [
    openstack_compute_instance_v2.infinity_manager,
    null_resource.wait_for_infinity_manager_http
  ]
}

// specify a TLS cert and private key
/*
resource "pexip_infinity_tls_certificate" "tls-cert-test2" {
  certificate = file("test-cert.pem")
  private_key = file("test-key.key")

  depends_on = [
    openstack_compute_instance_v2.infinity_manager,
    null_resource.wait_for_infinity_manager_http
  ]
}
*/

resource "pexip_infinity_licence" "license" {
  entitlement_id = var.license_key

  depends_on = [
    openstack_compute_instance_v2.infinity_manager,
    null_resource.wait_for_infinity_manager_http
  ]
}