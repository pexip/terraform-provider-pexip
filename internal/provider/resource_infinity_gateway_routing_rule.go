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
	_ resource.ResourceWithImportState = (*InfinityGatewayRoutingRuleResource)(nil)
)

type InfinityGatewayRoutingRuleResource struct {
	InfinityClient InfinityClient
}

type InfinityGatewayRoutingRuleResourceModel struct {
	ID               types.String `tfsdk:"id"`
	ResourceID       types.Int32  `tfsdk:"resource_id"`
	Name             types.String `tfsdk:"name"`
	Description      types.String `tfsdk:"description"`
	Priority         types.Int32  `tfsdk:"priority"`
	Enable           types.Bool   `tfsdk:"enable"`
	MatchString      types.String `tfsdk:"match_string"`
	ReplaceString    types.String `tfsdk:"replace_string"`
	CalledDeviceType types.String `tfsdk:"called_device_type"`
	OutgoingProtocol types.String `tfsdk:"outgoing_protocol"`
	CallType         types.String `tfsdk:"call_type"`
	IvrTheme         types.String `tfsdk:"ivr_theme"`
}

func (r *InfinityGatewayRoutingRuleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_infinity_gateway_routing_rule"
}

func (r *InfinityGatewayRoutingRuleResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *InfinityGatewayRoutingRuleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Resource URI for the gateway routing rule in Infinity",
			},
			"resource_id": schema.Int32Attribute{
				Computed:            true,
				MarkdownDescription: "The resource integer identifier for the gateway routing rule in Infinity",
			},
			"name": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(250),
				},
				MarkdownDescription: "The unique name of the gateway routing rule. Maximum length: 250 characters.",
			},
			"description": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(250),
				},
				MarkdownDescription: "A description of the gateway routing rule. Maximum length: 250 characters.",
			},
			"priority": schema.Int32Attribute{
				Required: true,
				Validators: []validator.Int32{
					int32validator.AtLeast(1),
				},
				MarkdownDescription: "The priority of the gateway routing rule (lower numbers have higher priority).",
			},
			"enable": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Whether the gateway routing rule is enabled. Defaults to true.",
			},
			"match_string": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(250),
				},
				MarkdownDescription: "Regular expression pattern to match incoming calls. Maximum length: 250 characters.",
			},
			"replace_string": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(250),
				},
				MarkdownDescription: "Pattern for outgoing call transformation. Maximum length: 250 characters.",
			},
			"called_device_type": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.OneOf("unknown", "conference", "gateway", "ip_pbx"),
				},
				MarkdownDescription: "Type of called device. Valid choices: unknown, conference, gateway, ip_pbx.",
			},
			"outgoing_protocol": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.OneOf("sip", "h323"),
				},
				MarkdownDescription: "Outgoing protocol. Valid choices: sip, h323.",
			},
			"call_type": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.OneOf("audio", "video"),
				},
				MarkdownDescription: "Call type. Valid choices: audio, video.",
			},
			"ivr_theme": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Reference to IVR theme resource URI.",
			},
		},
		MarkdownDescription: "Manages a gateway routing rule configuration with the Infinity service.",
	}
}

func (r *InfinityGatewayRoutingRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	plan := &InfinityGatewayRoutingRuleResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createRequest := &config.GatewayRoutingRuleCreateRequest{
		Name:        plan.Name.ValueString(),
		Priority:    int(plan.Priority.ValueInt32()),
		MatchString: plan.MatchString.ValueString(),
	}

	// Set optional fields
	if !plan.Description.IsNull() {
		createRequest.Description = plan.Description.ValueString()
	}
	if !plan.Enable.IsNull() {
		createRequest.Enable = plan.Enable.ValueBool()
	}
	if !plan.ReplaceString.IsNull() {
		createRequest.ReplaceString = plan.ReplaceString.ValueString()
	}
	if !plan.CalledDeviceType.IsNull() {
		createRequest.CalledDeviceType = plan.CalledDeviceType.ValueString()
	}
	if !plan.OutgoingProtocol.IsNull() {
		createRequest.OutgoingProtocol = plan.OutgoingProtocol.ValueString()
	}
	if !plan.CallType.IsNull() {
		createRequest.CallType = plan.CallType.ValueString()
	}
	if !plan.IvrTheme.IsNull() {
		ivrTheme := plan.IvrTheme.ValueString()
		createRequest.IVRTheme = &ivrTheme
	}

	createResponse, err := r.InfinityClient.Config().CreateGatewayRoutingRule(ctx, createRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Infinity gateway routing rule",
			fmt.Sprintf("Could not create Infinity gateway routing rule: %s", err),
		)
		return
	}

	resourceID, err := createResponse.ResourceID()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Retrieving Infinity gateway routing rule ID",
			fmt.Sprintf("Could not retrieve ID for created Infinity gateway routing rule: %s", err),
		)
		return
	}

	// Read the state from the API to get all computed values
	model, err := r.read(ctx, resourceID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Created Infinity gateway routing rule",
			fmt.Sprintf("Could not read created Infinity gateway routing rule with ID %d: %s", resourceID, err),
		)
		return
	}
	tflog.Trace(ctx, fmt.Sprintf("created Infinity gateway routing rule with ID: %s, name: %s", model.ID, model.Name))

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *InfinityGatewayRoutingRuleResource) read(ctx context.Context, resourceID int) (*InfinityGatewayRoutingRuleResourceModel, error) {
	var data InfinityGatewayRoutingRuleResourceModel

	srv, err := r.InfinityClient.Config().GetGatewayRoutingRule(ctx, resourceID)
	if err != nil {
		return nil, err
	}

	if srv.ResourceURI == "" {
		return nil, fmt.Errorf("gateway routing rule with ID %d not found", resourceID)
	}

	data.ID = types.StringValue(srv.ResourceURI)
	data.ResourceID = types.Int32Value(int32(resourceID)) // #nosec G115 -- API values are expected to be within int32 range
	data.Name = types.StringValue(srv.Name)
	data.Description = types.StringValue(srv.Description)
	data.Priority = types.Int32Value(int32(srv.Priority)) // #nosec G115 -- API values are expected to be within int32 range
	data.MatchString = types.StringValue(srv.MatchString)
	data.ReplaceString = types.StringValue(srv.ReplaceString)
	data.CalledDeviceType = types.StringValue(srv.CalledDeviceType)
	data.OutgoingProtocol = types.StringValue(srv.OutgoingProtocol)
	data.CallType = types.StringValue(srv.CallType)
	if srv.IVRTheme != nil {
		data.IvrTheme = types.StringValue(*srv.IVRTheme)
	} else {
		data.IvrTheme = types.StringNull()
	}

	// Set boolean field
	data.Enable = types.BoolValue(srv.Enable)

	return &data, nil
}

func (r *InfinityGatewayRoutingRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	state := &InfinityGatewayRoutingRuleResourceModel{}

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
			"Error Reading Infinity gateway routing rule",
			fmt.Sprintf("Could not read Infinity gateway routing rule: %s", err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *InfinityGatewayRoutingRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	plan := &InfinityGatewayRoutingRuleResourceModel{}
	state := &InfinityGatewayRoutingRuleResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resourceID := int(state.ResourceID.ValueInt32())

	updateRequest := &config.GatewayRoutingRuleUpdateRequest{
		Name:        plan.Name.ValueString(),
		MatchString: plan.MatchString.ValueString(),
	}

	// Set priority as pointer
	priority := int(plan.Priority.ValueInt32())
	updateRequest.Priority = &priority

	// Set optional fields
	if !plan.Description.IsNull() {
		updateRequest.Description = plan.Description.ValueString()
	}
	if !plan.Enable.IsNull() {
		enable := plan.Enable.ValueBool()
		updateRequest.Enable = &enable
	}
	if !plan.ReplaceString.IsNull() {
		updateRequest.ReplaceString = plan.ReplaceString.ValueString()
	}
	if !plan.CalledDeviceType.IsNull() {
		updateRequest.CalledDeviceType = plan.CalledDeviceType.ValueString()
	}
	if !plan.OutgoingProtocol.IsNull() {
		updateRequest.OutgoingProtocol = plan.OutgoingProtocol.ValueString()
	}
	if !plan.CallType.IsNull() {
		updateRequest.CallType = plan.CallType.ValueString()
	}
	if !plan.IvrTheme.IsNull() {
		ivrTheme := plan.IvrTheme.ValueString()
		updateRequest.IVRTheme = &ivrTheme
	}

	_, err := r.InfinityClient.Config().UpdateGatewayRoutingRule(ctx, resourceID, updateRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Infinity gateway routing rule",
			fmt.Sprintf("Could not update Infinity gateway routing rule with ID %d: %s", resourceID, err),
		)
		return
	}

	// Re-read the resource to get the latest state
	updatedModel, err := r.read(ctx, resourceID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Updated Infinity gateway routing rule",
			fmt.Sprintf("Could not read updated Infinity gateway routing rule with ID %d: %s", resourceID, err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, updatedModel)...)
}

func (r *InfinityGatewayRoutingRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	state := &InfinityGatewayRoutingRuleResourceModel{}

	tflog.Info(ctx, "Deleting Infinity gateway routing rule")

	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.InfinityClient.Config().DeleteGatewayRoutingRule(ctx, int(state.ResourceID.ValueInt32()))

	// Ignore 404 Not Found and Lookup errors on delete
	if err != nil && !isNotFoundError(err) && !isLookupError(err) {
		resp.Diagnostics.AddError(
			"Error Deleting Infinity gateway routing rule",
			fmt.Sprintf("Could not delete Infinity gateway routing rule with ID %s: %s", state.ID.ValueString(), err),
		)
		return
	}
}

func (r *InfinityGatewayRoutingRuleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resourceID, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Import ID must be a valid integer for the resource ID. Got: %s", req.ID),
		)
		return
	}

	tflog.Trace(ctx, fmt.Sprintf("Importing Infinity gateway routing rule with resource ID: %d", resourceID))

	// Read the resource from the API
	model, err := r.read(ctx, resourceID)
	if err != nil {
		// Check if the error is a 404 (not found)
		if isNotFoundError(err) {
			resp.Diagnostics.AddError(
				"Infinity Gateway Routing Rule Not Found",
				fmt.Sprintf("Infinity gateway routing rule with resource ID %d not found.", resourceID),
			)
			return
		}
		resp.Diagnostics.AddError(
			"Error Importing Infinity Gateway Routing Rule",
			fmt.Sprintf("Could not import Infinity gateway routing rule with resource ID %d: %s", resourceID, err),
		)
		return
	}

	// Set the state from the imported resource
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}
