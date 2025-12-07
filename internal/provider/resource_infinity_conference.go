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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/pexip/go-infinity-sdk/v38/config"
)

var (
	_ resource.ResourceWithImportState = (*InfinityConferenceResource)(nil)
)

type InfinityConferenceResource struct {
	InfinityClient InfinityClient
}

type InfinityConferenceResourceModel struct {
	ID                              types.String `tfsdk:"id"`
	ResourceID                      types.Int32  `tfsdk:"resource_id"`
	Name                            types.String `tfsdk:"name"`
	Aliases                         types.Set    `tfsdk:"aliases"`
	AllowGuests                     types.Bool   `tfsdk:"allow_guests"`
	AutomaticParticipants           types.Set    `tfsdk:"automatic_participants"`
	BreakoutRooms                   types.Bool   `tfsdk:"breakout_rooms"`
	CallType                        types.String `tfsdk:"call_type"`
	CryptoMode                      types.String `tfsdk:"crypto_mode"`
	DenoiseEnabled                  types.Bool   `tfsdk:"denoise_enabled"`
	Description                     types.String `tfsdk:"description"`
	DirectMedia                     types.String `tfsdk:"direct_media"`
	DirectMediaNotificationDuration types.Int32  `tfsdk:"direct_media_notification_duration"`
	EnableActiveSpeakerIndication   types.Bool   `tfsdk:"enable_active_speaker_indication"`
	EnableChat                      types.String `tfsdk:"enable_chat"`
	EnableOverlayText               types.Bool   `tfsdk:"enable_overlay_text"`
	ForcePresenterIntoMain          types.Bool   `tfsdk:"force_presenter_into_main"`
	GMSAccessToken                  types.String `tfsdk:"gms_access_token"`
	GuestIdentityProviderGroup      types.String `tfsdk:"guest_identity_provider_group"`
	GuestPIN                        types.String `tfsdk:"guest_pin"`
	GuestView                       types.String `tfsdk:"guest_view"`
	GuestsCanPresent                types.Bool   `tfsdk:"guests_can_present"`
	GuestsCanSeeGuests              types.String `tfsdk:"guests_can_see_guests"`
	HostIdentityProviderGroup       types.String `tfsdk:"host_identity_provider_group"`
	HostView                        types.String `tfsdk:"host_view"`
	IVRTheme                        types.String `tfsdk:"ivr_theme"`
	LiveCaptionsEnabled             types.String `tfsdk:"live_captions_enabled"`
	MatchString                     types.String `tfsdk:"match_string"`
	MaxCallRateIn                   types.Int32  `tfsdk:"max_callrate_in"`
	MaxCallRateOut                  types.Int32  `tfsdk:"max_callrate_out"`
	MaxPixelsPerSecond              types.String `tfsdk:"max_pixels_per_second"`
	MediaPlaylist                   types.String `tfsdk:"media_playlist"`
	MSSIPProxy                      types.String `tfsdk:"mssip_proxy"`
	MuteAllGuests                   types.Bool   `tfsdk:"mute_all_guests"`
	NonIdpParticipants              types.String `tfsdk:"non_idp_participants"`
	OnCompletion                    types.String `tfsdk:"on_completion"`
	ParticipantLimit                types.Int32  `tfsdk:"participant_limit"`
	PIN                             types.String `tfsdk:"pin"`
	PinningConfig                   types.String `tfsdk:"pinning_config"`
	PostMatchString                 types.String `tfsdk:"post_match_string"`
	PostReplaceString               types.String `tfsdk:"post_replace_string"`
	PrimaryOwnerEmailAddress        types.String `tfsdk:"primary_owner_email_address"`
	ReplaceString                   types.String `tfsdk:"replace_string"`
	ServiceType                     types.String `tfsdk:"service_type"`
	SoftmuteEnabled                 types.Bool   `tfsdk:"softmute_enabled"`
	SyncTag                         types.String `tfsdk:"sync_tag"`
	SystemLocation                  types.String `tfsdk:"system_location"`
	Tag                             types.String `tfsdk:"tag"`
	TeamsProxy                      types.String `tfsdk:"teams_proxy"`
	TwoStageDialType                types.String `tfsdk:"two_stage_dial_type"`
	//ScheduledConferences            types.Set    `tfsdk:"scheduled_conferences"`
	//ScheduledConferencesCount       types.Int32  `tfsdk:"scheduled_conferences_count"` # Read-only field
}

func (r *InfinityConferenceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_infinity_conference"
}

func (r *InfinityConferenceResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *InfinityConferenceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Resource URI for the conference in Infinity",
			},
			"resource_id": schema.Int32Attribute{
				Computed:            true,
				MarkdownDescription: "The resource integer identifier for the conference in Infinity",
			},
			"name": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(250),
				},
				MarkdownDescription: "The unique name used to refer to this conference. Maximum length: 250 characters.",
			},
			"aliases": schema.SetAttribute{
				Optional:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "The aliases associated with this conference.",
			},
			"allow_guests": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Whether Guest participants are allowed to join the conference. true: the conference will have two types of participants: Hosts and Guests. You must enter a PIN to be used by the Hosts. You can optionally enter a Guest PIN; if you do not enter a Guest PIN, Guests can join without a PIN, but the meeting will not start until the first Host has joined. false: all participants will have Host privileges.",
			},
			"automatic_participants": schema.SetAttribute{
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "When a conference begins, a call will be placed automatically to these selected participants.",
			},
			"breakout_rooms": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Allow this service to use different breakout rooms for participants.",
			},
			"call_type": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("video"),
				Validators: []validator.String{
					stringvalidator.OneOf("audio", "video", "video-only"),
				},
				MarkdownDescription: "Maximum media content of the conference. Participants will not be able to escalate beyond the selected capability.",
			},
			"crypto_mode": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.OneOf("besteffort", "on", "off"),
				},
				MarkdownDescription: "Controls the media encryption requirements for participants connecting to this service. Use global setting: Use the global media encryption setting (Platform > Global Settings). Required: All participants (including RTMP participants) must use media encryption. Best effort: Each participant will use media encryption if their device supports it, otherwise the connection will be unencrypted. No encryption: All H.323, SIP and MS-SIP participants must use unencrypted media. (RTMP participants will use encryption if their device supports it, otherwise the connection will be unencrypted.)",
			},
			"denoise_enabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "If enabled, all noisy participants will have noise removed from speech.",
			},
			"description": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(""),
				Validators: []validator.String{
					stringvalidator.LengthAtMost(250),
				},
				MarkdownDescription: "A description of the conference. Maximum length: 250 characters.",
			},
			"direct_media": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("never"),
				Validators: []validator.String{
					stringvalidator.OneOf("never", "best_effort", "always"),
				},
				MarkdownDescription: "Allows this VMR to use direct media between participants. When enabled, the VMR provides non-transcoded, encrypted, point-to-point calls between any two WebRTC participants.",
			},
			"direct_media_notification_duration": schema.Int32Attribute{
				Optional: true,
				Computed: true,
				Default:  int32default.StaticInt32(0),
				Validators: []validator.Int32{
					int32validator.Between(0, 30),
				},
				MarkdownDescription: "The number of seconds to show a notification before being escalated into a transcoded call, or de-escalated into a direct media call. Range: 0s to 30s. Default: 0s.",
			},
			"enable_active_speaker_indication": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Speaker display name or alias is shown across the bottom of their video image.",
			},
			"enable_chat": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("default"),
				Validators: []validator.String{
					stringvalidator.OneOf("default", "yes", "no"),
				},
				MarkdownDescription: "Enables relay of chat messages between conference participants using Lync / Skype for Business and Infinity Connect clients. You can use this option to override the global configuration setting.",
			},
			"enable_overlay_text": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "If enabled, the display name or alias will be shown for each participant.",
			},
			"force_presenter_into_main": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "If enabled, a Host sending a presentation stream will always hold the main video position, instead of being voice-switched.",
			},
			"gms_access_token": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Select an access token to use to resolve Google Meet meeting codes.",
			},
			"guests_can_present": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
				MarkdownDescription: "If enabled, Guests and Hosts can present into the conference. If disabled, only Hosts can present.",
			},
			"guest_identity_provider_group": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Select the set of Identity Providers to be offered to Guests to authenticate with, in order to join the conference. If this is blank, Guests will not be required to authenticate. ",
			},
			"guest_pin": schema.StringAttribute{
				Optional:  true,
				Computed:  true,
				Sensitive: true,
				Default:   stringdefault.StaticString(""),
				Validators: []validator.String{
					stringvalidator.LengthBetween(4, 20),
				},
				MarkdownDescription: "This optional field allows you to set a secure access code for Guest participants who dial in to the service. Length: 4-20 digits, including any terminal #.",
			},
			"guest_view": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.OneOf("one_main_zero_pips", "one_main_seven_pips", "one_main_twentyone_pips", "two_mains_twentyone_pips", "four_mains_zero_pips", "nine_mains_zero_pips", "sixteen_mains_zero_pips", "twentyfive_mains_zero_pips", "five_mains_seven_pips"),
				},
				MarkdownDescription: "The layout that Guests will see. Guests only see Host participants. one_main_zero_pips: full-screen main speaker only. one_main_seven_pips: large main speaker and up to 7 other participants. one_main_twentyone_pips: main speaker and up to 21 other participants. two_mains_twentyone_pips: two main speakers and up to 21 other participants. four_mains_zero_pips: up to four main speakers, in a 2x2 layout. nine_mains_zero_pips: up to nine main speakers, in a 3x3 layout. sixteen_mains_zero_pips: up to sixteen main speakers, in a 4x4 layout. twentyfive_mains_zero_pips: up to twenty five main speakers, in a 5x5 layout. five_mains_seven_pips: Adaptive Composition layout (does not apply to service_type of 'lecture').",
			},
			"guests_can_see_guests": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("no_hosts"),
				Validators: []validator.String{
					stringvalidator.OneOf("no_hosts", "never", "always"),
				},
				MarkdownDescription: "If enabled, when the Host leaves the conference Guests will see other Guests. If disabled, Guests will see a splash screen.",
			},
			"host_identity_provider_group": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Select the set of Identity Providers to be offered to Hosts to authenticate with, in order to join the conference. If this is blank, Hosts will not be required to authenticate.",
			},
			"host_view": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.OneOf("one_main_zero_pips", "one_main_seven_pips", "one_main_twentyone_pips", "two_mains_twentyone_pips", "four_mains_zero_pips", "nine_mains_zero_pips", "sixteen_mains_zero_pips", "twentyfive_mains_zero_pips", "five_mains_seven_pips"),
				},
				MarkdownDescription: "The layout that Hosts will see. one_main_zero_pips: full-screen main speaker only. one_main_seven_pips: large main speaker and up to 7 other participants. one_main_twentyone_pips: main speaker and up to 21 other participants. two_mains_twentyone_pips: two main speakers and up to 21 other participants. four_mains_zero_pips: up to four main speakers in a 2x2 layout. nine_mains_zero_pips: up to nine main speakers, in a 3x3 layout. sixteen_mains_zero_pips: up to sixteen main speakers, in a 4x4 layout. twentyfive_mains_zero_pips: up to twenty five main speakers, in a 5x5 layout. five_mains_seven_pips: Adaptive Composition layout (does not apply to service_type of 'lecture').",
			},
			"ivr_theme": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The IVR theme to use for this conference.",
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
			"match_string": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
				MarkdownDescription: "An optional regular expression used to match against the alias entered by the caller into the Virtual Reception. If the entered alias does not match the expression, the Virtual Reception will not route the call. If this field is left blank, any entered alias is permitted. Maximum length: 250 characters.",
			},
			"max_callrate_in": schema.Int32Attribute{
				Optional: true,
				Computed: true,
				Validators: []validator.Int32{
					int32validator.Between(128, 8192),
				},
				MarkdownDescription: "This optional field allows you to limit the bandwidth of media being received by Pexip Infinity from each individual participant dialed in to this service. Range: 128 to 8192.",
			},
			"max_callrate_out": schema.Int32Attribute{
				Optional: true,
				Computed: true,
				Validators: []validator.Int32{
					int32validator.Between(128, 8192),
				},
				MarkdownDescription: "This optional field allows you to limit the bandwidth of media being sent by Pexip Infinity to each individual participant dialed in to this service. Range: 128 to 8192.",
			},
			"max_pixels_per_second": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.OneOf("sd", "hd", "fullhd"),
				},
				MarkdownDescription: "Sets the maximum call quality for each participant. Valid choices: sd, hd, fullhd.",
			},
			"media_playlist": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The playlist to run when this Media Playback Service is used.",
			},
			"mssip_proxy": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The Lync / Skype for Business server to use to resolve the Lync / Skype for Business Conference ID entered by the user. You must then ensure that your Call Routing Rule that routes calls to your Lync / Skype for Business environment has Match against full alias URI selected and a Destination alias regex match in the style .+@example.com.",
			},
			"mute_all_guests": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "If enabled, all Guest participants will be muted by default.",
			},
			"non_idp_participants": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("disallow_all"),
				Validators: []validator.String{
					stringvalidator.OneOf("allow_if_trusted", "disallow_all"),
				},
				MarkdownDescription: "Determines whether participants attempting to join from devices other than the Infinity Connect apps (for example, SIP or H.323 endpoints) are permitted to join the conference when authentication is required. Disallow all: these devices may not join the conference. Allow if trusted: these devices may join the conference if they are locally registered.",
			},
			"on_completion": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "JSON format is used to specify what happens when the playlist finishes. If omitted, the last video's final frame remains in …the specified alias, for example, a VMR. Role is optional and can be auto, host, or guest. If omitted, the default is auto.",
			},
			"participant_limit": schema.Int32Attribute{
				Optional: true,
				Computed: true,
				Validators: []validator.Int32{
					int32validator.Between(0, 1000000),
				},
				MarkdownDescription: "This optional field allows you to limit the number of participants allowed to join this Virtual Meeting Room. Range: 0 to 1000000.",
			},
			"pin": schema.StringAttribute{
				Optional:  true,
				Computed:  true,
				Sensitive: true,
				Default:   stringdefault.StaticString(""),
				Validators: []validator.String{
					stringvalidator.LengthBetween(4, 20),
				},
				MarkdownDescription: "This optional field allows you to set a secure access code for participants who dial in to the service. Length: 4-20 digits, including any terminal #.",
			},
			"pinning_config": schema.StringAttribute{
				Optional: true,
				Computed: true,
				MarkdownDescription: "The layout pinning configuration that will be used for this conference.",
			},
			"post_match_string": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(""),
				Validators: []validator.String{
					stringvalidator.LengthAtMost(250),
				},
				MarkdownDescription: "An optional regular expression used to match against the meeting code, after the service lookup has been performed. Maximum length: 250 characters.",
			},
			"post_replace_string": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(""),
				Validators: []validator.String{
					stringvalidator.LengthAtMost(250),
				},
				MarkdownDescription: "An optional regular expression used to transform the meeting code so that, for example, it will match a Call Routing Rule fo…s if a Post-lookup regex match is also configured and the meeting code matches that regex.) Maximum length: 250 characters.",
			},
			"primary_owner_email_address": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(""),
				Validators: []validator.String{
					stringvalidator.LengthAtMost(100),
				},
				MarkdownDescription: "The email address of the owner of the VMR. Maximum length: 100 characters.",
			},
			"replace_string": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(""),
				Validators: []validator.String{
					stringvalidator.LengthAtMost(250),
				},
				MarkdownDescription: "An optional regular expression used to transform the alias entered by the caller into the Virtual Reception. (Only applies if a regex match string is also configured and the entered alias matches that regex.) Leave this field blank if you do not want to change the alias entered by the caller. Maximum length: 250 characters.",
			},
			// Field will change regularly, so commenting out for now
			//"scheduled_conferences": schema.ListAttribute{
			//	Optional:            true,
			//	Computed:            true,
			//	ElementType:         types.StringType,
			//	MarkdownDescription: "The scheduled conferences associated with this conference.",
			//},
			// Read-only field
			//"scheduled_conferences_count": schema.Int32Attribute{
			//	Optional:            true,
			//	Computed:            true,
			//	MarkdownDescription: "The number of scheduled conferences associated with this conference.",
			//},
			"service_type": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("conference"),
				Validators: []validator.String{
					stringvalidator.OneOf("conference", "lecture", "two_stage_dialing", "test_call", "media_playback"),
				},
				MarkdownDescription: "The type of conferencing service. Valid choices: conference, lecture, two_stage_dialing, test_call, media_playback.",
			},
			"softmute_enabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "If enabled, noisy participants will have reduced volume until starting to speak.",
			},
			"sync_tag": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(""),
				Validators: []validator.String{
					stringvalidator.LengthAtMost(250),
				},
				MarkdownDescription: "A read-only identifier used by the system to track synchronization of this Virtual Meeting Room with LDAP. Maximum length: 250 characters.",
			},
			"system_location": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "If selected, a Conferencing Node in this system location will perform the Lync / Skype for Business Conference ID lookup on the Lync / Skype for Business server. If a location is not selected, the IVR ingress node will perform the lookup.",
			},
			"tag": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(""),
				Validators: []validator.String{
					stringvalidator.LengthAtMost(250),
				},
				MarkdownDescription: "A unique identifier used to track usage of this service. Maximum length: 250 characters.",
			},
			"teams_proxy": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The Teams Connector to use to resolve the Conference ID entered by the user.",
			},
			"two_stage_dial_type": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("regular"),
				Validators: []validator.String{
					stringvalidator.OneOf("regular", "mssip", "gms", "teams"),
				},
				MarkdownDescription: "The type of this Virtual Reception. Select Lync / Skype for Business if this Virtual Reception is to act as an IVR gateway to scheduled and ad hoc Lync / Skype for Business meetings. Select Google Meet if this Virtual Reception is to act as an IVR gateway to Google Meet meetings. Skype for Business meetings. Otherwise, select Regular.",
			},
			// Add other fields as needed
		},
		MarkdownDescription: "Manages a conference configuration with the Infinity service.",
	}
}

func (r *InfinityConferenceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	plan := &InfinityConferenceResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// initialize create request with required fields, list fields, and fields with defaults
	createRequest := &config.ConferenceCreateRequest{
		Name:                            plan.Name.ValueString(),
		AllowGuests:                     plan.AllowGuests.ValueBool(),
		BreakoutRooms:                   plan.BreakoutRooms.ValueBool(),
		CallType:                        plan.CallType.ValueString(),
		DenoiseEnabled:                  plan.DenoiseEnabled.ValueBool(),
		Description:                     plan.Description.ValueString(),
		DirectMedia:                     plan.DirectMedia.ValueString(),
		DirectMediaNotificationDuration: int(plan.DirectMediaNotificationDuration.ValueInt32()),
		EnableActiveSpeakerIndication:   plan.EnableActiveSpeakerIndication.ValueBool(),
		EnableChat:                      plan.EnableChat.ValueString(),
		EnableOverlayText:               plan.EnableOverlayText.ValueBool(),
		ForcePresenterIntoMain:          plan.ForcePresenterIntoMain.ValueBool(),
		GuestPIN:                        plan.GuestPIN.ValueString(),
		GuestsCanPresent:                plan.GuestsCanPresent.ValueBool(),
		GuestsCanSeeGuests:              plan.GuestsCanSeeGuests.ValueString(),
		LiveCaptionsEnabled:             plan.LiveCaptionsEnabled.ValueString(),
		MatchString:                     plan.MatchString.ValueString(),
		MuteAllGuests:                   plan.MuteAllGuests.ValueBool(),
		NonIdpParticipants:              plan.NonIdpParticipants.ValueString(),
		PIN:                             plan.PIN.ValueString(),
		PostMatchString:                 plan.PostMatchString.ValueString(),
		PostReplaceString:               plan.PostReplaceString.ValueString(),
		PrimaryOwnerEmailAddress:        plan.PrimaryOwnerEmailAddress.ValueString(),
		ReplaceString:                   plan.ReplaceString.ValueString(),
		ServiceType:                     plan.ServiceType.ValueString(),
		SoftmuteEnabled:                 plan.SoftmuteEnabled.ValueBool(),
		SyncTag:                         plan.SyncTag.ValueString(),
		Tag:                             plan.Tag.ValueString(),
		TwoStageDialType:                plan.TwoStageDialType.ValueString(),
	}

	// All nullable fields
	// Only set optional fields if they are not null in the plan
	if !plan.Aliases.IsNull() && !plan.Aliases.IsUnknown() {
		aliases, diags := getStringList(ctx, plan.Aliases)
		resp.Diagnostics.Append(diags...)
		createRequest.Aliases = &aliases
	}
	if !plan.AutomaticParticipants.IsNull() && !plan.AutomaticParticipants.IsUnknown() {
		automaticParticipants, diags := getStringList(ctx, plan.AutomaticParticipants)
		resp.Diagnostics.Append(diags...)
		createRequest.AutomaticParticipants = automaticParticipants
	}
	if !plan.CryptoMode.IsNull() && !plan.CryptoMode.IsUnknown() {
		cryptoMode := plan.CryptoMode.ValueString()
		createRequest.CryptoMode = &cryptoMode
	}
	if !plan.GMSAccessToken.IsNull() && !plan.GMSAccessToken.IsUnknown() {
		gmsAccessToken := plan.GMSAccessToken.ValueString()
		createRequest.GMSAccessToken = &gmsAccessToken
	}
	if !plan.GuestIdentityProviderGroup.IsNull() && !plan.GuestIdentityProviderGroup.IsUnknown() {
		guestIdentityProviderGroup := plan.GuestIdentityProviderGroup.ValueString()
		createRequest.GuestIdentityProviderGroup = &guestIdentityProviderGroup
	}
	if !plan.HostIdentityProviderGroup.IsNull() && !plan.HostIdentityProviderGroup.IsUnknown() {
		hostIdentityProviderGroup := plan.HostIdentityProviderGroup.ValueString()
		createRequest.HostIdentityProviderGroup = &hostIdentityProviderGroup
	}
	if !plan.IVRTheme.IsNull() && !plan.IVRTheme.IsUnknown() {
		ivrTheme := plan.IVRTheme.ValueString()
		createRequest.IVRTheme = &ivrTheme
	}
	if !plan.MaxCallRateIn.IsNull() && !plan.MaxCallRateIn.IsUnknown() {
		maxCallrateIn := int(plan.MaxCallRateIn.ValueInt32())
		createRequest.MaxCallRateIn = &maxCallrateIn
	}
	if !plan.MaxCallRateOut.IsNull() && !plan.MaxCallRateOut.IsUnknown() {
		maxCallrateOut := int(plan.MaxCallRateOut.ValueInt32())
		createRequest.MaxCallRateOut = &maxCallrateOut
	}
	if !plan.MaxPixelsPerSecond.IsNull() && !plan.MaxPixelsPerSecond.IsUnknown() {
		maxPixelsPerSecond := plan.MaxPixelsPerSecond.ValueString()
		createRequest.MaxPixelsPerSecond = &maxPixelsPerSecond
	}
	if !plan.MediaPlaylist.IsNull() && !plan.MediaPlaylist.IsUnknown() {
		mediaPlaylist := plan.MediaPlaylist.ValueString()
		createRequest.MediaPlaylist = &mediaPlaylist
	}
	if !plan.MSSIPProxy.IsNull() && !plan.MSSIPProxy.IsUnknown() {
		mssipProxy := plan.MSSIPProxy.ValueString()
		createRequest.MSSIPProxy = &mssipProxy
	}
	if !plan.OnCompletion.IsNull() && !plan.OnCompletion.IsUnknown() {
		onCompletion := plan.OnCompletion.ValueString()
		createRequest.OnCompletion = &onCompletion
	}
	if !plan.ParticipantLimit.IsNull() && !plan.ParticipantLimit.IsUnknown() {
		participantLimit := int(plan.ParticipantLimit.ValueInt32())
		createRequest.ParticipantLimit = &participantLimit
	}
	if !plan.PinningConfig.IsNull() && !plan.PinningConfig.IsUnknown() {
		pinningConfig := plan.PinningConfig.ValueString()
		createRequest.PinningConfig = &pinningConfig
	}
	if !plan.SystemLocation.IsNull() && !plan.SystemLocation.IsUnknown() {
		systemLocation := plan.SystemLocation.ValueString()
		createRequest.SystemLocation = &systemLocation
	}
	if !plan.TeamsProxy.IsNull() && !plan.TeamsProxy.IsUnknown() {
		teamsProxy := plan.TeamsProxy.ValueString()
		createRequest.TeamsProxy = &teamsProxy
	}
	if !plan.GuestIdentityProviderGroup.IsNull() && !plan.GuestIdentityProviderGroup.IsUnknown() {
		guestIdentityProviderGroup := plan.GuestIdentityProviderGroup.ValueString()
		createRequest.GuestIdentityProviderGroup = &guestIdentityProviderGroup
	}
	if !plan.GuestView.IsNull() && !plan.GuestView.IsUnknown() {
		guestView := plan.GuestView.ValueString()
		createRequest.GuestView = &guestView
	}
	if !plan.HostView.IsNull() && !plan.HostView.IsUnknown() {
		hostView := plan.HostView.ValueString()
		createRequest.HostView = &hostView
	}
	if !plan.HostIdentityProviderGroup.IsNull() && !plan.HostIdentityProviderGroup.IsUnknown() {
		hostIdentityProviderGroup := plan.HostIdentityProviderGroup.ValueString()
		createRequest.HostIdentityProviderGroup = &hostIdentityProviderGroup
	}
	if plan.HostView.IsNull() && plan.HostView.IsUnknown() {
		hostView := plan.GuestView.ValueString()
		createRequest.HostView = &hostView
	}
	if plan.IVRTheme.IsNull() && plan.IVRTheme.IsUnknown() {
		ivrTheme := plan.IVRTheme.ValueString()
		createRequest.IVRTheme = &ivrTheme
	}
	if plan.MaxCallRateIn.IsNull() && plan.MaxCallRateIn.IsUnknown() {
		maxCallRateIn := int(plan.MaxCallRateIn.ValueInt32())
		createRequest.MaxCallRateIn = &maxCallRateIn
	}
	if plan.MaxCallRateOut.IsNull() && plan.MaxCallRateOut.IsUnknown() {
		maxCallRateOut := int(plan.MaxCallRateOut.ValueInt32())
		createRequest.MaxCallRateOut = &maxCallRateOut
	}
	if plan.MaxPixelsPerSecond.IsNull() && plan.MaxPixelsPerSecond.IsUnknown() {
		maxPixelsPerSecond := plan.MaxPixelsPerSecond.ValueString()
		createRequest.MaxPixelsPerSecond = &maxPixelsPerSecond
	}
	if plan.MediaPlaylist.IsNull() && plan.MediaPlaylist.IsUnknown() {
		mediaPlaylist := plan.MediaPlaylist.ValueString()
		createRequest.MediaPlaylist = &mediaPlaylist
	}
	if plan.MSSIPProxy.IsNull() && plan.MSSIPProxy.IsUnknown() {
		mssipProxy := plan.MSSIPProxy.ValueString()
		createRequest.MSSIPProxy = &mssipProxy
	}
	if plan.OnCompletion.IsNull() && plan.OnCompletion.IsUnknown() {
		onCompletion := plan.OnCompletion.ValueString()
		createRequest.OnCompletion = &onCompletion
	}
	if plan.ParticipantLimit.IsNull() && plan.ParticipantLimit.IsUnknown() {
		participantLimit := int(plan.ParticipantLimit.ValueInt32())
		createRequest.ParticipantLimit = &participantLimit
	}
	if plan.PinningConfig.IsNull() && plan.PinningConfig.IsUnknown() {
		pinningConfig := plan.PinningConfig.ValueString()
		createRequest.PinningConfig = &pinningConfig
	}
	if plan.SystemLocation.IsNull() && plan.SystemLocation.IsUnknown() {
		systemLocation := plan.SystemLocation.ValueString()
		createRequest.SystemLocation = &systemLocation
	}
	if plan.TeamsProxy.IsNull() && plan.TeamsProxy.IsUnknown() {
		teamsProxy := plan.TeamsProxy.ValueString()
		createRequest.TeamsProxy = &teamsProxy
	}

	createResponse, err := r.InfinityClient.Config().CreateConference(ctx, createRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Infinity conference",
			fmt.Sprintf("Could not create Infinity conference: %s", err),
		)
		return
	}

	resourceID, err := createResponse.ResourceID()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Retrieving Infinity conference ID",
			fmt.Sprintf("Could not retrieve ID for created Infinity conference: %s", err),
		)
		return
	}

	// Read the state from the API to get all computed values
	model, err := r.read(ctx, resourceID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Created Infinity conference",
			fmt.Sprintf("Could not read created Infinity conference with ID %d: %s", resourceID, err),
		)
		return
	}
	tflog.Trace(ctx, fmt.Sprintf("created Infinity conference with ID: %s, name: %s", model.ID, model.Name))

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *InfinityConferenceResource) read(ctx context.Context, resourceID int) (*InfinityConferenceResourceModel, error) {
	var data InfinityConferenceResourceModel

	srv, err := r.InfinityClient.Config().GetConference(ctx, resourceID)
	if err != nil {
		return nil, err
	}

	if srv.ResourceURI == "" {
		return nil, fmt.Errorf("conference with ID %d not found", resourceID)
	}

	data.ID = types.StringValue(srv.ResourceURI)
	data.ResourceID = types.Int32Value(int32(resourceID))
	data.Name = types.StringValue(srv.Name)
	data.AllowGuests = types.BoolValue(srv.AllowGuests)
	data.BreakoutRooms = types.BoolValue(srv.BreakoutRooms)
	data.CallType = types.StringValue(srv.CallType)
	data.CryptoMode = types.StringPointerValue(srv.CryptoMode)
	data.DenoiseEnabled = types.BoolValue(srv.DenoiseEnabled)
	data.Description = types.StringValue(srv.Description)
	data.DirectMedia = types.StringValue(srv.DirectMedia)
	data.DirectMediaNotificationDuration = types.Int32Value(int32(srv.DirectMediaNotificationDuration))
	data.EnableActiveSpeakerIndication = types.BoolValue(srv.EnableActiveSpeakerIndication)
	data.EnableChat = types.StringValue(srv.EnableChat)
	data.EnableOverlayText = types.BoolValue(srv.EnableOverlayText)
	data.ForcePresenterIntoMain = types.BoolValue(srv.ForcePresenterIntoMain)
	data.GMSAccessToken = types.StringPointerValue(srv.GMSAccessToken)
	data.GuestIdentityProviderGroup = types.StringPointerValue(srv.GuestIdentityProviderGroup)
	data.GuestPIN = types.StringValue(srv.GuestPIN)
	data.GuestView = types.StringPointerValue(srv.GuestView)
	data.GuestsCanPresent = types.BoolValue(srv.GuestsCanPresent)
	data.GuestsCanSeeGuests = types.StringValue(srv.GuestsCanSeeGuests)
	data.HostIdentityProviderGroup = types.StringPointerValue(srv.HostIdentityProviderGroup)
	data.HostView = types.StringPointerValue(srv.HostView)
	data.IVRTheme = types.StringPointerValue(srv.IVRTheme)
	data.LiveCaptionsEnabled = types.StringValue(srv.LiveCaptionsEnabled)
	data.MatchString = types.StringValue(srv.MatchString)
	data.MaxPixelsPerSecond = types.StringPointerValue(srv.MaxPixelsPerSecond)
	data.MediaPlaylist = types.StringPointerValue(srv.MediaPlaylist)
	data.MSSIPProxy = types.StringPointerValue(srv.MSSIPProxy)
	data.MuteAllGuests = types.BoolValue(srv.MuteAllGuests)
	data.NonIdpParticipants = types.StringValue(srv.NonIdpParticipants)
	data.OnCompletion = types.StringPointerValue(srv.OnCompletion)
	data.PIN = types.StringValue(srv.PIN)
	data.PinningConfig = types.StringPointerValue(srv.PinningConfig)
	data.PostMatchString = types.StringValue(srv.PostMatchString)
	data.PostReplaceString = types.StringValue(srv.PostReplaceString)
	data.PrimaryOwnerEmailAddress = types.StringValue(srv.PrimaryOwnerEmailAddress)
	data.ReplaceString = types.StringValue(srv.ReplaceString)
	data.ServiceType = types.StringValue(srv.ServiceType)
	data.SoftmuteEnabled = types.BoolValue(srv.SoftmuteEnabled)
	data.SyncTag = types.StringValue(srv.SyncTag)
	data.SystemLocation = types.StringPointerValue(srv.SystemLocation)
	data.Tag = types.StringValue(srv.Tag)
	data.TeamsProxy = types.StringPointerValue(srv.TeamsProxy)
	data.TwoStageDialType = types.StringValue(srv.TwoStageDialType)

	// Handle nullable integer fields
	if srv.MaxCallRateIn != nil {
		data.MaxCallRateIn = types.Int32Value(int32(*srv.MaxCallRateIn)) // #nosec G115 -- API values are expected to be within int32 range
	}
	if srv.MaxCallRateOut != nil {
		data.MaxCallRateOut = types.Int32Value(int32(*srv.MaxCallRateOut)) // #nosec G115 -- API values are expected to be within int32 range
	}
	if srv.ParticipantLimit != nil {
		data.ParticipantLimit = types.Int32Value(int32(*srv.ParticipantLimit)) // #nosec G115 -- API values are expected to be within int32 range
	}

	// Convert conference aliases from SDK to Terraform format
	var aliases []string
	if srv.Aliases != nil {
		for _, alias := range *srv.Aliases {
			aliases = append(aliases, fmt.Sprintf("/api/admin/configuration/v1/conference_alias/%d/", alias.ID))
		}
	}
	aliasesSetValue, diags := types.SetValueFrom(ctx, types.StringType, aliases)
	if diags.HasError() {
		return nil, fmt.Errorf("error converting aliases: %v", diags)
	}
	data.Aliases = aliasesSetValue

	// Convert automatic participants from SDK to Terraform format
	var participants []string
	if srv.AutomaticParticipants != nil {
		participants = append(participants, srv.AutomaticParticipants...)
	}
	participantsSetValue, diags := types.SetValueFrom(ctx, types.StringType, participants)
	if diags.HasError() {
		return nil, fmt.Errorf("error converting aliases: %v", diags)
	}
	data.AutomaticParticipants = participantsSetValue

	return &data, nil
}

func (r *InfinityConferenceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	state := &InfinityConferenceResourceModel{}

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
			"Error Reading Infinity conference",
			fmt.Sprintf("Could not read Infinity conference: %s", err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *InfinityConferenceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	plan := &InfinityConferenceResourceModel{}
	state := &InfinityConferenceResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resourceID := int(state.ResourceID.ValueInt32())

	// initialize create request with required fields, list fields, and fields with defaults
	updateRequest := &config.ConferenceUpdateRequest{
		Name:                            plan.Name.ValueString(),
		AllowGuests:                     plan.AllowGuests.ValueBool(),
		BreakoutRooms:                   plan.BreakoutRooms.ValueBool(),
		CallType:                        plan.CallType.ValueString(),
		DenoiseEnabled:                  plan.DenoiseEnabled.ValueBool(),
		Description:                     plan.Description.ValueString(),
		DirectMedia:                     plan.DirectMedia.ValueString(),
		DirectMediaNotificationDuration: int(plan.DirectMediaNotificationDuration.ValueInt32()),
		EnableActiveSpeakerIndication:   plan.EnableActiveSpeakerIndication.ValueBool(),
		EnableChat:                      plan.EnableChat.ValueString(),
		EnableOverlayText:               plan.EnableOverlayText.ValueBool(),
		ForcePresenterIntoMain:          plan.ForcePresenterIntoMain.ValueBool(),
		GuestPIN:                        plan.GuestPIN.ValueString(),
		GuestsCanPresent:                plan.GuestsCanPresent.ValueBool(),
		GuestsCanSeeGuests:              plan.GuestsCanSeeGuests.ValueString(),
		LiveCaptionsEnabled:             plan.LiveCaptionsEnabled.ValueString(),
		MatchString:                     plan.MatchString.ValueString(),
		MuteAllGuests:                   plan.MuteAllGuests.ValueBool(),
		NonIdpParticipants:              plan.NonIdpParticipants.ValueString(),
		PIN:                             plan.PIN.ValueString(),
		PostMatchString:                 plan.PostMatchString.ValueString(),
		PostReplaceString:               plan.PostReplaceString.ValueString(),
		PrimaryOwnerEmailAddress:        plan.PrimaryOwnerEmailAddress.ValueString(),
		ReplaceString:                   plan.ReplaceString.ValueString(),
		ServiceType:                     plan.ServiceType.ValueString(),
		SoftmuteEnabled:                 plan.SoftmuteEnabled.ValueBool(),
		SyncTag:                         plan.SyncTag.ValueString(),
		Tag:                             plan.Tag.ValueString(),
		TwoStageDialType:                plan.TwoStageDialType.ValueString(),
	}

	// All nullable fields
	// Only set optional fields if they are not null in the plan
	if !plan.Aliases.IsNull() && !plan.Aliases.IsUnknown() {
		aliases, diags := getStringList(ctx, plan.Aliases)
		resp.Diagnostics.Append(diags...)
		updateRequest.Aliases = &aliases
	}
	if !plan.AutomaticParticipants.IsNull() && !plan.AutomaticParticipants.IsUnknown() {
		automaticParticipants, diags := getStringList(ctx, plan.AutomaticParticipants)
		resp.Diagnostics.Append(diags...)
		updateRequest.AutomaticParticipants = automaticParticipants
	}
	if !plan.CryptoMode.IsNull() && !plan.CryptoMode.IsUnknown() {
		cryptoMode := plan.CryptoMode.ValueString()
		updateRequest.CryptoMode = &cryptoMode
	}
	if !plan.GMSAccessToken.IsNull() && !plan.GMSAccessToken.IsUnknown() {
		gmsAccessToken := plan.GMSAccessToken.ValueString()
		updateRequest.GMSAccessToken = &gmsAccessToken
	}
	if !plan.GuestIdentityProviderGroup.IsNull() && !plan.GuestIdentityProviderGroup.IsUnknown() {
		guestIdentityProviderGroup := plan.GuestIdentityProviderGroup.ValueString()
		updateRequest.GuestIdentityProviderGroup = &guestIdentityProviderGroup
	}
	if !plan.HostIdentityProviderGroup.IsNull() && !plan.HostIdentityProviderGroup.IsUnknown() {
		hostIdentityProviderGroup := plan.HostIdentityProviderGroup.ValueString()
		updateRequest.HostIdentityProviderGroup = &hostIdentityProviderGroup
	}
	if !plan.IVRTheme.IsNull() && !plan.IVRTheme.IsUnknown() {
		ivrTheme := plan.IVRTheme.ValueString()
		updateRequest.IVRTheme = &ivrTheme
	}
	if !plan.MaxCallRateIn.IsNull() && !plan.MaxCallRateIn.IsUnknown() {
		maxCallrateIn := int(plan.MaxCallRateIn.ValueInt32())
		updateRequest.MaxCallRateIn = &maxCallrateIn
	}
	if !plan.MaxCallRateOut.IsNull() && !plan.MaxCallRateOut.IsUnknown() {
		maxCallrateOut := int(plan.MaxCallRateOut.ValueInt32())
		updateRequest.MaxCallRateOut = &maxCallrateOut
	}
	if !plan.MaxPixelsPerSecond.IsNull() && !plan.MaxPixelsPerSecond.IsUnknown() {
		maxPixelsPerSecond := plan.MaxPixelsPerSecond.ValueString()
		updateRequest.MaxPixelsPerSecond = &maxPixelsPerSecond
	}
	if !plan.MediaPlaylist.IsNull() && !plan.MediaPlaylist.IsUnknown() {
		mediaPlaylist := plan.MediaPlaylist.ValueString()
		updateRequest.MediaPlaylist = &mediaPlaylist
	}
	if !plan.MSSIPProxy.IsNull() && !plan.MSSIPProxy.IsUnknown() {
		mssipProxy := plan.MSSIPProxy.ValueString()
		updateRequest.MSSIPProxy = &mssipProxy
	}
	if !plan.OnCompletion.IsNull() && !plan.OnCompletion.IsUnknown() {
		onCompletion := plan.OnCompletion.ValueString()
		updateRequest.OnCompletion = &onCompletion
	}
	if !plan.ParticipantLimit.IsNull() && !plan.ParticipantLimit.IsUnknown() {
		participantLimit := int(plan.ParticipantLimit.ValueInt32())
		updateRequest.ParticipantLimit = &participantLimit
	}
	if !plan.PinningConfig.IsNull() && !plan.PinningConfig.IsUnknown() {
		pinningConfig := plan.PinningConfig.ValueString()
		updateRequest.PinningConfig = &pinningConfig
	}
	if !plan.SystemLocation.IsNull() && !plan.SystemLocation.IsUnknown() {
		systemLocation := plan.SystemLocation.ValueString()
		updateRequest.SystemLocation = &systemLocation
	}
	if !plan.TeamsProxy.IsNull() && !plan.TeamsProxy.IsUnknown() {
		teamsProxy := plan.TeamsProxy.ValueString()
		updateRequest.TeamsProxy = &teamsProxy
	}
	if !plan.GuestIdentityProviderGroup.IsNull() && !plan.GuestIdentityProviderGroup.IsUnknown() {
		guestIdentityProviderGroup := plan.GuestIdentityProviderGroup.ValueString()
		updateRequest.GuestIdentityProviderGroup = &guestIdentityProviderGroup
	}
	if !plan.GuestView.IsNull() && !plan.GuestView.IsUnknown() {
		guestView := plan.GuestView.ValueString()
		updateRequest.GuestView = &guestView
	}
	if !plan.HostView.IsNull() && !plan.HostView.IsUnknown() {
		hostView := plan.HostView.ValueString()
		updateRequest.HostView = &hostView
	}
	if !plan.HostIdentityProviderGroup.IsNull() && !plan.HostIdentityProviderGroup.IsUnknown() {
		hostIdentityProviderGroup := plan.HostIdentityProviderGroup.ValueString()
		updateRequest.HostIdentityProviderGroup = &hostIdentityProviderGroup
	}
	if plan.HostView.IsNull() && plan.HostView.IsUnknown() {
		hostView := plan.GuestView.ValueString()
		updateRequest.HostView = &hostView
	}
	if plan.IVRTheme.IsNull() && plan.IVRTheme.IsUnknown() {
		ivrTheme := plan.IVRTheme.ValueString()
		updateRequest.IVRTheme = &ivrTheme
	}
	if plan.MaxCallRateIn.IsNull() && plan.MaxCallRateIn.IsUnknown() {
		maxCallRateIn := int(plan.MaxCallRateIn.ValueInt32())
		updateRequest.MaxCallRateIn = &maxCallRateIn
	}
	if plan.MaxCallRateOut.IsNull() && plan.MaxCallRateOut.IsUnknown() {
		maxCallRateOut := int(plan.MaxCallRateOut.ValueInt32())
		updateRequest.MaxCallRateOut = &maxCallRateOut
	}
	if plan.MaxPixelsPerSecond.IsNull() && plan.MaxPixelsPerSecond.IsUnknown() {
		maxPixelsPerSecond := plan.MaxPixelsPerSecond.ValueString()
		updateRequest.MaxPixelsPerSecond = &maxPixelsPerSecond
	}
	if plan.MediaPlaylist.IsNull() && plan.MediaPlaylist.IsUnknown() {
		mediaPlaylist := plan.MediaPlaylist.ValueString()
		updateRequest.MediaPlaylist = &mediaPlaylist
	}
	if plan.MSSIPProxy.IsNull() && plan.MSSIPProxy.IsUnknown() {
		mssipProxy := plan.MSSIPProxy.ValueString()
		updateRequest.MSSIPProxy = &mssipProxy
	}
	if plan.OnCompletion.IsNull() && plan.OnCompletion.IsUnknown() {
		onCompletion := plan.OnCompletion.ValueString()
		updateRequest.OnCompletion = &onCompletion
	}
	if plan.ParticipantLimit.IsNull() && plan.ParticipantLimit.IsUnknown() {
		participantLimit := int(plan.ParticipantLimit.ValueInt32())
		updateRequest.ParticipantLimit = &participantLimit
	}
	if plan.PinningConfig.IsNull() && plan.PinningConfig.IsUnknown() {
		pinningConfig := plan.PinningConfig.ValueString()
		updateRequest.PinningConfig = &pinningConfig
	}
	if plan.SystemLocation.IsNull() && plan.SystemLocation.IsUnknown() {
		systemLocation := plan.SystemLocation.ValueString()
		updateRequest.SystemLocation = &systemLocation
	}
	if plan.TeamsProxy.IsNull() && plan.TeamsProxy.IsUnknown() {
		teamsProxy := plan.TeamsProxy.ValueString()
		updateRequest.TeamsProxy = &teamsProxy
	}

	_, err := r.InfinityClient.Config().UpdateConference(ctx, resourceID, updateRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Infinity conference",
			fmt.Sprintf("Could not update Infinity conference with ID %d: %s", resourceID, err),
		)
		return
	}

	// Re-read the resource to get the latest state
	updatedModel, err := r.read(ctx, resourceID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Updated Infinity conference",
			fmt.Sprintf("Could not read updated Infinity conference with ID %d: %s", resourceID, err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, updatedModel)...)
}

func (r *InfinityConferenceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	state := &InfinityConferenceResourceModel{}

	tflog.Info(ctx, "Deleting Infinity conference")

	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.InfinityClient.Config().DeleteConference(ctx, int(state.ResourceID.ValueInt32()))

	// Ignore 404 Not Found and Lookup errors on delete
	if err != nil && !isNotFoundError(err) && !isLookupError(err) {
		resp.Diagnostics.AddError(
			"Error Deleting Infinity conference",
			fmt.Sprintf("Could not delete Infinity conference with ID %s: %s", state.ID.ValueString(), err),
		)
		return
	}
}

func (r *InfinityConferenceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resourceID, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Import ID must be a valid integer for the resource ID. Got: %s", req.ID),
		)
		return
	}

	tflog.Trace(ctx, fmt.Sprintf("Importing Infinity conference with resource ID: %d", resourceID))

	// Read the resource from the API
	model, err := r.read(ctx, resourceID)
	if err != nil {
		// Check if the error is a 404 (not found)
		if isNotFoundError(err) {
			resp.Diagnostics.AddError(
				"Infinity Conference Not Found",
				fmt.Sprintf("Infinity conference with resource ID %d not found.", resourceID),
			)
			return
		}
		resp.Diagnostics.AddError(
			"Error Importing Infinity Conference",
			fmt.Sprintf("Could not import Infinity conference with resource ID %d: %s", resourceID, err),
		)
		return
	}

	// Set the state from the imported resource
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}
