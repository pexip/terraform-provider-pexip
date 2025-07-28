resource "pexip_infinity_adfs_auth_server" "adfs_auth_server-test" {
  name = "adfs_auth_server-test"
  description = "Test ADFSAuthServer"
  client_id = "test-value"
  federation_service_name = "adfs_auth_server-test"
  federation_service_identifier = "test-value"
  relying_party_trust_identifier_url = "https://example.com"
}