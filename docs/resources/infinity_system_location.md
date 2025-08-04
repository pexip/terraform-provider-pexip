---
page_title: "pexip_infinity_system_location Resource - terraform-provider-pexip"
subcategory: ""
description: |-
  Manages a Pexip Infinity system location configuration.
---

# pexip_infinity_system_location (Resource)

Manages a Pexip Infinity system location configuration. A system location defines the infrastructure settings for a specific physical or logical location in your Pexip Infinity deployment, including DNS servers, NTP servers, syslog servers, and MTU settings.

## Example Usage

### Basic Usage

```terraform
resource "pexip_infinity_system_location" "main_location" {
  name        = "Main Office"
  description = "Primary office location"
}
```

### Full Configuration with External Services

```terraform
# Create DNS servers
resource "pexip_infinity_dns_server" "primary_dns" {
  address     = "8.8.8.8"
  description = "Primary DNS server"
}

resource "pexip_infinity_dns_server" "secondary_dns" {
  address     = "8.8.4.4"
  description = "Secondary DNS server"
}

# Create NTP servers
resource "pexip_infinity_ntp_server" "primary_ntp" {
  address     = "pool.ntp.org"
  description = "Primary NTP server"
}

# Create syslog server
resource "pexip_infinity_syslog_server" "central_logging" {
  address     = "syslog.company.com"
  port        = 514
  transport   = "udp"
  description = "Central syslog server"
  audit_log   = true
  support_log = true
  web_log     = true
}

# System location with all services
resource "pexip_infinity_system_location" "main_location" {
  name            = "Main Office"
  description     = "Primary office location with full logging and time sync"
  mtu             = 1500
  dns_servers     = [
    pexip_infinity_dns_server.primary_dns.id,
    pexip_infinity_dns_server.secondary_dns.id
  ]
  ntp_servers     = [
    pexip_infinity_ntp_server.primary_ntp.id
  ]
  syslog_servers  = [
    pexip_infinity_syslog_server.central_logging.id
  ]
}
```

### Multiple Locations for Different Sites

```terraform
# Branch office with minimal configuration
resource "pexip_infinity_system_location" "branch_office" {
  name        = "Branch Office"
  description = "Remote branch office location"
  mtu         = 1400  # Lower MTU for WAN links
  dns_servers = [
    pexip_infinity_dns_server.primary_dns.id
  ]
  ntp_servers = [
    pexip_infinity_ntp_server.primary_ntp.id
  ]
}

# Data center location with dedicated services
resource "pexip_infinity_system_location" "datacenter" {
  name        = "Data Center"
  description = "Primary data center location"
  mtu         = 9000  # Jumbo frames for internal networks
  dns_servers = [
    pexip_infinity_dns_server.primary_dns.id,
    pexip_infinity_dns_server.secondary_dns.id
  ]
  ntp_servers = [
    pexip_infinity_ntp_server.primary_ntp.id
  ]
  syslog_servers = [
    pexip_infinity_syslog_server.central_logging.id
  ]
}
```

## Schema

### Required

- `name` (String) - The name used to refer to this system location. Maximum length: 250 characters.

### Optional

- `description` (String) - A description of the system location. Maximum length: 250 characters.
- `dns_servers` (List of String) - List of DNS server resource URIs for this system location.
- `ntp_servers` (List of String) - List of NTP server resource URIs for this system location.
- `mtu` (Number) - Maximum Transmission Unit - the size of the largest packet that can be transmitted via the network interface for this system location. It depends on your network topology as to whether you may need to specify an MTU value here. Range: 512 to 1500.
- `syslog_servers` (List of String) - The Syslog servers to be used by Conferencing Nodes deployed in this Location.

### Read-Only

- `id` (String) - Resource URI for the system location in Infinity.
- `resource_id` (Number) - The resource integer identifier for the system location in Infinity.

## Import

Import is supported using the following syntax:

```shell
terraform import pexip_infinity_system_location.example 123
```

Where `123` is the numeric resource ID of the system location.

## Usage Notes

### System Location Purpose
- System locations define infrastructure settings for specific physical or logical locations
- Each location can have its own DNS, NTP, and syslog server configurations
- Conferencing nodes are deployed within system locations and inherit these settings

### DNS Server Configuration
- Specify DNS servers using their resource URIs (e.g., `/api/admin/configuration/v1/dns_server/1/`)
- Multiple DNS servers provide redundancy for name resolution
- Order in the list determines priority for DNS queries

### NTP Server Configuration
- NTP servers ensure accurate time synchronization across all nodes in the location
- Critical for proper operation of security certificates and logging
- Use reliable NTP sources like pool.ntp.org or internal NTP servers

### Syslog Server Configuration
- Centralized logging helps with monitoring and troubleshooting
- Different log types (audit, support, web) can be enabled per server
- Configure appropriate transport protocols (UDP, TCP, TLS) based on security requirements

### MTU Settings
- Default MTU is typically 1500 bytes for Ethernet networks
- Use jumbo frames (up to 9000 bytes) in data center environments for better performance
- Lower MTU may be needed for WAN connections or specific network configurations
- Ensure MTU settings match your network infrastructure

### Resource Dependencies
- DNS, NTP, and syslog servers must be created before referencing them in system locations
- Use Terraform dependencies to ensure proper creation order
- Changes to referenced servers automatically propagate to system locations

## Troubleshooting

### Common Issues

**System Location Creation Fails**
- Verify the location name is unique within the Infinity cluster
- Ensure referenced DNS, NTP, and syslog servers exist and are accessible
- Check that MTU value is within the valid range (512-1500)

**DNS Resolution Issues**
- Verify DNS server addresses are correct and reachable from the location
- Test DNS resolution from nodes deployed in the location
- Check firewall rules allow DNS traffic (UDP/TCP port 53)

**Time Synchronization Problems**
- Ensure NTP servers are reachable from the location
- Verify NTP server addresses or hostnames are correct
- Check that nodes can reach NTP servers on UDP port 123

**Logging Not Working**
- Verify syslog server configuration and reachability
- Check syslog server capacity and storage
- Ensure appropriate log types are enabled on the syslog server resource

**Network Performance Issues**
- Review MTU settings and ensure they match network infrastructure
- Test connectivity between nodes and external services
- Monitor for packet fragmentation if MTU is set too high

**Import Fails**
- Ensure you're using the numeric resource ID, not the name
- Verify the system location exists in the Infinity cluster
- Check provider authentication credentials have access to the resource

**Referenced Servers Not Found**
- Ensure DNS, NTP, and syslog server resources exist before referencing them
- Use proper resource URI format for server references
- Verify resource IDs are correct and servers are not deleted