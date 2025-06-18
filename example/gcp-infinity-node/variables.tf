variable "project_id" {
  type        = string
  description = "GCP Project ID"
}

variable "environment" {
  type        = string
  description = "Environment for the deployment (e.g., dev, prod)"
}

variable "location" {
  type        = string
  description = "location where resources will be created"
}

variable "network_id" {
  type        = string
  description = "ID of the GCP network to use"
}

variable "subnetwork_id" {
  type        = string
  description = "ID of the GCP subnetwork to use"
}

variable "vm_image_name" {
  type        = string
  description = "Pexip Infinity VM image to use"
}

variable "vm_image_project" {
  type        = string
  default     = "vc-operations"
  description = "Project ID where the VM image is located"
}

variable "machine_type" {
  type        = string
  default     = "n2d-standard-16"
  description = "Machine type for Infinity Node VM"
}

variable "cpu_platform" {
  type        = string
  default     = "AMD Milan"
  description = "Minimum CPU platform for Infinity Node VM"
}

variable "service_account_email" {
  type        = string
  description = "Email of the service account to use for the Infinity NOde VM"
}

variable "dns_zone_name" {
  type        = string
  description = "name of GCP DNS zone"
}

variable "tags" {
  type        = list(string)
  default     = []
  description = "List of tags to apply to the Infinity Node VM"
}

variable "index" {
  type        = number
  default     = 1
  description = "Index of the Infinity Node"
}
