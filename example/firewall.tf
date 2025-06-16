resource "google_compute_firewall" "allow_ssh" {
  name    = "allow-ssh"
  network = data.google_compute_network.default.name

  allow {
    protocol = "tcp"
    ports    = ["22"]
  }

  #source_ranges = ["0.0.0.0/0"] # Allow SSH from anywhere
  source_ranges = ["35.235.240.0/20"] # Allow SSH from GCP console
  target_tags   = ["allow-ssh"]
}

resource "google_compute_firewall" "allow_https" {
  name    = "allow-https"
  network = data.google_compute_network.default.name

  allow {
    protocol = "tcp"
    ports    = ["443"]
  }

  #source_ranges = ["0.0.0.0/0"] # Allow from anywhere
  source_ranges = ["35.235.240.0/20"] # Allow SSH from GCP console
  target_tags   = ["allow-https"]
}
