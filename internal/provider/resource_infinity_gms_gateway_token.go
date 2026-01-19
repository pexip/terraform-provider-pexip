/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/pexip/go-infinity-sdk/v38/config"
)

var (
	_ resource.ResourceWithImportState = (*InfinityGMSGatewayTokenResource)(nil)
)

type InfinityGMSGatewayTokenResource struct {
	InfinityClient InfinityClient
}

type InfinityGMSGatewayTokenResourceModel struct {
	ID                      types.String `tfsdk:"id"`
	Certificate             types.String `tfsdk:"certificate"`
	IntermediateCertificate types.String `tfsdk:"intermediate_certificate"`
	LeafCertificate         types.String `tfsdk:"leaf_certificate"`
	PrivateKey              types.String `tfsdk:"private_key"`
	SupportsDirectGuestJoin types.Bool   `tfsdk:"supports_direct_guest_join"`
	ResourceURI             types.String `tfsdk:"resource_uri"`
}

func (r *InfinityGMSGatewayTokenResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_infinity_gms_gateway_token"
}

func (r *InfinityGMSGatewayTokenResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *InfinityGMSGatewayTokenResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Resource URI for the GMS gateway token configuration in Infinity",
			},
			"certificate": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The Google Meet gateway token certificate.",
			},
			"intermediate_certificate": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The intermediate certificate for the Google Meet gateway token.",
			},
			"leaf_certificate": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The leaf certificate for the Google Meet gateway token.",
			},
			"private_key": schema.StringAttribute{
				Required:            true,
				Sensitive:           true,
				MarkdownDescription: "The private key for the Google Meet gateway token.",
			},
			"supports_direct_guest_join": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Whether the Google Meet gateway token supports direct guest join.",
			},
			"resource_uri": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Resource URI for the GMS gateway token configuration.",
			},
		},
		MarkdownDescription: "Manages the Google Meet gateway token configuration with the Infinity service. This is a singleton resource - only one GMS gateway token configuration exists per system.",
	}
}

func (r *InfinityGMSGatewayTokenResource) buildUpdateRequest(plan *InfinityGMSGatewayTokenResourceModel) *config.GMSGatewayTokenUpdateRequest {
	updateRequest := &config.GMSGatewayTokenUpdateRequest{
		Certificate: plan.Certificate.ValueString(),
	}

	// Handle optional pointer fields
	if !plan.IntermediateCertificate.IsNull() && !plan.IntermediateCertificate.IsUnknown() {
		val := plan.IntermediateCertificate.ValueString()
		updateRequest.IntermediateCertificate = &val
	}
	if !plan.LeafCertificate.IsNull() && !plan.LeafCertificate.IsUnknown() {
		val := plan.LeafCertificate.ValueString()
		updateRequest.LeafCertificate = &val
	}
	if !plan.PrivateKey.IsNull() && !plan.PrivateKey.IsUnknown() {
		val := plan.PrivateKey.ValueString()
		updateRequest.PrivateKey = &val
	}
	if !plan.SupportsDirectGuestJoin.IsNull() && !plan.SupportsDirectGuestJoin.IsUnknown() {
		val := plan.SupportsDirectGuestJoin.ValueBool()
		updateRequest.SupportsDirectGuestJoin = &val
	}

	return updateRequest
}

func (r *InfinityGMSGatewayTokenResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// For singleton resources, Create is actually Update since the resource always exists
	plan := &InfinityGMSGatewayTokenResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateRequest := r.buildUpdateRequest(plan)

	_, err := r.InfinityClient.Config().UpdateGMSGatewayToken(ctx, updateRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Infinity GMS gateway token",
			fmt.Sprintf("Could not update Infinity GMS gateway token: %s", err),
		)
		return
	}

	// Re-read the resource to get the latest state
	updatedModel, err := r.read(ctx, plan.Certificate.ValueString(), plan.PrivateKey.ValueStringPointer())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Updated Infinity GMS gateway token",
			fmt.Sprintf("Could not read updated Infinity GMS gateway token: %s", err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, updatedModel)...)
}

func (r *InfinityGMSGatewayTokenResource) read(ctx context.Context, cert string, privateKey *string) (*InfinityGMSGatewayTokenResourceModel, error) {
	var data InfinityGMSGatewayTokenResourceModel

	srv, err := r.InfinityClient.Config().GetGMSGatewayToken(ctx)
	if err != nil {
		return nil, err
	}

	if srv.ResourceURI == "" {
		return nil, fmt.Errorf("GMS gateway token configuration not found")
	}

	data.ID = types.StringValue(srv.ResourceURI)
	data.Certificate = types.StringValue(cert)
	data.IntermediateCertificate = types.StringPointerValue(srv.IntermediateCertificate)
	data.LeafCertificate = types.StringPointerValue(srv.LeafCertificate)
	// Preserve private key from state since it's not returned by the API
	data.PrivateKey = types.StringPointerValue(privateKey)
	if srv.SupportsDirectGuestJoin != nil {
		data.SupportsDirectGuestJoin = types.BoolValue(*srv.SupportsDirectGuestJoin)
	} else {
		data.SupportsDirectGuestJoin = types.BoolValue(false)
	}
	data.ResourceURI = types.StringValue(srv.ResourceURI)

	return &data, nil
}

func (r *InfinityGMSGatewayTokenResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	state := &InfinityGMSGatewayTokenResourceModel{}

	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	state, err := r.read(ctx, state.Certificate.ValueString(), state.PrivateKey.ValueStringPointer())
	if err != nil {
		// Check if the error is a 404 (not found) - unlikely for singleton resources
		if isNotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading Infinity GMS gateway token",
			fmt.Sprintf("Could not read Infinity GMS gateway token: %s", err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *InfinityGMSGatewayTokenResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	plan := &InfinityGMSGatewayTokenResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateRequest := r.buildUpdateRequest(plan)

	_, err := r.InfinityClient.Config().UpdateGMSGatewayToken(ctx, updateRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Infinity GMS gateway token",
			fmt.Sprintf("Could not update Infinity GMS gateway token: %s", err),
		)
		return
	}

	// Re-read the resource to get the latest state
	updatedModel, err := r.read(ctx, plan.Certificate.ValueString(), plan.PrivateKey.ValueStringPointer())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Updated Infinity GMS gateway token",
			fmt.Sprintf("Could not read updated Infinity GMS gateway token: %s", err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, updatedModel)...)
}

func (r *InfinityGMSGatewayTokenResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// For singleton resources, delete means resetting to default/minimal values
	tflog.Info(ctx, "The Infinity SDK does not yet support deleting the GMS gateway token. It will be removed from state.")

	// is this needed?
	//resp.State.RemoveResource(ctx)
}

func (r *InfinityGMSGatewayTokenResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// For singleton resources, the import ID doesn't matter since there's only one instance
	tflog.Trace(ctx, "Importing Infinity GMS gateway token")

	// Read the resource from the API
	model, err := r.read(ctx, "", nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Importing Infinity GMS Gateway Token",
			fmt.Sprintf("Could not import Infinity GMS gateway token: %s", err),
		)
		return
	}

	// Set the state from the imported resource
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}
