/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/pexip/go-infinity-sdk/v38/config"
)

var (
	_ resource.ResourceWithImportState    = (*InfinityAutobackupResource)(nil)
	_ resource.ResourceWithValidateConfig = (*InfinityAutobackupResource)(nil)
)

type InfinityAutobackupResource struct {
	InfinityClient InfinityClient
}

type InfinityAutobackupResourceModel struct {
	ID                       types.String `tfsdk:"id"`
	AutobackupEnabled        types.Bool   `tfsdk:"autobackup_enabled"`
	AutobackupInterval       types.Int32  `tfsdk:"autobackup_interval"`
	AutobackupPassphrase     types.String `tfsdk:"autobackup_passphrase"`
	AutobackupStartHour      types.Int32  `tfsdk:"autobackup_start_hour"`
	AutobackupUploadURL      types.String `tfsdk:"autobackup_upload_url"`
	AutobackupUploadUsername types.String `tfsdk:"autobackup_upload_username"`
	AutobackupUploadPassword types.String `tfsdk:"autobackup_upload_password"`
}

func (r *InfinityAutobackupResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_infinity_autobackup"
}

func (r *InfinityAutobackupResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *InfinityAutobackupResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages the Pexip Infinity automatic backup configuration. This is a singleton resource — only one instance exists.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Resource URI for the autobackup configuration in Infinity.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"autobackup_enabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Enable automatic creation of daily configuration backups.",
			},
			"autobackup_interval": schema.Int32Attribute{
				Optional: true,
				Computed: true,
				Default:  int32default.StaticInt32(24),
				Validators: []validator.Int32{
					int32validator.Between(1, 24),
				},
				MarkdownDescription: "The number of hours between running an automatic backup. Range: 1 to 24. Default: 24.",
			},
			"autobackup_passphrase": schema.StringAttribute{
				Optional:  true,
				Computed:  true,
				Sensitive: true,
				Default:   stringdefault.StaticString(""),
				Validators: []validator.String{
					stringvalidator.LengthAtMost(100),
				},
				MarkdownDescription: "The passphrase used to encrypt all automatically generated backup files. Maximum length: 100 characters.",
			},
			"autobackup_start_hour": schema.Int32Attribute{
				Optional: true,
				Computed: true,
				Default:  int32default.StaticInt32(1),
				Validators: []validator.Int32{
					int32validator.Between(0, 23),
				},
				MarkdownDescription: "The hour (in UTC time) to run the automatic backup. Range: 0 to 23. Default: 1.",
			},
			"autobackup_upload_url": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(""),
				Validators: []validator.String{
					stringvalidator.LengthAtMost(255),
				},
				MarkdownDescription: "The URL to which to upload automatic backups. Supported schemes: FTPS, FTP. Maximum length: 255 characters.",
			},
			"autobackup_upload_username": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(""),
				Validators: []validator.String{
					stringvalidator.LengthAtMost(100),
				},
				MarkdownDescription: "The username for the upload URL. Maximum length: 100 characters.",
			},
			"autobackup_upload_password": schema.StringAttribute{
				Optional:  true,
				Computed:  true,
				Sensitive: true,
				Default:   stringdefault.StaticString(""),
				Validators: []validator.String{
					stringvalidator.LengthAtMost(100),
				},
				MarkdownDescription: "The password for the upload URL. Maximum length: 100 characters.",
			},
		},
	}
}

func (r *InfinityAutobackupResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var data InfinityAutobackupResourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.AutobackupEnabled.ValueBool() && data.AutobackupPassphrase.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("autobackup_passphrase"),
			"Missing Required Attribute",
			"autobackup_passphrase must be set when autobackup_enabled is true.",
		)
	}
}

func (r *InfinityAutobackupResource) buildUpdateRequest(plan *InfinityAutobackupResourceModel) *config.AutobackupUpdateRequest {
	enabled := plan.AutobackupEnabled.ValueBool()
	interval := int(plan.AutobackupInterval.ValueInt32())
	startHour := int(plan.AutobackupStartHour.ValueInt32())

	return &config.AutobackupUpdateRequest{
		AutobackupEnabled:        &enabled,
		AutobackupInterval:       &interval,
		AutobackupPassphrase:     plan.AutobackupPassphrase.ValueString(),
		AutobackupStartHour:      &startHour,
		AutobackupUploadURL:      plan.AutobackupUploadURL.ValueString(),
		AutobackupUploadUsername: plan.AutobackupUploadUsername.ValueString(),
		AutobackupUploadPassword: plan.AutobackupUploadPassword.ValueString(),
	}
}

func (r *InfinityAutobackupResource) read(ctx context.Context, passphrase, uploadPassword string) (*InfinityAutobackupResourceModel, error) {
	var data InfinityAutobackupResourceModel

	srv, err := r.InfinityClient.Config().GetAutobackup(ctx)
	if err != nil {
		return nil, err
	}

	if srv.ResourceURI == "" {
		return nil, fmt.Errorf("autobackup configuration not found")
	}

	data.ID = types.StringValue(srv.ResourceURI)
	data.AutobackupEnabled = types.BoolValue(srv.AutobackupEnabled)
	data.AutobackupInterval = types.Int32Value(int32(srv.AutobackupInterval))
	// Passphrase is write-only — carry the value from plan/state
	data.AutobackupPassphrase = types.StringValue(passphrase)
	data.AutobackupStartHour = types.Int32Value(int32(srv.AutobackupStartHour))
	data.AutobackupUploadURL = types.StringValue(srv.AutobackupUploadURL)
	data.AutobackupUploadUsername = types.StringValue(srv.AutobackupUploadUsername)
	// Upload password is write-only — carry the value from plan/state
	data.AutobackupUploadPassword = types.StringValue(uploadPassword)

	return &data, nil
}

func (r *InfinityAutobackupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// For singleton resources, Create is actually Update since the resource always exists
	plan := &InfinityAutobackupResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateRequest := r.buildUpdateRequest(plan)

	_, err := r.InfinityClient.Config().UpdateAutobackup(ctx, updateRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Infinity autobackup configuration",
			fmt.Sprintf("Could not update Infinity autobackup configuration: %s", err),
		)
		return
	}

	updatedModel, err := r.read(ctx, plan.AutobackupPassphrase.ValueString(), plan.AutobackupUploadPassword.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Updated Infinity autobackup configuration",
			fmt.Sprintf("Could not read updated Infinity autobackup configuration: %s", err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, updatedModel)...)
}

func (r *InfinityAutobackupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	state := &InfinityAutobackupResourceModel{}

	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	state, err := r.read(ctx, state.AutobackupPassphrase.ValueString(), state.AutobackupUploadPassword.ValueString())
	if err != nil {
		if isNotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading Infinity autobackup configuration",
			fmt.Sprintf("Could not read Infinity autobackup configuration: %s", err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *InfinityAutobackupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	plan := &InfinityAutobackupResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateRequest := r.buildUpdateRequest(plan)

	_, err := r.InfinityClient.Config().UpdateAutobackup(ctx, updateRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Infinity autobackup configuration",
			fmt.Sprintf("Could not update Infinity autobackup configuration: %s", err),
		)
		return
	}

	updatedModel, err := r.read(ctx, plan.AutobackupPassphrase.ValueString(), plan.AutobackupUploadPassword.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Updated Infinity autobackup configuration",
			fmt.Sprintf("Could not read updated Infinity autobackup configuration: %s", err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, updatedModel)...)
}

func (r *InfinityAutobackupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Info(ctx, "Resetting Infinity autobackup configuration to defaults")

	enabled := false
	interval := 24
	startHour := 1

	updateRequest := &config.AutobackupUpdateRequest{
		AutobackupEnabled:        &enabled,
		AutobackupInterval:       &interval,
		AutobackupPassphrase:     "",
		AutobackupStartHour:      &startHour,
		AutobackupUploadURL:      "",
		AutobackupUploadUsername: "",
		AutobackupUploadPassword: "",
	}

	_, err := r.InfinityClient.Config().UpdateAutobackup(ctx, updateRequest)
	if err != nil && !isNotFoundError(err) && !isLookupError(err) {
		resp.Diagnostics.AddError(
			"Error Resetting Infinity autobackup configuration",
			fmt.Sprintf("Could not reset Infinity autobackup configuration: %s", err),
		)
		return
	}
}

func (r *InfinityAutobackupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// For singleton resources, the import ID doesn't matter since there's only one instance
	tflog.Trace(ctx, "Importing Infinity autobackup configuration")

	model, err := r.read(ctx, "", "")
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Importing Infinity Autobackup Configuration",
			fmt.Sprintf("Could not import Infinity autobackup configuration: %s", err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}
