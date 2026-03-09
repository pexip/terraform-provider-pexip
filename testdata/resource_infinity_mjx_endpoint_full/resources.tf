/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_mjx_endpoint" "cisco" {
  name                = "tf-test Cisco mjx-endpoint full"
  room_resource_email = "cisco-full@test.local"
  api_address         = "192.168.1.10"
  description         = "Full test MjxEndpoint"
  endpoint_type       = "CISCO"
  mjx_endpoint_group  = "test-group"
  api_username        = "admin"
  api_password        = "SuperSecretCiscoPassword456!"
  use_https           = "YES"
  verify_cert         = "YES"
  webex_device_id     = "webex-device-123"
}

resource "pexip_infinity_mjx_endpoint" "poly" {
  name                                = "tf-test Poly mjx-endpoint full"
  room_resource_email                 = "poly-full@test.local"
  api_address                         = "192.168.1.20"
  endpoint_type                       = "POLY"
  description                         = "Full test Poly MjxEndpoint"
  mjx_endpoint_group                  = "test-group"
  api_username                        = "admin"
  api_password                        = "SuperSecretPolyPassword456!"
  poly_username                       = "admin"
  poly_password                       = "SuperSecretPolyPassword789!"
  poly_raise_alarms_for_this_endpoint = true
  use_https                           = "YES"
  verify_cert                         = "YES"
}