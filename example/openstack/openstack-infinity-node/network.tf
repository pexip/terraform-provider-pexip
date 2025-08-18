/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

data "openstack_networking_network_v2" "cnf-private-network" {
  name = var.private_network_name
}

data "openstack_networking_subnet_v2" "cnf-private-subnet" {
  name = var.private_subnetwork_name
}

resource "openstack_networking_floatingip_v2" "infinity-cnf-fip" {
  pool = var.floating_ip_pool
}

resource "openstack_networking_port_v2" "infinity-cnf-port" {
  name                = local.hostname
  network_id          = data.openstack_networking_network_v2.cnf-private-network.id
  security_group_ids  = var.security_groups
  fixed_ip {
    subnet_id  = data.openstack_networking_subnet_v2.cnf-private-subnet.id
  }
}

resource "openstack_networking_floatingip_associate_v2" "infinity-cnf-fip_assoc" {
  floating_ip = openstack_networking_floatingip_v2.infinity-cnf-fip.address
  port_id     = openstack_networking_port_v2.infinity-cnf-port.id
}
