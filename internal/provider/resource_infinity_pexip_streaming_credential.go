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
	_ resource.ResourceWithImportState = (*InfinityPexipStreamingCredentialResource)(nil)
)

type InfinityPexipStreamingCredentialResource struct {
	InfinityClient InfinityClient
}

type InfinityPexipStreamingCredentialResourceModel struct {
	ID         types.String `tfsdk:"id"`
	ResourceID types.Int32  `tfsdk:"resource_id"`
	Kid        types.String `tfsdk:"kid"`
	PublicKey  types.String `tfsdk:"public_key"`
}

func (r *InfinityPexipStreamingCredentialResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_infinity_pexip_streaming_credential"
}

func (r *InfinityPexipStreamingCredentialResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *InfinityPexipStreamingCredentialResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Resource URI for the Pexip Streaming credential in Infinity",
			},
			"resource_id": schema.Int32Attribute{
				Computed:            true,
				MarkdownDescription: "The resource integer identifier for the Pexip Streaming credential in Infinity",
			},
			"kid": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					stringvalidator.LengthAtMost(100),
				},
				MarkdownDescription: "The key ID (kid) for the Pexip Streaming credential. Maximum length: 100 characters.",
			},
			"public_key": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
				MarkdownDescription: "The public key for the Pexip Streaming credential. This should be a valid public key in appropriate format.",
			},
		},
		MarkdownDescription: "Manages a Pexip Streaming credential with the Infinity service. Pexip Streaming credentials are used to authenticate and authorize streaming services, enabling secure content delivery and media streaming functionality within Pexip Infinity deployments.",
	}
}

func (r *InfinityPexipStreamingCredentialResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	plan := &InfinityPexipStreamingCredentialResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createRequest := &config.PexipStreamingCredentialCreateRequest{
		Kid:       plan.Kid.ValueString(),
		PublicKey: plan.PublicKey.ValueString(),
	}

	createResponse, err := r.InfinityClient.Config().CreatePexipStreamingCredential(ctx, createRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Infinity Pexip Streaming credential",
			fmt.Sprintf("Could not create Infinity Pexip Streaming credential: %s", err),
		)
		return
	}

	resourceID, err := createResponse.ResourceID()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Retrieving Infinity Pexip Streaming credential ID",
			fmt.Sprintf("Could not retrieve ID for created Infinity Pexip Streaming credential: %s", err),
		)
		return
	}

	// Read the state from the API to get all computed values
	model, err := r.read(ctx, resourceID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Created Infinity Pexip Streaming credential",
			fmt.Sprintf("Could not read created Infinity Pexip Streaming credential with ID %d: %s", resourceID, err),
		)
		return
	}
	tflog.Trace(ctx, fmt.Sprintf("created Infinity Pexip Streaming credential with ID: %s, kid: %s", model.ID, model.Kid))

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *InfinityPexipStreamingCredentialResource) read(ctx context.Context, resourceID int) (*InfinityPexipStreamingCredentialResourceModel, error) {
	var data InfinityPexipStreamingCredentialResourceModel

	srv, err := r.InfinityClient.Config().GetPexipStreamingCredential(ctx, resourceID)
	if err != nil {
		return nil, err
	}

	if srv.ResourceURI == "" {
		return nil, fmt.Errorf("pexip Streaming credential with ID %d not found", resourceID)
	}

	data.ID = types.StringValue(srv.ResourceURI)
	data.ResourceID = types.Int32Value(int32(resourceID)) // #nosec G115 -- API values are expected to be within int32 range
	data.Kid = types.StringValue(srv.Kid)
	data.PublicKey = types.StringValue(srv.PublicKey)

	return &data, nil
}

func (r *InfinityPexipStreamingCredentialResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	state := &InfinityPexipStreamingCredentialResourceModel{}

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
			"Error Reading Infinity Pexip Streaming credential",
			fmt.Sprintf("Could not read Infinity Pexip Streaming credential: %s", err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *InfinityPexipStreamingCredentialResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	plan := &InfinityPexipStreamingCredentialResourceModel{}
	state := &InfinityPexipStreamingCredentialResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateRequest := &config.PexipStreamingCredentialUpdateRequest{
		Kid:       plan.Kid.ValueString(),
		PublicKey: plan.PublicKey.ValueString(),
	}

	resourceID := int(state.ResourceID.ValueInt32())
	_, err := r.InfinityClient.Config().UpdatePexipStreamingCredential(ctx, resourceID, updateRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Infinity Pexip Streaming credential",
			fmt.Sprintf("Could not update Infinity Pexip Streaming credential: %s", err),
		)
		return
	}

	// Read the state from the API to get all computed values
	model, err := r.read(ctx, resourceID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Updated Infinity Pexip Streaming credential",
			fmt.Sprintf("Could not read updated Infinity Pexip Streaming credential with ID %d: %s", resourceID, err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *InfinityPexipStreamingCredentialResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	state := &InfinityPexipStreamingCredentialResourceModel{}

	tflog.Info(ctx, "Deleting Infinity Pexip Streaming credential")

	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.InfinityClient.Config().DeletePexipStreamingCredential(ctx, int(state.ResourceID.ValueInt32()))

	// Ignore 404 Not Found and Lookup errors on delete
	if err != nil && !isNotFoundError(err) && !isLookupError(err) {
		resp.Diagnostics.AddError(
			"Error Deleting Infinity Pexip Streaming credential",
			fmt.Sprintf("Could not delete Infinity Pexip Streaming credential with ID %s: %s", state.ID.ValueString(), err),
		)
		return
	}
}

func (r *InfinityPexipStreamingCredentialResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resourceID, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Import ID must be a valid integer for the resource ID. Got: %s", req.ID),
		)
		return
	}

	tflog.Trace(ctx, fmt.Sprintf("Importing Infinity Pexip Streaming credential with resource ID: %d", resourceID))

	// Read the resource from the API
	model, err := r.read(ctx, resourceID)
	if err != nil {
		// Check if the error is a 404 (not found)
		if isNotFoundError(err) {
			resp.Diagnostics.AddError(
				"Infinity Pexip Streaming Credential Not Found",
				fmt.Sprintf("Infinity Pexip Streaming credential with resource ID %d not found.", resourceID),
			)
			return
		}
		resp.Diagnostics.AddError(
			"Error Importing Infinity Pexip Streaming Credential",
			fmt.Sprintf("Could not import Infinity Pexip Streaming credential with resource ID %d: %s", resourceID, err),
		)
		return
	}

	// Set the state from the imported resource
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}
