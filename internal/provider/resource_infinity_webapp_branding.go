/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package provider

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
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
	Name         types.String `tfsdk:"name"`
	Description  types.String `tfsdk:"description"`
	UUID         types.String `tfsdk:"uuid"`
	WebappType   types.String `tfsdk:"webapp_type"`
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
			"name": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					stringvalidator.LengthAtMost(100),
				},
				MarkdownDescription: "The name of the webapp branding configuration. This is used as the identifier. Maximum length: 100 characters.",
			},
			"description": schema.StringAttribute{
				Computed: true,
				Optional: true,
				Default: stringdefault.StaticString(""),
				Validators: []validator.String{
					stringvalidator.LengthAtMost(250),
				},
				MarkdownDescription: "Description of the webapp branding configuration. Maximum length: 250 characters.",
			},
			"uuid": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The UUID for this branding configuration. If not provided, a UUID will be automatically generated.",
			},
			"webapp_type": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf("webapp1", "webapp2", "webapp3"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				MarkdownDescription: "The type of webapp this branding applies to. Valid values: pexapp, management, admin.",
			},
			"branding_file": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
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

	// Generate UUID if not provided
	uuidValue := plan.UUID.ValueString()
	if plan.UUID.IsNull() || uuidValue == "" {
		uuidValue = uuid.New().String()
		plan.UUID = types.StringValue(uuidValue)
		tflog.Debug(ctx, fmt.Sprintf("Generated UUID for webapp branding: %s", uuidValue))
	}

	createRequest := &config.WebappBrandingCreateRequest{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
		WebappType:  plan.WebappType.ValueString(),
	}

	// Open the branding_file file
	brandingFilePath := plan.BrandingFile.ValueString()
	brandingFile, err := os.Open(brandingFilePath) // #nosec G304 -- File path provided by user in Terraform configuration
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Opening Branding File",
			fmt.Sprintf("Could not open package file at path '%s': %s", brandingFilePath, err),
		)
		return
	}
	defer func() { _ = brandingFile.Close() }()

	// Extract filename from the path
	filename := filepath.Base(brandingFilePath)

	createResponse, err := r.InfinityClient.Config().CreateWebappBranding(ctx, createRequest, filename, brandingFile)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Infinity webapp branding",
			fmt.Sprintf("Could not create Infinity webapp branding: %s", err),
		)
		return
	}

	resourceUUID, err := createResponse.ResUUID()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Retrieving Infinity webapp branding UUID",
			fmt.Sprintf("Could not retrieve UUID for created Infinity webapp branding: %s", err),
		)
		return
	}

	// Read the state from the API to get all computed values
	model, err := r.read(ctx, resourceUUID, plan.BrandingFile.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Created Infinity webapp branding",
			fmt.Sprintf("Could not read created Infinity webapp branding UUID '%s': %s", resourceUUID, err),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *InfinityWebappBrandingResource) read(ctx context.Context, uuid string, brandingFile string) (*InfinityWebappBrandingResourceModel, error) {
	var data InfinityWebappBrandingResourceModel

	srv, err := r.InfinityClient.Config().GetWebappBranding(ctx, uuid)
	if err != nil {
		return nil, err
	}

	if srv.ResourceURI == "" {
		return nil, fmt.Errorf("webapp branding with UUID '%s' not found", uuid)
	}

	data.Name = types.StringValue(srv.Name)
	data.Description = types.StringValue(srv.Description)
	data.UUID = types.StringValue(srv.UUID)
	data.WebappType = types.StringValue(srv.WebappType)
	data.LastUpdated = types.StringValue(srv.LastUpdated.String())
	data.BrandingFile = types.StringValue(brandingFile)

	return &data, nil
}

func (r *InfinityWebappBrandingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	state := &InfinityWebappBrandingResourceModel{}

	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	state, err := r.read(ctx, state.UUID.ValueString(), state.BrandingFile.ValueString())
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
	var plan, state InfinityWebappBrandingResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateRequest := &config.WebappBrandingUpdateRequest{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
		WebappType:  plan.WebappType.ValueString(),
	}

	_, err := r.InfinityClient.Config().UpdateWebappBranding(ctx, state.UUID.ValueString(), updateRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Infinity webapp branding",
			fmt.Sprintf("Could not update Infinity webapp branding UUID '%s': %s", state.UUID.ValueString(), err),
		)
		return
	}

	// Read the updated resource from the API
	model, err := r.read(ctx, state.UUID.ValueString(), plan.BrandingFile.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Updated Infinity webapp branding",
			fmt.Sprintf("Could not read updated Infinity webapp branding UUID '%s': %s", state.UUID.ValueString(), err),
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

	err := r.InfinityClient.Config().DeleteWebappBranding(ctx, state.UUID.ValueString())

	// Ignore 404 Not Found and Lookup errors on delete
	if err != nil && !isNotFoundError(err) && !isLookupError(err) {
		resp.Diagnostics.AddError(
			"Error Deleting Infinity webapp branding",
			fmt.Sprintf("Could not delete Infinity webapp branding with UUID %s: %s", state.UUID.ValueString(), err),
		)
		return
	}
}

func (r *InfinityWebappBrandingResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	uuid := req.ID

	// Read the resource from the API
	model, err := r.read(ctx, uuid, "")
	if err != nil {
		// Check if the error is a 404 (not found)
		if isNotFoundError(err) {
			resp.Diagnostics.AddError(
				"Infinity Webapp Branding Not Found",
				fmt.Sprintf("Infinity webapp branding with UUID %s not found.", uuid),
			)
			return
		}
		resp.Diagnostics.AddError(
			"Error Importing Infinity Webapp Branding",
			fmt.Sprintf("Could not import Infinity webapp branding with UUID %s: %s", uuid, err),
		)
		return
	}

	// Set the state from the imported resource
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}
