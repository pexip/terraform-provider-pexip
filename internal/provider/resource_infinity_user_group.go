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
	_ resource.ResourceWithImportState = (*InfinityUserGroupResource)(nil)
)

type InfinityUserGroupResource struct {
	InfinityClient InfinityClient
}

type InfinityUserGroupResourceModel struct {
	ID                      types.String `tfsdk:"id"`
	ResourceID              types.Int32  `tfsdk:"resource_id"`
	Name                    types.String `tfsdk:"name"`
	Description             types.String `tfsdk:"description"`
	Users                   types.Set    `tfsdk:"users"`
	UserGroupEntityMappings types.Set    `tfsdk:"user_group_entity_mappings"`
}

func (r *InfinityUserGroupResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_infinity_user_group"
}

func (r *InfinityUserGroupResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *InfinityUserGroupResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Resource URI for the user group in Infinity",
			},
			"resource_id": schema.Int32Attribute{
				Computed:            true,
				MarkdownDescription: "The resource integer identifier for the user group in Infinity",
			},
			"name": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(250),
				},
				MarkdownDescription: "The unique name of the user group. Maximum length: 250 characters.",
			},
			"description": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(250),
				},
				MarkdownDescription: "A description of the user group. Maximum length: 250 characters.",
			},
			"users": schema.SetAttribute{
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "List of user resource URIs that belong to this group.",
			},
			"user_group_entity_mappings": schema.SetAttribute{
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "List of user group entity mapping resource URIs associated with this group.",
			},
		},
		MarkdownDescription: "Manages a user group configuration with the Infinity service.",
	}
}

func (r *InfinityUserGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	plan := &InfinityUserGroupResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createRequest := &config.UserGroupCreateRequest{
		Name: plan.Name.ValueString(),
	}

	// Only set optional fields if they are not null in the plan
	if !plan.Description.IsNull() {
		createRequest.Description = plan.Description.ValueString()
	}
	if !plan.Users.IsNull() && !plan.Users.IsUnknown() {
		var users []string
		resp.Diagnostics.Append(plan.Users.ElementsAs(ctx, &users, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		createRequest.Users = users
	}
	if !plan.UserGroupEntityMappings.IsNull() && !plan.UserGroupEntityMappings.IsUnknown() {
		var mappings []string
		resp.Diagnostics.Append(plan.UserGroupEntityMappings.ElementsAs(ctx, &mappings, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		createRequest.UserGroupEntityMappings = mappings
	}

	createResponse, err := r.InfinityClient.Config().CreateUserGroup(ctx, createRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Infinity user group",
			fmt.Sprintf("Could not create Infinity user group: %s", err),
		)
		return
	}

	resourceID, err := createResponse.ResourceID()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Retrieving Infinity user group ID",
			fmt.Sprintf("Could not retrieve ID for created Infinity user group: %s", err),
		)
		return
	}

	// Read the state from the API to get all computed values
	model, err := r.read(ctx, resourceID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Created Infinity user group",
			fmt.Sprintf("Could not read created Infinity user group with ID %d: %s", resourceID, err),
		)
		return
	}
	tflog.Trace(ctx, fmt.Sprintf("created Infinity user group with ID: %s, name: %s", model.ID, model.Name))

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *InfinityUserGroupResource) read(ctx context.Context, resourceID int) (*InfinityUserGroupResourceModel, error) {
	var data InfinityUserGroupResourceModel

	srv, err := r.InfinityClient.Config().GetUserGroup(ctx, resourceID)
	if err != nil {
		return nil, err
	}

	if srv.ResourceURI == "" {
		return nil, fmt.Errorf("user group with ID %d not found", resourceID)
	}

	data.ID = types.StringValue(srv.ResourceURI)
	data.ResourceID = types.Int32Value(int32(resourceID)) // #nosec G115 -- API values are expected to be within int32 range
	data.Name = types.StringValue(srv.Name)
	data.Description = types.StringValue(srv.Description)

	// Convert users to types.Set
	usersSet, diags := types.SetValueFrom(ctx, types.StringType, srv.Users)
	if diags.HasError() {
		return nil, fmt.Errorf("error converting users: %v", diags)
	}
	data.Users = usersSet

	// Convert user group entity mappings to types.Set (extract resource URIs)
	var mappingURIs []string
	if srv.UserGroupEntityMappings != nil {
		for _, mapping := range *srv.UserGroupEntityMappings {
			mappingURIs = append(mappingURIs, mapping.ResourceURI)
		}
	}
	mappingsSet, diags := types.SetValueFrom(ctx, types.StringType, mappingURIs)
	if diags.HasError() {
		return nil, fmt.Errorf("error converting user group entity mappings: %v", diags)
	}
	data.UserGroupEntityMappings = mappingsSet

	return &data, nil
}

func (r *InfinityUserGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	state := &InfinityUserGroupResourceModel{}

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
			"Error Reading Infinity user group",
			fmt.Sprintf("Could not read Infinity user group: %s", err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *InfinityUserGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	plan := &InfinityUserGroupResourceModel{}
	state := &InfinityUserGroupResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resourceID := int(state.ResourceID.ValueInt32())

	updateRequest := &config.UserGroupUpdateRequest{
		Name: plan.Name.ValueString(),
	}

	if !plan.Description.IsNull() {
		updateRequest.Description = plan.Description.ValueString()
	}
	if !plan.Users.IsNull() && !plan.Users.IsUnknown() {
		var users []string
		resp.Diagnostics.Append(plan.Users.ElementsAs(ctx, &users, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		updateRequest.Users = users
	}
	if !plan.UserGroupEntityMappings.IsNull() && !plan.UserGroupEntityMappings.IsUnknown() {
		var mappings []string
		resp.Diagnostics.Append(plan.UserGroupEntityMappings.ElementsAs(ctx, &mappings, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		updateRequest.UserGroupEntityMappings = mappings
	}

	_, err := r.InfinityClient.Config().UpdateUserGroup(ctx, resourceID, updateRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Infinity user group",
			fmt.Sprintf("Could not update Infinity user group with ID %d: %s", resourceID, err),
		)
		return
	}

	// Re-read the resource to get the latest state
	updatedModel, err := r.read(ctx, resourceID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Updated Infinity user group",
			fmt.Sprintf("Could not read updated Infinity user group with ID %d: %s", resourceID, err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, updatedModel)...)
}

func (r *InfinityUserGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	state := &InfinityUserGroupResourceModel{}

	tflog.Info(ctx, "Deleting Infinity user group")

	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.InfinityClient.Config().DeleteUserGroup(ctx, int(state.ResourceID.ValueInt32()))

	// Ignore 404 Not Found and Lookup errors on delete
	if err != nil && !isNotFoundError(err) && !isLookupError(err) {
		resp.Diagnostics.AddError(
			"Error Deleting Infinity user group",
			fmt.Sprintf("Could not delete Infinity user group with ID %s: %s", state.ID.ValueString(), err),
		)
		return
	}
}

func (r *InfinityUserGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resourceID, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Import ID must be a valid integer for the resource ID. Got: %s", req.ID),
		)
		return
	}

	tflog.Trace(ctx, fmt.Sprintf("Importing Infinity user group with resource ID: %d", resourceID))

	// Read the resource from the API
	model, err := r.read(ctx, resourceID)
	if err != nil {
		// Check if the error is a 404 (not found)
		if isNotFoundError(err) {
			resp.Diagnostics.AddError(
				"Infinity User Group Not Found",
				fmt.Sprintf("Infinity user group with resource ID %d not found.", resourceID),
			)
			return
		}
		resp.Diagnostics.AddError(
			"Error Importing Infinity User Group",
			fmt.Sprintf("Could not import Infinity user group with resource ID %d: %s", resourceID, err),
		)
		return
	}

	// Set the state from the imported resource
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}
