/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package main

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"

	"github.com/pexip/terraform-provider-pexip/internal/log"
	pexipProvider "github.com/pexip/terraform-provider-pexip/internal/provider"
	"github.com/pexip/terraform-provider-pexip/internal/version"
)

func createProvider() func() provider.Provider {
	return pexipProvider.New
}

func main() {
	ctx := context.Background()
	logger := log.NewTerraformLogger()
	path, err := os.Getwd()
	if err != nil {
		logger.Errorf(ctx, "failed to initialize provider: %s", err.Error())
	}

	logger.Infof(ctx, "%s", version.Version().String())
	logger.Infof(ctx, "%s", path)

	err = providerserver.Serve(ctx, createProvider(), providerserver.ServeOpts{
		Address: "pexip.com/pexip/pexip",
	})
	if err != nil {
		logger.Errorf(ctx, "failed to serve provider: %s", err.Error())
		os.Exit(1)
	}
}
