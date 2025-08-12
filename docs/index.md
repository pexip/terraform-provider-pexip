---
page_title: "Pexip Provider"
subcategory: ""
description: |-
  The Pexip provider enables you to manage Pexip Infinity infrastructure using Infrastructure as Code.
---

> **Beta Notice:** This Terraform provider is currently in Beta and is not recommended for production deployments. Please use with caution.

# Pexip Provider

The Pexip Terraform provider enables you to manage [Pexip Infinity](https://www.pexip.com/products/infinity/) infrastructure using Infrastructure as Code. Automate the provisioning and management of Pexip Infinity components including manager configurations and worker nodes.

## Features

- **Manager Configuration**: Generate bootstrap configurations for Pexip Infinity Manager
- **Node Management**: Register and manage Pexip Infinity worker nodes  
- **Infrastructure as Code**: Version control your Pexip infrastructure
- **Terraform Integration**: Native Terraform resource lifecycle management

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- Pexip Infinity Manager with API access

## Authentication

The provider supports basic authentication using username and password credentials for the Pexip Infinity Manager API.

### Provider Configuration

```terraform
provider "pexip" {
  address  = "https://manager.example.com"  # Required
  username = "admin"                        # Required, min 4 characters
  password = "secure_password"              # Required, min 4 characters, use variables
  insecure = true                           # Optional, defaults to false
}
```

### Environment Variables

You can also configure the provider using environment variables:

```bash
export PEXIP_ADDRESS="https://manager.example.com"
export PEXIP_USERNAME="admin"
export PEXIP_PASSWORD="secure_password"
export PEXIP_INSECURE="true"  # Optional: for development with self-signed certificates
```

### Provider Configuration Reference

| Argument | Description | Required | Environment Variable |
|----------|-------------|----------|---------------------|
| `address` | URL of the Pexip Infinity Manager API | Yes | `PEXIP_ADDRESS` |
| `username` | Username for authentication (minimum 4 characters) | Yes | `PEXIP_USERNAME` |
| `password` | Password for authentication (minimum 4 characters) | Yes | `PEXIP_PASSWORD` |
| `insecure` | Trust self-signed or otherwise invalid certificates | No | `PEXIP_INSECURE` |

## Example Usage

### Basic Configuration

```terraform
terraform {
  required_providers {
    pexip = {
      source  = "registry.terraform.io/pexip/pexip"
      version = "~> 0.1"
    }
  }
}

# Configure the Pexip Provider
provider "pexip" {
  address  = "https://manager.example.com"
  username = var.pexip_username
  password = var.pexip_password
  insecure = true  # Set to false in production with valid SSL
}

# Generate bootstrap configuration for Infinity Manager
data "pexip_infinity_manager_config" "primary" {
  hostname              = "pexip-mgr-01"
  domain                = "company.com"
  ip                    = "10.0.1.10"
  mask                  = "255.255.255.0"
  gw                    = "10.0.1.1"
  dns                   = "8.8.8.8"
  ntp                   = "pool.ntp.org"
  user                  = "admin"
  pass                  = var.manager_password
  admin_password        = var.admin_password
  error_reports         = false
  enable_analytics      = false
  contact_email_address = "admin@company.com"
}

# Create location first (required for worker VMs)
resource "pexip_infinity_location" "main" {
  name        = "Main Location"
  description = "Primary location for production workloads"
}

# Register worker VMs
resource "pexip_infinity_worker_vm" "worker_01" {
  name            = "pexip-worker-01"
  hostname        = "pexip-worker-01"
  domain          = "company.com"
  address         = "10.0.1.20"
  netmask         = "255.255.255.0"
  gateway         = "10.0.1.1"
  system_location = pexip_infinity_location.main.name
  
  depends_on = [pexip_infinity_location.main]
}

resource "pexip_infinity_worker_vm" "worker_02" {
  name            = "pexip-worker-02"
  hostname        = "pexip-worker-02"
  domain          = "company.com"
  address         = "10.0.1.21"
  netmask         = "255.255.255.0"
  gateway         = "10.0.1.1"
  system_location = pexip_infinity_location.main.name
  
  depends_on = [pexip_infinity_location.main]
}
```

### Using Variables for Security

```terraform
# variables.tf
variable "pexip_username" {
  description = "Pexip Infinity Manager username"
  type        = string
  sensitive   = true
}

variable "pexip_password" {
  description = "Pexip Infinity Manager password"
  type        = string
  sensitive   = true
}

variable "manager_password" {
  description = "Bootstrap password for manager"
  type        = string
  sensitive   = true
}

variable "admin_password" {
  description = "Admin password for manager"
  type        = string
  sensitive   = true
}
```

### Multiple Environment Setup

```terraform
# Production manager
data "pexip_infinity_manager_config" "prod" {
  hostname = "prod-manager"
  domain   = "prod.company.com"
  ip       = "10.0.1.10"
  mask     = "255.255.255.0"
  gw       = "10.0.1.1"
  dns      = "8.8.8.8"
  ntp      = "pool.ntp.org"
  user     = "admin"
  pass     = var.prod_manager_password
  admin_password        = var.prod_admin_password
  error_reports         = false
  enable_analytics      = false
  contact_email_address = "admin@company.com"
}

# Development manager  
data "pexip_infinity_manager_config" "dev" {
  hostname = "dev-manager"
  domain   = "dev.company.com"
  ip       = "10.0.2.10"
  mask     = "255.255.255.0"
  gw       = "10.0.2.1"
  dns      = "8.8.8.8"
  ntp      = "pool.ntp.org"
  user     = "admin"
  pass     = var.dev_manager_password
  admin_password        = var.dev_admin_password
  error_reports         = false
  enable_analytics      = false
  contact_email_address = "admin@company.com"
}

# Create locations for environments
resource "pexip_infinity_location" "production" {
  name        = "Production"
  description = "Production environment location"
}

resource "pexip_infinity_location" "development" {
  name        = "Development"
  description = "Development environment location"
}

# Production worker VMs
resource "pexip_infinity_worker_vm" "prod_workers" {
  count           = 3
  name            = "prod-worker-${count.index + 1}"
  hostname        = "prod-worker-${count.index + 1}"
  domain          = "prod.company.com"
  address         = "10.0.1.${20 + count.index}"
  netmask         = "255.255.255.0"
  gateway         = "10.0.1.1"
  node_type       = "conferencing"
  system_location = pexip_infinity_location.production.name
  
  depends_on = [pexip_infinity_location.production]
}

# Development worker VMs
resource "pexip_infinity_worker_vm" "dev_workers" {
  count           = 1
  name            = "dev-worker-${count.index + 1}"
  hostname        = "dev-worker-${count.index + 1}"
  domain          = "dev.company.com"
  address         = "10.0.2.${20 + count.index}"
  netmask         = "255.255.255.0"
  gateway         = "10.0.2.1"
  node_type       = "conferencing"
  system_location = pexip_infinity_location.development.name
  
  depends_on = [pexip_infinity_location.development]
}
```

## Complete Example

For a comprehensive example that includes Google Cloud Platform infrastructure deployment along with Pexip Infinity resources, see the `example/` directory in this repository. This example demonstrates:

- GCP VM instances for Pexip Infinity Manager and worker nodes
- Network configuration and firewall rules
- DNS setup and SSL certificates
- Service account and IAM permissions
- Integration between cloud infrastructure and Pexip resources

## Schema

### Provider Schema

- `address` (String, Required) - URL of the Infinity Manager API (e.g., `https://infinity.example.com`). Must be a valid URL.
- `username` (String, Required) - Pexip Infinity Manager username for authentication. Minimum length: 4 characters.
- `password` (String, Required, Sensitive) - Pexip Infinity Manager password for authentication. Minimum length: 4 characters.
- `insecure` (Boolean, Optional) - Trust self-signed or otherwise invalid certificates. Defaults to `false`.

## Resources and Data Sources

This provider includes the following resources and data sources:

### Data Sources

- [`pexip_infinity_manager_config`](data-sources/infinity_manager_config.md) - Generate bootstrap configuration for Pexip Infinity Manager

### Resources

- [`pexip_infinity_adfs_auth_server`](resources/infinity_adfs_auth_server.md) - Manage ADFS authentication servers
- [`pexip_infinity_authentication`](resources/infinity_authentication.md) - Manage global authentication configuration
- [`pexip_infinity_automatic_participant`](resources/infinity_automatic_participant.md) - Manage automatic participants
- [`pexip_infinity_azure_tenant`](resources/infinity_azure_tenant.md) - Manage Azure tenant configurations
- [`pexip_infinity_conference`](resources/infinity_conference.md) - Manage conference configurations
- [`pexip_infinity_conference_alias`](resources/infinity_conference_alias.md) - Manage conference aliases
- [`pexip_infinity_device`](resources/infinity_device.md) - Manage device configurations
- [`pexip_infinity_dns_server`](resources/infinity_dns_server.md) - Manage DNS server configurations
- [`pexip_infinity_end_user`](resources/infinity_end_user.md) - Manage end user accounts
- [`pexip_infinity_event_sink`](resources/infinity_event_sink.md) - Manage event sink configurations
- [`pexip_infinity_gateway_routing_rule`](resources/infinity_gateway_routing_rule.md) - Manage gateway routing rules
- [`pexip_infinity_global_configuration`](resources/infinity_global_configuration.md) - Manage global system configuration
- [`pexip_infinity_http_proxy`](resources/infinity_http_proxy.md) - Manage HTTP proxy configurations
- [`pexip_infinity_identity_provider`](resources/infinity_identity_provider.md) - Manage identity provider configurations
- [`pexip_infinity_ldap_sync_source`](resources/infinity_ldap_sync_source.md) - Manage LDAP synchronization sources
- [`pexip_infinity_licence`](resources/infinity_licence.md) - Manage license configurations
- [`pexip_infinity_location`](resources/infinity_location.md) - Manage location configurations
- [`pexip_infinity_management_vm`](resources/infinity_management_vm.md) - Manage management VM configurations
- [`pexip_infinity_mjx_integration`](resources/infinity_mjx_integration.md) - Manage MJX calendar integrations
- [`pexip_infinity_ntp_server`](resources/infinity_ntp_server.md) - Manage NTP server configurations
- [`pexip_infinity_oauth2_client`](resources/infinity_oauth2_client.md) - Manage OAuth2 client configurations
- [`pexip_infinity_policy_server`](resources/infinity_policy_server.md) - Manage policy server configurations
- [`pexip_infinity_recurring_conference`](resources/infinity_recurring_conference.md) - Manage recurring conference configurations
- [`pexip_infinity_role`](resources/infinity_role.md) - Manage role configurations
- [`pexip_infinity_scheduled_conference`](resources/infinity_scheduled_conference.md) - Manage scheduled conference configurations
- [`pexip_infinity_sip_proxy`](resources/infinity_sip_proxy.md) - Manage SIP proxy configurations
- [`pexip_infinity_smtp_server`](resources/infinity_smtp_server.md) - Manage SMTP server configurations
- [`pexip_infinity_snmp_network_management_system`](resources/infinity_snmp_network_management_system.md) - Manage SNMP network management systems
- [`pexip_infinity_ssh_authorized_key`](resources/infinity_ssh_authorized_key.md) - Manage SSH authorized keys
- [`pexip_infinity_static_route`](resources/infinity_static_route.md) - Manage static route configurations
- [`pexip_infinity_stun_server`](resources/infinity_stun_server.md) - Manage STUN server configurations
- [`pexip_infinity_syslog_server`](resources/infinity_syslog_server.md) - Manage syslog server configurations
- [`pexip_infinity_system_location`](resources/infinity_system_location.md) - Manage system location configurations
- [`pexip_infinity_teams_proxy`](resources/infinity_teams_proxy.md) - Manage Teams proxy configurations
- [`pexip_infinity_tls_certificate`](resources/infinity_tls_certificate.md) - Manage TLS certificate configurations
- [`pexip_infinity_turn_server`](resources/infinity_turn_server.md) - Manage TURN server configurations
- [`pexip_infinity_upgrade`](resources/infinity_upgrade.md) - Manage system upgrade operations
- [`pexip_infinity_user_group`](resources/infinity_user_group.md) - Manage user group configurations
- [`pexip_infinity_webapp_branding`](resources/infinity_webapp_branding.md) - Manage web app branding configurations
- [`pexip_infinity_worker_vm`](resources/infinity_worker_vm.md) - Manage worker VM configurations

## Common Issues

### Authentication Errors
- Verify your Pexip Manager URL, username, and password
- Ensure the API is accessible from your machine
- Check that your user has appropriate permissions

### SSL/TLS Errors
- Verify your Pexip Manager uses a valid SSL certificate
- For self-signed certificates in development, set `insecure = true` in the provider configuration
- For production, use proper SSL certificates and keep `insecure = false` (default)

### Network Connectivity
- Ensure your machine can reach the Pexip Manager on the configured port (typically 443)
- Check firewall rules and network connectivity

### Debug Logging

Enable debug logging for troubleshooting:

```bash
export TF_LOG=DEBUG
terraform plan
```

## Version Compatibility

| Provider Version | Terraform Version | Pexip Infinity Version | Go Version |
|------------------|-------------------|------------------------|------------|
| `~> 0.1` | `>= 1.0` | `>= v37` | `>= 1.21` |

## Known Limitations

- Some advanced Pexip Infinity features may require manual configuration
- SSL certificate validation is enforced by default (use `insecure = true` only for development)
- Provider requires Pexip Infinity version 37 or higher

## Advanced Usage

### Using with Modules

```terraform
# modules/pexip-environment/main.tf
variable "environment" {
  description = "Environment name"
  type        = string
}

variable "node_count" {
  description = "Number of worker nodes"
  type        = number
  default     = 2
}

variable "manager_password" {
  description = "Password for the Pexip manager"
  type        = string
  sensitive   = true
}

variable "admin_password" {
  description = "Admin password for the manager"
  type        = string
  sensitive   = true
}

data "pexip_infinity_manager_config" "config" {
  hostname              = "${var.environment}-manager"
  domain                = "${var.environment}.company.com"
  ip                    = "10.0.1.10"
  mask                  = "255.255.255.0"
  gw                    = "10.0.1.1"
  dns                   = "8.8.8.8"
  ntp                   = "pool.ntp.org"
  user                  = "admin"
  pass                  = var.manager_password
  admin_password        = var.admin_password
  error_reports         = false
  enable_analytics      = false
  contact_email_address = "admin@company.com"
}

# Create location for the environment
resource "pexip_infinity_location" "environment" {
  name        = var.environment
  description = "${title(var.environment)} environment location"
}

resource "pexip_infinity_worker_vm" "workers" {
  count           = var.node_count
  name            = "${var.environment}-worker-${count.index + 1}"
  hostname        = "${var.environment}-worker-${count.index + 1}"
  domain          = "${var.environment}.company.com"
  address         = "10.0.1.${20 + count.index}"
  netmask         = "255.255.255.0"
  gateway         = "10.0.1.1"
  node_type       = "conferencing"
  system_location = pexip_infinity_location.environment.name
  
  depends_on = [pexip_infinity_location.environment]
}
```

```terraform
# main.tf
module "production" {
  source = "./modules/pexip-environment"
  
  environment        = "prod"
  node_count         = 5
  manager_password   = var.prod_manager_password
  admin_password     = var.prod_admin_password
}

module "staging" {
  source = "./modules/pexip-environment"
  
  environment        = "staging"
  node_count         = 2
  manager_password   = var.staging_manager_password
  admin_password     = var.staging_admin_password
}
```

## Support

- **Documentation**: [Pexip Documentation](https://docs.pexip.com/)
- **Issues**: [GitHub Issues](https://github.com/pexip/terraform-provider-pexip/issues)
- **Security**: For security concerns, please email security@pexip.com
