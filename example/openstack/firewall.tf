/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

# --- Internode IPSec ---
resource "openstack_networking_secgroup_v2" "infinity-internode" {
  name        = "infinity-internode"
  description = "Security group for Infinity Internode IPSec Communication"
}

resource "openstack_networking_secgroup_rule_v2" "allow-esp" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "esp"
  remote_ip_prefix  = var.internode_ip_prefix
  security_group_id = openstack_networking_secgroup_v2.infinity-internode.id
}

resource "openstack_networking_secgroup_rule_v2" "allow-udp-500" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "udp"
  port_range_min    = 500
  port_range_max    = 500
  remote_ip_prefix  = var.internode_ip_prefix
  security_group_id = openstack_networking_secgroup_v2.infinity-internode.id
}

# --- Management ---
resource "openstack_networking_secgroup_v2" "infinity-management" {
  name        = "infinity-management"
  description = "Security group for Infinity Management"
}

resource "openstack_networking_secgroup_rule_v2" "allow-mgmt-ssh" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 22
  port_range_max    = 22
  remote_ip_prefix  = var.management_ip_prefix
  security_group_id = openstack_networking_secgroup_v2.infinity-management.id
}

resource "openstack_networking_secgroup_rule_v2" "all-mgmt-https" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 443
  port_range_max    = 443
  remote_ip_prefix  = var.management_ip_prefix
  security_group_id = openstack_networking_secgroup_v2.infinity-management.id
}

# --- Signalling and Media ---
resource "openstack_networking_secgroup_v2" "infinity-signalling-media" {
  name        = "infinity-signalling-media"
  description = "Security group for Infinity Signalling and Media"
}

resource "openstack_networking_secgroup_rule_v2" "allow-https" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 443
  port_range_max    = 443
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = openstack_networking_secgroup_v2.infinity-signalling-media.id
}

resource "openstack_networking_secgroup_rule_v2" "allow-sip" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 5060
  port_range_max    = 5061
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = openstack_networking_secgroup_v2.infinity-signalling-media.id
}

resource "openstack_networking_secgroup_rule_v2" "allow-h323" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 1720
  port_range_max    = 1720
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = openstack_networking_secgroup_v2.infinity-signalling-media.id
}

resource "openstack_networking_secgroup_rule_v2" "allow-media-tcp" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 40000
  port_range_max    = 49999
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = openstack_networking_secgroup_v2.infinity-signalling-media.id
}

resource "openstack_networking_secgroup_rule_v2" "allow-media-udp" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "udp"
  port_range_min    = 40000
  port_range_max    = 49999
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = openstack_networking_secgroup_v2.infinity-signalling-media.id
}
