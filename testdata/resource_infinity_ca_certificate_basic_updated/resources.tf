resource "pexip_infinity_ca_certificate" "ca_certificate-test" {
  certificate          = "updated-value" // Updated value
  trusted_intermediate = false           // Updated to false
}