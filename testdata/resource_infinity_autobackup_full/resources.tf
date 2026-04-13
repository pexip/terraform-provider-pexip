/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_autobackup" "autobackup-test" {
  autobackup_enabled         = true
  autobackup_interval        = 12
  autobackup_passphrase      = "SecretPassphrase123"
  autobackup_start_hour      = 2
  autobackup_upload_url      = "ftp://backup.example.com/pexip"
  autobackup_upload_username = "backupuser"
  autobackup_upload_password = "BackupPassword123"
}
