terraform {
  required_version = ">= 1.0"
  required_providers {
    pexip = {
      source  = "pexip.com/pexip/pexip"
      version = ">= 1.0"
    }
    google = {
      source  = "hashicorp/google"
      version = ">= 4.0"
    }
    random = {
      source  = "hashicorp/random"
      version = ">= 3.0"
    }
    null = {
      source  = "hashicorp/null"
      version = ">= 3.0"
    }
  }
}