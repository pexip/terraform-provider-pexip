package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/pexip/terraform-provider-pexip/internal/helpers"
	"math"
	"strconv"
)

var (
	_ resource.ResourceWithImportState = (*InfinityWebPasswordHashResource)(nil)
)

type InfinityWebPasswordHashResource struct{}

type InfinityWebPasswordHashResourceModel struct {
	ID       types.String `tfsdk:"id"`
	Password types.String `tfsdk:"password"`
	Salt     types.String `tfsdk:"salt"`
	Rounds   types.Int32  `tfsdk:"rounds"`
	Hash     types.String `tfsdk:"hash"`
}

func (r *InfinityWebPasswordHashResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_infinity_web_password_hash"
}

func (r *InfinityWebPasswordHashResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {

}

func (r *InfinityWebPasswordHashResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Resource identifier",
			},
			"password": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(2),
					stringvalidator.LengthAtMost(255),
				},
				MarkdownDescription: "The password to hash. This should be a strong password, ideally at least 12 characters long.",
			},
			"salt": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(12, 12),
				},
				MarkdownDescription: "The hostname of the Infinity node. This should be resolvable within the Infinity cluster.",
			},
			"rounds": schema.Int32Attribute{
				Optional: true,
				Computed: true,
				Validators: []validator.Int32{
					int32validator.Between(5000, math.MaxInt32),
				},
				Default:             int32default.StaticInt32(5000),
				MarkdownDescription: "The number of rounds to use for hashing the password. This is used to increase the security of the password hash.",
			},
			"hash": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The hash of the password, generated using the provided salt and rounds.",
			},
		},
		MarkdownDescription: "Registers a node with the Infinity service. This resource is used to manage the lifecycle of nodes in the Infinity cluster.",
	}
}

func (r *InfinityWebPasswordHashResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	plan := &InfinityWebPasswordHashResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.hashPassword(ctx, plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Hashing Password",
			fmt.Sprintf("Could not hash the password: %s", err),
		)
		return
	}

	plan.ID = plan.Salt
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *InfinityWebPasswordHashResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	state := &InfinityWebPasswordHashResourceModel{}

	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

// hashPassword hashes the password using SHA-512 with a salt and rounds.
func (r *InfinityWebPasswordHashResource) hashPassword(ctx context.Context, data *InfinityWebPasswordHashResourceModel) error {
	const saltLength = 12
	const rounds = 36000

	// Generate a random salt if not provided
	if data.Salt.IsNull() || len(data.Salt.ValueString()) == 0 {
		salt, err := helpers.GenerateRandomAlphanumeric(saltLength)
		if err != nil {
			return err
		}
		data.Salt = types.StringValue(salt)
	}

	// hash the password using the provided salt
	passwordHash, err := helpers.DjangoPasswordWithSalt(data.Password.ValueString(), data.Salt.ValueString(), rounds)
	if err != nil {
		return err
	}

	data.Hash = types.StringValue(passwordHash)
	tflog.Trace(ctx, fmt.Sprintf("hashed password with salt: %s, rounds: %d", data.Salt.ValueString(), data.Rounds.ValueInt32()))
	return nil
}

func (r *InfinityWebPasswordHashResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	plan := &InfinityWebPasswordHashResourceModel{}
	state := &InfinityWebPasswordHashResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !plan.Hash.Equal(state.Hash) || !plan.Salt.Equal(state.Salt) || !plan.Rounds.Equal(state.Rounds) || !plan.Password.Equal(state.Password) {
		err := r.hashPassword(ctx, plan)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Hashing Password",
				fmt.Sprintf("Could not hash the password: %s", err),
			)
			return
		}
	}

	plan.ID = state.ID
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *InfinityWebPasswordHashResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

}

func (r *InfinityWebPasswordHashResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Validate that the ID is a valid integer
	id, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Import ID must be a valid integer, got: %s", req.ID),
		)
		return
	}

	if id <= 0 {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Import ID must be a positive integer, got: %d", id),
		)
		return
	}

	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
