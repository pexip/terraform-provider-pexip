resource "pexip_infinity_dns_server" "cloudflare-dns" {
  address     = "1.1.1.1"
  description = "Cloudflare DNS"
}
