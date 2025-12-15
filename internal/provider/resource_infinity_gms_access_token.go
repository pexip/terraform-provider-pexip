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
	_ resource.ResourceWithImportState = (*InfinityGMSAccessTokenResource)(nil)
)

type InfinityGMSAccessTokenResource struct {
	InfinityClient InfinityClient
}

type InfinityGMSAccessTokenResourceModel struct {
	ID         types.String `tfsdk:"id"`
	ResourceID types.Int32  `tfsdk:"resource_id"`
	Name       types.String `tfsdk:"name"`
	Token      types.String `tfsdk:"token"`
}

func (r *InfinityGMSAccessTokenResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_infinity_gms_access_token"
}

func (r *InfinityGMSAccessTokenResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *InfinityGMSAccessTokenResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Resource URI for the GMS access token in Infinity",
			},
			"resource_id": schema.Int32Attribute{
				Computed:            true,
				MarkdownDescription: "The resource integer identifier for the GMS access token in Infinity",
			},
			"name": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					stringvalidator.LengthAtMost(250),
				},
				MarkdownDescription: "The name of the GMS access token. Maximum length: 250 characters.",
			},
			"token": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				MarkdownDescription: "The Google Meet access token. This field is sensitive and will not be displayed in logs.",
			},
		},
		MarkdownDescription: "Manages a Google Meet Service (GMS) access token configuration with the Infinity service. These tokens are used for Google Meet integration.",
	}
}

func (r *InfinityGMSAccessTokenResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	plan := &InfinityGMSAccessTokenResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createRequest := &config.GMSAccessTokenCreateRequest{
		Name:  plan.Name.ValueString(),
		Token: plan.Token.ValueString(),
	}

	createResponse, err := r.InfinityClient.Config().CreateGMSAccessToken(ctx, createRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Infinity GMS access token",
			fmt.Sprintf("Could not create Infinity GMS access token: %s", err),
		)
		return
	}

	resourceID, err := createResponse.ResourceID()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Retrieving Infinity GMS access token ID",
			fmt.Sprintf("Could not retrieve ID for created Infinity GMS access token: %s", err),
		)
		return
	}

	// Read the state from the API to get all computed values
	model, err := r.read(ctx, resourceID, plan.Token.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Created Infinity GMS access token",
			fmt.Sprintf("Could not read created Infinity GMS access token with ID %d: %s", resourceID, err),
		)
		return
	}
	tflog.Trace(ctx, fmt.Sprintf("created Infinity GMS access token with ID: %s, name: %s", model.ID, model.Name))

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *InfinityGMSAccessTokenResource) read(ctx context.Context, resourceID int, token string) (*InfinityGMSAccessTokenResourceModel, error) {
	var data InfinityGMSAccessTokenResourceModel

	srv, err := r.InfinityClient.Config().GetGMSAccessToken(ctx, resourceID)
	if err != nil {
		return nil, err
	}

	if srv.ResourceURI == "" {
		return nil, fmt.Errorf("GMS access token with ID %d not found", resourceID)
	}

	data.ID = types.StringValue(srv.ResourceURI)
	data.ResourceID = types.Int32Value(int32(resourceID)) // #nosec G115 -- API values are expected to be within int32 range
	data.Name = types.StringValue(srv.Name)
	// The token is sensitive and not returned by the API, so we use the value from the plan/state
	data.Token = types.StringValue(token)

	return &data, nil
}

func (r *InfinityGMSAccessTokenResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	state := &InfinityGMSAccessTokenResourceModel{}

	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resourceID := int(state.ResourceID.ValueInt32())
	state, err := r.read(ctx, resourceID, state.Token.ValueString())
	if err != nil {
		// Check if the error is a 404 (not found)
		if isNotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading Infinity GMS access token",
			fmt.Sprintf("Could not read Infinity GMS access token: %s", err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *InfinityGMSAccessTokenResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	plan := &InfinityGMSAccessTokenResourceModel{}
	state := &InfinityGMSAccessTokenResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resourceID := int(state.ResourceID.ValueInt32())

	updateRequest := &config.GMSAccessTokenUpdateRequest{
		Name:  plan.Name.ValueString(),
		Token: plan.Token.ValueString(),
	}

	_, err := r.InfinityClient.Config().UpdateGMSAccessToken(ctx, resourceID, updateRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Infinity GMS access token",
			fmt.Sprintf("Could not update Infinity GMS access token with ID %d: %s", resourceID, err),
		)
		return
	}

	// Re-read the resource to get the latest state
	updatedModel, err := r.read(ctx, resourceID, plan.Token.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Updated Infinity GMS access token",
			fmt.Sprintf("Could not read updated Infinity GMS access token with ID %d: %s", resourceID, err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, updatedModel)...)
}

func (r *InfinityGMSAccessTokenResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	state := &InfinityGMSAccessTokenResourceModel{}

	tflog.Info(ctx, "Deleting Infinity GMS access token")

	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.InfinityClient.Config().DeleteGMSAccessToken(ctx, int(state.ResourceID.ValueInt32()))

	// Ignore 404 Not Found and Lookup errors on delete
	if err != nil && !isNotFoundError(err) && !isLookupError(err) {
		resp.Diagnostics.AddError(
			"Error Deleting Infinity GMS access token",
			fmt.Sprintf("Could not delete Infinity GMS access token with ID %s: %s", state.ID.ValueString(), err),
		)
		return
	}
}

func (r *InfinityGMSAccessTokenResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resourceID, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Import ID must be a valid integer for the resource ID. Got: %s", req.ID),
		)
		return
	}

	tflog.Trace(ctx, fmt.Sprintf("Importing Infinity GMS access token with resource ID: %d", resourceID))

	// Read the resource from the API
	// Note: Token cannot be retrieved during import, so it will be empty
	model, err := r.read(ctx, resourceID, "")
	if err != nil {
		// Check if the error is a 404 (not found)
		if isNotFoundError(err) {
			resp.Diagnostics.AddError(
				"Infinity GMS Access Token Not Found",
				fmt.Sprintf("Infinity GMS access token with resource ID %d not found.", resourceID),
			)
			return
		}
		resp.Diagnostics.AddError(
			"Error Importing Infinity GMS Access Token",
			fmt.Sprintf("Could not import Infinity GMS access token with resource ID %d: %s", resourceID, err),
		)
		return
	}

	// Set the state from the imported resource
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}
