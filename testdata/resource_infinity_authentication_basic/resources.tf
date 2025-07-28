resource "pexip_infinity_authentication" "authentication-test" {
  source = "local"
  client_certificate = "disabled"
  oidc_client_secret = "kdfjggfkdjhfdvdd"
  ldap_bind_username = "authentication-test"
}