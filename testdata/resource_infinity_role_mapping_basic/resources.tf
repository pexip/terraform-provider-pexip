resource "pexip_infinity_role_mapping" "role_mapping-test" {
  name   = "role_mapping-test"
  source = "saml_attribute"
  value  = "test-value"
}