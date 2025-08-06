/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

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

variable "environment" {
  type        = string
  default     = "dev"
  description = "Environment for the deployment (e.g., dev, prod)"
}

variable "dns_zone_name" {
  type        = string
  default     = "dev-pexip-network"
  description = "name of GCP DNS zone"
}

variable "vm_image_manager_name" {
  type        = string
  default     = "pexip-mgr-v37"
  description = "Pexip Infinity VM image to use"
}

variable "vm_image_node_name" {
  type        = string
  default     = "pexip-cnf-v37"
  description = "VM image to use"
}

variable "infinity_license_key" {
  type        = string
  sensitive   = true
  description = "License key for Infinity"
}

variable "infinity_address" {
  type        = string
  default     = "https://manager.example.com"
  description = "Address of the Infinity Manager"
}

variable "infinity_ip_address" {
  type        = string
  default     = null
  description = "IP address of the Infinity Manager"
}

variable "infinity_manager_machine_type" {
  type        = string
  default     = "n2d-standard-16"
  description = "Machine type for Infinity Manager VM"
}

variable "infinity_node_machine_type" {
  type        = string
  default     = "n2d-standard-16"
  description = "Machine type for Infinity nodes"
}

variable "infinity_manager_cpu_platform" {
  type        = string
  default     = "AMD Milan"
  description = "Minimum CPU platform for Infinity Manager VM"
}

variable "infinity_node_cpu_platform" {
  type        = string
  default     = "AMD Milan"
  description = "Minimum CPU platform for Infinity Manager VM"
}

variable "infinity_username" {
  type        = string
  default     = "admin"
  description = "Username for Infinity Manager"
}

variable "infinity_password" {
  type        = string
  sensitive   = true
  description = "Password for Infinity Manager"
}

variable "infinity_primary_dns_server" {
  type        = string
  default     = "8.8.8.8"
  description = "Primary DNS server for Infinity Manager"
}

variable "infinity_ntp_server" {
  type        = string
  default     = "pool.ntp.org"
  description = "NTP server for Infinity Manager"
}

variable "infinity_report_errors" {
  type        = bool
  default     = false
  description = "Enable error reporting for Infinity Manager"
}

variable "infinity_enable_analytics" {
  type        = bool
  default     = false
  description = "Enable analytics for Infinity Manager"
}

variable "infinity_contact_email_address" {
  type        = string
  description = "Contact email address for Infinity Manager"
}

variable "infinity_node_count" {
  type        = number
  default     = 1
  description = "Number of Infinity nodes to deploy"
}
