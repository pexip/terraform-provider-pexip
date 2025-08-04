---
page_title: "pexip_infinity_global_configuration Resource - terraform-provider-pexip"
subcategory: ""
description: |-
  Manages the Pexip Infinity global system configuration.
---

# pexip_infinity_global_configuration (Resource)

Manages the Pexip Infinity global system configuration. This is a singleton resource that configures system-wide settings including protocol support, security parameters, cloud bursting, port ranges, timeouts, and administrative settings. Only one global configuration exists per Pexip Infinity cluster.

## Example Usage

### Basic Global Configuration

```terraform
resource "pexip_infinity_global_configuration" "basic" {
  enable_webrtc = true
  enable_sip    = true
  enable_h323   = false
  enable_rtmp   = true
  crypto_mode   = "besteffort"
}
```

### Enterprise Security Configuration

```terraform
resource "pexip_infinity_global_configuration" "enterprise" {
  enable_webrtc    = true
  enable_sip       = true
  enable_h323      = true
  enable_rtmp      = true
  crypto_mode      = "required"
  
  # Network configuration
  media_ports_start       = 40000
  media_ports_end         = 49999
  signalling_ports_start  = 5060
  signalling_ports_end    = 5080
  
  # Conference settings
  conference_create_permissions = "user_admin"
  conference_creation_mode      = "per_cluster"
  guests_only_timeout          = 30
  waiting_for_chair_timeout    = 10
  
  # Security and compliance
  enable_analytics      = true
  enable_error_reporting = false
  bandwidth_restrictions = "restricted"
  
  administrator_email = "admin@company.com"
}
```

### Cloud Bursting with AWS

```terraform
resource "pexip_infinity_global_configuration" "aws_bursting" {
  enable_webrtc    = true
  enable_sip       = true
  crypto_mode      = "besteffort"
  
  # Cloud bursting configuration
  bursting_enabled = true
  cloud_provider   = "aws"
  aws_access_key   = var.aws_access_key
  aws_secret_key   = var.aws_secret_key
  
  # Conference settings
  conference_create_permissions = "any_authenticated"
  conference_creation_mode      = "per_cluster"
  
  administrator_email = "pexip-admin@company.com"
}
```

### Azure Cloud Integration

```terraform
resource "pexip_infinity_global_configuration" "azure_integration" {
  enable_webrtc    = true
  enable_sip       = true
  enable_h323      = false
  enable_rtmp      = true
  crypto_mode      = "required"
  
  # Azure cloud bursting
  bursting_enabled = true
  cloud_provider   = "azure"
  azure_client_id  = var.azure_client_id
  azure_secret     = var.azure_secret
  
  # Media configuration
  max_pixels_per_second = "1920x1080x30"
  media_ports_start     = 36000
  media_ports_end       = 59999
  
  # Administrative settings
  enable_analytics       = true
  enable_error_reporting = true
  administrator_email    = "video-admin@company.com"
}
```

### High-Security Configuration

```terraform
resource "pexip_infinity_global_configuration" "high_security" {
  enable_webrtc    = true
  enable_sip       = true
  enable_h323      = false
  enable_rtmp      = false
  crypto_mode      = "required"
  
  # Restricted conference creation
  conference_create_permissions = "admin_only"
  conference_creation_mode      = "per_cluster"
  
  # Short timeouts for security
  guests_only_timeout       = 15
  waiting_for_chair_timeout = 5
  
  # Network restrictions
  bandwidth_restrictions = "restricted"
  media_ports_start      = 50000
  media_ports_end        = 50999
  signalling_ports_start = 5060
  signalling_ports_end   = 5061
  
  # Security monitoring
  enable_analytics       = true
  enable_error_reporting = false
  
  administrator_email = "security-admin@company.com"
  
  # Restrict conference creation to specific groups
  global_conference_create_groups = [
    "CN=Video Admins,OU=Groups,DC=company,DC=com",
    "CN=IT Managers,OU=Groups,DC=company,DC=com"
  ]
}
```

### Development Environment Configuration

```terraform
resource "pexip_infinity_global_configuration" "development" {
  enable_webrtc    = true
  enable_sip       = true
  enable_h323      = true
  enable_rtmp      = true
  crypto_mode      = "besteffort"
  
  # Relaxed conference creation for development
  conference_create_permissions = "any_authenticated"
  conference_creation_mode      = "per_node"
  
  # Extended timeouts for testing
  guests_only_timeout       = 120
  waiting_for_chair_timeout = 30
  
  # Analytics for development insights
  enable_analytics       = true
  enable_error_reporting = true
  
  administrator_email = "dev-team@company.com"
}
```

### Multi-Protocol Configuration

```terraform
resource "pexip_infinity_global_configuration" "multi_protocol" {
  # Enable all protocols
  enable_webrtc = true
  enable_sip    = true
  enable_h323   = true
  enable_rtmp   = true
  
  # Security settings
  crypto_mode = "besteffort"
  
  # Optimized port ranges for different protocols
  media_ports_start      = 32768
  media_ports_end        = 65535
  signalling_ports_start = 1024
  signalling_ports_end   = 65535
  
  # Video quality settings
  max_pixels_per_second = "3840x2160x30"  # 4K support
  
  # Conference management
  conference_create_permissions = "user_admin"
  conference_creation_mode      = "per_cluster"
  
  # Operational settings
  enable_analytics       = true
  enable_error_reporting = true
  bandwidth_restrictions = "none"
  
  administrator_email = "operations@company.com"
}
```

### Variable-Driven Configuration

```terraform
variable "global_config" {
  type = object({
    enable_webrtc                = bool
    enable_sip                   = bool
    enable_h323                  = bool
    enable_rtmp                  = bool
    crypto_mode                  = string
    conference_create_permissions = string
    conference_creation_mode     = string
    enable_analytics             = bool
    administrator_email          = string
  })
  default = {
    enable_webrtc                = true
    enable_sip                   = true
    enable_h323                  = false
    enable_rtmp                  = true
    crypto_mode                  = "besteffort"
    conference_create_permissions = "user_admin"
    conference_creation_mode     = "per_cluster"
    enable_analytics             = true
    administrator_email          = "admin@company.com"
  }
}

resource "pexip_infinity_global_configuration" "from_variable" {
  enable_webrtc                = var.global_config.enable_webrtc
  enable_sip                   = var.global_config.enable_sip
  enable_h323                  = var.global_config.enable_h323
  enable_rtmp                  = var.global_config.enable_rtmp
  crypto_mode                  = var.global_config.crypto_mode
  conference_create_permissions = var.global_config.conference_create_permissions
  conference_creation_mode     = var.global_config.conference_creation_mode
  enable_analytics             = var.global_config.enable_analytics
  administrator_email          = var.global_config.administrator_email
}
```

## Schema

### Optional

- `enable_webrtc` (Boolean) - Whether to enable WebRTC protocol support.
- `enable_sip` (Boolean) - Whether to enable SIP protocol support.
- `enable_h323` (Boolean) - Whether to enable H.323 protocol support.
- `enable_rtmp` (Boolean) - Whether to enable RTMP protocol support.
- `crypto_mode` (String) - Cryptographic mode for conferences. Valid values: `disabled`, `besteffort`, `required`.
- `max_pixels_per_second` (String) - Maximum pixels per second for video transmission (e.g., "1920x1080x30").
- `media_ports_start` (Number) - Starting port for media traffic. Valid range: 1024-65535.
- `media_ports_end` (Number) - Ending port for media traffic. Valid range: 1024-65535.
- `signalling_ports_start` (Number) - Starting port for signalling traffic. Valid range: 1024-65535.
- `signalling_ports_end` (Number) - Ending port for signalling traffic. Valid range: 1024-65535.
- `bursting_enabled` (Boolean) - Whether to enable cloud bursting functionality.
- `cloud_provider` (String) - Cloud provider for bursting. Valid values: `aws`, `azure`, `google`.
- `aws_access_key` (String, Sensitive) - AWS access key for cloud bursting.
- `aws_secret_key` (String, Sensitive) - AWS secret key for cloud bursting.
- `azure_client_id` (String, Sensitive) - Azure client ID for cloud bursting.
- `azure_secret` (String, Sensitive) - Azure secret for cloud bursting.
- `guests_only_timeout` (Number) - Timeout in minutes for guests-only conferences. Valid range: 0-1440.
- `waiting_for_chair_timeout` (Number) - Timeout in minutes when waiting for chair to join. Valid range: 0-1440.
- `conference_create_permissions` (String) - Who can create conferences. Valid values: `none`, `admin_only`, `user_admin`, `any_authenticated`.
- `conference_creation_mode` (String) - Conference creation mode. Valid values: `disabled`, `per_node`, `per_cluster`.
- `enable_analytics` (Boolean) - Whether to enable analytics collection.
- `enable_error_reporting` (Boolean) - Whether to enable automatic error reporting.
- `bandwidth_restrictions` (String) - Bandwidth restriction mode. Valid values: `none`, `restricted`.
- `administrator_email` (String) - Administrator email address for system notifications.
- `global_conference_create_groups` (List of String) - List of groups that can create conferences globally.

### Read-Only

- `id` (String) - Resource URI for the global configuration in Infinity.

## Import

Import is supported using any value as the import ID (since this is a singleton resource):

```shell
terraform import pexip_infinity_global_configuration.example global
```

## Usage Notes

### Singleton Resource

This is a singleton resource, meaning only one global configuration exists per Pexip Infinity cluster. Creating multiple instances of this resource in your Terraform configuration will result in conflicts.

### Protocol Support

- **WebRTC**: Required for browser-based video conferences and web applications
- **SIP**: Required for SIP-based endpoints and telephony integration
- **H.323**: Required for legacy H.323 video conferencing systems
- **RTMP**: Required for streaming and recording functionality

### Crypto Modes

- **disabled**: No encryption requirements (not recommended for production)
- **besteffort**: Encryption preferred but not required
- **required**: Encryption mandatory for all connections

### Port Configuration

- **Media Ports**: Used for RTP/RTCP media traffic (audio/video)
- **Signalling Ports**: Used for SIP and H.323 signalling
- **Range Planning**: Ensure adequate port ranges for expected concurrent calls
- **Firewall Coordination**: Coordinate port ranges with firewall and security teams

### Cloud Bursting

- **AWS Integration**: Requires AWS access key and secret key
- **Azure Integration**: Requires Azure client ID and secret
- **Google Cloud**: Supported for cloud bursting (credentials vary)
- **Scalability**: Enables automatic scaling to cloud resources during peak usage

### Conference Creation Permissions

- **none**: Conferences cannot be created dynamically
- **admin_only**: Only administrators can create conferences
- **user_admin**: Administrators and designated users can create conferences
- **any_authenticated**: Any authenticated user can create conferences

### Timeout Configuration

- **Guests Only**: Automatically ends conferences with only guest participants
- **Waiting for Chair**: Ends waiting rooms when no chair joins
- **Zero Values**: Setting to 0 disables the timeout

### Security Considerations

- **Encryption**: Use "required" crypto mode for sensitive environments
- **Access Control**: Restrict conference creation based on organizational needs
- **Monitoring**: Enable analytics for security monitoring and compliance
- **Credentials**: Store cloud credentials securely using Terraform variables

## Troubleshooting

### Common Issues

**Global Configuration Update Fails**
- Verify all enum values are from the valid options list
- Check that port ranges are valid and start <= end
- Ensure timeout values are within the valid range (0-1440)
- Verify cloud credentials are correctly formatted

**Protocol Support Issues**
- Ensure required protocols are enabled for your use case
- Check licensing requirements for certain protocols
- Verify client compatibility with enabled protocols
- Monitor protocol-specific logs for connection issues

**Port Range Conflicts**
- Ensure media and signalling port ranges don't overlap inappropriately
- Verify port ranges don't conflict with other system services
- Check firewall rules allow traffic on configured port ranges
- Monitor port utilization to ensure adequate range size

**Cloud Bursting Authentication Fails**
- Verify cloud provider credentials are correct and current
- Ensure cloud provider account has necessary permissions
- Check network connectivity to cloud provider APIs
- Monitor cloud provider service status for outages

**Conference Creation Issues**
- Verify conference_create_permissions setting allows intended users
- Check that global_conference_create_groups contains correct group DNs
- Ensure users have appropriate authentication and authorization
- Monitor conference creation logs for detailed error information

**Timeout Configuration Problems**
- Verify timeout values are appropriate for your use case
- Check that timeout settings don't conflict with user expectations
- Monitor conference duration to optimize timeout values
- Ensure timeout notifications are properly configured

**Analytics and Reporting Issues**
- Verify analytics collection is enabled if reporting is required
- Check that error reporting settings meet compliance requirements
- Ensure adequate storage for analytics data
- Monitor analytics pipeline for data collection issues

**Performance Degradation**
- Monitor system performance after configuration changes
- Check media quality settings against available bandwidth
- Verify port range sizes are adequate for concurrent users
- Monitor cloud bursting utilization and costs

**Import Issues**
- Use any string value as the import ID for this singleton resource
- Verify provider credentials have access to global configuration
- Check that the Pexip Infinity cluster is accessible
- Note that sensitive cloud credentials may not be imported

**Encryption Configuration Problems**
- Verify crypto_mode setting is compatible with all endpoints
- Check that "required" mode doesn't block legitimate connections
- Monitor encryption negotiation logs for compatibility issues
- Ensure proper certificate configuration for encrypted connections

**Email Notification Failures**
- Verify administrator_email address is valid and accessible
- Check SMTP configuration for email delivery
- Ensure email notifications are not blocked by spam filters
- Monitor email delivery logs for notification issues

**Bandwidth Restriction Issues**
- Verify bandwidth_restrictions setting matches network capacity
- Check that "restricted" mode doesn't overly limit call quality
- Monitor bandwidth utilization and quality metrics
- Adjust restrictions based on network performance data