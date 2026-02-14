# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

This is the official Terraform Provider for Pexip Infinity, a comprehensive Infrastructure as Code solution for managing Pexip video conferencing platforms. The provider is built using the HashiCorp Terraform Plugin Framework and integrates with the Pexip Infinity Manager API through the `go-infinity-sdk/v38` package.

@../go-infinity-sdk/CLAUDE.md

## Architecture

### Core Components

1. **Provider Entry Point** (`main.go`)
   - Minimal entry point that calls `providerserver.Serve()`
   - Delegates to `internal/provider.New()`

2. **Provider Implementation** (`internal/provider/provider.go`)
   - `PexipProvider` struct manages the Infinity SDK client
   - Uses a `sync.Mutex` for thread-safe operations
   - Configures authentication (address, username, password, insecure mode)
   - Implements `InfinityClient` interface wrapping the SDK client

3. **Resource Pattern** (`internal/provider/resource_infinity_*.go`)
   - 80+ resources following consistent naming: `resource_infinity_<name>.go`
   - Each resource implements Terraform Plugin Framework interfaces
   - Standard CRUD operations: Create, Read, Update, Delete
   - Import support via `ResourceWithImportState`
   - Resources use the SDK's `config.Service` for CRUD operations

4. **Testing**
   - **Unit Tests** (`*_test.go` with `tags=unit`): Use mocked SDK client
   - **Integration Tests** (`*_integration_test.go` with `tags=integration`): Require real Pexip environment with `TF_ACC=1`
   - Test fixtures in `testdata/resource_*` directories contain `.tf` files
   - Mock client from `infinity.NewClientMock()` with `stretchr/testify/mock`

5. **Test pattern**
   1. Refer to the go-infinity-sdk configuration schema
   2. The test folders should be named after the resource and end in `min` and `full`. For example `resource_infinity_worker_vm_full` and `resource_infinity_worker_vm_min`.
   3. If the resource has test folders that end in `basic` or `basic_updated` then they should be deleted.
   4. the terraform config for full should have every field populated
   5. the terraform config for min should only have the required fields populated
   6. there should be four tests. Test 1 creates the fully populated resource. Test 2 updates the resource using the min configuration then deletes the resource. Test 3 creates the resource using the min config and test 4 updates it using the full config
   7. Any field that is a related resource should also be created in Terraform and referenced by id, not hardcoded
   8. Any resource created in Terraform should have tf-test in the name
   9. if the integration test fails due to a field not being reset or cleared, check the go-infinity-sdk model for the resource and remove omitempty if it is configured on the attribute in the update method for that resource

6. **SDK Integration**
   - Primary dependency: `github.com/pexip/go-infinity-sdk/v38`
   - SDK provides typed API clients: `config.Service`, `status.Service`, `history.Service`, `command.Service`
   - Resources map Terraform schema to SDK request/response types
   - Error handling via SDK's HTTP response types

### Key Directories

- `internal/provider/` - All provider resources, data sources, and validators (230+ files)
- `internal/helpers/` - Utility functions (hashing, environment variables)
- `internal/version/` - Semantic versioning implementation
- `internal/test/` - Test utilities and helper functions
- `testdata/` - Terraform configuration fixtures for acceptance tests
- `integration_test/` - Full integration test configurations with Terraform files
- `example/` - Complete deployment examples (GCP/AWS/OpenStack)

### Resource Implementation Pattern

Each resource follows this structure:

```go
type InfinityXResource struct {
    InfinityClient InfinityClient
}

type InfinityXResourceModel struct {
    ID          types.String `tfsdk:"id"`
    ResourceID  types.Int32  `tfsdk:"resource_id"`
    Name        types.String `tfsdk:"name"`
    // ... other fields
}

// Implements: resource.Resource, resource.ResourceWithImportState
func (r *InfinityXResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse)
func (r *InfinityXResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse)
func (r *InfinityXResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse)
func (r *InfinityXResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse)
func (r *InfinityXResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse)
func (r *InfinityXResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse)
```

### Custom Validators

Located in `internal/provider/validators/`:
- `domain_validator.go` - Domain name validation
- `email_validator.go` - Email address validation
- `ip_address_validator.go` - IP address validation
- `netmask_validator.go` - Netmask validation
- `url_validator.go` - URL validation

### Testing Strategy

1. **Unit Tests** (Fast, Offline):
   - Mock the `InfinityClient` interface
   - Test schema validation, state management, and business logic
   - Use `infinity.NewClientMock()` with expectations
   - Example: `resource_infinity_worker_vm_test.go` shows complex mocking with related resources
   - Create an individual mock for each related resource. Keep things simple. The main goal is to get the tests to pass. Refactoring will be done at a later time.

2. **Integration Tests** (Slow, Requires Pexip):
   - Requires `TF_ACC=1` environment variable
   - Connects to real Pexip Infinity Manager
   - Tests actual API interactions
   - Uses fixtures from `testdata/` directories

3. **Terraform Validation**:
   - CI runs `terraform fmt -check -recursive ./example/`
   - CI runs `terraform validate` against examples

## Code Quality Requirements

### Linting Configuration

The project uses `golangci-lint` with strict rules (`.golangci.yml`):
- Enabled linters: bodyclose, depguard, errcheck, gosec, govet, staticcheck, and 15+ more
- **Dependency restrictions**: Only `tflog` logging allowed (no `logrus`)
- Cyclomatic complexity limit: 50
- Test files have relaxed linting (no errcheck, gosec, etc.)

### Coding Standards

1. **Terraform Plugin Framework**: All resources use Plugin Framework (not Plugin SDK v2)
2. **Logging**: Use `github.com/hashicorp/terraform-plugin-log/tflog` only
3. **Error Handling**: Always check and propagate errors appropriately
4. **State Management**: Use `resp.Diagnostics.AddError()` for user-facing errors
5. **Imports**: Use `goimports` with local prefix `github.com/pexip/terraform-provider-pexip`

## Common Development Patterns

### Adding a New Resource

1. Create `internal/provider/resource_infinity_<name>.go`
2. Create `internal/provider/resource_infinity_<name>_test.go` (unit tests with mocks)
3. Create `internal/provider/resource_infinity_<name>_integration_test.go` (if applicable)
4. Add test fixtures in `testdata/resource_infinity_<name>_basic/`
5. Register resource in `provider.go` Resources() method
6. Run tests: `make test && TF_ACC=1 make testacc`
7. Run linting: `make lint`

### Working with the SDK

The SDK client is accessed via the `InfinityClient` interface:

```go
// Create
resp, err := r.InfinityClient.Config().PostWithResponse(ctx, "configuration/v1/conference/", req, response)

// Read
err := r.InfinityClient.Config().GetJSON(ctx, "configuration/v1/conference/123/", nil, &resource)

// Update
err := r.InfinityClient.Config().PatchJSON(ctx, "configuration/v1/conference/123/", req, nil)

// Delete
err := r.InfinityClient.Config().DeleteJSON(ctx, "configuration/v1/conference/123/", nil)
```

### Test Data Fixtures

Each resource test uses fixtures in `testdata/`:
- `testdata/resource_infinity_<name>_basic/` - Minimal required configuration
- `testdata/resource_infinity_<name>_full/` - Full configuration with all fields
- `testdata/resource_infinity_<name>_basic_updated/` - Updated configuration for update tests

Each directory contains:
- `providers.tf` - Provider configuration (if needed)
- `resources.tf` - Resource definitions

## CI/CD Workflows

### GitHub Actions

- **test.yml** - Unit tests, linting, Terraform validation (runs on push/PR)
- **integration-test.yml** - Full integration tests against real Pexip environment
- **reusable_build.yml** - Build and package provider binary

### Test Matrix

- Go versions: 1.24
- Terraform versions: 1.5, 1.6, 1.7 (validation), 1.12 (tests)

## Important Considerations

### Environment Variables

When running tests locally or in CI:
- `TF_ACC=1` - Enable acceptance tests
- `PEXIP_ADDRESS` - Pexip Manager URL
- `PEXIP_USERNAME` - Authentication username
- `PEXIP_PASSWORD` - Authentication password
- `PEXIP_INSECURE` - Allow self-signed certificates (optional)
- `GOPRIVATE=github.com/pexip` - Required for private SDK dependency

## Documentation

- Main docs: `README.md` - Comprehensive user guide with examples
- Contributing: `CONTRIBUTING.md` - Development guidelines and processes
- Resource docs: `docs/resources/infinity_*.md` - Generated documentation per resource
- Data source docs: `docs/data-sources/infinity_*.md` - Data source documentation
