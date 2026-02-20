/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package validators

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// TimeValidator checks if a string is a valid time in HH:MM:SS format.
type TimeValidator struct{}

func (v TimeValidator) Description(ctx context.Context) string {
	return "Value must be a valid time in HH:MM:SS format"
}

func (v TimeValidator) MarkdownDescription(ctx context.Context) string {
	return "Value must be a **valid time in HH:MM:SS format**"
}

func (v TimeValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	if value := req.ConfigValue.ValueString(); value != "" {
		if !isValidTime(value) {
			resp.Diagnostics.AddAttributeError(
				req.Path,
				"Invalid Time Format",
				fmt.Sprintf("Time '%s' is not a valid time. Expected format: HH:MM:SS (e.g., 09:30:00).", value),
			)
		}
	}
}

// isValidTime validates a time string in HH:MM:SS format
func isValidTime(timeStr string) bool {
	// Regex pattern for HH:MM:SS format
	// Hours: 00-23, Minutes: 00-59, Seconds: 00-59
	timeRegex := regexp.MustCompile(`^([01][0-9]|2[0-3]):([0-5][0-9]):([0-5][0-9])$`)
	return timeRegex.MatchString(timeStr)
}

// Time returns an instance of the time validator.
func Time() validator.String {
	return TimeValidator{}
}
