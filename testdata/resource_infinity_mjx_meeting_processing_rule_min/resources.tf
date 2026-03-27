/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_mjx_meeting_processing_rule" "test" {
  name            = "tf-test mjx-meeting-processing-rule min"
  priority        = 1
  meeting_type    = "pexipinfinity"
  mjx_integration = "/api/admin/configuration/v1/mjx_integration/1/"
}
