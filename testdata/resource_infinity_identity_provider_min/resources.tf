/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "random_uuid4" "example" {
}

resource "pexip_infinity_identity_provider" "test" {
  name = "tf-test Identity Provider min"
  uuid = random_uuid4.example.result
  assertion_consumer_service_url = "https://test.com/samlconsumer/${random_uuid4.example.result}"
}