/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "google_service_account" "infinity-sa" {
  project      = var.project_id
  account_id   = "infinity-sa"
  description  = "Pexip Infinity Service Account"
  display_name = "Pexip Infinity Service Account"

}

