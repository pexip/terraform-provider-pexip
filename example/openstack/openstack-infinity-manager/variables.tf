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
  description = "Region where the manager will be created"
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
  default     = null
  description = "subnetwork gateway IP address"
}

variable "subnetwork_mask" {
  type        = string
  default     = null
  description = "subnetwork mask ('255.255.255.0' format)"
}

variable "floating_ip_pool" {
  type        = string
  description = "Name of the external network to allocate the floating IP from"
}

variable "image_id" {
  type        = string
  description = "ID of the Pexip Infinity VM image to use"
}

variable "flavor_name" {
  type        = string
  description = "The name of the desired flavor for the server"
}

variable "cpu_platform" {
  type        = string
  default     = "AMD Milan"
  description = "Minimum CPU platform for Infinity Manager VM"
}

variable "tags" {
  type        = list(string)
  default     = []
  description = "List of tags to apply to the Infinity Manager VM"
}

variable "dns_server" {
  type        = string
  description = "DNS server IP address for the Infinity Manager VM"
}

variable "ntp_server" {
  type        = string
  description = "NTP server IP address for the Infinity Manager VM"
}

variable "username" {
  type        = string
  default     = "admin"
  description = "Username for the Infinity Manager"
}

variable "password" {
  type        = string
  sensitive   = true
  description = "Password for the Infinity Manager"
}

variable "admin_password" {
  type        = string
  sensitive   = true
  description = "Admin password for the Infinity Manager"
}

variable "report_errors" {
  type        = bool
  default     = false
  description = "Enable error reporting for Infinity Manager"
}

variable "enable_analytics" {
  type        = bool
  default     = false
  description = "Enable analytics for Infinity Manager"
}

variable "contact_email_address" {
  type        = string
  description = "Contact email address for Infinity Manager notifications"
}

variable "management_ip_prefix" {
  type        = string
  description = "IP prefix for management access"
}

variable "domain" {
  type        = string
  description = "Domain name for the Infinity Manager"
}

variable "security_groups" {
  type        = list(string)
  default     = []
  description = "List of security groups to apply to the Infinity Manager VM"
}

variable "license_key" {
  type        = string
  sensitive   = true
  description = "Pexip Infinity license key"
}

variable "internode_ip_prefix" {
  type        = string
  description = "IP prefix for internode IPSec communication"
}
