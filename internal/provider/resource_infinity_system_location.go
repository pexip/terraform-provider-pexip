package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"sort"
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
	_ resource.ResourceWithImportState = (*InfinitySystemLocationResource)(nil)
)

type InfinitySystemLocationResource struct {
	InfinityClient InfinityClient
}

type InfinitySystemLocationResourceModel struct {
	ID            types.String `tfsdk:"id"`
	ResourceID    types.Int32  `tfsdk:"resource_id"`
	Name          types.String `tfsdk:"name"`
	Description   types.String `tfsdk:"description"`
	DNSServers    types.List   `tfsdk:"dns_servers"`
	NTPServers    types.List   `tfsdk:"ntp_servers"`
	MTU           types.Int32  `tfsdk:"mtu"`
	SyslogServers types.List   `tfsdk:"syslog_servers"`
}

// getSortedStringList is a generic helper to convert a types.List of strings to a sorted string slice.
func getSortedStringList(ctx context.Context, list types.List) ([]string, diag.Diagnostics) {
	if list.IsNull() || list.IsUnknown() {
		return nil, nil
	}
	var items []string
	diags := list.ElementsAs(ctx, &items, false)
	if diags.HasError() {
		return nil, diags
	}
	sort.Strings(items)
	return items, diags
}

func (m *InfinitySystemLocationResourceModel) GetDNSServers(ctx context.Context) ([]string, diag.Diagnostics) {
	return getSortedStringList(ctx, m.DNSServers)
}

func (m *InfinitySystemLocationResourceModel) GetNTPServers(ctx context.Context) ([]string, diag.Diagnostics) {
	return getSortedStringList(ctx, m.NTPServers)
}

func (m *InfinitySystemLocationResourceModel) GetSyslogServers(ctx context.Context) ([]string, diag.Diagnostics) {
	return getSortedStringList(ctx, m.SyslogServers)
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
			"dns_servers": schema.ListAttribute{
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "List of DNS server resource URIs for this system location.",
			},
			"ntp_servers": schema.ListAttribute{
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "List of NTP server resource URIs for this system location.",
			},
			"mtu": schema.Int32Attribute{
				Optional:            true,
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
			"syslog_servers": schema.ListAttribute{
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "The Syslog servers to be used by Conferencing Nodes deployed in this Location.",
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

	dnsServers, diags := plan.GetDNSServers(ctx)
	resp.Diagnostics.Append(diags...)
	ntpServers, diags := plan.GetNTPServers(ctx)
	resp.Diagnostics.Append(diags...)
	syslogServers, diags := plan.GetSyslogServers(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createRequest := &config.SystemLocationCreateRequest{
		Name:          plan.Name.ValueString(),
		DNSServers:    dnsServers,
		NTPServers:    ntpServers,
		SyslogServers: syslogServers,
	}

	// Only set optional fields if they are not null in the plan
	if !plan.Description.IsNull() {
		createRequest.Description = plan.Description.ValueString()
	}
	if !plan.MTU.IsNull() {
		createRequest.MTU = int(plan.MTU.ValueInt32())
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

	// Read the state from the API to get all computed values
	model, err := r.read(ctx, resourceID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Created Infinity system location",
			fmt.Sprintf("Could not read created Infinity system location with ID %d: %s", resourceID, err),
		)
		return
	}
	tflog.Trace(ctx, fmt.Sprintf("created Infinity system location with ID: %s, name: %s", model.ID, model.Name))

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
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
	data.MTU = types.Int32Value(int32(srv.MTU))
	data.Name = types.StringValue(srv.Name)

	// Convert DNS servers from SDK to Terraform format
	var dnsServers []string
	for _, dns := range srv.DNSServers {
		dnsServers = append(dnsServers, fmt.Sprintf("/api/admin/configuration/v1/dns_server/%d/", dns.ID))
	}
	sort.Strings(dnsServers)
	dnsListValue, diags := types.ListValueFrom(ctx, types.StringType, dnsServers)
	if diags.HasError() {
		return nil, fmt.Errorf("error converting DNS servers: %v", diags)
	}
	data.DNSServers = dnsListValue

	// Convert NTP servers from SDK to Terraform format
	var ntpServers []string
	for _, ntp := range srv.NTPServers {
		ntpServers = append(ntpServers, fmt.Sprintf("/api/admin/configuration/v1/ntp_server/%d/", ntp.ID))
	}
	sort.Strings(ntpServers)
	ntpListValue, diags := types.ListValueFrom(ctx, types.StringType, ntpServers)
	if diags.HasError() {
		return nil, fmt.Errorf("error converting NTP servers: %v", diags)
	}
	data.NTPServers = ntpListValue

	// Convert Syslog servers from SDK to Terraform format
	var syslogServers []string
	for _, syslog := range srv.SyslogServers {
		syslogServers = append(syslogServers, fmt.Sprintf("/api/admin/configuration/v1/syslog_server/%d/", syslog.ID))
	}
	sort.Strings(syslogServers)
	syslogListValue, diags := types.ListValueFrom(ctx, types.StringType, syslogServers)
	if diags.HasError() {
		return nil, fmt.Errorf("error converting Syslog servers: %v", diags)
	}
	data.SyslogServers = syslogListValue

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

	resourceID := int(state.ResourceID.ValueInt32())

	dnsServers, diags := plan.GetDNSServers(ctx)
	resp.Diagnostics.Append(diags...)
	ntpServers, diags := plan.GetNTPServers(ctx)
	resp.Diagnostics.Append(diags...)
	syslogServers, diags := plan.GetSyslogServers(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateRequest := &config.SystemLocationUpdateRequest{
		Name:          plan.Name.ValueString(),
		DNSServers:    dnsServers,
		NTPServers:    ntpServers,
		SyslogServers: syslogServers,
	}

	if !plan.Description.IsNull() {
		updateRequest.Description = plan.Description.ValueString()
	}
	if !plan.MTU.IsNull() {
		updateRequest.MTU = int(plan.MTU.ValueInt32())
	}

	_, err := r.InfinityClient.Config().UpdateSystemLocation(ctx, resourceID, updateRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Infinity system location",
			fmt.Sprintf("Could not update Infinity system location with ID %d: %s", resourceID, err),
		)
		return
	}

	// Re-read the resource to get the latest state
	updatedModel, err := r.read(ctx, resourceID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Updated Infinity system location",
			fmt.Sprintf("Could not read updated Infinity system location with ID %d: %s", resourceID, err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, updatedModel)...)
}

func (r *InfinitySystemLocationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	state := &InfinitySystemLocationResourceModel{}

	tflog.Info(ctx, "Deleting Infinity system location")

	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.InfinityClient.Config().DeleteSystemLocation(ctx, int(state.ResourceID.ValueInt32()))

	// Ignore 404 Not Found and Lookup errors on delete
	if err != nil && !isNotFoundError(err) && !isLookupError(err) {
		resp.Diagnostics.AddError(
			"Error Deleting Infinity system location",
			fmt.Sprintf("Could not delete Infinity system location with ID %s: %s", state.ID.ValueString(), err),
		)
		return
	}
}

func (r *InfinitySystemLocationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resourceID, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Import ID must be a valid integer for the resource ID. Got: %s", req.ID),
		)
		return
	}

	tflog.Trace(ctx, fmt.Sprintf("Importing Infinity system location with resource ID: %d", resourceID))

	// Read the resource from the API
	model, err := r.read(ctx, resourceID)
	if err != nil {
		// Check if the error is a 404 (not found)
		if isNotFoundError(err) {
			resp.Diagnostics.AddError(
				"Infinity System Location Not Found",
				fmt.Sprintf("Infinity system location with resource ID %d not found.", resourceID),
			)
			return
		}
		resp.Diagnostics.AddError(
			"Error Importing Infinity System Location",
			fmt.Sprintf("Could not import Infinity system location with resource ID %d: %s", resourceID, err),
		)
		return
	}

	// Set the state from the imported resource
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}
