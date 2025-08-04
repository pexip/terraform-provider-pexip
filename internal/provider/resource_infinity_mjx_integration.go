package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
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
				MarkdownDescription: "Resource URI for the MJX integration in Infinity",
			},
			"resource_id": schema.Int32Attribute{
				Computed:            true,
				MarkdownDescription: "The resource integer identifier for the MJX integration in Infinity",
			},
			"name": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					stringvalidator.LengthAtMost(100),
				},
				MarkdownDescription: "The name of the MJX integration. Maximum length: 100 characters.",
			},
			"description": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(500),
				},
				MarkdownDescription: "Description of the MJX integration. Maximum length: 500 characters.",
			},
			"display_upcoming_meetings": schema.Int64Attribute{
				Required: true,
				Validators: []validator.Int64{
					int64validator.Between(0, 24),
				},
				MarkdownDescription: "Number of hours ahead to display upcoming meetings. Valid range: 0-24.",
			},
			"enable_non_video_meetings": schema.BoolAttribute{
				Required:            true,
				MarkdownDescription: "Whether to enable non-video meetings in the integration.",
			},
			"enable_private_meetings": schema.BoolAttribute{
				Required:            true,
				MarkdownDescription: "Whether to enable private meetings in the integration.",
			},
			"end_buffer": schema.Int64Attribute{
				Required: true,
				Validators: []validator.Int64{
					int64validator.Between(0, 120),
				},
				MarkdownDescription: "End buffer time in minutes. Valid range: 0-120.",
			},
			"start_buffer": schema.Int64Attribute{
				Required: true,
				Validators: []validator.Int64{
					int64validator.Between(0, 120),
				},
				MarkdownDescription: "Start buffer time in minutes. Valid range: 0-120.",
			},
			"ep_username": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(100),
				},
				MarkdownDescription: "Endpoint username for authentication. Maximum length: 100 characters.",
			},
			"ep_password": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				MarkdownDescription: "Endpoint password for authentication. This field is sensitive.",
			},
			"ep_use_https": schema.BoolAttribute{
				Required:            true,
				MarkdownDescription: "Whether to use HTTPS for endpoint communication.",
			},
			"ep_verify_certificate": schema.BoolAttribute{
				Required:            true,
				MarkdownDescription: "Whether to verify SSL certificates for endpoint communication.",
			},
			"exchange_deployment": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Exchange deployment URI for Microsoft Exchange integration.",
			},
			"google_deployment": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Google deployment URI for Google Workspace integration.",
			},
			"graph_deployment": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Microsoft Graph deployment URI for Microsoft 365 integration.",
			},
			"process_alias_private_meetings": schema.BoolAttribute{
				Required:            true,
				MarkdownDescription: "Whether to process alias private meetings.",
			},
			"replace_empty_subject": schema.BoolAttribute{
				Required:            true,
				MarkdownDescription: "Whether to replace empty meeting subjects.",
			},
			"replace_subject_type": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf("none", "template", "alias"),
				},
				MarkdownDescription: "How to replace meeting subjects. Valid values: none, template, alias.",
			},
			"replace_subject_template": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(200),
				},
				MarkdownDescription: "Template for replacing meeting subjects. Maximum length: 200 characters.",
			},
			"use_webex": schema.BoolAttribute{
				Required:            true,
				MarkdownDescription: "Whether to enable Webex integration.",
			},
			"webex_api_domain": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(253),
				},
				MarkdownDescription: "Webex API domain. Maximum length: 253 characters.",
			},
			"webex_client_id": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(200),
				},
				MarkdownDescription: "Webex OAuth client ID. Maximum length: 200 characters.",
			},
			"webex_client_secret": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				MarkdownDescription: "Webex OAuth client secret. This field is sensitive.",
			},
			"webex_oauth_state": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Webex OAuth state (computed).",
			},
			"webex_redirect_uri": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(500),
				},
				MarkdownDescription: "Webex OAuth redirect URI. Maximum length: 500 characters.",
			},
			"webex_refresh_token": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				MarkdownDescription: "Webex OAuth refresh token. This field is sensitive.",
			},
			"endpoint_groups": schema.SetAttribute{
				ElementType:         types.StringType,
				Computed:            true,
				MarkdownDescription: "List of endpoint group URIs associated with this integration.",
			},
		},
		MarkdownDescription: "Manages an MJX integration with the Infinity service. MJX integrations provide comprehensive Microsoft Teams integration capabilities, including calendar synchronization, endpoint management, and multi-cloud deployment support for Exchange, Google Workspace, and Webex.",
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

	// Handle optional pointer fields
	if !plan.ExchangeDeployment.IsNull() && !plan.ExchangeDeployment.IsUnknown() {
		deployment := plan.ExchangeDeployment.ValueString()
		createRequest.ExchangeDeployment = &deployment
	}

	if !plan.GoogleDeployment.IsNull() && !plan.GoogleDeployment.IsUnknown() {
		deployment := plan.GoogleDeployment.ValueString()
		createRequest.GoogleDeployment = &deployment
	}

	if !plan.GraphDeployment.IsNull() && !plan.GraphDeployment.IsUnknown() {
		deployment := plan.GraphDeployment.ValueString()
		createRequest.GraphDeployment = &deployment
	}

	if !plan.WebexClientID.IsNull() && !plan.WebexClientID.IsUnknown() {
		clientID := plan.WebexClientID.ValueString()
		createRequest.WebexClientID = &clientID
	}

	if !plan.WebexClientSecret.IsNull() && !plan.WebexClientSecret.IsUnknown() {
		clientSecret := plan.WebexClientSecret.ValueString()
		createRequest.WebexClientSecret = &clientSecret
	}

	if !plan.WebexRedirectURI.IsNull() && !plan.WebexRedirectURI.IsUnknown() {
		redirectURI := plan.WebexRedirectURI.ValueString()
		createRequest.WebexRedirectURI = &redirectURI
	}

	if !plan.WebexRefreshToken.IsNull() && !plan.WebexRefreshToken.IsUnknown() {
		refreshToken := plan.WebexRefreshToken.ValueString()
		createRequest.WebexRefreshToken = &refreshToken
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

	// Read the state from the API to get all computed values
	model, err := r.read(ctx, resourceID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Created Infinity MJX integration",
			fmt.Sprintf("Could not read created Infinity MJX integration with ID %d: %s", resourceID, err),
		)
		return
	}
	tflog.Trace(ctx, fmt.Sprintf("created Infinity MJX integration with ID: %s, name: %s", model.ID, model.Name))

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *InfinityMjxIntegrationResource) read(ctx context.Context, resourceID int) (*InfinityMjxIntegrationResourceModel, error) {
	var data InfinityMjxIntegrationResourceModel

	srv, err := r.InfinityClient.Config().GetMjxIntegration(ctx, resourceID)
	if err != nil {
		return nil, err
	}

	if len(srv.ResourceURI) == 0 {
		return nil, fmt.Errorf("MJX integration with ID %d not found", resourceID)
	}

	data.ID = types.StringValue(srv.ResourceURI)
	data.ResourceID = types.Int32Value(int32(resourceID))
	data.Name = types.StringValue(srv.Name)
	data.Description = types.StringValue(srv.Description)
	data.DisplayUpcomingMeetings = types.Int64Value(int64(srv.DisplayUpcomingMeetings))
	data.EnableNonVideoMeetings = types.BoolValue(srv.EnableNonVideoMeetings)
	data.EnablePrivateMeetings = types.BoolValue(srv.EnablePrivateMeetings)
	data.EndBuffer = types.Int64Value(int64(srv.EndBuffer))
	data.StartBuffer = types.Int64Value(int64(srv.StartBuffer))
	data.EPUsername = types.StringValue(srv.EPUsername)
	data.EPPassword = types.StringValue(srv.EPPassword)
	data.EPUseHTTPS = types.BoolValue(srv.EPUseHTTPS)
	data.EPVerifyCertificate = types.BoolValue(srv.EPVerifyCertificate)
	data.ProcessAliasPrivateMeetings = types.BoolValue(srv.ProcessAliasPrivateMeetings)
	data.ReplaceEmptySubject = types.BoolValue(srv.ReplaceEmptySubject)
	data.ReplaceSubjectType = types.StringValue(srv.ReplaceSubjectType)
	data.ReplaceSubjectTemplate = types.StringValue(srv.ReplaceSubjectTemplate)
	data.UseWebex = types.BoolValue(srv.UseWebex)
	data.WebexAPIDomain = types.StringValue(srv.WebexAPIDomain)

	// Handle optional pointer fields
	if srv.ExchangeDeployment != nil {
		data.ExchangeDeployment = types.StringValue(*srv.ExchangeDeployment)
	} else {
		data.ExchangeDeployment = types.StringNull()
	}

	if srv.GoogleDeployment != nil {
		data.GoogleDeployment = types.StringValue(*srv.GoogleDeployment)
	} else {
		data.GoogleDeployment = types.StringNull()
	}

	if srv.GraphDeployment != nil {
		data.GraphDeployment = types.StringValue(*srv.GraphDeployment)
	} else {
		data.GraphDeployment = types.StringNull()
	}

	if srv.WebexClientID != nil {
		data.WebexClientID = types.StringValue(*srv.WebexClientID)
	} else {
		data.WebexClientID = types.StringNull()
	}

	if srv.WebexClientSecret != nil {
		data.WebexClientSecret = types.StringValue(*srv.WebexClientSecret)
	} else {
		data.WebexClientSecret = types.StringNull()
	}

	if srv.WebexOAuthState != nil {
		data.WebexOAuthState = types.StringValue(*srv.WebexOAuthState)
	} else {
		data.WebexOAuthState = types.StringNull()
	}

	if srv.WebexRedirectURI != nil {
		data.WebexRedirectURI = types.StringValue(*srv.WebexRedirectURI)
	} else {
		data.WebexRedirectURI = types.StringNull()
	}

	if srv.WebexRefreshToken != nil {
		data.WebexRefreshToken = types.StringValue(*srv.WebexRefreshToken)
	} else {
		data.WebexRefreshToken = types.StringNull()
	}

	// Handle list field
	if srv.EndpointGroups != nil {
		endpointGroupsSet, diags := types.SetValueFrom(ctx, types.StringType, srv.EndpointGroups)
		if diags.HasError() {
			return nil, fmt.Errorf("failed to convert endpoint groups: %s", diags.Errors())
		}
		data.EndpointGroups = endpointGroupsSet
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

	resourceID := int(state.ResourceID.ValueInt32())
	state, err := r.read(ctx, resourceID)
	if err != nil {
		// Check if the error is a 404 (not found)
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
		Name:                   plan.Name.ValueString(),
		Description:            plan.Description.ValueString(),
		EPUsername:             plan.EPUsername.ValueString(),
		EPPassword:             plan.EPPassword.ValueString(),
		ReplaceSubjectType:     plan.ReplaceSubjectType.ValueString(),
		ReplaceSubjectTemplate: plan.ReplaceSubjectTemplate.ValueString(),
		WebexAPIDomain:         plan.WebexAPIDomain.ValueString(),
	}

	// Handle optional pointer fields
	if !plan.DisplayUpcomingMeetings.IsNull() && !plan.DisplayUpcomingMeetings.IsUnknown() {
		meetings := int(plan.DisplayUpcomingMeetings.ValueInt64())
		updateRequest.DisplayUpcomingMeetings = &meetings
	}

	if !plan.EnableNonVideoMeetings.IsNull() && !plan.EnableNonVideoMeetings.IsUnknown() {
		enable := plan.EnableNonVideoMeetings.ValueBool()
		updateRequest.EnableNonVideoMeetings = &enable
	}

	if !plan.EnablePrivateMeetings.IsNull() && !plan.EnablePrivateMeetings.IsUnknown() {
		enable := plan.EnablePrivateMeetings.ValueBool()
		updateRequest.EnablePrivateMeetings = &enable
	}

	if !plan.EndBuffer.IsNull() && !plan.EndBuffer.IsUnknown() {
		buffer := int(plan.EndBuffer.ValueInt64())
		updateRequest.EndBuffer = &buffer
	}

	if !plan.StartBuffer.IsNull() && !plan.StartBuffer.IsUnknown() {
		buffer := int(plan.StartBuffer.ValueInt64())
		updateRequest.StartBuffer = &buffer
	}

	if !plan.EPUseHTTPS.IsNull() && !plan.EPUseHTTPS.IsUnknown() {
		useHTTPS := plan.EPUseHTTPS.ValueBool()
		updateRequest.EPUseHTTPS = &useHTTPS
	}

	if !plan.EPVerifyCertificate.IsNull() && !plan.EPVerifyCertificate.IsUnknown() {
		verify := plan.EPVerifyCertificate.ValueBool()
		updateRequest.EPVerifyCertificate = &verify
	}

	if !plan.ExchangeDeployment.IsNull() && !plan.ExchangeDeployment.IsUnknown() {
		deployment := plan.ExchangeDeployment.ValueString()
		updateRequest.ExchangeDeployment = &deployment
	}

	if !plan.GoogleDeployment.IsNull() && !plan.GoogleDeployment.IsUnknown() {
		deployment := plan.GoogleDeployment.ValueString()
		updateRequest.GoogleDeployment = &deployment
	}

	if !plan.GraphDeployment.IsNull() && !plan.GraphDeployment.IsUnknown() {
		deployment := plan.GraphDeployment.ValueString()
		updateRequest.GraphDeployment = &deployment
	}

	if !plan.ProcessAliasPrivateMeetings.IsNull() && !plan.ProcessAliasPrivateMeetings.IsUnknown() {
		process := plan.ProcessAliasPrivateMeetings.ValueBool()
		updateRequest.ProcessAliasPrivateMeetings = &process
	}

	if !plan.ReplaceEmptySubject.IsNull() && !plan.ReplaceEmptySubject.IsUnknown() {
		replace := plan.ReplaceEmptySubject.ValueBool()
		updateRequest.ReplaceEmptySubject = &replace
	}

	if !plan.UseWebex.IsNull() && !plan.UseWebex.IsUnknown() {
		useWebex := plan.UseWebex.ValueBool()
		updateRequest.UseWebex = &useWebex
	}

	if !plan.WebexClientID.IsNull() && !plan.WebexClientID.IsUnknown() {
		clientID := plan.WebexClientID.ValueString()
		updateRequest.WebexClientID = &clientID
	}

	if !plan.WebexClientSecret.IsNull() && !plan.WebexClientSecret.IsUnknown() {
		clientSecret := plan.WebexClientSecret.ValueString()
		updateRequest.WebexClientSecret = &clientSecret
	}

	if !plan.WebexRedirectURI.IsNull() && !plan.WebexRedirectURI.IsUnknown() {
		redirectURI := plan.WebexRedirectURI.ValueString()
		updateRequest.WebexRedirectURI = &redirectURI
	}

	if !plan.WebexRefreshToken.IsNull() && !plan.WebexRefreshToken.IsUnknown() {
		refreshToken := plan.WebexRefreshToken.ValueString()
		updateRequest.WebexRefreshToken = &refreshToken
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

	// Read the state from the API to get all computed values
	model, err := r.read(ctx, resourceID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Updated Infinity MJX integration",
			fmt.Sprintf("Could not read updated Infinity MJX integration with ID %d: %s", resourceID, err),
		)
		return
	}

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

	// Ignore 404 Not Found and Lookup errors on delete
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

	// Read the resource from the API
	model, err := r.read(ctx, resourceID)
	if err != nil {
		// Check if the error is a 404 (not found)
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

	// Set the state from the imported resource
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}
