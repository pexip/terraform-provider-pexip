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

// NetmaskValidator checks if a string is a valid netmask.
type NetmaskValidator struct{}

func (v NetmaskValidator) Description(ctx context.Context) string {
	return "Value must be a valid IPv4 netmask"
}

func (v NetmaskValidator) MarkdownDescription(ctx context.Context) string {
	return "Value must be a **valid IPv4 netmask**"
}

func (v NetmaskValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	if value := req.ConfigValue.ValueString(); value != "" {
		// Parse as IP address first
		ip := net.ParseIP(value)
		if ip == nil {
			resp.Diagnostics.AddAttributeError(
				req.Path,
				"Invalid Netmask",
				fmt.Sprintf("Netmask '%s' is not a valid IP address format.", value),
			)
			return
		}

		// Convert to IPv4 if possible
		ipv4 := ip.To4()
		if ipv4 == nil {
			resp.Diagnostics.AddAttributeError(
				req.Path,
				"Invalid Netmask",
				fmt.Sprintf("Netmask '%s' is not a valid IPv4 address.", value),
			)
			return
		}

		// Check if it's a valid netmask by converting to uint32 and checking bit pattern
		mask := uint32(ipv4[0])<<24 | uint32(ipv4[1])<<16 | uint32(ipv4[2])<<8 | uint32(ipv4[3])

		// A valid netmask has all 1s followed by all 0s
		// Find the first 0 bit
		foundZero := false
		for i := 31; i >= 0; i-- {
			bit := (mask >> i) & 1
			if bit == 0 {
				foundZero = true
			} else if foundZero {
				// Found a 1 bit after a 0 bit - invalid netmask
				resp.Diagnostics.AddAttributeError(
					req.Path,
					"Invalid Netmask",
					fmt.Sprintf("Netmask '%s' is not a valid subnet mask (must be contiguous 1s followed by 0s).", value),
				)
				return
			}
		}
	}
}

// Netmask returns an instance of the netmask validator.
func Netmask() validator.String {
	return NetmaskValidator{}
}
