---
page_title: "pexip_infinity_event_sink Resource - terraform-provider-pexip"
subcategory: ""
description: |-
  Manages a Pexip Infinity event sink configuration.
---

# pexip_infinity_event_sink (Resource)

Manages a Pexip Infinity event sink configuration. Event sinks are external systems that receive real-time events from Pexip Infinity, enabling integration with monitoring systems, analytics platforms, billing systems, and other third-party applications. Events include conference lifecycle events, participant join/leave events, media statistics, and system status updates.

## Example Usage

### Basic Event Sink Configuration

```terraform
resource "pexip_infinity_event_sink" "basic_sink" {
  name = "Basic Event Sink"
  url  = "https://events.example.com/webhook"
}
```

### Event Sink with Authentication

```terraform
resource "pexip_infinity_event_sink" "authenticated_sink" {
  name        = "Authenticated Event Sink"
  description = "Event sink with HTTP basic authentication"
  url         = "https://events.company.com/api/pexip-events"
  username    = "pexip-integration"
  password    = "secure-password-123"
  version     = 1
}
```

### Advanced Event Sink with Bulk Support

```terraform
resource "pexip_infinity_event_sink" "bulk_sink" {
  name                   = "Bulk Event Sink"
  description            = "High-volume event sink with bulk processing"
  url                    = "https://analytics.company.com/events/bulk"
  username               = "analytics-user"
  password               = var.analytics_password
  bulk_support           = true
  verify_tls_certificate = true
  version                = 2
}
```

### Multiple Event Sinks for Different Purposes

```terraform
# Analytics event sink
resource "pexip_infinity_event_sink" "analytics" {
  name                   = "Analytics Event Sink"
  description            = "Event sink for call analytics and reporting"
  url                    = "https://analytics.company.com/pexip/events"
  username               = "analytics-service"
  password               = var.analytics_password
  bulk_support           = true
  verify_tls_certificate = true
  version                = 2
}

# Monitoring event sink
resource "pexip_infinity_event_sink" "monitoring" {
  name        = "Monitoring Event Sink"
  description = "Event sink for real-time monitoring"
  url         = "https://monitoring.company.com/webhooks/pexip"
  username    = "monitoring-user"
  password    = var.monitoring_password
  version     = 1
}

# Billing event sink
resource "pexip_infinity_event_sink" "billing" {
  name                   = "Billing Event Sink"
  description            = "Event sink for usage tracking and billing"
  url                    = "https://billing.company.com/api/usage-events"
  username               = "billing-service"
  password               = var.billing_password
  bulk_support           = true
  verify_tls_certificate = true
  version                = 2
}
```

### Enterprise Event Integration

```terraform
variable "event_sinks" {
  type = list(object({
    name                   = string
    description            = string
    url                    = string
    username               = string
    bulk_support           = bool
    verify_tls_certificate = bool
    version                = number
  }))
  default = [
    {
      name                   = "SIEM Integration"
      description            = "Security information and event management"
      url                    = "https://siem.company.com/api/events"
      username               = "siem-integration"
      bulk_support           = true
      verify_tls_certificate = true
      version                = 2
    },
    {
      name                   = "Data Lake"
      description            = "Long-term storage for analytics"
      url                    = "https://datalake.company.com/ingest/pexip"
      username               = "datalake-service"
      bulk_support           = true
      verify_tls_certificate = true
      version                = 2
    }
  ]
}

resource "pexip_infinity_event_sink" "enterprise" {
  count                  = length(var.event_sinks)
  name                   = var.event_sinks[count.index].name
  description            = var.event_sinks[count.index].description
  url                    = var.event_sinks[count.index].url
  username               = var.event_sinks[count.index].username
  password               = var.event_sink_passwords[count.index]
  bulk_support           = var.event_sinks[count.index].bulk_support
  verify_tls_certificate = var.event_sinks[count.index].verify_tls_certificate
  version                = var.event_sinks[count.index].version
}
```

### Development and Testing Event Sink

```terraform
resource "pexip_infinity_event_sink" "development" {
  name                   = "Development Event Sink"
  description            = "Event sink for development and testing"
  url                    = "https://webhook.site/unique-id"
  verify_tls_certificate = false
  version                = 1
}
```

## Schema

### Required

- `name` (String) - The name used to refer to this event sink. Maximum length: 250 characters.
- `url` (String) - The URL for the event sink. Maximum length: 500 characters.

### Optional

- `description` (String) - A description of the event sink. Maximum length: 250 characters.
- `username` (String) - Username for authentication to the event sink. Maximum length: 100 characters.
- `password` (String, Sensitive) - Password for authentication to the event sink. Maximum length: 100 characters.
- `bulk_support` (Boolean) - Whether the event sink supports bulk operations. Defaults to `false`.
- `verify_tls_certificate` (Boolean) - Whether to verify TLS certificates when connecting to the event sink. Defaults to `false`.
- `version` (Number) - The version of the event sink API. Must be at least 1. Defaults to `1`.

### Read-Only

- `id` (String) - Resource URI for the event sink in Infinity.
- `resource_id` (Number) - The resource integer identifier for the event sink in Infinity.

## Import

Import is supported using the following syntax:

```shell
terraform import pexip_infinity_event_sink.example 123
```

Where `123` is the numeric resource ID of the event sink.

## Usage Notes

### Event Types

Pexip Infinity generates various types of events that can be sent to event sinks:

- **Conference Events**: Conference creation, start, end, and modification
- **Participant Events**: Participant join, leave, mute, unmute, and role changes
- **Media Events**: Media statistics, quality metrics, and codec information
- **System Events**: Node status changes, resource utilization, and health events
- **Security Events**: Authentication attempts, access control, and policy violations

### Event Sink Versions

- **Version 1**: Basic event format with essential information
- **Version 2**: Enhanced event format with additional metadata and statistics
- **Future Versions**: Support for new event types and enhanced data formats

### Bulk Support

- **Individual Events**: Each event is sent as a separate HTTP request
- **Bulk Events**: Multiple events are batched together in a single HTTP request
- **Performance**: Bulk support reduces network overhead and improves performance
- **Processing**: Event sink must be designed to handle bulk event arrays

### TLS Certificate Verification

- **Production**: Always enable TLS certificate verification in production environments
- **Development**: May disable for testing with self-signed certificates
- **Security**: Certificate verification prevents man-in-the-middle attacks
- **Trust Store**: Ensure proper certificate authority trust is configured

### Authentication

- **HTTP Basic Auth**: Username and password are sent with each request
- **API Keys**: Can be included in the URL or as custom headers
- **Token-Based**: Implement token-based authentication in your event sink
- **Security**: Always use HTTPS when authentication credentials are provided

### Event Delivery

- **Reliable Delivery**: Pexip attempts to deliver events reliably with retries
- **Failure Handling**: Failed events may be queued and retried
- **Monitoring**: Monitor event delivery success rates and error conditions
- **Timeouts**: Configure appropriate timeout values for event sink responses

### Performance Considerations

- **High Volume**: Consider bulk support for high-volume environments
- **Response Time**: Event sink should respond quickly to avoid blocking
- **Scalability**: Design event sink to handle peak conference loads
- **Queue Management**: Implement proper queue management for event processing

## Troubleshooting

### Common Issues

**Event Sink Creation Fails**
- Verify the URL format is correct and accessible
- Ensure the username and password don't exceed maximum length limits
- Check that the version number is valid (>= 1)
- Verify the description doesn't exceed the character limit

**Events Not Being Delivered**
- Verify the event sink URL is accessible from Pexip Infinity nodes
- Check firewall rules allow HTTPS traffic to the event sink
- Ensure the event sink is responding with appropriate HTTP status codes
- Monitor Pexip Infinity logs for event delivery errors

**Authentication Failures**
- Verify username and password are correct
- Check that the event sink properly handles HTTP Basic Authentication
- Ensure credentials have appropriate permissions on the event sink
- Monitor event sink logs for authentication errors

**TLS Connection Issues**
- Verify the event sink has a valid SSL/TLS certificate
- Check certificate expiration dates and trust chain
- Ensure TLS cipher suites are compatible
- Consider disabling certificate verification for testing (not recommended for production)

**Performance Issues**
- Enable bulk support if the event sink supports it
- Monitor event delivery latency and success rates
- Check event sink response times and processing capacity
- Consider implementing event queuing and async processing

**Event Processing Errors**
- Verify the event sink can handle the event format (version 1 or 2)
- Check event sink parsing and validation logic
- Monitor for malformed or unexpected event data
- Implement proper error handling and logging in the event sink

**High Event Volume Issues**
- Enable bulk support to reduce request overhead
- Implement proper queue management in the event sink
- Monitor event sink resource utilization and scaling
- Consider load balancing for high-volume scenarios

**Network Connectivity Problems**
- Test connectivity from Pexip Infinity nodes to the event sink URL
- Verify DNS resolution for the event sink hostname
- Check network routing and firewall configurations
- Monitor network latency and packet loss

**Certificate Verification Failures**
- Ensure the event sink certificate is issued by a trusted CA
- Check certificate validity period and expiration
- Verify the certificate hostname matches the event sink URL
- Update certificate trust store if using private CAs

**Import Fails**
- Ensure you're using the numeric resource ID, not the event sink name
- Verify the event sink exists in the Infinity cluster
- Check provider authentication credentials have access to the resource
- Note that password information may not be available during import

**Event Sink Unavailable**
- Implement proper error handling in the event sink application
- Monitor event sink uptime and availability
- Configure appropriate health checks and alerting
- Consider implementing event sink redundancy

**Data Processing Issues**
- Validate event data format and structure in the event sink
- Implement proper error handling for malformed events
- Monitor event processing rates and error conditions
- Ensure adequate storage and processing capacity