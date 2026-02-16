/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

# System locations
resource "pexip_infinity_system_location" "test1" {
  name = "tf-test 1"
}

resource "pexip_infinity_system_location" "test2" {
  name = "tf-test 2"
}

resource "pexip_infinity_system_location" "test3" {
  name = "tf-test 3"
}

# DNS Servers
resource "pexip_infinity_dns_server" "dns1" {
  address = "68.94.156.1"
}

resource "pexip_infinity_dns_server" "dns2" {
  address = "68.94.157.1"
}

# NTP Server
resource "pexip_infinity_ntp_server" "ntp1" {
  address = "time.google.com"
}

# STUN Servers
resource "pexip_infinity_stun_server" "stun1" {
  name    = "tf-test-stun1"
  address = "stun1.pexvclab.com"
  port    = 3478
}

resource "pexip_infinity_stun_server" "stun2" {
  name    = "tf-test-stun2"
  address = "stun2.pexvclab.com"
  port    = 3478
}

# TURN Servers
resource "pexip_infinity_turn_server" "turn1" {
  name           = "tf-test-turn1"
  address        = "turn1.pexvclab.com"
  server_type    = "namepsw"
  transport_type = "udp"
}

resource "pexip_infinity_turn_server" "turn2" {
  name           = "tf-test-turn2"
  address        = "turn2.pexvclab.com"
  server_type    = "namepsw"
  transport_type = "udp"
}

# Event Sink
resource "pexip_infinity_event_sink" "event1" {
  name                     = "tf-test-event-sink"
  url                      = "https://events.pexvclab.com/webhook"
  bulk_support             = false
  verify_tls_certificate   = false
  version                  = 2
}

# H.323 Gatekeeper
resource "pexip_infinity_h323_gatekeeper" "h323gk" {
  name    = "tf-test-h323-gk"
  address = "gatekeeper.pexvclab.com"
}

# HTTP Proxy
resource "pexip_infinity_http_proxy" "proxy" {
  name     = "tf-test-http-proxy"
  address  = "proxy.pexvclab.com"
  protocol = "http"
}

# MSSIP Proxy
resource "pexip_infinity_mssip_proxy" "mssip" {
  name      = "tf-test-mssip"
  address   = "sfb.pexvclab.com"
  transport = "tls"
}

# Policy Server
resource "pexip_infinity_policy_server" "policy" {
  name                                    = "tf-test-policy"
  url                                     = "https://policy.pexvclab.com"
  enable_avatar_lookup                    = false
  enable_directory_lookup                 = false
  enable_internal_media_location_policy   = false
  enable_internal_participant_policy      = false
  enable_internal_service_policy          = false
}

# SIP Proxy
resource "pexip_infinity_sip_proxy" "sip" {
  name      = "tf-test-sip"
  address   = "sip.pexvclab.com"
  transport = "tls"
}

# SNMP Network Management System
resource "pexip_infinity_snmp_network_management_system" "snmp" {
  name                 = "tf-test-snmp"
  description          = "Test SNMP System"
  address              = "snmp.pexvclab.com"
  port                 = 162
  snmp_trap_community  = "public"
}

# Azure Tenant for Teams Proxy
resource "pexip_infinity_azure_tenant" "azure" {
  name      = "tf-test-azure-tenant-system-location-full"
  tenant_id = "12345678-1234-1234-1234-123456789012"
}

# Teams Proxy
resource "pexip_infinity_teams_proxy" "teams" {
  name         = "tf-test-teams-proxy-for-location"
  address      = "teams.pexvclab.com"
  azure_tenant = pexip_infinity_azure_tenant.azure.id
}

# Main System Location
resource "pexip_infinity_system_location" "main-location" {
  name                           = "tf-test-system-location-full"
  description                    = "Full configuration test location"
  mtu                            = 1460
  media_qos                      = 46
  signalling_qos                 = 24
  local_mssip_domain             = "test-mssip.pexvclab.com"
  bdpm_pin_checks_enabled        = "ON"
  bdpm_scan_quarantine_enabled   = "ON"
  use_relay_candidates_only      = true

  # Related resources - created above and referenced by ID
  dns_servers         = [pexip_infinity_dns_server.dns1.id, pexip_infinity_dns_server.dns2.id]
  ntp_servers         = [pexip_infinity_ntp_server.ntp1.id]
  client_stun_servers = [pexip_infinity_stun_server.stun1.id, pexip_infinity_stun_server.stun2.id]
  client_turn_servers = [pexip_infinity_turn_server.turn1.id, pexip_infinity_turn_server.turn2.id]
  event_sinks         = [pexip_infinity_event_sink.event1.id]

  # to-one relationships
  h323_gatekeeper                = pexip_infinity_h323_gatekeeper.h323gk.id
  http_proxy                     = pexip_infinity_http_proxy.proxy.id
  mssip_proxy                    = pexip_infinity_mssip_proxy.mssip.id
  policy_server                  = pexip_infinity_policy_server.policy.id
  sip_proxy                      = pexip_infinity_sip_proxy.sip.id
  snmp_network_management_system = pexip_infinity_snmp_network_management_system.snmp.id
  stun_server                    = pexip_infinity_stun_server.stun1.id
  teams_proxy                    = pexip_infinity_teams_proxy.teams.id
  turn_server                    = pexip_infinity_turn_server.turn1.id

  # Circular reference fields - these reference other system_locations
  # Using hardcoded URIs since creating them would cause circular dependencies
  overflow_location1       = pexip_infinity_system_location.test1.id
  overflow_location2       = pexip_infinity_system_location.test2.id
  transcoding_location     = pexip_infinity_system_location.test3.id
  live_captions_dial_out_1 = pexip_infinity_system_location.test1.id
  live_captions_dial_out_2 = pexip_infinity_system_location.test2.id
  live_captions_dial_out_3 = pexip_infinity_system_location.test3.id
}
