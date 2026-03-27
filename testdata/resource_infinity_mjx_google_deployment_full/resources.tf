resource "pexip_infinity_mjx_google_deployment" "test" {
  name                           = "tf-test mjx-google-deployment full"
  description                    = "Test MJX Google deployment description"
  client_email                   = "test-service@my-project.iam.gserviceaccount.com"
  client_id                      = "123456789012345678901"
  client_secret                  = "test-client-secret"
  private_key                    = "test-private-key"
  use_user_consent               = true
  auth_endpoint                  = "https://accounts.google.com/o/oauth2/v2/auth"
  token_endpoint                 = "https://oauth2.googleapis.com/token"
  redirect_uri                   = "https://pexip.example.com/admin/platform/mjxgoogledeployment/oauth_redirect/"
  maximum_number_of_api_requests = 500000
}
