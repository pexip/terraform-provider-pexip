# Contributing to Terraform Provider for Pexip Infinity

Thank you for your interest in contributing to the Terraform Provider for Pexip Infinity! This document provides guidelines and information for contributors.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Making Changes](#making-changes)
- [Testing](#testing)
- [Code Quality](#code-quality)
- [Pull Request Process](#pull-request-process)
- [Release Process](#release-process)
- [Getting Help](#getting-help)

## Code of Conduct

This project and everyone participating in it is governed by our Code of Conduct. By participating, you are expected to uphold this code. Please report unacceptable behavior to the project maintainers.

## Getting Started

### Prerequisites

- [Go](https://golang.org/doc/install) >= 1.21
- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [golangci-lint](https://golangci-lint.run/usage/install/) for code quality checks
- Access to a Pexip Infinity Manager >= v38 for integration testing
- Git for version control

### Development Dependencies

```bash
# Install golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Install terraform-docs (optional, for documentation generation)
go install github.com/terraform-docs/terraform-docs@latest
```

## Development Setup

1. **Fork and Clone the Repository**
   ```bash
   git clone https://github.com/YOUR_USERNAME/terraform-provider-pexip.git
   cd terraform-provider-pexip
   ```

2. **Install Dependencies**
   ```bash
   go mod download
   go mod verify
   ```

3. **Set Up Local Development Environment**
   ```bash
   # Build and install locally
   make build-dev
   
   # Create ~/.terraformrc for local development
   cat > ~/.terraformrc << EOF
   provider_installation {
     dev_overrides {
       "pexip/pexip" = "$HOME/.terraform.d/plugins"
     }
     direct {}
   }
   EOF
   ```

4. **Verify Setup**
   ```bash
   make test
   ```

## Making Changes

### Branch Naming Convention

- `feature/description` - New features
- `fix/description` - Bug fixes
- `docs/description` - Documentation updates
- `refactor/description` - Code refactoring
- `test/description` - Test improvements

### Coding Standards

- **Go Style**: Follow [Effective Go](https://golang.org/doc/effective_go.html) and [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- **Terraform Plugin Framework**: Use the [Terraform Plugin Framework](https://developer.hashicorp.com/terraform/plugin/framework) patterns
- **Resource Naming**: Follow the pattern `resource_infinity_<resource_name>.go`
- **Test Naming**: Follow the pattern `resource_infinity_<resource_name>_test.go`

### File Structure

```
internal/
├── provider/               # Provider implementation
│   ├── resource_infinity_*.go     # Resource implementations
│   ├── data_source_*.go          # Data source implementations
│   ├── validators/               # Custom validators
│   └── *_test.go                # Unit tests
├── helpers/               # Utility functions
├── log/                  # Logging configuration
└── test/                 # Test utilities

testdata/                 # Test configuration files
├── resource_infinity_*/
│   ├── resources.tf      # Basic resource configuration
│   └── providers.tf      # Provider configuration

example/                  # Example configurations
docs/                    # Generated documentation
```

### Adding a New Resource

1. **Create the Resource File**
   ```bash
   # Example: Adding a new DNS server resource
   touch internal/provider/resource_infinity_dns_server.go
   ```

2. **Implement the Resource Interface**
   ```go
   type InfinityDNSServerResource struct {
       InfinityClient InfinityClient
   }

   func (r *InfinityDNSServerResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
       resp.TypeName = req.ProviderTypeName + "_infinity_dns_server"
   }

   func (r *InfinityDNSServerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
       // Define schema
   }

   func (r *InfinityDNSServerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
       // Implement create logic
   }

   // Implement Read, Update, Delete methods
   ```

3. **Register the Resource**
   Add to `internal/provider/provider.go`:
   ```go
   func (p *PexipProvider) Resources(ctx context.Context) []func() resource.Resource {
       return []func() resource.Resource{
           // ... existing resources
           NewInfinityDNSServerResource,
       }
   }
   ```

4. **Create Test Files**
   ```bash
   # Unit tests
   touch internal/provider/resource_infinity_dns_server_test.go
   
   # Integration tests
   touch internal/provider/resource_infinity_dns_server_integration_test.go
   
   # Test data
   mkdir -p testdata/resource_infinity_dns_server_basic
   touch testdata/resource_infinity_dns_server_basic/resources.tf
   touch testdata/resource_infinity_dns_server_basic/providers.tf
   ```

### Testing Requirements

All changes must include appropriate tests:

- **Unit Tests**: Test business logic with mocked dependencies
- **Integration Tests**: Test against real API (when possible)
- **Example Configuration**: Provide working Terraform configuration

## Testing

### Running Tests

```bash
# Run all unit tests
make test

# Run unit tests with verbose output
go test -v -tags=unit ./internal/provider

# Run specific test
go test -v -tags=unit ./internal/provider -run TestInfinityDNSServer

# Run integration tests (requires real Pexip environment)
make testacc

# Run integration tests for specific resource
TF_ACC=1 go test -v -tags=integration ./internal/provider -run TestInfinityDNSServer
```

### Writing Unit Tests

Unit tests should mock the Infinity client:

```go
package provider

import (
    "testing"
    
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
    "github.com/pexip/go-infinity-sdk/v38"
    "github.com/stretchr/testify/mock"
)

func TestInfinityDNSServer(t *testing.T) {
    t.Parallel()
    
    client := infinity.NewClientMock()
    // Set up mock expectations
    
    resource.Test(t, resource.TestCase{
        ProtoV5ProviderFactories: getTestProtoV5ProviderFactories(client),
        Steps: []resource.TestStep{
            {
                Config: test.LoadTestFolder(t, "resource_infinity_dns_server_basic"),
                Check: resource.ComposeTestCheckFunc(
                    resource.TestCheckResourceAttrSet("pexip_infinity_dns_server.test", "id"),
                ),
            },
        },
    })
}
```

### Writing Integration Tests

Integration tests should use the `//go:build integration` build tag and require `TF_ACC=1`:

```go
//go:build integration

package provider

import (
    "os"
    "testing"
    
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestInfinityDNSServerIntegration(t *testing.T) {
    if os.Getenv("TF_ACC") == "" {
        t.Skip("TF_ACC not set, skipping integration test")
    }
    
    resource.Test(t, resource.TestCase{
        ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
        Steps: []resource.TestStep{
            {
                Config: testAccInfinityDNSServerConfig,
                Check: resource.ComposeTestCheckFunc(
                    resource.TestCheckResourceAttr("pexip_infinity_dns_server.test", "address", "192.168.1.1"),
                ),
            },
        },
    })
}
```

## Code Quality

### Linting

We use golangci-lint for code quality checks:

```bash
# Run all linters
make lint

# Run specific linter
golangci-lint run --enable=gosec

# Fix auto-fixable issues
golangci-lint run --fix
```

### Code Formatting

```bash
# Format code
make fmt

# Check imports
goimports -w -local github.com/pexip/terraform-provider-pexip .
```

### Security

We use gosec for security analysis:

```bash
# Run security analysis
make sec

# Or directly
gosec ./...
```

## Pull Request Process

1. **Create a Feature Branch**
   ```bash
   git checkout -b feature/add-dns-server-resource
   ```

2. **Make Your Changes**
   - Follow coding standards
   - Add tests for new functionality
   - Update documentation if needed

3. **Test Your Changes**
   ```bash
   make check  # Runs lint and test
   ```

4. **Commit Your Changes**
   ```bash
   git add .
   git commit -m "feat: add DNS server resource
   
   - Add InfinityDNSServerResource with CRUD operations
   - Include unit and integration tests
   - Add example configuration
   
   Closes #123"
   ```

5. **Push and Create Pull Request**
   ```bash
   git push origin feature/add-dns-server-resource
   ```

6. **Pull Request Requirements**
   - Descriptive title and description
   - Reference related issues
   - All CI checks must pass
   - At least one maintainer approval
   - Up-to-date with main branch

### Pull Request Template

When creating a pull request, please include:

- **Description**: What changes were made and why
- **Type of Change**: Bug fix, new feature, documentation, etc.
- **Testing**: How the changes were tested
- **Checklist**: Confirm all requirements are met

## Release Process

Releases are managed by project maintainers following semantic versioning:

- **Major** (x.0.0): Breaking changes
- **Minor** (0.x.0): New features, backward compatible
- **Patch** (0.0.x): Bug fixes, backward compatible

## Getting Help

- **Issues**: Use GitHub Issues for bug reports and feature requests
- **Discussions**: Use GitHub Discussions for questions and general discussion
- **Documentation**: Check the [README.md](README.md) and [docs/](docs/) directory
- **Pexip Support**: For Pexip Infinity specific questions, consult [Pexip Documentation](https://docs.pexip.com/)

## License

By contributing to this project, you agree that your contributions will be licensed under the Apache License 2.0. See [LICENSE](LICENSE) for details.

## Additional Resources

- [Terraform Plugin Framework Documentation](https://developer.hashicorp.com/terraform/plugin/framework)
- [Pexip Infinity API Documentation](https://docs.pexip.com/admin/admin_api.htm)
- [Go Infinity SDK Documentation](https://github.com/pexip/go-infinity-sdk)
- [Terraform Provider Best Practices](https://developer.hashicorp.com/terraform/plugin/best-practices)

---

Thank you for contributing to the Terraform Provider for Pexip Infinity! 🎉