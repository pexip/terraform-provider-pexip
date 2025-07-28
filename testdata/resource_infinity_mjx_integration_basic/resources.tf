resource "pexip_infinity_mjx_integration" "mjx_integration-test" {
  name = "mjx_integration-test"
  description = "Test MjxIntegration"
  start_buffer = 30
  end_buffer = 30
  display_upcoming_meetings = 12
  enable_non_video_meetings = true
  enable_private_meetings = true
  ep_username = "mjx_integration-test"
  ep_password = "test-value"
  ep_use_https = true
  ep_verify_certificate = true
  exchange_deployment = "test-value"
  google_deployment = "test-value"
  graph_deployment = "test-value"
  process_alias_private_meetings = true
  replace_empty_subject = true
  replace_subject_type = "template"
  replace_subject_template = "test-value"
  use_webex = true
  webex_api_domain = "test-value"
  webex_client_id = "test-value"
  webex_client_secret = "test-value"
  webex_redirect_uri = "test-value"
  webex_refresh_token = "test-value"
}