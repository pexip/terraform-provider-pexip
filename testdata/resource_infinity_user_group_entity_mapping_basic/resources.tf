/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_user_group_entity_mapping" "user_group_entity_mapping-test" {
  description         = "Test UserGroupEntityMapping"
  entity_resource_uri = "test-value"
  user_group          = "test-value"
}