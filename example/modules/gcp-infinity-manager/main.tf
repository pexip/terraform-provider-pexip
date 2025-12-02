/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

locals {
  hostname         = "${var.environment}-manager"
  check_status_url = "https://${local.hostname}.${local.domain}/admin/login/"
}

data "google_compute_image" "pexip-infinity-manager-image" {
  name    = var.vm_image_name
  project = var.vm_image_project
}

resource "pexip_infinity_ssh_password_hash" "default" {
  password = var.admin_password
}

resource "pexip_infinity_web_password_hash" "default" {
  password = var.password
}

data "pexip_infinity_manager_config" "conf" {
  hostname              = local.hostname
  domain                = local.domain
  ip                    = google_compute_address.infinity_manager_private_ip.address
  mask                  = var.subnetwork_mask
  gw                    = var.gateway
  dns                   = var.dns_server
  ntp                   = var.ntp_server
  user                  = var.username
  pass                  = pexip_infinity_web_password_hash.default.hash
  admin_password        = pexip_infinity_ssh_password_hash.default.hash
  error_reports         = var.report_errors
  enable_analytics      = var.enable_analytics
  contact_email_address = var.contact_email_address
}

resource "random_string" "disk_encryption_key" {
  length  = 32
  special = true
  upper   = true
  lower   = true
  numeric = true
}

#tfsec:ignore:AVD-GCP-0030
#tfsec:ignore:AVD-GCP-0031
#tfsec:ignore:AVD-GCP-0037
resource "google_compute_instance" "infinity_manager" {
  name             = local.hostname
  zone             = "${var.location}-a"
  machine_type     = var.machine_type
  min_cpu_platform = var.cpu_platform


  metadata = {
    management_node_config = data.pexip_infinity_manager_config.conf.management_node_config
  }

  boot_disk {
    disk_encryption_key_raw = base64encode(random_string.disk_encryption_key.result)
    initialize_params {
      image = data.google_compute_image.pexip-infinity-manager-image.self_link
      type  = "pd-ssd"
    }
  }

  tags = var.tags

  network_interface {
    network    = var.private_network_id
    subnetwork = var.private_subnetwork_id
    network_ip = google_compute_address.infinity_manager_private_ip.address

    access_config {
      nat_ip = google_compute_address.infinity_manager_public_ip.address
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

resource "null_resource" "wait_for_infinity_manager_http" {
  depends_on = [google_compute_instance.infinity_manager]

  # Re‑run this null_resource whenever the instance is replaced
  triggers = {
    instance_id = google_compute_instance.infinity_manager.id
  }

  provisioner "local-exec" {
    command     = <<EOT
      echo "Waiting for Infinity Manager (HTTP 200 expected) ..."
      for i in $(seq 1 60); do
        status=$(curl --silent --show-error --insecure --location --output /dev/null --write-out "%%{http_code}" ${local.check_status_url})

        if [ "$status" -eq 200 ]; then
          sleep 10 # Wait for the service to stabilize
          echo "Infinity Manager is ready (HTTP 200)."
          exit 0
        fi

        sleep 10
      done

      echo "Timed out: Infinity Manager did not return HTTP 200" >&2
      exit 1
    EOT
    interpreter = ["/bin/bash", "-c"]
  }
}

resource "null_resource" "remove_metadata_key" {
  depends_on = [null_resource.wait_for_infinity_manager_http]

  provisioner "local-exec" {
    command = <<EOT
      gcloud compute instances remove-metadata ${google_compute_instance.infinity_manager.name} \
        --project ${var.project_id} \
        --zone ${google_compute_instance.infinity_manager.zone} \
        --keys management_node_config
    EOT
  }
}