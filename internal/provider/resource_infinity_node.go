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
	"sync"
)

var (
	_ resource.ResourceWithImportState = (*InfinityNodeResource)(nil)
)

type InfinityNodeResource struct {
	Mutex          *sync.Mutex
	InfinityClient *infinity.Client
}

type InfinityNodeResourceModel struct {
	ID   types.Int32  `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
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
	r.Mutex = p.Mutex
}

func (r *InfinityNodeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.Int32Attribute{
				Computed:            true,
				MarkdownDescription: "Resource identifier",
			},
			"name": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(3),
				},
				MarkdownDescription: "The name of the Infinity node. This should be unique within the Infinity cluster.",
			},
		},
		MarkdownDescription: "Registers a node with the Infinity service. This resource is used to manage the lifecycle of nodes in the Infinity cluster.",
	}

	//TODO implement me
	panic("implement me")
}

func (r *InfinityNodeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data InfinityNodeResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	vm, err := r.InfinityClient.Config.CreateWorkerVM(ctx, &config.WorkerVMCreateRequest{
		Name: data.Name.ValueString(), // TODO: Add additional configuration options
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Infinity node",
			fmt.Sprintf("Could not create Infinity node: %s", err),
		)
		return
	}

	data.ID = types.Int32Value(int32(vm.ID))
	tflog.Trace(ctx, fmt.Sprintf("created Infinity node with name: %s", data.Name.ValueString()))

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
		resp.Diagnostics.AddError(
			"Error reading Infinity node",
			fmt.Sprintf("Could not read Infinity node with ID %d: %s", data.ID.ValueInt32(), err),
		)
		return
	}

	// TODO: update data and save to state
	data.Name = types.StringValue(vm.Name)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *InfinityNodeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data InfinityNodeResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	vm, err := r.InfinityClient.Config.UpdateWorkerVM(ctx, int(data.ID.ValueInt32()), &config.WorkerVMUpdateRequest{
		Name: data.Name.ValueString(), // TODO: Add additional configuration options
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating Infinity node",
			fmt.Sprintf("Could not update Infinity node with ID %d: %s", data.ID.ValueInt32(), err),
		)
		return
	}

	// TODO: update data and save to state
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
			"Error deleting Infinity node",
			fmt.Sprintf("Could not delete Infinity node with ID %d: %s", data.ID.ValueInt32(), err),
		)
		return
	}
}

func (r *InfinityNodeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
