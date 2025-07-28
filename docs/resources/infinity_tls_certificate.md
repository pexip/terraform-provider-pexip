---
page_title: "pexip_infinity_tls_certificate Resource - terraform-provider-pexip"
subcategory: ""
description: |-
  Manages a Pexip Infinity TLS certificate configuration.
---

# pexip_infinity_tls_certificate (Resource)

Manages a Pexip Infinity TLS certificate configuration. TLS certificates are used to secure communications between Pexip Infinity components and external systems, including web browser connections, SIP over TLS, and API access. Certificates can be deployed to specific nodes or made available for system-wide use.

## Example Usage

### Basic TLS Certificate

```terraform
resource "pexip_infinity_tls_certificate" "example" {
  certificate = file("${path.module}/certificates/server.crt")
  private_key = file("${path.module}/certificates/server.key")
}
```

### TLS Certificate with Encrypted Private Key

```terraform
resource "pexip_infinity_tls_certificate" "encrypted_cert" {
  certificate              = file("${path.module}/certificates/encrypted-server.crt")
  private_key              = file("${path.module}/certificates/encrypted-server.key")
  private_key_passphrase   = var.private_key_passphrase
  parameters               = "SSL/TLS certificate for web interface"
}
```

### TLS Certificate Deployed to Specific Nodes

```terraform
# Get node resource URIs
data "pexip_infinity_node" "conferencing_nodes" {
  count = length(var.conferencing_node_ids)
  id    = var.conferencing_node_ids[count.index]
}

resource "pexip_infinity_tls_certificate" "node_specific" {
  certificate = file("${path.module}/certificates/conferencing.crt")
  private_key = file("${path.module}/certificates/conferencing.key")
  parameters  = "TLS certificate for conferencing nodes"
  nodes = [
    for node in data.pexip_infinity_node.conferencing_nodes : node.id
  ]
}
```

### Wildcard Certificate for Multiple Services

```terraform
resource "pexip_infinity_tls_certificate" "wildcard" {
  certificate = file("${path.module}/certificates/wildcard.company.com.crt")
  private_key = file("${path.module}/certificates/wildcard.company.com.key")
  parameters  = "Wildcard certificate for *.company.com"
}
```

### Let's Encrypt Certificate with Automated Renewal

```terraform
# Using external certificate generation (example with ACME provider)
resource "acme_certificate" "letsencrypt" {
  account_key_pem = acme_private_key.account.private_key_pem
  common_name     = "pexip.company.com"
  subject_alternative_names = [
    "pexip-web.company.com",
    "pexip-api.company.com"
  ]

  dns_challenge {
    provider = "route53"
  }
}

resource "pexip_infinity_tls_certificate" "letsencrypt" {
  certificate = "${acme_certificate.letsencrypt.certificate_pem}${acme_certificate.letsencrypt.issuer_pem}"
  private_key = acme_certificate.letsencrypt.private_key_pem
  parameters  = "Let's Encrypt certificate with automated renewal"
}
```

### Enterprise Certificate with Full Chain

```terraform
locals {
  certificate_chain = "${file("${path.module}/certificates/server.crt")}${file("${path.module}/certificates/intermediate.crt")}${file("${path.module}/certificates/root.crt")}"
}

resource "pexip_infinity_tls_certificate" "enterprise" {
  certificate = local.certificate_chain
  private_key = file("${path.module}/certificates/server.key")
  parameters  = "Enterprise certificate with full trust chain"
}
```

### Multiple Certificates for Different Purposes

```terraform
# Web interface certificate
resource "pexip_infinity_tls_certificate" "web_interface" {
  certificate = file("${path.module}/certificates/web.crt")
  private_key = file("${path.module}/certificates/web.key")
  parameters  = "Certificate for web interface (HTTPS)"
}

# API certificate
resource "pexip_infinity_tls_certificate" "api" {
  certificate = file("${path.module}/certificates/api.crt")
  private_key = file("${path.module}/certificates/api.key")
  parameters  = "Certificate for API access"
}

# SIP TLS certificate
resource "pexip_infinity_tls_certificate" "sip_tls" {
  certificate = file("${path.module}/certificates/sip.crt")
  private_key = file("${path.module}/certificates/sip.key")
  parameters  = "Certificate for SIP over TLS"
}
```

### Certificate with Variables for Sensitive Data

```terraform
variable "tls_certificate" {
  description = "TLS certificate in PEM format"
  type        = string
  sensitive   = true
}

variable "tls_private_key" {
  description = "TLS private key in PEM format"
  type        = string
  sensitive   = true
}

variable "private_key_passphrase" {
  description = "Passphrase for encrypted private key"
  type        = string
  sensitive   = true
  default     = ""
}

resource "pexip_infinity_tls_certificate" "from_variables" {
  certificate            = var.tls_certificate
  private_key            = var.tls_private_key
  private_key_passphrase = var.private_key_passphrase != "" ? var.private_key_passphrase : null
  parameters             = "Certificate deployed from Terraform variables"
}
```

## Schema

### Required

- `certificate` (String, Sensitive) - The PEM-encoded certificate. This can include the full certificate chain (server certificate + intermediate certificates + root certificate).
- `private_key` (String, Sensitive) - The PEM-encoded private key corresponding to the certificate.

### Optional

- `private_key_passphrase` (String, Sensitive) - The passphrase for the private key if it is encrypted. Maximum length: 100 characters.
- `parameters` (String) - Additional parameters or description for the certificate. Maximum length: 1000 characters.
- `nodes` (List of String) - List of node resource URIs where this certificate should be deployed. If not specified, the certificate is available for system-wide use.

### Read-Only

- `id` (String) - Resource URI for the TLS certificate in Infinity.
- `resource_id` (Number) - The resource integer identifier for the TLS certificate in Infinity.
- `start_date` (String) - The start date of the certificate validity period.
- `end_date` (String) - The end date of the certificate validity period.
- `subject_name` (String) - The subject name from the certificate.
- `subject_hash` (String) - The subject hash from the certificate.
- `subject_alt_names` (String) - The subject alternative names from the certificate.
- `raw_subject` (String) - The raw subject data from the certificate.
- `issuer_name` (String) - The issuer name from the certificate.
- `issuer_hash` (String) - The issuer hash from the certificate.
- `raw_issuer` (String) - The raw issuer data from the certificate.
- `serial_no` (String) - The serial number from the certificate.
- `key_id` (String) - The key identifier from the certificate.
- `issuer_key_id` (String) - The issuer key identifier from the certificate.
- `text` (String) - The text representation of the certificate.

## Import

Import is supported using the following syntax:

```shell
terraform import pexip_infinity_tls_certificate.example 123
```

Where `123` is the numeric resource ID of the TLS certificate.

**Note**: During import, the private key information will not be available from the API and must be provided manually after import.

## Usage Notes

### Certificate Format

- **PEM Format**: Certificates and private keys must be in PEM (Base64 encoded) format
- **Certificate Chain**: Include the full certificate chain (server cert + intermediates + root) for proper validation
- **Line Endings**: Ensure proper line endings (LF) in certificate and key files
- **Character Encoding**: Use UTF-8 encoding for certificate files

### Private Key Security

- **Encryption**: Consider using encrypted private keys for additional security
- **Key Storage**: Store private keys securely and limit access
- **Key Rotation**: Implement regular key rotation policies
- **Backup**: Securely backup private keys and certificates

### Certificate Validation

- **Expiration**: Monitor certificate expiration dates and implement renewal processes
- **Common Name**: Ensure the certificate common name matches the intended hostname
- **Subject Alternative Names**: Include all required hostnames in SAN fields
- **Chain Validation**: Verify the complete certificate chain is valid

### Node Deployment

- **System-wide**: Leave nodes empty for system-wide certificate availability
- **Node-specific**: Specify node URIs for targeted deployment
- **Load Balancing**: Deploy certificates to all nodes behind load balancers
- **High Availability**: Ensure certificates are deployed to all redundant nodes

### Certificate Types

- **Self-signed**: Suitable for development and internal testing
- **Internal CA**: Use for internal environments with corporate PKI
- **Public CA**: Use for production environments accessible from the internet
- **Wildcard**: Use for multiple subdomains under the same domain

### Performance Considerations

- **Key Size**: Balance security and performance (2048-bit RSA or 256-bit ECC)
- **Cipher Suites**: Configure appropriate cipher suites for security and performance
- **Session Resumption**: Enable TLS session resumption for improved performance
- **OCSP Stapling**: Configure OCSP stapling to reduce certificate validation overhead

### Security Best Practices

- **Key Management**: Use hardware security modules (HSMs) for production environments
- **Certificate Transparency**: Monitor Certificate Transparency logs for unauthorized certificates
- **Revocation**: Implement proper certificate revocation processes
- **Compliance**: Ensure certificates meet organizational and regulatory requirements

## Troubleshooting

### Common Issues

**Certificate Upload Fails**
- Verify the certificate is in valid PEM format
- Ensure the private key matches the certificate
- Check that the certificate chain is complete and in the correct order
- Verify the private key passphrase is correct if the key is encrypted

**Invalid Certificate Format**
- Ensure certificates start with `-----BEGIN CERTIFICATE-----` and end with `-----END CERTIFICATE-----`
- Verify private keys start with `-----BEGIN PRIVATE KEY-----` or `-----BEGIN RSA PRIVATE KEY-----`
- Check for proper line endings (Unix LF format)
- Remove any extra whitespace or characters

**Private Key Mismatch**
- Verify the private key corresponds to the uploaded certificate
- Check that the key was not corrupted during transfer
- Ensure the correct passphrase is provided for encrypted keys
- Validate the key format and encoding

**Certificate Chain Issues**
- Include intermediate certificates in the proper order
- Verify the certificate chain is complete from server to root CA
- Check that all certificates in the chain are valid and not expired
- Ensure proper certificate chain validation

**Encrypted Private Key Problems**
- Verify the passphrase is correct and properly encoded
- Check that the private key encryption format is supported
- Ensure the passphrase meets length requirements
- Consider using unencrypted keys for automated deployments

**Node Deployment Failures**
- Verify the specified node URIs exist and are accessible
- Check that nodes are online and responsive
- Ensure proper permissions for certificate deployment
- Monitor node-specific logs for deployment errors

**Certificate Validation Errors**
- Check certificate expiration dates
- Verify the certificate is issued by a trusted CA
- Ensure the certificate common name or SAN matches the hostname
- Validate certificate purpose and key usage extensions

**TLS Handshake Failures**
- Verify certificate and private key are properly paired
- Check TLS protocol version compatibility
- Ensure cipher suites are properly configured
- Monitor TLS handshake logs for specific error details

**Import Issues**
- Use the numeric resource ID, not the certificate common name
- Verify the certificate exists in the Infinity cluster
- Note that private key information is not available during import
- Manually configure the private key after successful import

**Certificate Renewal Problems**
- Monitor certificate expiration dates proactively
- Implement automated renewal processes where possible
- Test certificate renewal procedures regularly
- Ensure proper coordination between renewal and deployment

**Performance Issues**
- Consider using ECC certificates for better performance
- Monitor TLS handshake times and session establishment
- Check for proper session resumption configuration
- Optimize cipher suite selection for performance

**Compatibility Issues**
- Verify certificate compatibility with client systems
- Check TLS protocol version support
- Ensure cipher suite compatibility across all clients
- Test certificate validation with various client types

**Security Audit Failures**
- Ensure certificates meet organizational security standards
- Verify proper key sizes and algorithms are used
- Check certificate transparency logging compliance
- Implement proper certificate lifecycle management