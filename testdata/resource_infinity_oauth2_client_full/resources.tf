/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_oauth2_client" "oauth2_client-test" {
  client_name = "tf-test oauth2_client RW"
  role        = "/api/admin/configuration/v1/role/1/"
}