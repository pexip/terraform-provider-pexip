package validators

import (
	"context"
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// URLValidator checks if the string is a valid URL.
type URLValidator struct {
	https bool
}

func (v URLValidator) Description(_ context.Context) string {
	return "Value must be a valid URL"
}

func (v URLValidator) MarkdownDescription(_ context.Context) string {
	return "Value must be a **valid URL**"
}

func (v URLValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	raw := req.ConfigValue.ValueString()
	u, err := url.ParseRequestURI(raw)
	if err != nil {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid URL",
			fmt.Sprintf("The value %q is not a valid URL: %s", raw, err),
		)
		return
	}

	if v.https && u.Scheme != "https" {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid URL Scheme",
			fmt.Sprintf("The value %q must be an HTTPS URL, but scheme is %q", raw, u.Scheme),
		)
	}
}

func URL(https bool) validator.String {
	return URLValidator{
		https: https,
	}
}
