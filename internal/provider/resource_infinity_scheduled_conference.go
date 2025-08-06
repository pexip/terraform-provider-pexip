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
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/pexip/go-infinity-sdk/v38/config"

	"github.com/pexip/go-infinity-sdk/v38/util"
)

var (
	_ resource.ResourceWithImportState = (*InfinityScheduledConferenceResource)(nil)
)

type InfinityScheduledConferenceResource struct {
	InfinityClient InfinityClient
}

type InfinityScheduledConferenceResourceModel struct {
	ID                  types.String `tfsdk:"id"`
	ResourceID          types.Int32  `tfsdk:"resource_id"`
	Conference          types.String `tfsdk:"conference"`
	StartTime           types.String `tfsdk:"start_time"`
	EndTime             types.String `tfsdk:"end_time"`
	Subject             types.String `tfsdk:"subject"`
	EWSItemID           types.String `tfsdk:"ews_item_id"`
	EWSItemUID          types.String `tfsdk:"ews_item_uid"`
	RecurringConference types.String `tfsdk:"recurring_conference"`
	ScheduledAlias      types.String `tfsdk:"scheduled_alias"`
}

func (r *InfinityScheduledConferenceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_infinity_scheduled_conference"
}

func (r *InfinityScheduledConferenceResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *InfinityScheduledConferenceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Resource URI for the scheduled conference in Infinity",
			},
			"resource_id": schema.Int32Attribute{
				Computed:            true,
				MarkdownDescription: "The resource integer identifier for the scheduled conference in Infinity",
			},
			"conference": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(250),
				},
				MarkdownDescription: "The conference URI or reference. Maximum length: 250 characters.",
			},
			"start_time": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
				MarkdownDescription: "The start time of the scheduled conference in ISO 8601 format (e.g., '2024-01-01T10:00:00Z').",
			},
			"end_time": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
				MarkdownDescription: "The end time of the scheduled conference in ISO 8601 format (e.g., '2024-01-01T11:00:00Z').",
			},
			"subject": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(250),
				},
				MarkdownDescription: "The subject of the scheduled conference. Maximum length: 250 characters.",
			},
			"ews_item_id": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
				MarkdownDescription: "The Exchange Web Services (EWS) item ID for the conference.",
			},
			"ews_item_uid": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(250),
				},
				MarkdownDescription: "The Exchange Web Services (EWS) item UID for the conference. Maximum length: 250 characters.",
			},
			"recurring_conference": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Reference to recurring conference resource URI.",
			},
			"scheduled_alias": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Reference to scheduled alias resource URI.",
			},
		},
		MarkdownDescription: "Manages a scheduled conference configuration with the Infinity service.",
	}
}

func (r *InfinityScheduledConferenceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	plan := &InfinityScheduledConferenceResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Parse start and end times
	startTime, err := time.Parse(time.RFC3339, plan.StartTime.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Start Time Format",
			fmt.Sprintf("Could not parse start time '%s'. Expected ISO 8601 format (e.g., '2024-01-01T10:00:00Z'): %s", plan.StartTime.ValueString(), err),
		)
		return
	}

	endTime, err := time.Parse(time.RFC3339, plan.EndTime.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid End Time Format",
			fmt.Sprintf("Could not parse end time '%s'. Expected ISO 8601 format (e.g., '2024-01-01T11:00:00Z'): %s", plan.EndTime.ValueString(), err),
		)
		return
	}

	createRequest := &config.ScheduledConferenceCreateRequest{
		Conference: plan.Conference.ValueString(),
		StartTime:  util.InfinityTime{Time: startTime},
		EndTime:    util.InfinityTime{Time: endTime},
		EWSItemID:  plan.EWSItemID.ValueString(),
	}

	// Set optional fields
	if !plan.Subject.IsNull() {
		createRequest.Subject = plan.Subject.ValueString()
	}
	if !plan.EWSItemUID.IsNull() {
		createRequest.EWSItemUID = plan.EWSItemUID.ValueString()
	}
	if !plan.RecurringConference.IsNull() {
		recurringConference := plan.RecurringConference.ValueString()
		createRequest.RecurringConference = &recurringConference
	}
	if !plan.ScheduledAlias.IsNull() {
		scheduledAlias := plan.ScheduledAlias.ValueString()
		createRequest.ScheduledAlias = &scheduledAlias
	}

	createResponse, err := r.InfinityClient.Config().CreateScheduledConference(ctx, createRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Infinity scheduled conference",
			fmt.Sprintf("Could not create Infinity scheduled conference: %s", err),
		)
		return
	}

	resourceID, err := createResponse.ResourceID()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Retrieving Infinity scheduled conference ID",
			fmt.Sprintf("Could not retrieve ID for created Infinity scheduled conference: %s", err),
		)
		return
	}

	// Read the state from the API to get all computed values
	model, err := r.read(ctx, resourceID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Created Infinity scheduled conference",
			fmt.Sprintf("Could not read created Infinity scheduled conference with ID %d: %s", resourceID, err),
		)
		return
	}
	tflog.Trace(ctx, fmt.Sprintf("created Infinity scheduled conference with ID: %s, conference: %s", model.ID, model.Conference))

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *InfinityScheduledConferenceResource) read(ctx context.Context, resourceID int) (*InfinityScheduledConferenceResourceModel, error) {
	var data InfinityScheduledConferenceResourceModel

	srv, err := r.InfinityClient.Config().GetScheduledConference(ctx, resourceID)
	if err != nil {
		return nil, err
	}

	if srv.ResourceURI == "" {
		return nil, fmt.Errorf("scheduled conference with ID %d not found", resourceID)
	}

	data.ID = types.StringValue(srv.ResourceURI)
	data.ResourceID = types.Int32Value(int32(resourceID)) // #nosec G115 -- API values are expected to be within int32 range
	data.Conference = types.StringValue(srv.Conference)
	data.StartTime = types.StringValue(srv.StartTime.Format(time.RFC3339))
	data.EndTime = types.StringValue(srv.EndTime.Format(time.RFC3339))
	data.Subject = types.StringValue(srv.Subject)
	data.EWSItemID = types.StringValue(srv.EWSItemID)
	data.EWSItemUID = types.StringValue(srv.EWSItemUID)
	if srv.RecurringConference != nil {
		data.RecurringConference = types.StringValue(*srv.RecurringConference)
	} else {
		data.RecurringConference = types.StringNull()
	}
	if srv.ScheduledAlias != nil {
		data.ScheduledAlias = types.StringValue(*srv.ScheduledAlias)
	} else {
		data.ScheduledAlias = types.StringNull()
	}

	return &data, nil
}

func (r *InfinityScheduledConferenceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	state := &InfinityScheduledConferenceResourceModel{}

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
			"Error Reading Infinity scheduled conference",
			fmt.Sprintf("Could not read Infinity scheduled conference: %s", err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *InfinityScheduledConferenceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	plan := &InfinityScheduledConferenceResourceModel{}
	state := &InfinityScheduledConferenceResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resourceID := int(state.ResourceID.ValueInt32())

	updateRequest := &config.ScheduledConferenceUpdateRequest{
		Conference: plan.Conference.ValueString(),
		EWSItemID:  plan.EWSItemID.ValueString(),
	}

	// Parse and set start time
	if !plan.StartTime.IsNull() {
		startTime, err := time.Parse(time.RFC3339, plan.StartTime.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Invalid Start Time Format",
				fmt.Sprintf("Could not parse start time '%s'. Expected ISO 8601 format: %s", plan.StartTime.ValueString(), err),
			)
			return
		}
		startTimeUtil := util.InfinityTime{Time: startTime}
		updateRequest.StartTime = &startTimeUtil
	}

	// Parse and set end time
	if !plan.EndTime.IsNull() {
		endTime, err := time.Parse(time.RFC3339, plan.EndTime.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Invalid End Time Format",
				fmt.Sprintf("Could not parse end time '%s'. Expected ISO 8601 format: %s", plan.EndTime.ValueString(), err),
			)
			return
		}
		endTimeUtil := util.InfinityTime{Time: endTime}
		updateRequest.EndTime = &endTimeUtil
	}

	// Set optional fields
	if !plan.Subject.IsNull() {
		updateRequest.Subject = plan.Subject.ValueString()
	}
	if !plan.EWSItemUID.IsNull() {
		updateRequest.EWSItemUID = plan.EWSItemUID.ValueString()
	}
	if !plan.RecurringConference.IsNull() {
		recurringConference := plan.RecurringConference.ValueString()
		updateRequest.RecurringConference = &recurringConference
	}
	if !plan.ScheduledAlias.IsNull() {
		scheduledAlias := plan.ScheduledAlias.ValueString()
		updateRequest.ScheduledAlias = &scheduledAlias
	}

	_, err := r.InfinityClient.Config().UpdateScheduledConference(ctx, resourceID, updateRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Infinity scheduled conference",
			fmt.Sprintf("Could not update Infinity scheduled conference with ID %d: %s", resourceID, err),
		)
		return
	}

	// Re-read the resource to get the latest state
	updatedModel, err := r.read(ctx, resourceID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Updated Infinity scheduled conference",
			fmt.Sprintf("Could not read updated Infinity scheduled conference with ID %d: %s", resourceID, err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, updatedModel)...)
}

func (r *InfinityScheduledConferenceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	state := &InfinityScheduledConferenceResourceModel{}

	tflog.Info(ctx, "Deleting Infinity scheduled conference")

	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.InfinityClient.Config().DeleteScheduledConference(ctx, int(state.ResourceID.ValueInt32()))

	// Ignore 404 Not Found and Lookup errors on delete
	if err != nil && !isNotFoundError(err) && !isLookupError(err) {
		resp.Diagnostics.AddError(
			"Error Deleting Infinity scheduled conference",
			fmt.Sprintf("Could not delete Infinity scheduled conference with ID %s: %s", state.ID.ValueString(), err),
		)
		return
	}
}

func (r *InfinityScheduledConferenceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resourceID, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Import ID must be a valid integer for the resource ID. Got: %s", req.ID),
		)
		return
	}

	tflog.Trace(ctx, fmt.Sprintf("Importing Infinity scheduled conference with resource ID: %d", resourceID))

	// Read the resource from the API
	model, err := r.read(ctx, resourceID)
	if err != nil {
		// Check if the error is a 404 (not found)
		if isNotFoundError(err) {
			resp.Diagnostics.AddError(
				"Infinity Scheduled Conference Not Found",
				fmt.Sprintf("Infinity scheduled conference with resource ID %d not found.", resourceID),
			)
			return
		}
		resp.Diagnostics.AddError(
			"Error Importing Infinity Scheduled Conference",
			fmt.Sprintf("Could not import Infinity scheduled conference with resource ID %d: %s", resourceID, err),
		)
		return
	}

	// Set the state from the imported resource
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}
