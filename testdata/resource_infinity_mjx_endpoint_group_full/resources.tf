/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_mjx_endpoint_group" "test" {
  name            = "tf-test mjx-endpoint-integration-group full"
  description     = "Test MJX endpoint integration group"
  system_location = "/api/admin/configuration/v1/system_location/2/"
  mjx_integration = "/api/admin/configuration/v1/mjx_integration/1/"
  disable_proxy   = true
}
