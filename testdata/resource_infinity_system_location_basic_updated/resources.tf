/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_system_location" "main-location" {
  name                = "main"
  description         = "Main location for Pexip Infinity System - updated" # Updated description
  mtu                 = 1460
  dns_servers         = ["/api/admin/configuration/v1/dns_server/1/"] # Updated DNS servers
  ntp_servers         = ["/api/admin/configuration/v1/ntp_server/1/"]
  client_stun_servers = ["/api/admin/configuration/v1/stun_server/2/"] # Updated STUN server
  client_turn_servers = ["/api/admin/configuration/v1/turn_server/2/"] # Updated TURN server
  event_sinks         = ["/api/admin/configuration/v1/event_sink/1/"] # Updated event sink
  stun_server = "/api/admin/configuration/v1/stun_server/2/"
}
