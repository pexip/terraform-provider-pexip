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
	"github.com/pexip/go-infinity-sdk/v38"
	"github.com/pexip/go-infinity-sdk/v38/config"
	"github.com/pexip/terraform-provider-pexip/internal/provider/validators"
)

var (
	_ resource.ResourceWithImportState = (*InfinityDnsServerResource)(nil)
)

type InfinityDnsServerResource struct {
	InfinityClient *infinity.Client
}

type InfinityDnsServerResourceModel struct {
	ID                    types.Int32  `tfsdk:"id"`
	Address               types.String `tfsdk:"address"`
	Description           types.String `tfsdk:"description"`
}

func (r *InfinityDnsServerResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_infinity_dns_server"
}

func (r *InfinityDnsServerResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *InfinityDnsServerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.Int32Attribute{
				Computed:            true,
				MarkdownDescription: "Resource identifier",
			},
			"description": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(255),
				},
				MarkdownDescription: "A description of the DNS server. Maximum length: 250 characters.",
			},
			"address": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					validators.IPAddress(),
				},
				MarkdownDescription: "The IP address of the DNS server.",
			},
		},
		MarkdownDescription: "Registers a DNS server with the Infinity service.",
	}
}

func (r *InfinityDnsServerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	plan := &InfinityDnsServerResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createRequest := &config.DNSServerCreateRequest{}

	// Set name if provided, otherwise use generated name
	if !plan.Address.IsNull() && !plan.Address.IsUnknown() {
		createRequest.Address = plan.Address.ValueString()
		createRequest.Description = plan.Description.ValueString()
	}

	

	createResponse, err := r.InfinityClient.Config.CreateDNSServer(ctx, createRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Infinity DNS server",
			fmt.Sprintf("Could not create Infinity DNS server: %s", err),
		)
		return
	}

	resourceID, err := createResponse.ResourceID()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Retrieving Infinity DNS server ID",
			fmt.Sprintf("Could not retrieve ID for created Infinity DNS server: %s", err),
		)
		return
	}

	plan, err = r.read(ctx, resourceID)
	// don't think this is needed
	//plan.Config = types.StringValue(string(createResponse.Body))
	tflog.Trace(ctx, fmt.Sprintf("created Infinity DNS server with ID: %d, name: %s", plan.ID, plan.Address))

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *InfinityDnsServerResource) read(ctx context.Context, resourceID int) (*InfinityDnsServerResourceModel, error) {
	var data InfinityDnsServerResourceModel

	srv, err := r.InfinityClient.Config.GetDNSServer(ctx, resourceID)
	if err != nil {
		return nil, err
	}

	data.Address = types.StringValue(srv.Address)
	data.Description = types.StringValue(srv.Description)

	return &data, nil
}

func (r *InfinityDnsServerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	state := &InfinityDnsServerResourceModel{}

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
			"Error Reading Infinity DNS server",
			fmt.Sprintf("Could not read Infinity DNS server with ID %d: %s", state.ID.ValueInt32(), err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *InfinityDnsServerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan InfinityDnsServerResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateRequest := &config.DNSServerUpdateRequest{}

	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		updateRequest.Description = plan.Description.ValueString()
	}
	if !plan.Address.IsNull() && !plan.Address.IsUnknown() {
		updateRequest.Address = plan.Address.ValueString()
	}

	vm, err := r.InfinityClient.Config.UpdateDNSServer(ctx, int(plan.ID.ValueInt32()), updateRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Infinity DNS server",
			fmt.Sprintf("Could not update Infinity DNS server with ID %d: %s", plan.ID.ValueInt32(), err),
		)
		return
	}

	plan.Description = types.StringValue(vm.Description)
	plan.Address = types.StringValue(vm.Address)
	
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *InfinityDnsServerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state InfinityDnsServerResourceModel

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

func (r *InfinityDnsServerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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
//func isNotFoundError(err error) bool {
//	// This is a placeholder - you'll need to check the actual error types
//	// returned by the go-infinity-sdk to determine what constitutes a "not found" error
//	return strings.Contains(err.Error(), "404") ||
//		strings.Contains(err.Error(), "not found") ||
//		strings.Contains(err.Error(), "Not Found")
//}
