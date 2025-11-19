/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "google_compute_firewall" "allow_ssh" {
  name    = "allow-ssh-${var.project_id}"
  network = data.google_compute_network.default.name

  allow {
    protocol = "tcp"
    ports    = ["22"]
  }

  source_ranges = ["35.235.240.0/20"] # Allow SSH from GCP console
  target_tags = ["allow-ssh-${var.project_id}"]
}

resource "google_compute_firewall" "allow_https" {
  name    = "allow-https-${var.project_id}"
  network = data.google_compute_network.default.name

  allow {
    protocol = "tcp"
    ports    = ["443"]
  }

  source_ranges = ["0.0.0.0/0"] # Allow from anywhere
  target_tags = ["allow-https-${var.project_id}"]
}

resource "google_compute_firewall" "allow_inter_node" {
  name    = "allow-inter-node-${var.project_id}"
  network = data.google_compute_network.default.name

  allow {
    protocol = "udp"
    ports    = ["500"]
  }

  allow {
    protocol = "esp"
  }

  source_ranges = [data.google_compute_subnetwork.default.ip_cidr_range]
  target_tags   = ["allow-inter-node-${var.project_id}"]
}