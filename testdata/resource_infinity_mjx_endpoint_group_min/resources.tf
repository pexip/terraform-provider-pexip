/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_mjx_endpoint_group" "test" {
  name            = "tf-test mjx-endpoint-integration-group min"
  system_location = "/api/admin/configuration/v1/system_location/1/"
}
