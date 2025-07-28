---
page_title: "pexip_infinity_ntp_server Resource - terraform-provider-pexip"
subcategory: ""
description: |-
  Manages a Pexip Infinity NTP server configuration.
---

# pexip_infinity_ntp_server (Resource)

Manages a Pexip Infinity NTP server configuration. NTP (Network Time Protocol) servers provide accurate time synchronization for Pexip Infinity nodes, which is critical for proper operation of security certificates, logging, and other time-dependent features.

## Example Usage

### Basic Usage

```terraform
resource "pexip_infinity_ntp_server" "primary_ntp" {
  address = "pool.ntp.org"
}
```

### NTP Server with Description

```terraform
resource "pexip_infinity_ntp_server" "primary_ntp" {
  address     = "pool.ntp.org"
  description = "Primary NTP pool server"
}

resource "pexip_infinity_ntp_server" "secondary_ntp" {
  address     = "time.cloudflare.com"
  description = "Cloudflare NTP server"
}
```

### Corporate NTP Configuration

```terraform
resource "pexip_infinity_ntp_server" "internal_ntp_primary" {
  address     = "ntp1.company.com"
  description = "Primary internal NTP server"
}

resource "pexip_infinity_ntp_server" "internal_ntp_secondary" {
  address     = "ntp2.company.com"
  description = "Secondary internal NTP server"
}

# IP address based NTP server
resource "pexip_infinity_ntp_server" "ip_based_ntp" {
  address     = "10.0.1.100"
  description = "Internal NTP server by IP"
}
```

### Multiple NTP Servers for High Availability

```terraform
# Public NTP servers
resource "pexip_infinity_ntp_server" "public_ntp" {
  count       = length(var.ntp_servers)
  address     = var.ntp_servers[count.index]
  description = "Public NTP server ${count.index + 1}"
}

# Regional NTP servers
locals {
  regional_ntp_servers = {
    "us" = ["time1.google.com", "time2.google.com"]
    "eu" = ["ntp1.ntp.se", "ntp2.ntp.se"]
    "asia" = ["ntp.nict.jp", "ntp1.jst.mfeed.ad.jp"]
  }
}

resource "pexip_infinity_ntp_server" "regional_ntp" {
  for_each    = toset(local.regional_ntp_servers[var.region])
  address     = each.value
  description = "Regional NTP server for ${var.region} - ${each.value}"
}
```

### Stratum-based NTP Configuration

```terraform
# Stratum 1 NTP servers (high accuracy)
resource "pexip_infinity_ntp_server" "stratum1_ntp" {
  for_each = toset([
    "time.nist.gov",
    "time-a.nist.gov",
    "time-b.nist.gov"
  ])
  address     = each.value
  description = "NIST Stratum 1 NTP server - ${each.value}"
}
```

## Schema

### Required

- `address` (String) - The IP address of the NTP server. Maximum length: 255 characters.

### Optional

- `description` (String) - A description of the NTP server. Maximum length: 250 characters.

### Read-Only

- `id` (String) - Resource URI for the NTP server in Infinity.
- `resource_id` (Number) - The resource integer identifier for the NTP server in Infinity.

## Import

Import is supported using the following syntax:

```shell
terraform import pexip_infinity_ntp_server.example 123
```

Where `123` is the numeric resource ID of the NTP server.

## Usage Notes

### NTP Server Requirements
- NTP servers must be accessible from all system locations where they will be used
- Use reliable NTP sources with good stratum levels (Stratum 1 or 2 preferred)
- Ensure NTP servers are geographically diverse for redundancy
- Both hostnames and IP addresses are supported in the address field

### Time Synchronization Importance
- Accurate time is critical for TLS certificate validation
- Log correlation requires synchronized timestamps
- Authentication protocols depend on accurate time
- Media synchronization relies on precise timing

### Public vs Private NTP
- **Public NTP servers**: pool.ntp.org, time.google.com, time.cloudflare.com
- **Private NTP servers**: Corporate time servers, GPS-synchronized servers
- **Government time sources**: NIST, USNO for high-accuracy requirements
- Consider your network security policies when choosing between public and private sources

### System Location Assignment
- NTP servers are assigned to system locations, not directly to nodes
- All nodes deployed in a system location inherit the NTP configuration
- Multiple NTP servers provide redundancy and improved accuracy
- Nodes will synchronize with the most accurate available source

### Best Practices
- Use at least 3 NTP servers for proper synchronization algorithm operation
- Choose NTP servers with low network latency
- Monitor NTP synchronization status on all nodes
- Implement NTP server monitoring and alerting
- Use authenticated NTP where security is critical

## Troubleshooting

### Common Issues

**NTP Server Creation Fails**
- Verify the address format is correct (hostname or IP address)
- Ensure the NTP server address is reachable from the Infinity Manager
- Check that the description doesn't exceed the maximum length (250 characters)
- Verify the address doesn't exceed the maximum length (255 characters)

**Time Synchronization Not Working**
- Test NTP server connectivity from conferencing nodes using `ntpdate -q <server>`
- Verify firewall rules allow NTP traffic (UDP port 123)
- Check if the NTP server is responding to queries
- Ensure network routing allows access to NTP servers

**Clock Drift Issues**
- Monitor system clock accuracy on all nodes
- Check NTP server stratum levels and accuracy
- Verify network latency to NTP servers is consistent
- Consider using GPS-synchronized local NTP servers for high accuracy

**Certificate Validation Failures**
- Ensure system time is within certificate validity period
- Check that NTP synchronization is working correctly
- Verify time zone settings are correct
- Monitor for sudden time changes that could break TLS sessions

**Intermittent Time Sync Failures**
- Check NTP server availability and uptime
- Verify network connectivity is stable to NTP servers
- Monitor for NTP server maintenance windows
- Consider adding additional redundant NTP servers

**Import Fails**
- Ensure you're using the numeric resource ID, not the address
- Verify the NTP server exists in the Infinity cluster
- Check provider authentication credentials have access to the resource

**NTP Server Not Accessible**
- Verify the address is correct and reachable
- Check firewall rules and network routing to UDP port 123
- Test NTP queries manually using tools like `ntpdate` or `chrony`
- Ensure the NTP server service is running and configured correctly

**Corporate NTP Server Issues**
- Verify internal NTP servers are properly configured and synchronized
- Check that internal NTP servers have external time sources
- Ensure proper NTP hierarchy with stratum levels
- Monitor NTP server health and synchronization status

**Geographic Distribution Problems**
- Ensure NTP servers are distributed across different networks
- Avoid single points of failure in NTP infrastructure
- Consider regional time sources for global deployments
- Monitor network latency to all configured NTP servers