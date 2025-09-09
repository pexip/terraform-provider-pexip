/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package provider

import (
	"context"
	"fmt"
	"sort"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/pexip/go-infinity-sdk/v38/config"
)

var (
	_ resource.ResourceWithImportState = (*InfinitySystemLocationResource)(nil)
)

type InfinitySystemLocationResource struct {
	InfinityClient InfinityClient
}

type InfinitySystemLocationResourceModel struct {
	ID                          types.String `tfsdk:"id"`
	ResourceID                  types.Int32  `tfsdk:"resource_id"`
	Name                        types.String `tfsdk:"name"`
	BDPMPINChecksEnabled        types.String `tfsdk:"bdpm_pin_checks_enabled"`
	BDPMScanQuarantineEnabled   types.String `tfsdk:"bdpm_scan_quarantine_enabled"`
	Description                 types.String `tfsdk:"description"`
	LocalMSSIPDomain            types.String `tfsdk:"local_mssip_domain"`
	MTU                         types.Int32  `tfsdk:"mtu"`
	UseRelayCandidatesOnly      types.Bool   `tfsdk:"use_relay_candidates_only"`
	DNSServers                  types.Set    `tfsdk:"dns_servers"`
	NTPServers                  types.Set    `tfsdk:"ntp_servers"`
	SyslogServers               types.Set    `tfsdk:"syslog_servers"`
	H323GateKeeper              types.String `tfsdk:"h323_gatekeeper"`
	SNMPNetworkManagementSystem types.String `tfsdk:"snmp_network_management_system"`
	SIPProxy                    types.String `tfsdk:"sip_proxy"`
	HTTPProxy                   types.String `tfsdk:"http_proxy"`
	MSSIPProxy                  types.String `tfsdk:"mssip_proxy"`
	TeamsProxy                  types.String `tfsdk:"teams_proxy"`
	TURNServer                  types.String `tfsdk:"turn_server"`
	STUNServer                  types.String `tfsdk:"stun_server"`
	ClientTURNServers           types.Set    `tfsdk:"client_turn_servers"`
	ClientSTUNServers           types.Set    `tfsdk:"client_stun_servers"`
	MediaQOS                    types.Int32  `tfsdk:"media_qos"`
	SignallingQOS               types.Int32  `tfsdk:"signalling_qos"`
	TranscodingLocation         types.String `tfsdk:"transcoding_location"`
	OverflowLocation1           types.String `tfsdk:"overflow_location1"`
	OverflowLocation2           types.String `tfsdk:"overflow_location2"`
	PolicyServer                types.String `tfsdk:"policy_server"`
	EventSinks                  types.Set    `tfsdk:"event_sinks"`
	LiveCaptionsDialOut1        types.String `tfsdk:"live_captions_dial_out_1"`
	LiveCaptionsDialOut2        types.String `tfsdk:"live_captions_dial_out_2"`
	LiveCaptionsDialOut3        types.String `tfsdk:"live_captions_dial_out_3"`
}

func getStringList(ctx context.Context, set types.Set) ([]string, diag.Diagnostics) {
	if set.IsNull() || set.IsUnknown() {
		return nil, nil
	}
	var items []string
	diags := set.ElementsAs(ctx, &items, false)
	if diags.HasError() {
		return nil, diags
	}
	sort.Strings(items)
	return items, diags
}

func (r *InfinitySystemLocationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_infinity_system_location"
}

func (r *InfinitySystemLocationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *InfinitySystemLocationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Resource URI for the system location in Infinity",
			},
			"resource_id": schema.Int32Attribute{
				Computed:            true,
				MarkdownDescription: "The resource integer identifier for the system location in Infinity",
			},
			"name": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(250),
				},
				MarkdownDescription: "The name used to refer to this system location. Maximum length: 250 characters.",
			},
			"description": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(""),
				Validators: []validator.String{
					stringvalidator.LengthAtMost(250),
				},
				MarkdownDescription: "A description of the system location. Maximum length: 250 characters.",
			},
			"dns_servers": schema.SetAttribute{
				Optional:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "List of DNS server resource URIs for this system location.",
			},
			"ntp_servers": schema.SetAttribute{
				Optional:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "List of NTP server resource URIs for this system location.",
			},
			"mtu": schema.Int32Attribute{
				Optional:            true,
				Computed:            true,
				Default:             int32default.StaticInt32(1500),
				MarkdownDescription: "Maximum Transmission Unit for this system location. Range: 512 to 1500.",
			},
			"syslog_servers": schema.SetAttribute{
				Optional:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "The Syslog servers to be used by Conferencing Nodes deployed in this Location.",
			},
			"h323_gatekeeper": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "H.323 Gatekeeper resource URI.",
			},
			"snmp_network_management_system": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "SNMP Network Management System resource URI.",
			},
			"sip_proxy": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "SIP Proxy resource URI.",
			},
			"http_proxy": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "HTTP Proxy resource URI.",
			},
			"mssip_proxy": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Microsoft SIP Proxy resource URI.",
			},
			"teams_proxy": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Teams Proxy resource URI.",
			},
			"turn_server": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "TURN server resource URI.",
			},
			"stun_server": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "STUN server resource URI.",
			},
			"client_turn_servers": schema.SetAttribute{
				Optional:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "List of client TURN server URIs.",
			},
			"client_stun_servers": schema.SetAttribute{
				Optional:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "List of client STUN server URIs.",
			},
			"use_relay_candidates_only": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Whether to use relay candidates only.",
			},
			"media_qos": schema.Int32Attribute{
				Optional:            true,
				Computed:            true,
				Default:             int32default.StaticInt32(0),
				MarkdownDescription: "Media QoS value.",
			},
			"signalling_qos": schema.Int32Attribute{
				Optional:            true,
				Computed:            true,
				Default:             int32default.StaticInt32(0),
				MarkdownDescription: "Signalling QoS value.",
			},
			"transcoding_location": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Transcoding location resource URI.",
			},
			"overflow_location1": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Overflow location 1 resource URI.",
			},
			"overflow_location2": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Overflow location 2 resource URI.",
			},
			"local_mssip_domain": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
				MarkdownDescription: "Local Microsoft SIP domain.",
			},
			"policy_server": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Policy server resource URI.",
			},
			"event_sinks": schema.SetAttribute{
				Optional:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "List of event sink URIs.",
			},
			"bdpm_pin_checks_enabled": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("GLOBAL"),
				Validators: []validator.String{
					stringvalidator.OneOf("GLOBAL", "OFF", "ON"),
				},
				MarkdownDescription: "Whether BDPM PIN checks are enabled.",
			},
			"bdpm_scan_quarantine_enabled": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("GLOBAL"),
				Validators: []validator.String{
					stringvalidator.OneOf("GLOBAL", "OFF", "ON"),
				},
				MarkdownDescription: "Whether BDPM scan quarantine is enabled.",
			},
			"live_captions_dial_out_1": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Live captions dial out 1 URI.",
			},
			"live_captions_dial_out_2": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Live captions dial out 2 URI.",
			},
			"live_captions_dial_out_3": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Live captions dial out 3 URI.",
			},
		},
		MarkdownDescription: "Registers a system location with the Infinity service.",
	}
}

func (r *InfinitySystemLocationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	plan := &InfinitySystemLocationResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert List attributes to []string
	dnsServers, diags := getStringList(ctx, plan.DNSServers)
	resp.Diagnostics.Append(diags...)
	ntpServers, diags := getStringList(ctx, plan.NTPServers)
	resp.Diagnostics.Append(diags...)
	syslogServers, diags := getStringList(ctx, plan.SyslogServers)
	resp.Diagnostics.Append(diags...)
	clientTurnServers, diags := getStringList(ctx, plan.ClientTURNServers)
	resp.Diagnostics.Append(diags...)
	clientStunServers, diags := getStringList(ctx, plan.ClientSTUNServers)
	resp.Diagnostics.Append(diags...)
	eventSinks, diags := getStringList(ctx, plan.EventSinks)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// initialize create request with required fields, list fields, and fields with defaults
	createRequest := &config.SystemLocationCreateRequest{
		Name:                      plan.Name.ValueString(),
		DNSServers:                dnsServers,
		NTPServers:                ntpServers,
		SyslogServers:             syslogServers,
		ClientTURNServers:         clientTurnServers,
		ClientSTUNServers:         clientStunServers,
		EventSinks:                eventSinks,
		Description:               plan.Description.ValueString(),
		BDPMPinChecksEnabled:      plan.BDPMPINChecksEnabled.ValueString(),
		BDPMScanQuarantineEnabled: plan.BDPMScanQuarantineEnabled.ValueString(),
		LocalMSSIPDomain:          plan.LocalMSSIPDomain.ValueString(),
		MTU:                       int(plan.MTU.ValueInt32()),
		UseRelayCandidatesOnly:    plan.UseRelayCandidatesOnly.ValueBool(),
	}

	// All nullable fields
	// Only set optional fields if they are not null in the plan
	if !plan.H323GateKeeper.IsNull() && !plan.H323GateKeeper.IsUnknown() {
		value := plan.H323GateKeeper.ValueString()
		createRequest.H323Gatekeeper = &value
	}
	if !plan.SNMPNetworkManagementSystem.IsNull() && !plan.SNMPNetworkManagementSystem.IsUnknown() {
		value := plan.SNMPNetworkManagementSystem.ValueString()
		createRequest.SNMPNetworkManagementSystem = &value
	}
	if !plan.SIPProxy.IsNull() && !plan.SIPProxy.IsUnknown() {
		value := plan.SIPProxy.ValueString()
		createRequest.SIPProxy = &value
	}
	if !plan.MSSIPProxy.IsNull() && !plan.MSSIPProxy.IsUnknown() {
		value := plan.MSSIPProxy.ValueString()
		createRequest.MSSIPProxy = &value
	}
	if !plan.HTTPProxy.IsNull() && !plan.HTTPProxy.IsUnknown() {
		value := plan.HTTPProxy.ValueString()
		createRequest.HTTPProxy = &value
	}
	if !plan.TeamsProxy.IsNull() && !plan.TeamsProxy.IsUnknown() {
		value := plan.TeamsProxy.ValueString()
		createRequest.TeamsProxy = &value
	}
	if !plan.TURNServer.IsNull() && !plan.TURNServer.IsUnknown() {
		value := plan.TURNServer.ValueString()
		createRequest.TURNServer = &value
	}
	if !plan.STUNServer.IsNull() && !plan.STUNServer.IsUnknown() {
		value := plan.STUNServer.ValueString()
		createRequest.STUNServer = &value
	}
	if !plan.MediaQOS.IsNull() && !plan.MediaQOS.IsUnknown() {
		value := int(plan.MediaQOS.ValueInt32())
		createRequest.MediaQoS = &value
	}
	if !plan.SignallingQOS.IsNull() && !plan.SignallingQOS.IsUnknown() {
		value := int(plan.SignallingQOS.ValueInt32())
		createRequest.SignallingQoS = &value
	}
	if !plan.TranscodingLocation.IsNull() && !plan.TranscodingLocation.IsUnknown() {
		value := plan.TranscodingLocation.ValueString()
		createRequest.TranscodingLocation = &value
	}
	if !plan.OverflowLocation1.IsNull() && !plan.OverflowLocation1.IsUnknown() {
		value := plan.OverflowLocation1.ValueString()
		createRequest.OverflowLocation1 = &value
	}
	if !plan.OverflowLocation2.IsNull() && !plan.OverflowLocation2.IsUnknown() {
		value := plan.OverflowLocation2.ValueString()
		createRequest.OverflowLocation2 = &value
	}
	if !plan.PolicyServer.IsNull() && !plan.PolicyServer.IsUnknown() {
		value := plan.PolicyServer.ValueString()
		createRequest.PolicyServer = &value
	}
	if !plan.LiveCaptionsDialOut1.IsNull() && !plan.LiveCaptionsDialOut1.IsUnknown() {
		value := plan.LiveCaptionsDialOut1.ValueString()
		createRequest.LiveCaptionsDialOut1 = &value
	}
	if !plan.LiveCaptionsDialOut2.IsNull() && !plan.LiveCaptionsDialOut2.IsUnknown() {
		value := plan.LiveCaptionsDialOut2.ValueString()
		createRequest.LiveCaptionsDialOut2 = &value
	}
	if !plan.LiveCaptionsDialOut3.IsNull() && !plan.LiveCaptionsDialOut3.IsUnknown() {
		value := plan.LiveCaptionsDialOut3.ValueString()
		createRequest.LiveCaptionsDialOut3 = &value
	}

	createResponse, err := r.InfinityClient.Config().CreateSystemLocation(ctx, createRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Infinity system location",
			fmt.Sprintf("Could not create Infinity system location: %s", err),
		)
		return
	}

	resourceID, err := createResponse.ResourceID()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Retrieving Infinity system location ID",
			fmt.Sprintf("Could not retrieve ID for created Infinity system location: %s", err),
		)
		return
	}

	// Read the state from the API to get all computed values
	model, err := r.read(ctx, resourceID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Created Infinity system location",
			fmt.Sprintf("Could not read created Infinity system location with ID %d: %s", resourceID, err),
		)
		return
	}
	tflog.Trace(ctx, fmt.Sprintf("created Infinity system location with ID: %s, name: %s", model.ID, model.Name))

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *InfinitySystemLocationResource) read(ctx context.Context, resourceID int) (*InfinitySystemLocationResourceModel, error) {
	var data InfinitySystemLocationResourceModel

	srv, err := r.InfinityClient.Config().GetSystemLocation(ctx, resourceID)
	if err != nil {
		return nil, err
	}

	if srv.ResourceURI == "" {
		return nil, fmt.Errorf("system location with ID %d not found", resourceID)
	}

	data.ID = types.StringValue(srv.ResourceURI)
	data.ResourceID = types.Int32Value(int32(resourceID)) // #nosec G115 -- API values are expected to be within int32 range
	data.Description = types.StringValue(srv.Description)
	data.MTU = types.Int32Value(int32(srv.MTU)) // #nosec G115 -- API values are expected to be within int32 range
	data.Name = types.StringValue(srv.Name)
	data.H323GateKeeper = types.StringPointerValue(srv.H323Gatekeeper)
	data.SNMPNetworkManagementSystem = types.StringPointerValue(srv.SNMPNetworkManagementSystem)
	data.SIPProxy = types.StringPointerValue(srv.SIPProxy)
	data.HTTPProxy = types.StringPointerValue(srv.HTTPProxy)
	data.MSSIPProxy = types.StringPointerValue(srv.MSSIPProxy)
	data.TeamsProxy = types.StringPointerValue(srv.TeamsProxy)
	data.TURNServer = types.StringPointerValue(srv.TURNServer)
	data.STUNServer = types.StringPointerValue(srv.STUNServer)
	data.UseRelayCandidatesOnly = types.BoolValue(srv.UseRelayCandidatesOnly)
	data.TranscodingLocation = types.StringPointerValue(srv.TranscodingLocation)
	data.OverflowLocation1 = types.StringPointerValue(srv.OverflowLocation1)
	data.OverflowLocation2 = types.StringPointerValue(srv.OverflowLocation2)
	data.LocalMSSIPDomain = types.StringValue(srv.LocalMSSIPDomain)
	data.PolicyServer = types.StringPointerValue(srv.PolicyServer)
	data.BDPMPINChecksEnabled = types.StringValue(srv.BDPMPinChecksEnabled)
	data.BDPMScanQuarantineEnabled = types.StringValue(srv.BDPMScanQuarantineEnabled)
	data.LiveCaptionsDialOut1 = types.StringPointerValue(srv.LiveCaptionsDialOut1)
	data.LiveCaptionsDialOut2 = types.StringPointerValue(srv.LiveCaptionsDialOut2)
	data.LiveCaptionsDialOut3 = types.StringPointerValue(srv.LiveCaptionsDialOut3)

	// Handle nullable integer fields
	if srv.MediaQoS != nil {
		data.MediaQOS = types.Int32Value(int32(*srv.MediaQoS)) // #nosec G115 -- API values are expected to be within int32 range
	}
	if srv.SignallingQoS != nil {
		data.SignallingQOS = types.Int32Value(int32(*srv.SignallingQoS)) // #nosec G115 -- API values are expected to be within int32 range
	} else {
		data.SignallingQOS = types.Int32Null()
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

	// Event Sinks
	var eventSinks []string
	for _, sink := range srv.EventSinks {
		eventSinks = append(eventSinks, fmt.Sprintf("/api/admin/configuration/v1/event_sink/%d/", sink.ID))
	}
	eventSinksSet, diags := types.SetValueFrom(ctx, types.StringType, eventSinks)
	if diags.HasError() {
		return nil, fmt.Errorf("error converting event sinks: %v", diags)
	}
	data.EventSinks = eventSinksSet

	// Client TURN Servers
	var clientTurnServers []string
	clientTurnServers = append(clientTurnServers, srv.ClientTURNServers...)
	clientTurnSet, diags := types.SetValueFrom(ctx, types.StringType, clientTurnServers)
	if diags.HasError() {
		return nil, fmt.Errorf("error converting client TURN servers: %v", diags)
	}
	data.ClientTURNServers = clientTurnSet

	// Client STUN Servers
	var clientStunServers []string
	clientStunServers = append(clientStunServers, srv.ClientSTUNServers...)
	clientStunSet, diags := types.SetValueFrom(ctx, types.StringType, clientStunServers)
	if diags.HasError() {
		return nil, fmt.Errorf("error converting client STUN servers: %v", diags)
	}
	data.ClientSTUNServers = clientStunSet

	return &data, nil
}

func (r *InfinitySystemLocationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	state := &InfinitySystemLocationResourceModel{}

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
			"Error Reading Infinity system location",
			fmt.Sprintf("Could not read Infinity system location: %s", err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *InfinitySystemLocationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	plan := &InfinitySystemLocationResourceModel{}
	state := &InfinitySystemLocationResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resourceID := int(state.ResourceID.ValueInt32())

	// Convert List attributes to []string
	dnsServers, diags := getStringList(ctx, plan.DNSServers)
	resp.Diagnostics.Append(diags...)
	ntpServers, diags := getStringList(ctx, plan.NTPServers)
	resp.Diagnostics.Append(diags...)
	syslogServers, diags := getStringList(ctx, plan.SyslogServers)
	resp.Diagnostics.Append(diags...)
	clientTurnServers, diags := getStringList(ctx, plan.ClientTURNServers)
	resp.Diagnostics.Append(diags...)
	clientStunServers, diags := getStringList(ctx, plan.ClientSTUNServers)
	resp.Diagnostics.Append(diags...)
	eventSinks, diags := getStringList(ctx, plan.EventSinks)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// initialize update request with required fields, list fields, and fields with defaults
	updateRequest := &config.SystemLocationUpdateRequest{
		Name:                      plan.Name.ValueString(),
		DNSServers:                dnsServers,
		NTPServers:                ntpServers,
		SyslogServers:             syslogServers,
		ClientTURNServers:         clientTurnServers,
		ClientSTUNServers:         clientStunServers,
		EventSinks:                eventSinks,
		Description:               plan.Description.ValueString(),
		BDPMPinChecksEnabled:      plan.BDPMPINChecksEnabled.ValueString(),
		BDPMScanQuarantineEnabled: plan.BDPMScanQuarantineEnabled.ValueString(),
		LocalMSSIPDomain:          plan.LocalMSSIPDomain.ValueString(),
		MTU:                       int(plan.MTU.ValueInt32()),
		UseRelayCandidatesOnly:    plan.UseRelayCandidatesOnly.ValueBool(),
	}

	// All nullable fields
	// Only set optional fields if they are not null in the plan
	if !plan.H323GateKeeper.IsNull() && !plan.H323GateKeeper.IsUnknown() {
		value := plan.H323GateKeeper.ValueString()
		updateRequest.H323Gatekeeper = &value
	}
	if !plan.SNMPNetworkManagementSystem.IsNull() && !plan.SNMPNetworkManagementSystem.IsUnknown() {
		value := plan.SNMPNetworkManagementSystem.ValueString()
		updateRequest.SNMPNetworkManagementSystem = &value
	}
	if !plan.SIPProxy.IsNull() && !plan.SIPProxy.IsUnknown() {
		value := plan.SIPProxy.ValueString()
		updateRequest.SIPProxy = &value
	}
	if !plan.MSSIPProxy.IsNull() && !plan.MSSIPProxy.IsUnknown() {
		value := plan.MSSIPProxy.ValueString()
		updateRequest.MSSIPProxy = &value
	}
	if !plan.HTTPProxy.IsNull() && !plan.HTTPProxy.IsUnknown() {
		value := plan.HTTPProxy.ValueString()
		updateRequest.HTTPProxy = &value
	}
	if !plan.TeamsProxy.IsNull() && !plan.TeamsProxy.IsUnknown() {
		value := plan.TeamsProxy.ValueString()
		updateRequest.TeamsProxy = &value
	}
	if !plan.TURNServer.IsNull() && !plan.TURNServer.IsUnknown() {
		value := plan.TURNServer.ValueString()
		updateRequest.TURNServer = &value
	}
	if !plan.STUNServer.IsNull() && !plan.STUNServer.IsUnknown() {
		value := plan.STUNServer.ValueString()
		updateRequest.STUNServer = &value
	}
	if !plan.MediaQOS.IsNull() && !plan.MediaQOS.IsUnknown() {
		value := int(plan.MediaQOS.ValueInt32())
		updateRequest.MediaQoS = &value
	}
	if !plan.SignallingQOS.IsNull() && !plan.SignallingQOS.IsUnknown() {
		value := int(plan.SignallingQOS.ValueInt32())
		updateRequest.SignallingQoS = &value
	}
	if !plan.TranscodingLocation.IsNull() && !plan.TranscodingLocation.IsUnknown() {
		value := plan.TranscodingLocation.ValueString()
		updateRequest.TranscodingLocation = &value
	}
	if !plan.OverflowLocation1.IsNull() && !plan.OverflowLocation1.IsUnknown() {
		value := plan.OverflowLocation1.ValueString()
		updateRequest.OverflowLocation1 = &value
	}
	if !plan.OverflowLocation2.IsNull() && !plan.OverflowLocation2.IsUnknown() {
		value := plan.OverflowLocation2.ValueString()
		updateRequest.OverflowLocation2 = &value
	}
	if !plan.PolicyServer.IsNull() && !plan.PolicyServer.IsUnknown() {
		value := plan.PolicyServer.ValueString()
		updateRequest.PolicyServer = &value
	}
	if !plan.LiveCaptionsDialOut1.IsNull() && !plan.LiveCaptionsDialOut1.IsUnknown() {
		value := plan.LiveCaptionsDialOut1.ValueString()
		updateRequest.LiveCaptionsDialOut1 = &value
	}
	if !plan.LiveCaptionsDialOut2.IsNull() && !plan.LiveCaptionsDialOut2.IsUnknown() {
		value := plan.LiveCaptionsDialOut2.ValueString()
		updateRequest.LiveCaptionsDialOut2 = &value
	}
	if !plan.LiveCaptionsDialOut3.IsNull() && !plan.LiveCaptionsDialOut3.IsUnknown() {
		value := plan.LiveCaptionsDialOut3.ValueString()
		updateRequest.LiveCaptionsDialOut3 = &value
	}

	_, err := r.InfinityClient.Config().UpdateSystemLocation(ctx, resourceID, updateRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Infinity system location",
			fmt.Sprintf("Could not update Infinity system location with ID %d: %s", resourceID, err),
		)
		return
	}

	// Re-read the resource to get the latest state
	updatedModel, err := r.read(ctx, resourceID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Updated Infinity system location",
			fmt.Sprintf("Could not read updated Infinity system location with ID %d: %s", resourceID, err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, updatedModel)...)
}

func (r *InfinitySystemLocationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	state := &InfinitySystemLocationResourceModel{}

	tflog.Info(ctx, "Deleting Infinity system location")

	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.InfinityClient.Config().DeleteSystemLocation(ctx, int(state.ResourceID.ValueInt32()))

	// Ignore 404 Not Found and Lookup errors on delete
	if err != nil && !isNotFoundError(err) && !isLookupError(err) {
		resp.Diagnostics.AddError(
			"Error Deleting Infinity system location",
			fmt.Sprintf("Could not delete Infinity system location with ID %s: %s", state.ID.ValueString(), err),
		)
		return
	}
}

func (r *InfinitySystemLocationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resourceID, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Import ID must be a valid integer for the resource ID. Got: %s", req.ID),
		)
		return
	}

	tflog.Trace(ctx, fmt.Sprintf("Importing Infinity system location with resource ID: %d", resourceID))

	// Read the resource from the API
	model, err := r.read(ctx, resourceID)
	if err != nil {
		// Check if the error is a 404 (not found)
		if isNotFoundError(err) {
			resp.Diagnostics.AddError(
				"Infinity System Location Not Found",
				fmt.Sprintf("Infinity system location with resource ID %d not found.", resourceID),
			)
			return
		}
		resp.Diagnostics.AddError(
			"Error Importing Infinity System Location",
			fmt.Sprintf("Could not import Infinity system location with resource ID %d: %s", resourceID, err),
		)
		return
	}

	// Set the state from the imported resource
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}
