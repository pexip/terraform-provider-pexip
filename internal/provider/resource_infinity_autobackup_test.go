/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package provider

import (
	"os"
	"regexp"
	"testing"

	"github.com/pexip/go-infinity-sdk/v38/config"
	"github.com/stretchr/testify/mock"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/pexip/go-infinity-sdk/v38"

	"github.com/pexip/terraform-provider-pexip/internal/test"
)

func TestInfinityAutobackupValidator(t *testing.T) {
	t.Parallel()
	_ = os.Setenv("TF_ACC", "1")

	client := infinity.NewClientMock()

	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: getTestProtoV5ProviderFactories(client),
		Steps: []resource.TestStep{
			{
				Config: `resource "pexip_infinity_autobackup" "autobackup-test" {
  autobackup_enabled = true
}`,
				ExpectError: regexp.MustCompile(`autobackup_passphrase must be set when autobackup_enabled is true`),
			},
		},
	})
}

func TestInfinityAutobackup(t *testing.T) {
	t.Parallel()
	_ = os.Setenv("TF_ACC", "1")

	client := infinity.NewClientMock()

	mockState := &config.Autobackup{
		ResourceURI:              "/api/admin/configuration/v1/autobackup/1/",
		AutobackupEnabled:        false,
		AutobackupInterval:       24,
		AutobackupPassphrase:     "",
		AutobackupStartHour:      1,
		AutobackupUploadURL:      "",
		AutobackupUploadUsername: "",
		AutobackupUploadPassword: "",
	}

	client.On("GetJSON", mock.Anything, "configuration/v1/autobackup/1/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		autobackup := args.Get(3).(*config.Autobackup)
		*autobackup = *mockState
	}).Maybe()

	client.On("PatchJSON", mock.Anything, "configuration/v1/autobackup/1/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		updateReq := args.Get(2).(*config.AutobackupUpdateRequest)
		result := args.Get(3).(*config.Autobackup)

		if updateReq.AutobackupEnabled != nil {
			mockState.AutobackupEnabled = *updateReq.AutobackupEnabled
		}
		if updateReq.AutobackupInterval != nil {
			mockState.AutobackupInterval = *updateReq.AutobackupInterval
		}
		mockState.AutobackupPassphrase = updateReq.AutobackupPassphrase
		if updateReq.AutobackupStartHour != nil {
			mockState.AutobackupStartHour = *updateReq.AutobackupStartHour
		}
		mockState.AutobackupUploadURL = updateReq.AutobackupUploadURL
		mockState.AutobackupUploadUsername = updateReq.AutobackupUploadUsername
		mockState.AutobackupUploadPassword = updateReq.AutobackupUploadPassword

		*result = *mockState
	}).Maybe()

	testInfinityAutobackup(t, client)
}

func testInfinityAutobackup(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: getTestProtoV5ProviderFactories(client),
		Steps: []resource.TestStep{
			{
				// Step 1: Create with full config
				Config: test.LoadTestFolder(t, "resource_infinity_autobackup_full"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_autobackup.autobackup-test", "id"),
					resource.TestCheckResourceAttr("pexip_infinity_autobackup.autobackup-test", "autobackup_enabled", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_autobackup.autobackup-test", "autobackup_interval", "12"),
					resource.TestCheckResourceAttr("pexip_infinity_autobackup.autobackup-test", "autobackup_passphrase", "SecretPassphrase123"),
					resource.TestCheckResourceAttr("pexip_infinity_autobackup.autobackup-test", "autobackup_start_hour", "2"),
					resource.TestCheckResourceAttr("pexip_infinity_autobackup.autobackup-test", "autobackup_upload_url", "ftp://backup.example.com/pexip"),
					resource.TestCheckResourceAttr("pexip_infinity_autobackup.autobackup-test", "autobackup_upload_username", "backupuser"),
					resource.TestCheckResourceAttr("pexip_infinity_autobackup.autobackup-test", "autobackup_upload_password", "BackupPassword123"),
				),
			},
			{
				// Step 2: Update to min config
				Config: test.LoadTestFolder(t, "resource_infinity_autobackup_min"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_autobackup.autobackup-test", "id"),
					resource.TestCheckResourceAttr("pexip_infinity_autobackup.autobackup-test", "autobackup_enabled", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_autobackup.autobackup-test", "autobackup_interval", "24"),
					resource.TestCheckResourceAttr("pexip_infinity_autobackup.autobackup-test", "autobackup_passphrase", ""),
					resource.TestCheckResourceAttr("pexip_infinity_autobackup.autobackup-test", "autobackup_start_hour", "1"),
					resource.TestCheckResourceAttr("pexip_infinity_autobackup.autobackup-test", "autobackup_upload_url", ""),
					resource.TestCheckResourceAttr("pexip_infinity_autobackup.autobackup-test", "autobackup_upload_username", ""),
					resource.TestCheckResourceAttr("pexip_infinity_autobackup.autobackup-test", "autobackup_upload_password", ""),
				),
			},
			{
				// Step 3: Destroy (no-op for singleton, but included for consistency)
				Config:  test.LoadTestFolder(t, "resource_infinity_autobackup_min"),
				Destroy: true,
			},
			{
				// Step 4: Recreate with min config (actually just updates since it's a singleton)
				Config: test.LoadTestFolder(t, "resource_infinity_autobackup_min"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_autobackup.autobackup-test", "id"),
					resource.TestCheckResourceAttr("pexip_infinity_autobackup.autobackup-test", "autobackup_enabled", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_autobackup.autobackup-test", "autobackup_interval", "24"),
					resource.TestCheckResourceAttr("pexip_infinity_autobackup.autobackup-test", "autobackup_start_hour", "1"),
				),
			},
			{
				// Step 5: Update to full config
				Config: test.LoadTestFolder(t, "resource_infinity_autobackup_full"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_autobackup.autobackup-test", "id"),
					resource.TestCheckResourceAttr("pexip_infinity_autobackup.autobackup-test", "autobackup_enabled", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_autobackup.autobackup-test", "autobackup_interval", "12"),
					resource.TestCheckResourceAttr("pexip_infinity_autobackup.autobackup-test", "autobackup_passphrase", "SecretPassphrase123"),
					resource.TestCheckResourceAttr("pexip_infinity_autobackup.autobackup-test", "autobackup_start_hour", "2"),
					resource.TestCheckResourceAttr("pexip_infinity_autobackup.autobackup-test", "autobackup_upload_url", "ftp://backup.example.com/pexip"),
					resource.TestCheckResourceAttr("pexip_infinity_autobackup.autobackup-test", "autobackup_upload_username", "backupuser"),
					resource.TestCheckResourceAttr("pexip_infinity_autobackup.autobackup-test", "autobackup_upload_password", "BackupPassword123"),
				),
			},
		},
	})
}
