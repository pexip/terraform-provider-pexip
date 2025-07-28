---
page_title: "pexip_infinity_dns_server Resource - terraform-provider-pexip"
subcategory: ""
description: |-
  Manages a Pexip Infinity DNS server configuration.
---

# pexip_infinity_dns_server (Resource)

Manages a Pexip Infinity DNS server configuration. DNS servers are used by Pexip Infinity nodes for domain name resolution and are assigned to system locations to provide name resolution services for conferencing nodes deployed in those locations.

## Example Usage

### Basic Usage

```terraform
resource "pexip_infinity_dns_server" "primary_dns" {
  address = "8.8.8.8"
}
```

### DNS Server with Description

```terraform
resource "pexip_infinity_dns_server" "primary_dns" {
  address     = "8.8.8.8"
  description = "Primary Google DNS server"
}

resource "pexip_infinity_dns_server" "secondary_dns" {
  address     = "8.8.4.4"
  description = "Secondary Google DNS server"
}
```

### Corporate DNS Servers

```terraform
resource "pexip_infinity_dns_server" "internal_dns_primary" {
  address     = "10.0.1.10"
  description = "Primary internal DNS server"
}

resource "pexip_infinity_dns_server" "internal_dns_secondary" {
  address     = "10.0.1.11"
  description = "Secondary internal DNS server"
}

# IPv6 DNS server
resource "pexip_infinity_dns_server" "ipv6_dns" {
  address     = "2001:4860:4860::8888"
  description = "Google IPv6 DNS server"
}
```

### Multiple DNS Servers for Different Purposes

```terraform
# External DNS for internet resolution
resource "pexip_infinity_dns_server" "external_dns" {
  count       = length(var.external_dns_servers)
  address     = var.external_dns_servers[count.index]
  description = "External DNS server ${count.index + 1}"
}

# Internal DNS for corporate domains
resource "pexip_infinity_dns_server" "internal_dns" {
  count       = length(var.internal_dns_servers)
  address     = var.internal_dns_servers[count.index]
  description = "Internal DNS server ${count.index + 1}"
}
```

## Schema

### Required

- `address` (String) - The IP address of the DNS server.

### Optional

- `description` (String) - A description of the DNS server. Maximum length: 250 characters.

### Read-Only

- `id` (String) - Resource URI for the DNS server in Infinity.
- `resource_id` (Number) - The resource integer identifier for the DNS server in Infinity.

## Import

Import is supported using the following syntax:

```shell
terraform import pexip_infinity_dns_server.example 123
```

Where `123` is the numeric resource ID of the DNS server.

## Usage Notes

### DNS Server Requirements
- DNS servers must be accessible from all system locations where they will be used
- Both IPv4 and IPv6 addresses are supported
- DNS servers should be reliable and have high availability
- Consider using multiple DNS servers for redundancy

### Corporate vs Public DNS
- Use corporate DNS servers for internal domain resolution
- Public DNS servers (like Google's 8.8.8.8 or Cloudflare's 1.1.1.1) for external resolution
- Configure appropriate DNS servers based on your network security policies
- Consider DNS filtering and content blocking requirements

### System Location Assignment
- DNS servers are assigned to system locations, not directly to nodes
- All nodes deployed in a system location inherit the DNS configuration
- DNS servers are tried in the order they appear in the system location configuration
- Ensure DNS servers are reachable from all system locations where they're used

### Performance Considerations
- Use geographically close DNS servers for better performance
- Monitor DNS query response times and availability
- Consider using anycast DNS services for global deployments
- Implement DNS caching where appropriate

### Security Considerations
- Use DNS over HTTPS (DoH) or DNS over TLS (DoT) when supported
- Implement DNS filtering to block malicious domains
- Monitor DNS queries for security analysis
- Consider using private DNS zones for internal services

## Troubleshooting

### Common Issues

**DNS Server Creation Fails**
- Verify the IP address format is correct (IPv4 or IPv6)
- Ensure the DNS server address is reachable from the Infinity Manager
- Check that the description doesn't exceed the maximum length

**DNS Resolution Not Working**
- Test DNS server connectivity from conferencing nodes
- Verify firewall rules allow DNS traffic (UDP/TCP port 53)
- Check if the DNS server is responding to queries
- Ensure DNS server has proper zone configurations

**Slow DNS Resolution**
- Monitor DNS query response times
- Check network latency to DNS servers
- Consider using closer DNS servers geographically
- Verify DNS server is not overloaded

**Intermittent DNS Failures**
- Check DNS server availability and uptime
- Verify network connectivity is stable
- Monitor for DNS server maintenance windows
- Consider adding additional redundant DNS servers

**Import Fails**
- Ensure you're using the numeric resource ID, not the IP address
- Verify the DNS server exists in the Infinity cluster
- Check provider authentication credentials have access to the resource

**DNS Server Not Accessible**
- Verify the IP address is correct and reachable
- Check firewall rules and network routing
- Test DNS queries manually using tools like nslookup or dig
- Ensure the DNS server service is running and configured correctly

**Corporate DNS Issues**
- Verify internal DNS servers can resolve both internal and external domains
- Check DNS forwarders are configured correctly for external resolution
- Ensure proper zone delegation for subdomains
- Verify reverse DNS (PTR) records are configured if required