/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_policy_server" "tf-test-policy-server" {
  name                                      = "tf-test-policy-server"
  description                               = "tf-test Policy Server Description"
  url                                       = "https://policy.example.com"
  username                                  = "tf-test-user"
  password                                  = "tf-test-password"
  enable_service_lookup                     = true
  enable_participant_lookup                 = true
  enable_registration_lookup                = true
  enable_directory_lookup                   = true
  enable_avatar_lookup                      = true
  enable_media_location_lookup              = true
  enable_internal_service_policy            = true
  enable_internal_participant_policy        = true
  enable_internal_media_location_policy     = true
  prefer_local_avatar_configuration         = true
  internal_service_policy_template          = "tf-test service template"
  internal_participant_policy_template      = "tf-test participant template"
  internal_media_location_policy_template   = "tf-test media location template"
}
