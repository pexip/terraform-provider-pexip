resource "tls_private_key" "test_min" {
  algorithm = "RSA"
  rsa_bits  = 2048
}

resource "tls_private_key" "test_full" {
  algorithm = "RSA"
  rsa_bits  = 2048
}

resource "pexip_infinity_mjx_google_deployment" "test" {
  name         = "tf-test mjx-google-deployment min"
  client_email = "test-service@my-project.iam.gserviceaccount.com"
  private_key  = tls_private_key.test_min.private_key_pem
}
