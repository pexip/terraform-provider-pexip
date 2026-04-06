/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/pexip/go-infinity-sdk/v38/config"

	"github.com/pexip/terraform-provider-pexip/internal/provider/validators"
)

var (
	_ resource.ResourceWithImportState = (*InfinityAuthenticationResource)(nil)
)

type InfinityAuthenticationResource struct {
	InfinityClient InfinityClient
}

type InfinityAuthenticationResourceModel struct {
	ID                        types.String `tfsdk:"id"`
	Source                    types.String `tfsdk:"source"`
	ClientCertificate         types.String `tfsdk:"client_certificate"`
	ApiOauth2DisableBasic     types.Bool   `tfsdk:"api_oauth2_disable_basic"`
	ApiOauth2AllowAllPerms    types.Bool   `tfsdk:"api_oauth2_allow_all_perms"`
	ApiOauth2Expiration       types.Int64  `tfsdk:"api_oauth2_expiration"`
	LdapServer                types.String `tfsdk:"ldap_server"`
	LdapBaseDN                types.String `tfsdk:"ldap_base_dn"`
	LdapBindUsername          types.String `tfsdk:"ldap_bind_username"`
	LdapBindPassword          types.String `tfsdk:"ldap_bind_password"`
	LdapUserSearchDN          types.String `tfsdk:"ldap_user_search_dn"`
	LdapUserFilter            types.String `tfsdk:"ldap_user_filter"`
	LdapUserSearchFilter      types.String `tfsdk:"ldap_user_search_filter"`
	LdapUserGroupAttributes   types.String `tfsdk:"ldap_user_group_attributes"`
	LdapGroupSearchDN         types.String `tfsdk:"ldap_group_search_dn"`
	LdapGroupFilter           types.String `tfsdk:"ldap_group_filter"`
	LdapGroupMembershipFilter types.String `tfsdk:"ldap_group_membership_filter"`
	LdapUseGlobalCatalog      types.Bool   `tfsdk:"ldap_use_global_catalog"`
	LdapPermitNoTLS           types.Bool   `tfsdk:"ldap_permit_no_tls"`
	OidcMetadataURL           types.String `tfsdk:"oidc_metadata_url"`
	OidcMetadata              types.String `tfsdk:"oidc_metadata"`
	OidcClientID              types.String `tfsdk:"oidc_client_id"`
	OidcClientSecret          types.String `tfsdk:"oidc_client_secret"`
	OidcPrivateKey            types.String `tfsdk:"oidc_private_key"`
	OidcAuthMethod            types.String `tfsdk:"oidc_auth_method"`
	OidcScope                 types.String `tfsdk:"oidc_scope"`
	OidcAuthorizeURL          types.String `tfsdk:"oidc_authorize_url"`
	OidcTokenEndpointURL      types.String `tfsdk:"oidc_token_endpoint_url"`
	OidcUsernameField         types.String `tfsdk:"oidc_username_field"`
	OidcGroupsField           types.String `tfsdk:"oidc_groups_field"`
	OidcRequiredKey           types.String `tfsdk:"oidc_required_key"`
	OidcRequiredValue         types.String `tfsdk:"oidc_required_value"`
	OidcDomainHint            types.String `tfsdk:"oidc_domain_hint"`
	OidcLoginButton           types.String `tfsdk:"oidc_login_button"`
}

func (r *InfinityAuthenticationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_infinity_authentication"
}

func (r *InfinityAuthenticationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *InfinityAuthenticationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Resource URI for the authentication in Infinity.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"source": schema.StringAttribute{
				Computed: true,
				Optional: true,
				Default:  stringdefault.StaticString("LOCAL"),
				Validators: []validator.String{
					stringvalidator.OneOf("LOCAL", "LDAP", "LDAP+LOCAL", "OIDC", "OIDC+LOCAL", "LDAP+OIDC+LOCAL"),
				},
				MarkdownDescription: "The database to query for administrator authentication and authorization. Valid choices: LOCAL, LDAP, LDAP+LOCAL, OIDC, OIDC+LOCAL, LDAP+OIDC+LOCAL.",
			},
			"client_certificate": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("NO"),
				Validators: []validator.String{
					stringvalidator.OneOf("NO", "CN", "UPN"),
				},
				MarkdownDescription: "Whether to require a client TLS certificate for administrator authentication. Valid choices: NO, CN, UPN.",
			},
			"api_oauth2_disable_basic": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Disable basic authentication for management API clients. When this option is enabled, clients must use OAuth to access the management API. ",
			},
			"api_oauth2_allow_all_perms": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Allow management API clients authenticated using Oauth2 to use all permissions. ",
			},
			"api_oauth2_expiration": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				Validators: []validator.Int64{
					int64validator.AtLeast(60),
				},
				Default:             int64default.StaticInt64(3600),
				MarkdownDescription: "Specify the access token expiration time in seconds.",
			},
			"ldap_server": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(255),
				},
				Default:             stringdefault.StaticString(""),
				MarkdownDescription: "The hostname of the LDAP server. Enter a domain name for DNS SRV lookup or an FQDN for DNS A/AAAA lookup, and ensure that it is resolvable over DNS. Maximum length: 255 characters.",
			},
			"ldap_base_dn": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(""),
				Validators: []validator.String{
					stringvalidator.LengthAtMost(255),
				},
				MarkdownDescription: "The base DN of the LDAP forest to query (e.g. dc=example,dc=com). Maximum length: 255 characters.",
			},
			"ldap_bind_username": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(""),
				Validators: []validator.String{
					stringvalidator.LengthAtMost(255),
				},
				MarkdownDescription: "The username used to bind to the LDAP server. This should be a domain user service account. Maximum length: 255 characters.",
			},
			"ldap_bind_password": schema.StringAttribute{
				Optional:  true,
				Computed:  true,
				Sensitive: true,
				Default:   stringdefault.StaticString(""),
				Validators: []validator.String{
					stringvalidator.LengthAtMost(100),
				},
				MarkdownDescription: "The password used to bind to the LDAP server. Maximum length: 100 characters.",
			},
			"ldap_user_search_dn": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(""),
				Validators: []validator.String{
					stringvalidator.LengthAtMost(255),
				},
				MarkdownDescription: "The DN relative to the base DN to query for user records (e.g. ou=people). If omitted, the base DN is used. Maximum length: 255 characters.",
			},
			"ldap_user_filter": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("(&(objectclass=person)(!(objectclass=computer)))"),
				Validators: []validator.String{
					stringvalidator.LengthAtMost(1024),
				},
				MarkdownDescription: "The LDAP filter used to match user records in the directory. Default: (&(objectclass=person)(!(objectclass=computer))). Maximum length: 1024 characters.",
			},
			"ldap_user_search_filter": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("(|(uid={username})(sAMAccountName={username}))"),
				Validators: []validator.String{
					stringvalidator.LengthAtMost(1024),
				},
				MarkdownDescription: "The LDAP filter used to find user records when given the user name. The filter may contain the {username} token, which is replaced with the username, for example: (uid={username}). Maximum length: 1024 characters.",
			},
			"ldap_user_group_attributes": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("memberOf"),
				Validators: []validator.String{
					stringvalidator.LengthAtMost(100),
				},
				MarkdownDescription: "A comma-separated list of attributes in the LDAP user record to examine for group DNs when searching for the user's groups. Maximum length: 100 characters.",
			},
			"ldap_group_search_dn": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(""),
				Validators: []validator.String{
					stringvalidator.LengthAtMost(255),
				},
				MarkdownDescription: "The DN relative to the base DN to query for group records (e.g. ou=groups). If omitted, the base DN is used. Maximum length: 255 characters.",
			},
			"ldap_group_filter": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("(|(objectclass=group)(objectclass=groupOfNames)(objectclass=groupOfUniqueNames)(objectclass=posixGroup))"),
				Validators: []validator.String{
					stringvalidator.LengthAtMost(1024),
				},
				MarkdownDescription: "The LDAP filter used to match group records in the directory. Default: (|(objectclass=group)(objectclass=groupOfNames)(objectclass=groupOfUniqueNames)(objectclass=posixGroup)). Maximum length: 1024 characters.",
			},
			"ldap_group_membership_filter": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("(|(member={userdn})(uniquemember={userdn})(memberuid={useruid}))"),
				Validators: []validator.String{
					stringvalidator.LengthAtMost(1024),
				},
				MarkdownDescription: "The LDAP filter used to search for group membership of a user. The filter may contain the {userdn} and {useruid} tokens. Default: (|(member={userdn})(uniquemember={userdn})(memberuid={useruid})). Maximum length: 1024 characters.",
			},
			"ldap_use_global_catalog": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Search the Active Directory Global Catalog instead of traditional LDAP.",
			},
			"ldap_permit_no_tls": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Permit LDAP queries to be sent over an insecure connection.",
				Default:             booldefault.StaticBool(false),
			},
			"oidc_metadata_url": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					validators.URL(true),
				},
				Default:             stringdefault.StaticString(""),
				MarkdownDescription: "The URL of the OpenID Connect metadata document, copied from your OIDC provider.",
			},
			"oidc_metadata": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The OpenID Connect configuration metadata.  This will be loaded from the Metadata URL automatically.",
			},
			"oidc_client_id": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
				MarkdownDescription: "The OpenID Connect client ID.",
			},
			"oidc_client_secret": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Sensitive:           true,
				Default:             stringdefault.StaticString(""),
				MarkdownDescription: "The OpenID Connect client secret to use when authentication method is 'client_secret'. ",
			},
			"oidc_private_key": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Sensitive:           true,
				Default:             stringdefault.StaticString(""),
				MarkdownDescription: "The OpenID Connect private key to use when authentication method is 'private_key'. ",
			},
			"oidc_auth_method": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.OneOf("client_secret", "private_key"),
				},
				Default:             stringdefault.StaticString("client_secret"),
				MarkdownDescription: "The OpenID Connect authentication method. Valid choices: client_secret, private_key.",
			},
			"oidc_scope": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("openid email profile"),
				MarkdownDescription: "The OpenID Connection OAuth2 scope to request.",
			},
			"oidc_authorize_url": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					validators.URL(true),
				},
				Default:             stringdefault.StaticString(""),
				MarkdownDescription: "The OpenID Connect authorization URL.  This will be loaded from the Metadata URL automatically.",
			},
			"oidc_token_endpoint_url": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					validators.URL(true),
				},
				Default:             stringdefault.StaticString(""),
				MarkdownDescription: "The OpenID Connect token endpoint URL.  This will be loaded from the Metadata URL automatically.",
			},
			"oidc_username_field": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("preferred_username"),
				MarkdownDescription: "The field in the authentication token response to use as the username.",
			},
			"oidc_groups_field": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("groups"),
				MarkdownDescription: "The field in the authentication token response to use as the list of groups. ",
			},
			"oidc_required_key": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
				MarkdownDescription: "If there is a field in the authentication token response which must be present in order to grant access, enter the name of that field here.",
			},
			"oidc_required_value": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
				MarkdownDescription: "If you have specified a Required key, enter the value of the required key here.",
			},
			"oidc_domain_hint": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(""),
				Validators: []validator.String{
					stringvalidator.LengthAtMost(255),
				},
				MarkdownDescription: "A domain to pass to the OpenID Connect service as a hint to the expected login domain for this user. Maximum length: 255 characters.",
			},
			"oidc_login_button": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(128),
				},
				Default:             stringdefault.StaticString(""),
				MarkdownDescription: "The text to use for the OpenID Connect button on the login page of the Pexip Infinity web app.",
			},
		},
		MarkdownDescription: "Manages the authentication configuration with the Infinity service. This is a singleton resource - only one authentication configuration exists per system.",
	}
}

func (r *InfinityAuthenticationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// For singleton resources, Create is actually Update since the resource always exists
	plan := &InfinityAuthenticationResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateRequest := r.buildUpdateRequest(plan)

	_, err := r.InfinityClient.Config().UpdateAuthentication(ctx, updateRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Infinity authentication configuration",
			fmt.Sprintf("Could not update Infinity authentication configuration: %s", err),
		)
		return
	}

	// Re-read the resource to get the latest state
	updatedModel, err := r.read(ctx, plan.LdapBindPassword.ValueString(), plan.OidcClientSecret.ValueString(), plan.OidcPrivateKey.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Updated Infinity authentication configuration",
			fmt.Sprintf("Could not read updated Infinity authentication configuration: %s", err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, updatedModel)...)
}

func (r *InfinityAuthenticationResource) buildUpdateRequest(plan *InfinityAuthenticationResourceModel) *config.AuthenticationUpdateRequest {
	updateRequest := &config.AuthenticationUpdateRequest{
		Source:                    plan.Source.ValueString(),
		ClientCertificate:         plan.ClientCertificate.ValueString(),
		LdapServer:                plan.LdapServer.ValueString(),
		LdapBaseDN:                plan.LdapBaseDN.ValueString(),
		LdapBindUsername:          plan.LdapBindUsername.ValueString(),
		LdapBindPassword:          plan.LdapBindPassword.ValueString(),
		LdapUserSearchDN:          plan.LdapUserSearchDN.ValueString(),
		LdapUserFilter:            plan.LdapUserFilter.ValueString(),
		LdapUserSearchFilter:      plan.LdapUserSearchFilter.ValueString(),
		LdapUserGroupAttributes:   plan.LdapUserGroupAttributes.ValueString(),
		LdapGroupSearchDN:         plan.LdapGroupSearchDN.ValueString(),
		LdapGroupFilter:           plan.LdapGroupFilter.ValueString(),
		LdapGroupMembershipFilter: plan.LdapGroupMembershipFilter.ValueString(),
		OidcMetadataURL:           plan.OidcMetadataURL.ValueString(),
		OidcMetadata:              plan.OidcMetadata.ValueString(),
		OidcClientID:              plan.OidcClientID.ValueString(),
		OidcClientSecret:          plan.OidcClientSecret.ValueString(),
		OidcPrivateKey:            plan.OidcPrivateKey.ValueString(),
		OidcAuthMethod:            plan.OidcAuthMethod.ValueString(),
		OidcScope:                 plan.OidcScope.ValueString(),
		OidcAuthorizeURL:          plan.OidcAuthorizeURL.ValueString(),
		OidcTokenEndpointURL:      plan.OidcTokenEndpointURL.ValueString(),
		OidcUsernameField:         plan.OidcUsernameField.ValueString(),
		OidcGroupsField:           plan.OidcGroupsField.ValueString(),
		OidcRequiredKey:           plan.OidcRequiredKey.ValueString(),
		OidcRequiredValue:         plan.OidcRequiredValue.ValueString(),
		OidcDomainHint:            plan.OidcDomainHint.ValueString(),
		OidcLoginButton:           plan.OidcLoginButton.ValueString(),
	}

	apiOauth2DisableBasic := plan.ApiOauth2DisableBasic.ValueBool()
	updateRequest.ApiOauth2DisableBasic = &apiOauth2DisableBasic

	apiOauth2AllowAllPerms := plan.ApiOauth2AllowAllPerms.ValueBool()
	updateRequest.ApiOauth2AllowAllPerms = &apiOauth2AllowAllPerms

	ldapUseGlobalCatalog := plan.LdapUseGlobalCatalog.ValueBool()
	updateRequest.LdapUseGlobalCatalog = &ldapUseGlobalCatalog

	ldapPermitNoTLS := plan.LdapPermitNoTLS.ValueBool()
	updateRequest.LdapPermitNoTLS = &ldapPermitNoTLS

	apiOauth2Expiration := int(plan.ApiOauth2Expiration.ValueInt64())
	updateRequest.ApiOauth2Expiration = &apiOauth2Expiration

	return updateRequest
}

func (r *InfinityAuthenticationResource) read(ctx context.Context, ldapPass, oidcPass, oidcKey string) (*InfinityAuthenticationResourceModel, error) {
	var data InfinityAuthenticationResourceModel

	srv, err := r.InfinityClient.Config().GetAuthentication(ctx)
	if err != nil {
		return nil, err
	}

	if srv.ResourceURI == "" {
		return nil, fmt.Errorf("authentication configuration not found")
	}

	data.ID = types.StringValue(srv.ResourceURI)
	data.Source = types.StringValue(srv.Source)
	data.ClientCertificate = types.StringValue(srv.ClientCertificate)
	data.ApiOauth2DisableBasic = types.BoolValue(srv.ApiOauth2DisableBasic)
	data.ApiOauth2AllowAllPerms = types.BoolValue(srv.ApiOauth2AllowAllPerms)
	data.ApiOauth2Expiration = types.Int64Value(int64(srv.ApiOauth2Expiration))
	data.LdapServer = types.StringValue(srv.LdapServer)
	data.LdapBaseDN = types.StringValue(srv.LdapBaseDN)
	data.LdapBindUsername = types.StringValue(srv.LdapBindUsername)
	data.LdapBindPassword = types.StringValue(ldapPass)
	data.LdapUserSearchDN = types.StringValue(srv.LdapUserSearchDN)
	data.LdapUserFilter = types.StringValue(srv.LdapUserFilter)
	data.LdapUserSearchFilter = types.StringValue(srv.LdapUserSearchFilter)
	data.LdapUserGroupAttributes = types.StringValue(srv.LdapUserGroupAttributes)
	data.LdapGroupSearchDN = types.StringValue(srv.LdapGroupSearchDN)
	data.LdapGroupFilter = types.StringValue(srv.LdapGroupFilter)
	data.LdapGroupMembershipFilter = types.StringValue(srv.LdapGroupMembershipFilter)
	data.LdapUseGlobalCatalog = types.BoolValue(srv.LdapUseGlobalCatalog)
	data.LdapPermitNoTLS = types.BoolValue(srv.LdapPermitNoTLS)
	data.OidcMetadataURL = types.StringValue(srv.OidcMetadataURL)
	data.OidcMetadata = types.StringValue(srv.OidcMetadata)
	data.OidcClientID = types.StringValue(srv.OidcClientID)
	data.OidcClientSecret = types.StringValue(oidcPass)
	data.OidcPrivateKey = types.StringValue(oidcKey)
	data.OidcAuthMethod = types.StringValue(srv.OidcAuthMethod)
	data.OidcScope = types.StringValue(srv.OidcScope)
	data.OidcAuthorizeURL = types.StringValue(srv.OidcAuthorizeURL)
	data.OidcTokenEndpointURL = types.StringValue(srv.OidcTokenEndpointURL)
	data.OidcUsernameField = types.StringValue(srv.OidcUsernameField)
	data.OidcGroupsField = types.StringValue(srv.OidcGroupsField)
	data.OidcRequiredKey = types.StringValue(srv.OidcRequiredKey)
	data.OidcRequiredValue = types.StringValue(srv.OidcRequiredValue)
	data.OidcDomainHint = types.StringValue(srv.OidcDomainHint)
	data.OidcLoginButton = types.StringValue(srv.OidcLoginButton)

	return &data, nil
}

func (r *InfinityAuthenticationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	state := &InfinityAuthenticationResourceModel{}

	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	state, err := r.read(ctx, state.LdapBindPassword.ValueString(), state.OidcClientSecret.ValueString(), state.OidcPrivateKey.ValueString())
	if err != nil {
		// Check if the error is a 404 (not found) - unlikely for singleton resources
		if isNotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading Infinity authentication configuration",
			fmt.Sprintf("Could not read Infinity authentication configuration: %s", err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *InfinityAuthenticationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	plan := &InfinityAuthenticationResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateRequest := r.buildUpdateRequest(plan)

	_, err := r.InfinityClient.Config().UpdateAuthentication(ctx, updateRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Infinity authentication configuration",
			fmt.Sprintf("Could not update Infinity authentication configuration: %s", err),
		)
		return
	}

	// Re-read the resource to get the latest state
	updatedModel, err := r.read(ctx, plan.LdapBindPassword.ValueString(), plan.OidcClientSecret.ValueString(), plan.OidcPrivateKey.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Updated Infinity authentication configuration",
			fmt.Sprintf("Could not read updated Infinity authentication configuration: %s", err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, updatedModel)...)
}

func (r *InfinityAuthenticationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// For singleton resources, delete means resetting all fields to their API defaults.
	tflog.Info(ctx, "Resetting Infinity authentication configuration to defaults")

	falseVal := false
	expirationDefault := 3600

	updateRequest := &config.AuthenticationUpdateRequest{
		Source:                    "LOCAL",
		ClientCertificate:         "NO",
		ApiOauth2DisableBasic:     &falseVal,
		ApiOauth2AllowAllPerms:    &falseVal,
		ApiOauth2Expiration:       &expirationDefault,
		LdapServer:                "",
		LdapBaseDN:                "",
		LdapBindUsername:          "",
		LdapBindPassword:          "",
		LdapUserSearchDN:          "",
		LdapUserFilter:            "(&(objectclass=person)(!(objectclass=computer)))",
		LdapUserSearchFilter:      "(|(uid={username})(sAMAccountName={username}))",
		LdapUserGroupAttributes:   "memberOf",
		LdapGroupSearchDN:         "",
		LdapGroupFilter:           "(|(objectclass=group)(objectclass=groupOfNames)(objectclass=groupOfUniqueNames)(objectclass=posixGroup))",
		LdapGroupMembershipFilter: "(|(member={userdn})(uniquemember={userdn})(memberuid={useruid}))",
		LdapUseGlobalCatalog:      &falseVal,
		LdapPermitNoTLS:           &falseVal,
		OidcMetadataURL:           "",
		OidcMetadata:              "",
		OidcClientID:              "",
		OidcClientSecret:          "",
		OidcPrivateKey:            "",
		OidcAuthMethod:            "client_secret",
		OidcScope:                 "openid email profile",
		OidcAuthorizeURL:          "",
		OidcTokenEndpointURL:      "",
		OidcUsernameField:         "preferred_username",
		OidcGroupsField:           "groups",
		OidcRequiredKey:           "",
		OidcRequiredValue:         "",
		OidcDomainHint:            "",
		OidcLoginButton:           "",
	}

	_, err := r.InfinityClient.Config().UpdateAuthentication(ctx, updateRequest)
	if err != nil && !isNotFoundError(err) && !isLookupError(err) {
		resp.Diagnostics.AddError(
			"Error Resetting Infinity authentication configuration",
			fmt.Sprintf("Could not reset Infinity authentication configuration: %s", err),
		)
		return
	}
}

func (r *InfinityAuthenticationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// For singleton resources, the import ID doesn't matter since there's only one instance
	tflog.Trace(ctx, "Importing Infinity authentication configuration")

	// Read the resource from the API
	model, err := r.read(ctx, "", "", "")
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Importing Infinity Authentication Configuration",
			fmt.Sprintf("Could not import Infinity authentication configuration: %s", err),
		)
		return
	}

	// Set the state from the imported resource
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}
