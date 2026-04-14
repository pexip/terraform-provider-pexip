/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package provider

import (
	"os"
	"testing"

	"github.com/pexip/go-infinity-sdk/v38/config"
	"github.com/pexip/go-infinity-sdk/v38/types"
	"github.com/stretchr/testify/mock"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/pexip/go-infinity-sdk/v38"

	"github.com/pexip/terraform-provider-pexip/internal/test"
)

func TestInfinityMjxMeetingProcessingRule(t *testing.T) {
	t.Parallel()
	_ = os.Setenv("TF_ACC", "1")

	client := infinity.NewClientMock()

	mockState := &config.MjxMeetingProcessingRule{}

	createResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/mjx_meeting_processing_rule/123/",
	}
	client.On("PostWithResponse", mock.Anything, "configuration/v1/mjx_meeting_processing_rule/", mock.Anything, mock.Anything).Return(createResponse, nil).Run(func(args mock.Arguments) {
		createReq := args.Get(2).(*config.MjxMeetingProcessingRuleCreateRequest)
		*mockState = config.MjxMeetingProcessingRule{
			ID:                       123,
			ResourceURI:              "/api/admin/configuration/v1/mjx_meeting_processing_rule/123/",
			Name:                     createReq.Name,
			Description:              createReq.Description,
			Priority:                 createReq.Priority,
			Enabled:                  createReq.Enabled,
			MeetingType:              createReq.MeetingType,
			MjxIntegration:           createReq.MjxIntegration,
			MatchString:              createReq.MatchString,
			ReplaceString:            createReq.ReplaceString,
			TransformRule:            createReq.TransformRule,
			CustomTemplate:           createReq.CustomTemplate,
			Domain:                   createReq.Domain,
			CompanyID:                createReq.CompanyID,
			IncludePin:               createReq.IncludePin,
			DefaultProcessingEnabled: createReq.DefaultProcessingEnabled,
		}
	})

	client.On("GetJSON", mock.Anything, "configuration/v1/mjx_meeting_processing_rule/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		rule := args.Get(3).(*config.MjxMeetingProcessingRule)
		*rule = *mockState
	}).Maybe()

	client.On("PutJSON", mock.Anything, "configuration/v1/mjx_meeting_processing_rule/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		updateReq := args.Get(2).(*config.MjxMeetingProcessingRuleUpdateRequest)
		mockState.Name = updateReq.Name
		mockState.Description = updateReq.Description
		mockState.Priority = updateReq.Priority
		mockState.Enabled = updateReq.Enabled
		mockState.MeetingType = updateReq.MeetingType
		mockState.MjxIntegration = updateReq.MjxIntegration
		mockState.MatchString = updateReq.MatchString
		mockState.ReplaceString = updateReq.ReplaceString
		mockState.TransformRule = updateReq.TransformRule
		mockState.CustomTemplate = updateReq.CustomTemplate
		mockState.Domain = updateReq.Domain
		mockState.CompanyID = updateReq.CompanyID
		mockState.IncludePin = updateReq.IncludePin
		mockState.DefaultProcessingEnabled = updateReq.DefaultProcessingEnabled
	}).Maybe()

	client.On("DeleteJSON", mock.Anything, "configuration/v1/mjx_meeting_processing_rule/123/", mock.Anything).Return(nil)

	testInfinityMjxMeetingProcessingRule(t, client)
}

func testInfinityMjxMeetingProcessingRule(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: getTestProtoV6ProviderFactories(client),
		Steps: []resource.TestStep{
			// Step 1: Create with full config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_mjx_meeting_processing_rule_full"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_meeting_processing_rule.test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_meeting_processing_rule.test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_meeting_processing_rule.test", "name", "tf-test mjx-meeting-processing-rule full"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_meeting_processing_rule.test", "description", "Test MJX meeting processing rule"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_meeting_processing_rule.test", "priority", "10"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_meeting_processing_rule.test", "enabled", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_meeting_processing_rule.test", "meeting_type", "teams"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_meeting_processing_rule.test", "mjx_integration", "/api/admin/configuration/v1/mjx_integration/1/"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_meeting_processing_rule.test", "transform_rule", "{{ domain }}"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_meeting_processing_rule.test", "domain", "example.com"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_meeting_processing_rule.test", "company_id", "test-company-id"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_meeting_processing_rule.test", "include_pin", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_meeting_processing_rule.test", "default_processing_enabled", "false"),
				),
			},
			// Step 2: Update with min config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_mjx_meeting_processing_rule_min"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_meeting_processing_rule.test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_meeting_processing_rule.test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_meeting_processing_rule.test", "name", "tf-test mjx-meeting-processing-rule min"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_meeting_processing_rule.test", "description", ""),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_meeting_processing_rule.test", "priority", "1"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_meeting_processing_rule.test", "enabled", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_meeting_processing_rule.test", "meeting_type", "pexipinfinity"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_meeting_processing_rule.test", "mjx_integration", "/api/admin/configuration/v1/mjx_integration/1/"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_meeting_processing_rule.test", "match_string", ""),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_meeting_processing_rule.test", "replace_string", ""),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_meeting_processing_rule.test", "domain", ""),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_meeting_processing_rule.test", "company_id", ""),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_meeting_processing_rule.test", "include_pin", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_meeting_processing_rule.test", "default_processing_enabled", "true"),
				),
			},
			// Step 3: Destroy
			{
				Config:  test.LoadTestFolder(t, "resource_infinity_mjx_meeting_processing_rule_min"),
				Destroy: true,
			},
			// Step 4: Recreate with min config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_mjx_meeting_processing_rule_min"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_meeting_processing_rule.test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_meeting_processing_rule.test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_meeting_processing_rule.test", "name", "tf-test mjx-meeting-processing-rule min"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_meeting_processing_rule.test", "priority", "1"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_meeting_processing_rule.test", "meeting_type", "pexipinfinity"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_meeting_processing_rule.test", "mjx_integration", "/api/admin/configuration/v1/mjx_integration/1/"),
				),
			},
			// Step 5: Update to full config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_mjx_meeting_processing_rule_full"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_meeting_processing_rule.test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_meeting_processing_rule.test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_meeting_processing_rule.test", "name", "tf-test mjx-meeting-processing-rule full"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_meeting_processing_rule.test", "description", "Test MJX meeting processing rule"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_meeting_processing_rule.test", "priority", "10"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_meeting_processing_rule.test", "enabled", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_meeting_processing_rule.test", "meeting_type", "teams"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_meeting_processing_rule.test", "mjx_integration", "/api/admin/configuration/v1/mjx_integration/1/"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_meeting_processing_rule.test", "transform_rule", "{{ domain }}"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_meeting_processing_rule.test", "domain", "example.com"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_meeting_processing_rule.test", "company_id", "test-company-id"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_meeting_processing_rule.test", "include_pin", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_meeting_processing_rule.test", "default_processing_enabled", "false"),
				),
			},
		},
	})
}
