package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceInfinityManagerConfig() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceInfinityManagerConfigRead,
		Schema: map[string]*schema.Schema{
			"head": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceInfinityManagerConfigRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conf := meta.(providerConfiguration)

	repo := repository.Factory(conf.RepositoryType, conf.PrivateKey, conf.PrivateKeyPassword, conf.UseDefaultSSH)
	_, err := repo.Clone(conf.RepositoryURL, conf.RepositoryBranch)
	if err != nil {
		return diag.FromErr(err)
	}

	err = updateRepositoryResource(d, repo)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func updateRepositoryResource(d *schema.ResourceData, repo repository.Repository) error {
	head, err := repo.Head()
	if err != nil {
		return err
	}
	_ = d.Set("head", head)

	// Always set this
	d.SetId(repo.ID())

	return nil
}
