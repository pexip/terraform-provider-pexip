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
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"

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
	ID                         types.String `tfsdk:"id"`
	ResourceID                 types.Int32  `tfsdk:"resource_id"`
	Name                       types.String `tfsdk:"name"`
	Hostname                   types.String `tfsdk:"hostname"`
	Domain                     types.String `tfsdk:"domain"`
	Address                    types.String `tfsdk:"address"`
	Netmask                    types.String `tfsdk:"netmask"`
	Gateway                    types.String `tfsdk:"gateway"`
	IPv6Address                types.String `tfsdk:"ipv6_address"`
	IPv6Gateway                types.String `tfsdk:"ipv6_gateway"`
	VMCPUCount                 types.Int64  `tfsdk:"vm_cpu_count"`
	VMSystemMemory             types.Int64  `tfsdk:"vm_system_memory"`
	NodeType                   types.String `tfsdk:"node_type"`
	Transcoding                types.Bool   `tfsdk:"transcoding"`
	Password                   types.String `tfsdk:"password"`
	MaintenanceMode            types.Bool   `tfsdk:"maintenance_mode"`
	MaintenanceModeReason      types.String `tfsdk:"maintenance_mode_reason"`
	SystemLocation             types.String `tfsdk:"system_location"`
	AlternativeFQDN            types.String `tfsdk:"alternative_fqdn"`
	CloudBursting              types.Bool   `tfsdk:"cloud_bursting"`
	DeploymentType             types.String `tfsdk:"deployment_type"`
	Description                types.String `tfsdk:"description"`
	EnableDistributedDatabase  types.Bool   `tfsdk:"enable_distributed_database"`
	EnableSSH                  types.String `tfsdk:"enable_ssh"`
	Managed                    types.Bool   `tfsdk:"managed"`
	MediaPriorityWeight        types.Int64  `tfsdk:"media_priority_weight"`
	SecondaryAddress           types.String `tfsdk:"secondary_address"`
	SecondaryNetmask           types.String `tfsdk:"secondary_netmask"`
	ServiceManager             types.Bool   `tfsdk:"service_manager"`
	ServicePolicy              types.Bool   `tfsdk:"service_policy"`
	Signalling                 types.Bool   `tfsdk:"signalling"`
	SNMPAuthenticationPassword types.String `tfsdk:"snmp_authentication_password"`
	SNMPCommunity              types.String `tfsdk:"snmp_community"`
	SNMPMode                   types.String `tfsdk:"snmp_mode"`
	SNMPPrivacyPassword        types.String `tfsdk:"snmp_privacy_password"`
	SNMPSystemContact          types.String `tfsdk:"snmp_system_contact"`
	SNMPSystemLocation         types.String `tfsdk:"snmp_system_location"`
	SNMPUsername               types.String `tfsdk:"snmp_username"`
	SSHAuthorizedKeys          types.Set    `tfsdk:"ssh_authorized_keys"`
	SSHAuthorizedKeysUseCloud  types.Bool   `tfsdk:"ssh_authorized_keys_use_cloud"`
	StaticNATAddress           types.String `tfsdk:"static_nat_address"`
	StaticRoutes               types.Set    `tfsdk:"static_routes"`
	TLSCertificate             types.String `tfsdk:"tls_certificate"`

	Config types.String `tfsdk:"config"`
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
			"address": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					validators.IPAddress(),
				},
				MarkdownDescription: "The IPv4 address of the worker VM.",
			},
			"alternative_fqdn": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(""),
				Validators: []validator.String{
					stringvalidator.LengthAtMost(255),
				},
				MarkdownDescription: "An identity for this Conferencing Node, used in signaling SIP TLS Contact addresses",
			},
			"cloud_bursting": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Defines whether this Conference Node is a cloud bursting node.",
			},
			"deployment_type": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("MANUAL-PROVISION-ONLY"),
				MarkdownDescription: "The means by which this Conferencing Node will be deployed",
			},
			"description": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(""),
				Validators: []validator.String{
					stringvalidator.LengthAtMost(250),
				},
				MarkdownDescription: "A description of the Conferencing Node.",
			},
			"domain": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(192),
				},
				MarkdownDescription: "The domain of the worker VM. Maximum length: 250 characters.",
			},
			"enable_distributed_database": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
				MarkdownDescription: "This should usually be True for all nodes which are expected to be 'always on', and False for nodes which are expected to only be powered on some of the time (e.g. cloud bursting nodes that are likely to only be operational during peak times). Avoid frequent toggling of this setting.",
			},
			"enable_ssh": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("GLOBAL"),
				Validators: []validator.String{
					stringvalidator.OneOf("GLOBAL", "OFF", "ON"),
				},
				MarkdownDescription: "Allows an administrator to log in to this node over SSH. Valid values are: global, off, on. Defaults to global.",
			},
			"gateway": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					validators.IPAddress(),
				},
				MarkdownDescription: "The gateway address for the worker VM.",
			},
			"hostname": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(63),
				},
				MarkdownDescription: "The hostname for this Conferencing Node. Each Conferencing Node must have a unique DNS hostname. Maximum length: 63 characters.",
			},
			"ipv6_address": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(250),
				},
				MarkdownDescription: "The IPv6 address of the conferencing node. Maximum length: 250 characters.",
			},
			"ipv6_gateway": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(250),
				},
				MarkdownDescription: "The IPv6 gateway for the conferencing node. Maximum length: 250 characters.",
			},
			"maintenance_mode": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Whether the worker VM is in maintenance mode. Defaults to false.",
			},
			"maintenance_mode_reason": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(250),
				},
				Default:             stringdefault.StaticString(""),
				MarkdownDescription: "The reason for maintenance mode. Maximum length: 250 characters.",
			},
			"managed": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Whether the conferencing node is managed by the Infinity service. Defaults to false.",
			},
			"media_priority_weight": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(0),
				MarkdownDescription: "The relative priority of this node, used when determining the order of nodes to which Pexip Infinity will attempt to send media. A higher number represents a higher priority; the default is 0, i.e. the lowest priority.",
			},
			"name": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(32),
				},
				MarkdownDescription: "he name used to refer to this Conferencing Node. Each Conferencing Node must have a unique name. Maximum length: 32 characters.",
			},
			"netmask": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					validators.IPAddress(),
				},
				MarkdownDescription: "The IPv4 network mask for this Conferencing Node.",
			},
			"node_type": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("CONFERENCING"),
				Validators: []validator.String{
					stringvalidator.OneOf("CONFERENCING", "PROXYING"),
				},
				MarkdownDescription: "The role of this Conferencing Node. Valid choices: CONFERENCING, PROXYING. Defaults to CONFERENCING.",
			},
			"password": schema.StringAttribute{
				Optional:  true,
				Computed:  true,
				Sensitive: true,
				Default:   stringdefault.StaticString(""),
				Validators: []validator.String{
					stringvalidator.LengthAtMost(250),
				},
				MarkdownDescription: "The password to be used when logging in to this Conferencing Node over SSH. The username will always be admin.",
			},
			"secondary_address": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  nil,
				Validators: []validator.String{
					validators.IPAddress(),
				},
				MarkdownDescription: "The optional secondary interface IPv4 address for this Conferencing Node.",
			},
			"secondary_netmask": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  nil,
				Validators: []validator.String{
					validators.IPAddress(),
				},
				MarkdownDescription: "The optional secondary interface IPv4 netmask for this Conferencing Node.",
			},
			"service_manager": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
				MarkdownDescription: "Handle Service Manager.",
			},
			"service_policy": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
				MarkdownDescription: "Handle Service Policy.",
			},
			"signalling": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
				MarkdownDescription: "Handle signalling",
			},
			"snmp_authentication_password": schema.StringAttribute{
				Optional:  true,
				Computed:  true,
				Sensitive: false,
				Default:   stringdefault.StaticString(""),
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(8),
					stringvalidator.LengthAtMost(100),
				},
				MarkdownDescription: "The password used for SNMPv3 authentication. Minimum length: 8 characters. Maximum length: 100 characters.",
			},
			"snmp_community": schema.StringAttribute{
				Optional:  true,
				Computed:  true,
				Sensitive: true,
				Default:   stringdefault.StaticString("public"),
				Validators: []validator.String{
					stringvalidator.LengthAtMost(16),
				},
				MarkdownDescription: "The SNMP group to which this virtual machine belongs. Maximum length: 16 characters.",
			},
			"snmp_mode": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("DISABLED"),
				Validators: []validator.String{
					stringvalidator.OneOf("DISABLED", "STANDARD", "AUTHPRIV"),
				},
				MarkdownDescription: "The SNMP mode.",
			},
			"snmp_privacy_password": schema.StringAttribute{
				Optional:  true,
				Computed:  true,
				Sensitive: false,
				Default:   stringdefault.StaticString(""),
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(8),
					stringvalidator.LengthAtMost(100),
				},
				MarkdownDescription: "The password used for SNMPv3 privacy. Minimum length: 8 characters. Maximum length: 100 characters.",
			},
			"snmp_system_contact": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("admin@domain.com"),
				Validators: []validator.String{
					stringvalidator.LengthAtMost(70),
				},
				MarkdownDescription: "The SNMP contact for this Conferencing Node. Maximum length: 70 characters.",
			},
			"snmp_system_location": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("Virtual machine"),
				Validators: []validator.String{
					stringvalidator.LengthAtMost(70),
				},
				MarkdownDescription: "The SNMP location for this Conferencing Node. Maximum length: 70 characters.",
			},
			"snmp_username": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(""),
				Validators: []validator.String{
					stringvalidator.LengthAtMost(100),
				},
				MarkdownDescription: "The username used to authenticate SNMPv3 requests. Maximum length: 100 characters.",
			},
			"ssh_authorized_keys": schema.SetAttribute{
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
				Default:             nil,
				MarkdownDescription: "The selected authorized keys.",
			},
			"ssh_authorized_keys_use_cloud": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
				MarkdownDescription: "Allows use of SSH keys configured in the cloud service.",
			},
			"static_nat_address": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  nil,
				Validators: []validator.String{
					validators.IPAddress(),
				},
				MarkdownDescription: "The public IPv4 address used by the Conferencing Node when it is located behind a NAT device. Note that if you are using NAT, you must also configure your NAT device to route the Conferencing Node's IPv4 static NAT address to its IPv4 address.",
			},
			"static_routes": schema.SetAttribute{
				Optional:    true,
				Computed:    true,
				Default:     nil,
				ElementType: types.StringType,
				Description: "Additional configuration to permit routing of traffic to networks not accessible through the configured default gateway.",
			},
			"system_location": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The system location for this Conferencing Node. A system location should not contain a mixture of Proxying Edge Nodes and Transcoding Conferencing Nodes.",
			},
			"tls_certificate": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             nil,
				MarkdownDescription: "The TLS certificate to use on this node.",
			},
			"transcoding": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
				MarkdownDescription: "This determines the Conferencing Node's role. When transcoding is enabled, this node can handle all the media processing, protocol interworking, mixing and so on that is required in hosting Pexip Infinity calls and conferences. When transcoding is disabled, it becomes a Proxying Edge Node that can only handle the media and signaling connections with an endpoint or external device, and it then forwards the device's media on to a node that does have transcoding capabilities.",
			},
			"vm_cpu_count": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				Default:  int64default.StaticInt64(4),
				Validators: []validator.Int64{
					int64validator.Between(2, 128),
				},
				MarkdownDescription: "Enter the number of virtual CPUs to assign to this Conferencing Node. We do not recommend that you assign more virtual CPUs than there are physical cores on a single processor on the host server (unless you have enabled NUMA affinity). For example, if the host server has 2 processors each with 12 physical cores, we recommend that you assign no more than 12 virtual CPUs. Range: 2 to 128. Default: 4.",
			},
			"vm_system_memory": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				Default:  int64default.StaticInt64(4096),
				Validators: []validator.Int64{
					int64validator.Between(2000, 64000),
				},
				MarkdownDescription: "The amount of RAM (in megabytes) to assign to this Conferencing Node. Range: 2000 to 64000. Default: 4096.",
			},
			"config": schema.StringAttribute{
				Computed:            true,
				Sensitive:           true,
				MarkdownDescription: "Bootstrap configuration for the Infinity Node.",
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

	// Convert List attributes to []string
	sshAuthorizedKeys, diags := getStringList(ctx, plan.SSHAuthorizedKeys)
	resp.Diagnostics.Append(diags...)
	staticRoutes, diags := getStringList(ctx, plan.StaticRoutes)
	resp.Diagnostics.Append(diags...)

	// create request with required fields and list fields
	createRequest := &config.WorkerVMCreateRequest{
		Name:              plan.Name.ValueString(),
		Hostname:          plan.Hostname.ValueString(),
		Domain:            plan.Domain.ValueString(),
		Address:           plan.Address.ValueString(),
		Netmask:           plan.Netmask.ValueString(),
		Gateway:           plan.Gateway.ValueString(),
		SystemLocation:    plan.SystemLocation.ValueString(),
		SSHAuthorizedKeys: sshAuthorizedKeys,
		StaticRoutes:      staticRoutes,
	}

	// Set optional fields that have default values
	if !plan.AlternativeFQDN.IsNull() {
		createRequest.AlternativeFQDN = plan.AlternativeFQDN.ValueString()
	}
	if !plan.DeploymentType.IsNull() {
		createRequest.DeploymentType = plan.DeploymentType.ValueString()
	}
	if !plan.Description.IsNull() {
		createRequest.Description = plan.Description.ValueString()
	}
	if !plan.EnableDistributedDatabase.IsNull() {
		createRequest.EnableDistributedDatabase = plan.EnableDistributedDatabase.ValueBool()
	}
	if !plan.EnableSSH.IsNull() {
		createRequest.EnableSSH = plan.EnableSSH.ValueString()
	}
	if !plan.MaintenanceMode.IsNull() {
		createRequest.MaintenanceMode = plan.MaintenanceMode.ValueBool()
	}
	if !plan.MaintenanceModeReason.IsNull() {
		createRequest.MaintenanceModeReason = plan.MaintenanceModeReason.ValueString()
	}
	if !plan.NodeType.IsNull() {
		createRequest.NodeType = plan.NodeType.ValueString()
	}
	if !plan.Password.IsNull() {
		createRequest.Password = plan.Password.ValueString()
	}
	if !plan.SNMPAuthenticationPassword.IsNull() {
		createRequest.SNMPAuthenticationPassword = plan.SNMPAuthenticationPassword.ValueString()
	}
	if !plan.SNMPCommunity.IsNull() {
		createRequest.SNMPCommunity = plan.SNMPCommunity.ValueString()
	}
	if !plan.SNMPMode.IsNull() {
		createRequest.SNMPMode = plan.SNMPMode.ValueString()
	}
	if !plan.SNMPPrivacyPassword.IsNull() {
		createRequest.SNMPPrivacyPassword = plan.SNMPPrivacyPassword.ValueString()
	}
	if !plan.SNMPSystemContact.IsNull() {
		createRequest.SNMPSystemContact = plan.SNMPSystemContact.ValueString()
	}
	if !plan.SNMPSystemLocation.IsNull() {
		createRequest.SNMPSystemLocation = plan.SNMPSystemLocation.ValueString()
	}
	if !plan.SNMPUsername.IsNull() {
		createRequest.SNMPUsername = plan.SNMPUsername.ValueString()
	}
	if !plan.SSHAuthorizedKeysUseCloud.IsNull() {
		createRequest.SSHAuthorizedKeysUseCloud = plan.SSHAuthorizedKeysUseCloud.ValueBool()
	}
	if !plan.Transcoding.IsNull() {
		createRequest.Transcoding = plan.Transcoding.ValueBool()
	}
	if !plan.VMCPUCount.IsNull() {
		createRequest.VMCPUCount = int(plan.VMCPUCount.ValueInt64())
	}
	if !plan.VMSystemMemory.IsNull() {
		createRequest.VMSystemMemory = int(plan.VMSystemMemory.ValueInt64())
	}
	// Set optional fields that are nullable
	if !plan.IPv6Address.IsNull() && !plan.IPv6Address.IsUnknown() {
		value := plan.IPv6Address.ValueString()
		createRequest.IPv6Address = &value
	}
	if !plan.IPv6Gateway.IsNull() && !plan.IPv6Gateway.IsUnknown() {
		value := plan.IPv6Gateway.ValueString()
		createRequest.IPv6Gateway = &value
	}
	if !plan.MediaPriorityWeight.IsNull() && !plan.MediaPriorityWeight.IsUnknown() {
		value := int(plan.MediaPriorityWeight.ValueInt64())
		createRequest.MediaPriorityWeight = &value
	}
	if !plan.SecondaryAddress.IsNull() && !plan.SecondaryAddress.IsUnknown() {
		value := plan.SecondaryAddress.ValueString()
		createRequest.SecondaryAddress = &value
	}
	if !plan.SecondaryNetmask.IsNull() && !plan.SecondaryNetmask.IsUnknown() {
		value := plan.SecondaryNetmask.ValueString()
		createRequest.SecondaryNetmask = &value
	}
	if !plan.StaticNATAddress.IsNull() && !plan.StaticNATAddress.IsUnknown() {
		value := plan.StaticNATAddress.ValueString()
		createRequest.StaticNATAddress = &value
	}
	if !plan.TLSCertificate.IsNull() && !plan.TLSCertificate.IsUnknown() {
		value := plan.TLSCertificate.ValueString()
		createRequest.TLSCertificate = &value
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
	model, err := r.read(ctx, resourceID, string(createResponse.Body), plan.DeploymentType.ValueString(), plan.Password.ValueString(), plan.SNMPAuthenticationPassword.ValueString(), plan.SNMPPrivacyPassword.ValueString(), plan.VMSystemMemory.ValueInt64(), plan.VMCPUCount.ValueInt64())
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

func (r *InfinityWorkerVMResource) read(ctx context.Context, resourceID int, config string, deployType string, password string, snmpAuthPass string, snmpPrivPass string, vm_system_memory int64, vm_cpu_count int64) (*InfinityWorkerVMResourceModel, error) {
	var data InfinityWorkerVMResourceModel

	srv, err := r.InfinityClient.Config().GetWorkerVM(ctx, resourceID)
	if err != nil {
		return nil, err
	}

	if srv.ResourceURI == "" {
		return nil, fmt.Errorf("worker VM with ID %d not found", resourceID)
	}

	// Set required and default fields
	data.ID = types.StringValue(srv.ResourceURI)
	data.ResourceID = types.Int32Value(int32(resourceID))
	data.Config = types.StringValue(config)
	data.Name = types.StringValue(srv.Name)
	data.Hostname = types.StringValue(srv.Hostname)
	data.Domain = types.StringValue(srv.Domain)
	data.Address = types.StringValue(srv.Address)
	data.Netmask = types.StringValue(srv.Netmask)
	data.Gateway = types.StringValue(srv.Gateway)
	data.Description = types.StringValue(srv.Description)
	data.VMCPUCount = types.Int64Value(vm_cpu_count)
	data.VMSystemMemory = types.Int64Value(vm_system_memory)
	data.NodeType = types.StringValue(srv.NodeType)
	// value not returned by API
	data.DeploymentType = types.StringValue(deployType)
	data.Transcoding = types.BoolValue(srv.Transcoding)
	data.CloudBursting = types.BoolValue(srv.CloudBursting)
	data.Managed = types.BoolValue(srv.Managed)
	data.ServiceManager = types.BoolValue(srv.ServiceManager)
	data.ServicePolicy = types.BoolValue(srv.ServicePolicy)
	data.Signalling = types.BoolValue(srv.Signalling)
	// value not returned by API
	data.Password = types.StringValue(password)
	data.MaintenanceMode = types.BoolValue(srv.MaintenanceMode)
	data.MaintenanceModeReason = types.StringValue(srv.MaintenanceModeReason)
	data.SystemLocation = types.StringValue(srv.SystemLocation)
	data.AlternativeFQDN = types.StringValue(srv.AlternativeFQDN)
	data.EnableDistributedDatabase = types.BoolValue(srv.EnableDistributedDatabase)
	data.EnableSSH = types.StringValue(srv.EnableSSH)
	data.SSHAuthorizedKeysUseCloud = types.BoolValue(srv.SSHAuthorizedKeysUseCloud)
	// Ignore the hashed value from the API
	data.SNMPAuthenticationPassword = types.StringValue(snmpAuthPass)
	data.SNMPCommunity = types.StringValue(srv.SNMPCommunity)
	data.SNMPMode = types.StringValue(srv.SNMPMode)
	// ignore hashed value returned by API
	data.SNMPPrivacyPassword = types.StringValue(snmpPrivPass)
	data.SNMPSystemContact = types.StringValue(srv.SNMPSystemContact)
	data.SNMPSystemLocation = types.StringValue(srv.SNMPSystemLocation)
	data.SNMPUsername = types.StringValue(srv.SNMPUsername)

	// Convert SSH authorized keys from SDK to Terraform format
	var sshAuthorizedKeys []string
	for _, key := range srv.SSHAuthorizedKeys {
		sshAuthorizedKeys = append(sshAuthorizedKeys, fmt.Sprintf("/api/admin/configuration/v1/ssh_authorized_key/%s/", key))
	}
	sshSetValue, diags := types.SetValueFrom(ctx, types.StringType, sshAuthorizedKeys)
	if diags.HasError() {
		return nil, fmt.Errorf("error converting SSH authorized keys: %v", diags)
	}
	data.SSHAuthorizedKeys = sshSetValue

	// Convert static routes from SDK to Terraform format
	var staticRoutes []string
	for _, route := range srv.StaticRoutes {
		staticRoutes = append(staticRoutes, fmt.Sprintf("/api/admin/configuration/v1/static_route/%s/", route))
	}
	routeSetValue, diags := types.SetValueFrom(ctx, types.StringType, staticRoutes)
	if diags.HasError() {
		return nil, fmt.Errorf("error converting static routes: %v", diags)
	}
	data.StaticRoutes = routeSetValue

	// Set nullable fields
	data.IPv6Address = types.StringPointerValue(srv.IPv6Address)
	data.IPv6Gateway = types.StringPointerValue(srv.IPv6Gateway)
	data.TLSCertificate = types.StringPointerValue(srv.TLSCertificate)
	data.SecondaryAddress = types.StringPointerValue(srv.SecondaryAddress)
	data.SecondaryNetmask = types.StringPointerValue(srv.SecondaryNetmask)
	data.StaticNATAddress = types.StringPointerValue(srv.StaticNATAddress)
	data.MediaPriorityWeight = types.Int64PointerValue(srv.MediaPriorityWeight)

	return &data, nil
}

func (r *InfinityWorkerVMResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	state := &InfinityWorkerVMResourceModel{}

	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resourceID := int(state.ResourceID.ValueInt32())
	state, err := r.read(ctx, resourceID, state.Config.ValueString(), state.DeploymentType.ValueString(), state.Password.ValueString(), state.SNMPAuthenticationPassword.ValueString(), state.SNMPPrivacyPassword.ValueString(), state.VMSystemMemory.ValueInt64(), state.VMCPUCount.ValueInt64())
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

	// Convert List attributes to []string
	sshAuthorizedKeys, diags := getStringList(ctx, plan.SSHAuthorizedKeys)
	resp.Diagnostics.Append(diags...)
	staticRoutes, diags := getStringList(ctx, plan.StaticRoutes)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// create request with required fields and list fields
	updateRequest := &config.WorkerVMUpdateRequest{
		Name:              plan.Name.ValueString(),
		Hostname:          plan.Hostname.ValueString(),
		Domain:            plan.Domain.ValueString(),
		Address:           plan.Address.ValueString(),
		Netmask:           plan.Netmask.ValueString(),
		Gateway:           plan.Gateway.ValueString(),
		SystemLocation:    plan.SystemLocation.ValueString(),
		SSHAuthorizedKeys: sshAuthorizedKeys,
		StaticRoutes:      staticRoutes,
	}

	// Set optional fields that have default values
	if !plan.AlternativeFQDN.IsNull() {
		updateRequest.AlternativeFQDN = plan.AlternativeFQDN.ValueString()
	}
	if !plan.DeploymentType.IsNull() {
		updateRequest.DeploymentType = plan.DeploymentType.ValueString()
	}
	if !plan.Description.IsNull() {
		updateRequest.Description = plan.Description.ValueString()
	}
	if !plan.EnableDistributedDatabase.IsNull() {
		updateRequest.EnableDistributedDatabase = plan.EnableDistributedDatabase.ValueBool()
	}
	if !plan.EnableSSH.IsNull() {
		updateRequest.EnableSSH = plan.EnableSSH.ValueString()
	}
	if !plan.MaintenanceMode.IsNull() {
		updateRequest.MaintenanceMode = plan.MaintenanceMode.ValueBool()
	}
	if !plan.MaintenanceModeReason.IsNull() {
		updateRequest.MaintenanceModeReason = plan.MaintenanceModeReason.ValueString()
	}
	if !plan.NodeType.IsNull() {
		updateRequest.NodeType = plan.NodeType.ValueString()
	}
	if !plan.Password.IsNull() {
		updateRequest.Password = plan.Password.ValueString()
	}
	if !plan.SNMPCommunity.IsNull() {
		updateRequest.SNMPCommunity = plan.SNMPCommunity.ValueString()
	}
	if !plan.SNMPMode.IsNull() {
		updateRequest.SNMPMode = plan.SNMPMode.ValueString()
	}
	if !plan.SNMPSystemContact.IsNull() {
		updateRequest.SNMPSystemContact = plan.SNMPSystemContact.ValueString()
	}
	if !plan.SNMPSystemLocation.IsNull() {
		updateRequest.SNMPSystemLocation = plan.SNMPSystemLocation.ValueString()
	}
	if !plan.SNMPUsername.IsNull() {
		updateRequest.SNMPUsername = plan.SNMPUsername.ValueString()
	}
	if !plan.SSHAuthorizedKeysUseCloud.IsNull() {
		updateRequest.SSHAuthorizedKeysUseCloud = plan.SSHAuthorizedKeysUseCloud.ValueBool()
	}
	if !plan.Transcoding.IsNull() {
		updateRequest.Transcoding = plan.Transcoding.ValueBool()
	}
	if !plan.VMCPUCount.IsNull() {
		updateRequest.VMCPUCount = int(plan.VMCPUCount.ValueInt64())
	}
	if !plan.VMSystemMemory.IsNull() {
		updateRequest.VMSystemMemory = int(plan.VMSystemMemory.ValueInt64())
	}

	if !plan.SNMPAuthenticationPassword.IsNull() {
		updateRequest.SNMPAuthenticationPassword = plan.SNMPAuthenticationPassword.ValueString()
	}
	if !plan.SNMPPrivacyPassword.IsNull() {
		updateRequest.SNMPPrivacyPassword = plan.SNMPPrivacyPassword.ValueString()
	}

	// Set optional fields that are nullable
	if !plan.IPv6Address.IsNull() && !plan.IPv6Address.IsUnknown() {
		value := plan.IPv6Address.ValueString()
		updateRequest.IPv6Address = &value
	}
	if !plan.IPv6Gateway.IsNull() && !plan.IPv6Gateway.IsUnknown() {
		value := plan.IPv6Gateway.ValueString()
		updateRequest.IPv6Gateway = &value
	}
	if !plan.MediaPriorityWeight.IsNull() && !plan.MediaPriorityWeight.IsUnknown() {
		value := int(plan.MediaPriorityWeight.ValueInt64())
		updateRequest.MediaPriorityWeight = &value
	}
	if !plan.SecondaryAddress.IsNull() && !plan.SecondaryAddress.IsUnknown() {
		value := plan.SecondaryAddress.ValueString()
		updateRequest.SecondaryAddress = &value
	}
	if !plan.SecondaryNetmask.IsNull() && !plan.SecondaryNetmask.IsUnknown() {
		value := plan.SecondaryNetmask.ValueString()
		updateRequest.SecondaryNetmask = &value
	}
	if !plan.StaticNATAddress.IsNull() && !plan.StaticNATAddress.IsUnknown() {
		value := plan.StaticNATAddress.ValueString()
		updateRequest.StaticNATAddress = &value
	}
	if !plan.TLSCertificate.IsNull() && !plan.TLSCertificate.IsUnknown() {
		value := plan.TLSCertificate.ValueString()
		updateRequest.TLSCertificate = &value
	}
	// Send the update request to the management node
	_, err := r.InfinityClient.Config().UpdateWorkerVM(ctx, resourceID, updateRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Infinity worker VM",
			fmt.Sprintf("Could not update Infinity worker VM with ID %d: %s", resourceID, err),
		)
		return
	}

	// Re-read the resource to get the latest state
	updatedModel, err := r.read(ctx, resourceID, plan.Config.ValueString(), plan.DeploymentType.ValueString(), plan.Password.ValueString(), plan.SNMPAuthenticationPassword.ValueString(), plan.SNMPPrivacyPassword.ValueString(), plan.VMSystemMemory.ValueInt64(), plan.VMCPUCount.ValueInt64())
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
	model, err := r.read(ctx, resourceID, "", "", "", "", "", 0, 0)
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

// isNotFoundError checks if the error indicates a 404/not found response
func isNotFoundError(err error) bool {
	return strings.Contains(err.Error(), "404") ||
		strings.Contains(err.Error(), "not found") ||
		strings.Contains(err.Error(), "Not Found")
}

func isLookupError(err error) bool {
	return strings.Contains(err.Error(), "lookup")
}
