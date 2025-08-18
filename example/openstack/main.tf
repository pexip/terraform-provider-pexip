/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

locals {
  manager_hostname = "${var.environment}-manager"
}

module "openstack-infinity-manager" {
  source                = "./openstack-infinity-manager"
  flavor_name           = var.mgr_flavor_name
  environment           = var.environment
  region                = var.region
  image_id              = var.mgr_node_vm_image_id
  private_network_name  = var.private_network_name
  private_subnetwork_name = var.private_subnetwork_name_mgr
  domain                = var.domain
  management_ip_prefix  = var.management_ip_prefix
  floating_ip_pool      = var.mgr_floating_ip_pool
  internode_ip_prefix   = var.internode_ip_prefix
  security_groups       = [
    openstack_networking_secgroup_v2.infinity-management.id,
    openstack_networking_secgroup_v2.infinity-internode.id,
  ]
  license_key = var.infinity_license_key

  dns_server            = var.infinity_primary_dns_server
  ntp_server            = var.infinity_ntp_server
  username              = var.infinity_username
  password              = var.infinity_password
  admin_password        = var.infinity_password
  report_errors         = var.infinity_report_errors
  enable_analytics      = var.infinity_enable_analytics
  contact_email_address = var.infinity_contact_email_address

  depends_on = [ openstack_networking_secgroup_rule_v2.all-mgmt-https ]
}

module "openstack-infinity-node" {
  source                = "./openstack-infinity-node"
  flavor_name           = var.cnf_flavor_name
  count                 = var.infinity_node_count
  index                 = count.index + 1
  environment           = var.environment
  region                = var.region
  private_network_name  = var.private_network_name
  private_subnetwork_name = var.private_subnetwork_name_cnf
  security_groups = [
    openstack_networking_secgroup_v2.infinity-management.id,
    openstack_networking_secgroup_v2.infinity-internode.id,
    openstack_networking_secgroup_v2.infinity-signalling-media.id,
  ]
  floating_ip_pool = var.cnf_floating_ip_pool
  domain           = var.domain
  management_ip_prefix = var.management_ip_prefix
  internode_ip_prefix  = var.internode_ip_prefix
  image_id             = var.cnf_node_vm_image_id

  password        = var.infinity_password
  node_type       = "CONFERENCING"
  system_location = pexip_infinity_system_location.example-location-1.id

  depends_on = [
    module.openstack-infinity-manager
  ]
}