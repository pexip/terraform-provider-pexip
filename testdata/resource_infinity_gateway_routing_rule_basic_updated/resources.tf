/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_gateway_routing_rule" "gateway_routing_rule-test" {
  name               = "gateway_routing_rule-test"
  description        = "Updated Test GatewayRoutingRule" // Updated description
  priority           = 156                               // Updated value
  enable             = false                             // Updated to false
  match_string       = "updated-value"                   // Updated value
  replace_string     = "updated-value"                   // Updated value
  called_device_type = "external"                        // Updated value
  outgoing_protocol  = "h323"                            // Updated value
  call_type          = "audio"                           // Updated value
  //ivr_theme          = "updated-value"                   // Updated value
}