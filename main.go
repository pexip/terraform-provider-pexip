package main

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/pexip/terraform-provider-pexip/internal/log"
	"github.com/pexip/terraform-provider-pexip/internal/provider"
	"os"
)

var (
	version   string // build version number
	commit    string // sha1 revision used to build the program
	buildTime string // when the executable was built
	buildBy   string
)

func getVersionString(name string) string {
	return fmt.Sprintf("%s %s (%s at %s by %s)", name, version, commit, buildTime, buildBy)
}

func main() {
	ctx := context.Background()
	logger := log.NewTerraformLogger()
	path, err := os.Getwd()
	if err != nil {
		logger.Errorf(ctx, "failed to initialize provider: %s", err.Error())
	}

	logger.Infof(ctx, "%s", getVersionString("terraform-provider-pexip"))
	logger.Infof(ctx, "%s", path)
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() *schema.Provider {
			return provider.New()
		},
	})
}
