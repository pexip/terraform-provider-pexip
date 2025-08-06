package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

func TestInfinityManagerConfigModel(t *testing.T) {
	config := InfinityManagerConfigModel{
		ID:                  types.StringValue("test-id"),
		Hostname:            types.StringValue("manager-1"),
		Domain:              types.StringValue("example.com"),
		IP:                  types.StringValue("10.5.6.7"),
		Mask:                types.StringValue("255.255.255.0"),
		GW:                  types.StringValue("10.5.6.1"),
		DNS:                 types.StringValue("1.1.1.1"),
		NTP:                 types.StringValue("time.example.com"),
		User:                types.StringValue("admin"),
		Pass:                types.StringValue("password123"),
		AdminPassword:       types.StringValue("adminpass123"),
		ErrorReports:        types.BoolValue(true),
		EnableAnalytics:     types.BoolValue(true),
		ContactEmailAddress: types.StringValue("test@example.com"),
	}

	diags := config.validate(t.Context())
	require.False(t, diags.HasError(), "Expected no validation errors, got: %s", diags.Errors())
}
