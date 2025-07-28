# infinity_azure_tenant

Manages a Microsoft Teams Azure tenant with the Infinity service. Azure tenants are used to configure Microsoft Teams integration and allow Pexip Infinity to connect with specific Microsoft tenants for cloud-based collaboration.

## Example Usage

```hcl
resource "pexip_infinity_azure_tenant" "example" {
  name        = "Contoso Corporation"
  description = "Azure tenant for Contoso Corporation Teams integration"
  tenant_id   = "12345678-1234-1234-1234-123456789012"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Azure tenant. Maximum length: 100 characters.
* `tenant_id` - (Required) The Microsoft Azure tenant ID (GUID) for this tenant configuration.
* `description` - (Optional) Description of the Azure tenant. Maximum length: 500 characters.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Resource URI for the Azure tenant in Infinity.
* `resource_id` - The resource integer identifier for the Azure tenant in Infinity.

## Import

Azure tenants can be imported using their resource ID:

```bash
terraform import pexip_infinity_azure_tenant.example 123
```