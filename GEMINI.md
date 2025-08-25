# Project Overview

This project is a Terraform provider for Pexip Infinity. It allows users to manage their Pexip Infinity infrastructure as code, covering a wide range of features from basic node configuration to integrations with Microsoft 365 and Google Workspace. The provider is written in Go and uses the Pexip Go SDK to interact with the Pexip Infinity API.

## Key Technologies

*   **Go:** The programming language used to develop the Terraform provider.
*   **Terraform:** The infrastructure as code tool that this provider plugs into.
*   **Pexip Infinity API:** The API that this provider communicates with to manage Pexip Infinity resources.

## Architecture

The project follows the standard structure for a Terraform provider. The core logic is located in the `internal/provider` directory, with the main entry point being `main.go`. The provider defines a set of resources that correspond to Pexip Infinity objects, such as conferences, locations, and worker VMs.

# Building and Running

The project uses a `Makefile` to streamline the development process. Here are the key commands:

*   **Build the provider:**
    ```bash
    make build
    ```

*   **Run unit tests:**
    ```bash
    make test
    ```

*   **Run acceptance tests (requires a Pexip Infinity environment):**
    ```bash
    make testacc
    ```

*   **Lint the code:**
    ```bash
    make lint
    ```

*   **Format the code:**
    ```bash
    make fmt
    ```

*   **Install the provider locally:**
    ```bash
    make install
    ```

# Development Conventions

## Coding Style

The project follows standard Go coding conventions. Code should be formatted with `gofmt` before committing.

## Testing

The project has both unit and acceptance tests. Unit tests are located alongside the code they test and can be run with `make test`. Acceptance tests require a running Pexip Infinity instance and are located in the `internal/provider` directory. They can be run with `make testacc`.

## Contribution Guidelines

Contributions are welcome. Please see the `CONTRIBUTING.md` file for more details.
