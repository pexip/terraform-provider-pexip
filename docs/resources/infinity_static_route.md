---
page_title: "pexip_infinity_static_route Resource - terraform-provider-pexip"
subcategory: ""
description: |-
  Manages a Pexip Infinity static route configuration.
---

# pexip_infinity_static_route (Resource)

Manages a Pexip Infinity static route configuration. Static routes define custom network paths for traffic routing when the default routing table is insufficient. They are commonly used to direct traffic to specific network segments, route traffic through firewalls, or establish connectivity to remote networks that are not reachable through the default gateway.

## Example Usage

### Basic Static Route

```terraform
resource "pexip_infinity_static_route" "basic_route" {
  name    = "Internal Network Route"
  address = "10.1.0.0"
  prefix  = 16
  gateway = "192.168.1.1"
}
```

### Static Route to Remote Office

```terraform
resource "pexip_infinity_static_route" "remote_office" {
  name    = "Remote Office Network"
  address = "172.16.0.0"
  prefix  = 12
  gateway = "10.0.1.254"
}
```

### Multiple Static Routes for Different Networks

```terraform
# Route to internal corporate network
resource "pexip_infinity_static_route" "corporate_network" {
  name    = "Corporate Network"
  address = "10.0.0.0"
  prefix  = 8
  gateway = "192.168.100.1"
}

# Route to DMZ network
resource "pexip_infinity_static_route" "dmz_network" {
  name    = "DMZ Network"
  address = "172.20.0.0"
  prefix  = 16
  gateway = "192.168.100.2"
}

# Route to guest network
resource "pexip_infinity_static_route" "guest_network" {
  name    = "Guest Network"
  address = "192.168.50.0"
  prefix  = 24
  gateway = "192.168.100.3"
}
```

### Host-specific Static Route

```terraform
resource "pexip_infinity_static_route" "specific_server" {
  name    = "Database Server Route"
  address = "10.1.1.100"
  prefix  = 32
  gateway = "192.168.1.10"
}
```

### Default Route Override

```terraform
resource "pexip_infinity_static_route" "default_route" {
  name    = "Custom Default Route"
  address = "0.0.0.0"
  prefix  = 0
  gateway = "192.168.1.254"
}
```

### Enterprise Network Routing

```terraform
variable "static_routes" {
  type = list(object({
    name    = string
    address = string
    prefix  = number
    gateway = string
  }))
  default = [
    {
      name    = "Branch Office 1"
      address = "10.10.0.0"
      prefix  = 16
      gateway = "192.168.1.10"
    },
    {
      name    = "Branch Office 2"
      address = "10.20.0.0"
      prefix  = 16
      gateway = "192.168.1.11"
    },
    {
      name    = "Data Center"
      address = "172.16.0.0"
      prefix  = 12
      gateway = "192.168.1.20"
    }
  ]
}

resource "pexip_infinity_static_route" "enterprise_routes" {
  count   = length(var.static_routes)
  name    = var.static_routes[count.index].name
  address = var.static_routes[count.index].address
  prefix  = var.static_routes[count.index].prefix
  gateway = var.static_routes[count.index].gateway
}
```

### Cloud Network Integration

```terraform
# Route to AWS VPC
resource "pexip_infinity_static_route" "aws_vpc" {
  name    = "AWS VPC Network"
  address = "10.100.0.0"
  prefix  = 16
  gateway = "192.168.1.50"
}

# Route to Azure VNet
resource "pexip_infinity_static_route" "azure_vnet" {
  name    = "Azure VNet Network"
  address = "10.200.0.0"
  prefix  = 16
  gateway = "192.168.1.51"
}

# Route to Google Cloud VPC
resource "pexip_infinity_static_route" "gcp_vpc" {
  name    = "Google Cloud VPC"
  address = "10.300.0.0"
  prefix  = 16
  gateway = "192.168.1.52"
}
```

### VPN Gateway Routes

```terraform
# Route through Site-to-Site VPN
resource "pexip_infinity_static_route" "vpn_site_to_site" {
  name    = "Site-to-Site VPN"
  address = "172.30.0.0"
  prefix  = 16
  gateway = "192.168.1.100"
}

# Route through Client VPN concentrator
resource "pexip_infinity_static_route" "vpn_client" {
  name    = "Client VPN Network"
  address = "172.31.0.0"
  prefix  = 24
  gateway = "192.168.1.101"
}
```

### Conditional Static Routes

```terraform
# Only create routes if gateway is available
locals {
  create_routes = var.gateway_available
}

resource "pexip_infinity_static_route" "conditional_route" {
  count   = local.create_routes ? 1 : 0
  name    = "Conditional Network Route"
  address = "10.50.0.0"
  prefix  = 16
  gateway = var.conditional_gateway
}
```

## Schema

### Required

- `name` (String) - The name of the static route. Maximum length: 250 characters.
- `address` (String) - The destination network address for the static route. Must be a valid IPv4 address.
- `prefix` (Number) - The network prefix length (CIDR notation). Valid range: 0-32.
- `gateway` (String) - The gateway IP address for the static route. Must be a valid IPv4 address.

### Read-Only

- `id` (String) - Resource URI for the static route in Infinity.
- `resource_id` (Number) - The resource integer identifier for the static route in Infinity.

## Import

Import is supported using the following syntax:

```shell
terraform import pexip_infinity_static_route.example 123
```

Where `123` is the numeric resource ID of the static route.

## Usage Notes

### CIDR Notation

Static routes use CIDR (Classless Inter-Domain Routing) notation to define network segments:

- **Single Host**: Use `/32` prefix (e.g., `192.168.1.100/32`)
- **Subnet**: Use appropriate prefix length (e.g., `192.168.1.0/24` for 256 addresses)
- **Large Networks**: Use smaller prefix values (e.g., `10.0.0.0/8` for Class A networks)
- **Default Route**: Use `0.0.0.0/0` to match all traffic

### Route Priority

- **Most Specific**: Routes with longer prefixes (more specific) take precedence
- **Default Route**: Routes with `/0` prefix are used when no other routes match
- **Administrative Distance**: Static routes typically have lower administrative distance than dynamic routes
- **Route Conflicts**: Avoid overlapping routes that could cause routing loops

### Gateway Requirements

- **Reachability**: Gateway must be reachable from the Pexip Infinity nodes
- **Same Subnet**: Gateway should typically be on the same subnet as the Pexip nodes
- **Next Hop**: Gateway acts as the next hop for traffic destined to the specified network
- **High Availability**: Consider multiple gateways for redundancy

### Network Segmentation

- **Internal Networks**: Route internal corporate traffic through appropriate gateways
- **DMZ Access**: Direct DMZ traffic through security appliances
- **Internet Traffic**: Route internet-bound traffic through firewalls or proxy servers
- **Cloud Integration**: Route cloud network traffic through VPN gateways or direct connections

### Routing Best Practices

- **Hierarchical Design**: Use hierarchical addressing and routing design
- **Summarization**: Summarize routes where possible to reduce routing table size
- **Loop Prevention**: Avoid creating routing loops through careful route design
- **Documentation**: Maintain clear documentation of routing decisions and network topology

### Performance Considerations

- **Route Table Size**: Large routing tables can impact performance
- **Gateway Performance**: Ensure gateways can handle the expected traffic load
- **Latency**: Consider network latency when choosing gateway paths
- **Bandwidth**: Ensure adequate bandwidth on all routing paths

### Security Implications

- **Traffic Inspection**: Route sensitive traffic through security appliances
- **Access Control**: Use routing to enforce network access policies
- **Audit Trail**: Monitor and log routing changes for security auditing
- **Isolation**: Use routing to isolate different network segments

## Troubleshooting

### Common Issues

**Static Route Creation Fails**
- Verify the address is a valid IPv4 address format
- Ensure the prefix is within the valid range (0-32)
- Check that the gateway is a valid IPv4 address
- Verify the name doesn't exceed the maximum length

**Traffic Not Routing Through Static Route**
- Verify the destination address falls within the specified network range
- Check that the gateway is reachable from Pexip Infinity nodes
- Ensure no more specific routes are overriding the static route
- Verify the gateway is forwarding traffic correctly

**Gateway Unreachable**
- Test connectivity to the gateway from Pexip nodes using ping
- Check network configuration and firewall rules
- Verify the gateway is on the same subnet as Pexip nodes
- Ensure the gateway device is operational and configured correctly

**Routing Loops**
- Check for circular routing dependencies
- Verify static routes don't conflict with dynamic routing protocols
- Ensure proper route summarization and hierarchy
- Monitor routing tables for inconsistencies

**Network Performance Issues**
- Monitor network latency through the static route path
- Check bandwidth utilization on gateway links
- Verify Quality of Service (QoS) policies are applied correctly
- Consider optimizing route paths for better performance

**Route Conflicts**
- Check for overlapping network ranges in multiple routes
- Verify route priorities and administrative distances
- Ensure static routes don't conflict with DHCP-assigned routes
- Monitor routing table for duplicate or conflicting entries

**Security Access Issues**
- Verify firewall rules allow traffic through the gateway
- Check access control lists (ACLs) on routing devices
- Ensure proper NAT configuration if required
- Monitor security logs for blocked traffic

**High Availability Issues**
- Implement redundant gateways for critical routes
- Configure route failover mechanisms where supported
- Monitor gateway availability and automatic failover
- Test failover scenarios regularly

**Import Fails**
- Ensure you're using the numeric resource ID, not the route name
- Verify the static route exists in the Infinity cluster
- Check provider authentication credentials have access to the resource
- Confirm the route configuration is accessible

**Network Convergence Problems**
- Allow time for routing changes to propagate through the network
- Monitor dynamic routing protocol updates if applicable
- Check for routing protocol interactions with static routes
- Verify network topology changes are reflected in routing

**DNS Resolution Issues**
- Ensure DNS traffic can reach DNS servers through static routes
- Verify reverse DNS resolution for gateway addresses
- Check that DNS servers are reachable via the configured routes
- Monitor DNS query response times through static route paths

**Cloud Integration Problems**
- Verify VPN tunnel status for cloud network routes
- Check cloud network configuration and routing tables
- Ensure proper authentication for cloud network access
- Monitor cloud network connectivity and performance