/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_user_group_entity_mapping" "user_group_entity_mapping-test" {
  description         = "Updated Test UserGroupEntityMapping" // Updated description
  entity_resource_uri = "updated-value"                       // Updated value
  user_group          = "updated-value"                       // Updated value
}