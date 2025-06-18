package validators

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// EmailValidator checks if a string is a valid email address.
type EmailValidator struct{}

func (v EmailValidator) Description(ctx context.Context) string {
	return "Value must be a valid email address"
}

func (v EmailValidator) MarkdownDescription(ctx context.Context) string {
	return "Value must be a **valid email address**"
}

func (v EmailValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	if value := req.ConfigValue.ValueString(); value != "" {
		if !isValidEmail(value) {
			resp.Diagnostics.AddAttributeError(
				req.Path,
				"Invalid Email Address",
				fmt.Sprintf("Email '%s' is not a valid email address.", value),
			)
		}
	}
}

// isValidEmail validates an email address according to RFC 5322 standards
func isValidEmail(email string) bool {
	// Basic length check
	if len(email) < 3 || len(email) > 254 {
		return false
	}

	// Must contain exactly one @ symbol
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}

	localPart := parts[0]
	domainPart := parts[1]

	// Validate local part
	if !isValidLocalPart(localPart) {
		return false
	}

	// Validate domain part using existing domain validator logic
	return isValidEmailDomain(domainPart)
}

// isValidEmailDomain is a copy of isValidDomain for email validation
// to avoid circular dependencies
func isValidEmailDomain(domain string) bool {
	// Basic length check
	if len(domain) < 3 || len(domain) > 253 {
		return false
	}

	// Must contain at least one dot
	if !strings.Contains(domain, ".") {
		return false
	}

	// Remove trailing dot if present (FQDN)
	if strings.HasSuffix(domain, ".") {
		domain = domain[:len(domain)-1]
	}

	// Split into labels
	labels := strings.Split(domain, ".")
	if len(labels) < 2 {
		return false
	}

	// Validate each label
	labelRegex := regexp.MustCompile(`^[a-zA-Z0-9]([a-zA-Z0-9-]*[a-zA-Z0-9])?$`)

	for _, label := range labels {
		// Label length check
		if len(label) == 0 || len(label) > 63 {
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

// isValidLocalPart validates the local part (before @) of an email address
func isValidLocalPart(localPart string) bool {
	// Length check (RFC 5321 limits to 64 characters)
	if len(localPart) == 0 || len(localPart) > 64 {
		return false
	}

	// Cannot start or end with a dot
	if strings.HasPrefix(localPart, ".") || strings.HasSuffix(localPart, ".") {
		return false
	}

	// Cannot have consecutive dots
	if strings.Contains(localPart, "..") {
		return false
	}

	// Check for valid characters
	// Allow alphanumeric, dots, hyphens, underscores, and plus signs
	validChars := regexp.MustCompile(`^[a-zA-Z0-9._+-]+$`)
	return validChars.MatchString(localPart)
}

// Email returns an instance of the email validator.
func Email() validator.String {
	return EmailValidator{}
}
