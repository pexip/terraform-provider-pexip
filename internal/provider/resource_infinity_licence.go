package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/pexip/go-infinity-sdk/v38/config"
)

var (
	_ resource.ResourceWithImportState = (*InfinityLicenceResource)(nil)
)

type InfinityLicenceResource struct {
	InfinityClient InfinityClient
}

type InfinityLicenceResourceModel struct {
	ID                   types.String `tfsdk:"id"`
	FulfillmentID        types.String `tfsdk:"fulfillment_id"`
	EntitlementID        types.String `tfsdk:"entitlement_id"`
	FulfillmentType      types.String `tfsdk:"fulfillment_type"`
	ProductID            types.String `tfsdk:"product_id"`
	LicenseType          types.String `tfsdk:"license_type"`
	Features             types.String `tfsdk:"features"`
	Concurrent           types.Int64  `tfsdk:"concurrent"`
	ConcurrentOverdraft  types.Int64  `tfsdk:"concurrent_overdraft"`
	Activatable          types.Int64  `tfsdk:"activatable"`
	ActivatableOverdraft types.Int64  `tfsdk:"activatable_overdraft"`
	Hybrid               types.Int64  `tfsdk:"hybrid"`
	HybridOverdraft      types.Int64  `tfsdk:"hybrid_overdraft"`
	StartDate            types.String `tfsdk:"start_date"`
	ExpirationDate       types.String `tfsdk:"expiration_date"`
	Status               types.String `tfsdk:"status"`
	TrustFlags           types.Int64  `tfsdk:"trust_flags"`
	Repair               types.Int64  `tfsdk:"repair"`
	ServerChain          types.String `tfsdk:"server_chain"`
	OfflineMode          types.Bool   `tfsdk:"offline_mode"`
}

func (r *InfinityLicenceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_infinity_licence"
}

func (r *InfinityLicenceResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *InfinityLicenceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Resource URI for the licence in Infinity",
			},
			"fulfillment_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The fulfillment ID of the licence (used as the unique identifier)",
			},
			"entitlement_id": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
				MarkdownDescription: "The entitlement ID for the licence activation.",
			},
			"fulfillment_type": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The fulfillment type of the licence.",
			},
			"product_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The product ID associated with the licence.",
			},
			"license_type": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The type of the licence.",
			},
			"features": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The features enabled by this licence.",
			},
			"concurrent": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Number of concurrent sessions allowed.",
			},
			"concurrent_overdraft": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Number of concurrent sessions allowed in overdraft.",
			},
			"activatable": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Number of activatable licenses.",
			},
			"activatable_overdraft": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Number of activatable licenses allowed in overdraft.",
			},
			"hybrid": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Number of hybrid licenses.",
			},
			"hybrid_overdraft": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Number of hybrid licenses allowed in overdraft.",
			},
			"start_date": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The start date of the licence validity period.",
			},
			"expiration_date": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The expiration date of the licence.",
			},
			"status": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The current status of the licence.",
			},
			"trust_flags": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Trust flags for the licence.",
			},
			"repair": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Repair flag for the licence.",
			},
			"server_chain": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The server chain for the licence.",
			},
			"offline_mode": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Whether the licence should be activated in offline mode.",
			},
		},
		MarkdownDescription: "Manages a licence configuration with the Infinity service. This resource activates licences using entitlement IDs.",
	}
}

func (r *InfinityLicenceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	plan := &InfinityLicenceResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createRequest := &config.LicenceCreateRequest{
		EntitlementID: plan.EntitlementID.ValueString(),
		OfflineMode:   plan.OfflineMode.ValueBool(),
	}

	_, err := r.InfinityClient.Config().CreateLicence(ctx, createRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Infinity licence",
			fmt.Sprintf("Could not create Infinity licence: %s", err),
		)
		return
	}

	listResponse, err := r.InfinityClient.Config().ListLicences(ctx, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Listing Infinity licences",
			fmt.Sprintf("Could not list Infinity licences after creation: %s", err),
		)
		return
	}

	fulfillmentID := ""
	for _, licence := range listResponse.Objects {
		if licence.EntitlementID == strings.Replace(plan.EntitlementID.ValueString(), " ", "", -1) {
			fulfillmentID = licence.FulfillmentID
			break
		}
	}
	if fulfillmentID == "" {
		resp.Diagnostics.AddError(
			"Fulfillment ID Not Found",
			fmt.Sprintf("Could not find fulfillment ID for entitlement ID %s after creation", plan.EntitlementID.ValueString()),
		)
		return
	}

	// Read the state from the API to get all computed values
	model, err := r.read(ctx, fulfillmentID, plan.EntitlementID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Created Infinity licence",
			fmt.Sprintf("Could not read created Infinity licence with fulfillment ID %s: %s", fulfillmentID, err),
		)
		return
	}
	tflog.Trace(ctx, fmt.Sprintf("created Infinity licence with ID: %s, entitlement: %s, fulfillment: %s", model.ID, model.EntitlementID, model.FulfillmentID))

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *InfinityLicenceResource) read(ctx context.Context, fulfillmentID string, entitlementID string) (*InfinityLicenceResourceModel, error) {
	var data InfinityLicenceResourceModel

	srv, err := r.InfinityClient.Config().GetLicence(ctx, fulfillmentID)
	if err != nil {
		return nil, err
	}

	if len(srv.ResourceURI) == 0 {
		return nil, fmt.Errorf("licence with fulfillment ID %s not found", fulfillmentID)
	}

	if len(entitlementID) > 0 {
		data.EntitlementID = types.StringValue(entitlementID)
	} else {
		data.EntitlementID = types.StringValue(srv.EntitlementID)
	}

	data.ID = types.StringValue(srv.ResourceURI)
	data.FulfillmentID = types.StringValue(srv.FulfillmentID)
	data.FulfillmentType = types.StringValue(srv.FulfillmentType)
	data.ProductID = types.StringValue(srv.ProductID)
	data.LicenseType = types.StringValue(srv.LicenseType)
	data.Features = types.StringValue(srv.Features)
	data.Concurrent = types.Int64Value(int64(srv.Concurrent))
	data.ConcurrentOverdraft = types.Int64Value(int64(srv.ConcurrentOverdraft))
	data.Activatable = types.Int64Value(int64(srv.Activatable))
	data.ActivatableOverdraft = types.Int64Value(int64(srv.ActivatableOverdraft))
	data.Hybrid = types.Int64Value(int64(srv.Hybrid))
	data.HybridOverdraft = types.Int64Value(int64(srv.HybridOverdraft))
	data.StartDate = types.StringValue(srv.StartDate)
	data.ExpirationDate = types.StringValue(srv.ExpirationDate)
	data.Status = types.StringValue(srv.Status)
	data.TrustFlags = types.Int64Value(int64(srv.TrustFlags))
	data.Repair = types.Int64Value(int64(srv.Repair))
	data.ServerChain = types.StringValue(srv.ServerChain)
	data.OfflineMode = types.BoolValue(srv.OfflineMode)

	return &data, nil
}

func (r *InfinityLicenceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	state := &InfinityLicenceResourceModel{}

	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	fulfillmentID := state.FulfillmentID.ValueString()
	state, err := r.read(ctx, fulfillmentID, state.EntitlementID.ValueString())
	if err != nil {
		// Check if the error is a 404 (not found)
		if isNotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading Infinity licence",
			fmt.Sprintf("Could not read Infinity licence: %s", err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *InfinityLicenceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Licences cannot be updated - they are immutable once activated
	// Only deactivation (delete) and re-activation (create) are supported
	resp.Diagnostics.AddError(
		"Update Not Supported",
		"Licence resources cannot be updated. To change licence settings, you must delete and recreate the resource.",
	)
}

func (r *InfinityLicenceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	state := &InfinityLicenceResourceModel{}

	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	fulfillmentID := state.FulfillmentID.ValueString()
	tflog.Info(ctx, "Deleting Infinity licence", map[string]interface{}{"entitlement_id": state.EntitlementID.ValueString(), "fulfillment_id": fulfillmentID})
	err := r.InfinityClient.Config().DeleteLicence(ctx, fulfillmentID)

	// Ignore 404 Not Found and Lookup errors on delete
	if err != nil && !isNotFoundError(err) && !isLookupError(err) {
		resp.Diagnostics.AddError(
			"Error Deleting Infinity licence",
			fmt.Sprintf("Could not delete Infinity licence with fulfillment ID %s: %s", fulfillmentID, err),
		)
		return
	}
}

func (r *InfinityLicenceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	fulfillmentID := req.ID

	tflog.Trace(ctx, fmt.Sprintf("Importing Infinity licence with fulfillment ID: %s", fulfillmentID))

	// Read the resource from the API
	model, err := r.read(ctx, fulfillmentID, "")
	if err != nil {
		// Check if the error is a 404 (not found)
		if isNotFoundError(err) {
			resp.Diagnostics.AddError(
				"Infinity Licence Not Found",
				fmt.Sprintf("Infinity licence with fulfillment ID %s not found.", fulfillmentID),
			)
			return
		}
		resp.Diagnostics.AddError(
			"Error Importing Infinity Licence",
			fmt.Sprintf("Could not import Infinity licence with fulfillment ID %s: %s", fulfillmentID, err),
		)
		return
	}

	// Set the state from the imported resource
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}
