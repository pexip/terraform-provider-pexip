/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_system_location" "main-location" {
  name                = "main"
  description         = "Main location for Pexip Infinity System"
  mtu                 = 1460
  dns_servers         = ["/api/admin/configuration/v1/dns_server/2/", "/api/admin/configuration/v1/dns_server/1/"]
  ntp_servers         = ["/api/admin/configuration/v1/ntp_server/1/"]
  client_stun_servers = ["/api/admin/configuration/v1/stun_server/1/", "/api/admin/configuration/v1/stun_server/2/"]
  client_turn_servers = ["/api/admin/configuration/v1/turn_server/1/", "/api/admin/configuration/v1/turn_server/2/"]
  event_sinks         = ["/api/admin/configuration/v1/event_sink/1/"]
  bdpm_pin_checks_enabled = "ON"
  bdpm_scan_quarantine_enabled = "ON"
  h323_gatekeeper = "/api/admin/configuration/v1/h323_gatekeeper/1/"
  http_proxy = "/api/admin/configuration/v1/http_proxy/1/"
  live_captions_dial_out_1 = "/api/admin/configuration/v1/system_location/1/"
  live_captions_dial_out_2 = "/api/admin/configuration/v1/system_location/2/"
  live_captions_dial_out_3 = "/api/admin/configuration/v1/system_location/3/"
  local_mssip_domain = "test-mssip-domain.local"
  media_qos = "46"
  mssip_proxy = "/api/admin/configuration/v1/mssip_proxy/1/"
  overflow_location1 = "/api/admin/configuration/v1/system_location/1/"
  overflow_location2 = "/api/admin/configuration/v1/system_location/2/"
  policy_server = "/api/admin/configuration/v1/policy_server/1/"
  signalling_qos = "24"
  sip_proxy = "/api/admin/configuration/v1/sip_proxy/1/"
  snmp_network_management_system = "/api/admin/configuration/v1/snmp_network_management_system/2/"
  stun_server = "/api/admin/configuration/v1/stun_server/1/"
  syslog_servers = null
  teams_proxy = "/api/admin/configuration/v1/teams_proxy/1/"
  transcoding_location = "/api/admin/configuration/v1/system_location/3/"
  turn_server = "/api/admin/configuration/v1/turn_server/3/"
  use_relay_candidates_only = "true"
}
