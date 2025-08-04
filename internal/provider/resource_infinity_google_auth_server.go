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
	_ resource.ResourceWithImportState = (*InfinityGoogleAuthServerResource)(nil)
)

type InfinityGoogleAuthServerResource struct {
	InfinityClient InfinityClient
}

type InfinityGoogleAuthServerResourceModel struct {
	ID              types.String `tfsdk:"id"`
	ResourceID      types.Int32  `tfsdk:"resource_id"`
	Name            types.String `tfsdk:"name"`
	Description     types.String `tfsdk:"description"`
	ApplicationType types.String `tfsdk:"application_type"`
	ClientID        types.String `tfsdk:"client_id"`
	ClientSecret    types.String `tfsdk:"client_secret"`
}

func (r *InfinityGoogleAuthServerResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_infinity_google_auth_server"
}

func (r *InfinityGoogleAuthServerResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *InfinityGoogleAuthServerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Resource URI for the Google auth server in Infinity",
			},
			"resource_id": schema.Int32Attribute{
				Computed:            true,
				MarkdownDescription: "The resource integer identifier for the Google auth server in Infinity",
			},
			"name": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					stringvalidator.LengthAtMost(100),
				},
				MarkdownDescription: "The name of the Google auth server. Maximum length: 100 characters.",
			},
			"description": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(500),
				},
				MarkdownDescription: "Description of the Google auth server. Maximum length: 500 characters.",
			},
			"application_type": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf("web", "installed"),
				},
				MarkdownDescription: "The Google OAuth 2.0 application type. Valid values: web, installed.",
			},
			"client_id": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(200),
				},
				MarkdownDescription: "The Google OAuth 2.0 client ID. Maximum length: 200 characters.",
			},
			"client_secret": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				MarkdownDescription: "The Google OAuth 2.0 client secret. This field is sensitive.",
			},
		},
		MarkdownDescription: "Manages a Google OAuth 2.0 auth server with the Infinity service. Google auth servers enable OAuth 2.0 authentication integration with Google services for user authentication and authorization within Pexip Infinity.",
	}
}

func (r *InfinityGoogleAuthServerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	plan := &InfinityGoogleAuthServerResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createRequest := &config.GoogleAuthServerCreateRequest{
		Name:            plan.Name.ValueString(),
		Description:     plan.Description.ValueString(),
		ApplicationType: plan.ApplicationType.ValueString(),
		ClientSecret:    plan.ClientSecret.ValueString(),
	}

	// Handle optional pointer field
	if !plan.ClientID.IsNull() && !plan.ClientID.IsUnknown() {
		clientID := plan.ClientID.ValueString()
		createRequest.ClientID = &clientID
	}

	createResponse, err := r.InfinityClient.Config().CreateGoogleAuthServer(ctx, createRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Infinity Google auth server",
			fmt.Sprintf("Could not create Infinity Google auth server: %s", err),
		)
		return
	}

	resourceID, err := createResponse.ResourceID()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Retrieving Infinity Google auth server ID",
			fmt.Sprintf("Could not retrieve ID for created Infinity Google auth server: %s", err),
		)
		return
	}

	// Read the state from the API to get all computed values
	model, err := r.read(ctx, resourceID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Created Infinity Google auth server",
			fmt.Sprintf("Could not read created Infinity Google auth server with ID %d: %s", resourceID, err),
		)
		return
	}
	tflog.Trace(ctx, fmt.Sprintf("created Infinity Google auth server with ID: %s, name: %s", model.ID, model.Name))

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *InfinityGoogleAuthServerResource) read(ctx context.Context, resourceID int) (*InfinityGoogleAuthServerResourceModel, error) {
	var data InfinityGoogleAuthServerResourceModel

	srv, err := r.InfinityClient.Config().GetGoogleAuthServer(ctx, resourceID)
	if err != nil {
		return nil, err
	}

	if len(srv.ResourceURI) == 0 {
		return nil, fmt.Errorf("Google auth server with ID %d not found", resourceID)
	}

	data.ID = types.StringValue(srv.ResourceURI)
	data.ResourceID = types.Int32Value(int32(resourceID))
	data.Name = types.StringValue(srv.Name)
	data.Description = types.StringValue(srv.Description)
	data.ApplicationType = types.StringValue(srv.ApplicationType)
	data.ClientSecret = types.StringValue(srv.ClientSecret)

	// Handle optional pointer field
	if srv.ClientID != nil {
		data.ClientID = types.StringValue(*srv.ClientID)
	} else {
		data.ClientID = types.StringNull()
	}

	return &data, nil
}

func (r *InfinityGoogleAuthServerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	state := &InfinityGoogleAuthServerResourceModel{}

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
			"Error Reading Infinity Google auth server",
			fmt.Sprintf("Could not read Infinity Google auth server: %s", err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *InfinityGoogleAuthServerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	plan := &InfinityGoogleAuthServerResourceModel{}
	state := &InfinityGoogleAuthServerResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateRequest := &config.GoogleAuthServerUpdateRequest{
		Name:            plan.Name.ValueString(),
		Description:     plan.Description.ValueString(),
		ApplicationType: plan.ApplicationType.ValueString(),
		ClientSecret:    plan.ClientSecret.ValueString(),
	}

	// Handle optional pointer field
	if !plan.ClientID.IsNull() && !plan.ClientID.IsUnknown() {
		clientID := plan.ClientID.ValueString()
		updateRequest.ClientID = &clientID
	}

	resourceID := int(state.ResourceID.ValueInt32())
	_, err := r.InfinityClient.Config().UpdateGoogleAuthServer(ctx, resourceID, updateRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Infinity Google auth server",
			fmt.Sprintf("Could not update Infinity Google auth server: %s", err),
		)
		return
	}

	// Read the state from the API to get all computed values
	model, err := r.read(ctx, resourceID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Updated Infinity Google auth server",
			fmt.Sprintf("Could not read updated Infinity Google auth server with ID %d: %s", resourceID, err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *InfinityGoogleAuthServerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	state := &InfinityGoogleAuthServerResourceModel{}

	tflog.Info(ctx, "Deleting Infinity Google auth server")

	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.InfinityClient.Config().DeleteGoogleAuthServer(ctx, int(state.ResourceID.ValueInt32()))

	// Ignore 404 Not Found and Lookup errors on delete
	if err != nil && !isNotFoundError(err) && !isLookupError(err) {
		resp.Diagnostics.AddError(
			"Error Deleting Infinity Google auth server",
			fmt.Sprintf("Could not delete Infinity Google auth server with ID %s: %s", state.ID.ValueString(), err),
		)
		return
	}
}

func (r *InfinityGoogleAuthServerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resourceID, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Import ID must be a valid integer for the resource ID. Got: %s", req.ID),
		)
		return
	}

	tflog.Trace(ctx, fmt.Sprintf("Importing Infinity Google auth server with resource ID: %d", resourceID))

	// Read the resource from the API
	model, err := r.read(ctx, resourceID)
	if err != nil {
		// Check if the error is a 404 (not found)
		if isNotFoundError(err) {
			resp.Diagnostics.AddError(
				"Infinity Google Auth Server Not Found",
				fmt.Sprintf("Infinity Google auth server with resource ID %d not found.", resourceID),
			)
			return
		}
		resp.Diagnostics.AddError(
			"Error Importing Infinity Google Auth Server",
			fmt.Sprintf("Could not import Infinity Google auth server with resource ID %d: %s", resourceID, err),
		)
		return
	}

	// Set the state from the imported resource
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}
