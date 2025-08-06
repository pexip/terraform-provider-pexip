/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/pexip/go-infinity-sdk/v38/config"
)

var (
	_ resource.Resource = (*InfinityUpgradeResource)(nil)
)

type InfinityUpgradeResource struct {
	InfinityClient InfinityClient
}

type InfinityUpgradeResourceModel struct {
	ID        types.String `tfsdk:"id"`
	Package   types.String `tfsdk:"package"`
	Timestamp types.String `tfsdk:"timestamp"`
}

func (r *InfinityUpgradeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_infinity_upgrade"
}

func (r *InfinityUpgradeResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *InfinityUpgradeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Unique identifier for this upgrade trigger",
			},
			"package": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Specific upgrade package to use. If not specified, the system will use the default upgrade package.",
			},
			"timestamp": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Timestamp when the upgrade was triggered",
			},
		},
		MarkdownDescription: "Triggers a system upgrade on the Infinity service. This is an action resource that initiates an upgrade process. Note: This resource only supports creation and reading - upgrades cannot be updated or undone once initiated. The resource represents the upgrade trigger action, not the upgrade state itself.",
	}
}

func (r *InfinityUpgradeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	plan := &InfinityUpgradeResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createRequest := &config.UpgradeCreateRequest{}

	// Handle optional package field
	if !plan.Package.IsNull() && !plan.Package.IsUnknown() {
		pkg := plan.Package.ValueString()
		createRequest.Package = &pkg
	}

	createResponse, err := r.InfinityClient.Config().CreateUpgrade(ctx, createRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Triggering Infinity upgrade",
			fmt.Sprintf("Could not trigger Infinity upgrade: %s", err),
		)
		return
	}

	// Generate a unique ID for this upgrade trigger
	timestamp := time.Now().UTC()
	upgradeID := fmt.Sprintf("upgrade_%d", timestamp.Unix())

	// Since this is an action resource, we create a simple state representation
	model := &InfinityUpgradeResourceModel{
		ID:        types.StringValue(upgradeID),
		Timestamp: types.StringValue(timestamp.Format(time.RFC3339)),
	}

	if !plan.Package.IsNull() {
		model.Package = plan.Package
	} else {
		model.Package = types.StringNull()
	}

	// Log the upgrade trigger
	tflog.Info(ctx, fmt.Sprintf("triggered Infinity upgrade with ID: %s", upgradeID))
	if createResponse != nil {
		tflog.Trace(ctx, fmt.Sprintf("upgrade response: %+v", createResponse))
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *InfinityUpgradeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	state := &InfinityUpgradeResourceModel{}

	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// For an action resource like upgrade, we just maintain the state as-is
	// The upgrade may have completed, but the trigger record remains valid
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *InfinityUpgradeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError(
		"Update Not Supported",
		"Upgrade resources cannot be updated. To trigger a new upgrade, delete this resource and create a new one.",
	)
}

func (r *InfinityUpgradeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	state := &InfinityUpgradeResourceModel{}

	tflog.Info(ctx, "Deleting Infinity upgrade resource")

	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// For an action resource, deletion just removes the state record
	// The actual upgrade cannot be undone, but we can remove the trigger record
	tflog.Info(ctx, fmt.Sprintf("removed upgrade trigger record with ID: %s", state.ID.ValueString()))
}
