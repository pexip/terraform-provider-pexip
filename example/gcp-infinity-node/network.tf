resource "google_compute_address" "infinity_node_static_ip" {
  name   = "${local.hostname}-static-ip"
  region = var.location
}