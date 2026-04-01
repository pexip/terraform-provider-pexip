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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/pexip/go-infinity-sdk/v38/config"
)

var (
	_ resource.ResourceWithImportState = (*InfinityIdentityProviderResource)(nil)
)

type InfinityIdentityProviderResource struct {
	InfinityClient InfinityClient
}

type InfinityIdentityProviderResourceModel struct {
	ID                                  types.String `tfsdk:"id"`
	ResourceID                          types.Int32  `tfsdk:"resource_id"`
	Name                                types.String `tfsdk:"name"`
	Description                         types.String `tfsdk:"description"`
	IdpType                             types.String `tfsdk:"idp_type"`
	UUID                                types.String `tfsdk:"uuid"`
	SSOUrl                              types.String `tfsdk:"sso_url"`
	IdpEntityID                         types.String `tfsdk:"idp_entity_id"`
	IdpPublicKey                        types.String `tfsdk:"idp_public_key"`
	ServiceEntityID                     types.String `tfsdk:"service_entity_id"`
	ServicePublicKey                    types.String `tfsdk:"service_public_key"`
	ServicePrivateKey                   types.String `tfsdk:"service_private_key"`
	SignatureAlgorithm                  types.String `tfsdk:"signature_algorithm"`
	DigestAlgorithm                     types.String `tfsdk:"digest_algorithm"`
	DisplayNameAttributeName            types.String `tfsdk:"display_name_attribute_name"`
	RegistrationAliasAttributeName      types.String `tfsdk:"registration_alias_attribute_name"`
	AssertionConsumerServiceURL         types.String `tfsdk:"assertion_consumer_service_url"`
	AssertionConsumerServiceURL2        types.String `tfsdk:"assertion_consumer_service_url2"`
	AssertionConsumerServiceURL3        types.String `tfsdk:"assertion_consumer_service_url3"`
	AssertionConsumerServiceURL4        types.String `tfsdk:"assertion_consumer_service_url4"`
	AssertionConsumerServiceURL5        types.String `tfsdk:"assertion_consumer_service_url5"`
	AssertionConsumerServiceURL6        types.String `tfsdk:"assertion_consumer_service_url6"`
	AssertionConsumerServiceURL7        types.String `tfsdk:"assertion_consumer_service_url7"`
	AssertionConsumerServiceURL8        types.String `tfsdk:"assertion_consumer_service_url8"`
	AssertionConsumerServiceURL9        types.String `tfsdk:"assertion_consumer_service_url9"`
	AssertionConsumerServiceURL10       types.String `tfsdk:"assertion_consumer_service_url10"`
	WorkerFQDNACSURLs                   types.Bool   `tfsdk:"worker_fqdn_acs_urls"`
	DisablePopupFlow                    types.Bool   `tfsdk:"disable_popup_flow"`
	OidcFlow                            types.String `tfsdk:"oidc_flow"`
	OidcClientID                        types.String `tfsdk:"oidc_client_id"`
	OidcClientSecret                    types.String `tfsdk:"oidc_client_secret"`
	OidcTokenURL                        types.String `tfsdk:"oidc_token_url"`
	OidcUserInfoURL                     types.String `tfsdk:"oidc_user_info_url"`
	OidcJWKSURL                         types.String `tfsdk:"oidc_jwks_url"`
	OidcTokenEndpointAuthScheme         types.String `tfsdk:"oidc_token_endpoint_auth_scheme"`
	OidcTokenSignatureScheme            types.String `tfsdk:"oidc_token_signature_scheme"`
	OidcDisplayNameClaimName            types.String `tfsdk:"oidc_display_name_claim_name"`
	OidcRegistrationAliasClaimName      types.String `tfsdk:"oidc_registration_alias_claim_name"`
	OidcAdditionalScopes                types.String `tfsdk:"oidc_additional_scopes"`
	OidcFranceConnectRequiredEidasLevel types.String `tfsdk:"oidc_france_connect_required_eidas_level"`
	Attributes                          types.Set    `tfsdk:"attributes"`
}

func (r *InfinityIdentityProviderResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_infinity_identity_provider"
}

func (r *InfinityIdentityProviderResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *InfinityIdentityProviderResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Resource URI for the identity provider in Infinity.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"resource_id": schema.Int32Attribute{
				Computed:            true,
				MarkdownDescription: "The resource integer identifier for the identity provider in Infinity",
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(250),
				},
				MarkdownDescription: "The name used to refer to this Identity Provider. This name will be visible to end users, so you should use a name that will help users differentiate between Identity Providers without compromising security. Maximum length: 250 characters.",
			},
			"description": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(250),
				},
				MarkdownDescription: "A description of the Identity Provider. Maximum length: 250 characters.",
			},
			"idp_type": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("saml"),
				Validators: []validator.String{
					stringvalidator.OneOf("saml", "oidc"),
				},
				MarkdownDescription: "Select the protocol used by this Identity Provider. Valid choices: saml, oidc.",
			},
			"uuid": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "A unique identifier for the Identity Provider configuration. A value is automatically assigned and there is normally no need to modify it.",
			},
			"sso_url": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(255),
				},
				MarkdownDescription: "The URL to which users are sent when authenticating with this Identity Provider. Custom query string parameters may be appended, e.g. https://<url>?foo=bar. Maximum length: 255 characters. ",
			},
			"idp_entity_id": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(250),
				},
				MarkdownDescription: "The identifier for this Identity Provider integration. For SAML IdPs this is the Entity ID and for OpenID Connect IdPs this is the Issuer for returned JWTs.  Maximum length: 250 characters.",
			},
			"idp_public_key": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(4096),
				},
				MarkdownDescription: "The public key used  to verify assertions signed by this Identity Provider. Maximum length: 4096 characters.",
			},
			"service_entity_id": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(250),
				},
				MarkdownDescription: "The Entity ID for this SAML service. Maximum length: 250 characters.",
			},
			"service_public_key": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(4096),
				},
				MarkdownDescription: "Public key used by Pexip Infinity when communicating with the Identity Provider. This must be in PEM (certificate) format. Maximum length: 4096 characters.",
			},
			"service_private_key": schema.StringAttribute{
				Optional:  true,
				Computed:  true,
				Sensitive: true,
				Default:   stringdefault.StaticString(""),
				Validators: []validator.String{
					stringvalidator.LengthAtMost(12288),
				},
				MarkdownDescription: "Private key used by Pexip Infinity when communicating with the Identity Provider. Maximum length: 12288 characters.",
			},
			"signature_algorithm": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("http://www.w3.org/2001/04/xmldsig-more#rsa-sha256"),
				Validators: []validator.String{
					stringvalidator.OneOf("http://www.w3.org/2000/09/xmldsig#rsa-sha1", "http://www.w3.org/2001/04/xmldsig-more#rsa-sha256", "http://www.w3.org/2001/04/xmldsig-more#rsa-sha384", "http://www.w3.org/2001/04/xmldsig-more#rsa-sha512"),
				},
				MarkdownDescription: "Signature algorithm used to sign SAML authentication request messages and service metadata. Valid choices: http://www.w3.org/2000/09/xmldsig#rsa-sha1, http://www.w3.org/2001/04/xmldsig-more#rsa-sha256, http://www.w3.org/2001/04/xmldsig-more#rsa-sha384, http://www.w3.org/2001/04/xmldsig-more#rsa-sha512.",
			},
			"digest_algorithm": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("http://www.w3.org/2001/04/xmlenc#sha256"),
				Validators: []validator.String{
					stringvalidator.OneOf("http://www.w3.org/2000/09/xmldsig#sha1", "http://www.w3.org/2001/04/xmlenc#sha256", "http://www.w3.org/2001/04/xmldsig-more#sha384", "http://www.w3.org/2001/04/xmlenc#sha512"),
				},
				MarkdownDescription: "Digest algorithm used to sign SAML authentication request messages and service metadata. Valid choices: http://www.w3.org/2000/09/xmldsig#sha1, http://www.w3.org/2001/04/xmlenc#sha256, http://www.w3.org/2001/04/xmldsig-more#sha384, http://www.w3.org/2001/04/xmlenc#sha512.",
			},
			"display_name_attribute_name": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("NameId"),
				Validators: []validator.String{
					stringvalidator.LengthAtMost(250),
				},
				MarkdownDescription: "The SAML 2.0 attribute name from which the user's display name will be extracted. If one is not specified, participants are able to enter their own display name. Default: NameId. Maximum length: 250 characters.",
			},
			"registration_alias_attribute_name": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("NameId"),
				Validators: []validator.String{
					stringvalidator.LengthAtMost(250),
				},
				MarkdownDescription: "The SAML 2.0 attribute name from which the user's registration alias will be extracted. If one is not specified, the user will not be able to register. Default: NameId. Maximum length: 250 characters.",
			},
			"assertion_consumer_service_url": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(255),
				},
				MarkdownDescription: "A URL that can be used in the authentication process with this Identity Provider. For SAML2 this should be in the format https://<webapp_FQDN>/api/v1/samlconsumer/<uuid> and for OpenID Connect https://<webapp_FQDN>/api/v1/oidcconsumer/<uuid>. <webapp_FQDN> is the FQDN from which the web app is accessed, and <uuid> is the UUID shown in the field above. You should add one redirect URL for every web app FQDN in your deployment. Maximum length: 255 characters.",
			},
			"assertion_consumer_service_url2": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(255),
				},
				MarkdownDescription: "Enter any additional redirect URLs valid for use with this Identity Provider. Maximum length: 255 characters.",
			},
			"assertion_consumer_service_url3": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(255),
				},
				MarkdownDescription: "Enter any additional redirect URLs valid for use with this Identity Provider. Maximum length: 255 characters.",
			},
			"assertion_consumer_service_url4": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(255),
				},
				MarkdownDescription: "Enter any additional redirect URLs valid for use with this Identity Provider. Maximum length: 255 characters.",
			},
			"assertion_consumer_service_url5": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(255),
				},
				MarkdownDescription: "Enter any additional redirect URLs valid for use with this Identity Provider. Maximum length: 255 characters.",
			},
			"assertion_consumer_service_url6": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(255),
				},
				MarkdownDescription: "Enter any additional redirect URLs valid for use with this Identity Provider. Maximum length: 255 characters.",
			},
			"assertion_consumer_service_url7": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(255),
				},
				MarkdownDescription: "Enter any additional redirect URLs valid for use with this Identity Provider. Maximum length: 255 characters.",
			},
			"assertion_consumer_service_url8": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(255),
				},
				MarkdownDescription: "Enter any additional redirect URLs valid for use with this Identity Provider. Maximum length: 255 characters.",
			},
			"assertion_consumer_service_url9": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(255),
				},
				MarkdownDescription: "Enter any additional redirect URLs valid for use with this Identity Provider. Maximum length: 255 characters.",
			},
			"assertion_consumer_service_url10": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(255),
				},
				MarkdownDescription: "Enter any additional redirect URLs valid for use with this Identity Provider. Maximum length: 255 characters.",
			},
			"worker_fqdn_acs_urls": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Automatically generate allowed redirect URLs from the configured FQDNs for each Conferencing Node.",
			},
			"disable_popup_flow": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Disable pop-up windows used during Single Sign On",
			},
			"oidc_flow": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("code"),
				Validators: []validator.String{
					stringvalidator.OneOf("code", "implicit"),
				},
				MarkdownDescription: "The flow used by the OpenID Connect Identity Provider. Valid choices: implicit, code.",
			},
			"oidc_client_id": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(250),
				},
				MarkdownDescription: "The client identifier provided by the OpenID Connect Identity Provider. Maximum length: 250 characters.",
			},
			"oidc_client_secret": schema.StringAttribute{
				Optional:  true,
				Computed:  true,
				Sensitive: true,
				Default:   stringdefault.StaticString(""),
				Validators: []validator.String{
					stringvalidator.LengthAtMost(100),
				},
				MarkdownDescription: "The client secret provided by the OpenID Connect Identity Provider. Maximum length: 100 characters.",
			},
			"oidc_token_url": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(255),
				},
				MarkdownDescription: "OpenID Connect Token Endpoint URL used for exchanging codes for tokens in the Authorization Code Flow. Not required when using the Implicit Flow. Maximum length: 255 characters.",
			},
			"oidc_user_info_url": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(255),
				},
				MarkdownDescription: "You can optionally enter here the URL of an OpenID Connect UserInfo Endpoint if you wish to use this to retrieve information about the user. Maximum length: 255 characters.",
			},
			"oidc_jwks_url": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(255),
				},
				MarkdownDescription: "Download location for your Identity Provider's JSON Web Key Set (JWKS) to enable signature verification. Not required when using HS256 signatures. Maximum length: 255 characters.",
			},
			"oidc_token_endpoint_auth_scheme": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("client_secret_post"),
				Validators: []validator.String{
					stringvalidator.OneOf("client_secret_basic", "client_secret_post"),
				},
				MarkdownDescription: "The authentication method used by Infinity to authenticate when using the token endpoint. Valid choices: client_secret_basic, client_secret_post.",
			},
			"oidc_token_signature_scheme": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("rs256"),
				Validators: []validator.String{
					stringvalidator.OneOf("rs256", "hs256"),
				},
				MarkdownDescription: "The algorithm used by the Identity Provider to sign the contents of the token. Valid choices: rs256, hs256.",
			},
			"oidc_display_name_claim_name": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("name"),
				Validators: []validator.String{
					stringvalidator.LengthAtMost(250),
				},
				MarkdownDescription: "The claim name from which the user's display name will be extracted. This can come from either the JWT, or data from the UserInfo endpoint (if one is configured). Maximum length: 250 characters.",
			},
			"oidc_registration_alias_claim_name": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("sub"),
				Validators: []validator.String{
					stringvalidator.LengthAtMost(250),
				},
				MarkdownDescription: "The claim name from which the user's registration alias will be extracted. This can come from either the JWT, or data from the UserInfo endpoint (if one is configured). Maximum length: 250 characters.",
			},
			"oidc_additional_scopes": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(250),
				},
				MarkdownDescription: "Space-separated list of additional scopes to request from the OpenID Connect Identity Provider. Maximum length: 250 characters.",
			},
			"oidc_france_connect_required_eidas_level": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("disabled"),
				Validators: []validator.String{
					stringvalidator.OneOf("disabled", "eidas1", "eidas2", "eidas3"),
				},
				MarkdownDescription: "The eIDAS level to use in requests and responses. This should not be changed from the default \"Disabled\" unless advised by your Pexip support representative. Valid choices: disabled, eidas1, eidas2, eidas3.",
			},
			"attributes": schema.SetAttribute{
				Optional:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "List of identity provider attribute resource URIs.",
			},
		},
		MarkdownDescription: "Manages an identity provider configuration with the Infinity service.",
	}
}

func (r *InfinityIdentityProviderResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	plan := &InfinityIdentityProviderResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createRequest := &config.IdentityProviderCreateRequest{
		Name: plan.Name.ValueString(),
	}

	// Set optional fields
	if !plan.Description.IsNull() {
		createRequest.Description = plan.Description.ValueString()
	}
	if !plan.IdpType.IsNull() {
		createRequest.IdpType = plan.IdpType.ValueString()
	}
	if !plan.UUID.IsNull() {
		createRequest.UUID = plan.UUID.ValueString()
	}
	if !plan.SignatureAlgorithm.IsNull() {
		createRequest.SignatureAlgorithm = plan.SignatureAlgorithm.ValueString()
	}
	if !plan.DigestAlgorithm.IsNull() {
		createRequest.DigestAlgorithm = plan.DigestAlgorithm.ValueString()
	}
	if !plan.WorkerFQDNACSURLs.IsNull() {
		createRequest.WorkerFQDNACSURLs = plan.WorkerFQDNACSURLs.ValueBool()
	}
	if !plan.DisablePopupFlow.IsNull() {
		createRequest.DisablePopupFlow = plan.DisablePopupFlow.ValueBool()
	}
	if !plan.OidcTokenEndpointAuthScheme.IsNull() {
		createRequest.OidcTokenEndpointAuthScheme = plan.OidcTokenEndpointAuthScheme.ValueString()
	}
	if !plan.OidcTokenSignatureScheme.IsNull() {
		createRequest.OidcTokenSignatureScheme = plan.OidcTokenSignatureScheme.ValueString()
	}
	if !plan.OidcFranceConnectRequiredEidasLevel.IsNull() {
		createRequest.OidcFranceConnectRequiredEidasLevel = plan.OidcFranceConnectRequiredEidasLevel.ValueString()
	}
	if !plan.SSOUrl.IsNull() {
		createRequest.SSOUrl = plan.SSOUrl.ValueString()
	}
	if !plan.IdpEntityID.IsNull() {
		createRequest.IdpEntityID = plan.IdpEntityID.ValueString()
	}
	if !plan.IdpPublicKey.IsNull() {
		createRequest.IdpPublicKey = plan.IdpPublicKey.ValueString()
	}
	if !plan.ServiceEntityID.IsNull() {
		createRequest.ServiceEntityID = plan.ServiceEntityID.ValueString()
	}
	if !plan.ServicePublicKey.IsNull() {
		createRequest.ServicePublicKey = plan.ServicePublicKey.ValueString()
	}
	if !plan.ServicePrivateKey.IsNull() {
		createRequest.ServicePrivateKey = plan.ServicePrivateKey.ValueString()
	}
	if !plan.DisplayNameAttributeName.IsNull() {
		createRequest.DisplayNameAttributeName = plan.DisplayNameAttributeName.ValueString()
	}
	if !plan.RegistrationAliasAttributeName.IsNull() {
		createRequest.RegistrationAliasAttributeName = plan.RegistrationAliasAttributeName.ValueString()
	}
	if !plan.AssertionConsumerServiceURL.IsNull() && !plan.AssertionConsumerServiceURL.IsUnknown() {
		createRequest.AssertionConsumerServiceURL = plan.AssertionConsumerServiceURL.ValueString()
	} else {
		// Set a placeholder URL with zero UUID - API will update it with the actual UUID after creation
		createRequest.AssertionConsumerServiceURL = "https://localhost/samlconsumer/00000000-0000-0000-0000-000000000000"
	}
	if !plan.AssertionConsumerServiceURL2.IsNull() {
		createRequest.AssertionConsumerServiceURL2 = plan.AssertionConsumerServiceURL2.ValueString()
	}
	if !plan.AssertionConsumerServiceURL3.IsNull() {
		createRequest.AssertionConsumerServiceURL3 = plan.AssertionConsumerServiceURL3.ValueString()
	}
	if !plan.AssertionConsumerServiceURL4.IsNull() {
		createRequest.AssertionConsumerServiceURL4 = plan.AssertionConsumerServiceURL4.ValueString()
	}
	if !plan.AssertionConsumerServiceURL5.IsNull() {
		createRequest.AssertionConsumerServiceURL5 = plan.AssertionConsumerServiceURL5.ValueString()
	}
	if !plan.AssertionConsumerServiceURL6.IsNull() {
		createRequest.AssertionConsumerServiceURL6 = plan.AssertionConsumerServiceURL6.ValueString()
	}
	if !plan.AssertionConsumerServiceURL7.IsNull() {
		createRequest.AssertionConsumerServiceURL7 = plan.AssertionConsumerServiceURL7.ValueString()
	}
	if !plan.AssertionConsumerServiceURL8.IsNull() {
		createRequest.AssertionConsumerServiceURL8 = plan.AssertionConsumerServiceURL8.ValueString()
	}
	if !plan.AssertionConsumerServiceURL9.IsNull() {
		createRequest.AssertionConsumerServiceURL9 = plan.AssertionConsumerServiceURL9.ValueString()
	}
	if !plan.AssertionConsumerServiceURL10.IsNull() {
		createRequest.AssertionConsumerServiceURL10 = plan.AssertionConsumerServiceURL10.ValueString()
	}
	if !plan.OidcFlow.IsNull() {
		createRequest.OidcFlow = plan.OidcFlow.ValueString()
	}
	if !plan.OidcClientID.IsNull() {
		createRequest.OidcClientID = plan.OidcClientID.ValueString()
	}
	if !plan.OidcClientSecret.IsNull() {
		createRequest.OidcClientSecret = plan.OidcClientSecret.ValueString()
	}
	if !plan.OidcTokenURL.IsNull() {
		createRequest.OidcTokenURL = plan.OidcTokenURL.ValueString()
	}
	if !plan.OidcUserInfoURL.IsNull() {
		createRequest.OidcUserInfoURL = plan.OidcUserInfoURL.ValueString()
	}
	if !plan.OidcJWKSURL.IsNull() {
		createRequest.OidcJWKSURL = plan.OidcJWKSURL.ValueString()
	}
	if !plan.OidcDisplayNameClaimName.IsNull() {
		createRequest.OidcDisplayNameClaimName = plan.OidcDisplayNameClaimName.ValueString()
	}
	if !plan.OidcRegistrationAliasClaimName.IsNull() {
		createRequest.OidcRegistrationAliasClaimName = plan.OidcRegistrationAliasClaimName.ValueString()
	}
	if !plan.OidcAdditionalScopes.IsNull() {
		createRequest.OidcAdditionalScopes = plan.OidcAdditionalScopes.ValueString()
	}

	// Handle Attributes field
	if !plan.Attributes.IsNull() && len(plan.Attributes.Elements()) > 0 {
		var attributeURIs []string
		resp.Diagnostics.Append(plan.Attributes.ElementsAs(ctx, &attributeURIs, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		createRequest.Attributes = &attributeURIs
	}

	createResponse, err := r.InfinityClient.Config().CreateIdentityProvider(ctx, createRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Infinity identity provider",
			fmt.Sprintf("Could not create Infinity identity provider: %s", err),
		)
		return
	}

	resourceID, err := createResponse.ResourceID()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Retrieving Infinity identity provider ID",
			fmt.Sprintf("Could not retrieve ID for created Infinity identity provider: %s", err),
		)
		return
	}

	// Read the state from the API to get all computed values
	// Pass sensitive values from plan since API doesn't return them
	model, err := r.read(ctx, resourceID, plan.ServicePrivateKey.ValueString(), plan.OidcClientSecret.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Created Infinity identity provider",
			fmt.Sprintf("Could not read created Infinity identity provider with ID %d: %s", resourceID, err),
		)
		return
	}
	tflog.Trace(ctx, fmt.Sprintf("created Infinity identity provider with ID: %s, name: %s", model.ID, model.Name))

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *InfinityIdentityProviderResource) read(ctx context.Context, resourceID int, servicePrivateKey, oidcClientSecret string) (*InfinityIdentityProviderResourceModel, error) {
	var data InfinityIdentityProviderResourceModel

	srv, err := r.InfinityClient.Config().GetIdentityProvider(ctx, resourceID)
	if err != nil {
		return nil, err
	}

	if srv.ResourceURI == "" {
		return nil, fmt.Errorf("identity provider with ID %d not found", resourceID)
	}

	data.ID = types.StringValue(srv.ResourceURI)
	data.ResourceID = types.Int32Value(int32(resourceID)) // #nosec G115 -- API values are expected to be within int32 range
	data.Name = types.StringValue(srv.Name)
	data.Description = types.StringValue(srv.Description)
	data.IdpType = types.StringValue(srv.IdpType)
	data.UUID = types.StringValue(srv.UUID)
	data.SSOUrl = types.StringValue(srv.SSOUrl)
	data.IdpEntityID = types.StringValue(srv.IdpEntityID)
	data.IdpPublicKey = types.StringValue(srv.IdpPublicKey)
	data.ServiceEntityID = types.StringValue(srv.ServiceEntityID)
	data.ServicePublicKey = types.StringValue(srv.ServicePublicKey)
	data.ServicePrivateKey = types.StringValue(servicePrivateKey) // Use passed value instead of API response
	data.SignatureAlgorithm = types.StringValue(srv.SignatureAlgorithm)
	data.DigestAlgorithm = types.StringValue(srv.DigestAlgorithm)
	data.DisplayNameAttributeName = types.StringValue(srv.DisplayNameAttributeName)
	data.RegistrationAliasAttributeName = types.StringValue(srv.RegistrationAliasAttributeName)
	data.AssertionConsumerServiceURL = types.StringValue(srv.AssertionConsumerServiceURL)
	data.AssertionConsumerServiceURL2 = types.StringValue(srv.AssertionConsumerServiceURL2)
	data.AssertionConsumerServiceURL3 = types.StringValue(srv.AssertionConsumerServiceURL3)
	data.AssertionConsumerServiceURL4 = types.StringValue(srv.AssertionConsumerServiceURL4)
	data.AssertionConsumerServiceURL5 = types.StringValue(srv.AssertionConsumerServiceURL5)
	data.AssertionConsumerServiceURL6 = types.StringValue(srv.AssertionConsumerServiceURL6)
	data.AssertionConsumerServiceURL7 = types.StringValue(srv.AssertionConsumerServiceURL7)
	data.AssertionConsumerServiceURL8 = types.StringValue(srv.AssertionConsumerServiceURL8)
	data.AssertionConsumerServiceURL9 = types.StringValue(srv.AssertionConsumerServiceURL9)
	data.AssertionConsumerServiceURL10 = types.StringValue(srv.AssertionConsumerServiceURL10)
	data.WorkerFQDNACSURLs = types.BoolValue(srv.WorkerFQDNACSURLs)
	data.DisablePopupFlow = types.BoolValue(srv.DisablePopupFlow)
	data.OidcFlow = types.StringValue(srv.OidcFlow)
	data.OidcClientID = types.StringValue(srv.OidcClientID)
	data.OidcClientSecret = types.StringValue(oidcClientSecret) // Use passed value instead of API response
	data.OidcTokenURL = types.StringValue(srv.OidcTokenURL)
	data.OidcUserInfoURL = types.StringValue(srv.OidcUserInfoURL)
	data.OidcJWKSURL = types.StringValue(srv.OidcJWKSURL)
	data.OidcTokenEndpointAuthScheme = types.StringValue(srv.OidcTokenEndpointAuthScheme)
	data.OidcTokenSignatureScheme = types.StringValue(srv.OidcTokenSignatureScheme)
	data.OidcDisplayNameClaimName = types.StringValue(srv.OidcDisplayNameClaimName)
	data.OidcRegistrationAliasClaimName = types.StringValue(srv.OidcRegistrationAliasClaimName)
	data.OidcAdditionalScopes = types.StringValue(srv.OidcAdditionalScopes)
	data.OidcFranceConnectRequiredEidasLevel = types.StringValue(srv.OidcFranceConnectRequiredEidasLevel)

	// Convert attributes from SDK to Terraform format
	var attributes []string
	if srv.Attributes != nil {
		for _, attr := range *srv.Attributes {
			attributes = append(attributes, fmt.Sprintf("/api/admin/configuration/v1/identity_provider_attribute/%d/", attr.ID))
		}
	}
	attributesSetValue, diags := types.SetValueFrom(ctx, types.StringType, attributes)
	if diags.HasError() {
		return nil, fmt.Errorf("error converting attributes: %v", diags)
	}
	data.Attributes = attributesSetValue

	return &data, nil
}

func (r *InfinityIdentityProviderResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	state := &InfinityIdentityProviderResourceModel{}

	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resourceID := int(state.ResourceID.ValueInt32())
	// Pass sensitive values from state since API doesn't return them
	updatedState, err := r.read(ctx, resourceID, state.ServicePrivateKey.ValueString(), state.OidcClientSecret.ValueString())
	if err != nil {
		// Check if the error is a 404 (not found)
		if isNotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading Infinity identity provider",
			fmt.Sprintf("Could not read Infinity identity provider: %s", err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, updatedState)...)
}

func (r *InfinityIdentityProviderResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	plan := &InfinityIdentityProviderResourceModel{}
	state := &InfinityIdentityProviderResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resourceID := int(state.ResourceID.ValueInt32())

	updateRequest := &config.IdentityProviderUpdateRequest{
		Name:                        plan.Name.ValueString(),
		AssertionConsumerServiceURL: plan.AssertionConsumerServiceURL.ValueString(),
	}

	// Set boolean pointer fields for update
	if !plan.WorkerFQDNACSURLs.IsNull() {
		workerFQDN := plan.WorkerFQDNACSURLs.ValueBool()
		updateRequest.WorkerFQDNACSURLs = &workerFQDN
	}
	if !plan.DisablePopupFlow.IsNull() {
		disablePopup := plan.DisablePopupFlow.ValueBool()
		updateRequest.DisablePopupFlow = &disablePopup
	}

	// Set optional string fields
	if !plan.Description.IsNull() {
		updateRequest.Description = plan.Description.ValueString()
	}
	if !plan.IdpType.IsNull() {
		updateRequest.IdpType = plan.IdpType.ValueString()
	}
	if !plan.SignatureAlgorithm.IsNull() {
		updateRequest.SignatureAlgorithm = plan.SignatureAlgorithm.ValueString()
	}
	if !plan.DigestAlgorithm.IsNull() {
		updateRequest.DigestAlgorithm = plan.DigestAlgorithm.ValueString()
	}
	if !plan.OidcTokenEndpointAuthScheme.IsNull() {
		updateRequest.OidcTokenEndpointAuthScheme = plan.OidcTokenEndpointAuthScheme.ValueString()
	}
	if !plan.OidcTokenSignatureScheme.IsNull() {
		updateRequest.OidcTokenSignatureScheme = plan.OidcTokenSignatureScheme.ValueString()
	}
	if !plan.OidcFranceConnectRequiredEidasLevel.IsNull() {
		updateRequest.OidcFranceConnectRequiredEidasLevel = plan.OidcFranceConnectRequiredEidasLevel.ValueString()
	}
	if !plan.SSOUrl.IsNull() {
		updateRequest.SSOUrl = plan.SSOUrl.ValueString()
	}
	if !plan.IdpEntityID.IsNull() {
		updateRequest.IdpEntityID = plan.IdpEntityID.ValueString()
	}
	if !plan.IdpPublicKey.IsNull() {
		updateRequest.IdpPublicKey = plan.IdpPublicKey.ValueString()
	}
	if !plan.ServiceEntityID.IsNull() {
		updateRequest.ServiceEntityID = plan.ServiceEntityID.ValueString()
	}
	if !plan.ServicePublicKey.IsNull() {
		updateRequest.ServicePublicKey = plan.ServicePublicKey.ValueString()
	}
	if !plan.ServicePrivateKey.IsNull() {
		updateRequest.ServicePrivateKey = plan.ServicePrivateKey.ValueString()
	}
	if !plan.DisplayNameAttributeName.IsNull() {
		updateRequest.DisplayNameAttributeName = plan.DisplayNameAttributeName.ValueString()
	}
	if !plan.RegistrationAliasAttributeName.IsNull() {
		updateRequest.RegistrationAliasAttributeName = plan.RegistrationAliasAttributeName.ValueString()
	}
	// Always set additional assertion consumer URLs (no omitempty in UpdateRequest)
	updateRequest.AssertionConsumerServiceURL2 = plan.AssertionConsumerServiceURL2.ValueString()
	updateRequest.AssertionConsumerServiceURL3 = plan.AssertionConsumerServiceURL3.ValueString()
	updateRequest.AssertionConsumerServiceURL4 = plan.AssertionConsumerServiceURL4.ValueString()
	updateRequest.AssertionConsumerServiceURL5 = plan.AssertionConsumerServiceURL5.ValueString()
	updateRequest.AssertionConsumerServiceURL6 = plan.AssertionConsumerServiceURL6.ValueString()
	updateRequest.AssertionConsumerServiceURL7 = plan.AssertionConsumerServiceURL7.ValueString()
	updateRequest.AssertionConsumerServiceURL8 = plan.AssertionConsumerServiceURL8.ValueString()
	updateRequest.AssertionConsumerServiceURL9 = plan.AssertionConsumerServiceURL9.ValueString()
	updateRequest.AssertionConsumerServiceURL10 = plan.AssertionConsumerServiceURL10.ValueString()
	if !plan.OidcFlow.IsNull() {
		updateRequest.OidcFlow = plan.OidcFlow.ValueString()
	}
	if !plan.OidcClientID.IsNull() {
		updateRequest.OidcClientID = plan.OidcClientID.ValueString()
	}
	if !plan.OidcClientSecret.IsNull() {
		updateRequest.OidcClientSecret = plan.OidcClientSecret.ValueString()
	}
	if !plan.OidcTokenURL.IsNull() {
		updateRequest.OidcTokenURL = plan.OidcTokenURL.ValueString()
	}
	if !plan.OidcUserInfoURL.IsNull() {
		updateRequest.OidcUserInfoURL = plan.OidcUserInfoURL.ValueString()
	}
	if !plan.OidcJWKSURL.IsNull() {
		updateRequest.OidcJWKSURL = plan.OidcJWKSURL.ValueString()
	}
	if !plan.OidcDisplayNameClaimName.IsNull() {
		updateRequest.OidcDisplayNameClaimName = plan.OidcDisplayNameClaimName.ValueString()
	}
	if !plan.OidcRegistrationAliasClaimName.IsNull() {
		updateRequest.OidcRegistrationAliasClaimName = plan.OidcRegistrationAliasClaimName.ValueString()
	}
	if !plan.OidcAdditionalScopes.IsNull() {
		updateRequest.OidcAdditionalScopes = plan.OidcAdditionalScopes.ValueString()
	}

	// Handle Attributes field - always set it since it doesn't have omitempty in UpdateRequest
	if !plan.Attributes.IsNull() && len(plan.Attributes.Elements()) > 0 {
		var attributeURIs []string
		resp.Diagnostics.Append(plan.Attributes.ElementsAs(ctx, &attributeURIs, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		updateRequest.Attributes = &attributeURIs
	} else {
		// Set to nil to clear the attributes (sends null in JSON)
		updateRequest.Attributes = nil
	}

	_, err := r.InfinityClient.Config().UpdateIdentityProvider(ctx, resourceID, updateRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Infinity identity provider",
			fmt.Sprintf("Could not update Infinity identity provider with ID %d: %s", resourceID, err),
		)
		return
	}

	// Re-read the resource to get the latest state
	// Pass sensitive values from plan since API doesn't return them
	updatedModel, err := r.read(ctx, resourceID, plan.ServicePrivateKey.ValueString(), plan.OidcClientSecret.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Updated Infinity identity provider",
			fmt.Sprintf("Could not read updated Infinity identity provider with ID %d: %s", resourceID, err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, updatedModel)...)
}

func (r *InfinityIdentityProviderResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	state := &InfinityIdentityProviderResourceModel{}

	tflog.Info(ctx, "Deleting Infinity identity provider")

	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.InfinityClient.Config().DeleteIdentityProvider(ctx, int(state.ResourceID.ValueInt32()))

	// Ignore 404 Not Found and Lookup errors on delete
	if err != nil && !isNotFoundError(err) && !isLookupError(err) {
		resp.Diagnostics.AddError(
			"Error Deleting Infinity identity provider",
			fmt.Sprintf("Could not delete Infinity identity provider with ID %s: %s", state.ID.ValueString(), err),
		)
		return
	}
}

func (r *InfinityIdentityProviderResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resourceID, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Import ID must be a valid integer for the resource ID. Got: %s", req.ID),
		)
		return
	}

	tflog.Trace(ctx, fmt.Sprintf("Importing Infinity identity provider with resource ID: %d", resourceID))

	// Read the resource from the API
	// Pass empty strings for sensitive values since we don't have prior state on import
	// User will need to set these values in their config
	model, err := r.read(ctx, resourceID, "", "")
	if err != nil {
		// Check if the error is a 404 (not found)
		if isNotFoundError(err) {
			resp.Diagnostics.AddError(
				"Infinity Identity Provider Not Found",
				fmt.Sprintf("Infinity identity provider with resource ID %d not found.", resourceID),
			)
			return
		}
		resp.Diagnostics.AddError(
			"Error Importing Infinity Identity Provider",
			fmt.Sprintf("Could not import Infinity identity provider with resource ID %d: %s", resourceID, err),
		)
		return
	}

	// Set the state from the imported resource
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}
