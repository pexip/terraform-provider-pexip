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
  filename        = "${path.root}/infinity-manager.conf"
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
      type  = "pd-ssd"
    }
  }

  tags = ["allow-ssh-${var.project_id}", "allow-https-${var.project_id}"]

  network_interface {
    network    = data.google_compute_network.default.id
    subnetwork = data.google_compute_subnetwork.default.id

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
    # Google recommends custom service accounts that have cloud-appliance scope and permissions granted via IAM Roles.
    email  = google_service_account.infinity-sa.email
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
    command = <<EOT
echo "Waiting for Infinity Manager (HTTP 200 expected) ..."
for i in $(seq 1 30); do
  status=$(curl --silent --insecure --output /dev/null --write-out "%%{http_code}" \
    https://${var.hostname}.${data.google_dns_managed_zone.main.dns_name}/api/admin/health)

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

resource "pexip_infinity_node" "infinity-node-01" {
  name = "infinity-node-01"
  hostname = "infinity-node-01"

  depends_on = [
    google_compute_instance.infinity_manager,
    null_resource.wait_for_infinity_manager_http,
  ]
}