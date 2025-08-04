resource "pexip_infinity_registration" "registration-test" {
  enable                        = true
  refresh_strategy              = "adaptive"
  route_via_registrar           = true
  enable_push_notifications     = true
  enable_google_cloud_messaging = true
  push_token                    = "test-value"
}