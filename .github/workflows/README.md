# GitHub Actions Workflows

This directory contains the CI/CD workflows for the Terraform Provider for Pexip Infinity.

## Workflow Architecture

```
┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│   test.yml  │ -> │  build.yml  │ -> │  release    │
└─────────────┘    └─────────────┘    └─────────────┘
```

### test.yml - Test Workflow

**Purpose**: Runs all quality assurance checks before allowing builds.

**Triggers**:
- Push to main/master branches
- Pull requests to main/master
- Daily scheduled runs (2 AM UTC)
- Called by other workflows (`workflow_call`)

**Jobs**:
- **test**: Unit tests with Go 1.21-1.24 and Terraform 1.0-1.6 matrix
- **lint**: Code formatting, linting, and go mod tidy checks
- **security**: Security scanning with gosec and Nancy vulnerability scanner
- **terraform**: Terraform format and validation checks
- **acceptance**: Integration tests (main branch only, requires secrets)
- **dependency-review**: Dependency vulnerability review (PRs only)

### build.yml - Build and Release Workflow

**Purpose**: Builds binaries and handles releases after tests pass.

**Triggers**:
- Push to main/master branches
- Pull requests to main/master
- Tags starting with `v*` (for releases)
- Can be called by other workflows

**Jobs**:
- **test**: Calls test.yml workflow to ensure all checks pass
- **build**: Creates binaries for multiple OS/architecture combinations
- **release**: GoReleaser-based releases (tags only, requires GPG secrets)

**Dependencies**:
- Build job only runs if test workflow succeeds
- Release job only runs if both test and build succeed

## Secrets Required

### For Releases
- `GPG_PRIVATE_KEY`: GPG private key for signing releases
- `PASSPHRASE`: Passphrase for GPG key
- `GITHUB_TOKEN`: Automatically provided by GitHub

### For Acceptance Tests (Optional)
- `PEXIP_ADDRESS`: Pexip Infinity Manager URL
- `PEXIP_USERNAME`: Username for authentication
- `PEXIP_PASSWORD`: Password for authentication

## Configuration Files

- `.golangci.yml`: Linting configuration
- `.goreleaser.yml`: Release configuration
- `VERSION`: Version file for builds

## Workflow Execution Examples

### Pull Request
1. test.yml runs all jobs
2. build.yml calls test.yml then builds binaries
3. Results reported on PR

### Push to Main
1. test.yml runs all jobs (including acceptance tests if secrets configured)
2. build.yml builds and uploads artifacts

### Release (Git Tag)
1. test.yml runs all jobs
2. build.yml builds binaries
3. build.yml releases using GoReleaser with signed artifacts

## Local Development

To test workflows locally, you can use [act](https://github.com/nektos/act):

```bash
# Test the test workflow
act -W --container-architecture linux/amd64 .github/workflows/test.yml

# Test the build workflow (will call test workflow)
act -W  --container-architecture linux/amd64.github/workflows/build.yml
```

Note: Some jobs may need secrets or specific environments to run successfully.