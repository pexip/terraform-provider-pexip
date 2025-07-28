resource "pexip_infinity_dns_server" "dns-cloudflare" {
  address     = "1.1.1.1"
  description = "Cloud Flare"

  depends_on = [
    google_compute_instance.infinity_manager,
    null_resource.wait_for_infinity_manager_http
  ]
}

resource "pexip_infinity_dns_server" "dns-google-2" {
  address     = "8.8.4.4"
  description = "Google 2"

  depends_on = [
    google_compute_instance.infinity_manager,
    null_resource.wait_for_infinity_manager_http
  ]
}

resource "pexip_infinity_ntp_server" "ntp1" {
  address     = "1.pool.ntp.org"
  description = "1.pool.ntp.org"

  depends_on = [
    google_compute_instance.infinity_manager,
    null_resource.wait_for_infinity_manager_http
  ]
}

resource "pexip_infinity_system_location" "AMS" {
  name           = "AMS"
  description    = "AMS always on"
  mtu            = 1460
  dns_servers    = [pexip_infinity_dns_server.dns-cloudflare.id, pexip_infinity_dns_server.dns-google-2.id]
  ntp_servers    = [pexip_infinity_ntp_server.ntp1.id]
  syslog_servers = [pexip_infinity_syslog_server.syslog-server-test.id]

  depends_on = [
    google_compute_instance.infinity_manager,
    null_resource.wait_for_infinity_manager_http
  ]
}

resource "pexip_infinity_tls_certificate" "tls-cert-test" {
  certificate = tls_self_signed_cert.manager_cert.cert_pem
  private_key = tls_private_key.manager_private_key.private_key_pem
  nodes       = ["${local.hostname}.${local.domain}"]

  depends_on = [
    google_compute_instance.infinity_manager,
    null_resource.wait_for_infinity_manager_http
  ]
}

// specify a TLS cert and private key
/*
resource "pexip_infinity_tls_certificate" "tls-cert-test2" {
  certificate = file("gcp-infinity-manager/test-cert.pem")
  private_key = file("gcp-infinity-manager/test-key.key")

  depends_on = [
    google_compute_instance.infinity_manager,
    null_resource.wait_for_infinity_manager_http
  ]
}
*/

resource "pexip_infinity_licence" "license" {
  entitlement_id = var.license_key

  depends_on = [
    google_compute_instance.infinity_manager,
    null_resource.wait_for_infinity_manager_http
  ]
}

resource "pexip_infinity_azure_tenant" "azure-tenant-test" {
  tenant_id   = "12345678-1234-1234-1234-123456789012"
  name        = "Azure Tenant Test"
  description = "Test Azure Tenant"

  depends_on = [
    google_compute_instance.infinity_manager,
    null_resource.wait_for_infinity_manager_http,
    pexip_infinity_licence.license
  ]
}

resource "pexip_infinity_event_sink" "event-sink-test" {
  name        = "Event Sink Test"
  description = "Test Event Sink"
  url         = "https://example.local/event-sink"

  depends_on = [
    google_compute_instance.infinity_manager,
    null_resource.wait_for_infinity_manager_http
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
    google_compute_instance.infinity_manager,
    null_resource.wait_for_infinity_manager_http
  ]
}

// This is an example of a Teams Proxy with a queue, currently causes an error in the provider
/*
resource "pexip_infinity_teams_proxy" "teams-proxy-test-with-queue" {
  azure_tenant = pexip_infinity_azure_tenant.azure-tenant-test.id
  address     = "teams-test-proxy.local"
  port        = 443
  name        = "Teams Proxy Test with queue"
  description = "Test Teams Proxy"
  min_number_of_instances = 1
  //notifications_queue should be a secret
  notifications_enabled = true
  notifications_queue = "Endpoint=sb://test-fooboo-ehn.servicebus.windows.net/;SharedAccessKeyName=RootManageSharedAccessKey;SharedAccessKey=Y7cFX/7z5jzDpWBHkeJsrXZ+CqzleL9D4PjsR/CfaRQ="

  depends_on = [
    google_compute_instance.infinity_manager,
    null_resource.wait_for_infinity_manager_http
  ]
}
*/

resource "pexip_infinity_mssip_proxy" "mssip-proxy-test" {
  address     = "mssip-test-proxy.local"
  description = "Test MSSIP Proxy"
  port        = 5061
  name        = "MSSIP Proxy Test"
  transport   = "tls"

  depends_on = [
    google_compute_instance.infinity_manager,
    null_resource.wait_for_infinity_manager_http
  ]
}

resource "pexip_infinity_http_proxy" "http-proxy-test" {
  address  = "http-test-proxy.local"
  port     = 8080
  name     = "HTTP Proxy Test"
  protocol = "http"
  // password causes an error in the provider
  //password    = "supersecretpassword"

  depends_on = [
    google_compute_instance.infinity_manager,
    null_resource.wait_for_infinity_manager_http
  ]
}

resource "pexip_infinity_h323_gatekeeper" "h323-gatekeeper-test" {
  address     = "h323-gatekeeper-test.local"
  port        = 1719
  name        = "H323 Gatekeeper Test"
  description = "Test H323 Gatekeeper"

  depends_on = [
    google_compute_instance.infinity_manager,
    null_resource.wait_for_infinity_manager_http
  ]
}

resource "pexip_infinity_sip_proxy" "sip-proxy-test" {
  address     = "sip-test-proxy.local"
  port        = 5060
  name        = "SIP Proxy Test"
  description = "Test SIP Proxy"
  transport   = "tcp"

  depends_on = [
    google_compute_instance.infinity_manager,
    null_resource.wait_for_infinity_manager_http
  ]
}

resource "pexip_infinity_snmp_network_management_system" "snmp-nms-test" {
  address             = "snmp-nms-test.local"
  port                = 161
  name                = "SNMP NMS Test"
  description         = "Test SNMP NMS"
  snmp_trap_community = "public-test-trap"

  depends_on = [
    google_compute_instance.infinity_manager,
    null_resource.wait_for_infinity_manager_http
  ]
}

resource "pexip_infinity_stun_server" "stun-server-test1" {
  address     = "stun-server-test1.local"
  port        = 3478
  name        = "STUN Server Test 1"
  description = "Test STUN Server 1"

  depends_on = [
    google_compute_instance.infinity_manager,
    null_resource.wait_for_infinity_manager_http
  ]
}

resource "pexip_infinity_stun_server" "stun-server-test2" {
  address     = "stun-server-test2.local"
  port        = 3478
  name        = "STUN Server Test 2"
  description = "Test STUN Server 2"

  depends_on = [
    google_compute_instance.infinity_manager,
    null_resource.wait_for_infinity_manager_http
  ]
}

resource "pexip_infinity_turn_server" "turn-server-test1" {
  address     = "turn-server-test1.local"
  port        = 3478
  name        = "TURN Server Test 1"
  description = "Test TURN Server 1"

  depends_on = [
    google_compute_instance.infinity_manager,
    null_resource.wait_for_infinity_manager_http
  ]
}

resource "pexip_infinity_turn_server" "turn-server-test2" {
  address     = "turn-server-test2.local"
  port        = 3478
  name        = "TURN Server Test 2"
  description = "Test TURN Server 2"

  depends_on = [
    google_compute_instance.infinity_manager,
    null_resource.wait_for_infinity_manager_http
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
    google_compute_instance.infinity_manager,
    null_resource.wait_for_infinity_manager_http
  ]
}

resource "pexip_infinity_syslog_server" "syslog-server-test" {
  address     = "syslog-server-test.local"
  port        = 514
  description = "Test Syslog Server"
  transport   = "tls"
  support_log = true

  depends_on = [
    google_compute_instance.infinity_manager,
    null_resource.wait_for_infinity_manager_http
  ]
}