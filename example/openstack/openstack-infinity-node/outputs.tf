/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

output "hostname" {
  value = local.hostname
}

output "user_data" {
  value = pexip_infinity_worker_vm.worker.config
}

output "check_status_url" {
  value = local.check_status_url
}