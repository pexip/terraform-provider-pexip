/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_log_level" "log_level-test" {
  name  = "log_level-test"
  level = "INFO"
}