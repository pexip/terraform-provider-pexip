package provider

import (
	"context"
	"fmt"
	"sort"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
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
	Description                 types.String `tfsdk:"description"`
	DNSServers                  types.List   `tfsdk:"dns_servers"`
	NTPServers                  types.List   `tfsdk:"ntp_servers"`
	MTU                         types.Int32  `tfsdk:"mtu"`
	SyslogServers               types.List   `tfsdk:"syslog_servers"`
	H323GateKeeper              types.String `tfsdk:"h323_gatekeeper"`
	SNMPNetworkManagementSystem types.String `tfsdk:"snmp_network_management_system"`
	SIPProxy                    types.String `tfsdk:"sip_proxy"`
	HTTPProxy                   types.String `tfsdk:"http_proxy"`
	MSSIPProxy                  types.String `tfsdk:"mssip_proxy"`
	TeamsProxy                  types.String `tfsdk:"teams_proxy"`
	TURNServer                  types.String `tfsdk:"turn_server"`
	STUNServer                  types.String `tfsdk:"stun_server"`
	ClientTURNServers           types.List   `tfsdk:"client_turn_servers"`
	ClientSTUNServers           types.List   `tfsdk:"client_stun_servers"`
	UseRelayCandidatesOnly      types.Bool   `tfsdk:"use_relay_candidates_only"`
	MediaQOS                    types.Int32  `tfsdk:"media_qos"`
	SignallingQOS               types.Int32  `tfsdk:"signalling_qos"`
	TranscodingLocation         types.String `tfsdk:"transcoding_location"`
	OverflowLocationOne         types.String `tfsdk:"overflow_location1"`
	OverflowLocationTwo         types.String `tfsdk:"overflow_location2"`
	LocalMSSIPDomain            types.String `tfsdk:"local_mssip_domain"`
	PolicyServer                types.String `tfsdk:"policy_server"`
	EventSinks                  types.List   `tfsdk:"event_sinks"`
	BDPMPINChecksEnabled        types.Bool   `tfsdk:"bdpm_pin_checks_enabled"`
	BDPMScanQuarantineEnabled   types.Bool   `tfsdk:"bdpm_scan_quarantine_enabled"`
	LiveCaptionsDialOutOne      types.String `tfsdk:"live_captions_dial_out1"`
	LiveCaptionsDialOutTwo      types.String `tfsdk:"live_captions_dial_out2"`
	LiveCaptionsDialOutThree    types.String `tfsdk:"live_captions_dial_out3"`
}

// getSortedStringList is a generic helper to convert a types.List of strings to a sorted string slice.
func getSortedStringList(ctx context.Context, list types.List) ([]string, diag.Diagnostics) {
	if list.IsNull() || list.IsUnknown() {
		return nil, nil
	}
	var items []string
	diags := list.ElementsAs(ctx, &items, false)
	if diags.HasError() {
		return nil, diags
	}
	sort.Strings(items)
	return items, diags
}

func (m *InfinitySystemLocationResourceModel) GetDNSServers(ctx context.Context) ([]string, diag.Diagnostics) {
	return getSortedStringList(ctx, m.DNSServers)
}

func (m *InfinitySystemLocationResourceModel) GetNTPServers(ctx context.Context) ([]string, diag.Diagnostics) {
	return getSortedStringList(ctx, m.NTPServers)
}

func (m *InfinitySystemLocationResourceModel) GetSyslogServers(ctx context.Context) ([]string, diag.Diagnostics) {
	return getSortedStringList(ctx, m.SyslogServers)
}

func (m *InfinitySystemLocationResourceModel) GetClientTURNServers(ctx context.Context) ([]string, diag.Diagnostics) {
	return getSortedStringList(ctx, m.ClientTURNServers)
}

func (m *InfinitySystemLocationResourceModel) GetClientSTUNServers(ctx context.Context) ([]string, diag.Diagnostics) {
	return getSortedStringList(ctx, m.ClientSTUNServers)
}

func (m *InfinitySystemLocationResourceModel) GetEventSinks(ctx context.Context) ([]string, diag.Diagnostics) {
	return getSortedStringList(ctx, m.EventSinks)
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
			"description": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(250),
				},
				MarkdownDescription: "A description of the system location. Maximum length: 250 characters.",
			},
			"dns_servers": schema.ListAttribute{
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "List of DNS server resource URIs for this system location.",
			},
			"ntp_servers": schema.ListAttribute{
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "List of NTP server resource URIs for this system location.",
			},
			"mtu": schema.Int32Attribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Maximum Transmission Unit - the size of the largest packet that can be transmitted via the network interface for this system location. It depends on your network topology as to whether you may need to specify an MTU value here. Range: 512 to 1500.",
			},
			"name": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(250),
				},
				MarkdownDescription: "The name used to refer to this system location. Maximum length: 250 characters.",
			},
			"syslog_servers": schema.ListAttribute{
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "The Syslog servers to be used by Conferencing Nodes deployed in this Location.",
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
	dnsServers, diags := plan.GetDNSServers(ctx)
	resp.Diagnostics.Append(diags...)
	ntpServers, diags := plan.GetNTPServers(ctx)
	resp.Diagnostics.Append(diags...)
	syslogServers, diags := plan.GetSyslogServers(ctx)
	resp.Diagnostics.Append(diags...)
	clientTurnServers, diags := plan.GetClientTURNServers(ctx)
	resp.Diagnostics.Append(diags...)
	clientStunServers, diags := plan.GetClientSTUNServers(ctx)
	resp.Diagnostics.Append(diags...)
	eventSinks, diags := plan.GetEventSinks(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createRequest := &config.SystemLocationCreateRequest{
		Name:              plan.Name.ValueString(),
		DNSServers:        dnsServers,
		NTPServers:        ntpServers,
		SyslogServers:     syslogServers,
		ClientTURNServers: clientTurnServers,
		ClientSTUNServers: clientStunServers,
		EventSinks:        eventSinks,
	}

	// Only set optional fields if they are not null in the plan
	if !plan.Description.IsNull() {
		createRequest.Description = plan.Description.ValueString()
	}
	if !plan.MTU.IsNull() {
		createRequest.MTU = int(plan.MTU.ValueInt32())
	}
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
	if !plan.HTTPProxy.IsNull() && !plan.HTTPProxy.IsUnknown() {
		value := plan.HTTPProxy.ValueString()
		createRequest.HTTPProxy = &value
	}
	if !plan.MSSIPProxy.IsNull() && !plan.MSSIPProxy.IsUnknown() {
		value := plan.MSSIPProxy.ValueString()
		createRequest.MSSIPProxy = &value
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
	if !plan.UseRelayCandidatesOnly.IsNull() && !plan.UseRelayCandidatesOnly.IsUnknown() {
		createRequest.UseRelayCandidatesOnly = plan.UseRelayCandidatesOnly.ValueBool()
	}
	if !plan.MediaQoS.IsNull() && !plan.MediaQoS.IsUnknown() {
		value := int(plan.MediaQoS.ValueInt32())
		createRequest.MediaQoS = &value
	}
	if !plan.SignallingQoS.IsNull() && !plan.SignallingQoS.IsUnknown() {
		value := int(plan.SignallingQoS.ValueInt32())
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
	if !plan.LocalMSSIPDomain.IsNull() && !plan.LocalMSSIPDomain.IsUnknown() {
		createRequest.LocalMSSIPDomain = plan.LocalMSSIPDomain.ValueString()
	}
	if !plan.PolicyServer.IsNull() && !plan.PolicyServer.IsUnknown() {
		value := plan.PolicyServer.ValueString()
		createRequest.PolicyServer = &value
	}
	if !plan.BDPMPinChecksEnabled.IsNull() && !plan.BDPMPinChecksEnabled.IsUnknown() {
		createRequest.BDPMPinChecksEnabled = plan.BDPMPinChecksEnabled.ValueString()
	}
	if !plan.BDPMScanQuarantineEnabled.IsNull() && !plan.BDPMScanQuarantineEnabled.IsUnknown() {
		createRequest.BDPMScanQuarantineEnabled = plan.BDPMScanQuarantineEnabled.ValueString()
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

	if len(srv.ResourceURI) == 0 {
		return nil, fmt.Errorf("system location with ID %d not found", resourceID)
	}

	data.ID = types.StringValue(srv.ResourceURI)
	data.ResourceID = types.Int32Value(int32(resourceID))
	data.Description = types.StringValue(srv.Description)
	data.MTU = types.Int32Value(int32(srv.MTU))
	data.Name = types.StringValue(srv.Name)

	// Convert DNS servers from SDK to Terraform format
	var dnsServers []string
	for _, dns := range srv.DNSServers {
		dnsServers = append(dnsServers, fmt.Sprintf("/api/admin/configuration/v1/dns_server/%d/", dns.ID))
	}
	sort.Strings(dnsServers)
	dnsListValue, diags := types.ListValueFrom(ctx, types.StringType, dnsServers)
	if diags.HasError() {
		return nil, fmt.Errorf("error converting DNS servers: %v", diags)
	}
	data.DNSServers = dnsListValue

	// Convert NTP servers from SDK to Terraform format
	var ntpServers []string
	for _, ntp := range srv.NTPServers {
		ntpServers = append(ntpServers, fmt.Sprintf("/api/admin/configuration/v1/ntp_server/%d/", ntp.ID))
	}
	sort.Strings(ntpServers)
	ntpListValue, diags := types.ListValueFrom(ctx, types.StringType, ntpServers)
	if diags.HasError() {
		return nil, fmt.Errorf("error converting NTP servers: %v", diags)
	}
	data.NTPServers = ntpListValue

	// Convert Syslog servers from SDK to Terraform format
	var syslogServers []string
	for _, syslog := range srv.SyslogServers {
		syslogServers = append(syslogServers, fmt.Sprintf("/api/admin/configuration/v1/syslog_server/%d/", syslog.ID))
	}
	sort.Strings(syslogServers)
	syslogListValue, diags := types.ListValueFrom(ctx, types.StringType, syslogServers)
	if diags.HasError() {
		return nil, fmt.Errorf("error converting Syslog servers: %v", diags)
	}
	data.SyslogServers = syslogListValue

	// Event Sinks
	var eventSinks []string
	for _, sink := range srv.EventSinks {
		eventSinks = append(eventSinks, fmt.Sprintf("/api/admin/configuration/v1/event_sink/%d/", sink.ID))
	}
	eventSinksList, diags := types.ListValueFrom(ctx, types.StringType, eventSinks)
	if diags.HasError() {
		return nil, fmt.Errorf("error converting event sinks: %v", diags)
	}
	data.EventSinks = eventSinksList

	// Client TURN Servers
	var clientTurnServers []string
	for _, turn := range srv.ClientTURNServers {
		clientTurnServers = append(clientTurnServers, fmt.Sprintf("/api/admin/configuration/v1/turn_server/%d/", turn.ID))
	}
	clientTurnList, diags := types.ListValueFrom(ctx, types.StringType, clientTurnServers)
	if diags.HasError() {
		return nil, fmt.Errorf("error converting client TURN servers: %v", diags)
	}
	data.ClientTURNServers = clientTurnList

	// Client STUN Servers
	var clientStunServers []string
	for _, stun := range srv.ClientSTUNServers {
		clientStunServers = append(clientStunServers, fmt.Sprintf("/api/admin/configuration/v1/stun_server/%d/", stun.ID))
	}
	clientStunList, diags := types.ListValueFrom(ctx, types.StringType, clientStunServers)
	if diags.HasError() {
		return nil, fmt.Errorf("error converting client STUN servers: %v", diags)
	}
	data.ClientSTUNServers = clientStunList

	

	// The following fields are set if available in srv, otherwise set to Null
	data.H323GateKeeper = types.StringPointerValue(srv.H323Gatekeeper)
	data.SNMPNetworkManagementSystem = types.StringPointerValue(srv.SNMPNetworkManagementSystem)
	data.SIPProxy = types.StringPointerValue(srv.SIPProxy)
	data.HTTPProxy = types.StringPointerValue(srv.HTTPProxy)
	data.MSSIPProxy = types.StringPointerValue(srv.MSSIPProxy)
	data.TeamsProxy = types.StringPointerValue(srv.TeamsProxy)
	data.TURNServer = types.StringPointerValue(srv.TURNServer)
	data.STUNServer = types.StringPointerValue(srv.STUNServer)

	// Booleans and Ints
	data.UseRelayCandidatesOnly = types.BoolValue(srv.UseRelayCandidatesOnly)
	data.MediaQoS = types.Int32Value(int32(*srv.MediaQoS))
	data.SignallingQoS = types.Int32Value(int32(*srv.SignallingQoS))
	data.TranscodingLocation = types.StringPointerValue(srv.TranscodingLocation)
	data.OverflowLocation1 = types.StringPointerValue(srv.OverflowLocation1)
	data.OverflowLocation2 = types.StringPointerValue(srv.OverflowLocation2)
	data.LocalMSSIPDomain = types.StringValue(srv.LocalMSSIPDomain)
	data.PolicyServer = types.StringPointerValue(srv.PolicyServer)
	data.BDPMPinChecksEnabled = types.StringValue(srv.BDPMPinChecksEnabled)
	data.BDPMScanQuarantineEnabled = types.StringValue(srv.BDPMScanQuarantineEnabled)
	data.LiveCaptionsDialOut1 = types.StringPointerValue(srv.LiveCaptionsDialOut1)
	data.LiveCaptionsDialOut2 = types.StringPointerValue(srv.LiveCaptionsDialOut2)
	data.LiveCaptionsDialOut3 = types.StringPointerValue(srv.LiveCaptionsDialOut3)

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

	dnsServers, diags := plan.GetDNSServers(ctx)
	resp.Diagnostics.Append(diags...)
	ntpServers, diags := plan.GetNTPServers(ctx)
	resp.Diagnostics.Append(diags...)
	syslogServers, diags := plan.GetSyslogServers(ctx)
	resp.Diagnostics.Append(diags...)
	clientTurnServers, diags := plan.GetClientTURNServers(ctx)
	resp.Diagnostics.Append(diags...)
	clientStunServers, diags := plan.GetClientSTUNServers(ctx)
	resp.Diagnostics.Append(diags...)
	eventSinks, diags := plan.GetEventSinks(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateRequest := &config.SystemLocationUpdateRequest{
		Name:              plan.Name.ValueString(),
		DNSServers:        dnsServers,
		NTPServers:        ntpServers,
		SyslogServers:     syslogServers,
		ClientTURNServers: clientTurnServers,
		ClientSTUNServers: clientStunServers,
		EventSinks:        eventSinks,
	}

	if !plan.Description.IsNull() {
		updateRequest.Description = plan.Description.ValueString()
	}
	if !plan.MTU.IsNull() {
		updateRequest.MTU = int(plan.MTU.ValueInt32())
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
