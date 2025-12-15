/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_pexip_streaming_credential" "pexip_streaming_credential-test" {
  kid        = "test-key-2" // Updated value
  public_key = <<-EOT
-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA2Z3qX2BTLS4e7sxZdZHW
gXuZMqkj9vP8vEf0h1r8vQOGLbB9xHNdLRfEr5lZwV7IXYZ1aNKdBFQvGYcLc3mG
zqzxJQMqLtWd8QDFL9fGdEDWHVPvM+sC8jmZBYdxmJx0zRvZQYTGhUL1C4F5eWJP
vCnH2wN8KjJD5gX7dFjPUhLMvKvPQqZRxXnMzHx8VGQyXYqB7PzQxNVvRJdxLwYn
8HfGJcV4wQd5xLxR9C6fZNYmC7bVjPqHxYwT2fC8HvGJmN6nPxLvF4wQJdVzB8hL
YqC9xPfZNvC4QwBdHmJxR7C5nLvP8QdFxBhN4C6vZ8LwYH2fC9BdNqLxC4F5vQwJ
wQIDAQAB
-----END PUBLIC KEY-----
EOT
}