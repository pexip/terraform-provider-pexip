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
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/pexip/go-infinity-sdk/v38/config"
)

var (
	_ resource.ResourceWithImportState = (*InfinityMediaLibraryPlaylistResource)(nil)
)

type InfinityMediaLibraryPlaylistResource struct {
	InfinityClient InfinityClient
}

type InfinityMediaLibraryPlaylistResourceModel struct {
	ID              types.String `tfsdk:"id"`
	ResourceID      types.Int32  `tfsdk:"resource_id"`
	Name            types.String `tfsdk:"name"`
	Description     types.String `tfsdk:"description"`
	Loop            types.Bool   `tfsdk:"loop"`
	Shuffle         types.Bool   `tfsdk:"shuffle"`
	PlaylistEntries types.Set    `tfsdk:"playlist_entries"`
}

func (r *InfinityMediaLibraryPlaylistResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_infinity_media_library_playlist"
}

func (r *InfinityMediaLibraryPlaylistResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *InfinityMediaLibraryPlaylistResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Resource URI for the media library playlist in Infinity",
			},
			"resource_id": schema.Int32Attribute{
				Computed:            true,
				MarkdownDescription: "The resource integer identifier for the media library playlist in Infinity",
			},
			"name": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					stringvalidator.LengthAtMost(250),
				},
				MarkdownDescription: "The name of the media library playlist. Maximum length: 250 characters.",
			},
			"description": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(500),
				},
				MarkdownDescription: "Description of the media library playlist. Maximum length: 500 characters.",
			},
			"loop": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Whether the playlist should loop when it reaches the end.",
			},
			"shuffle": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Whether the playlist entries should be played in random order.",
			},
			"playlist_entries": schema.SetAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "List of media library entry URIs that make up this playlist.",
			},
		},
		MarkdownDescription: "Manages a media library playlist configuration with the Infinity service. Media library playlists organize media entries for sequential or randomized playback in conferences.",
	}
}

func (r *InfinityMediaLibraryPlaylistResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	plan := &InfinityMediaLibraryPlaylistResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createRequest := &config.MediaLibraryPlaylistCreateRequest{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
		Loop:        plan.Loop.ValueBool(),
		Shuffle:     plan.Shuffle.ValueBool(),
	}

	// Handle playlist entries
	if !plan.PlaylistEntries.IsNull() && !plan.PlaylistEntries.IsUnknown() {
		var entries []string
		resp.Diagnostics.Append(plan.PlaylistEntries.ElementsAs(ctx, &entries, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		createRequest.PlaylistEntries = entries
	}

	createResponse, err := r.InfinityClient.Config().CreateMediaLibraryPlaylist(ctx, createRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Infinity media library playlist",
			fmt.Sprintf("Could not create Infinity media library playlist: %s", err),
		)
		return
	}

	resourceID, err := createResponse.ResourceID()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Retrieving Infinity media library playlist ID",
			fmt.Sprintf("Could not retrieve ID for created Infinity media library playlist: %s", err),
		)
		return
	}

	// Read the state from the API to get all computed values
	model, err := r.read(ctx, resourceID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Created Infinity media library playlist",
			fmt.Sprintf("Could not read created Infinity media library playlist with ID %d: %s", resourceID, err),
		)
		return
	}
	tflog.Trace(ctx, fmt.Sprintf("created Infinity media library playlist with ID: %s, name: %s", model.ID, model.Name))

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *InfinityMediaLibraryPlaylistResource) read(ctx context.Context, resourceID int) (*InfinityMediaLibraryPlaylistResourceModel, error) {
	var data InfinityMediaLibraryPlaylistResourceModel

	srv, err := r.InfinityClient.Config().GetMediaLibraryPlaylist(ctx, resourceID)
	if err != nil {
		return nil, err
	}

	if srv.ResourceURI == "" {
		return nil, fmt.Errorf("media library playlist with ID %d not found", resourceID)
	}

	data.ID = types.StringValue(srv.ResourceURI)
	data.ResourceID = types.Int32Value(int32(resourceID)) // #nosec G115 -- API values are expected to be within int32 range
	data.Name = types.StringValue(srv.Name)
	data.Description = types.StringValue(srv.Description)
	data.Loop = types.BoolValue(srv.Loop)
	data.Shuffle = types.BoolValue(srv.Shuffle)

	// Handle playlist entries
	if srv.PlaylistEntries != nil {
		entriesSet, diags := types.SetValueFrom(ctx, types.StringType, srv.PlaylistEntries)
		if diags.HasError() {
			return nil, fmt.Errorf("failed to convert playlist entries: %s", diags.Errors())
		}
		data.PlaylistEntries = entriesSet
	} else {
		data.PlaylistEntries = types.SetNull(types.StringType)
	}

	return &data, nil
}

func (r *InfinityMediaLibraryPlaylistResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	state := &InfinityMediaLibraryPlaylistResourceModel{}

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
			"Error Reading Infinity media library playlist",
			fmt.Sprintf("Could not read Infinity media library playlist: %s", err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *InfinityMediaLibraryPlaylistResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	plan := &InfinityMediaLibraryPlaylistResourceModel{}
	state := &InfinityMediaLibraryPlaylistResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resourceID := int(state.ResourceID.ValueInt32())

	updateRequest := &config.MediaLibraryPlaylistUpdateRequest{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
	}

	// Handle boolean pointers
	if !plan.Loop.IsNull() {
		loop := plan.Loop.ValueBool()
		updateRequest.Loop = &loop
	}

	if !plan.Shuffle.IsNull() {
		shuffle := plan.Shuffle.ValueBool()
		updateRequest.Shuffle = &shuffle
	}

	// Handle playlist entries
	if !plan.PlaylistEntries.IsNull() && !plan.PlaylistEntries.IsUnknown() {
		var entries []string
		resp.Diagnostics.Append(plan.PlaylistEntries.ElementsAs(ctx, &entries, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		updateRequest.PlaylistEntries = entries
	}

	_, err := r.InfinityClient.Config().UpdateMediaLibraryPlaylist(ctx, resourceID, updateRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Infinity media library playlist",
			fmt.Sprintf("Could not update Infinity media library playlist with ID %d: %s", resourceID, err),
		)
		return
	}

	// Re-read the resource to get the latest state
	updatedModel, err := r.read(ctx, resourceID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Updated Infinity media library playlist",
			fmt.Sprintf("Could not read updated Infinity media library playlist with ID %d: %s", resourceID, err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, updatedModel)...)
}

func (r *InfinityMediaLibraryPlaylistResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	state := &InfinityMediaLibraryPlaylistResourceModel{}

	tflog.Info(ctx, "Deleting Infinity media library playlist")

	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.InfinityClient.Config().DeleteMediaLibraryPlaylist(ctx, int(state.ResourceID.ValueInt32()))

	// Ignore 404 Not Found and Lookup errors on delete
	if err != nil && !isNotFoundError(err) && !isLookupError(err) {
		resp.Diagnostics.AddError(
			"Error Deleting Infinity media library playlist",
			fmt.Sprintf("Could not delete Infinity media library playlist with ID %s: %s", state.ID.ValueString(), err),
		)
		return
	}
}

func (r *InfinityMediaLibraryPlaylistResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resourceID, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Import ID must be a valid integer for the resource ID. Got: %s", req.ID),
		)
		return
	}

	tflog.Trace(ctx, fmt.Sprintf("Importing Infinity media library playlist with resource ID: %d", resourceID))

	// Read the resource from the API
	model, err := r.read(ctx, resourceID)
	if err != nil {
		// Check if the error is a 404 (not found)
		if isNotFoundError(err) {
			resp.Diagnostics.AddError(
				"Infinity Media Library Playlist Not Found",
				fmt.Sprintf("Infinity media library playlist with resource ID %d not found.", resourceID),
			)
			return
		}
		resp.Diagnostics.AddError(
			"Error Importing Infinity Media Library Playlist",
			fmt.Sprintf("Could not import Infinity media library playlist with resource ID %d: %s", resourceID, err),
		)
		return
	}

	// Set the state from the imported resource
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}
