resource "pexip_infinity_mssip_proxy" "mssip-proxy-test" {
  name        = "mssip-proxy-test"
  description = "Updated Test MSSIP proxy"     // Updated description
  address     = "updated-mssip-proxy.dev.pexip.network"  // Updated address
  port        = 5061                          // Updated port
  transport   = "tls"                         // Updated transport
}