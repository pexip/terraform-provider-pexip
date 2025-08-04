# Terraform Provider for Pexip Infinity

[![Build Status](https://github.com/pexip/terraform-provider-pexip/actions/workflows/test.yml/badge.svg)](https://github.com/pexip/terraform-provider-pexip/actions/workflows/test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/pexip/terraform-provider-pexip)](https://goreportcard.com/report/github.com/pexip/terraform-provider-pexip)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

The official Terraform provider for [Pexip Infinity](https://www.pexip.com/products/infinity/) enables comprehensive Infrastructure as Code management of your Pexip video conferencing platform. Manage everything from basic node configurations to advanced integrations with Microsoft 365, Google Workspace, and external authentication systems.

## Features

- **🏗️ Complete Infrastructure Management**: 80+ resources covering 90% of Pexip Infinity API capabilities
- **🔐 Security & Authentication**: LDAP, Active Directory, OAuth 2.0, SAML, and certificate management
- **🎯 Conference Management**: Virtual meeting rooms, aliases, scheduled conferences, and participant controls
- **🌐 Network Configuration**: System locations, routing rules, proxies, and DNS management
- **📊 System Administration**: Licensing, logging, diagnostics, and scaling policies
- **🔗 Enterprise Integrations**: Microsoft 365, Google Workspace, Exchange, and telehealth platforms
- **📱 Media & Content**: IVR themes, media libraries, streaming credentials
- **⚡ Infrastructure as Code**: Full lifecycle management with Terraform's plan/apply workflow

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.21 (for development)
- Pexip Infinity Manager >= v37 with API access
- Valid authentication credentials for Pexip Infinity Manager

## Installation

### Terraform Registry (Recommended)

```hcl
terraform {
  required_version = ">= 1.0"
  required_providers {
    pexip = {
      source  = "registry.terraform.io/pexip/pexip"
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

### Basic Provider Configuration

```hcl
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
  address  = "https://manager.example.com"  # Required, must be valid URL
  username = var.pexip_username             # Required, min 4 characters
  password = var.pexip_password             # Required, min 4 characters
  insecure = false                          # Optional, defaults to false
}
```

### Essential Resources Example

```hcl
# Location for organizing resources (required for worker VMs)
resource "pexip_infinity_location" "datacenter_1" {
  name        = "Datacenter-1"
  description = "Primary data center location"
}

# Conference configuration
resource "pexip_infinity_conference" "team_meeting" {
  name         = "team-meeting"
  service_type = "conference"
  description  = "Weekly team meeting room"
  pin          = "123456"
  guest_pin    = "654321"
  allow_guests = true
}

# Conference alias for easy access
resource "pexip_infinity_conference_alias" "team_meeting_alias" {
  alias       = "team@company.com"
  conference  = pexip_infinity_conference.team_meeting.name
  description = "Email-style alias for team meetings"
}

# Worker VM registration
resource "pexip_infinity_worker_vm" "worker_01" {
  name            = "worker-01"
  hostname        = "worker-01"
  domain          = "company.com"
  address         = "10.0.1.101"
  netmask         = "255.255.255.0"
  gateway         = "10.0.1.1"
  node_type       = "conferencing"
  system_location = "Datacenter-1"
  transcoding     = true
  maintenance_mode = false
}

# Global system configuration
resource "pexip_infinity_global_configuration" "main" {
  bandwidth_limit                = 10000
  conference_name_max_length     = 100
  default_conference_pin_length  = 4
  enable_chat                    = true
  enable_cloud_burst             = false
  external_policy_server         = ""
}
```

### Microsoft 365 Integration Example

```hcl
# Azure tenant configuration
resource "pexip_infinity_azure_tenant" "company" {
  tenant_id           = "12345678-1234-1234-1234-123456789012"
  application_id      = "87654321-4321-4321-4321-210987654321"
  application_secret  = var.azure_app_secret
  domain              = "company.com"
}

# Microsoft Exchange connector
resource "pexip_infinity_ms_exchange_connector" "exchange" {
  name                = "company-exchange"
  description         = "Company Exchange connector"
  server_address      = "outlook.office365.com"
  username            = var.exchange_username
  password            = var.exchange_password
  domain              = "company.com"
  use_office365       = true
}

# MJX integration for Teams interoperability
resource "pexip_infinity_mjx_integration" "teams" {
  name                        = "teams-integration"
  description                 = "Microsoft Teams integration"
  display_upcoming_meetings   = 24
  enable_non_video_meetings   = true
  enable_private_meetings     = false
  start_buffer                = 5
  end_buffer                  = 5
  ep_use_https               = true
  ep_verify_certificate      = true
  replace_subject_type       = "none"
  use_webex                  = false
}
```

## Provider Configuration

### Authentication

The provider supports basic authentication with Pexip Infinity Manager:

```hcl
provider "pexip" {
  address  = "https://manager.example.com"  # Required: Manager API endpoint
  username = "admin"                        # Required: Username with API access (min 4 chars)
  password = "secure_password"              # Required: User password (min 4 chars)
  insecure = false                          # Optional: Trust self-signed certs (default: false)
}
```

### Environment Variables

Configure the provider using environment variables for CI/CD pipelines:

```bash
export PEXIP_ADDRESS="https://manager.example.com"
export PEXIP_USERNAME="admin"
export PEXIP_PASSWORD="secure_password"
export PEXIP_INSECURE="false"  # Optional
```

### Provider Arguments Reference

| Argument | Description | Required | Environment Variable | Default |
|----------|-------------|----------|---------------------|---------|
| `address` | URL of the Pexip Infinity Manager API (must be valid URL) | Yes | `PEXIP_ADDRESS` | - |
| `username` | Username for authentication (minimum 4 characters) | Yes | `PEXIP_USERNAME` | - |
| `password` | Password for authentication (minimum 4 characters) | Yes | `PEXIP_PASSWORD` | - |
| `insecure` | Trust self-signed or invalid certificates | No | `PEXIP_INSECURE` | `false` |

## Resource Categories

The provider includes 80+ resources organized into logical categories:

### 🔐 Security & Authentication (17 resources)
- **Identity Management**: `pexip_infinity_role`, `pexip_infinity_user_group`, `pexip_infinity_end_user`
- **Directory Integration**: `pexip_infinity_ldap_sync_source`, `pexip_infinity_adfs_auth_server`, `pexip_infinity_identity_provider`
- **Certificates**: `pexip_infinity_ca_certificate`, `pexip_infinity_tls_certificate`, `pexip_infinity_ssh_authorized_key`
- **OAuth & Auth**: `pexip_infinity_oauth2_client`, `pexip_infinity_google_auth_server`

### 🎯 Conference Management (9 resources)
- **Core Conference**: `pexip_infinity_conference`, `pexip_infinity_conference_alias`
- **Participants**: `pexip_infinity_automatic_participant`, `pexip_infinity_scheduled_conference`
- **Scheduling**: `pexip_infinity_recurring_conference`, `pexip_infinity_scheduled_alias`

### 🌐 Network & Infrastructure (8 resources)
- **Locations**: `pexip_infinity_system_location`, `pexip_infinity_location`
- **Networking**: `pexip_infinity_static_route`, `pexip_infinity_dns_server`, `pexip_infinity_ntp_server`
- **Monitoring**: `pexip_infinity_snmp_network_management_system`, `pexip_infinity_smtp_server`
- **Gateways**: `pexip_infinity_h323_gatekeeper`, `pexip_infinity_gateway_routing_rule`

### 🏢 Microsoft 365 & Office Integration (4 resources)
- **Azure Integration**: `pexip_infinity_azure_tenant`, `pexip_infinity_ms_exchange_connector`
- **Teams Integration**: `pexip_infinity_mjx_endpoint`, `pexip_infinity_mjx_integration`

### 📊 System Configuration (15 resources)
- **Core System**: `pexip_infinity_global_configuration`, `pexip_infinity_system_tuneable`
- **Licensing**: `pexip_infinity_licence`, `pexip_infinity_licence_request`
- **Monitoring**: `pexip_infinity_log_level`, `pexip_infinity_syslog_server`, `pexip_infinity_diagnostic_graph`
- **Scaling**: `pexip_infinity_scheduled_scaling`, `pexip_infinity_management_vm`

### 👥 User & Device Management (5 resources)
- **Devices**: `pexip_infinity_device`, `pexip_infinity_registration`, `pexip_infinity_sip_credential`
- **Infrastructure**: `pexip_infinity_worker_vm`, `pexip_infinity_management_vm`

### 📱 Media & Content (8 resources)
- **Themes & UI**: `pexip_infinity_ivr_theme`, `pexip_infinity_webapp_branding`
- **Media Libraries**: `pexip_infinity_media_library_entry`, `pexip_infinity_media_library_playlist`
- **Streaming**: `pexip_infinity_pexip_streaming_credential`, `pexip_infinity_media_processing_server`

### 🔗 External Integrations (4 resources)
- **Web Applications**: `pexip_infinity_external_webapp_host`, `pexip_infinity_webapp_alias`
- **Authentication**: `pexip_infinity_google_auth_server`
- **Access Tokens**: `pexip_infinity_gms_access_token`

## Data Sources

### `pexip_infinity_manager_config`

Generate bootstrap configuration for new Pexip Infinity Manager installations:

```hcl
data "pexip_infinity_manager_config" "bootstrap" {
  hostname              = "manager-01"
  domain                = "company.com"
  ip                    = "10.0.1.100"
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

# Use the generated configuration
output "manager_bootstrap_config" {
  value     = data.pexip_infinity_manager_config.bootstrap.rendered
  sensitive = true
}
```

## Advanced Usage Patterns

### Enterprise Environment Setup

```hcl
# Create locations for different sites
resource "pexip_infinity_location" "locations" {
  for_each = var.office_locations
  
  name        = each.key
  description = "Office location: ${each.value.description}"
}

# Deploy worker VMs across locations
resource "pexip_infinity_worker_vm" "workers" {
  for_each = { for idx, config in var.worker_configs : "${config.location}-${idx}" => config }
  
  name            = "worker-${each.key}"
  hostname        = each.value.hostname
  domain          = "company.com"
  address         = each.value.address
  netmask         = each.value.netmask
  gateway         = each.value.gateway
  node_type       = "conferencing"
  system_location = each.value.location
  transcoding     = true
  
  depends_on = [pexip_infinity_location.locations]
}

# Conference rooms per department
resource "pexip_infinity_conference" "department_rooms" {
  for_each = var.departments
  
  name         = "${each.key}-room"
  service_type = "conference"
  description  = "${each.value.name} department meeting room"
  pin          = each.value.pin
  allow_guests = each.value.allow_guests
}
```

### Microsoft 365 Complete Integration

```hcl
# Azure tenant setup
resource "pexip_infinity_azure_tenant" "main" {
  tenant_id          = var.azure_tenant_id
  application_id     = var.azure_app_id
  application_secret = var.azure_app_secret
  domain            = var.company_domain
}

# Exchange connector for calendar integration
resource "pexip_infinity_ms_exchange_connector" "main" {
  name           = "company-exchange"
  server_address = "outlook.office365.com"
  username       = var.exchange_service_account
  password       = var.exchange_service_password
  domain         = var.company_domain
  use_office365  = true
}

# Teams interoperability
resource "pexip_infinity_mjx_integration" "teams" {
  name                      = "teams-integration"
  display_upcoming_meetings = 24
  enable_non_video_meetings = true
  start_buffer              = 5
  end_buffer               = 5
  ep_use_https             = true
  replace_subject_type     = "template"
  replace_subject_template = "Pexip Meeting: {{subject}}"
}

# MJX endpoint for Teams rooms
resource "pexip_infinity_mjx_endpoint" "teams_rooms" {
  for_each = var.teams_rooms
  
  name                    = each.key
  description             = "Teams room: ${each.value.description}"
  mjx_integration         = pexip_infinity_mjx_integration.teams.id
  username               = each.value.username
  password               = each.value.password
  domain                 = var.company_domain
}
```

### Security & Compliance Configuration

```hcl
# LDAP integration for user authentication
resource "pexip_infinity_ldap_sync_source" "corporate_ad" {
  name            = "corporate-ad"
  address         = "ldap.company.com"
  port            = 636
  use_ssl         = true
  username        = var.ldap_bind_user
  password        = var.ldap_bind_password
  base_dn         = "DC=company,DC=com"
  user_search_dn  = "OU=Users,DC=company,DC=com"
  group_search_dn = "OU=Groups,DC=company,DC=com"
}

# Certificate management
resource "pexip_infinity_ca_certificate" "company_ca" {
  name        = "company-ca"
  certificate = file("${path.module}/certificates/company-ca.pem")
}

resource "pexip_infinity_tls_certificate" "manager_cert" {
  name        = "manager-certificate"
  certificate = file("${path.module}/certificates/manager.pem")
  private_key = file("${path.module}/certificates/manager-key.pem")
}

# Role-based access control
resource "pexip_infinity_role" "meeting_admin" {
  name        = "meeting-admin"
  description = "Meeting room administrators"
  
  # Define specific permissions for this role
}

resource "pexip_infinity_user_group" "admins" {
  name        = "meeting-admins"
  description = "Meeting administration group"
  ldap_source = pexip_infinity_ldap_sync_source.corporate_ad.id
  ldap_dn     = "CN=PexipAdmins,OU=Groups,DC=company,DC=com"
}
```

### Monitoring & Diagnostics Setup

```hcl
# Syslog configuration for centralized logging
resource "pexip_infinity_syslog_server" "central_logs" {
  address  = "syslog.company.com"
  port     = 514
  protocol = "UDP"
  facility = "LOCAL0"
}

# SNMP monitoring
resource "pexip_infinity_snmp_network_management_system" "monitoring" {
  name        = "company-monitoring"
  address     = "monitoring.company.com"
  port        = 161
  community   = var.snmp_community
  description = "Company SNMP monitoring system"
}

# Log level configuration
resource "pexip_infinity_log_level" "detailed" {
  component = "web"
  level     = "INFO"
}

# Diagnostic graphs for system monitoring
resource "pexip_infinity_diagnostic_graph" "system_health" {
  name        = "system-health"
  description = "System health monitoring graph"
  graph_type  = "system"
}
```

## Import Existing Resources

Import existing Pexip Infinity resources into Terraform management:

```bash
# Import a conference by its ID
terraform import pexip_infinity_conference.existing_room 123

# Import a worker VM by its ID  
terraform import pexip_infinity_worker_vm.existing_worker 456

# Import a location by its ID
terraform import pexip_infinity_location.datacenter 789
```

## Complete Example

See the [`example/`](./example/) directory for a comprehensive deployment example that includes:

- **Infrastructure**: GCP/AWS VM instances and networking
- **Pexip Core**: Manager and worker node configuration  
- **Security**: Certificate management and authentication
- **Integrations**: Microsoft 365 and Google Workspace setup
- **Monitoring**: Logging and diagnostic configuration

```bash
cd example/
cp terraform.tfvars.example terraform.tfvars
# Edit terraform.tfvars with your environment values
terraform init
terraform plan
terraform apply
```

## Development

### Prerequisites

- [Go](https://golang.org/doc/install) >= 1.21
- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Make](https://www.gnu.org/software/make/)
- Access to a Pexip Infinity Manager for testing

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
    "pexip/pexip" = "/path/to/your/terraform-provider-pexip/dist"
  }
  direct {}
}
```

2. Build and test locally:

```bash
make build
make test
```

### Available Make Targets

```bash
make build      # Build the provider binary
make install    # Build and install to local Terraform plugins directory  
make test       # Run unit tests
make testacc    # Run acceptance tests (requires Pexip environment)
make lint       # Run linting checks
make fmt        # Format Go code
make clean      # Clean build artifacts
make check      # Run all checks (lint + test)
```

### Testing

```bash
# Unit tests
make test

# Acceptance tests (requires real Pexip environment)
export TF_ACC=1
export PEXIP_ADDRESS="https://your-manager.example.com"
export PEXIP_USERNAME="admin" 
export PEXIP_PASSWORD="your-password"
make testacc

# Test specific resource
go test -v ./internal/provider -run TestAccInfinityConference
```

## Troubleshooting

### Common Issues

**Authentication Failures**
```bash
# Verify connectivity
curl -k -u username:password https://manager.example.com/api/admin/status/v1/system_summary/

# Check credentials
export TF_LOG=DEBUG
terraform plan
```

**SSL Certificate Issues**
```hcl
# For development/testing with self-signed certificates
provider "pexip" {
  address  = "https://manager.example.com"
  username = var.username
  password = var.password
  insecure = true  # Only for development!
}
```

**Network Connectivity**
```bash
# Test network access
telnet manager.example.com 443
nslookup manager.example.com
```

**Resource Import Issues**
```bash
# Get resource ID from Pexip Manager
curl -k -u admin:password "https://manager.example.com/api/admin/configuration/v1/conference/" | jq '.objects[] | {id, name}'

# Import with correct ID
terraform import pexip_infinity_conference.room 123
```

### Debug Logging

Enable comprehensive logging for troubleshooting:

```bash
export TF_LOG=DEBUG
export TF_LOG_PATH=./terraform-debug.log
terraform plan
```

### Performance Optimization

For large environments with many resources:

```hcl
# Use for_each instead of count for better performance
resource "pexip_infinity_conference" "rooms" {
  for_each = var.conference_rooms
  
  name        = each.key
  description = each.value.description
  # ... other configuration
}

# Parallelize independent resources
resource "pexip_infinity_system_location" "locations" {
  for_each = var.locations
  # ... configuration
}

resource "pexip_infinity_worker_vm" "workers" {
  for_each = var.worker_nodes
  
  name            = each.key
  hostname        = each.value.hostname
  domain          = each.value.domain
  address         = each.value.address
  netmask         = each.value.netmask
  gateway         = each.value.gateway
  system_location = each.value.location
  # ... other configuration
}
```

## Version Compatibility

| Provider Version | Terraform Version | Pexip Infinity Version | Go Version | Status |
|------------------|-------------------|------------------------|------------|---------|
| `~> 0.1` | `>= 1.0` | `>= v37` | `>= 1.21` | Active |
| `~> 0.2` | `>= 1.0` | `>= v38` | `>= 1.21` | Planned |

## Migration Guide

### From v0.0.x to v0.1.x

- Provider source changed from `pexip.com/pexip/pexip` to `pexip/pexip`
- Added `insecure` provider argument for SSL certificate handling
- Improved resource naming consistency
- Enhanced validation and error handling

```hcl
# Before (v0.0.x)
terraform {
  required_providers {
    pexip = {
      source = "pexip.com/pexip/pexip"
    }
  }
}

# After (v0.1.x)
terraform {
  required_providers {
    pexip = {
      source  = "registry.terraform.io/pexip/pexip"
      version = "~> 0.1"
    }
  }
}
```

## Best Practices

### Security
- **Never hardcode credentials** - Use variables and environment variables
- **Use HTTPS** - Always connect to Pexip Manager over HTTPS
- **Limit permissions** - Create dedicated service accounts with minimal required permissions
- **Rotate credentials** - Regularly rotate API credentials

### Resource Organization
- **Use modules** - Organize related resources into reusable modules
- **Consistent naming** - Follow a consistent naming convention across resources
- **Tag resources** - Use descriptions to document resource purpose
- **State management** - Use remote state storage for team environments

### Performance
- **Batch operations** - Use `for_each` for creating multiple similar resources
- **Minimize dependencies** - Avoid unnecessary resource dependencies
- **Parallel execution** - Structure resources to allow parallel creation/updates

## Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes following the coding standards
4. Add comprehensive tests for new functionality
5. Ensure all tests pass (`make test && make testacc`)
6. Run code quality checks (`make lint && make fmt`)
7. Update documentation as needed
8. Commit your changes (`git commit -m 'Add amazing feature'`)
9. Push to the branch (`git push origin feature/amazing-feature`)
10. Submit a pull request

### Development Guidelines

- Follow [Go best practices](https://golang.org/doc/effective_go.html)
- Use the [Terraform Plugin Framework](https://developer.hashicorp.com/terraform/plugin/framework)
- Write comprehensive tests with >80% coverage
- Document all public APIs and configuration options
- Ensure backward compatibility when possible
- Follow semantic versioning for releases

## Support

### Community Support
- [GitHub Discussions](https://github.com/pexip/terraform-provider-pexip/discussions) - Community Q&A
- [GitHub Issues](https://github.com/pexip/terraform-provider-pexip/issues) - Bug reports and feature requests

### Documentation
- [Pexip Infinity Documentation](https://docs.pexip.com/)
- [Terraform Provider Development](https://developer.hashicorp.com/terraform/plugin)
- [Terraform Best Practices](https://developer.hashicorp.com/terraform/docs/configuration)

### Professional Support
For enterprise support, contact [Pexip Support](https://www.pexip.com/support/).

## Security

For security vulnerabilities, please email [security@pexip.com](mailto:security@pexip.com) instead of using the public issue tracker.

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [Pexip Engineering Team](https://www.pexip.com/) for developing and maintaining Pexip Infinity
- [HashiCorp](https://www.hashicorp.com/) for the Terraform Plugin Framework
- [Go Community](https://golang.org/) for the excellent programming language and ecosystem

---

**Made with ❤️ by the Pexip team**