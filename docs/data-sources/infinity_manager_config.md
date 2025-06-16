---
page_title: "pexip_infinity_manager_config Data Source - terraform-provider-pexip"
subcategory: ""
description: |-
  Generate bootstrap configuration for Pexip Infinity Manager.
---

# pexip_infinity_manager_config (Data Source)

Generates bootstrap configuration for Pexip Infinity Manager. This data source creates a JSON configuration that can be used to bootstrap a new Pexip Infinity Manager instance.

## Example Usage

```terraform
data "pexip_infinity_manager_config" "config" {
  hostname              = "manager-01"
  domain                = "example.com"
  ip                    = "192.168.1.100"
  mask                  = "255.255.255.0"
  gw                    = "192.168.1.1"
  dns                   = "8.8.8.8"
  ntp                   = "pool.ntp.org"
  user                  = "admin"
  pass                  = var.manager_password
  admin_password        = var.admin_password
  error_reports         = false
  enable_analytics      = false
  contact_email_address = "admin@example.com"
}

# Use the rendered configuration
output "manager_config" {
  value     = data.pexip_infinity_manager_config.config.rendered
  sensitive = true
}
```

## Schema

### Required

- `hostname` (String) - Pexip Infinity Manager hostname, e.g. `manager-1`
- `domain` (String) - Pexip Infinity Manager domain, e.g. `example.com`
- `ip` (String) - Pexip Infinity Manager IP address
- `mask` (String) - Pexip Infinity Manager subnet mask (e.g. 255.255.255.0)
- `gw` (String) - Pexip Infinity Manager gateway IP address
- `dns` (String) - Pexip Infinity Manager DNS server IP address
- `ntp` (String) - Pexip Infinity Manager NTP server
- `user` (String) - Pexip Infinity Manager username for authentication
- `pass` (String, Sensitive) - Pexip Infinity Manager password for authentication
- `admin_password` (String, Sensitive) - Pexip Infinity Manager admin password for authentication
- `contact_email_address` (String) - Pexip Infinity Manager contact email address for notifications

### Optional

- `error_reports` (Boolean) - Pexip Infinity Manager error reports. Defaults to `false`.
- `enable_analytics` (Boolean) - Pexip Infinity Manager enable analytics. Defaults to `false`.

### Read-Only

- `id` (String) - CRC-32 checksum of `rendered` Pexip Infinity bootstrap config.
- `rendered` (String) - Rendered Pexip Infinity Manager bootstrap configuration in JSON format.

## Usage Notes

- The `rendered` output contains sensitive information and should be handled securely
- The generated configuration is in JSON format suitable for Pexip Infinity Manager bootstrap
- The `id` is a CRC-32 checksum that changes when any input parameters change
- Use variables for sensitive values like passwords rather than hardcoding them

## Security Considerations

- Always use Terraform variables for sensitive values like `pass` and `admin_password`
- Mark outputs containing the rendered configuration as `sensitive = true`
- Store Terraform state securely as it will contain sensitive configuration data
- Use strong passwords for both user and admin accounts