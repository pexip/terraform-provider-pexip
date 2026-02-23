/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_mjx_endpoint" "cisco" {
  name = "tf-test Cisco mjx-endpoint min"
  room_resource_email = "cisco@test.local"
  api_address = "192.168.1.10"
}

resource "pexip_infinity_mjx_endpoint" "poly" {
  name = "tf-test Poly mjx-endpoint min"
  room_resource_email = "poly@test.local"
  api_address = "192.168.1.20"
  endpoint_type = "POLY"
  # required for poly endpoint
  poly_username = "admin"
  poly_password = "SuperSecretPolyPassword789!"
}