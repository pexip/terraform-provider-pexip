/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_media_library_playlist_entry" "test" {
  entry_type = "MEDIA"
  media      = "/api/admin/configuration/v1/media_library_entry/1/"
  playlist   = "/api/admin/configuration/v1/media_library_playlist/2/"
  position   = 5
  playcount  = 3
}
