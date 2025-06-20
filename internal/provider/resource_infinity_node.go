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
	"github.com/pexip/terraform-provider-pexip/internal/provider/validators"
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
	ID                    types.Int32  `tfsdk:"id"`
	Name                  types.String `tfsdk:"name"`
	Hostname              types.String `tfsdk:"hostname"`
	Address               types.String `tfsdk:"address"`
	Netmask               types.String `tfsdk:"netmask"`
	Domain                types.String `tfsdk:"domain"`
	Gateway               types.String `tfsdk:"gateway"`
	Password              types.String `tfsdk:"password"`
	NodeType              types.String `tfsdk:"node_type"`
	SystemLocation        types.String `tfsdk:"system_location"`
	MaintenanceMode       types.Bool   `tfsdk:"maintenance_mode"`
	MaintenanceModeReason types.String `tfsdk:"maintenance_mode_reason"`
	Transcoding           types.Bool   `tfsdk:"transcoding"`
	VMCPUCount            types.Int64  `tfsdk:"vm_cpu_count"`
	VMSystemMemory        types.Int64  `tfsdk:"vm_system_memory"`

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
					stringvalidator.LengthAtLeast(2),
					stringvalidator.LengthAtMost(255),
				},
				MarkdownDescription: "The name of the Infinity node. This should be unique within the Infinity cluster.",
			},
			"hostname": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(2),
					stringvalidator.LengthAtMost(255),
				},
				MarkdownDescription: "The hostname of the Infinity node. This should be resolvable within the Infinity cluster.",
			},
			"address": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					validators.IPAddress(),
				},
				MarkdownDescription: "The IP address of the Infinity node. This should be reachable within the Infinity cluster.",
			},
			"netmask": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					validators.Netmask(),
				},
				MarkdownDescription: "The netmask for the Infinity node's network interface.",
			},
			"domain": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					validators.Domain(),
				},
				MarkdownDescription: "The domain name for the Infinity node. This is used for DNS resolution within the Infinity cluster.",
			},
			"gateway": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					validators.IPAddress(),
				},
				MarkdownDescription: "The gateway IP address for the Infinity node. This is used for routing traffic outside the Infinity cluster.",
			},
			"password": schema.StringAttribute{
				Required:  true,
				Sensitive: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(2),
					stringvalidator.LengthAtMost(255),
				},
				MarkdownDescription: "The password for the Infinity node. This is used for authentication and should be kept secure.",
			},
			"node_type": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf("CONFERENCING", "PROXYING"),
				},
				MarkdownDescription: "The type of the Infinity node. Valid values are `worker`, `controller`, or `transcoder`. This determines the role of the node in the Infinity cluster.",
			},
			"system_location": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(2),
					stringvalidator.LengthAtMost(20),
				},
				MarkdownDescription: "The system location for the Infinity node. This is used for geographical identification and should be a valid location string.",
			},
			"maintenance_mode": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Indicates whether the Infinity node is in maintenance mode. When set to `true`, the node will not accept new workloads and will be excluded from load balancing.",
			},
			"maintenance_mode_reason": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The reason for putting the Infinity node into maintenance mode. This is optional and can be used to provide context for the maintenance operation.",
			},
			"transcoding": schema.BoolAttribute{
				Required:            true,
				MarkdownDescription: "Indicates whether the Infinity node is capable of transcoding media streams. This should be set to `true` if the node has the necessary resources and software to perform transcoding operations.",
			},
			"vm_cpu_count": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "The number of CPUs used by the Infinity node",
			},
			"vm_system_memory": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "The amount of system memory (RAM) allocated to the Infinity node, in megabytes.",
			},
			"config": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Bootstrap configuration for the Infinity Node.",
			},
		},
		MarkdownDescription: "Registers a node with the Infinity service. This resource is used to manage the lifecycle of nodes in the Infinity cluster.",
	}
}

func (r *InfinityNodeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	plan := &InfinityNodeResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createRequest := &config.WorkerVMCreateRequest{}

	// Set name if provided, otherwise use generated name
	if !plan.Name.IsNull() && !plan.Name.IsUnknown() {
		createRequest.Name = plan.Name.ValueString()
		createRequest.Hostname = plan.Hostname.ValueString()
		createRequest.Address = plan.Address.ValueString()
		createRequest.Netmask = plan.Netmask.ValueString()
		createRequest.Domain = plan.Domain.ValueString()
		createRequest.Gateway = plan.Gateway.ValueString()
		createRequest.Password = plan.Password.ValueString()
		createRequest.NodeType = plan.NodeType.ValueString()
		createRequest.SystemLocation = plan.SystemLocation.ValueString()
		createRequest.MaintenanceMode = plan.MaintenanceMode.ValueBool()
		createRequest.MaintenanceModeReason = plan.MaintenanceModeReason.ValueString()
		createRequest.Transcoding = plan.Transcoding.ValueBool()
		createRequest.VMCPUCount = int(plan.VMCPUCount.ValueInt64())
		createRequest.VMSystemMemory = int(plan.VMSystemMemory.ValueInt64())
	}

	createResponse, err := r.InfinityClient.Config.CreateWorkerVM(ctx, createRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Infinity Node",
			fmt.Sprintf("Could not create Infinity node: %s", err),
		)
		return
	}

	resourceID, err := createResponse.ResourceID()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Retrieving Infinity Node ID",
			fmt.Sprintf("Could not retrieve ID for created Infinity node: %s", err),
		)
		return
	}

	plan, err = r.read(ctx, resourceID)
	plan.Config = types.StringValue(string(createResponse.Body))
	tflog.Trace(ctx, fmt.Sprintf("created Infinity node with ID: %d, name: %s", plan.ID, plan.Name))

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *InfinityNodeResource) read(ctx context.Context, resourceID int) (*InfinityNodeResourceModel, error) {
	var data InfinityNodeResourceModel

	vm, err := r.InfinityClient.Config.GetWorkerVM(ctx, resourceID)
	if err != nil {
		return nil, err
	}

	data.Name = types.StringValue(vm.Name)
	data.Hostname = types.StringValue(vm.Hostname)
	data.Address = types.StringValue(vm.Address)
	data.Netmask = types.StringValue(vm.Netmask)
	data.Domain = types.StringValue(vm.Domain)
	data.Gateway = types.StringValue(vm.Gateway)
	data.Password = types.StringValue(vm.Password)
	data.NodeType = types.StringValue(vm.NodeType)
	data.SystemLocation = types.StringValue(vm.SystemLocation)
	data.MaintenanceMode = types.BoolValue(vm.MaintenanceMode)
	data.MaintenanceModeReason = types.StringValue(vm.MaintenanceModeReason)
	data.Transcoding = types.BoolValue(vm.Transcoding)
	data.VMCPUCount = types.Int64Value(int64(vm.VMCPUCount))
	data.VMSystemMemory = types.Int64Value(int64(vm.VMSystemMemory))

	return &data, nil
}

func (r *InfinityNodeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	state := &InfinityNodeResourceModel{}

	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	state, err := r.read(ctx, int(state.ID.ValueInt32()))
	if err != nil {
		// Check if the error is a 404 (not found)
		if isNotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading Infinity Node",
			fmt.Sprintf("Could not read Infinity node with ID %d: %s", state.ID.ValueInt32(), err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *InfinityNodeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan InfinityNodeResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateRequest := &config.WorkerVMUpdateRequest{}

	if !plan.Name.IsNull() && !plan.Name.IsUnknown() {
		updateRequest.Name = plan.Name.ValueString()
	}
	if !plan.Hostname.IsNull() && !plan.Hostname.IsUnknown() {
		updateRequest.Hostname = plan.Hostname.ValueString()
	}
	if !plan.Address.IsNull() && !plan.Address.IsUnknown() {
		updateRequest.Address = plan.Address.ValueString()
	}
	if !plan.Netmask.IsNull() && !plan.Netmask.IsUnknown() {
		updateRequest.Netmask = plan.Netmask.ValueString()
	}
	if !plan.Domain.IsNull() && !plan.Domain.IsUnknown() {
		updateRequest.Domain = plan.Domain.ValueString()
	}
	if !plan.Gateway.IsNull() && !plan.Gateway.IsUnknown() {
		updateRequest.Gateway = plan.Gateway.ValueString()
	}
	if !plan.Password.IsNull() && !plan.Password.IsUnknown() {
		updateRequest.Password = plan.Password.ValueString()
	}
	if !plan.NodeType.IsNull() && !plan.NodeType.IsUnknown() {
		updateRequest.NodeType = plan.NodeType.ValueString()
	}
	if !plan.SystemLocation.IsNull() && !plan.SystemLocation.IsUnknown() {
		updateRequest.SystemLocation = plan.SystemLocation.ValueString()
	}
	if !plan.MaintenanceMode.IsNull() && !plan.MaintenanceMode.IsUnknown() {
		updateRequest.MaintenanceMode = plan.MaintenanceMode.ValueBool()
	}
	if !plan.MaintenanceModeReason.IsNull() && !plan.MaintenanceModeReason.IsUnknown() {
		updateRequest.MaintenanceModeReason = plan.MaintenanceModeReason.ValueString()
	}
	if !plan.Transcoding.IsNull() && !plan.Transcoding.IsUnknown() {
		updateRequest.Transcoding = plan.Transcoding.ValueBool()
	}
	if !plan.VMCPUCount.IsNull() && !plan.VMCPUCount.IsUnknown() {
		updateRequest.VMCPUCount = int(plan.VMCPUCount.ValueInt64())
	}
	if !plan.VMSystemMemory.IsNull() && !plan.VMSystemMemory.IsUnknown() {
		updateRequest.VMSystemMemory = int(plan.VMSystemMemory.ValueInt64())
	}

	vm, err := r.InfinityClient.Config.UpdateWorkerVM(ctx, int(plan.ID.ValueInt32()), updateRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Infinity Node",
			fmt.Sprintf("Could not update Infinity node with ID %d: %s", plan.ID.ValueInt32(), err),
		)
		return
	}

	plan.Name = types.StringValue(vm.Name)
	plan.Hostname = types.StringValue(vm.Hostname)
	plan.Address = types.StringValue(vm.Address)
	plan.Netmask = types.StringValue(vm.Netmask)
	plan.Domain = types.StringValue(vm.Domain)
	plan.Gateway = types.StringValue(vm.Gateway)
	plan.Password = types.StringValue(vm.Password)
	plan.NodeType = types.StringValue(vm.NodeType)
	plan.SystemLocation = types.StringValue(vm.SystemLocation)
	plan.MaintenanceMode = types.BoolValue(vm.MaintenanceMode)
	plan.MaintenanceModeReason = types.StringValue(vm.MaintenanceModeReason)
	plan.Transcoding = types.BoolValue(vm.Transcoding)
	plan.VMCPUCount = types.Int64Value(int64(vm.VMCPUCount))
	plan.VMSystemMemory = types.Int64Value(int64(vm.VMSystemMemory))

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *InfinityNodeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state InfinityNodeResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.InfinityClient.Config.DeleteWorkerVM(ctx, int(state.ID.ValueInt32()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Infinity Node",
			fmt.Sprintf("Could not delete Infinity node with ID %d: %s", state.ID.ValueInt32(), err),
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
