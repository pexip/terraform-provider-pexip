/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_media_library_playlist" "media_library_playlist-test" {
  name        = "media_library_playlist-test"
  description = "Test MediaLibraryPlaylist"
  loop        = true
  shuffle     = true
}