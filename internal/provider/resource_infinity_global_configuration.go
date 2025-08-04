package provider

import (
	"context"
	"fmt"

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
	_ resource.ResourceWithImportState = (*InfinityGlobalConfigurationResource)(nil)
)

type InfinityGlobalConfigurationResource struct {
	InfinityClient InfinityClient
}

type InfinityGlobalConfigurationResourceModel struct {
	ID                           types.String `tfsdk:"id"`
	EnableWebRTC                 types.Bool   `tfsdk:"enable_webrtc"`
	EnableSIP                    types.Bool   `tfsdk:"enable_sip"`
	EnableH323                   types.Bool   `tfsdk:"enable_h323"`
	EnableRTMP                   types.Bool   `tfsdk:"enable_rtmp"`
	CryptoMode                   types.String `tfsdk:"crypto_mode"`
	MaxPixelsPerSecond           types.String `tfsdk:"max_pixels_per_second"`
	MediaPortsStart              types.Int64  `tfsdk:"media_ports_start"`
	MediaPortsEnd                types.Int64  `tfsdk:"media_ports_end"`
	SignallingPortsStart         types.Int64  `tfsdk:"signalling_ports_start"`
	SignallingPortsEnd           types.Int64  `tfsdk:"signalling_ports_end"`
	BurstingEnabled              types.Bool   `tfsdk:"bursting_enabled"`
	CloudProvider                types.String `tfsdk:"cloud_provider"`
	AWSAccessKey                 types.String `tfsdk:"aws_access_key"`
	AWSSecretKey                 types.String `tfsdk:"aws_secret_key"`
	AzureClientID                types.String `tfsdk:"azure_client_id"`
	AzureSecret                  types.String `tfsdk:"azure_secret"`
	GuestsOnlyTimeout            types.Int64  `tfsdk:"guests_only_timeout"`
	WaitingForChairTimeout       types.Int64  `tfsdk:"waiting_for_chair_timeout"`
	ConferenceCreatePermissions  types.String `tfsdk:"conference_create_permissions"`
	ConferenceCreationMode       types.String `tfsdk:"conference_creation_mode"`
	EnableAnalytics              types.Bool   `tfsdk:"enable_analytics"`
	EnableErrorReporting         types.Bool   `tfsdk:"enable_error_reporting"`
	BandwidthRestrictions        types.String `tfsdk:"bandwidth_restrictions"`
	AdministratorEmail           types.String `tfsdk:"administrator_email"`
	GlobalConferenceCreateGroups types.Set    `tfsdk:"global_conference_create_groups"`
}

func (r *InfinityGlobalConfigurationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_infinity_global_configuration"
}

func (r *InfinityGlobalConfigurationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *InfinityGlobalConfigurationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Resource URI for the global configuration in Infinity",
			},
			"enable_webrtc": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Whether to enable WebRTC protocol support.",
			},
			"enable_sip": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Whether to enable SIP protocol support.",
			},
			"enable_h323": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Whether to enable H.323 protocol support.",
			},
			"enable_rtmp": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Whether to enable RTMP protocol support.",
			},
			"crypto_mode": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.OneOf("disabled", "besteffort", "required"),
				},
				MarkdownDescription: "Cryptographic mode for conferences. Valid values: disabled, besteffort, required.",
			},
			"max_pixels_per_second": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Maximum pixels per second for video transmission.",
			},
			"media_ports_start": schema.Int64Attribute{
				Optional: true,
				Validators: []validator.Int64{
					int64validator.Between(1024, 65535),
				},
				MarkdownDescription: "Starting port for media traffic. Valid range: 1024-65535.",
			},
			"media_ports_end": schema.Int64Attribute{
				Optional: true,
				Validators: []validator.Int64{
					int64validator.Between(1024, 65535),
				},
				MarkdownDescription: "Ending port for media traffic. Valid range: 1024-65535.",
			},
			"signalling_ports_start": schema.Int64Attribute{
				Optional: true,
				Validators: []validator.Int64{
					int64validator.Between(1024, 65535),
				},
				MarkdownDescription: "Starting port for signalling traffic. Valid range: 1024-65535.",
			},
			"signalling_ports_end": schema.Int64Attribute{
				Optional: true,
				Validators: []validator.Int64{
					int64validator.Between(1024, 65535),
				},
				MarkdownDescription: "Ending port for signalling traffic. Valid range: 1024-65535.",
			},
			"bursting_enabled": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Whether to enable cloud bursting functionality.",
			},
			"cloud_provider": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.OneOf("aws", "azure", "google"),
				},
				MarkdownDescription: "Cloud provider for bursting. Valid values: aws, azure, google.",
			},
			"aws_access_key": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				MarkdownDescription: "AWS access key for cloud bursting. This field is sensitive.",
			},
			"aws_secret_key": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				MarkdownDescription: "AWS secret key for cloud bursting. This field is sensitive.",
			},
			"azure_client_id": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				MarkdownDescription: "Azure client ID for cloud bursting. This field is sensitive.",
			},
			"azure_secret": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				MarkdownDescription: "Azure secret for cloud bursting. This field is sensitive.",
			},
			"guests_only_timeout": schema.Int64Attribute{
				Optional: true,
				Validators: []validator.Int64{
					int64validator.Between(0, 1440),
				},
				MarkdownDescription: "Timeout in minutes for guests-only conferences. Valid range: 0-1440.",
			},
			"waiting_for_chair_timeout": schema.Int64Attribute{
				Optional: true,
				Validators: []validator.Int64{
					int64validator.Between(0, 1440),
				},
				MarkdownDescription: "Timeout in minutes when waiting for chair to join. Valid range: 0-1440.",
			},
			"conference_create_permissions": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.OneOf("none", "admin_only", "user_admin", "any_authenticated"),
				},
				MarkdownDescription: "Who can create conferences. Valid values: none, admin_only, user_admin, any_authenticated.",
			},
			"conference_creation_mode": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.OneOf("disabled", "per_node", "per_cluster"),
				},
				MarkdownDescription: "Conference creation mode. Valid values: disabled, per_node, per_cluster.",
			},
			"enable_analytics": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Whether to enable analytics collection.",
			},
			"enable_error_reporting": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Whether to enable automatic error reporting.",
			},
			"bandwidth_restrictions": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.OneOf("none", "restricted"),
				},
				MarkdownDescription: "Bandwidth restriction mode. Valid values: none, restricted.",
			},
			"administrator_email": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					validators.Email(),
				},
				MarkdownDescription: "Administrator email address for system notifications.",
			},
			"global_conference_create_groups": schema.SetAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				MarkdownDescription: "List of groups that can create conferences globally.",
			},
		},
		MarkdownDescription: "Manages the global system configuration with the Infinity service. This is a singleton resource - only one global configuration exists per system.",
	}
}

func (r *InfinityGlobalConfigurationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// For singleton resources, Create is actually Update since the resource always exists
	plan := &InfinityGlobalConfigurationResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateRequest := &config.GlobalConfigurationUpdateRequest{
		CryptoMode:                  plan.CryptoMode.ValueString(),
		MaxPixelsPerSecond:          plan.MaxPixelsPerSecond.ValueString(),
		CloudProvider:               plan.CloudProvider.ValueString(),
		ConferenceCreatePermissions: plan.ConferenceCreatePermissions.ValueString(),
		ConferenceCreationMode:      plan.ConferenceCreationMode.ValueString(),
		BandwidthRestrictions:       plan.BandwidthRestrictions.ValueString(),
		AdministratorEmail:          plan.AdministratorEmail.ValueString(),
	}

	// Handle boolean fields
	if !plan.EnableWebRTC.IsNull() {
		enable := plan.EnableWebRTC.ValueBool()
		updateRequest.EnableWebRTC = &enable
	}

	if !plan.EnableSIP.IsNull() {
		enable := plan.EnableSIP.ValueBool()
		updateRequest.EnableSIP = &enable
	}

	if !plan.EnableH323.IsNull() {
		enable := plan.EnableH323.ValueBool()
		updateRequest.EnableH323 = &enable
	}

	if !plan.EnableRTMP.IsNull() {
		enable := plan.EnableRTMP.ValueBool()
		updateRequest.EnableRTMP = &enable
	}

	if !plan.BurstingEnabled.IsNull() {
		enable := plan.BurstingEnabled.ValueBool()
		updateRequest.BurstingEnabled = &enable
	}

	if !plan.EnableAnalytics.IsNull() {
		enable := plan.EnableAnalytics.ValueBool()
		updateRequest.EnableAnalytics = &enable
	}

	if !plan.EnableErrorReporting.IsNull() {
		enable := plan.EnableErrorReporting.ValueBool()
		updateRequest.EnableErrorReporting = &enable
	}

	// Handle integer fields
	if !plan.MediaPortsStart.IsNull() {
		port := int(plan.MediaPortsStart.ValueInt64())
		updateRequest.MediaPortsStart = &port
	}

	if !plan.MediaPortsEnd.IsNull() {
		port := int(plan.MediaPortsEnd.ValueInt64())
		updateRequest.MediaPortsEnd = &port
	}

	if !plan.SignallingPortsStart.IsNull() {
		port := int(plan.SignallingPortsStart.ValueInt64())
		updateRequest.SignallingPortsStart = &port
	}

	if !plan.SignallingPortsEnd.IsNull() {
		port := int(plan.SignallingPortsEnd.ValueInt64())
		updateRequest.SignallingPortsEnd = &port
	}

	if !plan.GuestsOnlyTimeout.IsNull() {
		timeout := int(plan.GuestsOnlyTimeout.ValueInt64())
		updateRequest.GuestsOnlyTimeout = &timeout
	}

	if !plan.WaitingForChairTimeout.IsNull() {
		timeout := int(plan.WaitingForChairTimeout.ValueInt64())
		updateRequest.WaitingForChairTimeout = &timeout
	}

	// Handle sensitive string fields
	if !plan.AWSAccessKey.IsNull() {
		key := plan.AWSAccessKey.ValueString()
		updateRequest.AWSAccessKey = &key
	}

	if !plan.AWSSecretKey.IsNull() {
		secret := plan.AWSSecretKey.ValueString()
		updateRequest.AWSSecretKey = &secret
	}

	if !plan.AzureClientID.IsNull() {
		clientID := plan.AzureClientID.ValueString()
		updateRequest.AzureClientID = &clientID
	}

	if !plan.AzureSecret.IsNull() {
		secret := plan.AzureSecret.ValueString()
		updateRequest.AzureSecret = &secret
	}

	// Handle list field
	if !plan.GlobalConferenceCreateGroups.IsNull() {
		var groups []string
		resp.Diagnostics.Append(plan.GlobalConferenceCreateGroups.ElementsAs(ctx, &groups, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		updateRequest.GlobalConferenceCreateGroups = groups
	}

	_, err := r.InfinityClient.Config().UpdateGlobalConfiguration(ctx, updateRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Infinity global configuration",
			fmt.Sprintf("Could not create Infinity global configuration: %s", err),
		)
		return
	}

	// Read the current state from the API to get all computed values
	model, err := r.read(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Created Infinity global configuration",
			fmt.Sprintf("Could not read created Infinity global configuration: %s", err),
		)
		return
	}
	tflog.Trace(ctx, fmt.Sprintf("created Infinity global configuration with ID: %s", model.ID))

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *InfinityGlobalConfigurationResource) read(ctx context.Context) (*InfinityGlobalConfigurationResourceModel, error) {
	var data InfinityGlobalConfigurationResourceModel

	srv, err := r.InfinityClient.Config().GetGlobalConfiguration(ctx)
	if err != nil {
		return nil, err
	}

	if len(srv.ResourceURI) == 0 {
		return nil, fmt.Errorf("global configuration not found")
	}

	data.ID = types.StringValue(srv.ResourceURI)
	data.EnableWebRTC = types.BoolValue(srv.EnableWebRTC)
	data.EnableSIP = types.BoolValue(srv.EnableSIP)
	data.EnableH323 = types.BoolValue(srv.EnableH323)
	data.EnableRTMP = types.BoolValue(srv.EnableRTMP)
	data.CryptoMode = types.StringValue(srv.CryptoMode)
	data.MaxPixelsPerSecond = types.StringValue(srv.MaxPixelsPerSecond)
	data.MediaPortsStart = types.Int64Value(int64(srv.MediaPortsStart))
	data.MediaPortsEnd = types.Int64Value(int64(srv.MediaPortsEnd))
	data.SignallingPortsStart = types.Int64Value(int64(srv.SignallingPortsStart))
	data.SignallingPortsEnd = types.Int64Value(int64(srv.SignallingPortsEnd))
	data.BurstingEnabled = types.BoolValue(srv.BurstingEnabled)
	data.CloudProvider = types.StringValue(srv.CloudProvider)
	data.GuestsOnlyTimeout = types.Int64Value(int64(srv.GuestsOnlyTimeout))
	data.WaitingForChairTimeout = types.Int64Value(int64(srv.WaitingForChairTimeout))
	data.ConferenceCreatePermissions = types.StringValue(srv.ConferenceCreatePermissions)
	data.ConferenceCreationMode = types.StringValue(srv.ConferenceCreationMode)
	data.EnableAnalytics = types.BoolValue(srv.EnableAnalytics)
	data.EnableErrorReporting = types.BoolValue(srv.EnableErrorReporting)
	data.BandwidthRestrictions = types.StringValue(srv.BandwidthRestrictions)
	data.AdministratorEmail = types.StringValue(srv.AdministratorEmail)

	// Handle optional pointer fields
	if srv.AWSAccessKey != nil {
		data.AWSAccessKey = types.StringValue(*srv.AWSAccessKey)
	} else {
		data.AWSAccessKey = types.StringNull()
	}

	if srv.AWSSecretKey != nil {
		data.AWSSecretKey = types.StringValue(*srv.AWSSecretKey)
	} else {
		data.AWSSecretKey = types.StringNull()
	}

	if srv.AzureClientID != nil {
		data.AzureClientID = types.StringValue(*srv.AzureClientID)
	} else {
		data.AzureClientID = types.StringNull()
	}

	if srv.AzureSecret != nil {
		data.AzureSecret = types.StringValue(*srv.AzureSecret)
	} else {
		data.AzureSecret = types.StringNull()
	}

	// Handle list field
	if srv.GlobalConferenceCreateGroups != nil {
		groupSet, diags := types.SetValueFrom(ctx, types.StringType, srv.GlobalConferenceCreateGroups)
		if diags.HasError() {
			return nil, fmt.Errorf("failed to convert global conference create groups: %s", diags.Errors())
		}
		data.GlobalConferenceCreateGroups = groupSet
	} else {
		data.GlobalConferenceCreateGroups = types.SetNull(types.StringType)
	}

	return &data, nil
}

func (r *InfinityGlobalConfigurationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	state, err := r.read(ctx)
	if err != nil {
		// Check if the error is a 404 (not found) - unlikely for singleton resources
		if isNotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading Infinity global configuration",
			fmt.Sprintf("Could not read Infinity global configuration: %s", err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *InfinityGlobalConfigurationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	plan := &InfinityGlobalConfigurationResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateRequest := &config.GlobalConfigurationUpdateRequest{
		CryptoMode:                  plan.CryptoMode.ValueString(),
		MaxPixelsPerSecond:          plan.MaxPixelsPerSecond.ValueString(),
		CloudProvider:               plan.CloudProvider.ValueString(),
		ConferenceCreatePermissions: plan.ConferenceCreatePermissions.ValueString(),
		ConferenceCreationMode:      plan.ConferenceCreationMode.ValueString(),
		BandwidthRestrictions:       plan.BandwidthRestrictions.ValueString(),
		AdministratorEmail:          plan.AdministratorEmail.ValueString(),
	}

	// Handle boolean fields
	if !plan.EnableWebRTC.IsNull() {
		enable := plan.EnableWebRTC.ValueBool()
		updateRequest.EnableWebRTC = &enable
	}

	if !plan.EnableSIP.IsNull() {
		enable := plan.EnableSIP.ValueBool()
		updateRequest.EnableSIP = &enable
	}

	if !plan.EnableH323.IsNull() {
		enable := plan.EnableH323.ValueBool()
		updateRequest.EnableH323 = &enable
	}

	if !plan.EnableRTMP.IsNull() {
		enable := plan.EnableRTMP.ValueBool()
		updateRequest.EnableRTMP = &enable
	}

	if !plan.BurstingEnabled.IsNull() {
		enable := plan.BurstingEnabled.ValueBool()
		updateRequest.BurstingEnabled = &enable
	}

	if !plan.EnableAnalytics.IsNull() {
		enable := plan.EnableAnalytics.ValueBool()
		updateRequest.EnableAnalytics = &enable
	}

	if !plan.EnableErrorReporting.IsNull() {
		enable := plan.EnableErrorReporting.ValueBool()
		updateRequest.EnableErrorReporting = &enable
	}

	// Handle integer fields
	if !plan.MediaPortsStart.IsNull() {
		port := int(plan.MediaPortsStart.ValueInt64())
		updateRequest.MediaPortsStart = &port
	}

	if !plan.MediaPortsEnd.IsNull() {
		port := int(plan.MediaPortsEnd.ValueInt64())
		updateRequest.MediaPortsEnd = &port
	}

	if !plan.SignallingPortsStart.IsNull() {
		port := int(plan.SignallingPortsStart.ValueInt64())
		updateRequest.SignallingPortsStart = &port
	}

	if !plan.SignallingPortsEnd.IsNull() {
		port := int(plan.SignallingPortsEnd.ValueInt64())
		updateRequest.SignallingPortsEnd = &port
	}

	if !plan.GuestsOnlyTimeout.IsNull() {
		timeout := int(plan.GuestsOnlyTimeout.ValueInt64())
		updateRequest.GuestsOnlyTimeout = &timeout
	}

	if !plan.WaitingForChairTimeout.IsNull() {
		timeout := int(plan.WaitingForChairTimeout.ValueInt64())
		updateRequest.WaitingForChairTimeout = &timeout
	}

	// Handle sensitive string fields
	if !plan.AWSAccessKey.IsNull() {
		key := plan.AWSAccessKey.ValueString()
		updateRequest.AWSAccessKey = &key
	}

	if !plan.AWSSecretKey.IsNull() {
		secret := plan.AWSSecretKey.ValueString()
		updateRequest.AWSSecretKey = &secret
	}

	if !plan.AzureClientID.IsNull() {
		clientID := plan.AzureClientID.ValueString()
		updateRequest.AzureClientID = &clientID
	}

	if !plan.AzureSecret.IsNull() {
		secret := plan.AzureSecret.ValueString()
		updateRequest.AzureSecret = &secret
	}

	// Handle list field
	if !plan.GlobalConferenceCreateGroups.IsNull() {
		var groups []string
		resp.Diagnostics.Append(plan.GlobalConferenceCreateGroups.ElementsAs(ctx, &groups, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		updateRequest.GlobalConferenceCreateGroups = groups
	}

	_, err := r.InfinityClient.Config().UpdateGlobalConfiguration(ctx, updateRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Infinity global configuration",
			fmt.Sprintf("Could not update Infinity global configuration: %s", err),
		)
		return
	}

	// Re-read the resource to get the latest state
	updatedModel, err := r.read(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Updated Infinity global configuration",
			fmt.Sprintf("Could not read updated Infinity global configuration: %s", err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, updatedModel)...)
}

func (r *InfinityGlobalConfigurationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// For singleton resources, delete means resetting to default values
	// We'll set minimal configuration to "delete" the customizations
	tflog.Info(ctx, "Deleting Infinity global configuration (resetting to defaults)")

	updateRequest := &config.GlobalConfigurationUpdateRequest{
		EnableWebRTC:    func() *bool { v := false; return &v }(),
		EnableSIP:       func() *bool { v := false; return &v }(),
		EnableH323:      func() *bool { v := false; return &v }(),
		EnableRTMP:      func() *bool { v := false; return &v }(),
		CryptoMode:      "disabled",
		CloudProvider:   "",
		BurstingEnabled: func() *bool { v := false; return &v }(),
	}

	_, err := r.InfinityClient.Config().UpdateGlobalConfiguration(ctx, updateRequest)
	if err != nil && !isNotFoundError(err) && !isLookupError(err) {
		resp.Diagnostics.AddError(
			"Error Deleting Infinity global configuration",
			fmt.Sprintf("Could not delete Infinity global configuration: %s", err),
		)
		return
	}
}

func (r *InfinityGlobalConfigurationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// For singleton resources, the import ID doesn't matter since there's only one instance
	tflog.Trace(ctx, "Importing Infinity global configuration")

	// Read the resource from the API
	model, err := r.read(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Importing Infinity Global Configuration",
			fmt.Sprintf("Could not import Infinity global configuration: %s", err),
		)
		return
	}

	// Set the state from the imported resource
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}
