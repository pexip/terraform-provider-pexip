/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package provider

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/pexip/go-infinity-sdk/v38/config"
)

var (
	_ resource.ResourceWithImportState = (*InfinityMediaLibraryEntryResource)(nil)
)

type InfinityMediaLibraryEntryResource struct {
	InfinityClient InfinityClient
}

type InfinityMediaLibraryEntryResourceModel struct {
	ID          types.String `tfsdk:"id"`
	ResourceID  types.Int32  `tfsdk:"resource_id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	UUID        types.String `tfsdk:"uuid"`
	FileName    types.String `tfsdk:"file_name"`
	MediaType   types.String `tfsdk:"media_type"`
	MediaFormat types.String `tfsdk:"media_format"`
	MediaSize   types.Int64  `tfsdk:"media_size"`
	MediaFile   types.String `tfsdk:"media_file"`
}

func (r *InfinityMediaLibraryEntryResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_infinity_media_library_entry"
}

func (r *InfinityMediaLibraryEntryResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *InfinityMediaLibraryEntryResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Resource URI for the media library entry in Infinity",
			},
			"resource_id": schema.Int32Attribute{
				Computed:            true,
				MarkdownDescription: "The resource integer identifier for the media library entry in Infinity",
			},
			"name": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(1024),
				},
				MarkdownDescription: "The name used to refer to this Media Library Item. Maximum length: 1024 characters.",
			},
			"description": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(""),
				Validators: []validator.String{
					stringvalidator.LengthAtMost(1024),
				},
				MarkdownDescription: "A description of this Media Library Item. Maximum length: 1024 characters.",
			},
			"uuid": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "UUID for the media library entry.",
			},
			"file_name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The file name of the uploaded media file.",
			},
			"media_type": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The media type of the uploaded file (e.g., video, audio, image).",
			},
			"media_format": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The media format of the uploaded file (e.g., mp4, wav, jpg).",
			},
			"media_size": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "The size of the media file in bytes.",
			},
			"media_file": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
				MarkdownDescription: "Path to the media file to upload (e.g., `media_file = \"path/to/video.mp4\"`).",
			},
		},
		MarkdownDescription: "Manages a media library entry configuration with the Infinity service. Media library entries are used for storing media files such as images, videos, and audio files that can be used in conferences.",
	}
}

func (r *InfinityMediaLibraryEntryResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	plan := &InfinityMediaLibraryEntryResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createRequest := &config.MediaLibraryEntryCreateRequest{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
		UUID:        plan.UUID.ValueString(),
	}

	// Open the media file
	mediaFilePath := plan.MediaFile.ValueString()
	mediaFile, err := os.Open(mediaFilePath) // #nosec G304 -- File path provided by user in Terraform configuration
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Opening Media File",
			fmt.Sprintf("Could not open media file at path '%s': %s", mediaFilePath, err),
		)
		return
	}
	defer func() { _ = mediaFile.Close() }()

	// Extract filename from the path
	filename := filepath.Base(mediaFilePath)

	createResponse, err := r.InfinityClient.Config().CreateMediaLibraryEntry(ctx, createRequest, filename, mediaFile)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Infinity media library entry",
			fmt.Sprintf("Could not create Infinity media library entry: %s", err),
		)
		return
	}

	resourceID, err := createResponse.ResourceID()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Retrieving Infinity media library entry ID",
			fmt.Sprintf("Could not retrieve ID for created Infinity media library entry: %s", err),
		)
		return
	}

	// Read the state from the API to get all computed values
	model, err := r.read(ctx, resourceID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Created Infinity media library entry",
			fmt.Sprintf("Could not read created Infinity media library entry with ID %d: %s", resourceID, err),
		)
		return
	}
	// Preserve the MediaFile value from the plan (cannot be retrieved from API)
	model.MediaFile = plan.MediaFile
	tflog.Trace(ctx, fmt.Sprintf("created Infinity media library entry with ID: %s, name: %s", model.ID, model.Name))

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *InfinityMediaLibraryEntryResource) read(ctx context.Context, resourceID int) (*InfinityMediaLibraryEntryResourceModel, error) {
	var data InfinityMediaLibraryEntryResourceModel

	srv, err := r.InfinityClient.Config().GetMediaLibraryEntry(ctx, resourceID)
	if err != nil {
		return nil, err
	}

	if srv.ResourceURI == "" {
		return nil, fmt.Errorf("media library entry with ID %d not found", resourceID)
	}

	data.ID = types.StringValue(srv.ResourceURI)
	data.ResourceID = types.Int32Value(int32(resourceID)) // #nosec G115 -- API values are expected to be within int32 range
	data.Name = types.StringValue(srv.Name)
	data.Description = types.StringValue(srv.Description)
	data.UUID = types.StringValue(srv.UUID)
	data.FileName = types.StringValue(srv.FileName)
	data.MediaType = types.StringValue(srv.MediaType)
	data.MediaFormat = types.StringValue(srv.MediaFormat)
	data.MediaSize = types.Int64Value(int64(srv.MediaSize))
	// Note: MediaFile is not returned by the API (binary content cannot be retrieved)
	// It will be preserved from plan/state by the caller

	return &data, nil
}

func (r *InfinityMediaLibraryEntryResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	state := &InfinityMediaLibraryEntryResourceModel{}

	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Preserve the current MediaFile value before refreshing state
	currentMediaFile := state.MediaFile

	resourceID := int(state.ResourceID.ValueInt32())
	state, err := r.read(ctx, resourceID)
	if err != nil {
		// Check if the error is a 404 (not found)
		if isNotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading Infinity media library entry",
			fmt.Sprintf("Could not read Infinity media library entry: %s", err),
		)
		return
	}

	// Restore the MediaFile value (cannot be retrieved from API)
	state.MediaFile = currentMediaFile

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *InfinityMediaLibraryEntryResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	plan := &InfinityMediaLibraryEntryResourceModel{}
	state := &InfinityMediaLibraryEntryResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resourceID := int(state.ResourceID.ValueInt32())

	updateRequest := &config.MediaLibraryEntryUpdateRequest{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
		UUID:        plan.UUID.ValueString(),
	}

	// Open the media file
	mediaFilePath := plan.MediaFile.ValueString()
	mediaFile, err := os.Open(mediaFilePath) // #nosec G304 -- File path provided by user in Terraform configuration
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Opening Media File",
			fmt.Sprintf("Could not open media file at path '%s': %s", mediaFilePath, err),
		)
		return
	}
	defer func() { _ = mediaFile.Close() }()

	// Extract filename from the path
	filename := filepath.Base(mediaFilePath)

	_, err = r.InfinityClient.Config().UpdateMediaLibraryEntry(ctx, resourceID, updateRequest, filename, mediaFile)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Infinity media library entry",
			fmt.Sprintf("Could not update Infinity media library entry with ID %d: %s", resourceID, err),
		)
		return
	}

	// Re-read the resource to get the latest state
	updatedModel, err := r.read(ctx, resourceID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Updated Infinity media library entry",
			fmt.Sprintf("Could not read updated Infinity media library entry with ID %d: %s", resourceID, err),
		)
		return
	}

	// Preserve the MediaFile value from the plan (cannot be retrieved from API)
	updatedModel.MediaFile = plan.MediaFile

	resp.Diagnostics.Append(resp.State.Set(ctx, updatedModel)...)
}

func (r *InfinityMediaLibraryEntryResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	state := &InfinityMediaLibraryEntryResourceModel{}

	tflog.Info(ctx, "Deleting Infinity media library entry")

	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.InfinityClient.Config().DeleteMediaLibraryEntry(ctx, int(state.ResourceID.ValueInt32()))

	// Ignore 404 Not Found and Lookup errors on delete
	if err != nil && !isNotFoundError(err) && !isLookupError(err) {
		resp.Diagnostics.AddError(
			"Error Deleting Infinity media library entry",
			fmt.Sprintf("Could not delete Infinity media library entry with ID %s: %s", state.ID.ValueString(), err),
		)
		return
	}
}

func (r *InfinityMediaLibraryEntryResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resourceID, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Import ID must be a valid integer for the resource ID. Got: %s", req.ID),
		)
		return
	}

	tflog.Trace(ctx, fmt.Sprintf("Importing Infinity media library entry with resource ID: %d", resourceID))

	// Read the resource from the API
	model, err := r.read(ctx, resourceID)
	if err != nil {
		// Check if the error is a 404 (not found)
		if isNotFoundError(err) {
			resp.Diagnostics.AddError(
				"Infinity Media Library Entry Not Found",
				fmt.Sprintf("Infinity media library entry with resource ID %d not found.", resourceID),
			)
			return
		}
		resp.Diagnostics.AddError(
			"Error Importing Infinity Media Library Entry",
			fmt.Sprintf("Could not import Infinity media library entry with resource ID %d: %s", resourceID, err),
		)
		return
	}

	// MediaFile content cannot be retrieved from the API
	// Set to empty string - user must provide the correct media file content in their configuration
	// and run `terraform apply` after import to upload the file if needed
	model.MediaFile = types.StringValue("")

	// Set the state from the imported resource
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}
