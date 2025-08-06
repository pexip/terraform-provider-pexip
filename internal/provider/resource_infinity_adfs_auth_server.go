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

	"github.com/pexip/terraform-provider-pexip/internal/helpers"
)

var (
	_ resource.ResourceWithImportState = (*InfinityADFSAuthServerResource)(nil)
)

type InfinityADFSAuthServerResource struct {
	InfinityClient InfinityClient
}

type InfinityADFSAuthServerResourceModel struct {
	ID                             types.String `tfsdk:"id"`
	ResourceID                     types.Int32  `tfsdk:"resource_id"`
	Name                           types.String `tfsdk:"name"`
	Description                    types.String `tfsdk:"description"`
	ClientID                       types.String `tfsdk:"client_id"`
	FederationServiceName          types.String `tfsdk:"federation_service_name"`
	FederationServiceIdentifier    types.String `tfsdk:"federation_service_identifier"`
	RelyingPartyTrustIdentifierURL types.String `tfsdk:"relying_party_trust_identifier_url"`
}

func (r *InfinityADFSAuthServerResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_infinity_adfs_auth_server"
}

func (r *InfinityADFSAuthServerResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *InfinityADFSAuthServerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Resource URI for the ADFS auth server in Infinity",
			},
			"resource_id": schema.Int32Attribute{
				Computed:            true,
				MarkdownDescription: "The resource integer identifier for the ADFS auth server in Infinity",
			},
			"name": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(250),
				},
				MarkdownDescription: "The unique name of the ADFS auth server. Maximum length: 250 characters.",
			},
			"description": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(250),
				},
				MarkdownDescription: "A description of the ADFS auth server. Maximum length: 250 characters.",
			},
			"client_id": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(250),
				},
				MarkdownDescription: "The client ID for the ADFS OAuth 2.0 client. Maximum length: 250 characters.",
			},
			"federation_service_name": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(250),
				},
				MarkdownDescription: "The federation service name. Maximum length: 250 characters.",
			},
			"federation_service_identifier": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(250),
				},
				MarkdownDescription: "The federation service identifier. Maximum length: 250 characters.",
			},
			"relying_party_trust_identifier_url": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(250),
				},
				MarkdownDescription: "The relying party trust identifier URL. Maximum length: 250 characters.",
			},
		},
		MarkdownDescription: "Manages an ADFS OAuth 2.0 auth server configuration with the Infinity service.",
	}
}

func (r *InfinityADFSAuthServerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	plan := &InfinityADFSAuthServerResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createRequest := &config.ADFSAuthServerCreateRequest{
		Name:                           plan.Name.ValueString(),
		ClientID:                       plan.ClientID.ValueString(),
		FederationServiceName:          plan.FederationServiceName.ValueString(),
		FederationServiceIdentifier:    plan.FederationServiceIdentifier.ValueString(),
		RelyingPartyTrustIdentifierURL: plan.RelyingPartyTrustIdentifierURL.ValueString(),
	}

	// Set optional fields
	if !plan.Description.IsNull() {
		createRequest.Description = plan.Description.ValueString()
	}

	createResponse, err := r.InfinityClient.Config().CreateADFSAuthServer(ctx, createRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Infinity ADFS auth server",
			fmt.Sprintf("Could not create Infinity ADFS auth server: %s", err),
		)
		return
	}

	resourceID, err := createResponse.ResourceID()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Retrieving Infinity ADFS auth server ID",
			fmt.Sprintf("Could not retrieve ID for created Infinity ADFS auth server: %s", err),
		)
		return
	}

	// Read the state from the API to get all computed values
	model, err := r.read(ctx, resourceID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Created Infinity ADFS auth server",
			fmt.Sprintf("Could not read created Infinity ADFS auth server with ID %d: %s", resourceID, err),
		)
		return
	}
	tflog.Trace(ctx, fmt.Sprintf("created Infinity ADFS auth server with ID: %s, name: %s", model.ID, model.Name))

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *InfinityADFSAuthServerResource) read(ctx context.Context, resourceID int) (*InfinityADFSAuthServerResourceModel, error) {
	var data InfinityADFSAuthServerResourceModel

	srv, err := r.InfinityClient.Config().GetADFSAuthServer(ctx, resourceID)
	if err != nil {
		return nil, err
	}

	if srv.ResourceURI == "" {
		return nil, fmt.Errorf("ADFS auth server with ID %d not found", resourceID)
	}

	data.ID = types.StringValue(srv.ResourceURI)
	resourceID32, err := helpers.SafeInt32(resourceID)
	if err != nil {
		return nil, fmt.Errorf("resource ID conversion error: %w", err)
	}
	data.ResourceID = types.Int32Value(resourceID32)
	data.Name = types.StringValue(srv.Name)
	data.Description = types.StringValue(srv.Description)
	data.ClientID = types.StringValue(srv.ClientID)
	data.FederationServiceName = types.StringValue(srv.FederationServiceName)
	data.FederationServiceIdentifier = types.StringValue(srv.FederationServiceIdentifier)
	data.RelyingPartyTrustIdentifierURL = types.StringValue(srv.RelyingPartyTrustIdentifierURL)

	return &data, nil
}

func (r *InfinityADFSAuthServerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	state := &InfinityADFSAuthServerResourceModel{}

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
			"Error Reading Infinity ADFS auth server",
			fmt.Sprintf("Could not read Infinity ADFS auth server: %s", err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *InfinityADFSAuthServerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	plan := &InfinityADFSAuthServerResourceModel{}
	state := &InfinityADFSAuthServerResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resourceID := int(state.ResourceID.ValueInt32())

	updateRequest := &config.ADFSAuthServerUpdateRequest{
		Name:                           plan.Name.ValueString(),
		ClientID:                       plan.ClientID.ValueString(),
		FederationServiceName:          plan.FederationServiceName.ValueString(),
		FederationServiceIdentifier:    plan.FederationServiceIdentifier.ValueString(),
		RelyingPartyTrustIdentifierURL: plan.RelyingPartyTrustIdentifierURL.ValueString(),
	}

	// Set optional fields
	if !plan.Description.IsNull() {
		updateRequest.Description = plan.Description.ValueString()
	}

	_, err := r.InfinityClient.Config().UpdateADFSAuthServer(ctx, resourceID, updateRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Infinity ADFS auth server",
			fmt.Sprintf("Could not update Infinity ADFS auth server with ID %d: %s", resourceID, err),
		)
		return
	}

	// Re-read the resource to get the latest state
	updatedModel, err := r.read(ctx, resourceID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Updated Infinity ADFS auth server",
			fmt.Sprintf("Could not read updated Infinity ADFS auth server with ID %d: %s", resourceID, err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, updatedModel)...)
}

func (r *InfinityADFSAuthServerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	state := &InfinityADFSAuthServerResourceModel{}

	tflog.Info(ctx, "Deleting Infinity ADFS auth server")

	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.InfinityClient.Config().DeleteADFSAuthServer(ctx, int(state.ResourceID.ValueInt32()))

	// Ignore 404 Not Found and Lookup errors on delete
	if err != nil && !isNotFoundError(err) && !isLookupError(err) {
		resp.Diagnostics.AddError(
			"Error Deleting Infinity ADFS auth server",
			fmt.Sprintf("Could not delete Infinity ADFS auth server with ID %s: %s", state.ID.ValueString(), err),
		)
		return
	}
}

func (r *InfinityADFSAuthServerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resourceID, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Import ID must be a valid integer for the resource ID. Got: %s", req.ID),
		)
		return
	}

	tflog.Trace(ctx, fmt.Sprintf("Importing Infinity ADFS auth server with resource ID: %d", resourceID))

	// Read the resource from the API
	model, err := r.read(ctx, resourceID)
	if err != nil {
		// Check if the error is a 404 (not found)
		if isNotFoundError(err) {
			resp.Diagnostics.AddError(
				"Infinity ADFS Auth Server Not Found",
				fmt.Sprintf("Infinity ADFS auth server with resource ID %d not found.", resourceID),
			)
			return
		}
		resp.Diagnostics.AddError(
			"Error Importing Infinity ADFS Auth Server",
			fmt.Sprintf("Could not import Infinity ADFS auth server with resource ID %d: %s", resourceID, err),
		)
		return
	}

	// Set the state from the imported resource
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}
