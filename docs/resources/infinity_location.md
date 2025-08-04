---
page_title: "pexip_infinity_location Resource - terraform-provider-pexip"
subcategory: ""
description: |-
  Manages a Pexip Infinity location configuration.
---

# pexip_infinity_location (Resource)

Manages a Pexip Infinity location configuration. Locations are used to group nodes and define logical network boundaries within the Pexip Infinity deployment.

## Example Usage

### Basic Usage

```terraform
resource "pexip_infinity_location" "main" {
  name = "Main Location"
}
```

### With Description

```terraform
resource "pexip_infinity_location" "datacenter" {
  name        = "Data Center 1"
  description = "Primary data center location for production workloads"
}
```

### Multiple Locations

```terraform
resource "pexip_infinity_location" "locations" {
  for_each = toset(["Production", "Staging", "Development"])
  
  name        = each.value
  description = "${each.value} environment location"
}
```

## Schema

### Required

- `name` (String) - The unique name of the location. Maximum length: 250 characters.

### Optional

- `description` (String) - A description of the location. Maximum length: 250 characters.

### Read-Only

- `id` (String) - Resource URI for the location in Infinity.
- `resource_id` (Number) - The resource integer identifier for the location in Infinity.

## Import

Import is supported using the following syntax:

```shell
terraform import pexip_infinity_location.example 123
```

Where `123` is the numeric resource ID of the location.

## Usage Notes

### Location Planning
- Locations should reflect your network topology and geographical distribution
- Consider bandwidth and latency when assigning nodes to locations
- Use descriptive names that clearly identify the location's purpose

### Node Assignment
- Worker VMs and other nodes reference locations via the `system_location` attribute
- Ensure locations exist before creating nodes that reference them
- A location can contain multiple nodes but should not mix node types inappropriately

### Naming Conventions
- Use consistent naming across your deployment
- Consider including environment identifiers (prod, staging, dev)
- Location names are case-sensitive and must be unique

## Troubleshooting

### Common Issues

**Location Creation Fails**
- Verify the location name is unique within the deployment
- Check that the name doesn't exceed 250 characters
- Ensure proper authentication credentials

**Cannot Delete Location**
- Verify no nodes are assigned to the location
- Check for any references in other resources
- Remove all dependent resources before deletion

**Import Fails**
- Ensure you're using the numeric resource ID, not the name
- Verify the location exists in the Infinity deployment
- Check provider authentication credentials