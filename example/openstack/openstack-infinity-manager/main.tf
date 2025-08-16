/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

locals {
  hostname         = "${var.environment}-manager"
  check_status_url = "https://${openstack_networking_floatingip_v2.infinity-mgr-fip.address}/admin/login/"
  user_data = jsonencode(data.pexip_infinity_manager_config.conf.management_node_config)
}

resource "pexip_infinity_ssh_password_hash" "default" {
  password = var.admin_password
}

resource "pexip_infinity_web_password_hash" "default" {
  password = var.password
}

data "pexip_infinity_manager_config" "conf" {
  hostname              = local.hostname
  domain                = var.domain
  ip                    = openstack_networking_port_v2.infinity-mgr-port.fixed_ip[0].ip_address
  mask                  = cidrnetmask(data.openstack_networking_subnet_v2.mgr-private-subnet.cidr)
  gw                    = data.openstack_networking_subnet_v2.mgr-private-subnet.gateway_ip
  dns                   = var.dns_server
  ntp                   = var.ntp_server
  user                  = var.username
  pass                  = pexip_infinity_web_password_hash.default.hash
  admin_password        = pexip_infinity_ssh_password_hash.default.hash
  error_reports         = var.report_errors
  enable_analytics      = var.enable_analytics
  contact_email_address = var.contact_email_address
}

resource "openstack_compute_instance_v2" "infinity_manager" {
  name            = local.hostname
  flavor_name     = var.flavor_name
  security_groups = var.security_groups
  user_data = "{\"management_node_config\":${local.user_data}}"

  block_device {
    uuid                  = var.image_id
    source_type           = "image"
    destination_type      = "volume"
    volume_size           = 99 # adjust as needed
    boot_index            = 0
    delete_on_termination = true
  }

  network {
    port = openstack_networking_port_v2.infinity-mgr-port.id
  }
}

resource "null_resource" "wait_for_infinity_manager_http" {
  depends_on = [openstack_compute_instance_v2.infinity_manager]

  # Re‑run this null_resource whenever the instance is replaced
  triggers = {
    instance_id = openstack_compute_instance_v2.infinity_manager.id
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