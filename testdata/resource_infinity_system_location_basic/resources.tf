resource "pexip_infinity_system_location" "main-location" {
  name                = "main"
  description         = "Main location for Pexip Infinity System"
  mtu                 = 1460
  dns_servers         = ["/api/admin/configuration/v1/dns_server/2/", "/api/admin/configuration/v1/dns_server/1/"]
  ntp_servers         = ["/api/admin/configuration/v1/ntp_server/1/"]
  client_stun_servers = ["/api/admin/configuration/v1/stun_server/1/", "/api/admin/configuration/v1/stun_server/2/"]
  client_turn_servers = ["/api/admin/configuration/v1/turn_server/1/", "/api/admin/configuration/v1/turn_server/2/"]
  event_sinks         = ["/api/admin/configuration/v1/event_sink/1/"]
}
