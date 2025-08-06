/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "google_compute_address" "infinity_node_public_ip" {
  name   = "${local.hostname}-public-ip"
  region = var.location
}

resource "google_compute_address" "infinity_node_private_ip" {
  name         = "${local.hostname}-private-ip"
  subnetwork   = var.private_subnetwork_id
  address_type = "INTERNAL"
  region       = var.location
}
