/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/pexip/go-infinity-sdk/v38/config"
)

var (
	_ resource.ResourceWithImportState    = (*InfinityGlobalConfigurationResource)(nil)
	_ resource.ResourceWithValidateConfig = (*InfinityGlobalConfigurationResource)(nil)
)

type InfinityGlobalConfigurationResource struct {
	InfinityClient InfinityClient
}

/*
The following fields have been deprecated in the Infinity API and have been removed from the model:
  * default_to_new_webapp - use 'default_webapp' instead
  * default_webapp - use 'default_webapp_alias' instead
  * enable_multiscreen (Deprecated field)
  * enable_push_notifications (Deprecated field)
  * legacy_api_http (Deprecated field)
  * live_captions_api_gateway - use 'media_processing_server' resource
  * live_captions_app_id - use 'media_processing_server' resource
  * live_captions_enabled - this field is deprecated and will be ignored
  * live_captions_public_jwt_key - use 'media_processing_server' resource
  * es_maximum_deferred_posts (Deprecated field)
  * pss_customer_id (Deprecated field)
  * pss_enabled (Deprecated field)
  * pss_gateway (Deprecated field)
  * pss_token (Deprecated field)
*/

type InfinityGlobalConfigurationResourceModel struct {
	ID                                  types.String `tfsdk:"id"`
	AWSAccessKey                        types.String `tfsdk:"aws_access_key"`
	AWSSecretKey                        types.String `tfsdk:"aws_secret_key"`
	AzureClientID                       types.String `tfsdk:"azure_client_id"`
	AzureSecret                         types.String `tfsdk:"azure_secret"`
	AzureSubscriptionID                 types.String `tfsdk:"azure_subscription_id"`
	AzureTenant                         types.String `tfsdk:"azure_tenant"`
	BdpmMaxPinFailuresPerWindow         types.Int64  `tfsdk:"bdpm_max_pin_failures_per_window"`
	BdpmMaxScanAttemptsPerWindow        types.Int64  `tfsdk:"bdpm_max_scan_attempts_per_window"`
	BdpmPinChecksEnabled                types.Bool   `tfsdk:"bdpm_pin_checks_enabled"`
	BdpmScanQuarantineEnabled           types.Bool   `tfsdk:"bdpm_scan_quarantine_enabled"`
	BurstingEnabled                     types.Bool   `tfsdk:"bursting_enabled"`
	BurstingMinLifetime                 types.Int64  `tfsdk:"bursting_min_lifetime"`
	BurstingThreshold                   types.Int64  `tfsdk:"bursting_threshold"`
	CloudProvider                       types.String `tfsdk:"cloud_provider"`
	ContactEmailAddress                 types.String `tfsdk:"contact_email_address"`
	ContentSecurityPolicyHeader         types.String `tfsdk:"content_security_policy_header"`
	ContentSecurityPolicyState          types.Bool   `tfsdk:"content_security_policy_state"`
	CryptoMode                          types.String `tfsdk:"crypto_mode"`
	DefaultTheme                        types.String `tfsdk:"default_theme"`
	DefaultWebappAlias                  types.String `tfsdk:"default_webapp_alias"`
	DeploymentUUID                      types.String `tfsdk:"deployment_uuid"`
	DisabledCodecs                      types.Set    `tfsdk:"disabled_codecs"`
	EjectLastParticipantBackstopTimeout types.Int64  `tfsdk:"eject_last_participant_backstop_timeout"`
	EnableAnalytics                     types.Bool   `tfsdk:"enable_analytics"`
	EnableApplicationAPI                types.Bool   `tfsdk:"enable_application_api"`
	EnableBreakoutRooms                 types.Bool   `tfsdk:"enable_breakout_rooms"`
	EnableChat                          types.Bool   `tfsdk:"enable_chat"`
	EnableClock                         types.Bool   `tfsdk:"enable_clock"`
	EnableClock                         types.Bool   `tfsdk:"enable_clock"`
	EnableDenoise                       types.Bool   `tfsdk:"enable_denoise"`
	EnableDialout                       types.Bool   `tfsdk:"enable_dialout"`
	EnableDirectory                     types.Bool   `tfsdk:"enable_directory"`
	EnableEdgeNonMesh                   types.Bool   `tfsdk:"enable_edge_non_mesh"`
	EnableFecc                          types.Bool   `tfsdk:"enable_fecc"`
	EnableH323                          types.Bool   `tfsdk:"enable_h323"`
	EnableLegacyDialoutAPI              types.Bool   `tfsdk:"enable_legacy_dialout_api"`
	EnableLyncAutoEscalate              types.Bool   `tfsdk:"enable_lync_auto_escalate"`
	EnableLyncVbss                      types.Bool   `tfsdk:"enable_lync_vbss"`
	EnableMlvad                         types.Bool   `tfsdk:"enable_mlvad"`
	EnableRTMP                          types.Bool   `tfsdk:"enable_rtmp"`
	EnableSIP                           types.Bool   `tfsdk:"enable_sip"`
	EnableSIPUDP                        types.Bool   `tfsdk:"enable_sip_udp"`
	EnableSoftmute                      types.Bool   `tfsdk:"enable_softmute"`
	EnableSSH                           types.Bool   `tfsdk:"enable_ssh"`
	EnableTurn443                       types.Bool   `tfsdk:"enable_turn_443"`
	EnableWebRTC                        types.Bool   `tfsdk:"enable_webrtc"`
	ErrorReportingEnabled               types.Bool   `tfsdk:"error_reporting_enabled"`
	ErrorReportingURL                   types.String `tfsdk:"error_reporting_url"`
	EsConnectionTimeout                 types.Int64  `tfsdk:"es_connection_timeout"`
	EsInitialRetryBackoff               types.Int64  `tfsdk:"es_initial_retry_backoff"`
	EsMaximumRetryBackoff               types.Int64  `tfsdk:"es_maximum_retry_backoff"`
	EsMediaStreamsWait                  types.Int64  `tfsdk:"es_media_streams_wait"`
	EsMetricsUpdateInterval             types.Int64  `tfsdk:"es_metrics_update_interval"`
	EsShortTermMemoryExpiration         types.Int64  `tfsdk:"es_short_term_memory_expiration"`
	ExternalParticipantAvatarLookup     types.Bool   `tfsdk:"external_participant_avatar_lookup"`
	GcpClientEmail                      types.String `tfsdk:"gcp_client_email"`
	GcpPrivateKey                       types.String `tfsdk:"gcp_private_key"`
	GcpProjectID                        types.String `tfsdk:"gcp_project_id"`
	GuestsOnlyTimeout                   types.Int64  `tfsdk:"guests_only_timeout"`
	LegacyAPIUsername                   types.String `tfsdk:"legacy_api_username"`
	LegacyAPIPassword                   types.String `tfsdk:"legacy_api_password"`
	LiveCaptionsVMRDefault              types.Bool   `tfsdk:"live_captions_vmr_default"`
	LiveviewShowConferences             types.Bool   `tfsdk:"liveview_show_conferences"`
	LocalMssipDomain                    types.String `tfsdk:"local_mssip_domain"`
	LogonBanner                         types.String `tfsdk:"logon_banner"`
	LogsMaxAge                          types.Int64  `tfsdk:"logs_max_age"`
	ManagementQos                       types.Int64  `tfsdk:"management_qos"`
	ManagementSessionTimeout            types.Int64  `tfsdk:"management_session_timeout"`
	ManagementStartPage                 types.String `tfsdk:"management_start_page"`
	MaxCallrateIn                       types.Int64  `tfsdk:"max_callrate_in"`
	MaxCallrateOut                      types.Int64  `tfsdk:"max_callrate_out"`
	MaxPixelsPerSecond                  types.String `tfsdk:"max_pixels_per_second"`
	MaxPresentationBandwidthRatio       types.Int64  `tfsdk:"max_presentation_bandwidth_ratio"`
	MediaPortsEnd                       types.Int64  `tfsdk:"media_ports_end"`
	MediaPortsStart                     types.Int64  `tfsdk:"media_ports_start"`
	OcspResponderURL                    types.String `tfsdk:"ocsp_responder_url"`
	OcspState                           types.String `tfsdk:"ocsp_state"`
	PinEntryTimeout                     types.Int64  `tfsdk:"pin_entry_timeout"`
	ResourceURI                         types.String `tfsdk:"resource_uri"`
	SessionTimeoutEnabled               types.Bool   `tfsdk:"session_timeout_enabled"`
	SignallingPortsEnd                  types.Int64  `tfsdk:"signalling_ports_end"`
	SignallingPortsStart                types.Int64  `tfsdk:"signalling_ports_start"`
	SipTLSCertVerifyMode                types.String `tfsdk:"sip_tls_cert_verify_mode"`
	SiteBanner                          types.String `tfsdk:"site_banner"`
	SiteBannerBg                        types.String `tfsdk:"site_banner_bg"`
	SiteBannerFg                        types.String `tfsdk:"site_banner_fg"`
	TeamsEnablePowerpointRender         types.Bool   `tfsdk:"teams_enable_powerpoint_render"`
	WaitingForChairTimeout              types.Int64  `tfsdk:"waiting_for_chair_timeout"`
}

func (r *InfinityGlobalConfigurationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_infinity_global_configuration"
}

func (r *InfinityGlobalConfigurationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *InfinityGlobalConfigurationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Resource URI for the global configuration in Infinity",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"aws_access_key": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(40),
				},
				MarkdownDescription: "The Amazon Web Services access key ID for the AWS user that the Pexip Infinity Management Node will use to log in to AWS and start and stop the node instances. Maximum length: 40 characters.",
			},
			"aws_secret_key": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				MarkdownDescription: "The Amazon Web Services secret access key that is associated with the AWS access key ID.",
			},
			"azure_client_id": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				MarkdownDescription: "The ID used to identify the client (sometimes referred to as Application ID).",
			},
			"azure_secret": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				MarkdownDescription: "The Azure secret key that is associated with the Azure client ID.",
			},
			"azure_subscription_id": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The ID of an Azure subscription.",
			},
			"azure_tenant": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The Azure tenant ID that is associated with the Azure client ID.",
			},
			"bdpm_max_pin_failures_per_window": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				Default:  int64default.StaticInt64(20),
				Validators: []validator.Int64{
					int64validator.Between(5, 200),
				},
				MarkdownDescription: "Sets the maximum number of PIN failures per service (e.g. VMR) in any sliding 10 minute windowed period, that are allowed from participants at unknown source addresses, before protective action is taken for that service. Range: 5 to 200.",
			},
			"bdpm_max_scan_attempts_per_window": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				Default:  int64default.StaticInt64(20),
				Validators: []validator.Int64{
					int64validator.Between(5, 200),
				},
				MarkdownDescription: "Sets the maximum number of incorrect alias dial attempts in any sliding 10-minute windowed period, that are allowed from an unknown source address, before protective action is taken against that address. Range: 5 to 200.",
			},
			"bdpm_pin_checks_enabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
				MarkdownDescription: "Select this option to instruct Pexip Infinity's Break-in Defense Policy Manager to temporarily block all access to a VMR that receives a significant number of incorrect PIN entry attempts (and thus may perhaps be under attack from a malicious actor). By default, this will block ALL new access attempts to a VMR for up to 10 minutes if more than 20 incorrect PIN entry attempts are made against that VMR in a 10 minute window.",
			},
			"bdpm_scan_quarantine_enabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
				MarkdownDescription: "Select this option to instruct Pexip Infinity's Break-in Defense Policy Manager to temporarily block service access attempts from any source IP address that dials a significant number of incorrect aliases in a short period (and thus may perhaps be attempting to scan your deployment to discover valid aliases).",
			},
			"bursting_enabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Select this option to instruct Pexip Infinity to monitor the system locations and start up / shutdown overflow Conferencing Nodes hosted in either Amazon Web Services (AWS) or Microsoft Azure when in need of extra capacity. For more information, see the Admin Guide section 'Dynamic bursting to a cloud service'.",
			},
			"bursting_min_lifetime": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				Default:  int64default.StaticInt64(50),
				Validators: []validator.Int64{
					int64validator.AtLeast(5),
				},
				MarkdownDescription: "The minimum number of minutes that a cloud bursting node is kept powered on. Note that newly started cloud Conferencing Nodes can take up to 5 minutes to fully startup. Minimum: 5.",
			},
			"bursting_threshold": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				Default:  int64default.StaticInt64(5),
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
				},
				MarkdownDescription: "The bursting threshold controls when your overflow Conferencing Nodes in the cloud are automatically started up so that they can provide additional conferencing capacity. Minimum: 1.",
			},
			"cloud_provider": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("AWS"),
				Validators: []validator.String{
					stringvalidator.OneOf("AWS", "AZURE", "GCP"),
				},
				MarkdownDescription: "Choose the cloud service provider to use for bursting. Valid values: `AWS`, `AZURE`, `GCP`.",
			},
			"contact_email_address": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(""),
				Validators: []validator.String{
					stringvalidator.LengthAtMost(100),
				},
				MarkdownDescription: "An email address to be added to incident reports to allow Pexip to contact the system administrator for further information. Maximum length: 100 characters.",
			},
			"content_security_policy_header": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("upgrade-insecure-requests; default-src 'self'; frame-ancestors 'self'; frame-src 'self' https://telemetryservice.firstpartyapps.oaspapps.com/telemetryservice/telemetryproxy.html https://*.microsoft.com https://*.office.com; style-src 'self' 'unsafe-inline' https://*.microsoft.com https://*.office.com; object-src 'self'; font-src 'self' https://*.microsoft.com https://*.office.com; img-src 'self' https://www.adobe.com data: blob:; script-src 'self' 'unsafe-inline' 'unsafe-eval' https://*.microsoft.com https://*.office.com https://ajax.aspnetcdn.com https://api.keen.io; media-src 'self' blob:; connect-src 'self' https://*.microsoft.com https://*.office.com https://example.com;"),
				Validators: []validator.String{
					stringvalidator.LengthAtMost(4096),
				},
				MarkdownDescription: "HTTP Content-Security-Policy header contents for Conferencing Nodes. Maximum length: 4096 characters.",
			},
			"content_security_policy_state": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
				MarkdownDescription: "Enable HTTP Content-Security-Policy for Conferencing Nodes.",
			},
			"crypto_mode": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("besteffort"),
				Validators: []validator.String{
					stringvalidator.OneOf("besteffort", "on", "off"),
				},
				MarkdownDescription: "Controls the media encryption requirements for participants connecting to Pexip Infinity services. `on`: All participants must use media encryption. `besteffort`: Each participant will use media encryption if their device supports it. `off`: All H.323, SIP and MS-SIP participants must use unencrypted media. You can override this global setting for each individual service.",
			},
			"default_theme": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The theme to use for services that have no specific theme selected.",
			},
			"default_webapp_alias": schema.StringAttribute{
				Optional: true,
				Computed: true,
				// The Infinity API schema erroneously shows the default as null
				Default:             stringdefault.StaticString("/api/admin/configuration/v1/webapp_alias/3/"),
				MarkdownDescription: "The web app path to use by default on conferencing nodes.",
			},
			// unique for each deployment, not update by users
			"deployment_uuid": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The ID of the deployment.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"disabled_codecs": schema.SetAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				Default: setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{
					types.StringValue("MP4A-LATM_128"),
					types.StringValue("H264_H_0"),
					types.StringValue("H264_H_1"),
				})),
				MarkdownDescription: "Choose codecs to disable.",
			},
			"eject_last_participant_backstop_timeout": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				Default:  int64default.StaticInt64(0),
				Validators: []validator.Int64{
					int64validator.Any(
						int64validator.OneOf(0),
						int64validator.Between(60, 86400),
					),
				},
				MarkdownDescription: "The length of time (in seconds) for which a conference will continue with only one participant remaining (independent of Host/Guest role). Must be 0 (never eject) or between 60 and 86400. Default: 0.",
			},
			"enable_analytics": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Select this option to allow submission of deployment and usage statistics to Pexip. This will help us improve the product.",
			},
			"enable_application_api": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
				MarkdownDescription: "Enable or disable support for Pexip Infinity Client API. This is required for integration with Infinity Connect browser-based and desktop clients, the Pexip Mobile App for iOS and Android, and any other third-party applications that use the client API, as well as for integration with Microsoft Teams.",
			},
			"enable_breakout_rooms": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Enable the Breakout Rooms feature on VMRs.",
			},
			"enable_chat": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
				MarkdownDescription: "Enables relay of chat messages between conference participants using Skype for Business and Infinity Connect clients. You can also configure this setting on individual Virtual Meeting Rooms and Virtual Auditoriums.",
			},
			"enable_clock": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Enables support for displaying an in-conference timer or countdown clock.",
			},
			"enable_denoise": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
				MarkdownDescription: "Enable server side denoising for speech from noisy participants (see documentation for ways to enable it for a VMR).",
			},
			"enable_dialout": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
				MarkdownDescription: "Enables calls via the Distributed Gateway, and allows users of Pexip Infinity Connect, the Pexip Mobile Apps for iOS and Android, and the Pexip management web interface to add participants to a conference.",
			},
			"enable_directory": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
				MarkdownDescription: "When disabled, Infinity Connect clients will display aliases from their own call history only. When enabled, registered Infinity Connect clients will additionally display the aliases of VMRs, Virtual Auditoriums, Virtual Receptions, and devices registered to the Pexip Infinity deployment.",
			},
			"enable_edge_non_mesh": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
				MarkdownDescription: "Enable the restricted IPsec network routing requirements of Proxying Edge Nodes. When enabled, if a location only contains Proxying Edge Nodes, then those nodes only require IPsec connectivity with other nodes in that location, the transcoding location, the primary and secondary overflow locations, and with the Management Node.",
			},
			"enable_fecc": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
				MarkdownDescription: "Enables Connect apps and SIP/H.323 endpoints to send Far-End Camera Control (FECC) signals to supporting endpoints, in order to pan, tilt and zoom the device's camera.",
			},
			"enable_h323": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
				MarkdownDescription: "Enable the H323 protocol on all Conferencing Nodes.",
			},
			"enable_legacy_dialout_api": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Enables outbound calls from a VMR using the legacy dialout API. When disabled, outbound calls are only permitted by following Call Routing Rules.",
			},
			"enable_lync_auto_escalate": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Determines whether a Skype for Business audio call is automatically escalated so that it receives video from a conference.",
			},
			"enable_lync_vbss": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Determines whether Video-based Screen Sharing (VbSS) is enabled for Skype for Business calls.",
			},
			"enable_mlvad": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Enable Voice Focus for advanced voice activity detection.",
			},
			"enable_rtmp": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
				MarkdownDescription: "Enables RTMP calls on all Conferencing Nodes. This allows Infinity Connect clients that use RTMP to access Pexip Infinity services, and allows conference content to be output to streaming and recording services.",
			},
			"enable_sip": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
				MarkdownDescription: "Enable the SIP protocol over TCP and TLS on all Conferencing Nodes.",
			},
			"enable_sip_udp": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Enable incoming calls using the SIP protocol over UDP on all Conferencing Nodes. If changing from enabled to disabled, all Conferencing Nodes must be rebooted.",
			},
			"enable_softmute": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
				MarkdownDescription: "Enable Softmute for advance speech-aware audio gating (see documentation for ways to enable it for a VMR). Note that this does not remove any noise from the audio.",
			},
			"enable_ssh": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
				MarkdownDescription: "Allows an administrator to log in to the Management and Conferencing Nodes over SSH. This setting can be overridden on individual nodes.",
			},
			"enable_turn_443": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Enable media relay on TCP port 443 for WebRTC clients as a fallback.",
			},
			"enable_webrtc": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
				MarkdownDescription: "Enables WebRTC calls on all Conferencing Nodes. This allows access to Pexip Infinity services from Infinity Connect clients that use WebRTC, including Google Chrome, Microsoft Edge, Firefox, Opera and Safari (version 11 onwards) browsers, and the Infinity Connect desktop client.",
			},
			"error_reporting_enabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Select this option to permit submission of incident reports.",
			},
			"error_reporting_url": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("https://acr.pexip.com"),
				Validators: []validator.String{
					stringvalidator.LengthAtMost(255),
				},
				MarkdownDescription: "The URL to which incident reports will be sent. Maximum length: 255 characters.",
			},
			"es_connection_timeout": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(7),
				MarkdownDescription: "Maximum number of seconds allowed to connect, send, and wait for a response.",
			},
			"es_initial_retry_backoff": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(1),
				MarkdownDescription: "Initial time, in seconds, for the first retry attempt when an event cannot be delivered.",
			},
			"es_maximum_retry_backoff": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(1800),
				MarkdownDescription: "Maximum number of seconds allowed for the retry backoff before raising an alarm and stopping the event publisher.",
			},
			"es_media_streams_wait": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(1),
				MarkdownDescription: "Maximum time, in seconds, to wait for an end-of-call media stream message.",
			},
			"es_metrics_update_interval": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(60),
				MarkdownDescription: "Time between metrics updates. To disable eventsink metrics, enter 0.",
			},
			"es_short_term_memory_expiration": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(2),
				MarkdownDescription: "Internal cache expiration time in seconds. Used to briefly store 'participant_disconnected' events in order to gather end-of-call media statistics.",
			},
			"external_participant_avatar_lookup": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
				MarkdownDescription: "Determines whether or not avatars for external participants will be retrieved using the method appropriate for the external meeting type.",
			},
			"gcp_client_email": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The GCP service account ID.",
			},
			"gcp_private_key": schema.StringAttribute{
				Optional:  true,
				Computed:  true,
				Sensitive: false,
				Default:   stringdefault.StaticString(""),
				Validators: []validator.String{
					stringvalidator.LengthAtMost(12288),
				},
				MarkdownDescription: "The private key for the Google Cloud Platform service account user that the Pexip Infinity Management Node will use to log in to GCP and start and stop the node instances. Maximum length: 12288 characters.",
			},
			"gcp_project_id": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The ID of the GCP project containing bursting nodes.",
			},
			"guests_only_timeout": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				Default:  int64default.StaticInt64(60),
				Validators: []validator.Int64{
					int64validator.Between(0, 86400),
				},
				MarkdownDescription: "The length of time (in seconds) for which a conference will continue with only Guest participants, after all Host participants have left. Range: 0 to 86400. Default: 60.",
			},
			"legacy_api_username": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(""),
				Validators: []validator.String{
					stringvalidator.LengthAtMost(100),
				},
				MarkdownDescription: "The username presented to Pexip Infinity by external systems attempting to authenticate with it. Maximum length: 100 characters.",
			},
			"legacy_api_password": schema.StringAttribute{
				Optional: true,
				//Sensitive:           true,
				Computed: true,
				Default:  stringdefault.StaticString(""),
				Validators: []validator.String{
					stringvalidator.LengthAtMost(100),
				},
				MarkdownDescription: "The password presented to Pexip Infinity by external systems attempting to authenticate with it. Maximum length: 100 characters.",
			},
			"live_captions_vmr_default": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "This option controls whether live captions are enabled by default on all VMRs, Virtual Auditoriums and Call Routing Rules. You can override this setting on each service individually.",
			},
			"liveview_show_conferences": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
				MarkdownDescription: "Whether to show conferences and backplanes in Live View.",
			},
			"local_mssip_domain": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(""),
				Validators: []validator.String{
					stringvalidator.LengthAtMost(255),
				},
				MarkdownDescription: "The name of the SIP domain that is routed from Skype for Business to Pexip Infinity, either as a static route or via federation. It is also used as the default domain in the From address for outgoing SIP gateway calls and outbound SIP calls from conferences without a valid SIP URI as an alias. Maximum length: 255 characters.",
			},
			"logon_banner": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(""),
				Validators: []validator.String{
					stringvalidator.LengthAtMost(4096),
				},
				MarkdownDescription: "Text of the message to display on the login page of the Pexip Infinity administrator web interface. Maximum length: 4096 characters.",
			},
			"logs_max_age": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				Default:  int64default.StaticInt64(0),
				Validators: []validator.Int64{
					int64validator.Between(0, 3650),
				},
				MarkdownDescription: "The maximum number of days of logs and call history to retain on Pexip nodes. 0 to disable. Range: 0 to 3650 days. Default: 0 (disabled).",
			},
			"management_qos": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				Default:  int64default.StaticInt64(0),
				Validators: []validator.Int64{
					int64validator.Between(0, 63),
				},
				MarkdownDescription: "The DSCP value for management traffic sent from the Management Node and Conferencing Nodes. Range: 0 to 63.",
			},
			"management_session_timeout": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				Default:  int64default.StaticInt64(30),
				Validators: []validator.Int64{
					int64validator.Between(5, 1440),
				},
				MarkdownDescription: "The number of minutes a browser session may remain idle before the user is logged out of the Management Node administration interface. Range: 5 to 1440. Default: 30 minutes.",
			},
			"management_start_page": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("/admin/conferencingstatus/deploymentgraph/deployment_graph/"),
				MarkdownDescription: "The first page you see after logging into the Management Web.",
			},
			"max_callrate_in": schema.Int64Attribute{
				Optional: true,
				Validators: []validator.Int64{
					int64validator.Between(128, 8192),
				},
				MarkdownDescription: "This optional field allows you to limit the bandwidth of media being received by Pexip Infinity from individual participants, for calls where bandwidth limits have not otherwise been specified. Range: 128 to 8192.",
			},
			"max_callrate_out": schema.Int64Attribute{
				Optional: true,
				Validators: []validator.Int64{
					int64validator.Between(128, 8192),
				},
				MarkdownDescription: "This optional field allows you to limit the bandwidth of media being sent by Pexip Infinity to individual participants, for calls where bandwidth limits have not otherwise been specified. Range: 128 to 8192. Default: 4128.",
			},
			"max_pixels_per_second": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("hd"),
				Validators: []validator.String{
					stringvalidator.OneOf("sd", "hd", "fullhd"),
				},
				MarkdownDescription: "Sets the maximum call quality for participants connecting to Pexip Infinity services (VMRs, gateway calls etc.). You can also override this setting on individual services and call routing rules. Valid values: `sd`, `hd`, `fullhd`.",
			},
			"max_presentation_bandwidth_ratio": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				Default:  int64default.StaticInt64(75),
				Validators: []validator.Int64{
					int64validator.Between(25, 75),
				},
				MarkdownDescription: "The maximum percentage of call bandwidth to be allocated to sending presentation. Range: 25 to 75. Default: 75.",
			},
			"media_ports_end": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				Default:  int64default.StaticInt64(49999),
				Validators: []validator.Int64{
					int64validator.Between(10000, 49999),
				},
				MarkdownDescription: "The end value for the range of ports (UDP and TCP) that all Conferencing Nodes will use to send media (for all call protocols). The media port range must contain at least 100 ports. Range: 10000 to 49999. Default: 49999.",
			},
			"media_ports_start": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				Default:  int64default.StaticInt64(40000),
				Validators: []validator.Int64{
					int64validator.Between(10000, 49999),
				},
				MarkdownDescription: "The start value for the range of ports (UDP and TCP) that all Conferencing Nodes will use to send media (for all call protocols). The media port range must contain at least 100 ports. Range: 10000 to 49999. Default: 40000.",
			},
			"ocsp_responder_url": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(""),
				Validators: []validator.String{
					stringvalidator.LengthAtMost(255),
				},
				MarkdownDescription: "The URL to which OCSP requests will be sent either if the OCSP state is set to Override, or if the OCSP state is set to On but there is no URL specified in the TLS certificate. Maximum length: 255 characters.",
			},
			"ocsp_state": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("OFF"),
				Validators: []validator.String{
					stringvalidator.OneOf("OFF", "ON", "OVERRIDE"),
				},
				MarkdownDescription: "Whether to use OCSP when checking the validity of TLS certificates. `ON`: An OCSP request will be sent to the URL specified in the TLS certificate. `OVERRIDE`: An OCSP request will be sent to the URL specified in the OCSP responder URL field. Valid values: `OFF`, `ON`, `OVERRIDE`.",
			},
			"pin_entry_timeout": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				Default:  int64default.StaticInt64(120),
				Validators: []validator.Int64{
					int64validator.Between(30, 86400),
				},
				MarkdownDescription: "The length of time (in seconds) for which a participant will be permitted to remain at the PIN entry screen before being disconnected. Range: 30 to 86400. Default: 120.",
			},
			"resource_uri": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The URI that identifies this resource.",
			},
			"session_timeout_enabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
				MarkdownDescription: "Determines whether inactive users are automatically logged out of the Management Node administration interface after a period of time. If disabled, users of the administrator interface are never timed out.",
			},
			"signalling_ports_end": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				Default:  int64default.StaticInt64(39999),
				Validators: []validator.Int64{
					int64validator.Between(10000, 49999),
				},
				MarkdownDescription: "The end value for the range of ports (UDP and TCP) that all Conferencing Nodes will use to send signaling (for H.323, H.245 and SIP). Range: 10000 to 49999. Default: 39999.",
			},
			"signalling_ports_start": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				Default:  int64default.StaticInt64(33000),
				Validators: []validator.Int64{
					int64validator.Between(10000, 49999),
				},
				MarkdownDescription: "The start value for the range of ports (UDP and TCP) that all Conferencing Nodes will use to send signaling (for H.323, H.245 and SIP). Range: 10000 to 49999. Default: 33000.",
			},
			"sip_tls_cert_verify_mode": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("OFF"),
				Validators: []validator.String{
					stringvalidator.OneOf("OFF", "ON"),
				},
				MarkdownDescription: "Determines whether to verify the peer certificate for connections over SIP TLS. `OFF`: the peer certificate will not be verified; all connections will be allowed. `ON`: the peer certificate will be verified, and the peer's remote identities will be compared against the Application Unique String (AUS). Valid values: `OFF`, `ON`.",
			},
			"site_banner": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(""),
				Validators: []validator.String{
					stringvalidator.LengthAtMost(255),
				},
				MarkdownDescription: "Text of the banner to display on the top of every page of this Pexip Infinity administrator web interface. Maximum length: 255 characters.",
			},
			"site_banner_bg": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("#c0c0c0"),
				MarkdownDescription: "The background color for the site banner.",
			},
			"site_banner_fg": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("#000000"),
				MarkdownDescription: "The text color for the site banner.",
			},
			"teams_enable_powerpoint_render": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "This setting is intended for future use to enable PowerPoint Live content in Microsoft Teams calls. Check the online documentation for the latest status for this feature.",
			},
			"waiting_for_chair_timeout": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				Default:  int64default.StaticInt64(900),
				Validators: []validator.Int64{
					int64validator.Between(0, 86400),
				},
				MarkdownDescription: "The length of time (in seconds) for which a Guest participant will remain at the waiting screen if a Host does not join, before being disconnected. Range: 0 to 86400. Default: 900.",
			},
		},
		MarkdownDescription: "Manages the global system configuration with the Infinity service. This is a singleton resource - only one global configuration exists per system.",
	}
}

func (r *InfinityGlobalConfigurationResource) buildUpdateRequest(plan *InfinityGlobalConfigurationResourceModel) *config.GlobalConfigurationUpdateRequest {
	updateRequest := &config.GlobalConfigurationUpdateRequest{
		CloudProvider:                       plan.CloudProvider.ValueString(),
		ContactEmailAddress:                 plan.ContactEmailAddress.ValueString(),
		ContentSecurityPolicyHeader:         plan.ContentSecurityPolicyHeader.ValueString(),
		ContentSecurityPolicyState:          plan.ContentSecurityPolicyState.ValueBool(),
		CryptoMode:                          plan.CryptoMode.ValueString(),
		DeploymentUUID:                      plan.DeploymentUUID.ValueString(),
		ErrorReportingURL:                   plan.ErrorReportingURL.ValueString(),
		EnableSIPUDP:                        plan.EnableSIPUDP.ValueBool(),
		LegacyAPIUsername:                   plan.LegacyAPIUsername.ValueString(),
		LegacyAPIPassword:                   plan.LegacyAPIPassword.ValueString(),
		LiveviewShowConferences:             plan.LiveviewShowConferences.ValueBool(),
		LocalMssipDomain:                    plan.LocalMssipDomain.ValueString(),
		LogonBanner:                         plan.LogonBanner.ValueString(),
		ManagementStartPage:                 plan.ManagementStartPage.ValueString(),
		MaxPixelsPerSecond:                  plan.MaxPixelsPerSecond.ValueString(),
		OcspResponderURL:                    plan.OcspResponderURL.ValueString(),
		OcspState:                           plan.OcspState.ValueString(),
		SipTLSCertVerifyMode:                plan.SipTLSCertVerifyMode.ValueString(),
		SiteBanner:                          plan.SiteBanner.ValueString(),
		SiteBannerBg:                        plan.SiteBannerBg.ValueString(),
		SiteBannerFg:                        plan.SiteBannerFg.ValueString(),
		TeamsEnablePowerpointRender:         plan.TeamsEnablePowerpointRender.ValueBool(),
		EnableWebRTC:                        plan.EnableWebRTC.ValueBool(),
		EnableSIP:                           plan.EnableSIP.ValueBool(),
		EnableH323:                          plan.EnableH323.ValueBool(),
		EnableRTMP:                          plan.EnableRTMP.ValueBool(),
		EnableAnalytics:                     plan.EnableAnalytics.ValueBool(),
		EnableApplicationAPI:                plan.EnableApplicationAPI.ValueBool(),
		EnableBreakoutRooms:                 plan.EnableBreakoutRooms.ValueBool(),
		EnableChat:                          plan.EnableChat.ValueBool(),
		EnableClock:                         plan.EnableClock.ValueBool(),
		EnableClock:                         plan.EnableClock.ValueBool(),
		EnableDenoise:                       plan.EnableDenoise.ValueBool(),
		EnableDialout:                       plan.EnableDialout.ValueBool(),
		EnableDirectory:                     plan.EnableDirectory.ValueBool(),
		EnableEdgeNonMesh:                   plan.EnableEdgeNonMesh.ValueBool(),
		EnableFecc:                          plan.EnableFecc.ValueBool(),
		EnableLegacyDialoutAPI:              plan.EnableLegacyDialoutAPI.ValueBool(),
		EnableLyncAutoEscalate:              plan.EnableLyncAutoEscalate.ValueBool(),
		EnableLyncVbss:                      plan.EnableLyncVbss.ValueBool(),
		EnableMlvad:                         plan.EnableMlvad.ValueBool(),
		EnableSoftmute:                      plan.EnableSoftmute.ValueBool(),
		EnableSSH:                           plan.EnableSSH.ValueBool(),
		EnableTurn443:                       plan.EnableTurn443.ValueBool(),
		ErrorReportingEnabled:               plan.ErrorReportingEnabled.ValueBool(),
		EsConnectionTimeout:                 int(plan.EsConnectionTimeout.ValueInt64()),
		EsInitialRetryBackoff:               int(plan.EsInitialRetryBackoff.ValueInt64()),
		EsMaximumRetryBackoff:               int(plan.EsMaximumRetryBackoff.ValueInt64()),
		EsMediaStreamsWait:                  int(plan.EsMediaStreamsWait.ValueInt64()),
		EsMetricsUpdateInterval:             int(plan.EsMetricsUpdateInterval.ValueInt64()),
		EsShortTermMemoryExpiration:         int(plan.EsShortTermMemoryExpiration.ValueInt64()),
		ExternalParticipantAvatarLookup:     plan.ExternalParticipantAvatarLookup.ValueBool(),
		GuestsOnlyTimeout:                   int(plan.GuestsOnlyTimeout.ValueInt64()),
		LiveCaptionsVMRDefault:              plan.LiveCaptionsVMRDefault.ValueBool(),
		LogsMaxAge:                          int(plan.LogsMaxAge.ValueInt64()),
		ManagementSessionTimeout:            int(plan.ManagementSessionTimeout.ValueInt64()),
		SessionTimeoutEnabled:               plan.SessionTimeoutEnabled.ValueBool(),
		WaitingForChairTimeout:              int(plan.WaitingForChairTimeout.ValueInt64()),
		EjectLastParticipantBackstopTimeout: int(plan.EjectLastParticipantBackstopTimeout.ValueInt64()),
		MaxPresentationBandwidthRatio:       int(plan.MaxPresentationBandwidthRatio.ValueInt64()),
		MediaPortsStart:                     int(plan.MediaPortsStart.ValueInt64()),
		MediaPortsEnd:                       int(plan.MediaPortsEnd.ValueInt64()),
		SignallingPortsStart:                int(plan.SignallingPortsStart.ValueInt64()),
		SignallingPortsEnd:                  int(plan.SignallingPortsEnd.ValueInt64()),
		PinEntryTimeout:                     int(plan.PinEntryTimeout.ValueInt64()),
		BdpmMaxPinFailuresPerWindow:         int(plan.BdpmMaxPinFailuresPerWindow.ValueInt64()),
		BdpmMaxScanAttemptsPerWindow:        int(plan.BdpmMaxScanAttemptsPerWindow.ValueInt64()),
		BdpmPinChecksEnabled:                plan.BdpmPinChecksEnabled.ValueBool(),
		BdpmScanQuarantineEnabled:           plan.BdpmScanQuarantineEnabled.ValueBool(),
		BurstingEnabled:                     plan.BurstingEnabled.ValueBool(),
	}

	if !plan.DisabledCodecs.IsNull() {
		var disabledCodecs []config.CodecValue
		for _, v := range plan.DisabledCodecs.Elements() {
			disabledCodecs = append(disabledCodecs, config.CodecValue{Value: v.(types.String).ValueString()})
		}
		updateRequest.DisabledCodecs = disabledCodecs
	}

	// handle pointers
	if !plan.AWSAccessKey.IsNull() && !plan.AWSAccessKey.IsUnknown() {
		val := plan.AWSAccessKey.ValueString()
		updateRequest.AWSAccessKey = &val
	}
	if !plan.AWSSecretKey.IsNull() && !plan.AWSSecretKey.IsUnknown() {
		val := plan.AWSSecretKey.ValueString()
		updateRequest.AWSSecretKey = &val
	}
	if !plan.AzureClientID.IsNull() && !plan.AzureClientID.IsUnknown() {
		val := plan.AzureClientID.ValueString()
		updateRequest.AzureClientID = &val
	}
	if !plan.AzureSecret.IsNull() && !plan.AzureSecret.IsUnknown() {
		val := plan.AzureSecret.ValueString()
		updateRequest.AzureSecret = &val
	}
	if !plan.AzureSubscriptionID.IsNull() && !plan.AzureSubscriptionID.IsUnknown() {
		val := plan.AzureSubscriptionID.ValueString()
		updateRequest.AzureSubscriptionID = &val
	}
	if !plan.AzureTenant.IsNull() && !plan.AzureTenant.IsUnknown() {
		val := plan.AzureTenant.ValueString()
		updateRequest.AzureTenant = &val
	}
	if !plan.BurstingMinLifetime.IsNull() && !plan.BurstingMinLifetime.IsUnknown() {
		val := int(plan.BurstingMinLifetime.ValueInt64())
		updateRequest.BurstingMinLifetime = &val
	}
	if !plan.BurstingThreshold.IsNull() && !plan.BurstingThreshold.IsUnknown() {
		val := int(plan.BurstingThreshold.ValueInt64())
		updateRequest.BurstingThreshold = &val
	}
	if !plan.DefaultTheme.IsNull() && !plan.DefaultTheme.IsUnknown() {
		val := plan.DefaultTheme.ValueString()
		updateRequest.DefaultTheme = &config.IVRTheme{Name: val}
	}
	if !plan.DefaultWebappAlias.IsNull() && !plan.DefaultWebappAlias.IsUnknown() {
		val := plan.DefaultWebappAlias.ValueString()
		updateRequest.DefaultWebappAlias = &val
	}
	if !plan.GcpClientEmail.IsNull() && !plan.GcpClientEmail.IsUnknown() {
		val := plan.GcpClientEmail.ValueString()
		updateRequest.GcpClientEmail = &val
	}
	if !plan.GcpPrivateKey.IsNull() && !plan.GcpPrivateKey.IsUnknown() {
		val := plan.GcpPrivateKey.ValueString()
		updateRequest.GcpPrivateKey = &val
	}
	if !plan.GcpProjectID.IsNull() && !plan.GcpProjectID.IsUnknown() {
		val := plan.GcpProjectID.ValueString()
		updateRequest.GcpProjectID = &val
	}
	if !plan.ManagementQos.IsNull() && !plan.ManagementQos.IsUnknown() {
		val := int(plan.ManagementQos.ValueInt64())
		updateRequest.ManagementQos = &val
	}
	if !plan.MaxCallrateIn.IsNull() && !plan.MaxCallrateIn.IsUnknown() {
		val := int(plan.MaxCallrateIn.ValueInt64())
		updateRequest.MaxCallrateIn = &val
	}
	if !plan.MaxCallrateOut.IsNull() && !plan.MaxCallrateOut.IsUnknown() {
		val := int(plan.MaxCallrateOut.ValueInt64())
		updateRequest.MaxCallrateOut = &val
	}

	return updateRequest
}

func (r *InfinityGlobalConfigurationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// For singleton resources, Create is actually Update since the resource always exists
	plan := &InfinityGlobalConfigurationResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateRequest := r.buildUpdateRequest(plan)

	_, err := r.InfinityClient.Config().UpdateGlobalConfiguration(ctx, updateRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Infinity global configuration",
			fmt.Sprintf("Could not update Infinity global configuration: %s", err),
		)
		return
	}

	// Re-read the resource to get the latest state
	updatedModel, err := r.read(ctx, plan.AWSSecretKey.ValueStringPointer(), plan.AzureSecret.ValueStringPointer(), plan.GcpPrivateKey.ValueStringPointer(), plan.LegacyAPIPassword.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Updated Infinity global configuration",
			fmt.Sprintf("Could not read updated Infinity global configuration: %s", err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, updatedModel)...)
}

func (r *InfinityGlobalConfigurationResource) read(ctx context.Context, awsSecretKey, azureSecret, gcpPrivateKey *string, legacyAPIPassword string) (*InfinityGlobalConfigurationResourceModel, error) {
	var data InfinityGlobalConfigurationResourceModel

	srv, err := r.InfinityClient.Config().GetGlobalConfiguration(ctx)
	if err != nil {
		return nil, err
	}

	if srv.ResourceURI == "" {
		return nil, fmt.Errorf("global configuration not found")
	}

	data.ID = types.StringValue(srv.ResourceURI)
	data.AWSAccessKey = types.StringPointerValue(srv.AWSAccessKey)
	data.AWSSecretKey = types.StringPointerValue(awsSecretKey)
	data.AzureClientID = types.StringPointerValue(srv.AzureClientID)
	data.AzureSecret = types.StringPointerValue(azureSecret)
	data.AzureSubscriptionID = types.StringPointerValue(srv.AzureSubscriptionID)
	data.AzureTenant = types.StringPointerValue(srv.AzureTenant)
	data.BdpmMaxPinFailuresPerWindow = types.Int64Value(int64(srv.BdpmMaxPinFailuresPerWindow))
	data.BdpmMaxScanAttemptsPerWindow = types.Int64Value(int64(srv.BdpmMaxScanAttemptsPerWindow))
	data.BdpmPinChecksEnabled = types.BoolValue(srv.BdpmPinChecksEnabled)
	data.BdpmScanQuarantineEnabled = types.BoolValue(srv.BdpmScanQuarantineEnabled)
	data.BurstingEnabled = types.BoolValue(srv.BurstingEnabled)
	data.CloudProvider = types.StringValue(srv.CloudProvider)
	data.ContactEmailAddress = types.StringValue(srv.ContactEmailAddress)
	data.ContentSecurityPolicyHeader = types.StringValue(srv.ContentSecurityPolicyHeader)
	data.ContentSecurityPolicyState = types.BoolValue(srv.ContentSecurityPolicyState)
	data.CryptoMode = types.StringValue(srv.CryptoMode)
	data.DefaultWebappAlias = types.StringPointerValue(srv.DefaultWebappAlias)
	data.DeploymentUUID = types.StringValue(srv.DeploymentUUID)
	data.ErrorReportingURL = types.StringValue(srv.ErrorReportingURL)
	data.LegacyAPIUsername = types.StringValue(srv.LegacyAPIUsername)
	data.LegacyAPIPassword = types.StringValue(legacyAPIPassword)
	data.LiveviewShowConferences = types.BoolValue(srv.LiveviewShowConferences)
	data.LocalMssipDomain = types.StringValue(srv.LocalMssipDomain)
	data.LogonBanner = types.StringValue(srv.LogonBanner)
	data.ManagementStartPage = types.StringValue(srv.ManagementStartPage)
	data.MaxPixelsPerSecond = types.StringValue(srv.MaxPixelsPerSecond)
	data.OcspResponderURL = types.StringValue(srv.OcspResponderURL)
	data.OcspState = types.StringValue(srv.OcspState)
	data.SipTLSCertVerifyMode = types.StringValue(srv.SipTLSCertVerifyMode)
	data.SiteBanner = types.StringValue(srv.SiteBanner)
	data.SiteBannerBg = types.StringValue(srv.SiteBannerBg)
	data.SiteBannerFg = types.StringValue(srv.SiteBannerFg)
	data.TeamsEnablePowerpointRender = types.BoolValue(srv.TeamsEnablePowerpointRender)
	data.EnableWebRTC = types.BoolValue(srv.EnableWebRTC)
	data.EnableSIP = types.BoolValue(srv.EnableSIP)
	data.EnableH323 = types.BoolValue(srv.EnableH323)
	data.EnableRTMP = types.BoolValue(srv.EnableRTMP)
	data.EnableAnalytics = types.BoolValue(srv.EnableAnalytics)
	data.EnableApplicationAPI = types.BoolValue(srv.EnableApplicationAPI)
	data.EnableBreakoutRooms = types.BoolValue(srv.EnableBreakoutRooms)
	data.EnableChat = types.BoolValue(srv.EnableChat)
	data.EnableClock = types.BoolValue(srv.EnableClock)
	data.EnableClock = types.BoolValue(srv.EnableClock)
	data.EnableDenoise = types.BoolValue(srv.EnableDenoise)
	data.EnableDialout = types.BoolValue(srv.EnableDialout)
	data.EnableDirectory = types.BoolValue(srv.EnableDirectory)
	data.EnableEdgeNonMesh = types.BoolValue(srv.EnableEdgeNonMesh)
	data.EnableFecc = types.BoolValue(srv.EnableFecc)
	data.EnableLegacyDialoutAPI = types.BoolValue(srv.EnableLegacyDialoutAPI)
	data.EnableLyncAutoEscalate = types.BoolValue(srv.EnableLyncAutoEscalate)
	data.EnableLyncVbss = types.BoolValue(srv.EnableLyncVbss)
	data.EnableMlvad = types.BoolValue(srv.EnableMlvad)
	data.EnableSIPUDP = types.BoolValue(srv.EnableSIPUDP)
	data.EnableSoftmute = types.BoolValue(srv.EnableSoftmute)
	data.EnableSSH = types.BoolValue(srv.EnableSSH)
	data.EnableTurn443 = types.BoolValue(srv.EnableTurn443)
	data.ErrorReportingEnabled = types.BoolValue(srv.ErrorReportingEnabled)
	data.EsConnectionTimeout = types.Int64Value(int64(srv.EsConnectionTimeout))
	data.EsInitialRetryBackoff = types.Int64Value(int64(srv.EsInitialRetryBackoff))
	data.EsMaximumRetryBackoff = types.Int64Value(int64(srv.EsMaximumRetryBackoff))
	data.EsMediaStreamsWait = types.Int64Value(int64(srv.EsMediaStreamsWait))
	data.EsMetricsUpdateInterval = types.Int64Value(int64(srv.EsMetricsUpdateInterval))
	data.EsShortTermMemoryExpiration = types.Int64Value(int64(srv.EsShortTermMemoryExpiration))
	data.ExternalParticipantAvatarLookup = types.BoolValue(srv.ExternalParticipantAvatarLookup)
	data.GcpProjectID = types.StringPointerValue(srv.GcpProjectID)
	data.GcpClientEmail = types.StringPointerValue(srv.GcpClientEmail)
	data.GcpPrivateKey = types.StringPointerValue(gcpPrivateKey)
	data.GuestsOnlyTimeout = types.Int64Value(int64(srv.GuestsOnlyTimeout))
	data.LiveCaptionsVMRDefault = types.BoolValue(srv.LiveCaptionsVMRDefault)
	data.LogsMaxAge = types.Int64Value(int64(srv.LogsMaxAge))
	data.ManagementSessionTimeout = types.Int64Value(int64(srv.ManagementSessionTimeout))
	data.SessionTimeoutEnabled = types.BoolValue(srv.SessionTimeoutEnabled)
	data.WaitingForChairTimeout = types.Int64Value(int64(srv.WaitingForChairTimeout))
	data.EjectLastParticipantBackstopTimeout = types.Int64Value(int64(srv.EjectLastParticipantBackstopTimeout))
	data.MaxPresentationBandwidthRatio = types.Int64Value(int64(srv.MaxPresentationBandwidthRatio))
	data.MediaPortsStart = types.Int64Value(int64(srv.MediaPortsStart))
	data.MediaPortsEnd = types.Int64Value(int64(srv.MediaPortsEnd))
	data.SignallingPortsStart = types.Int64Value(int64(srv.SignallingPortsStart))
	data.SignallingPortsEnd = types.Int64Value(int64(srv.SignallingPortsEnd))
	data.PinEntryTimeout = types.Int64Value(int64(srv.PinEntryTimeout))
	data.BdpmMaxPinFailuresPerWindow = types.Int64Value(int64(srv.BdpmMaxPinFailuresPerWindow))
	data.BdpmMaxScanAttemptsPerWindow = types.Int64Value(int64(srv.BdpmMaxScanAttemptsPerWindow))
	data.BdpmPinChecksEnabled = types.BoolValue(srv.BdpmPinChecksEnabled)
	data.BdpmScanQuarantineEnabled = types.BoolValue(srv.BdpmScanQuarantineEnabled)
	data.BurstingEnabled = types.BoolValue(srv.BurstingEnabled)

	// Convert default theme from SDK to Terraform format
	if srv.DefaultTheme != nil {
		data.DefaultTheme = types.StringValue(srv.DefaultTheme.Name)
	}

	if srv.ManagementQos != nil {
		data.ManagementQos = types.Int64Value(int64(*srv.ManagementQos))
	}
	if srv.BurstingMinLifetime != nil {
		data.BurstingMinLifetime = types.Int64Value(int64(*srv.BurstingMinLifetime))
	}
	if srv.BurstingThreshold != nil {
		data.BurstingThreshold = types.Int64Value(int64(*srv.BurstingThreshold))
	}
	if srv.MaxCallrateIn != nil {
		data.MaxCallrateIn = types.Int64Value(int64(*srv.MaxCallrateIn))
	}
	if srv.MaxCallrateOut != nil {
		data.MaxCallrateOut = types.Int64Value(int64(*srv.MaxCallrateOut))
	}

	var disabledCodecs []attr.Value
	for _, v := range srv.DisabledCodecs {
		disabledCodecs = append(disabledCodecs, types.StringValue(v.Value))
	}
	data.DisabledCodecs, _ = types.SetValue(types.StringType, disabledCodecs)

	return &data, nil
}

func (r *InfinityGlobalConfigurationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	state := &InfinityGlobalConfigurationResourceModel{}

	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	state, err := r.read(ctx, state.AWSSecretKey.ValueStringPointer(), state.AzureSecret.ValueStringPointer(), state.GcpPrivateKey.ValueStringPointer(), state.LegacyAPIPassword.ValueString())
	if err != nil {
		// Check if the error is a 404 (not found) - unlikely for singleton resources
		if isNotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading Infinity global configuration",
			fmt.Sprintf("Could not read Infinity global configuration: %s", err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *InfinityGlobalConfigurationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	plan := &InfinityGlobalConfigurationResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateRequest := r.buildUpdateRequest(plan)

	_, err := r.InfinityClient.Config().UpdateGlobalConfiguration(ctx, updateRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Infinity global configuration",
			fmt.Sprintf("Could not update Infinity global configuration: %s", err),
		)
		return
	}

	// Re-read the resource to get the latest state
	updatedModel, err := r.read(ctx, plan.AWSSecretKey.ValueStringPointer(), plan.AzureSecret.ValueStringPointer(), plan.GcpPrivateKey.ValueStringPointer(), plan.LegacyAPIPassword.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Updated Infinity global configuration",
			fmt.Sprintf("Could not read updated Infinity global configuration: %s", err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, updatedModel)...)
}

func (r *InfinityGlobalConfigurationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// For singleton resources, delete means resetting all fields to their schema defaults.
	tflog.Info(ctx, "Deleting Infinity global configuration (resetting to defaults)")

	burstingMinLifetimeDefault := 50
	burstingThresholdDefault := 5
	managementQosDefault := 0
	gcpPrivateKeyDefault := ""

	updateRequest := &config.GlobalConfigurationUpdateRequest{
		// Nullable fields with null defaults — cleared to nil
		AWSAccessKey:        nil,
		AWSSecretKey:        nil,
		AzureClientID:       nil,
		AzureSecret:         nil,
		AzureSubscriptionID: nil,
		AzureTenant:         nil,
		DefaultTheme:        nil,
		DefaultWebappAlias:  nil,
		GcpClientEmail:      nil,
		GcpProjectID:        nil,
		MaxCallrateIn:       nil,
		MaxCallrateOut:      nil,

		// Nullable fields with non-null defaults — must send the default value explicitly
		BurstingMinLifetime: &burstingMinLifetimeDefault,
		BurstingThreshold:   &burstingThresholdDefault,
		GcpPrivateKey:       &gcpPrivateKeyDefault,
		ManagementQos:       &managementQosDefault,

		// Non-nullable fields — schema defaults
		BdpmMaxPinFailuresPerWindow:         20,
		BdpmMaxScanAttemptsPerWindow:        20,
		BdpmPinChecksEnabled:                true,
		BdpmScanQuarantineEnabled:           true,
		BurstingEnabled:                     false,
		CloudProvider:                       "AWS",
		ContactEmailAddress:                 "",
		ContentSecurityPolicyHeader:         "upgrade-insecure-requests; default-src 'self'; frame-ancestors 'self'; frame-src 'self' https://telemetryservice.firstpartyapps.oaspapps.com/telemetryservice/telemetryproxy.html https://*.microsoft.com https://*.office.com; style-src 'self' 'unsafe-inline' https://*.microsoft.com https://*.office.com; object-src 'self'; font-src 'self' https://*.microsoft.com https://*.office.com; img-src 'self' https://www.adobe.com data: blob:; script-src 'self' 'unsafe-inline' 'unsafe-eval' https://*.microsoft.com https://*.office.com https://ajax.aspnetcdn.com https://api.keen.io; media-src 'self' blob:; connect-src 'self' https://*.microsoft.com https://*.office.com https://example.com;",
		ContentSecurityPolicyState:          true,
		CryptoMode:                          "besteffort",
		EjectLastParticipantBackstopTimeout: 0,
		EnableAnalytics:                     false,
		EnableApplicationAPI:                true,
		EnableBreakoutRooms:                 false,
		EnableChat:                          true,
		EnableClock:                         false,
		EnableDenoise:                       true,
		EnableDialout:                       true,
		EnableDirectory:                     true,
		EnableEdgeNonMesh:                   true,
		EnableFecc:                          true,
		EnableH323:                          true,
		EnableLegacyDialoutAPI:              false,
		EnableLyncAutoEscalate:              false,
		EnableLyncVbss:                      false,
		EnableMlvad:                         false,
		EnableRTMP:                          true,
		EnableSIP:                           true,
		EnableSIPUDP:                        false,
		EnableSoftmute:                      true,
		EnableSSH:                           true,
		EnableTurn443:                       false,
		EnableWebRTC:                        true,
		ErrorReportingEnabled:               false,
		ErrorReportingURL:                   "https://acr.pexip.com",
		EsConnectionTimeout:                 7,
		EsInitialRetryBackoff:               1,
		EsMaximumDeferredPosts:              1000,
		EsMaximumRetryBackoff:               1800,
		EsMediaStreamsWait:                  1,
		EsMetricsUpdateInterval:             60,
		EsShortTermMemoryExpiration:         2,
		ExternalParticipantAvatarLookup:     true,
		GuestsOnlyTimeout:                   60,
		LegacyAPIUsername:                   "",
		LegacyAPIPassword:                   "",
		LiveCaptionsVMRDefault:              false,
		LiveviewShowConferences:             true,
		LocalMssipDomain:                    "",
		LogonBanner:                         "",
		LogsMaxAge:                          0,
		ManagementSessionTimeout:            30,
		ManagementStartPage:                 "/admin/conferencingstatus/deploymentgraph/deployment_graph/",
		MaxPixelsPerSecond:                  "hd",
		MaxPresentationBandwidthRatio:       75,
		MediaPortsEnd:                       49999,
		MediaPortsStart:                     40000,
		OcspResponderURL:                    "",
		OcspState:                           "OFF",
		PinEntryTimeout:                     120,
		SessionTimeoutEnabled:               true,
		SignallingPortsEnd:                  39999,
		SignallingPortsStart:                33000,
		SipTLSCertVerifyMode:                "OFF",
		SiteBanner:                          "",
		SiteBannerBg:                        "#c0c0c0",
		SiteBannerFg:                        "#000000",
		TeamsEnablePowerpointRender:         false,
		WaitingForChairTimeout:              900,
	}

	_, err := r.InfinityClient.Config().UpdateGlobalConfiguration(ctx, updateRequest)
	if err != nil && !isNotFoundError(err) && !isLookupError(err) {
		resp.Diagnostics.AddError(
			"Error Deleting Infinity global configuration",
			fmt.Sprintf("Could not delete Infinity global configuration: %s", err),
		)
		return
	}
}

func (r *InfinityGlobalConfigurationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// For singleton resources, the import ID doesn't matter since there's only one instance
	tflog.Trace(ctx, "Importing Infinity global configuration")

	// Read the resource from the API
	model, err := r.read(ctx, nil, nil, nil, "")
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Importing Infinity Global Configuration",
			fmt.Sprintf("Could not import Infinity global configuration: %s", err),
		)
		return
	}

	// Set the state from the imported resource
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *InfinityGlobalConfigurationResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var data InfinityGlobalConfigurationResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Skip validation if bursting_enabled or cloud_provider are unknown (known after apply).
	if data.BurstingEnabled.IsUnknown() || data.CloudProvider.IsUnknown() {
		return
	}

	// cloud_provider defaults to "AWS", so treat null (not set) the same as "AWS".
	isAWSProvider := data.CloudProvider.IsNull() || data.CloudProvider.ValueString() == "AWS"
	if data.BurstingEnabled.ValueBool() && isAWSProvider {
		if data.AWSAccessKey.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("aws_access_key"),
				"Missing AWS Access Key",
				"aws_access_key must be configured when bursting_enabled is true and cloud_provider is \"AWS\".",
			)
		}
		if data.AWSSecretKey.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("aws_secret_key"),
				"Missing AWS Secret Key",
				"aws_secret_key must be configured when bursting_enabled is true and cloud_provider is \"AWS\".",
			)
		}
	}

	if data.BurstingEnabled.ValueBool() && data.CloudProvider.ValueString() == "GCP" {
		if data.GcpClientEmail.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("gcp_client_email"),
				"Missing GCP Service Account ID",
				"gcp_client_email must be configured when bursting_enabled is true and cloud_provider is \"GCP\".",
			)
		}
		if data.GcpPrivateKey.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("gcp_private_key"),
				"Missing GCP Private Key",
				"gcp_private_key must be configured when bursting_enabled is true and cloud_provider is \"GCP\".",
			)
		}
		if data.GcpProjectID.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("gcp_project_id"),
				"Missing GCP Project ID",
				"gcp_project_id must be configured when bursting_enabled is true and cloud_provider is \"GCP\".",
			)
		}
	}

	if data.BurstingEnabled.ValueBool() && data.CloudProvider.ValueString() == "AZURE" {
		if data.AzureClientID.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("azure_client_id"),
				"Missing Azure Client ID",
				"azure_client_id must be configured when bursting_enabled is true and cloud_provider is \"AZURE\".",
			)
		}
		if data.AzureSecret.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("azure_secret"),
				"Missing Azure Secret Key",
				"azure_secret must be configured when bursting_enabled is true and cloud_provider is \"AZURE\".",
			)
		}
		if data.AzureSubscriptionID.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("azure_subscription_id"),
				"Missing Azure Subscription ID",
				"azure_subscription_id must be configured when bursting_enabled is true and cloud_provider is \"AZURE\".",
			)
		}
		if data.AzureTenant.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("azure_tenant"),
				"Missing Azure Tenant ID",
				"azure_tenant must be configured when bursting_enabled is true and cloud_provider is \"AZURE\".",
			)
		}
	}
}
