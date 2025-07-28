resource "pexip_infinity_smtp_server" "smtp_server-test" {
  name                = "smtp_server-test"
  description         = "Updated Test SMTPServer"    // Updated description
  address             = "updated-server.example.com" // Updated address
  port                = 465                          // Updated port
  username            = "smtp_server-test"
  password            = "updated-value"       // Updated value
  from_email_address  = "updated@example.com" // Updated email
  connection_security = "ssl_tls"             // Updated value
}