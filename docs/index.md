---
page_title: "Pexip Provider"
subcategory: ""
description: |-
  The Pexip provider enables you to manage Pexip Infinity infrastructure using Infrastructure as Code.
---

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
  username = "admin"                        # Required
  password = "secure_password"              # Required, use variables
}
```

### Environment Variables

You can also configure the provider using environment variables:

```bash
export PEXIP_ADDRESS="https://manager.example.com"
export PEXIP_USERNAME="admin"
export PEXIP_PASSWORD="secure_password"
```

## Example Usage

### Basic Configuration

```terraform
terraform {
  required_providers {
    pexip = {
      source  = "pexip/pexip"
      version = "~> 1.0"
    }
  }
}

# Configure the Pexip Provider
provider "pexip" {
  address  = "https://manager.example.com"
  username = var.pexip_username
  password = var.pexip_password
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

# Register worker nodes
resource "pexip_infinity_node" "worker_01" {
  name   = "pexip-worker-01"
  config = data.pexip_infinity_manager_config.primary.rendered
}

resource "pexip_infinity_node" "worker_02" {
  name   = "pexip-worker-02"
  config = data.pexip_infinity_manager_config.primary.rendered
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

# Production nodes
resource "pexip_infinity_node" "prod_workers" {
  count  = 3
  name   = "prod-worker-${count.index + 1}"
  config = data.pexip_infinity_manager_config.prod.rendered
}

# Development nodes
resource "pexip_infinity_node" "dev_workers" {
  count  = 1
  name   = "dev-worker-${count.index + 1}"
  config = data.pexip_infinity_manager_config.dev.rendered
}
```

## Schema

### Provider Schema

- `address` (String, Required) - URL of the Infinity Manager API (e.g., `https://infinity.example.com`)
- `username` (String, Required) - Pexip Infinity Manager username for authentication
- `password` (String, Required, Sensitive) - Pexip Infinity Manager password for authentication

## Resources and Data Sources

This provider includes the following resources and data sources:

### Data Sources

- [`pexip_infinity_manager_config`](data-sources/infinity_manager_config.md) - Generate bootstrap configuration for Pexip Infinity Manager

### Resources

- [`pexip_infinity_node`](resources/infinity_node.md) - Manage Pexip Infinity worker nodes

## Common Issues

### Authentication Errors
- Verify your Pexip Manager URL, username, and password
- Ensure the API is accessible from your machine
- Check that your user has appropriate permissions

### SSL/TLS Errors
- Verify your Pexip Manager uses a valid SSL certificate
- For self-signed certificates, you may need to configure your system's trust store

### Network Connectivity
- Ensure your machine can reach the Pexip Manager on the configured port (typically 443)
- Check firewall rules and network connectivity

### Debug Logging

Enable debug logging for troubleshooting:

```bash
export TF_LOG=DEBUG
terraform plan
```

## Support

- **Documentation**: [Pexip Documentation](https://docs.pexip.com/)
- **Issues**: [GitHub Issues](https://github.com/pexip/terraform-provider-pexip/issues)
