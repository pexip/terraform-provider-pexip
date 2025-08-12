/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "google_project_iam_binding" "compute-instance-admin" {
  project = var.project_id
  role    = "roles/compute.instanceAdmin.v1"
  members = [
    google_service_account.infinity-sa.member,
  ]
}

resource "google_project_iam_binding" "secret-manager-secret-accessor" {
  project = var.project_id
  role    = "roles/secretmanager.secretAccessor"
  members = [
    google_service_account.infinity-sa.member,
  ]
}

resource "google_project_iam_binding" "storage-object-admin" {
  project = var.project_id
  role    = "roles/storage.objectAdmin"
  members = [
    google_service_account.infinity-sa.member,
  ]
}
