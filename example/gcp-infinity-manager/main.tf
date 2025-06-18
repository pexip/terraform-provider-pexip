locals {
  hostname = "${var.environment}-manager"
}

data "google_compute_image" "pexip-infinity-manager-image" {
  name    = var.vm_image_name
  project = var.vm_image_project
}

data "pexip_infinity_manager_config" "conf" {
  hostname              = local.hostname
  domain                = local.domain
  ip                    = google_compute_address.infinity_manager_static_ip.address
  mask                  = var.subnetwork_mask
  gw                    = var.gateway
  dns                   = var.dns_server
  ntp                   = var.ntp_server
  user                  = var.username
  pass                  = var.password
  admin_password        = var.password
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

resource "google_compute_instance" "infinity_manager" {
  name             = local.hostname
  zone             = "${var.location}-a"
  machine_type     = var.machine_type
  min_cpu_platform = var.cpu_platform

  metadata = {
    user-data = data.pexip_infinity_manager_config.conf.rendered
    fqdn      = "${local.hostname}.${local.domain}"
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
    network    = var.network_id
    subnetwork = var.subnetwork_id

    access_config {
      nat_ip = google_compute_address.infinity_manager_static_ip.address
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

resource "null_resource" "wait_for_infinity_manager_http" {
  depends_on = [google_compute_instance.infinity_manager]

  # Re‑run this null_resource whenever the instance is replaced
  triggers = {
    instance_id = google_compute_instance.infinity_manager.id
  }

  provisioner "local-exec" {
    command     = <<EOT
echo "Waiting for Infinity Manager (HTTP 200 expected) ..."
for i in $(seq 1 30); do
  status=$(curl --silent --insecure --output /dev/null --write-out "%%{http_code}" \
    https://${local.hostname}.${local.domain}/api/admin/configuration/v1/)

  if [ "$status" -eq 200 ]; then
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
