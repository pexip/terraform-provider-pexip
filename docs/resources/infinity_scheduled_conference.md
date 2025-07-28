# infinity_scheduled_conference

Manages a scheduled conference configuration with the Infinity service. Scheduled conferences allow you to configure time-based conference reservations with specific start and end times.

## Example Usage

```hcl
resource "pexip_infinity_scheduled_conference" "example" {
  conference   = "/configuration/v1/conference/1/"
  start_time   = "2024-12-01T10:00:00Z"
  end_time     = "2024-12-01T11:00:00Z"
  subject      = "Weekly Team Meeting"
  ews_item_id  = "AAMkAGVm...ExampleEwsItemId"
  ews_item_uid = "040000008200E00074C5B7101A82E00800000000"
}
```

## Argument Reference

The following arguments are supported:

* `conference` - (Required) The conference URI or reference. Maximum length: 250 characters.
* `start_time` - (Required) The start time of the scheduled conference in ISO 8601 format (e.g., '2024-01-01T10:00:00Z').
* `end_time` - (Required) The end time of the scheduled conference in ISO 8601 format (e.g., '2024-01-01T11:00:00Z').
* `ews_item_id` - (Required) The Exchange Web Services (EWS) item ID for the conference.
* `subject` - (Optional) The subject of the scheduled conference. Maximum length: 250 characters.
* `ews_item_uid` - (Optional) The Exchange Web Services (EWS) item UID for the conference. Maximum length: 250 characters.
* `recurring_conference` - (Optional) Reference to recurring conference resource URI.
* `scheduled_alias` - (Optional) Reference to scheduled alias resource URI.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Resource URI for the scheduled conference in Infinity.
* `resource_id` - The resource integer identifier for the scheduled conference in Infinity.

## Import

Scheduled conferences can be imported using their resource ID:

```bash
terraform import pexip_infinity_scheduled_conference.example 123
```

## Time Format

The `start_time` and `end_time` fields must be specified in ISO 8601 format with timezone information. Examples:
- `2024-01-01T10:00:00Z` (UTC)
- `2024-01-01T10:00:00+01:00` (CET)
- `2024-01-01T10:00:00-05:00` (EST)