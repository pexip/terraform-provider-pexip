resource "pexip_infinity_mjx_integration" "mjx_integration-test" {
  name = "mjx_integration-test"
  description = "Updated Test MjxIntegration"  // Updated description
  start_buffer = 60  // Updated value
  end_buffer = 60  // Updated value
  display_upcoming_meetings = 24  // Updated value
  enable_non_video_meetings = false  // Updated to false
  enable_private_meetings = false  // Updated to false
  ep_username = "mjx_integration-test"
  ep_password = "updated-value"  // Updated value
  ep_use_https = false  // Updated to false
  ep_verify_certificate = false  // Updated to false
  exchange_deployment = "updated-value"  // Updated value
  google_deployment = "updated-value"  // Updated value
  graph_deployment = "updated-value"  // Updated value
  process_alias_private_meetings = false  // Updated to false
  replace_empty_subject = false  // Updated to false
  replace_subject_type = "alias"  // Updated value
  replace_subject_template = "updated-value"  // Updated value
  use_webex = false  // Updated to false
  webex_api_domain = "updated-value"  // Updated value
  webex_client_id = "updated-value"  // Updated value
  webex_client_secret = "updated-value"  // Updated value
  webex_redirect_uri = "updated-value"  // Updated value
  webex_refresh_token = "updated-value"  // Updated value
}