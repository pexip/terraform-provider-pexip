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

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"

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
	_ resource.ResourceWithImportState = (*InfinityTeamsProxyResource)(nil)
)

type InfinityTeamsProxyResource struct {
	InfinityClient InfinityClient
}

type InfinityTeamsProxyResourceModel struct {
	ID                   types.String `tfsdk:"id"`
	ResourceID           types.Int32  `tfsdk:"resource_id"`
	Name                 types.String `tfsdk:"name"`
	Description          types.String `tfsdk:"description"`
	Address              types.String `tfsdk:"address"`
	Port                 types.Int32  `tfsdk:"port"`
	AzureTenant          types.String `tfsdk:"azure_tenant"`
	EventhubID           types.String `tfsdk:"eventhub_id"`
	MinNumberOfInstances types.Int32  `tfsdk:"min_number_of_instances"`
	NotificationsEnabled types.Bool   `tfsdk:"notifications_enabled"`
	NotificationsQueue   types.String `tfsdk:"notifications_queue"`
}

func (r *InfinityTeamsProxyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_infinity_teams_proxy"
}

func (r *InfinityTeamsProxyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *InfinityTeamsProxyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Resource URI for the Teams proxy in Infinity",
			},
			"resource_id": schema.Int32Attribute{
				Computed:            true,
				MarkdownDescription: "The resource integer identifier for the Teams proxy in Infinity",
			},
			"name": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(250),
				},
				MarkdownDescription: "The name used to refer to this Teams proxy. Maximum length: 250 characters.",
			},
			"description": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(""),
				Validators: []validator.String{
					stringvalidator.LengthAtMost(250),
				},
				MarkdownDescription: "A description of the Teams proxy. Maximum length: 250 characters.",
			},
			"address": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(255),
				},
				MarkdownDescription: "The address or hostname of the Teams proxy. Maximum length: 255 characters.",
			},
			"port": schema.Int32Attribute{
				Optional: true,
				Computed: true,
				Default:  int32default.StaticInt32(443),
				Validators: []validator.Int32{
					int32validator.Between(1, 65535),
				},
				MarkdownDescription: "The port number for the Teams proxy. Range: 1 to 65535.",
			},
			"azure_tenant": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(255),
				},
				MarkdownDescription: "The Azure tenant ID for the Teams proxy. Maximum length: 255 characters.",
			},
			"eventhub_id": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(255),
				},
				MarkdownDescription: "The event hub identifier for the Teams proxy. Maximum length: 255 characters.",
			},
			"min_number_of_instances": schema.Int32Attribute{
				Optional:            true,
				Computed:            true,
				Default:             int32default.StaticInt32(1),
				Validators: []validator.Int32{
					int32validator.AtLeast(0),
				},
				MarkdownDescription: "The minimum number of instances for the Teams proxy.",
			},
			"notifications_enabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Whether notifications are enabled for the Teams proxy.",
			},
			"notifications_queue": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Sensitive:           true,
				MarkdownDescription: "The Connection string primary key for the Azure Event Hub (standard access policy). This is in the format Endpoint=sb://examplevmss-tzfk6222uo-ehn.servicebus.windows.net/;SharedAccessKeyName=standard_access_policy;SharedAccessKey=[string]/[string]/[string]=;",
			},
		},
		MarkdownDescription: "Manages a Teams proxy configuration with the Infinity service.",
	}
}

func (r *InfinityTeamsProxyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	plan := &InfinityTeamsProxyResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createRequest := &config.TeamsProxyCreateRequest{
		Name:                 plan.Name.ValueString(),
		Description:          plan.Description.ValueString(),
		Address:              plan.Address.ValueString(),
		Port:                 int(plan.Port.ValueInt32()),
		AzureTenant:          plan.AzureTenant.ValueString(),
		MinNumberOfInstances: int(plan.MinNumberOfInstances.ValueInt32()),
	}

	// Only set optional fields if they are not null in the plan
	if !plan.EventhubID.IsNull() && !plan.EventhubID.IsUnknown() {
		eventhubID := plan.EventhubID.ValueString()
		createRequest.EventhubID = &eventhubID
	}
	if !plan.NotificationsEnabled.IsNull() && !plan.NotificationsEnabled.IsUnknown() {
		createRequest.NotificationsEnabled = plan.NotificationsEnabled.ValueBool()
	}
	if !plan.NotificationsQueue.IsNull() && !plan.NotificationsQueue.IsUnknown() {
		notificationsQueue := plan.NotificationsQueue.ValueString()
		createRequest.NotificationsQueue = &notificationsQueue
	}

	createResponse, err := r.InfinityClient.Config().CreateTeamsProxy(ctx, createRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Infinity Teams proxy",
			fmt.Sprintf("Could not create Infinity Teams proxy: %s", err),
		)
		return
	}

	resourceID, err := createResponse.ResourceID()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Retrieving Infinity Teams proxy ID",
			fmt.Sprintf("Could not retrieve ID for created Infinity Teams proxy: %s", err),
		)
		return
	}

	// Read the state from the API to get all computed values
	model, err := r.read(ctx, resourceID, plan.NotificationsQueue.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Created Infinity Teams proxy",
			fmt.Sprintf("Could not read created Infinity Teams proxy with ID %d: %s", resourceID, err),
		)
		return
	}
	tflog.Trace(ctx, fmt.Sprintf("created Infinity Teams proxy with ID: %s, name: %s", model.ID, model.Name))

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *InfinityTeamsProxyResource) read(ctx context.Context, resourceID int, notificationsQueue string) (*InfinityTeamsProxyResourceModel, error) {
	var data InfinityTeamsProxyResourceModel

	srv, err := r.InfinityClient.Config().GetTeamsProxy(ctx, resourceID)
	if err != nil {
		return nil, err
	}

	if srv.ResourceURI == "" {
		return nil, fmt.Errorf("teams proxy with ID %d not found", resourceID)
	}

	data.ID = types.StringValue(srv.ResourceURI)
	data.ResourceID = types.Int32Value(int32(resourceID)) // #nosec G115 -- API values are expected to be within int32 range
	data.Name = types.StringValue(srv.Name)
	data.Description = types.StringValue(srv.Description)
	data.Address = types.StringValue(srv.Address)
	data.Port = types.Int32Value(int32(srv.Port)) // #nosec G115 -- API values are expected to be within int32 range
	data.AzureTenant = types.StringValue(srv.AzureTenant)
	data.MinNumberOfInstances = types.Int32Value(int32(srv.MinNumberOfInstances)) // #nosec G115 -- API values are expected to be within int32 range
	data.NotificationsEnabled = types.BoolValue(srv.NotificationsEnabled)
	data.EventhubID = types.StringPointerValue(srv.EventhubID)
	data.NotificationsQueue = types.StringValue(notificationsQueue) // The server does not return the notifications queue, so we use the provided one

	return &data, nil
}

func (r *InfinityTeamsProxyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	state := &InfinityTeamsProxyResourceModel{}

	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resourceID := int(state.ResourceID.ValueInt32())
	state, err := r.read(ctx, resourceID, state.NotificationsQueue.ValueString())
	if err != nil {
		// Check if the error is a 404 (not found)
		if isNotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading Infinity Teams proxy",
			fmt.Sprintf("Could not read Infinity Teams proxy: %s", err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *InfinityTeamsProxyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	plan := &InfinityTeamsProxyResourceModel{}
	state := &InfinityTeamsProxyResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resourceID := int(state.ResourceID.ValueInt32())

	updateRequest := &config.TeamsProxyUpdateRequest{
		Name:                 plan.Name.ValueString(),
		Description:          plan.Description.ValueString(),
		Address:              plan.Address.ValueString(),
		Port:                 int(plan.Port.ValueInt32()),
		AzureTenant:          plan.AzureTenant.ValueString(),
		MinNumberOfInstances: int(plan.MinNumberOfInstances.ValueInt32()),
		NotificationsEnabled: plan.NotificationsEnabled.ValueBool(),
	}

	// Only set optional fields if they are not null in the plan
	if !plan.EventhubID.IsNull() && !plan.EventhubID.IsUnknown() {
		eventhubID := plan.EventhubID.ValueString()
		updateRequest.EventhubID = &eventhubID
	}
	if !plan.NotificationsEnabled.IsNull() && !plan.NotificationsEnabled.IsUnknown() {
		updateRequest.NotificationsEnabled = plan.NotificationsEnabled.ValueBool()
	}
	if !plan.NotificationsQueue.IsNull() && !plan.NotificationsQueue.IsUnknown() {
		notificationsQueue := plan.NotificationsQueue.ValueString()
		updateRequest.NotificationsQueue = &notificationsQueue
	}

	_, err := r.InfinityClient.Config().UpdateTeamsProxy(ctx, resourceID, updateRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Infinity Teams proxy",
			fmt.Sprintf("Could not update Infinity Teams proxy with ID %d: %s", resourceID, err),
		)
		return
	}

	// Re-read the resource to get the latest state
	updatedModel, err := r.read(ctx, resourceID, plan.NotificationsQueue.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Updated Infinity Teams proxy",
			fmt.Sprintf("Could not read updated Infinity Teams proxy with ID %d: %s", resourceID, err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, updatedModel)...)
}

func (r *InfinityTeamsProxyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	state := &InfinityTeamsProxyResourceModel{}

	tflog.Info(ctx, "Deleting Infinity Teams proxy")

	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.InfinityClient.Config().DeleteTeamsProxy(ctx, int(state.ResourceID.ValueInt32()))

	// Ignore 404 Not Found and Lookup errors on delete
	if err != nil && !isNotFoundError(err) && !isLookupError(err) {
		resp.Diagnostics.AddError(
			"Error Deleting Infinity Teams proxy",
			fmt.Sprintf("Could not delete Infinity Teams proxy with ID %s: %s", state.ID.ValueString(), err),
		)
		return
	}
}

func (r *InfinityTeamsProxyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resourceID, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Import ID must be a valid integer for the resource ID. Got: %s", req.ID),
		)
		return
	}

	tflog.Trace(ctx, fmt.Sprintf("Importing Infinity Teams proxy with resource ID: %d", resourceID))

	// Read the resource from the API
	model, err := r.read(ctx, resourceID, "")
	if err != nil {
		// Check if the error is a 404 (not found)
		if isNotFoundError(err) {
			resp.Diagnostics.AddError(
				"Infinity Teams Proxy Not Found",
				fmt.Sprintf("Infinity Teams proxy with resource ID %d not found.", resourceID),
			)
			return
		}
		resp.Diagnostics.AddError(
			"Error Importing Infinity Teams Proxy",
			fmt.Sprintf("Could not import Infinity Teams proxy with resource ID %d: %s", resourceID, err),
		)
		return
	}

	// Set the state from the imported resource
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}
