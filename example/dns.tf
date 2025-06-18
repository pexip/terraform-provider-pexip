locals {
  domain = data.google_dns_managed_zone.main.dns_name
}

data "google_dns_managed_zone" "main" {
  name = var.dns_zone_name
}
