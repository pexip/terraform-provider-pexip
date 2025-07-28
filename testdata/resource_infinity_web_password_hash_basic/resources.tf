resource "pexip_infinity_web_password_hash" "web_password_hash-test" {
  password = "test-value"
  salt = "abcdefghijkl"  // Exactly 12 characters
  rounds = 5000  // Minimum valid value
}