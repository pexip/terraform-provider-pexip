resource "pexip_infinity_device" "device-test" {
  alias                           = "device-test"
  description                     = "Test Device"
  username                        = "deviceuser"
  password                        = "devicepass"
  primary_owner_email_address     = "owner@example.com"
  enable_sip                      = true
  enable_h323                     = false
  enable_infinity_connect_non_sso = true
  enable_infinity_connect_sso     = false
  enable_standard_sso             = false
  tag                             = "test-tag"
  sync_tag                        = "sync-tag"
}