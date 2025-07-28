---
page_title: "pexip_infinity_sip_proxy Resource - terraform-provider-pexip"
subcategory: ""
description: |-
  Manages a Pexip Infinity SIP proxy configuration.
---

# pexip_infinity_sip_proxy (Resource)

Manages a Pexip Infinity SIP proxy configuration. SIP proxies are used to route SIP traffic through intermediate servers for load balancing, security, or network topology reasons. They enable Pexip Infinity to communicate with SIP endpoints that may be behind firewalls or in different network segments.

## Example Usage

### Basic SIP Proxy Configuration

```terraform
resource "pexip_infinity_sip_proxy" "primary_proxy" {
  name      = "Primary SIP Proxy"
  address   = "sip-proxy.example.com"
  transport = "tcp"
}
```

### SIP Proxy with Custom Port

```terraform
resource "pexip_infinity_sip_proxy" "custom_port_proxy" {
  name        = "Custom Port SIP Proxy"
  description = "SIP proxy with custom port configuration"
  address     = "192.168.1.100"
  port        = 5080
  transport   = "udp"
}
```

### TLS-Secured SIP Proxy

```terraform
resource "pexip_infinity_sip_proxy" "secure_proxy" {
  name        = "Secure SIP Proxy"
  description = "TLS-secured SIP proxy for sensitive communications"
  address     = "secure-sip.example.com"
  port        = 5061
  transport   = "tls"
}
```

### Multiple SIP Proxies for Load Balancing

```terraform
variable "sip_proxy_servers" {
  type = list(object({
    name      = string
    address   = string
    port      = number
    transport = string
  }))
  default = [
    {
      name      = "SIP Proxy 1"
      address   = "sip-proxy-1.example.com"
      port      = 5060
      transport = "tcp"
    },
    {
      name      = "SIP Proxy 2"
      address   = "sip-proxy-2.example.com"
      port      = 5060
      transport = "tcp"
    }
  ]
}

resource "pexip_infinity_sip_proxy" "load_balanced_proxies" {
  count       = length(var.sip_proxy_servers)
  name        = var.sip_proxy_servers[count.index].name
  description = "Load balanced SIP proxy ${count.index + 1}"
  address     = var.sip_proxy_servers[count.index].address
  port        = var.sip_proxy_servers[count.index].port
  transport   = var.sip_proxy_servers[count.index].transport
}
```

### Enterprise SIP Infrastructure

```terraform
# Internal SIP proxy for corporate network
resource "pexip_infinity_sip_proxy" "internal_sip_proxy" {
  name        = "Internal SIP Proxy"
  description = "Internal corporate SIP proxy server"
  address     = "10.1.1.50"
  port        = 5060
  transport   = "tcp"
}

# External SIP proxy for internet traffic
resource "pexip_infinity_sip_proxy" "external_sip_proxy" {
  name        = "External SIP Proxy"
  description = "External SIP proxy for internet-facing traffic"
  address     = "sip.company.com"
  port        = 5061
  transport   = "tls"
}

# SIP proxy for legacy systems
resource "pexip_infinity_sip_proxy" "legacy_sip_proxy" {
  name        = "Legacy SIP Proxy"
  description = "SIP proxy for legacy phone systems"
  address     = "legacy-pbx.internal.com"
  port        = 5060
  transport   = "udp"
}
```

## Schema

### Required

- `name` (String) - The name used to refer to this SIP proxy. Maximum length: 250 characters.
- `address` (String) - The address or hostname of the SIP proxy. Maximum length: 255 characters.
- `transport` (String) - The transport protocol for the SIP proxy. Valid values: `tcp`, `udp`, `tls`.

### Optional

- `description` (String) - A description of the SIP proxy. Maximum length: 250 characters.
- `port` (Number) - The port number for the SIP proxy. Range: 1 to 65535. Defaults to standard SIP ports based on transport.

### Read-Only

- `id` (String) - Resource URI for the SIP proxy in Infinity.
- `resource_id` (Number) - The resource integer identifier for the SIP proxy in Infinity.

## Import

Import is supported using the following syntax:

```shell
terraform import pexip_infinity_sip_proxy.example 123
```

Where `123` is the numeric resource ID of the SIP proxy.

## Usage Notes

### Transport Protocol Selection

- **TCP**: Reliable transport, good for most enterprise environments. Provides connection-oriented communication with error checking.
- **UDP**: Faster transport with less overhead, suitable for high-volume environments where some packet loss is acceptable.
- **TLS**: Encrypted transport for secure communications. Recommended for external-facing proxies or sensitive environments.

### Port Configuration

- **Standard Ports**: SIP typically uses port 5060 for TCP/UDP and 5061 for TLS
- **Custom Ports**: Use custom ports when required by network policies or to avoid conflicts
- **Firewall Considerations**: Ensure the configured port is accessible through firewalls and security groups

### Network Topology

- **Internal Proxies**: Use private IP addresses for proxies within your corporate network
- **External Proxies**: Use public IP addresses or FQDNs for internet-facing proxies
- **DMZ Deployment**: Consider placing SIP proxies in a DMZ for security isolation
- **Load Balancing**: Deploy multiple proxies with load balancers for high availability

### SIP Proxy vs Direct Connection

- **Proxy Benefits**: Centralized routing, security enforcement, protocol translation, NAT traversal
- **Direct Connection**: Lower latency, simpler configuration, fewer points of failure
- **Use Cases**: Use proxies when you need centralized control, security, or network segmentation

### Integration with Pexip Infinity

- SIP proxies are used for outbound SIP connections from Pexip Infinity
- Configure multiple proxies for redundancy and load distribution
- Proxies can be assigned to specific routing rules or gateways
- Monitor proxy availability and performance for optimal call routing

## Troubleshooting

### Common Issues

**SIP Proxy Creation Fails**
- Verify the address format is correct (IP address or FQDN)
- Ensure the transport protocol is one of the valid values (tcp, udp, tls)
- Check that the port number is within the valid range (1-65535)
- Verify the name and description don't exceed maximum length limits

**SIP Traffic Not Routing Through Proxy**
- Verify the proxy is configured correctly in Pexip Infinity routing rules
- Check that the proxy server is accessible from Pexip Infinity nodes
- Ensure firewall rules allow traffic on the configured port
- Test connectivity using network tools like telnet or nc

**TLS Connection Issues**
- Verify the proxy supports TLS on the configured port
- Check certificate validity and trust chain
- Ensure TLS cipher suites are compatible between Pexip and the proxy
- Monitor TLS handshake logs for detailed error information

**High Latency or Call Quality Issues**
- Monitor network latency between Pexip Infinity and the SIP proxy
- Check proxy server performance and resource utilization
- Verify Quality of Service (QoS) policies are applied correctly
- Consider deploying proxies closer to Pexip Infinity nodes geographically

**Authentication Failures**
- Verify SIP authentication credentials are configured correctly
- Check that the proxy supports the authentication methods used by Pexip
- Monitor SIP authentication logs on both the proxy and Pexip sides
- Ensure time synchronization between systems for time-based authentication

**NAT and Firewall Traversal Issues**
- Configure appropriate NAT settings on the SIP proxy
- Ensure RTP ports are opened in firewalls for media traffic
- Use STUN/TURN servers if required for NAT traversal
- Consider using Session Border Controllers (SBCs) for complex network topologies

**Import Fails**
- Ensure you're using the numeric resource ID, not the proxy name or address
- Verify the SIP proxy exists in the Infinity cluster
- Check provider authentication credentials have access to the resource
- Confirm the resource ID format matches the expected pattern

**Load Balancing Issues**
- Verify all proxy servers in the pool are accessible and healthy
- Check load balancer configuration for proper health checks
- Monitor traffic distribution across proxy servers
- Ensure session affinity is configured correctly if required

**Legacy System Integration**
- Verify protocol compatibility between legacy systems and the proxy
- Check codec support and transcoding capabilities
- Configure appropriate SIP headers and parameters for legacy compatibility
- Test call flows end-to-end with legacy systems