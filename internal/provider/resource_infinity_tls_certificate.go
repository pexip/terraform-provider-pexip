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

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/pexip/go-infinity-sdk/v38/config"
)

var (
	_ resource.ResourceWithImportState = (*InfinityTLSCertificateResource)(nil)
)

type InfinityTLSCertificateResource struct {
	InfinityClient InfinityClient
}

type InfinityTLSCertificateResourceModel struct {
	ID                   types.String `tfsdk:"id"`
	ResourceID           types.Int32  `tfsdk:"resource_id"`
	Certificate          types.String `tfsdk:"certificate"`
	PrivateKey           types.String `tfsdk:"private_key"`
	PrivateKeyPassphrase types.String `tfsdk:"private_key_passphrase"`
	Parameters           types.String `tfsdk:"parameters"`
	Nodes                types.Set    `tfsdk:"nodes"`
	StartDate            types.String `tfsdk:"start_date"`
	EndDate              types.String `tfsdk:"end_date"`
	SubjectName          types.String `tfsdk:"subject_name"`
	SubjectHash          types.String `tfsdk:"subject_hash"`
	SubjectAltNames      types.String `tfsdk:"subject_alt_names"`
	RawSubject           types.String `tfsdk:"raw_subject"`
	IssuerName           types.String `tfsdk:"issuer_name"`
	IssuerHash           types.String `tfsdk:"issuer_hash"`
	RawIssuer            types.String `tfsdk:"raw_issuer"`
	SerialNo             types.String `tfsdk:"serial_no"`
	KeyID                types.String `tfsdk:"key_id"`
	IssuerKeyID          types.String `tfsdk:"issuer_key_id"`
	Text                 types.String `tfsdk:"text"`
}

func (r *InfinityTLSCertificateResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_infinity_tls_certificate"
}

func (r *InfinityTLSCertificateResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *InfinityTLSCertificateResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Resource URI for the TLS certificate in Infinity",
			},
			"resource_id": schema.Int32Attribute{
				Computed:            true,
				MarkdownDescription: "The resource integer identifier for the TLS certificate in Infinity",
			},
			"certificate": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
				MarkdownDescription: "The PEM-encoded certificate. This field is sensitive.",
			},
			"private_key": schema.StringAttribute{
				Required:  true,
				Sensitive: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
				MarkdownDescription: "The PEM-encoded private key. This field is sensitive.",
			},
			"private_key_passphrase": schema.StringAttribute{
				Optional:  true,
				Computed:  true,
				Sensitive: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(100),
				},
				MarkdownDescription: "The passphrase for the private key if it is encrypted. This field is sensitive. Maximum length: 100 characters.",
			},
			"parameters": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(1000),
				},
				MarkdownDescription: "Additional parameters for the certificate. Maximum length: 1000 characters.",
			},
			"nodes": schema.SetAttribute{
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
				Default:             setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{})),
				MarkdownDescription: "List of node resource URIs where this certificate should be deployed.",
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
				MarkdownDescription: "The subject name from the certificate.",
			},
			"subject_hash": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The subject hash from the certificate.",
			},
			"subject_alt_names": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The subject alternative names from the certificate.",
			},
			"raw_subject": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The raw subject data from the certificate.",
			},
			"issuer_name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The issuer name from the certificate.",
			},
			"issuer_hash": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The issuer hash from the certificate.",
			},
			"raw_issuer": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The raw issuer data from the certificate.",
			},
			"serial_no": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The serial number from the certificate.",
			},
			"key_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The key identifier from the certificate.",
			},
			"issuer_key_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The issuer key identifier from the certificate.",
			},
			"text": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The text representation of the certificate.",
			},
		},
		MarkdownDescription: "Manages a TLS certificate configuration with the Infinity service.",
	}
}

func (r *InfinityTLSCertificateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	plan := &InfinityTLSCertificateResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createRequest := &config.TLSCertificateCreateRequest{
		Certificate: plan.Certificate.ValueString(),
		PrivateKey:  plan.PrivateKey.ValueString(),
	}

	// Only set optional fields if they are not null in the plan
	if !plan.PrivateKeyPassphrase.IsNull() {
		createRequest.PrivateKeyPassphrase = plan.PrivateKeyPassphrase.ValueString()
	}
	if !plan.Parameters.IsNull() {
		createRequest.Parameters = plan.Parameters.ValueString()
	}
	if !plan.Nodes.IsNull() {
		var nodes []string
		resp.Diagnostics.Append(plan.Nodes.ElementsAs(ctx, &nodes, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		createRequest.Nodes = nodes
	}

	createResponse, err := r.InfinityClient.Config().CreateTLSCertificate(ctx, createRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Infinity TLS certificate",
			fmt.Sprintf("Could not create Infinity TLS certificate: %s", err),
		)
		return
	}

	resourceID, err := createResponse.ResourceID()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Retrieving Infinity TLS certificate ID",
			fmt.Sprintf("Could not retrieve ID for created Infinity TLS certificate: %s", err),
		)
		return
	}

	// Read the state from the API to get all computed values
	model, err := r.read(ctx, resourceID, plan.PrivateKey.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Created Infinity TLS certificate",
			fmt.Sprintf("Could not read created Infinity TLS certificate with ID %d: %s", resourceID, err),
		)
		return
	}
	tflog.Trace(ctx, fmt.Sprintf("created Infinity TLS certificate with ID: %s", model.ID))

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *InfinityTLSCertificateResource) read(ctx context.Context, resourceID int, privateKey string) (*InfinityTLSCertificateResourceModel, error) {
	var data InfinityTLSCertificateResourceModel

	srv, err := r.InfinityClient.Config().GetTLSCertificate(ctx, resourceID)
	if err != nil {
		return nil, err
	}

	if srv.ResourceURI == "" {
		return nil, fmt.Errorf("TLS certificate with ID %d not found", resourceID)
	}

	data.ID = types.StringValue(srv.ResourceURI)
	data.ResourceID = types.Int32Value(int32(resourceID)) // #nosec G115 -- API values are expected to be within int32 range
	data.Certificate = types.StringValue(srv.Certificate)
	data.PrivateKey = types.StringValue(privateKey) // The privateKey property of the TLS certificate is not returned by the API, so we need to set it manually
	data.PrivateKeyPassphrase = types.StringValue(srv.PrivateKeyPassphrase)
	data.Parameters = types.StringValue(srv.Parameters)
	data.StartDate = types.StringValue(srv.StartDate.String())
	data.EndDate = types.StringValue(srv.EndDate.String())
	data.SubjectName = types.StringValue(srv.SubjectName)
	data.SubjectHash = types.StringValue(srv.SubjectHash)
	data.SubjectAltNames = types.StringValue(srv.SubjectAltNames)
	data.RawSubject = types.StringValue(srv.RawSubject)
	data.IssuerName = types.StringValue(srv.IssuerName)
	data.IssuerHash = types.StringValue(srv.IssuerHash)
	data.RawIssuer = types.StringValue(srv.RawIssuer)
	data.SerialNo = types.StringValue(srv.SerialNo)
	data.Text = types.StringValue(srv.Text)

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

	// Convert nodes to types.Set
	nodesSetValue, diags := types.SetValueFrom(ctx, types.StringType, srv.Nodes)
	if diags.HasError() {
		return nil, fmt.Errorf("error converting nodes: %v", diags)
	}
	data.Nodes = nodesSetValue

	return &data, nil
}

func (r *InfinityTLSCertificateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	state := &InfinityTLSCertificateResourceModel{}

	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resourceID := int(state.ResourceID.ValueInt32())
	state, err := r.read(ctx, resourceID, state.PrivateKey.ValueString())
	if err != nil {
		// Check if the error is a 404 (not found)
		if isNotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading Infinity TLS certificate",
			fmt.Sprintf("Could not read Infinity TLS certificate: %s", err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *InfinityTLSCertificateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	plan := &InfinityTLSCertificateResourceModel{}
	state := &InfinityTLSCertificateResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resourceID := int(state.ResourceID.ValueInt32())

	updateRequest := &config.TLSCertificateUpdateRequest{
		Certificate: plan.Certificate.ValueString(),
		PrivateKey:  plan.PrivateKey.ValueString(),
	}

	if !plan.PrivateKeyPassphrase.IsNull() {
		updateRequest.PrivateKeyPassphrase = plan.PrivateKeyPassphrase.ValueString()
	}
	if !plan.Parameters.IsNull() {
		updateRequest.Parameters = plan.Parameters.ValueString()
	}
	if !plan.Nodes.IsNull() {
		var nodes []string
		resp.Diagnostics.Append(plan.Nodes.ElementsAs(ctx, &nodes, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		updateRequest.Nodes = nodes
	}

	_, err := r.InfinityClient.Config().UpdateTLSCertificate(ctx, resourceID, updateRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Infinity TLS certificate",
			fmt.Sprintf("Could not update Infinity TLS certificate with ID %d: %s", resourceID, err),
		)
		return
	}

	// Re-read the resource to get the latest state
	updatedModel, err := r.read(ctx, resourceID, plan.PrivateKey.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Updated Infinity TLS certificate",
			fmt.Sprintf("Could not read updated Infinity TLS certificate with ID %d: %s", resourceID, err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, updatedModel)...)
}

func (r *InfinityTLSCertificateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	state := &InfinityTLSCertificateResourceModel{}

	tflog.Info(ctx, "Deleting Infinity TLS certificate")

	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.InfinityClient.Config().DeleteTLSCertificate(ctx, int(state.ResourceID.ValueInt32()))

	// Ignore 404 Not Found and Lookup errors on delete
	if err != nil && !isNotFoundError(err) && !isLookupError(err) {
		resp.Diagnostics.AddError(
			"Error Deleting Infinity TLS certificate",
			fmt.Sprintf("Could not delete Infinity TLS certificate with ID %s: %s", state.ID.ValueString(), err),
		)
		return
	}
}

func (r *InfinityTLSCertificateResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resourceID, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Import ID must be a valid integer for the resource ID. Got: %s", req.ID),
		)
		return
	}

	tflog.Trace(ctx, fmt.Sprintf("Importing Infinity TLS certificate with resource ID: %d", resourceID))

	// Read the resource from the API
	model, err := r.read(ctx, resourceID, "")
	if err != nil {
		// Check if the error is a 404 (not found)
		if isNotFoundError(err) {
			resp.Diagnostics.AddError(
				"Infinity TLS Certificate Not Found",
				fmt.Sprintf("Infinity TLS certificate with resource ID %d not found.", resourceID),
			)
			return
		}
		resp.Diagnostics.AddError(
			"Error Importing Infinity TLS Certificate",
			fmt.Sprintf("Could not import Infinity TLS certificate with resource ID %d: %s", resourceID, err),
		)
		return
	}

	// Set the state from the imported resource
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}
