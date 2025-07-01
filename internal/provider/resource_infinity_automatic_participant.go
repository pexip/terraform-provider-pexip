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
	Conference          types.String `tfsdk:"conference"`
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
				MarkdownDescription: "Resource URI for the automatic participant in Infinity",
			},
			"resource_id": schema.Int32Attribute{
				Computed:            true,
				MarkdownDescription: "The resource integer identifier for the automatic participant in Infinity",
			},
			"alias": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(250),
				},
				MarkdownDescription: "The unique alias of the automatic participant. Maximum length: 250 characters.",
			},
			"description": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(250),
				},
				MarkdownDescription: "A description of the automatic participant. Maximum length: 250 characters.",
			},
			"conference": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(250),
				},
				MarkdownDescription: "The conference URI or reference. Maximum length: 250 characters.",
			},
			"protocol": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf("sip", "h323", "rtmp", "webrtc"),
				},
				MarkdownDescription: "The protocol for the automatic participant. Valid choices: sip, h323, rtmp, webrtc.",
			},
			"call_type": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf("audio", "video"),
				},
				MarkdownDescription: "The call type. Valid choices: audio, video.",
			},
			"role": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf("guest", "chair"),
				},
				MarkdownDescription: "The role of the automatic participant. Valid choices: guest, chair.",
			},
			"dtmf_sequence": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(250),
				},
				MarkdownDescription: "DTMF sequence to send when connecting. Maximum length: 250 characters.",
			},
			"keep_conference_alive": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf("keep_conference_alive", "end_conference_when_alone"),
				},
				MarkdownDescription: "Conference behavior when only this participant remains. Valid choices: keep_conference_alive, end_conference_when_alone.",
			},
			"routing": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf("auto", "manual"),
				},
				MarkdownDescription: "The routing type. Valid choices: auto, manual.",
			},
			"system_location": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Reference to system location resource URI.",
			},
			"streaming": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Whether streaming is enabled. Defaults to false.",
			},
			"remote_display_name": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(250),
				},
				MarkdownDescription: "The remote display name. Maximum length: 250 characters.",
			},
			"presentation_url": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(250),
				},
				MarkdownDescription: "The presentation URL. Maximum length: 250 characters.",
			},
			"creation_time": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The creation timestamp of the automatic participant.",
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

	createRequest := &config.AutomaticParticipantCreateRequest{
		Alias:               plan.Alias.ValueString(),
		Conference:          plan.Conference.ValueString(),
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
	if !plan.SystemLocation.IsNull() {
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

	if len(srv.ResourceURI) == 0 {
		return nil, fmt.Errorf("automatic participant with ID %d not found", resourceID)
	}

	data.ID = types.StringValue(srv.ResourceURI)
	data.ResourceID = types.Int32Value(int32(resourceID))
	data.Alias = types.StringValue(srv.Alias)
	data.Description = types.StringValue(srv.Description)
	data.Conference = types.StringValue(srv.Conference)
	data.Protocol = types.StringValue(srv.Protocol)
	data.CallType = types.StringValue(srv.CallType)
	data.Role = types.StringValue(srv.Role)
	data.DTMFSequence = types.StringValue(srv.DTMFSequence)
	data.KeepConferenceAlive = types.StringValue(srv.KeepConferenceAlive)
	data.Routing = types.StringValue(srv.Routing)
	if srv.SystemLocation != nil {
		data.SystemLocation = types.StringValue(*srv.SystemLocation)
	} else {
		data.SystemLocation = types.StringNull()
	}
	data.Streaming = types.BoolValue(srv.Streaming)
	data.RemoteDisplayName = types.StringValue(srv.RemoteDisplayName)
	data.PresentationURL = types.StringValue(srv.PresentationURL)
	data.CreationTime = types.StringValue(srv.CreationTime.String())

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

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resourceID := int(state.ResourceID.ValueInt32())

	updateRequest := &config.AutomaticParticipantUpdateRequest{
		Alias:               plan.Alias.ValueString(),
		Conference:          plan.Conference.ValueString(),
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
	if !plan.SystemLocation.IsNull() {
		systemLocation := plan.SystemLocation.ValueString()
		updateRequest.SystemLocation = &systemLocation
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
