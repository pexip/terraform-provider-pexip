/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

output "hostname" {
  value = local.hostname
}

output "user_data" {
  value = data.pexip_infinity_manager_config.conf.rendered
}

output "check_status_url" {
  value = local.check_status_url
}

output "manager_cert" {
  value = pexip_infinity_tls_certificate.tls-cert-test
}

output "mgr-public-ip" {
  value = google_compute_address.infinity_manager_public_ip.address
}