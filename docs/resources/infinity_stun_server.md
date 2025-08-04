---
page_title: "pexip_infinity_stun_server Resource - terraform-provider-pexip"
subcategory: ""
description: |-
  Manages a Pexip Infinity STUN server configuration.
---

# pexip_infinity_stun_server (Resource)

Manages a STUN server configuration with the Infinity service. STUN (Session Traversal Utilities for NAT) servers help clients discover their public IP address and determine the type of NAT they are behind. This information is essential for establishing direct peer-to-peer media connections and optimizing call quality by avoiding unnecessary media relay.

## Example Usage

### Basic STUN Server

```terraform
resource "pexip_infinity_stun_server" "basic_stun" {
  name    = "Basic STUN Server"
  address = "stun.company.com"
  port    = 3478
}
```

### Corporate STUN Server

```terraform
resource "pexip_infinity_stun_server" "corporate_stun" {
  name        = "Corporate STUN Server"
  description = "Internal STUN server for NAT traversal"
  address     = "stun.company.com"
  port        = 3478
}
```

### Public STUN Services

```terraform
# Google STUN server
resource "pexip_infinity_stun_server" "google_stun" {
  name        = "Google STUN"
  description = "Google public STUN server"
  address     = "stun.l.google.com"
  port        = 19302
}

# Cloudflare STUN server
resource "pexip_infinity_stun_server" "cloudflare_stun" {
  name        = "Cloudflare STUN"
  description = "Cloudflare public STUN server"
  address     = "stun.cloudflare.com"
  port        = 3478
}
```

### Multiple STUN Servers for Redundancy

```terraform
# Primary internal STUN server
resource "pexip_infinity_stun_server" "internal_stun_primary" {
  name        = "Internal STUN Primary"
  description = "Primary internal STUN server"
  address     = "stun1.company.com"
  port        = 3478
}

# Secondary internal STUN server
resource "pexip_infinity_stun_server" "internal_stun_secondary" {
  name        = "Internal STUN Secondary"
  description = "Secondary internal STUN server for redundancy"
  address     = "stun2.company.com"
  port        = 3478
}

# External backup STUN server
resource "pexip_infinity_stun_server" "external_stun_backup" {
  name        = "External STUN Backup"
  description = "External STUN server for backup"
  address     = "stun.stunprotocol.org"
  port        = 3478
}
```

### Regional STUN Servers

```terraform
# STUN servers for different regions
locals {
  regional_stun_servers = {
    "us-east" = {
      address = "stun-us-east.company.com"
      port    = 3478
    }
    "us-west" = {
      address = "stun-us-west.company.com"
      port    = 3478
    }
    "europe" = {
      address = "stun-eu.company.com"
      port    = 3478
    }
    "asia" = {
      address = "stun-asia.company.com"
      port    = 3478
    }
  }
}

resource "pexip_infinity_stun_server" "regional_stun" {
  for_each = local.regional_stun_servers
  
  name        = "STUN Server - ${each.key}"
  description = "Regional STUN server for ${each.key}"
  address     = each.value.address
  port        = each.value.port
}
```

### Load-Balanced STUN Configuration

```terraform
# STUN servers behind load balancer
resource "pexip_infinity_stun_server" "lb_stun" {
  count = var.stun_server_count
  
  name        = "Load Balanced STUN ${count.index + 1}"
  description = "STUN server ${count.index + 1} behind load balancer"
  address     = "stun-lb.company.com"
  port        = 3478 + count.index
}
```

### Mixed Internal and External STUN

```terraform
# Internal STUN for office networks
resource "pexip_infinity_stun_server" "internal_office_stun" {
  name        = "Office Internal STUN"
  description = "STUN server for internal office network"
  address     = "10.0.1.100"
  port        = 3478
}

# External STUN for remote users
resource "pexip_infinity_stun_server" "external_remote_stun" {
  name        = "Remote User STUN"
  description = "External STUN server for remote users"
  address     = "stun.external.company.com"
  port        = 3478
}
```

## Schema

### Required

- `name` (String) - The name used to refer to this STUN server. Maximum length: 250 characters.
- `address` (String) - The address or hostname of the STUN server. Maximum length: 255 characters.
- `port` (Number) - The port number for the STUN server. Range: 1 to 65535.

### Optional

- `description` (String) - A description of the STUN server. Maximum length: 250 characters.

### Read-Only

- `id` (String) - Resource URI for the STUN server in Infinity.
- `resource_id` (Number) - The resource integer identifier for the STUN server in Infinity.

## Import

Import is supported using the following syntax:

```shell
terraform import pexip_infinity_stun_server.example 123
```

Where `123` is the numeric resource ID of the STUN server.

## Usage Notes

### STUN Server Purpose
- Helps clients discover their public IP address and port mappings
- Determines NAT type and behavior for optimal media routing
- Enables direct peer-to-peer connections when possible
- Reduces media latency by avoiding unnecessary relay

### Standard STUN Ports
- **Port 3478**: Standard STUN port (UDP)
- **Port 19302**: Alternative port used by some providers (Google)
- **Custom Ports**: Can be configured for specific deployment needs

### Public vs Private STUN Servers
- **Public STUN**: Free services like Google, Cloudflare, or stunprotocol.org
- **Private STUN**: Internal servers for better control and security
- **Hybrid**: Combination of internal and external for redundancy

### Geographic Distribution
- Deploy STUN servers close to users for better performance
- Use multiple servers in different regions for global deployments
- Consider network latency when choosing STUN server locations
- Implement failover between regional servers

### Security Considerations
- STUN protocol itself doesn't require authentication
- Monitor STUN server usage to detect potential abuse
- Consider firewall rules to restrict STUN server access
- Use internal STUN servers for sensitive environments

### Performance Optimization
- Use multiple STUN servers for load distribution
- Monitor STUN server response times and availability
- Choose servers with low latency to your users
- Implement health checks for STUN server monitoring

## Troubleshooting

### Common Issues

**STUN Server Creation Fails**
- Verify the STUN server name is unique
- Ensure address format is correct (hostname or IP address)
- Check that port number is within valid range (1-65535)
- Verify description doesn't exceed maximum length

**STUN Server Not Responding**
- Test STUN server connectivity using tools like `stunclient` or online STUN testers
- Verify firewall rules allow UDP traffic on the specified port
- Check DNS resolution if using hostname instead of IP address
- Ensure STUN server service is running and properly configured

**NAT Discovery Issues**
- Verify STUN requests and responses are properly formatted
- Check if STUN server supports the required NAT discovery methods
- Ensure network path allows bidirectional UDP traffic
- Test from different network locations and NAT types

**Performance Problems**
- Monitor STUN server response times and availability
- Check for network congestion affecting STUN traffic
- Verify STUN server has adequate resources for concurrent requests
- Consider geographic distribution of STUN servers

**Firewall and Network Issues**
- Ensure UDP port is open for STUN traffic
- Check that symmetric NAT isn't blocking STUN responses
- Verify routing allows traffic to reach STUN server
- Test connectivity from various network segments

**DNS Resolution Problems**
- Verify STUN server hostname resolves correctly
- Check DNS server configuration and reachability
- Ensure DNS TTL values are appropriate
- Test DNS resolution from client networks

**Import Fails**
- Ensure you're using the numeric resource ID, not the name
- Verify the STUN server exists in the Infinity cluster
- Check provider authentication credentials have access to the resource

**Load Balancing Issues**
- Verify load balancer health checks for STUN servers
- Ensure proper session affinity configuration if needed
- Check that all STUN servers behind load balancer are functional
- Test failover behavior between STUN servers

**Public STUN Service Issues**
- Monitor availability of public STUN services
- Have backup STUN servers configured for redundancy
- Check for rate limiting or usage restrictions
- Verify public service terms of use and limitations

**Regional Connectivity Problems**
- Test STUN server access from different geographic locations
- Monitor cross-region network latency and reliability
- Ensure regional STUN servers are properly configured
- Verify routing policies don't block regional traffic

**Client-Side STUN Problems**
- Verify client applications support STUN protocol correctly
- Check client firewall and security software settings
- Ensure client network allows outbound UDP traffic
- Test with different client types and versions

**Monitoring and Diagnostics**
- Implement STUN server health monitoring
- Set up alerts for STUN server unavailability
- Monitor STUN request/response patterns
- Log STUN server usage for capacity planning

**Network Architecture Considerations**
- Ensure STUN servers are accessible from all client networks
- Consider placement relative to NAT devices and firewalls
- Plan for STUN server scalability and redundancy
- Document STUN server network requirements and dependencies