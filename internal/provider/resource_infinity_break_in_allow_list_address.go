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

	"github.com/pexip/terraform-provider-pexip/internal/provider/validators"
)

var (
	_ resource.ResourceWithImportState = (*InfinityBreakInAllowListAddressResource)(nil)
)

type InfinityBreakInAllowListAddressResource struct {
	InfinityClient InfinityClient
}

type InfinityBreakInAllowListAddressResourceModel struct {
	ID                     types.String `tfsdk:"id"`
	ResourceID             types.Int32  `tfsdk:"resource_id"`
	Name                   types.String `tfsdk:"name"`
	Description            types.String `tfsdk:"description"`
	Address                types.String `tfsdk:"address"`
	Prefix                 types.Int64  `tfsdk:"prefix"`
	AllowlistEntryType     types.String `tfsdk:"allowlist_entry_type"`
	IgnoreIncorrectAliases types.Bool   `tfsdk:"ignore_incorrect_aliases"`
	IgnoreIncorrectPins    types.Bool   `tfsdk:"ignore_incorrect_pins"`
}

func (r *InfinityBreakInAllowListAddressResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_infinity_break_in_allow_list_address"
}

func (r *InfinityBreakInAllowListAddressResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *InfinityBreakInAllowListAddressResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Resource URI for the break-in allow list address in Infinity",
			},
			"resource_id": schema.Int32Attribute{
				Computed:            true,
				MarkdownDescription: "The resource integer identifier for the break-in allow list address in Infinity",
			},
			"name": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					stringvalidator.LengthAtMost(100),
				},
				MarkdownDescription: "The name of the break-in allow list entry. Maximum length: 100 characters.",
			},
			"description": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(500),
				},
				MarkdownDescription: "Description of the break-in allow list entry. Maximum length: 500 characters.",
			},
			"address": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					validators.IPAddress(),
				},
				MarkdownDescription: "The IP address for this allow list entry.",
			},
			"prefix": schema.Int64Attribute{
				Required: true,
				Validators: []validator.Int64{
					int64validator.Between(0, 128),
				},
				MarkdownDescription: "The network prefix length (CIDR notation). Valid range: 0-128 (supports both IPv4 and IPv6).",
			},
			"allowlist_entry_type": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf("temporary", "permanent"),
				},
				MarkdownDescription: "The type of allow list entry. Valid values: temporary, permanent.",
			},
			"ignore_incorrect_aliases": schema.BoolAttribute{
				Required:            true,
				MarkdownDescription: "Whether to ignore incorrect alias attempts from this address range.",
			},
			"ignore_incorrect_pins": schema.BoolAttribute{
				Required:            true,
				MarkdownDescription: "Whether to ignore incorrect PIN attempts from this address range.",
			},
		},
		MarkdownDescription: "Manages a break-in allow list address with the Infinity service. Break-in allow list addresses define IP address ranges that are exempt from certain security restrictions, allowing specified networks to bypass break-in attempt detection for specific scenarios.",
	}
}

func (r *InfinityBreakInAllowListAddressResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	plan := &InfinityBreakInAllowListAddressResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createRequest := &config.BreakInAllowListAddressCreateRequest{
		Name:                   plan.Name.ValueString(),
		Description:            plan.Description.ValueString(),
		Address:                plan.Address.ValueString(),
		Prefix:                 int(plan.Prefix.ValueInt64()),
		AllowlistEntryType:     plan.AllowlistEntryType.ValueString(),
		IgnoreIncorrectAliases: plan.IgnoreIncorrectAliases.ValueBool(),
		IgnoreIncorrectPins:    plan.IgnoreIncorrectPins.ValueBool(),
	}

	createResponse, err := r.InfinityClient.Config().CreateBreakInAllowListAddress(ctx, createRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Infinity break-in allow list address",
			fmt.Sprintf("Could not create Infinity break-in allow list address: %s", err),
		)
		return
	}

	resourceID, err := createResponse.ResourceID()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Retrieving Infinity break-in allow list address ID",
			fmt.Sprintf("Could not retrieve ID for created Infinity break-in allow list address: %s", err),
		)
		return
	}

	// Read the state from the API to get all computed values
	model, err := r.read(ctx, resourceID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Created Infinity break-in allow list address",
			fmt.Sprintf("Could not read created Infinity break-in allow list address with ID %d: %s", resourceID, err),
		)
		return
	}
	tflog.Trace(ctx, fmt.Sprintf("created Infinity break-in allow list address with ID: %s, name: %s", model.ID, model.Name))

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *InfinityBreakInAllowListAddressResource) read(ctx context.Context, resourceID int) (*InfinityBreakInAllowListAddressResourceModel, error) {
	var data InfinityBreakInAllowListAddressResourceModel

	srv, err := r.InfinityClient.Config().GetBreakInAllowListAddress(ctx, resourceID)
	if err != nil {
		return nil, err
	}

	if srv.ResourceURI == "" {
		return nil, fmt.Errorf("break-in allow list address with ID %d not found", resourceID)
	}

	data.ID = types.StringValue(srv.ResourceURI)
	data.ResourceID = types.Int32Value(int32(resourceID)) // #nosec G115 -- API values are expected to be within int32 range
	data.Name = types.StringValue(srv.Name)
	data.Description = types.StringValue(srv.Description)
	data.Address = types.StringValue(srv.Address)
	data.Prefix = types.Int64Value(int64(srv.Prefix))
	data.AllowlistEntryType = types.StringValue(srv.AllowlistEntryType)
	data.IgnoreIncorrectAliases = types.BoolValue(srv.IgnoreIncorrectAliases)
	data.IgnoreIncorrectPins = types.BoolValue(srv.IgnoreIncorrectPins)

	return &data, nil
}

func (r *InfinityBreakInAllowListAddressResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	state := &InfinityBreakInAllowListAddressResourceModel{}

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
			"Error Reading Infinity break-in allow list address",
			fmt.Sprintf("Could not read Infinity break-in allow list address: %s", err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *InfinityBreakInAllowListAddressResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	plan := &InfinityBreakInAllowListAddressResourceModel{}
	state := &InfinityBreakInAllowListAddressResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateRequest := &config.BreakInAllowListAddressUpdateRequest{
		Name:               plan.Name.ValueString(),
		Description:        plan.Description.ValueString(),
		Address:            plan.Address.ValueString(),
		AllowlistEntryType: plan.AllowlistEntryType.ValueString(),
	}

	// Handle optional pointer fields
	if !plan.Prefix.IsNull() && !plan.Prefix.IsUnknown() {
		prefix := int(plan.Prefix.ValueInt64())
		updateRequest.Prefix = &prefix
	}

	if !plan.IgnoreIncorrectAliases.IsNull() && !plan.IgnoreIncorrectAliases.IsUnknown() {
		ignoreAliases := plan.IgnoreIncorrectAliases.ValueBool()
		updateRequest.IgnoreIncorrectAliases = &ignoreAliases
	}

	if !plan.IgnoreIncorrectPins.IsNull() && !plan.IgnoreIncorrectPins.IsUnknown() {
		ignorePins := plan.IgnoreIncorrectPins.ValueBool()
		updateRequest.IgnoreIncorrectPins = &ignorePins
	}

	resourceID := int(state.ResourceID.ValueInt32())
	_, err := r.InfinityClient.Config().UpdateBreakInAllowListAddress(ctx, resourceID, updateRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Infinity break-in allow list address",
			fmt.Sprintf("Could not update Infinity break-in allow list address: %s", err),
		)
		return
	}

	// Read the state from the API to get all computed values
	model, err := r.read(ctx, resourceID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Updated Infinity break-in allow list address",
			fmt.Sprintf("Could not read updated Infinity break-in allow list address with ID %d: %s", resourceID, err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *InfinityBreakInAllowListAddressResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	state := &InfinityBreakInAllowListAddressResourceModel{}

	tflog.Info(ctx, "Deleting Infinity break-in allow list address")

	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.InfinityClient.Config().DeleteBreakInAllowListAddress(ctx, int(state.ResourceID.ValueInt32()))

	// Ignore 404 Not Found and Lookup errors on delete
	if err != nil && !isNotFoundError(err) && !isLookupError(err) {
		resp.Diagnostics.AddError(
			"Error Deleting Infinity break-in allow list address",
			fmt.Sprintf("Could not delete Infinity break-in allow list address with ID %s: %s", state.ID.ValueString(), err),
		)
		return
	}
}

func (r *InfinityBreakInAllowListAddressResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resourceID, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Import ID must be a valid integer for the resource ID. Got: %s", req.ID),
		)
		return
	}

	tflog.Trace(ctx, fmt.Sprintf("Importing Infinity break-in allow list address with resource ID: %d", resourceID))

	// Read the resource from the API
	model, err := r.read(ctx, resourceID)
	if err != nil {
		// Check if the error is a 404 (not found)
		if isNotFoundError(err) {
			resp.Diagnostics.AddError(
				"Infinity Break-in Allow List Address Not Found",
				fmt.Sprintf("Infinity break-in allow list address with resource ID %d not found.", resourceID),
			)
			return
		}
		resp.Diagnostics.AddError(
			"Error Importing Infinity Break-in Allow List Address",
			fmt.Sprintf("Could not import Infinity break-in allow list address with resource ID %d: %s", resourceID, err),
		)
		return
	}

	// Set the state from the imported resource
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}
