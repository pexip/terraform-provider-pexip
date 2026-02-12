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
)

var (
	_ resource.Resource = (*InfinityDeleteDefaultTLSCertificateActionResource)(nil)
)

type InfinityDeleteDefaultTLSCertificateActionResource struct {
	InfinityClient InfinityClient
}

type InfinityDeleteDefaultTLSCertificateActionResourceModel struct {
	ID        types.String `tfsdk:"id"`
	Timestamp types.String `tfsdk:"timestamp"`
}

func (r *InfinityDeleteDefaultTLSCertificateActionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_infinity_delete_default_tls_certificate_action"
}

func (r *InfinityDeleteDefaultTLSCertificateActionResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *InfinityDeleteDefaultTLSCertificateActionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Unique identifier for this deletion action",
			},
			"timestamp": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Timestamp when the default TLS certificate was deleted",
			},
		},
		MarkdownDescription: "Deletes the default TLS certificate (ID 1) that is automatically created when a management node is deployed. This is an action resource that performs a one-time deletion. Note: This resource only supports creation and reading - the deletion cannot be updated or undone once performed. The resource represents the deletion action trigger, not the certificate state itself.",
	}
}

func (r *InfinityDeleteDefaultTLSCertificateActionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	plan := &InfinityDeleteDefaultTLSCertificateActionResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete the default TLS certificate (ID 1)
	const defaultCertificateID = 1
	err := r.InfinityClient.Config().DeleteTLSCertificate(ctx, defaultCertificateID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Default TLS Certificate",
			fmt.Sprintf("Could not delete default TLS certificate (ID %d): %s", defaultCertificateID, err),
		)
		return
	}

	// Generate a unique ID for this deletion action
	timestamp := time.Now().UTC()
	actionID := fmt.Sprintf("delete_default_tls_cert_%d", timestamp.Unix())

	// Since this is an action resource, we create a simple state representation
	model := &InfinityDeleteDefaultTLSCertificateActionResourceModel{
		ID:        types.StringValue(actionID),
		Timestamp: types.StringValue(timestamp.Format(time.RFC3339)),
	}

	// Log the deletion action
	tflog.Info(ctx, fmt.Sprintf("deleted default TLS certificate (ID %d) with action ID: %s", defaultCertificateID, actionID))

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *InfinityDeleteDefaultTLSCertificateActionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	state := &InfinityDeleteDefaultTLSCertificateActionResourceModel{}

	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// For an action resource like this deletion, we just maintain the state as-is
	// The certificate deletion is permanent, but the action record remains valid
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *InfinityDeleteDefaultTLSCertificateActionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError(
		"Update Not Supported",
		"Delete default TLS certificate action resources cannot be updated. To trigger a new deletion (if a new default certificate was created), delete this resource and create a new one.",
	)
}

func (r *InfinityDeleteDefaultTLSCertificateActionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	state := &InfinityDeleteDefaultTLSCertificateActionResourceModel{}

	tflog.Info(ctx, "Deleting default TLS certificate deletion action resource")

	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// For an action resource, deletion just removes the state record
	// The actual certificate deletion cannot be undone, but we can remove the action record
	tflog.Info(ctx, fmt.Sprintf("removed delete default TLS certificate action record with ID: %s", state.ID.ValueString()))
}
