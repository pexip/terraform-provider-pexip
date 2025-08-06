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

	"github.com/pexip/terraform-provider-pexip/internal/provider/validators"
)

var (
	_ resource.ResourceWithImportState = (*InfinitySMTPServerResource)(nil)
)

type InfinitySMTPServerResource struct {
	InfinityClient InfinityClient
}

type InfinitySMTPServerResourceModel struct {
	ID                 types.String `tfsdk:"id"`
	ResourceID         types.Int32  `tfsdk:"resource_id"`
	Name               types.String `tfsdk:"name"`
	Description        types.String `tfsdk:"description"`
	Address            types.String `tfsdk:"address"`
	Port               types.Int64  `tfsdk:"port"`
	Username           types.String `tfsdk:"username"`
	Password           types.String `tfsdk:"password"`
	FromEmailAddress   types.String `tfsdk:"from_email_address"`
	ConnectionSecurity types.String `tfsdk:"connection_security"`
}

func (r *InfinitySMTPServerResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_infinity_smtp_server"
}

func (r *InfinitySMTPServerResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *InfinitySMTPServerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Resource URI for the SMTP server in Infinity",
			},
			"resource_id": schema.Int32Attribute{
				Computed:            true,
				MarkdownDescription: "The resource integer identifier for the SMTP server in Infinity",
			},
			"name": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					stringvalidator.LengthAtMost(250),
				},
				MarkdownDescription: "The name of the SMTP server. Maximum length: 250 characters.",
			},
			"description": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(500),
				},
				MarkdownDescription: "Description of the SMTP server. Maximum length: 500 characters.",
			},
			"address": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
				MarkdownDescription: "The IP address or hostname of the SMTP server.",
			},
			"port": schema.Int64Attribute{
				Required: true,
				Validators: []validator.Int64{
					int64validator.Between(1, 65535),
				},
				MarkdownDescription: "The port number for SMTP communications. Valid range: 1-65535.",
			},
			"username": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Username for SMTP authentication (optional).",
			},
			"password": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				MarkdownDescription: "Password for SMTP authentication (optional). This field is sensitive.",
			},
			"from_email_address": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					validators.Email(),
				},
				MarkdownDescription: "The from email address used when sending emails through this SMTP server.",
			},
			"connection_security": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf("none", "starttls", "ssl_tls"),
				},
				MarkdownDescription: "Connection security method for SMTP. Valid values: none, starttls, ssl_tls.",
			},
		},
		MarkdownDescription: "Manages an SMTP server with the Infinity service. SMTP servers are used for sending email notifications and alerts from Pexip Infinity.",
	}
}

func (r *InfinitySMTPServerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	plan := &InfinitySMTPServerResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createRequest := &config.SMTPServerCreateRequest{
		Name:               plan.Name.ValueString(),
		Description:        plan.Description.ValueString(),
		Address:            plan.Address.ValueString(),
		Port:               int(plan.Port.ValueInt64()),
		Username:           plan.Username.ValueString(),
		Password:           plan.Password.ValueString(),
		FromEmailAddress:   plan.FromEmailAddress.ValueString(),
		ConnectionSecurity: plan.ConnectionSecurity.ValueString(),
	}

	createResponse, err := r.InfinityClient.Config().CreateSMTPServer(ctx, createRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Infinity SMTP server",
			fmt.Sprintf("Could not create Infinity SMTP server: %s", err),
		)
		return
	}

	resourceID, err := createResponse.ResourceID()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Retrieving Infinity SMTP server ID",
			fmt.Sprintf("Could not retrieve ID for created Infinity SMTP server: %s", err),
		)
		return
	}

	// Read the state from the API to get all computed values
	model, err := r.read(ctx, resourceID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Created Infinity SMTP server",
			fmt.Sprintf("Could not read created Infinity SMTP server with ID %d: %s", resourceID, err),
		)
		return
	}
	tflog.Trace(ctx, fmt.Sprintf("created Infinity SMTP server with ID: %s, name: %s", model.ID, model.Name))

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *InfinitySMTPServerResource) read(ctx context.Context, resourceID int) (*InfinitySMTPServerResourceModel, error) {
	var data InfinitySMTPServerResourceModel

	srv, err := r.InfinityClient.Config().GetSMTPServer(ctx, resourceID)
	if err != nil {
		return nil, err
	}

	if srv.ResourceURI == "" {
		return nil, fmt.Errorf("SMTP server with ID %d not found", resourceID)
	}

	data.ID = types.StringValue(srv.ResourceURI)
	data.ResourceID = types.Int32Value(int32(resourceID)) // #nosec G115 -- API values are expected to be within int32 range
	data.Name = types.StringValue(srv.Name)
	data.Description = types.StringValue(srv.Description)
	data.Address = types.StringValue(srv.Address)
	data.Port = types.Int64Value(int64(srv.Port))
	data.Username = types.StringValue(srv.Username)
	data.Password = types.StringValue(srv.Password)
	data.FromEmailAddress = types.StringValue(srv.FromEmailAddress)
	data.ConnectionSecurity = types.StringValue(srv.ConnectionSecurity)

	return &data, nil
}

func (r *InfinitySMTPServerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	state := &InfinitySMTPServerResourceModel{}

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
			"Error Reading Infinity SMTP server",
			fmt.Sprintf("Could not read Infinity SMTP server: %s", err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *InfinitySMTPServerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	plan := &InfinitySMTPServerResourceModel{}
	state := &InfinitySMTPServerResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateRequest := &config.SMTPServerUpdateRequest{
		Name:               plan.Name.ValueString(),
		Description:        plan.Description.ValueString(),
		Address:            plan.Address.ValueString(),
		Username:           plan.Username.ValueString(),
		Password:           plan.Password.ValueString(),
		FromEmailAddress:   plan.FromEmailAddress.ValueString(),
		ConnectionSecurity: plan.ConnectionSecurity.ValueString(),
	}

	// Handle optional pointer field for port
	if !plan.Port.IsNull() && !plan.Port.IsUnknown() {
		port := int(plan.Port.ValueInt64())
		updateRequest.Port = &port
	}

	resourceID := int(state.ResourceID.ValueInt32())
	_, err := r.InfinityClient.Config().UpdateSMTPServer(ctx, resourceID, updateRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Infinity SMTP server",
			fmt.Sprintf("Could not update Infinity SMTP server: %s", err),
		)
		return
	}

	// Read the state from the API to get all computed values
	model, err := r.read(ctx, resourceID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Updated Infinity SMTP server",
			fmt.Sprintf("Could not read updated Infinity SMTP server with ID %d: %s", resourceID, err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *InfinitySMTPServerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	state := &InfinitySMTPServerResourceModel{}

	tflog.Info(ctx, "Deleting Infinity SMTP server")

	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.InfinityClient.Config().DeleteSMTPServer(ctx, int(state.ResourceID.ValueInt32()))

	// Ignore 404 Not Found and Lookup errors on delete
	if err != nil && !isNotFoundError(err) && !isLookupError(err) {
		resp.Diagnostics.AddError(
			"Error Deleting Infinity SMTP server",
			fmt.Sprintf("Could not delete Infinity SMTP server with ID %s: %s", state.ID.ValueString(), err),
		)
		return
	}
}

func (r *InfinitySMTPServerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resourceID, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Import ID must be a valid integer for the resource ID. Got: %s", req.ID),
		)
		return
	}

	tflog.Trace(ctx, fmt.Sprintf("Importing Infinity SMTP server with resource ID: %d", resourceID))

	// Read the resource from the API
	model, err := r.read(ctx, resourceID)
	if err != nil {
		// Check if the error is a 404 (not found)
		if isNotFoundError(err) {
			resp.Diagnostics.AddError(
				"Infinity SMTP Server Not Found",
				fmt.Sprintf("Infinity SMTP server with resource ID %d not found.", resourceID),
			)
			return
		}
		resp.Diagnostics.AddError(
			"Error Importing Infinity SMTP Server",
			fmt.Sprintf("Could not import Infinity SMTP server with resource ID %d: %s", resourceID, err),
		)
		return
	}

	// Set the state from the imported resource
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}
