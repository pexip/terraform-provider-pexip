resource "pexip_infinity_automatic_participant" "automatic-participant-test" {
  alias = "automatic-participant-test"
  description = "Test AutomaticParticipant"
  conference = "test-conference"
  protocol = "sip"
  call_type = "video"
  role = "guest"
  dtmf_sequence = "123#"
  keep_conference_alive = "keep_conference_alive"
  routing = "auto"
  system_location = "test-location"
  streaming = true
  remote_display_name = "automatic_participant-test"
  presentation_url = "https://example.com"
}