/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package main

import (
	"context"
	"log"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"

	pexipProvider "github.com/pexip/terraform-provider-pexip/internal/provider"
)

func createProvider() func() provider.Provider {
	return pexipProvider.New
}

func main() {
	err := providerserver.Serve(context.Background(), createProvider(), providerserver.ServeOpts{
		Address: "registry.terraform.io/pexip/pexip",
	})
	if err != nil {
		log.Printf("failed to serve provider: %s", err)
		os.Exit(1)
	}
}
