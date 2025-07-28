resource "pexip_infinity_certificate_signing_request" "certificate_signing_request-test" {
  subject_name = "certificate_signing_request-test"
  dn = "test-value"
  additional_subject_alt_names = "certificate_signing_request-test"
  private_key_type = "rsa2048"
  private_key = "test-value"
  private_key_passphrase = "test-value"
  ad_compatible = true
  tls_certificate = "test-value"
}