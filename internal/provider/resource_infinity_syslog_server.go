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
	_ resource.ResourceWithImportState = (*InfinitySyslogServerResource)(nil)
)

type InfinitySyslogServerResource struct {
	InfinityClient InfinityClient
}

type InfinitySyslogServerResourceModel struct {
	ID          types.String `tfsdk:"id"`
	ResourceID  types.Int32  `tfsdk:"resource_id"`
	Address     types.String `tfsdk:"address"`
	Description types.String `tfsdk:"description"`
	Port        types.Int64  `tfsdk:"port"`
	Transport   types.String `tfsdk:"transport"`
	AuditLog    types.Bool   `tfsdk:"audit_log"`
	SupportLog  types.Bool   `tfsdk:"support_log"`
	WebLog      types.Bool   `tfsdk:"web_log"`
}

func (r *InfinitySyslogServerResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_infinity_syslog_server"
}

func (r *InfinitySyslogServerResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *InfinitySyslogServerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Resource URI for the syslog server in Infinity",
			},
			"resource_id": schema.Int32Attribute{
				Computed:            true,
				MarkdownDescription: "The resource integer identifier for the syslog server in Infinity",
			},
			"address": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					validators.IPAddress(),
				},
				MarkdownDescription: "The IP address of the syslog server.",
			},
			"description": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(500),
				},
				MarkdownDescription: "Description of the syslog server. Maximum length: 500 characters.",
			},
			"port": schema.Int64Attribute{
				Required: true,
				Validators: []validator.Int64{
					int64validator.Between(1, 65535),
				},
				MarkdownDescription: "The port number for syslog communications. Valid range: 1-65535.",
			},
			"transport": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf("udp", "tcp", "tls"),
				},
				MarkdownDescription: "Transport protocol for syslog. Valid values: udp, tcp, tls.",
			},
			"audit_log": schema.BoolAttribute{
				Required:            true,
				MarkdownDescription: "Whether to send audit logs to this syslog server.",
			},
			"support_log": schema.BoolAttribute{
				Required:            true,
				MarkdownDescription: "Whether to send support logs to this syslog server.",
			},
			"web_log": schema.BoolAttribute{
				Required:            true,
				MarkdownDescription: "Whether to send web logs to this syslog server.",
			},
		},
		MarkdownDescription: "Manages a syslog server with the Infinity service. Syslog servers receive system logs and audit information from Pexip Infinity for centralized logging and monitoring.",
	}
}

func (r *InfinitySyslogServerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	plan := &InfinitySyslogServerResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createRequest := &config.SyslogServerCreateRequest{
		Address:     plan.Address.ValueString(),
		Description: plan.Description.ValueString(),
		Port:        int(plan.Port.ValueInt64()),
		Transport:   plan.Transport.ValueString(),
		AuditLog:    plan.AuditLog.ValueBool(),
		SupportLog:  plan.SupportLog.ValueBool(),
		WebLog:      plan.WebLog.ValueBool(),
	}

	createResponse, err := r.InfinityClient.Config().CreateSyslogServer(ctx, createRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Infinity syslog server",
			fmt.Sprintf("Could not create Infinity syslog server: %s", err),
		)
		return
	}

	resourceID, err := createResponse.ResourceID()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Retrieving Infinity syslog server ID",
			fmt.Sprintf("Could not retrieve ID for created Infinity syslog server: %s", err),
		)
		return
	}

	// Read the state from the API to get all computed values
	model, err := r.read(ctx, resourceID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Created Infinity syslog server",
			fmt.Sprintf("Could not read created Infinity syslog server with ID %d: %s", resourceID, err),
		)
		return
	}
	tflog.Trace(ctx, fmt.Sprintf("created Infinity syslog server with ID: %s, address: %s", model.ID, model.Address))

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *InfinitySyslogServerResource) read(ctx context.Context, resourceID int) (*InfinitySyslogServerResourceModel, error) {
	var data InfinitySyslogServerResourceModel

	srv, err := r.InfinityClient.Config().GetSyslogServer(ctx, resourceID)
	if err != nil {
		return nil, err
	}

	if len(srv.ResourceURI) == 0 {
		return nil, fmt.Errorf("syslog server with ID %d not found", resourceID)
	}

	data.ID = types.StringValue(srv.ResourceURI)
	data.ResourceID = types.Int32Value(int32(resourceID))
	data.Address = types.StringValue(srv.Address)
	data.Description = types.StringValue(srv.Description)
	data.Port = types.Int64Value(int64(srv.Port))
	data.Transport = types.StringValue(srv.Transport)
	data.AuditLog = types.BoolValue(srv.AuditLog)
	data.SupportLog = types.BoolValue(srv.SupportLog)
	data.WebLog = types.BoolValue(srv.WebLog)

	return &data, nil
}

func (r *InfinitySyslogServerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	state := &InfinitySyslogServerResourceModel{}

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
			"Error Reading Infinity syslog server",
			fmt.Sprintf("Could not read Infinity syslog server: %s", err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *InfinitySyslogServerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	plan := &InfinitySyslogServerResourceModel{}
	state := &InfinitySyslogServerResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateRequest := &config.SyslogServerUpdateRequest{
		Address:     plan.Address.ValueString(),
		Description: plan.Description.ValueString(),
		Port:        int(plan.Port.ValueInt64()),
		Transport:   plan.Transport.ValueString(),
	}

	// Handle optional pointer fields for booleans
	if !plan.AuditLog.IsNull() && !plan.AuditLog.IsUnknown() {
		auditLog := plan.AuditLog.ValueBool()
		updateRequest.AuditLog = &auditLog
	}

	if !plan.SupportLog.IsNull() && !plan.SupportLog.IsUnknown() {
		supportLog := plan.SupportLog.ValueBool()
		updateRequest.SupportLog = &supportLog
	}

	if !plan.WebLog.IsNull() && !plan.WebLog.IsUnknown() {
		webLog := plan.WebLog.ValueBool()
		updateRequest.WebLog = &webLog
	}

	resourceID := int(state.ResourceID.ValueInt32())
	_, err := r.InfinityClient.Config().UpdateSyslogServer(ctx, resourceID, updateRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Infinity syslog server",
			fmt.Sprintf("Could not update Infinity syslog server: %s", err),
		)
		return
	}

	// Read the state from the API to get all computed values
	model, err := r.read(ctx, resourceID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Updated Infinity syslog server",
			fmt.Sprintf("Could not read updated Infinity syslog server with ID %d: %s", resourceID, err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *InfinitySyslogServerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	state := &InfinitySyslogServerResourceModel{}

	tflog.Info(ctx, "Deleting Infinity syslog server")

	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.InfinityClient.Config().DeleteSyslogServer(ctx, int(state.ResourceID.ValueInt32()))

	// Ignore 404 Not Found and Lookup errors on delete
	if err != nil && !isNotFoundError(err) && !isLookupError(err) {
		resp.Diagnostics.AddError(
			"Error Deleting Infinity syslog server",
			fmt.Sprintf("Could not delete Infinity syslog server with ID %s: %s", state.ID.ValueString(), err),
		)
		return
	}
}

func (r *InfinitySyslogServerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resourceID, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Import ID must be a valid integer for the resource ID. Got: %s", req.ID),
		)
		return
	}

	tflog.Trace(ctx, fmt.Sprintf("Importing Infinity syslog server with resource ID: %d", resourceID))

	// Read the resource from the API
	model, err := r.read(ctx, resourceID)
	if err != nil {
		// Check if the error is a 404 (not found)
		if isNotFoundError(err) {
			resp.Diagnostics.AddError(
				"Infinity Syslog Server Not Found",
				fmt.Sprintf("Infinity syslog server with resource ID %d not found.", resourceID),
			)
			return
		}
		resp.Diagnostics.AddError(
			"Error Importing Infinity Syslog Server",
			fmt.Sprintf("Could not import Infinity syslog server with resource ID %d: %s", resourceID, err),
		)
		return
	}

	// Set the state from the imported resource
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}
