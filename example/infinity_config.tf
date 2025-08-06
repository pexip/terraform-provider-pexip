/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_dns_server" "dns-cloudflare" {
  address     = "1.1.1.1"
  description = "Cloud Flare"

  depends_on = [
    module.gcp-infinity-manager,
  ]
}

resource "pexip_infinity_dns_server" "dns-google-2" {
  address     = "8.8.4.4"
  description = "Google 2"

  depends_on = [
    module.gcp-infinity-manager,
  ]
}

resource "pexip_infinity_ntp_server" "ntp1" {
  address     = "1.pool.ntp.org"
  description = "1.pool.ntp.org"

  depends_on = [
    module.gcp-infinity-manager,
  ]
}

resource "pexip_infinity_syslog_server" "syslog-server-test" {
  address     = "syslog-server-test.local"
  port        = 514
  description = "Test Syslog Server"
  transport   = "tls"
  support_log = true

  depends_on = [
    module.gcp-infinity-manager,
  ]
}

resource "pexip_infinity_system_location" "AMS" {
  name        = "AMS"
  description = "AMS always on"
  mtu         = 1460
  dns_servers = [pexip_infinity_dns_server.dns-cloudflare.id, pexip_infinity_dns_server.dns-google-2.id]
  ntp_servers = [pexip_infinity_ntp_server.ntp1.id]
  // need a syslog server to for worker vm to register properly
  //syslog_servers = [pexip_infinity_syslog_server.syslog-server-test.id]

  depends_on = [
    module.gcp-infinity-manager,
  ]
}

resource "pexip_infinity_system_location" "TEST-AMS" {
  name        = "TEST AMS"
  description = "Set every parameter to test the system location"
  mtu         = 1460
  dns_servers = [pexip_infinity_dns_server.dns-cloudflare.id, pexip_infinity_dns_server.dns-google-2.id]
  ntp_servers = [pexip_infinity_ntp_server.ntp1.id]
  // need a syslog server to for worker vm to register properly
  syslog_servers                 = [pexip_infinity_syslog_server.syslog-server-test.id]
  stun_server                    = pexip_infinity_stun_server.stun-server-test1.id
  turn_server                    = pexip_infinity_turn_server.turn-server-test1.id
  client_turn_servers            = [pexip_infinity_turn_server.turn-server-test1.id, pexip_infinity_turn_server.turn-server-test1.id]
  client_stun_servers            = [pexip_infinity_stun_server.stun-server-test1.id, pexip_infinity_stun_server.stun-server-test1.id]
  sip_proxy                      = pexip_infinity_sip_proxy.sip-proxy-test.id
  h323_gatekeeper                = pexip_infinity_h323_gatekeeper.h323-gatekeeper-test.id
  http_proxy                     = pexip_infinity_http_proxy.http-proxy-test.id
  mssip_proxy                    = pexip_infinity_mssip_proxy.mssip-proxy-test.id
  teams_proxy                    = pexip_infinity_teams_proxy.teams-proxy-test-no-queue.id
  event_sinks                    = [pexip_infinity_event_sink.event-sink-test.id]
  media_qos                      = 46
  signalling_qos                 = 24
  local_mssip_domain             = "test-mssip-domain.local"
  bdpm_pin_checks_enabled        = "ON"
  bdpm_scan_quarantine_enabled   = "ON"
  use_relay_candidates_only      = true
  snmp_network_management_system = pexip_infinity_snmp_network_management_system.snmp-nms-test2.id
  policy_server                  = pexip_infinity_policy_server.policy-server-test.id
  live_captions_dial_out1        = pexip_infinity_system_location.AMS.id
  //live_captions_dial_out2 = pexip_infinity_system_location.AMS.id
  //live_captions_dial_out3 = pexip_infinity_system_location.AMS.id

  // These locations must already exist before being set
  //transcoding_location = pexip_infinity_system_location.AMS.id
  //overflow_location1 = pexip_infinity_system_location.AMS.id
  //overflow_location2 = pexip_infinity_system_location.AMS.id

  depends_on = [
    module.gcp-infinity-manager,
  ]
}

resource "pexip_infinity_turn_server" "turn-server-test1" {
  address     = "turn-server-test1.local"
  port        = 3478
  name        = "TURN Server Test 1"
  description = "Test TURN Server 1"

  depends_on = [
    module.gcp-infinity-manager,
  ]
}

resource "pexip_infinity_turn_server" "turn-server-test2" {
  address     = "turn-server-test2.local"
  port        = 3478
  name        = "TURN Server Test 2"
  description = "Test TURN Server 2"

  depends_on = [
    module.gcp-infinity-manager,
  ]
}

resource "pexip_infinity_turn_server" "turn-server-test3" {
  address     = "turn-server-test3.local"
  port        = 3478
  name        = "TURN Server Test 3"
  description = "Test TURN Server 3"
  username    = "turnuser3"
  // secret
  password = "secretturnpassword3"

  depends_on = [
    module.gcp-infinity-manager,
  ]
}

resource "pexip_infinity_stun_server" "stun-server-test1" {
  address     = "stun-server-test1.local"
  port        = 3478
  name        = "STUN Server Test 1"
  description = "Test STUN Server 1"

  depends_on = [
    module.gcp-infinity-manager,
  ]
}

resource "pexip_infinity_stun_server" "stun-server-test2" {
  address     = "stun-server-test2.local"
  port        = 3478
  name        = "STUN Server Test 2"
  description = "Test STUN Server 2"

  depends_on = [
    module.gcp-infinity-manager,
  ]
}

resource "pexip_infinity_sip_proxy" "sip-proxy-test" {
  address     = "sip-test-proxy.local"
  port        = 5060
  name        = "SIP Proxy Test"
  description = "Test SIP Proxy"
  transport   = "tcp"

  depends_on = [
    module.gcp-infinity-manager,
  ]
}

resource "pexip_infinity_h323_gatekeeper" "h323-gatekeeper-test" {
  address     = "h323-gatekeeper-test.local"
  port        = 1719
  name        = "H323 Gatekeeper Test"
  description = "Test H323 Gatekeeper"

  depends_on = [
    module.gcp-infinity-manager,
  ]
}

resource "pexip_infinity_http_proxy" "http-proxy-test" {
  address  = "http-test-proxy.local"
  port     = 8080
  name     = "HTTP Proxy Test"
  protocol = "http"
  username = "testuser"
  password = "supersecretpassword"

  depends_on = [
    module.gcp-infinity-manager,
  ]
}

resource "pexip_infinity_mssip_proxy" "mssip-proxy-test" {
  address     = "mssip-test-proxy.local"
  description = "Test MSSIP Proxy"
  port        = 5061
  name        = "MSSIP Proxy Test"
  transport   = "tls"

  depends_on = [
    module.gcp-infinity-manager,
  ]
}

resource "pexip_infinity_teams_proxy" "teams-proxy-test-no-queue" {
  azure_tenant            = pexip_infinity_azure_tenant.azure-tenant-test.id
  address                 = "teams-test-proxy.local"
  port                    = 443
  name                    = "Teams Proxy Test no queue"
  description             = "Test Teams Proxy"
  min_number_of_instances = 1

  depends_on = [
    module.gcp-infinity-manager,
  ]
}

resource "pexip_infinity_teams_proxy" "teams-proxy-test-with-queue" {
  azure_tenant            = pexip_infinity_azure_tenant.azure-tenant-test.id
  address                 = "teams-test-proxy.local"
  port                    = 443
  name                    = "Teams Proxy Test with queue"
  description             = "Test Teams Proxy"
  min_number_of_instances = 1
  //notifications_queue should be a secret
  notifications_enabled = true
  notifications_queue   = "Endpoint=sb://test-fooboo-ehn.servicebus.windows.net/;SharedAccessKeyName=RootManageSharedAccessKey;SharedAccessKey=Y7cFX/7z5jzDpWBHkeJsrXZ+CqzleL9D4PjsR/CfaRQ="

  depends_on = [
    module.gcp-infinity-manager,
  ]
}

resource "pexip_infinity_event_sink" "event-sink-test" {
  name        = "Event Sink Test"
  description = "Test Event Sink"
  url         = "https://example.local/event-sink"

  depends_on = [
    module.gcp-infinity-manager,
  ]
}

resource "pexip_infinity_azure_tenant" "azure-tenant-test" {
  tenant_id   = "12345678-1234-1234-1234-123456789012"
  name        = "Azure Tenant Test"
  description = "Test Azure Tenant"

  depends_on = [
    module.gcp-infinity-manager,
  ]
}

resource "pexip_infinity_policy_server" "policy-server-test" {
  name        = "Policy Server Test"
  description = "Test Policy Server"
  url         = "https://policy-server-test.local"
  username    = "policyuser"
  password    = "policypassword"

  depends_on = [
    module.gcp-infinity-manager,
  ]

}

resource "pexip_infinity_snmp_network_management_system" "snmp-nms-test2" {
  address             = "snmp-nms-test2.local"
  port                = 161
  name                = "SNMP NMS Test 2"
  description         = "Test SNMP NMS 2"
  snmp_trap_community = "public-test-trap"

  depends_on = [
    module.gcp-infinity-manager
  ]
}