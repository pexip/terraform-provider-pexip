# infinity_gateway_routing_rule

Manages a gateway routing rule configuration with the Infinity service. Gateway routing rules define how incoming calls are routed and transformed when passing through Pexip Infinity gateways.

## Example Usage

```hcl
resource "pexip_infinity_gateway_routing_rule" "example" {
  name               = "SIP Gateway Rule"
  description        = "Route calls to external SIP gateway"
  priority           = 100
  enable             = true
  match_string       = "^\\+1(.*)$"
  replace_string     = "$1"
  called_device_type = "gateway"
  outgoing_protocol  = "sip"
  call_type          = "video"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The unique name of the gateway routing rule. Maximum length: 250 characters.
* `priority` - (Required) The priority of the gateway routing rule (lower numbers have higher priority).
* `match_string` - (Required) Regular expression pattern to match incoming calls. Maximum length: 250 characters.
* `description` - (Optional) A description of the gateway routing rule. Maximum length: 250 characters.
* `enable` - (Optional) Whether the gateway routing rule is enabled. Defaults to true.
* `replace_string` - (Optional) Pattern for outgoing call transformation. Maximum length: 250 characters.
* `called_device_type` - (Optional) Type of called device. Valid choices: `unknown`, `conference`, `gateway`, `ip_pbx`.
* `outgoing_protocol` - (Optional) Outgoing protocol. Valid choices: `sip`, `h323`.
* `call_type` - (Optional) Call type. Valid choices: `audio`, `video`.
* `ivr_theme` - (Optional) Reference to IVR theme resource URI.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Resource URI for the gateway routing rule in Infinity.
* `resource_id` - The resource integer identifier for the gateway routing rule in Infinity.

## Import

Gateway routing rules can be imported using their resource ID:

```bash
terraform import pexip_infinity_gateway_routing_rule.example 123
```