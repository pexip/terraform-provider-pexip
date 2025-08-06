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
	_ resource.ResourceWithImportState = (*InfinityScheduledScalingResource)(nil)
)

type InfinityScheduledScalingResource struct {
	InfinityClient InfinityClient
}

type InfinityScheduledScalingResourceModel struct {
	ID                 types.String `tfsdk:"id"`
	ResourceID         types.Int32  `tfsdk:"resource_id"`
	PolicyName         types.String `tfsdk:"policy_name"`
	PolicyType         types.String `tfsdk:"policy_type"`
	ResourceIdentifier types.String `tfsdk:"resource_identifier"`
	Enabled            types.Bool   `tfsdk:"enabled"`
	InstancesToAdd     types.Int64  `tfsdk:"instances_to_add"`
	MinutesInAdvance   types.Int64  `tfsdk:"minutes_in_advance"`
	LocalTimezone      types.String `tfsdk:"local_timezone"`
	StartDate          types.String `tfsdk:"start_date"`
	TimeFrom           types.String `tfsdk:"time_from"`
	TimeTo             types.String `tfsdk:"time_to"`
	Mon                types.Bool   `tfsdk:"mon"`
	Tue                types.Bool   `tfsdk:"tue"`
	Wed                types.Bool   `tfsdk:"wed"`
	Thu                types.Bool   `tfsdk:"thu"`
	Fri                types.Bool   `tfsdk:"fri"`
	Sat                types.Bool   `tfsdk:"sat"`
	Sun                types.Bool   `tfsdk:"sun"`
	Updated            types.String `tfsdk:"updated"`
}

func (r *InfinityScheduledScalingResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_infinity_scheduled_scaling"
}

func (r *InfinityScheduledScalingResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *InfinityScheduledScalingResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Resource URI for the scheduled scaling policy in Infinity",
			},
			"resource_id": schema.Int32Attribute{
				Computed:            true,
				MarkdownDescription: "The resource integer identifier for the scheduled scaling policy in Infinity",
			},
			"policy_name": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					stringvalidator.LengthAtMost(100),
				},
				MarkdownDescription: "The name of the scheduled scaling policy. Maximum length: 100 characters.",
			},
			"policy_type": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf("worker_vm", "management_vm"),
				},
				MarkdownDescription: "The type of resource to scale. Valid values: worker_vm, management_vm.",
			},
			"resource_identifier": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
				MarkdownDescription: "The identifier of the resource group or deployment to scale.",
			},
			"enabled": schema.BoolAttribute{
				Required:            true,
				MarkdownDescription: "Whether this scheduled scaling policy is enabled.",
			},
			"instances_to_add": schema.Int64Attribute{
				Required: true,
				Validators: []validator.Int64{
					int64validator.Between(0, 100),
				},
				MarkdownDescription: "The number of instances to add during scaling. Valid range: 0-100.",
			},
			"minutes_in_advance": schema.Int64Attribute{
				Required: true,
				Validators: []validator.Int64{
					int64validator.Between(0, 1440),
				},
				MarkdownDescription: "How many minutes in advance to start scaling. Valid range: 0-1440 (24 hours).",
			},
			"local_timezone": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
				MarkdownDescription: "The local timezone for scheduling (e.g., 'America/New_York', 'Europe/London').",
			},
			"start_date": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
				MarkdownDescription: "The start date for this scaling policy in YYYY-MM-DD format.",
			},
			"time_from": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
				MarkdownDescription: "The start time for scaling in HH:MM format.",
			},
			"time_to": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
				MarkdownDescription: "The end time for scaling in HH:MM format.",
			},
			"mon": schema.BoolAttribute{
				Required:            true,
				MarkdownDescription: "Whether to apply scaling on Monday.",
			},
			"tue": schema.BoolAttribute{
				Required:            true,
				MarkdownDescription: "Whether to apply scaling on Tuesday.",
			},
			"wed": schema.BoolAttribute{
				Required:            true,
				MarkdownDescription: "Whether to apply scaling on Wednesday.",
			},
			"thu": schema.BoolAttribute{
				Required:            true,
				MarkdownDescription: "Whether to apply scaling on Thursday.",
			},
			"fri": schema.BoolAttribute{
				Required:            true,
				MarkdownDescription: "Whether to apply scaling on Friday.",
			},
			"sat": schema.BoolAttribute{
				Required:            true,
				MarkdownDescription: "Whether to apply scaling on Saturday.",
			},
			"sun": schema.BoolAttribute{
				Required:            true,
				MarkdownDescription: "Whether to apply scaling on Sunday.",
			},
			"updated": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The timestamp when this scheduled scaling policy was last updated.",
			},
		},
		MarkdownDescription: "Manages a scheduled scaling policy with the Infinity service. Scheduled scaling policies automatically add or remove VM instances based on time-based schedules, helping optimize resource utilization for predictable workload patterns.",
	}
}

func (r *InfinityScheduledScalingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	plan := &InfinityScheduledScalingResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createRequest := &config.ScheduledScalingCreateRequest{
		PolicyName:         plan.PolicyName.ValueString(),
		PolicyType:         plan.PolicyType.ValueString(),
		ResourceIdentifier: plan.ResourceIdentifier.ValueString(),
		Enabled:            plan.Enabled.ValueBool(),
		InstancesToAdd:     int(plan.InstancesToAdd.ValueInt64()),
		MinutesInAdvance:   int(plan.MinutesInAdvance.ValueInt64()),
		LocalTimezone:      plan.LocalTimezone.ValueString(),
		StartDate:          plan.StartDate.ValueString(),
		TimeFrom:           plan.TimeFrom.ValueString(),
		TimeTo:             plan.TimeTo.ValueString(),
		Mon:                plan.Mon.ValueBool(),
		Tue:                plan.Tue.ValueBool(),
		Wed:                plan.Wed.ValueBool(),
		Thu:                plan.Thu.ValueBool(),
		Fri:                plan.Fri.ValueBool(),
		Sat:                plan.Sat.ValueBool(),
		Sun:                plan.Sun.ValueBool(),
	}

	createResponse, err := r.InfinityClient.Config().CreateScheduledScaling(ctx, createRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Infinity scheduled scaling policy",
			fmt.Sprintf("Could not create Infinity scheduled scaling policy: %s", err),
		)
		return
	}

	resourceID, err := createResponse.ResourceID()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Retrieving Infinity scheduled scaling policy ID",
			fmt.Sprintf("Could not retrieve ID for created Infinity scheduled scaling policy: %s", err),
		)
		return
	}

	// Read the state from the API to get all computed values
	model, err := r.read(ctx, resourceID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Created Infinity scheduled scaling policy",
			fmt.Sprintf("Could not read created Infinity scheduled scaling policy with ID %d: %s", resourceID, err),
		)
		return
	}
	tflog.Trace(ctx, fmt.Sprintf("created Infinity scheduled scaling policy with ID: %s, name: %s", model.ID, model.PolicyName))

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *InfinityScheduledScalingResource) read(ctx context.Context, resourceID int) (*InfinityScheduledScalingResourceModel, error) {
	var data InfinityScheduledScalingResourceModel

	srv, err := r.InfinityClient.Config().GetScheduledScaling(ctx, resourceID)
	if err != nil {
		return nil, err
	}

	if srv.ResourceURI == "" {
		return nil, fmt.Errorf("scheduled scaling policy with ID %d not found", resourceID)
	}

	data.ID = types.StringValue(srv.ResourceURI)
	data.ResourceID = types.Int32Value(int32(resourceID)) // #nosec G115 -- API values are expected to be within int32 range
	data.PolicyName = types.StringValue(srv.PolicyName)
	data.PolicyType = types.StringValue(srv.PolicyType)
	data.ResourceIdentifier = types.StringValue(srv.ResourceIdentifier)
	data.Enabled = types.BoolValue(srv.Enabled)
	data.InstancesToAdd = types.Int64Value(int64(srv.InstancesToAdd))
	data.MinutesInAdvance = types.Int64Value(int64(srv.MinutesInAdvance))
	data.LocalTimezone = types.StringValue(srv.LocalTimezone)
	data.StartDate = types.StringValue(srv.StartDate)
	data.TimeFrom = types.StringValue(srv.TimeFrom)
	data.TimeTo = types.StringValue(srv.TimeTo)
	data.Mon = types.BoolValue(srv.Mon)
	data.Tue = types.BoolValue(srv.Tue)
	data.Wed = types.BoolValue(srv.Wed)
	data.Thu = types.BoolValue(srv.Thu)
	data.Fri = types.BoolValue(srv.Fri)
	data.Sat = types.BoolValue(srv.Sat)
	data.Sun = types.BoolValue(srv.Sun)

	if srv.Updated != nil {
		data.Updated = types.StringValue(srv.Updated.String())
	} else {
		data.Updated = types.StringNull()
	}

	return &data, nil
}

func (r *InfinityScheduledScalingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	state := &InfinityScheduledScalingResourceModel{}

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
			"Error Reading Infinity scheduled scaling policy",
			fmt.Sprintf("Could not read Infinity scheduled scaling policy: %s", err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *InfinityScheduledScalingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	plan := &InfinityScheduledScalingResourceModel{}
	state := &InfinityScheduledScalingResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateRequest := &config.ScheduledScalingUpdateRequest{
		PolicyName:         plan.PolicyName.ValueString(),
		PolicyType:         plan.PolicyType.ValueString(),
		ResourceIdentifier: plan.ResourceIdentifier.ValueString(),
		LocalTimezone:      plan.LocalTimezone.ValueString(),
		StartDate:          plan.StartDate.ValueString(),
		TimeFrom:           plan.TimeFrom.ValueString(),
		TimeTo:             plan.TimeTo.ValueString(),
	}

	// Handle optional pointer fields
	if !plan.Enabled.IsNull() && !plan.Enabled.IsUnknown() {
		enabled := plan.Enabled.ValueBool()
		updateRequest.Enabled = &enabled
	}

	if !plan.InstancesToAdd.IsNull() && !plan.InstancesToAdd.IsUnknown() {
		instances := int(plan.InstancesToAdd.ValueInt64())
		updateRequest.InstancesToAdd = &instances
	}

	if !plan.MinutesInAdvance.IsNull() && !plan.MinutesInAdvance.IsUnknown() {
		minutes := int(plan.MinutesInAdvance.ValueInt64())
		updateRequest.MinutesInAdvance = &minutes
	}

	if !plan.Mon.IsNull() && !plan.Mon.IsUnknown() {
		mon := plan.Mon.ValueBool()
		updateRequest.Mon = &mon
	}

	if !plan.Tue.IsNull() && !plan.Tue.IsUnknown() {
		tue := plan.Tue.ValueBool()
		updateRequest.Tue = &tue
	}

	if !plan.Wed.IsNull() && !plan.Wed.IsUnknown() {
		wed := plan.Wed.ValueBool()
		updateRequest.Wed = &wed
	}

	if !plan.Thu.IsNull() && !plan.Thu.IsUnknown() {
		thu := plan.Thu.ValueBool()
		updateRequest.Thu = &thu
	}

	if !plan.Fri.IsNull() && !plan.Fri.IsUnknown() {
		fri := plan.Fri.ValueBool()
		updateRequest.Fri = &fri
	}

	if !plan.Sat.IsNull() && !plan.Sat.IsUnknown() {
		sat := plan.Sat.ValueBool()
		updateRequest.Sat = &sat
	}

	if !plan.Sun.IsNull() && !plan.Sun.IsUnknown() {
		sun := plan.Sun.ValueBool()
		updateRequest.Sun = &sun
	}

	resourceID := int(state.ResourceID.ValueInt32())
	_, err := r.InfinityClient.Config().UpdateScheduledScaling(ctx, resourceID, updateRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Infinity scheduled scaling policy",
			fmt.Sprintf("Could not update Infinity scheduled scaling policy: %s", err),
		)
		return
	}

	// Read the state from the API to get all computed values
	model, err := r.read(ctx, resourceID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Updated Infinity scheduled scaling policy",
			fmt.Sprintf("Could not read updated Infinity scheduled scaling policy with ID %d: %s", resourceID, err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *InfinityScheduledScalingResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	state := &InfinityScheduledScalingResourceModel{}

	tflog.Info(ctx, "Deleting Infinity scheduled scaling policy")

	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.InfinityClient.Config().DeleteScheduledScaling(ctx, int(state.ResourceID.ValueInt32()))

	// Ignore 404 Not Found and Lookup errors on delete
	if err != nil && !isNotFoundError(err) && !isLookupError(err) {
		resp.Diagnostics.AddError(
			"Error Deleting Infinity scheduled scaling policy",
			fmt.Sprintf("Could not delete Infinity scheduled scaling policy with ID %s: %s", state.ID.ValueString(), err),
		)
		return
	}
}

func (r *InfinityScheduledScalingResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resourceID, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Import ID must be a valid integer for the resource ID. Got: %s", req.ID),
		)
		return
	}

	tflog.Trace(ctx, fmt.Sprintf("Importing Infinity scheduled scaling policy with resource ID: %d", resourceID))

	// Read the resource from the API
	model, err := r.read(ctx, resourceID)
	if err != nil {
		// Check if the error is a 404 (not found)
		if isNotFoundError(err) {
			resp.Diagnostics.AddError(
				"Infinity Scheduled Scaling Policy Not Found",
				fmt.Sprintf("Infinity scheduled scaling policy with resource ID %d not found.", resourceID),
			)
			return
		}
		resp.Diagnostics.AddError(
			"Error Importing Infinity Scheduled Scaling Policy",
			fmt.Sprintf("Could not import Infinity scheduled scaling policy with resource ID %d: %s", resourceID, err),
		)
		return
	}

	// Set the state from the imported resource
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}
