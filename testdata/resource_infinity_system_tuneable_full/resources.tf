/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_system_tuneable" "system_tuneable-test" {
  name    = "tf-test-system-tuneable-full"
  setting = "full-test-value"
}
