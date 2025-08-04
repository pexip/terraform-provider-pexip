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
	_ resource.ResourceWithImportState = (*InfinityConferenceResource)(nil)
)

type InfinityConferenceResource struct {
	InfinityClient InfinityClient
}

type InfinityConferenceResourceModel struct {
	ID                 types.String `tfsdk:"id"`
	ResourceID         types.Int32  `tfsdk:"resource_id"`
	Name               types.String `tfsdk:"name"`
	Description        types.String `tfsdk:"description"`
	ServiceType        types.String `tfsdk:"service_type"`
	PIN                types.String `tfsdk:"pin"`
	GuestPIN           types.String `tfsdk:"guest_pin"`
	AllowGuests        types.Bool   `tfsdk:"allow_guests"`
	GuestsMuted        types.Bool   `tfsdk:"guests_muted"`
	HostsCanUnmute     types.Bool   `tfsdk:"hosts_can_unmute"`
	MaxPixelsPerSecond types.Int32  `tfsdk:"max_pixels_per_second"`
	Tag                types.String `tfsdk:"tag"`
}

func (r *InfinityConferenceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_infinity_conference"
}

func (r *InfinityConferenceResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *InfinityConferenceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Resource URI for the conference in Infinity",
			},
			"resource_id": schema.Int32Attribute{
				Computed:            true,
				MarkdownDescription: "The resource integer identifier for the conference in Infinity",
			},
			"name": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(250),
				},
				MarkdownDescription: "The unique name used to refer to this conference. Maximum length: 250 characters.",
			},
			"description": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(250),
				},
				MarkdownDescription: "A description of the conference. Maximum length: 250 characters.",
			},
			"service_type": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf("conference", "lecture", "two_stage_dialing", "test_call", "media_playback"),
				},
				MarkdownDescription: "The type of conferencing service. Valid choices: conference, lecture, two_stage_dialing, test_call, media_playback.",
			},
			"pin": schema.StringAttribute{
				Optional:  true,
				Computed:  true,
				Sensitive: true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(4, 20),
				},
				MarkdownDescription: "Secure access code for participants. Length: 4-20 digits, including any terminal #.",
			},
			"guest_pin": schema.StringAttribute{
				Optional:  true,
				Computed:  true,
				Sensitive: true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(4, 20),
				},
				MarkdownDescription: "Optional secure access code for Guest participants. Length: 4-20 digits.",
			},
			"allow_guests": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Whether Guest participants are allowed to join. Defaults to false.",
			},
			"guests_muted": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Whether Guest participants are muted by default. Defaults to false.",
			},
			"hosts_can_unmute": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Whether Host participants can unmute Guest participants. Defaults to false.",
			},
			"max_pixels_per_second": schema.Int32Attribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Maximum pixels per second for video quality.",
			},
			"tag": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(250),
				},
				MarkdownDescription: "A unique identifier used to track usage. Maximum length: 250 characters.",
			},
		},
		MarkdownDescription: "Manages a conference configuration with the Infinity service.",
	}
}

func (r *InfinityConferenceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	plan := &InfinityConferenceResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createRequest := &config.ConferenceCreateRequest{
		Name:        plan.Name.ValueString(),
		ServiceType: plan.ServiceType.ValueString(),
	}

	// Set optional string fields
	if !plan.Description.IsNull() {
		createRequest.Description = plan.Description.ValueString()
	}
	if !plan.PIN.IsNull() {
		createRequest.PIN = plan.PIN.ValueString()
	}
	if !plan.GuestPIN.IsNull() {
		createRequest.GuestPIN = plan.GuestPIN.ValueString()
	}
	if !plan.Tag.IsNull() {
		createRequest.Tag = plan.Tag.ValueString()
	}

	// Set boolean fields (required fields need default values)
	if !plan.AllowGuests.IsNull() {
		createRequest.AllowGuests = plan.AllowGuests.ValueBool()
	} else {
		createRequest.AllowGuests = false
	}

	if !plan.GuestsMuted.IsNull() {
		createRequest.GuestsMuted = plan.GuestsMuted.ValueBool()
	} else {
		createRequest.GuestsMuted = false
	}

	if !plan.HostsCanUnmute.IsNull() {
		createRequest.HostsCanUnmute = plan.HostsCanUnmute.ValueBool()
	} else {
		createRequest.HostsCanUnmute = false
	}

	// Set optional integer fields
	if !plan.MaxPixelsPerSecond.IsNull() {
		createRequest.MaxPixelsPerSecond = int(plan.MaxPixelsPerSecond.ValueInt32())
	}

	createResponse, err := r.InfinityClient.Config().CreateConference(ctx, createRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Infinity conference",
			fmt.Sprintf("Could not create Infinity conference: %s", err),
		)
		return
	}

	resourceID, err := createResponse.ResourceID()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Retrieving Infinity conference ID",
			fmt.Sprintf("Could not retrieve ID for created Infinity conference: %s", err),
		)
		return
	}

	// Read the state from the API to get all computed values
	model, err := r.read(ctx, resourceID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Created Infinity conference",
			fmt.Sprintf("Could not read created Infinity conference with ID %d: %s", resourceID, err),
		)
		return
	}
	tflog.Trace(ctx, fmt.Sprintf("created Infinity conference with ID: %s, name: %s", model.ID, model.Name))

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *InfinityConferenceResource) read(ctx context.Context, resourceID int) (*InfinityConferenceResourceModel, error) {
	var data InfinityConferenceResourceModel

	srv, err := r.InfinityClient.Config().GetConference(ctx, resourceID)
	if err != nil {
		return nil, err
	}

	if len(srv.ResourceURI) == 0 {
		return nil, fmt.Errorf("conference with ID %d not found", resourceID)
	}

	data.ID = types.StringValue(srv.ResourceURI)
	data.ResourceID = types.Int32Value(int32(resourceID))
	data.Name = types.StringValue(srv.Name)
	data.Description = types.StringValue(srv.Description)
	data.ServiceType = types.StringValue(srv.ServiceType)
	data.PIN = types.StringValue(srv.PIN)
	data.GuestPIN = types.StringValue(srv.GuestPIN)
	data.Tag = types.StringValue(srv.Tag)

	// Set boolean fields
	data.AllowGuests = types.BoolValue(srv.AllowGuests)
	data.GuestsMuted = types.BoolValue(srv.GuestsMuted)
	data.HostsCanUnmute = types.BoolValue(srv.HostsCanUnmute)

	// Set integer fields
	if srv.MaxPixelsPerSecond != 0 {
		data.MaxPixelsPerSecond = types.Int32Value(int32(srv.MaxPixelsPerSecond))
	} else {
		data.MaxPixelsPerSecond = types.Int32Null()
	}

	return &data, nil
}

func (r *InfinityConferenceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	state := &InfinityConferenceResourceModel{}

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
			"Error Reading Infinity conference",
			fmt.Sprintf("Could not read Infinity conference: %s", err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *InfinityConferenceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	plan := &InfinityConferenceResourceModel{}
	state := &InfinityConferenceResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resourceID := int(state.ResourceID.ValueInt32())

	updateRequest := &config.ConferenceUpdateRequest{
		Name: plan.Name.ValueString(),
	}

	// Set optional string fields
	if !plan.Description.IsNull() {
		updateRequest.Description = plan.Description.ValueString()
	}
	if !plan.PIN.IsNull() {
		updateRequest.PIN = plan.PIN.ValueString()
	}
	if !plan.GuestPIN.IsNull() {
		updateRequest.GuestPIN = plan.GuestPIN.ValueString()
	}
	if !plan.Tag.IsNull() {
		updateRequest.Tag = plan.Tag.ValueString()
	}

	// Set optional boolean fields (use pointers for update requests)
	if !plan.AllowGuests.IsNull() {
		allowGuests := plan.AllowGuests.ValueBool()
		updateRequest.AllowGuests = &allowGuests
	}
	if !plan.GuestsMuted.IsNull() {
		guestsMuted := plan.GuestsMuted.ValueBool()
		updateRequest.GuestsMuted = &guestsMuted
	}
	if !plan.HostsCanUnmute.IsNull() {
		hostsCanUnmute := plan.HostsCanUnmute.ValueBool()
		updateRequest.HostsCanUnmute = &hostsCanUnmute
	}

	// Set optional integer fields
	if !plan.MaxPixelsPerSecond.IsNull() {
		maxPixelsPerSecond := int(plan.MaxPixelsPerSecond.ValueInt32())
		updateRequest.MaxPixelsPerSecond = &maxPixelsPerSecond
	}

	_, err := r.InfinityClient.Config().UpdateConference(ctx, resourceID, updateRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Infinity conference",
			fmt.Sprintf("Could not update Infinity conference with ID %d: %s", resourceID, err),
		)
		return
	}

	// Re-read the resource to get the latest state
	updatedModel, err := r.read(ctx, resourceID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Updated Infinity conference",
			fmt.Sprintf("Could not read updated Infinity conference with ID %d: %s", resourceID, err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, updatedModel)...)
}

func (r *InfinityConferenceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	state := &InfinityConferenceResourceModel{}

	tflog.Info(ctx, "Deleting Infinity conference")

	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.InfinityClient.Config().DeleteConference(ctx, int(state.ResourceID.ValueInt32()))

	// Ignore 404 Not Found and Lookup errors on delete
	if err != nil && !isNotFoundError(err) && !isLookupError(err) {
		resp.Diagnostics.AddError(
			"Error Deleting Infinity conference",
			fmt.Sprintf("Could not delete Infinity conference with ID %s: %s", state.ID.ValueString(), err),
		)
		return
	}
}

func (r *InfinityConferenceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resourceID, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Import ID must be a valid integer for the resource ID. Got: %s", req.ID),
		)
		return
	}

	tflog.Trace(ctx, fmt.Sprintf("Importing Infinity conference with resource ID: %d", resourceID))

	// Read the resource from the API
	model, err := r.read(ctx, resourceID)
	if err != nil {
		// Check if the error is a 404 (not found)
		if isNotFoundError(err) {
			resp.Diagnostics.AddError(
				"Infinity Conference Not Found",
				fmt.Sprintf("Infinity conference with resource ID %d not found.", resourceID),
			)
			return
		}
		resp.Diagnostics.AddError(
			"Error Importing Infinity Conference",
			fmt.Sprintf("Could not import Infinity conference with resource ID %d: %s", resourceID, err),
		)
		return
	}

	// Set the state from the imported resource
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}
