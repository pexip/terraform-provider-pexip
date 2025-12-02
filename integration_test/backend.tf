terraform {
  backend "gcs" {
    bucket = "px-eng-terraform-provider-pexip-tf-state"
    prefix = "integration-test"
  }
}
