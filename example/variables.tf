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

variable "vm_image_project" {
  type        = string
  default     = "vc-operations"
  description = "Project ID where the VM image is located"
}

variable "vm_image_manager_name" {
  type        = string
  default     = "pexip-mgr-v37"
  description = "VM image to use"
}

variable "vm_image_node_name" {
  type        = string
  default     = "pexip-cnf-v37"
  description = "VM image to use"
}

variable "infinity_address" {
  type        = string
  default     = "https://manager.example.com"
  description = "Address of the Infinity Manager"
}

variable "infinity_hostname" {
  type        = string
  default     = "manager-01"
  description = "Hostname for Infinity Manager"
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

variable "infinity_secondary_dns_server" {
  type        = string
  default     = "8.8.4.4"
  description = "Secondary DNS server for Infinity Manager"
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
