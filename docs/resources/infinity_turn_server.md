---
page_title: "pexip_infinity_turn_server Resource - terraform-provider-pexip"
subcategory: ""
description: |-
  Manages a Pexip Infinity TURN server configuration.
---

# pexip_infinity_turn_server (Resource)

Manages a TURN server configuration with the Infinity service. TURN (Traversal Using Relays around NAT) servers help clients behind NAT and firewalls establish media connections by providing relay services. This is essential for participants connecting from restrictive network environments.

## Example Usage

### Basic TURN Server

```terraform
resource "pexip_infinity_turn_server" "basic_turn" {
  name    = "Basic TURN Server"
  address = "turn.company.com"
}
```

### TURN Server with Authentication

```terraform
resource "pexip_infinity_turn_server" "authenticated_turn" {
  name        = "Corporate TURN Server"
  description = "TURN server with username/password authentication"
  address     = "turn.company.com"
  port        = 3478
  server_type = "namepsw"
  transport_type = "udp"
  username    = "turnuser"
  password    = var.turn_password
}
```

### Secure TURN Server with TLS

```terraform
resource "pexip_infinity_turn_server" "secure_turn" {
  name           = "Secure TURN Server"
  description    = "TLS-encrypted TURN server"
  address        = "turns.company.com"
  port           = 5349
  server_type    = "namepsw"
  transport_type = "tls"
  username       = "secure-turn-user"
  password       = var.secure_turn_password
}
```

### Shared Secret TURN Server

```terraform
resource "pexip_infinity_turn_server" "shared_secret_turn" {
  name           = "Shared Secret TURN"
  description    = "TURN server using shared secret authentication"
  address        = "turn-shared.company.com"
  port           = 3478
  server_type    = "coturn_shared"
  transport_type = "udp"
  secret_key     = var.turn_shared_secret
}
```

### Multiple TURN Servers for Load Distribution

```terraform
# Primary TURN server
resource "pexip_infinity_turn_server" "turn_primary" {
  name           = "Primary TURN Server"
  description    = "Primary TURN server for main office"
  address        = "turn1.company.com"
  port           = 3478
  server_type    = "namepsw"
  transport_type = "udp"
  username       = "primary-turn"
  password       = var.primary_turn_password
}

# Secondary TURN server
resource "pexip_infinity_turn_server" "turn_secondary" {
  name           = "Secondary TURN Server"
  description    = "Secondary TURN server for redundancy"
  address        = "turn2.company.com"
  port           = 3478
  server_type    = "namepsw"
  transport_type = "tcp"
  username       = "secondary-turn"
  password       = var.secondary_turn_password
}
```

### Regional TURN Servers

```terraform
# TURN servers for different regions
resource "pexip_infinity_turn_server" "regional_turn" {
  for_each = var.regional_turn_servers
  
  name           = "TURN Server - ${each.key}"
  description    = "Regional TURN server for ${each.key}"
  address        = each.value.address
  port           = each.value.port
  server_type    = "namepsw"
  transport_type = each.value.transport
  username       = each.value.username
  password       = each.value.password
}
```

### Public TURN Service Integration

```terraform
resource "pexip_infinity_turn_server" "public_turn" {
  name           = "Public TURN Service"
  description    = "External TURN service provider"
  address        = "global-turn.example.com"
  port           = 3478
  server_type    = "namepsw"
  transport_type = "udp"
  username       = var.public_turn_username
  password       = var.public_turn_password
}
```

## Schema

### Required

- `name` (String) - The name used to refer to this TURN server. Maximum length: 250 characters.
- `address` (String) - The address or hostname of the TURN server. Maximum length: 255 characters.

### Optional

- `description` (String) - A description of the TURN server. Maximum length: 250 characters.
- `port` (Number) - The port number for the TURN server. Range: 1 to 65535.
- `server_type` (String) - The type of TURN server. Valid values: namepsw, coturn_shared. Defaults to `"namepsw"`.
- `transport_type` (String) - The transport type for the TURN server. Valid values: udp, tcp, tls. Defaults to `"udp"`.
- `username` (String) - Username for authentication to the TURN server. Maximum length: 100 characters.
- `password` (String, Sensitive) - Password for authentication to the TURN server. Maximum length: 100 characters.
- `secret_key` (String, Sensitive) - Secret key for shared secret TURN servers. Maximum length: 256 characters.

### Read-Only

- `id` (String) - Resource URI for the TURN server in Infinity.
- `resource_id` (Number) - The resource integer identifier for the TURN server in Infinity.

## Import

Import is supported using the following syntax:

```shell
terraform import pexip_infinity_turn_server.example 123
```

Where `123` is the numeric resource ID of the TURN server.

## Usage Notes

### TURN Server Types
- **namepsw**: Standard TURN server with username/password authentication
- **coturn_shared**: TURN server using shared secret authentication (e.g., coturn)

### Transport Types
- **UDP (3478)**: Default and most common, lowest overhead
- **TCP (3478)**: More reliable for restricted networks, higher overhead
- **TLS (5349)**: Encrypted transport for secure environments

### Authentication Methods
- **Username/Password**: Traditional static credentials for namepsw servers
- **Shared Secret**: Dynamic credential generation for coturn_shared servers
- **Credentials**: Should be regularly rotated for security

### Network Considerations
- TURN servers should be accessible from client networks
- Consider geographic distribution for global deployments
- Ensure adequate bandwidth for media relay
- Monitor TURN server load and capacity

### Security Best Practices
- Use TLS transport in security-sensitive environments
- Regularly rotate TURN server credentials
- Monitor TURN server usage and access patterns
- Implement proper firewall rules for TURN traffic

### Performance Optimization
- Deploy TURN servers close to users geographically
- Use UDP for better performance when possible
- Monitor TURN server resource utilization
- Implement load balancing for high-traffic scenarios

## Troubleshooting

### Common Issues

**TURN Server Creation Fails**
- Verify the TURN server name is unique
- Ensure address format is correct (hostname or IP)
- Check that port number is within valid range (1-65535)
- Verify server_type is one of: namepsw, coturn_shared

**TURN Server Not Accessible**
- Test TURN server connectivity from client networks
- Verify firewall rules allow TURN traffic on specified port
- Check DNS resolution for TURN server hostname
- Ensure TURN server service is running and configured

**Authentication Failures**
- Verify username and password are correct for namepsw servers
- Check shared secret configuration for coturn_shared servers
- Ensure credentials haven't expired
- Test authentication independently using TURN testing tools

**Media Relay Issues**
- Monitor TURN server bandwidth and capacity
- Check for port allocation limits on TURN server
- Verify media firewall rules allow TURN-relayed traffic
- Ensure TURN server has adequate resources for concurrent sessions

**Transport Protocol Problems**
- **UDP**: Check for packet loss or firewall blocking
- **TCP**: Verify TCP connections can be established
- **TLS**: Ensure valid certificates and TLS configuration

**Performance Problems**
- Monitor TURN server CPU and memory usage
- Check network bandwidth utilization
- Verify TURN server isn't overloaded with concurrent sessions
- Consider adding additional TURN servers for load distribution

**Certificate Issues (TLS)**
- Verify TLS certificate is valid and not expired
- Check certificate subject name matches server address
- Ensure certificate chain and root CA trust
- Verify TLS version compatibility

**Import Fails**
- Ensure you're using the numeric resource ID, not the name
- Verify the TURN server exists in the Infinity cluster
- Check provider authentication credentials have access to the resource

**Connectivity Testing**
- Use STUN/TURN testing tools to verify server functionality
- Test from different network locations and NAT scenarios
- Verify allocation and binding requests work correctly
- Check relay functionality with sample media streams

**Shared Secret Configuration**
- Verify shared secret matches TURN server configuration
- Check time-based credential generation algorithm
- Ensure clocks are synchronized between systems
- Test credential generation and validation

**Geographic Distribution Issues**
- Monitor latency from different client locations
- Ensure TURN servers are optimally placed
- Consider anycast or load balancing for global access
- Test failover between regional TURN servers

**NAT Traversal Problems**
- Verify TURN server can handle symmetric NAT scenarios
- Check allocation lifetime and refresh mechanisms
- Ensure proper STUN binding maintenance
- Test with various NAT types and configurations

**Load Balancing Considerations**
- Implement health checks for TURN server availability
- Configure proper session affinity if required
- Monitor load distribution across multiple servers
- Plan for TURN server failover scenarios