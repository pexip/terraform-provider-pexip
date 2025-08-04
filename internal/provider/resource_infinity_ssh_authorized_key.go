package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/pexip/go-infinity-sdk/v38/config"
)

var (
	_ resource.ResourceWithImportState = (*InfinitySSHAuthorizedKeyResource)(nil)
)

type InfinitySSHAuthorizedKeyResource struct {
	InfinityClient InfinityClient
}

type InfinitySSHAuthorizedKeyResourceModel struct {
	ID         types.String `tfsdk:"id"`
	ResourceID types.Int32  `tfsdk:"resource_id"`
	Keytype    types.String `tfsdk:"keytype"`
	Key        types.String `tfsdk:"key"`
	Comment    types.String `tfsdk:"comment"`
	Nodes      types.Set    `tfsdk:"nodes"`
}

func (r *InfinitySSHAuthorizedKeyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_infinity_ssh_authorized_key"
}

func (r *InfinitySSHAuthorizedKeyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *InfinitySSHAuthorizedKeyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Resource URI for the SSH authorized key in Infinity",
			},
			"resource_id": schema.Int32Attribute{
				Computed:            true,
				MarkdownDescription: "The resource integer identifier for the SSH authorized key in Infinity",
			},
			"keytype": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf("ssh-rsa", "ssh-dss", "ssh-ed25519", "ecdsa-sha2-nistp256", "ecdsa-sha2-nistp384", "ecdsa-sha2-nistp521"),
				},
				MarkdownDescription: "The SSH key type. Valid choices: ssh-rsa, ssh-dss, ssh-ed25519, ecdsa-sha2-nistp256, ecdsa-sha2-nistp384, ecdsa-sha2-nistp521.",
			},
			"key": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
				MarkdownDescription: "The SSH public key content (base64 encoded key data).",
			},
			"comment": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(250),
				},
				MarkdownDescription: "A comment for the SSH key. Maximum length: 250 characters.",
			},
			"nodes": schema.SetAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "List of node resource URIs where this SSH key is authorized.",
			},
		},
		MarkdownDescription: "Manages an SSH authorized key configuration with the Infinity service.",
	}
}

func (r *InfinitySSHAuthorizedKeyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	plan := &InfinitySSHAuthorizedKeyResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createRequest := &config.SSHAuthorizedKeyCreateRequest{
		Keytype: plan.Keytype.ValueString(),
		Key:     plan.Key.ValueString(),
	}

	// Set optional fields
	if !plan.Comment.IsNull() {
		createRequest.Comment = plan.Comment.ValueString()
	}
	if !plan.Nodes.IsNull() && !plan.Nodes.IsUnknown() {
		var nodes []string
		resp.Diagnostics.Append(plan.Nodes.ElementsAs(ctx, &nodes, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		createRequest.Nodes = nodes
	}

	createResponse, err := r.InfinityClient.Config().CreateSSHAuthorizedKey(ctx, createRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Infinity SSH authorized key",
			fmt.Sprintf("Could not create Infinity SSH authorized key: %s", err),
		)
		return
	}

	resourceID, err := createResponse.ResourceID()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Retrieving Infinity SSH authorized key ID",
			fmt.Sprintf("Could not retrieve ID for created Infinity SSH authorized key: %s", err),
		)
		return
	}

	// Read the state from the API to get all computed values
	model, err := r.read(ctx, resourceID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Created Infinity SSH authorized key",
			fmt.Sprintf("Could not read created Infinity SSH authorized key with ID %d: %s", resourceID, err),
		)
		return
	}
	tflog.Trace(ctx, fmt.Sprintf("created Infinity SSH authorized key with ID: %s, keytype: %s", model.ID, model.Keytype))

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *InfinitySSHAuthorizedKeyResource) read(ctx context.Context, resourceID int) (*InfinitySSHAuthorizedKeyResourceModel, error) {
	var data InfinitySSHAuthorizedKeyResourceModel

	srv, err := r.InfinityClient.Config().GetSSHAuthorizedKey(ctx, resourceID)
	if err != nil {
		return nil, err
	}

	if len(srv.ResourceURI) == 0 {
		return nil, fmt.Errorf("SSH authorized key with ID %d not found", resourceID)
	}

	data.ID = types.StringValue(srv.ResourceURI)
	data.ResourceID = types.Int32Value(int32(resourceID))
	data.Keytype = types.StringValue(srv.Keytype)
	data.Key = types.StringValue(srv.Key)
	data.Comment = types.StringValue(srv.Comment)

	// Convert nodes slice to types.Set
	nodeElements := make([]types.String, len(srv.Nodes))
	for i, node := range srv.Nodes {
		nodeElements[i] = types.StringValue(node)
	}
	nodesSet, diags := types.SetValueFrom(ctx, types.StringType, nodeElements)
	if diags.HasError() {
		return nil, fmt.Errorf("failed to convert nodes to set: %v", diags)
	}
	data.Nodes = nodesSet

	return &data, nil
}

func (r *InfinitySSHAuthorizedKeyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	state := &InfinitySSHAuthorizedKeyResourceModel{}

	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resourceID := int(state.ResourceID.ValueInt32())
	state, err := r.read(ctx, resourceID)
	if err != nil {
		// Check if the error is a 404 (not found)
		if isNotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading Infinity SSH authorized key",
			fmt.Sprintf("Could not read Infinity SSH authorized key: %s", err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *InfinitySSHAuthorizedKeyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	plan := &InfinitySSHAuthorizedKeyResourceModel{}
	state := &InfinitySSHAuthorizedKeyResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resourceID := int(state.ResourceID.ValueInt32())

	updateRequest := &config.SSHAuthorizedKeyUpdateRequest{
		Keytype: plan.Keytype.ValueString(),
		Key:     plan.Key.ValueString(),
	}

	// Set optional fields
	if !plan.Comment.IsNull() {
		updateRequest.Comment = plan.Comment.ValueString()
	}
	if !plan.Nodes.IsNull() && !plan.Nodes.IsUnknown() {
		var nodes []string
		resp.Diagnostics.Append(plan.Nodes.ElementsAs(ctx, &nodes, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		updateRequest.Nodes = nodes
	}

	_, err := r.InfinityClient.Config().UpdateSSHAuthorizedKey(ctx, resourceID, updateRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Infinity SSH authorized key",
			fmt.Sprintf("Could not update Infinity SSH authorized key with ID %d: %s", resourceID, err),
		)
		return
	}

	// Re-read the resource to get the latest state
	updatedModel, err := r.read(ctx, resourceID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Updated Infinity SSH authorized key",
			fmt.Sprintf("Could not read updated Infinity SSH authorized key with ID %d: %s", resourceID, err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, updatedModel)...)
}

func (r *InfinitySSHAuthorizedKeyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	state := &InfinitySSHAuthorizedKeyResourceModel{}

	tflog.Info(ctx, "Deleting Infinity SSH authorized key")

	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.InfinityClient.Config().DeleteSSHAuthorizedKey(ctx, int(state.ResourceID.ValueInt32()))

	// Ignore 404 Not Found and Lookup errors on delete
	if err != nil && !isNotFoundError(err) && !isLookupError(err) {
		resp.Diagnostics.AddError(
			"Error Deleting Infinity SSH authorized key",
			fmt.Sprintf("Could not delete Infinity SSH authorized key with ID %s: %s", state.ID.ValueString(), err),
		)
		return
	}
}

func (r *InfinitySSHAuthorizedKeyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resourceID, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Import ID must be a valid integer for the resource ID. Got: %s", req.ID),
		)
		return
	}

	tflog.Trace(ctx, fmt.Sprintf("Importing Infinity SSH authorized key with resource ID: %d", resourceID))

	// Read the resource from the API
	model, err := r.read(ctx, resourceID)
	if err != nil {
		// Check if the error is a 404 (not found)
		if isNotFoundError(err) {
			resp.Diagnostics.AddError(
				"Infinity SSH Authorized Key Not Found",
				fmt.Sprintf("Infinity SSH authorized key with resource ID %d not found.", resourceID),
			)
			return
		}
		resp.Diagnostics.AddError(
			"Error Importing Infinity SSH Authorized Key",
			fmt.Sprintf("Could not import Infinity SSH authorized key with resource ID %d: %s", resourceID, err),
		)
		return
	}

	// Set the state from the imported resource
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}
