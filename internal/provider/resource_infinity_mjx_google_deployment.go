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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/pexip/go-infinity-sdk/v38/config"
)

var (
	_ resource.ResourceWithImportState = (*InfinityMjxGoogleDeploymentResource)(nil)
)

type InfinityMjxGoogleDeploymentResource struct {
	InfinityClient InfinityClient
}

type InfinityMjxGoogleDeploymentResourceModel struct {
	ID                         types.String `tfsdk:"id"`
	ResourceID                 types.Int32  `tfsdk:"resource_id"`
	Name                       types.String `tfsdk:"name"`
	Description                types.String `tfsdk:"description"`
	ClientEmail                types.String `tfsdk:"client_email"`
	ClientID                   types.String `tfsdk:"client_id"`
	ClientSecret               types.String `tfsdk:"client_secret"`
	PrivateKey                 types.String `tfsdk:"private_key"`
	UseUserConsent             types.Bool   `tfsdk:"use_user_consent"`
	AuthEndpoint               types.String `tfsdk:"auth_endpoint"`
	TokenEndpoint              types.String `tfsdk:"token_endpoint"`
	RedirectURI                types.String `tfsdk:"redirect_uri"`
	RefreshToken               types.String `tfsdk:"refresh_token"`
	OAuthState                 types.String `tfsdk:"oauth_state"`
	MaximumNumberOfAPIRequests types.Int64  `tfsdk:"maximum_number_of_api_requests"`
	MjxIntegrations            types.Set    `tfsdk:"mjx_integrations"`
}

func (r *InfinityMjxGoogleDeploymentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_infinity_mjx_google_deployment"
}

func (r *InfinityMjxGoogleDeploymentResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *InfinityMjxGoogleDeploymentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Resource URI for the MJX Google deployment in Infinity.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"resource_id": schema.Int32Attribute{
				Computed:            true,
				MarkdownDescription: "The resource integer identifier for the MJX Google deployment in Infinity.",
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					stringvalidator.LengthAtMost(250),
				},
				MarkdownDescription: "The name of the OTJ Google Workspace Integration. Maximum length: 250 characters.",
			},
			"description": schema.StringAttribute{
				Computed: true,
				Optional: true,
				Default:  stringdefault.StaticString(""),
				Validators: []validator.String{
					stringvalidator.LengthAtMost(250),
				},
				MarkdownDescription: "A description of this OTJ Google Workspace Integration. Maximum length: 250 characters.",
			},
			"client_email": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(256),
				},
				MarkdownDescription: "The email address of the service account (or authorization user account) used by this OTJ Google Workspace Integration when logging in to Google Workspace to read room calendars. Maximum length: 256 characters.",
			},
			"client_id": schema.StringAttribute{
				Computed: true,
				Optional: true,
				Default:  stringdefault.StaticString(""),
				Validators: []validator.String{
					stringvalidator.LengthAtMost(250),
				},
				MarkdownDescription: "The client ID of the application you created in the Google API Console, for use by OTJ. Maximum length: 250 characters.",
			},
			"client_secret": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				MarkdownDescription: "The client secret for the application you created in the Google API Console, for use by OTJ.",
			},
			"private_key": schema.StringAttribute{
				Required:  true,
				Sensitive: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(12288),
				},
				MarkdownDescription: "The private key used by OTJ to authenticate the service account when logging in to Google Workspace to read the room calendars. Maximum length: 12288 characters.",
			},
			"use_user_consent": schema.BoolAttribute{
				Computed:            true,
				Optional:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Leave this option disabled to use the recommended method of a service account to access room calendars. Enable this option to use an authorization user, authenticated via OAuth, to access room calendars.",
			},
			"auth_endpoint": schema.StringAttribute{
				Computed: true,
				Optional: true,
				Default:  stringdefault.StaticString("https://accounts.google.com/o/oauth2/v2/auth"),
				Validators: []validator.String{
					stringvalidator.LengthAtMost(255),
				},
				MarkdownDescription: "The URI of the Google OAuth 2.0 endpoint. Maximum length: 255 characters.",
			},
			"token_endpoint": schema.StringAttribute{
				Computed: true,
				Optional: true,
				Default:  stringdefault.StaticString("https://oauth2.googleapis.com/token"),
				Validators: []validator.String{
					stringvalidator.LengthAtMost(255),
				},
				MarkdownDescription: "The URI of the Google authorization server. Maximum length: 255 characters.",
			},
			"redirect_uri": schema.StringAttribute{
				Computed: true,
				Optional: true,
				Default:  stringdefault.StaticString(""),
				Validators: []validator.String{
					stringvalidator.LengthAtMost(255),
				},
				MarkdownDescription: "The redirect URI you configured in the Google API Console Credentials. It must be in the format 'https://[Management Node Address]/admin/platform/mjxgoogledeployment/oauth_redirect/'. Maximum length: 255 characters.",
			},
			"refresh_token": schema.StringAttribute{
				Computed:            true,
				Sensitive:           true,
				MarkdownDescription: "The OAuth refresh token which is obtained after successfully finishing the authorization for accessing Google API flow.",
			},
			"oauth_state": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "A unique state which is used during the OAuth sign-in flow.",
			},
			"maximum_number_of_api_requests": schema.Int64Attribute{
				Computed: true,
				Optional: true,
				Default:  int64default.StaticInt64(900000),
				Validators: []validator.Int64{
					int64validator.AtLeast(10000),
				},
				MarkdownDescription: "The maximum number of API requests that can be made by OTJ to your Google Workspace Domain in a 24-hour period. Minimum: 10000. Default: 900000.",
			},
			"mjx_integrations": schema.SetAttribute{
				Computed:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "The OTJ Google Workspace Integration associated with this One-Touch Join Profile.",
			},
		},
		MarkdownDescription: "Manages an MJX Google deployment in Infinity. An MJX Google deployment provides integration with Google Workspace, enabling OTJ (One-Touch Join) functionality for Google Calendar-based meeting management.",
	}
}

func (r *InfinityMjxGoogleDeploymentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	plan := &InfinityMjxGoogleDeploymentResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createRequest := &config.MjxGoogleDeploymentCreateRequest{
		Name:                       plan.Name.ValueString(),
		Description:                plan.Description.ValueString(),
		ClientEmail:                plan.ClientEmail.ValueString(),
		ClientID:                   plan.ClientID.ValueString(),
		PrivateKey:                 plan.PrivateKey.ValueString(),
		UseUserConsent:             plan.UseUserConsent.ValueBool(),
		AuthEndpoint:               plan.AuthEndpoint.ValueString(),
		TokenEndpoint:              plan.TokenEndpoint.ValueString(),
		RedirectURI:                plan.RedirectURI.ValueString(),
		MaximumNumberOfAPIRequests: int(plan.MaximumNumberOfAPIRequests.ValueInt64()),
	}

	if !plan.ClientSecret.IsNull() && !plan.ClientSecret.IsUnknown() {
		createRequest.ClientSecret = plan.ClientSecret.ValueString()
	}

	if !plan.OAuthState.IsNull() && !plan.OAuthState.IsUnknown() {
		oauthState := plan.OAuthState.ValueString()
		createRequest.OAuthState = &oauthState
	}

	createResponse, err := r.InfinityClient.Config().CreateMjxGoogleDeployment(ctx, createRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Infinity MJX Google deployment",
			fmt.Sprintf("Could not create Infinity MJX Google deployment: %s", err),
		)
		return
	}

	resourceID, err := createResponse.ResourceID()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Retrieving Infinity MJX Google deployment ID",
			fmt.Sprintf("Could not retrieve ID for created Infinity MJX Google deployment: %s", err),
		)
		return
	}

	model, err := r.read(ctx, resourceID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Created Infinity MJX Google deployment",
			fmt.Sprintf("Could not read created Infinity MJX Google deployment with ID %d: %s", resourceID, err),
		)
		return
	}

	// Preserve sensitive fields from plan as they are not returned by the API
	model.PrivateKey = plan.PrivateKey
	model.ClientSecret = plan.ClientSecret

	tflog.Trace(ctx, fmt.Sprintf("created Infinity MJX Google deployment with ID: %s, name: %s", model.ID, model.Name))

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *InfinityMjxGoogleDeploymentResource) read(ctx context.Context, resourceID int) (*InfinityMjxGoogleDeploymentResourceModel, error) {
	var data InfinityMjxGoogleDeploymentResourceModel

	srv, err := r.InfinityClient.Config().GetMjxGoogleDeployment(ctx, resourceID)
	if err != nil {
		return nil, err
	}

	if srv.ResourceURI == "" {
		return nil, fmt.Errorf("MJX Google deployment with ID %d not found", resourceID)
	}

	data.ID = types.StringValue(srv.ResourceURI)
	data.ResourceID = types.Int32Value(int32(resourceID)) // #nosec G115 -- API values are expected to be within int32 range
	data.Name = types.StringValue(srv.Name)
	data.Description = types.StringValue(srv.Description)
	data.ClientEmail = types.StringValue(srv.ClientEmail)
	data.ClientID = types.StringValue(srv.ClientID)
	data.UseUserConsent = types.BoolValue(srv.UseUserConsent)
	data.AuthEndpoint = types.StringValue(srv.AuthEndpoint)
	data.TokenEndpoint = types.StringValue(srv.TokenEndpoint)
	data.RedirectURI = types.StringValue(srv.RedirectURI)
	data.RefreshToken = types.StringValue(srv.RefreshToken)
	data.MaximumNumberOfAPIRequests = types.Int64Value(int64(srv.MaximumNumberOfAPIRequests))

	// PrivateKey and ClientSecret are not returned by the API and will be preserved from plan/state
	data.PrivateKey = types.StringNull()
	data.ClientSecret = types.StringNull()

	if srv.OAuthState != nil {
		data.OAuthState = types.StringValue(*srv.OAuthState)
	} else {
		data.OAuthState = types.StringNull()
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

func (r *InfinityMjxGoogleDeploymentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	state := &InfinityMjxGoogleDeploymentResourceModel{}

	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Preserve sensitive fields from existing state
	privateKey := state.PrivateKey
	clientSecret := state.ClientSecret

	resourceID := int(state.ResourceID.ValueInt32())
	state, err := r.read(ctx, resourceID)
	if err != nil {
		if isNotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading Infinity MJX Google deployment",
			fmt.Sprintf("Could not read Infinity MJX Google deployment: %s", err),
		)
		return
	}

	// Restore sensitive fields as they are not returned by the API
	state.PrivateKey = privateKey
	state.ClientSecret = clientSecret

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *InfinityMjxGoogleDeploymentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	plan := &InfinityMjxGoogleDeploymentResourceModel{}
	state := &InfinityMjxGoogleDeploymentResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	useUserConsent := plan.UseUserConsent.ValueBool()

	updateRequest := &config.MjxGoogleDeploymentUpdateRequest{
		Name:                       plan.Name.ValueString(),
		Description:                plan.Description.ValueString(),
		ClientEmail:                plan.ClientEmail.ValueString(),
		ClientID:                   plan.ClientID.ValueString(),
		PrivateKey:                 plan.PrivateKey.ValueString(),
		UseUserConsent:             &useUserConsent,
		AuthEndpoint:               plan.AuthEndpoint.ValueString(),
		TokenEndpoint:              plan.TokenEndpoint.ValueString(),
		RedirectURI:                plan.RedirectURI.ValueString(),
		MaximumNumberOfAPIRequests: int(plan.MaximumNumberOfAPIRequests.ValueInt64()),
	}

	if !plan.ClientSecret.IsNull() && !plan.ClientSecret.IsUnknown() {
		updateRequest.ClientSecret = plan.ClientSecret.ValueString()
	}

	if !plan.OAuthState.IsNull() && !plan.OAuthState.IsUnknown() {
		oauthState := plan.OAuthState.ValueString()
		updateRequest.OAuthState = &oauthState
	}

	resourceID := int(state.ResourceID.ValueInt32())
	_, err := r.InfinityClient.Config().UpdateMjxGoogleDeployment(ctx, resourceID, updateRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Infinity MJX Google deployment",
			fmt.Sprintf("Could not update Infinity MJX Google deployment: %s", err),
		)
		return
	}

	model, err := r.read(ctx, resourceID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Updated Infinity MJX Google deployment",
			fmt.Sprintf("Could not read updated Infinity MJX Google deployment with ID %d: %s", resourceID, err),
		)
		return
	}

	// Preserve sensitive fields from plan as they are not returned by the API
	model.PrivateKey = plan.PrivateKey
	model.ClientSecret = plan.ClientSecret

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *InfinityMjxGoogleDeploymentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	state := &InfinityMjxGoogleDeploymentResourceModel{}

	tflog.Info(ctx, "Deleting Infinity MJX Google deployment")

	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.InfinityClient.Config().DeleteMjxGoogleDeployment(ctx, int(state.ResourceID.ValueInt32()))

	if err != nil && !isNotFoundError(err) && !isLookupError(err) {
		resp.Diagnostics.AddError(
			"Error Deleting Infinity MJX Google deployment",
			fmt.Sprintf("Could not delete Infinity MJX Google deployment with ID %s: %s", state.ID.ValueString(), err),
		)
		return
	}
}

func (r *InfinityMjxGoogleDeploymentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resourceID, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Import ID must be a valid integer for the resource ID. Got: %s", req.ID),
		)
		return
	}

	tflog.Trace(ctx, fmt.Sprintf("Importing Infinity MJX Google deployment with resource ID: %d", resourceID))

	model, err := r.read(ctx, resourceID)
	if err != nil {
		if isNotFoundError(err) {
			resp.Diagnostics.AddError(
				"Infinity MJX Google Deployment Not Found",
				fmt.Sprintf("Infinity MJX Google deployment with resource ID %d not found.", resourceID),
			)
			return
		}
		resp.Diagnostics.AddError(
			"Error Importing Infinity MJX Google Deployment",
			fmt.Sprintf("Could not import Infinity MJX Google deployment with resource ID %d: %s", resourceID, err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}
