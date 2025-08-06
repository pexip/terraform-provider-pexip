package provider

import (
	"context"
	"fmt"
	"strconv"

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
	_ resource.ResourceWithImportState = (*InfinityMsExchangeConnectorResource)(nil)
)

type InfinityMsExchangeConnectorResource struct {
	InfinityClient InfinityClient
}

type InfinityMsExchangeConnectorResourceModel struct {
	ID                                         types.String `tfsdk:"id"`
	ResourceID                                 types.Int32  `tfsdk:"resource_id"`
	Name                                       types.String `tfsdk:"name"`
	Description                                types.String `tfsdk:"description"`
	RoomMailboxEmailAddress                    types.String `tfsdk:"room_mailbox_email_address"`
	RoomMailboxName                            types.String `tfsdk:"room_mailbox_name"`
	URL                                        types.String `tfsdk:"url"`
	Username                                   types.String `tfsdk:"username"`
	Password                                   types.String `tfsdk:"password"`
	AuthenticationMethod                       types.String `tfsdk:"authentication_method"`
	AuthProvider                               types.String `tfsdk:"auth_provider"`
	UUID                                       types.String `tfsdk:"uuid"`
	ScheduledAliasPrefix                       types.String `tfsdk:"scheduled_alias_prefix"`
	ScheduledAliasDomain                       types.String `tfsdk:"scheduled_alias_domain"`
	ScheduledAliasSuffixLength                 types.Int64  `tfsdk:"scheduled_alias_suffix_length"`
	MeetingBufferBefore                        types.Int64  `tfsdk:"meeting_buffer_before"`
	MeetingBufferAfter                         types.Int64  `tfsdk:"meeting_buffer_after"`
	EnableDynamicVmrs                          types.Bool   `tfsdk:"enable_dynamic_vmrs"`
	EnablePersonalVmrs                         types.Bool   `tfsdk:"enable_personal_vmrs"`
	AllowNewUsers                              types.Bool   `tfsdk:"allow_new_users"`
	DisableProxy                               types.Bool   `tfsdk:"disable_proxy"`
	UseCustomAddInSources                      types.Bool   `tfsdk:"use_custom_add_in_sources"`
	EnableAddinDebugLogs                       types.Bool   `tfsdk:"enable_addin_debug_logs"`
	OauthClientID                              types.String `tfsdk:"oauth_client_id"`
	OauthClientSecret                          types.String `tfsdk:"oauth_client_secret"`
	OauthAuthEndpoint                          types.String `tfsdk:"oauth_auth_endpoint"`
	OauthTokenEndpoint                         types.String `tfsdk:"oauth_token_endpoint"`
	OauthRedirectURI                           types.String `tfsdk:"oauth_redirect_uri"`
	OauthRefreshToken                          types.String `tfsdk:"oauth_refresh_token"`
	OauthState                                 types.String `tfsdk:"oauth_state"`
	KerberosRealm                              types.String `tfsdk:"kerberos_realm"`
	KerberosKdc                                types.String `tfsdk:"kerberos_kdc"`
	KerberosKdcHttpsProxy                      types.String `tfsdk:"kerberos_kdc_https_proxy"`
	KerberosExchangeSpn                        types.String `tfsdk:"kerberos_exchange_spn"`
	KerberosEnableTls                          types.Bool   `tfsdk:"kerberos_enable_tls"`
	KerberosAuthEveryRequest                   types.Bool   `tfsdk:"kerberos_auth_every_request"`
	KerberosVerifyTlsUsingCustomCa             types.Bool   `tfsdk:"kerberos_verify_tls_using_custom_ca"`
	AddinServerDomain                          types.String `tfsdk:"addin_server_domain"`
	AddinDisplayName                           types.String `tfsdk:"addin_display_name"`
	AddinDescription                           types.String `tfsdk:"addin_description"`
	AddinProviderName                          types.String `tfsdk:"addin_provider_name"`
	AddinButtonLabel                           types.String `tfsdk:"addin_button_label"`
	AddinGroupLabel                            types.String `tfsdk:"addin_group_label"`
	AddinSupertipTitle                         types.String `tfsdk:"addin_supertip_title"`
	AddinSupertipDescription                   types.String `tfsdk:"addin_supertip_description"`
	AddinApplicationID                         types.String `tfsdk:"addin_application_id"`
	AddinAuthorityURL                          types.String `tfsdk:"addin_authority_url"`
	AddinOidcMetadataURL                       types.String `tfsdk:"addin_oidc_metadata_url"`
	AddinAuthenticationMethod                  types.String `tfsdk:"addin_authentication_method"`
	AddinNaaWebApiApplicationID                types.String `tfsdk:"addin_naa_web_api_application_id"`
	PersonalVmrOauthClientID                   types.String `tfsdk:"personal_vmr_oauth_client_id"`
	PersonalVmrOauthClientSecret               types.String `tfsdk:"personal_vmr_oauth_client_secret"`
	PersonalVmrOauthAuthEndpoint               types.String `tfsdk:"personal_vmr_oauth_auth_endpoint"`
	PersonalVmrOauthTokenEndpoint              types.String `tfsdk:"personal_vmr_oauth_token_endpoint"`
	PersonalVmrAdfsRelyingPartyTrustIdentifier types.String `tfsdk:"personal_vmr_adfs_relying_party_trust_identifier"`
	OfficeJsURL                                types.String `tfsdk:"office_js_url"`
	MicrosoftFabricURL                         types.String `tfsdk:"microsoft_fabric_url"`
	MicrosoftFabricComponentsURL               types.String `tfsdk:"microsoft_fabric_components_url"`
	AdditionalAddInScriptSources               types.String `tfsdk:"additional_add_in_script_sources"`
	Domains                                    types.String `tfsdk:"domains"`
	HostIdentityProviderGroup                  types.String `tfsdk:"host_identity_provider_group"`
	IvrTheme                                   types.String `tfsdk:"ivr_theme"`
	NonIdpParticipants                         types.String `tfsdk:"non_idp_participants"`
	PrivateKey                                 types.String `tfsdk:"private_key"`
	PublicKey                                  types.String `tfsdk:"public_key"`
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
				Validators: []validator.String{
					stringvalidator.LengthAtMost(500),
				},
				MarkdownDescription: "Description of the Microsoft Exchange connector. Maximum length: 500 characters.",
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
				Optional:            true,
				MarkdownDescription: "Username for Exchange authentication.",
			},
			"password": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				MarkdownDescription: "Password for Exchange authentication. This field is sensitive.",
			},
			"authentication_method": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf("ntlm", "basic", "kerberos", "oauth2"),
				},
				MarkdownDescription: "Authentication method for Exchange. Valid values: ntlm, basic, kerberos, oauth2.",
			},
			"auth_provider": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf("azure", "adfs", "exchange"),
				},
				MarkdownDescription: "Authentication provider. Valid values: azure, adfs, exchange.",
			},
			"uuid": schema.StringAttribute{
				Required:            true,
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
				Required:            true,
				MarkdownDescription: "Length of the suffix for scheduled conference aliases.",
			},
			"meeting_buffer_before": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "Buffer time before meetings in minutes.",
			},
			"meeting_buffer_after": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "Buffer time after meetings in minutes.",
			},
			"enable_dynamic_vmrs": schema.BoolAttribute{
				Required:            true,
				MarkdownDescription: "Whether to enable dynamic Virtual Meeting Rooms.",
			},
			"enable_personal_vmrs": schema.BoolAttribute{
				Required:            true,
				MarkdownDescription: "Whether to enable personal Virtual Meeting Rooms.",
			},
			"allow_new_users": schema.BoolAttribute{
				Required:            true,
				MarkdownDescription: "Whether to allow new users to be created.",
			},
			"disable_proxy": schema.BoolAttribute{
				Required:            true,
				MarkdownDescription: "Whether to disable proxy for connections.",
			},
			"use_custom_add_in_sources": schema.BoolAttribute{
				Required:            true,
				MarkdownDescription: "Whether to use custom add-in sources.",
			},
			"enable_addin_debug_logs": schema.BoolAttribute{
				Required:            true,
				MarkdownDescription: "Whether to enable debug logs for add-ins.",
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
				MarkdownDescription: "OAuth authorization endpoint URL.",
			},
			"oauth_token_endpoint": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "OAuth token endpoint URL.",
			},
			"oauth_redirect_uri": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "OAuth redirect URI for authentication flow.",
			},
			"oauth_refresh_token": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				MarkdownDescription: "OAuth refresh token. This field is sensitive.",
			},
			"oauth_state": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "OAuth state parameter for security.",
			},
			"kerberos_realm": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Kerberos realm for authentication.",
			},
			"kerberos_kdc": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Kerberos Key Distribution Center.",
			},
			"kerberos_kdc_https_proxy": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "HTTPS proxy for Kerberos KDC connections.",
			},
			"kerberos_exchange_spn": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Kerberos Service Principal Name for Exchange.",
			},
			"kerberos_enable_tls": schema.BoolAttribute{
				Required:            true,
				MarkdownDescription: "Whether to enable TLS for Kerberos connections.",
			},
			"kerberos_auth_every_request": schema.BoolAttribute{
				Required:            true,
				MarkdownDescription: "Whether to authenticate every Kerberos request.",
			},
			"kerberos_verify_tls_using_custom_ca": schema.BoolAttribute{
				Required:            true,
				MarkdownDescription: "Whether to verify TLS using custom CA for Kerberos.",
			},
			"addin_server_domain": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Server domain for Exchange add-in.",
			},
			"addin_display_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Display name for the Exchange add-in.",
			},
			"addin_description": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Description for the Exchange add-in.",
			},
			"addin_provider_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Provider name for the Exchange add-in.",
			},
			"addin_button_label": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Button label for the Exchange add-in.",
			},
			"addin_group_label": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Group label for the Exchange add-in.",
			},
			"addin_supertip_title": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Supertip title for the Exchange add-in.",
			},
			"addin_supertip_description": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Supertip description for the Exchange add-in.",
			},
			"addin_application_id": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Application ID for the Exchange add-in.",
			},
			"addin_authority_url": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Authority URL for add-in authentication.",
			},
			"addin_oidc_metadata_url": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "OIDC metadata URL for add-in authentication.",
			},
			"addin_authentication_method": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf("web_api", "naa"),
				},
				MarkdownDescription: "Authentication method for add-in. Valid values: web_api, naa.",
			},
			"addin_naa_web_api_application_id": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "NAA Web API application ID for add-in authentication.",
			},
			"personal_vmr_oauth_client_id": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "OAuth client ID for personal VMR integration.",
			},
			"personal_vmr_oauth_client_secret": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				MarkdownDescription: "OAuth client secret for personal VMR integration. This field is sensitive.",
			},
			"personal_vmr_oauth_auth_endpoint": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "OAuth authorization endpoint for personal VMR integration.",
			},
			"personal_vmr_oauth_token_endpoint": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "OAuth token endpoint for personal VMR integration.",
			},
			"personal_vmr_adfs_relying_party_trust_identifier": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "ADFS relying party trust identifier for personal VMR integration.",
			},
			"office_js_url": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Office.js library URL for add-in functionality.",
			},
			"microsoft_fabric_url": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Microsoft Fabric URL for add-in styling.",
			},
			"microsoft_fabric_components_url": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Microsoft Fabric Components URL for add-in styling.",
			},
			"additional_add_in_script_sources": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Additional script sources for add-in functionality.",
			},
			"domains": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "URI reference to associated domains resource.",
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
				Optional:            true,
				MarkdownDescription: "Configuration for non-IDP participants.",
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
		PersonalVmrAdfsRelyingPartyTrustIdentifier: plan.PersonalVmrAdfsRelyingPartyTrustIdentifier.ValueString(),
		OfficeJsURL:                  plan.OfficeJsURL.ValueString(),
		MicrosoftFabricURL:           plan.MicrosoftFabricURL.ValueString(),
		MicrosoftFabricComponentsURL: plan.MicrosoftFabricComponentsURL.ValueString(),
		AdditionalAddInScriptSources: plan.AdditionalAddInScriptSources.ValueString(),
		NonIdpParticipants:           plan.NonIdpParticipants.ValueString(),
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
		domains := plan.Domains.ValueString()
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
	data.Password = types.StringValue(srv.Password)
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
	data.OauthClientSecret = types.StringValue(srv.OauthClientSecret)
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
	data.PersonalVmrOauthClientSecret = types.StringValue(srv.PersonalVmrOauthClientSecret)
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
		data.Domains = types.StringValue(*srv.Domains)
	} else {
		data.Domains = types.StringNull()
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

	return &data, nil
}

func (r *InfinityMsExchangeConnectorResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	state := &InfinityMsExchangeConnectorResourceModel{}

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
			"Error Reading Infinity Microsoft Exchange connector",
			fmt.Sprintf("Could not read Infinity Microsoft Exchange connector: %s", err),
		)
		return
	}

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
		Name:                          plan.Name.ValueString(),
		Description:                   plan.Description.ValueString(),
		RoomMailboxName:               plan.RoomMailboxName.ValueString(),
		URL:                           plan.URL.ValueString(),
		Username:                      plan.Username.ValueString(),
		Password:                      plan.Password.ValueString(),
		AuthenticationMethod:          plan.AuthenticationMethod.ValueString(),
		AuthProvider:                  plan.AuthProvider.ValueString(),
		UUID:                          plan.UUID.ValueString(),
		ScheduledAliasDomain:          plan.ScheduledAliasDomain.ValueString(),
		OauthClientSecret:             plan.OauthClientSecret.ValueString(),
		OauthAuthEndpoint:             plan.OauthAuthEndpoint.ValueString(),
		OauthTokenEndpoint:            plan.OauthTokenEndpoint.ValueString(),
		OauthRedirectURI:              plan.OauthRedirectURI.ValueString(),
		OauthRefreshToken:             plan.OauthRefreshToken.ValueString(),
		KerberosRealm:                 plan.KerberosRealm.ValueString(),
		KerberosKdc:                   plan.KerberosKdc.ValueString(),
		KerberosKdcHttpsProxy:         plan.KerberosKdcHttpsProxy.ValueString(),
		KerberosExchangeSpn:           plan.KerberosExchangeSpn.ValueString(),
		AddinServerDomain:             plan.AddinServerDomain.ValueString(),
		AddinDisplayName:              plan.AddinDisplayName.ValueString(),
		AddinDescription:              plan.AddinDescription.ValueString(),
		AddinProviderName:             plan.AddinProviderName.ValueString(),
		AddinButtonLabel:              plan.AddinButtonLabel.ValueString(),
		AddinGroupLabel:               plan.AddinGroupLabel.ValueString(),
		AddinSupertipTitle:            plan.AddinSupertipTitle.ValueString(),
		AddinSupertipDescription:      plan.AddinSupertipDescription.ValueString(),
		AddinAuthorityURL:             plan.AddinAuthorityURL.ValueString(),
		AddinOidcMetadataURL:          plan.AddinOidcMetadataURL.ValueString(),
		AddinAuthenticationMethod:     plan.AddinAuthenticationMethod.ValueString(),
		PersonalVmrOauthClientSecret:  plan.PersonalVmrOauthClientSecret.ValueString(),
		PersonalVmrOauthAuthEndpoint:  plan.PersonalVmrOauthAuthEndpoint.ValueString(),
		PersonalVmrOauthTokenEndpoint: plan.PersonalVmrOauthTokenEndpoint.ValueString(),
		PersonalVmrAdfsRelyingPartyTrustIdentifier: plan.PersonalVmrAdfsRelyingPartyTrustIdentifier.ValueString(),
		OfficeJsURL:                  plan.OfficeJsURL.ValueString(),
		MicrosoftFabricURL:           plan.MicrosoftFabricURL.ValueString(),
		MicrosoftFabricComponentsURL: plan.MicrosoftFabricComponentsURL.ValueString(),
		AdditionalAddInScriptSources: plan.AdditionalAddInScriptSources.ValueString(),
		NonIdpParticipants:           plan.NonIdpParticipants.ValueString(),
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
		domains := plan.Domains.ValueString()
		updateRequest.Domains = &domains
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
