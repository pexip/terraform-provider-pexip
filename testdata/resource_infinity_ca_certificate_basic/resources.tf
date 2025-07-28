resource "pexip_infinity_ca_certificate" "ca_certificate-test" {
  certificate          = "test-value"
  trusted_intermediate = true
}