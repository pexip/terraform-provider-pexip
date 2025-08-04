# infinity_webapp_branding

Manages webapp branding configuration with the Infinity service. Webapp branding allows customization of the user interface for different Pexip web applications including the management interface, admin interface, and client applications.

## Example Usage

```hcl
resource "pexip_infinity_webapp_branding" "example" {
  name          = "Corporate Branding"
  description   = "Corporate branding for Pexip web applications"
  uuid          = "12345678-1234-1234-1234-123456789012"
  webapp_type   = "pexapp"
  is_default    = true
  branding_file = "/path/to/branding/package.zip"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the webapp branding configuration. This is used as the identifier. Maximum length: 100 characters.
* `uuid` - (Required) The UUID for this branding configuration.
* `webapp_type` - (Required) The type of webapp this branding applies to. Valid values: `pexapp`, `management`, `admin`.
* `is_default` - (Required) Whether this is the default branding configuration for the webapp type.
* `branding_file` - (Required) The path or identifier for the branding file to use for customization.
* `description` - (Optional) Description of the webapp branding configuration. Maximum length: 500 characters.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Resource URI for the webapp branding in Infinity.
* `last_updated` - Timestamp when this branding configuration was last updated.

## Import

Webapp branding configurations can be imported using their name:

```bash
terraform import pexip_infinity_webapp_branding.example "Corporate Branding"
```

## Webapp Types

The different webapp types that can be customized:

- `pexapp` - The client web application used by end users
- `management` - The management interface for administrators
- `admin` - The administrative interface for system configuration