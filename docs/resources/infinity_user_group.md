# infinity_user_group

Manages a user group configuration with the Infinity service. User groups allow organizing users for management and permission purposes.

## Example Usage

```hcl
resource "pexip_infinity_user_group" "example" {
  name        = "Administrators"
  description = "System administrators group"
  users = [
    "/configuration/v1/end_user/1/",
    "/configuration/v1/end_user/2/"
  ]
  user_group_entity_mappings = [
    "/configuration/v1/user_group_entity_mapping/1/"
  ]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The unique name of the user group. Maximum length: 250 characters.
* `description` - (Optional) A description of the user group. Maximum length: 250 characters.
* `users` - (Optional) List of user resource URIs that belong to this group.
* `user_group_entity_mappings` - (Optional) List of user group entity mapping resource URIs associated with this group.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Resource URI for the user group in Infinity.
* `resource_id` - The resource integer identifier for the user group in Infinity.

## Import

User groups can be imported using their resource ID:

```bash
terraform import pexip_infinity_user_group.example 123
```