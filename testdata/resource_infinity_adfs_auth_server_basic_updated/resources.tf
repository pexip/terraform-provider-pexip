resource "pexip_infinity_adfs_auth_server" "adfs_auth_server-test" {
  name = "adfs_auth_server-test"
  description = "Updated Test ADFSAuthServer"  // Updated description
  client_id = "updated-value"  // Updated value
  federation_service_name = "adfs_auth_server-test"
  federation_service_identifier = "updated-value"  // Updated value
  relying_party_trust_identifier_url = "https://updated.example.com"  // Updated URL
}