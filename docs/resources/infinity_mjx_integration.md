# infinity_mjx_integration

Manages an MJX (Meeting Join Experience) integration with the Infinity service. MJX integrations provide calendar integration and meeting join experiences for various platforms including Exchange, Google Calendar, and Webex.

## Example Usage

```hcl
resource "pexip_infinity_mjx_integration" "example" {
  name                          = "Corporate MJX Integration"
  description                   = "MJX integration for corporate calendar systems"
  display_upcoming_meetings     = 5
  enable_non_video_meetings     = true
  enable_private_meetings       = false
  end_buffer                    = 5
  start_buffer                  = 10
  ep_username                   = "pexip-service@example.com"
  ep_password                   = "secure-password"
  ep_use_https                  = true
  ep_verify_certificate         = true
  exchange_deployment           = "/configuration/v1/ms_exchange_connector/1/"
  process_alias_private_meetings = false
  replace_empty_subject         = true
  replace_subject_type          = "template"
  replace_subject_template      = "Meeting in {{location}}"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the MJX integration.
* `description` - (Optional) Description of the MJX integration.
* `display_upcoming_meetings` - (Optional) Number of upcoming meetings to display.
* `enable_non_video_meetings` - (Optional) Whether to enable non-video meetings.
* `enable_private_meetings` - (Optional) Whether to enable private meetings.
* `end_buffer` - (Optional) End buffer time in minutes.
* `start_buffer` - (Optional) Start buffer time in minutes.
* `ep_username` - (Optional) Endpoint username for authentication.
* `ep_password` - (Optional) Endpoint password for authentication. This field is sensitive.
* `ep_use_https` - (Optional) Whether to use HTTPS for endpoint communication.
* `ep_verify_certificate` - (Optional) Whether to verify SSL certificates.
* `exchange_deployment` - (Optional) Reference to Exchange deployment resource URI.
* `google_deployment` - (Optional) Reference to Google deployment resource URI.
* `graph_deployment` - (Optional) Reference to Microsoft Graph deployment resource URI.
* `process_alias_private_meetings` - (Optional) Whether to process alias private meetings.
* `replace_empty_subject` - (Optional) Whether to replace empty meeting subjects.
* `replace_subject_type` - (Optional) Type of subject replacement.
* `replace_subject_template` - (Optional) Template for subject replacement.
* `use_webex` - (Optional) Whether to use Webex integration.
* `webex_api_domain` - (Optional) Webex API domain.
* `webex_client_id` - (Optional) Webex client ID.
* `webex_client_secret` - (Optional) Webex client secret. This field is sensitive.
* `webex_oauth_state` - (Optional) Webex OAuth state.
* `webex_redirect_uri` - (Optional) Webex redirect URI.
* `webex_refresh_token` - (Optional) Webex refresh token. This field is sensitive.
* `endpoint_groups` - (Optional) List of endpoint groups associated with this integration.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Resource URI for the MJX integration in Infinity.
* `resource_id` - The resource integer identifier for the MJX integration in Infinity.

## Import

MJX integrations can be imported using their resource ID:

```bash
terraform import pexip_infinity_mjx_integration.example 123
```