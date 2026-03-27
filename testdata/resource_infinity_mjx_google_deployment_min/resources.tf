resource "pexip_infinity_mjx_google_deployment" "test" {
  name         = "tf-test mjx-google-deployment min"
  client_email = "test-service@my-project.iam.gserviceaccount.com"
  private_key  = "test-private-key"
}
