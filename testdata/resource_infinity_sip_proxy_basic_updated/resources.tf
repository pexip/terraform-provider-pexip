resource "pexip_infinity_sip_proxy" "sip-proxy-test" {
  name        = "test-sip-proxy"
  description = "Test SIP Proxy"
  address     = "test-sip-proxy.dev.pexip.network"
  port        = 8081
  transport   = "tls"
}