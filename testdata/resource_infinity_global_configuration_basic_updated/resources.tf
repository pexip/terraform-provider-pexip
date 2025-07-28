resource "pexip_infinity_global_configuration" "global_configuration-test" {
  enable_webrtc = false  // Updated to false
  enable_sip = false  // Updated to false
  enable_h323 = false  // Updated to false
  enable_rtmp = false  // Updated to false
  crypto_mode = "required"  // Updated value
  max_pixels_per_second = "1080000"  // Updated value
  bursting_enabled = false  // Updated to false
  cloud_provider = "azure"  // Updated value
  aws_access_key = "updated-key"  // Updated value
  aws_secret_key = "updated-secret"  // Updated value
  azure_client_id = "updated-client"  // Updated value
  azure_secret = "updated-secret"  // Updated value
  conference_create_permissions = "admin_only"  // Updated value
  conference_creation_mode = "per_node"  // Updated value
  enable_analytics = false  // Updated to false
  enable_error_reporting = false  // Updated to false
  bandwidth_restrictions = "none"  // Updated value
  administrator_email = "updated@example.com"  // Updated email
  media_ports_start = 50000  // Updated value
  media_ports_end = 50100  // Updated value
  signalling_ports_start = 5080  // Updated value
  signalling_ports_end = 5090  // Updated value
  guests_only_timeout = 600  // Updated value
  waiting_for_chair_timeout = 900  // Updated value
}