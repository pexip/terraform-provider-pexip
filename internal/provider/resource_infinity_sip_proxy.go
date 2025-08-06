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

	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/pexip/go-infinity-sdk/v38/config"
)

var (
	_ resource.ResourceWithImportState = (*InfinitySIPProxyResource)(nil)
)

type InfinitySIPProxyResource struct {
	InfinityClient InfinityClient
}

type InfinitySIPProxyResourceModel struct {
	ID          types.String `tfsdk:"id"`
	ResourceID  types.Int32  `tfsdk:"resource_id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Address     types.String `tfsdk:"address"`
	Port        types.Int32  `tfsdk:"port"`
	Transport   types.String `tfsdk:"transport"`
}

func (r *InfinitySIPProxyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_infinity_sip_proxy"
}

func (r *InfinitySIPProxyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *InfinitySIPProxyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Resource URI for the SIP proxy in Infinity",
			},
			"resource_id": schema.Int32Attribute{
				Computed:            true,
				MarkdownDescription: "The resource integer identifier for the SIP proxy in Infinity",
			},
			"name": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(250),
				},
				MarkdownDescription: "The name used to refer to this SIP proxy. Maximum length: 250 characters.",
			},
			"description": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(250),
				},
				MarkdownDescription: "A description of the SIP proxy. Maximum length: 250 characters.",
			},
			"address": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(255),
				},
				MarkdownDescription: "The address or hostname of the SIP proxy. Maximum length: 255 characters.",
			},
			"port": schema.Int32Attribute{
				Optional: true,
				Computed: true,
				Validators: []validator.Int32{
					int32validator.Between(1, 65535),
				},
				MarkdownDescription: "The port number for the SIP proxy. Range: 1 to 65535.",
			},
			"transport": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf("tcp", "udp", "tls"),
				},
				MarkdownDescription: "The transport protocol for the SIP proxy. Valid values: tcp, udp, tls.",
			},
		},
		MarkdownDescription: "Manages a SIP proxy configuration with the Infinity service.",
	}
}

func (r *InfinitySIPProxyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	plan := &InfinitySIPProxyResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createRequest := &config.SIPProxyCreateRequest{
		Name:      plan.Name.ValueString(),
		Address:   plan.Address.ValueString(),
		Transport: plan.Transport.ValueString(),
	}

	// Only set optional fields if they are not null in the plan
	if !plan.Description.IsNull() {
		createRequest.Description = plan.Description.ValueString()
	}
	if !plan.Port.IsNull() {
		port := int(plan.Port.ValueInt32())
		createRequest.Port = &port
	}

	createResponse, err := r.InfinityClient.Config().CreateSIPProxy(ctx, createRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Infinity SIP proxy",
			fmt.Sprintf("Could not create Infinity SIP proxy: %s", err),
		)
		return
	}

	resourceID, err := createResponse.ResourceID()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Retrieving Infinity SIP proxy ID",
			fmt.Sprintf("Could not retrieve ID for created Infinity SIP proxy: %s", err),
		)
		return
	}

	// Read the state from the API to get all computed values
	model, err := r.read(ctx, resourceID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Created Infinity SIP proxy",
			fmt.Sprintf("Could not read created Infinity SIP proxy with ID %d: %s", resourceID, err),
		)
		return
	}
	tflog.Trace(ctx, fmt.Sprintf("created Infinity SIP proxy with ID: %s, name: %s", model.ID, model.Name))

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *InfinitySIPProxyResource) read(ctx context.Context, resourceID int) (*InfinitySIPProxyResourceModel, error) {
	var data InfinitySIPProxyResourceModel

	srv, err := r.InfinityClient.Config().GetSIPProxy(ctx, resourceID)
	if err != nil {
		return nil, err
	}

	if srv.ResourceURI == "" {
		return nil, fmt.Errorf("SIP proxy with ID %d not found", resourceID)
	}

	data.ID = types.StringValue(srv.ResourceURI)
	data.ResourceID = types.Int32Value(int32(resourceID)) // #nosec G115 -- API values are expected to be within int32 range
	data.Name = types.StringValue(srv.Name)
	data.Description = types.StringValue(srv.Description)
	data.Address = types.StringValue(srv.Address)
	data.Transport = types.StringValue(srv.Transport)

	if srv.Port != nil {
		data.Port = types.Int32Value(int32(*srv.Port)) // #nosec G115 -- API values are expected to be within int32 range
	} else {
		data.Port = types.Int32Null()
	}

	return &data, nil
}

func (r *InfinitySIPProxyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	state := &InfinitySIPProxyResourceModel{}

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
			"Error Reading Infinity SIP proxy",
			fmt.Sprintf("Could not read Infinity SIP proxy: %s", err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *InfinitySIPProxyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	plan := &InfinitySIPProxyResourceModel{}
	state := &InfinitySIPProxyResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resourceID := int(state.ResourceID.ValueInt32())

	updateRequest := &config.SIPProxyUpdateRequest{
		Name:      plan.Name.ValueString(),
		Address:   plan.Address.ValueString(),
		Transport: plan.Transport.ValueString(),
	}

	if !plan.Description.IsNull() {
		updateRequest.Description = plan.Description.ValueString()
	}
	if !plan.Port.IsNull() {
		port := int(plan.Port.ValueInt32())
		updateRequest.Port = &port
	}

	_, err := r.InfinityClient.Config().UpdateSIPProxy(ctx, resourceID, updateRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Infinity SIP proxy",
			fmt.Sprintf("Could not update Infinity SIP proxy with ID %d: %s", resourceID, err),
		)
		return
	}

	// Re-read the resource to get the latest state
	updatedModel, err := r.read(ctx, resourceID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Updated Infinity SIP proxy",
			fmt.Sprintf("Could not read updated Infinity SIP proxy with ID %d: %s", resourceID, err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, updatedModel)...)
}

func (r *InfinitySIPProxyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	state := &InfinitySIPProxyResourceModel{}

	tflog.Info(ctx, "Deleting Infinity SIP proxy")

	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.InfinityClient.Config().DeleteSIPProxy(ctx, int(state.ResourceID.ValueInt32()))

	// Ignore 404 Not Found and Lookup errors on delete
	if err != nil && !isNotFoundError(err) && !isLookupError(err) {
		resp.Diagnostics.AddError(
			"Error Deleting Infinity SIP proxy",
			fmt.Sprintf("Could not delete Infinity SIP proxy with ID %s: %s", state.ID.ValueString(), err),
		)
		return
	}
}

func (r *InfinitySIPProxyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resourceID, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Import ID must be a valid integer for the resource ID. Got: %s", req.ID),
		)
		return
	}

	tflog.Trace(ctx, fmt.Sprintf("Importing Infinity SIP proxy with resource ID: %d", resourceID))

	// Read the resource from the API
	model, err := r.read(ctx, resourceID)
	if err != nil {
		// Check if the error is a 404 (not found)
		if isNotFoundError(err) {
			resp.Diagnostics.AddError(
				"Infinity SIP Proxy Not Found",
				fmt.Sprintf("Infinity SIP proxy with resource ID %d not found.", resourceID),
			)
			return
		}
		resp.Diagnostics.AddError(
			"Error Importing Infinity SIP Proxy",
			fmt.Sprintf("Could not import Infinity SIP proxy with resource ID %d: %s", resourceID, err),
		)
		return
	}

	// Set the state from the imported resource
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}
