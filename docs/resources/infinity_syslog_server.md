---
page_title: "pexip_infinity_syslog_server Resource - terraform-provider-pexip"
subcategory: ""
description: |-
  Manages a Pexip Infinity syslog server configuration.
---

# pexip_infinity_syslog_server (Resource)

Manages a syslog server with the Infinity service. Syslog servers receive system logs and audit information from Pexip Infinity for centralized logging and monitoring. This enables comprehensive monitoring, troubleshooting, and compliance auditing of your Pexip Infinity deployment.

## Example Usage

### Basic UDP Syslog Server

```terraform
resource "pexip_infinity_syslog_server" "basic_syslog" {
  address   = "syslog.company.com"
  port      = 514
  transport = "udp"
}
```

### Comprehensive Syslog Configuration

```terraform
resource "pexip_infinity_syslog_server" "central_logging" {
  address     = "logs.company.com"
  description = "Central logging server for all Pexip logs"
  port        = 514
  transport   = "udp"
  audit_log   = true
  support_log = true
  web_log     = true
}
```

### Secure TLS Syslog Server

```terraform
resource "pexip_infinity_syslog_server" "secure_syslog" {
  address     = "secure-logs.company.com"
  description = "Secure syslog server with TLS encryption"
  port        = 6514
  transport   = "tls"
  audit_log   = true
  support_log = false
  web_log     = false
}
```

### Multiple Syslog Servers for Different Log Types

```terraform
# Audit logs to compliance system
resource "pexip_infinity_syslog_server" "audit_logs" {
  address     = "audit.company.com"
  description = "Compliance audit logging system"
  port        = 514
  transport   = "tcp"
  audit_log   = true
  support_log = false
  web_log     = false
}

# Support logs to monitoring system
resource "pexip_infinity_syslog_server" "support_logs" {
  address     = "monitoring.company.com"
  description = "Technical support and monitoring logs"
  port        = 514
  transport   = "udp"
  audit_log   = false
  support_log = true
  web_log     = true
}

# High-priority logs to SIEM
resource "pexip_infinity_syslog_server" "siem_logs" {
  address     = "siem.company.com"
  description = "Security Information and Event Management"
  port        = 1514
  transport   = "tls"
  audit_log   = true
  support_log = true
  web_log     = true
}
```

### Regional Syslog Servers

```terraform
# Different syslog servers per region
resource "pexip_infinity_syslog_server" "regional_syslog" {
  for_each = var.regional_syslog_config
  
  address     = each.value.address
  description = "Regional syslog server for ${each.key}"
  port        = each.value.port
  transport   = each.value.transport
  audit_log   = true
  support_log = true
  web_log     = false
}
```

### Load-Balanced Syslog Configuration

```terraform
resource "pexip_infinity_syslog_server" "lb_syslog" {
  count = length(var.syslog_servers)
  
  address     = var.syslog_servers[count.index].address
  description = "Load-balanced syslog server ${count.index + 1}"
  port        = var.syslog_servers[count.index].port
  transport   = "udp"
  audit_log   = true
  support_log = true
  web_log     = true
}
```

## Schema

### Required

- `address` (String) - The IP address or FQDN of the remote syslog server. Maximum length: 255 characters.
- `port` (Number) - The port number for syslog communications. Valid range: 1-65535.
- `transport` (String) - Transport protocol for syslog. Valid values: udp, tcp, tls.

### Optional

- `description` (String) - Description of the syslog server. Maximum length: 500 characters.
- `audit_log` (Boolean) - Whether to send audit logs to this syslog server. Defaults to `false`.
- `support_log` (Boolean) - Whether to send support logs to this syslog server. Defaults to `false`.
- `web_log` (Boolean) - Whether to send web logs to this syslog server. Defaults to `false`.

### Read-Only

- `id` (String) - Resource URI for the syslog server in Infinity.
- `resource_id` (Number) - The resource integer identifier for the syslog server in Infinity.

## Import

Import is supported using the following syntax:

```shell
terraform import pexip_infinity_syslog_server.example 123
```

Where `123` is the numeric resource ID of the syslog server.

## Usage Notes

### Log Types
- **Audit Logs**: Security events, authentication, authorization, configuration changes
- **Support Logs**: System events, errors, warnings, diagnostic information
- **Web Logs**: HTTP access logs, web interface usage, API calls

### Transport Protocols
- **UDP (514)**: Fastest but no delivery guarantee, most common for syslog
- **TCP (514)**: Reliable delivery, better for critical logs
- **TLS (6514)**: Encrypted transmission, required for sensitive environments

### Syslog Server Requirements
- Ensure syslog servers can handle the expected log volume
- Configure appropriate log retention policies
- Implement log rotation to prevent disk space issues
- Set up monitoring and alerting for syslog server health

### System Location Assignment
- Syslog servers are assigned to system locations
- All nodes in a system location send logs to the configured syslog servers
- Multiple syslog servers can be configured for redundancy
- Different log types can be sent to different servers

### Performance Considerations
- Use UDP for high-volume logging when some log loss is acceptable
- Use TCP or TLS for critical logs that must be delivered
- Monitor network bandwidth usage for syslog traffic
- Consider log compression for high-volume environments

### Security Considerations
- Use TLS encryption for sensitive log data
- Implement proper access controls on syslog servers
- Consider log integrity and tamper detection
- Ensure syslog servers are properly secured and patched

## Troubleshooting

### Common Issues

**Syslog Server Creation Fails**
- Verify the address format is correct (hostname or IP address)
- Ensure the port number is within the valid range (1-65535)
- Check that the transport protocol is one of: udp, tcp, tls
- Verify the description doesn't exceed 500 characters

**Logs Not Being Received**
- Test syslog server connectivity from conferencing nodes
- Verify firewall rules allow syslog traffic on the specified port
- Check if the syslog server service is running and configured correctly
- Ensure the correct log types are enabled (audit_log, support_log, web_log)

**Transport Protocol Issues**
- **UDP**: Check for packet loss, firewall blocking, or server overload
- **TCP**: Verify TCP connections can be established and maintained
- **TLS**: Ensure certificates are valid and TLS configuration is correct

**Performance Problems**
- Monitor syslog server CPU and memory usage
- Check disk space and I/O performance on log storage
- Verify network bandwidth is sufficient for log volume
- Consider implementing log rate limiting if necessary

**Log Volume Issues**
- Monitor daily log volume and storage requirements
- Implement log rotation and archival policies
- Consider filtering or sampling for very high-volume environments
- Use compression to reduce storage requirements

**Network Connectivity Problems**
- Verify syslog server is reachable from all system locations
- Check network routing and firewall rules
- Test connectivity using tools like telnet or nc
- Monitor for network congestion affecting log delivery

**Import Fails**
- Ensure you're using the numeric resource ID, not the address
- Verify the syslog server exists in the Infinity cluster
- Check provider authentication credentials have access to the resource

**TLS Configuration Issues**
- Verify TLS certificates are valid and not expired
- Check that certificate subject names match the server address
- Ensure proper TLS version and cipher suite compatibility
- Verify certificate chain and root CA trust

**Log Format Problems**
- Check syslog server configuration accepts RFC 3164 or RFC 5424 format
- Verify timestamp format compatibility
- Ensure character encoding is handled correctly
- Test log parsing and field extraction

**Compliance and Retention Issues**
- Implement appropriate log retention periods for compliance requirements
- Ensure log integrity and tamper detection mechanisms
- Configure proper backup and archival procedures
- Monitor compliance with data protection regulations

**Monitoring and Alerting Setup**
- Configure alerts for syslog server unavailability
- Monitor log ingestion rates and patterns
- Set up alerts for unusual log patterns or errors
- Implement log analysis and correlation tools