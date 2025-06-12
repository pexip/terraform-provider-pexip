data "google_compute_image" "pexip-infinity-image" {
  name = var.vm_image_name
  project = var.vm_image_project
}

# https://github.com/hashicorp/terraform-provider-cloudinit/blob/main/internal/provider/data_source_cloudinit_config.go
data "infinity_manager_config" conf {

}

# data "cloudinit_config" "conf" {
#   gzip          = false
#   base64_encode = false
#
#   part {
#     content_type = "text/cloud-config"
#     content = templatefile("cloud-init.yml.tmpl", {
#       caddy_domain = "${var.hostname}.${trimsuffix(data.google_dns_managed_zone.main.dns_name, ".")}"
#       caddy_acme_ca   = var.caddy_acme_ca
#       caddy_acme_email = var.caddy_acme_email
#       caddy_use_custom_tls = var.caddy_use_custom_tls
#       oauth2_proxy_email_domain = var.oauth2_proxy_email_domain
#       oauth2_proxy_auth_provider = var.oauth2_proxy_auth_provider
#       oauth2_proxy_oidc_issuer_url = var.oauth2_proxy_oidc_issuer_url
#       oauth2_proxy_client_id = var.oauth2_proxy_client_id
#       oauth2_proxy_client_secret = google_secret_manager_secret.oauth2_proxy_client_secret.name
#       oauth2_proxy_cookie_secret = google_secret_manager_secret.oauth2_proxy_cookie_secret.name
#       backup_bucket_name = google_storage_bucket.pexip_backups.name
#       restore_on_first_boot = var.restore_on_first_boot
#       vector_config = file("${path.module}/vector.yml")
#       prometheus_config = file("${path.module}/prometheus.yml")
#       blackbox_config = file("${path.module}/blackbox.yml")
#       pexip_metric_config = file("${path.module}/pexip-metric-exporter.yml")
#     })
#     filename = "cloud-init.yml"
#   }
# }

resource "local_file" "infinity_manager_config" {
  file_permission = "0640"
  filename        = "${path.module}/infinity-manager.conf"
  content         = data.infinity_manager_config.conf.rendered
}

resource "random_string" "disk_encryption_key" {
  length   = 32
  special  = true
  upper    = true
  lower    = true
  numeric  = true
}

resource "google_compute_instance" "infinity_manager" {
  name             = var.hostname
  zone             = "${var.location}-a"
  machine_type     = "n2d-standard-4"
  min_cpu_platform = "AMD Milan"

  metadata = {
    user-data = data.infinity_manager_config.conf.rendered
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