terraform {
  required_providers {
    pexip = {
      source  = "pexip.com/pexip/pexip"
      version = "0.0.1"
    }
    google = {
      source  = "hashicorp/google"
      version = "6.25.0"
    }
    random = {
      source  = "hashicorp/random"
      version = "3.7.1"
    }
    local = {
      source  = "hashicorp/local"
      version = "2.4.0"
    }
    null = {
      source  = "hashicorp/null"
      version = "3.2.1"
    }
  }
}

provider "google" {
  project = var.project_id
}

provider "pexip" {
  address  = "https://${local.manager_hostname}.${local.domain}"
  username = var.infinity_username
  password = var.infinity_password
}