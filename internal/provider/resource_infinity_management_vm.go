/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/pexip/go-infinity-sdk/v38/config"

	"github.com/pexip/terraform-provider-pexip/internal/provider/validators"
)

var (
	_ resource.ResourceWithImportState = (*InfinityManagementVMResource)(nil)
)

type InfinityManagementVMResource struct {
	InfinityClient InfinityClient
}

type InfinityManagementVMResourceModel struct {
	ID                          types.String `tfsdk:"id"`
	ResourceID                  types.Int32  `tfsdk:"resource_id"`
	Name                        types.String `tfsdk:"name"`
	Description                 types.String `tfsdk:"description"`
	Address                     types.String `tfsdk:"address"`
	Netmask                     types.String `tfsdk:"netmask"`
	Gateway                     types.String `tfsdk:"gateway"`
	Hostname                    types.String `tfsdk:"hostname"`
	Domain                      types.String `tfsdk:"domain"`
	AlternativeFQDN             types.String `tfsdk:"alternative_fqdn"`
	IPV6Address                 types.String `tfsdk:"ipv6_address"`
	IPV6Gateway                 types.String `tfsdk:"ipv6_gateway"`
	MTU                         types.Int32  `tfsdk:"mtu"`
	StaticNATAddress            types.String `tfsdk:"static_nat_address"`
	DNSServers                  types.Set    `tfsdk:"dns_servers"`
	NTPServers                  types.Set    `tfsdk:"ntp_servers"`
	SyslogServers               types.Set    `tfsdk:"syslog_servers"`
	StaticRoutes                types.Set    `tfsdk:"static_routes"`
	EventSinks                  types.Set    `tfsdk:"event_sinks"`
	HTTPProxy                   types.String `tfsdk:"http_proxy"`
	TLSCertificate              types.String `tfsdk:"tls_certificate"`
	EnableSSH                   types.String `tfsdk:"enable_ssh"`
	SSHAuthorizedKeys           types.Set    `tfsdk:"ssh_authorized_keys"`
	SSHAuthorizedKeysUseCloud   types.Bool   `tfsdk:"ssh_authorized_keys_use_cloud"`
	SecondaryConfigPassphrase   types.String `tfsdk:"secondary_config_passphrase"`
	SNMPMode                    types.String `tfsdk:"snmp_mode"`
	SNMPCommunity               types.String `tfsdk:"snmp_community"`
	SNMPUsername                types.String `tfsdk:"snmp_username"`
	SNMPAuthenticationPassword  types.String `tfsdk:"snmp_authentication_password"`
	SNMPPrivacyPassword         types.String `tfsdk:"snmp_privacy_password"`
	SNMPSystemContact           types.String `tfsdk:"snmp_system_contact"`
	SNMPSystemLocation          types.String `tfsdk:"snmp_system_location"`
	SNMPNetworkManagementSystem types.String `tfsdk:"snmp_network_management_system"`
	Initializing                types.Bool   `tfsdk:"initializing"`
	Primary                     types.Bool   `tfsdk:"primary"`
}

func (r *InfinityManagementVMResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_infinity_management_vm"
}

func (r *InfinityManagementVMResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	p, ok := req.ProviderData.(*PexipProvider)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *PexipProvider, got: %T. Please report this issue to the provider developers", req.ProviderData),
		)
		return
	}

	r.InfinityClient = p.client
}

func (r *InfinityManagementVMResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Resource URI for the management VM in Infinity",
			},
			"resource_id": schema.Int32Attribute{
				Computed:            true,
				MarkdownDescription: "The resource integer identifier for the management VM in Infinity",
			},
			"name": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					stringvalidator.LengthAtMost(250),
				},
				MarkdownDescription: "The name of the management VM. Maximum length: 250 characters.",
			},
			"description": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(""),
				Validators: []validator.String{
					stringvalidator.LengthAtMost(250),
				},
				MarkdownDescription: "Description of the management VM. Maximum length: 250 characters.",
			},
			"address": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					validators.IPAddress(),
				},
				MarkdownDescription: "The IP address of the management VM.",
			},
			"netmask": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					validators.Netmask(),
				},
				MarkdownDescription: "The network mask for the management VM.",
			},
			"gateway": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					validators.IPAddress(),
				},
				MarkdownDescription: "The gateway IP address for the management VM.",
			},
			"hostname": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					stringvalidator.LengthAtMost(253),
				},
				MarkdownDescription: "The hostname of the management VM. Maximum length: 253 characters.",
			},
			"domain": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					validators.Domain(),
				},
				MarkdownDescription: "The domain name for the management VM.",
			},
			"alternative_fqdn": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(""),
				Validators: []validator.String{
					stringvalidator.LengthAtMost(255),
				},
				MarkdownDescription: "Alternative fully qualified domain name for the management VM.",
			},
			"ipv6_address": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The IPv6 address of the management VM.",
			},
			"ipv6_gateway": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The IPv6 gateway for the management VM.",
			},
			"mtu": schema.Int32Attribute{
				Optional: true,
				Validators: []validator.Int32{
					int32validator.Between(576, 9000),
				},
				MarkdownDescription: "Maximum Transmission Unit (MTU) size. Valid range: 576-9000.",
			},
			"static_nat_address": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					validators.IPAddress(),
				},
				MarkdownDescription: "Static NAT address for the management VM.",
			},
			"dns_servers": schema.SetAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				MarkdownDescription: "List of DNS server URIs for the management VM.",
			},
			"ntp_servers": schema.SetAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				MarkdownDescription: "List of NTP server URIs for the management VM.",
			},
			"syslog_servers": schema.SetAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				MarkdownDescription: "List of syslog server URIs for the management VM.",
			},
			"static_routes": schema.SetAttribute{
				Optional:    true,
				Computed:    true,
				Default:     nil,
				ElementType: types.StringType,
				Description: "Additional configuration to permit routing of traffic to networks not accessible through the configured default gateway.",
			},
			"event_sinks": schema.SetAttribute{
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "List of event sink URIs for the management VM.",
			},
			"http_proxy": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "HTTP proxy URI for the management VM.",
			},
			"tls_certificate": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "TLS certificate URI for the management VM.",
			},
			"enable_ssh": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("GLOBAL"),
				Validators: []validator.String{
					stringvalidator.OneOf("GLOBAL", "OFF", "ON"),
				},
				MarkdownDescription: "Allows an administrator to log in to this node over SSH. Valid values are: global, off, on. Defaults to global.",
			},
			"ssh_authorized_keys": schema.SetAttribute{
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
				Default:             nil,
				MarkdownDescription: "The selected authorized keys.",
			},
			"ssh_authorized_keys_use_cloud": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
				MarkdownDescription: "Allows use of SSH keys configured in the cloud service.",
			},
			"secondary_config_passphrase": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
				MarkdownDescription: "Secondary configuration passphrase. This field is sensitive.",
			},
			"snmp_mode": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("DISABLED"),
				Validators: []validator.String{
					stringvalidator.OneOf("DISABLED", "STANDARD", "AUTHPRIV"),
				},
				MarkdownDescription: "The SNMP mode.",
			},
			"snmp_community": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				Computed:            true,
				Default:             stringdefault.StaticString("public"),
				MarkdownDescription: "SNMP community string. This field is sensitive.",
			},
			"snmp_username": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(""),
				Validators: []validator.String{
					stringvalidator.LengthAtMost(100),
				},
				MarkdownDescription: "The username used to authenticate SNMPv3 requests. Maximum length: 100 characters.",
			},
			"snmp_authentication_password": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				Computed:            true,
				MarkdownDescription: "The password used for SNMPv3 privacy. Minimum length: 8 characters. Maximum length: 100 characters.",
			},
			"snmp_privacy_password": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				Computed:            true,
				MarkdownDescription: "The password used for SNMPv3 privacy. Minimum length: 8 characters. Maximum length: 100 characters.",
			},
			"snmp_system_contact": schema.StringAttribute{
				Computed:            true,
				Optional:            true,
				Default:             stringdefault.StaticString("admin@domain.com"),
				MarkdownDescription: "SNMP system contact information.",
			},
			"snmp_system_location": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("Virtual machine"),
				MarkdownDescription: "SNMP system location information.",
			},
			"snmp_network_management_system": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "SNMP network management system URI.",
			},
			"initializing": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Whether the management VM is in initializing state.",
			},
			"primary": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Whether this is the primary management VM.",
			},
		},
		MarkdownDescription: "Manages a management VM configuration with the Infinity service. Management VMs are Pexip Infinity Manager nodes that control the platform. Note: This resource supports Create, Read, and Delete operations only - updates are not supported.",
	}
}

func (r *InfinityManagementVMResource) buildUpdateRequest(plan *InfinityManagementVMResourceModel) *config.ManagementVMUpdateRequest {

	updateRequest := &config.ManagementVMUpdateRequest{
		Name:        plan.Name.ValueString(),
		Address:     plan.Address.ValueString(),
		Netmask:     plan.Netmask.ValueString(),
		Gateway:     plan.Gateway.ValueString(),
		Hostname:    plan.Hostname.ValueString(),
		Domain:      plan.Domain.ValueString(),
		Description: plan.Description.ValueString(),
	}

	// Handle optional pointer fields
	if !plan.MTU.IsNull() {
		updateRequest.MTU = int(plan.MTU.ValueInt32())
	}
	if !plan.IPV6Address.IsNull() && !plan.IPV6Address.IsUnknown() {
		addr := plan.IPV6Address.ValueString()
		updateRequest.IPV6Address = &addr
	}

	if !plan.IPV6Gateway.IsNull() && !plan.IPV6Gateway.IsUnknown() {
		gateway := plan.IPV6Gateway.ValueString()
		updateRequest.IPV6Gateway = &gateway
	}

	if !plan.StaticNATAddress.IsNull() && !plan.StaticNATAddress.IsUnknown() {
		addr := plan.StaticNATAddress.ValueString()
		updateRequest.StaticNATAddress = &addr
	}

	if !plan.HTTPProxy.IsNull() && !plan.HTTPProxy.IsUnknown() {
		proxy := plan.HTTPProxy.ValueString()
		updateRequest.HTTPProxy = &proxy
	}

	if !plan.TLSCertificate.IsNull() && !plan.TLSCertificate.IsUnknown() {
		cert := plan.TLSCertificate.ValueString()
		updateRequest.TLSCertificate = &cert
	}
	if !plan.EnableSSH.IsNull() && !plan.EnableSSH.IsUnknown() {
		updateRequest.EnableSSH = plan.EnableSSH.ValueString()
	}
	if !plan.SNMPNetworkManagementSystem.IsNull() && !plan.SNMPNetworkManagementSystem.IsUnknown() {
		nms := plan.SNMPNetworkManagementSystem.ValueString()
		updateRequest.SNMPNetworkManagementSystem = &nms
	}
	if !plan.SNMPMode.IsNull() && !plan.SNMPMode.IsUnknown() {
		updateRequest.SNMPMode = plan.SNMPMode.ValueString()
	}
	if !plan.SNMPCommunity.IsNull() && !plan.SNMPCommunity.IsUnknown() {
		updateRequest.SNMPCommunity = plan.SNMPCommunity.ValueString()
	}
	if !plan.SNMPUsername.IsNull() && !plan.SNMPUsername.IsUnknown() {
		updateRequest.SNMPUsername = plan.SNMPUsername.ValueString()
	}
	if !plan.SNMPAuthenticationPassword.IsNull() && !plan.SNMPAuthenticationPassword.IsUnknown() {
		updateRequest.SNMPAuthenticationPassword = plan.SNMPAuthenticationPassword.ValueString()
	}
	if !plan.SNMPPrivacyPassword.IsNull() && !plan.SNMPPrivacyPassword.IsUnknown() {
		updateRequest.SNMPPrivacyPassword = plan.SNMPPrivacyPassword.ValueString()
	}
	if !plan.SNMPSystemContact.IsNull() && !plan.SNMPSystemContact.IsUnknown() {
		updateRequest.SNMPSystemContact = plan.SNMPSystemContact.ValueString()
	}
	if !plan.SNMPSystemLocation.IsNull() && !plan.SNMPSystemLocation.IsUnknown() {
		updateRequest.SNMPSystemLocation = plan.SNMPSystemLocation.ValueString()
	}
	if !plan.SSHAuthorizedKeysUseCloud.IsNull() {
		updateRequest.SSHAuthorizedKeysUseCloud = plan.SSHAuthorizedKeysUseCloud.ValueBool()
	}

	return updateRequest
}

func (r *InfinityManagementVMResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// For singleton resources, Create is actually Update since the resource always exists
	plan := &InfinityManagementVMResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the list of DNS, NTP, and Syslog servers
	dnsServers, diags := getStringList(ctx, plan.DNSServers)
	resp.Diagnostics.Append(diags...)
	ntpServers, diags := getStringList(ctx, plan.NTPServers)
	resp.Diagnostics.Append(diags...)
	syslogServers, diags := getStringList(ctx, plan.SyslogServers)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateRequest := r.buildUpdateRequest(plan)

	updateRequest.DNSServers = dnsServers
	updateRequest.NTPServers = ntpServers
	updateRequest.SyslogServers = syslogServers

	_, err := r.InfinityClient.Config().UpdateManagementVM(ctx, updateRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Infinity Management VM",
			fmt.Sprintf("Could not update Infinity Management VM: %s", err),
		)
		return
	}

	// Re-read the resource to get the latest state
	updatedModel, err := r.read(ctx, 1, plan.SNMPCommunity.ValueString(), plan.SNMPAuthenticationPassword.ValueString(), plan.SNMPPrivacyPassword.ValueString(), plan.SecondaryConfigPassphrase.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Updated Infinity management VM",
			fmt.Sprintf("Could not read updated Infinity management VM: %s", err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, updatedModel)...)
}

func (r *InfinityManagementVMResource) read(ctx context.Context, resourceID int, snmpCommunity, snmpAuthPass, snmpPrivPass, secondaryConfigPass string) (*InfinityManagementVMResourceModel, error) {
	var data InfinityManagementVMResourceModel

	srv, err := r.InfinityClient.Config().GetManagementVM(ctx)
	if err != nil {
		return nil, err
	}

	if srv.ResourceURI == "" {
		return nil, fmt.Errorf("management VM with ID %d not found", resourceID)
	}

	data.ID = types.StringValue(srv.ResourceURI)
	data.ResourceID = types.Int32Value(int32(resourceID)) // #nosec G115 -- API values are expected to be within int32 range
	data.Name = types.StringValue(srv.Name)
	data.Description = types.StringValue(srv.Description)
	data.Address = types.StringValue(srv.Address)
	data.Netmask = types.StringValue(srv.Netmask)
	data.Gateway = types.StringValue(srv.Gateway)
	data.Hostname = types.StringValue(srv.Hostname)
	data.Domain = types.StringValue(srv.Domain)
	data.AlternativeFQDN = types.StringValue(srv.AlternativeFQDN)
	data.MTU = types.Int32Value(int32(srv.MTU)) // #nosec G115 -- API values are expected to be within int32 range
	data.EnableSSH = types.StringValue(srv.EnableSSH)
	data.SSHAuthorizedKeysUseCloud = types.BoolValue(srv.SSHAuthorizedKeysUseCloud)
	data.SecondaryConfigPassphrase = types.StringValue(srv.SecondaryConfigPassphrase)
	data.SNMPMode = types.StringValue(srv.SNMPMode)
	data.SNMPCommunity = types.StringValue(srv.SNMPCommunity)
	data.SNMPUsername = types.StringValue(srv.SNMPUsername)
	data.SNMPAuthenticationPassword = types.StringValue(srv.SNMPAuthenticationPassword)
	data.SNMPPrivacyPassword = types.StringValue(srv.SNMPPrivacyPassword)
	data.SNMPSystemContact = types.StringValue(srv.SNMPSystemContact)
	data.SNMPSystemLocation = types.StringValue(srv.SNMPSystemLocation)
	data.Initializing = types.BoolValue(srv.Initializing)
	data.Primary = types.BoolValue(srv.Primary)

	// Handle optional pointer fields
	if srv.IPV6Address != nil {
		data.IPV6Address = types.StringValue(*srv.IPV6Address)
	} else {
		data.IPV6Address = types.StringNull()
	}

	if srv.IPV6Gateway != nil {
		data.IPV6Gateway = types.StringValue(*srv.IPV6Gateway)
	} else {
		data.IPV6Gateway = types.StringNull()
	}

	if srv.StaticNATAddress != nil {
		data.StaticNATAddress = types.StringValue(*srv.StaticNATAddress)
	} else {
		data.StaticNATAddress = types.StringNull()
	}

	if srv.HTTPProxy != nil {
		data.HTTPProxy = types.StringValue(*srv.HTTPProxy)
	} else {
		data.HTTPProxy = types.StringNull()
	}

	if srv.TLSCertificate != nil {
		data.TLSCertificate = types.StringValue(*srv.TLSCertificate)
	} else {
		data.TLSCertificate = types.StringNull()
	}

	if srv.SNMPNetworkManagementSystem != nil {
		data.SNMPNetworkManagementSystem = types.StringValue(*srv.SNMPNetworkManagementSystem)
	} else {
		data.SNMPNetworkManagementSystem = types.StringNull()
	}

	// Convert DNS servers from SDK to Terraform format
	var dnsServers []string
	for _, dns := range srv.DNSServers {
		dnsServers = append(dnsServers, fmt.Sprintf("/api/admin/configuration/v1/dns_server/%d/", dns.ID))
	}
	dnsSetValue, diags := types.SetValueFrom(ctx, types.StringType, dnsServers)
	if diags.HasError() {
		return nil, fmt.Errorf("error converting DNS servers: %v", diags)
	}
	data.DNSServers = dnsSetValue

	// Convert NTP servers from SDK to Terraform format
	var ntpServers []string
	for _, ntp := range srv.NTPServers {
		ntpServers = append(ntpServers, fmt.Sprintf("/api/admin/configuration/v1/ntp_server/%d/", ntp.ID))
	}
	ntpSetValue, diags := types.SetValueFrom(ctx, types.StringType, ntpServers)
	if diags.HasError() {
		return nil, fmt.Errorf("error converting NTP servers: %v", diags)
	}
	data.NTPServers = ntpSetValue

	// Convert Syslog servers from SDK to Terraform format
	var syslogServers []string
	for _, syslog := range srv.SyslogServers {
		syslogServers = append(syslogServers, fmt.Sprintf("/api/admin/configuration/v1/syslog_server/%d/", syslog.ID))
	}
	syslogSetValue, diags := types.SetValueFrom(ctx, types.StringType, syslogServers)
	if diags.HasError() {
		return nil, fmt.Errorf("error converting Syslog servers: %v", diags)
	}
	data.SyslogServers = syslogSetValue

	if srv.StaticRoutes != nil {
		routesSet, diags := types.SetValueFrom(ctx, types.StringType, srv.StaticRoutes)
		if diags.HasError() {
			return nil, fmt.Errorf("failed to convert static routes: %s", diags.Errors())
		}
		data.StaticRoutes = routesSet
	} else {
		data.StaticRoutes = types.SetNull(types.StringType)
	}

	if srv.EventSinks != nil {
		sinksSet, diags := types.SetValueFrom(ctx, types.StringType, srv.EventSinks)
		if diags.HasError() {
			return nil, fmt.Errorf("failed to convert event sinks: %s", diags.Errors())
		}
		data.EventSinks = sinksSet
	} else {
		data.EventSinks = types.SetNull(types.StringType)
	}

	if srv.SSHAuthorizedKeys != nil {
		keysSet, diags := types.SetValueFrom(ctx, types.StringType, srv.SSHAuthorizedKeys)
		if diags.HasError() {
			return nil, fmt.Errorf("failed to convert SSH authorized keys: %s", diags.Errors())
		}
		data.SSHAuthorizedKeys = keysSet
	} else {
		data.SSHAuthorizedKeys = types.SetNull(types.StringType)
	}

	return &data, nil
}

func (r *InfinityManagementVMResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	state := &InfinityManagementVMResourceModel{}

	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resourceID := int(state.ResourceID.ValueInt32())
	state, err := r.read(ctx, resourceID, state.SNMPCommunity.ValueString(), state.SNMPAuthenticationPassword.ValueString(), state.SNMPPrivacyPassword.ValueString(), state.SecondaryConfigPassphrase.ValueString())
	if err != nil {
		// Check if the error is a 404 (not found)
		if isNotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading Infinity management VM",
			fmt.Sprintf("Could not read Infinity management VM: %s", err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *InfinityManagementVMResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	plan := &InfinityManagementVMResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the list of DNS, NTP, and Syslog servers
	dnsServers, diags := getStringList(ctx, plan.DNSServers)
	resp.Diagnostics.Append(diags...)
	ntpServers, diags := getStringList(ctx, plan.NTPServers)
	resp.Diagnostics.Append(diags...)
	syslogServers, diags := getStringList(ctx, plan.SyslogServers)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateRequest := r.buildUpdateRequest(plan)

	updateRequest.DNSServers = dnsServers
	updateRequest.NTPServers = ntpServers
	updateRequest.SyslogServers = syslogServers
	_, err := r.InfinityClient.Config().UpdateManagementVM(ctx, updateRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Infinity Management VM",
			fmt.Sprintf("Could not update Infinity Management VM: %s", err),
		)
		return
	}

	// Re-read the resource to get the latest state
	updatedModel, err := r.read(ctx, 1, plan.SNMPCommunity.ValueString(), plan.SNMPAuthenticationPassword.ValueString(), plan.SNMPPrivacyPassword.ValueString(), plan.SecondaryConfigPassphrase.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Updated Infinity management VM",
			fmt.Sprintf("Could not read updated Infinity management VM: %s", err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, updatedModel)...)
}

func (r *InfinityManagementVMResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Info(ctx, "Setting Infinity management VM to empty state")

	updateRequest := &config.ManagementVMUpdateRequest{
		DNSServers:                  []string{},
		NTPServers:                  []string{},
		SyslogServers:               []string{},
		SSHAuthorizedKeys:           []string{},
		StaticRoutes:                []string{},
		TLSCertificate:              nil,
		SNMPNetworkManagementSystem: nil,
	}

	_, err := r.InfinityClient.Config().UpdateManagementVM(ctx, updateRequest)
	if err != nil && !isNotFoundError(err) && !isLookupError(err) {
		resp.Diagnostics.AddError(
			"Error Resetting Infinity management VM configuration",
			fmt.Sprintf("Could not reset Infinity management VM configuration: %s", err),
		)
		return
	}
}

func (r *InfinityManagementVMResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resourceID, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Import ID must be a valid integer for the resource ID. Got: %s", req.ID),
		)
		return
	}

	tflog.Trace(ctx, fmt.Sprintf("Importing Infinity management VM with resource ID: %d", resourceID))

	// Read the resource from the API
	model, err := r.read(ctx, resourceID, "", "", "", "")
	if err != nil {
		// Check if the error is a 404 (not found)
		if isNotFoundError(err) {
			resp.Diagnostics.AddError(
				"Infinity Management VM Not Found",
				fmt.Sprintf("Infinity management VM with resource ID %d not found.", resourceID),
			)
			return
		}
		resp.Diagnostics.AddError(
			"Error Importing Infinity Management VM",
			fmt.Sprintf("Could not import Infinity management VM with resource ID %d: %s", resourceID, err),
		)
		return
	}

	// Set the state from the imported resource
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}
