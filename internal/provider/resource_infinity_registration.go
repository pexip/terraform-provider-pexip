/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

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
)

var (
	_ resource.ResourceWithImportState = (*InfinityRegistrationResource)(nil)
)

type InfinityRegistrationResource struct {
	InfinityClient InfinityClient
}

type InfinityRegistrationResourceModel struct {
	ID                         types.String `tfsdk:"id"`
	Enable                     types.Bool   `tfsdk:"enable"`
	RefreshStrategy            types.String `tfsdk:"refresh_strategy"`
	AdaptiveMinRefresh         types.Int64  `tfsdk:"adaptive_min_refresh"`
	AdaptiveMaxRefresh         types.Int64  `tfsdk:"adaptive_max_refresh"`
	MaximumMinRefresh          types.Int64  `tfsdk:"maximum_min_refresh"`
	MaximumMaxRefresh          types.Int64  `tfsdk:"maximum_max_refresh"`
	NattedMinRefresh           types.Int64  `tfsdk:"natted_min_refresh"`
	NattedMaxRefresh           types.Int64  `tfsdk:"natted_max_refresh"`
	RouteViaRegistrar          types.Bool   `tfsdk:"route_via_registrar"`
	EnablePushNotifications    types.Bool   `tfsdk:"enable_push_notifications"`
	EnableGoogleCloudMessaging types.Bool   `tfsdk:"enable_google_cloud_messaging"`
	PushToken                  types.String `tfsdk:"push_token"`
}

func (r *InfinityRegistrationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_infinity_registration"
}

func (r *InfinityRegistrationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *InfinityRegistrationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Resource URI for the registration configuration in Infinity",
			},
			"enable": schema.BoolAttribute{
				Required:            true,
				MarkdownDescription: "Whether to enable registration functionality.",
			},
			"refresh_strategy": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.OneOf("adaptive", "maximum", "natted"),
				},
				MarkdownDescription: "The refresh strategy to use. Valid values: adaptive, maximum, natted.",
			},
			"adaptive_min_refresh": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				Validators: []validator.Int64{
					int64validator.Between(30, 3600),
				},
				MarkdownDescription: "Minimum refresh interval for adaptive strategy in seconds. Valid range: 30-3600.",
			},
			"adaptive_max_refresh": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				Validators: []validator.Int64{
					int64validator.Between(30, 3600),
				},
				MarkdownDescription: "Maximum refresh interval for adaptive strategy in seconds. Valid range: 30-3600.",
			},
			"maximum_min_refresh": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				Validators: []validator.Int64{
					int64validator.Between(30, 3600),
				},
				MarkdownDescription: "Minimum refresh interval for maximum strategy in seconds. Valid range: 30-3600.",
			},
			"maximum_max_refresh": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				Validators: []validator.Int64{
					int64validator.Between(30, 3600),
				},
				MarkdownDescription: "Maximum refresh interval for maximum strategy in seconds. Valid range: 30-3600.",
			},
			"natted_min_refresh": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				Validators: []validator.Int64{
					int64validator.Between(30, 3600),
				},
				MarkdownDescription: "Minimum refresh interval for NATted connections in seconds. Valid range: 30-3600.",
			},
			"natted_max_refresh": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				Validators: []validator.Int64{
					int64validator.Between(30, 3600),
				},
				MarkdownDescription: "Maximum refresh interval for NATted connections in seconds. Valid range: 30-3600.",
			},
			"route_via_registrar": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Whether to route calls via the registrar.",
			},
			"enable_push_notifications": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Whether to enable push notifications for mobile clients.",
			},
			"enable_google_cloud_messaging": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Whether to enable Google Cloud Messaging for push notifications.",
			},
			"push_token": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				MarkdownDescription: "Push notification token for mobile clients. This field is sensitive.",
			},
		},
		MarkdownDescription: "Manages the registration configuration with the Infinity service. This is a singleton resource - only one registration configuration exists per system.",
	}
}

func (r *InfinityRegistrationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// For singleton resources, Create is actually Update since the resource always exists
	plan := &InfinityRegistrationResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateRequest := &config.RegistrationUpdateRequest{
		RefreshStrategy: plan.RefreshStrategy.ValueString(),
		PushToken:       plan.PushToken.ValueString(),
	}

	if !plan.Enable.IsNull() {
		enable := plan.Enable.ValueBool()
		updateRequest.Enable = &enable
	}

	if !plan.AdaptiveMinRefresh.IsNull() {
		refresh := int(plan.AdaptiveMinRefresh.ValueInt64())
		updateRequest.AdaptiveMinRefresh = &refresh
	}

	if !plan.AdaptiveMaxRefresh.IsNull() {
		refresh := int(plan.AdaptiveMaxRefresh.ValueInt64())
		updateRequest.AdaptiveMaxRefresh = &refresh
	}

	if !plan.MaximumMinRefresh.IsNull() {
		refresh := int(plan.MaximumMinRefresh.ValueInt64())
		updateRequest.MaximumMinRefresh = &refresh
	}

	if !plan.MaximumMaxRefresh.IsNull() {
		refresh := int(plan.MaximumMaxRefresh.ValueInt64())
		updateRequest.MaximumMaxRefresh = &refresh
	}

	if !plan.NattedMinRefresh.IsNull() {
		refresh := int(plan.NattedMinRefresh.ValueInt64())
		updateRequest.NattedMinRefresh = &refresh
	}

	if !plan.NattedMaxRefresh.IsNull() {
		refresh := int(plan.NattedMaxRefresh.ValueInt64())
		updateRequest.NattedMaxRefresh = &refresh
	}

	if !plan.RouteViaRegistrar.IsNull() {
		route := plan.RouteViaRegistrar.ValueBool()
		updateRequest.RouteViaRegistrar = &route
	}

	if !plan.EnablePushNotifications.IsNull() {
		enable := plan.EnablePushNotifications.ValueBool()
		updateRequest.EnablePushNotifications = &enable
	}

	if !plan.EnableGoogleCloudMessaging.IsNull() {
		enable := plan.EnableGoogleCloudMessaging.ValueBool()
		updateRequest.EnableGoogleCloudMessaging = &enable
	}

	_, err := r.InfinityClient.Config().UpdateRegistration(ctx, updateRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Infinity registration configuration",
			fmt.Sprintf("Could not create Infinity registration configuration: %s", err),
		)
		return
	}

	// Read the current state from the API to get all computed values
	model, err := r.read(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Created Infinity registration configuration",
			fmt.Sprintf("Could not read created Infinity registration configuration: %s", err),
		)
		return
	}
	tflog.Trace(ctx, fmt.Sprintf("created Infinity registration configuration with ID: %s", model.ID))

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *InfinityRegistrationResource) read(ctx context.Context) (*InfinityRegistrationResourceModel, error) {
	var data InfinityRegistrationResourceModel

	srv, err := r.InfinityClient.Config().GetRegistration(ctx)
	if err != nil {
		return nil, err
	}

	if srv.ResourceURI == "" {
		return nil, fmt.Errorf("registration configuration not found")
	}

	data.ID = types.StringValue(srv.ResourceURI)
	data.Enable = types.BoolValue(srv.Enable)
	data.RefreshStrategy = types.StringValue(srv.RefreshStrategy)

	// Set refresh values based on strategy, null for unused strategies
	switch srv.RefreshStrategy {
	case "adaptive":
		data.AdaptiveMinRefresh = types.Int64Value(int64(srv.AdaptiveMinRefresh))
		data.AdaptiveMaxRefresh = types.Int64Value(int64(srv.AdaptiveMaxRefresh))
		data.MaximumMinRefresh = types.Int64Null()
		data.MaximumMaxRefresh = types.Int64Null()
	case "maximum":
		data.MaximumMinRefresh = types.Int64Value(int64(srv.MaximumMinRefresh))
		data.MaximumMaxRefresh = types.Int64Value(int64(srv.MaximumMaxRefresh))
		data.AdaptiveMinRefresh = types.Int64Null()
		data.AdaptiveMaxRefresh = types.Int64Null()
	default:
		// For other strategies or when not set, all strategy-specific fields are null
		data.AdaptiveMinRefresh = types.Int64Null()
		data.AdaptiveMaxRefresh = types.Int64Null()
		data.MaximumMinRefresh = types.Int64Null()
		data.MaximumMaxRefresh = types.Int64Null()
	}

	// Natted fields are strategy-independent but may be null if not configured
	if srv.NattedMinRefresh > 0 {
		data.NattedMinRefresh = types.Int64Value(int64(srv.NattedMinRefresh))
	} else {
		data.NattedMinRefresh = types.Int64Null()
	}
	if srv.NattedMaxRefresh > 0 {
		data.NattedMaxRefresh = types.Int64Value(int64(srv.NattedMaxRefresh))
	} else {
		data.NattedMaxRefresh = types.Int64Null()
	}

	data.RouteViaRegistrar = types.BoolValue(srv.RouteViaRegistrar)
	data.EnablePushNotifications = types.BoolValue(srv.EnablePushNotifications)
	data.EnableGoogleCloudMessaging = types.BoolValue(srv.EnableGoogleCloudMessaging)
	data.PushToken = types.StringValue(srv.PushToken)

	return &data, nil
}

func (r *InfinityRegistrationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	state, err := r.read(ctx)
	if err != nil {
		// Check if the error is a 404 (not found) - unlikely for singleton resources
		if isNotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading Infinity registration configuration",
			fmt.Sprintf("Could not read Infinity registration configuration: %s", err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *InfinityRegistrationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	plan := &InfinityRegistrationResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateRequest := &config.RegistrationUpdateRequest{
		RefreshStrategy: plan.RefreshStrategy.ValueString(),
		PushToken:       plan.PushToken.ValueString(),
	}

	if !plan.Enable.IsNull() {
		enable := plan.Enable.ValueBool()
		updateRequest.Enable = &enable
	}

	if !plan.AdaptiveMinRefresh.IsNull() {
		refresh := int(plan.AdaptiveMinRefresh.ValueInt64())
		updateRequest.AdaptiveMinRefresh = &refresh
	}

	if !plan.AdaptiveMaxRefresh.IsNull() {
		refresh := int(plan.AdaptiveMaxRefresh.ValueInt64())
		updateRequest.AdaptiveMaxRefresh = &refresh
	}

	if !plan.MaximumMinRefresh.IsNull() {
		refresh := int(plan.MaximumMinRefresh.ValueInt64())
		updateRequest.MaximumMinRefresh = &refresh
	}

	if !plan.MaximumMaxRefresh.IsNull() {
		refresh := int(plan.MaximumMaxRefresh.ValueInt64())
		updateRequest.MaximumMaxRefresh = &refresh
	}

	if !plan.NattedMinRefresh.IsNull() {
		refresh := int(plan.NattedMinRefresh.ValueInt64())
		updateRequest.NattedMinRefresh = &refresh
	}

	if !plan.NattedMaxRefresh.IsNull() {
		refresh := int(plan.NattedMaxRefresh.ValueInt64())
		updateRequest.NattedMaxRefresh = &refresh
	}

	if !plan.RouteViaRegistrar.IsNull() {
		route := plan.RouteViaRegistrar.ValueBool()
		updateRequest.RouteViaRegistrar = &route
	}

	if !plan.EnablePushNotifications.IsNull() {
		enable := plan.EnablePushNotifications.ValueBool()
		updateRequest.EnablePushNotifications = &enable
	}

	if !plan.EnableGoogleCloudMessaging.IsNull() {
		enable := plan.EnableGoogleCloudMessaging.ValueBool()
		updateRequest.EnableGoogleCloudMessaging = &enable
	}

	_, err := r.InfinityClient.Config().UpdateRegistration(ctx, updateRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Infinity registration configuration",
			fmt.Sprintf("Could not update Infinity registration configuration: %s", err),
		)
		return
	}

	// Re-read the resource to get the latest state
	updatedModel, err := r.read(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Updated Infinity registration configuration",
			fmt.Sprintf("Could not read updated Infinity registration configuration: %s", err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, updatedModel)...)
}

func (r *InfinityRegistrationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// For singleton resources, delete means resetting to default values
	// We'll disable registration to "delete" the configuration
	tflog.Info(ctx, "Deleting Infinity registration configuration (disabling)")

	updateRequest := &config.RegistrationUpdateRequest{
		Enable: func() *bool { v := false; return &v }(),
	}

	_, err := r.InfinityClient.Config().UpdateRegistration(ctx, updateRequest)
	if err != nil && !isNotFoundError(err) && !isLookupError(err) {
		resp.Diagnostics.AddError(
			"Error Deleting Infinity registration configuration",
			fmt.Sprintf("Could not delete Infinity registration configuration: %s", err),
		)
		return
	}
}

func (r *InfinityRegistrationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// For singleton resources, the import ID doesn't matter since there's only one instance
	tflog.Trace(ctx, "Importing Infinity registration configuration")

	// Read the resource from the API
	model, err := r.read(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Importing Infinity Registration Configuration",
			fmt.Sprintf("Could not import Infinity registration configuration: %s", err),
		)
		return
	}

	// Set the state from the imported resource
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}
