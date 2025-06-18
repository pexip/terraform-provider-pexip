output "hostname" {
  value = local.hostname
}

output "user_data" {
  value = pexip_infinity_node.worker.config
}

output "check_status_url" {
  value = local.check_status_url
}