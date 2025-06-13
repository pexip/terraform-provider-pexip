# Pexip Terraform Provider

This is the Pexip Terraform provider, which allows you to manage Pexip Infinity using
Terraform and to automate the provisioning of Pexip Infinity.


## Example usage

```terraform
terraform {
  required_providers {
    pexip = "=> 1.0.0"
  }
}

provider "pexip" {
  address = "https://manager.example.com"
  username = "admin"
  password = "password"
}

# Creating a Pexip Infinity Manager bootstrap config
data "pexip_infinity_manager_config" "config" {
  hostname              = "test-mgr1"
  domain                = "dev.vcops.tech"
  ip                    = "10.0.0.40"
  mask                  = "255.255.255.0"
  gw                    = "10.0.0.1"
  dns                   = "1.1.1.1"
  ntp                   = "pool.ntp.org"
  user                  = "admin"
  pass                  = "admin_password"
  admin_password        = "admin_password"
  error_reports         = false
  enable_analytics      = false
  contact_email_address = "vcops@pexip.com"
}

# Registering a Pexip Infinity node with the Pexip Infinity Manager
resource "pexip_infinity_node" "test-node-1" {
  name = "test-node-1"
  config = "<config from manager>"
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