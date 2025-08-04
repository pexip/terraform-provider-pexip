resource "pexip_infinity_automatic_participant" "automatic-participant-test" {
  alias                 = "automatic-participant-updated"     // Updated value
  description           = "Updated Test AutomaticParticipant" // Updated description
  conference            = "updated-conference"                // Updated value
  protocol              = "h323"                              // Updated value
  call_type             = "audio"                             // Updated value
  role                  = "chair"                             // Updated value
  dtmf_sequence         = "456*"                              // Updated value
  keep_conference_alive = "end_conference_when_alone"         // Updated value
  routing               = "manual"                            // Updated value
  system_location       = "updated-location"                  // Updated value
  streaming             = false                               // Updated to false
  remote_display_name   = "automatic_participant-test"
  presentation_url      = "https://updated.example.com" // Updated URL
}