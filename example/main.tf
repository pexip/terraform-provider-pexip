data "google_compute_image" "pexip-infinity-image" {
  name    = var.vm_image_name
  project = var.vm_image_project
}

data "pexip_infinity_manager_config" "conf" {
  hostname              = var.infinity_hostname
  domain                = data.google_dns_managed_zone.main.dns_name
  ip                    = google_compute_address.infinity_manager_static_ip.address
  mask                  = "255.255.255.0"
  gw                    = data.google_compute_subnetwork.default.gateway_address
  dns                   = var.infinity_primary_dns_server
  ntp                   = var.infinity_ntp_server
  user                  = var.infinity_username
  pass                  = var.infinity_password
  admin_password        = var.infinity_password
  error_reports         = var.infinity_report_errors
  enable_analytics      = var.infinity_enable_analytics
  contact_email_address = var.infinity_contact_email_address
}

resource "local_file" "pexip_infinity_manager_config" {
  file_permission = "0640"
  filename        = "${path.module}/infinity-manager.conf"
  content         = data.pexip_infinity_manager_config.conf.rendered
}

resource "random_string" "disk_encryption_key" {
  length  = 32
  special = true
  upper   = true
  lower   = true
  numeric = true
}

resource "google_compute_instance" "infinity_manager" {
  name             = var.hostname
  zone             = "${var.location}-a"
  machine_type     = "n2d-standard-16"
  min_cpu_platform = "AMD Milan"

  metadata = {
    user-data = data.pexip_infinity_manager_config.conf.rendered
    fqdn      = "${var.hostname}.${data.google_dns_managed_zone.main.dns_name}"
  }

  boot_disk {
    disk_encryption_key_raw = base64encode(random_string.disk_encryption_key.result)
    initialize_params {
      image = data.google_compute_image.pexip-infinity-image.self_link
    }
  }

  tags = ["allow-ssh", "allow-http", "allow-https"]

  network_interface {
    network    = data.google_compute_network.default.id
    subnetwork = data.google_compute_subnetwork.default.id

    access_config {
      nat_ip = google_compute_address.infinity_manager_static_ip.address
    }
  }

  shielded_instance_config {
    enable_secure_boot          = true
    enable_vtpm                 = true
    enable_integrity_monitoring = true
  }

  service_account {
    # Google recommends custom service accounts that have cloud-appliance scope and permissions granted via IAM Roles.
    email  = google_service_account.infinity-sa.email
    scopes = ["cloud-platform"]
  }
}

resource "pexip_infinity_node" "infinity-node-01" {
  name = "infinity-node-01"
  hostname = "infinity-node-01"
}