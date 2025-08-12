/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

locals {
  manager_hostname = "${var.environment}-manager"
}

module "gcp-infinity-manager" {
  source                = "./gcp-infinity-manager"
  license_key           = var.infinity_license_key
  vm_image_name         = var.vm_image_manager_name
  machine_type          = var.infinity_manager_machine_type
  cpu_platform          = var.infinity_manager_cpu_platform
  environment           = var.environment
  project_id            = var.project_id
  location              = var.location
  service_account_email = google_service_account.infinity-sa.email
  private_network_id    = data.google_compute_network.default.id
  private_subnetwork_id = data.google_compute_subnetwork.default.id
  dns_zone_name         = var.dns_zone_name
  tags = concat(
    tolist(google_compute_firewall.allow_ssh.target_tags),
    tolist(google_compute_firewall.allow_https.target_tags),
    tolist(google_compute_firewall.allow_inter_node.target_tags),
  )

  ip_address            = var.infinity_ip_address
  gateway               = data.google_compute_subnetwork.default.gateway_address
  subnetwork_mask       = "255.255.255.255" // Use /32 for single IP address for GCP
  dns_server            = var.infinity_primary_dns_server
  ntp_server            = var.infinity_ntp_server
  username              = var.infinity_username
  password              = var.infinity_password
  admin_password        = var.infinity_password
  report_errors         = var.infinity_report_errors
  enable_analytics      = var.infinity_enable_analytics
  contact_email_address = var.infinity_contact_email_address
}

module "gcp-infinity-node" {
  source                = "./gcp-infinity-node"
  count                 = var.infinity_node_count
  index                 = count.index + 1
  vm_image_name         = var.vm_image_node_name
  machine_type          = var.infinity_node_machine_type
  cpu_platform          = var.infinity_node_cpu_platform
  environment           = var.environment
  project_id            = var.project_id
  location              = var.location
  service_account_email = google_service_account.infinity-sa.email
  private_network_id    = data.google_compute_network.default.id
  private_subnetwork_id = data.google_compute_subnetwork.default.id
  dns_zone_name         = var.dns_zone_name
  tags = concat(
    tolist(google_compute_firewall.allow_ssh.target_tags),
    tolist(google_compute_firewall.allow_https.target_tags),
    tolist(google_compute_firewall.allow_inter_node.target_tags)
  )

  gateway         = data.google_compute_subnetwork.default.gateway_address
  subnetwork_mask = "255.255.255.255" // Use /32 for single IP address for GCP
  password        = var.infinity_password
  node_type       = "CONFERENCING"
  system_location = pexip_infinity_system_location.example-location-1.id
  tls_certificate = module.gcp-infinity-manager.manager_cert.id
  depends_on = [
    module.gcp-infinity-manager
  ]
}