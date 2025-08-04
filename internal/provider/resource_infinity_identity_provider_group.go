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
	_ resource.ResourceWithImportState = (*InfinityIdentityProviderGroupResource)(nil)
)

type InfinityIdentityProviderGroupResource struct {
	InfinityClient InfinityClient
}

type InfinityIdentityProviderGroupResourceModel struct {
	ID               types.String `tfsdk:"id"`
	ResourceID       types.Int32  `tfsdk:"resource_id"`
	Name             types.String `tfsdk:"name"`
	Description      types.String `tfsdk:"description"`
	IdentityProvider types.Set    `tfsdk:"identity_provider"`
}

func (r *InfinityIdentityProviderGroupResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_infinity_identity_provider_group"
}

func (r *InfinityIdentityProviderGroupResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *InfinityIdentityProviderGroupResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Resource URI for the identity provider group in Infinity",
			},
			"resource_id": schema.Int32Attribute{
				Computed:            true,
				MarkdownDescription: "The resource integer identifier for the identity provider group in Infinity",
			},
			"name": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					stringvalidator.LengthAtMost(250),
				},
				MarkdownDescription: "The name of the identity provider group. Maximum length: 250 characters.",
			},
			"description": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(500),
				},
				MarkdownDescription: "Description of the identity provider group. Maximum length: 500 characters.",
			},
			"identity_provider": schema.SetAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				MarkdownDescription: "List of identity provider URIs associated with this group.",
			},
		},
		MarkdownDescription: "Manages an identity provider group with the Infinity service. Identity provider groups organize multiple identity providers for authentication and authorization management.",
	}
}

func (r *InfinityIdentityProviderGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	plan := &InfinityIdentityProviderGroupResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createRequest := &config.IdentityProviderGroupCreateRequest{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
	}

	// Handle list field
	if !plan.IdentityProvider.IsNull() {
		var providers []string
		resp.Diagnostics.Append(plan.IdentityProvider.ElementsAs(ctx, &providers, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		createRequest.IdentityProvider = providers
	}

	createResponse, err := r.InfinityClient.Config().CreateIdentityProviderGroup(ctx, createRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Infinity identity provider group",
			fmt.Sprintf("Could not create Infinity identity provider group: %s", err),
		)
		return
	}

	resourceID, err := createResponse.ResourceID()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Retrieving Infinity identity provider group ID",
			fmt.Sprintf("Could not retrieve ID for created Infinity identity provider group: %s", err),
		)
		return
	}

	// Read the state from the API to get all computed values
	model, err := r.read(ctx, resourceID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Created Infinity identity provider group",
			fmt.Sprintf("Could not read created Infinity identity provider group with ID %d: %s", resourceID, err),
		)
		return
	}
	tflog.Trace(ctx, fmt.Sprintf("created Infinity identity provider group with ID: %s, name: %s", model.ID, model.Name))

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *InfinityIdentityProviderGroupResource) read(ctx context.Context, resourceID int) (*InfinityIdentityProviderGroupResourceModel, error) {
	var data InfinityIdentityProviderGroupResourceModel

	srv, err := r.InfinityClient.Config().GetIdentityProviderGroup(ctx, resourceID)
	if err != nil {
		return nil, err
	}

	if len(srv.ResourceURI) == 0 {
		return nil, fmt.Errorf("identity provider group with ID %d not found", resourceID)
	}

	data.ID = types.StringValue(srv.ResourceURI)
	data.ResourceID = types.Int32Value(int32(resourceID))
	data.Name = types.StringValue(srv.Name)
	data.Description = types.StringValue(srv.Description)

	// Handle list field
	if srv.IdentityProvider != nil {
		providersSet, diags := types.SetValueFrom(ctx, types.StringType, srv.IdentityProvider)
		if diags.HasError() {
			return nil, fmt.Errorf("failed to convert identity providers: %s", diags.Errors())
		}
		data.IdentityProvider = providersSet
	} else {
		data.IdentityProvider = types.SetNull(types.StringType)
	}

	return &data, nil
}

func (r *InfinityIdentityProviderGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	state := &InfinityIdentityProviderGroupResourceModel{}

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
			"Error Reading Infinity identity provider group",
			fmt.Sprintf("Could not read Infinity identity provider group: %s", err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *InfinityIdentityProviderGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	plan := &InfinityIdentityProviderGroupResourceModel{}
	state := &InfinityIdentityProviderGroupResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateRequest := &config.IdentityProviderGroupUpdateRequest{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
	}

	// Handle list field
	if !plan.IdentityProvider.IsNull() {
		var providers []string
		resp.Diagnostics.Append(plan.IdentityProvider.ElementsAs(ctx, &providers, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		updateRequest.IdentityProvider = providers
	}

	resourceID := int(state.ResourceID.ValueInt32())
	_, err := r.InfinityClient.Config().UpdateIdentityProviderGroup(ctx, resourceID, updateRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Infinity identity provider group",
			fmt.Sprintf("Could not update Infinity identity provider group: %s", err),
		)
		return
	}

	// Read the state from the API to get all computed values
	model, err := r.read(ctx, resourceID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Updated Infinity identity provider group",
			fmt.Sprintf("Could not read updated Infinity identity provider group with ID %d: %s", resourceID, err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *InfinityIdentityProviderGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	state := &InfinityIdentityProviderGroupResourceModel{}

	tflog.Info(ctx, "Deleting Infinity identity provider group")

	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.InfinityClient.Config().DeleteIdentityProviderGroup(ctx, int(state.ResourceID.ValueInt32()))

	// Ignore 404 Not Found and Lookup errors on delete
	if err != nil && !isNotFoundError(err) && !isLookupError(err) {
		resp.Diagnostics.AddError(
			"Error Deleting Infinity identity provider group",
			fmt.Sprintf("Could not delete Infinity identity provider group with ID %s: %s", state.ID.ValueString(), err),
		)
		return
	}
}

func (r *InfinityIdentityProviderGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resourceID, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Import ID must be a valid integer for the resource ID. Got: %s", req.ID),
		)
		return
	}

	tflog.Trace(ctx, fmt.Sprintf("Importing Infinity identity provider group with resource ID: %d", resourceID))

	// Read the resource from the API
	model, err := r.read(ctx, resourceID)
	if err != nil {
		// Check if the error is a 404 (not found)
		if isNotFoundError(err) {
			resp.Diagnostics.AddError(
				"Infinity Identity Provider Group Not Found",
				fmt.Sprintf("Infinity identity provider group with resource ID %d not found.", resourceID),
			)
			return
		}
		resp.Diagnostics.AddError(
			"Error Importing Infinity Identity Provider Group",
			fmt.Sprintf("Could not import Infinity identity provider group with resource ID %d: %s", resourceID, err),
		)
		return
	}

	// Set the state from the imported resource
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}
