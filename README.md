# Terraform Provider for Pexip Infinity

[![Build Status](https://github.com/pexip/terraform-provider-pexip/actions/workflows/test.yml/badge.svg)](https://github.com/pexip/terraform-provider-pexip/actions/workflows/test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/pexip/terraform-provider-pexip)](https://goreportcard.com/report/github.com/pexip/terraform-provider-pexip)

This Terraform provider enables you to manage [Pexip Infinity](https://www.pexip.com/products/infinity/) infrastructure using Infrastructure as Code. Automate the provisioning and management of Pexip Infinity components including manager configurations and worker nodes.

## Features

- **Manager Configuration**: Generate bootstrap configurations for Pexip Infinity Manager
- **Node Management**: Register and manage Pexip Infinity worker nodes  
- **Infrastructure as Code**: Version control your Pexip infrastructure
- **Terraform Integration**: Native Terraform resource lifecycle management

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.21 (for development)
- Pexip Infinity Manager with API access

## Installation

### Terraform Registry (Recommended)

```hcl
terraform {
  required_providers {
    pexip = {
      source  = "pexip.com/pexip/pexip"
      version = "~> 0.1"
    }
  }
}
```

### Manual Installation

1. Download the latest release from [GitHub Releases](https://github.com/pexip/terraform-provider-pexip/releases)
2. Extract the binary to your Terraform plugins directory
3. Configure Terraform to use the local provider

## Quick Start

### Basic Configuration

```hcl
terraform {
  required_providers {
    pexip = {
      source  = "pexip.com/pexip/pexip"
      version = "~> 0.1"
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
  name     = "pexip-worker-01"
  hostname = "pexip-worker-01"
}

resource "pexip_infinity_node" "worker_02" {
  name     = "pexip-worker-02"
  hostname = "pexip-worker-02"
}
```

### Variables Example

```hcl
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

## Provider Configuration

### Authentication

The provider supports the following authentication methods:

```hcl
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

### Provider Configuration Reference

| Argument | Description | Required | Environment Variable |
|----------|-------------|----------|---------------------|
| `address` | URL of the Pexip Infinity Manager API | Yes | `PEXIP_ADDRESS` |
| `username` | Username for authentication | Yes | `PEXIP_USERNAME` |
| `password` | Password for authentication | Yes | `PEXIP_PASSWORD` |

## Resources and Data Sources

### Data Sources

#### `pexip_infinity_manager_config`

Generates bootstrap configuration for Pexip Infinity Manager.

**Example:**

```hcl
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
```

**Attributes:**
- `hostname` (Required) - Manager hostname
- `domain` (Required) - DNS domain
- `ip` (Required) - IP address
- `mask` (Required) - Subnet mask
- `gw` (Required) - Gateway IP
- `dns` (Required) - DNS server IP
- `ntp` (Required) - NTP server
- `user` (Required) - Username
- `pass` (Required, Sensitive) - Password
- `admin_password` (Required, Sensitive) - Admin password
- `error_reports` (Optional) - Enable error reporting
- `enable_analytics` (Optional) - Enable analytics
- `contact_email_address` (Required) - Contact email

**Computed Attributes:**
- `rendered` - Generated configuration JSON
- `id` - CRC32 checksum of configuration

### Resources

#### `pexip_infinity_node`

Manages Pexip Infinity worker nodes.

**Example:**

```hcl
resource "pexip_infinity_node" "worker" {
  name                    = "worker-node-01"
  hostname                = "worker-node-01"
  address                 = "192.168.1.101"
  netmask                 = "255.255.255.0"
  domain                  = "example.com"
  gateway                 = "192.168.1.1"
  password                = var.node_password
  node_type               = "CONFERENCING"
  system_location         = "Data Center 1"
  maintenance_mode        = false
  transcoding             = true
}
```

**Arguments:**
- `name` (Required) - Node name
- `hostname` (Required) - Node hostname
- `address` (Optional) - IP address of the node
- `netmask` (Optional) - Network mask
- `domain` (Optional) - DNS domain
- `gateway` (Optional) - Gateway IP address
- `password` (Required, Sensitive) - Node password
- `node_type` (Optional) - Node type (default: "worker")
- `system_location` (Optional) - Physical location description
- `maintenance_mode` (Optional) - Enable maintenance mode (default: false)
- `maintenance_mode_reason` (Optional) - Reason for maintenance mode
- `transcoding` (Optional) - Enable transcoding (default: true)
- `vm_cpu_count` (Optional) - Number of CPUs allocated to VM
- `vm_system_memory` (Optional) - System memory allocated to VM

**Computed Attributes:**
- `id` - Node ID
- `config` - Generated node configuration

**Import:**

```bash
terraform import pexip_infinity_node.worker 123
```

## Complete Example

For a comprehensive example that includes Google Cloud Platform infrastructure deployment along with Pexip Infinity resources, see the [`example/`](./example/) directory in this repository. This example demonstrates:

- GCP VM instances for Pexip Infinity Manager and worker nodes
- Network configuration and firewall rules
- DNS setup and SSL certificates
- Service account and IAM permissions
- Integration between cloud infrastructure and Pexip resources

To run the example:

```bash
cd example/
cp terraform.tfvars.example terraform.tfvars
# Edit terraform.tfvars with your values
terraform init
terraform plan
terraform apply
```

## Advanced Usage

### Multiple Manager Configurations

```hcl
# Production manager
data "pexip_infinity_manager_config" "prod" {
  hostname = "prod-manager"
  domain   = "prod.company.com"
  # ... other configuration
}

# Development manager  
data "pexip_infinity_manager_config" "dev" {
  hostname = "dev-manager"
  domain   = "dev.company.com"
  # ... other configuration
}

# Production nodes
resource "pexip_infinity_node" "prod_workers" {
  count           = 3
  name            = "prod-worker-${count.index + 1}"
  hostname        = "prod-worker-${count.index + 1}"
  node_type       = "CONFERENCING"
  system_location = "Production"
}

# Development nodes
resource "pexip_infinity_node" "dev_workers" {
  count           = 1
  name            = "dev-worker-${count.index + 1}"
  hostname        = "dev-worker-${count.index + 1}"
  node_type       = "CONFERENCING"
  system_location = "Development"
}
```

### Using with Modules

```hcl
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

data "pexip_infinity_manager_config" "config" {
  hostname = "${var.environment}-manager"
  domain   = "${var.environment}.company.com"
  # ... other configuration
}

resource "pexip_infinity_node" "workers" {
  count           = var.node_count
  name            = "${var.environment}-worker-${count.index + 1}"
  hostname        = "${var.environment}-worker-${count.index + 1}"
  node_type       = "CONFERENCING"
  system_location = var.environment
}
```

```hcl
# main.tf
module "production" {
  source = "./modules/pexip-environment"
  
  environment = "prod"
  node_count  = 5
}

module "staging" {
  source = "./modules/pexip-environment"
  
  environment = "staging"
  node_count  = 2
}
```

## Development

### Prerequisites

- [Go](https://golang.org/doc/install) >= 1.21
- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Make](https://www.gnu.org/software/make/)

### Building from Source

```bash
git clone https://github.com/pexip/terraform-provider-pexip.git
cd terraform-provider-pexip
make build
```

### Local Development Setup

1. Create a `.terraformrc` file in your home directory:

```hcl
provider_installation {
  dev_overrides {
    "pexip.com/pexip/pexip" = "</path/to/your/terraform-plugins>/pexip.com/pexip/pexip"
  }
  filesystem_mirror {
    path = "</path/to/your/terraform-plugins>"
    include = ["pexip.com/*/*"]
  }
  direct {
    exclude = ["pexip.com/*/*"]
  }
}
```

2. Build and install locally:

```bash
make install
```

### Running Tests

```bash
# Unit tests
make test

# Acceptance tests (requires running Pexip environment)
make testacc

# Linting
make lint

# Format code
make fmt
```

### Testing with Real Infrastructure

For acceptance tests, you'll need access to a Pexip Infinity Manager. Set these environment variables:

```bash
export TF_ACC=1
export PEXIP_ADDRESS="https://your-manager.example.com"
export PEXIP_USERNAME="admin"
export PEXIP_PASSWORD="your-password"
```

## Troubleshooting

### Common Issues

**Authentication Errors**
- Verify your Pexip Manager URL, username, and password
- Ensure the API is accessible from your machine
- Check that your user has appropriate permissions

**SSL/TLS Errors**
- Verify your Pexip Manager uses a valid SSL certificate
- For self-signed certificates, you may need to configure your system's trust store

**Network Connectivity**
- Ensure your machine can reach the Pexip Manager on the configured port (typically 443)
- Check firewall rules and network connectivity

### Debug Logging

Enable debug logging for troubleshooting:

```bash
export TF_LOG=DEBUG
terraform plan
```

### Getting Help

- [Pexip Documentation](https://docs.pexip.com/)
- [Terraform Documentation](https://developer.hashicorp.com/terraform/docs)
- [GitHub Issues](https://github.com/pexip/terraform-provider-pexip/issues)

## Version Compatibility

| Provider Version | Terraform Version | Pexip Infinity Version | Go Version |
|------------------|-------------------|------------------------|------------|
| `~> 0.1` | `>= 1.0` | `>= v37` | `>= 1.21` |

## Known Limitations

- Provider currently supports basic node management operations
- Advanced Pexip Infinity features may require manual configuration
- SSL certificate validation is enforced (use proper certificates in production)

## Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass (`make test && make testacc`)
6. Run linting (`make lint && make fmt`)
7. Commit your changes (`git commit -m 'Add amazing feature'`)
8. Push to the branch (`git push origin feature/amazing-feature`)
9. Submit a pull request

### Development Guidelines

- Follow [Go best practices](https://golang.org/doc/effective_go.html)
- Use the [Terraform Plugin Framework](https://developer.hashicorp.com/terraform/plugin/framework)
- Write comprehensive tests for all new features
- Update documentation for any changes
- Ensure backward compatibility when possible

## License

For project license see the [LICENSE](LICENSE) file for details.

## Changelog

See [Releases](https://github.com/pexip/terraform-provider-pexip/releases) for release notes and version history.

## Security

For security concerns, please email security@pexip.com instead of using the public issue tracker.