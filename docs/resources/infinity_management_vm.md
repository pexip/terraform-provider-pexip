---
page_title: "pexip_infinity_management_vm Resource - terraform-provider-pexip"
subcategory: ""
description: |-
  Manages a Pexip Infinity Management VM configuration.
---

# pexip_infinity_management_vm (Resource)

Manages a Pexip Infinity Management VM configuration. Management VMs are Pexip Infinity Manager nodes that control the entire Pexip Infinity platform. They provide the administrative interface, API access, and coordination services for the conferencing infrastructure. This resource supports Create, Read, and Delete operations only - updates are not supported and require recreation of the resource.

## Example Usage

### Basic Management VM

```terraform
resource "pexip_infinity_management_vm" "primary" {
  name        = "Primary Management VM"
  description = "Primary Pexip Infinity Manager"
  
  # Network configuration
  address = "192.168.1.10"
  netmask = "255.255.255.0"
  gateway = "192.168.1.1"
  mtu     = 1500
  
  # Hostname configuration
  hostname = "pexip-mgr01"
  domain   = "company.com"
  
  # SSH configuration
  enable_ssh = "yes"
  
  # SNMP configuration
  snmp_mode = "disabled"
  
  # Management VM state
  initializing = false
}
```

### High Availability Management VM Pair

```terraform
# Primary Management VM
resource "pexip_infinity_management_vm" "primary" {
  name        = "Primary Management VM"
  description = "Primary Pexip Infinity Manager"
  
  address = "10.1.1.10"
  netmask = "255.255.255.0"
  gateway = "10.1.1.1"
  mtu     = 1500
  
  hostname = "pexip-mgr01"
  domain   = "enterprise.com"
  
  # IPv6 support
  ipv6_address = "2001:db8::10"
  ipv6_gateway = "2001:db8::1"
  
  # SSH with key-only access
  enable_ssh = "keys_only"
  ssh_authorized_keys = [
    data.pexip_infinity_ssh_key.admin_key.id
  ]
  
  # SNMP v3 configuration
  snmp_mode                     = "v3"
  snmp_username                 = "pexip-monitor"
  snmp_authentication_password  = var.snmp_auth_password
  snmp_privacy_password        = var.snmp_privacy_password
  snmp_system_contact          = "IT Operations <it-ops@enterprise.com>"
  snmp_system_location         = "Data Center 1, Rack A1"
  
  # Network services
  dns_servers = [
    data.pexip_infinity_dns_server.primary.id,
    data.pexip_infinity_dns_server.secondary.id
  ]
  
  ntp_servers = [
    data.pexip_infinity_ntp_server.primary.id,
    data.pexip_infinity_ntp_server.secondary.id
  ]
  
  # Logging and monitoring
  syslog_servers = [
    data.pexip_infinity_syslog_server.central.id
  ]
  
  event_sinks = [
    data.pexip_infinity_event_sink.analytics.id,
    data.pexip_infinity_event_sink.monitoring.id
  ]
  
  # Security
  tls_certificate = data.pexip_infinity_tls_certificate.management.id
  
  initializing = false
}

# Secondary Management VM
resource "pexip_infinity_management_vm" "secondary" {
  name        = "Secondary Management VM"
  description = "Secondary Pexip Infinity Manager for HA"
  
  address = "10.1.1.11"
  netmask = "255.255.255.0"
  gateway = "10.1.1.1"
  mtu     = 1500
  
  hostname = "pexip-mgr02"
  domain   = "enterprise.com"
  
  # IPv6 support
  ipv6_address = "2001:db8::11"
  ipv6_gateway = "2001:db8::1"
  
  # SSH with key-only access
  enable_ssh = "keys_only"
  ssh_authorized_keys = [
    data.pexip_infinity_ssh_key.admin_key.id
  ]
  
  # SNMP v3 configuration
  snmp_mode                     = "v3"
  snmp_username                 = "pexip-monitor"
  snmp_authentication_password  = var.snmp_auth_password
  snmp_privacy_password        = var.snmp_privacy_password
  snmp_system_contact          = "IT Operations <it-ops@enterprise.com>"
  snmp_system_location         = "Data Center 2, Rack B1"
  
  # Network services
  dns_servers = [
    data.pexip_infinity_dns_server.primary.id,
    data.pexip_infinity_dns_server.secondary.id
  ]
  
  ntp_servers = [
    data.pexip_infinity_ntp_server.primary.id,
    data.pexip_infinity_ntp_server.secondary.id
  ]
  
  # Logging and monitoring
  syslog_servers = [
    data.pexip_infinity_syslog_server.central.id
  ]
  
  event_sinks = [
    data.pexip_infinity_event_sink.analytics.id,
    data.pexip_infinity_event_sink.monitoring.id
  ]
  
  # Security
  tls_certificate = data.pexip_infinity_tls_certificate.management.id
  
  initializing = false
}
```

### Enterprise Management VM with NAT

```terraform
resource "pexip_infinity_management_vm" "enterprise" {
  name        = "Enterprise Management VM"
  description = "Enterprise Pexip Infinity Manager with NAT"
  
  # Network configuration with NAT
  address           = "172.16.1.10"
  netmask           = "255.255.255.0"
  gateway           = "172.16.1.1"
  static_nat_address = "203.0.113.10"
  mtu               = 1500
  
  # Hostname configuration
  hostname         = "pexip-mgr"
  domain          = "company.com"
  alternative_fqdn = "pexip-external.company.com"
  
  # IPv6 configuration
  ipv6_address = "2001:db8:100::10"
  ipv6_gateway = "2001:db8:100::1"
  
  # SSH configuration
  enable_ssh                    = "keys_only"
  ssh_authorized_keys_use_cloud = true
  
  # HTTP proxy for external access
  http_proxy = data.pexip_infinity_http_proxy.corporate.id
  
  # Advanced SNMP configuration
  snmp_mode                      = "v3"
  snmp_username                  = "enterprise-monitor"
  snmp_authentication_password   = var.snmp_auth_password
  snmp_privacy_password         = var.snmp_privacy_password
  snmp_system_contact           = "Enterprise IT <enterprise-it@company.com>"
  snmp_system_location          = "Enterprise Data Center"
  snmp_network_management_system = data.pexip_infinity_snmp_nms.enterprise.id
  
  # Comprehensive network services
  dns_servers = [
    data.pexip_infinity_dns_server.internal_primary.id,
    data.pexip_infinity_dns_server.internal_secondary.id,
    data.pexip_infinity_dns_server.external.id
  ]
  
  ntp_servers = [
    data.pexip_infinity_ntp_server.internal.id,
    data.pexip_infinity_ntp_server.pool.id
  ]
  
  syslog_servers = [
    data.pexip_infinity_syslog_server.security.id,
    data.pexip_infinity_syslog_server.operations.id
  ]
  
  # Static routing
  static_routes = [
    data.pexip_infinity_static_route.corporate_network.id,
    data.pexip_infinity_static_route.dmz_network.id
  ]
  
  # Event collection
  event_sinks = [
    data.pexip_infinity_event_sink.siem.id,
    data.pexip_infinity_event_sink.analytics.id,
    data.pexip_infinity_event_sink.billing.id
  ]
  
  # Security certificates
  tls_certificate = data.pexip_infinity_tls_certificate.wildcard.id
  
  # Secondary configuration
  secondary_config_passphrase = var.secondary_config_passphrase
  
  initializing = false
}
```

### Cloud-Based Management VM

```terraform
resource "pexip_infinity_management_vm" "cloud" {
  name        = "Cloud Management VM"
  description = "Pexip Infinity Manager in cloud environment"
  
  # Cloud network configuration
  address = "10.0.1.10"
  netmask = "255.255.255.0"
  gateway = "10.0.1.1"
  mtu     = 9000  # Jumbo frames for cloud
  
  hostname = "pexip-cloud-mgr"
  domain   = "cloud.company.com"
  
  # Cloud-specific SSH configuration
  enable_ssh                    = "keys_only"
  ssh_authorized_keys_use_cloud = true
  
  # Cloud proxy for external services
  http_proxy = data.pexip_infinity_http_proxy.cloud.id
  
  # Basic SNMP for cloud monitoring
  snmp_mode                     = "v1v2c"
  snmp_community               = var.cloud_snmp_community
  snmp_system_contact          = "Cloud Operations"
  snmp_system_location         = "AWS us-east-1"
  
  # Cloud DNS and NTP
  dns_servers = [
    data.pexip_infinity_dns_server.cloud_primary.id,
    data.pexip_infinity_dns_server.cloud_secondary.id
  ]
  
  ntp_servers = [
    data.pexip_infinity_ntp_server.cloud.id
  ]
  
  # Cloud logging
  syslog_servers = [
    data.pexip_infinity_syslog_server.cloudwatch.id
  ]
  
  # Cloud event collection
  event_sinks = [
    data.pexip_infinity_event_sink.cloud_analytics.id
  ]
  
  # Cloud certificate
  tls_certificate = data.pexip_infinity_tls_certificate.cloud.id
  
  initializing = false
}
```

### Development Management VM

```terraform
resource "pexip_infinity_management_vm" "development" {
  name        = "Development Management VM"
  description = "Development environment Pexip Infinity Manager"
  
  address = "192.168.100.10"
  netmask = "255.255.255.0"
  gateway = "192.168.100.1"
  mtu     = 1500
  
  hostname = "pexip-dev"
  domain   = "dev.company.local"
  
  # Relaxed SSH for development
  enable_ssh = "yes"
  
  # Simple SNMP for development
  snmp_mode      = "v1v2c"
  snmp_community = "public"
  
  # Minimal services for development
  dns_servers = [
    data.pexip_infinity_dns_server.dev.id
  ]
  
  ntp_servers = [
    data.pexip_infinity_ntp_server.dev.id
  ]
  
  initializing = false
}
```

## Schema

### Required

- `name` (String) - The name of the management VM. Maximum length: 250 characters.
- `address` (String) - The IP address of the management VM.
- `netmask` (String) - The network mask for the management VM.
- `gateway` (String) - The gateway IP address for the management VM.
- `hostname` (String) - The hostname of the management VM. Maximum length: 253 characters.
- `domain` (String) - The domain name for the management VM.
- `mtu` (Number) - Maximum Transmission Unit (MTU) size. Valid range: 576-9000.
- `enable_ssh` (String) - SSH access configuration. Valid values: `yes`, `no`, `keys_only`.
- `snmp_mode` (String) - SNMP mode configuration. Valid values: `disabled`, `v1v2c`, `v3`.
- `initializing` (Boolean) - Whether the management VM is in initializing state.

### Optional

- `description` (String) - Description of the management VM. Maximum length: 500 characters.
- `alternative_fqdn` (String) - Alternative fully qualified domain name for the management VM.
- `ipv6_address` (String) - The IPv6 address of the management VM.
- `ipv6_gateway` (String) - The IPv6 gateway for the management VM.
- `static_nat_address` (String) - Static NAT address for the management VM.
- `dns_servers` (List of String) - List of DNS server URIs for the management VM.
- `ntp_servers` (List of String) - List of NTP server URIs for the management VM.
- `syslog_servers` (List of String) - List of syslog server URIs for the management VM.
- `static_routes` (List of String) - List of static route URIs for the management VM.
- `event_sinks` (List of String) - List of event sink URIs for the management VM.
- `http_proxy` (String) - HTTP proxy URI for the management VM.
- `tls_certificate` (String) - TLS certificate URI for the management VM.
- `ssh_authorized_keys` (List of String) - List of SSH authorized key URIs for the management VM.
- `ssh_authorized_keys_use_cloud` (Boolean) - Whether to use cloud-based SSH authorized keys.
- `secondary_config_passphrase` (String, Sensitive) - Secondary configuration passphrase.
- `snmp_community` (String, Sensitive) - SNMP community string.
- `snmp_username` (String) - SNMP username for v3 authentication.
- `snmp_authentication_password` (String, Sensitive) - SNMP authentication password.
- `snmp_privacy_password` (String, Sensitive) - SNMP privacy password.
- `snmp_system_contact` (String) - SNMP system contact information.
- `snmp_system_location` (String) - SNMP system location information.
- `snmp_network_management_system` (String) - SNMP network management system URI.

### Read-Only

- `id` (String) - Resource URI for the management VM in Infinity.
- `resource_id` (Number) - The resource integer identifier for the management VM in Infinity.
- `primary` (Boolean) - Whether this is the primary management VM.

## Import

Import is supported using the following syntax:

```shell
terraform import pexip_infinity_management_vm.example 123
```

Where `123` is the numeric resource ID of the management VM.

## Usage Notes

### Update Limitations

Management VM resources **do not support update operations**. If you need to change any configuration parameters, you must delete and recreate the resource. This is due to the critical nature of management VMs in the Pexip Infinity infrastructure.

### High Availability

- Deploy multiple management VMs for high availability
- Only one management VM can be primary at a time
- Secondary management VMs provide failover capability
- Ensure network connectivity between all management VMs

### Network Configuration

- **IPv4**: Required for basic connectivity
- **IPv6**: Optional dual-stack support
- **NAT**: Use static_nat_address for external access
- **MTU**: Consider jumbo frames (9000) for high-performance networks

### SSH Access Configuration

- **yes**: Allow SSH with password and key authentication
- **no**: Disable SSH access entirely
- **keys_only**: Allow SSH with key authentication only (recommended)

### SNMP Configuration

- **disabled**: No SNMP access
- **v1v2c**: SNMP versions 1 and 2c with community strings
- **v3**: SNMP version 3 with user-based security

### Certificate Management

- Use TLS certificates for secure web interface access
- Certificates should include management VM hostname and alternative FQDN
- Ensure certificate trust chain is properly configured

### Cloud Integration

- Enable ssh_authorized_keys_use_cloud for cloud-based key management
- Configure HTTP proxy for external service access
- Use appropriate MTU settings for cloud environments

### Monitoring and Logging

- Configure syslog servers for centralized logging
- Set up event sinks for real-time monitoring
- Use SNMP for infrastructure monitoring integration

## Troubleshooting

### Common Issues

**Management VM Creation Fails**
- Verify network configuration (address, netmask, gateway)
- Ensure hostname and domain follow proper naming conventions
- Check that MTU value is within valid range (576-9000)
- Verify all referenced URIs (DNS servers, NTP servers, etc.) exist

**Network Connectivity Issues**
- Verify IP address is not in use by another device
- Check that gateway is reachable and properly configured
- Ensure netmask is correct for the network segment
- Test DNS resolution for the specified domain

**SSH Access Problems**
- Verify enable_ssh setting allows the intended access method
- Check that SSH authorized keys exist and are properly configured
- Ensure SSH service is running on the management VM
- Verify firewall rules allow SSH traffic (port 22)

**SNMP Configuration Issues**
- Verify SNMP mode matches your monitoring requirements
- Check SNMP community string is correctly configured
- For SNMPv3, ensure username and passwords are properly set
- Verify SNMP network management system URI is accessible

**IPv6 Configuration Problems**
- Ensure IPv6 is properly configured on the network
- Verify IPv6 gateway is reachable
- Check that IPv6 addressing scheme is correct
- Test IPv6 connectivity from other network segments

**Certificate Issues**
- Verify TLS certificate URI exists and is accessible
- Check certificate validity and expiration dates
- Ensure certificate includes management VM hostname
- Verify certificate trust chain is properly configured

**NAT Configuration Problems**
- Verify static NAT address is properly configured on firewall/router
- Check that NAT rules include necessary ports (HTTPS, SSH, etc.)
- Ensure external access can reach the NAT address
- Test connectivity from both internal and external networks

**High Availability Issues**
- Ensure only one management VM is marked as primary
- Verify network connectivity between management VMs
- Check that failover mechanisms are properly configured
- Test failover scenarios regularly

**Service Integration Problems**
- Verify all service URIs (DNS, NTP, syslog) are valid and accessible
- Check that services are running and properly configured
- Ensure network connectivity to all referenced services
- Monitor service logs for integration errors

**Import Issues**
- Use the numeric resource ID, not the management VM name
- Verify the management VM exists in the Infinity cluster
- Check provider authentication credentials have access to the resource
- Note that sensitive information may not be available during import

**Initialization Problems**
- Set initializing to false only after VM is fully configured
- Monitor management VM startup logs for errors
- Ensure all dependencies (network, services) are available
- Allow adequate time for initialization to complete

**Update Restrictions**
- Remember that updates are not supported for management VMs
- Plan configuration changes carefully before creation
- Use version control to track configuration changes
- Test configuration in development environment before production deployment