data "google_compute_network" "default" {
  name = "default"
}

data "google_compute_subnetwork" "default" {
  name   = "default"
  region = var.location
}

resource "google_compute_address" "infinity_manager_static_ip" {
  name   = "infinity-manager-static-ip"
  region = var.location
}

resource "google_compute_address" "infinity_workers_static_ip" {
  count  = var.infinity_node_count
  name   = "infinity-worker-${count.index + 1}-static-ip"
  region = var.location
}