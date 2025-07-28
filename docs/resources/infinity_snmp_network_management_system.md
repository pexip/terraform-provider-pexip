# infinity_snmp_network_management_system

Manages an SNMP network management system with the Infinity service. SNMP network management systems receive SNMP traps and notifications from Pexip Infinity for monitoring and alerting purposes.

## Example Usage

```hcl
resource "pexip_infinity_snmp_network_management_system" "example" {
  name                 = "Primary NMS"
  description          = "Primary network management system for monitoring"
  address              = "nms.example.com"
  port                 = 162
  snmp_trap_community  = "public"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the SNMP network management system. Maximum length: 250 characters.
* `address` - (Required) The IP address or FQDN of the SNMP Network Management System. Maximum length: 255 characters.
* `port` - (Required) The port number for SNMP communications. Valid range: 1-65535.
* `snmp_trap_community` - (Required) The SNMP trap community string for authentication. This field is sensitive.
* `description` - (Optional) Description of the SNMP network management system. Maximum length: 500 characters.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Resource URI for the SNMP network management system in Infinity.
* `resource_id` - The resource integer identifier for the SNMP network management system in Infinity.

## Import

SNMP network management systems can be imported using their resource ID:

```bash
terraform import pexip_infinity_snmp_network_management_system.example 123
```

## Security Notes

- The `snmp_trap_community` field is marked as sensitive and will not be displayed in Terraform output.
- Use strong community strings to secure SNMP communications.
- Consider using SNMPv3 with encryption when supported by your monitoring infrastructure.