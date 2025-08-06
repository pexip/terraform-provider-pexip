/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/mail"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/pexip/terraform-provider-pexip/internal/helpers"
)

type InfinityManagerConfigModel struct {
	ID                           types.String `tfsdk:"id"`
	Hostname                     types.String `tfsdk:"hostname"`
	Domain                       types.String `tfsdk:"domain"`
	IP                           types.String `tfsdk:"ip"`
	Mask                         types.String `tfsdk:"mask"`
	GW                           types.String `tfsdk:"gw"`
	DNS                          types.String `tfsdk:"dns"`
	NTP                          types.String `tfsdk:"ntp"`
	User                         types.String `tfsdk:"user"`
	Pass                         types.String `tfsdk:"pass"`
	AdminPassword                types.String `tfsdk:"admin_password"`
	ErrorReports                 types.Bool   `tfsdk:"error_reports"`
	EnableAnalytics              types.Bool   `tfsdk:"enable_analytics"`
	ContactEmailAddress          types.String `tfsdk:"contact_email_address"`
	Rendered                     types.String `tfsdk:"rendered"`
	RenderedManagementNodeConfig types.String `tfsdk:"management_node_config"`
}

type outerInfinityManagerConfig struct {
	ManagementNodeConfig innerInfinityBootstrapConfig `json:"management_node_config"`
}

type innerInfinityBootstrapConfig struct {
	Hostname            string `json:"hostname"`
	Domain              string `json:"domain"`
	IP                  string `json:"ip"`
	Mask                string `json:"mask"`
	GW                  string `json:"gw"`
	DNS                 string `json:"dns"`
	NTP                 string `json:"ntp"`
	User                string `json:"user"`
	Pass                string `json:"pass"`
	AdminPassword       string `json:"admin_password"`
	ErrorReports        bool   `json:"error_reports"`
	EnableAnalytics     bool   `json:"enable_analytics"`
	ContactEmailAddress string `json:"contact_email_address"`
}

func (c *InfinityManagerConfigModel) toOuterConfig() outerInfinityManagerConfig {
	return outerInfinityManagerConfig{
		ManagementNodeConfig: innerInfinityBootstrapConfig{
			Hostname:            c.Hostname.ValueString(),
			Domain:              c.Domain.ValueString(),
			IP:                  c.IP.ValueString(),
			Mask:                c.Mask.ValueString(),
			GW:                  c.GW.ValueString(),
			DNS:                 c.DNS.ValueString(),
			NTP:                 c.NTP.ValueString(),
			User:                c.User.ValueString(),
			Pass:                c.Pass.ValueString(),
			AdminPassword:       c.AdminPassword.ValueString(),
			ErrorReports:        c.ErrorReports.ValueBool(),
			EnableAnalytics:     c.EnableAnalytics.ValueBool(),
			ContactEmailAddress: c.ContactEmailAddress.ValueString(),
		},
	}
}

func (c *InfinityManagerConfigModel) validate(ctx context.Context) diag.Diagnostics {
	var diags diag.Diagnostics

	c.setDefaults()

	if c.Hostname.IsNull() && !c.Hostname.IsUnknown() {
		diags.AddAttributeError(
			path.Root("hostname"),
			"Hostname is required",
			"Hostname must be set to a valid value.",
		)
	}
	if c.Domain.IsNull() && !c.Domain.IsUnknown() {
		diags.AddAttributeError(
			path.Root("domain"),
			"Domain is required",
			"Domain must be set to a valid value.",
		)
	}
	if c.IP.IsNull() && !c.IP.IsUnknown() {
		diags.AddAttributeError(
			path.Root("ip"),
			"IP is required",
			"IP must be set to a valid value.",
		)
	} else if !c.IP.IsUnknown() && c.IP.ValueString() != "" {
		if net.ParseIP(c.IP.ValueString()) == nil {
			diags.AddAttributeError(
				path.Root("ip"),
				"Invalid IP Address",
				fmt.Sprintf("IP '%s' is not a valid IP address.", c.IP.ValueString()),
			)
		}
	}
	if c.Mask.IsNull() && !c.Mask.IsUnknown() {
		diags.AddAttributeError(
			path.Root("mask"),
			"Mask is required",
			"Mask must be set to a valid value.",
		)
	} else if !c.Mask.IsUnknown() && c.Mask.ValueString() != "" {
		if net.ParseIP(c.Mask.ValueString()) == nil {
			diags.AddAttributeError(
				path.Root("mask"),
				"Invalid Mask",
				fmt.Sprintf("Mask '%s' is not a valid IP subnet.", c.Mask.ValueString()),
			)
		}
	}
	if c.GW.IsNull() && !c.GW.IsUnknown() {
		diags.AddAttributeError(
			path.Root("gw"),
			"Gateway is required",
			"Gateway must be set to a valid value.",
		)
	} else if !c.GW.IsUnknown() && c.GW.ValueString() != "" {
		if net.ParseIP(c.GW.ValueString()) == nil {
			diags.AddAttributeError(
				path.Root("gw"),
				"Invalid Gateway IP Address",
				fmt.Sprintf("Gateway '%s' is not a valid IP address.", c.GW.ValueString()),
			)
		}
	}
	if c.DNS.IsNull() && !c.DNS.IsUnknown() {
		diags.AddAttributeError(
			path.Root("dns"),
			"DNS is required",
			"DNS must be set to a valid value.",
		)
	} else if !c.DNS.IsUnknown() && c.DNS.ValueString() != "" {
		if net.ParseIP(c.DNS.ValueString()) == nil {
			diags.AddAttributeError(
				path.Root("dns"),
				"Invalid DNS IP Address",
				fmt.Sprintf("DNS '%s' is not a valid IP address.", c.DNS.ValueString()),
			)
		}
	}
	if c.NTP.IsNull() && !c.NTP.IsUnknown() {
		diags.AddAttributeError(
			path.Root("ntp"),
			"NTP is required",
			"NTP must be set to a valid value.",
		)
	}
	if c.User.IsNull() && !c.User.IsUnknown() {
		diags.AddAttributeError(
			path.Root("user"),
			"User is required",
			"User must be set to a valid value.",
		)
	}
	if c.Pass.IsNull() && !c.Pass.IsUnknown() {
		diags.AddAttributeError(
			path.Root("pass"),
			"Password is required",
			"Password must be set to a valid value.",
		)
	}
	if c.AdminPassword.IsNull() && !c.AdminPassword.IsUnknown() {
		diags.AddAttributeError(
			path.Root("admin_password"),
			"Admin Password is required",
			"Admin Password must be set to a valid value.",
		)
	}
	if c.ContactEmailAddress.IsNull() && !c.ContactEmailAddress.IsUnknown() {
		diags.AddAttributeError(
			path.Root("contact_email_address"),
			"Contact Email Address is required",
			"Contact Email Address must be set to a valid value.",
		)
	} else if !c.ContactEmailAddress.IsUnknown() && c.ContactEmailAddress.ValueString() != "" {
		if _, err := mail.ParseAddress(c.ContactEmailAddress.ValueString()); err != nil {
			diags.AddAttributeError(
				path.Root("contact_email_address"),
				"Invalid Email Address",
				fmt.Sprintf("Contact Email Address '%s' is not a valid email address.", c.ContactEmailAddress.ValueString()),
			)
		}
	}

	return diags
}

func (c *InfinityManagerConfigModel) setDefaults() {
	if c.ErrorReports.IsNull() {
		c.ErrorReports = types.BoolValue(false)
	}
	if c.EnableAnalytics.IsNull() {
		c.EnableAnalytics = types.BoolValue(false)
	}
}

func prepareOutput(output []byte) string {
	return strings.ReplaceAll(string(output), "\n", "")
}

func (c *InfinityManagerConfigModel) update() diag.Diagnostics {
	var diags diag.Diagnostics

	fullOutput, err := json.MarshalIndent(c.toOuterConfig(), "", "")
	if err != nil {
		diags.AddError("Failed to marshal infinity manager config for output", err.Error())
		return diags
	}

	mgmtNodeConfigOutput, err := json.MarshalIndent(c.toOuterConfig().ManagementNodeConfig, "", "")
	if err != nil {
		diags.AddError("Failed to marshal infinity manager management node config for output", err.Error())
		return diags
	}

	c.ID = types.StringValue(strconv.Itoa(helpers.String(prepareOutput(fullOutput))))
	c.Rendered = types.StringValue(prepareOutput(fullOutput))
	c.RenderedManagementNodeConfig = types.StringValue(prepareOutput(mgmtNodeConfigOutput))

	return diags
}
