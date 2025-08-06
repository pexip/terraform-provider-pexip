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

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
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
	MTU                         types.Int64  `tfsdk:"mtu"`
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
				Validators: []validator.String{
					stringvalidator.LengthAtMost(500),
				},
				MarkdownDescription: "Description of the management VM. Maximum length: 500 characters.",
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
				Optional:            true,
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
			"mtu": schema.Int64Attribute{
				Required: true,
				Validators: []validator.Int64{
					int64validator.Between(576, 9000),
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
				ElementType:         types.StringType,
				Optional:            true,
				MarkdownDescription: "List of static route URIs for the management VM.",
			},
			"event_sinks": schema.SetAttribute{
				ElementType:         types.StringType,
				Optional:            true,
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
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf("yes", "no", "keys_only"),
				},
				MarkdownDescription: "SSH access configuration. Valid values: yes, no, keys_only.",
			},
			"ssh_authorized_keys": schema.SetAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				MarkdownDescription: "List of SSH authorized key URIs for the management VM.",
			},
			"ssh_authorized_keys_use_cloud": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Whether to use cloud-based SSH authorized keys.",
			},
			"secondary_config_passphrase": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				MarkdownDescription: "Secondary configuration passphrase. This field is sensitive.",
			},
			"snmp_mode": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf("disabled", "v1v2c", "v3"),
				},
				MarkdownDescription: "SNMP mode configuration. Valid values: disabled, v1v2c, v3.",
			},
			"snmp_community": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				MarkdownDescription: "SNMP community string. This field is sensitive.",
			},
			"snmp_username": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "SNMP username for v3 authentication.",
			},
			"snmp_authentication_password": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				MarkdownDescription: "SNMP authentication password. This field is sensitive.",
			},
			"snmp_privacy_password": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				MarkdownDescription: "SNMP privacy password. This field is sensitive.",
			},
			"snmp_system_contact": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "SNMP system contact information.",
			},
			"snmp_system_location": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "SNMP system location information.",
			},
			"snmp_network_management_system": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "SNMP network management system URI.",
			},
			"initializing": schema.BoolAttribute{
				Required:            true,
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

func (r *InfinityManagementVMResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	plan := &InfinityManagementVMResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createRequest := &config.ManagementVMCreateRequest{
		Name:                       plan.Name.ValueString(),
		Description:                plan.Description.ValueString(),
		Address:                    plan.Address.ValueString(),
		Netmask:                    plan.Netmask.ValueString(),
		Gateway:                    plan.Gateway.ValueString(),
		Hostname:                   plan.Hostname.ValueString(),
		Domain:                     plan.Domain.ValueString(),
		AlternativeFQDN:            plan.AlternativeFQDN.ValueString(),
		MTU:                        int(plan.MTU.ValueInt64()),
		EnableSSH:                  plan.EnableSSH.ValueString(),
		SSHAuthorizedKeysUseCloud:  plan.SSHAuthorizedKeysUseCloud.ValueBool(),
		SecondaryConfigPassphrase:  plan.SecondaryConfigPassphrase.ValueString(),
		SNMPMode:                   plan.SNMPMode.ValueString(),
		SNMPCommunity:              plan.SNMPCommunity.ValueString(),
		SNMPUsername:               plan.SNMPUsername.ValueString(),
		SNMPAuthenticationPassword: plan.SNMPAuthenticationPassword.ValueString(),
		SNMPPrivacyPassword:        plan.SNMPPrivacyPassword.ValueString(),
		SNMPSystemContact:          plan.SNMPSystemContact.ValueString(),
		SNMPSystemLocation:         plan.SNMPSystemLocation.ValueString(),
		Initializing:               plan.Initializing.ValueBool(),
	}

	// Handle optional pointer fields
	if !plan.IPV6Address.IsNull() && !plan.IPV6Address.IsUnknown() {
		addr := plan.IPV6Address.ValueString()
		createRequest.IPV6Address = &addr
	}

	if !plan.IPV6Gateway.IsNull() && !plan.IPV6Gateway.IsUnknown() {
		gateway := plan.IPV6Gateway.ValueString()
		createRequest.IPV6Gateway = &gateway
	}

	if !plan.StaticNATAddress.IsNull() && !plan.StaticNATAddress.IsUnknown() {
		addr := plan.StaticNATAddress.ValueString()
		createRequest.StaticNATAddress = &addr
	}

	if !plan.HTTPProxy.IsNull() && !plan.HTTPProxy.IsUnknown() {
		proxy := plan.HTTPProxy.ValueString()
		createRequest.HTTPProxy = &proxy
	}

	if !plan.TLSCertificate.IsNull() && !plan.TLSCertificate.IsUnknown() {
		cert := plan.TLSCertificate.ValueString()
		createRequest.TLSCertificate = &cert
	}

	if !plan.SNMPNetworkManagementSystem.IsNull() && !plan.SNMPNetworkManagementSystem.IsUnknown() {
		nms := plan.SNMPNetworkManagementSystem.ValueString()
		createRequest.SNMPNetworkManagementSystem = &nms
	}

	// Handle list fields
	if !plan.DNSServers.IsNull() {
		var servers []string
		resp.Diagnostics.Append(plan.DNSServers.ElementsAs(ctx, &servers, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		createRequest.DNSServers = servers
	}

	if !plan.NTPServers.IsNull() {
		var servers []string
		resp.Diagnostics.Append(plan.NTPServers.ElementsAs(ctx, &servers, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		createRequest.NTPServers = servers
	}

	if !plan.SyslogServers.IsNull() {
		var servers []string
		resp.Diagnostics.Append(plan.SyslogServers.ElementsAs(ctx, &servers, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		createRequest.SyslogServers = servers
	}

	if !plan.StaticRoutes.IsNull() {
		var routes []string
		resp.Diagnostics.Append(plan.StaticRoutes.ElementsAs(ctx, &routes, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		createRequest.StaticRoutes = routes
	}

	if !plan.EventSinks.IsNull() {
		var sinks []string
		resp.Diagnostics.Append(plan.EventSinks.ElementsAs(ctx, &sinks, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		createRequest.EventSinks = sinks
	}

	if !plan.SSHAuthorizedKeys.IsNull() {
		var keys []string
		resp.Diagnostics.Append(plan.SSHAuthorizedKeys.ElementsAs(ctx, &keys, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		createRequest.SSHAuthorizedKeys = keys
	}

	createResponse, err := r.InfinityClient.Config().CreateManagementVM(ctx, createRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Infinity management VM",
			fmt.Sprintf("Could not create Infinity management VM: %s", err),
		)
		return
	}

	resourceID, err := createResponse.ResourceID()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Retrieving Infinity management VM ID",
			fmt.Sprintf("Could not retrieve ID for created Infinity management VM: %s", err),
		)
		return
	}

	// Read the state from the API to get all computed values
	model, err := r.read(ctx, resourceID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Created Infinity management VM",
			fmt.Sprintf("Could not read created Infinity management VM with ID %d: %s", resourceID, err),
		)
		return
	}
	tflog.Trace(ctx, fmt.Sprintf("created Infinity management VM with ID: %s, name: %s", model.ID, model.Name))

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *InfinityManagementVMResource) read(ctx context.Context, resourceID int) (*InfinityManagementVMResourceModel, error) {
	var data InfinityManagementVMResourceModel

	srv, err := r.InfinityClient.Config().GetManagementVM(ctx, resourceID)
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
	data.MTU = types.Int64Value(int64(srv.MTU))
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

	// Handle list fields
	if srv.DNSServers != nil {
		serversSet, diags := types.SetValueFrom(ctx, types.StringType, srv.DNSServers)
		if diags.HasError() {
			return nil, fmt.Errorf("failed to convert DNS servers: %s", diags.Errors())
		}
		data.DNSServers = serversSet
	} else {
		data.DNSServers = types.SetNull(types.StringType)
	}

	if srv.NTPServers != nil {
		serversSet, diags := types.SetValueFrom(ctx, types.StringType, srv.NTPServers)
		if diags.HasError() {
			return nil, fmt.Errorf("failed to convert NTP servers: %s", diags.Errors())
		}
		data.NTPServers = serversSet
	} else {
		data.NTPServers = types.SetNull(types.StringType)
	}

	if srv.SyslogServers != nil {
		serversSet, diags := types.SetValueFrom(ctx, types.StringType, srv.SyslogServers)
		if diags.HasError() {
			return nil, fmt.Errorf("failed to convert syslog servers: %s", diags.Errors())
		}
		data.SyslogServers = serversSet
	} else {
		data.SyslogServers = types.SetNull(types.StringType)
	}

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
	state, err := r.read(ctx, resourceID)
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
	// Management VMs do not support update operations
	resp.Diagnostics.AddError(
		"Update Not Supported",
		"Management VM resources cannot be updated. To change management VM settings, you must delete and recreate the resource.",
	)
}

func (r *InfinityManagementVMResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	state := &InfinityManagementVMResourceModel{}

	tflog.Info(ctx, "Deleting Infinity management VM")

	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.InfinityClient.Config().DeleteManagementVM(ctx, int(state.ResourceID.ValueInt32()))

	// Ignore 404 Not Found and Lookup errors on delete
	if err != nil && !isNotFoundError(err) && !isLookupError(err) {
		resp.Diagnostics.AddError(
			"Error Deleting Infinity management VM",
			fmt.Sprintf("Could not delete Infinity management VM with ID %s: %s", state.ID.ValueString(), err),
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
	model, err := r.read(ctx, resourceID)
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
