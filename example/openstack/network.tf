/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

data "openstack_networking_network_v2" "netops-private-network" {
  name = "netops-private-network"
}

data "openstack_networking_subnet_v2" "netops-private-subnet" {
  name   = "netops-private-subnet"
  network_id = data.openstack_networking_network_v2.netops-private-network.id
}
