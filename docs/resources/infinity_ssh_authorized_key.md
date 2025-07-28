# infinity_ssh_authorized_key

Manages an SSH authorized key configuration with the Infinity service. SSH authorized keys allow secure key-based authentication for SSH access to Pexip Infinity nodes.

## Example Usage

```hcl
resource "pexip_infinity_ssh_authorized_key" "example" {
  keytype = "ssh-rsa"
  key     = "AAAAB3NzaC1yc2EAAAADAQABAAABgQC7vbqajDhA..."
  comment = "admin@example.com"
  nodes = [
    "/configuration/v1/worker_vm/1/",
    "/configuration/v1/management_vm/1/"
  ]
}
```

## Argument Reference

The following arguments are supported:

* `keytype` - (Required) The SSH key type. Valid choices: `ssh-rsa`, `ssh-dss`, `ssh-ed25519`, `ecdsa-sha2-nistp256`, `ecdsa-sha2-nistp384`, `ecdsa-sha2-nistp521`.
* `key` - (Required) The SSH public key content (base64 encoded key data).
* `comment` - (Optional) A comment for the SSH key. Maximum length: 250 characters.
* `nodes` - (Optional) List of node resource URIs where this SSH key is authorized.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Resource URI for the SSH authorized key in Infinity.
* `resource_id` - The resource integer identifier for the SSH authorized key in Infinity.

## Import

SSH authorized keys can be imported using their resource ID:

```bash
terraform import pexip_infinity_ssh_authorized_key.example 123
```

## Key Format

The SSH public key should be provided in the standard OpenSSH format without the key type prefix or comment suffix. For example:
- **Correct**: `AAAAB3NzaC1yc2EAAAADAQABAAABgQC7vbqajDhA...`
- **Incorrect**: `ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQC7vbqajDhA... user@hostname`

The key type and comment are specified separately in the `keytype` and `comment` fields respectively.