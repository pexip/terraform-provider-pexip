resource "pexip_infinity_conference" "conference-test" {
  name                  = "conference-test"
  description           = "Updated Test Conference" // Updated description
  service_type          = "conference"              // Keep same service type (not updatable)
  pin                   = "9876"                    // Updated PIN
  guest_pin             = "4321"                    // Updated guest PIN
  allow_guests          = false                     // Updated to false
  guests_muted          = true                      // Updated to true
  hosts_can_unmute      = false                     // Updated to false
  max_pixels_per_second = 1280000                   // Updated pixels per second
  tag                   = "updated-tag"             // Updated tag
}