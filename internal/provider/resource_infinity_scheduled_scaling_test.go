package provider

import (
	"os"
	"testing"

	"github.com/pexip/go-infinity-sdk/v38/config"
	"github.com/pexip/go-infinity-sdk/v38/types"
	"github.com/stretchr/testify/mock"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/pexip/go-infinity-sdk/v38"

	"github.com/pexip/terraform-provider-pexip/internal/test"
)

func TestInfinityScheduledScaling(t *testing.T) {
	_ = os.Setenv("TF_ACC", "1")

	// Create a mock client and set up expectations
	client := infinity.NewClientMock()

	// Mock the CreateScheduledscaling API call
	createResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/scheduled_scaling/123/",
	}
	client.On("PostWithResponse", mock.Anything, "configuration/v1/scheduled_scaling/", mock.Anything, mock.Anything).Return(createResponse, nil)

	// Shared state for mocking
	mockState := &config.ScheduledScaling{
		ID:                 123,
		ResourceURI:        "/api/admin/configuration/v1/scheduled_scaling/123/",
		PolicyName:         "scheduled_scaling-test",
		PolicyType:         "worker_vm",
		ResourceIdentifier: "test-value",
		Enabled:            true,
		LocalTimezone:      "test-value",
		StartDate:          "test-value",
		TimeFrom:           "test-value",
		TimeTo:             "test-value",
		InstancesToAdd:     2,
		MinutesInAdvance:   15,
		Mon:                true,
		Tue:                true,
		Wed:                true,
		Thu:                true,
		Fri:                true,
		Sat:                true,
		Sun:                true,
	}

	// Mock the GetScheduledscaling API call for Read operations
	client.On("GetJSON", mock.Anything, "configuration/v1/scheduled_scaling/123/", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		scheduled_scaling := args.Get(2).(*config.ScheduledScaling)
		*scheduled_scaling = *mockState
	}).Maybe()

	// Mock the UpdateScheduledscaling API call
	client.On("PutJSON", mock.Anything, "configuration/v1/scheduled_scaling/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		updateReq := args.Get(2).(*config.ScheduledScalingUpdateRequest)
		scheduled_scaling := args.Get(3).(*config.ScheduledScaling)

		// Update mock state based on request
		if updateReq.PolicyType != "" {
			mockState.PolicyType = updateReq.PolicyType
		}
		if updateReq.ResourceIdentifier != "" {
			mockState.ResourceIdentifier = updateReq.ResourceIdentifier
		}
		if updateReq.Enabled != nil {
			mockState.Enabled = *updateReq.Enabled
		}
		if updateReq.LocalTimezone != "" {
			mockState.LocalTimezone = updateReq.LocalTimezone
		}
		if updateReq.StartDate != "" {
			mockState.StartDate = updateReq.StartDate
		}
		if updateReq.TimeFrom != "" {
			mockState.TimeFrom = updateReq.TimeFrom
		}
		if updateReq.TimeTo != "" {
			mockState.TimeTo = updateReq.TimeTo
		}
		if updateReq.InstancesToAdd != nil {
			mockState.InstancesToAdd = *updateReq.InstancesToAdd
		}
		if updateReq.MinutesInAdvance != nil {
			mockState.MinutesInAdvance = *updateReq.MinutesInAdvance
		}
		if updateReq.Mon != nil {
			mockState.Mon = *updateReq.Mon
		}
		if updateReq.Tue != nil {
			mockState.Tue = *updateReq.Tue
		}
		if updateReq.Wed != nil {
			mockState.Wed = *updateReq.Wed
		}
		if updateReq.Thu != nil {
			mockState.Thu = *updateReq.Thu
		}
		if updateReq.Fri != nil {
			mockState.Fri = *updateReq.Fri
		}
		if updateReq.Sat != nil {
			mockState.Sat = *updateReq.Sat
		}
		if updateReq.Sun != nil {
			mockState.Sun = *updateReq.Sun
		}

		// Return updated state
		*scheduled_scaling = *mockState
	}).Maybe()

	// Mock the DeleteScheduledscaling API call
	client.On("DeleteJSON", mock.Anything, mock.MatchedBy(func(path string) bool {
		return path == "configuration/v1/scheduled_scaling/123/"
	}), mock.Anything).Return(nil)

	testInfinityScheduledScaling(t, client)
}

func testInfinityScheduledScaling(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: getTestProtoV5ProviderFactories(client),
		Steps: []resource.TestStep{
			{
				Config: test.LoadTestFolder(t, "resource_infinity_scheduled_scaling_basic"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_scheduled_scaling.scheduled_scaling-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_scheduled_scaling.scheduled_scaling-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_scheduled_scaling.scheduled_scaling-test", "policy_name", "scheduled_scaling-test"),
					resource.TestCheckResourceAttr("pexip_infinity_scheduled_scaling.scheduled_scaling-test", "enabled", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_scheduled_scaling.scheduled_scaling-test", "mon", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_scheduled_scaling.scheduled_scaling-test", "tue", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_scheduled_scaling.scheduled_scaling-test", "wed", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_scheduled_scaling.scheduled_scaling-test", "thu", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_scheduled_scaling.scheduled_scaling-test", "fri", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_scheduled_scaling.scheduled_scaling-test", "sat", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_scheduled_scaling.scheduled_scaling-test", "sun", "true"),
				),
			},
			{
				Config: test.LoadTestFolder(t, "resource_infinity_scheduled_scaling_basic_updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_scheduled_scaling.scheduled_scaling-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_scheduled_scaling.scheduled_scaling-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_scheduled_scaling.scheduled_scaling-test", "policy_name", "scheduled_scaling-test"),
					resource.TestCheckResourceAttr("pexip_infinity_scheduled_scaling.scheduled_scaling-test", "policy_type", "management_vm"),
					resource.TestCheckResourceAttr("pexip_infinity_scheduled_scaling.scheduled_scaling-test", "resource_identifier", "updated-value"),
					resource.TestCheckResourceAttr("pexip_infinity_scheduled_scaling.scheduled_scaling-test", "enabled", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_scheduled_scaling.scheduled_scaling-test", "instances_to_add", "3"),
					resource.TestCheckResourceAttr("pexip_infinity_scheduled_scaling.scheduled_scaling-test", "minutes_in_advance", "30"),
					resource.TestCheckResourceAttr("pexip_infinity_scheduled_scaling.scheduled_scaling-test", "mon", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_scheduled_scaling.scheduled_scaling-test", "tue", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_scheduled_scaling.scheduled_scaling-test", "wed", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_scheduled_scaling.scheduled_scaling-test", "thu", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_scheduled_scaling.scheduled_scaling-test", "fri", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_scheduled_scaling.scheduled_scaling-test", "sat", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_scheduled_scaling.scheduled_scaling-test", "sun", "false"),
				),
			},
		},
	})
}
