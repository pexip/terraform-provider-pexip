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
  alias              = "test-scheduled-alias"
  alias_number       = 1234567890
  numeric_alias      = "123456"
  uuid               = "11111111-1111-1111-1111-111111111111"
  exchange_connector = pexip_infinity_ms_exchange_connector.test-connector.id
  is_used            = true
  ews_item_uid       = "test-ews-uid"

  depends_on = [
    pexip_infinity_ms_exchange_connector.test-connector
  ]
}