# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

This is a Terraform provider for Pexip Infinity that allows managing Pexip Infinity infrastructure through Terraform. The provider uses the Pexip go-infinity-sdk v38 to interact with Pexip Infinity Manager APIs.

## Development Commands

### Building and Installing
- `make install` - Build and install the provider locally to `~/.terraform.d/plugins`
- `make build` - Build the provider binary to `dist/` directory
- `make build-dev` - Build directly to Terraform plugins directory for development

### Testing
- `make test` - Run unit tests (tagged as `unit`)
- `make testacc` - Run Terraform acceptance tests (sets `TF_ACC=true`, tagged as `integration`)
- `make check` - Run linting and tests together

### Code Quality
- `make lint` - Run golangci-lint with specific checks (govet, ineffassign, staticcheck, deadcode, unused)
- `make fmt` - Format Go code using `go fmt`
- `make sec` - Run security analysis with gosec

### Local Development Setup
Requires creating `~/.terraformrc` with dev_overrides pointing to local plugin directory as shown in README.md.

## Architecture

### Provider Structure
- **Main entry point**: `main.go` - Standard Terraform plugin serve setup
- **Provider core**: `internal/provider/provider.go` - Defines PexipProvider with Infinity client integration
- **Authentication**: Uses basic auth with username/password against Pexip Infinity Manager API
- **Concurrency**: Provider includes mutex for thread-safe operations

### Resources and Data Sources
- **InfinityNodeResource**: Manages Pexip Infinity nodes
- **InfinityManagerConfigDataSource**: Generates bootstrap configuration for Pexip Infinity Manager
- Models are defined in separate files (e.g., `infinity_manager_config_model.go`)

### Key Dependencies
- Terraform Plugin Framework v1.15.0 (newer framework, not SDK v2)
- Pexip go-infinity-sdk v38 for API interactions
- Custom validators in `internal/provider/validators/` for IP addresses and URLs

### Directory Structure
- `internal/provider/` - Core provider implementation
- `internal/helpers/` - Utility functions and hashing
- `internal/log/` - Terraform-specific logging setup
- `testdata/` - Terraform configuration files for testing
- `example/` - Complete example Terraform configurations

### Testing Patterns
- Unit tests use `unit` build tag
- Acceptance tests use `integration` build tag and require `TF_ACC=true`
- Test files follow `*_test.go` naming convention