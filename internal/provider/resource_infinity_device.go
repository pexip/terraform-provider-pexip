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
	_ resource.ResourceWithImportState = (*InfinityDeviceResource)(nil)
)

type InfinityDeviceResource struct {
	InfinityClient InfinityClient
}

type InfinityDeviceResourceModel struct {
	ID                          types.String `tfsdk:"id"`
	ResourceID                  types.Int32  `tfsdk:"resource_id"`
	Alias                       types.String `tfsdk:"alias"`
	Description                 types.String `tfsdk:"description"`
	Username                    types.String `tfsdk:"username"`
	Password                    types.String `tfsdk:"password"`
	PrimaryOwnerEmailAddress    types.String `tfsdk:"primary_owner_email_address"`
	EnableSIP                   types.Bool   `tfsdk:"enable_sip"`
	EnableH323                  types.Bool   `tfsdk:"enable_h323"`
	EnableInfinityConnectNonSSO types.Bool   `tfsdk:"enable_infinity_connect_non_sso"`
	EnableInfinityConnectSSO    types.Bool   `tfsdk:"enable_infinity_connect_sso"`
	EnableStandardSSO           types.Bool   `tfsdk:"enable_standard_sso"`
	SSOIdentityProviderGroup    types.String `tfsdk:"sso_identity_provider_group"`
	Tag                         types.String `tfsdk:"tag"`
	SyncTag                     types.String `tfsdk:"sync_tag"`
}

func (r *InfinityDeviceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_infinity_device"
}

func (r *InfinityDeviceResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *InfinityDeviceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Resource URI for the device in Infinity",
			},
			"resource_id": schema.Int32Attribute{
				Computed:            true,
				MarkdownDescription: "The resource integer identifier for the device in Infinity",
			},
			"alias": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(250),
				},
				MarkdownDescription: "The unique alias name of the device. Maximum length: 250 characters.",
			},
			"description": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(250),
				},
				MarkdownDescription: "A description of the device. Maximum length: 250 characters.",
			},
			"username": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(250),
				},
				MarkdownDescription: "The username for device authentication. Maximum length: 250 characters.",
			},
			"password": schema.StringAttribute{
				Optional:  true,
				Computed:  true,
				Sensitive: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(100),
				},
				MarkdownDescription: "The password for device authentication. Maximum length: 100 characters.",
			},
			"primary_owner_email_address": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(100),
				},
				MarkdownDescription: "Email address of the device owner. Maximum length: 100 characters.",
			},
			"enable_sip": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Whether SIP is enabled for this device. Defaults to false.",
			},
			"enable_h323": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Whether H.323 is enabled for this device. Defaults to false.",
			},
			"enable_infinity_connect_non_sso": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Whether Infinity Connect without SSO is enabled. Defaults to false.",
			},
			"enable_infinity_connect_sso": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Whether Infinity Connect with SSO is enabled. Defaults to false.",
			},
			"enable_standard_sso": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Whether standard SSO is enabled. Defaults to false.",
			},
			"sso_identity_provider_group": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "SSO identity provider group for authentication.",
			},
			"tag": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(250),
				},
				MarkdownDescription: "A tag for categorizing the device. Maximum length: 250 characters.",
			},
			"sync_tag": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(250),
				},
				MarkdownDescription: "A sync tag for external system integration. Maximum length: 250 characters.",
			},
		},
		MarkdownDescription: "Manages a device configuration with the Infinity service.",
	}
}

func (r *InfinityDeviceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	plan := &InfinityDeviceResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createRequest := &config.DeviceCreateRequest{
		Alias: plan.Alias.ValueString(),
	}

	// Set optional string fields
	if !plan.Description.IsNull() {
		createRequest.Description = plan.Description.ValueString()
	}
	if !plan.Username.IsNull() {
		createRequest.Username = plan.Username.ValueString()
	}
	if !plan.Password.IsNull() {
		createRequest.Password = plan.Password.ValueString()
	}
	if !plan.PrimaryOwnerEmailAddress.IsNull() {
		createRequest.PrimaryOwnerEmailAddress = plan.PrimaryOwnerEmailAddress.ValueString()
	}
	if !plan.Tag.IsNull() {
		createRequest.Tag = plan.Tag.ValueString()
	}
	if !plan.SyncTag.IsNull() {
		createRequest.SyncTag = plan.SyncTag.ValueString()
	}

	// Set boolean fields (required fields need default values)
	if !plan.EnableSIP.IsNull() {
		createRequest.EnableSIP = plan.EnableSIP.ValueBool()
	} else {
		createRequest.EnableSIP = false
	}

	if !plan.EnableH323.IsNull() {
		createRequest.EnableH323 = plan.EnableH323.ValueBool()
	} else {
		createRequest.EnableH323 = false
	}

	if !plan.EnableInfinityConnectNonSSO.IsNull() {
		createRequest.EnableInfinityConnectNonSSO = plan.EnableInfinityConnectNonSSO.ValueBool()
	} else {
		createRequest.EnableInfinityConnectNonSSO = false
	}

	if !plan.EnableInfinityConnectSSO.IsNull() {
		createRequest.EnableInfinityConnectSSO = plan.EnableInfinityConnectSSO.ValueBool()
	} else {
		createRequest.EnableInfinityConnectSSO = false
	}

	if !plan.EnableStandardSSO.IsNull() {
		createRequest.EnableStandardSSO = plan.EnableStandardSSO.ValueBool()
	} else {
		createRequest.EnableStandardSSO = false
	}

	// Set optional pointer fields
	if !plan.SSOIdentityProviderGroup.IsNull() {
		ssoGroup := plan.SSOIdentityProviderGroup.ValueString()
		createRequest.SSOIdentityProviderGroup = &ssoGroup
	}

	createResponse, err := r.InfinityClient.Config().CreateDevice(ctx, createRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Infinity device",
			fmt.Sprintf("Could not create Infinity device: %s", err),
		)
		return
	}

	resourceID, err := createResponse.ResourceID()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Retrieving Infinity device ID",
			fmt.Sprintf("Could not retrieve ID for created Infinity device: %s", err),
		)
		return
	}

	// Read the state from the API to get all computed values
	model, err := r.read(ctx, resourceID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Created Infinity device",
			fmt.Sprintf("Could not read created Infinity device with ID %d: %s", resourceID, err),
		)
		return
	}
	tflog.Trace(ctx, fmt.Sprintf("created Infinity device with ID: %s, alias: %s", model.ID, model.Alias))

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *InfinityDeviceResource) read(ctx context.Context, resourceID int) (*InfinityDeviceResourceModel, error) {
	var data InfinityDeviceResourceModel

	srv, err := r.InfinityClient.Config().GetDevice(ctx, resourceID)
	if err != nil {
		return nil, err
	}

	if srv.ResourceURI == "" {
		return nil, fmt.Errorf("device with ID %d not found", resourceID)
	}

	data.ID = types.StringValue(srv.ResourceURI)
	data.ResourceID = types.Int32Value(int32(resourceID)) // #nosec G115 -- API values are expected to be within int32 range
	data.Alias = types.StringValue(srv.Alias)
	data.Description = types.StringValue(srv.Description)
	data.Username = types.StringValue(srv.Username)
	data.Password = types.StringValue(srv.Password)
	data.PrimaryOwnerEmailAddress = types.StringValue(srv.PrimaryOwnerEmailAddress)
	data.Tag = types.StringValue(srv.Tag)
	data.SyncTag = types.StringValue(srv.SyncTag)

	// Set boolean fields
	data.EnableSIP = types.BoolValue(srv.EnableSIP)
	data.EnableH323 = types.BoolValue(srv.EnableH323)
	data.EnableInfinityConnectNonSSO = types.BoolValue(srv.EnableInfinityConnectNonSSO)
	data.EnableInfinityConnectSSO = types.BoolValue(srv.EnableInfinityConnectSSO)
	data.EnableStandardSSO = types.BoolValue(srv.EnableStandardSSO)

	// Handle pointer field
	if srv.SSOIdentityProviderGroup != nil {
		data.SSOIdentityProviderGroup = types.StringValue(*srv.SSOIdentityProviderGroup)
	} else {
		data.SSOIdentityProviderGroup = types.StringNull()
	}

	return &data, nil
}

func (r *InfinityDeviceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	state := &InfinityDeviceResourceModel{}

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
			"Error Reading Infinity device",
			fmt.Sprintf("Could not read Infinity device: %s", err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *InfinityDeviceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	plan := &InfinityDeviceResourceModel{}
	state := &InfinityDeviceResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resourceID := int(state.ResourceID.ValueInt32())

	updateRequest := &config.DeviceUpdateRequest{
		Alias: plan.Alias.ValueString(),
	}

	// Set optional string fields
	if !plan.Description.IsNull() {
		updateRequest.Description = plan.Description.ValueString()
	}
	if !plan.Username.IsNull() {
		updateRequest.Username = plan.Username.ValueString()
	}
	if !plan.Password.IsNull() {
		updateRequest.Password = plan.Password.ValueString()
	}
	if !plan.PrimaryOwnerEmailAddress.IsNull() {
		updateRequest.PrimaryOwnerEmailAddress = plan.PrimaryOwnerEmailAddress.ValueString()
	}
	if !plan.Tag.IsNull() {
		updateRequest.Tag = plan.Tag.ValueString()
	}
	if !plan.SyncTag.IsNull() {
		updateRequest.SyncTag = plan.SyncTag.ValueString()
	}

	// Set optional boolean fields (use pointers for update requests)
	if !plan.EnableSIP.IsNull() {
		enableSIP := plan.EnableSIP.ValueBool()
		updateRequest.EnableSIP = &enableSIP
	}
	if !plan.EnableH323.IsNull() {
		enableH323 := plan.EnableH323.ValueBool()
		updateRequest.EnableH323 = &enableH323
	}
	if !plan.EnableInfinityConnectNonSSO.IsNull() {
		enableInfinityConnectNonSSO := plan.EnableInfinityConnectNonSSO.ValueBool()
		updateRequest.EnableInfinityConnectNonSSO = &enableInfinityConnectNonSSO
	}
	if !plan.EnableInfinityConnectSSO.IsNull() {
		enableInfinityConnectSSO := plan.EnableInfinityConnectSSO.ValueBool()
		updateRequest.EnableInfinityConnectSSO = &enableInfinityConnectSSO
	}
	if !plan.EnableStandardSSO.IsNull() {
		enableStandardSSO := plan.EnableStandardSSO.ValueBool()
		updateRequest.EnableStandardSSO = &enableStandardSSO
	}

	// Set optional pointer fields
	if !plan.SSOIdentityProviderGroup.IsNull() {
		ssoGroup := plan.SSOIdentityProviderGroup.ValueString()
		updateRequest.SSOIdentityProviderGroup = &ssoGroup
	}

	_, err := r.InfinityClient.Config().UpdateDevice(ctx, resourceID, updateRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Infinity device",
			fmt.Sprintf("Could not update Infinity device with ID %d: %s", resourceID, err),
		)
		return
	}

	// Re-read the resource to get the latest state
	updatedModel, err := r.read(ctx, resourceID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Updated Infinity device",
			fmt.Sprintf("Could not read updated Infinity device with ID %d: %s", resourceID, err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, updatedModel)...)
}

func (r *InfinityDeviceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	state := &InfinityDeviceResourceModel{}

	tflog.Info(ctx, "Deleting Infinity device")

	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.InfinityClient.Config().DeleteDevice(ctx, int(state.ResourceID.ValueInt32()))

	// Ignore 404 Not Found and Lookup errors on delete
	if err != nil && !isNotFoundError(err) && !isLookupError(err) {
		resp.Diagnostics.AddError(
			"Error Deleting Infinity device",
			fmt.Sprintf("Could not delete Infinity device with ID %s: %s", state.ID.ValueString(), err),
		)
		return
	}
}

func (r *InfinityDeviceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resourceID, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Import ID must be a valid integer for the resource ID. Got: %s", req.ID),
		)
		return
	}

	tflog.Trace(ctx, fmt.Sprintf("Importing Infinity device with resource ID: %d", resourceID))

	// Read the resource from the API
	model, err := r.read(ctx, resourceID)
	if err != nil {
		// Check if the error is a 404 (not found)
		if isNotFoundError(err) {
			resp.Diagnostics.AddError(
				"Infinity Device Not Found",
				fmt.Sprintf("Infinity device with resource ID %d not found.", resourceID),
			)
			return
		}
		resp.Diagnostics.AddError(
			"Error Importing Infinity Device",
			fmt.Sprintf("Could not import Infinity device with resource ID %d: %s", resourceID, err),
		)
		return
	}

	// Set the state from the imported resource
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}
