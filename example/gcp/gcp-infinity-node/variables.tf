/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

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

variable "private_network_id" {
  type        = string
  description = "ID of the GCP private network to use"
}

variable "private_subnetwork_id" {
  type        = string
  description = "ID of the GCP private subnetwork to use"
}

variable "gateway" {
  type        = string
  description = "subnetwork gateway IP address"
}

variable "subnetwork_mask" {
  type        = string
  description = "subnetwork mask ('255.255.255.0' format)"
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
  description = "Email of the service account to use for the Infinity Node VM"
}

variable "dns_zone_name" {
  type        = string
  description = "name of GCP DNS zone"
}

variable "password" {
  type        = string
  sensitive   = true
  description = "Password for the Infinity Node"
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

variable "node_type" {
  type        = string
  description = "Type of Infinity Node (e.g., CONFERENCING, PROXYING)"
  validation {
    condition     = contains(["CONFERENCING", "PROXYING"], var.node_type)
    error_message = "Valid values for node_type are CONFERENCING or PROXYING."
  }
}

variable "system_location" {
  type        = string
  description = "Location of the Infinity Node system"
}

variable "maintenance_mode" {
  type        = bool
  default     = false
  description = "Enable maintenance mode for the Infinity Node"
}

variable "maintenance_mode_reason" {
  type        = string
  default     = ""
  description = "Reason for enabling maintenance mode"
}

variable "transcoding" {
  type        = bool
  default     = true
  description = "Enable transcoding for the Infinity Node"
}

variable "vm_cpu_count" {
  type        = number
  default     = 16
  description = "Number of vCPUs for the Infinity Node VM"
}

variable "vm_system_memory" {
  type        = number
  default     = 4096
  description = "Amount of system memory (in MB) for the Infinity Node VM"
}

variable "tls_certificate" {
  type        = string
  default     = ""
  description = "TLS certificate for the Infinity Node"
}
