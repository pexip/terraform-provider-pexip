locals {
  manager_hostname = "${var.environment}-manager"
}

module "gcp-infinity-manager" {
  source                = "./gcp-infinity-manager"
  vm_image_name         = var.vm_image_manager_name
  machine_type          = var.infinity_manager_machine_type
  cpu_platform          = var.infinity_manager_cpu_platform
  environment           = var.environment
  project_id            = var.project_id
  location              = var.location
  service_account_email = google_service_account.infinity-sa.email
  network_id            = data.google_compute_network.default.id
  subnetwork_id         = data.google_compute_subnetwork.default.id
  dns_zone_name         = var.dns_zone_name
  tags                  = ["allow-ssh-${var.project_id}", "allow-https-${var.project_id}"]

  gateway               = data.google_compute_subnetwork.default.gateway_address
  subnetwork_mask       = cidrnetmask(data.google_compute_subnetwork.default.ip_cidr_range)
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
  index                 = count.index
  vm_image_name         = var.vm_image_node_name
  machine_type          = var.infinity_node_machine_type
  cpu_platform          = var.infinity_node_cpu_platform
  environment           = var.environment
  project_id            = var.project_id
  location              = var.location
  service_account_email = google_service_account.infinity-sa.email
  network_id            = data.google_compute_network.default.id
  subnetwork_id         = data.google_compute_subnetwork.default.id
  dns_zone_name         = var.dns_zone_name
  tags                  = ["allow-ssh-${var.project_id}", "allow-https-${var.project_id}"]

  gateway         = data.google_compute_subnetwork.default.gateway_address
  subnetwork_mask = cidrnetmask(data.google_compute_subnetwork.default.ip_cidr_range)
  password        = var.infinity_password
  node_type       = "CONFERENCING"
  system_location = "OSL"
}
