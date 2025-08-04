---
page_title: "pexip_infinity_conference Resource - terraform-provider-pexip"
subcategory: ""
description: |-
  Manages a Pexip Infinity conference configuration.
---

# pexip_infinity_conference (Resource)

Manages a Pexip Infinity conference configuration. This resource creates and configures virtual meeting rooms (VMRs) with specific settings for participants, security, and media handling.

## Example Usage

### Basic Conference

```terraform
resource "pexip_infinity_conference" "meeting_room" {
  name         = "executive-meeting"
  service_type = "conference"
  pin          = "1234567"
}
```

### Lecture Mode Conference

```terraform
resource "pexip_infinity_conference" "webinar" {
  name           = "company-webinar"
  service_type   = "lecture"
  description    = "Monthly company-wide webinar"
  pin            = "987654321"
  allow_guests   = true
  guests_muted   = true
  hosts_can_unmute = true
}
```

### Full Configuration

```terraform
resource "pexip_infinity_conference" "secure_meeting" {
  name        = "board-meeting"
  description = "Secure board meeting room"
  service_type = "conference"
  
  # Security configuration
  pin                   = var.meeting_pin
  guest_pin            = var.guest_pin
  allow_guests         = true
  guests_muted         = false
  hosts_can_unmute     = true
  
  # Media configuration
  max_pixels_per_second = 1920000
  
  # Tracking
  tag = "executive-meetings"
}
```

### Test Call Service

```terraform
resource "pexip_infinity_conference" "test_service" {
  name         = "test-call"
  service_type = "test_call"
  description  = "Service for testing audio and video"
}
```

## Schema

### Required

- `name` (String) - The unique name used to refer to this conference. Maximum length: 250 characters.
- `service_type` (String) - The type of conferencing service. Valid choices: `conference`, `lecture`, `two_stage_dialing`, `test_call`, `media_playback`.

### Optional

- `allow_guests` (Boolean) - Whether Guest participants are allowed to join. Defaults to `false`.
- `description` (String) - A description of the conference. Maximum length: 250 characters.
- `guest_pin` (String, Sensitive) - Optional secure access code for Guest participants. Length: 4-20 digits.
- `guests_muted` (Boolean) - Whether Guest participants are muted by default. Defaults to `false`.
- `hosts_can_unmute` (Boolean) - Whether Host participants can unmute Guest participants. Defaults to `false`.
- `max_pixels_per_second` (Number) - Maximum pixels per second for video quality.
- `pin` (String, Sensitive) - Secure access code for participants. Length: 4-20 digits, including any terminal #.
- `tag` (String) - A unique identifier used to track usage. Maximum length: 250 characters.

### Read-Only

- `id` (String) - Resource URI for the conference in Infinity.
- `resource_id` (Number) - The resource integer identifier for the conference in Infinity.

## Service Types

### conference
Standard multi-party conference where all participants can see and hear each other.

### lecture
One-to-many presentation mode where hosts present to many participants. Typically used for webinars and large meetings.

### two_stage_dialing
Allows participants to dial out to external numbers from within the conference.

### test_call
Service for testing audio and video quality before joining actual meetings.

### media_playback
Service for playing back recorded media content.

## Import

Import is supported using the following syntax:

```shell
terraform import pexip_infinity_conference.example 123
```

Where `123` is the numeric resource ID of the conference.

## Usage Notes

### Security Configuration
- Always use strong PINs with sufficient length and complexity
- Consider separate PINs for hosts and guests when needed
- Guest access should be carefully controlled based on meeting requirements

### Audio Management
- Use `guests_muted` for large meetings to reduce background noise
- Enable `hosts_can_unmute` to allow presenters to manage audio
- Consider service type when configuring audio settings

### Video Quality
- Set `max_pixels_per_second` based on network capacity and meeting requirements
- Higher values provide better quality but require more bandwidth
- Consider participant device capabilities when setting limits

### Usage Tracking
- Use the `tag` field for reporting and analytics
- Tags help organize conferences by department, project, or purpose
- Consistent tagging enables better usage insights

## Troubleshooting

### Common Issues

**Conference Creation Fails**
- Verify the conference name is unique within the deployment
- Check that PIN requirements are met (4-20 digits)
- Ensure service_type is valid

**Participants Cannot Join**
- Verify PIN is correct and properly shared
- Check that allow_guests is enabled if guests need access
- Ensure the conference is properly configured and active

**Audio/Video Issues**
- Review max_pixels_per_second settings for video quality
- Check guests_muted and hosts_can_unmute settings
- Verify service_type matches intended usage pattern

**Import Fails**
- Ensure you're using the numeric resource ID, not the name
- Verify the conference exists in the Infinity deployment
- Check provider authentication credentials

### PIN Security Best Practices
- Use PINs of at least 6-8 digits for better security
- Avoid sequential or repeated numbers
- Rotate PINs regularly for sensitive meetings
- Use different PINs for hosts and guests when appropriate