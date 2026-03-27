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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/pexip/go-infinity-sdk/v38/config"
)

var (
	_ resource.ResourceWithImportState = (*InfinityMjxMeetingProcessingRuleResource)(nil)
)

type InfinityMjxMeetingProcessingRuleResource struct {
	InfinityClient InfinityClient
}

type InfinityMjxMeetingProcessingRuleResourceModel struct {
	ID                       types.String `tfsdk:"id"`
	ResourceID               types.Int32  `tfsdk:"resource_id"`
	Name                     types.String `tfsdk:"name"`
	Description              types.String `tfsdk:"description"`
	Priority                 types.Int64  `tfsdk:"priority"`
	Enabled                  types.Bool   `tfsdk:"enabled"`
	MeetingType              types.String `tfsdk:"meeting_type"`
	MjxIntegration           types.String `tfsdk:"mjx_integration"`
	MatchString              types.String `tfsdk:"match_string"`
	ReplaceString            types.String `tfsdk:"replace_string"`
	TransformRule            types.String `tfsdk:"transform_rule"`
	CustomTemplate           types.String `tfsdk:"custom_template"`
	Domain                   types.String `tfsdk:"domain"`
	CompanyID                types.String `tfsdk:"company_id"`
	IncludePin               types.Bool   `tfsdk:"include_pin"`
	DefaultProcessingEnabled types.Bool   `tfsdk:"default_processing_enabled"`
}

func (r *InfinityMjxMeetingProcessingRuleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_infinity_mjx_meeting_processing_rule"
}

func (r *InfinityMjxMeetingProcessingRuleResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *InfinityMjxMeetingProcessingRuleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Resource URI for the MJX meeting processing rule in Infinity.",
			},
			"resource_id": schema.Int32Attribute{
				Computed:            true,
				MarkdownDescription: "The resource integer identifier for the MJX meeting processing rule in Infinity.",
			},
			"name": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					stringvalidator.LengthAtMost(250),
				},
				MarkdownDescription: "The name of this Meeting Processing Rule. Maximum length: 250 characters.",
			},
			"description": schema.StringAttribute{
				Computed: true,
				Optional: true,
				Default:  stringdefault.StaticString(""),
				Validators: []validator.String{
					stringvalidator.LengthAtMost(250),
				},
				MarkdownDescription: "The description of this Meeting Processing Rule. Maximum length: 250 characters.",
			},
			"priority": schema.Int64Attribute{
				Required: true,
				Validators: []validator.Int64{
					int64validator.Between(1, 200),
				},
				MarkdownDescription: "The priority of this rule. Rules are checked in ascending priority order until the first matching rule is found, and it is then applied. Range: 1 to 200.",
			},
			"enabled": schema.BoolAttribute{
				Computed:            true,
				Optional:            true,
				Default:             booldefault.StaticBool(true),
				MarkdownDescription: "Determines whether or not the rule is enabled. Any disabled rules still appear in the rules list but are ignored. Use this setting to test configuration changes, or to temporarily disable specific rules.",
			},
			"meeting_type": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						"pexipinfinity", "pexipservice", "teams", "teamssipguestjoin",
						"polyteamsbody", "ciscoteamsbody", "pexipserviceteamsbody", "pexipinfinityteamsbody",
						"hangouts", "googlemeetsipguestjoin", "s4b", "polys4bbody",
						"webex", "zoom", "gotomeeting", "domain", "regex", "custom",
					),
				},
				MarkdownDescription: "The meeting type of this Meeting Processing Rule. Valid values: `pexipinfinity`, `pexipservice`, `teams`, `teamssipguestjoin`, `polyteamsbody`, `ciscoteamsbody`, `pexipserviceteamsbody`, `pexipinfinityteamsbody`, `hangouts`, `googlemeetsipguestjoin`, `s4b`, `polys4bbody`, `webex`, `zoom`, `gotomeeting`, `domain`, `regex`, `custom`.",
			},
			"mjx_integration": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The One-Touch Join Profile associated with this Meeting Processing Rule.",
			},
			"match_string": schema.StringAttribute{
				Computed: true,
				Optional: true,
				Default:  stringdefault.StaticString(""),
				Validators: []validator.String{
					stringvalidator.LengthAtMost(250),
				},
				MarkdownDescription: "The regular expression that defines the string to search for in the invitation. Maximum length: 250 characters.",
			},
			"replace_string": schema.StringAttribute{
				Computed: true,
				Optional: true,
				Default:  stringdefault.StaticString(""),
				Validators: []validator.String{
					stringvalidator.LengthAtMost(250),
				},
				MarkdownDescription: "A regular expression that defines how to transform the matched string into the alias to dial. Maximum length: 250 characters.",
			},
			"transform_rule": schema.StringAttribute{
				Computed: true,
				Optional: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(512),
				},
				MarkdownDescription: "A Jinja2 template that is used to process the meeting information in order to extract the meeting alias. Maximum length: 512 characters.",
			},
			"custom_template": schema.StringAttribute{
				Computed: true,
				Optional: true,
				Default:  stringdefault.StaticString(""),
				Validators: []validator.String{
					stringvalidator.LengthAtMost(10240),
				},
				MarkdownDescription: "A Jinja2 template which is used to process the meeting information from calendar events in order to extract the meeting alias. Maximum length: 10240 characters.",
			},
			"domain": schema.StringAttribute{
				Computed: true,
				Optional: true,
				Default:  stringdefault.StaticString(""),
				Validators: []validator.String{
					stringvalidator.LengthAtMost(255),
				},
				MarkdownDescription: "The domain associated with this meeting invitation. Maximum length: 255 characters.",
			},
			"company_id": schema.StringAttribute{
				Computed: true,
				Optional: true,
				Default:  stringdefault.StaticString(""),
				Validators: []validator.String{
					stringvalidator.LengthAtMost(36),
				},
				MarkdownDescription: "For a Meeting type of Teams SIP Guest Join or Google Meet SIP Guest Join: the Pexip Service Customer ID that OTJ will add to the dial string for CDRs to appear in Pexip Control Center. This field is required unless your SIP endpoints are registered to the Pexip cloud service. Maximum length: 36 characters.",
			},
			"include_pin": schema.BoolAttribute{
				Computed:            true,
				Optional:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Append the meeting password to the alias, so that users do not have to enter the password themselves.",
			},
			"default_processing_enabled": schema.BoolAttribute{
				Computed:            true,
				Optional:            true,
				Default:             booldefault.StaticBool(true),
				MarkdownDescription: "Apply default meeting processing rules for this meeting type.",
			},
		},
		MarkdownDescription: "Manages an MJX meeting processing rule in Infinity. Meeting processing rules define how OTJ (One-Touch Join) extracts the meeting alias from calendar invitations based on meeting type and pattern matching.",
	}
}

func (r *InfinityMjxMeetingProcessingRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	plan := &InfinityMjxMeetingProcessingRuleResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createRequest := &config.MjxMeetingProcessingRuleCreateRequest{
		Name:                     plan.Name.ValueString(),
		Description:              plan.Description.ValueString(),
		Priority:                 int(plan.Priority.ValueInt64()),
		Enabled:                  plan.Enabled.ValueBool(),
		MeetingType:              plan.MeetingType.ValueString(),
		MjxIntegration:           plan.MjxIntegration.ValueString(),
		MatchString:              plan.MatchString.ValueString(),
		ReplaceString:            plan.ReplaceString.ValueString(),
		TransformRule:            plan.TransformRule.ValueString(),
		CustomTemplate:           plan.CustomTemplate.ValueString(),
		Domain:                   plan.Domain.ValueString(),
		CompanyID:                plan.CompanyID.ValueString(),
		IncludePin:               plan.IncludePin.ValueBool(),
		DefaultProcessingEnabled: plan.DefaultProcessingEnabled.ValueBool(),
	}

	createResponse, err := r.InfinityClient.Config().CreateMjxMeetingProcessingRule(ctx, createRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Infinity MJX meeting processing rule",
			fmt.Sprintf("Could not create Infinity MJX meeting processing rule: %s", err),
		)
		return
	}

	resourceID, err := createResponse.ResourceID()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Retrieving Infinity MJX meeting processing rule ID",
			fmt.Sprintf("Could not retrieve ID for created Infinity MJX meeting processing rule: %s", err),
		)
		return
	}

	model, err := r.read(ctx, resourceID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Created Infinity MJX meeting processing rule",
			fmt.Sprintf("Could not read created Infinity MJX meeting processing rule with ID %d: %s", resourceID, err),
		)
		return
	}

	tflog.Trace(ctx, fmt.Sprintf("created Infinity MJX meeting processing rule with ID: %s, name: %s", model.ID, model.Name))

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *InfinityMjxMeetingProcessingRuleResource) read(ctx context.Context, resourceID int) (*InfinityMjxMeetingProcessingRuleResourceModel, error) {
	var data InfinityMjxMeetingProcessingRuleResourceModel

	srv, err := r.InfinityClient.Config().GetMjxMeetingProcessingRule(ctx, resourceID)
	if err != nil {
		return nil, err
	}

	if srv.ResourceURI == "" {
		return nil, fmt.Errorf("MJX meeting processing rule with ID %d not found", resourceID)
	}

	data.ID = types.StringValue(srv.ResourceURI)
	data.ResourceID = types.Int32Value(int32(resourceID)) // #nosec G115 -- API values are expected to be within int32 range
	data.Name = types.StringValue(srv.Name)
	data.Description = types.StringValue(srv.Description)
	data.Priority = types.Int64Value(int64(srv.Priority))
	data.Enabled = types.BoolValue(srv.Enabled)
	data.MeetingType = types.StringValue(srv.MeetingType)
	data.MjxIntegration = types.StringValue(srv.MjxIntegration)
	data.MatchString = types.StringValue(srv.MatchString)
	data.ReplaceString = types.StringValue(srv.ReplaceString)
	data.TransformRule = types.StringValue(srv.TransformRule)
	data.CustomTemplate = types.StringValue(srv.CustomTemplate)
	data.Domain = types.StringValue(srv.Domain)
	data.CompanyID = types.StringValue(srv.CompanyID)
	data.IncludePin = types.BoolValue(srv.IncludePin)
	data.DefaultProcessingEnabled = types.BoolValue(srv.DefaultProcessingEnabled)

	return &data, nil
}

func (r *InfinityMjxMeetingProcessingRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	state := &InfinityMjxMeetingProcessingRuleResourceModel{}

	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resourceID := int(state.ResourceID.ValueInt32())
	state, err := r.read(ctx, resourceID)
	if err != nil {
		if isNotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading Infinity MJX meeting processing rule",
			fmt.Sprintf("Could not read Infinity MJX meeting processing rule: %s", err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *InfinityMjxMeetingProcessingRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	plan := &InfinityMjxMeetingProcessingRuleResourceModel{}
	state := &InfinityMjxMeetingProcessingRuleResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateRequest := &config.MjxMeetingProcessingRuleUpdateRequest{
		Name:                     plan.Name.ValueString(),
		Description:              plan.Description.ValueString(),
		Priority:                 int(plan.Priority.ValueInt64()),
		Enabled:                  plan.Enabled.ValueBool(),
		MeetingType:              plan.MeetingType.ValueString(),
		MjxIntegration:           plan.MjxIntegration.ValueString(),
		MatchString:              plan.MatchString.ValueString(),
		ReplaceString:            plan.ReplaceString.ValueString(),
		TransformRule:            plan.TransformRule.ValueString(),
		CustomTemplate:           plan.CustomTemplate.ValueString(),
		Domain:                   plan.Domain.ValueString(),
		CompanyID:                plan.CompanyID.ValueString(),
		IncludePin:               plan.IncludePin.ValueBool(),
		DefaultProcessingEnabled: plan.DefaultProcessingEnabled.ValueBool(),
	}

	resourceID := int(state.ResourceID.ValueInt32())
	_, err := r.InfinityClient.Config().UpdateMjxMeetingProcessingRule(ctx, resourceID, updateRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Infinity MJX meeting processing rule",
			fmt.Sprintf("Could not update Infinity MJX meeting processing rule: %s", err),
		)
		return
	}

	model, err := r.read(ctx, resourceID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Updated Infinity MJX meeting processing rule",
			fmt.Sprintf("Could not read updated Infinity MJX meeting processing rule with ID %d: %s", resourceID, err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *InfinityMjxMeetingProcessingRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	state := &InfinityMjxMeetingProcessingRuleResourceModel{}

	tflog.Info(ctx, "Deleting Infinity MJX meeting processing rule")

	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.InfinityClient.Config().DeleteMjxMeetingProcessingRule(ctx, int(state.ResourceID.ValueInt32()))

	if err != nil && !isNotFoundError(err) && !isLookupError(err) {
		resp.Diagnostics.AddError(
			"Error Deleting Infinity MJX meeting processing rule",
			fmt.Sprintf("Could not delete Infinity MJX meeting processing rule with ID %s: %s", state.ID.ValueString(), err),
		)
		return
	}
}

func (r *InfinityMjxMeetingProcessingRuleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resourceID, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Import ID must be a valid integer for the resource ID. Got: %s", req.ID),
		)
		return
	}

	tflog.Trace(ctx, fmt.Sprintf("Importing Infinity MJX meeting processing rule with resource ID: %d", resourceID))

	model, err := r.read(ctx, resourceID)
	if err != nil {
		if isNotFoundError(err) {
			resp.Diagnostics.AddError(
				"Infinity MJX Meeting Processing Rule Not Found",
				fmt.Sprintf("Infinity MJX meeting processing rule with resource ID %d not found.", resourceID),
			)
			return
		}
		resp.Diagnostics.AddError(
			"Error Importing Infinity MJX Meeting Processing Rule",
			fmt.Sprintf("Could not import Infinity MJX meeting processing rule with resource ID %d: %s", resourceID, err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}
