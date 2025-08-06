/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package validators

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Helper function to test string validators
func testStringValidator(t *testing.T, v validator.String, value string, expectError bool) {
	t.Helper()

	ctx := context.Background()

	req := validator.StringRequest{
		Path:        path.Root("test"),
		ConfigValue: types.StringValue(value),
	}

	resp := &validator.StringResponse{}

	v.ValidateString(ctx, req, resp)

	hasErrors := resp.Diagnostics.HasError()
	if hasErrors != expectError {
		t.Errorf("Expected error: %v, got error: %v, errors: %v", expectError, hasErrors, resp.Diagnostics.Errors())
	}
}

// Helper function to test null and unknown values
func testNullAndUnknown(t *testing.T, v validator.String) {
	t.Helper()

	ctx := context.Background()

	// Test null value
	req := validator.StringRequest{
		Path:        path.Root("test"),
		ConfigValue: types.StringNull(),
	}
	resp := &validator.StringResponse{}
	v.ValidateString(ctx, req, resp)

	if resp.Diagnostics.HasError() {
		t.Error("Expected no error for null value")
	}

	// Test unknown value
	req = validator.StringRequest{
		Path:        path.Root("test"),
		ConfigValue: types.StringUnknown(),
	}
	resp = &validator.StringResponse{}
	v.ValidateString(ctx, req, resp)

	if resp.Diagnostics.HasError() {
		t.Error("Expected no error for unknown value")
	}
}

func TestEmailValidator(t *testing.T) {
	t.Parallel()
	v := Email()

	tests := []struct {
		name        string
		value       string
		expectError bool
	}{
		{"valid simple email", "test@example.com", false},
		{"valid email with subdomain", "user@mail.example.com", false},
		{"valid email with plus", "user+tag@example.com", false},
		{"valid email with dash", "user-name@example.com", false},
		{"valid email with underscore", "user_name@example.com", false},
		{"valid email with dots", "first.last@example.com", false},
		{"valid email with numbers", "user123@example123.com", false},
		{"invalid email - no @", "userexample.com", true},
		{"invalid email - multiple @", "user@@example.com", true},
		{"invalid email - no domain", "user@", true},
		{"invalid email - no local part", "@example.com", true},
		{"invalid email - starts with dot", ".user@example.com", true},
		{"invalid email - ends with dot", "user.@example.com", true},
		{"invalid email - consecutive dots", "user..name@example.com", true},
		{"invalid email - invalid domain", "user@invalid", true},
		{"empty email", "", false}, // Empty values are allowed (handled by required validation)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testStringValidator(t, v, tt.value, tt.expectError)
		})
	}

	t.Run("null and unknown", func(t *testing.T) {
		testNullAndUnknown(t, v)
	})
}

func TestIPAddressValidator(t *testing.T) {
	t.Parallel()
	v := IPAddress()

	tests := []struct {
		name        string
		value       string
		expectError bool
	}{
		{"valid IPv4 address", "192.168.1.1", false},
		{"valid IPv4 address - localhost", "127.0.0.1", false},
		{"valid IPv4 address - zero", "0.0.0.0", false},
		{"valid IPv4 address - broadcast", "255.255.255.255", false},
		{"valid IPv6 address - full", "2001:0db8:85a3:0000:0000:8a2e:0370:7334", false},
		{"valid IPv6 address - compressed", "2001:db8:85a3::8a2e:370:7334", false},
		{"valid IPv6 address - localhost", "::1", false},
		{"valid IPv6 address - zero", "::", false},
		{"invalid IPv4 - too many octets", "192.168.1.1.1", true},
		{"invalid IPv4 - octet out of range", "192.168.1.256", true},
		{"invalid IPv4 - missing octets", "192.168.1", true},
		{"invalid IPv4 - non-numeric", "192.168.1.abc", true},
		{"invalid IPv6 - too many groups", "2001:0db8:85a3:0000:0000:8a2e:0370:7334:1234", true},
		{"invalid IPv6 - invalid characters", "2001:0db8:85a3::8a2e:370g:7334", true},
		{"invalid - not an IP", "not.an.ip.address", true},
		{"invalid - hostname", "example.com", true},
		{"empty IP", "", false}, // Empty values are allowed (handled by required validation)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testStringValidator(t, v, tt.value, tt.expectError)
		})
	}

	t.Run("null and unknown", func(t *testing.T) {
		testNullAndUnknown(t, v)
	})
}

func TestDomainValidator(t *testing.T) {
	t.Parallel()
	v := Domain()

	tests := []struct {
		name        string
		value       string
		expectError bool
	}{
		{"valid domain", "example.com", false},
		{"valid subdomain", "sub.example.com", false},
		{"valid domain with multiple subdomains", "mail.sub.example.com", false},
		{"valid domain with FQDN", "example.com.", false},
		{"valid domain with numbers", "test123.example.com", false},
		{"valid domain with hyphens", "test-domain.example.com", false},
		{"valid international domain", "example.co.uk", false},
		{"valid long TLD", "example.museum", false},
		{"invalid domain - no dot", "localhost", true},
		{"invalid domain - starts with dot", ".example.com", true},
		{"invalid domain - consecutive dots", "example..com", true},
		{"invalid domain - starts with hyphen", "-example.com", true},
		{"invalid domain - ends with hyphen", "example-.com", true},
		{"invalid domain - too short", "a.b", false}, // Actually valid according to RFC
		{"invalid domain - numeric TLD", "example.123", true},
		{"invalid domain - special characters", "example@.com", true},
		{"invalid domain - underscore", "example_.com", true},
		{"empty domain", "", false}, // Empty values are allowed (handled by required validation)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testStringValidator(t, v, tt.value, tt.expectError)
		})
	}

	t.Run("null and unknown", func(t *testing.T) {
		testNullAndUnknown(t, v)
	})
}

func TestURLValidator(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		value       string
		httpsOnly   bool
		expectError bool
	}{
		{"valid HTTP URL", "http://example.com", false, false},
		{"valid HTTPS URL", "https://example.com", false, false},
		{"valid HTTPS URL with HTTPS required", "https://example.com", true, false},
		{"valid URL with path", "https://example.com/path/to/resource", false, false},
		{"valid URL with query parameters", "https://example.com/search?q=test&type=web", false, false},
		{"valid URL with fragment", "https://example.com/page#section", false, false},
		{"valid URL with port", "https://example.com:8080", false, false},
		{"valid URL with subdomain", "https://api.example.com", false, false},
		{"valid FTP URL", "ftp://files.example.com", false, false},
		{"HTTP URL with HTTPS required", "http://example.com", true, true},
		{"FTP URL with HTTPS required", "ftp://files.example.com", true, true},
		{"invalid URL - no scheme", "example.com", false, true},
		{"invalid URL - malformed", "http://", false, false}, // url.ParseRequestURI allows this
		{"empty URL", "", false, true},                       // Empty URLs are invalid
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := URL(tt.httpsOnly)
			testStringValidator(t, v, tt.value, tt.expectError)
		})
	}

	t.Run("null and unknown", func(t *testing.T) {
		v := URL(false)
		testNullAndUnknown(t, v)
	})
}

func TestNetmaskValidator(t *testing.T) {
	t.Parallel()
	v := Netmask()

	tests := []struct {
		name        string
		value       string
		expectError bool
	}{
		{"valid netmask /24", "255.255.255.0", false},
		{"valid netmask /16", "255.255.0.0", false},
		{"valid netmask /8", "255.0.0.0", false},
		{"valid netmask /30", "255.255.255.252", false},
		{"valid netmask /31", "255.255.255.254", false},
		{"valid netmask /32", "255.255.255.255", false},
		{"valid netmask /0", "0.0.0.0", false},
		{"valid netmask /28", "255.255.255.240", false},
		{"valid netmask /20", "255.255.240.0", false},
		{"valid netmask /12", "255.240.0.0", false},
		{"invalid netmask - not contiguous", "255.255.255.129", true},
		{"invalid netmask - gaps", "255.255.0.255", true},
		{"invalid netmask - reversed pattern", "0.255.255.255", true},
		{"invalid netmask - random pattern", "255.128.255.0", true},
		{"invalid netmask - not IP format", "255.255.255", true},
		{"invalid netmask - out of range", "255.255.255.256", true},
		{"invalid netmask - IPv6", "ffff:ffff:ffff:ffff::", true},
		{"invalid netmask - non-numeric", "255.255.255.abc", true},
		{"invalid netmask - hostname", "example.com", true},
		{"empty netmask", "", false}, // Empty values are allowed (handled by required validation)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testStringValidator(t, v, tt.value, tt.expectError)
		})
	}

	t.Run("null and unknown", func(t *testing.T) {
		testNullAndUnknown(t, v)
	})
}
