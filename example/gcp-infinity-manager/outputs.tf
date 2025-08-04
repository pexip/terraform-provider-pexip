output "hostname" {
  value = local.hostname
}

output "user_data" {
  value = data.pexip_infinity_manager_config.conf.rendered
}

output "check_status_url" {
  value = local.check_status_url
}