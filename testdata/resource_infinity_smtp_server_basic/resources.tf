resource "pexip_infinity_smtp_server" "smtp_server-test" {
  name = "smtp_server-test"
  description = "Test SMTPServer"
  address = "test-server.example.com"
  port = 587
  username = "smtp_server-test"
  password = "test-value"
  from_email_address = "test@example.com"
  connection_security = "starttls"
}