# infinity_webapp_branding

Manages webapp branding configuration with the Infinity service. Webapp branding allows customization of the user interface for different Pexip web applications including the management interface, admin interface, and client applications.

## Example Usage

### With Auto-Generated UUID

```hcl
resource "pexip_infinity_webapp_branding" "example" {
  name          = "Corporate Branding"
  description   = "Corporate branding for Pexip web applications"
  webapp_type   = "webapp1"
  is_default    = true
  branding_file = "/path/to/branding/package.zip"
  # uuid will be automatically generated if not provided
}
```

### With Custom UUID

```hcl
resource "pexip_infinity_webapp_branding" "example_custom_uuid" {
  name          = "Corporate Branding"
  description   = "Corporate branding for Pexip web applications"
  uuid          = "12345678-1234-1234-1234-123456789012"
  webapp_type   = "webapp2"
  is_default    = true
  branding_file = "/path/to/branding/package.zip"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the webapp branding configuration. This is used as the identifier. Maximum length: 100 characters. **Note:** Changing this after creation will force replacement of the resource.
* `uuid` - (Optional) The UUID for this branding configuration. If not provided, a UUID will be automatically generated. Must be a valid RFC 4122 UUID format (e.g., `550e8400-e29b-41d4-a716-446655440000`). **Note:** Changing this after creation will force replacement of the resource.
* `webapp_type` - (Required) The type of webapp this branding applies to. Valid values: `webapp1`, `webapp2`, `webapp3`. **Note:** Changing this after creation will force replacement of the resource.
* `is_default` - (Optional) Whether this is the default branding configuration for the webapp type. Defaults to computed value. **Note:** Changing this after creation will force replacement of the resource.
* `branding_file` - (Required) The path to the branding file (ZIP archive) to use for customization. **Note:** Changing this after creation will force replacement of the resource.
* `description` - (Optional) Description of the webapp branding configuration. Maximum length: 500 characters. **Note:** Changing this after creation will force replacement of the resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Resource URI for the webapp branding in Infinity.
* `last_updated` - Timestamp when this branding configuration was last updated.

## Import

Webapp branding configurations can be imported using their UUID:

```bash
terraform import pexip_infinity_webapp_branding.example "12345678-1234-1234-1234-123456789012"
```

## Webapp Types

The different webapp types that can be customized:

- `webapp1` - Webapp type 1
- `webapp2` - Webapp type 2
- `webapp3` - Webapp type 3