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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/pexip/go-infinity-sdk/v38/config"
)

var (
	_ resource.ResourceWithImportState = (*InfinityMjxExchangeDeploymentResource)(nil)
)

type InfinityMjxExchangeDeploymentResource struct {
	InfinityClient InfinityClient
}

type InfinityMjxExchangeDeploymentResourceModel struct {
	ID                             types.String `tfsdk:"id"`
	ResourceID                     types.Int32  `tfsdk:"resource_id"`
	Name                           types.String `tfsdk:"name"`
	Description                    types.String `tfsdk:"description"`
	ServiceAccountUsername         types.String `tfsdk:"service_account_username"`
	ServiceAccountPassword         types.String `tfsdk:"service_account_password"`
	AuthenticationMethod           types.String `tfsdk:"authentication_method"`
	EWSURL                         types.String `tfsdk:"ews_url"`
	DisableProxy                   types.Bool   `tfsdk:"disable_proxy"`
	FindItemsRequestQuota          types.Int64  `tfsdk:"find_items_request_quota"`
	KerberosRealm                  types.String `tfsdk:"kerberos_realm"`
	KerberosKDC                    types.String `tfsdk:"kerberos_kdc"`
	KerberosExchangeSPN            types.String `tfsdk:"kerberos_exchange_spn"`
	KerberosAuthEveryRequest       types.Bool   `tfsdk:"kerberos_auth_every_request"`
	KerberosEnableTLS              types.Bool   `tfsdk:"kerberos_enable_tls"`
	KerberosKDCHTTPSProxy          types.String `tfsdk:"kerberos_kdc_https_proxy"`
	KerberosVerifyTLSUsingCustomCA types.Bool   `tfsdk:"kerberos_verify_tls_using_custom_ca"`
	OAuthClientID                  types.String `tfsdk:"oauth_client_id"`
	OAuthAuthEndpoint              types.String `tfsdk:"oauth_auth_endpoint"`
	OAuthTokenEndpoint             types.String `tfsdk:"oauth_token_endpoint"`
	OAuthRedirectURI               types.String `tfsdk:"oauth_redirect_uri"`
	OAuthRefreshToken              types.String `tfsdk:"oauth_refresh_token"`
	OAuthState                     types.String `tfsdk:"oauth_state"`
	AutodiscoverURLs               types.Set    `tfsdk:"autodiscover_urls"`
	MjxIntegrations                types.Set    `tfsdk:"mjx_integrations"`
}

func (r *InfinityMjxExchangeDeploymentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_infinity_mjx_exchange_deployment"
}

func (r *InfinityMjxExchangeDeploymentResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *InfinityMjxExchangeDeploymentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Resource URI for the MJX Exchange deployment in Infinity.",
			},
			"resource_id": schema.Int32Attribute{
				Computed:            true,
				MarkdownDescription: "The resource integer identifier for the MJX Exchange deployment in Infinity.",
			},
			"name": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					stringvalidator.LengthAtMost(250),
				},
				MarkdownDescription: "The name of this One-Touch Join Exchange Integration. Maximum length: 250 characters.",
			},
			"description": schema.StringAttribute{
				Computed: true,
				Optional: true,
				Default:  stringdefault.StaticString(""),
				Validators: []validator.String{
					stringvalidator.LengthAtMost(250),
				},
				MarkdownDescription: "An optional description of this One-Touch Join Exchange Integration. Maximum length: 250 characters.",
			},
			"service_account_username": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(100),
				},
				MarkdownDescription: "The username of the service account to be used by the One-Touch Join Exchange Integration. Maximum length: 100 characters.",
			},
			"service_account_password": schema.StringAttribute{
				Required:  true,
				Sensitive: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(100),
				},
				MarkdownDescription: "The password of the service account to be used by the One-Touch Join Exchange Integration. Maximum length: 100 characters.",
			},
			"authentication_method": schema.StringAttribute{
				Computed: true,
				Optional: true,
				Default:  stringdefault.StaticString("BASIC"),
				Validators: []validator.String{
					stringvalidator.OneOf("BASIC", "NTLM", "KERBEROS", "OAUTH"),
				},
				MarkdownDescription: "The method used to authenticate to Exchange. Valid values: `BASIC`, `NTLM`, `KERBEROS`, `OAUTH`. Default: `BASIC`.",
			},
			"ews_url": schema.StringAttribute{
				Computed: true,
				Optional: true,
				Default:  stringdefault.StaticString(""),
				Validators: []validator.String{
					stringvalidator.LengthAtMost(255),
				},
				MarkdownDescription: "The URL used to connect to Exchange Web Services (EWS) on the Exchange server. Maximum length: 255 characters.",
			},
			"disable_proxy": schema.BoolAttribute{
				Computed:            true,
				Optional:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Bypass the web proxy (where configured for the system location) for outbound requests sent from this integration.",
			},
			"find_items_request_quota": schema.Int64Attribute{
				Computed: true,
				Optional: true,
				Default:  int64default.StaticInt64(1000000),
				Validators: []validator.Int64{
					int64validator.Between(10000, 10000000),
				},
				MarkdownDescription: "The number of Find Item requests that can be made by OTJ to your Exchange Server in a 24-hour period. Minimum: 10000. Maximum: 10000000. Default: 1000000.",
			},
			"kerberos_realm": schema.StringAttribute{
				Computed: true,
				Optional: true,
				Default:  stringdefault.StaticString(""),
				Validators: []validator.String{
					stringvalidator.LengthAtMost(250),
				},
				MarkdownDescription: "The Kerberos Realm, which is usually your domain in upper-case. Maximum length: 250 characters.",
			},
			"kerberos_kdc": schema.StringAttribute{
				Computed: true,
				Optional: true,
				Default:  stringdefault.StaticString(""),
				Validators: []validator.String{
					stringvalidator.LengthAtMost(255),
				},
				MarkdownDescription: "The address of the Kerberos key distribution center (KDC). Maximum length: 255 characters.",
			},
			"kerberos_exchange_spn": schema.StringAttribute{
				Computed: true,
				Optional: true,
				Default:  stringdefault.StaticString(""),
				Validators: []validator.String{
					stringvalidator.LengthAtMost(255),
				},
				MarkdownDescription: "The Exchange Service Principal Name (SPN). Maximum length: 255 characters.",
			},
			"kerberos_auth_every_request": schema.BoolAttribute{
				Computed:            true,
				Optional:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "When Kerberos authentication is enabled, send a Kerberos Authorization header in every request to the Exchange server.",
			},
			"kerberos_enable_tls": schema.BoolAttribute{
				Computed:            true,
				Optional:            true,
				Default:             booldefault.StaticBool(true),
				MarkdownDescription: "If enabled, all communication to the KDC will go through an HTTPS proxy and all traffic to the KDC will be encrypted using TLS.",
			},
			"kerberos_kdc_https_proxy": schema.StringAttribute{
				Computed: true,
				Optional: true,
				Default:  stringdefault.StaticString(""),
				Validators: []validator.String{
					stringvalidator.LengthAtMost(255),
				},
				MarkdownDescription: "The URL of the Kerberos key distribution center (KDC) HTTPS proxy. Maximum length: 255 characters.",
			},
			"kerberos_verify_tls_using_custom_ca": schema.BoolAttribute{
				Computed:            true,
				Optional:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "If enabled, use the configured Root Trust CA Certificates to verify the KDC HTTPS proxy SSL certificate. If disabled, the HTTPS proxy SSL certificate is verified using the system-wide default set of trusted certificates.",
			},
			"oauth_client_id": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The Application ID which was generated when creating an App Registration in Azure Active Directory.",
			},
			"oauth_auth_endpoint": schema.StringAttribute{
				Computed: true,
				Optional: true,
				Default:  stringdefault.StaticString(""),
				Validators: []validator.String{
					stringvalidator.LengthAtMost(255),
				},
				MarkdownDescription: "The URI of the OAuth authorization endpoint. This should be copied from the 'Endpoints' section in Azure Active Directory App Registrations. Maximum length: 255 characters.",
			},
			"oauth_token_endpoint": schema.StringAttribute{
				Computed: true,
				Optional: true,
				Default:  stringdefault.StaticString(""),
				Validators: []validator.String{
					stringvalidator.LengthAtMost(255),
				},
				MarkdownDescription: "The URI of the OAuth token endpoint. This should be copied from the 'Endpoints' section in Azure Active Directory App Registrations. Maximum length: 255 characters.",
			},
			"oauth_redirect_uri": schema.StringAttribute{
				Computed: true,
				Optional: true,
				Default:  stringdefault.StaticString(""),
				Validators: []validator.String{
					stringvalidator.LengthAtMost(255),
				},
				MarkdownDescription: "The redirect URI you entered when creating an App Registration in Azure Active Directory. It should be in the format 'https://[Management Node Address]/admin/platform/mjxexchangedeployment/oauth_redirect/'. Maximum length: 255 characters.",
			},
			"oauth_refresh_token": schema.StringAttribute{
				Computed:            true,
				Sensitive:           true,
				MarkdownDescription: "The OAuth refresh token which is obtained after successfully signing in via the OAuth flow.",
			},
			"oauth_state": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "A unique state which is used during the OAuth sign-in flow.",
			},
			"autodiscover_urls": schema.SetAttribute{
				Computed:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "The Autodiscover URLs associated with this One-Touch Join Exchange Integration.",
			},
			"mjx_integrations": schema.SetAttribute{
				Optional:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "The One-Touch Join Profiles associated with this OTJ Exchange Integration.",
			},
		},
		MarkdownDescription: "Manages an MJX Exchange deployment in Infinity. An MJX Exchange deployment provides integration with Microsoft Exchange, enabling OTJ (One-Touch Join) functionality for Exchange calendar-based meeting management.",
	}
}

func (r *InfinityMjxExchangeDeploymentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	plan := &InfinityMjxExchangeDeploymentResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createRequest := &config.MjxExchangeDeploymentCreateRequest{
		Name:                           plan.Name.ValueString(),
		Description:                    plan.Description.ValueString(),
		ServiceAccountUsername:         plan.ServiceAccountUsername.ValueString(),
		ServiceAccountPassword:         plan.ServiceAccountPassword.ValueString(),
		AuthenticationMethod:           plan.AuthenticationMethod.ValueString(),
		EWSURL:                         plan.EWSURL.ValueString(),
		DisableProxy:                   plan.DisableProxy.ValueBool(),
		FindItemsRequestQuota:          int(plan.FindItemsRequestQuota.ValueInt64()),
		KerberosRealm:                  plan.KerberosRealm.ValueString(),
		KerberosKDC:                    plan.KerberosKDC.ValueString(),
		KerberosExchangeSPN:            plan.KerberosExchangeSPN.ValueString(),
		KerberosAuthEveryRequest:       plan.KerberosAuthEveryRequest.ValueBool(),
		KerberosEnableTLS:              plan.KerberosEnableTLS.ValueBool(),
		KerberosKDCHTTPSProxy:          plan.KerberosKDCHTTPSProxy.ValueString(),
		KerberosVerifyTLSUsingCustomCA: plan.KerberosVerifyTLSUsingCustomCA.ValueBool(),
		OAuthAuthEndpoint:              plan.OAuthAuthEndpoint.ValueString(),
		OAuthTokenEndpoint:             plan.OAuthTokenEndpoint.ValueString(),
		OAuthRedirectURI:               plan.OAuthRedirectURI.ValueString(),
	}

	if !plan.OAuthClientID.IsNull() && !plan.OAuthClientID.IsUnknown() {
		v := plan.OAuthClientID.ValueString()
		createRequest.OAuthClientID = &v
	}

	if !plan.OAuthState.IsNull() && !plan.OAuthState.IsUnknown() {
		v := plan.OAuthState.ValueString()
		createRequest.OAuthState = &v
	}

	if !plan.MjxIntegrations.IsNull() && !plan.MjxIntegrations.IsUnknown() {
		integrations, diags := getStringList(ctx, plan.MjxIntegrations)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		createRequest.MjxIntegrations = &integrations
	}

	createResponse, err := r.InfinityClient.Config().CreateMjxExchangeDeployment(ctx, createRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Infinity MJX Exchange deployment",
			fmt.Sprintf("Could not create Infinity MJX Exchange deployment: %s", err),
		)
		return
	}

	resourceID, err := createResponse.ResourceID()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Retrieving Infinity MJX Exchange deployment ID",
			fmt.Sprintf("Could not retrieve ID for created Infinity MJX Exchange deployment: %s", err),
		)
		return
	}

	model, err := r.read(ctx, resourceID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Created Infinity MJX Exchange deployment",
			fmt.Sprintf("Could not read created Infinity MJX Exchange deployment with ID %d: %s", resourceID, err),
		)
		return
	}

	// Preserve service_account_password from plan as it is not returned by the API
	model.ServiceAccountPassword = plan.ServiceAccountPassword

	tflog.Trace(ctx, fmt.Sprintf("created Infinity MJX Exchange deployment with ID: %s, name: %s", model.ID, model.Name))

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *InfinityMjxExchangeDeploymentResource) read(ctx context.Context, resourceID int) (*InfinityMjxExchangeDeploymentResourceModel, error) {
	var data InfinityMjxExchangeDeploymentResourceModel

	srv, err := r.InfinityClient.Config().GetMjxExchangeDeployment(ctx, resourceID)
	if err != nil {
		return nil, err
	}

	if srv.ResourceURI == "" {
		return nil, fmt.Errorf("MJX Exchange deployment with ID %d not found", resourceID)
	}

	data.ID = types.StringValue(srv.ResourceURI)
	data.ResourceID = types.Int32Value(int32(resourceID)) // #nosec G115 -- API values are expected to be within int32 range
	data.Name = types.StringValue(srv.Name)
	data.Description = types.StringValue(srv.Description)
	data.ServiceAccountUsername = types.StringValue(srv.ServiceAccountUsername)
	data.AuthenticationMethod = types.StringValue(srv.AuthenticationMethod)
	data.EWSURL = types.StringValue(srv.EWSURL)
	data.DisableProxy = types.BoolValue(srv.DisableProxy)
	data.FindItemsRequestQuota = types.Int64Value(int64(srv.FindItemsRequestQuota))
	data.KerberosRealm = types.StringValue(srv.KerberosRealm)
	data.KerberosKDC = types.StringValue(srv.KerberosKDC)
	data.KerberosExchangeSPN = types.StringValue(srv.KerberosExchangeSPN)
	data.KerberosAuthEveryRequest = types.BoolValue(srv.KerberosAuthEveryRequest)
	data.KerberosEnableTLS = types.BoolValue(srv.KerberosEnableTLS)
	data.KerberosKDCHTTPSProxy = types.StringValue(srv.KerberosKDCHTTPSProxy)
	data.KerberosVerifyTLSUsingCustomCA = types.BoolValue(srv.KerberosVerifyTLSUsingCustomCA)
	data.OAuthAuthEndpoint = types.StringValue(srv.OAuthAuthEndpoint)
	data.OAuthTokenEndpoint = types.StringValue(srv.OAuthTokenEndpoint)
	data.OAuthRedirectURI = types.StringValue(srv.OAuthRedirectURI)
	data.OAuthRefreshToken = types.StringValue(srv.OAuthRefreshToken)

	// ServiceAccountPassword is not returned by the API and will be preserved from plan/state
	data.ServiceAccountPassword = types.StringNull()

	if srv.OAuthClientID != nil {
		data.OAuthClientID = types.StringValue(*srv.OAuthClientID)
	} else {
		data.OAuthClientID = types.StringNull()
	}

	if srv.OAuthState != nil {
		data.OAuthState = types.StringValue(*srv.OAuthState)
	} else {
		data.OAuthState = types.StringNull()
	}

	if srv.AutodiscoverURLs != nil && len(*srv.AutodiscoverURLs) > 0 {
		uris := make([]string, len(*srv.AutodiscoverURLs))
		for i, ref := range *srv.AutodiscoverURLs {
			uris[i] = ref.ResourceURI
		}
		urls, diags := types.SetValueFrom(ctx, types.StringType, uris)
		if diags.HasError() {
			return nil, fmt.Errorf("error converting autodiscover URLs: %v", diags)
		}
		data.AutodiscoverURLs = urls
	} else {
		data.AutodiscoverURLs = types.SetNull(types.StringType)
	}

	if srv.MjxIntegrations != nil && len(*srv.MjxIntegrations) > 0 {
		integrations, diags := types.SetValueFrom(ctx, types.StringType, *srv.MjxIntegrations)
		if diags.HasError() {
			return nil, fmt.Errorf("error converting MJX integrations: %v", diags)
		}
		data.MjxIntegrations = integrations
	} else {
		data.MjxIntegrations = types.SetNull(types.StringType)
	}

	return &data, nil
}

func (r *InfinityMjxExchangeDeploymentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	state := &InfinityMjxExchangeDeploymentResourceModel{}

	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Preserve fields not returned by the API
	serviceAccountPassword := state.ServiceAccountPassword
	oauthRefreshToken := state.OAuthRefreshToken

	resourceID := int(state.ResourceID.ValueInt32())
	state, err := r.read(ctx, resourceID)
	if err != nil {
		if isNotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading Infinity MJX Exchange deployment",
			fmt.Sprintf("Could not read Infinity MJX Exchange deployment: %s", err),
		)
		return
	}

	// Restore fields not returned consistently by the API
	state.ServiceAccountPassword = serviceAccountPassword
	state.OAuthRefreshToken = oauthRefreshToken

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *InfinityMjxExchangeDeploymentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	plan := &InfinityMjxExchangeDeploymentResourceModel{}
	state := &InfinityMjxExchangeDeploymentResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	disableProxy := plan.DisableProxy.ValueBool()
	kerberosAuthEveryRequest := plan.KerberosAuthEveryRequest.ValueBool()
	kerberosEnableTLS := plan.KerberosEnableTLS.ValueBool()
	kerberosVerifyTLS := plan.KerberosVerifyTLSUsingCustomCA.ValueBool()

	updateRequest := &config.MjxExchangeDeploymentUpdateRequest{
		Name:                           plan.Name.ValueString(),
		Description:                    plan.Description.ValueString(),
		ServiceAccountUsername:         plan.ServiceAccountUsername.ValueString(),
		ServiceAccountPassword:         plan.ServiceAccountPassword.ValueString(),
		AuthenticationMethod:           plan.AuthenticationMethod.ValueString(),
		EWSURL:                         plan.EWSURL.ValueString(),
		DisableProxy:                   &disableProxy,
		FindItemsRequestQuota:          int(plan.FindItemsRequestQuota.ValueInt64()),
		KerberosRealm:                  plan.KerberosRealm.ValueString(),
		KerberosKDC:                    plan.KerberosKDC.ValueString(),
		KerberosExchangeSPN:            plan.KerberosExchangeSPN.ValueString(),
		KerberosAuthEveryRequest:       &kerberosAuthEveryRequest,
		KerberosEnableTLS:              &kerberosEnableTLS,
		KerberosKDCHTTPSProxy:          plan.KerberosKDCHTTPSProxy.ValueString(),
		KerberosVerifyTLSUsingCustomCA: &kerberosVerifyTLS,
		OAuthAuthEndpoint:              plan.OAuthAuthEndpoint.ValueString(),
		OAuthTokenEndpoint:             plan.OAuthTokenEndpoint.ValueString(),
		OAuthRedirectURI:               plan.OAuthRedirectURI.ValueString(),
	}

	if !plan.OAuthClientID.IsNull() && !plan.OAuthClientID.IsUnknown() {
		v := plan.OAuthClientID.ValueString()
		updateRequest.OAuthClientID = &v
	}

	if !plan.OAuthState.IsNull() && !plan.OAuthState.IsUnknown() {
		v := plan.OAuthState.ValueString()
		updateRequest.OAuthState = &v
	}

	if !plan.MjxIntegrations.IsNull() && !plan.MjxIntegrations.IsUnknown() {
		integrations, diags := getStringList(ctx, plan.MjxIntegrations)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		updateRequest.MjxIntegrations = &integrations
	}

	resourceID := int(state.ResourceID.ValueInt32())
	_, err := r.InfinityClient.Config().UpdateMjxExchangeDeployment(ctx, resourceID, updateRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Infinity MJX Exchange deployment",
			fmt.Sprintf("Could not update Infinity MJX Exchange deployment: %s", err),
		)
		return
	}

	model, err := r.read(ctx, resourceID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Updated Infinity MJX Exchange deployment",
			fmt.Sprintf("Could not read updated Infinity MJX Exchange deployment with ID %d: %s", resourceID, err),
		)
		return
	}

	// Preserve service_account_password from plan as it is not returned by the API
	model.ServiceAccountPassword = plan.ServiceAccountPassword

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *InfinityMjxExchangeDeploymentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	state := &InfinityMjxExchangeDeploymentResourceModel{}

	tflog.Info(ctx, "Deleting Infinity MJX Exchange deployment")

	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.InfinityClient.Config().DeleteMjxExchangeDeployment(ctx, int(state.ResourceID.ValueInt32()))

	if err != nil && !isNotFoundError(err) && !isLookupError(err) {
		resp.Diagnostics.AddError(
			"Error Deleting Infinity MJX Exchange deployment",
			fmt.Sprintf("Could not delete Infinity MJX Exchange deployment with ID %s: %s", state.ID.ValueString(), err),
		)
		return
	}
}

func (r *InfinityMjxExchangeDeploymentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resourceID, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Import ID must be a valid integer for the resource ID. Got: %s", req.ID),
		)
		return
	}

	tflog.Trace(ctx, fmt.Sprintf("Importing Infinity MJX Exchange deployment with resource ID: %d", resourceID))

	model, err := r.read(ctx, resourceID)
	if err != nil {
		if isNotFoundError(err) {
			resp.Diagnostics.AddError(
				"Infinity MJX Exchange Deployment Not Found",
				fmt.Sprintf("Infinity MJX Exchange deployment with resource ID %d not found.", resourceID),
			)
			return
		}
		resp.Diagnostics.AddError(
			"Error Importing Infinity MJX Exchange Deployment",
			fmt.Sprintf("Could not import Infinity MJX Exchange deployment with resource ID %d: %s", resourceID, err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}
