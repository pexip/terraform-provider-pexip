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
	_ resource.ResourceWithImportState = (*InfinityUserGroupEntityMappingResource)(nil)
)

type InfinityUserGroupEntityMappingResource struct {
	InfinityClient InfinityClient
}

type InfinityUserGroupEntityMappingResourceModel struct {
	ID                types.String `tfsdk:"id"`
	ResourceID        types.Int32  `tfsdk:"resource_id"`
	Description       types.String `tfsdk:"description"`
	EntityResourceURI types.String `tfsdk:"entity_resource_uri"`
	UserGroup         types.String `tfsdk:"user_group"`
}

func (r *InfinityUserGroupEntityMappingResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_infinity_user_group_entity_mapping"
}

func (r *InfinityUserGroupEntityMappingResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *InfinityUserGroupEntityMappingResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Resource URI for the user group entity mapping in Infinity",
			},
			"resource_id": schema.Int32Attribute{
				Computed:            true,
				MarkdownDescription: "The resource integer identifier for the user group entity mapping in Infinity",
			},
			"description": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(500),
				},
				MarkdownDescription: "Description of the user group entity mapping. Maximum length: 500 characters.",
			},
			"entity_resource_uri": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
				MarkdownDescription: "The entity resource URI that this mapping applies to. This should be a valid resource URI in the system.",
			},
			"user_group": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
				MarkdownDescription: "The user group that should be mapped to the entity. This should be a valid user group URI or identifier.",
			},
		},
		MarkdownDescription: "Manages a user group entity mapping configuration with the Infinity service. User group entity mappings define relationships between user groups and system entities for access control and permissions.",
	}
}

func (r *InfinityUserGroupEntityMappingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	plan := &InfinityUserGroupEntityMappingResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createRequest := &config.UserGroupEntityMappingCreateRequest{
		Description:       plan.Description.ValueString(),
		EntityResourceURI: plan.EntityResourceURI.ValueString(),
		UserGroup:         plan.UserGroup.ValueString(),
	}

	createResponse, err := r.InfinityClient.Config().CreateUserGroupEntityMapping(ctx, createRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Infinity user group entity mapping",
			fmt.Sprintf("Could not create Infinity user group entity mapping: %s", err),
		)
		return
	}

	resourceID, err := createResponse.ResourceID()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Retrieving Infinity user group entity mapping ID",
			fmt.Sprintf("Could not retrieve ID for created Infinity user group entity mapping: %s", err),
		)
		return
	}

	// Read the state from the API to get all computed values
	model, err := r.read(ctx, resourceID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Created Infinity user group entity mapping",
			fmt.Sprintf("Could not read created Infinity user group entity mapping with ID %d: %s", resourceID, err),
		)
		return
	}
	tflog.Trace(ctx, fmt.Sprintf("created Infinity user group entity mapping with ID: %s, user_group: %s", model.ID, model.UserGroup))

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *InfinityUserGroupEntityMappingResource) read(ctx context.Context, resourceID int) (*InfinityUserGroupEntityMappingResourceModel, error) {
	var data InfinityUserGroupEntityMappingResourceModel

	srv, err := r.InfinityClient.Config().GetUserGroupEntityMapping(ctx, resourceID)
	if err != nil {
		return nil, err
	}

	if srv.ResourceURI == "" {
		return nil, fmt.Errorf("user group entity mapping with ID %d not found", resourceID)
	}

	data.ID = types.StringValue(srv.ResourceURI)
	data.ResourceID = types.Int32Value(int32(resourceID)) // #nosec G115 -- API values are expected to be within int32 range
	data.Description = types.StringValue(srv.Description)
	data.EntityResourceURI = types.StringValue(srv.EntityResourceURI)
	data.UserGroup = types.StringValue(srv.UserGroup)

	return &data, nil
}

func (r *InfinityUserGroupEntityMappingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	state := &InfinityUserGroupEntityMappingResourceModel{}

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
			"Error Reading Infinity user group entity mapping",
			fmt.Sprintf("Could not read Infinity user group entity mapping: %s", err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *InfinityUserGroupEntityMappingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	plan := &InfinityUserGroupEntityMappingResourceModel{}
	state := &InfinityUserGroupEntityMappingResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resourceID := int(state.ResourceID.ValueInt32())

	updateRequest := &config.UserGroupEntityMappingUpdateRequest{
		Description:       plan.Description.ValueString(),
		EntityResourceURI: plan.EntityResourceURI.ValueString(),
		UserGroup:         plan.UserGroup.ValueString(),
	}

	_, err := r.InfinityClient.Config().UpdateUserGroupEntityMapping(ctx, resourceID, updateRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Infinity user group entity mapping",
			fmt.Sprintf("Could not update Infinity user group entity mapping with ID %d: %s", resourceID, err),
		)
		return
	}

	// Re-read the resource to get the latest state
	updatedModel, err := r.read(ctx, resourceID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Updated Infinity user group entity mapping",
			fmt.Sprintf("Could not read updated Infinity user group entity mapping with ID %d: %s", resourceID, err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, updatedModel)...)
}

func (r *InfinityUserGroupEntityMappingResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	state := &InfinityUserGroupEntityMappingResourceModel{}

	tflog.Info(ctx, "Deleting Infinity user group entity mapping")

	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.InfinityClient.Config().DeleteUserGroupEntityMapping(ctx, int(state.ResourceID.ValueInt32()))

	// Ignore 404 Not Found and Lookup errors on delete
	if err != nil && !isNotFoundError(err) && !isLookupError(err) {
		resp.Diagnostics.AddError(
			"Error Deleting Infinity user group entity mapping",
			fmt.Sprintf("Could not delete Infinity user group entity mapping with ID %s: %s", state.ID.ValueString(), err),
		)
		return
	}
}

func (r *InfinityUserGroupEntityMappingResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resourceID, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Import ID must be a valid integer for the resource ID. Got: %s", req.ID),
		)
		return
	}

	tflog.Trace(ctx, fmt.Sprintf("Importing Infinity user group entity mapping with resource ID: %d", resourceID))

	// Read the resource from the API
	model, err := r.read(ctx, resourceID)
	if err != nil {
		// Check if the error is a 404 (not found)
		if isNotFoundError(err) {
			resp.Diagnostics.AddError(
				"Infinity User Group Entity Mapping Not Found",
				fmt.Sprintf("Infinity user group entity mapping with resource ID %d not found.", resourceID),
			)
			return
		}
		resp.Diagnostics.AddError(
			"Error Importing Infinity User Group Entity Mapping",
			fmt.Sprintf("Could not import Infinity user group entity mapping with resource ID %d: %s", resourceID, err),
		)
		return
	}

	// Set the state from the imported resource
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}
