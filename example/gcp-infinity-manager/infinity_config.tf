resource "pexip_infinity_dns_server" "dns-cloudflare" {
  address     = "1.1.1.1"
  description = "Cloud Flare"

  depends_on = [
    google_compute_instance.infinity_manager,
    null_resource.wait_for_infinity_manager_http
  ]
}

resource "pexip_infinity_dns_server" "dns-google-2" {
  address     = "8.8.4.4"
  description = "Google 2"

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
  name        = "AMS"
  description = "AMS always on"
  mtu         = 1460
  dns_servers = [pexip_infinity_dns_server.dns-cloudflare.id, pexip_infinity_dns_server.dns-google-2.id]
  ntp_servers = [pexip_infinity_ntp_server.ntp1.id]

  depends_on = [
    google_compute_instance.infinity_manager,
    null_resource.wait_for_infinity_manager_http
  ]
}

resource "pexip_infinity_tls_certificate" "tls-cert-test" {
  certificate = tls_self_signed_cert.manager_cert.cert_pem
  private_key = tls_private_key.manager_private_key.private_key_pem
  nodes = ["${local.hostname}.${local.domain}"]

  depends_on = [
    google_compute_instance.infinity_manager,
    null_resource.wait_for_infinity_manager_http
  ]
}

/*
resource "pexip_infinity_tls_certificate" "tls-cert-test2" {
  certificate = file("gcp-infinity-manager/test-cert.pem")
  private_key = file("gcp-infinity-manager/test-key.key")

  depends_on = [
    google_compute_instance.infinity_manager,
    null_resource.wait_for_infinity_manager_http
  ]
}
*/