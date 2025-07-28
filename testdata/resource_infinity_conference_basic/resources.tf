resource "pexip_infinity_conference" "conference-test" {
  name                  = "conference-test"
  description           = "Test Conference"
  service_type          = "conference"
  pin                   = "1234"
  guest_pin             = "5678"
  allow_guests          = true
  guests_muted          = false
  hosts_can_unmute      = true
  max_pixels_per_second = 1920000
  tag                   = "test-tag"
}