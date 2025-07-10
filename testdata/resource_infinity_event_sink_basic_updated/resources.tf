resource "pexip_infinity_event_sink" "event-sink-test" {
  name        = "test-event-sink"
  description = "Test Event Sink"
  url         = "https://test-event-sink.dev.pexip.network"
  username    = "testuser"
  password    = "updatedpassword"
}