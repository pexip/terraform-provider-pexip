# GitHub Actions Workflows

This directory contains the CI/CD workflows for the Terraform Provider for Pexip Infinity.

## Workflow Architecture

```
┌─────────────┐    ┌───────────────┐    ┌────────────────────────┐
│   test.yml  │ -> │  release.yml  │ -> │  publish-assets.yml    │
└─────────────┘    └───────────────┘    └────────────────────────┘
```

## Local Development

To test workflows locally, you can use [act](https://github.com/nektos/act):

```bash
# Test the test workflow
act --container-architecture linux/amd64 -W .github/workflows/test.yml

# Test the build workflow (will call test workflow)
act --container-architecture linux/amd64 -W .github/workflows/publish-assets.yml
```

Note: Some jobs may need secrets or specific environments to run successfully.