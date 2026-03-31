/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_mjx_endpoint" "test" {
  name                                = "tf-test mjx-endpoint full"
  description                         = "Test MJX endpoint description"
  endpoint_type                       = "CISCO"
  room_resource_email                 = "room2@example.com"
  api_address                         = "192.168.1.101"
  api_port                            = 443
  api_username                        = "apiuser"
  api_password                        = "apipassword"
  use_https                           = "YES"
  verify_cert                         = "YES"
  poly_username                       = "polyuser"
  poly_password                       = "polypassword"
  poly_raise_alarms_for_this_endpoint = false
  webex_device_id                     = "device-123"
}
