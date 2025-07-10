resource "tls_private_key" "manager_private_key" {
  algorithm = "RSA"
  rsa_bits  = 2048
}

resource "tls_self_signed_cert" "manager_cert" {
  private_key_pem = tls_private_key.manager_private_key.private_key_pem

  subject {
    common_name  = local.hostname
    organization = "Pexip AS"
  }

  validity_period_hours = 8760  # 1 year
  allowed_uses = [
    "key_encipherment",
    "digital_signature",
    "server_auth",
  ]

  dns_names = [local.hostname]
}