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
	_ resource.ResourceWithImportState = (*InfinityRecurringConferenceResource)(nil)
)

type InfinityRecurringConferenceResource struct {
	InfinityClient InfinityClient
}

type InfinityRecurringConferenceResourceModel struct {
	ID             types.String `tfsdk:"id"`
	ResourceID     types.Int32  `tfsdk:"resource_id"`
	Conference     types.String `tfsdk:"conference"`
	CurrentIndex   types.Int64  `tfsdk:"current_index"`
	EWSItemID      types.String `tfsdk:"ews_item_id"`
	IsDepleted     types.Bool   `tfsdk:"is_depleted"`
	Subject        types.String `tfsdk:"subject"`
	ScheduledAlias types.String `tfsdk:"scheduled_alias"`
}

func (r *InfinityRecurringConferenceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_infinity_recurring_conference"
}

func (r *InfinityRecurringConferenceResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *InfinityRecurringConferenceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Resource URI for the recurring conference in Infinity",
			},
			"resource_id": schema.Int32Attribute{
				Computed:            true,
				MarkdownDescription: "The resource integer identifier for the recurring conference in Infinity",
			},
			"conference": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
				MarkdownDescription: "The conference identifier or URI associated with this recurring conference.",
			},
			"current_index": schema.Int64Attribute{
				Required: true,
				Validators: []validator.Int64{
					int64validator.AtLeast(0),
				},
				MarkdownDescription: "The current index of the recurring conference series.",
			},
			"ews_item_id": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
				MarkdownDescription: "The Exchange Web Services (EWS) item identifier for this recurring conference.",
			},
			"is_depleted": schema.BoolAttribute{
				Required:            true,
				MarkdownDescription: "Whether the recurring conference series is depleted (no more occurrences).",
			},
			"subject": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(500),
				},
				MarkdownDescription: "The subject or title of the recurring conference. Maximum length: 500 characters.",
			},
			"scheduled_alias": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The scheduled alias for the recurring conference.",
			},
		},
		MarkdownDescription: "Manages a recurring conference configuration with the Infinity service. Recurring conferences are used for scheduled conference series.",
	}
}

func (r *InfinityRecurringConferenceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	plan := &InfinityRecurringConferenceResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createRequest := &config.RecurringConferenceCreateRequest{
		Conference:   plan.Conference.ValueString(),
		CurrentIndex: int(plan.CurrentIndex.ValueInt64()),
		EWSItemID:    plan.EWSItemID.ValueString(),
		IsDepleted:   plan.IsDepleted.ValueBool(),
		Subject:      plan.Subject.ValueString(),
	}

	// Handle optional scheduled_alias field
	if !plan.ScheduledAlias.IsNull() && !plan.ScheduledAlias.IsUnknown() {
		alias := plan.ScheduledAlias.ValueString()
		createRequest.ScheduledAlias = &alias
	}

	createResponse, err := r.InfinityClient.Config().CreateRecurringConference(ctx, createRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Infinity recurring conference",
			fmt.Sprintf("Could not create Infinity recurring conference: %s", err),
		)
		return
	}

	resourceID, err := createResponse.ResourceID()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Retrieving Infinity recurring conference ID",
			fmt.Sprintf("Could not retrieve ID for created Infinity recurring conference: %s", err),
		)
		return
	}

	// Read the state from the API to get all computed values
	model, err := r.read(ctx, resourceID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Created Infinity recurring conference",
			fmt.Sprintf("Could not read created Infinity recurring conference with ID %d: %s", resourceID, err),
		)
		return
	}
	tflog.Trace(ctx, fmt.Sprintf("created Infinity recurring conference with ID: %s, conference: %s", model.ID, model.Conference))

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *InfinityRecurringConferenceResource) read(ctx context.Context, resourceID int) (*InfinityRecurringConferenceResourceModel, error) {
	var data InfinityRecurringConferenceResourceModel

	srv, err := r.InfinityClient.Config().GetRecurringConference(ctx, resourceID)
	if err != nil {
		return nil, err
	}

	if srv.ResourceURI == "" {
		return nil, fmt.Errorf("recurring conference with ID %d not found", resourceID)
	}

	data.ID = types.StringValue(srv.ResourceURI)
	data.ResourceID = types.Int32Value(int32(resourceID)) // #nosec G115 -- API values are expected to be within int32 range
	data.Conference = types.StringValue(srv.Conference)
	data.CurrentIndex = types.Int64Value(int64(srv.CurrentIndex))
	data.EWSItemID = types.StringValue(srv.EWSItemID)
	data.IsDepleted = types.BoolValue(srv.IsDepleted)
	data.Subject = types.StringValue(srv.Subject)

	// Handle optional scheduled_alias field
	if srv.ScheduledAlias != nil {
		data.ScheduledAlias = types.StringValue(*srv.ScheduledAlias)
	} else {
		data.ScheduledAlias = types.StringNull()
	}

	return &data, nil
}

func (r *InfinityRecurringConferenceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	state := &InfinityRecurringConferenceResourceModel{}

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
			"Error Reading Infinity recurring conference",
			fmt.Sprintf("Could not read Infinity recurring conference: %s", err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *InfinityRecurringConferenceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	plan := &InfinityRecurringConferenceResourceModel{}
	state := &InfinityRecurringConferenceResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resourceID := int(state.ResourceID.ValueInt32())

	updateRequest := &config.RecurringConferenceUpdateRequest{
		Conference: plan.Conference.ValueString(),
		EWSItemID:  plan.EWSItemID.ValueString(),
		Subject:    plan.Subject.ValueString(),
	}

	if !plan.CurrentIndex.IsNull() {
		index := int(plan.CurrentIndex.ValueInt64())
		updateRequest.CurrentIndex = &index
	}

	if !plan.IsDepleted.IsNull() {
		depleted := plan.IsDepleted.ValueBool()
		updateRequest.IsDepleted = &depleted
	}

	// Handle optional scheduled_alias field
	if !plan.ScheduledAlias.IsNull() && !plan.ScheduledAlias.IsUnknown() {
		alias := plan.ScheduledAlias.ValueString()
		updateRequest.ScheduledAlias = &alias
	}

	_, err := r.InfinityClient.Config().UpdateRecurringConference(ctx, resourceID, updateRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Infinity recurring conference",
			fmt.Sprintf("Could not update Infinity recurring conference with ID %d: %s", resourceID, err),
		)
		return
	}

	// Re-read the resource to get the latest state
	updatedModel, err := r.read(ctx, resourceID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Updated Infinity recurring conference",
			fmt.Sprintf("Could not read updated Infinity recurring conference with ID %d: %s", resourceID, err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, updatedModel)...)
}

func (r *InfinityRecurringConferenceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	state := &InfinityRecurringConferenceResourceModel{}

	tflog.Info(ctx, "Deleting Infinity recurring conference")

	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.InfinityClient.Config().DeleteRecurringConference(ctx, int(state.ResourceID.ValueInt32()))

	// Ignore 404 Not Found and Lookup errors on delete
	if err != nil && !isNotFoundError(err) && !isLookupError(err) {
		resp.Diagnostics.AddError(
			"Error Deleting Infinity recurring conference",
			fmt.Sprintf("Could not delete Infinity recurring conference with ID %s: %s", state.ID.ValueString(), err),
		)
		return
	}
}

func (r *InfinityRecurringConferenceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resourceID, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Import ID must be a valid integer for the resource ID. Got: %s", req.ID),
		)
		return
	}

	tflog.Trace(ctx, fmt.Sprintf("Importing Infinity recurring conference with resource ID: %d", resourceID))

	// Read the resource from the API
	model, err := r.read(ctx, resourceID)
	if err != nil {
		// Check if the error is a 404 (not found)
		if isNotFoundError(err) {
			resp.Diagnostics.AddError(
				"Infinity Recurring Conference Not Found",
				fmt.Sprintf("Infinity recurring conference with resource ID %d not found.", resourceID),
			)
			return
		}
		resp.Diagnostics.AddError(
			"Error Importing Infinity Recurring Conference",
			fmt.Sprintf("Could not import Infinity recurring conference with resource ID %d: %s", resourceID, err),
		)
		return
	}

	// Set the state from the imported resource
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}
