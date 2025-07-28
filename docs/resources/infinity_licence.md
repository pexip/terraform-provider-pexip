# infinity_licence

Manages a licence configuration with the Infinity service. This resource activates licences using entitlement IDs.

## Example Usage

```hcl
resource "pexip_infinity_licence" "example" {
  entitlement_id = "ENT123456-7890-ABCD-EF01-234567890123"
}
```

## Argument Reference

The following arguments are supported:

* `entitlement_id` - (Required) The entitlement ID for the licence activation.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Resource URI for the licence in Infinity.
* `fulfillment_id` - The fulfillment ID of the licence (used as the unique identifier).
* `fulfillment_type` - The fulfillment type of the licence.
* `product_id` - The product ID associated with the licence.
* `license_type` - The type of the licence.
* `features` - The features enabled by this licence.
* `concurrent` - Number of concurrent sessions allowed.
* `concurrent_overdraft` - Number of concurrent sessions allowed in overdraft.
* `activatable` - Number of activatable licenses.
* `activatable_overdraft` - Number of activatable licenses allowed in overdraft.
* `hybrid` - Number of hybrid licenses.
* `hybrid_overdraft` - Number of hybrid licenses allowed in overdraft.
* `start_date` - The start date of the licence validity period.
* `expiration_date` - The expiration date of the licence.
* `status` - The current status of the licence.
* `trust_flags` - Trust flags for the licence.
* `repair` - Repair flag for the licence.
* `server_chain` - The server chain for the licence.
* `offline_mode` - Whether the licence should be activated in offline mode.

## Import

Licences can be imported using their fulfillment ID:

```bash
terraform import pexip_infinity_licence.example FUL123456-7890-ABCD-EF01-234567890123
```

## Notes

- Licence resources are immutable once activated. Updates are not supported.
- To change licence settings, you must delete and recreate the resource.
- The licence will be deactivated when the resource is destroyed.