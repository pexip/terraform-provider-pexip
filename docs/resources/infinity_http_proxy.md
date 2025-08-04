---
page_title: "pexip_infinity_http_proxy Resource - terraform-provider-pexip"
subcategory: ""
description: |-
  Manages a Pexip Infinity HTTP proxy configuration.
---

# pexip_infinity_http_proxy (Resource)

Manages a Pexip Infinity HTTP proxy configuration. HTTP proxies are used to route outbound HTTP and HTTPS traffic through an intermediate server for purposes such as content filtering, security scanning, bandwidth management, or accessing resources in restricted network environments. This is commonly required in enterprise environments where direct internet access is not permitted.

## Example Usage

### Basic HTTP Proxy

```terraform
resource "pexip_infinity_http_proxy" "corporate_proxy" {
  name     = "Corporate HTTP Proxy"
  address  = "proxy.company.com"
  port     = 8080
  protocol = "http"
}
```

### HTTPS Proxy with Authentication

```terraform
resource "pexip_infinity_http_proxy" "secure_proxy" {
  name     = "Secure HTTP Proxy"
  address  = "secure-proxy.company.com"
  port     = 8443
  protocol = "https"
  username = "pexip-service"
  password = var.proxy_password
}
```

### Multiple HTTP Proxies for Different Environments

```terraform
# Production proxy
resource "pexip_infinity_http_proxy" "production" {
  name     = "Production HTTP Proxy"
  address  = "proxy-prod.company.com"
  port     = 3128
  protocol = "http"
  username = "prod-pexip"
  password = var.prod_proxy_password
}

# Development proxy
resource "pexip_infinity_http_proxy" "development" {
  name     = "Development HTTP Proxy"
  address  = "proxy-dev.company.com"
  port     = 3128
  protocol = "http"
  username = "dev-pexip"
  password = var.dev_proxy_password
}

# Test proxy without authentication
resource "pexip_infinity_http_proxy" "test" {
  name     = "Test HTTP Proxy"
  address  = "proxy-test.company.com"
  port     = 8080
  protocol = "http"
}
```

### Enterprise Proxy Configuration

```terraform
variable "proxy_config" {
  type = object({
    address  = string
    port     = number
    protocol = string
    username = string
  })
  default = {
    address  = "enterprise-proxy.company.com"
    port     = 8080
    protocol = "http"
    username = "pexip-infinity"
  }
}

resource "pexip_infinity_http_proxy" "enterprise" {
  name     = "Enterprise HTTP Proxy"
  address  = var.proxy_config.address
  port     = var.proxy_config.port
  protocol = var.proxy_config.protocol
  username = var.proxy_config.username
  password = var.enterprise_proxy_password
}
```

### Proxy with Standard Ports

```terraform
# HTTP proxy using standard port (implied port 80)
resource "pexip_infinity_http_proxy" "standard_http" {
  name     = "Standard HTTP Proxy"
  address  = "http-proxy.company.com"
  protocol = "http"
}

# HTTPS proxy using standard port (implied port 443)
resource "pexip_infinity_http_proxy" "standard_https" {
  name     = "Standard HTTPS Proxy"
  address  = "https-proxy.company.com"
  protocol = "https"
  username = "secure-user"
  password = var.secure_proxy_password
}
```

### Conditional Proxy Configuration

```terraform
locals {
  use_proxy = var.environment == "production" || var.environment == "staging"
}

resource "pexip_infinity_http_proxy" "conditional" {
  count    = local.use_proxy ? 1 : 0
  name     = "${title(var.environment)} HTTP Proxy"
  address  = var.proxy_address
  port     = var.proxy_port
  protocol = var.proxy_protocol
  username = var.proxy_username
  password = var.proxy_password
}
```

### Regional Proxy Configuration

```terraform
variable "regional_proxies" {
  type = map(object({
    name     = string
    address  = string
    port     = number
    username = string
  }))
  default = {
    us-east = {
      name     = "US East HTTP Proxy"
      address  = "proxy-useast.company.com"
      port     = 8080
      username = "useast-pexip"
    }
    us-west = {
      name     = "US West HTTP Proxy"
      address  = "proxy-uswest.company.com"
      port     = 8080
      username = "uswest-pexip"
    }
    europe = {
      name     = "Europe HTTP Proxy"
      address  = "proxy-eu.company.com"
      port     = 8080
      username = "eu-pexip"
    }
  }
}

resource "pexip_infinity_http_proxy" "regional" {
  for_each = var.regional_proxies
  name     = each.value.name
  address  = each.value.address
  port     = each.value.port
  protocol = "http"
  username = each.value.username
  password = var.regional_proxy_passwords[each.key]
}
```

### Failover Proxy Configuration

```terraform
# Primary proxy
resource "pexip_infinity_http_proxy" "primary" {
  name     = "Primary HTTP Proxy"
  address  = "proxy1.company.com"
  port     = 8080
  protocol = "http"
  username = "pexip-primary"
  password = var.primary_proxy_password
}

# Secondary proxy for failover
resource "pexip_infinity_http_proxy" "secondary" {
  name     = "Secondary HTTP Proxy"
  address  = "proxy2.company.com"
  port     = 8080
  protocol = "http"
  username = "pexip-secondary"
  password = var.secondary_proxy_password
}
```

## Schema

### Required

- `name` (String) - The name used to refer to this HTTP proxy. Maximum length: 250 characters.
- `address` (String) - The address or hostname of the HTTP proxy. Maximum length: 255 characters.
- `protocol` (String) - The protocol for the HTTP proxy. Valid values: `http`, `https`.

### Optional

- `port` (Number) - The port number for the HTTP proxy. Range: 1 to 65535. Defaults to standard ports (80 for HTTP, 443 for HTTPS).
- `username` (String) - Username for authentication to the HTTP proxy. Maximum length: 100 characters.
- `password` (String, Sensitive) - Password for authentication to the HTTP proxy. Maximum length: 100 characters.

### Read-Only

- `id` (String) - Resource URI for the HTTP proxy in Infinity.
- `resource_id` (Number) - The resource integer identifier for the HTTP proxy in Infinity.

## Import

Import is supported using the following syntax:

```shell
terraform import pexip_infinity_http_proxy.example 123
```

Where `123` is the numeric resource ID of the HTTP proxy.

## Usage Notes

### Protocol Selection

- **HTTP**: Standard unencrypted proxy protocol, suitable for internal networks
- **HTTPS**: Encrypted proxy protocol, required when proxy server uses TLS/SSL
- **Port Defaults**: HTTP defaults to port 80, HTTPS defaults to port 443
- **Security**: Use HTTPS when transmitting sensitive data through the proxy

### Authentication Methods

- **Basic Authentication**: Username and password are supported for proxy authentication
- **No Authentication**: Leave username and password empty for unauthenticated proxies
- **Credential Security**: Store passwords securely using Terraform variables or external secret management
- **Service Accounts**: Use dedicated service accounts for proxy authentication

### Proxy Server Requirements

- **Protocol Support**: Proxy server must support HTTP CONNECT method for HTTPS tunneling
- **Access Control**: Configure proxy server to allow Pexip Infinity traffic
- **Performance**: Ensure proxy server has adequate capacity for expected traffic volume
- **Reliability**: Use high-availability proxy configurations for production environments

### Network Configuration

- **Firewall Rules**: Configure firewall to allow traffic from Pexip nodes to proxy server
- **DNS Resolution**: Ensure proxy server hostname resolves correctly from Pexip nodes
- **Routing**: Verify network routing allows connectivity to proxy server
- **Port Access**: Ensure proxy server port is accessible and not blocked

### Enterprise Integration

- **Policy Compliance**: Ensure proxy configuration meets corporate security policies
- **Content Filtering**: Proxy may filter or block certain types of traffic
- **Logging**: Proxy servers typically log all traffic for security and compliance
- **Bandwidth Management**: Proxy may implement bandwidth limiting or QoS policies

### Performance Considerations

- **Latency**: Proxy adds network latency to all HTTP/HTTPS requests
- **Throughput**: Proxy server capacity affects overall throughput
- **Connection Pooling**: Modern proxies support connection pooling for better performance
- **Caching**: Some proxies provide caching to improve performance

### Security Best Practices

- **Encrypted Connections**: Use HTTPS proxy protocol when possible
- **Credential Management**: Rotate proxy credentials regularly
- **Access Logging**: Monitor proxy access logs for security events
- **Network Segmentation**: Place proxy servers in appropriate network segments

## Troubleshooting

### Common Issues

**HTTP Proxy Creation Fails**
- Verify the address format is correct (hostname or IP address)
- Ensure the protocol is either "http" or "https"
- Check that the port number is within the valid range (1-65535)
- Verify username and password don't exceed maximum length limits

**Proxy Connection Failures**
- Test connectivity from Pexip nodes to the proxy server address and port
- Verify firewall rules allow traffic to the proxy server
- Check that the proxy server is running and accepting connections
- Ensure DNS resolution works for the proxy server hostname

**Authentication Failures**
- Verify the username and password are correct
- Check that the proxy server is configured to accept the provided credentials
- Ensure the proxy server supports HTTP Basic Authentication
- Monitor proxy server logs for authentication errors

**Traffic Not Going Through Proxy**
- Verify Pexip Infinity is configured to use the HTTP proxy
- Check that the proxy configuration is applied to the correct services
- Ensure no bypass rules are preventing proxy usage
- Monitor network traffic to confirm routing through proxy

**Slow Performance Through Proxy**
- Monitor proxy server performance and resource utilization
- Check network latency between Pexip nodes and proxy server
- Verify proxy server has adequate bandwidth capacity
- Consider proxy server optimization or load balancing

**HTTPS Proxy Certificate Issues**
- Verify the proxy server has a valid SSL/TLS certificate
- Check certificate trust chain and expiration dates
- Ensure certificate hostname matches the proxy server address
- Monitor TLS handshake logs for certificate errors

**Proxy Server Unreachable**
- Test network connectivity using ping or telnet from Pexip nodes
- Verify proxy server is running and listening on the configured port
- Check network routing and firewall configurations
- Ensure proxy server is not overloaded or experiencing issues

**Authentication Method Not Supported**
- Verify the proxy server supports HTTP Basic Authentication
- Check if the proxy requires different authentication methods
- Ensure authentication credentials are properly formatted
- Monitor proxy server logs for authentication method errors

**Content Filtering Issues**
- Check if proxy server is blocking required Pexip traffic
- Verify proxy content filtering policies allow necessary protocols
- Review proxy server logs for blocked requests
- Configure proxy server to allow Pexip-specific traffic patterns

**Import Fails**
- Ensure you're using the numeric resource ID, not the proxy name
- Verify the HTTP proxy exists in the Infinity cluster
- Check provider authentication credentials have access to the resource
- Note that password information may not be available during import

**High Latency Issues**
- Monitor network latency to the proxy server
- Check proxy server processing time for requests
- Verify network path optimization between Pexip and proxy
- Consider deploying proxy servers closer to Pexip nodes

**Proxy Failover Problems**
- Implement multiple proxy configurations for redundancy
- Monitor proxy server availability and health
- Configure automatic failover mechanisms where supported
- Test failover scenarios regularly

**Protocol Compatibility Issues**
- Verify proxy server supports required HTTP/HTTPS protocols
- Check for HTTP version compatibility (HTTP/1.1, HTTP/2)
- Ensure proxy supports CONNECT method for HTTPS tunneling
- Monitor protocol negotiation in proxy server logs