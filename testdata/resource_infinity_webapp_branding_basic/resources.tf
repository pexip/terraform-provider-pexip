resource "pexip_infinity_webapp_branding" "webapp_branding-test" {
  name          = "webapp_branding-test"
  description   = "Test WebappBranding"
  uuid          = "test-value"
  webapp_type   = "pexapp"
  is_default    = true
  branding_file = "test-value"
}