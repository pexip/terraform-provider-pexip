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
	_ resource.ResourceWithImportState = (*InfinityLdapSyncFieldResource)(nil)
)

type InfinityLdapSyncFieldResource struct {
	InfinityClient InfinityClient
}

type InfinityLdapSyncFieldResourceModel struct {
	ID                   types.String `tfsdk:"id"`
	ResourceID           types.Int32  `tfsdk:"resource_id"`
	Name                 types.String `tfsdk:"name"`
	Description          types.String `tfsdk:"description"`
	TemplateVariableName types.String `tfsdk:"template_variable_name"`
	IsBinary             types.Bool   `tfsdk:"is_binary"`
}

func (r *InfinityLdapSyncFieldResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_infinity_ldap_sync_field"
}

func (r *InfinityLdapSyncFieldResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *InfinityLdapSyncFieldResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Resource URI for the LDAP sync field in Infinity",
			},
			"resource_id": schema.Int32Attribute{
				Computed:            true,
				MarkdownDescription: "The resource integer identifier for the LDAP sync field in Infinity",
			},
			"name": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					stringvalidator.LengthAtMost(100),
				},
				MarkdownDescription: "The name of the LDAP sync field. This should match the LDAP attribute name. Maximum length: 100 characters.",
			},
			"description": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(500),
				},
				MarkdownDescription: "Description of the LDAP sync field. Maximum length: 500 characters.",
			},
			"template_variable_name": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					stringvalidator.LengthAtMost(100),
				},
				MarkdownDescription: "The template variable name used to reference this field in synchronization templates. Maximum length: 100 characters.",
			},
			"is_binary": schema.BoolAttribute{
				Required:            true,
				MarkdownDescription: "Whether this LDAP field contains binary data (e.g., photos, certificates) or text data.",
			},
		},
		MarkdownDescription: "Manages an LDAP sync field with the Infinity service. LDAP sync fields define how specific LDAP directory attributes are mapped and synchronized with Pexip Infinity user and configuration data during LDAP directory synchronization operations.",
	}
}

func (r *InfinityLdapSyncFieldResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	plan := &InfinityLdapSyncFieldResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createRequest := &config.LdapSyncFieldCreateRequest{
		Name:                 plan.Name.ValueString(),
		Description:          plan.Description.ValueString(),
		TemplateVariableName: plan.TemplateVariableName.ValueString(),
		IsBinary:             plan.IsBinary.ValueBool(),
	}

	createResponse, err := r.InfinityClient.Config().CreateLdapSyncField(ctx, createRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Infinity LDAP sync field",
			fmt.Sprintf("Could not create Infinity LDAP sync field: %s", err),
		)
		return
	}

	resourceID, err := createResponse.ResourceID()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Retrieving Infinity LDAP sync field ID",
			fmt.Sprintf("Could not retrieve ID for created Infinity LDAP sync field: %s", err),
		)
		return
	}

	// Read the state from the API to get all computed values
	model, err := r.read(ctx, resourceID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Created Infinity LDAP sync field",
			fmt.Sprintf("Could not read created Infinity LDAP sync field with ID %d: %s", resourceID, err),
		)
		return
	}
	tflog.Trace(ctx, fmt.Sprintf("created Infinity LDAP sync field with ID: %s, name: %s", model.ID, model.Name))

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *InfinityLdapSyncFieldResource) read(ctx context.Context, resourceID int) (*InfinityLdapSyncFieldResourceModel, error) {
	var data InfinityLdapSyncFieldResourceModel

	srv, err := r.InfinityClient.Config().GetLdapSyncField(ctx, resourceID)
	if err != nil {
		return nil, err
	}

	if srv.ResourceURI == "" {
		return nil, fmt.Errorf("LDAP sync field with ID %d not found", resourceID)
	}

	data.ID = types.StringValue(srv.ResourceURI)
	data.ResourceID = types.Int32Value(int32(resourceID)) // #nosec G115 -- API values are expected to be within int32 range
	data.Name = types.StringValue(srv.Name)
	data.Description = types.StringValue(srv.Description)
	data.TemplateVariableName = types.StringValue(srv.TemplateVariableName)
	data.IsBinary = types.BoolValue(srv.IsBinary)

	return &data, nil
}

func (r *InfinityLdapSyncFieldResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	state := &InfinityLdapSyncFieldResourceModel{}

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
			"Error Reading Infinity LDAP sync field",
			fmt.Sprintf("Could not read Infinity LDAP sync field: %s", err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *InfinityLdapSyncFieldResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	plan := &InfinityLdapSyncFieldResourceModel{}
	state := &InfinityLdapSyncFieldResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateRequest := &config.LdapSyncFieldUpdateRequest{
		Name:                 plan.Name.ValueString(),
		Description:          plan.Description.ValueString(),
		TemplateVariableName: plan.TemplateVariableName.ValueString(),
	}

	// Handle optional pointer field for is_binary
	if !plan.IsBinary.IsNull() && !plan.IsBinary.IsUnknown() {
		isBinary := plan.IsBinary.ValueBool()
		updateRequest.IsBinary = &isBinary
	}

	resourceID := int(state.ResourceID.ValueInt32())
	_, err := r.InfinityClient.Config().UpdateLdapSyncField(ctx, resourceID, updateRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Infinity LDAP sync field",
			fmt.Sprintf("Could not update Infinity LDAP sync field: %s", err),
		)
		return
	}

	// Read the state from the API to get all computed values
	model, err := r.read(ctx, resourceID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Updated Infinity LDAP sync field",
			fmt.Sprintf("Could not read updated Infinity LDAP sync field with ID %d: %s", resourceID, err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *InfinityLdapSyncFieldResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	state := &InfinityLdapSyncFieldResourceModel{}

	tflog.Info(ctx, "Deleting Infinity LDAP sync field")

	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.InfinityClient.Config().DeleteLdapSyncField(ctx, int(state.ResourceID.ValueInt32()))

	// Ignore 404 Not Found and Lookup errors on delete
	if err != nil && !isNotFoundError(err) && !isLookupError(err) {
		resp.Diagnostics.AddError(
			"Error Deleting Infinity LDAP sync field",
			fmt.Sprintf("Could not delete Infinity LDAP sync field with ID %s: %s", state.ID.ValueString(), err),
		)
		return
	}
}

func (r *InfinityLdapSyncFieldResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resourceID, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Import ID must be a valid integer for the resource ID. Got: %s", req.ID),
		)
		return
	}

	tflog.Trace(ctx, fmt.Sprintf("Importing Infinity LDAP sync field with resource ID: %d", resourceID))

	// Read the resource from the API
	model, err := r.read(ctx, resourceID)
	if err != nil {
		// Check if the error is a 404 (not found)
		if isNotFoundError(err) {
			resp.Diagnostics.AddError(
				"Infinity LDAP Sync Field Not Found",
				fmt.Sprintf("Infinity LDAP sync field with resource ID %d not found.", resourceID),
			)
			return
		}
		resp.Diagnostics.AddError(
			"Error Importing Infinity LDAP Sync Field",
			fmt.Sprintf("Could not import Infinity LDAP sync field with resource ID %d: %s", resourceID, err),
		)
		return
	}

	// Set the state from the imported resource
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}
