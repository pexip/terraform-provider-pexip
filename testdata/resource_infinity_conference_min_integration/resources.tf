/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_global_configuration" "config" {
  enable_breakout_rooms = true
}

resource "pexip_infinity_conference" "tf-test-conference" {
  name = "tf-test-conference"
}
