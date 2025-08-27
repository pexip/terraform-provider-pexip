/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

variable "environment" {
  type        = string
  description = "Environment for the deployment (e.g., dev, prod)"
}

variable "region" {
  type        = string
  default     = null
  description = "Region where the conferencing node will be created"
}

variable "private_network_name" {
  type        = string
  description = "Name of the Openstack private network to use"
}

variable "private_subnetwork_name" {
  type        = string
  description = "Name of the Openstack private subnetwork to use"
}

variable "floating_ip_pool" {
  type        = string
  description = "Name of the external network to allocate the floating IP from"
}

variable "gateway" {
  type        = string
  default     = null
  description = "subnetwork gateway IP address"
}

variable "subnetwork_mask" {
  type        = string
  default     = null
  description = "subnetwork mask ('255.255.255.0' format)"
}

variable "image_id" {
  type        = string
  description = "ID of the Pexip Infinity VM image to use"
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

variable "domain" {
  type        = string
  description = "Domain name for the Infinity Node"
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

variable "management_ip_prefix" {
  type        = string
  description = "IP prefix for management access"
}

variable "flavor_name" {
  type        = string
  description = "The name of the desired flavor for the server"
}

variable "security_groups" {
  type        = list(string)
  default     = []
  description = "List of security groups to apply to the Infinity Node VM"
}

variable "internode_ip_prefix" {
  type        = string
  description = "IP prefix for internode IPSec communication"
}

variable "mgr_public_ip" {
  type        = string
  description = "Public IP address of the Infinity Manager"
}

variable "web_username" {
  type        = string
  default     = "admin"
  description = "Username for Infinity Web interface"
}

variable "web_password" {
  type        = string
  sensitive   = true
  description = "Password for Infinity Web interface user"
}

variable "tls_certificate" {
  type        = string
  default     = ""
  description = "TLS certificate for the Infinity Node"
}