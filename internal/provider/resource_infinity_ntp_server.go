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
	_ resource.ResourceWithImportState = (*InfinityNtpServerResource)(nil)
)

type InfinityNtpServerResource struct {
	InfinityClient InfinityClient
}

type InfinityNtpServerResourceModel struct {
	ID          types.String `tfsdk:"id"`
	ResourceID  types.Int32  `tfsdk:"resource_id"`
	Address     types.String `tfsdk:"address"`
	Description types.String `tfsdk:"description"`
}

func (r *InfinityNtpServerResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_infinity_ntp_server"
}

func (r *InfinityNtpServerResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *InfinityNtpServerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Resource URI for the NTP server in Infinity",
			},
			"resource_id": schema.Int32Attribute{
				Computed:            true,
				MarkdownDescription: "The resource integer identifier for the NTP server in Infinity",
			},
			"description": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(255),
				},
				MarkdownDescription: "A description of the NTP server. Maximum length: 250 characters.",
			},
			"address": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(255),
				},
				MarkdownDescription: "The IP address of the NTP server.",
			},
		},
		MarkdownDescription: "Registers a NTP server with the Infinity service.",
	}
}

func (r *InfinityNtpServerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	plan := &InfinityNtpServerResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createRequest := &config.NTPServerCreateRequest{
		Address:     plan.Address.ValueString(),
		Description: plan.Description.ValueString(),
	}

	createResponse, err := r.InfinityClient.Config().CreateNTPServer(ctx, createRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Infinity NTP server",
			fmt.Sprintf("Could not create Infinity NTP server: %s", err),
		)
		return
	}

	resourceID, err := createResponse.ResourceID()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Retrieving Infinity NTP server ID",
			fmt.Sprintf("Could not retrieve ID for created Infinity NTP server: %s", err),
		)
		return
	}

	plan, err = r.read(ctx, resourceID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Created Infinity NTP server",
			fmt.Sprintf("Could not read created Infinity NTP server with ID %d: %s", resourceID, err),
		)
		return
	}
	tflog.Trace(ctx, fmt.Sprintf("created Infinity NTP server with ID: %s, name: %s", plan.ID, plan.Address))

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *InfinityNtpServerResource) read(ctx context.Context, resourceID int) (*InfinityNtpServerResourceModel, error) {
	var data InfinityNtpServerResourceModel

	srv, err := r.InfinityClient.Config().GetNTPServer(ctx, resourceID)
	if err != nil {
		return nil, err
	}

	if len(srv.ResourceURI) == 0 {
		return nil, fmt.Errorf("NTP server with ID %d not found", resourceID)
	}

	data.ID = types.StringValue(srv.ResourceURI)
	data.ResourceID = types.Int32Value(int32(resourceID))
	data.Address = types.StringValue(srv.Address)
	data.Description = types.StringValue(srv.Description)

	return &data, nil
}

func (r *InfinityNtpServerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	state := &InfinityNtpServerResourceModel{}

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
			"Error Reading Infinity NTP server",
			fmt.Sprintf("Could not read Infinity NTP server: %s", err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *InfinityNtpServerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	plan := &InfinityNtpServerResourceModel{}
	state := &InfinityNtpServerResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// The resource ID is required for the update API call.
	resourceID := int(state.ResourceID.ValueInt32())

	// Prepare the update request from the plan
	updateRequest := &config.NTPServerUpdateRequest{
		Address:     plan.Address.ValueString(),
		Description: plan.Description.ValueString(),
	}
	_, err := r.InfinityClient.Config().UpdateNTPServer(ctx, resourceID, updateRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Infinity NTP server",
			fmt.Sprintf("Could not update Infinity NTP server with ID %s: %s", plan.ID.ValueString(), err),
		)
		return
	}

	plan.ID = state.ID
	plan.ResourceID = state.ResourceID
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *InfinityNtpServerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	state := &InfinityNtpServerResourceModel{}

	tflog.Info(ctx, "Deleting Infinity NTP server")

	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.InfinityClient.Config().DeleteNTPServer(ctx, int(state.ResourceID.ValueInt32()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Infinity NTP server",
			fmt.Sprintf("Could not delete Infinity NTP server with ID %s: %s", state.ID.ValueString(), err),
		)
		return
	}
}

func (r *InfinityNtpServerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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
