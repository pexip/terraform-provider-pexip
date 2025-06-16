package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/pexip/go-infinity-sdk/v38"
	"github.com/pexip/go-infinity-sdk/v38/config"
	"strconv"
	"strings"
)

var (
	_ resource.ResourceWithImportState = (*InfinityNodeResource)(nil)
)

type InfinityNodeResource struct {
	InfinityClient *infinity.Client
}

type InfinityNodeResourceModel struct {
	ID     types.Int32  `tfsdk:"id"`
	Name   types.String `tfsdk:"name"`
	Config types.String `tfsdk:"config"`
}

func (r *InfinityNodeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_infinity_node"
}

func (r *InfinityNodeResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.InfinityClient = p.InfinityClient
}

func (r *InfinityNodeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.Int32Attribute{
				Computed:            true,
				MarkdownDescription: "Resource identifier",
			},
			"name": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(3),
				},
				MarkdownDescription: "The name of the Infinity node. This should be unique within the Infinity cluster.",
			},
			"config": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(10),
				},
				MarkdownDescription: "Bootstrap configuration for the Infinity node.",
			},
		},
		MarkdownDescription: "Registers a node with the Infinity service. This resource is used to manage the lifecycle of nodes in the Infinity cluster.",
	}
}

func (r *InfinityNodeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data InfinityNodeResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createRequest := &config.WorkerVMCreateRequest{}

	// Set name if provided, otherwise use generated name
	if !data.Name.IsNull() && !data.Name.IsUnknown() {
		createRequest.Name = data.Name.ValueString()
	}

	vm, err := r.InfinityClient.Config.CreateWorkerVM(ctx, createRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Infinity Node",
			fmt.Sprintf("Could not create Infinity node: %s", err),
		)
		return
	}

	data.ID = types.Int32Value(int32(vm.ID))
	data.Name = types.StringValue(vm.Name)
	tflog.Trace(ctx, fmt.Sprintf("created Infinity node with ID: %d, name: %s", vm.ID, vm.Name))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *InfinityNodeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data InfinityNodeResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	vm, err := r.InfinityClient.Config.GetWorkerVM(ctx, int(data.ID.ValueInt32()))
	if err != nil {
		// Check if the error is a 404 (not found)
		if isNotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading Infinity Node",
			fmt.Sprintf("Could not read Infinity node with ID %d: %s", data.ID.ValueInt32(), err),
		)
		return
	}

	data.Name = types.StringValue(vm.Name)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *InfinityNodeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data InfinityNodeResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateRequest := &config.WorkerVMUpdateRequest{}

	// Set name if provided
	if !data.Name.IsNull() && !data.Name.IsUnknown() {
		updateRequest.Name = data.Name.ValueString()
	}

	vm, err := r.InfinityClient.Config.UpdateWorkerVM(ctx, int(data.ID.ValueInt32()), updateRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Infinity Node",
			fmt.Sprintf("Could not update Infinity node with ID %d: %s", data.ID.ValueInt32(), err),
		)
		return
	}

	data.Name = types.StringValue(vm.Name)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *InfinityNodeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data InfinityNodeResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.InfinityClient.Config.DeleteWorkerVM(ctx, int(data.ID.ValueInt32()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Infinity Node",
			fmt.Sprintf("Could not delete Infinity node with ID %d: %s", data.ID.ValueInt32(), err),
		)
		return
	}
}

func (r *InfinityNodeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Validate that the ID is a valid integer
	id, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Import ID must be a valid integer, got: %s", req.ID),
		)
		return
	}

	if id <= 0 {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Import ID must be a positive integer, got: %d", id),
		)
		return
	}

	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// isNotFoundError checks if the error indicates a 404/not found response
func isNotFoundError(err error) bool {
	// This is a placeholder - you'll need to check the actual error types
	// returned by the go-infinity-sdk to determine what constitutes a "not found" error
	return strings.Contains(err.Error(), "404") ||
		strings.Contains(err.Error(), "not found") ||
		strings.Contains(err.Error(), "Not Found")
}
