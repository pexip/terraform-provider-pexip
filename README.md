# Pexip Terraform Provider

This is the Pexip Terraform provider, which allows you to manage Pexip Infinity using
Terraform and to automate the provisioning of Pexip Infinity.


## Example usage

```terraform
terraform {
  required_providers {
    pexip = "~> 1.0"
  }
}

provider "pexip" {
  # information about the Infinity manager
  address = "https://manager.example.com"
  username = "admin"
  password = "password"
}

data "infinity_manager_config" "config" {
  # set config values here and fetch the rendered config from the .render attribute
  # this is used as input to the compute resource metadata field.
}

data "infinity_manager" "manager" {
  address = "https://manager.example.com"
  username = "admin"
  password = "password"
}

resource "infinity_node" "node" {
  manager = data.infinity_manager.manager.id
}
```

### Multiple Provider Configurations
You can optionally define multiple configurations for the same provider,
and select which one to use on a per-resource or per-module basis. The primary reason
for this is to support multiple regions for a cloud platform; other examples include
targeting multiple Docker hosts, multiple Consul hosts, etc.

To create multiple configurations for a given provider, include multiple provider blocks with the
same provider name. For each additional non-default configuration, use the alias meta-argument to
provide an extra name segment

## How to build
To build an test the plugin locally first create a `~/.terraformrc` file

```shell
provider_installation {

  dev_overrides {
    "<username>/pexip" = "/Users/<username>/.terraform.d/plugins"
  }
  direct {}
}
```

Then build and install the plugin locally using

```shell
make install
```

## Running tests
To run the internal unit tests run test `test` make target

```shell
make test
```

To run terraform acceptance tests, the `TF_ACC` env variable must be set to true before making the
`test` make target, or the `testacc` make target can be used

```shell
make testacc
```