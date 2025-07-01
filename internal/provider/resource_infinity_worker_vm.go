package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/pexip/go-infinity-sdk/v38/config"
	"github.com/pexip/terraform-provider-pexip/internal/provider/validators"
)

var (
	_ resource.ResourceWithImportState = (*InfinityWorkerVMResource)(nil)
)

type InfinityWorkerVMResource struct {
	InfinityClient InfinityClient
}

type InfinityWorkerVMResourceModel struct {
	ID                    types.String `tfsdk:"id"`
	ResourceID            types.Int32  `tfsdk:"resource_id"`
	Name                  types.String `tfsdk:"name"`
	Hostname              types.String `tfsdk:"hostname"`
	Domain                types.String `tfsdk:"domain"`
	Address               types.String `tfsdk:"address"`
	Netmask               types.String `tfsdk:"netmask"`
	Gateway               types.String `tfsdk:"gateway"`
	IPv6Address           types.String `tfsdk:"ipv6_address"`
	IPv6Gateway           types.String `tfsdk:"ipv6_gateway"`
	VMCPUCount            types.Int64  `tfsdk:"vm_cpu_count"`
	VMSystemMemory        types.Int64  `tfsdk:"vm_system_memory"`
	NodeType              types.String `tfsdk:"node_type"`
	Transcoding           types.Bool   `tfsdk:"transcoding"`
	Password              types.String `tfsdk:"password"`
	MaintenanceMode       types.Bool   `tfsdk:"maintenance_mode"`
	MaintenanceModeReason types.String `tfsdk:"maintenance_mode_reason"`
	SystemLocation        types.String `tfsdk:"system_location"`
}

func (r *InfinityWorkerVMResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_infinity_worker_vm"
}

func (r *InfinityWorkerVMResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *InfinityWorkerVMResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Resource URI for the worker VM in Infinity",
			},
			"resource_id": schema.Int32Attribute{
				Computed:            true,
				MarkdownDescription: "The resource integer identifier for the worker VM in Infinity",
			},
			"name": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(250),
				},
				MarkdownDescription: "The name of the worker VM. Maximum length: 250 characters.",
			},
			"hostname": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(250),
				},
				MarkdownDescription: "The hostname of the worker VM. Maximum length: 250 characters.",
			},
			"domain": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(250),
				},
				MarkdownDescription: "The domain of the worker VM. Maximum length: 250 characters.",
			},
			"address": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					validators.IPAddress(),
				},
				MarkdownDescription: "The IPv4 address of the worker VM.",
			},
			"netmask": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					validators.IPAddress(),
				},
				MarkdownDescription: "The netmask for the worker VM.",
			},
			"gateway": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					validators.IPAddress(),
				},
				MarkdownDescription: "The gateway address for the worker VM.",
			},
			"ipv6_address": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(250),
				},
				MarkdownDescription: "The IPv6 address of the worker VM. Maximum length: 250 characters.",
			},
			"ipv6_gateway": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(250),
				},
				MarkdownDescription: "The IPv6 gateway for the worker VM. Maximum length: 250 characters.",
			},
			"vm_cpu_count": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The number of CPUs for the VM. Defaults to system default.",
			},
			"vm_system_memory": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The amount of system memory (MB) for the VM. Defaults to system default.",
			},
			"node_type": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.OneOf("transcoding", "conferencing", "proxying"),
				},
				MarkdownDescription: "The node type. Valid choices: transcoding, conferencing, proxying.",
			},
			"transcoding": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Whether transcoding is enabled. Defaults to false.",
			},
			"password": schema.StringAttribute{
				Optional:  true,
				Computed:  true,
				Sensitive: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(250),
				},
				MarkdownDescription: "The password for the worker VM. Maximum length: 250 characters.",
			},
			"maintenance_mode": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Whether the worker VM is in maintenance mode. Defaults to false.",
			},
			"maintenance_mode_reason": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(250),
				},
				MarkdownDescription: "The reason for maintenance mode. Maximum length: 250 characters.",
			},
			"system_location": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Reference to system location resource URI.",
			},
		},
		MarkdownDescription: "Manages a worker VM configuration with the Infinity service.",
	}
}

func (r *InfinityWorkerVMResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	plan := &InfinityWorkerVMResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createRequest := &config.WorkerVMCreateRequest{
		Name:           plan.Name.ValueString(),
		Hostname:       plan.Hostname.ValueString(),
		Domain:         plan.Domain.ValueString(),
		Address:        plan.Address.ValueString(),
		Netmask:        plan.Netmask.ValueString(),
		Gateway:        plan.Gateway.ValueString(),
		SystemLocation: plan.SystemLocation.ValueString(),
	}

	// Set optional fields
	if !plan.IPv6Address.IsNull() {
		ipv6Address := plan.IPv6Address.ValueString()
		createRequest.IPv6Address = &ipv6Address
	}
	if !plan.IPv6Gateway.IsNull() {
		ipv6Gateway := plan.IPv6Gateway.ValueString()
		createRequest.IPv6Gateway = &ipv6Gateway
	}
	if !plan.VMCPUCount.IsNull() {
		createRequest.VMCPUCount = int(plan.VMCPUCount.ValueInt64())
	}
	if !plan.VMSystemMemory.IsNull() {
		createRequest.VMSystemMemory = int(plan.VMSystemMemory.ValueInt64())
	}
	if !plan.NodeType.IsNull() {
		createRequest.NodeType = plan.NodeType.ValueString()
	}
	if !plan.Transcoding.IsNull() {
		createRequest.Transcoding = plan.Transcoding.ValueBool()
	}
	if !plan.Password.IsNull() {
		createRequest.Password = plan.Password.ValueString()
	}
	if !plan.MaintenanceMode.IsNull() {
		createRequest.MaintenanceMode = plan.MaintenanceMode.ValueBool()
	}
	if !plan.MaintenanceModeReason.IsNull() {
		createRequest.MaintenanceModeReason = plan.MaintenanceModeReason.ValueString()
	}

	createResponse, err := r.InfinityClient.Config().CreateWorkerVM(ctx, createRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Infinity worker VM",
			fmt.Sprintf("Could not create Infinity worker VM: %s", err),
		)
		return
	}

	resourceID, err := createResponse.ResourceID()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Retrieving Infinity worker VM ID",
			fmt.Sprintf("Could not retrieve ID for created Infinity worker VM: %s", err),
		)
		return
	}

	// Read the state from the API to get all computed values
	model, err := r.read(ctx, resourceID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Created Infinity worker VM",
			fmt.Sprintf("Could not read created Infinity worker VM with ID %d: %s", resourceID, err),
		)
		return
	}
	tflog.Trace(ctx, fmt.Sprintf("created Infinity worker VM with ID: %s, name: %s", model.ID, model.Name))

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *InfinityWorkerVMResource) read(ctx context.Context, resourceID int) (*InfinityWorkerVMResourceModel, error) {
	var data InfinityWorkerVMResourceModel

	srv, err := r.InfinityClient.Config().GetWorkerVM(ctx, resourceID)
	if err != nil {
		return nil, err
	}

	if len(srv.ResourceURI) == 0 {
		return nil, fmt.Errorf("worker VM with ID %d not found", resourceID)
	}

	data.ID = types.StringValue(srv.ResourceURI)
	data.ResourceID = types.Int32Value(int32(resourceID))
	data.Name = types.StringValue(srv.Name)
	data.Hostname = types.StringValue(srv.Hostname)
	data.Domain = types.StringValue(srv.Domain)
	data.Address = types.StringValue(srv.Address)
	data.Netmask = types.StringValue(srv.Netmask)
	data.Gateway = types.StringValue(srv.Gateway)
	if srv.IPv6Address != nil {
		data.IPv6Address = types.StringValue(*srv.IPv6Address)
	} else {
		data.IPv6Address = types.StringNull()
	}
	if srv.IPv6Gateway != nil {
		data.IPv6Gateway = types.StringValue(*srv.IPv6Gateway)
	} else {
		data.IPv6Gateway = types.StringNull()
	}
	data.VMCPUCount = types.Int64Value(int64(srv.VMCPUCount))
	data.VMSystemMemory = types.Int64Value(int64(srv.VMSystemMemory))
	data.NodeType = types.StringValue(srv.NodeType)
	data.Transcoding = types.BoolValue(srv.Transcoding)
	data.Password = types.StringValue(srv.Password)
	data.MaintenanceMode = types.BoolValue(srv.MaintenanceMode)
	data.MaintenanceModeReason = types.StringValue(srv.MaintenanceModeReason)
	data.SystemLocation = types.StringValue(srv.SystemLocation)

	return &data, nil
}

func (r *InfinityWorkerVMResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	state := &InfinityWorkerVMResourceModel{}

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
			"Error Reading Infinity worker VM",
			fmt.Sprintf("Could not read Infinity worker VM: %s", err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *InfinityWorkerVMResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	plan := &InfinityWorkerVMResourceModel{}
	state := &InfinityWorkerVMResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resourceID := int(state.ResourceID.ValueInt32())

	updateRequest := &config.WorkerVMUpdateRequest{
		Name:           plan.Name.ValueString(),
		Hostname:       plan.Hostname.ValueString(),
		Domain:         plan.Domain.ValueString(),
		Address:        plan.Address.ValueString(),
		Netmask:        plan.Netmask.ValueString(),
		Gateway:        plan.Gateway.ValueString(),
		SystemLocation: plan.SystemLocation.ValueString(),
	}

	// Set optional fields
	if !plan.IPv6Address.IsNull() {
		ipv6Address := plan.IPv6Address.ValueString()
		updateRequest.IPv6Address = &ipv6Address
	}
	if !plan.IPv6Gateway.IsNull() {
		ipv6Gateway := plan.IPv6Gateway.ValueString()
		updateRequest.IPv6Gateway = &ipv6Gateway
	}
	if !plan.VMCPUCount.IsNull() {
		updateRequest.VMCPUCount = int(plan.VMCPUCount.ValueInt64())
	}
	if !plan.VMSystemMemory.IsNull() {
		updateRequest.VMSystemMemory = int(plan.VMSystemMemory.ValueInt64())
	}
	if !plan.NodeType.IsNull() {
		updateRequest.NodeType = plan.NodeType.ValueString()
	}
	if !plan.Transcoding.IsNull() {
		updateRequest.Transcoding = plan.Transcoding.ValueBool()
	}
	if !plan.Password.IsNull() {
		updateRequest.Password = plan.Password.ValueString()
	}
	if !plan.MaintenanceMode.IsNull() {
		updateRequest.MaintenanceMode = plan.MaintenanceMode.ValueBool()
	}
	if !plan.MaintenanceModeReason.IsNull() {
		updateRequest.MaintenanceModeReason = plan.MaintenanceModeReason.ValueString()
	}

	_, err := r.InfinityClient.Config().UpdateWorkerVM(ctx, resourceID, updateRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Infinity worker VM",
			fmt.Sprintf("Could not update Infinity worker VM with ID %d: %s", resourceID, err),
		)
		return
	}

	// Re-read the resource to get the latest state
	updatedModel, err := r.read(ctx, resourceID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Updated Infinity worker VM",
			fmt.Sprintf("Could not read updated Infinity worker VM with ID %d: %s", resourceID, err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, updatedModel)...)
}

func (r *InfinityWorkerVMResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	state := &InfinityWorkerVMResourceModel{}

	tflog.Info(ctx, "Deleting Infinity worker VM")

	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.InfinityClient.Config().DeleteWorkerVM(ctx, int(state.ResourceID.ValueInt32()))

	// Ignore 404 Not Found and Lookup errors on delete
	if err != nil && !isNotFoundError(err) && !isLookupError(err) {
		resp.Diagnostics.AddError(
			"Error Deleting Infinity worker VM",
			fmt.Sprintf("Could not delete Infinity worker VM with ID %s: %s", state.ID.ValueString(), err),
		)
		return
	}
}

func (r *InfinityWorkerVMResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resourceID, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Import ID must be a valid integer for the resource ID. Got: %s", req.ID),
		)
		return
	}

	tflog.Trace(ctx, fmt.Sprintf("Importing Infinity worker VM with resource ID: %d", resourceID))

	// Read the resource from the API
	model, err := r.read(ctx, resourceID)
	if err != nil {
		// Check if the error is a 404 (not found)
		if isNotFoundError(err) {
			resp.Diagnostics.AddError(
				"Infinity Worker VM Not Found",
				fmt.Sprintf("Infinity worker VM with resource ID %d not found.", resourceID),
			)
			return
		}
		resp.Diagnostics.AddError(
			"Error Importing Infinity Worker VM",
			fmt.Sprintf("Could not import Infinity worker VM with resource ID %d: %s", resourceID, err),
		)
		return
	}

	// Set the state from the imported resource
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}
