resource "pexip_infinity_dns_server" "dns-cloudflare" {
  address     = "1.1.1.1"
  description = "Cloud Flare"

  depends_on = [
    google_compute_instance.infinity_manager,
    null_resource.wait_for_infinity_manager_http
  ]
}

resource "pexip_infinity_ntp_server" "ntp1" {
  address     = "1.pool.ntp.org"
  description = "1.pool.ntp.org"

  depends_on = [
    google_compute_instance.infinity_manager,
    null_resource.wait_for_infinity_manager_http
  ]
}

resource "pexip_infinity_system_location" "AMS" {
  name     = "AMS"
  description = "AMS always on"
  mtu= 1460

  depends_on = [
    google_compute_instance.infinity_manager,
    null_resource.wait_for_infinity_manager_http
  ]
}