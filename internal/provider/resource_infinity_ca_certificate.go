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
	_ resource.ResourceWithImportState = (*InfinityCACertificateResource)(nil)
)

type InfinityCACertificateResource struct {
	InfinityClient InfinityClient
}

type InfinityCACertificateResourceModel struct {
	ID                  types.String `tfsdk:"id"`
	ResourceID          types.Int32  `tfsdk:"resource_id"`
	Certificate         types.String `tfsdk:"certificate"`
	TrustedIntermediate types.Bool   `tfsdk:"trusted_intermediate"`
	StartDate           types.String `tfsdk:"start_date"`
	EndDate             types.String `tfsdk:"end_date"`
	SubjectName         types.String `tfsdk:"subject_name"`
	SubjectHash         types.String `tfsdk:"subject_hash"`
	RawSubject          types.String `tfsdk:"raw_subject"`
	IssuerName          types.String `tfsdk:"issuer_name"`
	IssuerHash          types.String `tfsdk:"issuer_hash"`
	RawIssuer           types.String `tfsdk:"raw_issuer"`
	SerialNo            types.String `tfsdk:"serial_no"`
	KeyID               types.String `tfsdk:"key_id"`
	IssuerKeyID         types.String `tfsdk:"issuer_key_id"`
	Text                types.String `tfsdk:"text"`
}

func (r *InfinityCACertificateResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_infinity_ca_certificate"
}

func (r *InfinityCACertificateResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *InfinityCACertificateResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Resource URI for the CA certificate in Infinity",
			},
			"resource_id": schema.Int32Attribute{
				Computed:            true,
				MarkdownDescription: "The resource integer identifier for the CA certificate in Infinity",
			},
			"certificate": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
				MarkdownDescription: "The PEM-encoded CA certificate content.",
			},
			"trusted_intermediate": schema.BoolAttribute{
				Computed:            true,
				Optional:            true,
				MarkdownDescription: "Whether this CA certificate is trusted as an intermediate certificate.",
			},
			"start_date": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The start date of the certificate validity period.",
			},
			"end_date": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The end date of the certificate validity period.",
			},
			"subject_name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The subject name of the certificate.",
			},
			"subject_hash": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The subject hash of the certificate.",
			},
			"raw_subject": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The raw subject of the certificate.",
			},
			"issuer_name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The issuer name of the certificate.",
			},
			"issuer_hash": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The issuer hash of the certificate.",
			},
			"raw_issuer": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The raw issuer of the certificate.",
			},
			"serial_no": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The serial number of the certificate.",
			},
			"key_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The key identifier of the certificate.",
			},
			"issuer_key_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The issuer key identifier of the certificate.",
			},
			"text": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The text representation of the certificate.",
			},
		},
		MarkdownDescription: "Manages a CA certificate configuration with the Infinity service.",
	}
}

func (r *InfinityCACertificateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	plan := &InfinityCACertificateResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createRequest := &config.CACertificateCreateRequest{
		Certificate:         plan.Certificate.ValueString(),
	}

	createResponse, err := r.InfinityClient.Config().CreateCACertificate(ctx, createRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Infinity CA certificate",
			fmt.Sprintf("Could not create Infinity CA certificate: %s", err),
		)
		return
	}

	resourceID, err := createResponse.ResourceID()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Retrieving Infinity CA certificate ID",
			fmt.Sprintf("Could not retrieve ID for created Infinity CA certificate: %s", err),
		)
		return
	}

	// Read the state from the API to get all computed values
	model, err := r.read(ctx, resourceID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Created Infinity CA certificate",
			fmt.Sprintf("Could not read created Infinity CA certificate with ID %d: %s", resourceID, err),
		)
		return
	}
	tflog.Trace(ctx, fmt.Sprintf("created Infinity CA certificate with ID: %s, subject: %s", model.ID, model.SubjectName))

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *InfinityCACertificateResource) read(ctx context.Context, resourceID int) (*InfinityCACertificateResourceModel, error) {
	var data InfinityCACertificateResourceModel

	srv, err := r.InfinityClient.Config().GetCACertificate(ctx, resourceID)
	if err != nil {
		return nil, err
	}

	if srv.ResourceURI == "" {
		return nil, fmt.Errorf("CA certificate with ID %d not found", resourceID)
	}

	data.ID = types.StringValue(srv.ResourceURI)
	data.ResourceID = types.Int32Value(int32(resourceID)) // #nosec G115 -- API values are expected to be within int32 range
	data.Certificate = types.StringValue(srv.Certificate)
	data.TrustedIntermediate = types.BoolValue(srv.TrustedIntermediate)
	data.StartDate = types.StringValue(srv.StartDate.String())
	data.EndDate = types.StringValue(srv.EndDate.String())
	data.SubjectName = types.StringValue(srv.SubjectName)
	data.SubjectHash = types.StringValue(srv.SubjectHash)
	data.RawSubject = types.StringValue(srv.RawSubject)
	data.IssuerName = types.StringValue(srv.IssuerName)
	data.IssuerHash = types.StringValue(srv.IssuerHash)
	data.RawIssuer = types.StringValue(srv.RawIssuer)
	data.SerialNo = types.StringValue(srv.SerialNo)
	if srv.KeyID != nil {
		data.KeyID = types.StringValue(*srv.KeyID)
	} else {
		data.KeyID = types.StringNull()
	}
	if srv.IssuerKeyID != nil {
		data.IssuerKeyID = types.StringValue(*srv.IssuerKeyID)
	} else {
		data.IssuerKeyID = types.StringNull()
	}
	data.Text = types.StringValue(srv.Text)

	return &data, nil
}

func (r *InfinityCACertificateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	state := &InfinityCACertificateResourceModel{}

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
			"Error Reading Infinity CA certificate",
			fmt.Sprintf("Could not read Infinity CA certificate: %s", err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *InfinityCACertificateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	plan := &InfinityCACertificateResourceModel{}
	state := &InfinityCACertificateResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resourceID := int(state.ResourceID.ValueInt32())

	updateRequest := &config.CACertificateUpdateRequest{
		Certificate: plan.Certificate.ValueString(),
	}

	// Set boolean pointer field for update
	if !plan.TrustedIntermediate.IsNull() {
		trustedIntermediate := plan.TrustedIntermediate.ValueBool()
		updateRequest.TrustedIntermediate = &trustedIntermediate
	}

	_, err := r.InfinityClient.Config().UpdateCACertificate(ctx, resourceID, updateRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Infinity CA certificate",
			fmt.Sprintf("Could not update Infinity CA certificate with ID %d: %s", resourceID, err),
		)
		return
	}

	// Re-read the resource to get the latest state
	updatedModel, err := r.read(ctx, resourceID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Updated Infinity CA certificate",
			fmt.Sprintf("Could not read updated Infinity CA certificate with ID %d: %s", resourceID, err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, updatedModel)...)
}

func (r *InfinityCACertificateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	state := &InfinityCACertificateResourceModel{}

	tflog.Info(ctx, "Deleting Infinity CA certificate")

	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.InfinityClient.Config().DeleteCACertificate(ctx, int(state.ResourceID.ValueInt32()))

	// Ignore 404 Not Found and Lookup errors on delete
	if err != nil && !isNotFoundError(err) && !isLookupError(err) {
		resp.Diagnostics.AddError(
			"Error Deleting Infinity CA certificate",
			fmt.Sprintf("Could not delete Infinity CA certificate with ID %s: %s", state.ID.ValueString(), err),
		)
		return
	}
}

func (r *InfinityCACertificateResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resourceID, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Import ID must be a valid integer for the resource ID. Got: %s", req.ID),
		)
		return
	}

	tflog.Trace(ctx, fmt.Sprintf("Importing Infinity CA certificate with resource ID: %d", resourceID))

	// Read the resource from the API
	model, err := r.read(ctx, resourceID)
	if err != nil {
		// Check if the error is a 404 (not found)
		if isNotFoundError(err) {
			resp.Diagnostics.AddError(
				"Infinity CA Certificate Not Found",
				fmt.Sprintf("Infinity CA certificate with resource ID %d not found.", resourceID),
			)
			return
		}
		resp.Diagnostics.AddError(
			"Error Importing Infinity CA Certificate",
			fmt.Sprintf("Could not import Infinity CA certificate with resource ID %d: %s", resourceID, err),
		)
		return
	}

	// Set the state from the imported resource
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}
