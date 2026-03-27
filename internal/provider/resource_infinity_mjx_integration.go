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
	_ resource.ResourceWithImportState = (*InfinityMjxIntegrationResource)(nil)
)

type InfinityMjxIntegrationResource struct {
	InfinityClient InfinityClient
}

type InfinityMjxIntegrationResourceModel struct {
	ID                          types.String `tfsdk:"id"`
	ResourceID                  types.Int32  `tfsdk:"resource_id"`
	Name                        types.String `tfsdk:"name"`
	Description                 types.String `tfsdk:"description"`
	DisplayUpcomingMeetings     types.Int64  `tfsdk:"display_upcoming_meetings"`
	EnableNonVideoMeetings      types.Bool   `tfsdk:"enable_non_video_meetings"`
	EnablePrivateMeetings       types.Bool   `tfsdk:"enable_private_meetings"`
	EndBuffer                   types.Int64  `tfsdk:"end_buffer"`
	StartBuffer                 types.Int64  `tfsdk:"start_buffer"`
	EPUsername                  types.String `tfsdk:"ep_username"`
	EPPassword                  types.String `tfsdk:"ep_password"`
	EPUseHTTPS                  types.Bool   `tfsdk:"ep_use_https"`
	EPVerifyCertificate         types.Bool   `tfsdk:"ep_verify_certificate"`
	ExchangeDeployment          types.String `tfsdk:"exchange_deployment"`
	GoogleDeployment            types.String `tfsdk:"google_deployment"`
	GraphDeployment             types.String `tfsdk:"graph_deployment"`
	ProcessAliasPrivateMeetings types.Bool   `tfsdk:"process_alias_private_meetings"`
	ReplaceEmptySubject         types.Bool   `tfsdk:"replace_empty_subject"`
	ReplaceSubjectType          types.String `tfsdk:"replace_subject_type"`
	ReplaceSubjectTemplate      types.String `tfsdk:"replace_subject_template"`
	UseWebex                    types.Bool   `tfsdk:"use_webex"`
	WebexAPIDomain              types.String `tfsdk:"webex_api_domain"`
	WebexClientID               types.String `tfsdk:"webex_client_id"`
	WebexClientSecret           types.String `tfsdk:"webex_client_secret"`
	WebexOAuthState             types.String `tfsdk:"webex_oauth_state"`
	WebexRedirectURI            types.String `tfsdk:"webex_redirect_uri"`
	WebexRefreshToken           types.String `tfsdk:"webex_refresh_token"`
	EndpointGroups              types.Set    `tfsdk:"endpoint_groups"`
}

func (r *InfinityMjxIntegrationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_infinity_mjx_integration"
}

func (r *InfinityMjxIntegrationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *InfinityMjxIntegrationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Resource URI for the MJX integration in Infinity.",
			},
			"resource_id": schema.Int32Attribute{
				Computed:            true,
				MarkdownDescription: "The resource integer identifier for the MJX integration in Infinity.",
			},
			"name": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					stringvalidator.LengthAtMost(250),
				},
				MarkdownDescription: "The name of this One-Touch Join Profile. Maximum length: 250 characters.",
			},
			"description": schema.StringAttribute{
				Computed: true,
				Optional: true,
				Default:  stringdefault.StaticString(""),
				Validators: []validator.String{
					stringvalidator.LengthAtMost(250),
				},
				MarkdownDescription: "An optional description of this One-Touch Join Profile. Maximum length: 250 characters.",
			},
			"display_upcoming_meetings": schema.Int64Attribute{
				Computed: true,
				Optional: true,
				Default:  int64default.StaticInt64(7),
				Validators: []validator.Int64{
					int64validator.Between(0, 365),
				},
				MarkdownDescription: "The number of days of upcoming One-Touch Join meetings to be shown on endpoints. Range: 0 to 365. Default: 7.",
			},
			"enable_non_video_meetings": schema.BoolAttribute{
				Computed:            true,
				Optional:            true,
				Default:             booldefault.StaticBool(true),
				MarkdownDescription: "When enabled, if the invitation has no valid video address the meeting will still appear on the endpoint as a scheduled meeting, but the Join button will not appear.",
			},
			"enable_private_meetings": schema.BoolAttribute{
				Computed:            true,
				Optional:            true,
				Default:             booldefault.StaticBool(true),
				MarkdownDescription: "Determines whether or not meetings flagged as private are processed by the OTJ service.",
			},
			"end_buffer": schema.Int64Attribute{
				Computed: true,
				Optional: true,
				Default:  int64default.StaticInt64(0),
				Validators: []validator.Int64{
					int64validator.Between(0, 180),
				},
				MarkdownDescription: "The number of minutes after the meeting's scheduled end time that the Join button on the endpoint will remain enabled. Range: 0 to 180. Default: 0.",
			},
			"start_buffer": schema.Int64Attribute{
				Computed: true,
				Optional: true,
				Default:  int64default.StaticInt64(5),
				Validators: []validator.Int64{
					int64validator.Between(0, 180),
				},
				MarkdownDescription: "The number of minutes before the meeting's scheduled start time that the Join button on the endpoint will become enabled. Range: 0 to 180. Default: 5.",
			},
			"ep_username": schema.StringAttribute{
				Computed: true,
				Optional: true,
				Default:  stringdefault.StaticString(""),
				Validators: []validator.String{
					stringvalidator.LengthAtMost(100),
				},
				MarkdownDescription: "The username used by OTJ to access a Cisco OBTP endpoint's API; only used if the endpoint's username is left blank. Maximum length: 100 characters.",
			},
			"ep_password": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(100),
				},
				MarkdownDescription: "The password used by OTJ to access a Cisco OBTP endpoint's API; only used if the endpoint's password is left blank. Maximum length: 100 characters.",
			},
			"ep_use_https": schema.BoolAttribute{
				Computed:            true,
				Optional:            true,
				Default:             booldefault.StaticBool(true),
				MarkdownDescription: "Whether or not to use HTTPS by default when accessing a Cisco OBTP endpoint's API. Can be overridden per endpoint.",
			},
			"ep_verify_certificate": schema.BoolAttribute{
				Computed:            true,
				Optional:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Whether or not to verify the TLS certificate of a Cisco OBTP endpoint by default when accessing its API. Can be overridden per endpoint.",
			},
			"exchange_deployment": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The OTJ Exchange Integration associated with this One-Touch Join Profile.",
			},
			"google_deployment": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The OTJ Google Workspace Integration associated with this One-Touch Join Profile.",
			},
			"graph_deployment": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The OTJ O365 Graph Integration associated with this One-Touch Join Profile.",
			},
			"process_alias_private_meetings": schema.BoolAttribute{
				Computed:            true,
				Optional:            true,
				Default:             booldefault.StaticBool(true),
				MarkdownDescription: "When enabled, the meeting alias for private meetings will be extracted from the invitation in the usual way. When disabled, the meeting alias will not appear on the endpoint and therefore the Join button will be disabled.",
			},
			"replace_empty_subject": schema.BoolAttribute{
				Computed:            true,
				Optional:            true,
				Default:             booldefault.StaticBool(true),
				MarkdownDescription: "For meetings that do not have a subject, use the organizer's name in place of the subject.",
			},
			"replace_subject_type": schema.StringAttribute{
				Computed: true,
				Optional: true,
				Default:  stringdefault.StaticString("PRIVATE"),
				Validators: []validator.String{
					stringvalidator.OneOf("ALL", "NEVER", "PRIVATE"),
				},
				MarkdownDescription: "Whether the meeting subject should be replaced. When enabled, the subject will be replaced with the name of the organizer unless you specify an alternative in the Replace subject string field. Valid values: `ALL`, `NEVER`, `PRIVATE`. Default: `PRIVATE`.",
			},
			"replace_subject_template": schema.StringAttribute{
				Computed: true,
				Optional: true,
				Default:  stringdefault.StaticString(""),
				Validators: []validator.String{
					stringvalidator.LengthAtMost(512),
				},
				MarkdownDescription: "A Jinja2 snippet that defines how the subject should be replaced (when this has been enabled). If this field is left blank, the subject will be replaced with the name of the organizer. Maximum length: 512 characters.",
			},
			"use_webex": schema.BoolAttribute{
				Computed:            true,
				Optional:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Enable OTJ to connect to Webex endpoints via Webex Cloud.",
			},
			"webex_api_domain": schema.StringAttribute{
				Computed: true,
				Optional: true,
				Default:  stringdefault.StaticString("webexapis.com"),
				Validators: []validator.String{
					stringvalidator.LengthAtMost(192),
				},
				MarkdownDescription: "The FQDN to use when connecting to the Webex API. Maximum length: 192 characters.",
			},
			"webex_client_id": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(100),
				},
				MarkdownDescription: "The Client ID that was generated when creating a Webex Integration for OTJ. Maximum length: 100 characters.",
			},
			"webex_client_secret": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				MarkdownDescription: "The Client Secret that was generated when creating a Webex Integration for OTJ.",
			},
			"webex_oauth_state": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The OAuth State parameter used to verify the OAuth endpoint's API.",
			},
			"webex_redirect_uri": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(255),
				},
				MarkdownDescription: "The redirect URI you entered when creating a Webex Integration for OTJ. It must be in the format 'https://[Management Node Address]/admin/platform/mjxintegration/oauth_redirect/'. Maximum length: 255 characters.",
			},
			"webex_refresh_token": schema.StringAttribute{
				Computed:            true,
				Sensitive:           true,
				MarkdownDescription: "The Webex Refresh token for your Webex integration.",
			},
			"endpoint_groups": schema.SetAttribute{
				Computed:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "The Endpoint Groups used by this One-Touch Join Profile.",
			},
		},
		MarkdownDescription: "Manages an MJX integration (One-Touch Join Profile) in Infinity. An MJX integration connects calendar deployments (Exchange, Google, or Graph) with endpoints to enable One-Touch Join functionality.",
	}
}

func (r *InfinityMjxIntegrationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	plan := &InfinityMjxIntegrationResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createRequest := &config.MjxIntegrationCreateRequest{
		Name:                        plan.Name.ValueString(),
		Description:                 plan.Description.ValueString(),
		DisplayUpcomingMeetings:     int(plan.DisplayUpcomingMeetings.ValueInt64()),
		EnableNonVideoMeetings:      plan.EnableNonVideoMeetings.ValueBool(),
		EnablePrivateMeetings:       plan.EnablePrivateMeetings.ValueBool(),
		EndBuffer:                   int(plan.EndBuffer.ValueInt64()),
		StartBuffer:                 int(plan.StartBuffer.ValueInt64()),
		EPUsername:                  plan.EPUsername.ValueString(),
		EPUseHTTPS:                  plan.EPUseHTTPS.ValueBool(),
		EPVerifyCertificate:         plan.EPVerifyCertificate.ValueBool(),
		ProcessAliasPrivateMeetings: plan.ProcessAliasPrivateMeetings.ValueBool(),
		ReplaceEmptySubject:         plan.ReplaceEmptySubject.ValueBool(),
		ReplaceSubjectType:          plan.ReplaceSubjectType.ValueString(),
		ReplaceSubjectTemplate:      plan.ReplaceSubjectTemplate.ValueString(),
		UseWebex:                    plan.UseWebex.ValueBool(),
		WebexAPIDomain:              plan.WebexAPIDomain.ValueString(),
	}

	if !plan.EPPassword.IsNull() && !plan.EPPassword.IsUnknown() {
		createRequest.EPPassword = plan.EPPassword.ValueString()
	}

	if !plan.ExchangeDeployment.IsNull() && !plan.ExchangeDeployment.IsUnknown() {
		v := plan.ExchangeDeployment.ValueString()
		createRequest.ExchangeDeployment = &v
	}

	if !plan.GoogleDeployment.IsNull() && !plan.GoogleDeployment.IsUnknown() {
		v := plan.GoogleDeployment.ValueString()
		createRequest.GoogleDeployment = &v
	}

	if !plan.GraphDeployment.IsNull() && !plan.GraphDeployment.IsUnknown() {
		v := plan.GraphDeployment.ValueString()
		createRequest.GraphDeployment = &v
	}

	if !plan.WebexClientID.IsNull() && !plan.WebexClientID.IsUnknown() {
		v := plan.WebexClientID.ValueString()
		createRequest.WebexClientID = &v
	}

	if !plan.WebexClientSecret.IsNull() && !plan.WebexClientSecret.IsUnknown() {
		v := plan.WebexClientSecret.ValueString()
		createRequest.WebexClientSecret = &v
	}

	if !plan.WebexOAuthState.IsNull() && !plan.WebexOAuthState.IsUnknown() {
		v := plan.WebexOAuthState.ValueString()
		createRequest.WebexOAuthState = &v
	}

	if !plan.WebexRedirectURI.IsNull() && !plan.WebexRedirectURI.IsUnknown() {
		v := plan.WebexRedirectURI.ValueString()
		createRequest.WebexRedirectURI = &v
	}

	createResponse, err := r.InfinityClient.Config().CreateMjxIntegration(ctx, createRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Infinity MJX integration",
			fmt.Sprintf("Could not create Infinity MJX integration: %s", err),
		)
		return
	}

	resourceID, err := createResponse.ResourceID()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Retrieving Infinity MJX integration ID",
			fmt.Sprintf("Could not retrieve ID for created Infinity MJX integration: %s", err),
		)
		return
	}

	model, err := r.read(ctx, resourceID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Created Infinity MJX integration",
			fmt.Sprintf("Could not read created Infinity MJX integration with ID %d: %s", resourceID, err),
		)
		return
	}

	// Preserve sensitive fields not returned by the API
	model.EPPassword = plan.EPPassword
	model.WebexClientSecret = plan.WebexClientSecret

	tflog.Trace(ctx, fmt.Sprintf("created Infinity MJX integration with ID: %s, name: %s", model.ID, model.Name))

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *InfinityMjxIntegrationResource) read(ctx context.Context, resourceID int) (*InfinityMjxIntegrationResourceModel, error) {
	var data InfinityMjxIntegrationResourceModel

	srv, err := r.InfinityClient.Config().GetMjxIntegration(ctx, resourceID)
	if err != nil {
		return nil, err
	}

	if srv.ResourceURI == "" {
		return nil, fmt.Errorf("MJX integration with ID %d not found", resourceID)
	}

	data.ID = types.StringValue(srv.ResourceURI)
	data.ResourceID = types.Int32Value(int32(resourceID)) // #nosec G115 -- API values are expected to be within int32 range
	data.Name = types.StringValue(srv.Name)
	data.Description = types.StringValue(srv.Description)
	data.DisplayUpcomingMeetings = types.Int64Value(int64(srv.DisplayUpcomingMeetings))
	data.EnableNonVideoMeetings = types.BoolValue(srv.EnableNonVideoMeetings)
	data.EnablePrivateMeetings = types.BoolValue(srv.EnablePrivateMeetings)
	data.EndBuffer = types.Int64Value(int64(srv.EndBuffer))
	data.StartBuffer = types.Int64Value(int64(srv.StartBuffer))
	data.EPUsername = types.StringValue(srv.EPUsername)
	data.EPUseHTTPS = types.BoolValue(srv.EPUseHTTPS)
	data.EPVerifyCertificate = types.BoolValue(srv.EPVerifyCertificate)
	data.ProcessAliasPrivateMeetings = types.BoolValue(srv.ProcessAliasPrivateMeetings)
	data.ReplaceEmptySubject = types.BoolValue(srv.ReplaceEmptySubject)
	data.ReplaceSubjectType = types.StringValue(srv.ReplaceSubjectType)
	data.ReplaceSubjectTemplate = types.StringValue(srv.ReplaceSubjectTemplate)
	data.UseWebex = types.BoolValue(srv.UseWebex)
	data.WebexAPIDomain = types.StringValue(srv.WebexAPIDomain)

	// Sensitive fields not returned by the API — preserved from plan/state
	data.EPPassword = types.StringNull()
	data.WebexClientSecret = types.StringNull()

	if srv.ExchangeDeployment != nil {
		data.ExchangeDeployment = types.StringValue(srv.ExchangeDeployment.ResourceURI)
	} else {
		data.ExchangeDeployment = types.StringNull()
	}

	if srv.GoogleDeployment != nil {
		data.GoogleDeployment = types.StringValue(srv.GoogleDeployment.ResourceURI)
	} else {
		data.GoogleDeployment = types.StringNull()
	}

	if srv.GraphDeployment != nil {
		data.GraphDeployment = types.StringValue(srv.GraphDeployment.ResourceURI)
	} else {
		data.GraphDeployment = types.StringNull()
	}

	if srv.WebexClientID != nil && *srv.WebexClientID != "" {
		data.WebexClientID = types.StringValue(*srv.WebexClientID)
	} else {
		data.WebexClientID = types.StringNull()
	}

	if srv.WebexOAuthState != nil && *srv.WebexOAuthState != "" {
		data.WebexOAuthState = types.StringValue(*srv.WebexOAuthState)
	} else {
		data.WebexOAuthState = types.StringNull()
	}

	if srv.WebexRedirectURI != nil && *srv.WebexRedirectURI != "" {
		data.WebexRedirectURI = types.StringValue(*srv.WebexRedirectURI)
	} else {
		data.WebexRedirectURI = types.StringNull()
	}

	if srv.WebexRefreshToken != nil && *srv.WebexRefreshToken != "" {
		data.WebexRefreshToken = types.StringValue(*srv.WebexRefreshToken)
	} else {
		data.WebexRefreshToken = types.StringNull()
	}

	if len(srv.EndpointGroups) > 0 {
		uris := make([]string, len(srv.EndpointGroups))
		for i, ref := range srv.EndpointGroups {
			uris[i] = ref.ResourceURI
		}
		groups, diags := types.SetValueFrom(ctx, types.StringType, uris)
		if diags.HasError() {
			return nil, fmt.Errorf("error converting endpoint groups: %v", diags)
		}
		data.EndpointGroups = groups
	} else {
		data.EndpointGroups = types.SetNull(types.StringType)
	}

	return &data, nil
}

func (r *InfinityMjxIntegrationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	state := &InfinityMjxIntegrationResourceModel{}

	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Preserve fields not returned consistently by the API
	epPassword := state.EPPassword
	webexClientSecret := state.WebexClientSecret
	webexRefreshToken := state.WebexRefreshToken

	resourceID := int(state.ResourceID.ValueInt32())
	state, err := r.read(ctx, resourceID)
	if err != nil {
		if isNotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading Infinity MJX integration",
			fmt.Sprintf("Could not read Infinity MJX integration: %s", err),
		)
		return
	}

	state.EPPassword = epPassword
	state.WebexClientSecret = webexClientSecret
	state.WebexRefreshToken = webexRefreshToken

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *InfinityMjxIntegrationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	plan := &InfinityMjxIntegrationResourceModel{}
	state := &InfinityMjxIntegrationResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateRequest := &config.MjxIntegrationUpdateRequest{
		Name:                        plan.Name.ValueString(),
		Description:                 plan.Description.ValueString(),
		DisplayUpcomingMeetings:     int(plan.DisplayUpcomingMeetings.ValueInt64()),
		EnableNonVideoMeetings:      plan.EnableNonVideoMeetings.ValueBool(),
		EnablePrivateMeetings:       plan.EnablePrivateMeetings.ValueBool(),
		EndBuffer:                   int(plan.EndBuffer.ValueInt64()),
		StartBuffer:                 int(plan.StartBuffer.ValueInt64()),
		EPUsername:                  plan.EPUsername.ValueString(),
		EPPassword:                  plan.EPPassword.ValueString(),
		EPUseHTTPS:                  plan.EPUseHTTPS.ValueBool(),
		EPVerifyCertificate:         plan.EPVerifyCertificate.ValueBool(),
		ProcessAliasPrivateMeetings: plan.ProcessAliasPrivateMeetings.ValueBool(),
		ReplaceEmptySubject:         plan.ReplaceEmptySubject.ValueBool(),
		ReplaceSubjectType:          plan.ReplaceSubjectType.ValueString(),
		ReplaceSubjectTemplate:      plan.ReplaceSubjectTemplate.ValueString(),
		UseWebex:                    plan.UseWebex.ValueBool(),
		WebexAPIDomain:              plan.WebexAPIDomain.ValueString(),
	}

	// Nullable fields: nil pointer sends JSON null to clear the value
	if !plan.ExchangeDeployment.IsNull() && !plan.ExchangeDeployment.IsUnknown() {
		v := plan.ExchangeDeployment.ValueString()
		updateRequest.ExchangeDeployment = &v
	}

	if !plan.GoogleDeployment.IsNull() && !plan.GoogleDeployment.IsUnknown() {
		v := plan.GoogleDeployment.ValueString()
		updateRequest.GoogleDeployment = &v
	}

	if !plan.GraphDeployment.IsNull() && !plan.GraphDeployment.IsUnknown() {
		v := plan.GraphDeployment.ValueString()
		updateRequest.GraphDeployment = &v
	}

	if !plan.WebexClientID.IsNull() && !plan.WebexClientID.IsUnknown() {
		v := plan.WebexClientID.ValueString()
		updateRequest.WebexClientID = &v
	}

	if !plan.WebexClientSecret.IsNull() && !plan.WebexClientSecret.IsUnknown() {
		v := plan.WebexClientSecret.ValueString()
		updateRequest.WebexClientSecret = &v
	}

	if !plan.WebexOAuthState.IsNull() && !plan.WebexOAuthState.IsUnknown() {
		v := plan.WebexOAuthState.ValueString()
		updateRequest.WebexOAuthState = &v
	}

	if !plan.WebexRedirectURI.IsNull() && !plan.WebexRedirectURI.IsUnknown() {
		v := plan.WebexRedirectURI.ValueString()
		updateRequest.WebexRedirectURI = &v
	}

	resourceID := int(state.ResourceID.ValueInt32())
	_, err := r.InfinityClient.Config().UpdateMjxIntegration(ctx, resourceID, updateRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Infinity MJX integration",
			fmt.Sprintf("Could not update Infinity MJX integration: %s", err),
		)
		return
	}

	model, err := r.read(ctx, resourceID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Updated Infinity MJX integration",
			fmt.Sprintf("Could not read updated Infinity MJX integration with ID %d: %s", resourceID, err),
		)
		return
	}

	// Preserve sensitive fields not returned by the API
	model.EPPassword = plan.EPPassword
	model.WebexClientSecret = plan.WebexClientSecret
	model.WebexRefreshToken = state.WebexRefreshToken

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *InfinityMjxIntegrationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	state := &InfinityMjxIntegrationResourceModel{}

	tflog.Info(ctx, "Deleting Infinity MJX integration")

	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.InfinityClient.Config().DeleteMjxIntegration(ctx, int(state.ResourceID.ValueInt32()))

	if err != nil && !isNotFoundError(err) && !isLookupError(err) {
		resp.Diagnostics.AddError(
			"Error Deleting Infinity MJX integration",
			fmt.Sprintf("Could not delete Infinity MJX integration with ID %s: %s", state.ID.ValueString(), err),
		)
		return
	}
}

func (r *InfinityMjxIntegrationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resourceID, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Import ID must be a valid integer for the resource ID. Got: %s", req.ID),
		)
		return
	}

	tflog.Trace(ctx, fmt.Sprintf("Importing Infinity MJX integration with resource ID: %d", resourceID))

	model, err := r.read(ctx, resourceID)
	if err != nil {
		if isNotFoundError(err) {
			resp.Diagnostics.AddError(
				"Infinity MJX Integration Not Found",
				fmt.Sprintf("Infinity MJX integration with resource ID %d not found.", resourceID),
			)
			return
		}
		resp.Diagnostics.AddError(
			"Error Importing Infinity MJX Integration",
			fmt.Sprintf("Could not import Infinity MJX integration with resource ID %d: %s", resourceID, err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}
