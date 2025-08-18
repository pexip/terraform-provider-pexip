/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

variable "auth_url" {
  type        = string
  default     = null
  description = "OpenStack Auth URL"
}

variable "tenant_name" {
  type        = string
  default     = null
  description = "OpenStack Tenant Name"
}

variable "project_name" {
  type        = string
  default     = null
  description = "OpenStack Project Name"
}

variable "region" {
  type        = string
  default     = null
  description = "The region of the OpenStack cloud to use"
}

variable "environment" {
  type        = string
  default     = "dev"
  description = "Environment for the deployment (e.g., dev, prod)"
}

variable "domain" {
  type        = string
  default     = "dev-pexip-network"
  description = "VM host domain"
}

variable "management_ip_prefix" {
  type        = string
  description = "IP prefix for management access"
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

variable "mgr_floating_ip_pool" {
  type        = string
  description = "Name of the external network to allocate the manager floating IP from"
}

variable "cnf_floating_ip_pool" {
  type        = string
  description = "Name of the external network to allocate the conferencing node floating IP from"
}

variable "internode_ip_prefix" {
  type        = string
  description = "IP prefix for internode IPSec communication"
}

variable "mgr_flavor_name" {
  type        = string
  default     = "4c4g"
  description = "The name of the desired flavor for the Infinity Manager server"
}

variable "cnf_flavor_name" {
  type        = string
  default     = "4c4g"
  description = "The name of the desired flavor for the Infinity Conferencing Node servers"
}

variable "private_network_id" {
  type        = string
  description = "ID of the private network for Infinity Manager and nodes"
}

variable "private_subnetwork_id_mgr" {
  type        = string
  description = "ID of the private subnetwork for Infinity Manager"
}

variable "private_subnetwork_id_cnf" {
  type        = string
  description = "ID of the private subnetwork for Infinity Conferencing Nodes"
}

variable "openstack_username" {
  type        = string
  description = "OpenStack username"
}

variable "openstack_password" {
  type        = string
  sensitive   = true
  description = "OpenStack password"
}

variable "openstack_tenant_name" {
  type        = string
  description = "OpenStack tenant name (project name)"
}

variable "openstack_domain_name" {
  type        = string
  description = "OpenStack domain name"
}