# GitHub Actions Workflows

This directory contains the CI/CD workflows for the Terraform Provider for Pexip Infinity.

## Workflow Files

### Main Workflows

- **test.yml**: Runs on push to master and pull requests. Executes unit tests, terraform validation, linting, and dependency review.
- **release.yml**: Manual workflow to create releases and trigger asset publishing.
- **publish-assets.yml**: Builds and publishes provider binaries for multiple platforms.
- **integration-test.yml**: Runs integration tests against a live Pexip Infinity environment. Triggered by cron schedule, manual dispatch, or workflow call.

### Reusable Workflows

- **reusable_build.yml**: Builds the Terraform provider for specified OS/architecture combinations.
- **reusable_install_provider.yml**: Installs the built provider for testing.
- **reusable_deploy_infinity.yml**: Deploys a Pexip Infinity instance for integration testing.
- **reusable_destroy_infinity.yml**: Cleans up Pexip Infinity instances after integration testing.

## Workflow Architecture

```
CI/CD Pipeline:
┌─────────────┐    ┌───────────────┐    ┌────────────────────────┐
│   test.yml  │    │  release.yml  │ -> │  publish-assets.yml    │
└─────────────┘    └───────────────┘    └────────────────────────┘
       │                                            │
       v                                            v
┌─────────────────┐                    ┌─────────────────────────┐
│ reusable_build  │ <----------------> │  Multi-platform builds  │
└─────────────────┘                    └─────────────────────────┘

Integration Testing (Scheduled/Manual):
┌────────────────────────┐
│  integration-test.yml  │
└────────────────────────┘
           │
           ├─> reusable_build.yml
           ├─> reusable_install_provider.yml
           ├─> reusable_deploy_infinity.yml
           ├─> [Run Integration Tests]
           └─> reusable_destroy_infinity.yml
```

## Triggers

- **Push/PR**: test.yml runs automatically
- **Manual**: release.yml, publish-assets.yml, integration-test.yml can be triggered manually
- **Scheduled**: integration-test.yml runs weekly (Sunday at midnight)
- **Workflow Call**: All main workflows except release.yml support being called by other workflows

## Local Development

To test workflows locally, you can use [act](https://github.com/nektos/act):

```bash
# Test the test workflow
act --container-architecture linux/amd64 -W .github/workflows/test.yml

# Test the build workflow
act --container-architecture linux/amd64 -W .github/workflows/publish-assets.yml
```

Note: Some jobs may need secrets or specific environments to run successfully.