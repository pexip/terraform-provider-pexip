# Pexip Terraform Provider

This is the Pexip Terraform provider, which allows you to manage Pexip Infinity using
Terraform. The provider is designed tobe used to automate the provisioning of Pexip Infinity.


## Example usage

```terraform
terraform {
  required_providers {
    pexip = "~> 1.0"
  }
}

provider "pexip" {
}

data "infinity_manager_config" "config" {
  # set config values here and fetch the rendered config from the .render attribute
  # this is used as input to the compute resource metadata field.
}

data "infinity_manager" "manager" {
  # information about the Infinity manager
  name = "my-infinity-manager"
  address = "manager.example.com"
  username = "admin"
  password = "password"
}

resource "infinity_node" "node" {
  manager = data.infinity_manager.manager.id
}
```