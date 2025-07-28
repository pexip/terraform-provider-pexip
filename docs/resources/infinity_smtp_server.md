# infinity_smtp_server

Manages an SMTP server with the Infinity service. SMTP servers are used for sending email notifications and alerts from Pexip Infinity.

## Example Usage

```hcl
resource "pexip_infinity_smtp_server" "example" {
  name                = "Corporate SMTP Server"
  description         = "SMTP server for email notifications"
  address             = "smtp.example.com"
  port                = 587
  username            = "pexip@example.com"
  password            = "secure-password"
  from_email_address  = "noreply@example.com"
  connection_security = "starttls"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the SMTP server. Maximum length: 250 characters.
* `address` - (Required) The IP address or hostname of the SMTP server.
* `port` - (Required) The port number for SMTP communications. Valid range: 1-65535.
* `from_email_address` - (Required) The from email address used when sending emails through this SMTP server.
* `connection_security` - (Required) Connection security method for SMTP. Valid values: `none`, `starttls`, `ssl_tls`.
* `description` - (Optional) Description of the SMTP server. Maximum length: 500 characters.
* `username` - (Optional) Username for SMTP authentication.
* `password` - (Optional) Password for SMTP authentication. This field is sensitive.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Resource URI for the SMTP server in Infinity.
* `resource_id` - The resource integer identifier for the SMTP server in Infinity.

## Import

SMTP servers can be imported using their resource ID:

```bash
terraform import pexip_infinity_smtp_server.example 123
```

## Security Notes

- The `password` field is marked as sensitive and will not be displayed in Terraform output.
- Use `starttls` or `ssl_tls` for `connection_security` when possible to ensure encrypted email transmission.
- Store credentials securely and consider using environment variables or a secrets management system.