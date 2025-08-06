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
	_ resource.ResourceWithImportState = (*InfinityCertificateSigningRequestResource)(nil)
)

type InfinityCertificateSigningRequestResource struct {
	InfinityClient InfinityClient
}

type InfinityCertificateSigningRequestResourceModel struct {
	ID                        types.String `tfsdk:"id"`
	ResourceID                types.Int32  `tfsdk:"resource_id"`
	SubjectName               types.String `tfsdk:"subject_name"`
	DN                        types.String `tfsdk:"dn"`
	AdditionalSubjectAltNames types.String `tfsdk:"additional_subject_alt_names"`
	PrivateKeyType            types.String `tfsdk:"private_key_type"`
	PrivateKey                types.String `tfsdk:"private_key"`
	PrivateKeyPassphrase      types.String `tfsdk:"private_key_passphrase"`
	AdCompatible              types.Bool   `tfsdk:"ad_compatible"`
	TLSCertificate            types.String `tfsdk:"tls_certificate"`
	CSR                       types.String `tfsdk:"csr"`
	Certificate               types.String `tfsdk:"certificate"`
}

func (r *InfinityCertificateSigningRequestResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_infinity_certificate_signing_request"
}

func (r *InfinityCertificateSigningRequestResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *InfinityCertificateSigningRequestResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Resource URI for the certificate signing request in Infinity",
			},
			"resource_id": schema.Int32Attribute{
				Computed:            true,
				MarkdownDescription: "The resource integer identifier for the certificate signing request in Infinity",
			},
			"subject_name": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(250),
				},
				MarkdownDescription: "The subject name for the certificate. Maximum length: 250 characters.",
			},
			"dn": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(500),
				},
				MarkdownDescription: "The distinguished name for the certificate. Maximum length: 500 characters.",
			},
			"additional_subject_alt_names": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(500),
				},
				MarkdownDescription: "Additional subject alternative names for the certificate. Maximum length: 500 characters.",
			},
			"private_key_type": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf("rsa2048", "rsa4096", "ecdsa256", "ecdsa384"),
				},
				MarkdownDescription: "The private key type. Valid choices: rsa2048, rsa4096, ecdsa256, ecdsa384.",
			},
			"private_key": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Sensitive:           true,
				MarkdownDescription: "The private key content (PEM format). Will be generated if not provided.",
			},
			"private_key_passphrase": schema.StringAttribute{
				Optional:  true,
				Computed:  true,
				Sensitive: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(250),
				},
				MarkdownDescription: "The passphrase for the private key. Maximum length: 250 characters.",
			},
			"ad_compatible": schema.BoolAttribute{
				Required:            true,
				MarkdownDescription: "Whether the certificate should be Active Directory compatible.",
			},
			"tls_certificate": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Reference to the TLS certificate resource URI.",
			},
			"csr": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The generated certificate signing request (PEM format).",
			},
			"certificate": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
				MarkdownDescription: "The signed certificate content (PEM format).",
			},
		},
		MarkdownDescription: "Manages a certificate signing request configuration with the Infinity service.",
	}
}

func (r *InfinityCertificateSigningRequestResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	plan := &InfinityCertificateSigningRequestResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createRequest := &config.CertificateSigningRequestCreateRequest{
		SubjectName:    plan.SubjectName.ValueString(),
		PrivateKeyType: plan.PrivateKeyType.ValueString(),
		AdCompatible:   plan.AdCompatible.ValueBool(),
	}

	// Set optional fields
	if !plan.DN.IsNull() {
		createRequest.DN = plan.DN.ValueString()
	}
	if !plan.AdditionalSubjectAltNames.IsNull() {
		createRequest.AdditionalSubjectAltNames = plan.AdditionalSubjectAltNames.ValueString()
	}
	if !plan.PrivateKey.IsNull() {
		privateKey := plan.PrivateKey.ValueString()
		createRequest.PrivateKey = &privateKey
	}
	if !plan.PrivateKeyPassphrase.IsNull() {
		createRequest.PrivateKeyPassphrase = plan.PrivateKeyPassphrase.ValueString()
	}
	if !plan.TLSCertificate.IsNull() {
		tlsCertificate := plan.TLSCertificate.ValueString()
		createRequest.TLSCertificate = &tlsCertificate
	}

	createResponse, err := r.InfinityClient.Config().CreateCertificateSigningRequest(ctx, createRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Infinity certificate signing request",
			fmt.Sprintf("Could not create Infinity certificate signing request: %s", err),
		)
		return
	}

	resourceID, err := createResponse.ResourceID()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Retrieving Infinity certificate signing request ID",
			fmt.Sprintf("Could not retrieve ID for created Infinity certificate signing request: %s", err),
		)
		return
	}

	// Read the state from the API to get all computed values
	model, err := r.read(ctx, resourceID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Created Infinity certificate signing request",
			fmt.Sprintf("Could not read created Infinity certificate signing request with ID %d: %s", resourceID, err),
		)
		return
	}
	tflog.Trace(ctx, fmt.Sprintf("created Infinity certificate signing request with ID: %s, subject: %s", model.ID, model.SubjectName))

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *InfinityCertificateSigningRequestResource) read(ctx context.Context, resourceID int) (*InfinityCertificateSigningRequestResourceModel, error) {
	var data InfinityCertificateSigningRequestResourceModel

	srv, err := r.InfinityClient.Config().GetCertificateSigningRequest(ctx, resourceID)
	if err != nil {
		return nil, err
	}

	if srv.ResourceURI == "" {
		return nil, fmt.Errorf("certificate signing request with ID %d not found", resourceID)
	}

	data.ID = types.StringValue(srv.ResourceURI)
	data.ResourceID = types.Int32Value(int32(resourceID)) // #nosec G115 -- API values are expected to be within int32 range
	data.SubjectName = types.StringValue(srv.SubjectName)
	data.DN = types.StringValue(srv.DN)
	data.AdditionalSubjectAltNames = types.StringValue(srv.AdditionalSubjectAltNames)
	data.PrivateKeyType = types.StringValue(srv.PrivateKeyType)
	if srv.PrivateKey != nil {
		data.PrivateKey = types.StringValue(*srv.PrivateKey)
	} else {
		data.PrivateKey = types.StringNull()
	}
	data.PrivateKeyPassphrase = types.StringValue(srv.PrivateKeyPassphrase)
	data.AdCompatible = types.BoolValue(srv.AdCompatible)
	if srv.TLSCertificate != nil {
		data.TLSCertificate = types.StringValue(*srv.TLSCertificate)
	} else {
		data.TLSCertificate = types.StringNull()
	}
	data.CSR = types.StringValue(srv.CSR)
	data.Certificate = types.StringValue(srv.Certificate)

	return &data, nil
}

func (r *InfinityCertificateSigningRequestResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	state := &InfinityCertificateSigningRequestResourceModel{}

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
			"Error Reading Infinity certificate signing request",
			fmt.Sprintf("Could not read Infinity certificate signing request: %s", err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *InfinityCertificateSigningRequestResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	plan := &InfinityCertificateSigningRequestResourceModel{}
	state := &InfinityCertificateSigningRequestResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resourceID := int(state.ResourceID.ValueInt32())

	updateRequest := &config.CertificateSigningRequestUpdateRequest{
		SubjectName:    plan.SubjectName.ValueString(),
		PrivateKeyType: plan.PrivateKeyType.ValueString(),
	}

	// Set boolean pointer field for update
	if !plan.AdCompatible.IsNull() {
		adCompatible := plan.AdCompatible.ValueBool()
		updateRequest.AdCompatible = &adCompatible
	}

	// Set optional fields
	if !plan.DN.IsNull() {
		updateRequest.DN = plan.DN.ValueString()
	}
	if !plan.AdditionalSubjectAltNames.IsNull() {
		updateRequest.AdditionalSubjectAltNames = plan.AdditionalSubjectAltNames.ValueString()
	}
	if !plan.PrivateKey.IsNull() {
		privateKey := plan.PrivateKey.ValueString()
		updateRequest.PrivateKey = &privateKey
	}
	if !plan.PrivateKeyPassphrase.IsNull() {
		updateRequest.PrivateKeyPassphrase = plan.PrivateKeyPassphrase.ValueString()
	}
	if !plan.TLSCertificate.IsNull() {
		tlsCertificate := plan.TLSCertificate.ValueString()
		updateRequest.TLSCertificate = &tlsCertificate
	}
	if !plan.Certificate.IsNull() {
		updateRequest.Certificate = plan.Certificate.ValueString()
	}

	_, err := r.InfinityClient.Config().UpdateCertificateSigningRequest(ctx, resourceID, updateRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Infinity certificate signing request",
			fmt.Sprintf("Could not update Infinity certificate signing request with ID %d: %s", resourceID, err),
		)
		return
	}

	// Re-read the resource to get the latest state
	updatedModel, err := r.read(ctx, resourceID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Updated Infinity certificate signing request",
			fmt.Sprintf("Could not read updated Infinity certificate signing request with ID %d: %s", resourceID, err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, updatedModel)...)
}

func (r *InfinityCertificateSigningRequestResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	state := &InfinityCertificateSigningRequestResourceModel{}

	tflog.Info(ctx, "Deleting Infinity certificate signing request")

	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.InfinityClient.Config().DeleteCertificateSigningRequest(ctx, int(state.ResourceID.ValueInt32()))

	// Ignore 404 Not Found and Lookup errors on delete
	if err != nil && !isNotFoundError(err) && !isLookupError(err) {
		resp.Diagnostics.AddError(
			"Error Deleting Infinity certificate signing request",
			fmt.Sprintf("Could not delete Infinity certificate signing request with ID %s: %s", state.ID.ValueString(), err),
		)
		return
	}
}

func (r *InfinityCertificateSigningRequestResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resourceID, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Import ID must be a valid integer for the resource ID. Got: %s", req.ID),
		)
		return
	}

	tflog.Trace(ctx, fmt.Sprintf("Importing Infinity certificate signing request with resource ID: %d", resourceID))

	// Read the resource from the API
	model, err := r.read(ctx, resourceID)
	if err != nil {
		// Check if the error is a 404 (not found)
		if isNotFoundError(err) {
			resp.Diagnostics.AddError(
				"Infinity Certificate Signing Request Not Found",
				fmt.Sprintf("Infinity certificate signing request with resource ID %d not found.", resourceID),
			)
			return
		}
		resp.Diagnostics.AddError(
			"Error Importing Infinity Certificate Signing Request",
			fmt.Sprintf("Could not import Infinity certificate signing request with resource ID %d: %s", resourceID, err),
		)
		return
	}

	// Set the state from the imported resource
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}
