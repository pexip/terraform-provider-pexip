variable "project_id" {
  type        = string
  default     = "trusted-lemur-145810"
  description = "GCP Project ID"
}

variable "location" {
  type        = string
  default     = "europe-west4"
  description = "location where resources will be created"
}

variable "hostname" {
  type       = string
  default    = "infinity-manager"
  description = "Hostname of the VM instance"
}

variable "dns_zone_name" {
  type        = string
  default     = "dev-pexip-network"
  description = "name of GCP DNS zone"
}

variable "vm_image_project" {
    type        = string
    default     = "trusted-lemur-145810"
    description = "Project ID where the VM image is located"
}

variable "vm_image_name" {
  type        = string
  default     = "pexip-infiniy-manager-v38"
  description = "VM image to use"
}