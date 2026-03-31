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

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/pexip/go-infinity-sdk/v38/config"
)

var (
	_ resource.ResourceWithImportState = (*InfinityMjxGraphDeploymentResource)(nil)
)

type InfinityMjxGraphDeploymentResource struct {
	InfinityClient InfinityClient
}

type InfinityMjxGraphDeploymentResourceModel struct {
	ID              types.String `tfsdk:"id"`
	ResourceID      types.Int32  `tfsdk:"resource_id"`
	Name            types.String `tfsdk:"name"`
	Description     types.String `tfsdk:"description"`
	ClientID        types.String `tfsdk:"client_id"`
	ClientSecret    types.String `tfsdk:"client_secret"`
	OAuthTokenURL   types.String `tfsdk:"oauth_token_url"`
	GraphAPIDomain  types.String `tfsdk:"graph_api_domain"`
	RequestQuota    types.Int64  `tfsdk:"request_quota"`
	DisableProxy    types.Bool   `tfsdk:"disable_proxy"`
	MjxIntegrations types.Set    `tfsdk:"mjx_integrations"`
}

func (r *InfinityMjxGraphDeploymentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_infinity_mjx_graph_deployment"
}

func (r *InfinityMjxGraphDeploymentResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *InfinityMjxGraphDeploymentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Resource URI for the MJX Graph deployment in Infinity.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"resource_id": schema.Int32Attribute{
				Computed:            true,
				MarkdownDescription: "The resource integer identifier for the MJX Graph deployment in Infinity.",
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					stringvalidator.LengthAtMost(250),
				},
				MarkdownDescription: "The name of the OTJ O365 Graph Integration. Maximum length: 250 characters.",
			},
			"description": schema.StringAttribute{
				Computed: true,
				Optional: true,
				Default:  stringdefault.StaticString(""),
				Validators: []validator.String{
					stringvalidator.LengthAtMost(250),
				},
				MarkdownDescription: "A description of this OTJ O365 Graph Integration. Maximum length: 250 characters.",
			},
			"client_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The Application (client) ID which was generated when creating an App Registration in Azure Active Directory.",
			},
			"client_secret": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(100),
				},
				MarkdownDescription: "The client secret of the application you created in the Azure Portal, for use by OTJ. Maximum length: 100 characters.",
			},
			"oauth_token_url": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(255),
				},
				MarkdownDescription: "The URI of the OAuth 2.0 (v2) token endpoint. This should be copied from the 'Endpoints' section in Azure Active Directory App Registrations. Maximum length: 255 characters.",
			},
			"graph_api_domain": schema.StringAttribute{
				Computed: true,
				Optional: true,
				Default:  stringdefault.StaticString("graph.microsoft.com"),
				Validators: []validator.String{
					stringvalidator.LengthAtMost(192),
				},
				MarkdownDescription: "The FQDN to use when connecting to the Graph API. Maximum length: 192 characters.",
			},
			"request_quota": schema.Int64Attribute{
				Computed: true,
				Optional: true,
				Default:  int64default.StaticInt64(1000000),
				Validators: []validator.Int64{
					int64validator.Between(10000, 10000000),
				},
				MarkdownDescription: "The maximum number of API requests that can be made by OTJ to the Microsoft Graph API in a 24-hour period. Minimum: 10000. Maximum: 10000000. Default: 1000000.",
			},
			"disable_proxy": schema.BoolAttribute{
				Computed:            true,
				Optional:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Bypass the web proxy (where configured for the system location) for outbound requests sent from this integration.",
			},
			"mjx_integrations": schema.SetAttribute{
				Computed:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "The One-Touch Join Profiles associated with this OTJ O365 Graph Integration.",
			},
		},
		MarkdownDescription: "Manages an MJX Graph deployment in Infinity. An MJX Graph deployment provides integration with the Microsoft Graph API, enabling OTJ (One-Touch Join) functionality for Microsoft 365 calendar-based meeting management.",
	}
}

func (r *InfinityMjxGraphDeploymentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	plan := &InfinityMjxGraphDeploymentResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createRequest := &config.MjxGraphDeploymentCreateRequest{
		Name:           plan.Name.ValueString(),
		Description:    plan.Description.ValueString(),
		ClientID:       plan.ClientID.ValueString(),
		OAuthTokenURL:  plan.OAuthTokenURL.ValueString(),
		GraphAPIDomain: plan.GraphAPIDomain.ValueString(),
		RequestQuota:   int(plan.RequestQuota.ValueInt64()),
		DisableProxy:   plan.DisableProxy.ValueBool(),
	}

	if !plan.ClientSecret.IsNull() && !plan.ClientSecret.IsUnknown() {
		createRequest.ClientSecret = plan.ClientSecret.ValueString()
	}

	createResponse, err := r.InfinityClient.Config().CreateMjxGraphDeployment(ctx, createRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Infinity MJX Graph deployment",
			fmt.Sprintf("Could not create Infinity MJX Graph deployment: %s", err),
		)
		return
	}

	resourceID, err := createResponse.ResourceID()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Retrieving Infinity MJX Graph deployment ID",
			fmt.Sprintf("Could not retrieve ID for created Infinity MJX Graph deployment: %s", err),
		)
		return
	}

	model, err := r.read(ctx, resourceID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Created Infinity MJX Graph deployment",
			fmt.Sprintf("Could not read created Infinity MJX Graph deployment with ID %d: %s", resourceID, err),
		)
		return
	}

	// Preserve client_secret from plan as it is not returned by the API
	model.ClientSecret = plan.ClientSecret

	tflog.Trace(ctx, fmt.Sprintf("created Infinity MJX Graph deployment with ID: %s, name: %s", model.ID, model.Name))

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *InfinityMjxGraphDeploymentResource) read(ctx context.Context, resourceID int) (*InfinityMjxGraphDeploymentResourceModel, error) {
	var data InfinityMjxGraphDeploymentResourceModel

	srv, err := r.InfinityClient.Config().GetMjxGraphDeployment(ctx, resourceID)
	if err != nil {
		return nil, err
	}

	if srv.ResourceURI == "" {
		return nil, fmt.Errorf("MJX Graph deployment with ID %d not found", resourceID)
	}

	data.ID = types.StringValue(srv.ResourceURI)
	data.ResourceID = types.Int32Value(int32(resourceID)) // #nosec G115 -- API values are expected to be within int32 range
	data.Name = types.StringValue(srv.Name)
	data.Description = types.StringValue(srv.Description)
	data.ClientID = types.StringValue(srv.ClientID)
	data.OAuthTokenURL = types.StringValue(srv.OAuthTokenURL)
	data.GraphAPIDomain = types.StringValue(srv.GraphAPIDomain)
	data.RequestQuota = types.Int64Value(int64(srv.RequestQuota))
	data.DisableProxy = types.BoolValue(srv.DisableProxy)

	// Note: ClientSecret is not returned by the API and will be preserved from plan/state

	if srv.MjxIntegrations != nil && len(*srv.MjxIntegrations) > 0 {
		integrations, diags := types.SetValueFrom(ctx, types.StringType, *srv.MjxIntegrations)
		if diags.HasError() {
			return nil, fmt.Errorf("error converting MJX integrations: %v", diags)
		}
		data.MjxIntegrations = integrations
	} else {
		data.MjxIntegrations = types.SetNull(types.StringType)
	}

	return &data, nil
}

func (r *InfinityMjxGraphDeploymentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	state := &InfinityMjxGraphDeploymentResourceModel{}

	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Preserve client_secret from existing state
	clientSecret := state.ClientSecret

	resourceID := int(state.ResourceID.ValueInt32())
	state, err := r.read(ctx, resourceID)
	if err != nil {
		if isNotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading Infinity MJX Graph deployment",
			fmt.Sprintf("Could not read Infinity MJX Graph deployment: %s", err),
		)
		return
	}

	// Restore client_secret as it is not returned by the API
	state.ClientSecret = clientSecret

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *InfinityMjxGraphDeploymentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	plan := &InfinityMjxGraphDeploymentResourceModel{}
	state := &InfinityMjxGraphDeploymentResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	requestQuota := int(plan.RequestQuota.ValueInt64())
	disableProxy := plan.DisableProxy.ValueBool()

	updateRequest := &config.MjxGraphDeploymentUpdateRequest{
		Name:           plan.Name.ValueString(),
		Description:    plan.Description.ValueString(),
		ClientID:       plan.ClientID.ValueString(),
		OAuthTokenURL:  plan.OAuthTokenURL.ValueString(),
		GraphAPIDomain: plan.GraphAPIDomain.ValueString(),
		RequestQuota:   &requestQuota,
		DisableProxy:   &disableProxy,
	}

	if !plan.ClientSecret.IsNull() && !plan.ClientSecret.IsUnknown() {
		updateRequest.ClientSecret = plan.ClientSecret.ValueString()
	}

	resourceID := int(state.ResourceID.ValueInt32())
	_, err := r.InfinityClient.Config().UpdateMjxGraphDeployment(ctx, resourceID, updateRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Infinity MJX Graph deployment",
			fmt.Sprintf("Could not update Infinity MJX Graph deployment: %s", err),
		)
		return
	}

	model, err := r.read(ctx, resourceID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Updated Infinity MJX Graph deployment",
			fmt.Sprintf("Could not read updated Infinity MJX Graph deployment with ID %d: %s", resourceID, err),
		)
		return
	}

	// Preserve client_secret from plan as it is not returned by the API
	model.ClientSecret = plan.ClientSecret

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *InfinityMjxGraphDeploymentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	state := &InfinityMjxGraphDeploymentResourceModel{}

	tflog.Info(ctx, "Deleting Infinity MJX Graph deployment")

	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.InfinityClient.Config().DeleteMjxGraphDeployment(ctx, int(state.ResourceID.ValueInt32()))

	if err != nil && !isNotFoundError(err) && !isLookupError(err) {
		resp.Diagnostics.AddError(
			"Error Deleting Infinity MJX Graph deployment",
			fmt.Sprintf("Could not delete Infinity MJX Graph deployment with ID %s: %s", state.ID.ValueString(), err),
		)
		return
	}
}

func (r *InfinityMjxGraphDeploymentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resourceID, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Import ID must be a valid integer for the resource ID. Got: %s", req.ID),
		)
		return
	}

	tflog.Trace(ctx, fmt.Sprintf("Importing Infinity MJX Graph deployment with resource ID: %d", resourceID))

	model, err := r.read(ctx, resourceID)
	if err != nil {
		if isNotFoundError(err) {
			resp.Diagnostics.AddError(
				"Infinity MJX Graph Deployment Not Found",
				fmt.Sprintf("Infinity MJX Graph deployment with resource ID %d not found.", resourceID),
			)
			return
		}
		resp.Diagnostics.AddError(
			"Error Importing Infinity MJX Graph Deployment",
			fmt.Sprintf("Could not import Infinity MJX Graph deployment with resource ID %d: %s", resourceID, err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}
