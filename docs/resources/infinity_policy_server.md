---
page_title: "pexip_infinity_policy_server Resource - terraform-provider-pexip"
subcategory: ""
description: |-
  Manages a Pexip Infinity policy server configuration.
---

# pexip_infinity_policy_server (Resource)

Manages a policy server configuration with the Infinity service. Policy servers provide external policy control for Pexip Infinity, enabling custom business logic for call routing, user permissions, participant management, and service configuration through external HTTP-based policy engines.

## Example Usage

### Basic Policy Server

```terraform
resource "pexip_infinity_policy_server" "basic_policy" {
  name = "Basic Policy Server"
  url  = "https://policy.company.com/pexip"
}
```

### Policy Server with Authentication

```terraform
resource "pexip_infinity_policy_server" "authenticated_policy" {
  name        = "Authenticated Policy Server"
  description = "Policy server with HTTP basic authentication"
  url         = "https://policy.company.com/pexip/api"
  username    = "pexip-service"
  password    = var.policy_server_password
}
```

### Service and Participant Lookup Policy

```terraform
resource "pexip_infinity_policy_server" "lookup_policy" {
  name                        = "Directory Lookup Policy"
  description                 = "External directory and service lookup"
  url                         = "https://directory.company.com/policy"
  username                    = "directory-user"
  password                    = var.directory_password
  
  # Enable lookup services
  enable_service_lookup       = true
  enable_participant_lookup   = true
  enable_directory_lookup     = true
  enable_avatar_lookup        = true
}
```

### Internal Policy Configuration

```terraform
resource "pexip_infinity_policy_server" "internal_policy" {
  name                               = "Internal Policy Engine"
  description                        = "Internal policy with custom templates"
  url                                = "https://internal-policy.company.com"
  
  # Enable internal policy features
  enable_internal_service_policy     = true
  enable_internal_participant_policy = true
  enable_internal_media_location_policy = true
  
  # Custom policy templates
  service_configuration_template     = jsonencode({
    default_bandwidth = "2048"
    recording_enabled = true
    encryption_required = true
  })
  
  participant_configuration_template = jsonencode({
    default_role = "guest"
    max_participants = 50
    allow_dial_out = false
  })
  
  media_location_configuration_template = jsonencode({
    preferred_location = "primary"
    failover_enabled = true
    quality_threshold = 0.8
  })
}
```

### Multi-Purpose Policy Server

```terraform
resource "pexip_infinity_policy_server" "comprehensive_policy" {
  name        = "Comprehensive Policy Server"
  description = "Full-featured policy server with all capabilities"
  url         = "https://policy-engine.company.com/api/v2"
  username    = "pexip-integration"
  password    = var.comprehensive_policy_password
  
  # Enable all lookup types
  enable_service_lookup             = true
  enable_participant_lookup         = true
  enable_registration_lookup        = true
  enable_directory_lookup           = true
  enable_avatar_lookup              = true
  enable_media_location_lookup      = true
  
  # Enable internal policies
  enable_internal_service_policy       = true
  enable_internal_participant_policy   = true
  enable_internal_media_location_policy = true
  
  # Avatar preferences
  prefer_local_avatar_configuration = false
  
  # Configuration templates
  service_configuration_template = file("${path.module}/templates/service-policy.json")
  participant_configuration_template = file("${path.module}/templates/participant-policy.json")
  media_location_configuration_template = file("${path.module}/templates/media-location-policy.json")
}
```

### Regional Policy Servers

```terraform
# Different policy servers for different regions
resource "pexip_infinity_policy_server" "regional_policy" {
  for_each = var.regional_policy_servers
  
  name                          = "Policy Server - ${each.key}"
  description                   = "Regional policy server for ${each.key}"
  url                           = each.value.url
  username                      = each.value.username
  password                      = each.value.password
  
  enable_service_lookup         = true
  enable_participant_lookup     = true
  enable_directory_lookup       = each.value.enable_directory
  enable_avatar_lookup          = each.value.enable_avatars
}
```

## Schema

### Required

- `name` (String) - The name used to refer to this policy server. Maximum length: 250 characters.

### Optional

- `description` (String) - A description of the policy server. Maximum length: 250 characters.
- `url` (String) - The URL for the policy server. Maximum length: 500 characters.
- `username` (String) - Username for authentication to the policy server. Maximum length: 100 characters.
- `password` (String, Sensitive) - Password for authentication to the policy server. Maximum length: 100 characters.
- `enable_service_lookup` (Boolean) - Whether to enable service lookup on this policy server. Defaults to `false`.
- `enable_participant_lookup` (Boolean) - Whether to enable participant lookup on this policy server. Defaults to `false`.
- `enable_registration_lookup` (Boolean) - Whether to enable registration lookup on this policy server. Defaults to `false`.
- `enable_directory_lookup` (Boolean) - Whether to enable directory lookup on this policy server. Defaults to `false`.
- `enable_avatar_lookup` (Boolean) - Whether to enable avatar lookup on this policy server. Defaults to `false`.
- `enable_media_location_lookup` (Boolean) - Whether to enable media location lookup on this policy server. Defaults to `false`.
- `enable_internal_service_policy` (Boolean) - Whether to enable internal service policy on this policy server. Defaults to `false`.
- `enable_internal_participant_policy` (Boolean) - Whether to enable internal participant policy on this policy server. Defaults to `false`.
- `enable_internal_media_location_policy` (Boolean) - Whether to enable internal media location policy on this policy server. Defaults to `false`.
- `prefer_local_avatar_configuration` (Boolean) - Whether to prefer local avatar configuration over policy server configuration. Defaults to `false`.
- `service_configuration_template` (String) - Service configuration template. Maximum length: 1000 characters.
- `participant_configuration_template` (String) - Participant configuration template. Maximum length: 1000 characters.
- `registration_configuration_template` (String) - Registration configuration template. Maximum length: 1000 characters.
- `directory_search_template` (String) - Directory search template. Maximum length: 1000 characters.
- `avatar_configuration_template` (String) - Avatar configuration template. Maximum length: 1000 characters.
- `media_location_configuration_template` (String) - Media location configuration template. Maximum length: 1000 characters.

### Read-Only

- `id` (String) - Resource URI for the policy server in Infinity.
- `resource_id` (Number) - The resource integer identifier for the policy server in Infinity.

## Import

Import is supported using the following syntax:

```shell
terraform import pexip_infinity_policy_server.example 123
```

Where `123` is the numeric resource ID of the policy server.

## Usage Notes

### Policy Server Types
- **Lookup Services**: Provide external directory integration and service resolution
- **Internal Policies**: Define business rules using templates and logic
- **Registration Control**: Manage device and user registration policies
- **Media Location**: Control media routing and location selection

### Lookup Service Types
- **Service Lookup**: Resolve conference aliases and virtual meeting rooms
- **Participant Lookup**: Validate and configure participant access
- **Directory Lookup**: Search corporate directories for users and contacts
- **Avatar Lookup**: Retrieve user profile pictures and avatars
- **Media Location Lookup**: Determine optimal media processing locations

### Authentication
- HTTP Basic Authentication with username/password
- Ensure credentials have appropriate permissions on policy server
- Use HTTPS URLs to protect authentication credentials
- Consider using service accounts with limited permissions

### Policy Templates
- JSON-formatted configuration templates for different policy types
- Templates can include variables and conditional logic
- Used by internal policy engines for dynamic configuration
- Should be validated for proper JSON syntax and required fields

### Performance Considerations
- Policy servers are called for each relevant event
- Ensure policy servers can handle expected request volume
- Implement proper timeout handling and fallback mechanisms
- Monitor policy server response times and availability

### Security Considerations
- Use HTTPS for all policy server communications
- Implement proper authentication and authorization
- Validate all policy server responses
- Log policy decisions for audit and troubleshooting

## Troubleshooting

### Common Issues

**Policy Server Creation Fails**
- Verify the policy server name is unique
- Ensure URL format is correct and accessible
- Check that all field lengths are within specified limits
- Verify authentication credentials if provided

**Policy Server Not Responding**
- Test policy server URL accessibility from Pexip nodes
- Verify HTTP/HTTPS connectivity and certificate validity
- Check authentication credentials and permissions
- Ensure policy server is running and configured correctly

**Authentication Failures**
- Verify username and password are correct
- Check that authentication method is supported by policy server
- Ensure credentials have appropriate permissions
- Test authentication independently of Pexip

**Lookup Service Issues**
- Verify appropriate lookup types are enabled
- Check policy server endpoints for lookup requests
- Ensure proper request/response format implementation
- Test lookup functionality independently

**Template Configuration Problems**
- Validate JSON syntax in configuration templates
- Ensure required fields are present in templates
- Check template variable substitution
- Test template rendering with sample data

**Performance Problems**
- Monitor policy server response times
- Check for policy server resource constraints
- Verify network connectivity and latency
- Consider implementing caching where appropriate

**Import Fails**
- Ensure you're using the numeric resource ID, not the name
- Verify the policy server exists in the Infinity cluster
- Check provider authentication credentials have access to the resource

**SSL/TLS Issues**
- Verify certificate validity and trust chain
- Check certificate subject name matches URL
- Ensure proper SSL/TLS version support
- Verify certificate is not expired or revoked

**Policy Logic Errors**
- Review policy server logs for errors and warnings
- Test policy logic with various input scenarios
- Verify policy responses match expected format
- Check for proper error handling in policy code

**Integration Problems**
- Verify policy server API compatibility with Pexip requirements
- Check request/response formats match specifications
- Ensure proper HTTP status codes are returned
- Test integration with sample data

**Failover and Redundancy**
- Implement multiple policy servers for high availability
- Configure proper failover mechanisms
- Test policy server failure scenarios
- Monitor policy server health and availability

**Debugging Policy Decisions**
- Enable detailed logging on policy server
- Review Pexip logs for policy-related events
- Use policy server test interfaces where available
- Implement proper error reporting and diagnostics