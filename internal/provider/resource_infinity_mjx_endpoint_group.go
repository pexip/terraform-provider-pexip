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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/pexip/go-infinity-sdk/v38/config"
)

var (
	_ resource.ResourceWithImportState = (*InfinityMjxEndpointGroupResource)(nil)
)

type InfinityMjxEndpointGroupResource struct {
	InfinityClient InfinityClient
}

type InfinityMjxEndpointGroupResourceModel struct {
	ID             types.String `tfsdk:"id"`
	ResourceID     types.Int32  `tfsdk:"resource_id"`
	Name           types.String `tfsdk:"name"`
	Description    types.String `tfsdk:"description"`
	SystemLocation types.String `tfsdk:"system_location"`
	MjxIntegration types.String `tfsdk:"mjx_integration"`
	DisableProxy   types.Bool   `tfsdk:"disable_proxy"`
	Endpoints      types.Set    `tfsdk:"endpoints"`
}

func (r *InfinityMjxEndpointGroupResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_infinity_mjx_endpoint_group"
}

func (r *InfinityMjxEndpointGroupResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *InfinityMjxEndpointGroupResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Resource URI for the MJX endpoint group in Infinity.",
			},
			"resource_id": schema.Int32Attribute{
				Computed:            true,
				MarkdownDescription: "The resource integer identifier for the MJX endpoint group in Infinity.",
			},
			"name": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					stringvalidator.LengthAtMost(250),
				},
				MarkdownDescription: "The name of this OTJ Endpoint Group. Maximum length: 250 characters.",
			},
			"description": schema.StringAttribute{
				Computed: true,
				Optional: true,
				Default:  stringdefault.StaticString(""),
				Validators: []validator.String{
					stringvalidator.LengthAtMost(250),
				},
				MarkdownDescription: "An optional description of this OTJ Endpoint Group. Maximum length: 250 characters.",
			},
			"system_location": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The system location of the Conferencing Nodes which will provide One-Touch Join services for this Endpoint Group.",
			},
			"mjx_integration": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The One-Touch Join Profile to which this Endpoint Group belongs.",
			},
			"disable_proxy": schema.BoolAttribute{
				Computed:            true,
				Optional:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Bypass the web proxy when sending requests to Cisco OBTP Endpoints in this OTJ Endpoint Group.",
			},
			"endpoints": schema.SetAttribute{
				Computed:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "The endpoints that belong to this Endpoint Group.",
			},
		},
		MarkdownDescription: "Manages an MJX endpoint group in Infinity. An MJX endpoint group defines the system location and One-Touch Join Profile for a set of OTJ endpoints.",
	}
}

func (r *InfinityMjxEndpointGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	plan := &InfinityMjxEndpointGroupResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	systemLocation := plan.SystemLocation.ValueString()
	createRequest := &config.MjxEndpointGroupCreateRequest{
		Name:           plan.Name.ValueString(),
		Description:    plan.Description.ValueString(),
		SystemLocation: &systemLocation,
		DisableProxy:   plan.DisableProxy.ValueBool(),
	}

	if !plan.MjxIntegration.IsNull() && !plan.MjxIntegration.IsUnknown() {
		v := plan.MjxIntegration.ValueString()
		createRequest.MjxIntegration = &v
	}

	createResponse, err := r.InfinityClient.Config().CreateMjxEndpointGroup(ctx, createRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Infinity MJX endpoint group",
			fmt.Sprintf("Could not create Infinity MJX endpoint group: %s", err),
		)
		return
	}

	resourceID, err := createResponse.ResourceID()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Retrieving Infinity MJX endpoint group ID",
			fmt.Sprintf("Could not retrieve ID for created Infinity MJX endpoint group: %s", err),
		)
		return
	}

	model, err := r.read(ctx, resourceID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Created Infinity MJX endpoint group",
			fmt.Sprintf("Could not read created Infinity MJX endpoint group with ID %d: %s", resourceID, err),
		)
		return
	}

	tflog.Trace(ctx, fmt.Sprintf("created Infinity MJX endpoint group with ID: %s, name: %s", model.ID, model.Name))

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *InfinityMjxEndpointGroupResource) read(ctx context.Context, resourceID int) (*InfinityMjxEndpointGroupResourceModel, error) {
	var data InfinityMjxEndpointGroupResourceModel

	srv, err := r.InfinityClient.Config().GetMjxEndpointGroup(ctx, resourceID)
	if err != nil {
		return nil, err
	}

	if srv.ResourceURI == "" {
		return nil, fmt.Errorf("MJX endpoint group with ID %d not found", resourceID)
	}

	data.ID = types.StringValue(srv.ResourceURI)
	data.ResourceID = types.Int32Value(int32(resourceID)) // #nosec G115 -- API values are expected to be within int32 range
	data.Name = types.StringValue(srv.Name)
	data.Description = types.StringValue(srv.Description)
	data.DisableProxy = types.BoolValue(srv.DisableProxy)

	if srv.SystemLocation != nil {
		data.SystemLocation = types.StringValue(*srv.SystemLocation)
	} else {
		data.SystemLocation = types.StringNull()
	}

	if srv.MjxIntegration != nil {
		data.MjxIntegration = types.StringValue(*srv.MjxIntegration)
	} else {
		data.MjxIntegration = types.StringNull()
	}

	if len(srv.Endpoints) > 0 {
		endpointURIs := make([]string, len(srv.Endpoints))
		for i, ep := range srv.Endpoints {
			endpointURIs[i] = ep.ResourceURI
		}
		endpoints, diags := types.SetValueFrom(ctx, types.StringType, endpointURIs)
		if diags.HasError() {
			return nil, fmt.Errorf("error converting endpoints: %v", diags)
		}
		data.Endpoints = endpoints
	} else {
		data.Endpoints = types.SetNull(types.StringType)
	}

	return &data, nil
}

func (r *InfinityMjxEndpointGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	state := &InfinityMjxEndpointGroupResourceModel{}

	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resourceID := int(state.ResourceID.ValueInt32())
	state, err := r.read(ctx, resourceID)
	if err != nil {
		if isNotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading Infinity MJX endpoint group",
			fmt.Sprintf("Could not read Infinity MJX endpoint group: %s", err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *InfinityMjxEndpointGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	plan := &InfinityMjxEndpointGroupResourceModel{}
	state := &InfinityMjxEndpointGroupResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	systemLocation := plan.SystemLocation.ValueString()
	updateRequest := &config.MjxEndpointGroupUpdateRequest{
		Name:           plan.Name.ValueString(),
		Description:    plan.Description.ValueString(),
		SystemLocation: &systemLocation,
		DisableProxy:   plan.DisableProxy.ValueBool(),
	}

	if !plan.MjxIntegration.IsNull() && !plan.MjxIntegration.IsUnknown() {
		v := plan.MjxIntegration.ValueString()
		updateRequest.MjxIntegration = &v
	}

	resourceID := int(state.ResourceID.ValueInt32())
	_, err := r.InfinityClient.Config().UpdateMjxEndpointGroup(ctx, resourceID, updateRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Infinity MJX endpoint group",
			fmt.Sprintf("Could not update Infinity MJX endpoint group: %s", err),
		)
		return
	}

	model, err := r.read(ctx, resourceID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Updated Infinity MJX endpoint group",
			fmt.Sprintf("Could not read updated Infinity MJX endpoint group with ID %d: %s", resourceID, err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *InfinityMjxEndpointGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	state := &InfinityMjxEndpointGroupResourceModel{}

	tflog.Info(ctx, "Deleting Infinity MJX endpoint group")

	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.InfinityClient.Config().DeleteMjxEndpointGroup(ctx, int(state.ResourceID.ValueInt32()))

	if err != nil && !isNotFoundError(err) && !isLookupError(err) {
		resp.Diagnostics.AddError(
			"Error Deleting Infinity MJX endpoint group",
			fmt.Sprintf("Could not delete Infinity MJX endpoint group with ID %s: %s", state.ID.ValueString(), err),
		)
		return
	}
}

func (r *InfinityMjxEndpointGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resourceID, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Import ID must be a valid integer for the resource ID. Got: %s", req.ID),
		)
		return
	}

	tflog.Trace(ctx, fmt.Sprintf("Importing Infinity MJX endpoint group with resource ID: %d", resourceID))

	model, err := r.read(ctx, resourceID)
	if err != nil {
		if isNotFoundError(err) {
			resp.Diagnostics.AddError(
				"Infinity MJX endpoint group Not Found",
				fmt.Sprintf("Infinity MJX endpoint group with resource ID %d not found.", resourceID),
			)
			return
		}
		resp.Diagnostics.AddError(
			"Error Importing Infinity MJX endpoint group",
			fmt.Sprintf("Could not import Infinity MJX endpoint group with resource ID %d: %s", resourceID, err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}
