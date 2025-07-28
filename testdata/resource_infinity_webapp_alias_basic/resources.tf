resource "pexip_infinity_webapp_alias" "webapp_alias-test" {
  slug        = "test-value"
  description = "Test WebappAlias"
  webapp_type = "pexapp"
  is_enabled  = true
  bundle      = "test-value"
  branding    = "test-value"
}