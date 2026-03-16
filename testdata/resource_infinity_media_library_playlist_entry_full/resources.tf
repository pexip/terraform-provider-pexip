/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_media_library_entry" "test" {
  name        = "tf-test-media-entry"
  description = "Test media entry for playlist"
  media_file  = "${path.module}/earth.mp4"
}

resource "pexip_infinity_media_library_playlist" "test" {
  name        = "tf-test-playlist"
  description = "Test playlist"
  loop        = false
  shuffle     = false
}

resource "pexip_infinity_media_library_playlist_entry" "test" {
  entry_type = "MEDIA"
  media      = pexip_infinity_media_library_entry.test.id
  playlist   = pexip_infinity_media_library_playlist.test.id
  position   = 5
  playcount  = 3
}
