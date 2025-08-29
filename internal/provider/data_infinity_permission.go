/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/pexip/go-infinity-sdk/v38/config"
)

type InfinityPermissionDataSource struct {
	InfinityClient InfinityClient
}

type InfinityPermissionModel struct {
	Codename    types.String `tfsdk:"codename"`
	ID          types.String `tfsdk:"id"`
	ResourceID  types.Int32  `tfsdk:"resource_id"`
	Name        types.String `tfsdk:"name"`
	ResourceURI types.String `tfsdk:"resource_uri"`
}

func (d *InfinityPermissionDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_infinity_permission"
}

func (d *InfinityPermissionDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	p, ok := req.ProviderData.(*PexipProvider)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected DataSource Configure Type",
			fmt.Sprintf("Expected *PexipProvider, got: %T. Please report this issue to the provider developers", req.ProviderData),
		)
		return
	}
	d.InfinityClient = p.client
}

func (d *InfinityPermissionDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"codename": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The codename of the permission.",
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Resource URI for the permission in Infinity",
			},
			"resource_id": schema.Int32Attribute{
				Computed:            true,
				MarkdownDescription: "The resource integer identifier for the permission in Infinity",
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the permission.",
			},
			"resource_uri": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The API resource URI for the permission.",
			},
		},
		MarkdownDescription: "Reads a permission configuration from the Infinity service.",
	}
}

func (d *InfinityPermissionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state InfinityPermissionModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Lookup by name
	list, err := d.InfinityClient.Config().ListPermissions(ctx, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Listing Infinity permissions",
			fmt.Sprintf("Could not list permissions: %s", err),
		)
		return
	}

	var perm *config.Permission
	for _, p := range list.Objects {
		if p.Name == state.Name.ValueString() {
			perm = &p
			break
		}
	}

	if perm == nil {
		resp.Diagnostics.AddError(
			"Permission Not Found",
			fmt.Sprintf("Permission with name '%s' not found.", state.Name.ValueString()),
		)
		return
	}

	state.ID = types.StringValue(fmt.Sprintf("/api/admin/configuration/v1/permission/%d/", perm.ID))
	state.ResourceID = types.Int32Value(int32(perm.ID))
	state.Name = types.StringValue(perm.Name)
	state.Codename = types.StringValue(perm.Codename)
	state.ResourceURI = types.StringValue(perm.ResourceURI)

	resp.State.Set(ctx, &state)
}
