# Terraform Provider Pexip - Plugin Framework Migration Analysis

## Current State Analysis

### ✅ Already Migrated to Plugin Framework:
1. **Provider Core** (`internal/provider/provider.go`) - ✅ COMPLETE
   - Uses Plugin Framework interfaces
   - Provider configuration and schema defined with Framework
   - Client setup and authentication working

2. **Main Entry Point** (`main.go`) - ✅ COMPLETE
   - Uses `providerserver.Serve` 
   - Plugin Framework ready

3. **Resource Implementations** - ✅ COMPLETE
   - All 77+ resources have been migrated to Plugin Framework
   - Using `resource.Resource` interface
   - Schema definitions using Plugin Framework syntax
   - CRUD operations implemented with Framework

4. **Data Sources** - ✅ COMPLETE
   - `InfinityManagerConfigDataSource` migrated

### ❌ Still Using Plugin SDK v2:
1. **All Test Files** - 75 test files still use SDK v2
   - Using `github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource`
   - Test configurations and assertions use SDK patterns

### Dependencies Analysis:
- `go.mod` contains both frameworks:
  - ✅ `github.com/hashicorp/terraform-plugin-framework v1.15.1` (target)
  - ❌ `github.com/hashicorp/terraform-plugin-sdk/v2 v2.37.0` (legacy - only used by tests)

## Migration Strategy

### Phase 1: Test Migration (Primary Focus - 80% effort)
The main work remaining is migrating 75+ test files from SDK v2 to Plugin Framework testing patterns.

#### Test Migration Approach:
1. **Unit Tests**: Convert to use Plugin Framework's `resource.TestCase`
2. **Integration Tests**: Update provider factory and test configurations
3. **Test Utilities**: Update test helpers and utilities

#### Key Changes Needed:
- Replace `github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource` 
- With `github.com/hashicorp/terraform-plugin-testing/helper/resource`
- Update provider factory functions
- Update test configurations and assertions

### Phase 2: Dependency Cleanup (20% effort)
1. Remove SDK v2 dependency from `go.mod`
2. Clean up any remaining SDK v2 imports
3. Update CI/build configurations if needed

## Resource Coverage
All resources are Framework-ready:
- 77+ resources migrated
- 1 data source migrated
- Provider core migrated
- Authentication and client setup working

## Risk Assessment: LOW
- Provider core already working with Framework
- Resources already migrated and functional
- Only tests need migration
- No breaking changes to end users