---
page_title: "pexip_infinity_device Resource - terraform-provider-pexip"
subcategory: ""
description: |-
  Manages a Pexip Infinity device configuration.
---

# pexip_infinity_device (Resource)

Manages a device configuration with the Infinity service. Devices represent endpoints that can connect to Pexip Infinity conferences, including video conferencing systems, SIP phones, and software clients. Device resources allow you to pre-configure authentication and connectivity settings for known endpoints.

## Example Usage

### Basic Device Configuration

```terraform
resource "pexip_infinity_device" "conference_room_a" {
  alias = "conference-room-a.company.com"
}
```

### SIP Device with Authentication

```terraform
resource "pexip_infinity_device" "sip_phone" {
  alias                = "sip-phone-101"
  description          = "Reception desk SIP phone"
  username             = "phone101"
  password             = var.sip_phone_password
  primary_owner_email_address = "reception@company.com"
  enable_sip           = true
  enable_h323          = false
}
```

### Video Conferencing Room System

```terraform
resource "pexip_infinity_device" "boardroom_system" {
  alias                = "boardroom.company.com"
  description          = "Main boardroom video system"
  username             = "boardroom"
  password             = var.room_system_password
  primary_owner_email_address = "facilities@company.com"
  enable_sip           = true
  enable_h323          = true
  tag                  = "video-systems"
}
```

### Infinity Connect Device with SSO

```terraform
resource "pexip_infinity_device" "infinity_connect_device" {
  alias                           = "connect-device-01"
  description                     = "Infinity Connect device with SSO"
  primary_owner_email_address     = "user@company.com"
  enable_infinity_connect_sso     = true
  enable_standard_sso             = true
  sso_identity_provider_group     = "corporate-users"
  tag                             = "infinity-connect"
}
```

### Multiple Devices for Department

```terraform
# Sales team devices
resource "pexip_infinity_device" "sales_devices" {
  count = length(var.sales_team_devices)
  
  alias                = var.sales_team_devices[count.index].alias
  description          = "Sales team device ${count.index + 1}"
  username             = var.sales_team_devices[count.index].username
  password             = var.sales_team_devices[count.index].password
  primary_owner_email_address = var.sales_team_devices[count.index].email
  enable_sip           = true
  enable_h323          = false
  tag                  = "sales-team"
  sync_tag             = "ldap-sales"
}
```

### Device with Full Configuration

```terraform
resource "pexip_infinity_device" "executive_room" {
  alias                           = "executive-suite.company.com"
  description                     = "Executive conference room with full capabilities"
  username                        = "executive_room"
  password                        = var.executive_room_password
  primary_owner_email_address     = "executive-assistant@company.com"
  
  # Protocol support
  enable_sip                      = true
  enable_h323                     = true
  enable_infinity_connect_non_sso = true
  enable_infinity_connect_sso     = true
  enable_standard_sso             = true
  
  # SSO configuration
  sso_identity_provider_group     = "executives"
  
  # Organizational tags
  tag                             = "executive-rooms"
  sync_tag                        = "ldap-executives"
}
```

## Schema

### Required

- `alias` (String) - The unique alias name of the device. Maximum length: 250 characters.

### Optional

- `description` (String) - A description of the device. Maximum length: 250 characters.
- `username` (String) - The username for device authentication. Maximum length: 250 characters.
- `password` (String, Sensitive) - The password for device authentication. Maximum length: 100 characters.
- `primary_owner_email_address` (String) - Email address of the device owner. Maximum length: 100 characters.
- `enable_sip` (Boolean) - Whether SIP is enabled for this device. Defaults to false.
- `enable_h323` (Boolean) - Whether H.323 is enabled for this device. Defaults to false.
- `enable_infinity_connect_non_sso` (Boolean) - Whether Infinity Connect without SSO is enabled. Defaults to false.
- `enable_infinity_connect_sso` (Boolean) - Whether Infinity Connect with SSO is enabled. Defaults to false.
- `enable_standard_sso` (Boolean) - Whether standard SSO is enabled. Defaults to false.
- `sso_identity_provider_group` (String) - SSO identity provider group for authentication.
- `tag` (String) - A tag for categorizing the device. Maximum length: 250 characters.
- `sync_tag` (String) - A sync tag for external system integration. Maximum length: 250 characters.

### Read-Only

- `id` (String) - Resource URI for the device in Infinity.
- `resource_id` (Number) - The resource integer identifier for the device in Infinity.

## Import

Import is supported using the following syntax:

```shell
terraform import pexip_infinity_device.example 123
```

Where `123` is the numeric resource ID of the device.

## Usage Notes

### Device Types and Protocols
- **SIP Devices**: IP phones, software clients, room systems supporting SIP
- **H.323 Devices**: Legacy video conferencing systems and terminals
- **Infinity Connect**: Pexip's native client application
- **Mixed Protocol**: Devices supporting multiple protocols for flexibility

### Authentication Methods
- **Username/Password**: Traditional credential-based authentication
- **SSO Integration**: Single sign-on with identity providers
- **Certificate-based**: Using client certificates (configured separately)

### Device Alias Requirements
- Must be unique across the Pexip Infinity deployment
- Can be hostname (device.company.com) or identifier (room-101)
- Used for dialing and device identification
- Should follow consistent naming conventions

### SSO Configuration
- `enable_standard_sso`: General SSO authentication
- `enable_infinity_connect_sso`: SSO specifically for Infinity Connect clients
- `sso_identity_provider_group`: Maps device to specific identity provider groups
- Requires proper SSO configuration in authentication settings

### Infinity Connect Options
- `enable_infinity_connect_non_sso`: Traditional username/password for Infinity Connect
- `enable_infinity_connect_sso`: SSO authentication for Infinity Connect
- Both can be enabled to support different authentication methods

### Device Management
- Use tags for organizational grouping and policy application
- Sync tags enable integration with external directory services
- Primary owner email helps with device management and notifications
- Device descriptions aid in identification and troubleshooting

## Troubleshooting

### Common Issues

**Device Creation Fails**
- Verify the alias is unique across the Infinity deployment
- Ensure alias follows proper naming conventions
- Check that all field lengths are within specified limits
- Verify password meets security requirements if specified

**Device Authentication Issues**
- Verify username and password are correct and not expired
- Check that the appropriate protocol is enabled (SIP/H.323)
- Ensure device supports the enabled authentication methods
- Verify SSO configuration if SSO is enabled

**SIP Registration Problems**
- Verify `enable_sip` is set to true
- Check SIP proxy configuration on the device
- Ensure network connectivity between device and Pexip nodes
- Verify SIP credentials match the device configuration

**H.323 Connection Issues**
- Verify `enable_h323` is set to true
- Check H.323 gatekeeper settings on the device
- Ensure proper network routing for H.323 traffic
- Verify device supports required H.323 protocol versions

**Infinity Connect Problems**
- Verify appropriate Infinity Connect options are enabled
- Check client version compatibility
- Ensure proper DNS resolution for Infinity services
- Verify SSO configuration if SSO is enabled

**SSO Authentication Failures**
- Verify SSO is properly configured in authentication settings
- Check identity provider group configuration
- Ensure user exists in the specified identity provider group
- Verify SSO certificates and trust relationships

**Device Not Found in Calls**
- Verify device alias is correct and unique
- Check that device is properly registered
- Ensure device is not in maintenance or disabled state
- Verify routing rules allow calls to the device

**Import Fails**
- Ensure you're using the numeric resource ID, not the alias
- Verify the device exists in the Infinity cluster
- Check provider authentication credentials have access to the resource

**Protocol Compatibility Issues**
- Verify device supports the enabled protocols
- Check codec compatibility between device and Pexip
- Ensure proper bandwidth and QoS settings
- Verify firewall rules allow required protocol traffic

**Organizational Management Problems**
- Use consistent tagging strategies for device management
- Implement proper sync tag integration with directory services
- Ensure primary owner email addresses are current and valid
- Regularly audit device configurations and ownership

**External System Integration**
- Verify sync tags match external directory attributes
- Check integration credentials and permissions
- Ensure proper mapping between external and Pexip identities
- Monitor synchronization processes for errors or conflicts