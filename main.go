package main

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
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

	err = providerserver.Serve(ctx, provider.New, providerserver.ServeOpts{
		Address: "registry.terraform.io/pexip/pexip",
	})
	if err != nil {
		logger.Errorf(ctx, "failed to serve provider: %s", err.Error())
		os.Exit(1)
	}
}
