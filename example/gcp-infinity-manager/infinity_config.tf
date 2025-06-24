resource "pexip_infinity_dns_server" "pexip-dns-cloudflare" {
  address     = "1.1.1.1"
  description = "Cloud Flare DNS"

  depends_on = [
    google_compute_instance.infinity_manager,
    null_resource.wait_for_infinity_manager_http
  ]
}