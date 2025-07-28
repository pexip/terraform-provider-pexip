# infinity_recurring_conference

Manages a recurring conference configuration with the Infinity service. Recurring conferences are used for scheduled conference series that repeat over time.

## Example Usage

```hcl
resource "pexip_infinity_recurring_conference" "example" {
  conference      = "/configuration/v1/conference/1/"
  current_index   = 0
  ews_item_id     = "AAMkAGVm...ExampleRecurringEwsItemId"
  is_depleted     = false
  subject         = "Weekly Team Standup"
  scheduled_alias = "/configuration/v1/scheduled_alias/1/"
}
```

## Argument Reference

The following arguments are supported:

* `conference` - (Required) The conference identifier or URI associated with this recurring conference.
* `current_index` - (Required) The current index of the recurring conference series (minimum: 0).
* `ews_item_id` - (Required) The Exchange Web Services (EWS) item identifier for this recurring conference.
* `is_depleted` - (Required) Whether the recurring conference series is depleted (no more occurrences).
* `subject` - (Optional) The subject or title of the recurring conference. Maximum length: 500 characters.
* `scheduled_alias` - (Optional) The scheduled alias for the recurring conference.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Resource URI for the recurring conference in Infinity.
* `resource_id` - The resource integer identifier for the recurring conference in Infinity.

## Import

Recurring conferences can be imported using their resource ID:

```bash
terraform import pexip_infinity_recurring_conference.example 123
```