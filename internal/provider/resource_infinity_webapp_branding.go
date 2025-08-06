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
	_ resource.ResourceWithImportState = (*InfinityWebappBrandingResource)(nil)
)

type InfinityWebappBrandingResource struct {
	InfinityClient InfinityClient
}

type InfinityWebappBrandingResourceModel struct {
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Description  types.String `tfsdk:"description"`
	UUID         types.String `tfsdk:"uuid"`
	WebappType   types.String `tfsdk:"webapp_type"`
	IsDefault    types.Bool   `tfsdk:"is_default"`
	BrandingFile types.String `tfsdk:"branding_file"`
	LastUpdated  types.String `tfsdk:"last_updated"`
}

func (r *InfinityWebappBrandingResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_infinity_webapp_branding"
}

func (r *InfinityWebappBrandingResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *InfinityWebappBrandingResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Resource URI for the webapp branding in Infinity",
			},
			"name": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					stringvalidator.LengthAtMost(100),
				},
				MarkdownDescription: "The name of the webapp branding configuration. This is used as the identifier. Maximum length: 100 characters.",
			},
			"description": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(500),
				},
				MarkdownDescription: "Description of the webapp branding configuration. Maximum length: 500 characters.",
			},
			"uuid": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
				MarkdownDescription: "The UUID for this branding configuration.",
			},
			"webapp_type": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf("pexapp", "management", "admin"),
				},
				MarkdownDescription: "The type of webapp this branding applies to. Valid values: pexapp, management, admin.",
			},
			"is_default": schema.BoolAttribute{
				Required:            true,
				MarkdownDescription: "Whether this is the default branding configuration for the webapp type.",
			},
			"branding_file": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
				MarkdownDescription: "The path or identifier for the branding file to use for customization.",
			},
			"last_updated": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Timestamp when this branding configuration was last updated.",
			},
		},
		MarkdownDescription: "Manages webapp branding configuration with the Infinity service. Webapp branding allows customization of the user interface for different Pexip web applications including the management interface, admin interface, and client applications.",
	}
}

func (r *InfinityWebappBrandingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	plan := &InfinityWebappBrandingResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createRequest := &config.WebappBrandingCreateRequest{
		Name:         plan.Name.ValueString(),
		Description:  plan.Description.ValueString(),
		UUID:         plan.UUID.ValueString(),
		WebappType:   plan.WebappType.ValueString(),
		IsDefault:    plan.IsDefault.ValueBool(),
		BrandingFile: plan.BrandingFile.ValueString(),
	}

	_, err := r.InfinityClient.Config().CreateWebappBranding(ctx, createRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Infinity webapp branding",
			fmt.Sprintf("Could not create Infinity webapp branding: %s", err),
		)
		return
	}

	// The resource uses name as the identifier
	brandingName := plan.Name.ValueString()

	// Read the state from the API to get all computed values
	model, err := r.read(ctx, brandingName)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Created Infinity webapp branding",
			fmt.Sprintf("Could not read created Infinity webapp branding with name %s: %s", brandingName, err),
		)
		return
	}
	tflog.Trace(ctx, fmt.Sprintf("created Infinity webapp branding with ID: %s, name: %s", model.ID, model.Name))

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *InfinityWebappBrandingResource) read(ctx context.Context, name string) (*InfinityWebappBrandingResourceModel, error) {
	var data InfinityWebappBrandingResourceModel

	srv, err := r.InfinityClient.Config().GetWebappBranding(ctx, name)
	if err != nil {
		return nil, err
	}

	if srv.ResourceURI == "" {
		return nil, fmt.Errorf("webapp branding with name %s not found", name)
	}

	data.ID = types.StringValue(srv.ResourceURI)
	data.Name = types.StringValue(srv.Name)
	data.Description = types.StringValue(srv.Description)
	data.UUID = types.StringValue(srv.UUID)
	data.WebappType = types.StringValue(srv.WebappType)
	data.IsDefault = types.BoolValue(srv.IsDefault)
	data.BrandingFile = types.StringValue(srv.BrandingFile)
	data.LastUpdated = types.StringValue(srv.LastUpdated.String())

	return &data, nil
}

func (r *InfinityWebappBrandingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	state := &InfinityWebappBrandingResourceModel{}

	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	name := state.Name.ValueString()
	state, err := r.read(ctx, name)
	if err != nil {
		// Check if the error is a 404 (not found)
		if isNotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading Infinity webapp branding",
			fmt.Sprintf("Could not read Infinity webapp branding: %s", err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *InfinityWebappBrandingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	plan := &InfinityWebappBrandingResourceModel{}
	state := &InfinityWebappBrandingResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateRequest := &config.WebappBrandingUpdateRequest{
		Name:         plan.Name.ValueString(),
		Description:  plan.Description.ValueString(),
		UUID:         plan.UUID.ValueString(),
		WebappType:   plan.WebappType.ValueString(),
		BrandingFile: plan.BrandingFile.ValueString(),
	}

	// Handle optional pointer field for is_default
	if !plan.IsDefault.IsNull() && !plan.IsDefault.IsUnknown() {
		isDefault := plan.IsDefault.ValueBool()
		updateRequest.IsDefault = &isDefault
	}

	name := state.Name.ValueString()
	_, err := r.InfinityClient.Config().UpdateWebappBranding(ctx, name, updateRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Infinity webapp branding",
			fmt.Sprintf("Could not update Infinity webapp branding: %s", err),
		)
		return
	}

	// Read the state from the API to get all computed values
	model, err := r.read(ctx, name)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Updated Infinity webapp branding",
			fmt.Sprintf("Could not read updated Infinity webapp branding with name %s: %s", name, err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *InfinityWebappBrandingResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	state := &InfinityWebappBrandingResourceModel{}

	tflog.Info(ctx, "Deleting Infinity webapp branding")

	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.InfinityClient.Config().DeleteWebappBranding(ctx, state.Name.ValueString())

	// Ignore 404 Not Found and Lookup errors on delete
	if err != nil && !isNotFoundError(err) && !isLookupError(err) {
		resp.Diagnostics.AddError(
			"Error Deleting Infinity webapp branding",
			fmt.Sprintf("Could not delete Infinity webapp branding with ID %s: %s", state.ID.ValueString(), err),
		)
		return
	}
}

func (r *InfinityWebappBrandingResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	name := req.ID

	tflog.Trace(ctx, fmt.Sprintf("Importing Infinity webapp branding with name: %s", name))

	// Read the resource from the API
	model, err := r.read(ctx, name)
	if err != nil {
		// Check if the error is a 404 (not found)
		if isNotFoundError(err) {
			resp.Diagnostics.AddError(
				"Infinity Webapp Branding Not Found",
				fmt.Sprintf("Infinity webapp branding with name %s not found.", name),
			)
			return
		}
		resp.Diagnostics.AddError(
			"Error Importing Infinity Webapp Branding",
			fmt.Sprintf("Could not import Infinity webapp branding with name %s: %s", name, err),
		)
		return
	}

	// Set the state from the imported resource
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}
