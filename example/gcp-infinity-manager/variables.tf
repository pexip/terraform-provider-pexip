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
  description = "Machine type for Infinity Manager VM"
}

variable "cpu_platform" {
  type        = string
  default     = "AMD Milan"
  description = "Minimum CPU platform for Infinity Manager VM"
}

variable "service_account_email" {
  type        = string
  description = "Email of the service account to use for the Infinity Manager VM"
}

variable "dns_zone_name" {
  type        = string
  description = "name of GCP DNS zone"
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

