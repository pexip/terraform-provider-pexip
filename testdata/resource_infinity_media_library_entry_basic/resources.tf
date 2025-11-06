/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_media_library_entry" "media_library_entry-test" {
  name        = "media_library_entry-test"
  description = "Test MediaLibraryEntry"
  uuid        = "test-value"
  media_file  = "test.mp4"
}