---
page_title: "pexip_infinity_end_user Resource - terraform-provider-pexip"
subcategory: ""
description: |-
  Manages a Pexip Infinity end user account configuration.
---

# pexip_infinity_end_user (Resource)

Manages an end user account with the Infinity service. End users represent individuals who can participate in conferences, create meetings, and access Pexip Infinity services. This resource allows you to manage user profiles, contact information, and group memberships for integration with directory services and user management systems.

## Example Usage

### Basic End User

```terraform
resource "pexip_infinity_end_user" "john_doe" {
  primary_email_address = "john.doe@company.com"
}
```

### Complete User Profile

```terraform
resource "pexip_infinity_end_user" "jane_smith" {
  primary_email_address = "jane.smith@company.com"
  first_name            = "Jane"
  last_name             = "Smith"
  display_name          = "Jane Smith"
  telephone_number      = "+1-555-0123"
  mobile_number         = "+1-555-0124"
  title                 = "Senior Engineer"
  department            = "Engineering"
  avatar_url            = "https://avatars.company.com/jane.smith.jpg"
}
```

### User with Group Memberships

```terraform
# Assume user groups are defined elsewhere
resource "pexip_infinity_end_user" "manager_user" {
  primary_email_address = "manager@company.com"
  first_name            = "John"
  last_name             = "Manager"
  display_name          = "John Manager"
  title                 = "Engineering Manager"
  department            = "Engineering"
  user_groups          = [
    "/api/admin/configuration/v1/user_group/1/",  # Managers group
    "/api/admin/configuration/v1/user_group/3/",  # Engineering group
  ]
}
```

### Exchange Integration User

```terraform
resource "pexip_infinity_end_user" "exchange_user" {
  primary_email_address = "exchange.user@company.com"
  first_name            = "Exchange"
  last_name             = "User"
  display_name          = "Exchange User"
  department            = "IT"
  ms_exchange_guid      = "a1b2c3d4-e5f6-7890-abcd-123456789abc"
  sync_tag              = "exchange-sync"
}
```

### Bulk User Creation

```terraform
# Create multiple users from variable
resource "pexip_infinity_end_user" "department_users" {
  for_each = var.department_users
  
  primary_email_address = each.value.email
  first_name            = each.value.first_name
  last_name             = each.value.last_name
  display_name          = "${each.value.first_name} ${each.value.last_name}"
  telephone_number      = each.value.phone
  department            = each.value.department
  title                 = each.value.title
  sync_tag              = "ldap-import"
}
```

### Users with Avatar and Contact Information

```terraform
resource "pexip_infinity_end_user" "executive_user" {
  primary_email_address = "ceo@company.com"
  first_name            = "Chief"
  last_name             = "Executive"
  display_name          = "CEO"
  title                 = "Chief Executive Officer"
  department            = "Executive"
  telephone_number      = "+1-555-0100"
  mobile_number         = "+1-555-0101"
  avatar_url            = "https://cdn.company.com/avatars/ceo.jpg"
  user_groups          = [
    "/api/admin/configuration/v1/user_group/1/",  # Executives
    "/api/admin/configuration/v1/user_group/2/",  # Board Members
  ]
}
```

## Schema

### Required

- `primary_email_address` (String) - The unique primary email address for the end user. Maximum length: 100 characters.

### Optional

- `first_name` (String) - The first name of the end user. Maximum length: 250 characters.
- `last_name` (String) - The last name of the end user. Maximum length: 250 characters.
- `display_name` (String) - The display name of the end user. Maximum length: 250 characters.
- `telephone_number` (String) - The telephone number of the end user. Maximum length: 100 characters.
- `mobile_number` (String) - The mobile number of the end user. Maximum length: 100 characters.
- `title` (String) - The job title of the end user. Maximum length: 128 characters.
- `department` (String) - The department of the end user. Maximum length: 100 characters.
- `avatar_url` (String) - The avatar URL for the end user. Maximum length: 255 characters.
- `user_groups` (List of String) - List of user group resource URIs that this user belongs to.
- `ms_exchange_guid` (String) - Exchange Mailbox ID. Maximum length: 100 characters.
- `sync_tag` (String) - LDAP sync identifier. Maximum length: 250 characters.

### Read-Only

- `id` (String) - Resource URI for the end user in Infinity.
- `resource_id` (Number) - The resource integer identifier for the end user in Infinity.
- `user_oid` (String) - Microsoft 365 Object ID (read-only).
- `exchange_user_id` (String) - Exchange User ID (read-only).

## Import

Import is supported using the following syntax:

```shell
terraform import pexip_infinity_end_user.example 123
```

Where `123` is the numeric resource ID of the end user.

## Usage Notes

### User Identity Management
- Primary email address must be unique across all end users
- Display name is used in conferencing interfaces and notifications
- Contact information enables communication and calendar integration
- User groups define permissions and access levels

### Directory Service Integration
- Use `sync_tag` for LDAP/AD synchronization identification
- `ms_exchange_guid` enables Exchange calendar integration
- `user_oid` and `exchange_user_id` are automatically populated for Microsoft 365 users
- External sync can update user information automatically

### Avatar Configuration
- Avatar URLs should be publicly accessible or use internal CDN
- Recommended image size is 128x128 pixels
- Supported formats typically include JPG, PNG, and GIF
- Consider caching and load balancing for avatar services

### User Group Memberships
- User groups control access to conferences and features
- Groups are referenced by their resource URI
- Multiple group memberships enable fine-grained permissions
- Group changes immediately affect user access

### Contact Information Best Practices
- Use standardized phone number formats
- Ensure email addresses are valid and monitored
- Keep contact information current for emergency situations
- Consider privacy requirements for personal information

## Troubleshooting

### Common Issues

**End User Creation Fails**
- Verify the primary email address is unique
- Ensure email address format is valid
- Check that all field lengths are within specified limits
- Verify required fields are provided

**Email Address Conflicts**
- Check for existing users with the same email address
- Verify email address is not used by devices or other resources
- Ensure proper email domain configuration
- Use consistent email address formatting

**User Group Assignment Issues**
- Verify user group URIs are correct and exist
- Check that user groups are properly configured
- Ensure user has appropriate permissions for group membership
- Verify group membership limits are not exceeded

**Directory Synchronization Problems**
- Verify sync_tag matches external directory attribute
- Check LDAP/AD connection and authentication
- Ensure proper attribute mapping configuration
- Monitor synchronization logs for errors

**Exchange Integration Issues**
- Verify ms_exchange_guid format is correct
- Check Exchange server connectivity and permissions
- Ensure proper Exchange Web Services configuration
- Verify calendar integration settings

**Avatar Display Problems**
- Verify avatar URL is accessible from clients
- Check image format and size requirements
- Ensure proper CORS headers for cross-origin requests
- Test avatar URL accessibility from different networks

**Authentication Issues**
- Verify user exists in authentication source (local/LDAP/SSO)
- Check user account status and permissions
- Ensure proper authentication configuration
- Verify user credentials and account lockout status

**Import Fails**
- Ensure you're using the numeric resource ID, not the email
- Verify the end user exists in the Infinity cluster
- Check provider authentication credentials have access to the resource

**Contact Information Formatting**
- Use consistent phone number formats across all users
- Ensure telephone numbers include proper country codes
- Validate email address formats before creation
- Consider international number formatting standards

**User Profile Synchronization**
- Monitor external directory changes for user updates
- Implement proper conflict resolution for synchronized data
- Ensure proper mapping between external and Pexip attributes
- Regularly audit user information for accuracy

**Performance with Large User Sets**
- Consider batch operations for bulk user creation
- Implement proper indexing for user searches
- Monitor system performance with large user databases
- Use pagination for user management interfaces

**Compliance and Data Protection**
- Ensure user data handling complies with privacy regulations
- Implement proper data retention policies
- Consider user consent for profile information
- Audit user data access and modifications

**Microsoft 365 Integration**
- Verify user_oid is populated for Office 365 users
- Check Azure AD synchronization status
- Ensure proper permissions for Exchange integration
- Monitor Microsoft Graph API integration