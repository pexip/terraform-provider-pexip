/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/action/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ action.Action              = (*InfinityDeleteDefaultMgrTLSCertificateAction)(nil)
	_ action.ActionWithConfigure = (*InfinityDeleteDefaultMgrTLSCertificateAction)(nil)
)

type InfinityDeleteDefaultMgrTLSCertificateAction struct {
	InfinityClient InfinityClient
}

func (a *InfinityDeleteDefaultMgrTLSCertificateAction) Metadata(ctx context.Context, req action.MetadataRequest, resp *action.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_delete_default_mgr_tls_certificate"
}

func (a *InfinityDeleteDefaultMgrTLSCertificateAction) Configure(ctx context.Context, req action.ConfigureRequest, resp *action.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	p, ok := req.ProviderData.(*PexipProvider)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Action Configure Type",
			fmt.Sprintf("Expected *PexipProvider, got: %T. Please report this issue to the provider developers", req.ProviderData),
		)
		return
	}

	a.InfinityClient = p.client
}

func (a *InfinityDeleteDefaultMgrTLSCertificateAction) Schema(ctx context.Context, req action.SchemaRequest, resp *action.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Deletes the default management node TLS certificate (ID 1) that is automatically created when a management node is deployed. This is an action that performs a one-time deletion operation.",
	}
}

func (a *InfinityDeleteDefaultMgrTLSCertificateAction) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	// Delete the default TLS certificate (ID 1)
	const defaultCertificateID = 1

	tflog.Info(ctx, fmt.Sprintf("Deleting default TLS certificate (ID %d)", defaultCertificateID))

	err := a.InfinityClient.Config().DeleteTLSCertificate(ctx, defaultCertificateID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Default TLS Certificate",
			fmt.Sprintf("Could not delete default TLS certificate (ID %d): %s", defaultCertificateID, err),
		)
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Successfully deleted default TLS certificate (ID %d)", defaultCertificateID))

	// Report progress that the deletion is complete
	resp.SendProgress(action.InvokeProgressEvent{
		Message: fmt.Sprintf("Successfully deleted default TLS certificate (ID %d)", defaultCertificateID),
	})
}
