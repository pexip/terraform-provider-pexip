/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_user_group" "user-group-test" {
  name        = "user-group-test"
  description = "Updated Test User Group" // Updated description
}