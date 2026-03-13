/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_role" "tf-test-role" {
  name = "tf-test-role"
  permissions = [
    "/api/admin/configuration/v1/permission/1/",
    "/api/admin/configuration/v1/permission/2/",
  ]
}
