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
	_ resource.ResourceWithImportState = (*InfinityScheduledAliasResource)(nil)
)

type InfinityScheduledAliasResource struct {
	InfinityClient InfinityClient
}

type InfinityScheduledAliasResourceModel struct {
	ID                     types.String `tfsdk:"id"`
	ResourceID             types.Int32  `tfsdk:"resource_id"`
	Alias                  types.String `tfsdk:"alias"`
	AliasNumber            types.Int64  `tfsdk:"alias_number"`
	NumericAlias           types.String `tfsdk:"numeric_alias"`
	UUID                   types.String `tfsdk:"uuid"`
	ExchangeConnector      types.String `tfsdk:"exchange_connector"`
	IsUsed                 types.Bool   `tfsdk:"is_used"`
	EWSItemUID             types.String `tfsdk:"ews_item_uid"`
	CreationTime           types.String `tfsdk:"creation_time"`
	ConferenceDeletionTime types.String `tfsdk:"conference_deletion_time"`
}

func (r *InfinityScheduledAliasResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_infinity_scheduled_alias"
}

func (r *InfinityScheduledAliasResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *InfinityScheduledAliasResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Resource URI for the scheduled alias in Infinity",
			},
			"resource_id": schema.Int32Attribute{
				Computed:            true,
				MarkdownDescription: "The resource integer identifier for the scheduled alias in Infinity",
			},
			"alias": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					stringvalidator.LengthAtMost(250),
				},
				MarkdownDescription: "The alias name for scheduled conferences. Maximum length: 250 characters.",
			},
			"alias_number": schema.Int64Attribute{
				Required: true,
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
				},
				MarkdownDescription: "The numeric identifier for this alias.",
			},
			"numeric_alias": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					stringvalidator.LengthAtMost(50),
				},
				MarkdownDescription: "The numeric alias string for dial-in access. Maximum length: 50 characters.",
			},
			"uuid": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
				MarkdownDescription: "The UUID for this scheduled alias.",
			},
			"exchange_connector": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
				MarkdownDescription: "The Exchange connector URI associated with this scheduled alias.",
			},
			"is_used": schema.BoolAttribute{
				Required:            true,
				MarkdownDescription: "Whether this scheduled alias is currently in use.",
			},
			"ews_item_uid": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(200),
				},
				MarkdownDescription: "The Exchange Web Services item UID. Maximum length: 200 characters.",
			},
			"creation_time": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The timestamp when this scheduled alias was created.",
			},
			"conference_deletion_time": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The timestamp when the associated conference was deleted, if applicable.",
			},
		},
		MarkdownDescription: "Manages a scheduled alias with the Infinity service. Scheduled aliases are used for Microsoft Exchange integration to provide consistent conference aliases for scheduled meetings and calendar integration.",
	}
}

func (r *InfinityScheduledAliasResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	plan := &InfinityScheduledAliasResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createRequest := &config.ScheduledAliasCreateRequest{
		Alias:             plan.Alias.ValueString(),
		AliasNumber:       int(plan.AliasNumber.ValueInt64()),
		NumericAlias:      plan.NumericAlias.ValueString(),
		UUID:              plan.UUID.ValueString(),
		ExchangeConnector: plan.ExchangeConnector.ValueString(),
		IsUsed:            plan.IsUsed.ValueBool(),
	}

	// Handle optional pointer field
	if !plan.EWSItemUID.IsNull() && !plan.EWSItemUID.IsUnknown() {
		ewsUID := plan.EWSItemUID.ValueString()
		createRequest.EWSItemUID = &ewsUID
	}

	createResponse, err := r.InfinityClient.Config().CreateScheduledAlias(ctx, createRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Infinity scheduled alias",
			fmt.Sprintf("Could not create Infinity scheduled alias: %s", err),
		)
		return
	}

	resourceID, err := createResponse.ResourceID()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Retrieving Infinity scheduled alias ID",
			fmt.Sprintf("Could not retrieve ID for created Infinity scheduled alias: %s", err),
		)
		return
	}

	// Read the state from the API to get all computed values
	model, err := r.read(ctx, resourceID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Created Infinity scheduled alias",
			fmt.Sprintf("Could not read created Infinity scheduled alias with ID %d: %s", resourceID, err),
		)
		return
	}
	tflog.Trace(ctx, fmt.Sprintf("created Infinity scheduled alias with ID: %s, alias: %s", model.ID, model.Alias))

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *InfinityScheduledAliasResource) read(ctx context.Context, resourceID int) (*InfinityScheduledAliasResourceModel, error) {
	var data InfinityScheduledAliasResourceModel

	srv, err := r.InfinityClient.Config().GetScheduledAlias(ctx, resourceID)
	if err != nil {
		return nil, err
	}

	if srv.ResourceURI == "" {
		return nil, fmt.Errorf("scheduled alias with ID %d not found", resourceID)
	}

	data.ID = types.StringValue(srv.ResourceURI)
	data.ResourceID = types.Int32Value(int32(resourceID)) // #nosec G115 -- API values are expected to be within int32 range
	data.Alias = types.StringValue(srv.Alias)
	data.AliasNumber = types.Int64Value(int64(srv.AliasNumber))
	data.NumericAlias = types.StringValue(srv.NumericAlias)
	data.UUID = types.StringValue(srv.UUID)
	data.ExchangeConnector = types.StringValue(srv.ExchangeConnector)
	data.IsUsed = types.BoolValue(srv.IsUsed)
	data.CreationTime = types.StringValue(srv.CreationTime.String())

	// Handle optional pointer fields
	if srv.EWSItemUID != nil {
		data.EWSItemUID = types.StringValue(*srv.EWSItemUID)
	} else {
		data.EWSItemUID = types.StringNull()
	}

	if srv.ConferenceDeletionTime != nil {
		data.ConferenceDeletionTime = types.StringValue(srv.ConferenceDeletionTime.String())
	} else {
		data.ConferenceDeletionTime = types.StringNull()
	}

	return &data, nil
}

func (r *InfinityScheduledAliasResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	state := &InfinityScheduledAliasResourceModel{}

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
			"Error Reading Infinity scheduled alias",
			fmt.Sprintf("Could not read Infinity scheduled alias: %s", err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *InfinityScheduledAliasResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	plan := &InfinityScheduledAliasResourceModel{}
	state := &InfinityScheduledAliasResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateRequest := &config.ScheduledAliasUpdateRequest{
		Alias:             plan.Alias.ValueString(),
		NumericAlias:      plan.NumericAlias.ValueString(),
		UUID:              plan.UUID.ValueString(),
		ExchangeConnector: plan.ExchangeConnector.ValueString(),
	}

	// Handle optional pointer fields
	if !plan.AliasNumber.IsNull() && !plan.AliasNumber.IsUnknown() {
		aliasNumber := int(plan.AliasNumber.ValueInt64())
		updateRequest.AliasNumber = &aliasNumber
	}

	if !plan.IsUsed.IsNull() && !plan.IsUsed.IsUnknown() {
		isUsed := plan.IsUsed.ValueBool()
		updateRequest.IsUsed = &isUsed
	}

	if !plan.EWSItemUID.IsNull() && !plan.EWSItemUID.IsUnknown() {
		ewsUID := plan.EWSItemUID.ValueString()
		updateRequest.EWSItemUID = &ewsUID
	}

	resourceID := int(state.ResourceID.ValueInt32())
	_, err := r.InfinityClient.Config().UpdateScheduledAlias(ctx, resourceID, updateRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Infinity scheduled alias",
			fmt.Sprintf("Could not update Infinity scheduled alias: %s", err),
		)
		return
	}

	// Read the state from the API to get all computed values
	model, err := r.read(ctx, resourceID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Updated Infinity scheduled alias",
			fmt.Sprintf("Could not read updated Infinity scheduled alias with ID %d: %s", resourceID, err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *InfinityScheduledAliasResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	state := &InfinityScheduledAliasResourceModel{}

	tflog.Info(ctx, "Deleting Infinity scheduled alias")

	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.InfinityClient.Config().DeleteScheduledAlias(ctx, int(state.ResourceID.ValueInt32()))

	// Ignore 404 Not Found and Lookup errors on delete
	if err != nil && !isNotFoundError(err) && !isLookupError(err) {
		resp.Diagnostics.AddError(
			"Error Deleting Infinity scheduled alias",
			fmt.Sprintf("Could not delete Infinity scheduled alias with ID %s: %s", state.ID.ValueString(), err),
		)
		return
	}
}

func (r *InfinityScheduledAliasResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resourceID, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Import ID must be a valid integer for the resource ID. Got: %s", req.ID),
		)
		return
	}

	tflog.Trace(ctx, fmt.Sprintf("Importing Infinity scheduled alias with resource ID: %d", resourceID))

	// Read the resource from the API
	model, err := r.read(ctx, resourceID)
	if err != nil {
		// Check if the error is a 404 (not found)
		if isNotFoundError(err) {
			resp.Diagnostics.AddError(
				"Infinity Scheduled Alias Not Found",
				fmt.Sprintf("Infinity scheduled alias with resource ID %d not found.", resourceID),
			)
			return
		}
		resp.Diagnostics.AddError(
			"Error Importing Infinity Scheduled Alias",
			fmt.Sprintf("Could not import Infinity scheduled alias with resource ID %d: %s", resourceID, err),
		)
		return
	}

	// Set the state from the imported resource
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}
