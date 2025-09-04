/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"

	//"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/pexip/go-infinity-sdk/v38/config"
)

var (
	_ resource.ResourceWithImportState = (*InfinityGlobalConfigurationResource)(nil)
)

type InfinityGlobalConfigurationResource struct {
	InfinityClient InfinityClient
}

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
	DefaultToNewWebapp                  types.Bool   `tfsdk:"default_to_new_webapp"`
	DefaultWebapp                       types.String `tfsdk:"default_webapp"`
	DefaultWebappAlias                  types.String `tfsdk:"default_webapp_alias"`
	DeploymentUUID                      types.String `tfsdk:"deployment_uuid"`
	DisabledCodecs                      types.Set    `tfsdk:"disabled_codecs"`
	EjectLastParticipantBackstopTimeout types.Int64  `tfsdk:"eject_last_participant_backstop_timeout"`
	EnableAnalytics                     types.Bool   `tfsdk:"enable_analytics"`
	EnableApplicationAPI                types.Bool   `tfsdk:"enable_application_api"`
	EnableBreakoutRooms                 types.Bool   `tfsdk:"enable_breakout_rooms"`
	EnableChat                          types.Bool   `tfsdk:"enable_chat"`
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
	EnableMultiscreen                   types.Bool   `tfsdk:"enable_multiscreen"`
	EnablePushNotifications             types.Bool   `tfsdk:"enable_push_notifications"`
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
	EsMaximumDeferredPosts              types.Int64  `tfsdk:"es_maximum_deferred_posts"`
	EsMaximumRetryBackoff               types.Int64  `tfsdk:"es_maximum_retry_backoff"`
	EsMediaStreamsWait                  types.Int64  `tfsdk:"es_media_streams_wait"`
	EsMetricsUpdateInterval             types.Int64  `tfsdk:"es_metrics_update_interval"`
	EsShortTermMemoryExpiration         types.Int64  `tfsdk:"es_short_term_memory_expiration"`
	ExternalParticipantAvatarLookup     types.Bool   `tfsdk:"external_participant_avatar_lookup"`
	GcpClientEmail                      types.String `tfsdk:"gcp_client_email"`
	GcpPrivateKey                       types.String `tfsdk:"gcp_private_key"`
	GcpProjectID                        types.String `tfsdk:"gcp_project_id"`
	GuestsOnlyTimeout                   types.Int64  `tfsdk:"guests_only_timeout"`
	LegacyAPIHTTP                       types.Bool   `tfsdk:"legacy_api_http"`
	LegacyAPIUsername                   types.String `tfsdk:"legacy_api_username"`
	LegacyAPIPassword                   types.String `tfsdk:"legacy_api_password"`
	LiveCaptionsAPIGateway              types.String `tfsdk:"live_captions_api_gateway"`
	LiveCaptionsAppID                   types.String `tfsdk:"live_captions_app_id"`
	LiveCaptionsEnabled                 types.Bool   `tfsdk:"live_captions_enabled"`
	LiveCaptionsPublicJWTKey            types.String `tfsdk:"live_captions_public_jwt_key"`
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
	PssCustomerID                       types.String `tfsdk:"pss_customer_id"`
	PssEnabled                          types.Bool   `tfsdk:"pss_enabled"`
	PssGateway                          types.String `tfsdk:"pss_gateway"`
	PssToken                            types.String `tfsdk:"pss_token"`
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
			},
			"aws_access_key": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The Amazon Web Services access key ID for the AWS user.",
			},
			"aws_secret_key": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				MarkdownDescription: "The Amazon Web Services secret access key.",
			},
			"azure_client_id": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				MarkdownDescription: "The Azure client (Application) ID.",
			},
			"azure_secret": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				MarkdownDescription: "The Azure secret key.",
			},
			"azure_subscription_id": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The Azure subscription ID.",
			},
			"azure_tenant": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The Azure tenant ID.",
			},
			"bdpm_max_pin_failures_per_window": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(20),
				MarkdownDescription: "Maximum PIN failures per window for Break-in Defense Policy Manager.",
			},
			"bdpm_max_scan_attempts_per_window": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(20),
				MarkdownDescription: "Maximum scan attempts per window for Break-in Defense Policy Manager.",
			},
			"bdpm_pin_checks_enabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
				MarkdownDescription: "Enable PIN checks for Break-in Defense Policy Manager.",
			},
			"bdpm_scan_quarantine_enabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
				MarkdownDescription: "Enable scan quarantine for Break-in Defense Policy Manager.",
			},
			"bursting_enabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Enable cloud bursting functionality.",
			},
			"bursting_min_lifetime": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(50),
				MarkdownDescription: "Minimum lifetime (minutes) for bursting nodes.",
			},
			"bursting_threshold": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(5),
				MarkdownDescription: "Threshold for bursting node activation.",
			},
			"cloud_provider": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("AWS"),
				Validators: []validator.String{
					stringvalidator.OneOf("AWS", "AZURE", "GCP"),
				},
				MarkdownDescription: "Cloud provider for bursting (AWS, AZURE, GCP).",
			},
			"contact_email_address": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
				MarkdownDescription: "Contact email address for incident reports.",
			},
			"content_security_policy_header": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("upgrade-insecure-requests; default-src 'self'; frame-src 'self' https://telemetryservice.firstpartyapps.oaspapps.com/telemetryservice/telemetryproxy.html https://*.microsoft.com https://*.office.com; style-src 'self' 'unsafe-inline' https://*.microsoft.com https://*.office.com; object-src 'self'; font-src 'self' https://*.microsoft.com https://*.office.com; img-src 'self' https://www.adobe.com data: blob:; script-src 'self' 'unsafe-inline' 'unsafe-eval' https://*.microsoft.com https://*.office.com https://ajax.aspnetcdn.com https://api.keen.io; media-src 'self' blob:; connect-src 'self' https://*.microsoft.com https://*.office.com https://example.com; frame-ancestors 'self';"),
				MarkdownDescription: "HTTP Content-Security-Policy header contents.",
			},
			"content_security_policy_state": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
				MarkdownDescription: "Enable HTTP Content-Security-Policy.",
			},
			"crypto_mode": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("besteffort"),
				Validators: []validator.String{
					stringvalidator.OneOf("besteffort", "on", "off"),
				},
				MarkdownDescription: "Controls media encryption requirements.",
			},
			"default_theme": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Default theme for services.",
			},
			"default_to_new_webapp": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
				MarkdownDescription: "Deprecated field - use 'default_webapp' instead.",
			},
			"default_webapp": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("latest"),
				MarkdownDescription: "Deprecated field - use 'default_webapp_alias' instead.",
			},
			"default_webapp_alias": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Default web app path for conferencing nodes.",
			},
			// unique for each deployment, not update by users
			"deployment_uuid": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Deployment UUID.",
			},
			"disabled_codecs": schema.SetAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				//Default: setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{
				//	types.StringValue("MP4A-LATM_128"),
				//	types.StringValue("H264_H_0"),
				//	types.StringValue("H264_H_1"),
				//})),
				MarkdownDescription: "Codecs to disable.",
			},
			"eject_last_participant_backstop_timeout": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(0),
				MarkdownDescription: "Timeout for ejecting last participant.",
			},
			"enable_analytics": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Enable analytics collection.",
			},
			"enable_application_api": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
				MarkdownDescription: "Enable Infinity Client API.",
			},
			"enable_breakout_rooms": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Enable Breakout Rooms feature.",
			},
			"enable_chat": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
				MarkdownDescription: "Enable chat relay between participants.",
			},
			"enable_denoise": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
				MarkdownDescription: "Enable server-side denoising.",
			},
			"enable_dialout": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
				MarkdownDescription: "Enable dialout functionality.",
			},
			"enable_directory": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
				MarkdownDescription: "Enable directory for Infinity Connect clients.",
			},
			"enable_edge_non_mesh": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
				MarkdownDescription: "Enable restricted IPsec routing for Edge Nodes.",
			},
			"enable_fecc": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
				MarkdownDescription: "Enable Far-End Camera Control (FECC).",
			},
			"enable_h323": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
				MarkdownDescription: "Enable H.323 protocol.",
			},
			"enable_legacy_dialout_api": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Enable legacy dialout API.",
			},
			"enable_lync_auto_escalate": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Enable Lync auto escalate.",
			},
			"enable_lync_vbss": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Enable Lync VbSS.",
			},
			"enable_mlvad": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Enable advanced voice activity detection.",
			},
			"enable_multiscreen": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
				MarkdownDescription: "Enable dual screen layouts.",
			},
			"enable_push_notifications": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Enable push notifications (deprecated).",
			},
			"enable_rtmp": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
				MarkdownDescription: "Enable RTMP protocol.",
			},
			"enable_sip": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
				MarkdownDescription: "Enable SIP protocol.",
			},
			"enable_sip_udp": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Enable SIP UDP protocol.",
			},
			"enable_softmute": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
				MarkdownDescription: "Enable Softmute for audio gating.",
			},
			"enable_ssh": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
				MarkdownDescription: "Enable SSH access.",
			},
			"enable_turn_443": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Enable TURN on TCP port 443 for WebRTC.",
			},
			"enable_webrtc": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
				MarkdownDescription: "Enable WebRTC protocol.",
			},
			"error_reporting_enabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Enable error reporting.",
			},
			"error_reporting_url": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("https://acr.pexip.com"),
				MarkdownDescription: "URL for error reporting. Default https://acr.pexip.com",
			},
			"es_connection_timeout": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(7),
				MarkdownDescription: "Connection timeout for event sink.",
			},
			"es_initial_retry_backoff": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(1),
				MarkdownDescription: "Initial retry backoff for event sink.",
			},
			"es_maximum_deferred_posts": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(1000),
				MarkdownDescription: "Maximum deferred posts for event sink.",
			},
			"es_maximum_retry_backoff": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(1800),
				MarkdownDescription: "Maximum retry backoff for event sink.",
			},
			"es_media_streams_wait": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(1),
				MarkdownDescription: "Media streams wait time for event sink.",
			},
			"es_metrics_update_interval": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(60),
				MarkdownDescription: "Metrics update interval for event sink.",
			},
			"es_short_term_memory_expiration": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(2),
				MarkdownDescription: "Short term memory expiration for event sink.",
			},
			"external_participant_avatar_lookup": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
				MarkdownDescription: "Enable external participant avatar lookup.",
			},
			"gcp_client_email": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "GCP service account email.",
			},
			"gcp_private_key": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Sensitive:           false,
				Default:             stringdefault.StaticString(""),
				MarkdownDescription: "GCP service account private key.",
			},
			"gcp_project_id": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "GCP project ID.",
			},
			"guests_only_timeout": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(60),
				MarkdownDescription: "Timeout for guests-only conferences.",
			},
			"legacy_api_http": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Enable legacy API HTTP access.",
			},
			"legacy_api_username": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
				MarkdownDescription: "Legacy API username.",
			},
			"legacy_api_password": schema.StringAttribute{
				Optional:            true,
				//Sensitive:           true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
				MarkdownDescription: "Legacy API password.",
			},
			"live_captions_api_gateway": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
				MarkdownDescription: "Live captions API gateway.",
			},
			"live_captions_app_id": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
				MarkdownDescription: "Live captions App ID.",
			},
			"live_captions_enabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
				MarkdownDescription: "Enable live captions.",
			},
			"live_captions_public_jwt_key": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
				MarkdownDescription: "Live captions public JWT key.",
			},
			"live_captions_vmr_default": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Enable live captions by default for VMRs.",
			},
			"liveview_show_conferences": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
				MarkdownDescription: "Show conferences in Live View.",
			},
			"local_mssip_domain": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
				MarkdownDescription: "Local MSSIP domain.",
			},
			"logon_banner": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
				MarkdownDescription: "Logon banner text.",
			},
			"logs_max_age": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(0),
				MarkdownDescription: "Maximum age of logs (days).",
			},
			"management_qos": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(0),
				MarkdownDescription: "DSCP value for management traffic.",
			},
			"management_session_timeout": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(30),
				MarkdownDescription: "Session timeout for management interface (minutes).",
			},
			"management_start_page": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("/admin/conferencingstatus/deploymentgraph/deployment_graph/"),
				MarkdownDescription: "Start page for management web.",
			},
			"max_callrate_in": schema.Int64Attribute{
				Optional:            true,
				MarkdownDescription: "Maximum callrate in (kbps).",
			},
			"max_callrate_out": schema.Int64Attribute{
				Optional:            true,
				MarkdownDescription: "Maximum callrate out (kbps).",
			},
			"max_pixels_per_second": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("hd"),
				Validators: []validator.String{
					stringvalidator.OneOf("sd", "hd", "fullhd"),
				},
				MarkdownDescription: "Maximum pixels per second.",
			},
			"max_presentation_bandwidth_ratio": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(75),
				MarkdownDescription: "Maximum percentage of bandwidth for presentation.",
			},
			"media_ports_end": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(49999),
				MarkdownDescription: "End port for media traffic.",
			},
			"media_ports_start": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(40000),
				MarkdownDescription: "Start port for media traffic.",
			},
			"ocsp_responder_url": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
				MarkdownDescription: "OCSP responder URL.",
			},
			"ocsp_state": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("OFF"),
				Validators: []validator.String{
					stringvalidator.OneOf("OFF", "ON", "OVERRIDE"),
				},
				MarkdownDescription: "OCSP state.",
			},
			"pin_entry_timeout": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(120),
				MarkdownDescription: "Timeout for PIN entry (seconds).",
			},
			"pss_customer_id": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
				MarkdownDescription: "Pexip Private Cloud customer ID.",
			},
			"pss_enabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Enable Pexip Private Cloud connection.",
			},
			"pss_gateway": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
				MarkdownDescription: "Pexip Private Cloud gateway URL.",
			},
			"pss_token": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
				MarkdownDescription: "Pexip Private Cloud token.",
			},
			"resource_uri": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Resource URI for the global configuration.",
			},
			"session_timeout_enabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
				MarkdownDescription: "Enable session timeout for management interface.",
			},
			"signalling_ports_end": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(39999),
				MarkdownDescription: "End port for signalling traffic.",
			},
			"signalling_ports_start": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(34000),
				MarkdownDescription: "Start port for signalling traffic.",
			},
			"sip_tls_cert_verify_mode": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("OFF"),
				Validators: []validator.String{
					stringvalidator.OneOf("OFF", "ON"),
				},
				MarkdownDescription: "SIP TLS certificate verify mode.",
			},
			"site_banner": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
				MarkdownDescription: "Site banner text.",
			},
			"site_banner_bg": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("#c0c0c0"),
				MarkdownDescription: "Site banner background color.",
			},
			"site_banner_fg": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("#000000"),
				MarkdownDescription: "Site banner foreground color.",
			},
			"teams_enable_powerpoint_render": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Enable PowerPoint Live content in Teams calls.",
			},
			"waiting_for_chair_timeout": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(900),
				MarkdownDescription: "Timeout for waiting for chair (seconds).",
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
		DefaultWebapp:                       plan.DefaultWebapp.ValueString(),
		DeploymentUUID:                      plan.DeploymentUUID.ValueString(),
		ErrorReportingURL:                   plan.ErrorReportingURL.ValueString(),
		LegacyAPIUsername:                   plan.LegacyAPIUsername.ValueString(),
		LegacyAPIPassword:                   plan.LegacyAPIPassword.ValueString(),
		LiveCaptionsAPIGateway:              plan.LiveCaptionsAPIGateway.ValueString(),
		LiveCaptionsAppID:                   plan.LiveCaptionsAppID.ValueString(),
		LiveCaptionsPublicJWTKey:            plan.LiveCaptionsPublicJWTKey.ValueString(),
		LiveviewShowConferences:             plan.LiveviewShowConferences.ValueBool(),
		LocalMssipDomain:                    plan.LocalMssipDomain.ValueString(),
		LogonBanner:                         plan.LogonBanner.ValueString(),
		ManagementStartPage:                 plan.ManagementStartPage.ValueString(),
		MaxPixelsPerSecond:                  plan.MaxPixelsPerSecond.ValueString(),
		OcspResponderURL:                    plan.OcspResponderURL.ValueString(),
		OcspState:                           plan.OcspState.ValueString(),
		PssCustomerID:                       plan.PssCustomerID.ValueString(),
		PssGateway:                          plan.PssGateway.ValueString(),
		PssToken:                            plan.PssToken.ValueString(),
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
		EnableDenoise:                       plan.EnableDenoise.ValueBool(),
		EnableDialout:                       plan.EnableDialout.ValueBool(),
		EnableDirectory:                     plan.EnableDirectory.ValueBool(),
		EnableEdgeNonMesh:                   plan.EnableEdgeNonMesh.ValueBool(),
		EnableFecc:                          plan.EnableFecc.ValueBool(),
		EnableLegacyDialoutAPI:              plan.EnableLegacyDialoutAPI.ValueBool(),
		EnableLyncAutoEscalate:              plan.EnableLyncAutoEscalate.ValueBool(),
		EnableLyncVbss:                      plan.EnableLyncVbss.ValueBool(),
		EnableMlvad:                         plan.EnableMlvad.ValueBool(),
		EnableMultiscreen:                   plan.EnableMultiscreen.ValueBool(),
		EnablePushNotifications:             plan.EnablePushNotifications.ValueBool(),
		EnableSoftmute:                      plan.EnableSoftmute.ValueBool(),
		EnableSSH:                           plan.EnableSSH.ValueBool(),
		EnableTurn443:                       plan.EnableTurn443.ValueBool(),
		ErrorReportingEnabled:               plan.ErrorReportingEnabled.ValueBool(),
		EsConnectionTimeout:                 int(plan.EsConnectionTimeout.ValueInt64()),
		EsInitialRetryBackoff:               int(plan.EsInitialRetryBackoff.ValueInt64()),
		EsMaximumDeferredPosts:              int(plan.EsMaximumDeferredPosts.ValueInt64()),
		EsMaximumRetryBackoff:               int(plan.EsMaximumRetryBackoff.ValueInt64()),
		EsMediaStreamsWait:                  int(plan.EsMediaStreamsWait.ValueInt64()),
		EsMetricsUpdateInterval:             int(plan.EsMetricsUpdateInterval.ValueInt64()),
		EsShortTermMemoryExpiration:         int(plan.EsShortTermMemoryExpiration.ValueInt64()),
		ExternalParticipantAvatarLookup:     plan.ExternalParticipantAvatarLookup.ValueBool(),
		GuestsOnlyTimeout:                   int(plan.GuestsOnlyTimeout.ValueInt64()),
		LegacyAPIHTTP:                       plan.LegacyAPIHTTP.ValueBool(),
		LiveCaptionsEnabled:                 plan.LiveCaptionsEnabled.ValueBool(),
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
			disabledCodecs = append(disabledCodecs, config.CodecValue{Value: v.String()})
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
		updateRequest.DefaultTheme = &val
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
			"Error Creating Infinity global configuration",
			fmt.Sprintf("Could not create Infinity global configuration: %s", err),
		)
		return
	}

	// Read the current state from the API to get all computed values
	model, err := r.read(ctx, plan.AWSSecretKey.ValueStringPointer(), plan.AzureSecret.ValueStringPointer(), plan.GcpPrivateKey.ValueStringPointer(), plan.LegacyAPIPassword.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Created Infinity global configuration",
			fmt.Sprintf("Could not read created Infinity global configuration: %s", err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
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
	data.DefaultToNewWebapp = types.BoolValue(srv.DefaultToNewWebapp)
	data.DefaultWebapp = types.StringValue(srv.DefaultWebapp)
	data.DeploymentUUID = types.StringValue(srv.DeploymentUUID)
	data.ErrorReportingURL = types.StringValue(srv.ErrorReportingURL)
	data.LegacyAPIUsername = types.StringValue(srv.LegacyAPIUsername)
	data.LegacyAPIPassword = types.StringValue(legacyAPIPassword)
	data.LiveCaptionsAPIGateway = types.StringValue(srv.LiveCaptionsAPIGateway)
	data.LiveCaptionsAppID = types.StringValue(srv.LiveCaptionsAppID)
	data.LiveCaptionsPublicJWTKey = types.StringValue(srv.LiveCaptionsPublicJWTKey)
	data.LiveviewShowConferences = types.BoolValue(srv.LiveviewShowConferences)
	data.LocalMssipDomain = types.StringValue(srv.LocalMssipDomain)
	data.LogonBanner = types.StringValue(srv.LogonBanner)
	data.ManagementStartPage = types.StringValue(srv.ManagementStartPage)
	data.MaxPixelsPerSecond = types.StringValue(srv.MaxPixelsPerSecond)
	data.OcspResponderURL = types.StringValue(srv.OcspResponderURL)
	data.OcspState = types.StringValue(srv.OcspState)
	data.PssCustomerID = types.StringValue(srv.PssCustomerID)
	data.PssEnabled = types.BoolValue(srv.PssEnabled)
	data.PssGateway = types.StringValue(srv.PssGateway)
	data.PssToken = types.StringValue(srv.PssToken)
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
	data.EnableDenoise = types.BoolValue(srv.EnableDenoise)
	data.EnableDialout = types.BoolValue(srv.EnableDialout)
	data.EnableDirectory = types.BoolValue(srv.EnableDirectory)
	data.EnableEdgeNonMesh = types.BoolValue(srv.EnableEdgeNonMesh)
	data.EnableFecc = types.BoolValue(srv.EnableFecc)
	data.EnableLegacyDialoutAPI = types.BoolValue(srv.EnableLegacyDialoutAPI)
	data.EnableLyncAutoEscalate = types.BoolValue(srv.EnableLyncAutoEscalate)
	data.EnableLyncVbss = types.BoolValue(srv.EnableLyncVbss)
	data.EnableMlvad = types.BoolValue(srv.EnableMlvad)
	data.EnableMultiscreen = types.BoolValue(srv.EnableMultiscreen)
	data.EnablePushNotifications = types.BoolValue(srv.EnablePushNotifications)
	data.EnableSIPUDP = types.BoolValue(srv.EnableSIPUDP)
	data.EnableSoftmute = types.BoolValue(srv.EnableSoftmute)
	data.EnableSSH = types.BoolValue(srv.EnableSSH)
	data.EnableTurn443 = types.BoolValue(srv.EnableTurn443)
	data.ErrorReportingEnabled = types.BoolValue(srv.ErrorReportingEnabled)
	data.EsConnectionTimeout = types.Int64Value(int64(srv.EsConnectionTimeout))
	data.EsInitialRetryBackoff = types.Int64Value(int64(srv.EsInitialRetryBackoff))
	data.EsMaximumDeferredPosts = types.Int64Value(int64(srv.EsMaximumDeferredPosts))
	data.EsMaximumRetryBackoff = types.Int64Value(int64(srv.EsMaximumRetryBackoff))
	data.EsMediaStreamsWait = types.Int64Value(int64(srv.EsMediaStreamsWait))
	data.EsMetricsUpdateInterval = types.Int64Value(int64(srv.EsMetricsUpdateInterval))
	data.EsShortTermMemoryExpiration = types.Int64Value(int64(srv.EsShortTermMemoryExpiration))
	data.ExternalParticipantAvatarLookup = types.BoolValue(srv.ExternalParticipantAvatarLookup)
	data.GcpProjectID = types.StringPointerValue(srv.GcpProjectID)
	data.GcpClientEmail = types.StringPointerValue(srv.GcpClientEmail)
	data.GcpPrivateKey = types.StringPointerValue(gcpPrivateKey)
	data.GuestsOnlyTimeout = types.Int64Value(int64(srv.GuestsOnlyTimeout))
	data.LegacyAPIHTTP = types.BoolValue(srv.LegacyAPIHTTP)
	data.LiveCaptionsEnabled = types.BoolValue(srv.LiveCaptionsEnabled)
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

	if srv.ManagementQos != nil {
		data.ManagementQos = types.Int64Value(int64(*srv.ManagementQos))
	}
	//else {
	//	data.ManagementQos = types.Int64Null()
	//}

	if srv.BurstingMinLifetime != nil {
		data.BurstingMinLifetime = types.Int64Value(int64(*srv.BurstingMinLifetime))
	}
	if srv.BurstingThreshold != nil {
		data.BurstingThreshold = types.Int64Value(int64(*srv.BurstingThreshold))
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

	tflog.Debug(ctx, "Debug disabled codecs", map[string]interface{}{
		"disabled_codecs": updateRequest.DisabledCodecs,
	})

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
	// For singleton resources, delete means resetting to default values
	// We'll set minimal configuration to "delete" the customizations
	tflog.Info(ctx, "Deleting Infinity global configuration (resetting to defaults)")

	// only need to unset related fields
	updateRequest := &config.GlobalConfigurationUpdateRequest{
		DefaultTheme:       nil,
		DefaultWebappAlias: nil,
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
