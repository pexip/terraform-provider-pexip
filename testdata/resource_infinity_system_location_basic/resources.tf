resource "pexip_infinity_system_location" "main-location" {
  name        = "main"
  description = "Main location for Pexip Infinity System"
  mtu         = 1460
  dns_servers = ["/api/admin/configuration/v1/dns_server/1/", "/api/admin/configuration/v1/dns_server/2/"]
  ntp_servers = ["/api/admin/configuration/v1/ntp_server/1/"]
}
