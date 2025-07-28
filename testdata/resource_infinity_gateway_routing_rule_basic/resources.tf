resource "pexip_infinity_gateway_routing_rule" "gateway_routing_rule-test" {
  name               = "gateway_routing_rule-test"
  description        = "Test GatewayRoutingRule"
  priority           = 123
  enable             = true
  match_string       = "test-value"
  replace_string     = "test-value"
  called_device_type = "unknown"
  outgoing_protocol  = "sip"
  call_type          = "video"
  ivr_theme          = "test-value"
}