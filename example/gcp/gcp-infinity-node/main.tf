/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

locals {
  hostname         = "${var.environment}-worker-${var.index}"
  check_status_url = "https://${var.mgr_public_ip}/api/admin/status/v1/worker_vm/${pexip_infinity_worker_vm.worker.resource_id}/"
}

data "google_compute_image" "pexip-infinity-node-image" {
  name    = var.vm_image_name
  project = var.vm_image_project
}

resource "pexip_infinity_ssh_password_hash" "default" {
  password = var.password
}

resource "pexip_infinity_worker_vm" "worker" {
  name               = local.hostname
  hostname           = local.hostname
  address            = google_compute_address.infinity_node_private_ip.address
  netmask            = var.subnetwork_mask
  domain             = local.domain
  gateway            = var.gateway
  password           = pexip_infinity_ssh_password_hash.default.hash
  node_type          = var.node_type
  system_location    = var.system_location
  tls_certificate    = var.tls_certificate
  static_nat_address = google_compute_address.infinity_node_public_ip.address

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
    conferencing_node_config = pexip_infinity_worker_vm.worker.config
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
    network    = var.private_network_id
    subnetwork = var.private_subnetwork_id
    network_ip = google_compute_address.infinity_node_private_ip.address

    access_config {
      nat_ip = google_compute_address.infinity_node_public_ip.address
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

  lifecycle {
    ignore_changes = [
      metadata,
      shielded_instance_config
    ]
  }
}

resource "null_resource" "wait_for_infinity_node_http" {
  depends_on = [google_compute_instance.infinity_worker]

  # Re‑run this null_resource whenever the instance is replaced
  triggers = {
    instance_id = google_compute_instance.infinity_worker.id
  }

  provisioner "local-exec" {

    environment = {
      PASSWORD = var.web_password
    }

    command     = <<EOT
      "Waiting for Infinity Node to sync ..."
      for i in $(seq 1 60); do
        status=$(curl --silent --show-error --insecure --location -u ${var.web_username}:$PASSWORD ${local.check_status_url} | jq -r '.sync_status')

        if [ "$status" = "SYNCED" ]; then
          echo "Infinity Node is synced."
          exit 0
        else
          echo "Infinity Node not synced, status: $status"
        fi

        sleep 10
      done

      echo "Timed out: unable to connect to Management Node" >&2
      exit 1
    EOT
    interpreter = ["/bin/bash", "-c"]
  }
}

resource "null_resource" "remove_metadata_key" {
  depends_on = [null_resource.wait_for_infinity_node_http]

  provisioner "local-exec" {
    command = <<EOT
      gcloud compute instances remove-metadata ${google_compute_instance.infinity_worker.name} \
        --project ${var.project_id} \
        --zone ${google_compute_instance.infinity_worker.zone} \
        --keys conferencing_node_config
    EOT
  }
}