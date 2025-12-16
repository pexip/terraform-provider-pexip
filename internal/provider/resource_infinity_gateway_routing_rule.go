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
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/pexip/go-infinity-sdk/v38/config"
)

var (
	_ resource.ResourceWithImportState = (*InfinityGatewayRoutingRuleResource)(nil)
)

type InfinityGatewayRoutingRuleResource struct {
	InfinityClient InfinityClient
}

type InfinityGatewayRoutingRuleResourceModel struct {
	ID                              types.String `tfsdk:"id"`
	ResourceID                      types.Int32  `tfsdk:"resource_id"`
	CallType                        types.String `tfsdk:"call_type"`
	CalledDeviceType                types.String `tfsdk:"called_device_type"`
	CryptoMode                      types.String `tfsdk:"crypto_mode"`
	DenoiseAudio                    types.Bool   `tfsdk:"denoise_audio"`
	Description                     types.String `tfsdk:"description"`
	DisabledCodecs                  types.Set    `tfsdk:"disabled_codecs"`
	Enable                          types.Bool   `tfsdk:"enable"`
	ExternalParticipantAvatarLookup types.String `tfsdk:"enable_participant_avatar_lookup"`
	GMSAccessToken                  types.String `tfsdk:"gms_access_token"`
	H323Gatekeeper                  types.String `tfsdk:"h323_gatekeeper"`
	IVRTheme                        types.String `tfsdk:"ivr_theme"`
	LiveCaptionsEnabled             types.String `tfsdk:"live_captions_enabled"`
	MatchIncomingCalls              types.Bool   `tfsdk:"match_incoming_calls"`
	MatchIncomingH323               types.Bool   `tfsdk:"match_incoming_h323"`
	MatchIncomingMSSIP              types.Bool   `tfsdk:"match_incoming_mssip"`
	MatchIncomingOnlyIfRegistered   types.Bool   `tfsdk:"match_incoming_only_if_registered"`
	MatchIncomingSIP                types.Bool   `tfsdk:"match_incoming_sip"`
	MatchIncomingTeams              types.Bool   `tfsdk:"match_incoming_teams"`
	MatchIncomingWebRTC             types.Bool   `tfsdk:"match_incoming_webrtc"`
	MatchOutgoingCalls              types.Bool   `tfsdk:"match_outgoing_calls"`
	MatchSourceLocation             types.String `tfsdk:"match_source_location"`
	MatchString                     types.String `tfsdk:"match_string"`
	MatchStringFull                 types.Bool   `tfsdk:"match_string_full"`
	MaxCallrateIn                   types.Int32  `tfsdk:"max_callrate_in"`
	MaxCallrateOut                  types.Int32  `tfsdk:"max_callrate_out"`
	MaxPixelsPerSecond              types.String `tfsdk:"max_pixels_per_second"`
	MSSIPProxy                      types.String `tfsdk:"mssip_proxy"`
	Name                            types.String `tfsdk:"name"`
	OutgoingLocation                types.String `tfsdk:"outgoing_location"`
	OutgoingProtocol                types.String `tfsdk:"outgoing_protocol"`
	Priority                        types.Int32  `tfsdk:"priority"`
	ReplaceString                   types.String `tfsdk:"replace_string"`
	SIPProxy                        types.String `tfsdk:"sip_proxy"`
	STUNServer                      types.String `tfsdk:"stun_server"`
	Tag                             types.String `tfsdk:"tag"`
	TeamsProxy                      types.String `tfsdk:"teams_proxy"`
	TelehealthProfile               types.String `tfsdk:"telehealth_profile"`
	TreatAsTrusted                  types.Bool   `tfsdk:"treat_as_trusted"`
	TURNServer                      types.String `tfsdk:"turn_server"`
}

func (r *InfinityGatewayRoutingRuleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_infinity_gateway_routing_rule"
}

func (r *InfinityGatewayRoutingRuleResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *InfinityGatewayRoutingRuleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Resource URI for the gateway routing rule in Infinity",
			},
			"resource_id": schema.Int32Attribute{
				Computed:            true,
				MarkdownDescription: "The resource integer identifier for the gateway routing rule in Infinity",
			},
			"call_type": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("video"),
				Validators: []validator.String{
					stringvalidator.OneOf("audio", "video", "video-only", "auto"),
				},
				MarkdownDescription: "Maximum media content of the call. The participant being called will not be able to escalate beyond the selected capability.",
			},
			"called_device_type": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("external"),
				Validators: []validator.String{
					stringvalidator.OneOf("external", "registration", "mssip_conference_id", "mssip_server", "gms_conference", "teams_conference", "teams_user", "telehealth_profile"),
				},
				MarkdownDescription: "The device or system to which the call is routed. The options are: Registered device or external system: routes the call to …ms Room: routes the call to a Microsoft Teams Room. Epic Telehealth meeting: routes the call to an Epic Telehealth meeting.",
			},
			"crypto_mode": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(""),
				Validators: []validator.String{
					stringvalidator.OneOf("besteffort", "on", "off"),
				},
				MarkdownDescription: "Controls the media encryption requirements for participants connecting to this service. Use global setting: Use the global media encryption setting. Encrypted media: Require encrypted media. Unencrypted media: Allow unencrypted media. (RTMP participants will use encryption if their device supports it, otherwise the connection will be unencrypted.)",
			},
			"denoise_audio": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
				MarkdownDescription: "Whether to remove background noise from audio streams as they pass through the infrastructure.",
			},
			"description": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(250),
				},
				MarkdownDescription: "A description of the call routing rule. Maximum length: 250 characters.",
			},
			"disabled_codecs": schema.SetAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Set of disabled codecs.",
			},
			"enable": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
				MarkdownDescription: "Determines whether or not the rule is enabled. Any disabled rules still appear in the rules list but are ignored. Use this setting to test configuration changes, or to temporarily disable specific rules.",
			},
			"enable_participant_avatar_lookup": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("default"),
				Validators: []validator.String{
					stringvalidator.OneOf("default", "yes", "no"),
				},
				MarkdownDescription: "Determines whether or not avatars for external participants will be retrieved using the method appropriate for the external meeting type. You can use this option to override the global configuration setting.",
			},
			"gms_access_token": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The access token to use when resolving Google Meet meeting codes.",
			},
			"h323_gatekeeper": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "When calling an external system, this is the H.323 gatekeeper to use for outbound H.323 calls. Select Use DNS to try to use normal H.323 resolution procedures to route the call.",
			},
			"ivr_theme": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "When calling an external system, this is the H.323 gatekeeper to use for outbound H.323 calls. Select Use DNS to try to use normal H.323 resolution procedures to route the call.",
			},
			"live_captions_enabled": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("default"),
				Validators: []validator.String{
					stringvalidator.OneOf("default", "yes", "no"),
				},
				MarkdownDescription: "Select whether to enable, disable, or use the default global live captions setting for this service.",
			},
			"match_incoming_calls": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
				MarkdownDescription: "Applies this rule to incoming calls that have not been routed to a Virtual Meeting Room or Virtual Reception, and should be routed via the Pexip Distributed Gateway service.",
			},
			"match_incoming_h323": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
				MarkdownDescription: "Select whether this rule should apply to incoming H.323 calls.",
			},
			"match_incoming_mssip": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
				MarkdownDescription: "Select whether this rule should apply to incoming Lync / Skype for Business (MS-SIP) calls.",
			},
			"match_incoming_only_if_registered": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Only apply this rule to incoming calls from devices, videoconferencing endpoints, soft clients or Infinity Connect clients t…aced from non-registered clients or devices, or from the Infinity web app will not be routed by this rule if it is enabled.",
			},
			"match_incoming_sip": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
				MarkdownDescription: "Only apply this rule to incoming calls from devices, videoconferencing endpoints, soft clients or Infinity Connect clients t…aced from non-registered clients or devices, or from the Infinity web app will not be routed by this rule if it is enabled.",
			},
			"match_incoming_teams": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Select whether this rule should apply to incoming Teams calls.",
			},
			"match_incoming_webrtc": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
				MarkdownDescription: "Select whether this rule should apply to incoming calls from Infinity Connect clients (WebRTC / RTMP).",
			},
			"match_outgoing_calls": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Applies this rule to outgoing calls placed from a conference service (e.g. when adding a participant to a Virtual Meeting Room) where Automatic routing has been selected.",
			},
			"match_source_location": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Applies the rule only if the incoming call is being handled by a Conferencing Node in the selected location or the outgoing call is being initiated from the selected location. To apply the rule regardless of the location, select Any Location.",
			},
			"match_string": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(250),
				},
				MarkdownDescription: "The regular expression that the destination alias (the alias that was dialed) is checked against to see if this rule applies to this call. Maximum length: 250 characters.",
			},
			"match_string_full": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "This setting is for advanced use cases and will not normally be required. By default, Pexip Infinity matches against a parsed version of the destination alias (e.g. if the alias is \"alice@example.com\", by default the rule will match against \"alice@example.com\". Select this option to match against the full, unparsed alias instead.",
			},
			"max_callrate_in": schema.Int32Attribute{
				Optional: true,
				Computed: true,
				Validators: []validator.Int32{
					int32validator.Between(128, 8192),
				},
				MarkdownDescription: "This optional field allows you to limit the bandwidth of media being received by Pexip Infinity from each individual participant dialed in via this Call Routing Rule. Range: 128 to 8192.",
			},
			"max_callrate_out": schema.Int32Attribute{
				Optional: true,
				Computed: true,
				Validators: []validator.Int32{
					int32validator.Between(128, 8192),
				},
				MarkdownDescription: "This optional field allows you to limit the bandwidth of media being sent by Pexip Infinity to each individual participant dialed out from this Call Routing Rule. Range: 128 to 8192. Default: 4128.",
			},
			"max_pixels_per_second": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.OneOf("sd", "hd", "fullhd"),
				},
				MarkdownDescription: "Sets the maximum call quality for each participant.",
			},
			"mssip_proxy": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "When calling an external system, this is the Lync / Skype for Business server to use for outbound Lync / Skype for Business (MS-SIP) calls. Select Use DNS to try to use normal Lync / Skype for Business (MS-SIP) resolution procedures to route the call. When calling a Lync / Skype for Business meeting, this is the Lync / Skype for Business server to use for the Conference ID lookup and to place the call.",
			},
			"name": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(250),
				},
				MarkdownDescription: "The name used to refer to this Call Routing Rule. Maximum length: 250 characters.",
			},
			"outgoing_location": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "When calling an external system, this forces the outgoing call to be handled by a Conferencing Node in a specific location. When calling a Lync / Skype for Business meeting, a Conferencing Node in this location will handle the outgoing call, and - for 'Lync / Skype for Business meeting direct' targets - perform the Conference ID lookup on the Lync / Skype for Business server. Select Automatic to allow Pexip Infinity to automatically select which Conferencing Node to use.",
			},
			"outgoing_protocol": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("sip"),
				Validators: []validator.String{
					stringvalidator.OneOf("h323", "mssip", "sip", "rtmp", "gms", "teams"),
				},
				MarkdownDescription: "When calling an external system, this is the protocol to use when placing the outbound call. Note that if the call is to a registered device, Pexip Infinity will instead use the protocol that the device used to make the registration.",
			},
			"priority": schema.Int32Attribute{
				Required: true,
				Validators: []validator.Int32{
					int32validator.AtLeast(1),
				},
				MarkdownDescription: "The priority of this rule. Rules are checked in ascending priority order until the first matching rule is found, and it is then applied. Range: 1 to 200.",
			},
			"replace_string": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(""),
				Validators: []validator.String{
					stringvalidator.LengthAtMost(250),
				},
				MarkdownDescription: "The regular expression string used to transform the originally dialed alias (if a match was found). Leave blank to leave the originally dialed alias unchanged. Maximum length: 250 characters.",
			},
			"sip_proxy": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "When calling an external system, this is the SIP proxy to use for outbound SIP calls. Select Use DNS to try to use normal SIP resolution procedures to route the call.",
			},
			"stun_server": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The STUN server to be used for outbound Lync / Skype for Business (MS-SIP) calls (where applicable).",
			},
			"tag": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
				MarkdownDescription: "A unique identifier used to track usage of this Call Routing Rule. Maximum length: 250 characters.",
			},
			"teams_proxy": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The Teams Connector to use for the Teams meeting. If you do not specify anything, the Teams Connector associated with the outgoing location is used.",
			},
			"telehealth_profile": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The Telehealth Profile to use for the meeting.",
			},
			"treat_as_trusted": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "This indicates the target of this routing rule will treat the caller as part of the target organization for trust purposes.",
			},
			"turn_server": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The TURN server to be used for outbound Lync / Skype for Business (MS-SIP) calls (where applicable).",
			},
		},
		MarkdownDescription: "Manages a gateway routing rule configuration with the Infinity service.",
	}
}

func (r *InfinityGatewayRoutingRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	plan := &InfinityGatewayRoutingRuleResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createRequest := &config.GatewayRoutingRuleCreateRequest{
		Name:                          plan.Name.ValueString(),
		DenoiseAudio:                  plan.DenoiseAudio.ValueBool(),
		Description:                   plan.Description.ValueString(),
		Enable:                        plan.Enable.ValueBool(),
		CalledDeviceType:              plan.CalledDeviceType.ValueString(),
		CallType:                      plan.CallType.ValueString(),
		LiveCaptionsEnabled:           plan.LiveCaptionsEnabled.ValueString(),
		MatchIncomingCalls:            plan.MatchIncomingCalls.ValueBool(),
		MatchOutgoingCalls:            plan.MatchOutgoingCalls.ValueBool(),
		MatchIncomingSIP:              plan.MatchIncomingSIP.ValueBool(),
		MatchIncomingH323:             plan.MatchIncomingH323.ValueBool(),
		MatchIncomingMSSIP:            plan.MatchIncomingMSSIP.ValueBool(),
		MatchIncomingWebRTC:           plan.MatchIncomingWebRTC.ValueBool(),
		MatchIncomingTeams:            plan.MatchIncomingTeams.ValueBool(),
		MatchIncomingOnlyIfRegistered: plan.MatchIncomingOnlyIfRegistered.ValueBool(),
		MatchString:                   plan.MatchString.ValueString(),
		MatchStringFull:               plan.MatchStringFull.ValueBool(),
		OutgoingProtocol:              plan.OutgoingProtocol.ValueString(),
		Priority:                      int(plan.Priority.ValueInt32()),
		ReplaceString:                 plan.ReplaceString.ValueString(),
		Tag:                           plan.Tag.ValueString(),
		TreatAsTrusted:                plan.TreatAsTrusted.ValueBool(),
	}

	// All nullable fields
	// Only set optional fields if they are not null in the plan
	if !plan.CryptoMode.IsNull() && !plan.CryptoMode.IsUnknown() {
		value := plan.CryptoMode.ValueString()
		createRequest.CryptoMode = &value
	}
	if !plan.DisabledCodecs.IsNull() && !plan.DisabledCodecs.IsUnknown() {
		var disabledCodecs []config.CodecValue
		for _, v := range plan.DisabledCodecs.Elements() {
			disabledCodecs = append(disabledCodecs, config.CodecValue{Value: v.(types.String).ValueString()})
		}
		createRequest.DisabledCodecs = &disabledCodecs
	}
	if !plan.ExternalParticipantAvatarLookup.IsNull() && !plan.ExternalParticipantAvatarLookup.IsUnknown() {
		value := plan.ExternalParticipantAvatarLookup.ValueString()
		createRequest.ExternalParticipantAvatarLookup = &value
	}
	if !plan.GMSAccessToken.IsNull() && !plan.GMSAccessToken.IsUnknown() {
		value := plan.GMSAccessToken.ValueString()
		createRequest.GMSAccessToken = &value
	}
	if !plan.H323Gatekeeper.IsNull() && !plan.H323Gatekeeper.IsUnknown() {
		value := plan.H323Gatekeeper.ValueString()
		createRequest.H323Gatekeeper = &value
	}
	if !plan.IVRTheme.IsNull() && !plan.IVRTheme.IsUnknown() {
		value := plan.IVRTheme.ValueString()
		createRequest.IVRTheme = &value
	}
	if !plan.MatchSourceLocation.IsNull() && !plan.MatchSourceLocation.IsUnknown() {
		value := plan.MatchSourceLocation.ValueString()
		createRequest.MatchSourceLocation = &value
	}
	if !plan.MaxPixelsPerSecond.IsNull() && !plan.MaxPixelsPerSecond.IsUnknown() {
		value := plan.MaxPixelsPerSecond.ValueString()
		createRequest.MaxPixelsPerSecond = &value
	}
	if !plan.MaxCallrateIn.IsNull() && !plan.MaxCallrateIn.IsUnknown() {
		value := int(plan.MaxCallrateIn.ValueInt32())
		createRequest.MaxCallrateIn = &value
	}
	if !plan.MaxCallrateOut.IsNull() && !plan.MaxCallrateOut.IsUnknown() {
		value := int(plan.MaxCallrateOut.ValueInt32())
		createRequest.MaxCallrateOut = &value
	}
	if !plan.MSSIPProxy.IsNull() && !plan.MSSIPProxy.IsUnknown() {
		value := plan.MSSIPProxy.ValueString()
		createRequest.MSSIPProxy = &value
	}
	if !plan.OutgoingLocation.IsNull() && !plan.OutgoingLocation.IsUnknown() {
		value := plan.OutgoingLocation.ValueString()
		createRequest.OutgoingLocation = &value
	}
	if !plan.SIPProxy.IsNull() && !plan.SIPProxy.IsUnknown() {
		value := plan.SIPProxy.ValueString()
		createRequest.SIPProxy = &value
	}
	if !plan.STUNServer.IsNull() && !plan.STUNServer.IsUnknown() {
		value := plan.STUNServer.ValueString()
		createRequest.STUNServer = &value
	}
	if !plan.TeamsProxy.IsNull() && !plan.TeamsProxy.IsUnknown() {
		value := plan.TeamsProxy.ValueString()
		createRequest.TeamsProxy = &value
	}
	if !plan.TelehealthProfile.IsNull() && !plan.TelehealthProfile.IsUnknown() {
		value := plan.TelehealthProfile.ValueString()
		createRequest.TelehealthProfile = &value
	}
	if !plan.TURNServer.IsNull() && !plan.TURNServer.IsUnknown() {
		value := plan.TURNServer.ValueString()
		createRequest.TURNServer = &value
	}

	createResponse, err := r.InfinityClient.Config().CreateGatewayRoutingRule(ctx, createRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Infinity gateway routing rule",
			fmt.Sprintf("Could not create Infinity gateway routing rule: %s", err),
		)
		return
	}

	resourceID, err := createResponse.ResourceID()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Retrieving Infinity gateway routing rule ID",
			fmt.Sprintf("Could not retrieve ID for created Infinity gateway routing rule: %s", err),
		)
		return
	}

	// Read the state from the API to get all computed values
	model, err := r.read(ctx, resourceID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Created Infinity gateway routing rule",
			fmt.Sprintf("Could not read created Infinity gateway routing rule with ID %d: %s", resourceID, err),
		)
		return
	}
	tflog.Trace(ctx, fmt.Sprintf("created Infinity gateway routing rule with ID: %s, name: %s", model.ID, model.Name))

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *InfinityGatewayRoutingRuleResource) read(ctx context.Context, resourceID int) (*InfinityGatewayRoutingRuleResourceModel, error) {
	var data InfinityGatewayRoutingRuleResourceModel

	srv, err := r.InfinityClient.Config().GetGatewayRoutingRule(ctx, resourceID)
	if err != nil {
		return nil, err
	}

	if srv.ResourceURI == "" {
		return nil, fmt.Errorf("gateway routing rule with ID %d not found", resourceID)
	}

	data.ID = types.StringValue(srv.ResourceURI)
	data.ResourceID = types.Int32Value(int32(resourceID))
	data.CallType = types.StringValue(srv.CallType)
	data.CalledDeviceType = types.StringValue(srv.CalledDeviceType)
	data.CryptoMode = types.StringPointerValue(srv.CryptoMode)
	data.DenoiseAudio = types.BoolValue(srv.DenoiseAudio)
	data.Description = types.StringValue(srv.Description)
	data.Enable = types.BoolValue(srv.Enable)
	data.ExternalParticipantAvatarLookup = types.StringPointerValue(srv.ExternalParticipantAvatarLookup)
	data.GMSAccessToken = types.StringPointerValue(srv.GMSAccessToken)
	data.H323Gatekeeper = types.StringPointerValue(srv.H323Gatekeeper)
	data.IVRTheme = types.StringPointerValue(srv.IVRTheme)
	data.LiveCaptionsEnabled = types.StringValue(srv.LiveCaptionsEnabled)
	data.MatchIncomingCalls = types.BoolValue(srv.MatchIncomingCalls)
	data.MatchIncomingH323 = types.BoolValue(srv.MatchIncomingH323)
	data.MatchIncomingMSSIP = types.BoolValue(srv.MatchIncomingMSSIP)
	data.MatchIncomingOnlyIfRegistered = types.BoolValue(srv.MatchIncomingOnlyIfRegistered)
	data.MatchIncomingSIP = types.BoolValue(srv.MatchIncomingSIP)
	data.MatchIncomingTeams = types.BoolValue(srv.MatchIncomingTeams)
	data.MatchIncomingWebRTC = types.BoolValue(srv.MatchIncomingWebRTC)
	data.MatchOutgoingCalls = types.BoolValue(srv.MatchOutgoingCalls)
	data.MatchSourceLocation = types.StringPointerValue(srv.MatchSourceLocation)
	data.MatchString = types.StringValue(srv.MatchString)
	data.MatchStringFull = types.BoolValue(srv.MatchStringFull)
	data.MaxPixelsPerSecond = types.StringPointerValue(srv.MaxPixelsPerSecond)
	data.MSSIPProxy = types.StringPointerValue(srv.MSSIPProxy)
	data.Name = types.StringValue(srv.Name)
	data.OutgoingLocation = types.StringPointerValue(srv.OutgoingLocation)
	data.OutgoingProtocol = types.StringValue(srv.OutgoingProtocol)
	data.Priority = types.Int32Value(int32(srv.Priority))
	data.ReplaceString = types.StringValue(srv.ReplaceString)
	data.SIPProxy = types.StringPointerValue(srv.SIPProxy)
	data.STUNServer = types.StringPointerValue(srv.STUNServer)
	data.Tag = types.StringValue(srv.Tag)
	data.TeamsProxy = types.StringPointerValue(srv.TeamsProxy)
	data.TelehealthProfile = types.StringPointerValue(srv.TelehealthProfile)
	data.TreatAsTrusted = types.BoolValue(srv.TreatAsTrusted)
	data.TURNServer = types.StringPointerValue(srv.TURNServer)

	// Handle nullable integer fields
	if srv.MaxCallrateIn != nil {
		data.MaxCallrateIn = types.Int32Value(int32(*srv.MaxCallrateIn)) // #nosec G115 -- API values are expected to be within int32 range
	}
	if srv.MaxCallrateOut != nil {
		data.MaxCallrateOut = types.Int32Value(int32(*srv.MaxCallrateOut)) // #nosec G115 -- API values are expected to be within int32 range
	}

	// Convert disabled codecs from SDK to Terraform format
	var disabledCodecs []attr.Value
	if srv.DisabledCodecs != nil {
		for _, v := range *srv.DisabledCodecs {
			disabledCodecs = append(disabledCodecs, types.StringValue(v.Value))
		}
	}
	data.DisabledCodecs, _ = types.SetValue(types.StringType, disabledCodecs)

	return &data, nil
}

func (r *InfinityGatewayRoutingRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	state := &InfinityGatewayRoutingRuleResourceModel{}

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
			"Error Reading Infinity gateway routing rule",
			fmt.Sprintf("Could not read Infinity gateway routing rule: %s", err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *InfinityGatewayRoutingRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	plan := &InfinityGatewayRoutingRuleResourceModel{}
	state := &InfinityGatewayRoutingRuleResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resourceID := int(state.ResourceID.ValueInt32())

	updateRequest := &config.GatewayRoutingRuleUpdateRequest{
		Name:                          plan.Name.ValueString(),
		DenoiseAudio:                  plan.DenoiseAudio.ValueBool(),
		Description:                   plan.Description.ValueString(),
		Enable:                        plan.Enable.ValueBool(),
		CalledDeviceType:              plan.CalledDeviceType.ValueString(),
		CallType:                      plan.CallType.ValueString(),
		LiveCaptionsEnabled:           plan.LiveCaptionsEnabled.ValueString(),
		MatchIncomingCalls:            plan.MatchIncomingCalls.ValueBool(),
		MatchOutgoingCalls:            plan.MatchOutgoingCalls.ValueBool(),
		MatchIncomingSIP:              plan.MatchIncomingSIP.ValueBool(),
		MatchIncomingH323:             plan.MatchIncomingH323.ValueBool(),
		MatchIncomingMSSIP:            plan.MatchIncomingMSSIP.ValueBool(),
		MatchIncomingWebRTC:           plan.MatchIncomingWebRTC.ValueBool(),
		MatchIncomingTeams:            plan.MatchIncomingTeams.ValueBool(),
		MatchIncomingOnlyIfRegistered: plan.MatchIncomingOnlyIfRegistered.ValueBool(),
		MatchString:                   plan.MatchString.ValueString(),
		MatchStringFull:               plan.MatchStringFull.ValueBool(),
		OutgoingProtocol:              plan.OutgoingProtocol.ValueString(),
		Priority:                      int(plan.Priority.ValueInt32()),
		ReplaceString:                 plan.ReplaceString.ValueString(),
		Tag:                           plan.Tag.ValueString(),
		TreatAsTrusted:                plan.TreatAsTrusted.ValueBool(),
	}

	// All nullable fields
	// Only set optional fields if they are not null in the plan
	if !plan.CryptoMode.IsNull() && !plan.CryptoMode.IsUnknown() {
		value := plan.CryptoMode.ValueString()
		updateRequest.CryptoMode = &value
	}
	if !plan.DisabledCodecs.IsNull() && !plan.DisabledCodecs.IsUnknown() {
		var disabledCodecs []config.CodecValue
		for _, v := range plan.DisabledCodecs.Elements() {
			disabledCodecs = append(disabledCodecs, config.CodecValue{Value: v.(types.String).ValueString()})
		}
		updateRequest.DisabledCodecs = &disabledCodecs
	}
	if !plan.ExternalParticipantAvatarLookup.IsNull() && !plan.ExternalParticipantAvatarLookup.IsUnknown() {
		value := plan.ExternalParticipantAvatarLookup.ValueString()
		updateRequest.ExternalParticipantAvatarLookup = &value
	}
	if !plan.GMSAccessToken.IsNull() && !plan.GMSAccessToken.IsUnknown() {
		value := plan.GMSAccessToken.ValueString()
		updateRequest.GMSAccessToken = &value
	}
	if !plan.H323Gatekeeper.IsNull() && !plan.H323Gatekeeper.IsUnknown() {
		value := plan.H323Gatekeeper.ValueString()
		updateRequest.H323Gatekeeper = &value
	}
	if !plan.IVRTheme.IsNull() && !plan.IVRTheme.IsUnknown() {
		value := plan.IVRTheme.ValueString()
		updateRequest.IVRTheme = &value
	}
	if !plan.MatchSourceLocation.IsNull() && !plan.MatchSourceLocation.IsUnknown() {
		value := plan.MatchSourceLocation.ValueString()
		updateRequest.MatchSourceLocation = &value
	}
	if !plan.MaxPixelsPerSecond.IsNull() && !plan.MaxPixelsPerSecond.IsUnknown() {
		value := plan.MaxPixelsPerSecond.ValueString()
		updateRequest.MaxPixelsPerSecond = &value
	}
	if !plan.MaxCallrateIn.IsNull() && !plan.MaxCallrateIn.IsUnknown() {
		value := int(plan.MaxCallrateIn.ValueInt32())
		updateRequest.MaxCallrateIn = &value
	}
	if !plan.MaxCallrateOut.IsNull() && !plan.MaxCallrateOut.IsUnknown() {
		value := int(plan.MaxCallrateOut.ValueInt32())
		updateRequest.MaxCallrateOut = &value
	}
	if !plan.MSSIPProxy.IsNull() && !plan.MSSIPProxy.IsUnknown() {
		value := plan.MSSIPProxy.ValueString()
		updateRequest.MSSIPProxy = &value
	}
	if !plan.OutgoingLocation.IsNull() && !plan.OutgoingLocation.IsUnknown() {
		value := plan.OutgoingLocation.ValueString()
		updateRequest.OutgoingLocation = &value
	}
	if !plan.SIPProxy.IsNull() && !plan.SIPProxy.IsUnknown() {
		value := plan.SIPProxy.ValueString()
		updateRequest.SIPProxy = &value
	}
	if !plan.STUNServer.IsNull() && !plan.STUNServer.IsUnknown() {
		value := plan.STUNServer.ValueString()
		updateRequest.STUNServer = &value
	}
	if !plan.TeamsProxy.IsNull() && !plan.TeamsProxy.IsUnknown() {
		value := plan.TeamsProxy.ValueString()
		updateRequest.TeamsProxy = &value
	}
	if !plan.TelehealthProfile.IsNull() && !plan.TelehealthProfile.IsUnknown() {
		value := plan.TelehealthProfile.ValueString()
		updateRequest.TelehealthProfile = &value
	}
	if !plan.TURNServer.IsNull() && !plan.TURNServer.IsUnknown() {
		value := plan.TURNServer.ValueString()
		updateRequest.TURNServer = &value
	}

	_, err := r.InfinityClient.Config().UpdateGatewayRoutingRule(ctx, resourceID, updateRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Infinity gateway routing rule",
			fmt.Sprintf("Could not update Infinity gateway routing rule with ID %d: %s", resourceID, err),
		)
		return
	}

	// Re-read the resource to get the latest state
	updatedModel, err := r.read(ctx, resourceID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Updated Infinity gateway routing rule",
			fmt.Sprintf("Could not read updated Infinity gateway routing rule with ID %d: %s", resourceID, err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, updatedModel)...)
}

func (r *InfinityGatewayRoutingRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	state := &InfinityGatewayRoutingRuleResourceModel{}

	tflog.Info(ctx, "Deleting Infinity gateway routing rule")

	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.InfinityClient.Config().DeleteGatewayRoutingRule(ctx, int(state.ResourceID.ValueInt32()))

	// Ignore 404 Not Found and Lookup errors on delete
	if err != nil && !isNotFoundError(err) && !isLookupError(err) {
		resp.Diagnostics.AddError(
			"Error Deleting Infinity gateway routing rule",
			fmt.Sprintf("Could not delete Infinity gateway routing rule with ID %s: %s", state.ID.ValueString(), err),
		)
		return
	}
}

func (r *InfinityGatewayRoutingRuleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resourceID, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Import ID must be a valid integer for the resource ID. Got: %s", req.ID),
		)
		return
	}

	tflog.Trace(ctx, fmt.Sprintf("Importing Infinity gateway routing rule with resource ID: %d", resourceID))

	// Read the resource from the API
	model, err := r.read(ctx, resourceID)
	if err != nil {
		// Check if the error is a 404 (not found)
		if isNotFoundError(err) {
			resp.Diagnostics.AddError(
				"Infinity Gateway Routing Rule Not Found",
				fmt.Sprintf("Infinity gateway routing rule with resource ID %d not found.", resourceID),
			)
			return
		}
		resp.Diagnostics.AddError(
			"Error Importing Infinity Gateway Routing Rule",
			fmt.Sprintf("Could not import Infinity gateway routing rule with resource ID %d: %s", resourceID, err),
		)
		return
	}

	// Set the state from the imported resource
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}
