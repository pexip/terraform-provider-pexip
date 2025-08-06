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
	_ resource.ResourceWithImportState = (*InfinitySnmpNetworkManagementSystemResource)(nil)
)

type InfinitySnmpNetworkManagementSystemResource struct {
	InfinityClient InfinityClient
}

type InfinitySnmpNetworkManagementSystemResourceModel struct {
	ID                types.String `tfsdk:"id"`
	ResourceID        types.Int32  `tfsdk:"resource_id"`
	Name              types.String `tfsdk:"name"`
	Description       types.String `tfsdk:"description"`
	Address           types.String `tfsdk:"address"`
	Port              types.Int64  `tfsdk:"port"`
	SnmpTrapCommunity types.String `tfsdk:"snmp_trap_community"`
}

func (r *InfinitySnmpNetworkManagementSystemResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_infinity_snmp_network_management_system"
}

func (r *InfinitySnmpNetworkManagementSystemResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *InfinitySnmpNetworkManagementSystemResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Resource URI for the SNMP network management system in Infinity",
			},
			"resource_id": schema.Int32Attribute{
				Computed:            true,
				MarkdownDescription: "The resource integer identifier for the SNMP network management system in Infinity",
			},
			"name": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					stringvalidator.LengthAtMost(250),
				},
				MarkdownDescription: "The name of the SNMP network management system. Maximum length: 250 characters.",
			},
			"description": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(500),
				},
				MarkdownDescription: "Description of the SNMP network management system. Maximum length: 500 characters.",
			},
			"address": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(255),
				},
				MarkdownDescription: "The IP address or FQDN of the SNMP Network Management System. Maximum length: 255 characters.",
			},
			"port": schema.Int64Attribute{
				Required: true,
				Validators: []validator.Int64{
					int64validator.Between(1, 65535),
				},
				MarkdownDescription: "The port number for SNMP communications. Valid range: 1-65535.",
			},
			"snmp_trap_community": schema.StringAttribute{
				Required:            true,
				Sensitive:           true,
				MarkdownDescription: "The SNMP trap community string for authentication. This field is sensitive.",
			},
		},
		MarkdownDescription: "Manages an SNMP network management system with the Infinity service. SNMP network management systems receive SNMP traps and notifications from Pexip Infinity for monitoring and alerting purposes.",
	}
}

func (r *InfinitySnmpNetworkManagementSystemResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	plan := &InfinitySnmpNetworkManagementSystemResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createRequest := &config.SnmpNetworkManagementSystemCreateRequest{
		Name:              plan.Name.ValueString(),
		Description:       plan.Description.ValueString(),
		Address:           plan.Address.ValueString(),
		Port:              int(plan.Port.ValueInt64()),
		SnmpTrapCommunity: plan.SnmpTrapCommunity.ValueString(),
	}

	createResponse, err := r.InfinityClient.Config().CreateSnmpNetworkManagementSystem(ctx, createRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Infinity SNMP network management system",
			fmt.Sprintf("Could not create Infinity SNMP network management system: %s", err),
		)
		return
	}

	resourceID, err := createResponse.ResourceID()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Retrieving Infinity SNMP network management system ID",
			fmt.Sprintf("Could not retrieve ID for created Infinity SNMP network management system: %s", err),
		)
		return
	}

	// Read the state from the API to get all computed values
	model, err := r.read(ctx, resourceID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Created Infinity SNMP network management system",
			fmt.Sprintf("Could not read created Infinity SNMP network management system with ID %d: %s", resourceID, err),
		)
		return
	}
	tflog.Trace(ctx, fmt.Sprintf("created Infinity SNMP network management system with ID: %s, name: %s", model.ID, model.Name))

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *InfinitySnmpNetworkManagementSystemResource) read(ctx context.Context, resourceID int) (*InfinitySnmpNetworkManagementSystemResourceModel, error) {
	var data InfinitySnmpNetworkManagementSystemResourceModel

	srv, err := r.InfinityClient.Config().GetSnmpNetworkManagementSystem(ctx, resourceID)
	if err != nil {
		return nil, err
	}

	if srv.ResourceURI == "" {
		return nil, fmt.Errorf("SNMP network management system with ID %d not found", resourceID)
	}

	data.ID = types.StringValue(srv.ResourceURI)
	data.ResourceID = types.Int32Value(int32(resourceID)) // #nosec G115 -- API values are expected to be within int32 range
	data.Name = types.StringValue(srv.Name)
	data.Description = types.StringValue(srv.Description)
	data.Address = types.StringValue(srv.Address)
	data.Port = types.Int64Value(int64(srv.Port))
	data.SnmpTrapCommunity = types.StringValue(srv.SnmpTrapCommunity)

	return &data, nil
}

func (r *InfinitySnmpNetworkManagementSystemResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	state := &InfinitySnmpNetworkManagementSystemResourceModel{}

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
			"Error Reading Infinity SNMP network management system",
			fmt.Sprintf("Could not read Infinity SNMP network management system: %s", err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *InfinitySnmpNetworkManagementSystemResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	plan := &InfinitySnmpNetworkManagementSystemResourceModel{}
	state := &InfinitySnmpNetworkManagementSystemResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateRequest := &config.SnmpNetworkManagementSystemUpdateRequest{
		Name:              plan.Name.ValueString(),
		Description:       plan.Description.ValueString(),
		Address:           plan.Address.ValueString(),
		SnmpTrapCommunity: plan.SnmpTrapCommunity.ValueString(),
	}

	// Handle optional pointer field for port
	if !plan.Port.IsNull() && !plan.Port.IsUnknown() {
		port := int(plan.Port.ValueInt64())
		updateRequest.Port = &port
	}

	resourceID := int(state.ResourceID.ValueInt32())
	_, err := r.InfinityClient.Config().UpdateSnmpNetworkManagementSystem(ctx, resourceID, updateRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Infinity SNMP network management system",
			fmt.Sprintf("Could not update Infinity SNMP network management system: %s", err),
		)
		return
	}

	// Read the state from the API to get all computed values
	model, err := r.read(ctx, resourceID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Updated Infinity SNMP network management system",
			fmt.Sprintf("Could not read updated Infinity SNMP network management system with ID %d: %s", resourceID, err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *InfinitySnmpNetworkManagementSystemResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	state := &InfinitySnmpNetworkManagementSystemResourceModel{}

	tflog.Info(ctx, "Deleting Infinity SNMP network management system")

	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.InfinityClient.Config().DeleteSnmpNetworkManagementSystem(ctx, int(state.ResourceID.ValueInt32()))

	// Ignore 404 Not Found and Lookup errors on delete
	if err != nil && !isNotFoundError(err) && !isLookupError(err) {
		resp.Diagnostics.AddError(
			"Error Deleting Infinity SNMP network management system",
			fmt.Sprintf("Could not delete Infinity SNMP network management system with ID %s: %s", state.ID.ValueString(), err),
		)
		return
	}
}

func (r *InfinitySnmpNetworkManagementSystemResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resourceID, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Import ID must be a valid integer for the resource ID. Got: %s", req.ID),
		)
		return
	}

	tflog.Trace(ctx, fmt.Sprintf("Importing Infinity SNMP network management system with resource ID: %d", resourceID))

	// Read the resource from the API
	model, err := r.read(ctx, resourceID)
	if err != nil {
		// Check if the error is a 404 (not found)
		if isNotFoundError(err) {
			resp.Diagnostics.AddError(
				"Infinity SNMP Network Management System Not Found",
				fmt.Sprintf("Infinity SNMP network management system with resource ID %d not found.", resourceID),
			)
			return
		}
		resp.Diagnostics.AddError(
			"Error Importing Infinity SNMP Network Management System",
			fmt.Sprintf("Could not import Infinity SNMP network management system with resource ID %d: %s", resourceID, err),
		)
		return
	}

	// Set the state from the imported resource
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}
