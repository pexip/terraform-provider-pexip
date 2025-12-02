/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

locals {
  domain = trimsuffix(data.google_dns_managed_zone.main.dns_name, ".")
}

data "google_dns_managed_zone" "main" {
  name = var.dns_zone_name
}

resource "google_dns_record_set" "infinity_node_dns" {
  name         = "${local.hostname}.${local.domain}."
  managed_zone = data.google_dns_managed_zone.main.name
  type         = "A"
  ttl          = 60
  rrdatas      = [google_compute_address.infinity_node_public_ip.address]
}