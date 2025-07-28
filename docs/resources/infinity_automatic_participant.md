---
page_title: "pexip_infinity_automatic_participant Resource - terraform-provider-pexip"
subcategory: ""
description: |-
  Manages a Pexip Infinity automatic participant configuration.
---

# pexip_infinity_automatic_participant (Resource)

Manages a Pexip Infinity automatic participant configuration. Automatic participants are endpoints that automatically join conferences when they start, providing functionality such as recording, streaming, content sharing, or acting as persistent conference hosts. They can be configured to join with different protocols, roles, and behaviors.

## Example Usage

### Basic Recording Participant

```terraform
resource "pexip_infinity_automatic_participant" "recorder" {
  alias               = "conference-recorder"
  description         = "Automatic recording participant"
  conference          = data.pexip_infinity_conference.meeting_room.id
  protocol            = "sip"
  call_type           = "video"
  role                = "guest"
  keep_conference_alive = "end_conference_when_alone"
  routing             = "auto"
  streaming           = true
  remote_display_name = "Conference Recorder"
}
```

### Chair Participant for Persistent Conferences

```terraform
resource "pexip_infinity_automatic_participant" "chair_bot" {
  alias               = "meeting-chair"
  description         = "Automatic chair to keep conference active"
  conference          = data.pexip_infinity_conference.persistent_meeting.id
  protocol            = "webrtc"
  call_type           = "audio"
  role                = "chair"
  keep_conference_alive = "keep_conference_alive"
  routing             = "auto"
  streaming           = false
  remote_display_name = "Meeting Host"
}
```

### RTMP Streaming Participant

```terraform
resource "pexip_infinity_automatic_participant" "rtmp_streamer" {
  alias               = "live-stream"
  description         = "RTMP streaming to social media"
  conference          = data.pexip_infinity_conference.webinar.id
  protocol            = "rtmp"
  call_type           = "video"
  role                = "guest"
  keep_conference_alive = "end_conference_when_alone"
  routing             = "manual"
  system_location     = data.pexip_infinity_system_location.streaming.id
  streaming           = true
  remote_display_name = "Live Stream"
  presentation_url    = "rtmp://streaming.example.com/live/stream-key"
}
```

### H.323 Gateway Participant

```terraform
resource "pexip_infinity_automatic_participant" "h323_gateway" {
  alias               = "legacy-system"
  description         = "H.323 gateway for legacy video systems"
  conference          = data.pexip_infinity_conference.hybrid_meeting.id
  protocol            = "h323"
  call_type           = "video"
  role                = "guest"
  keep_conference_alive = "end_conference_when_alone"
  routing             = "manual"
  system_location     = data.pexip_infinity_system_location.gateway.id
  dtmf_sequence       = "12345#"
  remote_display_name = "Legacy Video System"
}
```

### Multiple Automatic Participants

```terraform
# Recording participant
resource "pexip_infinity_automatic_participant" "recorder" {
  alias               = "recorder-${var.conference_name}"
  description         = "Recording for ${var.conference_name}"
  conference          = var.conference_id
  protocol            = "sip"
  call_type           = "video"
  role                = "guest"
  keep_conference_alive = "end_conference_when_alone"
  routing             = "auto"
  streaming           = true
  remote_display_name = "Recording Service"
}

# Streaming participant
resource "pexip_infinity_automatic_participant" "streamer" {
  alias               = "streamer-${var.conference_name}"
  description         = "Streaming for ${var.conference_name}"
  conference          = var.conference_id
  protocol            = "rtmp"
  call_type           = "video"
  role                = "guest"
  keep_conference_alive = "end_conference_when_alone"
  routing             = "auto"
  streaming           = true
  remote_display_name = "Live Stream"
  presentation_url    = var.rtmp_stream_url
}

# Chair participant for persistent hosting
resource "pexip_infinity_automatic_participant" "host" {
  alias               = "host-${var.conference_name}"
  description         = "Persistent host for ${var.conference_name}"
  conference          = var.conference_id
  protocol            = "webrtc"
  call_type           = "audio"
  role                = "chair"
  keep_conference_alive = "keep_conference_alive"
  routing             = "auto"
  streaming           = false
  remote_display_name = "Conference Host"
}
```

### Enterprise Conference Automation

```terraform
variable "automatic_participants" {
  type = list(object({
    alias               = string
    description         = string
    protocol            = string
    call_type           = string
    role                = string
    keep_conference_alive = string
    routing             = string
    streaming           = bool
    remote_display_name = string
    presentation_url    = string
  }))
  default = [
    {
      alias               = "enterprise-recorder"
      description         = "Enterprise recording service"
      protocol            = "sip"
      call_type           = "video"
      role                = "guest"
      keep_conference_alive = "end_conference_when_alone"
      routing             = "auto"
      streaming           = true
      remote_display_name = "Enterprise Recorder"
      presentation_url    = ""
    },
    {
      alias               = "compliance-monitor"
      description         = "Compliance monitoring participant"
      protocol            = "webrtc"
      call_type           = "audio"
      role                = "guest"
      keep_conference_alive = "end_conference_when_alone"
      routing             = "auto"
      streaming           = false
      remote_display_name = "Compliance Monitor"
      presentation_url    = ""
    }
  ]
}

resource "pexip_infinity_automatic_participant" "enterprise" {
  count               = length(var.automatic_participants)
  alias               = var.automatic_participants[count.index].alias
  description         = var.automatic_participants[count.index].description
  conference          = data.pexip_infinity_conference.enterprise_meeting.id
  protocol            = var.automatic_participants[count.index].protocol
  call_type           = var.automatic_participants[count.index].call_type
  role                = var.automatic_participants[count.index].role
  keep_conference_alive = var.automatic_participants[count.index].keep_conference_alive
  routing             = var.automatic_participants[count.index].routing
  streaming           = var.automatic_participants[count.index].streaming
  remote_display_name = var.automatic_participants[count.index].remote_display_name
  presentation_url    = var.automatic_participants[count.index].presentation_url != "" ? var.automatic_participants[count.index].presentation_url : null
}
```

### Webinar Automation

```terraform
# Main presentation stream
resource "pexip_infinity_automatic_participant" "webinar_stream" {
  alias               = "webinar-main-stream"
  description         = "Main webinar stream to CDN"
  conference          = data.pexip_infinity_conference.webinar.id
  protocol            = "rtmp"
  call_type           = "video"
  role                = "guest"
  keep_conference_alive = "end_conference_when_alone"
  routing             = "manual"
  system_location     = data.pexip_infinity_system_location.streaming_servers.id
  streaming           = true
  remote_display_name = "Webinar Stream"
  presentation_url    = "rtmp://cdn.example.com/live/${var.webinar_stream_key}"
}

# Backup recording
resource "pexip_infinity_automatic_participant" "webinar_backup" {
  alias               = "webinar-backup-recorder"
  description         = "Backup recording for webinar"
  conference          = data.pexip_infinity_conference.webinar.id
  protocol            = "sip"
  call_type           = "video"
  role                = "guest"
  keep_conference_alive = "end_conference_when_alone"
  routing             = "auto"
  streaming           = true
  remote_display_name = "Backup Recorder"
}

# Webinar host bot
resource "pexip_infinity_automatic_participant" "webinar_host" {
  alias               = "webinar-host-bot"
  description         = "Automated webinar host"
  conference          = data.pexip_infinity_conference.webinar.id
  protocol            = "webrtc"
  call_type           = "audio"
  role                = "chair"
  keep_conference_alive = "keep_conference_alive"
  routing             = "auto"
  streaming           = false
  remote_display_name = "Webinar Host"
}
```

### Development and Testing

```terraform
resource "pexip_infinity_automatic_participant" "test_participant" {
  alias               = "test-${random_id.test.hex}"
  description         = "Test automatic participant"
  conference          = data.pexip_infinity_conference.test_conference.id
  protocol            = "webrtc"
  call_type           = "video"
  role                = "guest"
  keep_conference_alive = "end_conference_when_alone"
  routing             = "auto"
  streaming           = false
  remote_display_name = "Test Participant"
}

resource "random_id" "test" {
  byte_length = 4
}
```

## Schema

### Required

- `alias` (String) - The unique alias of the automatic participant. Maximum length: 250 characters.
- `conference` (String) - The conference URI or reference. Maximum length: 250 characters.
- `protocol` (String) - The protocol for the automatic participant. Valid values: `sip`, `h323`, `rtmp`, `webrtc`.
- `call_type` (String) - The call type. Valid values: `audio`, `video`.
- `role` (String) - The role of the automatic participant. Valid values: `guest`, `chair`.
- `keep_conference_alive` (String) - Conference behavior when only this participant remains. Valid values: `keep_conference_alive`, `end_conference_when_alone`.
- `routing` (String) - The routing type. Valid values: `auto`, `manual`.

### Optional

- `description` (String) - A description of the automatic participant. Maximum length: 250 characters.
- `dtmf_sequence` (String) - DTMF sequence to send when connecting. Maximum length: 250 characters.
- `system_location` (String) - Reference to system location resource URI.
- `streaming` (Boolean) - Whether streaming is enabled. Defaults to `false`.
- `remote_display_name` (String) - The remote display name. Maximum length: 250 characters.
- `presentation_url` (String) - The presentation URL. Maximum length: 250 characters.

### Read-Only

- `id` (String) - Resource URI for the automatic participant in Infinity.
- `resource_id` (Number) - The resource integer identifier for the automatic participant in Infinity.
- `creation_time` (String) - The creation timestamp of the automatic participant.

## Import

Import is supported using the following syntax:

```shell
terraform import pexip_infinity_automatic_participant.example 123
```

Where `123` is the numeric resource ID of the automatic participant.

## Usage Notes

### Protocol Selection

- **SIP**: Standard for VoIP and video calling, widely supported
- **H.323**: Legacy protocol for video conferencing systems
- **RTMP**: Real-Time Messaging Protocol for streaming applications
- **WebRTC**: Web-based real-time communication, modern browser support

### Role Configuration

- **Guest**: Standard participant role with limited permissions
- **Chair**: Conference chair role with administrative permissions
- **Automatic Promotion**: Some configurations may automatically promote participants

### Conference Behavior

- **keep_conference_alive**: Conference continues even when only automatic participants remain
- **end_conference_when_alone**: Conference ends when only automatic participants remain
- **Use Cases**: Choose based on whether the participant is meant to be persistent

### Routing Options

- **auto**: Pexip Infinity automatically selects the best routing path
- **manual**: Specific system location or routing configuration required
- **Performance**: Auto routing typically provides better load distribution

### Streaming Configuration

- **Enabled**: Participant can receive and process media streams
- **Disabled**: Participant operates in a limited mode without media processing
- **Resource Impact**: Streaming participants consume more system resources

### System Location

- Specify when using manual routing
- Must reference a valid system location URI
- Affects which conferencing nodes handle the participant
- Important for geographic distribution and performance

### DTMF Sequences

- Used for automated authentication or menu navigation
- Commonly used with legacy phone systems or IVR systems
- Format: digits (0-9), *, #, and pause characters
- Timing considerations for system response

### Presentation URLs

- Used primarily with RTMP protocol for streaming destinations
- Can specify external streaming endpoints
- Include authentication tokens or keys as needed
- Validate URL format and accessibility

## Troubleshooting

### Common Issues

**Automatic Participant Creation Fails**
- Verify the alias is unique within the system
- Check that the conference URI exists and is accessible
- Ensure the protocol is supported and properly configured
- Verify system location exists if using manual routing

**Participant Doesn't Join Conference**
- Check that the conference is active and accepting participants
- Verify network connectivity for the specified protocol
- Ensure authentication credentials are correct if required
- Monitor conference logs for connection attempts

**Protocol Connection Issues**
- **SIP**: Verify SIP proxy configuration and registration
- **H.323**: Check H.323 gatekeeper settings and registration
- **RTMP**: Validate streaming endpoint accessibility and credentials
- **WebRTC**: Ensure proper certificate configuration and firewall rules

**Audio/Video Quality Problems**
- Check network bandwidth and latency to the participant endpoint
- Verify codec compatibility between participant and conference
- Monitor system resource utilization (CPU, memory, network)
- Adjust media quality settings based on available resources

**Streaming Issues**
- Verify streaming endpoints are accessible and configured correctly
- Check bandwidth requirements for streaming quality
- Monitor streaming server capacity and performance
- Validate streaming URLs and authentication tokens

**Role and Permission Problems**
- Verify the role assignment matches the intended participant behavior
- Check conference settings for role-based restrictions
- Ensure chair participants have necessary permissions
- Monitor conference access control logs

**System Location Routing Issues**
- Verify system location exists and is properly configured
- Check that manual routing is configured correctly
- Ensure network connectivity between system locations
- Monitor load balancing and resource allocation

**DTMF Sequence Problems**
- Verify DTMF sequence format and timing
- Check that receiving system supports DTMF input
- Test DTMF sequence manually to validate expected behavior
- Monitor audio logs for DTMF transmission

**Conference Lifecycle Issues**
- Verify keep_conference_alive setting matches intended behavior
- Check that automatic participants are properly joining conferences
- Monitor conference start and end events
- Ensure proper cleanup when conferences end

**Performance and Scaling Issues**
- Monitor system resources when multiple automatic participants are active
- Check concurrent participant limits and licensing
- Verify network capacity for multiple streams
- Consider load balancing across multiple system locations

**Import Issues**
- Use the numeric resource ID, not the participant alias
- Verify the automatic participant exists in the Infinity cluster
- Check provider authentication credentials have access to the resource
- Confirm the participant configuration is accessible

**Integration Issues**
- Verify external service integration (recording, streaming, monitoring)
- Check API credentials and connectivity to external services
- Monitor integration logs for error conditions
- Test failover scenarios for critical integrations

**Security and Compliance Issues**
- Ensure automatic participants comply with organizational security policies
- Verify recording and streaming meet privacy requirements
- Check access controls and audit logging
- Monitor for unauthorized participant access attempts