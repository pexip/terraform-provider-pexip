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

	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/pexip/go-infinity-sdk/v38/config"
)

var (
	_ resource.ResourceWithImportState = (*InfinityMediaLibraryPlaylistEntryResource)(nil)
)

type InfinityMediaLibraryPlaylistEntryResource struct {
	InfinityClient InfinityClient
}

type InfinityMediaLibraryPlaylistEntryResourceModel struct {
	ID         types.String `tfsdk:"id"`
	ResourceID types.Int32  `tfsdk:"resource_id"`
	EntryType  types.String `tfsdk:"entry_type"`
	Media      types.String `tfsdk:"media"`
	Playlist   types.String `tfsdk:"playlist"`
	Position   types.Int32  `tfsdk:"position"`
	Playcount  types.Int32  `tfsdk:"playcount"`
}

func (r *InfinityMediaLibraryPlaylistEntryResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_infinity_media_library_playlist_entry"
}

func (r *InfinityMediaLibraryPlaylistEntryResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *InfinityMediaLibraryPlaylistEntryResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Resource URI for the media library playlist entry in Infinity",
			},
			"resource_id": schema.Int32Attribute{
				Computed:            true,
				MarkdownDescription: "The resource integer identifier for the media library playlist entry in Infinity",
			},
			"entry_type": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("MEDIA"),
				Validators: []validator.String{
					stringvalidator.OneOf("MEDIA"),
				},
				MarkdownDescription: "Type of entry referred to by playlist entry. Valid choices: MEDIA. Default: MEDIA.",
			},
			"media": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Resource URI of the media library entry to include in the playlist.",
			},
			"playlist": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Resource URI of the playlist this entry belongs to.",
			},
			"position": schema.Int32Attribute{
				Optional: true,
				Computed: true,
				Default:  int32default.StaticInt32(1),
				Validators: []validator.Int32{
					int32validator.AtLeast(1),
				},
				MarkdownDescription: "Every item must have a unique position specified. This refers to the order in which the Media Library Item is played relative to the other items in the playlist. Please note that Position is only used when shuffle is off. Default: 1.",
			},
			"playcount": schema.Int32Attribute{
				Optional: true,
				Computed: true,
				Default:  int32default.StaticInt32(1),
				Validators: []validator.Int32{
					int32validator.AtLeast(0),
				},
				MarkdownDescription: "The number of times the item is played, the default is 1. If you set the value to 0 the item plays repeatedly until the user disconnects themselves or until the call is terminated programmatically via the management API. Default: 1.",
			},
		},
		MarkdownDescription: "Manages a media library playlist entry configuration with the Infinity service. Playlist entries define the media items and their playback order within a playlist.",
	}
}

func (r *InfinityMediaLibraryPlaylistEntryResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	plan := &InfinityMediaLibraryPlaylistEntryResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createRequest := &config.MediaLibraryPlaylistEntryCreateRequest{
		EntryType: plan.EntryType.ValueString(),
		Position:  int(plan.Position.ValueInt32()),
		Playcount: int(plan.Playcount.ValueInt32()),
	}

	// Handle optional media field
	if !plan.Media.IsNull() && !plan.Media.IsUnknown() {
		media := plan.Media.ValueString()
		createRequest.Media = &media
	}

	// Handle required playlist field
	if !plan.Playlist.IsNull() && !plan.Playlist.IsUnknown() {
		playlist := plan.Playlist.ValueString()
		createRequest.Playlist = &playlist
	}

	createResponse, err := r.InfinityClient.Config().CreateMediaLibraryPlaylistEntry(ctx, createRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Infinity media library playlist entry",
			fmt.Sprintf("Could not create Infinity media library playlist entry: %s", err),
		)
		return
	}

	resourceID, err := createResponse.ResourceID()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Retrieving Infinity media library playlist entry ID",
			fmt.Sprintf("Could not retrieve ID for created Infinity media library playlist entry: %s", err),
		)
		return
	}

	// Read the state from the API to get all computed values
	model, err := r.read(ctx, resourceID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Created Infinity media library playlist entry",
			fmt.Sprintf("Could not read created Infinity media library playlist entry with ID %d: %s", resourceID, err),
		)
		return
	}
	tflog.Trace(ctx, fmt.Sprintf("created Infinity media library playlist entry with ID: %s", model.ID))

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *InfinityMediaLibraryPlaylistEntryResource) read(ctx context.Context, resourceID int) (*InfinityMediaLibraryPlaylistEntryResourceModel, error) {
	var data InfinityMediaLibraryPlaylistEntryResourceModel

	srv, err := r.InfinityClient.Config().GetMediaLibraryPlaylistEntry(ctx, resourceID)
	if err != nil {
		return nil, err
	}

	if srv.ResourceURI == "" {
		return nil, fmt.Errorf("media library playlist entry with ID %d not found", resourceID)
	}

	data.ID = types.StringValue(srv.ResourceURI)
	data.ResourceID = types.Int32Value(int32(resourceID)) // #nosec G115 -- API values are expected to be within int32 range
	data.EntryType = types.StringValue(srv.EntryType)
	data.Position = types.Int32Value(int32(srv.Position)) // #nosec G115 -- API values are expected to be within int32 range
	data.Playcount = types.Int32Value(int32(srv.Playcount)) // #nosec G115 -- API values are expected to be within int32 range

	// Handle optional media field
	if srv.Media != nil {
		data.Media = types.StringValue(*srv.Media)
	} else {
		data.Media = types.StringNull()
	}

	// Handle optional playlist field
	if srv.Playlist != nil {
		data.Playlist = types.StringValue(*srv.Playlist)
	} else {
		data.Playlist = types.StringNull()
	}

	return &data, nil
}

func (r *InfinityMediaLibraryPlaylistEntryResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	state := &InfinityMediaLibraryPlaylistEntryResourceModel{}

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
			"Error Reading Infinity media library playlist entry",
			fmt.Sprintf("Could not read Infinity media library playlist entry: %s", err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *InfinityMediaLibraryPlaylistEntryResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	plan := &InfinityMediaLibraryPlaylistEntryResourceModel{}
	state := &InfinityMediaLibraryPlaylistEntryResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resourceID := int(state.ResourceID.ValueInt32())

	updateRequest := &config.MediaLibraryPlaylistEntryUpdateRequest{
		EntryType: plan.EntryType.ValueString(),
	}

	// Handle optional media field
	if !plan.Media.IsNull() && !plan.Media.IsUnknown() {
		media := plan.Media.ValueString()
		updateRequest.Media = &media
	}

	// Handle optional playlist field
	if !plan.Playlist.IsNull() && !plan.Playlist.IsUnknown() {
		playlist := plan.Playlist.ValueString()
		updateRequest.Playlist = &playlist
	}

	// Handle position
	if !plan.Position.IsNull() {
		position := int(plan.Position.ValueInt32())
		updateRequest.Position = &position
	}

	// Handle playcount
	if !plan.Playcount.IsNull() {
		playcount := int(plan.Playcount.ValueInt32())
		updateRequest.Playcount = &playcount
	}

	_, err := r.InfinityClient.Config().UpdateMediaLibraryPlaylistEntry(ctx, resourceID, updateRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Infinity media library playlist entry",
			fmt.Sprintf("Could not update Infinity media library playlist entry with ID %d: %s", resourceID, err),
		)
		return
	}

	// Re-read the resource to get the latest state
	updatedModel, err := r.read(ctx, resourceID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Updated Infinity media library playlist entry",
			fmt.Sprintf("Could not read updated Infinity media library playlist entry with ID %d: %s", resourceID, err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, updatedModel)...)
}

func (r *InfinityMediaLibraryPlaylistEntryResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	state := &InfinityMediaLibraryPlaylistEntryResourceModel{}

	tflog.Info(ctx, "Deleting Infinity media library playlist entry")

	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.InfinityClient.Config().DeleteMediaLibraryPlaylistEntry(ctx, int(state.ResourceID.ValueInt32()))

	// Ignore 404 Not Found and Lookup errors on delete
	if err != nil && !isNotFoundError(err) && !isLookupError(err) {
		resp.Diagnostics.AddError(
			"Error Deleting Infinity media library playlist entry",
			fmt.Sprintf("Could not delete Infinity media library playlist entry with ID %s: %s", state.ID.ValueString(), err),
		)
		return
	}
}

func (r *InfinityMediaLibraryPlaylistEntryResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resourceID, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Import ID must be a valid integer for the resource ID. Got: %s", req.ID),
		)
		return
	}

	tflog.Trace(ctx, fmt.Sprintf("Importing Infinity media library playlist entry with resource ID: %d", resourceID))

	// Read the resource from the API
	model, err := r.read(ctx, resourceID)
	if err != nil {
		// Check if the error is a 404 (not found)
		if isNotFoundError(err) {
			resp.Diagnostics.AddError(
				"Infinity Media Library Playlist Entry Not Found",
				fmt.Sprintf("Infinity media library playlist entry with resource ID %d not found.", resourceID),
			)
			return
		}
		resp.Diagnostics.AddError(
			"Error Importing Infinity Media Library Playlist Entry",
			fmt.Sprintf("Could not import Infinity media library playlist entry with resource ID %d: %s", resourceID, err),
		)
		return
	}

	// Set the state from the imported resource
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}
