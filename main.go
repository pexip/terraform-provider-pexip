package main

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/pexip/terraform-provider-pexip/internal/log"
	"github.com/pexip/terraform-provider-pexip/internal/provider"
	"github.com/pexip/terraform-provider-pexip/internal/version"
	"os"
)

func main() {
	ctx := context.Background()
	logger := log.NewTerraformLogger()
	path, err := os.Getwd()
	if err != nil {
		logger.Errorf(ctx, "failed to initialize provider: %s", err.Error())
	}

	logger.Infof(ctx, "%s", version.Version().String())
	logger.Infof(ctx, "%s", path)

	err = providerserver.Serve(ctx, provider.New, providerserver.ServeOpts{
		Address: "pexip.com/pexip/pexip",
	})
	if err != nil {
		logger.Errorf(ctx, "failed to serve provider: %s", err.Error())
		os.Exit(1)
	}
}
