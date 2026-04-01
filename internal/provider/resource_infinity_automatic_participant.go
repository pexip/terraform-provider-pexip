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
	_ resource.ResourceWithImportState = (*InfinityAutomaticParticipantResource)(nil)
)

type InfinityAutomaticParticipantResource struct {
	InfinityClient InfinityClient
}

type InfinityAutomaticParticipantResourceModel struct {
	ID                  types.String `tfsdk:"id"`
	ResourceID          types.Int32  `tfsdk:"resource_id"`
	Alias               types.String `tfsdk:"alias"`
	Description         types.String `tfsdk:"description"`
	Conference          types.Set    `tfsdk:"conference"`
	Protocol            types.String `tfsdk:"protocol"`
	CallType            types.String `tfsdk:"call_type"`
	Role                types.String `tfsdk:"role"`
	DTMFSequence        types.String `tfsdk:"dtmf_sequence"`
	KeepConferenceAlive types.String `tfsdk:"keep_conference_alive"`
	Routing             types.String `tfsdk:"routing"`
	SystemLocation      types.String `tfsdk:"system_location"`
	Streaming           types.Bool   `tfsdk:"streaming"`
	RemoteDisplayName   types.String `tfsdk:"remote_display_name"`
	PresentationURL     types.String `tfsdk:"presentation_url"`
	CreationTime        types.String `tfsdk:"creation_time"`
}

func (r *InfinityAutomaticParticipantResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_infinity_automatic_participant"
}

func (r *InfinityAutomaticParticipantResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *InfinityAutomaticParticipantResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Resource URI for the automatic participant in Infinity.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"resource_id": schema.Int32Attribute{
				Computed:            true,
				MarkdownDescription: "The resource integer identifier for the automatic participant in Infinity",
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.UseStateForUnknown(),
				},
			},
			"alias": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(250),
				},
				MarkdownDescription: "The alias of the participant that is to be dialed when the conference starts. Maximum length: 250 characters.",
			},
			"description": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(250),
				},
				MarkdownDescription: "An optional description of the Automatically Dialed Participant. Maximum length: 250 characters.",
			},
			"conference": schema.SetAttribute{
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "The conference to which the Automatically Dialed Participant belongs.",
			},
			"protocol": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("sip"),
				Validators: []validator.String{
					stringvalidator.OneOf("teams", "gms", "h323", "sip", "mssip", "rtmp"),
				},
				MarkdownDescription: "The protocol to use when dialing the participant. Note that if the call is to a registered device, Pexip Infinity will instead use the protocol that the device used to make the registration. Valid choices: teams, gms, h323, sip, mssip, rtmp.",
			},
			"call_type": schema.StringAttribute{
				Computed: true,
				Optional: true,
				Default:  stringdefault.StaticString("video"),
				Validators: []validator.String{
					stringvalidator.OneOf("audio", "video", "video-only"),
				},
				MarkdownDescription: "Maximum media content of the call. The participant being called will not be able to escalate beyond the selected capability. Valid choices: audio, video, video-only.",
			},
			"role": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("guest"),
				Validators: []validator.String{
					stringvalidator.OneOf("guest", "chair"),
				},
				MarkdownDescription: "The level of privileges the participant will have in the conference. host: The participant will have full privileges. guest: The participant will have restricted privileges. Valid choices: guest, chair.",
			},
			"dtmf_sequence": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(250),
				},
				MarkdownDescription: "The DTMF sequence to be transmitted after the call to the automatically dialed participant starts. Insert a comma for a 2 second pause. Maximum length: 250 characters.",
			},
			"keep_conference_alive": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("keep_conference_alive_if_multiple"),
				Validators: []validator.String{
					stringvalidator.OneOf("keep_conference_alive", "keep_conference_alive_if_multiple", "keep_conference_alive_never"),
				},
				MarkdownDescription: "Determines whether the conference will continue when all other participants have disconnected. Yes: the conference will continue to run until this participant has disconnected (applies to Hosts only). If multiple: the conference will continue to run as long as there are two or more If multiple participants and at least one of them is a Host. No: the conference will be terminated automatically if this is the only remaining participant.",
			},
			"routing": schema.StringAttribute{
				Computed: true,
				Optional: true,
				Default:  stringdefault.StaticString("manual"),
				Validators: []validator.String{
					stringvalidator.OneOf("manual", "routing_rule"),
				},
				MarkdownDescription: "Route this call manually using the defaults for the specified location - or route this call automatically using Call Routing Rules. Valid choices: manual, routing_rule.",
			},
			"system_location": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "For manually routed Automatically Dialed Participants (ADPs), this is the location of the Conferencing Node from which the call to the ADP will be initiated. For automatically routed ADPs, this is the notional source location used when considering if a routing rule applies or not - however the routing rule itself determines the location of the node that dials the ADP. To allow Pexip Infinity to automatically select the Conferencing Node to initiate the outgoing call, select Automatic.",
			},
			"streaming": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Identify the dialed participant as a streaming or recording device.",
			},
			"remote_display_name": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(250),
				},
				MarkdownDescription: "The optional user-facing display name for this participant, which will be shown in the participant lists. Maximum length: 250 characters.",
			},
			"presentation_url": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(250),
				},
				MarkdownDescription: "The optional RTMP URL for the second (presentation) stream. Maximum length: 250 characters.",
			},
			"creation_time": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The time at which the Automatically Dialed Participant was created.",
			},
		},
		MarkdownDescription: "Manages an automatic participant configuration with the Infinity service.",
	}
}

func (r *InfinityAutomaticParticipantResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	plan := &InfinityAutomaticParticipantResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	conference, diags := getStringList(ctx, plan.Conference)
	resp.Diagnostics.Append(diags...)

	createRequest := &config.AutomaticParticipantCreateRequest{
		Alias:               plan.Alias.ValueString(),
		Conference:          conference,
		Protocol:            plan.Protocol.ValueString(),
		CallType:            plan.CallType.ValueString(),
		Role:                plan.Role.ValueString(),
		KeepConferenceAlive: plan.KeepConferenceAlive.ValueString(),
		Routing:             plan.Routing.ValueString(),
		Streaming:           plan.Streaming.ValueBool(),
	}

	// Set optional fields
	if !plan.Description.IsNull() {
		createRequest.Description = plan.Description.ValueString()
	}
	if !plan.DTMFSequence.IsNull() {
		createRequest.DTMFSequence = plan.DTMFSequence.ValueString()
	}
	if !plan.SystemLocation.IsNull() && plan.SystemLocation.ValueString() != "" {
		systemLocation := plan.SystemLocation.ValueString()
		createRequest.SystemLocation = &systemLocation
	}
	if !plan.RemoteDisplayName.IsNull() {
		createRequest.RemoteDisplayName = plan.RemoteDisplayName.ValueString()
	}
	if !plan.PresentationURL.IsNull() {
		createRequest.PresentationURL = plan.PresentationURL.ValueString()
	}

	createResponse, err := r.InfinityClient.Config().CreateAutomaticParticipant(ctx, createRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Infinity automatic participant",
			fmt.Sprintf("Could not create Infinity automatic participant: %s", err),
		)
		return
	}

	resourceID, err := createResponse.ResourceID()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Retrieving Infinity automatic participant ID",
			fmt.Sprintf("Could not retrieve ID for created Infinity automatic participant: %s", err),
		)
		return
	}

	// Read the state from the API to get all computed values
	model, err := r.read(ctx, resourceID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Created Infinity automatic participant",
			fmt.Sprintf("Could not read created Infinity automatic participant with ID %d: %s", resourceID, err),
		)
		return
	}
	tflog.Trace(ctx, fmt.Sprintf("created Infinity automatic participant with ID: %s, alias: %s", model.ID, model.Alias))

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *InfinityAutomaticParticipantResource) read(ctx context.Context, resourceID int) (*InfinityAutomaticParticipantResourceModel, error) {
	var data InfinityAutomaticParticipantResourceModel

	srv, err := r.InfinityClient.Config().GetAutomaticParticipant(ctx, resourceID)
	if err != nil {
		return nil, err
	}

	if srv.ResourceURI == "" {
		return nil, fmt.Errorf("automatic participant with ID %d not found", resourceID)
	}

	data.ID = types.StringValue(srv.ResourceURI)
	data.ResourceID = types.Int32Value(int32(resourceID)) // #nosec G115 -- API values are expected to be within int32 range
	data.Alias = types.StringValue(srv.Alias)
	data.Description = types.StringValue(srv.Description)
	data.Protocol = types.StringValue(srv.Protocol)
	data.CallType = types.StringValue(srv.CallType)
	data.Role = types.StringValue(srv.Role)
	data.DTMFSequence = types.StringValue(srv.DTMFSequence)
	data.KeepConferenceAlive = types.StringValue(srv.KeepConferenceAlive)
	data.Routing = types.StringValue(srv.Routing)
	data.Streaming = types.BoolValue(srv.Streaming)
	data.RemoteDisplayName = types.StringValue(srv.RemoteDisplayName)
	data.PresentationURL = types.StringValue(srv.PresentationURL)
	data.CreationTime = types.StringValue(srv.CreationTime.String())
	if srv.SystemLocation != nil {
		data.SystemLocation = types.StringValue(*srv.SystemLocation)
	} else {
		data.SystemLocation = types.StringNull()
	}

	var conferences []string
	for _, c := range srv.Conference {
		conferences = append(conferences, c)
	}
	confSetValue, diags := types.SetValueFrom(ctx, types.StringType, conferences)
	if diags.HasError() {
		return nil, fmt.Errorf("error converting conferences: %v", diags)
	}
	data.Conference = confSetValue

	return &data, nil
}

func (r *InfinityAutomaticParticipantResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	state := &InfinityAutomaticParticipantResourceModel{}

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
			"Error Reading Infinity automatic participant",
			fmt.Sprintf("Could not read Infinity automatic participant: %s", err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *InfinityAutomaticParticipantResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	plan := &InfinityAutomaticParticipantResourceModel{}
	state := &InfinityAutomaticParticipantResourceModel{}
	rawConfig := &InfinityAutomaticParticipantResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	resp.Diagnostics.Append(req.Config.Get(ctx, rawConfig)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resourceID := int(state.ResourceID.ValueInt32())

	conference, diags := getStringList(ctx, plan.Conference)
	resp.Diagnostics.Append(diags...)

	updateRequest := &config.AutomaticParticipantUpdateRequest{
		Alias:               plan.Alias.ValueString(),
		Conference:          conference,
		Protocol:            plan.Protocol.ValueString(),
		CallType:            plan.CallType.ValueString(),
		Role:                plan.Role.ValueString(),
		KeepConferenceAlive: plan.KeepConferenceAlive.ValueString(),
		Routing:             plan.Routing.ValueString(),
	}

	// Set boolean pointer field for update
	if !plan.Streaming.IsNull() {
		streaming := plan.Streaming.ValueBool()
		updateRequest.Streaming = &streaming
	}

	// Set optional fields
	if !plan.Description.IsNull() {
		updateRequest.Description = plan.Description.ValueString()
	}
	if !plan.DTMFSequence.IsNull() {
		updateRequest.DTMFSequence = plan.DTMFSequence.ValueString()
	}
	// SystemLocation must always be set (even to nil) to allow clearing
	// Check the config to see if the user specified it, rather than the plan which might have computed values
	if !rawConfig.SystemLocation.IsNull() && rawConfig.SystemLocation.ValueString() != "" {
		systemLocation := plan.SystemLocation.ValueString()
		updateRequest.SystemLocation = &systemLocation
	} else {
		updateRequest.SystemLocation = nil
	}
	if !plan.RemoteDisplayName.IsNull() {
		updateRequest.RemoteDisplayName = plan.RemoteDisplayName.ValueString()
	}
	if !plan.PresentationURL.IsNull() {
		updateRequest.PresentationURL = plan.PresentationURL.ValueString()
	}

	_, err := r.InfinityClient.Config().UpdateAutomaticParticipant(ctx, resourceID, updateRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Infinity automatic participant",
			fmt.Sprintf("Could not update Infinity automatic participant with ID %d: %s", resourceID, err),
		)
		return
	}

	// Re-read the resource to get the latest state
	updatedModel, err := r.read(ctx, resourceID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Updated Infinity automatic participant",
			fmt.Sprintf("Could not read updated Infinity automatic participant with ID %d: %s", resourceID, err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, updatedModel)...)
}

func (r *InfinityAutomaticParticipantResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	state := &InfinityAutomaticParticipantResourceModel{}

	tflog.Info(ctx, "Deleting Infinity automatic participant")

	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.InfinityClient.Config().DeleteAutomaticParticipant(ctx, int(state.ResourceID.ValueInt32()))

	// Ignore 404 Not Found and Lookup errors on delete
	if err != nil && !isNotFoundError(err) && !isLookupError(err) {
		resp.Diagnostics.AddError(
			"Error Deleting Infinity automatic participant",
			fmt.Sprintf("Could not delete Infinity automatic participant with ID %s: %s", state.ID.ValueString(), err),
		)
		return
	}
}

func (r *InfinityAutomaticParticipantResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resourceID, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Import ID must be a valid integer for the resource ID. Got: %s", req.ID),
		)
		return
	}

	tflog.Trace(ctx, fmt.Sprintf("Importing Infinity automatic participant with resource ID: %d", resourceID))

	// Read the resource from the API
	model, err := r.read(ctx, resourceID)
	if err != nil {
		// Check if the error is a 404 (not found)
		if isNotFoundError(err) {
			resp.Diagnostics.AddError(
				"Infinity Automatic Participant Not Found",
				fmt.Sprintf("Infinity automatic participant with resource ID %d not found.", resourceID),
			)
			return
		}
		resp.Diagnostics.AddError(
			"Error Importing Infinity Automatic Participant",
			fmt.Sprintf("Could not import Infinity automatic participant with resource ID %d: %s", resourceID, err),
		)
		return
	}

	// Set the state from the imported resource
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}
