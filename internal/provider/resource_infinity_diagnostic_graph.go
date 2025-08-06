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
	_ resource.ResourceWithImportState = (*InfinityDiagnosticGraphResource)(nil)
)

type InfinityDiagnosticGraphResource struct {
	InfinityClient InfinityClient
}

type InfinityDiagnosticGraphResourceModel struct {
	ID         types.String `tfsdk:"id"`
	ResourceID types.Int32  `tfsdk:"resource_id"`
	Title      types.String `tfsdk:"title"`
	Order      types.Int64  `tfsdk:"order"`
	Datasets   types.Set    `tfsdk:"datasets"`
}

func (r *InfinityDiagnosticGraphResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_infinity_diagnostic_graph"
}

func (r *InfinityDiagnosticGraphResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *InfinityDiagnosticGraphResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Resource URI for the diagnostic graph in Infinity",
			},
			"resource_id": schema.Int32Attribute{
				Computed:            true,
				MarkdownDescription: "The resource integer identifier for the diagnostic graph in Infinity",
			},
			"title": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					stringvalidator.LengthAtMost(250),
				},
				MarkdownDescription: "The title of the diagnostic graph. Maximum length: 250 characters.",
			},
			"order": schema.Int64Attribute{
				Required: true,
				Validators: []validator.Int64{
					int64validator.AtLeast(0),
				},
				MarkdownDescription: "The display order of the diagnostic graph. Lower values appear first.",
			},
			"datasets": schema.SetAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				MarkdownDescription: "List of dataset identifiers to include in this diagnostic graph.",
			},
		},
		MarkdownDescription: "Manages a diagnostic graph with the Infinity service. Diagnostic graphs provide visual monitoring and troubleshooting capabilities for system health and performance metrics.",
	}
}

func (r *InfinityDiagnosticGraphResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	plan := &InfinityDiagnosticGraphResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createRequest := &config.DiagnosticGraphCreateRequest{
		Title: plan.Title.ValueString(),
		Order: int(plan.Order.ValueInt64()),
	}

	// Handle list field
	if !plan.Datasets.IsNull() {
		var datasets []string
		resp.Diagnostics.Append(plan.Datasets.ElementsAs(ctx, &datasets, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		createRequest.Datasets = datasets
	}

	createResponse, err := r.InfinityClient.Config().CreateDiagnosticGraph(ctx, createRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Infinity diagnostic graph",
			fmt.Sprintf("Could not create Infinity diagnostic graph: %s", err),
		)
		return
	}

	resourceID, err := createResponse.ResourceID()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Retrieving Infinity diagnostic graph ID",
			fmt.Sprintf("Could not retrieve ID for created Infinity diagnostic graph: %s", err),
		)
		return
	}

	// Read the state from the API to get all computed values
	model, err := r.read(ctx, resourceID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Created Infinity diagnostic graph",
			fmt.Sprintf("Could not read created Infinity diagnostic graph with ID %d: %s", resourceID, err),
		)
		return
	}
	tflog.Trace(ctx, fmt.Sprintf("created Infinity diagnostic graph with ID: %s, title: %s", model.ID, model.Title))

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *InfinityDiagnosticGraphResource) read(ctx context.Context, resourceID int) (*InfinityDiagnosticGraphResourceModel, error) {
	var data InfinityDiagnosticGraphResourceModel

	srv, err := r.InfinityClient.Config().GetDiagnosticGraph(ctx, resourceID)
	if err != nil {
		return nil, err
	}

	if srv.ResourceURI == "" {
		return nil, fmt.Errorf("diagnostic graph with ID %d not found", resourceID)
	}

	data.ID = types.StringValue(srv.ResourceURI)
	data.ResourceID = types.Int32Value(int32(resourceID)) // #nosec G115 -- API values are expected to be within int32 range
	data.Title = types.StringValue(srv.Title)
	data.Order = types.Int64Value(int64(srv.Order))

	// Handle list field
	if srv.Datasets != nil {
		datasetsSet, diags := types.SetValueFrom(ctx, types.StringType, srv.Datasets)
		if diags.HasError() {
			return nil, fmt.Errorf("failed to convert datasets: %s", diags.Errors())
		}
		data.Datasets = datasetsSet
	} else {
		data.Datasets = types.SetNull(types.StringType)
	}

	return &data, nil
}

func (r *InfinityDiagnosticGraphResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	state := &InfinityDiagnosticGraphResourceModel{}

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
			"Error Reading Infinity diagnostic graph",
			fmt.Sprintf("Could not read Infinity diagnostic graph: %s", err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *InfinityDiagnosticGraphResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	plan := &InfinityDiagnosticGraphResourceModel{}
	state := &InfinityDiagnosticGraphResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateRequest := &config.DiagnosticGraphUpdateRequest{
		Title: plan.Title.ValueString(),
	}

	// Handle optional pointer field for order
	if !plan.Order.IsNull() && !plan.Order.IsUnknown() {
		order := int(plan.Order.ValueInt64())
		updateRequest.Order = &order
	}

	// Handle list field
	if !plan.Datasets.IsNull() {
		var datasets []string
		resp.Diagnostics.Append(plan.Datasets.ElementsAs(ctx, &datasets, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		updateRequest.Datasets = datasets
	}

	resourceID := int(state.ResourceID.ValueInt32())
	_, err := r.InfinityClient.Config().UpdateDiagnosticGraph(ctx, resourceID, updateRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Infinity diagnostic graph",
			fmt.Sprintf("Could not update Infinity diagnostic graph: %s", err),
		)
		return
	}

	// Read the state from the API to get all computed values
	model, err := r.read(ctx, resourceID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Updated Infinity diagnostic graph",
			fmt.Sprintf("Could not read updated Infinity diagnostic graph with ID %d: %s", resourceID, err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *InfinityDiagnosticGraphResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	state := &InfinityDiagnosticGraphResourceModel{}

	tflog.Info(ctx, "Deleting Infinity diagnostic graph")

	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.InfinityClient.Config().DeleteDiagnosticGraph(ctx, int(state.ResourceID.ValueInt32()))

	// Ignore 404 Not Found and Lookup errors on delete
	if err != nil && !isNotFoundError(err) && !isLookupError(err) {
		resp.Diagnostics.AddError(
			"Error Deleting Infinity diagnostic graph",
			fmt.Sprintf("Could not delete Infinity diagnostic graph with ID %s: %s", state.ID.ValueString(), err),
		)
		return
	}
}

func (r *InfinityDiagnosticGraphResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resourceID, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Import ID must be a valid integer for the resource ID. Got: %s", req.ID),
		)
		return
	}

	tflog.Trace(ctx, fmt.Sprintf("Importing Infinity diagnostic graph with resource ID: %d", resourceID))

	// Read the resource from the API
	model, err := r.read(ctx, resourceID)
	if err != nil {
		// Check if the error is a 404 (not found)
		if isNotFoundError(err) {
			resp.Diagnostics.AddError(
				"Infinity Diagnostic Graph Not Found",
				fmt.Sprintf("Infinity diagnostic graph with resource ID %d not found.", resourceID),
			)
			return
		}
		resp.Diagnostics.AddError(
			"Error Importing Infinity Diagnostic Graph",
			fmt.Sprintf("Could not import Infinity diagnostic graph with resource ID %d: %s", resourceID, err),
		)
		return
	}

	// Set the state from the imported resource
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}
