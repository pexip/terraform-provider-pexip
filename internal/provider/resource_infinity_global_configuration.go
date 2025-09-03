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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
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
				Sensitive:           true,
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
				Default:             "",
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
				MarkdownDescription: "Deployment UUID.",
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
				Default:             int64default.StaticInt64(0),
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
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("hd"),
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
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("OFF"),
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
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("OFF"),
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
			"global_conference_create_groups": schema.SetAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				MarkdownDescription: "Groups that can create conferences globally.",
			},
		},
		MarkdownDescription: "Manages the global system configuration with the Infinity service. This is a singleton resource - only one global configuration exists per system.",
	}
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
	model, err := r.read(ctx, plan.AWSSecretKey.ValueString(), plan.AzureSecret.ValueString(), plan.GcpPrivateKey.ValueString(), plan.LegacyAPIPassword.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Created Infinity global configuration",
			fmt.Sprintf("Could not read created Infinity global configuration: %s", err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *InfinityGlobalConfigurationResource) buildUpdateRequest(plan *InfinityGlobalConfigurationResourceModel) *config.GlobalConfigurationUpdateRequest {
	updateRequest := &config.GlobalConfigurationUpdateRequest{
		CloudProvider:               plan.CloudProvider.ValueString(),
		ContactEmailAddress:         plan.ContactEmailAddress.ValueString(),
		ContentSecurityPolicyHeader: plan.ContentSecurityPolicyHeader.ValueString(),
		CryptoMode:                  plan.CryptoMode.ValueString(),
		DefaultWebapp:               plan.DefaultWebapp.ValueString(),
		DeploymentUUID:              plan.DeploymentUUID.ValueString(),
		ErrorReportingURL:           plan.ErrorReportingURL.ValueString(),
		LegacyAPIUsername:           plan.LegacyAPIUsername.ValueString(),
		LegacyAPIPassword:           plan.LegacyAPIPassword.ValueString(),
		LiveCaptionsAPIGateway:      plan.LiveCaptionsAPIGateway.ValueString(),
		LiveCaptionsAppID:           plan.LiveCaptionsAppID.ValueString(),
		LiveCaptionsPublicJWTKey:    plan.LiveCaptionsPublicJWTKey.ValueString(),
		LocalMssipDomain:            plan.LocalMssipDomain.ValueString(),
		LogonBanner:                 plan.LogonBanner.ValueString(),
		ManagementStartPage:         plan.ManagementStartPage.ValueString(),
		MaxPixelsPerSecond:          plan.MaxPixelsPerSecond.ValueString(),
		OcspResponderURL:            plan.OcspResponderURL.ValueString(),
		OcspState:                   plan.OcspState.ValueString(),
		PssCustomerID:               plan.PssCustomerID.ValueString(),
		PssGateway:                  plan.PssGateway.ValueString(),
		PssToken:                    plan.PssToken.ValueString(),
		SipTLSCertVerifyMode:        plan.SipTLSCertVerifyMode.ValueString(),
		SiteBanner:                  plan.SiteBanner.ValueString(),
		SiteBannerBg:                plan.SiteBannerBg.ValueString(),
		SiteBannerFg:                plan.SiteBannerFg.ValueString(),
	}

	// handle pointers
	if !plan.AWSSecretKey.IsNull() {
		val := plan.AWSSecretKey.ValueString()
		updateRequest.AWSSecretKey = &val
	}
	if !plan.AzureClientID.IsNull() {
		val := plan.AzureClientID.ValueString()
		updateRequest.AzureClientID = &val
	}
	if !plan.AzureSecret.IsNull() {
		val := plan.AzureSecret.ValueString()
		updateRequest.AzureSecret = &val
	}
	if !plan.AzureSubscriptionID.IsNull() {
		val := plan.AzureSubscriptionID.ValueString()
		updateRequest.AzureSubscriptionID = &val
	}
	if !plan.AzureTenant.IsNull() {
		val := plan.AzureTenant.ValueString()
		updateRequest.AzureTenant = &val
	}
	if !plan.BdpmMaxPinFailuresPerWindow.IsNull() {
		val := int(plan.BdpmMaxPinFailuresPerWindow.ValueInt64())
		updateRequest.BdpmMaxPinFailuresPerWindow = &val
	}
	if !plan.BdpmMaxScanAttemptsPerWindow.IsNull() {
		val := int(plan.BdpmMaxScanAttemptsPerWindow.ValueInt64())
		updateRequest.BdpmMaxScanAttemptsPerWindow = &val
	}
	if !plan.BdpmPinChecksEnabled.IsNull() {
		val := plan.BdpmPinChecksEnabled.ValueBool()
		updateRequest.BdpmPinChecksEnabled = &val
	}
	if !plan.BdpmScanQuarantineEnabled.IsNull() {
		val := plan.BdpmScanQuarantineEnabled.ValueBool()
		updateRequest.BdpmScanQuarantineEnabled = &val
	}
	if !plan.BurstingEnabled.IsNull() {
		val := plan.BurstingEnabled.ValueBool()
		updateRequest.BurstingEnabled = &val
	}
	if !plan.BurstingMinLifetime.IsNull() {
		val := int(plan.BurstingMinLifetime.ValueInt64())
		updateRequest.BurstingMinLifetime = &val
	}
	if !plan.BurstingThreshold.IsNull() {
		val := int(plan.BurstingThreshold.ValueInt64())
		updateRequest.BurstingThreshold = &val
	}
	if !plan.ContentSecurityPolicyState.IsNull() {
		val := plan.ContentSecurityPolicyState.ValueBool()
		updateRequest.ContentSecurityPolicyState = &val
	}
	if !plan.DefaultTheme.IsNull() {
		val := plan.DefaultTheme.ValueString()
		updateRequest.DefaultTheme = &val
	}
	if !plan.DefaultToNewWebapp.IsNull() {
		val := plan.DefaultToNewWebapp.ValueBool()
		updateRequest.DefaultToNewWebapp = &val
	}
	if !plan.DefaultWebappAlias.IsNull() {
		val := plan.DefaultWebappAlias.ValueString()
		updateRequest.DefaultWebappAlias = &val
	}
	if !plan.DisabledCodecs.IsNull() {
		var disabledCodecs []config.CodecValue
		for _, v := range plan.DisabledCodecs.Elements() {
			disabledCodecs = append(disabledCodecs, config.CodecValue{Value: v.String()})
		}
		updateRequest.DisabledCodecs = disabledCodecs
	}
	if !plan.EjectLastParticipantBackstopTimeout.IsNull() {
		val := int(plan.EjectLastParticipantBackstopTimeout.ValueInt64())
		updateRequest.EjectLastParticipantBackstopTimeout = &val
	}
	if !plan.EnableAnalytics.IsNull() {
		val := plan.EnableAnalytics.ValueBool()
		updateRequest.EnableAnalytics = &val
	}
	if !plan.EnableApplicationAPI.IsNull() {
		val := plan.EnableApplicationAPI.ValueBool()
		updateRequest.EnableApplicationAPI = &val
	}
	if !plan.EnableBreakoutRooms.IsNull() {
		val := plan.EnableBreakoutRooms.ValueBool()
		updateRequest.EnableBreakoutRooms = &val
	}
	if !plan.EnableChat.IsNull() {
		val := plan.EnableChat.ValueBool()
		updateRequest.EnableChat = &val
	}
	if !plan.EnableDenoise.IsNull() {
		val := plan.EnableDenoise.ValueBool()
		updateRequest.EnableDenoise = &val
	}
	if !plan.EnableDialout.IsNull() {
		val := plan.EnableDialout.ValueBool()
		updateRequest.EnableDialout = &val
	}
	if !plan.EnableDirectory.IsNull() {
		val := plan.EnableDirectory.ValueBool()
		updateRequest.EnableDirectory = &val
	}
	if !plan.EnableEdgeNonMesh.IsNull() {
		val := plan.EnableEdgeNonMesh.ValueBool()
		updateRequest.EnableEdgeNonMesh = &val
	}
	if !plan.EnableFecc.IsNull() {
		val := plan.EnableFecc.ValueBool()
		updateRequest.EnableFecc = &val
	}
	if !plan.EnableH323.IsNull() {
		val := plan.EnableH323.ValueBool()
		updateRequest.EnableH323 = &val
	}
	if !plan.EnableLegacyDialoutAPI.IsNull() {
		val := plan.EnableLegacyDialoutAPI.ValueBool()
		updateRequest.EnableLegacyDialoutAPI = &val
	}
	if !plan.EnableLyncAutoEscalate.IsNull() {
		val := plan.EnableLyncAutoEscalate.ValueBool()
		updateRequest.EnableLyncAutoEscalate = &val
	}
	if !plan.EnableLyncVbss.IsNull() {
		val := plan.EnableLyncVbss.ValueBool()
		updateRequest.EnableLyncVbss = &val
	}
	if !plan.EnableMlvad.IsNull() {
		val := plan.EnableMlvad.ValueBool()
		updateRequest.EnableMlvad = &val
	}
	if !plan.EnableMultiscreen.IsNull() {
		val := plan.EnableMultiscreen.ValueBool()
		updateRequest.EnableMultiscreen = &val
	}
	if !plan.EnablePushNotifications.IsNull() {
		val := plan.EnablePushNotifications.ValueBool()
		updateRequest.EnablePushNotifications = &val
	}
	if !plan.EnableRTMP.IsNull() {
		val := plan.EnableRTMP.ValueBool()
		updateRequest.EnableRTMP = &val
	}
	if !plan.EnableSIP.IsNull() {
		val := plan.EnableSIP.ValueBool()
		updateRequest.EnableSIP = &val
	}
	if !plan.EnableSIPUDP.IsNull() {
		val := plan.EnableSIPUDP.ValueBool()
		updateRequest.EnableSIPUDP = &val
	}
	if !plan.EnableSoftmute.IsNull() {
		val := plan.EnableSoftmute.ValueBool()
		updateRequest.EnableSoftmute = &val
	}
	if !plan.EnableSSH.IsNull() {
		val := plan.EnableSSH.ValueBool()
		updateRequest.EnableSSH = &val
	}
	if !plan.EnableTurn443.IsNull() {
		val := plan.EnableTurn443.ValueBool()
		updateRequest.EnableTurn443 = &val
	}
	if !plan.EnableWebRTC.IsNull() {
		val := plan.EnableWebRTC.ValueBool()
		updateRequest.EnableWebRTC = &val
	}
	if !plan.ErrorReportingEnabled.IsNull() {
		val := plan.ErrorReportingEnabled.ValueBool()
		updateRequest.ErrorReportingEnabled = &val
	}
	if !plan.EsConnectionTimeout.IsNull() {
		val := int(plan.EsConnectionTimeout.ValueInt64())
		updateRequest.EsConnectionTimeout = &val
	}
	if !plan.EsInitialRetryBackoff.IsNull() {
		val := int(plan.EsInitialRetryBackoff.ValueInt64())
		updateRequest.EsInitialRetryBackoff = &val
	}
	if !plan.EsMaximumDeferredPosts.IsNull() {
		val := int(plan.EsMaximumDeferredPosts.ValueInt64())
		updateRequest.EsMaximumDeferredPosts = &val
	}
	if !plan.EsMaximumRetryBackoff.IsNull() {
		val := int(plan.EsMaximumRetryBackoff.ValueInt64())
		updateRequest.EsMaximumRetryBackoff = &val
	}
	if !plan.EsMediaStreamsWait.IsNull() {
		val := int(plan.EsMediaStreamsWait.ValueInt64())
		updateRequest.EsMediaStreamsWait = &val
	}
	if !plan.EsMetricsUpdateInterval.IsNull() {
		val := int(plan.EsMetricsUpdateInterval.ValueInt64())
		updateRequest.EsMetricsUpdateInterval = &val
	}
	if !plan.EsShortTermMemoryExpiration.IsNull() {
		val := int(plan.EsShortTermMemoryExpiration.ValueInt64())
		updateRequest.EsShortTermMemoryExpiration = &val
	}
	if !plan.ExternalParticipantAvatarLookup.IsNull() {
		val := plan.ExternalParticipantAvatarLookup.ValueBool()
		updateRequest.ExternalParticipantAvatarLookup = &val
	}
	if !plan.GcpClientEmail.IsNull() {
		val := plan.GcpClientEmail.ValueString()
		updateRequest.GcpClientEmail = &val
	}
	if !plan.GcpPrivateKey.IsNull() {
		val := plan.GcpPrivateKey.ValueString()
		updateRequest.GcpPrivateKey = &val
	}
	if !plan.GcpProjectID.IsNull() {
		val := plan.GcpProjectID.ValueString()
		updateRequest.GcpProjectID = &val
	}
	if !plan.GuestsOnlyTimeout.IsNull() {
		val := int(plan.GuestsOnlyTimeout.ValueInt64())
		updateRequest.GuestsOnlyTimeout = &val
	}
	if !plan.LegacyAPIHTTP.IsNull() {
		val := plan.LegacyAPIHTTP.ValueBool()
		updateRequest.LegacyAPIHTTP = &val
	}
	if !plan.LiveCaptionsEnabled.IsNull() {
		val := plan.LiveCaptionsEnabled.ValueBool()
		updateRequest.LiveCaptionsEnabled = &val
	}
	if !plan.LiveCaptionsVMRDefault.IsNull() {
		val := plan.LiveCaptionsVMRDefault.ValueBool()
		updateRequest.LiveCaptionsVMRDefault = &val
	}
	if !plan.LiveviewShowConferences.IsNull() {
		val := plan.LiveviewShowConferences.ValueBool()
		updateRequest.LiveviewShowConferences = &val
	}
	if !plan.LogsMaxAge.IsNull() {
		val := int(plan.LogsMaxAge.ValueInt64())
		updateRequest.LogsMaxAge = &val
	}
	if !plan.ManagementQos.IsNull() {
		val := int(plan.ManagementQos.ValueInt64())
		updateRequest.ManagementQos = &val
	}
	if !plan.ManagementSessionTimeout.IsNull() {
		val := int(plan.ManagementSessionTimeout.ValueInt64())
		updateRequest.ManagementSessionTimeout = &val
	}
	if !plan.MaxCallrateIn.IsNull() {
		val := int(plan.MaxCallrateIn.ValueInt64())
		updateRequest.MaxCallrateIn = &val
	}
	if !plan.MaxCallrateOut.IsNull() {
		val := int(plan.MaxCallrateOut.ValueInt64())
		updateRequest.MaxCallrateOut = &val
	}
	if !plan.MaxPresentationBandwidthRatio.IsNull() {
		val := int(plan.MaxPresentationBandwidthRatio.ValueInt64())
		updateRequest.MaxPresentationBandwidthRatio = &val
	}
	if !plan.MediaPortsEnd.IsNull() {
		val := int(plan.MediaPortsEnd.ValueInt64())
		updateRequest.MediaPortsEnd = &val
	}
	if !plan.MediaPortsStart.IsNull() {
		val := int(plan.MediaPortsStart.ValueInt64())
		updateRequest.MediaPortsStart = &val
	}
	if !plan.PinEntryTimeout.IsNull() {
		val := int(plan.PinEntryTimeout.ValueInt64())
		updateRequest.PinEntryTimeout = &val
	}
	if !plan.PssEnabled.IsNull() {
		val := plan.PssEnabled.ValueBool()
		updateRequest.PssEnabled = &val
	}
	if !plan.SessionTimeoutEnabled.IsNull() {
		val := plan.SessionTimeoutEnabled.ValueBool()
		updateRequest.SessionTimeoutEnabled = &val
	}
	if !plan.SignallingPortsEnd.IsNull() {
		val := int(plan.SignallingPortsEnd.ValueInt64())
		updateRequest.SignallingPortsEnd = &val
	}
	if !plan.SignallingPortsStart.IsNull() {
		val := int(plan.SignallingPortsStart.ValueInt64())
		updateRequest.SignallingPortsStart = &val
	}
	if !plan.TeamsEnablePowerpointRender.IsNull() {
		val := plan.TeamsEnablePowerpointRender.ValueBool()
		updateRequest.TeamsEnablePowerpointRender = &val
	}
	if !plan.WaitingForChairTimeout.IsNull() {
		val := int(plan.WaitingForChairTimeout.ValueInt64())
		updateRequest.WaitingForChairTimeout = &val
	}

	return updateRequest
}

func (r *InfinityGlobalConfigurationResource) read(ctx context.Context, awsSecretKey, azureSecret, gcpPrivateKey, legacyAPIPassword string) (*InfinityGlobalConfigurationResourceModel, error) {
	var data InfinityGlobalConfigurationResourceModel

	srv, err := r.InfinityClient.Config().GetGlobalConfiguration(ctx)
	if err != nil {
		return nil, err
	}

	if srv.ResourceURI == "" {
		return nil, fmt.Errorf("global configuration not found")
	}

	data.ID = types.StringValue(srv.ResourceURI)
	data.EnableWebRTC = types.BoolValue(srv.EnableWebRTC)
	data.EnableSIP = types.BoolValue(srv.EnableSIP)
	data.EnableH323 = types.BoolValue(srv.EnableH323)
	data.EnableRTMP = types.BoolValue(srv.EnableRTMP)
	data.CryptoMode = types.StringValue(srv.CryptoMode)
	data.MaxPixelsPerSecond = types.StringValue(srv.MaxPixelsPerSecond)
	data.MediaPortsStart = types.Int64Value(int64(srv.MediaPortsStart))
	data.MediaPortsEnd = types.Int64Value(int64(srv.MediaPortsEnd))
	data.SignallingPortsStart = types.Int64Value(int64(srv.SignallingPortsStart))
	data.SignallingPortsEnd = types.Int64Value(int64(srv.SignallingPortsEnd))
	data.BurstingEnabled = types.BoolValue(srv.BurstingEnabled)
	data.CloudProvider = types.StringValue(srv.CloudProvider)
	data.GuestsOnlyTimeout = types.Int64Value(int64(srv.GuestsOnlyTimeout))
	data.WaitingForChairTimeout = types.Int64Value(int64(srv.WaitingForChairTimeout))
	data.EnableAnalytics = types.BoolValue(srv.EnableAnalytics)
	data.ErrorReportingEnabled = types.BoolValue(srv.ErrorReportingEnabled)
	data.ContactEmailAddress = types.StringValue(srv.ContactEmailAddress)

	// Handle optional pointer fields
	if srv.AWSAccessKey != nil {
		data.AWSAccessKey = types.StringValue(*srv.AWSAccessKey)
	} else {
		data.AWSAccessKey = types.StringNull()
	}

	if srv.AWSSecretKey != nil {
		data.AWSSecretKey = types.StringValue(*srv.AWSSecretKey)
	} else {
		data.AWSSecretKey = types.StringNull()
	}

	if srv.AzureClientID != nil {
		data.AzureClientID = types.StringValue(*srv.AzureClientID)
	} else {
		data.AzureClientID = types.StringNull()
	}

	if srv.AzureSecret != nil {
		data.AzureSecret = types.StringValue(*srv.AzureSecret)
	} else {
		data.AzureSecret = types.StringNull()
	}

	return &data, nil
}

func (r *InfinityGlobalConfigurationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	state := &InfinityGlobalConfigurationResourceModel{}

	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	state, err := r.read(ctx, state.AWSSecretKey.ValueString(), state.AzureSecret.ValueString(), state.GcpPrivateKey.ValueString(), state.LegacyAPIPassword.ValueString())
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

	updateRequest := &config.GlobalConfigurationUpdateRequest{
		CryptoMode:          plan.CryptoMode.ValueString(),
		MaxPixelsPerSecond:  plan.MaxPixelsPerSecond.ValueString(),
		CloudProvider:       plan.CloudProvider.ValueString(),
		ContactEmailAddress: plan.ContactEmailAddress.ValueString(),
	}

	// Handle boolean fields
	if !plan.EnableWebRTC.IsNull() {
		enable := plan.EnableWebRTC.ValueBool()
		updateRequest.EnableWebRTC = &enable
	}

	if !plan.EnableSIP.IsNull() {
		enable := plan.EnableSIP.ValueBool()
		updateRequest.EnableSIP = &enable
	}

	if !plan.EnableH323.IsNull() {
		enable := plan.EnableH323.ValueBool()
		updateRequest.EnableH323 = &enable
	}

	if !plan.EnableRTMP.IsNull() {
		enable := plan.EnableRTMP.ValueBool()
		updateRequest.EnableRTMP = &enable
	}

	if !plan.BurstingEnabled.IsNull() {
		enable := plan.BurstingEnabled.ValueBool()
		updateRequest.BurstingEnabled = &enable
	}

	if !plan.EnableAnalytics.IsNull() {
		enable := plan.EnableAnalytics.ValueBool()
		updateRequest.EnableAnalytics = &enable
	}

	if !plan.ErrorReportingEnabled.IsNull() {
		enable := plan.ErrorReportingEnabled.ValueBool()
		updateRequest.ErrorReportingEnabled = &enable
	}

	// Handle integer fields
	if !plan.MediaPortsStart.IsNull() {
		port := int(plan.MediaPortsStart.ValueInt64())
		updateRequest.MediaPortsStart = &port
	}

	if !plan.MediaPortsEnd.IsNull() {
		port := int(plan.MediaPortsEnd.ValueInt64())
		updateRequest.MediaPortsEnd = &port
	}

	if !plan.SignallingPortsStart.IsNull() {
		port := int(plan.SignallingPortsStart.ValueInt64())
		updateRequest.SignallingPortsStart = &port
	}

	if !plan.SignallingPortsEnd.IsNull() {
		port := int(plan.SignallingPortsEnd.ValueInt64())
		updateRequest.SignallingPortsEnd = &port
	}

	if !plan.GuestsOnlyTimeout.IsNull() {
		timeout := int(plan.GuestsOnlyTimeout.ValueInt64())
		updateRequest.GuestsOnlyTimeout = &timeout
	}

	if !plan.WaitingForChairTimeout.IsNull() {
		timeout := int(plan.WaitingForChairTimeout.ValueInt64())
		updateRequest.WaitingForChairTimeout = &timeout
	}

	// Handle sensitive string fields
	if !plan.AWSAccessKey.IsNull() {
		key := plan.AWSAccessKey.ValueString()
		updateRequest.AWSAccessKey = &key
	}

	if !plan.AWSSecretKey.IsNull() {
		secret := plan.AWSSecretKey.ValueString()
		updateRequest.AWSSecretKey = &secret
	}

	if !plan.AzureClientID.IsNull() {
		clientID := plan.AzureClientID.ValueString()
		updateRequest.AzureClientID = &clientID
	}

	if !plan.AzureSecret.IsNull() {
		secret := plan.AzureSecret.ValueString()
		updateRequest.AzureSecret = &secret
	}

	_, err := r.InfinityClient.Config().UpdateGlobalConfiguration(ctx, updateRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Infinity global configuration",
			fmt.Sprintf("Could not update Infinity global configuration: %s", err),
		)
		return
	}

	// Re-read the resource to get the latest state
	updatedModel, err := r.read(ctx, plan.AWSSecretKey.ValueString(), plan.AzureSecret.ValueString(), plan.GcpPrivateKey.ValueString(), plan.LegacyAPIPassword.ValueString())
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

	updateRequest := &config.GlobalConfigurationUpdateRequest{
		EnableWebRTC:    func() *bool { v := false; return &v }(),
		EnableSIP:       func() *bool { v := false; return &v }(),
		EnableH323:      func() *bool { v := false; return &v }(),
		EnableRTMP:      func() *bool { v := false; return &v }(),
		CryptoMode:      "disabled",
		CloudProvider:   "",
		BurstingEnabled: func() *bool { v := false; return &v }(),
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
	model, err := r.read(ctx, "", "", "", "")
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
