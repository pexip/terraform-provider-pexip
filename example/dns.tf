data "google_dns_managed_zone" "main" {
  name        = var.dns_zone_name
}

resource "google_dns_record_set" "infinity_manager_dns" {
  name         = "${var.hostname}.${data.google_dns_managed_zone.main.dns_name}"
  managed_zone = data.google_dns_managed_zone.main.name
  type         = "A"
  ttl          = 60
  rrdatas      = [google_compute_address.infinity_manager_static_ip.address]
}