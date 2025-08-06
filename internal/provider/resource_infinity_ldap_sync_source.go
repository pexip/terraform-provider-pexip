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
	_ resource.ResourceWithImportState = (*InfinityLdapSyncSourceResource)(nil)
)

type InfinityLdapSyncSourceResource struct {
	InfinityClient InfinityClient
}

type InfinityLdapSyncSourceResourceModel struct {
	ID                   types.String `tfsdk:"id"`
	ResourceID           types.Int32  `tfsdk:"resource_id"`
	Name                 types.String `tfsdk:"name"`
	Description          types.String `tfsdk:"description"`
	LdapServer           types.String `tfsdk:"ldap_server"`
	LdapBaseDN           types.String `tfsdk:"ldap_base_dn"`
	LdapBindUsername     types.String `tfsdk:"ldap_bind_username"`
	LdapBindPassword     types.String `tfsdk:"ldap_bind_password"`
	LdapUseGlobalCatalog types.Bool   `tfsdk:"ldap_use_global_catalog"`
	LdapPermitNoTLS      types.Bool   `tfsdk:"ldap_permit_no_tls"`
}

func (r *InfinityLdapSyncSourceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_infinity_ldap_sync_source"
}

func (r *InfinityLdapSyncSourceResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *InfinityLdapSyncSourceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Resource URI for the LDAP sync source in Infinity",
			},
			"resource_id": schema.Int32Attribute{
				Computed:            true,
				MarkdownDescription: "The resource integer identifier for the LDAP sync source in Infinity",
			},
			"name": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(250),
				},
				MarkdownDescription: "The unique name of the LDAP synchronization source. Maximum length: 250 characters.",
			},
			"description": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(250),
				},
				MarkdownDescription: "A description of the LDAP synchronization source. Maximum length: 250 characters.",
			},
			"ldap_server": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(255),
				},
				MarkdownDescription: "The hostname of the LDAP server. Enter a domain name for DNS SRV lookup or an FQDN for DNS A/AAAA lookup. Maximum length: 255 characters.",
			},
			"ldap_base_dn": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(255),
				},
				MarkdownDescription: "The base DN of the LDAP forest to query (e.g. dc=example,dc=com). Maximum length: 255 characters.",
			},
			"ldap_bind_username": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(255),
				},
				MarkdownDescription: "The username used to bind to the LDAP server. This should be a domain user service account. Maximum length: 255 characters.",
			},
			"ldap_bind_password": schema.StringAttribute{
				Required:  true,
				Sensitive: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(100),
				},
				MarkdownDescription: "The password used to bind to the LDAP server. Maximum length: 100 characters.",
			},
			"ldap_use_global_catalog": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Search the Active Directory Global Catalog instead of traditional LDAP. Defaults to false.",
			},
			"ldap_permit_no_tls": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Permit LDAP queries to be sent over an insecure connection. Defaults to false.",
			},
		},
		MarkdownDescription: "Manages an LDAP synchronization source configuration with the Infinity service.",
	}
}

func (r *InfinityLdapSyncSourceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	plan := &InfinityLdapSyncSourceResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createRequest := &config.LdapSyncSourceCreateRequest{
		Name:             plan.Name.ValueString(),
		LdapServer:       plan.LdapServer.ValueString(),
		LdapBaseDN:       plan.LdapBaseDN.ValueString(),
		LdapBindUsername: plan.LdapBindUsername.ValueString(),
		LdapBindPassword: plan.LdapBindPassword.ValueString(),
	}

	// Set optional fields
	if !plan.Description.IsNull() {
		createRequest.Description = plan.Description.ValueString()
	}
	if !plan.LdapUseGlobalCatalog.IsNull() {
		createRequest.LdapUseGlobalCatalog = plan.LdapUseGlobalCatalog.ValueBool()
	}
	if !plan.LdapPermitNoTLS.IsNull() {
		createRequest.LdapPermitNoTLS = plan.LdapPermitNoTLS.ValueBool()
	}

	createResponse, err := r.InfinityClient.Config().CreateLdapSyncSource(ctx, createRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Infinity LDAP sync source",
			fmt.Sprintf("Could not create Infinity LDAP sync source: %s", err),
		)
		return
	}

	resourceID, err := createResponse.ResourceID()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Retrieving Infinity LDAP sync source ID",
			fmt.Sprintf("Could not retrieve ID for created Infinity LDAP sync source: %s", err),
		)
		return
	}

	// Read the state from the API to get all computed values
	model, err := r.read(ctx, resourceID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Created Infinity LDAP sync source",
			fmt.Sprintf("Could not read created Infinity LDAP sync source with ID %d: %s", resourceID, err),
		)
		return
	}
	tflog.Trace(ctx, fmt.Sprintf("created Infinity LDAP sync source with ID: %s, name: %s", model.ID, model.Name))

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *InfinityLdapSyncSourceResource) read(ctx context.Context, resourceID int) (*InfinityLdapSyncSourceResourceModel, error) {
	var data InfinityLdapSyncSourceResourceModel

	srv, err := r.InfinityClient.Config().GetLdapSyncSource(ctx, resourceID)
	if err != nil {
		return nil, err
	}

	if srv.ResourceURI == "" {
		return nil, fmt.Errorf("LDAP sync source with ID %d not found", resourceID)
	}

	data.ID = types.StringValue(srv.ResourceURI)
	data.ResourceID = types.Int32Value(int32(resourceID)) // #nosec G115 -- API values are expected to be within int32 range
	data.Name = types.StringValue(srv.Name)
	data.Description = types.StringValue(srv.Description)
	data.LdapServer = types.StringValue(srv.LdapServer)
	data.LdapBaseDN = types.StringValue(srv.LdapBaseDN)
	data.LdapBindUsername = types.StringValue(srv.LdapBindUsername)
	data.LdapBindPassword = types.StringValue(srv.LdapBindPassword)

	// Set boolean fields
	data.LdapUseGlobalCatalog = types.BoolValue(srv.LdapUseGlobalCatalog)
	data.LdapPermitNoTLS = types.BoolValue(srv.LdapPermitNoTLS)

	return &data, nil
}

func (r *InfinityLdapSyncSourceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	state := &InfinityLdapSyncSourceResourceModel{}

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
			"Error Reading Infinity LDAP sync source",
			fmt.Sprintf("Could not read Infinity LDAP sync source: %s", err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *InfinityLdapSyncSourceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	plan := &InfinityLdapSyncSourceResourceModel{}
	state := &InfinityLdapSyncSourceResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resourceID := int(state.ResourceID.ValueInt32())

	updateRequest := &config.LdapSyncSourceUpdateRequest{
		Name:             plan.Name.ValueString(),
		LdapServer:       plan.LdapServer.ValueString(),
		LdapBaseDN:       plan.LdapBaseDN.ValueString(),
		LdapBindUsername: plan.LdapBindUsername.ValueString(),
		LdapBindPassword: plan.LdapBindPassword.ValueString(),
	}

	if !plan.Description.IsNull() {
		updateRequest.Description = plan.Description.ValueString()
	}
	if !plan.LdapUseGlobalCatalog.IsNull() {
		useGlobalCatalog := plan.LdapUseGlobalCatalog.ValueBool()
		updateRequest.LdapUseGlobalCatalog = &useGlobalCatalog
	}
	if !plan.LdapPermitNoTLS.IsNull() {
		permitNoTLS := plan.LdapPermitNoTLS.ValueBool()
		updateRequest.LdapPermitNoTLS = &permitNoTLS
	}

	_, err := r.InfinityClient.Config().UpdateLdapSyncSource(ctx, resourceID, updateRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Infinity LDAP sync source",
			fmt.Sprintf("Could not update Infinity LDAP sync source with ID %d: %s", resourceID, err),
		)
		return
	}

	// Re-read the resource to get the latest state
	updatedModel, err := r.read(ctx, resourceID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Updated Infinity LDAP sync source",
			fmt.Sprintf("Could not read updated Infinity LDAP sync source with ID %d: %s", resourceID, err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, updatedModel)...)
}

func (r *InfinityLdapSyncSourceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	state := &InfinityLdapSyncSourceResourceModel{}

	tflog.Info(ctx, "Deleting Infinity LDAP sync source")

	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.InfinityClient.Config().DeleteLdapSyncSource(ctx, int(state.ResourceID.ValueInt32()))

	// Ignore 404 Not Found and Lookup errors on delete
	if err != nil && !isNotFoundError(err) && !isLookupError(err) {
		resp.Diagnostics.AddError(
			"Error Deleting Infinity LDAP sync source",
			fmt.Sprintf("Could not delete Infinity LDAP sync source with ID %s: %s", state.ID.ValueString(), err),
		)
		return
	}
}

func (r *InfinityLdapSyncSourceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resourceID, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Import ID must be a valid integer for the resource ID. Got: %s", req.ID),
		)
		return
	}

	tflog.Trace(ctx, fmt.Sprintf("Importing Infinity LDAP sync source with resource ID: %d", resourceID))

	// Read the resource from the API
	model, err := r.read(ctx, resourceID)
	if err != nil {
		// Check if the error is a 404 (not found)
		if isNotFoundError(err) {
			resp.Diagnostics.AddError(
				"Infinity LDAP Sync Source Not Found",
				fmt.Sprintf("Infinity LDAP sync source with resource ID %d not found.", resourceID),
			)
			return
		}
		resp.Diagnostics.AddError(
			"Error Importing Infinity LDAP Sync Source",
			fmt.Sprintf("Could not import Infinity LDAP sync source with resource ID %d: %s", resourceID, err),
		)
		return
	}

	// Set the state from the imported resource
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}
