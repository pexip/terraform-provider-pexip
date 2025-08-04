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
)

var (
	_ resource.ResourceWithImportState = (*InfinityMediaProcessingServerResource)(nil)
)

type InfinityMediaProcessingServerResource struct {
	InfinityClient InfinityClient
}

type InfinityMediaProcessingServerResourceModel struct {
	ID           types.String `tfsdk:"id"`
	ResourceID   types.Int32  `tfsdk:"resource_id"`
	FQDN         types.String `tfsdk:"fqdn"`
	AppID        types.String `tfsdk:"app_id"`
	PublicJWTKey types.String `tfsdk:"public_jwt_key"`
}

func (r *InfinityMediaProcessingServerResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_infinity_media_processing_server"
}

func (r *InfinityMediaProcessingServerResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *InfinityMediaProcessingServerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Resource URI for the media processing server in Infinity",
			},
			"resource_id": schema.Int32Attribute{
				Computed:            true,
				MarkdownDescription: "The resource integer identifier for the media processing server in Infinity",
			},
			"fqdn": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					stringvalidator.LengthAtMost(253),
				},
				MarkdownDescription: "The fully qualified domain name (FQDN) of the media processing server. Maximum length: 253 characters.",
			},
			"app_id": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(100),
				},
				MarkdownDescription: "The application ID for the media processing server. Maximum length: 100 characters.",
			},
			"public_jwt_key": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
				MarkdownDescription: "The public JWT key used for authentication with the media processing server.",
			},
		},
		MarkdownDescription: "Manages a media processing server with the Infinity service. Media processing servers provide advanced media handling capabilities for conferencing, such as transcoding, recording, and streaming services that extend beyond the core Pexip Infinity functionality.",
	}
}

func (r *InfinityMediaProcessingServerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	plan := &InfinityMediaProcessingServerResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createRequest := &config.MediaProcessingServerCreateRequest{
		FQDN:         plan.FQDN.ValueString(),
		PublicJWTKey: plan.PublicJWTKey.ValueString(),
	}

	// Handle optional pointer field
	if !plan.AppID.IsNull() && !plan.AppID.IsUnknown() {
		appID := plan.AppID.ValueString()
		createRequest.AppID = &appID
	}

	createResponse, err := r.InfinityClient.Config().CreateMediaProcessingServer(ctx, createRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Infinity media processing server",
			fmt.Sprintf("Could not create Infinity media processing server: %s", err),
		)
		return
	}

	resourceID, err := createResponse.ResourceID()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Retrieving Infinity media processing server ID",
			fmt.Sprintf("Could not retrieve ID for created Infinity media processing server: %s", err),
		)
		return
	}

	// Read the state from the API to get all computed values
	model, err := r.read(ctx, resourceID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Created Infinity media processing server",
			fmt.Sprintf("Could not read created Infinity media processing server with ID %d: %s", resourceID, err),
		)
		return
	}
	tflog.Trace(ctx, fmt.Sprintf("created Infinity media processing server with ID: %s, FQDN: %s", model.ID, model.FQDN))

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *InfinityMediaProcessingServerResource) read(ctx context.Context, resourceID int) (*InfinityMediaProcessingServerResourceModel, error) {
	var data InfinityMediaProcessingServerResourceModel

	srv, err := r.InfinityClient.Config().GetMediaProcessingServer(ctx, resourceID)
	if err != nil {
		return nil, err
	}

	if len(srv.ResourceURI) == 0 {
		return nil, fmt.Errorf("media processing server with ID %d not found", resourceID)
	}

	data.ID = types.StringValue(srv.ResourceURI)
	data.ResourceID = types.Int32Value(int32(resourceID))
	data.FQDN = types.StringValue(srv.FQDN)
	data.PublicJWTKey = types.StringValue(srv.PublicJWTKey)

	// Handle optional pointer field
	if srv.AppID != nil {
		data.AppID = types.StringValue(*srv.AppID)
	} else {
		data.AppID = types.StringNull()
	}

	return &data, nil
}

func (r *InfinityMediaProcessingServerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	state := &InfinityMediaProcessingServerResourceModel{}

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
			"Error Reading Infinity media processing server",
			fmt.Sprintf("Could not read Infinity media processing server: %s", err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *InfinityMediaProcessingServerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	plan := &InfinityMediaProcessingServerResourceModel{}
	state := &InfinityMediaProcessingServerResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateRequest := &config.MediaProcessingServerUpdateRequest{
		FQDN:         plan.FQDN.ValueString(),
		PublicJWTKey: plan.PublicJWTKey.ValueString(),
	}

	// Handle optional pointer field
	if !plan.AppID.IsNull() && !plan.AppID.IsUnknown() {
		appID := plan.AppID.ValueString()
		updateRequest.AppID = &appID
	}

	resourceID := int(state.ResourceID.ValueInt32())
	_, err := r.InfinityClient.Config().UpdateMediaProcessingServer(ctx, resourceID, updateRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Infinity media processing server",
			fmt.Sprintf("Could not update Infinity media processing server: %s", err),
		)
		return
	}

	// Read the state from the API to get all computed values
	model, err := r.read(ctx, resourceID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Updated Infinity media processing server",
			fmt.Sprintf("Could not read updated Infinity media processing server with ID %d: %s", resourceID, err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *InfinityMediaProcessingServerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	state := &InfinityMediaProcessingServerResourceModel{}

	tflog.Info(ctx, "Deleting Infinity media processing server")

	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.InfinityClient.Config().DeleteMediaProcessingServer(ctx, int(state.ResourceID.ValueInt32()))

	// Ignore 404 Not Found and Lookup errors on delete
	if err != nil && !isNotFoundError(err) && !isLookupError(err) {
		resp.Diagnostics.AddError(
			"Error Deleting Infinity media processing server",
			fmt.Sprintf("Could not delete Infinity media processing server with ID %s: %s", state.ID.ValueString(), err),
		)
		return
	}
}

func (r *InfinityMediaProcessingServerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resourceID, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Import ID must be a valid integer for the resource ID. Got: %s", req.ID),
		)
		return
	}

	tflog.Trace(ctx, fmt.Sprintf("Importing Infinity media processing server with resource ID: %d", resourceID))

	// Read the resource from the API
	model, err := r.read(ctx, resourceID)
	if err != nil {
		// Check if the error is a 404 (not found)
		if isNotFoundError(err) {
			resp.Diagnostics.AddError(
				"Infinity Media Processing Server Not Found",
				fmt.Sprintf("Infinity media processing server with resource ID %d not found.", resourceID),
			)
			return
		}
		resp.Diagnostics.AddError(
			"Error Importing Infinity Media Processing Server",
			fmt.Sprintf("Could not import Infinity media processing server with resource ID %d: %s", resourceID, err),
		)
		return
	}

	// Set the state from the imported resource
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}
