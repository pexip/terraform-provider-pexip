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
	_ resource.ResourceWithImportState = (*InfinityEndUserResource)(nil)
)

type InfinityEndUserResource struct {
	InfinityClient InfinityClient
}

type InfinityEndUserResourceModel struct {
	ID                  types.String `tfsdk:"id"`
	ResourceID          types.Int32  `tfsdk:"resource_id"`
	PrimaryEmailAddress types.String `tfsdk:"primary_email_address"`
	FirstName           types.String `tfsdk:"first_name"`
	LastName            types.String `tfsdk:"last_name"`
	DisplayName         types.String `tfsdk:"display_name"`
	TelephoneNumber     types.String `tfsdk:"telephone_number"`
	MobileNumber        types.String `tfsdk:"mobile_number"`
	Title               types.String `tfsdk:"title"`
	Department          types.String `tfsdk:"department"`
	AvatarURL           types.String `tfsdk:"avatar_url"`
	UserGroups          types.Set    `tfsdk:"user_groups"`
	UserOID             types.String `tfsdk:"user_oid"`
	ExchangeUserID      types.String `tfsdk:"exchange_user_id"`
	MSExchangeGUID      types.String `tfsdk:"ms_exchange_guid"`
	SyncTag             types.String `tfsdk:"sync_tag"`
}

func (r *InfinityEndUserResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_infinity_end_user"
}

func (r *InfinityEndUserResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *InfinityEndUserResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Resource URI for the end user in Infinity",
			},
			"resource_id": schema.Int32Attribute{
				Computed:            true,
				MarkdownDescription: "The resource integer identifier for the end user in Infinity",
			},
			"primary_email_address": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(100),
				},
				MarkdownDescription: "The unique primary email address for the end user. Maximum length: 100 characters.",
			},
			"first_name": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(250),
				},
				MarkdownDescription: "The first name of the end user. Maximum length: 250 characters.",
			},
			"last_name": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(250),
				},
				MarkdownDescription: "The last name of the end user. Maximum length: 250 characters.",
			},
			"display_name": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(250),
				},
				MarkdownDescription: "The display name of the end user. Maximum length: 250 characters.",
			},
			"telephone_number": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(100),
				},
				MarkdownDescription: "The telephone number of the end user. Maximum length: 100 characters.",
			},
			"mobile_number": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(100),
				},
				MarkdownDescription: "The mobile number of the end user. Maximum length: 100 characters.",
			},
			"title": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(128),
				},
				MarkdownDescription: "The job title of the end user. Maximum length: 128 characters.",
			},
			"department": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(100),
				},
				MarkdownDescription: "The department of the end user. Maximum length: 100 characters.",
			},
			"avatar_url": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(255),
				},
				MarkdownDescription: "The avatar URL for the end user. Maximum length: 255 characters.",
			},
			"user_groups": schema.SetAttribute{
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "List of user group resource URIs that this user belongs to.",
			},
			"user_oid": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Microsoft 365 Object ID (read-only).",
			},
			"exchange_user_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Exchange User ID (read-only).",
			},
			"ms_exchange_guid": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(100),
				},
				MarkdownDescription: "Exchange Mailbox ID. Maximum length: 100 characters.",
			},
			"sync_tag": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(250),
				},
				MarkdownDescription: "LDAP sync identifier. Maximum length: 250 characters.",
			},
		},
		MarkdownDescription: "Manages an end user account with the Infinity service.",
	}
}

func (r *InfinityEndUserResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	plan := &InfinityEndUserResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createRequest := &config.EndUserCreateRequest{
		PrimaryEmailAddress: plan.PrimaryEmailAddress.ValueString(),
	}

	// Set optional string fields
	if !plan.FirstName.IsNull() {
		createRequest.FirstName = plan.FirstName.ValueString()
	}
	if !plan.LastName.IsNull() {
		createRequest.LastName = plan.LastName.ValueString()
	}
	if !plan.DisplayName.IsNull() {
		createRequest.DisplayName = plan.DisplayName.ValueString()
	}
	if !plan.TelephoneNumber.IsNull() {
		createRequest.TelephoneNumber = plan.TelephoneNumber.ValueString()
	}
	if !plan.MobileNumber.IsNull() {
		createRequest.MobileNumber = plan.MobileNumber.ValueString()
	}
	if !plan.Title.IsNull() {
		createRequest.Title = plan.Title.ValueString()
	}
	if !plan.Department.IsNull() {
		createRequest.Department = plan.Department.ValueString()
	}
	if !plan.AvatarURL.IsNull() {
		createRequest.AvatarURL = plan.AvatarURL.ValueString()
	}
	if !plan.MSExchangeGUID.IsNull() {
		msExchangeGUID := plan.MSExchangeGUID.ValueString()
		createRequest.MSExchangeGUID = &msExchangeGUID
	}
	if !plan.SyncTag.IsNull() {
		createRequest.SyncTag = plan.SyncTag.ValueString()
	}

	// Set user groups if provided
	if !plan.UserGroups.IsNull() && !plan.UserGroups.IsUnknown() {
		var userGroups []string
		resp.Diagnostics.Append(plan.UserGroups.ElementsAs(ctx, &userGroups, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		createRequest.UserGroups = userGroups
	}

	createResponse, err := r.InfinityClient.Config().CreateEndUser(ctx, createRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Infinity end user",
			fmt.Sprintf("Could not create Infinity end user: %s", err),
		)
		return
	}

	resourceID, err := createResponse.ResourceID()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Retrieving Infinity end user ID",
			fmt.Sprintf("Could not retrieve ID for created Infinity end user: %s", err),
		)
		return
	}

	// Read the state from the API to get all computed values
	model, err := r.read(ctx, resourceID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Created Infinity end user",
			fmt.Sprintf("Could not read created Infinity end user with ID %d: %s", resourceID, err),
		)
		return
	}
	tflog.Trace(ctx, fmt.Sprintf("created Infinity end user with ID: %s, email: %s", model.ID, model.PrimaryEmailAddress))

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *InfinityEndUserResource) read(ctx context.Context, resourceID int) (*InfinityEndUserResourceModel, error) {
	var data InfinityEndUserResourceModel

	srv, err := r.InfinityClient.Config().GetEndUser(ctx, resourceID)
	if err != nil {
		return nil, err
	}

	if srv.ResourceURI == "" {
		return nil, fmt.Errorf("end user with ID %d not found", resourceID)
	}

	data.ID = types.StringValue(srv.ResourceURI)
	data.ResourceID = types.Int32Value(int32(resourceID)) // #nosec G115 -- API values are expected to be within int32 range
	data.PrimaryEmailAddress = types.StringValue(srv.PrimaryEmailAddress)
	data.FirstName = types.StringValue(srv.FirstName)
	data.LastName = types.StringValue(srv.LastName)
	data.DisplayName = types.StringValue(srv.DisplayName)
	data.TelephoneNumber = types.StringValue(srv.TelephoneNumber)
	data.MobileNumber = types.StringValue(srv.MobileNumber)
	data.Title = types.StringValue(srv.Title)
	data.Department = types.StringValue(srv.Department)
	data.AvatarURL = types.StringValue(srv.AvatarURL)
	data.SyncTag = types.StringValue(srv.SyncTag)

	// Handle pointer fields
	if srv.UserOID != nil {
		data.UserOID = types.StringValue(*srv.UserOID)
	} else {
		data.UserOID = types.StringNull()
	}

	if srv.ExchangeUserID != nil {
		data.ExchangeUserID = types.StringValue(*srv.ExchangeUserID)
	} else {
		data.ExchangeUserID = types.StringNull()
	}

	if srv.MSExchangeGUID != nil {
		data.MSExchangeGUID = types.StringValue(*srv.MSExchangeGUID)
	} else {
		data.MSExchangeGUID = types.StringNull()
	}

	// Convert user groups to types.Set
	userGroupsSet, diags := types.SetValueFrom(ctx, types.StringType, srv.UserGroups)
	if diags.HasError() {
		return nil, fmt.Errorf("error converting user groups: %v", diags)
	}
	data.UserGroups = userGroupsSet

	return &data, nil
}

func (r *InfinityEndUserResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	state := &InfinityEndUserResourceModel{}

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
			"Error Reading Infinity end user",
			fmt.Sprintf("Could not read Infinity end user: %s", err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *InfinityEndUserResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	plan := &InfinityEndUserResourceModel{}
	state := &InfinityEndUserResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resourceID := int(state.ResourceID.ValueInt32())

	updateRequest := &config.EndUserUpdateRequest{
		PrimaryEmailAddress: plan.PrimaryEmailAddress.ValueString(),
	}

	// Set optional string fields
	if !plan.FirstName.IsNull() {
		updateRequest.FirstName = plan.FirstName.ValueString()
	}
	if !plan.LastName.IsNull() {
		updateRequest.LastName = plan.LastName.ValueString()
	}
	if !plan.DisplayName.IsNull() {
		updateRequest.DisplayName = plan.DisplayName.ValueString()
	}
	if !plan.TelephoneNumber.IsNull() {
		updateRequest.TelephoneNumber = plan.TelephoneNumber.ValueString()
	}
	if !plan.MobileNumber.IsNull() {
		updateRequest.MobileNumber = plan.MobileNumber.ValueString()
	}
	if !plan.Title.IsNull() {
		updateRequest.Title = plan.Title.ValueString()
	}
	if !plan.Department.IsNull() {
		updateRequest.Department = plan.Department.ValueString()
	}
	if !plan.AvatarURL.IsNull() {
		updateRequest.AvatarURL = plan.AvatarURL.ValueString()
	}
	if !plan.MSExchangeGUID.IsNull() {
		msExchangeGUID := plan.MSExchangeGUID.ValueString()
		updateRequest.MSExchangeGUID = &msExchangeGUID
	}
	if !plan.SyncTag.IsNull() {
		updateRequest.SyncTag = plan.SyncTag.ValueString()
	}

	// Set user groups if provided
	if !plan.UserGroups.IsNull() && !plan.UserGroups.IsUnknown() {
		var userGroups []string
		resp.Diagnostics.Append(plan.UserGroups.ElementsAs(ctx, &userGroups, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		updateRequest.UserGroups = userGroups
	}

	_, err := r.InfinityClient.Config().UpdateEndUser(ctx, resourceID, updateRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Infinity end user",
			fmt.Sprintf("Could not update Infinity end user with ID %d: %s", resourceID, err),
		)
		return
	}

	// Re-read the resource to get the latest state
	updatedModel, err := r.read(ctx, resourceID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Updated Infinity end user",
			fmt.Sprintf("Could not read updated Infinity end user with ID %d: %s", resourceID, err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, updatedModel)...)
}

func (r *InfinityEndUserResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	state := &InfinityEndUserResourceModel{}

	tflog.Info(ctx, "Deleting Infinity end user")

	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.InfinityClient.Config().DeleteEndUser(ctx, int(state.ResourceID.ValueInt32()))

	// Ignore 404 Not Found and Lookup errors on delete
	if err != nil && !isNotFoundError(err) && !isLookupError(err) {
		resp.Diagnostics.AddError(
			"Error Deleting Infinity end user",
			fmt.Sprintf("Could not delete Infinity end user with ID %s: %s", state.ID.ValueString(), err),
		)
		return
	}
}

func (r *InfinityEndUserResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resourceID, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Import ID must be a valid integer for the resource ID. Got: %s", req.ID),
		)
		return
	}

	tflog.Trace(ctx, fmt.Sprintf("Importing Infinity end user with resource ID: %d", resourceID))

	// Read the resource from the API
	model, err := r.read(ctx, resourceID)
	if err != nil {
		// Check if the error is a 404 (not found)
		if isNotFoundError(err) {
			resp.Diagnostics.AddError(
				"Infinity End User Not Found",
				fmt.Sprintf("Infinity end user with resource ID %d not found.", resourceID),
			)
			return
		}
		resp.Diagnostics.AddError(
			"Error Importing Infinity End User",
			fmt.Sprintf("Could not import Infinity end user with resource ID %d: %s", resourceID, err),
		)
		return
	}

	// Set the state from the imported resource
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}
