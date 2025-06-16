package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/pexip/terraform-provider-pexip/internal/provider/validators"
)

var (
	_ datasource.DataSourceWithValidateConfig = (*InfinityManagerConfigDataSource)(nil)
)

type InfinityManagerConfigDataSource struct{}

func (d *InfinityManagerConfigDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_infinity_manager_config"
}

func (d *InfinityManagerConfigDataSource) ValidateConfig(ctx context.Context, req datasource.ValidateConfigRequest, resp *datasource.ValidateConfigResponse) {
	var infinityManagerConfig InfinityManagerConfigModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &infinityManagerConfig)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(infinityManagerConfig.validate(ctx)...)
}

func (d *InfinityManagerConfigDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"hostname": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(3),
				},
				MarkdownDescription: "Pexip Infinity Manager hostname, e.g. `manager-1`",
			},
			"domain": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(3),
				},
				MarkdownDescription: "Pexip Infinity Manager domain, e.g. `example.com`",
			},
			"ip": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					validators.IPAddress(),
				},
				MarkdownDescription: "Pexip Infinity Manager IP address",
			},
			"mask": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					validators.IPAddress(),
				},
				MarkdownDescription: "Pexip Infinity Manager subnet mask (e.g. 255.255.255.0)",
			},
			"gw": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					validators.IPAddress(),
				},
				MarkdownDescription: "Pexip Infinity Manager gateway IP address",
			},
			"dns": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					validators.IPAddress(),
				},
				MarkdownDescription: "Pexip Infinity Manager DNS server IP address",
			},
			"ntp": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(3),
				},
				MarkdownDescription: "Pexip Infinity Manager NTP server",
			},
			"user": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
				MarkdownDescription: "Pexip Infinity Manager username for authentication",
			},
			"pass": schema.StringAttribute{
				Required:  true,
				Sensitive: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
				MarkdownDescription: "Pexip Infinity Manager password for authentication",
			},
			"admin_password": schema.StringAttribute{
				Required:  true,
				Sensitive: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
				MarkdownDescription: "Pexip Infinity Manager admin password for authentication",
			},
			"error_reports": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Pexip Infinity Manager error reports",
			},
			"enable_analytics": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Pexip Infinity Manager enable analytics",
			},
			"contact_email_address": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(3),
				},
				MarkdownDescription: "Pexip Infinity Manager contact email address for notifications",
			},
			"rendered": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Rendered Pexip Infinity Manager bootstrap configuration.",
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "[CRC-32](https://pkg.go.dev/hash/crc32) checksum of `rendered` Pexip Infinity bootstrap config.",
			},
		},
		MarkdownDescription: "Renders Pexip Infinity Manager bootstrap configuration.",
	}
}

func (d *InfinityManagerConfigDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var model InfinityManagerConfigModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(model.update(ctx)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}
