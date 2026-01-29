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

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/pexip/go-infinity-sdk/v38/config"

	"github.com/pexip/terraform-provider-pexip/internal/provider/validators"
)

var (
	_ resource.ResourceWithImportState = (*InfinityMsExchangeConnectorResource)(nil)
)

type InfinityMsExchangeConnectorResource struct {
	InfinityClient InfinityClient
}

type InfinityMsExchangeConnectorResourceModel struct {
	ID                                               types.String `tfsdk:"id"`
	ResourceID                                       types.Int32  `tfsdk:"resource_id"`
	Name                                             types.String `tfsdk:"name"`
	Description                                      types.String `tfsdk:"description"`
	RoomMailboxEmailAddress                          types.String `tfsdk:"room_mailbox_email_address"`
	RoomMailboxName                                  types.String `tfsdk:"room_mailbox_name"`
	URL                                              types.String `tfsdk:"url"`
	Username                                         types.String `tfsdk:"username"`
	Password                                         types.String `tfsdk:"password"`
	AuthenticationMethod                             types.String `tfsdk:"authentication_method"`
	AuthProvider                                     types.String `tfsdk:"auth_provider"`
	UUID                                             types.String `tfsdk:"uuid"`
	ScheduledAliasPrefix                             types.String `tfsdk:"scheduled_alias_prefix"`
	ScheduledAliasDomain                             types.String `tfsdk:"scheduled_alias_domain"`
	ScheduledAliasSuffixLength                       types.Int64  `tfsdk:"scheduled_alias_suffix_length"`
	MeetingBufferBefore                              types.Int64  `tfsdk:"meeting_buffer_before"`
	MeetingBufferAfter                               types.Int64  `tfsdk:"meeting_buffer_after"`
	EnableDynamicVmrs                                types.Bool   `tfsdk:"enable_dynamic_vmrs"`
	EnablePersonalVmrs                               types.Bool   `tfsdk:"enable_personal_vmrs"`
	AllowNewUsers                                    types.Bool   `tfsdk:"allow_new_users"`
	DisableProxy                                     types.Bool   `tfsdk:"disable_proxy"`
	UseCustomAddInSources                            types.Bool   `tfsdk:"use_custom_add_in_sources"`
	EnableAddinDebugLogs                             types.Bool   `tfsdk:"enable_addin_debug_logs"`
	OauthClientID                                    types.String `tfsdk:"oauth_client_id"`
	OauthClientSecret                                types.String `tfsdk:"oauth_client_secret"`
	OauthAuthEndpoint                                types.String `tfsdk:"oauth_auth_endpoint"`
	OauthTokenEndpoint                               types.String `tfsdk:"oauth_token_endpoint"`
	OauthRedirectURI                                 types.String `tfsdk:"oauth_redirect_uri"`
	OauthRefreshToken                                types.String `tfsdk:"oauth_refresh_token"`
	OauthState                                       types.String `tfsdk:"oauth_state"`
	KerberosRealm                                    types.String `tfsdk:"kerberos_realm"`
	KerberosKdc                                      types.String `tfsdk:"kerberos_kdc"`
	KerberosKdcHttpsProxy                            types.String `tfsdk:"kerberos_kdc_https_proxy"`
	KerberosExchangeSpn                              types.String `tfsdk:"kerberos_exchange_spn"`
	KerberosEnableTls                                types.Bool   `tfsdk:"kerberos_enable_tls"`
	KerberosAuthEveryRequest                         types.Bool   `tfsdk:"kerberos_auth_every_request"`
	KerberosVerifyTlsUsingCustomCa                   types.Bool   `tfsdk:"kerberos_verify_tls_using_custom_ca"`
	AddinServerDomain                                types.String `tfsdk:"addin_server_domain"`
	AddinDisplayName                                 types.String `tfsdk:"addin_display_name"`
	AddinDescription                                 types.String `tfsdk:"addin_description"`
	AddinProviderName                                types.String `tfsdk:"addin_provider_name"`
	AddinButtonLabel                                 types.String `tfsdk:"addin_button_label"`
	AddinGroupLabel                                  types.String `tfsdk:"addin_group_label"`
	AddinSupertipTitle                               types.String `tfsdk:"addin_supertip_title"`
	AddinSupertipDescription                         types.String `tfsdk:"addin_supertip_description"`
	AddinApplicationID                               types.String `tfsdk:"addin_application_id"`
	AddinAuthorityURL                                types.String `tfsdk:"addin_authority_url"`
	AddinOidcMetadataURL                             types.String `tfsdk:"addin_oidc_metadata_url"`
	AddinAuthenticationMethod                        types.String `tfsdk:"addin_authentication_method"`
	AddinNaaWebApiApplicationID                      types.String `tfsdk:"addin_naa_web_api_application_id"`
	PersonalVmrOauthClientID                         types.String `tfsdk:"personal_vmr_oauth_client_id"`
	PersonalVmrOauthClientSecret                     types.String `tfsdk:"personal_vmr_oauth_client_secret"`
	PersonalVmrOauthAuthEndpoint                     types.String `tfsdk:"personal_vmr_oauth_auth_endpoint"`
	PersonalVmrOauthTokenEndpoint                    types.String `tfsdk:"personal_vmr_oauth_token_endpoint"`
	PersonalVmrAdfsRelyingPartyTrustIdentifier       types.String `tfsdk:"personal_vmr_adfs_relying_party_trust_identifier"`
	OfficeJsURL                                      types.String `tfsdk:"office_js_url"`
	MicrosoftFabricURL                               types.String `tfsdk:"microsoft_fabric_url"`
	MicrosoftFabricComponentsURL                     types.String `tfsdk:"microsoft_fabric_components_url"`
	AdditionalAddInScriptSources                     types.String `tfsdk:"additional_add_in_script_sources"`
	Domains                                          types.Set    `tfsdk:"domains"`
	HostIdentityProviderGroup                        types.String `tfsdk:"host_identity_provider_group"`
	IvrTheme                                         types.String `tfsdk:"ivr_theme"`
	NonIdpParticipants                               types.String `tfsdk:"non_idp_participants"`
	PrivateKey                                       types.String `tfsdk:"private_key"`
	PublicKey                                        types.String `tfsdk:"public_key"`
	AcceptEditedOccurrenceTemplate                   types.String `tfsdk:"accept_edited_occurrence_template"`
	AcceptEditedRecurringSeriesTemplate              types.String `tfsdk:"accept_edited_recurring_series_template"`
	AcceptEditedSingleMeetingTemplate                types.String `tfsdk:"accept_edited_single_meeting_template"`
	AcceptNewRecurringSeriesTemplate                 types.String `tfsdk:"accept_new_recurring_series_template"`
	AcceptNewSingleMeetingTemplate                   types.String `tfsdk:"accept_new_single_meeting_template"`
	ConferenceDescriptionTemplate                    types.String `tfsdk:"conference_description_template"`
	ConferenceNameTemplate                           types.String `tfsdk:"conference_name_template"`
	ConferenceSubjectTemplate                        types.String `tfsdk:"conference_subject_template"`
	MeetingInstructionsTemplate                      types.String `tfsdk:"meeting_instructions_template"`
	PersonalVmrDescriptionTemplate                   types.String `tfsdk:"personal_vmr_description_template"`
	PersonalVmrInstructionsTemplate                  types.String `tfsdk:"personal_vmr_instructions_template"`
	PersonalVmrLocationTemplate                      types.String `tfsdk:"personal_vmr_location_template"`
	PersonalVmrNameTemplate                          types.String `tfsdk:"personal_vmr_name_template"`
	PlaceholderInstructionsTemplate                  types.String `tfsdk:"placeholder_instructions_template"`
	RejectAliasConflictTemplate                      types.String `tfsdk:"reject_alias_conflict_template"`
	RejectAliasDeletedTemplate                       types.String `tfsdk:"reject_alias_deleted_template"`
	RejectGeneralErrorTemplate                       types.String `tfsdk:"reject_general_error_template"`
	RejectInvalidAliasIDTemplate                     types.String `tfsdk:"reject_invalid_alias_id_template"`
	RejectRecurringSeriesPastTemplate                types.String `tfsdk:"reject_recurring_series_past_template"`
	RejectSingleMeetingPast                          types.String `tfsdk:"reject_single_meeting_past"`
	ScheduledAliasDescriptionTemplate                types.String `tfsdk:"scheduled_alias_description_template"`
	AddinPaneAlreadyVideoMeetingHeading              types.String `tfsdk:"addin_pane_already_video_meeting_heading"`
	AddinPaneAlreadyVideoMeetingMessage              types.String `tfsdk:"addin_pane_already_video_meeting_message"`
	AddinPaneButtonTitle                             types.String `tfsdk:"addin_pane_button_title"`
	AddinPaneDescription                             types.String `tfsdk:"addin_pane_description"`
	AddinPaneGeneralErrorHeading                     types.String `tfsdk:"addin_pane_general_error_heading"`
	AddinPaneGeneralErrorMessage                     types.String `tfsdk:"addin_pane_general_error_message"`
	AddinPaneManagementNodeDownHeading               types.String `tfsdk:"addin_pane_management_node_down_heading"`
	AddinPaneManagementNodeDownMessage               types.String `tfsdk:"addin_pane_management_node_down_message"`
	AddinPanePersonalVmrAddButton                    types.String `tfsdk:"addin_pane_personal_vmr_add_button"`
	AddinPanePersonalVmrErrorGettingMessage          types.String `tfsdk:"addin_pane_personal_vmr_error_getting_message"`
	AddinPanePersonalVmrErrorInsertingMeetingMessage types.String `tfsdk:"addin_pane_personal_vmr_error_inserting_meeting_message"`
	AddinPanePersonalVmrErrorSigningInMessage        types.String `tfsdk:"addin_pane_personal_vmr_error_signing_in_message"`
	AddinPanePersonalVmrNoneMessage                  types.String `tfsdk:"addin_pane_personal_vmr_none_message"`
	AddinPanePersonalVmrSelectMessage                types.String `tfsdk:"addin_pane_personal_vmr_select_message"`
	AddinPanePersonalVmrSignInButton                 types.String `tfsdk:"addin_pane_personal_vmr_sign_in_button"`
	AddinPaneSuccessHeading                          types.String `tfsdk:"addin_pane_success_heading"`
	AddinPaneSuccessMessage                          types.String `tfsdk:"addin_pane_success_message"`
	AddinPaneTitle                                   types.String `tfsdk:"addin_pane_title"`
}

func (r *InfinityMsExchangeConnectorResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_infinity_ms_exchange_connector"
}

func (r *InfinityMsExchangeConnectorResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *InfinityMsExchangeConnectorResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Resource URI for the Microsoft Exchange connector in Infinity",
			},
			"resource_id": schema.Int32Attribute{
				Computed:            true,
				MarkdownDescription: "The resource integer identifier for the Microsoft Exchange connector in Infinity",
			},
			"name": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					stringvalidator.LengthAtMost(250),
				},
				MarkdownDescription: "The name of the Microsoft Exchange connector. Maximum length: 250 characters.",
			},
			"description": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(""),
				Validators: []validator.String{
					stringvalidator.LengthAtMost(250),
				},
				MarkdownDescription: "An optional description of the Secure Scheduler for Exchange Integration. Maximum length: 250 characters.",
			},
			"room_mailbox_email_address": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Room mailbox email address for Exchange integration.",
			},
			"room_mailbox_name": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Room mailbox name for Exchange integration.",
			},
			"url": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					validators.URL(false),
				},
				MarkdownDescription: "Exchange server URL for connectivity.",
			},
			"username": schema.StringAttribute{
				Computed:            true,
				Optional:            true,
				Default:             stringdefault.StaticString(""),
				MarkdownDescription: "Username for Exchange authentication.",
			},
			"password": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				MarkdownDescription: "Password for Exchange authentication. This field is sensitive.",
			},
			"authentication_method": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("BASIC"),
				Validators: []validator.String{
					stringvalidator.OneOf("BASIC", "NTLM", "KERBEROS", "OAUTH", "APP_PERM"),
				},
				MarkdownDescription: "The method used to authenticate to Exchange",
			},
			"auth_provider": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("ADFS"),
				Validators: []validator.String{
					stringvalidator.OneOf("ADFS", "AZURE"),
				},
				MarkdownDescription: "The method by which users will sign into the Outlook add-in.",
			},
			"uuid": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "UUID for the Exchange connector.",
			},
			"scheduled_alias_prefix": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Prefix for scheduled conference aliases.",
			},
			"scheduled_alias_domain": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					validators.Domain(),
				},
				MarkdownDescription: "Domain for scheduled conference aliases.",
			},
			"scheduled_alias_suffix_length": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(6),
				MarkdownDescription: "The length of the random number suffix part of aliases used for scheduled conferences. Range: 5 to 15. Default: 6.",
			},
			"meeting_buffer_before": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(30),
				MarkdownDescription: "The number of minutes before the meeting's scheduled start time that participants will be able to join the VMR. Range: 0 to 180. Default: 30.",
			},
			"meeting_buffer_after": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(60),
				MarkdownDescription: "The number of minutes after the meeting's scheduled end of a conference participants will be able to join the VMR. Range: 0 to 180. Default: 60.",
			},
			"enable_dynamic_vmrs": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
				MarkdownDescription: "Enable this option to allow Outlook users to schedule meetings in single-use (randomly generated) VMRs.",
			},
			"enable_personal_vmrs": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Enable this option to allow Outlook users to schedule meetings in their personal VMRs.",
			},
			"allow_new_users": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
				MarkdownDescription: "Disable this option to allow only those users with an existing User record to access the Outlook add-in.",
			},
			"disable_proxy": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Disable the usage of any web proxy which may have been configured on the Management Node by this Secure Scheduler for Exchange Integration.",
			},
			"use_custom_add_in_sources": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Enable this to specify custom locations to serve add-in JavaScript and CSS from. This can be used to support offline deployments.",
			},
			"enable_addin_debug_logs": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Enable this option to view debug logs within the add-in side pane. Note that these logs will appear for all users of this add-in.",
			},
			"oauth_client_id": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "OAuth client ID for authentication.",
			},
			"oauth_client_secret": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				MarkdownDescription: "OAuth client secret for authentication. This field is sensitive.",
			},
			"oauth_auth_endpoint": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
				MarkdownDescription: "OAuth authorization endpoint URL.",
			},
			"oauth_token_endpoint": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "OAuth token endpoint URL.",
			},
			"oauth_redirect_uri": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
				MarkdownDescription: "OAuth redirect URI for authentication flow.",
			},
			"oauth_refresh_token": schema.StringAttribute{
				Computed:            true,
				Sensitive:           true,
				MarkdownDescription: "OAuth refresh token. This field is sensitive.",
			},
			"oauth_state": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "OAuth state parameter for security.",
			},
			"kerberos_realm": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
				MarkdownDescription: "Kerberos realm for authentication.",
			},
			"kerberos_kdc": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
				MarkdownDescription: "Kerberos Key Distribution Center.",
			},
			"kerberos_kdc_https_proxy": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
				MarkdownDescription: "HTTPS proxy for Kerberos KDC connections.",
			},
			"kerberos_exchange_spn": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
				MarkdownDescription: "Kerberos Service Principal Name for Exchange.",
			},
			"kerberos_enable_tls": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
				MarkdownDescription: "If enabled, all communication to the KDC will go through an HTTPS proxy and all traffic to the KDC will be encrypted using TLS.",
			},
			"kerberos_auth_every_request": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "When Kerberos authentication is enabled, send a Kerberos Authorization header in every request to the Exchange server.",
			},
			"kerberos_verify_tls_using_custom_ca": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "If enabled, use the configured Root Trust CA Certificates to verify the KDC HTTPS proxy SSL certificate. If disabled, the HTTPS proxy SSL certificate is verified using the system-wide default set of trusted certificates.",
			},
			"addin_server_domain": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The FQDN of the reverse proxy or Conferencing Node that provides the add-in content. The FQDN must have a valid certificate. Maximum length: 192 characters.",
			},
			"addin_display_name": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("Pexip Scheduling Service"),
				MarkdownDescription: "The display name of the add-in. Maximum length: 250 characters.",
			},
			"addin_description": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("Turns meetings into Pexip meetings"),
				MarkdownDescription: "The description of the add-in. Maximum length: 250 characters.",
			},
			"addin_provider_name": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("Pexip"),
				MarkdownDescription: "The name of the organization which provides the add-in. Maximum length: 250 characters.",
			},
			"addin_button_label": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("Create a Pexip meeting"),
				MarkdownDescription: "The label for the add-in button on desktop clients. Maximum length: 250 characters.",
			},
			"addin_group_label": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("Pexip meeting"),
				MarkdownDescription: "The name of the group in which to place the add-in button on desktop clients. Maximum length: 250 characters.",
			},
			"addin_supertip_title": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("Makes this a Pexip meeting"),
				MarkdownDescription: "The title of the supertip help text for the add-in button on desktop clients. Maximum length: 250 characters.",
			},
			"addin_supertip_description": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("Turns this meeting into an audio or video conference hosted in a Pexip VMR. The meeting is not scheduled until you select Send."),
				MarkdownDescription: "The text of the supertip for the add-in button on desktop clients. Maximum length: 250 characters.",
			},
			"addin_application_id": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Application ID for the Exchange add-in.",
			},
			"addin_authority_url": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
				MarkdownDescription: "The Authority URL copied from the App Registration created in Microsoft Entra for add-in authentication. Maximum length: 255 characters.",
			},
			"addin_oidc_metadata_url": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
				MarkdownDescription: "The OpenID Connect metadata document copied from the App Registration created in Microsoft Entra for add-in authentication. Maximum length: 255 characters.",
			},
			"addin_authentication_method": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("EXCHANGE_USER_ID_TOKEN"),
				Validators: []validator.String{
					stringvalidator.OneOf("EXCHANGE_USER_ID_TOKEN", "SSO_TOKEN", "NAA_TOKEN"),
				},
				MarkdownDescription: "The type of token the Outlook add-in uses to authenticate to Pexip",
			},
			"addin_naa_web_api_application_id": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "NAA Web API application ID for add-in authentication.",
			},
			"personal_vmr_oauth_client_id": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("4189c2b4-92ca-416c-b7ea-bc3cfab3d0f0"),
				MarkdownDescription: "The client ID of the OAuth application used to authenticate users when signing in to the Outlook add-in.",
			},
			"personal_vmr_oauth_client_secret": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
				Sensitive:           true,
				MarkdownDescription: "The client secret of the OAuth application created for signing in users in the Outlook add-in.",
			},
			"personal_vmr_oauth_auth_endpoint": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
				MarkdownDescription: "The authorization URI of the OAuth application used to authenticate users when signing in to the Outlook add-in. Maximum length: 255 characters.",
			},
			"personal_vmr_oauth_token_endpoint": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
				MarkdownDescription: "The token URI of the OAuth application used to authenticate users when signing in to the Outlook add-in. Maximum length: 255 characters.",
			},
			"personal_vmr_adfs_relying_party_trust_identifier": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
				MarkdownDescription: "The URL which identifies the OAuth 2.0 resource on AD FS. Maximum length: 255 characters.",
			},
			"office_js_url": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("https://appsforoffice.microsoft.com/lib/1/hosted/office.js"),
				MarkdownDescription: "The URL used to download the Office.js JavaScript library. Maximum length: 255 characters.",
			},
			"microsoft_fabric_url": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("https://appsforoffice.microsoft.com/fabric/1.0/fabric.min.css"),
				MarkdownDescription: "The URL used to download the Microsoft Fabric CSS. Maximum length: 255 characters.",
			},
			"microsoft_fabric_components_url": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("https://appsforoffice.microsoft.com/fabric/1.0/fabric.components.min.css"),
				MarkdownDescription: "The URL used to download the Microsoft Fabric Components CSS. Maximum length: 255 characters.",
			},
			"additional_add_in_script_sources": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
				MarkdownDescription: "Optionally specify additional URLs to download JavaScript script files. Each URL must be entered on a separate line. Maximum length: 4096 characters.",
			},
			"domains": schema.SetAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				MarkdownDescription: "List of Exchange Metadata Domain URIs associated with this connector.",
			},
			"host_identity_provider_group": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "URI reference to host identity provider group resource.",
			},
			"ivr_theme": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "URI reference to IVR theme resource.",
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
			"private_key": schema.StringAttribute{
				Computed:            true,
				Sensitive:           true,
				MarkdownDescription: "Private key for Exchange connector. This field is sensitive and computed.",
			},
			"public_key": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Public key for Exchange connector. This field is computed.",
			},
			"accept_edited_occurrence_template": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("<div style=\"font-size:11.0pt; color:#000000; font-family:Calibri,Arial,Helvetica,sans-serif;\">\r\nThis meeting occurrence in a recurring series has been successfully rescheduled using the aliases: {{alias}} and {{numeric_alias}}.<br>\r\n</div>"),
				MarkdownDescription: "A Jinja2 template that is used to produce the message sent to meeting organizers once the scheduling service successfully schedules an edited occurrence in a recurring series. Maximum length: 12288 characters.",
			},
			"accept_edited_recurring_series_template": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("<div style=\"font-size:11.0pt; color:#000000; font-family:Calibri,Arial,Helvetica,sans-serif;\">\r\nThis recurring meeting series has been successfully rescheduled.<br>\r\nAll meetings in this series will use the aliases: {{alias}} and {{numeric_alias}}.<br>\r\n</div>"),
				MarkdownDescription: "A Jinja2 template that is used to produce the message sent to meeting organizers once the scheduling service successfully schedules an edited recurring meeting. Maximum length: 12288 characters.",
			},
			"accept_edited_single_meeting_template": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("<div style=\"font-size:11.0pt; color:#000000; font-family:Calibri,Arial,Helvetica,sans-serif;\">\r\nThis meeting has been successfully rescheduled using the aliases: {{alias}} and {{numeric_alias}}.<br>\r\n</div>"),
				MarkdownDescription: "A Jinja2 template that is used to produce the message sent to meeting organizers once the scheduling service successfully schedules an edited single meeting. Maximum length: 12288 characters.",
			},
			"accept_new_recurring_series_template": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("<div style=\"font-size:11.0pt; color:#000000; font-family:Calibri,Arial,Helvetica,sans-serif;\">\r\nThis recurring meeting series has been successfully scheduled.<br>\r\nAll meetings in this series will use the aliases: {{alias}} and {{numeric_alias}}.<br>\r\n</div>"),
				MarkdownDescription: "A Jinja2 template that is used to produce the message sent to meeting organizers once the scheduling service successfully schedules a new recurring meeting. Maximum length: 12288 characters.",
			},
			"accept_new_single_meeting_template": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("<div style=\"font-size:11.0pt; color:#000000; font-family:Calibri,Arial,Helvetica,sans-serif;\">\r\nThis meeting has been successfully scheduled using the aliases: {{alias}} and {{numeric_alias}}.<br>\r\n</div>"),
				MarkdownDescription: "A Jinja2 template that is used to produce the message sent to meeting organizers once the scheduling service successfully schedules a new single meeting. Maximum length: 12288 characters.",
			},
			"conference_description_template": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("Scheduled Conference booked by {{organizer_email}}"),
				MarkdownDescription: "A Jinja2 template that is used to produce the description of scheduled conferences. Maximum length: 12288 characters.",
			},
			"conference_name_template": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("{{subject}} ({{organizer_name}})"),
				MarkdownDescription: "A Jinja2 template that is used to produce the name of scheduled conferences. Please note conference names must be unique so a random number may be appended if the name that is generated is already in use by another service. Maximum length: 12288 characters.",
			},
			"conference_subject_template": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("{{subject}}"),
				MarkdownDescription: "A Jinja2 template that is used to produce the subject field of scheduled conferences. By default this will use the subject line of the meeting invitation but this field can be deleted or amended if you do not want the subject to be visible to administrators. Maximum length: 12288 characters.",
			},
			"meeting_instructions_template": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("<br>\r\n<div style=\"font-size:11.0pt; color:#000000; font-family:Calibri,Arial,Helvetica,sans-serif;\">\r\n<b>Please join my Pexip Virtual Meeting Room in one of the following ways:</b><br>\r\n<br>\r\nFrom a VC endpoint or a Skype/Lync client:<br>\r\n{{alias}}<br>\r\n<br>\r\nFrom a web browser:<br>\r\n<a href=\"https://{{addin_server_domain}}/webapp/#/?conference={{alias}}\">https://{{addin_server_domain}}/webapp/#/?conference={{alias}}</a><br>\r\n<br>\r\nFrom a Pexip Infinity Connect client:<br>\r\npexip://{{alias}}<br>\r\n<br>\r\nFrom a telephone:<br>\r\n[Your number], then {{numeric_alias}} #<br>\r\n<br>\r\n{{alias_uuid}}<br>\r\n</div>"),
				MarkdownDescription: "A Jinja2 template that is used to generate the instructions added by the scheduling service to the body of the meeting request when a single-use VMR is being used. Maximum length: 12288 characters.",
			},
			"personal_vmr_description_template": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("{{description}}"),
				MarkdownDescription: "A Jinja2 template that is used to generate the description of the personal VMR, shown to users when they hover over the button. Maximum length: 12288 characters.",
			},
			"personal_vmr_instructions_template": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("{% if domain_aliases %}\r\n    {% set alias = domain_aliases[0] %}\r\n{% elif other_aliases %}\r\n    {% set alias = other_aliases[0] %}\r\n{% else %}\r\n    {% set alias = numeric_aliases[0] %}\r\n{% endif %}\r\n{% if (not allow_guests) and pin %}\r\n    {% set meeting_pin = pin %}\r\n{% elif allow_guests and guest_pin %}\r\n    {% set meeting_pin = guest_pin %}\r\n{% else %}\r\n    {% set meeting_pin = \"\" %}\r\n{% endif %}\r\n<br>\r\n<div style=\"font-size:11.0pt; color:#000000; font-family:Calibri,Arial,Helvetica,sans-serif;\">\r\n<b>Please join my Pexip Virtual Meeting Room in one of the following ways:</b><br>\r\n<br>\r\nFrom a VC endpoint or a Skype/Lync client:<br>\r\n{{alias}}<br>\r\n<br>\r\nFrom a web browser:<br>\r\n<a href=\"https://{{addin_server_domain}}/webapp/#/?conference={{alias}}\">https://{{addin_server_domain}}/webapp/#/?conference={{alias}}</a><br>\r\n<br>\r\nFrom a Pexip Infinity Connect client:<br>\r\npexip://{{alias}}<br>\r\n<br>\r\n{% if numeric_aliases %}\r\nFrom a telephone:<br>\r\n[Your number], then {{numeric_aliases[0]}} #<br>\r\n<br>\r\n{% endif %}\r\n{% if meeting_pin %}\r\nPlease join using the PIN <b>{{meeting_pin}}</b><br>\r\n<br>\r\n{% endif %}\r\n</div>"),
				MarkdownDescription: "A Jinja2 template that is used to produce the joining instructions added by the scheduling service to the body of the meeting request when a personal VMR is being used. Maximum length: 12288 characters.",
			},
			"personal_vmr_location_template": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("{% if domain_aliases %}\r\n    {% set alias = domain_aliases[0] %}\r\n{% elif other_aliases %}\r\n    {% set alias = other_aliases[0] %}\r\n{% else %}\r\n    {% set alias = numeric_aliases[0] %}\r\n{% endif %}\r\nhttps://{{addin_server_domain}}/webapp/#/?conference={{alias}}"),
				MarkdownDescription: "A Jinja2 template that is used to generate the text that will be inserted into the Location field of the meeting request when a personal VMR is being used. Maximum length: 12288 characters.",
			},
			"personal_vmr_name_template": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("{{name}}"),
				MarkdownDescription: "A Jinja2 template that is used to generate the name of the personal VMR, as it appears on the button offered to users. Maximum length: 12288 characters.",
			},
			"placeholder_instructions_template": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("<div style=\"font-size:11.0pt; color:#000000; font-family:Calibri,Arial,Helvetica,sans-serif;\">\r\nThis meeting will be hosted in a Virtual Meeting Room. Joining instructions will be<br>\r\nsent to you soon in a separate email.<br>\r\n</div>"),
				MarkdownDescription: "The text that is added by the scheduling service to email messages when the actual joining instructions cannot be obtained. Maximum length: 12288 characters.",
			},
			"reject_alias_conflict_template": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("<div style=\"font-size:11.0pt; color:#000000; font-family:Calibri,Arial,Helvetica,sans-serif;\">\r\nWe are unable to schedule this meeting because the alias: {{alias}} is already <br>\r\nin use by another Pexip Virtual Meeting Room. Please try creating a new meeting.<br>\r\n</div>"),
				MarkdownDescription: "A Jinja2 template that is used to produce the message sent to meeting organizers when the scheduling service fails to schedule a meeting because the alias conflicts with an existing alias. Maximum length: 12288 characters.",
			},
			"reject_alias_deleted_template": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("<div style=\"font-size:11.0pt; color:#000000; font-family:Calibri,Arial,Helvetica,sans-serif;\">\r\nWe are unable to schedule this meeting because its alias has been deleted.<br>\r\nPlease try creating a new meeting.<br>\r\n</div>"),
				MarkdownDescription: "The text that is sent to meeting organizers when the scheduling service fails to schedule a meeting because the alias for this meeting has been deleted. Maximum length: 12288 characters.",
			},
			"reject_general_error_template": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("<div style=\"font-size:11.0pt; color:#000000; font-family:Calibri,Arial,Helvetica,sans-serif;\">\r\nWe are unable to schedule this meeting. Please try creating a new meeting.<br>\r\nIf this issue continues, please forward this message to your system administrator, including the following ID:<br>\r\nCorrelationID=\"{{correlation_id}}\".<br>\r\n</div>"),
				MarkdownDescription: "A Jinja2 template that is used to produce the message sent to meeting organizers when the scheduling service fails to schedule a meeting because a general error occurred. Maximum length: 12288 characters.",
			},
			"reject_invalid_alias_id_template": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("<div style=\"font-size:11.0pt; color:#000000; font-family:Calibri,Arial,Helvetica,sans-serif;\">\r\nThis meeting request does not contain currently valid scheduling data, and therefore cannot be processed.<br>\r\nPlease use the add-in to create a new meeting request, without editing any of the content that is inserted by the add-in.<br>\r\nIf this issue continues, please contact your system administrator.<br>\r\n</div>"),
				MarkdownDescription: "The text that is sent to meeting organizers when the scheduling service fails to schedule a meeting because the alias ID in the meeting email is invalid. Maximum length: 12288 characters.",
			},
			"reject_recurring_series_past_template": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("<div style=\"font-size:11.0pt; color:#000000; font-family:Calibri,Arial,Helvetica,sans-serif;\">\r\nThis recurring series cannot be scheduled because all<br>\r\noccurrences happen in the past.<br>\r\n</div>"),
				MarkdownDescription: "The text that is sent to meeting organizers when the scheduling service fails to schedule a recurring meeting because all occurrences occur in the past. Maximum length: 12288 characters.",
			},
			"reject_single_meeting_past": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("<div style=\"font-size:11.0pt; color:#000000; font-family:Calibri,Arial,Helvetica,sans-serif;\">\r\nThis meeting cannot be scheduled because it occurs in the past.<br>\r\n</div>"),
				MarkdownDescription: "The text that is sent to meeting organizers when the scheduling service fails to schedule a meeting because it occurs in the past. Maximum length: 12288 characters.",
			},
			"scheduled_alias_description_template": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("Scheduled Conference booked by {{organizer_email}}"),
				MarkdownDescription: "A Jinja2 template that is used to produce the description of scheduled conference aliases. Maximum length: 12288 characters.",
			},
			"addin_pane_already_video_meeting_heading": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("VMR already assigned"),
				MarkdownDescription: "The heading that appears on the side pane when the add-in is activated after an alias has already been obtained for the meeting. Maximum length: 250 characters.",
			},
			"addin_pane_already_video_meeting_message": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("It looks like this meeting has already been set up to be hosted in a Virtual Meeting Room. If this is a new meeting, select Send to schedule the conference."),
				MarkdownDescription: "The message that appears on the side pane when the add-in is activated after an alias has already been obtained for the meeting. Maximum length: 250 characters.",
			},
			"addin_pane_button_title": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("Add a Single-use VMR"),
				MarkdownDescription: "The label of the button on the side pane used to add a single-use VMR. Maximum length: 250 characters.",
			},
			"addin_pane_description": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("This assigns a Virtual Meeting Room for your meeting"),
				MarkdownDescription: "The description of the add-in on the side pane. Maximum length: 250 characters.",
			},
			"addin_pane_general_error_heading": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("Error"),
				MarkdownDescription: "The heading that appears on the side pane when an error occurs trying to add the joining instructions. Maximum length: 250 characters.",
			},
			"addin_pane_general_error_message": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("There was a problem adding the joining instructions. Please try again."),
				MarkdownDescription: "The message that appears on the side pane when an error occurs trying to add the joining instructions of the single-use VMR. Maximum length: 250 characters.",
			},
			"addin_pane_management_node_down_heading": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("Cannot assign a VMR right now"),
				MarkdownDescription: "The heading that appears on the side pane when the Management Node cannot be contacted to obtain an alias. Maximum length: 250 characters.",
			},
			"addin_pane_management_node_down_message": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("Sorry, we are unable to assign a Virtual Meeting Room at this time. Select Send to schedule the meeting, and all attendees will be sent joining instructions later."),
				MarkdownDescription: "The message that appears on the side pane when the Management Node cannot be contacted to obtain an alias. Maximum length: 250 characters.",
			},
			"addin_pane_personal_vmr_add_button": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("Add a Personal VMR"),
				MarkdownDescription: "The label of the button on the side pane used to add a personal VMR. Maximum length: 250 characters.",
			},
			"addin_pane_personal_vmr_error_getting_message": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("There was a problem getting your personal VMRs. Please try again."),
				MarkdownDescription: "The message that appears on the side pane when an error occurs trying to obtain a list of the user's personal VMRs. Maximum length: 250 characters.",
			},
			"addin_pane_personal_vmr_error_inserting_meeting_message": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("There was a problem adding the joining instructions. Please try again."),
				MarkdownDescription: "The message that appears on the side pane when an error occurs trying to add the personal VMR details to the meeting. Maximum length: 250 characters.",
			},
			"addin_pane_personal_vmr_error_signing_in_message": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("There was a problem signing you in. Please try again."),
				MarkdownDescription: "The message that appears on the side pane when an error occurs trying to sign the user in. Maximum length: 250 characters.",
			},
			"addin_pane_personal_vmr_none_message": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("You do not have any personal VMRs"),
				MarkdownDescription: "The message that appears on the side pane when the user has no personal VMRs. Maximum length: 250 characters.",
			},
			"addin_pane_personal_vmr_select_message": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("Select the VMR you want to add to the meeting"),
				MarkdownDescription: "The message that appears on the side pane requesting users to select a personal VMR to use for the meeting. Maximum length: 250 characters.",
			},
			"addin_pane_personal_vmr_sign_in_button": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("Sign In"),
				MarkdownDescription: "The label of the button on the side pane requesting users to sign in to obtain a list of their personal VMRs. Maximum length: 250 characters.",
			},
			"addin_pane_success_heading": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("Success"),
				MarkdownDescription: "The heading that appears on the side pane when when an alias has been obtained successfully from the Management Node. Maximum length: 250 characters.",
			},
			"addin_pane_success_message": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("This meeting is now set up to be hosted as an audio or video conference in a Virtual Meeting Room. Please note this conference is not scheduled until you select Send."),
				MarkdownDescription: "The message that appears on the side pane when when an alias has been obtained successfully from the Management Node. Maximum length: 250 characters.",
			},
			"addin_pane_title": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("Add a VMR"),
				MarkdownDescription: "The title of the add-in on the side pane. Maximum length: 250 characters.",
			},
		},
		MarkdownDescription: "Manages a Microsoft Exchange connector with the Infinity service. Exchange connectors enable integration with Microsoft Exchange/Office 365 for calendar and meeting management.",
	}
}

func (r *InfinityMsExchangeConnectorResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	plan := &InfinityMsExchangeConnectorResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert List attributes to []string
	exchangeDomains, diags := getStringList(ctx, plan.Domains)
	resp.Diagnostics.Append(diags...)

	createRequest := &config.MsExchangeConnectorCreateRequest{
		Name:                           plan.Name.ValueString(),
		Description:                    plan.Description.ValueString(),
		RoomMailboxName:                plan.RoomMailboxName.ValueString(),
		URL:                            plan.URL.ValueString(),
		Username:                       plan.Username.ValueString(),
		Password:                       plan.Password.ValueString(),
		AuthenticationMethod:           plan.AuthenticationMethod.ValueString(),
		AuthProvider:                   plan.AuthProvider.ValueString(),
		UUID:                           plan.UUID.ValueString(),
		ScheduledAliasDomain:           plan.ScheduledAliasDomain.ValueString(),
		ScheduledAliasSuffixLength:     int(plan.ScheduledAliasSuffixLength.ValueInt64()),
		MeetingBufferBefore:            int(plan.MeetingBufferBefore.ValueInt64()),
		MeetingBufferAfter:             int(plan.MeetingBufferAfter.ValueInt64()),
		EnableDynamicVmrs:              plan.EnableDynamicVmrs.ValueBool(),
		EnablePersonalVmrs:             plan.EnablePersonalVmrs.ValueBool(),
		AllowNewUsers:                  plan.AllowNewUsers.ValueBool(),
		DisableProxy:                   plan.DisableProxy.ValueBool(),
		UseCustomAddInSources:          plan.UseCustomAddInSources.ValueBool(),
		EnableAddinDebugLogs:           plan.EnableAddinDebugLogs.ValueBool(),
		OauthClientSecret:              plan.OauthClientSecret.ValueString(),
		OauthAuthEndpoint:              plan.OauthAuthEndpoint.ValueString(),
		OauthTokenEndpoint:             plan.OauthTokenEndpoint.ValueString(),
		OauthRedirectURI:               plan.OauthRedirectURI.ValueString(),
		OauthRefreshToken:              plan.OauthRefreshToken.ValueString(),
		KerberosRealm:                  plan.KerberosRealm.ValueString(),
		KerberosKdc:                    plan.KerberosKdc.ValueString(),
		KerberosKdcHttpsProxy:          plan.KerberosKdcHttpsProxy.ValueString(),
		KerberosExchangeSpn:            plan.KerberosExchangeSpn.ValueString(),
		KerberosEnableTls:              plan.KerberosEnableTls.ValueBool(),
		KerberosAuthEveryRequest:       plan.KerberosAuthEveryRequest.ValueBool(),
		KerberosVerifyTlsUsingCustomCa: plan.KerberosVerifyTlsUsingCustomCa.ValueBool(),
		AddinServerDomain:              plan.AddinServerDomain.ValueString(),
		AddinDisplayName:               plan.AddinDisplayName.ValueString(),
		AddinDescription:               plan.AddinDescription.ValueString(),
		AddinProviderName:              plan.AddinProviderName.ValueString(),
		AddinButtonLabel:               plan.AddinButtonLabel.ValueString(),
		AddinGroupLabel:                plan.AddinGroupLabel.ValueString(),
		AddinSupertipTitle:             plan.AddinSupertipTitle.ValueString(),
		AddinSupertipDescription:       plan.AddinSupertipDescription.ValueString(),
		AddinAuthorityURL:              plan.AddinAuthorityURL.ValueString(),
		AddinOidcMetadataURL:           plan.AddinOidcMetadataURL.ValueString(),
		AddinAuthenticationMethod:      plan.AddinAuthenticationMethod.ValueString(),
		PersonalVmrOauthClientSecret:   plan.PersonalVmrOauthClientSecret.ValueString(),
		PersonalVmrOauthAuthEndpoint:   plan.PersonalVmrOauthAuthEndpoint.ValueString(),
		PersonalVmrOauthTokenEndpoint:  plan.PersonalVmrOauthTokenEndpoint.ValueString(),
		PersonalVmrAdfsRelyingPartyTrustIdentifier:         plan.PersonalVmrAdfsRelyingPartyTrustIdentifier.ValueString(),
		OfficeJsURL:                                        plan.OfficeJsURL.ValueString(),
		MicrosoftFabricURL:                                 plan.MicrosoftFabricURL.ValueString(),
		MicrosoftFabricComponentsURL:                       plan.MicrosoftFabricComponentsURL.ValueString(),
		AdditionalAddInScriptSources:                       plan.AdditionalAddInScriptSources.ValueString(),
		NonIdpParticipants:                                 plan.NonIdpParticipants.ValueString(),
		AcceptEditedOccurrenceTemplate:                     plan.AcceptEditedOccurrenceTemplate.ValueString(),
		AcceptEditedRecurringSeriesTemplate:                plan.AcceptEditedRecurringSeriesTemplate.ValueString(),
		AcceptEditedSingleMeetingTemplate:                  plan.AcceptEditedSingleMeetingTemplate.ValueString(),
		AcceptNewRecurringSeriesTemplate:                   plan.AcceptNewRecurringSeriesTemplate.ValueString(),
		AcceptNewSingleMeetingTemplate:                     plan.AcceptNewSingleMeetingTemplate.ValueString(),
		ConferenceDescriptionTemplate:                      plan.ConferenceDescriptionTemplate.ValueString(),
		ConferenceNameTemplate:                             plan.ConferenceNameTemplate.ValueString(),
		ConferenceSubjectTemplate:                          plan.ConferenceSubjectTemplate.ValueString(),
		MeetingInstructionsTemplate:                        plan.MeetingInstructionsTemplate.ValueString(),
		PersonalVmrDescriptionTemplate:                     plan.PersonalVmrDescriptionTemplate.ValueString(),
		PersonalVmrInstructionsTemplate:                    plan.PersonalVmrInstructionsTemplate.ValueString(),
		PersonalVmrLocationTemplate:                        plan.PersonalVmrLocationTemplate.ValueString(),
		PersonalVmrNameTemplate:                            plan.PersonalVmrNameTemplate.ValueString(),
		PlaceholderInstructionsTemplate:                    plan.PlaceholderInstructionsTemplate.ValueString(),
		RejectAliasConflictTemplate:                        plan.RejectAliasConflictTemplate.ValueString(),
		RejectAliasDeletedTemplate:                         plan.RejectAliasDeletedTemplate.ValueString(),
		RejectGeneralErrorTemplate:                         plan.RejectGeneralErrorTemplate.ValueString(),
		RejectInvalidAliasIDTemplate:                       plan.RejectInvalidAliasIDTemplate.ValueString(),
		RejectRecurringSeriesPastTemplate:                  plan.RejectRecurringSeriesPastTemplate.ValueString(),
		RejectSingleMeetingPast:                            plan.RejectSingleMeetingPast.ValueString(),
		ScheduledAliasDescriptionTemplate:                  plan.ScheduledAliasDescriptionTemplate.ValueString(),
		AddinPaneAlreadyVideoMeetingHeading:                plan.AddinPaneAlreadyVideoMeetingHeading.ValueString(),
		AddinPaneAlreadyVideoMeetingMessage:                plan.AddinPaneAlreadyVideoMeetingMessage.ValueString(),
		AddinPaneButtonTitle:                               plan.AddinPaneButtonTitle.ValueString(),
		AddinPaneDescription:                               plan.AddinPaneDescription.ValueString(),
		AddinPaneGeneralErrorHeading:                       plan.AddinPaneGeneralErrorHeading.ValueString(),
		AddinPaneGeneralErrorMessage:                       plan.AddinPaneGeneralErrorMessage.ValueString(),
		AddinPaneManagementNodeDownHeading:                 plan.AddinPaneManagementNodeDownHeading.ValueString(),
		AddinPaneManagementNodeDownMessage:                 plan.AddinPaneManagementNodeDownMessage.ValueString(),
		AddinPanePersonalVmrAddButton:                      plan.AddinPanePersonalVmrAddButton.ValueString(),
		AddinPanePersonalVmrErrorGettingMessage:            plan.AddinPanePersonalVmrErrorGettingMessage.ValueString(),
		AddinPanePersonalVmrErrorInsertingMeetingMessage:   plan.AddinPanePersonalVmrErrorInsertingMeetingMessage.ValueString(),
		AddinPanePersonalVmrErrorSigningInMessage:          plan.AddinPanePersonalVmrErrorSigningInMessage.ValueString(),
		AddinPanePersonalVmrNoneMessage:                    plan.AddinPanePersonalVmrNoneMessage.ValueString(),
		AddinPanePersonalVmrSelectMessage:                  plan.AddinPanePersonalVmrSelectMessage.ValueString(),
		AddinPanePersonalVmrSignInButton:                   plan.AddinPanePersonalVmrSignInButton.ValueString(),
		AddinPaneSuccessHeading:                            plan.AddinPaneSuccessHeading.ValueString(),
		AddinPaneSuccessMessage:                            plan.AddinPaneSuccessMessage.ValueString(),
		AddinPaneTitle:                                     plan.AddinPaneTitle.ValueString(),
	}

	// Handle optional pointer fields
	if !plan.RoomMailboxEmailAddress.IsNull() && !plan.RoomMailboxEmailAddress.IsUnknown() {
		email := plan.RoomMailboxEmailAddress.ValueString()
		createRequest.RoomMailboxEmailAddress = &email
	}

	if !plan.ScheduledAliasPrefix.IsNull() && !plan.ScheduledAliasPrefix.IsUnknown() {
		prefix := plan.ScheduledAliasPrefix.ValueString()
		createRequest.ScheduledAliasPrefix = &prefix
	}

	if !plan.OauthClientID.IsNull() && !plan.OauthClientID.IsUnknown() {
		clientID := plan.OauthClientID.ValueString()
		createRequest.OauthClientID = &clientID
	}

	if !plan.AddinApplicationID.IsNull() && !plan.AddinApplicationID.IsUnknown() {
		appID := plan.AddinApplicationID.ValueString()
		createRequest.AddinApplicationID = &appID
	}

	if !plan.AddinNaaWebApiApplicationID.IsNull() && !plan.AddinNaaWebApiApplicationID.IsUnknown() {
		apiAppID := plan.AddinNaaWebApiApplicationID.ValueString()
		createRequest.AddinNaaWebApiApplicationID = &apiAppID
	}

	if !plan.PersonalVmrOauthClientID.IsNull() && !plan.PersonalVmrOauthClientID.IsUnknown() {
		vmrClientID := plan.PersonalVmrOauthClientID.ValueString()
		createRequest.PersonalVmrOauthClientID = &vmrClientID
	}

	if !plan.Domains.IsNull() && !plan.Domains.IsUnknown() {
		domains := exchangeDomains
		createRequest.Domains = &domains
	}

	if !plan.HostIdentityProviderGroup.IsNull() && !plan.HostIdentityProviderGroup.IsUnknown() {
		hostIdp := plan.HostIdentityProviderGroup.ValueString()
		createRequest.HostIdentityProviderGroup = &hostIdp
	}

	if !plan.IvrTheme.IsNull() && !plan.IvrTheme.IsUnknown() {
		theme := plan.IvrTheme.ValueString()
		createRequest.IvrTheme = &theme
	}

	createResponse, err := r.InfinityClient.Config().CreateMsExchangeConnector(ctx, createRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Infinity Microsoft Exchange connector",
			fmt.Sprintf("Could not create Infinity Microsoft Exchange connector: %s", err),
		)
		return
	}

	resourceID, err := createResponse.ResourceID()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Retrieving Infinity Microsoft Exchange connector ID",
			fmt.Sprintf("Could not retrieve ID for created Infinity Microsoft Exchange connector: %s", err),
		)
		return
	}

	// Read the state from the API to get all computed values
	model, err := r.read(ctx, resourceID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Created Infinity Microsoft Exchange connector",
			fmt.Sprintf("Could not read created Infinity Microsoft Exchange connector with ID %d: %s", resourceID, err),
		)
		return
	}

	// Preserve write-only sensitive fields from plan
	model.Password = plan.Password
	model.OauthClientSecret = plan.OauthClientSecret
	model.PersonalVmrOauthClientSecret = plan.PersonalVmrOauthClientSecret

	tflog.Trace(ctx, fmt.Sprintf("created Infinity Microsoft Exchange connector with ID: %s, name: %s", model.ID, model.Name))

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *InfinityMsExchangeConnectorResource) read(ctx context.Context, resourceID int) (*InfinityMsExchangeConnectorResourceModel, error) {
	var data InfinityMsExchangeConnectorResourceModel

	srv, err := r.InfinityClient.Config().GetMsExchangeConnector(ctx, resourceID)
	if err != nil {
		return nil, err
	}

	if srv.ResourceURI == "" {
		return nil, fmt.Errorf("microsoft Exchange connector with ID %d not found", resourceID)
	}

	data.ID = types.StringValue(srv.ResourceURI)
	data.ResourceID = types.Int32Value(int32(resourceID)) // #nosec G115 -- API values are expected to be within int32 range
	data.Name = types.StringValue(srv.Name)
	data.Description = types.StringValue(srv.Description)
	data.RoomMailboxName = types.StringValue(srv.RoomMailboxName)
	data.URL = types.StringValue(srv.URL)
	data.Username = types.StringValue(srv.Username)
	// Password is write-only and not returned by the API, will be preserved from state/plan
	data.AuthenticationMethod = types.StringValue(srv.AuthenticationMethod)
	data.AuthProvider = types.StringValue(srv.AuthProvider)
	data.UUID = types.StringValue(srv.UUID)
	data.ScheduledAliasDomain = types.StringValue(srv.ScheduledAliasDomain)
	data.ScheduledAliasSuffixLength = types.Int64Value(int64(srv.ScheduledAliasSuffixLength))
	data.MeetingBufferBefore = types.Int64Value(int64(srv.MeetingBufferBefore))
	data.MeetingBufferAfter = types.Int64Value(int64(srv.MeetingBufferAfter))
	data.EnableDynamicVmrs = types.BoolValue(srv.EnableDynamicVmrs)
	data.EnablePersonalVmrs = types.BoolValue(srv.EnablePersonalVmrs)
	data.AllowNewUsers = types.BoolValue(srv.AllowNewUsers)
	data.DisableProxy = types.BoolValue(srv.DisableProxy)
	data.UseCustomAddInSources = types.BoolValue(srv.UseCustomAddInSources)
	data.EnableAddinDebugLogs = types.BoolValue(srv.EnableAddinDebugLogs)
	// OauthClientSecret is write-only and not returned by the API, will be preserved from state/plan
	data.OauthAuthEndpoint = types.StringValue(srv.OauthAuthEndpoint)
	data.OauthTokenEndpoint = types.StringValue(srv.OauthTokenEndpoint)
	data.OauthRedirectURI = types.StringValue(srv.OauthRedirectURI)
	data.OauthRefreshToken = types.StringValue(srv.OauthRefreshToken)
	data.KerberosRealm = types.StringValue(srv.KerberosRealm)
	data.KerberosKdc = types.StringValue(srv.KerberosKdc)
	data.KerberosKdcHttpsProxy = types.StringValue(srv.KerberosKdcHttpsProxy)
	data.KerberosExchangeSpn = types.StringValue(srv.KerberosExchangeSpn)
	data.KerberosEnableTls = types.BoolValue(srv.KerberosEnableTls)
	data.KerberosAuthEveryRequest = types.BoolValue(srv.KerberosAuthEveryRequest)
	data.KerberosVerifyTlsUsingCustomCa = types.BoolValue(srv.KerberosVerifyTlsUsingCustomCa)
	data.AddinServerDomain = types.StringValue(srv.AddinServerDomain)
	data.AddinDisplayName = types.StringValue(srv.AddinDisplayName)
	data.AddinDescription = types.StringValue(srv.AddinDescription)
	data.AddinProviderName = types.StringValue(srv.AddinProviderName)
	data.AddinButtonLabel = types.StringValue(srv.AddinButtonLabel)
	data.AddinGroupLabel = types.StringValue(srv.AddinGroupLabel)
	data.AddinSupertipTitle = types.StringValue(srv.AddinSupertipTitle)
	data.AddinSupertipDescription = types.StringValue(srv.AddinSupertipDescription)
	data.AddinAuthorityURL = types.StringValue(srv.AddinAuthorityURL)
	data.AddinOidcMetadataURL = types.StringValue(srv.AddinOidcMetadataURL)
	data.AddinAuthenticationMethod = types.StringValue(srv.AddinAuthenticationMethod)
	// PersonalVmrOauthClientSecret is write-only and not returned by the API, will be preserved from state/plan
	data.PersonalVmrOauthAuthEndpoint = types.StringValue(srv.PersonalVmrOauthAuthEndpoint)
	data.PersonalVmrOauthTokenEndpoint = types.StringValue(srv.PersonalVmrOauthTokenEndpoint)
	data.PersonalVmrAdfsRelyingPartyTrustIdentifier = types.StringValue(srv.PersonalVmrAdfsRelyingPartyTrustIdentifier)
	data.OfficeJsURL = types.StringValue(srv.OfficeJsURL)
	data.MicrosoftFabricURL = types.StringValue(srv.MicrosoftFabricURL)
	data.MicrosoftFabricComponentsURL = types.StringValue(srv.MicrosoftFabricComponentsURL)
	data.AdditionalAddInScriptSources = types.StringValue(srv.AdditionalAddInScriptSources)
	data.NonIdpParticipants = types.StringValue(srv.NonIdpParticipants)
	data.PublicKey = types.StringValue(srv.PublicKey)

	// Handle optional pointer fields
	if srv.RoomMailboxEmailAddress != nil {
		data.RoomMailboxEmailAddress = types.StringValue(*srv.RoomMailboxEmailAddress)
	} else {
		data.RoomMailboxEmailAddress = types.StringNull()
	}

	if srv.ScheduledAliasPrefix != nil {
		data.ScheduledAliasPrefix = types.StringValue(*srv.ScheduledAliasPrefix)
	} else {
		data.ScheduledAliasPrefix = types.StringNull()
	}

	if srv.OauthClientID != nil {
		data.OauthClientID = types.StringValue(*srv.OauthClientID)
	} else {
		data.OauthClientID = types.StringNull()
	}

	if srv.OauthState != nil {
		data.OauthState = types.StringValue(*srv.OauthState)
	} else {
		data.OauthState = types.StringNull()
	}

	if srv.AddinApplicationID != nil {
		data.AddinApplicationID = types.StringValue(*srv.AddinApplicationID)
	} else {
		data.AddinApplicationID = types.StringNull()
	}

	if srv.AddinNaaWebApiApplicationID != nil {
		data.AddinNaaWebApiApplicationID = types.StringValue(*srv.AddinNaaWebApiApplicationID)
	} else {
		data.AddinNaaWebApiApplicationID = types.StringNull()
	}

	if srv.PersonalVmrOauthClientID != nil {
		data.PersonalVmrOauthClientID = types.StringValue(*srv.PersonalVmrOauthClientID)
	} else {
		data.PersonalVmrOauthClientID = types.StringNull()
	}

	if srv.Domains != nil {
		// Convert DNS servers from SDK to Terraform format
		var exchangeDomains []string
		for _, domain := range *srv.Domains {
			exchangeDomains = append(exchangeDomains, fmt.Sprintf("/api/admin/configuration/v1/exchange_domain/%d/", domain.ID))
		}
		domainSetValue, diags := types.SetValueFrom(ctx, types.StringType, exchangeDomains)
		if diags.HasError() {
			return nil, fmt.Errorf("error converting DNS servers: %v", diags)
		}
		data.Domains = domainSetValue
	} else {
		data.Domains = types.SetNull(types.StringType)
	}

	if srv.HostIdentityProviderGroup != nil {
		data.HostIdentityProviderGroup = types.StringValue(*srv.HostIdentityProviderGroup)
	} else {
		data.HostIdentityProviderGroup = types.StringNull()
	}

	if srv.IvrTheme != nil {
		data.IvrTheme = types.StringValue(*srv.IvrTheme)
	} else {
		data.IvrTheme = types.StringNull()
	}

	if srv.PrivateKey != nil {
		data.PrivateKey = types.StringValue(*srv.PrivateKey)
	} else {
		data.PrivateKey = types.StringNull()
	}

	// Template fields - these have default values from the API
	data.AcceptEditedOccurrenceTemplate = types.StringValue(srv.AcceptEditedOccurrenceTemplate)
	data.AcceptEditedRecurringSeriesTemplate = types.StringValue(srv.AcceptEditedRecurringSeriesTemplate)
	data.AcceptEditedSingleMeetingTemplate = types.StringValue(srv.AcceptEditedSingleMeetingTemplate)
	data.AcceptNewRecurringSeriesTemplate = types.StringValue(srv.AcceptNewRecurringSeriesTemplate)
	data.AcceptNewSingleMeetingTemplate = types.StringValue(srv.AcceptNewSingleMeetingTemplate)
	data.ConferenceDescriptionTemplate = types.StringValue(srv.ConferenceDescriptionTemplate)
	data.ConferenceNameTemplate = types.StringValue(srv.ConferenceNameTemplate)
	data.ConferenceSubjectTemplate = types.StringValue(srv.ConferenceSubjectTemplate)
	data.MeetingInstructionsTemplate = types.StringValue(srv.MeetingInstructionsTemplate)
	data.PersonalVmrDescriptionTemplate = types.StringValue(srv.PersonalVmrDescriptionTemplate)
	data.PersonalVmrInstructionsTemplate = types.StringValue(srv.PersonalVmrInstructionsTemplate)
	data.PersonalVmrLocationTemplate = types.StringValue(srv.PersonalVmrLocationTemplate)
	data.PersonalVmrNameTemplate = types.StringValue(srv.PersonalVmrNameTemplate)
	data.PlaceholderInstructionsTemplate = types.StringValue(srv.PlaceholderInstructionsTemplate)
	data.RejectAliasConflictTemplate = types.StringValue(srv.RejectAliasConflictTemplate)
	data.RejectAliasDeletedTemplate = types.StringValue(srv.RejectAliasDeletedTemplate)
	data.RejectGeneralErrorTemplate = types.StringValue(srv.RejectGeneralErrorTemplate)
	data.RejectInvalidAliasIDTemplate = types.StringValue(srv.RejectInvalidAliasIDTemplate)
	data.RejectRecurringSeriesPastTemplate = types.StringValue(srv.RejectRecurringSeriesPastTemplate)
	data.RejectSingleMeetingPast = types.StringValue(srv.RejectSingleMeetingPast)
	data.ScheduledAliasDescriptionTemplate = types.StringValue(srv.ScheduledAliasDescriptionTemplate)
	data.AddinPaneAlreadyVideoMeetingHeading = types.StringValue(srv.AddinPaneAlreadyVideoMeetingHeading)
	data.AddinPaneAlreadyVideoMeetingMessage = types.StringValue(srv.AddinPaneAlreadyVideoMeetingMessage)
	data.AddinPaneButtonTitle = types.StringValue(srv.AddinPaneButtonTitle)
	data.AddinPaneDescription = types.StringValue(srv.AddinPaneDescription)
	data.AddinPaneGeneralErrorHeading = types.StringValue(srv.AddinPaneGeneralErrorHeading)
	data.AddinPaneGeneralErrorMessage = types.StringValue(srv.AddinPaneGeneralErrorMessage)
	data.AddinPaneManagementNodeDownHeading = types.StringValue(srv.AddinPaneManagementNodeDownHeading)
	data.AddinPaneManagementNodeDownMessage = types.StringValue(srv.AddinPaneManagementNodeDownMessage)
	data.AddinPanePersonalVmrAddButton = types.StringValue(srv.AddinPanePersonalVmrAddButton)
	data.AddinPanePersonalVmrErrorGettingMessage = types.StringValue(srv.AddinPanePersonalVmrErrorGettingMessage)
	data.AddinPanePersonalVmrErrorInsertingMeetingMessage = types.StringValue(srv.AddinPanePersonalVmrErrorInsertingMeetingMessage)
	data.AddinPanePersonalVmrErrorSigningInMessage = types.StringValue(srv.AddinPanePersonalVmrErrorSigningInMessage)
	data.AddinPanePersonalVmrNoneMessage = types.StringValue(srv.AddinPanePersonalVmrNoneMessage)
	data.AddinPanePersonalVmrSelectMessage = types.StringValue(srv.AddinPanePersonalVmrSelectMessage)
	data.AddinPanePersonalVmrSignInButton = types.StringValue(srv.AddinPanePersonalVmrSignInButton)
	data.AddinPaneSuccessHeading = types.StringValue(srv.AddinPaneSuccessHeading)
	data.AddinPaneSuccessMessage = types.StringValue(srv.AddinPaneSuccessMessage)
	data.AddinPaneTitle = types.StringValue(srv.AddinPaneTitle)

	return &data, nil
}

func (r *InfinityMsExchangeConnectorResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	state := &InfinityMsExchangeConnectorResourceModel{}

	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Preserve write-only sensitive fields from prior state
	priorPassword := state.Password
	priorOauthClientSecret := state.OauthClientSecret
	priorPersonalVmrOauthClientSecret := state.PersonalVmrOauthClientSecret

	resourceID := int(state.ResourceID.ValueInt32())
	state, err := r.read(ctx, resourceID)
	if err != nil {
		// Check if the error is a 404 (not found)
		if isNotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading Infinity Microsoft Exchange connector",
			fmt.Sprintf("Could not read Infinity Microsoft Exchange connector: %s", err),
		)
		return
	}

	// Restore write-only sensitive fields from prior state
	state.Password = priorPassword
	state.OauthClientSecret = priorOauthClientSecret
	state.PersonalVmrOauthClientSecret = priorPersonalVmrOauthClientSecret

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *InfinityMsExchangeConnectorResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	plan := &InfinityMsExchangeConnectorResourceModel{}
	state := &InfinityMsExchangeConnectorResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateRequest := &config.MsExchangeConnectorUpdateRequest{
		Name:                                               plan.Name.ValueString(),
		Description:                                        plan.Description.ValueString(),
		RoomMailboxName:                                    plan.RoomMailboxName.ValueString(),
		URL:                                                plan.URL.ValueString(),
		Username:                                           plan.Username.ValueString(),
		Password:                                           plan.Password.ValueString(),
		AuthenticationMethod:                               plan.AuthenticationMethod.ValueString(),
		AuthProvider:                                       plan.AuthProvider.ValueString(),
		UUID:                                               plan.UUID.ValueString(),
		ScheduledAliasDomain:                               plan.ScheduledAliasDomain.ValueString(),
		OauthClientSecret:                                  plan.OauthClientSecret.ValueString(),
		OauthAuthEndpoint:                                  plan.OauthAuthEndpoint.ValueString(),
		OauthTokenEndpoint:                                 plan.OauthTokenEndpoint.ValueString(),
		OauthRedirectURI:                                   plan.OauthRedirectURI.ValueString(),
		OauthRefreshToken:                                  plan.OauthRefreshToken.ValueString(),
		KerberosRealm:                                      plan.KerberosRealm.ValueString(),
		KerberosKdc:                                        plan.KerberosKdc.ValueString(),
		KerberosKdcHttpsProxy:                              plan.KerberosKdcHttpsProxy.ValueString(),
		KerberosExchangeSpn:                                plan.KerberosExchangeSpn.ValueString(),
		AddinServerDomain:                                  plan.AddinServerDomain.ValueString(),
		AddinDisplayName:                                   plan.AddinDisplayName.ValueString(),
		AddinDescription:                                   plan.AddinDescription.ValueString(),
		AddinProviderName:                                  plan.AddinProviderName.ValueString(),
		AddinButtonLabel:                                   plan.AddinButtonLabel.ValueString(),
		AddinGroupLabel:                                    plan.AddinGroupLabel.ValueString(),
		AddinSupertipTitle:                                 plan.AddinSupertipTitle.ValueString(),
		AddinSupertipDescription:                           plan.AddinSupertipDescription.ValueString(),
		AddinAuthorityURL:                                  plan.AddinAuthorityURL.ValueString(),
		AddinOidcMetadataURL:                               plan.AddinOidcMetadataURL.ValueString(),
		AddinAuthenticationMethod:                          plan.AddinAuthenticationMethod.ValueString(),
		PersonalVmrOauthClientSecret:                       plan.PersonalVmrOauthClientSecret.ValueString(),
		PersonalVmrOauthAuthEndpoint:                       plan.PersonalVmrOauthAuthEndpoint.ValueString(),
		PersonalVmrOauthTokenEndpoint:                      plan.PersonalVmrOauthTokenEndpoint.ValueString(),
		PersonalVmrAdfsRelyingPartyTrustIdentifier:         plan.PersonalVmrAdfsRelyingPartyTrustIdentifier.ValueString(),
		OfficeJsURL:                                        plan.OfficeJsURL.ValueString(),
		MicrosoftFabricURL:                                 plan.MicrosoftFabricURL.ValueString(),
		MicrosoftFabricComponentsURL:                       plan.MicrosoftFabricComponentsURL.ValueString(),
		AdditionalAddInScriptSources:                       plan.AdditionalAddInScriptSources.ValueString(),
		NonIdpParticipants:                                 plan.NonIdpParticipants.ValueString(),
		AcceptEditedOccurrenceTemplate:                     plan.AcceptEditedOccurrenceTemplate.ValueString(),
		AcceptEditedRecurringSeriesTemplate:                plan.AcceptEditedRecurringSeriesTemplate.ValueString(),
		AcceptEditedSingleMeetingTemplate:                  plan.AcceptEditedSingleMeetingTemplate.ValueString(),
		AcceptNewRecurringSeriesTemplate:                   plan.AcceptNewRecurringSeriesTemplate.ValueString(),
		AcceptNewSingleMeetingTemplate:                     plan.AcceptNewSingleMeetingTemplate.ValueString(),
		ConferenceDescriptionTemplate:                      plan.ConferenceDescriptionTemplate.ValueString(),
		ConferenceNameTemplate:                             plan.ConferenceNameTemplate.ValueString(),
		ConferenceSubjectTemplate:                          plan.ConferenceSubjectTemplate.ValueString(),
		MeetingInstructionsTemplate:                        plan.MeetingInstructionsTemplate.ValueString(),
		PersonalVmrDescriptionTemplate:                     plan.PersonalVmrDescriptionTemplate.ValueString(),
		PersonalVmrInstructionsTemplate:                    plan.PersonalVmrInstructionsTemplate.ValueString(),
		PersonalVmrLocationTemplate:                        plan.PersonalVmrLocationTemplate.ValueString(),
		PersonalVmrNameTemplate:                            plan.PersonalVmrNameTemplate.ValueString(),
		PlaceholderInstructionsTemplate:                    plan.PlaceholderInstructionsTemplate.ValueString(),
		RejectAliasConflictTemplate:                        plan.RejectAliasConflictTemplate.ValueString(),
		RejectAliasDeletedTemplate:                         plan.RejectAliasDeletedTemplate.ValueString(),
		RejectGeneralErrorTemplate:                         plan.RejectGeneralErrorTemplate.ValueString(),
		RejectInvalidAliasIDTemplate:                       plan.RejectInvalidAliasIDTemplate.ValueString(),
		RejectRecurringSeriesPastTemplate:                  plan.RejectRecurringSeriesPastTemplate.ValueString(),
		RejectSingleMeetingPast:                            plan.RejectSingleMeetingPast.ValueString(),
		ScheduledAliasDescriptionTemplate:                  plan.ScheduledAliasDescriptionTemplate.ValueString(),
		AddinPaneAlreadyVideoMeetingHeading:                plan.AddinPaneAlreadyVideoMeetingHeading.ValueString(),
		AddinPaneAlreadyVideoMeetingMessage:                plan.AddinPaneAlreadyVideoMeetingMessage.ValueString(),
		AddinPaneButtonTitle:                               plan.AddinPaneButtonTitle.ValueString(),
		AddinPaneDescription:                               plan.AddinPaneDescription.ValueString(),
		AddinPaneGeneralErrorHeading:                       plan.AddinPaneGeneralErrorHeading.ValueString(),
		AddinPaneGeneralErrorMessage:                       plan.AddinPaneGeneralErrorMessage.ValueString(),
		AddinPaneManagementNodeDownHeading:                 plan.AddinPaneManagementNodeDownHeading.ValueString(),
		AddinPaneManagementNodeDownMessage:                 plan.AddinPaneManagementNodeDownMessage.ValueString(),
		AddinPanePersonalVmrAddButton:                      plan.AddinPanePersonalVmrAddButton.ValueString(),
		AddinPanePersonalVmrErrorGettingMessage:            plan.AddinPanePersonalVmrErrorGettingMessage.ValueString(),
		AddinPanePersonalVmrErrorInsertingMeetingMessage:   plan.AddinPanePersonalVmrErrorInsertingMeetingMessage.ValueString(),
		AddinPanePersonalVmrErrorSigningInMessage:          plan.AddinPanePersonalVmrErrorSigningInMessage.ValueString(),
		AddinPanePersonalVmrNoneMessage:                    plan.AddinPanePersonalVmrNoneMessage.ValueString(),
		AddinPanePersonalVmrSelectMessage:                  plan.AddinPanePersonalVmrSelectMessage.ValueString(),
		AddinPanePersonalVmrSignInButton:                   plan.AddinPanePersonalVmrSignInButton.ValueString(),
		AddinPaneSuccessHeading:                            plan.AddinPaneSuccessHeading.ValueString(),
		AddinPaneSuccessMessage:                            plan.AddinPaneSuccessMessage.ValueString(),
		AddinPaneTitle:                                     plan.AddinPaneTitle.ValueString(),
	}

	// Handle optional pointer fields for update
	if !plan.RoomMailboxEmailAddress.IsNull() && !plan.RoomMailboxEmailAddress.IsUnknown() {
		email := plan.RoomMailboxEmailAddress.ValueString()
		updateRequest.RoomMailboxEmailAddress = &email
	}

	if !plan.ScheduledAliasPrefix.IsNull() && !plan.ScheduledAliasPrefix.IsUnknown() {
		prefix := plan.ScheduledAliasPrefix.ValueString()
		updateRequest.ScheduledAliasPrefix = &prefix
	}

	// Handle optional pointer fields for integers and booleans
	if !plan.ScheduledAliasSuffixLength.IsNull() && !plan.ScheduledAliasSuffixLength.IsUnknown() {
		suffixLength := int(plan.ScheduledAliasSuffixLength.ValueInt64())
		updateRequest.ScheduledAliasSuffixLength = &suffixLength
	}

	if !plan.MeetingBufferBefore.IsNull() && !plan.MeetingBufferBefore.IsUnknown() {
		bufferBefore := int(plan.MeetingBufferBefore.ValueInt64())
		updateRequest.MeetingBufferBefore = &bufferBefore
	}

	if !plan.MeetingBufferAfter.IsNull() && !plan.MeetingBufferAfter.IsUnknown() {
		bufferAfter := int(plan.MeetingBufferAfter.ValueInt64())
		updateRequest.MeetingBufferAfter = &bufferAfter
	}

	if !plan.EnableDynamicVmrs.IsNull() && !plan.EnableDynamicVmrs.IsUnknown() {
		enableDynamic := plan.EnableDynamicVmrs.ValueBool()
		updateRequest.EnableDynamicVmrs = &enableDynamic
	}

	if !plan.EnablePersonalVmrs.IsNull() && !plan.EnablePersonalVmrs.IsUnknown() {
		enablePersonal := plan.EnablePersonalVmrs.ValueBool()
		updateRequest.EnablePersonalVmrs = &enablePersonal
	}

	if !plan.AllowNewUsers.IsNull() && !plan.AllowNewUsers.IsUnknown() {
		allowNew := plan.AllowNewUsers.ValueBool()
		updateRequest.AllowNewUsers = &allowNew
	}

	if !plan.DisableProxy.IsNull() && !plan.DisableProxy.IsUnknown() {
		disableProxy := plan.DisableProxy.ValueBool()
		updateRequest.DisableProxy = &disableProxy
	}

	if !plan.UseCustomAddInSources.IsNull() && !plan.UseCustomAddInSources.IsUnknown() {
		useCustom := plan.UseCustomAddInSources.ValueBool()
		updateRequest.UseCustomAddInSources = &useCustom
	}

	if !plan.EnableAddinDebugLogs.IsNull() && !plan.EnableAddinDebugLogs.IsUnknown() {
		enableDebug := plan.EnableAddinDebugLogs.ValueBool()
		updateRequest.EnableAddinDebugLogs = &enableDebug
	}

	if !plan.OauthClientID.IsNull() && !plan.OauthClientID.IsUnknown() {
		clientID := plan.OauthClientID.ValueString()
		updateRequest.OauthClientID = &clientID
	}

	if !plan.KerberosEnableTls.IsNull() && !plan.KerberosEnableTls.IsUnknown() {
		enableTls := plan.KerberosEnableTls.ValueBool()
		updateRequest.KerberosEnableTls = &enableTls
	}

	if !plan.KerberosAuthEveryRequest.IsNull() && !plan.KerberosAuthEveryRequest.IsUnknown() {
		authEvery := plan.KerberosAuthEveryRequest.ValueBool()
		updateRequest.KerberosAuthEveryRequest = &authEvery
	}

	if !plan.KerberosVerifyTlsUsingCustomCa.IsNull() && !plan.KerberosVerifyTlsUsingCustomCa.IsUnknown() {
		verifyTls := plan.KerberosVerifyTlsUsingCustomCa.ValueBool()
		updateRequest.KerberosVerifyTlsUsingCustomCa = &verifyTls
	}

	if !plan.AddinApplicationID.IsNull() && !plan.AddinApplicationID.IsUnknown() {
		appID := plan.AddinApplicationID.ValueString()
		updateRequest.AddinApplicationID = &appID
	}

	if !plan.AddinNaaWebApiApplicationID.IsNull() && !plan.AddinNaaWebApiApplicationID.IsUnknown() {
		apiAppID := plan.AddinNaaWebApiApplicationID.ValueString()
		updateRequest.AddinNaaWebApiApplicationID = &apiAppID
	}

	if !plan.PersonalVmrOauthClientID.IsNull() && !plan.PersonalVmrOauthClientID.IsUnknown() {
		vmrClientID := plan.PersonalVmrOauthClientID.ValueString()
		updateRequest.PersonalVmrOauthClientID = &vmrClientID
	}

	if !plan.Domains.IsNull() && !plan.Domains.IsUnknown() {
		// Convert List attributes to []string
		exchangeDomains, diags := getStringList(ctx, plan.Domains)
		resp.Diagnostics.Append(diags...)
		updateRequest.Domains = &exchangeDomains
	}

	if !plan.HostIdentityProviderGroup.IsNull() && !plan.HostIdentityProviderGroup.IsUnknown() {
		hostIdp := plan.HostIdentityProviderGroup.ValueString()
		updateRequest.HostIdentityProviderGroup = &hostIdp
	}

	if !plan.IvrTheme.IsNull() && !plan.IvrTheme.IsUnknown() {
		theme := plan.IvrTheme.ValueString()
		updateRequest.IvrTheme = &theme
	}

	resourceID := int(state.ResourceID.ValueInt32())
	_, err := r.InfinityClient.Config().UpdateMsExchangeConnector(ctx, resourceID, updateRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Infinity Microsoft Exchange connector",
			fmt.Sprintf("Could not update Infinity Microsoft Exchange connector: %s", err),
		)
		return
	}

	// Read the state from the API to get all computed values
	model, err := r.read(ctx, resourceID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Updated Infinity Microsoft Exchange connector",
			fmt.Sprintf("Could not read updated Infinity Microsoft Exchange connector with ID %d: %s", resourceID, err),
		)
		return
	}

	// Preserve write-only sensitive fields from plan
	model.Password = plan.Password
	model.OauthClientSecret = plan.OauthClientSecret
	model.PersonalVmrOauthClientSecret = plan.PersonalVmrOauthClientSecret

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *InfinityMsExchangeConnectorResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	state := &InfinityMsExchangeConnectorResourceModel{}

	tflog.Info(ctx, "Deleting Infinity Microsoft Exchange connector")

	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.InfinityClient.Config().DeleteMsExchangeConnector(ctx, int(state.ResourceID.ValueInt32()))

	// Ignore 404 Not Found and Lookup errors on delete
	if err != nil && !isNotFoundError(err) && !isLookupError(err) {
		resp.Diagnostics.AddError(
			"Error Deleting Infinity Microsoft Exchange connector",
			fmt.Sprintf("Could not delete Infinity Microsoft Exchange connector with ID %s: %s", state.ID.ValueString(), err),
		)
		return
	}
}

func (r *InfinityMsExchangeConnectorResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resourceID, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Import ID must be a valid integer for the resource ID. Got: %s", req.ID),
		)
		return
	}

	tflog.Trace(ctx, fmt.Sprintf("Importing Infinity Microsoft Exchange connector with resource ID: %d", resourceID))

	// Read the resource from the API
	model, err := r.read(ctx, resourceID)
	if err != nil {
		// Check if the error is a 404 (not found)
		if isNotFoundError(err) {
			resp.Diagnostics.AddError(
				"Infinity Microsoft Exchange Connector Not Found",
				fmt.Sprintf("Infinity Microsoft Exchange connector with resource ID %d not found.", resourceID),
			)
			return
		}
		resp.Diagnostics.AddError(
			"Error Importing Infinity Microsoft Exchange Connector",
			fmt.Sprintf("Could not import Infinity Microsoft Exchange connector with resource ID %d: %s", resourceID, err),
		)
		return
	}

	// Set the state from the imported resource
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}
