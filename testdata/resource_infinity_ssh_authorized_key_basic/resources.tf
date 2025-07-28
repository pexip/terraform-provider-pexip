resource "pexip_infinity_ssh_authorized_key" "ssh_authorized_key-test" {
  keytype = "ssh-rsa"
  key = "test-value"
  comment = "test-value"
}