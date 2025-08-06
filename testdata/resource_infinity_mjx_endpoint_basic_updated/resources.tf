/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_mjx_endpoint" "mjx_endpoint-test" {
  name                                = "mjx_endpoint-test"
  description                         = "Updated Test MjxEndpoint"   // Updated description
  endpoint_type                       = "cisco"                      // Updated value
  room_resource_email                 = "updated@example.com"        // Updated email
  mjx_endpoint_group                  = "updated-value"              // Updated value
  api_address                         = "updated-server.example.com" // Updated address
  api_username                        = "mjx_endpoint-test"
  api_password                        = "updated-value" // Updated value
  use_https                           = "no"            // Updated value
  verify_cert                         = "no"            // Updated value
  poly_username                       = "mjx_endpoint-test"
  poly_password                       = "updated-value" // Updated value
  poly_raise_alarms_for_this_endpoint = false           // Updated to false
  webex_device_id                     = "updated-value" // Updated value
}