/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package validators

import (
	"context"
	"fmt"
	"net"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// IPAddressValidator checks if a string is a valid IP address.
type IPAddressValidator struct{}

func (v IPAddressValidator) Description(ctx context.Context) string {
	return "Value must be a valid IPv4 or IPv6 address"
}

func (v IPAddressValidator) MarkdownDescription(ctx context.Context) string {
	return "Value must be a **valid IPv4 or IPv6 address**"
}

func (v IPAddressValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	if value := req.ConfigValue.ValueString(); value != "" {
		ip := net.ParseIP(value)
		if ip == nil {
			resp.Diagnostics.AddAttributeError(
				req.Path,
				"Invalid IP Address",
				fmt.Sprintf("IP '%s' is not a valid IP address.", value),
			)
		}
	}
}

// IPAddress returns an instance of the IP address validator.
func IPAddress() validator.String {
	return IPAddressValidator{}
}
