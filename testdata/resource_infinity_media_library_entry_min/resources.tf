/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_media_library_entry" "media_library_entry-test" {
  name       = "tf-test-media-library-entry"
  media_file = "${path.module}/earth.mp4"
}
