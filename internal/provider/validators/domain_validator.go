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
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// DomainValidator checks if a string is a valid domain name.
type DomainValidator struct{}

func (v DomainValidator) Description(ctx context.Context) string {
	return "Value must be a valid domain name"
}

func (v DomainValidator) MarkdownDescription(ctx context.Context) string {
	return "Value must be a **valid domain name**"
}

func (v DomainValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	if value := req.ConfigValue.ValueString(); value != "" {
		if !isValidDomain(value) {
			resp.Diagnostics.AddAttributeError(
				req.Path,
				"Invalid Domain Name",
				fmt.Sprintf("Domain '%s' is not a valid domain name.", value),
			)
		}
	}
}

// isValidDomain validates a domain name according to RFC standards
func isValidDomain(domain string) bool {
	// Basic length check
	if len(domain) < 3 || len(domain) > 253 {
		return false
	}

	// Must contain at least one dot
	if !strings.Contains(domain, ".") {
		return false
	}

	// Remove trailing dot if present (FQDN)
	domain = strings.TrimSuffix(domain, ".")

	// Split into labels
	labels := strings.Split(domain, ".")
	if len(labels) < 2 {
		return false
	}

	// Validate each label
	labelRegex := regexp.MustCompile(`^[a-zA-Z0-9]([a-zA-Z0-9-]*[a-zA-Z0-9])?$`)

	for _, label := range labels {
		// Label length check
		if label == "" || len(label) > 63 {
			return false
		}

		// Label format check
		if !labelRegex.MatchString(label) {
			return false
		}
	}

	// TLD (last label) must contain at least one letter
	tld := labels[len(labels)-1]
	hasLetter := false
	for _, char := range tld {
		if (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') {
			hasLetter = true
			break
		}
	}

	return hasLetter
}

// Domain returns an instance of the domain validator.
func Domain() validator.String {
	return DomainValidator{}
}
