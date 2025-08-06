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
	_ resource.ResourceWithImportState = (*InfinityWebappAliasResource)(nil)
)

type InfinityWebappAliasResource struct {
	InfinityClient InfinityClient
}

type InfinityWebappAliasResourceModel struct {
	ID          types.String `tfsdk:"id"`
	ResourceID  types.Int32  `tfsdk:"resource_id"`
	Slug        types.String `tfsdk:"slug"`
	Description types.String `tfsdk:"description"`
	WebappType  types.String `tfsdk:"webapp_type"`
	IsEnabled   types.Bool   `tfsdk:"is_enabled"`
	Bundle      types.String `tfsdk:"bundle"`
	Branding    types.String `tfsdk:"branding"`
}

func (r *InfinityWebappAliasResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_infinity_webapp_alias"
}

func (r *InfinityWebappAliasResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *InfinityWebappAliasResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Resource URI for the webapp alias in Infinity",
			},
			"resource_id": schema.Int32Attribute{
				Computed:            true,
				MarkdownDescription: "The resource integer identifier for the webapp alias in Infinity",
			},
			"slug": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					stringvalidator.LengthAtMost(100),
				},
				MarkdownDescription: "The slug (URL path component) for this webapp alias. Maximum length: 100 characters.",
			},
			"description": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(500),
				},
				MarkdownDescription: "Description of the webapp alias. Maximum length: 500 characters.",
			},
			"webapp_type": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf("pexapp", "management", "admin"),
				},
				MarkdownDescription: "The type of webapp this alias serves. Valid values: pexapp, management, admin.",
			},
			"is_enabled": schema.BoolAttribute{
				Required:            true,
				MarkdownDescription: "Whether this webapp alias is enabled and active.",
			},
			"bundle": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(200),
				},
				MarkdownDescription: "The bundle URI associated with this webapp alias. Maximum length: 200 characters.",
			},
			"branding": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(200),
				},
				MarkdownDescription: "The branding URI associated with this webapp alias. Maximum length: 200 characters.",
			},
		},
		MarkdownDescription: "Manages a webapp alias with the Infinity service. Webapp aliases provide alternative URL paths to access different web applications within Pexip Infinity, allowing for customized branding and user experiences.",
	}
}

func (r *InfinityWebappAliasResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	plan := &InfinityWebappAliasResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createRequest := &config.WebappAliasCreateRequest{
		Slug:        plan.Slug.ValueString(),
		Description: plan.Description.ValueString(),
		WebappType:  plan.WebappType.ValueString(),
		IsEnabled:   plan.IsEnabled.ValueBool(),
	}

	// Handle optional pointer fields
	if !plan.Bundle.IsNull() && !plan.Bundle.IsUnknown() {
		bundle := plan.Bundle.ValueString()
		createRequest.Bundle = &bundle
	}

	if !plan.Branding.IsNull() && !plan.Branding.IsUnknown() {
		branding := plan.Branding.ValueString()
		createRequest.Branding = &branding
	}

	createResponse, err := r.InfinityClient.Config().CreateWebappAlias(ctx, createRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Infinity webapp alias",
			fmt.Sprintf("Could not create Infinity webapp alias: %s", err),
		)
		return
	}

	resourceID, err := createResponse.ResourceID()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Retrieving Infinity webapp alias ID",
			fmt.Sprintf("Could not retrieve ID for created Infinity webapp alias: %s", err),
		)
		return
	}

	// Read the state from the API to get all computed values
	model, err := r.read(ctx, resourceID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Created Infinity webapp alias",
			fmt.Sprintf("Could not read created Infinity webapp alias with ID %d: %s", resourceID, err),
		)
		return
	}
	tflog.Trace(ctx, fmt.Sprintf("created Infinity webapp alias with ID: %s, slug: %s", model.ID, model.Slug))

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *InfinityWebappAliasResource) read(ctx context.Context, resourceID int) (*InfinityWebappAliasResourceModel, error) {
	var data InfinityWebappAliasResourceModel

	srv, err := r.InfinityClient.Config().GetWebappAlias(ctx, resourceID)
	if err != nil {
		return nil, err
	}

	if srv.ResourceURI == "" {
		return nil, fmt.Errorf("webapp alias with ID %d not found", resourceID)
	}

	data.ID = types.StringValue(srv.ResourceURI)
	data.ResourceID = types.Int32Value(int32(resourceID)) // #nosec G115 -- API values are expected to be within int32 range
	data.Slug = types.StringValue(srv.Slug)
	data.Description = types.StringValue(srv.Description)
	data.WebappType = types.StringValue(srv.WebappType)
	data.IsEnabled = types.BoolValue(srv.IsEnabled)

	// Handle optional pointer fields
	if srv.Bundle != nil {
		data.Bundle = types.StringValue(*srv.Bundle)
	} else {
		data.Bundle = types.StringNull()
	}

	if srv.Branding != nil {
		data.Branding = types.StringValue(*srv.Branding)
	} else {
		data.Branding = types.StringNull()
	}

	return &data, nil
}

func (r *InfinityWebappAliasResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	state := &InfinityWebappAliasResourceModel{}

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
			"Error Reading Infinity webapp alias",
			fmt.Sprintf("Could not read Infinity webapp alias: %s", err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *InfinityWebappAliasResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	plan := &InfinityWebappAliasResourceModel{}
	state := &InfinityWebappAliasResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateRequest := &config.WebappAliasUpdateRequest{
		Slug:        plan.Slug.ValueString(),
		Description: plan.Description.ValueString(),
		WebappType:  plan.WebappType.ValueString(),
	}

	// Handle optional pointer fields
	if !plan.IsEnabled.IsNull() && !plan.IsEnabled.IsUnknown() {
		isEnabled := plan.IsEnabled.ValueBool()
		updateRequest.IsEnabled = &isEnabled
	}

	if !plan.Bundle.IsNull() && !plan.Bundle.IsUnknown() {
		bundle := plan.Bundle.ValueString()
		updateRequest.Bundle = &bundle
	}

	if !plan.Branding.IsNull() && !plan.Branding.IsUnknown() {
		branding := plan.Branding.ValueString()
		updateRequest.Branding = &branding
	}

	resourceID := int(state.ResourceID.ValueInt32())
	_, err := r.InfinityClient.Config().UpdateWebappAlias(ctx, resourceID, updateRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Infinity webapp alias",
			fmt.Sprintf("Could not update Infinity webapp alias: %s", err),
		)
		return
	}

	// Read the state from the API to get all computed values
	model, err := r.read(ctx, resourceID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Updated Infinity webapp alias",
			fmt.Sprintf("Could not read updated Infinity webapp alias with ID %d: %s", resourceID, err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *InfinityWebappAliasResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	state := &InfinityWebappAliasResourceModel{}

	tflog.Info(ctx, "Deleting Infinity webapp alias")

	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.InfinityClient.Config().DeleteWebappAlias(ctx, int(state.ResourceID.ValueInt32()))

	// Ignore 404 Not Found and Lookup errors on delete
	if err != nil && !isNotFoundError(err) && !isLookupError(err) {
		resp.Diagnostics.AddError(
			"Error Deleting Infinity webapp alias",
			fmt.Sprintf("Could not delete Infinity webapp alias with ID %s: %s", state.ID.ValueString(), err),
		)
		return
	}
}

func (r *InfinityWebappAliasResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resourceID, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Import ID must be a valid integer for the resource ID. Got: %s", req.ID),
		)
		return
	}

	tflog.Trace(ctx, fmt.Sprintf("Importing Infinity webapp alias with resource ID: %d", resourceID))

	// Read the resource from the API
	model, err := r.read(ctx, resourceID)
	if err != nil {
		// Check if the error is a 404 (not found)
		if isNotFoundError(err) {
			resp.Diagnostics.AddError(
				"Infinity Webapp Alias Not Found",
				fmt.Sprintf("Infinity webapp alias with resource ID %d not found.", resourceID),
			)
			return
		}
		resp.Diagnostics.AddError(
			"Error Importing Infinity Webapp Alias",
			fmt.Sprintf("Could not import Infinity webapp alias with resource ID %d: %s", resourceID, err),
		)
		return
	}

	// Set the state from the imported resource
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}
