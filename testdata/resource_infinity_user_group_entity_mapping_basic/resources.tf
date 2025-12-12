/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_user_group" "test-group" {
  name        = "test-user-group"
  description = "Test User Group"
}

resource "pexip_infinity_conference" "test-conference" {
  name         = "test-conference"
  description  = "Test Conference"
  service_type = "conference"
}

resource "pexip_infinity_user_group_entity_mapping" "user_group_entity_mapping-test" {
  description         = "Test UserGroupEntityMapping"
  entity_resource_uri = pexip_infinity_conference.test-conference.id
  user_group          = pexip_infinity_user_group.test-group.id

  depends_on = [
    pexip_infinity_user_group.test-group,
    pexip_infinity_conference.test-conference
  ]
}