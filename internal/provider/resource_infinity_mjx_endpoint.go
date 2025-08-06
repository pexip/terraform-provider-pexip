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

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/pexip/go-infinity-sdk/v38/config"
)

var (
	_ resource.ResourceWithImportState = (*InfinityMjxEndpointResource)(nil)
)

type InfinityMjxEndpointResource struct {
	InfinityClient InfinityClient
}

type InfinityMjxEndpointResourceModel struct {
	ID                             types.String `tfsdk:"id"`
	ResourceID                     types.Int32  `tfsdk:"resource_id"`
	Name                           types.String `tfsdk:"name"`
	Description                    types.String `tfsdk:"description"`
	EndpointType                   types.String `tfsdk:"endpoint_type"`
	RoomResourceEmail              types.String `tfsdk:"room_resource_email"`
	MjxEndpointGroup               types.String `tfsdk:"mjx_endpoint_group"`
	APIAddress                     types.String `tfsdk:"api_address"`
	APIPort                        types.Int64  `tfsdk:"api_port"`
	APIUsername                    types.String `tfsdk:"api_username"`
	APIPassword                    types.String `tfsdk:"api_password"`
	UseHTTPS                       types.String `tfsdk:"use_https"`
	VerifyCert                     types.String `tfsdk:"verify_cert"`
	PolyUsername                   types.String `tfsdk:"poly_username"`
	PolyPassword                   types.String `tfsdk:"poly_password"`
	PolyRaiseAlarmsForThisEndpoint types.Bool   `tfsdk:"poly_raise_alarms_for_this_endpoint"`
	WebexDeviceID                  types.String `tfsdk:"webex_device_id"`
}

func (r *InfinityMjxEndpointResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_infinity_mjx_endpoint"
}

func (r *InfinityMjxEndpointResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *InfinityMjxEndpointResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Resource URI for the MJX endpoint in Infinity",
			},
			"resource_id": schema.Int32Attribute{
				Computed:            true,
				MarkdownDescription: "The resource integer identifier for the MJX endpoint in Infinity",
			},
			"name": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					stringvalidator.LengthAtMost(100),
				},
				MarkdownDescription: "The name of the MJX endpoint. Maximum length: 100 characters.",
			},
			"description": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(500),
				},
				MarkdownDescription: "Description of the MJX endpoint. Maximum length: 500 characters.",
			},
			"endpoint_type": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf("polycom", "cisco", "webex"),
				},
				MarkdownDescription: "The type of MJX endpoint. Valid values: polycom, cisco, webex.",
			},
			"room_resource_email": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(254),
				},
				MarkdownDescription: "The email address of the room resource associated with this endpoint.",
			},
			"mjx_endpoint_group": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The MJX endpoint group URI this endpoint belongs to.",
			},
			"api_address": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(100),
				},
				MarkdownDescription: "The API address for the endpoint management interface.",
			},
			"api_port": schema.Int64Attribute{
				Optional: true,
				Validators: []validator.Int64{
					int64validator.Between(1, 65535),
				},
				MarkdownDescription: "The API port for the endpoint management interface. Valid range: 1-65535.",
			},
			"api_username": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(100),
				},
				MarkdownDescription: "The username for API authentication to the endpoint.",
			},
			"api_password": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				MarkdownDescription: "The password for API authentication to the endpoint. This field is sensitive.",
			},
			"use_https": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf("yes", "no"),
				},
				MarkdownDescription: "Whether to use HTTPS for API communication. Valid values: yes, no.",
			},
			"verify_cert": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf("yes", "no"),
				},
				MarkdownDescription: "Whether to verify SSL certificates. Valid values: yes, no.",
			},
			"poly_username": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(100),
				},
				MarkdownDescription: "The username for Polycom-specific authentication.",
			},
			"poly_password": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				MarkdownDescription: "The password for Polycom-specific authentication. This field is sensitive.",
			},
			"poly_raise_alarms_for_this_endpoint": schema.BoolAttribute{
				Required:            true,
				MarkdownDescription: "Whether to raise alarms for this Polycom endpoint.",
			},
			"webex_device_id": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(100),
				},
				MarkdownDescription: "The Webex device ID for Webex endpoints.",
			},
		},
		MarkdownDescription: "Manages an MJX endpoint with the Infinity service. MJX endpoints represent Microsoft Teams integrated endpoints such as Polycom, Cisco, and Webex devices that can be managed and monitored through Pexip Infinity for hybrid Teams deployments.",
	}
}

func (r *InfinityMjxEndpointResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	plan := &InfinityMjxEndpointResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createRequest := &config.MjxEndpointCreateRequest{
		Name:                           plan.Name.ValueString(),
		Description:                    plan.Description.ValueString(),
		EndpointType:                   plan.EndpointType.ValueString(),
		RoomResourceEmail:              plan.RoomResourceEmail.ValueString(),
		UseHTTPS:                       plan.UseHTTPS.ValueString(),
		VerifyCert:                     plan.VerifyCert.ValueString(),
		PolyRaiseAlarmsForThisEndpoint: plan.PolyRaiseAlarmsForThisEndpoint.ValueBool(),
	}

	// Handle optional pointer fields
	if !plan.MjxEndpointGroup.IsNull() && !plan.MjxEndpointGroup.IsUnknown() {
		group := plan.MjxEndpointGroup.ValueString()
		createRequest.MjxEndpointGroup = &group
	}

	if !plan.APIAddress.IsNull() && !plan.APIAddress.IsUnknown() {
		address := plan.APIAddress.ValueString()
		createRequest.APIAddress = &address
	}

	if !plan.APIPort.IsNull() && !plan.APIPort.IsUnknown() {
		port := int(plan.APIPort.ValueInt64())
		createRequest.APIPort = &port
	}

	if !plan.APIUsername.IsNull() && !plan.APIUsername.IsUnknown() {
		username := plan.APIUsername.ValueString()
		createRequest.APIUsername = &username
	}

	if !plan.APIPassword.IsNull() && !plan.APIPassword.IsUnknown() {
		password := plan.APIPassword.ValueString()
		createRequest.APIPassword = &password
	}

	if !plan.PolyUsername.IsNull() && !plan.PolyUsername.IsUnknown() {
		username := plan.PolyUsername.ValueString()
		createRequest.PolyUsername = &username
	}

	if !plan.PolyPassword.IsNull() && !plan.PolyPassword.IsUnknown() {
		password := plan.PolyPassword.ValueString()
		createRequest.PolyPassword = &password
	}

	if !plan.WebexDeviceID.IsNull() && !plan.WebexDeviceID.IsUnknown() {
		deviceID := plan.WebexDeviceID.ValueString()
		createRequest.WebexDeviceID = &deviceID
	}

	createResponse, err := r.InfinityClient.Config().CreateMjxEndpoint(ctx, createRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Infinity MJX endpoint",
			fmt.Sprintf("Could not create Infinity MJX endpoint: %s", err),
		)
		return
	}

	resourceID, err := createResponse.ResourceID()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Retrieving Infinity MJX endpoint ID",
			fmt.Sprintf("Could not retrieve ID for created Infinity MJX endpoint: %s", err),
		)
		return
	}

	// Read the state from the API to get all computed values
	model, err := r.read(ctx, resourceID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Created Infinity MJX endpoint",
			fmt.Sprintf("Could not read created Infinity MJX endpoint with ID %d: %s", resourceID, err),
		)
		return
	}
	tflog.Trace(ctx, fmt.Sprintf("created Infinity MJX endpoint with ID: %s, name: %s", model.ID, model.Name))

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *InfinityMjxEndpointResource) read(ctx context.Context, resourceID int) (*InfinityMjxEndpointResourceModel, error) {
	var data InfinityMjxEndpointResourceModel

	srv, err := r.InfinityClient.Config().GetMjxEndpoint(ctx, resourceID)
	if err != nil {
		return nil, err
	}

	if srv.ResourceURI == "" {
		return nil, fmt.Errorf("MJX endpoint with ID %d not found", resourceID)
	}

	data.ID = types.StringValue(srv.ResourceURI)
	data.ResourceID = types.Int32Value(int32(resourceID)) // #nosec G115 -- API values are expected to be within int32 range
	data.Name = types.StringValue(srv.Name)
	data.Description = types.StringValue(srv.Description)
	data.EndpointType = types.StringValue(srv.EndpointType)
	data.RoomResourceEmail = types.StringValue(srv.RoomResourceEmail)
	data.UseHTTPS = types.StringValue(srv.UseHTTPS)
	data.VerifyCert = types.StringValue(srv.VerifyCert)
	data.PolyRaiseAlarmsForThisEndpoint = types.BoolValue(srv.PolyRaiseAlarmsForThisEndpoint)

	// Handle optional pointer fields
	if srv.MjxEndpointGroup != nil {
		data.MjxEndpointGroup = types.StringValue(*srv.MjxEndpointGroup)
	} else {
		data.MjxEndpointGroup = types.StringNull()
	}

	if srv.APIAddress != nil {
		data.APIAddress = types.StringValue(*srv.APIAddress)
	} else {
		data.APIAddress = types.StringNull()
	}

	if srv.APIPort != nil {
		data.APIPort = types.Int64Value(int64(*srv.APIPort))
	} else {
		data.APIPort = types.Int64Null()
	}

	if srv.APIUsername != nil {
		data.APIUsername = types.StringValue(*srv.APIUsername)
	} else {
		data.APIUsername = types.StringNull()
	}

	if srv.APIPassword != nil {
		data.APIPassword = types.StringValue(*srv.APIPassword)
	} else {
		data.APIPassword = types.StringNull()
	}

	if srv.PolyUsername != nil {
		data.PolyUsername = types.StringValue(*srv.PolyUsername)
	} else {
		data.PolyUsername = types.StringNull()
	}

	if srv.PolyPassword != nil {
		data.PolyPassword = types.StringValue(*srv.PolyPassword)
	} else {
		data.PolyPassword = types.StringNull()
	}

	if srv.WebexDeviceID != nil {
		data.WebexDeviceID = types.StringValue(*srv.WebexDeviceID)
	} else {
		data.WebexDeviceID = types.StringNull()
	}

	return &data, nil
}

func (r *InfinityMjxEndpointResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	state := &InfinityMjxEndpointResourceModel{}

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
			"Error Reading Infinity MJX endpoint",
			fmt.Sprintf("Could not read Infinity MJX endpoint: %s", err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *InfinityMjxEndpointResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	plan := &InfinityMjxEndpointResourceModel{}
	state := &InfinityMjxEndpointResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateRequest := &config.MjxEndpointUpdateRequest{
		Name:              plan.Name.ValueString(),
		Description:       plan.Description.ValueString(),
		EndpointType:      plan.EndpointType.ValueString(),
		RoomResourceEmail: plan.RoomResourceEmail.ValueString(),
		UseHTTPS:          plan.UseHTTPS.ValueString(),
		VerifyCert:        plan.VerifyCert.ValueString(),
	}

	// Handle optional pointer fields
	if !plan.MjxEndpointGroup.IsNull() && !plan.MjxEndpointGroup.IsUnknown() {
		group := plan.MjxEndpointGroup.ValueString()
		updateRequest.MjxEndpointGroup = &group
	}

	if !plan.APIAddress.IsNull() && !plan.APIAddress.IsUnknown() {
		address := plan.APIAddress.ValueString()
		updateRequest.APIAddress = &address
	}

	if !plan.APIPort.IsNull() && !plan.APIPort.IsUnknown() {
		port := int(plan.APIPort.ValueInt64())
		updateRequest.APIPort = &port
	}

	if !plan.APIUsername.IsNull() && !plan.APIUsername.IsUnknown() {
		username := plan.APIUsername.ValueString()
		updateRequest.APIUsername = &username
	}

	if !plan.APIPassword.IsNull() && !plan.APIPassword.IsUnknown() {
		password := plan.APIPassword.ValueString()
		updateRequest.APIPassword = &password
	}

	if !plan.PolyUsername.IsNull() && !plan.PolyUsername.IsUnknown() {
		username := plan.PolyUsername.ValueString()
		updateRequest.PolyUsername = &username
	}

	if !plan.PolyPassword.IsNull() && !plan.PolyPassword.IsUnknown() {
		password := plan.PolyPassword.ValueString()
		updateRequest.PolyPassword = &password
	}

	if !plan.PolyRaiseAlarmsForThisEndpoint.IsNull() && !plan.PolyRaiseAlarmsForThisEndpoint.IsUnknown() {
		alarms := plan.PolyRaiseAlarmsForThisEndpoint.ValueBool()
		updateRequest.PolyRaiseAlarmsForThisEndpoint = &alarms
	}

	if !plan.WebexDeviceID.IsNull() && !plan.WebexDeviceID.IsUnknown() {
		deviceID := plan.WebexDeviceID.ValueString()
		updateRequest.WebexDeviceID = &deviceID
	}

	resourceID := int(state.ResourceID.ValueInt32())
	_, err := r.InfinityClient.Config().UpdateMjxEndpoint(ctx, resourceID, updateRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Infinity MJX endpoint",
			fmt.Sprintf("Could not update Infinity MJX endpoint: %s", err),
		)
		return
	}

	// Read the state from the API to get all computed values
	model, err := r.read(ctx, resourceID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Updated Infinity MJX endpoint",
			fmt.Sprintf("Could not read updated Infinity MJX endpoint with ID %d: %s", resourceID, err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *InfinityMjxEndpointResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	state := &InfinityMjxEndpointResourceModel{}

	tflog.Info(ctx, "Deleting Infinity MJX endpoint")

	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.InfinityClient.Config().DeleteMjxEndpoint(ctx, int(state.ResourceID.ValueInt32()))

	// Ignore 404 Not Found and Lookup errors on delete
	if err != nil && !isNotFoundError(err) && !isLookupError(err) {
		resp.Diagnostics.AddError(
			"Error Deleting Infinity MJX endpoint",
			fmt.Sprintf("Could not delete Infinity MJX endpoint with ID %s: %s", state.ID.ValueString(), err),
		)
		return
	}
}

func (r *InfinityMjxEndpointResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resourceID, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Import ID must be a valid integer for the resource ID. Got: %s", req.ID),
		)
		return
	}

	tflog.Trace(ctx, fmt.Sprintf("Importing Infinity MJX endpoint with resource ID: %d", resourceID))

	// Read the resource from the API
	model, err := r.read(ctx, resourceID)
	if err != nil {
		// Check if the error is a 404 (not found)
		if isNotFoundError(err) {
			resp.Diagnostics.AddError(
				"Infinity MJX Endpoint Not Found",
				fmt.Sprintf("Infinity MJX endpoint with resource ID %d not found.", resourceID),
			)
			return
		}
		resp.Diagnostics.AddError(
			"Error Importing Infinity MJX Endpoint",
			fmt.Sprintf("Could not import Infinity MJX endpoint with resource ID %d: %s", resourceID, err),
		)
		return
	}

	// Set the state from the imported resource
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}
