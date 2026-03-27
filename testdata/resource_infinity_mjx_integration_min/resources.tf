/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_mjx_integration" "test" {
  name             = "tf-test mjx-integration min"
  graph_deployment = "/api/admin/configuration/v1/mjx_graph_deployment/1/"
}
