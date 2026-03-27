/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/pexip/go-infinity-sdk/v38/config"
)

var (
	_ resource.ResourceWithImportState = (*InfinityMjxExchangeAutodiscoverURLResource)(nil)
)

type InfinityMjxExchangeAutodiscoverURLResource struct {
	InfinityClient InfinityClient
}

type InfinityMjxExchangeAutodiscoverURLResourceModel struct {
	ID                 types.String `tfsdk:"id"`
	ResourceID         types.Int32  `tfsdk:"resource_id"`
	Name               types.String `tfsdk:"name"`
	Description        types.String `tfsdk:"description"`
	URL                types.String `tfsdk:"url"`
	ExchangeDeployment types.String `tfsdk:"exchange_deployment"`
}

func (r *InfinityMjxExchangeAutodiscoverURLResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_infinity_mjx_exchange_autodiscover_url"
}

func (r *InfinityMjxExchangeAutodiscoverURLResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *InfinityMjxExchangeAutodiscoverURLResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Resource URI for the MJX Exchange Autodiscover URL in Infinity.",
			},
			"resource_id": schema.Int32Attribute{
				Computed:            true,
				MarkdownDescription: "The resource integer identifier for the MJX Exchange Autodiscover URL in Infinity.",
			},
			"name": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					stringvalidator.LengthAtMost(250),
				},
				MarkdownDescription: "The name of this Exchange Autodiscover URL. Maximum length: 250 characters.",
			},
			"description": schema.StringAttribute{
				Computed: true,
				Optional: true,
				Default:  stringdefault.StaticString(""),
				Validators: []validator.String{
					stringvalidator.LengthAtMost(250),
				},
				MarkdownDescription: "An optional description of this Exchange Autodiscover URL. Maximum length: 250 characters.",
			},
			"url": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(255),
				},
				MarkdownDescription: "The URL used to connect to the Autodiscover service on the Exchange deployment. Maximum length: 255 characters.",
			},
			"exchange_deployment": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The OTJ Exchange Integration this Autodiscover URL belongs to.",
			},
		},
		MarkdownDescription: "Manages an MJX Exchange Autodiscover URL in Infinity. An Autodiscover URL is used by the OTJ Exchange Integration to discover Exchange Web Services endpoints.",
	}
}

func (r *InfinityMjxExchangeAutodiscoverURLResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	plan := &InfinityMjxExchangeAutodiscoverURLResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	exchangeDeployment := plan.ExchangeDeployment.ValueString()
	createRequest := &config.MjxExchangeAutodiscoverURLCreateRequest{
		Name:               plan.Name.ValueString(),
		Description:        plan.Description.ValueString(),
		URL:                plan.URL.ValueString(),
		ExchangeDeployment: &exchangeDeployment,
	}

	createResponse, err := r.InfinityClient.Config().CreateMjxExchangeAutodiscoverURL(ctx, createRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Infinity MJX Exchange Autodiscover URL",
			fmt.Sprintf("Could not create Infinity MJX Exchange Autodiscover URL: %s", err),
		)
		return
	}

	resourceID, err := createResponse.ResourceID()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Retrieving Infinity MJX Exchange Autodiscover URL ID",
			fmt.Sprintf("Could not retrieve ID for created Infinity MJX Exchange Autodiscover URL: %s", err),
		)
		return
	}

	model, err := r.read(ctx, resourceID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Created Infinity MJX Exchange Autodiscover URL",
			fmt.Sprintf("Could not read created Infinity MJX Exchange Autodiscover URL with ID %d: %s", resourceID, err),
		)
		return
	}

	tflog.Trace(ctx, fmt.Sprintf("created Infinity MJX Exchange Autodiscover URL with ID: %s, name: %s", model.ID, model.Name))

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *InfinityMjxExchangeAutodiscoverURLResource) read(ctx context.Context, resourceID int) (*InfinityMjxExchangeAutodiscoverURLResourceModel, error) {
	var data InfinityMjxExchangeAutodiscoverURLResourceModel

	srv, err := r.InfinityClient.Config().GetMjxExchangeAutodiscoverURL(ctx, resourceID)
	if err != nil {
		return nil, err
	}

	if srv.ResourceURI == "" {
		return nil, fmt.Errorf("MJX Exchange Autodiscover URL with ID %d not found", resourceID)
	}

	data.ID = types.StringValue(srv.ResourceURI)
	data.ResourceID = types.Int32Value(int32(resourceID)) // #nosec G115 -- API values are expected to be within int32 range
	data.Name = types.StringValue(srv.Name)
	data.Description = types.StringValue(srv.Description)
	data.URL = types.StringValue(srv.URL)

	if srv.ExchangeDeployment != nil {
		data.ExchangeDeployment = types.StringValue(*srv.ExchangeDeployment)
	} else {
		data.ExchangeDeployment = types.StringNull()
	}

	return &data, nil
}

func (r *InfinityMjxExchangeAutodiscoverURLResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	state := &InfinityMjxExchangeAutodiscoverURLResourceModel{}

	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resourceID := int(state.ResourceID.ValueInt32())
	model, err := r.read(ctx, resourceID)
	if err != nil {
		if isNotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading Infinity MJX Exchange Autodiscover URL",
			fmt.Sprintf("Could not read Infinity MJX Exchange Autodiscover URL: %s", err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *InfinityMjxExchangeAutodiscoverURLResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	plan := &InfinityMjxExchangeAutodiscoverURLResourceModel{}
	state := &InfinityMjxExchangeAutodiscoverURLResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	exchangeDeployment := plan.ExchangeDeployment.ValueString()
	updateRequest := &config.MjxExchangeAutodiscoverURLUpdateRequest{
		Name:               plan.Name.ValueString(),
		Description:        plan.Description.ValueString(),
		URL:                plan.URL.ValueString(),
		ExchangeDeployment: &exchangeDeployment,
	}

	resourceID := int(state.ResourceID.ValueInt32())
	_, err := r.InfinityClient.Config().UpdateMjxExchangeAutodiscoverURL(ctx, resourceID, updateRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Infinity MJX Exchange Autodiscover URL",
			fmt.Sprintf("Could not update Infinity MJX Exchange Autodiscover URL: %s", err),
		)
		return
	}

	model, err := r.read(ctx, resourceID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Updated Infinity MJX Exchange Autodiscover URL",
			fmt.Sprintf("Could not read updated Infinity MJX Exchange Autodiscover URL with ID %d: %s", resourceID, err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *InfinityMjxExchangeAutodiscoverURLResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	state := &InfinityMjxExchangeAutodiscoverURLResourceModel{}

	tflog.Info(ctx, "Deleting Infinity MJX Exchange Autodiscover URL")

	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.InfinityClient.Config().DeleteMjxExchangeAutodiscoverURL(ctx, int(state.ResourceID.ValueInt32()))

	if err != nil && !isNotFoundError(err) && !isLookupError(err) {
		resp.Diagnostics.AddError(
			"Error Deleting Infinity MJX Exchange Autodiscover URL",
			fmt.Sprintf("Could not delete Infinity MJX Exchange Autodiscover URL with ID %s: %s", state.ID.ValueString(), err),
		)
		return
	}
}

func (r *InfinityMjxExchangeAutodiscoverURLResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resourceID, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Import ID must be a valid integer for the resource ID. Got: %s", req.ID),
		)
		return
	}

	tflog.Trace(ctx, fmt.Sprintf("Importing Infinity MJX Exchange Autodiscover URL with resource ID: %d", resourceID))

	model, err := r.read(ctx, resourceID)
	if err != nil {
		if isNotFoundError(err) {
			resp.Diagnostics.AddError(
				"Infinity MJX Exchange Autodiscover URL Not Found",
				fmt.Sprintf("Infinity MJX Exchange Autodiscover URL with resource ID %d not found.", resourceID),
			)
			return
		}
		resp.Diagnostics.AddError(
			"Error Importing Infinity MJX Exchange Autodiscover URL",
			fmt.Sprintf("Could not import Infinity MJX Exchange Autodiscover URL with resource ID %d: %s", resourceID, err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}
