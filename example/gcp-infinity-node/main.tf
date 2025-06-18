locals {
  hostname = "${var.environment}-worker-${var.index}"
}

data "google_compute_image" "pexip-infinity-node-image" {
  name    = var.vm_image_name
  project = var.vm_image_project
}

resource "pexip_infinity_node" "worker" {
  name                    = local.hostname
  hostname                = local.hostname
  address                 = google_compute_address.infinity_node_static_ip.address
  netmask                 = var.subnetwork_mask
  domain                  = local.domain
  gateway                 = var.gateway
  password                = var.password
  node_type               = var.node_type
  system_location         = var.system_location
  maintenance_mode        = var.maintenance_mode
  maintenance_mode_reason = var.maintenance_mode_reason
  transcoding             = var.transcoding
  vm_cpu_count            = var.vm_cpu_count
  vm_system_memory        = var.vm_system_memory
}

resource "random_string" "disk_encryption_key" {
  length  = 32
  special = true
  upper   = true
  lower   = true
  numeric = true
}

resource "google_compute_instance" "infinity_worker" {
  name             = local.hostname
  zone             = "${var.location}-a"
  machine_type     = var.machine_type
  min_cpu_platform = var.cpu_platform

  metadata = {
    user-data = pexip_infinity_node.worker.config
    fqdn      = "${local.hostname}.${local.domain}"
  }

  boot_disk {
    disk_encryption_key_raw = base64encode(random_string.disk_encryption_key.result)
    initialize_params {
      image = data.google_compute_image.pexip-infinity-node-image.self_link
      type  = "pd-ssd"
    }
  }

  tags = var.tags

  network_interface {
    network    = var.network_id
    subnetwork = var.subnetwork_id

    access_config {
      nat_ip = google_compute_address.infinity_node_static_ip.address
    }
  }

  shielded_instance_config {
    enable_secure_boot          = false
    enable_vtpm                 = false
    enable_integrity_monitoring = false
  }

  service_account {
    email  = var.service_account_email
    scopes = ["cloud-platform"]
  }
}