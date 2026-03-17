/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */
resource "pexip_infinity_user_group" "tf-test-user-group" {
  name        = "tf-test-user-group-full"
  description = "tf-test user group description"
}
