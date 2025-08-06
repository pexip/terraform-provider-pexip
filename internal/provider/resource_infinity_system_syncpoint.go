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

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/pexip/go-infinity-sdk/v38/config"
)

var (
	_ resource.ResourceWithImportState = (*InfinitySystemSyncpointResource)(nil)
)

type InfinitySystemSyncpointResource struct {
	InfinityClient InfinityClient
}

type InfinitySystemSyncpointResourceModel struct {
	ID           types.String `tfsdk:"id"`
	ResourceID   types.Int32  `tfsdk:"resource_id"`
	CreationTime types.String `tfsdk:"creation_time"`
}

func (r *InfinitySystemSyncpointResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_infinity_system_syncpoint"
}

func (r *InfinitySystemSyncpointResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *InfinitySystemSyncpointResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Resource URI for the system syncpoint in Infinity",
			},
			"resource_id": schema.Int32Attribute{
				Computed:            true,
				MarkdownDescription: "The resource integer identifier for the system syncpoint in Infinity",
			},
			"creation_time": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The timestamp when this system syncpoint was created",
			},
		},
		MarkdownDescription: "Manages a system syncpoint with the Infinity service. System syncpoints are critical for multi-site deployments and provide system synchronization points for coordinated operations. Note: This resource only supports creation and reading - syncpoints cannot be updated or deleted once created.",
	}
}

func (r *InfinitySystemSyncpointResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	plan := &InfinitySystemSyncpointResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createRequest := &config.SystemSyncpointCreateRequest{}

	createResponse, err := r.InfinityClient.Config().CreateSystemSyncpoint(ctx, createRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Infinity system syncpoint",
			fmt.Sprintf("Could not create Infinity system syncpoint: %s", err),
		)
		return
	}

	resourceID, err := createResponse.ResourceID()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Retrieving Infinity system syncpoint ID",
			fmt.Sprintf("Could not retrieve ID for created Infinity system syncpoint: %s", err),
		)
		return
	}

	// Read the state from the API to get all computed values
	model, err := r.read(ctx, resourceID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Created Infinity system syncpoint",
			fmt.Sprintf("Could not read created Infinity system syncpoint with ID %d: %s", resourceID, err),
		)
		return
	}
	tflog.Trace(ctx, fmt.Sprintf("created Infinity system syncpoint with ID: %s", model.ID))

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *InfinitySystemSyncpointResource) read(ctx context.Context, resourceID int) (*InfinitySystemSyncpointResourceModel, error) {
	var data InfinitySystemSyncpointResourceModel

	srv, err := r.InfinityClient.Config().GetSystemSyncpoint(ctx, resourceID)
	if err != nil {
		return nil, err
	}

	if srv.ResourceURI == "" {
		return nil, fmt.Errorf("system syncpoint with ID %d not found", resourceID)
	}

	data.ID = types.StringValue(srv.ResourceURI)
	data.ResourceID = types.Int32Value(int32(resourceID)) // #nosec G115 -- API values are expected to be within int32 range
	data.CreationTime = types.StringValue(srv.CreationTime.String())

	return &data, nil
}

func (r *InfinitySystemSyncpointResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	state := &InfinitySystemSyncpointResourceModel{}

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
			"Error Reading Infinity system syncpoint",
			fmt.Sprintf("Could not read Infinity system syncpoint: %s", err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *InfinitySystemSyncpointResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError(
		"Update Not Supported",
		"System syncpoint resources cannot be updated. System syncpoints are immutable once created.",
	)
}

func (r *InfinitySystemSyncpointResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	resp.Diagnostics.AddError(
		"Delete Not Supported",
		"System syncpoint resources cannot be deleted. System syncpoints are permanent once created for system synchronization purposes.",
	)
}

func (r *InfinitySystemSyncpointResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resourceID, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Import ID must be a valid integer for the resource ID. Got: %s", req.ID),
		)
		return
	}

	tflog.Trace(ctx, fmt.Sprintf("Importing Infinity system syncpoint with resource ID: %d", resourceID))

	// Read the resource from the API
	model, err := r.read(ctx, resourceID)
	if err != nil {
		// Check if the error is a 404 (not found)
		if isNotFoundError(err) {
			resp.Diagnostics.AddError(
				"Infinity System Syncpoint Not Found",
				fmt.Sprintf("Infinity system syncpoint with resource ID %d not found.", resourceID),
			)
			return
		}
		resp.Diagnostics.AddError(
			"Error Importing Infinity System Syncpoint",
			fmt.Sprintf("Could not import Infinity system syncpoint with resource ID %d: %s", resourceID, err),
		)
		return
	}

	// Set the state from the imported resource
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}
