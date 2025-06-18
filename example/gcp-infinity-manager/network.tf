resource "google_compute_address" "infinity_manager_static_ip" {
  name   = "${local.hostname}-static-ip"
  region = var.location
}
