/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

locals {
  hostname         = "${var.environment}-worker-${var.index}"
  check_status_url = "https://${var.mgr_public_ip}/api/admin/status/v1/worker_vm/${pexip_infinity_worker_vm.worker.resource_id}/"
  user_data        = jsonencode(pexip_infinity_worker_vm.worker.config)
}

resource "pexip_infinity_ssh_password_hash" "default" {
  password = var.password
}

resource "pexip_infinity_worker_vm" "worker" {
  name               = local.hostname
  hostname           = local.hostname
  address            = openstack_networking_port_v2.infinity-cnf-port.all_fixed_ips.0
  netmask            = cidrnetmask(data.openstack_networking_subnet_v2.cnf-private-subnet.cidr)
  gateway            = data.openstack_networking_subnet_v2.cnf-private-subnet.gateway_ip
  domain             = var.domain
  password           = pexip_infinity_ssh_password_hash.default.hash
  node_type          = var.node_type
  system_location    = var.system_location
  transcoding        = var.transcoding
  tls_certificate    = var.tls_certificate
  static_nat_address = openstack_networking_floatingip_v2.infinity-cnf-fip.address
}

resource "openstack_compute_instance_v2" "infinity" {
  name        = local.hostname
  flavor_name = var.flavor_name
  user_data   = "{\"conferencing_node_config\":${local.user_data}}"

  block_device {
    uuid                  = var.image_id
    source_type           = "image"
    destination_type      = "volume"
    volume_size           = 49
    boot_index            = 0
    delete_on_termination = true
  }

  network {
    port = openstack_networking_port_v2.infinity-cnf-port.id
  }

  lifecycle {
    ignore_changes = [
      user_data,
    ]
  }
}

resource "null_resource" "wait_for_infinity_node_http" {
  depends_on = [openstack_compute_instance_v2.infinity]

  # Re‑run this null_resource whenever the instance is replaced
  triggers = {
    instance_id = openstack_compute_instance_v2.infinity.id
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