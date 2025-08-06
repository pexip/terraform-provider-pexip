/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/pexip/go-infinity-sdk/v38/config"
)

var (
	_ resource.ResourceWithImportState = (*InfinityLicenceRequestResource)(nil)
)

type InfinityLicenceRequestResource struct {
	InfinityClient InfinityClient
}

type InfinityLicenceRequestResourceModel struct {
	ID             types.String `tfsdk:"id"`
	SequenceNumber types.String `tfsdk:"sequence_number"`
	Reference      types.String `tfsdk:"reference"`
	Actions        types.String `tfsdk:"actions"`
	GenerationTime types.String `tfsdk:"generation_time"`
	Status         types.String `tfsdk:"status"`
	ResponseXML    types.String `tfsdk:"response_xml"`
}

func (r *InfinityLicenceRequestResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_infinity_licence_request"
}

func (r *InfinityLicenceRequestResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *InfinityLicenceRequestResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Resource URI for the licence request in Infinity",
			},
			"sequence_number": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The sequence number generated for this licence request",
			},
			"reference": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					stringvalidator.LengthAtMost(100),
				},
				MarkdownDescription: "A reference identifier for this licence request. Maximum length: 100 characters.",
			},
			"actions": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
				MarkdownDescription: "The actions XML for the licence request defining what licenses are being requested.",
			},
			"generation_time": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The timestamp when this licence request was generated.",
			},
			"status": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The current status of the licence request.",
			},
			"response_xml": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The response XML from the license server, if available.",
			},
		},
		MarkdownDescription: "Manages a licence request with the Infinity service. Licence requests are used to generate license activation files that can be submitted to Pexip for license provisioning. Note: This resource only supports creation and reading - licence requests cannot be updated or deleted once created.",
	}
}

func (r *InfinityLicenceRequestResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	plan := &InfinityLicenceRequestResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createRequest := &config.LicenceRequestCreateRequest{
		Reference: plan.Reference.ValueString(),
		Actions:   plan.Actions.ValueString(),
	}

	createResponse, err := r.InfinityClient.Config().CreateLicenceRequest(ctx, createRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Infinity licence request",
			fmt.Sprintf("Could not create Infinity licence request: %s", err),
		)
		return
	}

	// For licence requests, the sequence number is returned in the response location
	sequenceNumber, err := createResponse.ResourceID()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Retrieving Infinity licence request sequence number",
			fmt.Sprintf("Could not retrieve sequence number for created Infinity licence request: %s", err),
		)
		return
	}

	// The resource uses sequence number as string identifier
	sequenceNumberStr := fmt.Sprintf("%d", sequenceNumber)

	// Read the state from the API to get all computed values
	model, err := r.read(ctx, sequenceNumberStr)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Created Infinity licence request",
			fmt.Sprintf("Could not read created Infinity licence request with sequence number %s: %s", sequenceNumberStr, err),
		)
		return
	}
	tflog.Trace(ctx, fmt.Sprintf("created Infinity licence request with ID: %s, sequence: %s", model.ID, model.SequenceNumber))

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *InfinityLicenceRequestResource) read(ctx context.Context, sequenceNumber string) (*InfinityLicenceRequestResourceModel, error) {
	var data InfinityLicenceRequestResourceModel

	srv, err := r.InfinityClient.Config().GetLicenceRequest(ctx, sequenceNumber)
	if err != nil {
		return nil, err
	}

	if srv.ResourceURI == "" {
		return nil, fmt.Errorf("licence request with sequence number %s not found", sequenceNumber)
	}

	data.ID = types.StringValue(srv.ResourceURI)
	data.SequenceNumber = types.StringValue(srv.SequenceNumber)
	data.Reference = types.StringValue(srv.Reference)
	data.Actions = types.StringValue(srv.Actions)
	data.GenerationTime = types.StringValue(srv.GenerationTime)
	data.Status = types.StringValue(srv.Status)

	if srv.ResponseXML != nil {
		data.ResponseXML = types.StringValue(*srv.ResponseXML)
	} else {
		data.ResponseXML = types.StringNull()
	}

	return &data, nil
}

func (r *InfinityLicenceRequestResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	state := &InfinityLicenceRequestResourceModel{}

	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	sequenceNumber := state.SequenceNumber.ValueString()
	state, err := r.read(ctx, sequenceNumber)
	if err != nil {
		// Check if the error is a 404 (not found)
		if isNotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading Infinity licence request",
			fmt.Sprintf("Could not read Infinity licence request: %s", err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *InfinityLicenceRequestResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError(
		"Update Not Supported",
		"Licence request resources cannot be updated. Licence requests are immutable once created.",
	)
}

func (r *InfinityLicenceRequestResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	resp.State.RemoveResource(ctx)
}

func (r *InfinityLicenceRequestResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	sequenceNumber := req.ID

	tflog.Trace(ctx, fmt.Sprintf("Importing Infinity licence request with sequence number: %s", sequenceNumber))

	// Read the resource from the API
	model, err := r.read(ctx, sequenceNumber)
	if err != nil {
		// Check if the error is a 404 (not found)
		if isNotFoundError(err) {
			resp.Diagnostics.AddError(
				"Infinity Licence Request Not Found",
				fmt.Sprintf("Infinity licence request with sequence number %s not found.", sequenceNumber),
			)
			return
		}
		resp.Diagnostics.AddError(
			"Error Importing Infinity Licence Request",
			fmt.Sprintf("Could not import Infinity licence request with sequence number %s: %s", sequenceNumber, err),
		)
		return
	}

	// Set the state from the imported resource
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}
