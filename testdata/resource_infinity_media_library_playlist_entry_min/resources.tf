/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_media_library_entry" "test" {
  name       = "tf-test-media-entry"
  media_file = "${path.module}/earth.mp4"
}

resource "pexip_infinity_media_library_playlist" "test" {
  name = "tf-test-playlist"
}

resource "pexip_infinity_media_library_playlist_entry" "test" {
  playlist = pexip_infinity_media_library_playlist.test.id
}
