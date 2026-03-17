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

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/pexip/go-infinity-sdk/v38/config"
)

var (
	_ resource.ResourceWithImportState = (*InfinityPolicyServerResource)(nil)
)

type InfinityPolicyServerResource struct {
	InfinityClient InfinityClient
}

type InfinityPolicyServerResourceModel struct {
	ID                                  types.String `tfsdk:"id"`
	ResourceID                          types.Int32  `tfsdk:"resource_id"`
	Name                                types.String `tfsdk:"name"`
	Description                         types.String `tfsdk:"description"`
	URL                                 types.String `tfsdk:"url"`
	Username                            types.String `tfsdk:"username"`
	Password                            types.String `tfsdk:"password"`
	EnableServiceLookup                 types.Bool   `tfsdk:"enable_service_lookup"`
	EnableParticipantLookup             types.Bool   `tfsdk:"enable_participant_lookup"`
	EnableRegistrationLookup            types.Bool   `tfsdk:"enable_registration_lookup"`
	EnableDirectoryLookup               types.Bool   `tfsdk:"enable_directory_lookup"`
	EnableAvatarLookup                  types.Bool   `tfsdk:"enable_avatar_lookup"`
	EnableMediaLocationLookup           types.Bool   `tfsdk:"enable_media_location_lookup"`
	EnableInternalServicePolicy         types.Bool   `tfsdk:"enable_internal_service_policy"`
	EnableInternalParticipantPolicy     types.Bool   `tfsdk:"enable_internal_participant_policy"`
	EnableInternalMediaLocationPolicy   types.Bool   `tfsdk:"enable_internal_media_location_policy"`
	PreferLocalAvatarConfiguration      types.Bool   `tfsdk:"prefer_local_avatar_configuration"`
	InternalServicePolicyTemplate       types.String `tfsdk:"internal_service_policy_template"`
	InternalParticipantPolicyTemplate   types.String `tfsdk:"internal_participant_policy_template"`
	InternalMediaLocationPolicyTemplate types.String `tfsdk:"internal_media_location_policy_template"`
}

func (r *InfinityPolicyServerResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_infinity_policy_server"
}

func (r *InfinityPolicyServerResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *InfinityPolicyServerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Resource URI for the policy server in Infinity.",
			},
			"resource_id": schema.Int32Attribute{
				Computed:            true,
				MarkdownDescription: "The resource integer identifier for the policy server in Infinity",
			},
			"name": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(250),
				},
				MarkdownDescription: "The name used to refer to this policy profile. Maximum length: 250 characters.",
			},
			"description": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(""),
				Validators: []validator.String{
					stringvalidator.LengthAtMost(250),
				},
				MarkdownDescription: "A description of the policy profile. Maximum length: 250 characters.",
			},
			"url": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(""),
				Validators: []validator.String{
					stringvalidator.LengthAtMost(255),
				},
				MarkdownDescription: "The URL of the external policy server. Maximum length: 255 characters. ",
			},
			"username": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(""),
				Validators: []validator.String{
					stringvalidator.LengthAtMost(100),
				},
				MarkdownDescription: "The username used when accessing the policy server. Maximum length: 100 characters.",
			},
			"password": schema.StringAttribute{
				Optional:  true,
				Computed:  true,
				Sensitive: true,
				Default:   stringdefault.StaticString(""),
				Validators: []validator.String{
					stringvalidator.LengthAtMost(100),
				},
				MarkdownDescription: "The password used when accessing the policy server. Maximum length: 100 characters.",
			},
			"enable_service_lookup": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "If enabled, requests are sent to the external policy server to fetch service configuration details. ",
			},
			"enable_participant_lookup": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "If enabled, requests are sent to the external policy server to fetch participant configuration details. ",
			},
			"enable_registration_lookup": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "If enabled, requests are sent to the external policy server to determine whether the device is permitted to register to Pexip Infinity. ",
			},
			"enable_directory_lookup": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "If enabled, requests are sent to the external policy server to fetch directory listing results. ",
			},
			"enable_avatar_lookup": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "If enabled, requests are sent to the external policy server to fetch avatar images for participants. ",
			},
			"enable_media_location_lookup": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "If enabled, requests are sent to the external policy server to fetch the system location to use for media processing. ",
			},
			"enable_internal_service_policy": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "If enabled, service configuration retrieved from the local database or an external policy server can be transformed using Jinja2 templating before being returned to the endpoint. ",
			},
			"enable_internal_participant_policy": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "If enabled, the participant configuration based on information from the local database or an external policy server can be transformed using Jinja2 templating before being returned to the endpoint. ",
			},
			"enable_internal_media_location_policy": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "If enabled, the media location configuration based on information from the local database or an external policy server can be transformed using Jinja2 templating before being returned to the endpoint. ",
			},
			"prefer_local_avatar_configuration": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "If enabled, requests are sent to the Avatar URL configured on the user when logging into the Pexip Infinity web app. Note: Avatar URL for Pexip Infinity web app (conference control) is not currently supported, so avatars are not displayed in the conference roster or in self-view. ",
			},
			"internal_service_policy_template": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(""),
				Validators: []validator.String{
					stringvalidator.LengthAtMost(49152),
				},
				MarkdownDescription: "A Jinja2 script that takes the existing service configuration (if any) and optionally modifies or overrides service settings. Maximum length: 49152 characters.",
			},
			"internal_participant_policy_template": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(""),
				Validators: []validator.String{
					stringvalidator.LengthAtMost(49152),
				},
				MarkdownDescription: "A Jinja2 script that takes the existing participant configuration and optionally modifies or overrides participant settings. Maximum length: 49152 characters.",
			},
			"internal_media_location_policy_template": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(""),
				Validators: []validator.String{
					stringvalidator.LengthAtMost(49152),
				},
				MarkdownDescription: "A Jinja2 script that takes the existing media location configuration and optionally modifies or overrides location settings. Maximum length: 49152 characters.",
			},
		},
		MarkdownDescription: "Manages a policy server configuration with the Infinity service.",
	}
}

func (r *InfinityPolicyServerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	plan := &InfinityPolicyServerResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createRequest := &config.PolicyServerCreateRequest{
		Name:                              plan.Name.ValueString(),
		EnableServiceLookup:               plan.EnableServiceLookup.ValueBool(),
		EnableParticipantLookup:           plan.EnableParticipantLookup.ValueBool(),
		EnableRegistrationLookup:          plan.EnableRegistrationLookup.ValueBool(),
		EnableDirectoryLookup:             plan.EnableDirectoryLookup.ValueBool(),
		EnableAvatarLookup:                plan.EnableAvatarLookup.ValueBool(),
		EnableMediaLocationLookup:         plan.EnableMediaLocationLookup.ValueBool(),
		EnableInternalServicePolicy:       plan.EnableInternalServicePolicy.ValueBool(),
		EnableInternalParticipantPolicy:   plan.EnableInternalParticipantPolicy.ValueBool(),
		EnableInternalMediaLocationPolicy: plan.EnableInternalMediaLocationPolicy.ValueBool(),
		PreferLocalAvatarConfiguration:    plan.PreferLocalAvatarConfiguration.ValueBool(),
	}

	// Set optional string fields
	if !plan.Description.IsNull() {
		createRequest.Description = plan.Description.ValueString()
	}
	if !plan.URL.IsNull() {
		createRequest.URL = plan.URL.ValueString()
	}
	if !plan.Username.IsNull() {
		createRequest.Username = plan.Username.ValueString()
	}
	if !plan.Password.IsNull() {
		createRequest.Password = plan.Password.ValueString()
	}
	if !plan.InternalServicePolicyTemplate.IsNull() {
		createRequest.InternalServicePolicyTemplate = plan.InternalServicePolicyTemplate.ValueString()
	}
	if !plan.InternalParticipantPolicyTemplate.IsNull() {
		createRequest.InternalParticipantPolicyTemplate = plan.InternalParticipantPolicyTemplate.ValueString()
	}
	if !plan.InternalMediaLocationPolicyTemplate.IsNull() {
		createRequest.InternalMediaLocationPolicyTemplate = plan.InternalMediaLocationPolicyTemplate.ValueString()
	}

	createResponse, err := r.InfinityClient.Config().CreatePolicyServer(ctx, createRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Infinity policy server",
			fmt.Sprintf("Could not create Infinity policy server: %s", err),
		)
		return
	}

	resourceID, err := createResponse.ResourceID()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Retrieving Infinity policy server ID",
			fmt.Sprintf("Could not retrieve ID for created Infinity policy server: %s", err),
		)
		return
	}

	// Read the state from the API to get all computed values
	model, err := r.read(ctx, resourceID, plan.Password.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Created Infinity policy server",
			fmt.Sprintf("Could not read created Infinity policy server with ID %d: %s", resourceID, err),
		)
		return
	}
	tflog.Trace(ctx, fmt.Sprintf("created Infinity policy server with ID: %s, name: %s", model.ID, model.Name))

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *InfinityPolicyServerResource) read(ctx context.Context, resourceID int, password string) (*InfinityPolicyServerResourceModel, error) {
	var data InfinityPolicyServerResourceModel

	srv, err := r.InfinityClient.Config().GetPolicyServer(ctx, resourceID)
	if err != nil {
		return nil, err
	}

	if srv.ResourceURI == "" {
		return nil, fmt.Errorf("policy server with ID %d not found", resourceID)
	}

	data.ID = types.StringValue(srv.ResourceURI)
	data.ResourceID = types.Int32Value(int32(resourceID)) // #nosec G115 -- API values are expected to be within int32 range
	data.Name = types.StringValue(srv.Name)
	data.Description = types.StringValue(srv.Description)
	data.URL = types.StringValue(srv.URL)
	data.Username = types.StringValue(srv.Username)
	data.Password = types.StringValue(password) // The password property of the policy server is returned in hashed format, so we need to ignore it by setting it to the input string
	data.EnableServiceLookup = types.BoolValue(srv.EnableServiceLookup)
	data.EnableParticipantLookup = types.BoolValue(srv.EnableParticipantLookup)
	data.EnableRegistrationLookup = types.BoolValue(srv.EnableRegistrationLookup)
	data.EnableDirectoryLookup = types.BoolValue(srv.EnableDirectoryLookup)
	data.EnableAvatarLookup = types.BoolValue(srv.EnableAvatarLookup)
	data.EnableMediaLocationLookup = types.BoolValue(srv.EnableMediaLocationLookup)
	data.EnableInternalServicePolicy = types.BoolValue(srv.EnableInternalServicePolicy)
	data.EnableInternalParticipantPolicy = types.BoolValue(srv.EnableInternalParticipantPolicy)
	data.EnableInternalMediaLocationPolicy = types.BoolValue(srv.EnableInternalMediaLocationPolicy)
	data.PreferLocalAvatarConfiguration = types.BoolValue(srv.PreferLocalAvatarConfiguration)
	data.InternalServicePolicyTemplate = types.StringValue(srv.InternalServicePolicyTemplate)
	data.InternalParticipantPolicyTemplate = types.StringValue(srv.InternalParticipantPolicyTemplate)
	data.InternalMediaLocationPolicyTemplate = types.StringValue(srv.InternalMediaLocationPolicyTemplate)

	return &data, nil
}

func (r *InfinityPolicyServerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	state := &InfinityPolicyServerResourceModel{}

	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resourceID := int(state.ResourceID.ValueInt32())
	state, err := r.read(ctx, resourceID, state.Password.ValueString())
	if err != nil {
		// Check if the error is a 404 (not found)
		if isNotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading Infinity policy server",
			fmt.Sprintf("Could not read Infinity policy server: %s", err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *InfinityPolicyServerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	plan := &InfinityPolicyServerResourceModel{}
	state := &InfinityPolicyServerResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resourceID := int(state.ResourceID.ValueInt32())

	enableServiceLookup := plan.EnableServiceLookup.ValueBool()
	enableParticipantLookup := plan.EnableParticipantLookup.ValueBool()
	enableRegistrationLookup := plan.EnableRegistrationLookup.ValueBool()
	enableDirectoryLookup := plan.EnableDirectoryLookup.ValueBool()
	enableAvatarLookup := plan.EnableAvatarLookup.ValueBool()
	enableMediaLocationLookup := plan.EnableMediaLocationLookup.ValueBool()
	enableInternalServicePolicy := plan.EnableInternalServicePolicy.ValueBool()
	enableInternalParticipantPolicy := plan.EnableInternalParticipantPolicy.ValueBool()
	enableInternalMediaLocationPolicy := plan.EnableInternalMediaLocationPolicy.ValueBool()
	preferLocalAvatarConfiguration := plan.PreferLocalAvatarConfiguration.ValueBool()

	updateRequest := &config.PolicyServerUpdateRequest{
		Name:                              plan.Name.ValueString(),
		EnableServiceLookup:               &enableServiceLookup,
		EnableParticipantLookup:           &enableParticipantLookup,
		EnableRegistrationLookup:          &enableRegistrationLookup,
		EnableDirectoryLookup:             &enableDirectoryLookup,
		EnableAvatarLookup:                &enableAvatarLookup,
		EnableMediaLocationLookup:         &enableMediaLocationLookup,
		EnableInternalServicePolicy:       &enableInternalServicePolicy,
		EnableInternalParticipantPolicy:   &enableInternalParticipantPolicy,
		EnableInternalMediaLocationPolicy: &enableInternalMediaLocationPolicy,
		PreferLocalAvatarConfiguration:    &preferLocalAvatarConfiguration,
	}

	// Set optional string fields
	if !plan.Description.IsNull() {
		updateRequest.Description = plan.Description.ValueString()
	}
	if !plan.URL.IsNull() {
		updateRequest.URL = plan.URL.ValueString()
	}
	if !plan.Username.IsNull() {
		updateRequest.Username = plan.Username.ValueString()
	}
	if !plan.Password.IsNull() {
		updateRequest.Password = plan.Password.ValueString()
	}
	if !plan.InternalServicePolicyTemplate.IsNull() {
		updateRequest.InternalServicePolicyTemplate = plan.InternalServicePolicyTemplate.ValueString()
	}
	if !plan.InternalParticipantPolicyTemplate.IsNull() {
		updateRequest.InternalParticipantPolicyTemplate = plan.InternalParticipantPolicyTemplate.ValueString()
	}
	if !plan.InternalMediaLocationPolicyTemplate.IsNull() {
		updateRequest.InternalMediaLocationPolicyTemplate = plan.InternalMediaLocationPolicyTemplate.ValueString()
	}

	_, err := r.InfinityClient.Config().UpdatePolicyServer(ctx, resourceID, updateRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Infinity policy server",
			fmt.Sprintf("Could not update Infinity policy server with ID %d: %s", resourceID, err),
		)
		return
	}

	// Re-read the resource to get the latest state
	updatedModel, err := r.read(ctx, resourceID, plan.Password.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Updated Infinity policy server",
			fmt.Sprintf("Could not read updated Infinity policy server with ID %d: %s", resourceID, err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, updatedModel)...)
}

func (r *InfinityPolicyServerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	state := &InfinityPolicyServerResourceModel{}

	tflog.Info(ctx, "Deleting Infinity policy server")

	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.InfinityClient.Config().DeletePolicyServer(ctx, int(state.ResourceID.ValueInt32()))

	// Ignore 404 Not Found and Lookup errors on delete
	if err != nil && !isNotFoundError(err) && !isLookupError(err) {
		resp.Diagnostics.AddError(
			"Error Deleting Infinity policy server",
			fmt.Sprintf("Could not delete Infinity policy server with ID %s: %s", state.ID.ValueString(), err),
		)
		return
	}
}

func (r *InfinityPolicyServerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resourceID, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Import ID must be a valid integer for the resource ID. Got: %s", req.ID),
		)
		return
	}

	tflog.Trace(ctx, fmt.Sprintf("Importing Infinity policy server with resource ID: %d", resourceID))

	// Read the resource from the API
	model, err := r.read(ctx, resourceID, "")
	if err != nil {
		// Check if the error is a 404 (not found)
		if isNotFoundError(err) {
			resp.Diagnostics.AddError(
				"Infinity Policy Server Not Found",
				fmt.Sprintf("Infinity policy server with resource ID %d not found.", resourceID),
			)
			return
		}
		resp.Diagnostics.AddError(
			"Error Importing Infinity Policy Server",
			fmt.Sprintf("Could not import Infinity policy server with resource ID %d: %s", resourceID, err),
		)
		return
	}

	// Set the state from the imported resource
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}
