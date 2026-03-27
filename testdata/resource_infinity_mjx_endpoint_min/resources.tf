/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_mjx_endpoint" "test" {
  name                = "tf-test mjx-endpoint min"
  room_resource_email = "room@example.com"
  api_address         = "192.168.1.100"
}
