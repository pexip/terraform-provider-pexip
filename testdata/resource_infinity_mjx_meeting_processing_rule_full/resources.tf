/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_mjx_meeting_processing_rule" "test" {
  name                       = "tf-test mjx-meeting-processing-rule full"
  description                = "Test MJX meeting processing rule"
  priority                   = 10
  enabled                    = false
  meeting_type               = "teams"
  mjx_integration            = "/api/admin/configuration/v1/mjx_integration/1/"
  transform_rule             = "{{ domain }}"
  domain                     = "example.com"
  company_id                 = "test-company-id"
  include_pin                = true
  default_processing_enabled = false
}
