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
