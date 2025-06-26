package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/pexip/go-infinity-sdk/v38/config"
)

var (
	_ resource.ResourceWithImportState = (*InfinitySystemLocationResource)(nil)
)

type InfinitySystemLocationResource struct {
	InfinityClient InfinityClient
}

type InfinitySystemLocationResourceModel struct {
	ID          types.String `tfsdk:"id"`
	ResourceID  types.Int32  `tfsdk:"resource_id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	DNSServers  types.String `tfsdk:"dns_servers"`
	//NTPServers  types.String `tfsdk:"ntp_servers"`
	MTU         types.Int32  `tfsdk:"mtu"`

}

func (r *InfinitySystemLocationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_infinity_system_location"
}

func (r *InfinitySystemLocationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *InfinitySystemLocationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Resource URI for the system location in Infinity",
			},
			"resource_id": schema.Int32Attribute{
				Computed:            true,
				MarkdownDescription: "The resource integer identifier for the system location in Infinity",
			},
			
			"description": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(250),
				},
				MarkdownDescription: "A description of the system location. Maximum length: 250 characters.",
			},
			"dns_servers": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(250),
				},
				MarkdownDescription: "The DNS servers to be used by Conferencing Nodes deployed in this Location.",
			},
			"mtu": schema.Int32Attribute{
				Optional: true,
				Computed:            true,
				MarkdownDescription: "Maximum Transmission Unit - the size of the largest packet that can be transmitted via the network interface for this system location. It depends on your network topology as to whether you may need to specify an MTU value here. Range: 512 to 1500.",
			},
			"name": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(250),
				},
				MarkdownDescription: "The name used to refer to this system location. Maximum length: 250 characters.",
			},
		},
		MarkdownDescription: "Registers a system location with the Infinity service.",
	}
}

func (r *InfinitySystemLocationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	plan := &InfinitySystemLocationResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createRequest := &config.SystemLocationCreateRequest{
		Description: plan.Description.ValueString(),
		Name:		 plan.Name.ValueString(),
		MTU:		 int(plan.MTU.ValueInt32()),
	}

	createResponse, err := r.InfinityClient.Config().CreateSystemLocation(ctx, createRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Infinity system location",
			fmt.Sprintf("Could not create Infinity system location: %s", err),
		)
		return
	}

	resourceID, err := createResponse.ResourceID()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Retrieving Infinity system location ID",
			fmt.Sprintf("Could not retrieve ID for created Infinity system location: %s", err),
		)
		return
	}

	plan, err = r.read(ctx, resourceID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Created Infinity system location",
			fmt.Sprintf("Could not read created Infinity system location with ID %d: %s", resourceID, err),
		)
		return
	}
	tflog.Trace(ctx, fmt.Sprintf("created Infinity system location with ID: %s, name: %s", plan.ID, plan.Name))

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *InfinitySystemLocationResource) read(ctx context.Context, resourceID int) (*InfinitySystemLocationResourceModel, error) {
	var data InfinitySystemLocationResourceModel

	srv, err := r.InfinityClient.Config().GetSystemLocation(ctx, resourceID)
	if err != nil {
		return nil, err
	}

	if len(srv.ResourceURI) == 0 {
		return nil, fmt.Errorf("system location with ID %d not found", resourceID)
	}

	data.ID = types.StringValue(srv.ResourceURI)
	data.ResourceID = types.Int32Value(int32(resourceID))
	data.Description = types.StringValue(srv.Description)
	data.DNSServers = 
	data.MTU = types.Int32Value(int32(srv.MTU))
	data.Name = types.StringValue(srv.Name)

	return &data, nil
}

func (r *InfinitySystemLocationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	state := &InfinitySystemLocationResourceModel{}

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
			"Error Reading Infinity system location",
			fmt.Sprintf("Could not read Infinity system location: %s", err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *InfinitySystemLocationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	plan := &InfinitySystemLocationResourceModel{}
	state := &InfinitySystemLocationResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// The resource ID is required for the update API call.
	resourceID := int(state.ResourceID.ValueInt32())

	// Prepare the update request from the plan
	updateRequest := &config.SystemLocationUpdateRequest{
		Description: plan.Description.ValueString(),
		Name: plan.Name.ValueString(),
		MTU: int(plan.MTU.ValueInt32()),
	}
	_, err := r.InfinityClient.Config().UpdateSystemLocation(ctx, resourceID, updateRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Infinity system location",
			fmt.Sprintf("Could not update Infinity system location with ID %s: %s", plan.ID.ValueString(), err),
		)
		return
	}

	plan.ID = state.ID
	plan.ResourceID = state.ResourceID
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *InfinitySystemLocationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	state := &InfinitySystemLocationResourceModel{}

	tflog.Info(ctx, "Deleting Infinity system location")

	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.InfinityClient.Config().DeleteSystemLocation(ctx, int(state.ResourceID.ValueInt32()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Infinity system location",
			fmt.Sprintf("Could not delete Infinity system location with ID %s: %s", state.ID.ValueString(), err),
		)
		return
	}
}

func (r *InfinitySystemLocationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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
