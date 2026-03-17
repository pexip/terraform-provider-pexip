/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_break_in_allow_list_address" "tf-test-break-in-allow-list-address" {
  name                     = "tf-test-break-in-allow-list-address"
  description              = "Full test configuration for break-in allow list address"
  address                  = "10.0.0.0"
  prefix                   = 16
  allowlist_entry_type     = "proxy"
  ignore_incorrect_aliases = false
  ignore_incorrect_pins    = false
}
