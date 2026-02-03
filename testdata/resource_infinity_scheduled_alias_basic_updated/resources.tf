/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_ms_exchange_connector" "test-connector" {
  name                  = "test-exchange-connector"
  description           = "Test Exchange Connector"
  addin_server_domain   = "test-domain"
}

resource "pexip_infinity_scheduled_alias" "scheduled_alias-test" {
  alias              = "updated-scheduled-alias"              // Updated value
  alias_number       = 9876543210                             // Updated value
  numeric_alias      = "987654"                               // Updated value
  uuid               = "22222222-2222-2222-2222-222222222222" // Updated value
  exchange_connector = pexip_infinity_ms_exchange_connector.test-connector.id
  is_used            = false             // Updated to false
  ews_item_uid       = "updated-ews-uid" // Updated value

  depends_on = [
    pexip_infinity_ms_exchange_connector.test-connector
  ]
}