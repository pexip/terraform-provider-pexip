# infinity_role

Manages a role configuration with the Infinity service. Roles define sets of permissions that can be assigned to users or applications.

## Example Usage

```hcl
resource "pexip_infinity_role" "example" {
  name = "Conference Manager"
  permissions = [
    "conference.create",
    "conference.read", 
    "conference.update",
    "conference.delete"
  ]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The unique name of the role. Maximum length: 250 characters.
* `permissions` - (Optional) List of permissions assigned to this role.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Resource URI for the role in Infinity.
* `resource_id` - The resource integer identifier for the role in Infinity.

## Import

Roles can be imported using their resource ID:

```bash
terraform import pexip_infinity_role.example 123
```