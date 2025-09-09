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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/pexip/go-infinity-sdk/v38/config"
)

var (
	_ resource.ResourceWithImportState = (*InfinityHTTPProxyResource)(nil)
)

type InfinityHTTPProxyResource struct {
	InfinityClient InfinityClient
}

type InfinityHTTPProxyResourceModel struct {
	ID         types.String `tfsdk:"id"`
	ResourceID types.Int32  `tfsdk:"resource_id"`
	Name       types.String `tfsdk:"name"`
	Address    types.String `tfsdk:"address"`
	Port       types.Int32  `tfsdk:"port"`
	Protocol   types.String `tfsdk:"protocol"`
	Username   types.String `tfsdk:"username"`
	Password   types.String `tfsdk:"password"`
}

func (r *InfinityHTTPProxyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_infinity_http_proxy"
}

func (r *InfinityHTTPProxyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *InfinityHTTPProxyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Resource URI for the HTTP proxy in Infinity",
			},
			"resource_id": schema.Int32Attribute{
				Computed:            true,
				MarkdownDescription: "The resource integer identifier for the HTTP proxy in Infinity",
			},
			"name": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(250),
				},
				MarkdownDescription: "The name used to refer to this HTTP proxy. Maximum length: 250 characters.",
			},
			"address": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(255),
				},
				MarkdownDescription: "The address or hostname of the HTTP proxy. Maximum length: 255 characters.",
			},
			"port": schema.Int32Attribute{
				Optional: true,
				Computed: true,
				Default:  int32default.StaticInt32(8080),
				Validators: []validator.Int32{
					int32validator.Between(1, 65535),
				},
				MarkdownDescription: "The port number for the HTTP proxy. Range: 1 to 65535.",
			},
			"protocol": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("http"),
				Validators: []validator.String{
					stringvalidator.OneOf("http"),
				},
				MarkdownDescription: "The protocol for the HTTP proxy. Valid values: http.",
			},
			"username": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(100),
				},
				MarkdownDescription: "Username for authentication to the HTTP proxy. Maximum length: 100 characters.",
			},
			"password": schema.StringAttribute{
				Optional:  true,
				Computed:  true,
				Sensitive: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(100),
				},
				MarkdownDescription: "Password for authentication to the HTTP proxy. Maximum length: 100 characters.",
			},
		},
		MarkdownDescription: "Manages an HTTP proxy configuration with the Infinity service.",
	}
}

func (r *InfinityHTTPProxyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	plan := &InfinityHTTPProxyResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createRequest := &config.HTTPProxyCreateRequest{
		Name:     plan.Name.ValueString(),
		Address:  plan.Address.ValueString(),
		Protocol: plan.Protocol.ValueString(),
	}

	// Only set optional fields if they are not null in the plan
	if !plan.Port.IsNull() {
		port := int(plan.Port.ValueInt32())
		createRequest.Port = &port
	}
	if !plan.Username.IsNull() {
		createRequest.Username = plan.Username.ValueString()
	}
	if !plan.Password.IsNull() {
		createRequest.Password = plan.Password.ValueString()
	}

	createResponse, err := r.InfinityClient.Config().CreateHTTPProxy(ctx, createRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Infinity HTTP proxy",
			fmt.Sprintf("Could not create Infinity HTTP proxy: %s", err),
		)
		return
	}

	resourceID, err := createResponse.ResourceID()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Retrieving Infinity HTTP proxy ID",
			fmt.Sprintf("Could not retrieve ID for created Infinity HTTP proxy: %s", err),
		)
		return
	}

	// Read the state from the API to get all computed values
	model, err := r.read(ctx, resourceID, plan.Password.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Created Infinity HTTP proxy",
			fmt.Sprintf("Could not read created Infinity HTTP proxy with ID %d: %s", resourceID, err),
		)
		return
	}
	tflog.Trace(ctx, fmt.Sprintf("created Infinity HTTP proxy with ID: %s, name: %s", model.ID, model.Name))

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *InfinityHTTPProxyResource) read(ctx context.Context, resourceID int, password string) (*InfinityHTTPProxyResourceModel, error) {
	var data InfinityHTTPProxyResourceModel

	srv, err := r.InfinityClient.Config().GetHTTPProxy(ctx, resourceID)
	if err != nil {
		return nil, err
	}

	if srv.ResourceURI == "" {
		return nil, fmt.Errorf("HTTP proxy with ID %d not found", resourceID)
	}

	data.ID = types.StringValue(srv.ResourceURI)
	data.ResourceID = types.Int32Value(int32(resourceID)) // #nosec G115 -- API values are expected to be within int32 range
	data.Name = types.StringValue(srv.Name)
	data.Address = types.StringValue(srv.Address)
	data.Protocol = types.StringValue(srv.Protocol)
	data.Username = types.StringValue(srv.Username)
	data.Password = types.StringValue(password) // The server does not return the password, so we use the provided one

	if srv.Port != nil {
		data.Port = types.Int32Value(int32(*srv.Port)) // #nosec G115 -- API values are expected to be within int32 range
	} else {
		data.Port = types.Int32Null()
	}

	return &data, nil
}

func (r *InfinityHTTPProxyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	state := &InfinityHTTPProxyResourceModel{}

	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resourceID := int(state.ResourceID.ValueInt32())
	state, err := r.read(ctx, resourceID, state.Password.ValueString())
	if err != nil {
		// Check if the error is a 404 (not found)
		if isNotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading Infinity HTTP proxy",
			fmt.Sprintf("Could not read Infinity HTTP proxy: %s", err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *InfinityHTTPProxyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	plan := &InfinityHTTPProxyResourceModel{}
	state := &InfinityHTTPProxyResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resourceID := int(state.ResourceID.ValueInt32())

	updateRequest := &config.HTTPProxyUpdateRequest{
		Name:     plan.Name.ValueString(),
		Address:  plan.Address.ValueString(),
		Protocol: plan.Protocol.ValueString(),
	}

	if !plan.Port.IsNull() {
		port := int(plan.Port.ValueInt32())
		updateRequest.Port = &port
	}
	if !plan.Username.IsNull() {
		updateRequest.Username = plan.Username.ValueString()
	}
	if !plan.Password.IsNull() {
		updateRequest.Password = plan.Password.ValueString()
	}

	_, err := r.InfinityClient.Config().UpdateHTTPProxy(ctx, resourceID, updateRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Infinity HTTP proxy",
			fmt.Sprintf("Could not update Infinity HTTP proxy with ID %d: %s", resourceID, err),
		)
		return
	}

	// Re-read the resource to get the latest state
	updatedModel, err := r.read(ctx, resourceID, plan.Password.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Updated Infinity HTTP proxy",
			fmt.Sprintf("Could not read updated Infinity HTTP proxy with ID %d: %s", resourceID, err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, updatedModel)...)
}

func (r *InfinityHTTPProxyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	state := &InfinityHTTPProxyResourceModel{}

	tflog.Info(ctx, "Deleting Infinity HTTP proxy")

	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.InfinityClient.Config().DeleteHTTPProxy(ctx, int(state.ResourceID.ValueInt32()))

	// Ignore 404 Not Found and Lookup errors on delete
	if err != nil && !isNotFoundError(err) && !isLookupError(err) {
		resp.Diagnostics.AddError(
			"Error Deleting Infinity HTTP proxy",
			fmt.Sprintf("Could not delete Infinity HTTP proxy with ID %s: %s", state.ID.ValueString(), err),
		)
		return
	}
}

func (r *InfinityHTTPProxyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resourceID, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Import ID must be a valid integer for the resource ID. Got: %s", req.ID),
		)
		return
	}

	tflog.Trace(ctx, fmt.Sprintf("Importing Infinity HTTP proxy with resource ID: %d", resourceID))

	// Read the resource from the API
	model, err := r.read(ctx, resourceID, "")
	if err != nil {
		// Check if the error is a 404 (not found)
		if isNotFoundError(err) {
			resp.Diagnostics.AddError(
				"Infinity HTTP Proxy Not Found",
				fmt.Sprintf("Infinity HTTP proxy with resource ID %d not found.", resourceID),
			)
			return
		}
		resp.Diagnostics.AddError(
			"Error Importing Infinity HTTP Proxy",
			fmt.Sprintf("Could not import Infinity HTTP proxy with resource ID %d: %s", resourceID, err),
		)
		return
	}

	// Set the state from the imported resource
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}
