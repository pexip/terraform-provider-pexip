package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/pexip/go-infinity-sdk/v38/config"
)

var (
	_ resource.ResourceWithImportState = (*InfinityEventSinkResource)(nil)
)

type InfinityEventSinkResource struct {
	InfinityClient InfinityClient
}

type InfinityEventSinkResourceModel struct {
	ID                   types.String `tfsdk:"id"`
	ResourceID           types.Int32  `tfsdk:"resource_id"`
	Name                 types.String `tfsdk:"name"`
	Description          types.String `tfsdk:"description"`
	URL                  types.String `tfsdk:"url"`
	Username             types.String `tfsdk:"username"`
	Password             types.String `tfsdk:"password"`
	BulkSupport          types.Bool   `tfsdk:"bulk_support"`
	VerifyTLSCertificate types.Bool   `tfsdk:"verify_tls_certificate"`
	Version              types.Int32  `tfsdk:"version"`
}

func (r *InfinityEventSinkResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_infinity_event_sink"
}

func (r *InfinityEventSinkResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *InfinityEventSinkResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Resource URI for the event sink in Infinity",
			},
			"resource_id": schema.Int32Attribute{
				Computed:            true,
				MarkdownDescription: "The resource integer identifier for the event sink in Infinity",
			},
			"name": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(250),
				},
				MarkdownDescription: "The name used to refer to this event sink. Maximum length: 250 characters.",
			},
			"description": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(250),
				},
				MarkdownDescription: "A description of the event sink. Maximum length: 250 characters.",
			},
			"url": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(500),
				},
				MarkdownDescription: "The URL for the event sink. Maximum length: 500 characters.",
			},
			"username": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(100),
				},
				MarkdownDescription: "Username for authentication to the event sink. Maximum length: 100 characters.",
			},
			"password": schema.StringAttribute{
				Optional:  true,
				Computed:  true,
				Sensitive: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(100),
				},
				MarkdownDescription: "Password for authentication to the event sink. Maximum length: 100 characters.",
			},
			"bulk_support": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Whether the event sink supports bulk operations.",
			},
			"verify_tls_certificate": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Whether to verify TLS certificates when connecting to the event sink.",
			},
			"version": schema.Int32Attribute{
				Optional: true,
				Computed: true,
				Validators: []validator.Int32{
					int32validator.AtLeast(1),
				},
				Default:             int32default.StaticInt32(1),
				MarkdownDescription: "The version of the event sink API. Must be at least 1.",
			},
		},
		MarkdownDescription: "Manages an event sink configuration with the Infinity service.",
	}
}

func (r *InfinityEventSinkResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	plan := &InfinityEventSinkResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createRequest := &config.EventSinkCreateRequest{
		Name:                 plan.Name.ValueString(),
		URL:                  plan.URL.ValueString(),
		BulkSupport:          plan.BulkSupport.ValueBool(),
		VerifyTLSCertificate: plan.VerifyTLSCertificate.ValueBool(),
		Version:              int(plan.Version.ValueInt32()),
	}

	// Only set optional fields if they are not null in the plan
	if !plan.Description.IsNull() {
		description := plan.Description.ValueString()
		createRequest.Description = &description
	}
	if !plan.Username.IsNull() {
		username := plan.Username.ValueString()
		createRequest.Username = &username
	}
	if !plan.Password.IsNull() {
		password := plan.Password.ValueString()
		createRequest.Password = &password
	}

	createResponse, err := r.InfinityClient.Config().CreateEventSink(ctx, createRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Infinity event sink",
			fmt.Sprintf("Could not create Infinity event sink: %s", err),
		)
		return
	}

	resourceID, err := createResponse.ResourceID()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Retrieving Infinity event sink ID",
			fmt.Sprintf("Could not retrieve ID for created Infinity event sink: %s", err),
		)
		return
	}

	// Read the state from the API to get all computed values
	model, err := r.read(ctx, resourceID, plan.Password.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Created Infinity event sink",
			fmt.Sprintf("Could not read created Infinity event sink with ID %d: %s", resourceID, err),
		)
		return
	}
	tflog.Trace(ctx, fmt.Sprintf("created Infinity event sink with ID: %s, name: %s", model.ID, model.Name))

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *InfinityEventSinkResource) read(ctx context.Context, resourceID int, password string) (*InfinityEventSinkResourceModel, error) {
	var data InfinityEventSinkResourceModel

	srv, err := r.InfinityClient.Config().GetEventSink(ctx, resourceID)
	if err != nil {
		return nil, err
	}

	if len(srv.ResourceURI) == 0 {
		return nil, fmt.Errorf("event sink with ID %d not found", resourceID)
	}

	data.ID = types.StringValue(srv.ResourceURI)
	data.ResourceID = types.Int32Value(int32(resourceID))
	data.Name = types.StringValue(srv.Name)
	data.URL = types.StringValue(srv.URL)
	data.BulkSupport = types.BoolValue(srv.BulkSupport)
	data.VerifyTLSCertificate = types.BoolValue(srv.VerifyTLSCertificate)
	data.Version = types.Int32Value(int32(srv.Version))
	data.Password = types.StringValue(password)

	if srv.Description != nil {
		data.Description = types.StringValue(*srv.Description)
	} else {
		data.Description = types.StringNull()
	}

	if srv.Username != nil {
		data.Username = types.StringValue(*srv.Username)
	} else {
		data.Username = types.StringNull()
	}

	return &data, nil
}

func (r *InfinityEventSinkResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	state := &InfinityEventSinkResourceModel{}

	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resourceID := int(state.ResourceID.ValueInt32())
	state, err := r.read(ctx, resourceID, state.Password.ValueString())
	if err != nil {
		// Check if the error is a 404 (not found)
		if isNotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading Infinity event sink",
			fmt.Sprintf("Could not read Infinity event sink: %s", err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *InfinityEventSinkResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	plan := &InfinityEventSinkResourceModel{}
	state := &InfinityEventSinkResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resourceID := int(state.ResourceID.ValueInt32())

	version := int(plan.Version.ValueInt32())
	bulkSupport := plan.BulkSupport.ValueBool()
	verifyTLS := plan.VerifyTLSCertificate.ValueBool()
	updateRequest := &config.EventSinkUpdateRequest{
		Name:                 plan.Name.ValueString(),
		URL:                  plan.URL.ValueString(),
		BulkSupport:          &bulkSupport,
		VerifyTLSCertificate: &verifyTLS,
		Version:              &version,
	}

	if !plan.Description.IsNull() {
		description := plan.Description.ValueString()
		updateRequest.Description = &description
	}
	if !plan.Username.IsNull() {
		username := plan.Username.ValueString()
		updateRequest.Username = &username
	}
	if !plan.Password.IsNull() {
		password := plan.Password.ValueString()
		updateRequest.Password = &password
	}

	_, err := r.InfinityClient.Config().UpdateEventSink(ctx, resourceID, updateRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Infinity event sink",
			fmt.Sprintf("Could not update Infinity event sink with ID %d: %s", resourceID, err),
		)
		return
	}

	// Re-read the resource to get the latest state
	updatedModel, err := r.read(ctx, resourceID, plan.Password.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Updated Infinity event sink",
			fmt.Sprintf("Could not read updated Infinity event sink with ID %d: %s", resourceID, err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, updatedModel)...)
}

func (r *InfinityEventSinkResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	state := &InfinityEventSinkResourceModel{}

	tflog.Info(ctx, "Deleting Infinity event sink")

	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.InfinityClient.Config().DeleteEventSink(ctx, int(state.ResourceID.ValueInt32()))

	// Ignore 404 Not Found and Lookup errors on delete
	if err != nil && !isNotFoundError(err) && !isLookupError(err) {
		resp.Diagnostics.AddError(
			"Error Deleting Infinity event sink",
			fmt.Sprintf("Could not delete Infinity event sink with ID %s: %s", state.ID.ValueString(), err),
		)
		return
	}
}

func (r *InfinityEventSinkResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resourceID, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Import ID must be a valid integer for the resource ID. Got: %s", req.ID),
		)
		return
	}

	tflog.Trace(ctx, fmt.Sprintf("Importing Infinity event sink with resource ID: %d", resourceID))

	// Read the resource from the API
	model, err := r.read(ctx, resourceID, "")
	if err != nil {
		// Check if the error is a 404 (not found)
		if isNotFoundError(err) {
			resp.Diagnostics.AddError(
				"Infinity Event Sink Not Found",
				fmt.Sprintf("Infinity event sink with resource ID %d not found.", resourceID),
			)
			return
		}
		resp.Diagnostics.AddError(
			"Error Importing Infinity Event Sink",
			fmt.Sprintf("Could not import Infinity event sink with resource ID %d: %s", resourceID, err),
		)
		return
	}

	// Set the state from the imported resource
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}
