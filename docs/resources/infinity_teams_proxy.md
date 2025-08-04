---
page_title: "pexip_infinity_teams_proxy Resource - terraform-provider-pexip"
subcategory: ""
description: |-
  Manages a Pexip Infinity Teams proxy configuration.
---

# pexip_infinity_teams_proxy (Resource)

Manages a Pexip Infinity Teams proxy configuration. Teams proxies enable Pexip Infinity to integrate with Microsoft Teams, providing seamless interoperability between Teams and other video conferencing platforms. The proxy handles authentication, protocol translation, and service discovery for Teams integration.

## Example Usage

### Basic Teams Proxy Configuration

```terraform
resource "pexip_infinity_teams_proxy" "primary_teams_proxy" {
  name                     = "Primary Teams Proxy"
  address                  = "teams-proxy.example.com"
  port                     = 443
  azure_tenant             = "contoso.onmicrosoft.com"
  min_number_of_instances  = 2
}
```

### Teams Proxy with Event Hub Integration

```terraform
resource "pexip_infinity_teams_proxy" "teams_proxy_with_events" {
  name                     = "Teams Proxy with Events"
  description              = "Teams proxy with event hub integration"
  address                  = "teams-proxy.company.com"
  port                     = 443
  azure_tenant             = "company.onmicrosoft.com"
  eventhub_id              = "teams-events-hub"
  min_number_of_instances  = 3
  notifications_enabled    = true
  notifications_queue      = "teams-notifications-queue"
}
```

### High Availability Teams Proxy

```terraform
resource "pexip_infinity_teams_proxy" "ha_teams_proxy" {
  name                     = "HA Teams Proxy"
  description              = "High availability Teams proxy configuration"
  address                  = "teams-ha.example.com"
  port                     = 443
  azure_tenant             = "example.onmicrosoft.com"
  min_number_of_instances  = 5
  notifications_enabled    = true
}
```

### Multi-Tenant Teams Proxy Setup

```terraform
variable "teams_tenants" {
  type = list(object({
    name         = string
    azure_tenant = string
    instances    = number
  }))
  default = [
    {
      name         = "Tenant A Teams Proxy"
      azure_tenant = "tenanta.onmicrosoft.com"
      instances    = 2
    },
    {
      name         = "Tenant B Teams Proxy"
      azure_tenant = "tenantb.onmicrosoft.com"
      instances    = 3
    }
  ]
}

resource "pexip_infinity_teams_proxy" "multi_tenant" {
  count                    = length(var.teams_tenants)
  name                     = var.teams_tenants[count.index].name
  description              = "Teams proxy for ${var.teams_tenants[count.index].azure_tenant}"
  address                  = "teams-proxy-${count.index + 1}.example.com"
  port                     = 443
  azure_tenant             = var.teams_tenants[count.index].azure_tenant
  min_number_of_instances  = var.teams_tenants[count.index].instances
  notifications_enabled    = true
}
```

### Enterprise Teams Integration

```terraform
# Primary Teams proxy for production
resource "pexip_infinity_teams_proxy" "production_teams" {
  name                     = "Production Teams Proxy"
  description              = "Production Teams proxy for enterprise"
  address                  = "teams-prod.enterprise.com"
  port                     = 443
  azure_tenant             = "enterprise.onmicrosoft.com"
  eventhub_id              = "production-teams-events"
  min_number_of_instances  = 10
  notifications_enabled    = true
  notifications_queue      = "prod-teams-notifications"
}

# Development Teams proxy
resource "pexip_infinity_teams_proxy" "development_teams" {
  name                     = "Development Teams Proxy"
  description              = "Development Teams proxy for testing"
  address                  = "teams-dev.enterprise.com"
  port                     = 443
  azure_tenant             = "enterprise-dev.onmicrosoft.com"
  min_number_of_instances  = 1
  notifications_enabled    = false
}
```

## Schema

### Required

- `name` (String) - The name used to refer to this Teams proxy. Maximum length: 250 characters.
- `address` (String) - The address or hostname of the Teams proxy. Maximum length: 255 characters.
- `port` (Number) - The port number for the Teams proxy. Range: 1 to 65535.
- `azure_tenant` (String) - The Azure tenant ID for the Teams proxy. Maximum length: 255 characters.
- `min_number_of_instances` (Number) - The minimum number of instances for the Teams proxy. Must be at least 1.

### Optional

- `description` (String) - A description of the Teams proxy. Maximum length: 250 characters.
- `eventhub_id` (String) - The event hub identifier for the Teams proxy. Maximum length: 255 characters.
- `notifications_enabled` (Boolean) - Whether notifications are enabled for the Teams proxy. Defaults to `false`.
- `notifications_queue` (String, Sensitive) - The notification queue name for the Teams proxy. Maximum length: 255 characters.

### Read-Only

- `id` (String) - Resource URI for the Teams proxy in Infinity.
- `resource_id` (Number) - The resource integer identifier for the Teams proxy in Infinity.

## Import

Import is supported using the following syntax:

```shell
terraform import pexip_infinity_teams_proxy.example 123
```

Where `123` is the numeric resource ID of the Teams proxy.

## Usage Notes

### Azure Tenant Configuration

- **Tenant ID Format**: Use the full Microsoft tenant domain (e.g., "company.onmicrosoft.com")
- **Multi-Tenant Support**: Each Teams proxy can only serve one Azure tenant
- **Tenant Validation**: Ensure the Azure tenant exists and is properly configured for Teams integration
- **Permissions**: The Azure tenant must have appropriate permissions for Pexip integration

### Instance Scaling

- **Minimum Instances**: Configure based on expected concurrent Teams users
- **Performance**: Each instance can typically handle multiple concurrent Teams sessions
- **High Availability**: Use minimum 2 instances for production environments
- **Auto-Scaling**: Pexip Infinity can automatically scale beyond the minimum based on load

### Event Hub Integration

- **Event Streaming**: Event Hub integration enables real-time Teams event streaming
- **Analytics**: Use events for call analytics, monitoring, and reporting
- **Compliance**: Event data can be used for compliance and audit requirements
- **Azure Service Bus**: Ensure Azure Service Bus is properly configured

### Notification Configuration

- **Queue Integration**: Notifications queue enables integration with external systems
- **Real-Time Updates**: Receive real-time updates about Teams proxy status and events
- **Monitoring**: Use notifications for monitoring and alerting systems
- **Security**: Notification queue credentials are stored securely

### Network Requirements

- **HTTPS**: Teams proxy communication requires HTTPS (typically port 443)
- **Certificate**: Ensure valid SSL/TLS certificates are configured
- **Firewall**: Configure firewall rules to allow Teams proxy traffic
- **DNS**: Proper DNS resolution is critical for Teams integration

### Security Considerations

- **OAuth2**: Teams integration uses OAuth2 for authentication
- **Certificate Validation**: Ensure proper certificate validation is enabled
- **Network Isolation**: Consider network isolation for Teams proxy servers
- **Access Control**: Implement appropriate access control policies

## Troubleshooting

### Common Issues

**Teams Proxy Creation Fails**
- Verify the Azure tenant ID format is correct
- Ensure the address is a valid hostname or IP address
- Check that the port number is within the valid range
- Verify minimum instances is at least 1

**Teams Integration Not Working**
- Verify Azure tenant is correctly configured for Pexip integration
- Check OAuth2 application registration in Azure AD
- Ensure Teams proxy is accessible from Teams infrastructure
- Verify firewall rules allow HTTPS traffic on the configured port

**Authentication Failures**
- Check Azure AD application permissions and consent
- Verify OAuth2 client credentials are configured correctly
- Ensure the Azure tenant has appropriate licensing for Teams integration
- Monitor Azure AD sign-in logs for authentication errors

**Instance Scaling Issues**
- Monitor Teams proxy instance utilization and performance
- Verify minimum instance configuration meets demand
- Check Pexip Infinity licensing for Teams integration
- Monitor auto-scaling behavior and thresholds

**Event Hub Connection Problems**
- Verify Event Hub exists and is properly configured in Azure
- Check Event Hub connection string and authentication
- Ensure network connectivity between Teams proxy and Azure Event Hub
- Monitor Event Hub metrics for connection and throughput issues

**Notification Queue Issues**
- Verify notification queue exists and is accessible
- Check queue authentication credentials and permissions
- Monitor queue message processing and error rates
- Ensure notification queue has sufficient capacity

**High Latency or Performance Issues**
- Monitor network latency between Teams infrastructure and proxy
- Check Teams proxy server performance and resource utilization
- Verify adequate bandwidth for Teams traffic
- Consider geographic proximity of Teams proxy to users

**TLS/SSL Certificate Issues**
- Verify SSL certificate validity and expiration dates
- Check certificate chain and root CA trust
- Ensure certificate matches the Teams proxy hostname
- Monitor TLS handshake logs for certificate errors

**Multi-Tenant Configuration Problems**
- Verify each Azure tenant is configured independently
- Check tenant isolation and permissions
- Ensure proper routing for multi-tenant scenarios
- Monitor cross-tenant access attempts

**Import Fails**
- Ensure you're using the numeric resource ID, not the proxy name
- Verify the Teams proxy exists in the Infinity cluster
- Check provider authentication credentials have access to the resource
- Confirm the Azure tenant configuration is accessible

**Teams Client Connection Issues**
- Verify Teams clients can reach the configured proxy address
- Check Teams client policy configuration for external access
- Ensure proper Teams admin center configuration
- Monitor Teams client logs for connection attempts

**Licensing and Compliance Issues**
- Verify appropriate Teams and Pexip licensing is in place
- Check compliance requirements for Teams integration
- Ensure proper data handling and privacy configurations
- Monitor usage for license compliance