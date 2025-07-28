---
page_title: "pexip_infinity_conference_alias Resource - terraform-provider-pexip"
subcategory: ""
description: |-
  Manages a Pexip Infinity conference alias configuration.
---

# pexip_infinity_conference_alias (Resource)

Manages a Pexip Infinity conference alias configuration. Conference aliases provide alternative names or addresses that participants can use to join conferences. They enable flexible routing and user-friendly access to conferences through memorable names, phone numbers, or domain-specific identifiers. Aliases are particularly useful for creating multiple entry points to the same conference or for providing backward compatibility during conference migrations.

## Example Usage

### Basic Conference Alias

```terraform
resource "pexip_infinity_conference_alias" "meeting_room_alias" {
  alias       = "boardroom"
  description = "Alias for main boardroom conference"
  conference  = data.pexip_infinity_conference.main_boardroom.id
}
```

### Phone Number Alias

```terraform
resource "pexip_infinity_conference_alias" "phone_number" {
  alias       = "8005551234"
  description = "Toll-free number for customer support meetings"
  conference  = data.pexip_infinity_conference.customer_support.id
}
```

### Multiple Aliases for One Conference

```terraform
# Primary conference
data "pexip_infinity_conference" "weekly_standup" {
  name = "Weekly Team Standup"
}

# Short alias for easy access
resource "pexip_infinity_conference_alias" "standup_short" {
  alias       = "standup"
  description = "Short alias for weekly standup"
  conference  = data.pexip_infinity_conference.weekly_standup.id
}

# Department-specific alias
resource "pexip_infinity_conference_alias" "dev_standup" {
  alias       = "dev.standup"
  description = "Development team standup alias"
  conference  = data.pexip_infinity_conference.weekly_standup.id
}

# Numeric alias for phone users
resource "pexip_infinity_conference_alias" "standup_number" {
  alias       = "12345"
  description = "Numeric alias for phone access"
  conference  = data.pexip_infinity_conference.weekly_standup.id
}
```

### Department Conference Aliases

```terraform
variable "departments" {
  type = list(object({
    name        = string
    alias       = string
    conference  = string
    description = string
  }))
  default = [
    {
      name        = "Engineering"
      alias       = "engineering"
      conference  = "/api/admin/configuration/v1/conference/42/"
      description = "Engineering team conference"
    },
    {
      name        = "Marketing"
      alias       = "marketing"
      conference  = "/api/admin/configuration/v1/conference/43/"
      description = "Marketing team conference"
    },
    {
      name        = "Sales"
      alias       = "sales"
      conference  = "/api/admin/configuration/v1/conference/44/"
      description = "Sales team conference"
    }
  ]
}

resource "pexip_infinity_conference_alias" "department_aliases" {
  count       = length(var.departments)
  alias       = var.departments[count.index].alias
  description = var.departments[count.index].description
  conference  = var.departments[count.index].conference
}
```

### Geographic Conference Aliases

```terraform
# US East Coast conference room
resource "pexip_infinity_conference_alias" "us_east_room" {
  alias       = "us-east.conference"
  description = "US East Coast conference room"
  conference  = data.pexip_infinity_conference.us_east.id
}

# US West Coast conference room
resource "pexip_infinity_conference_alias" "us_west_room" {
  alias       = "us-west.conference"
  description = "US West Coast conference room"
  conference  = data.pexip_infinity_conference.us_west.id
}

# European conference room
resource "pexip_infinity_conference_alias" "europe_room" {
  alias       = "europe.conference"
  description = "European conference room"
  conference  = data.pexip_infinity_conference.europe.id
}

# Asia Pacific conference room
resource "pexip_infinity_conference_alias" "apac_room" {
  alias       = "apac.conference"
  description = "Asia Pacific conference room"
  conference  = data.pexip_infinity_conference.apac.id
}
```

### Customer-Facing Conference Aliases

```terraform
# Customer support conference
resource "pexip_infinity_conference_alias" "support_main" {
  alias       = "support"
  description = "Main customer support conference"
  conference  = data.pexip_infinity_conference.customer_support.id
}

# Sales demo conference
resource "pexip_infinity_conference_alias" "sales_demo" {
  alias       = "demo"
  description = "Sales demonstration conference"
  conference  = data.pexip_infinity_conference.sales_demo.id
}

# Training sessions
resource "pexip_infinity_conference_alias" "training" {
  alias       = "training"
  description = "Customer training sessions"
  conference  = data.pexip_infinity_conference.training_room.id
}

# Webinar alias
resource "pexip_infinity_conference_alias" "webinar" {
  alias       = "webinar"
  description = "Monthly product webinar"
  conference  = data.pexip_infinity_conference.monthly_webinar.id
}
```

### Legacy System Integration

```terraform
# Legacy phone system integration
resource "pexip_infinity_conference_alias" "legacy_bridge" {
  alias       = "9001"
  description = "Bridge number for legacy phone system"
  conference  = data.pexip_infinity_conference.bridge_conference.id
}

# H.323 system alias
resource "pexip_infinity_conference_alias" "h323_room" {
  alias       = "h323.room1"
  description = "H.323 compatible room alias"
  conference  = data.pexip_infinity_conference.video_room.id
}

# SIP endpoint alias
resource "pexip_infinity_conference_alias" "sip_endpoint" {
  alias       = "sip.meeting"
  description = "SIP endpoint for enterprise calls"
  conference  = data.pexip_infinity_conference.enterprise_meeting.id
}
```

### Event and Meeting Aliases

```terraform
# All-hands meeting
resource "pexip_infinity_conference_alias" "all_hands" {
  alias       = "all-hands"
  description = "Monthly all-hands meeting"
  conference  = data.pexip_infinity_conference.all_hands_meeting.id
}

# Board meeting
resource "pexip_infinity_conference_alias" "board_meeting" {
  alias       = "board"
  description = "Executive board meeting"
  conference  = data.pexip_infinity_conference.executive_board.id
}

# Customer event
resource "pexip_infinity_conference_alias" "customer_event" {
  alias       = "customer-event-2024"
  description = "Annual customer event 2024"
  conference  = data.pexip_infinity_conference.annual_event.id
}
```

### Temporary and Project Aliases

```terraform
# Project-specific conference
resource "pexip_infinity_conference_alias" "project_alpha" {
  alias       = "project-alpha"
  description = "Project Alpha team meetings"
  conference  = data.pexip_infinity_conference.project_room.id
}

# Temporary event alias
resource "pexip_infinity_conference_alias" "temp_event" {
  alias       = "temp-${formatdate("YYYY-MM-DD", timestamp())}"
  description = "Temporary conference for today's event"
  conference  = data.pexip_infinity_conference.temp_conference.id
}

# Contractor access
resource "pexip_infinity_conference_alias" "contractor" {
  alias       = "contractor-access"
  description = "External contractor conference access"
  conference  = data.pexip_infinity_conference.contractor_room.id
}
```

### Multi-Domain Aliases

```terraform
# Corporate domain aliases
resource "pexip_infinity_conference_alias" "corp_main" {
  alias       = "corp.main"
  description = "Corporate main conference"
  conference  = data.pexip_infinity_conference.corporate_main.id
}

# Subsidiary domain aliases
resource "pexip_infinity_conference_alias" "subsidiary_a" {
  alias       = "suba.meeting"
  description = "Subsidiary A conference"
  conference  = data.pexip_infinity_conference.subsidiary_a.id
}

# Partner domain aliases
resource "pexip_infinity_conference_alias" "partner_collab" {
  alias       = "partner.collaboration"
  description = "Partner collaboration conference"
  conference  = data.pexip_infinity_conference.partner_room.id
}
```

## Schema

### Required

- `alias` (String) - The unique alias for the conference. Maximum length: 250 characters.
- `conference` (String) - Reference to the conference resource URI that this alias points to.

### Optional

- `description` (String) - A description of the conference alias. Maximum length: 250 characters.

### Read-Only

- `id` (String) - Resource URI for the conference alias in Infinity.
- `resource_id` (Number) - The resource integer identifier for the conference alias in Infinity.

## Import

Import is supported using the following syntax:

```shell
terraform import pexip_infinity_conference_alias.example 123
```

Where `123` is the numeric resource ID of the conference alias.

## Usage Notes

### Alias Naming Conventions

- **Alphanumeric Characters**: Use letters, numbers, and common separators (hyphens, periods, underscores)
- **Case Sensitivity**: Aliases are typically case-insensitive but maintain consistency
- **Length Limits**: Keep aliases reasonably short for ease of use
- **Uniqueness**: Each alias must be unique across the entire Pexip Infinity deployment

### Conference References

- **URI Format**: Conference references use the full API URI format
- **Resource IDs**: References point to existing conference resources
- **Data Sources**: Use data sources to reference existing conferences
- **Validation**: Ensure referenced conferences exist before creating aliases

### Use Cases

- **User-Friendly Names**: Provide memorable names for frequently used conferences
- **Phone Access**: Create numeric aliases for traditional phone system integration
- **Department Access**: Organize aliases by department, team, or function
- **Geographic Distribution**: Create location-specific aliases for global organizations
- **Legacy Integration**: Support existing naming conventions during system migrations

### Routing and Load Balancing

- **Automatic Routing**: Aliases inherit routing behavior from the target conference
- **Load Distribution**: Multiple aliases to the same conference share the load
- **Geographic Routing**: Combine with system locations for geographic optimization
- **Failover**: Aliases provide redundant access paths to conferences

### Security Considerations

- **Access Control**: Aliases inherit security settings from the target conference
- **Enumeration**: Be mindful of alias patterns that might enable conference enumeration
- **Temporary Access**: Use temporary aliases for time-limited access requirements
- **Audit Logging**: Monitor alias usage for security analysis

### Integration Patterns

- **SIP Integration**: Use domain-style aliases for SIP endpoint integration
- **H.323 Integration**: Support legacy H.323 naming conventions
- **Phone Systems**: Provide numeric aliases for PSTN and PBX integration
- **Web Applications**: Create user-friendly aliases for web-based access

### Performance Considerations

- **Alias Resolution**: Large numbers of aliases may impact resolution performance
- **Caching**: Pexip Infinity caches alias resolutions for performance
- **Distribution**: Distribute aliases across different conferences for load balancing
- **Monitoring**: Monitor alias usage patterns for optimization opportunities

## Troubleshooting

### Common Issues

**Conference Alias Creation Fails**
- Verify the alias is unique across the entire Pexip Infinity deployment
- Check that the conference URI exists and is accessible
- Ensure the alias follows valid naming conventions
- Verify the description doesn't exceed the maximum length

**Alias Resolution Failures**
- Confirm the alias exists and is properly configured
- Verify the target conference is active and accessible
- Check for typos in the alias name
- Ensure proper case sensitivity handling

**Conference Not Found Through Alias**
- Verify the conference URI reference is correct and current
- Check that the target conference hasn't been deleted or moved
- Ensure the conference is in an active state
- Verify permissions to access the target conference

**Duplicate Alias Errors**
- Check for existing aliases with the same name
- Verify uniqueness across all conference types (permanent, scheduled, etc.)
- Consider using namespacing or prefixes to avoid conflicts
- Review alias naming conventions and standards

**Phone System Integration Issues**
- Verify numeric aliases are properly formatted for phone systems
- Check that phone system routing is configured correctly
- Ensure proper integration between Pexip and telephony infrastructure
- Test dial-in functionality from various phone types

**SIP Integration Problems**
- Verify SIP domain configuration matches alias expectations
- Check SIP proxy settings and routing rules
- Ensure proper DNS resolution for domain-based aliases
- Monitor SIP registration and routing logs

**H.323 Integration Issues**
- Verify H.323 gatekeeper configuration supports the alias format
- Check H.323 zone and prefix configurations
- Ensure proper registration and routing through gatekeepers
- Test H.323 endpoint connectivity using aliases

**Performance Issues**
- Monitor alias resolution times and frequency
- Check for excessive alias creation or deletion patterns
- Verify system resources aren't constrained
- Consider alias caching and optimization strategies

**Update and Migration Problems**
- Plan alias changes carefully to avoid service disruption
- Use temporary aliases during migration periods
- Communicate alias changes to users in advance
- Test new aliases before removing old ones

**Import Issues**
- Use the numeric resource ID, not the alias name
- Verify the conference alias exists in the Infinity cluster
- Check provider authentication credentials have access to the resource
- Confirm the alias configuration is accessible

**Security and Access Control Issues**
- Verify that alias access matches intended security policies
- Check that aliases don't inadvertently expose restricted conferences
- Monitor for unauthorized alias usage or enumeration attempts
- Ensure proper audit logging for alias access

**Legacy System Compatibility**
- Test alias formats with legacy video conferencing systems
- Verify compatibility with existing phone system configurations
- Check for character encoding issues with international systems
- Ensure proper protocol translation for different endpoint types

**Monitoring and Alerting Issues**
- Set up monitoring for alias resolution failures
- Track alias usage patterns for capacity planning
- Monitor for suspicious alias access patterns
- Implement alerting for alias configuration changes