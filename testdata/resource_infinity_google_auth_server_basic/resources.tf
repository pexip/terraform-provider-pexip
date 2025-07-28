resource "pexip_infinity_google_auth_server" "google_auth_server-test" {
  name = "google_auth_server-test"
  description = "Test GoogleAuthServer"
  application_type = "web"
  client_id = "test-value"
  client_secret = "test-value"
}