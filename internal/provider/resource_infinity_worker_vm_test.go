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

func TestInfinityWorkerVM(t *testing.T) {
	t.Parallel()
	_ = os.Setenv("TF_ACC", "1")

	// Create a mock client and set up expectations
	client := infinity.NewClientMock()

	// Mock the CreateSystemLocation API call (needed because worker VM references it)
	systemLocationCreateResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/system_location/1/",
	}
	client.On("PostWithResponse", mock.Anything, "configuration/v1/system_location/", mock.Anything, mock.Anything).Return(systemLocationCreateResponse, nil)

	// Mock the GetSystemLocation API call for Read operations
	mockSystemLocation := &config.SystemLocation{
		ID:                        1,
		ResourceURI:               "/api/admin/configuration/v1/system_location/1/",
		Name:                      "provider test system location",
		MTU:                       1500,
		MediaQoS:                  test.IntPtr(0),
		SignallingQoS:             test.IntPtr(0),
		BDPMPinChecksEnabled:      "GLOBAL",
		BDPMScanQuarantineEnabled: "GLOBAL",
		UseRelayCandidatesOnly:    false,
	}
	client.On("GetJSON", mock.Anything, "configuration/v1/system_location/1/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		systemLocation := args.Get(3).(*config.SystemLocation)
		*systemLocation = *mockSystemLocation
	}).Maybe()

	// Mock the DeleteSystemLocation API call
	client.On("DeleteJSON", mock.Anything, "configuration/v1/system_location/1/", mock.Anything).Return(nil)

	// Mock the CreateSSHAuthorizedKey API call (needed because worker VM references it)
	sshKeyCreateResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/ssh_authorized_key/3/",
	}
	client.On("PostWithResponse", mock.Anything, "configuration/v1/ssh_authorized_key/", mock.Anything, mock.Anything).Return(sshKeyCreateResponse, nil)

	// Mock the GetSSHAuthorizedKey API call for Read operations
	mockSSHKey := &config.SSHAuthorizedKey{
		ID:          3,
		ResourceURI: "/api/admin/configuration/v1/ssh_authorized_key/3/",
		Keytype:     "ssh-rsa",
		Key:         "AAAAB3NzaC1yc2EAAAADAQABAAABgQC7",
		Comment:     "Test SSH key for worker VM",
	}
	client.On("GetJSON", mock.Anything, "configuration/v1/ssh_authorized_key/3/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		sshKey := args.Get(3).(*config.SSHAuthorizedKey)
		*sshKey = *mockSSHKey
	}).Maybe()

	// Mock the DeleteSSHAuthorizedKey API call
	client.On("DeleteJSON", mock.Anything, "configuration/v1/ssh_authorized_key/3/", mock.Anything).Return(nil)

	// Mock the CreateStaticRoute API call (needed because worker VM references it)
	staticRouteCreateResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/static_route/4/",
	}
	client.On("PostWithResponse", mock.Anything, "configuration/v1/static_route/", mock.Anything, mock.Anything).Return(staticRouteCreateResponse, nil)

	// Mock the GetStaticRoute API call for Read operations
	mockStaticRoute := &config.StaticRoute{
		ID:          4,
		ResourceURI: "/api/admin/configuration/v1/static_route/4/",
		Name:        "test-static-route",
		Address:     "10.0.0.0",
		Prefix:      24,
		Gateway:     "192.168.1.254",
	}
	client.On("GetJSON", mock.Anything, "configuration/v1/static_route/4/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		staticRoute := args.Get(3).(*config.StaticRoute)
		*staticRoute = *mockStaticRoute
	}).Maybe()

	// Mock the DeleteStaticRoute API call
	client.On("DeleteJSON", mock.Anything, "configuration/v1/static_route/4/", mock.Anything).Return(nil)

	// Mock the CreateTLSCertificate API call (needed because worker VM references it)
	// Note: We need to store/return certificate data for Terraform state consistency,
	// but we don't validate it - that's tested in the dedicated TLS certificate tests
	tlsCertificateCreateResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/tls_certificate/2/",
	}
	var mockTLSCertState *config.TLSCertificate

	client.On("PostWithResponse", mock.Anything, "configuration/v1/tls_certificate/", mock.Anything, mock.Anything).Return(tlsCertificateCreateResponse, nil).Run(func(args mock.Arguments) {
		createReq := args.Get(2).(*config.TLSCertificateCreateRequest)
		mockTLSCertState = &config.TLSCertificate{
			ID:          2,
			ResourceURI: "/api/admin/configuration/v1/tls_certificate/2/",
			Certificate: createReq.Certificate,
			PrivateKey:  createReq.PrivateKey,
			Nodes:       []string{},
		}
	})

	client.On("GetJSON", mock.Anything, "configuration/v1/tls_certificate/2/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		tlsCert := args.Get(3).(*config.TLSCertificate)
		if mockTLSCertState != nil {
			*tlsCert = *mockTLSCertState
		}
	}).Maybe()

	client.On("PatchJSON", mock.Anything, "configuration/v1/tls_certificate/2/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		updateReq := args.Get(2).(*config.TLSCertificateUpdateRequest)
		tlsCert := args.Get(3).(*config.TLSCertificate)
		if mockTLSCertState != nil {
			if updateReq.Certificate != "" {
				mockTLSCertState.Certificate = updateReq.Certificate
			}
			if updateReq.PrivateKey != "" {
				mockTLSCertState.PrivateKey = updateReq.PrivateKey
			}
			if updateReq.Nodes != nil {
				mockTLSCertState.Nodes = updateReq.Nodes
			}
			*tlsCert = *mockTLSCertState
		}
	}).Maybe()

	client.On("DeleteJSON", mock.Anything, "configuration/v1/tls_certificate/2/", mock.Anything).Return(nil)

	// Shared mock state - will be initialized on first create
	mockState := &config.WorkerVM{}

	// Mock the CreateWorkervm API call
	createResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/worker_vm/123/",
	}
	client.On("PostWithResponse", mock.Anything, "configuration/v1/worker_vm/", mock.Anything, mock.Anything).Return(createResponse, nil).Run(func(args mock.Arguments) {
		createReq := args.Get(2).(*config.WorkerVMCreateRequest)
		// Reinitialize mockState from create request (important for destroy/recreate cycles)
		*mockState = config.WorkerVM{
			ID:                        123,
			ResourceURI:               "/api/admin/configuration/v1/worker_vm/123/",
			Name:                      createReq.Name,
			Hostname:                  createReq.Hostname,
			Domain:                    createReq.Domain,
			AlternativeFQDN:           createReq.AlternativeFQDN,
			Address:                   createReq.Address,
			Netmask:                   createReq.Netmask,
			Gateway:                   createReq.Gateway,
			SystemLocation:            createReq.SystemLocation,
			TLSCertificate:            createReq.TLSCertificate,
			Description:               createReq.Description,
			IPv6Address:               createReq.IPv6Address,
			IPv6Gateway:               createReq.IPv6Gateway,
			NodeType:                  createReq.NodeType,
			DeploymentType:            createReq.DeploymentType,
			Password:                  createReq.Password,
			MaintenanceMode:           createReq.MaintenanceMode,
			MaintenanceModeReason:     createReq.MaintenanceModeReason,
			VMCPUCount:                createReq.VMCPUCount,
			VMSystemMemory:            createReq.VMSystemMemory,
			SecondaryAddress:          createReq.SecondaryAddress,
			SecondaryNetmask:          createReq.SecondaryNetmask,
			MediaPriorityWeight:       createReq.MediaPriorityWeight,
			SSHAuthorizedKeysUseCloud: createReq.SSHAuthorizedKeysUseCloud,
			StaticNATAddress:          createReq.StaticNATAddress,
			SNMPAuthenticationPassword: createReq.SNMPAuthenticationPassword,
			SNMPCommunity:             createReq.SNMPCommunity,
			SNMPMode:                  createReq.SNMPMode,
			SNMPPrivacyPassword:       createReq.SNMPPrivacyPassword,
			SNMPSystemContact:         createReq.SNMPSystemContact,
			SNMPSystemLocation:        createReq.SNMPSystemLocation,
			SNMPUsername:              createReq.SNMPUsername,
			EnableSSH:                 createReq.EnableSSH,
			EnableDistributedDatabase: createReq.EnableDistributedDatabase,
			SSHAuthorizedKeys:         createReq.SSHAuthorizedKeys,
			CloudBursting:             false, // Always false per API behavior
		}
		// Convert StaticRoutes URIs to objects (API returns objects, not strings)
		if len(createReq.StaticRoutes) > 0 {
			mockState.StaticRoutes = []config.StaticRoute{
				{
					ID:          mockStaticRoute.ID,
					ResourceURI: mockStaticRoute.ResourceURI,
					Name:        mockStaticRoute.Name,
					Address:     mockStaticRoute.Address,
					Prefix:      mockStaticRoute.Prefix,
					Gateway:     mockStaticRoute.Gateway,
				},
			}
		} else {
			mockState.StaticRoutes = []config.StaticRoute{}
		}
		// Compute transcoding based on node_type (API behavior - this field is computed)
		if createReq.NodeType == "CONFERENCING" || createReq.NodeType == "TRANSCODING" {
			mockState.Transcoding = true
		} else {
			mockState.Transcoding = false
		}
	})

	// Mock the GetWorkervm API call for Read operations
	client.On("GetJSON", mock.Anything, "configuration/v1/worker_vm/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		workerVM := args.Get(3).(*config.WorkerVM)
		*workerVM = *mockState
	}).Maybe()

	// Mock the UpdateWorkervm API call
	client.On("PutJSON", mock.Anything, "configuration/v1/worker_vm/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		updateReq := args.Get(2).(*config.WorkerVMUpdateRequest)
		workerVM := args.Get(3).(*config.WorkerVM)

		// Update mock state based on request - required fields
		mockState.Name = updateReq.Name
		mockState.Hostname = updateReq.Hostname
		mockState.Domain = updateReq.Domain
		mockState.Address = updateReq.Address
		mockState.Netmask = updateReq.Netmask
		mockState.Gateway = updateReq.Gateway
		mockState.SystemLocation = updateReq.SystemLocation
		mockState.EnableDistributedDatabase = updateReq.EnableDistributedDatabase

		// Optional string fields with defaults - use request value or default
		mockState.AlternativeFQDN = updateReq.AlternativeFQDN
		mockState.Description = updateReq.Description
		mockState.Password = updateReq.Password
		mockState.NodeType = updateReq.NodeType
		mockState.DeploymentType = updateReq.DeploymentType
		mockState.EnableSSH = updateReq.EnableSSH
		mockState.SNMPMode = updateReq.SNMPMode
		mockState.MaintenanceModeReason = updateReq.MaintenanceModeReason
		mockState.SNMPAuthenticationPassword = updateReq.SNMPAuthenticationPassword
		mockState.SNMPCommunity = updateReq.SNMPCommunity
		mockState.SNMPPrivacyPassword = updateReq.SNMPPrivacyPassword
		mockState.SNMPSystemContact = updateReq.SNMPSystemContact
		mockState.SNMPSystemLocation = updateReq.SNMPSystemLocation
		mockState.SNMPUsername = updateReq.SNMPUsername

		// Boolean fields
		mockState.MaintenanceMode = updateReq.MaintenanceMode
		mockState.SSHAuthorizedKeysUseCloud = updateReq.SSHAuthorizedKeysUseCloud

		// Compute transcoding based on node_type (API behavior)
		if updateReq.NodeType == "CONFERENCING" || updateReq.NodeType == "TRANSCODING" {
			mockState.Transcoding = true
		} else {
			mockState.Transcoding = false
		}

		// Integer fields
		mockState.VMCPUCount = updateReq.VMCPUCount
		mockState.VMSystemMemory = updateReq.VMSystemMemory

		// Fields with defaults not in update request
		mockState.CloudBursting = false

		// Nullable pointer fields - set to nil if not provided in request
		mockState.TLSCertificate = updateReq.TLSCertificate
		mockState.IPv6Address = updateReq.IPv6Address
		mockState.IPv6Gateway = updateReq.IPv6Gateway
		mockState.SecondaryAddress = updateReq.SecondaryAddress
		mockState.SecondaryNetmask = updateReq.SecondaryNetmask
		mockState.MediaPriorityWeight = updateReq.MediaPriorityWeight
		mockState.StaticNATAddress = updateReq.StaticNATAddress

		// List fields
		mockState.SSHAuthorizedKeys = updateReq.SSHAuthorizedKeys
		// Convert StaticRoutes URIs to objects (API returns objects, not strings)
		if len(updateReq.StaticRoutes) > 0 {
			mockState.StaticRoutes = []config.StaticRoute{
				{
					ID:          mockStaticRoute.ID,
					ResourceURI: mockStaticRoute.ResourceURI,
					Name:        mockStaticRoute.Name,
					Address:     mockStaticRoute.Address,
					Prefix:      mockStaticRoute.Prefix,
					Gateway:     mockStaticRoute.Gateway,
				},
			}
		} else {
			mockState.StaticRoutes = []config.StaticRoute{}
		}

		// Return updated state
		*workerVM = *mockState
	}).Maybe()

	// Mock the DeleteWorkervm API call
	client.On("DeleteJSON", mock.Anything, mock.MatchedBy(func(path string) bool {
		return path == "configuration/v1/worker_vm/123/"
	}), mock.Anything).Return(nil)

	testInfinityWorkerVM(t, client)
}

func testInfinityWorkerVM(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: getTestProtoV5ProviderFactories(client),
		ExternalProviders: map[string]resource.ExternalProvider{
			"tls": {
				Source: "hashicorp/tls",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: test.LoadTestFolder(t, "resource_infinity_worker_vm_full"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_worker_vm.worker-vm-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_worker_vm.worker-vm-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "name", "tf-test-full_worker-vm"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "hostname", "worker-vm-test"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "domain", "test-value"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "address", "192.168.1.10"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "netmask", "255.255.255.0"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "gateway", "192.168.1.1"),
					resource.TestCheckResourceAttrPair(
						"pexip_infinity_worker_vm.worker-vm-test", "system_location",
						"pexip_infinity_system_location.test", "id",
					),
					resource.TestCheckResourceAttrPair(
						"pexip_infinity_worker_vm.worker-vm-test", "tls_certificate",
						"pexip_infinity_tls_certificate.test", "id",
					),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "alternative_fqdn", "alt.example.com"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "ipv6_address", "2001:db8::1"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "ipv6_gateway", "2001:db8::fe"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "description", "initial description"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "transcoding", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "maintenance_mode", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "maintenance_mode_reason", "test-value"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "vm_cpu_count", "4"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "vm_system_memory", "4096"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "secondary_address", "172.16.0.10"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "secondary_netmask", "255.255.255.0"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "media_priority_weight", "10"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "static_nat_address", "203.0.113.2"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "snmp_system_location", "test-value"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "snmp_authentication_password", "auth-password1"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "snmp_community", "public1"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "snmp_mode", "STANDARD"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "snmp_privacy_password", "privacy-password1"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "snmp_system_contact", "snmpcontact1@domain.com"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "snmp_username", "snmp-user1"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "enable_ssh", "ON"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "password", "password-initial"),
					resource.TestCheckTypeSetElemAttrPair("pexip_infinity_worker_vm.worker-vm-test", "ssh_authorized_keys.*", "pexip_infinity_ssh_authorized_key.test", "id"),
					resource.TestCheckTypeSetElemAttrPair("pexip_infinity_worker_vm.worker-vm-test", "static_routes.*", "pexip_infinity_static_route.test", "id"),
				),
			},
			{
				Config: test.LoadTestFolder(t, "resource_infinity_worker_vm_min"),
				Check: resource.ComposeTestCheckFunc(
					// IDs and required fields
					resource.TestCheckResourceAttrSet("pexip_infinity_worker_vm.worker-vm-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_worker_vm.worker-vm-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "name", "tf-test-min_worker-vm"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "hostname", "worker-vm-test"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "domain", "test-value"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "address", "192.168.1.10"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "netmask", "255.255.255.0"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "gateway", "192.168.1.1"),
					resource.TestCheckResourceAttrSet("pexip_infinity_worker_vm.worker-vm-test", "system_location"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "password", "password-initial"),

					// RequiresReplace fields kept
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "ipv6_address", "2001:db8::1"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "ipv6_gateway", "2001:db8::fe"),

					// Optional fields cleared - verify defaults
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "alternative_fqdn", ""),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "description", ""),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "maintenance_mode", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "maintenance_mode_reason", ""),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "node_type", "CONFERENCING"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "deployment_type", "MANUAL-PROVISION-ONLY"),
					resource.TestCheckResourceAttrSet("pexip_infinity_worker_vm.worker-vm-test", "transcoding"), // Computed only, just verify it exists
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "enable_ssh", "GLOBAL"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "snmp_mode", "DISABLED"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "snmp_system_contact", "admin@domain.com"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "snmp_system_location", "Virtual machine"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "snmp_community", "public"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "snmp_authentication_password", ""),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "snmp_privacy_password", ""),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "snmp_username", ""),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "vm_cpu_count", "4"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "vm_system_memory", "4096"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "enable_distributed_database", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "ssh_authorized_keys_use_cloud", "true"),

					// Fields with defaults that should be reset to default when cleared
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "media_priority_weight", "0"),

					// Nullable fields cleared via UPDATE - API should return null
					resource.TestCheckNoResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "tls_certificate"),
					resource.TestCheckNoResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "secondary_address"),
					resource.TestCheckNoResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "secondary_netmask"),
					resource.TestCheckNoResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "static_nat_address"),
				),
			},
			{
				// Step 3: Destroy and recreate with minimal config
				Config:  test.LoadTestFolder(t, "resource_infinity_worker_vm_min"),
				Destroy: true,
			},
			{
				// Step 4: Recreate with minimal config (after destroy)
				Config: test.LoadTestFolder(t, "resource_infinity_worker_vm_min"),
				Check: resource.ComposeTestCheckFunc(
					// IDs and required fields
					resource.TestCheckResourceAttrSet("pexip_infinity_worker_vm.worker-vm-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_worker_vm.worker-vm-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "name", "tf-test-min_worker-vm"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "hostname", "worker-vm-test"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "domain", "test-value"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "address", "192.168.1.10"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "netmask", "255.255.255.0"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "gateway", "192.168.1.1"),
					resource.TestCheckResourceAttrSet("pexip_infinity_worker_vm.worker-vm-test", "system_location"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "password", "password-initial"),

					// RequiresReplace fields kept
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "ipv6_address", "2001:db8::1"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "ipv6_gateway", "2001:db8::fe"),

					// Optional fields cleared - verify defaults
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "alternative_fqdn", ""),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "description", ""),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "maintenance_mode", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "maintenance_mode_reason", ""),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "node_type", "CONFERENCING"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "deployment_type", "MANUAL-PROVISION-ONLY"),
					resource.TestCheckResourceAttrSet("pexip_infinity_worker_vm.worker-vm-test", "transcoding"), // Computed only, just verify it exists
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "enable_ssh", "GLOBAL"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "snmp_mode", "DISABLED"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "snmp_system_contact", "admin@domain.com"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "snmp_system_location", "Virtual machine"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "snmp_community", "public"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "snmp_authentication_password", ""),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "snmp_privacy_password", ""),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "snmp_username", ""),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "vm_cpu_count", "4"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "vm_system_memory", "4096"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "enable_distributed_database", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "ssh_authorized_keys_use_cloud", "true"),

					// Fields with defaults that should be reset to default when cleared
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "media_priority_weight", "0"),

					// Nullable fields not specified in CREATE - not all asserted
					// Note: API may assign defaults on CREATE (e.g., tls_certificate)
					// We only assert null for fields cleared via UPDATE (see Step 2)
					resource.TestCheckNoResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "secondary_address"),
					resource.TestCheckNoResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "secondary_netmask"),
					resource.TestCheckNoResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "static_nat_address"),
				),
			},
			{
				// Step 5: Update to full config
				Config: test.LoadTestFolder(t, "resource_infinity_worker_vm_full"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_worker_vm.worker-vm-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_worker_vm.worker-vm-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "name", "tf-test-full_worker-vm"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "hostname", "worker-vm-test"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "domain", "test-value"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "address", "192.168.1.10"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "netmask", "255.255.255.0"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "gateway", "192.168.1.1"),
					resource.TestCheckResourceAttrPair(
						"pexip_infinity_worker_vm.worker-vm-test", "system_location",
						"pexip_infinity_system_location.test", "id",
					),
					resource.TestCheckResourceAttrPair(
						"pexip_infinity_worker_vm.worker-vm-test", "tls_certificate",
						"pexip_infinity_tls_certificate.test", "id",
					),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "alternative_fqdn", "alt.example.com"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "ipv6_address", "2001:db8::1"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "ipv6_gateway", "2001:db8::fe"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "description", "initial description"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "transcoding", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "maintenance_mode", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "maintenance_mode_reason", "test-value"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "vm_cpu_count", "4"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "vm_system_memory", "4096"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "secondary_address", "172.16.0.10"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "secondary_netmask", "255.255.255.0"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "media_priority_weight", "10"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "static_nat_address", "203.0.113.2"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "snmp_system_location", "test-value"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "snmp_authentication_password", "auth-password1"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "snmp_community", "public1"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "snmp_mode", "STANDARD"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "snmp_privacy_password", "privacy-password1"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "snmp_system_contact", "snmpcontact1@domain.com"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "snmp_username", "snmp-user1"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "enable_ssh", "ON"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "password", "password-initial"),
					resource.TestCheckTypeSetElemAttrPair("pexip_infinity_worker_vm.worker-vm-test", "ssh_authorized_keys.*", "pexip_infinity_ssh_authorized_key.test", "id"),
					resource.TestCheckTypeSetElemAttrPair("pexip_infinity_worker_vm.worker-vm-test", "static_routes.*", "pexip_infinity_static_route.test", "id"),
				),
			},
		},
	})
}
