---
page_title: "pexip_infinity_worker_vm Resource - terraform-provider-pexip"
subcategory: ""
description: |-
  Manages a Pexip Infinity worker VM configuration.
---

# pexip_infinity_worker_vm (Resource)

Manages a Pexip Infinity worker VM configuration. This resource configures the worker VM settings for a Pexip Infinity node in the cluster.

## Example Usage

### Basic Usage

```terraform
resource "pexip_infinity_worker_vm" "worker" {
  name            = "worker-vm-01"
  hostname        = "worker-vm-01"
  domain          = "company.com"
  address         = "10.0.1.20"
  netmask         = "255.255.255.0"
  gateway         = "10.0.1.1"
  system_location = "Main Location"
}
```

### Full Configuration

```terraform
resource "pexip_infinity_worker_vm" "worker" {
  name                  = "worker-vm-01"
  hostname              = "worker-vm-01"
  domain                = "company.com" 
  address               = "10.0.1.20"
  netmask               = "255.255.255.0"
  gateway               = "10.0.1.1"
  system_location       = "Main Location"
  
  # VM specifications
  vm_cpu_count          = 8
  vm_system_memory      = 8192
  
  # Node configuration
  node_type             = "conferencing"
  transcoding           = true
  maintenance_mode      = false
  maintenance_mode_reason = ""
  
  # Network configuration
  ipv6_address          = "2001:db8::20"
  ipv6_gateway          = "2001:db8::1"
  secondary_address     = "10.0.2.20"
  secondary_netmask     = "255.255.255.0"
  static_nat_address    = "203.0.113.20"
  
  # Security
  password              = var.vm_password
  enable_ssh            = "on"
  
  # SNMP configuration
  snmp_mode             = "standard"
  snmp_community        = "private"
  snmp_system_contact   = "admin@company.com"
  snmp_system_location  = "Data Center 1"
}
```

### Multiple Worker VMs

```terraform
resource "pexip_infinity_worker_vm" "workers" {
  count           = 3
  name            = "worker-vm-${count.index + 1}"
  hostname        = "worker-vm-${count.index + 1}"
  domain          = "company.com"
  address         = "10.0.1.${20 + count.index}"
  netmask         = "255.255.255.0"
  gateway         = "10.0.1.1"
  system_location = "Main Location"
}
```

## Schema

### Required

- `address` (String) - The IPv4 address of the worker VM.
- `domain` (String) - The domain of the worker VM. Maximum length: 192 characters.
- `gateway` (String) - The gateway address for the worker VM.
- `hostname` (String) - The hostname for this Conferencing Node. Each Conferencing Node must have a unique DNS hostname. Maximum length: 63 characters.
- `name` (String) - The name used to refer to this Conferencing Node. Each Conferencing Node must have a unique name. Maximum length: 32 characters.
- `netmask` (String) - The IPv4 network mask for this Conferencing Node.
- `system_location` (String) - The system location for this Conferencing Node. A system location should not contain a mixture of Proxying Edge Nodes and Transcoding Conferencing Nodes.

### Optional

- `alternative_fqdn` (String) - An identity for this Conferencing Node, used in signaling SIP TLS Contact addresses. Maximum length: 255 characters. Defaults to `""`.
- `cloud_bursting` (Boolean) - Defines whether this Conference Node is a cloud bursting node. Defaults to `false`.
- `deployment_type` (String) - The means by which this Conferencing Node will be deployed. Defaults to `"MANUAL-PROVISION-ONLY"`.
- `description` (String) - A description of the Conferencing Node. Maximum length: 250 characters. Defaults to `""`.
- `enable_distributed_database` (Boolean) - This should usually be True for all nodes which are expected to be 'always on', and False for nodes which are expected to only be powered on some of the time (e.g. cloud bursting nodes that are likely to only be operational during peak times). Avoid frequent toggling of this setting. Defaults to `true`.
- `enable_ssh` (String) - Allows an administrator to log in to this node over SSH. Valid values are: `global`, `off`, `on`. Defaults to `global`.
- `ipv6_address` (String) - The IPv6 address of the conferencing node. Maximum length: 250 characters.
- `ipv6_gateway` (String) - The IPv6 gateway for the conferencing node. Maximum length: 250 characters.
- `maintenance_mode` (Boolean) - Whether the worker VM is in maintenance mode. Defaults to `false`.
- `maintenance_mode_reason` (String) - The reason for maintenance mode. Maximum length: 250 characters. Defaults to `""`.
- `managed` (Boolean) - Whether the conferencing node is managed by the Infinity service. Defaults to `false`.
- `media_priority_weight` (Number) - The relative priority of this node, used when determining the order of nodes to which Pexip Infinity will attempt to send media. A higher number represents a higher priority; the default is 0, i.e. the lowest priority. Defaults to `0`.
- `node_type` (String) - The role of this Conferencing Node. Valid choices: `conferencing`, `proxying`. Defaults to `conferencing`.
- `password` (String, Sensitive) - The password to be used when logging in to this Conferencing Node over SSH. The username will always be admin. Maximum length: 250 characters. Defaults to `""`.
- `secondary_address` (String) - The optional secondary interface IPv4 address for this Conferencing Node.
- `secondary_netmask` (String) - The optional secondary interface IPv4 netmask for this Conferencing Node.
- `service_manager` (Boolean) - Handle Service Manager. Defaults to `true`.
- `service_policy` (Boolean) - Handle Service Policy. Defaults to `true`.
- `signalling` (Boolean) - Handle signalling. Defaults to `true`.
- `snmp_authentication_password` (String, Sensitive) - The password used for SNMPv3 authentication. Minimum length: 8 characters. Maximum length: 100 characters. Defaults to `""`.
- `snmp_community` (String, Sensitive) - The SNMP group to which this virtual machine belongs. Maximum length: 16 characters. Defaults to `"public"`.
- `snmp_mode` (String) - The SNMP mode. Valid values: `disabled`, `standard`, `authpriv`. Defaults to `"disabled"`.
- `snmp_privacy_password` (String, Sensitive) - The password used for SNMPv3 privacy. Minimum length: 8 characters. Maximum length: 100 characters. Defaults to `""`.
- `snmp_system_contact` (String) - The SNMP contact for this Conferencing Node. Maximum length: 70 characters. Defaults to `"admin@domain.com"`.
- `snmp_system_location` (String) - The SNMP location for this Conferencing Node. Maximum length: 70 characters. Defaults to `"Virtual machine"`.
- `snmp_username` (String) - The username used to authenticate SNMPv3 requests. Maximum length: 100 characters. Defaults to `""`.
- `ssh_authorized_keys` (List of String) - The selected authorized keys.
- `ssh_authorized_keys_use_cloud` (Boolean) - Allows use of SSH keys configured in the cloud service. Defaults to `true`.
- `static_nat_address` (String) - The public IPv4 address used by the Conferencing Node when it is located behind a NAT device. Note that if you are using NAT, you must also configure your NAT device to route the Conferencing Node's IPv4 static NAT address to its IPv4 address.
- `static_routes` (List of String) - Additional configuration to permit routing of traffic to networks not accessible through the configured default gateway.
- `tls_certificate` (String) - The TLS certificate to use on this node.
- `transcoding` (Boolean) - This determines the Conferencing Node's role. When transcoding is enabled, this node can handle all the media processing, protocol interworking, mixing and so on that is required in hosting Pexip Infinity calls and conferences. When transcoding is disabled, it becomes a Proxying Edge Node that can only handle the media and signaling connections with an endpoint or external device, and it then forwards the device's media on to a node that does have transcoding capabilities. Defaults to `true`.
- `vm_cpu_count` (Number) - Enter the number of virtual CPUs to assign to this Conferencing Node. We do not recommend that you assign more virtual CPUs than there are physical cores on a single processor on the host server (unless you have enabled NUMA affinity). For example, if the host server has 2 processors each with 12 physical cores, we recommend that you assign no more than 12 virtual CPUs. Range: 2 to 128. Defaults to `4`.
- `vm_system_memory` (Number) - The amount of RAM (in megabytes) to assign to this Conferencing Node. Range: 2000 to 64000. Defaults to `4096`.

### Read-Only

- `id` (String) - Resource URI for the worker VM in Infinity.
- `resource_id` (Number) - The resource integer identifier for the worker VM in Infinity.

## Import

Import is supported using the following syntax:

```shell
terraform import pexip_infinity_worker_vm.example 123
```

Where `123` is the numeric resource ID of the worker VM.

## Usage Notes

### VM Configuration
- The worker VM must have a unique name and hostname within the Infinity cluster
- CPU and memory settings depend on expected workload and available host resources
- Network configuration must be valid for your infrastructure

### Node Types
- `conferencing`: Full media processing and transcoding capabilities
- `proxying`: Edge node that forwards media to transcoding nodes

### Maintenance Mode
- Use maintenance mode to temporarily remove a node from service
- Always provide a reason when enabling maintenance mode
- Disable maintenance mode when the node is ready to serve traffic again

### SNMP Configuration
- SNMP can be disabled, use standard community strings, or use v3 authentication
- Configure appropriate contact and location information for monitoring

### Import Considerations
- Use the numeric resource ID (not the name) for import
- Ensure the worker VM exists in the Infinity cluster before importing
- Verify authentication credentials have access to the resource

## Troubleshooting

### Common Issues

**Worker VM Creation Fails**
- Verify the IP address is available and not in use
- Ensure the system location exists in the Infinity configuration
- Check that hostname and name are unique within the cluster

**Network Configuration Issues**
- Verify IP addresses, netmasks, and gateways are valid for your network
- For NAT configurations, ensure the static NAT address is properly configured
- IPv6 addresses must be valid if specified

**Performance Issues**
- Adjust CPU and memory allocations based on actual usage
- Monitor media priority weight to ensure proper load distribution
- Consider transcoding vs. proxying roles based on deployment needs

**Import Fails**
- Ensure you're using the numeric resource ID, not the name
- Verify the worker VM exists in the Infinity cluster
- Check provider authentication credentials