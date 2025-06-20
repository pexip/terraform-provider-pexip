resource "google_compute_address" "infinity_manager_public_ip" {
  name   = "${local.hostname}-public-ip"
  region = var.location
}

resource "google_compute_address" "infinity_manager_private_ip" {
  name         = "${local.hostname}-private-ip"
  subnetwork   = var.private_subnetwork_id
  address_type = "INTERNAL"
  region       = var.location
  address      = var.ip_address
}
