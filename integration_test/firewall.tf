/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "google_compute_firewall" "allow_https" {
  name    = "allow-https-${var.project_id}"
  network = data.google_compute_network.default.name

  allow {
    protocol = "tcp"
    ports    = ["443"]
  }

  source_ranges = ["0.0.0.0/0"] # Allow from anywhere
  target_tags   = ["allow-https-${var.project_id}"]
}