/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "random_uuid4" "example" {
}

locals {
  uuid = "988d1247-7997-46e9-a89a-5a148b5c5f29"
}

# Create identity provider attributes for testing
resource "pexip_infinity_identity_provider_attribute" "attr1" {
  name        = "tf-test-displayName"
  description = "Test attribute for display name"
}

resource "pexip_infinity_identity_provider_attribute" "attr2" {
  name        = "tf-test-email"
  description = "Test attribute for email"
}

resource "pexip_infinity_identity_provider" "test" {
  # Required fields
  name                                = "tf-test Identity Provider full"
  uuid                                = local.uuid
  assertion_consumer_service_url      = "https://test.example.com/oidcconsumer/${local.uuid}"

  # Optional basic fields
  description                         = "Full test Identity Provider with all fields"
  idp_type                            = "oidc"

  # SAML specific fields
  sso_url                             = "https://idp.example.com/sso"
  idp_entity_id                       = "https://idp.example.com/entity"
  idp_public_key                      = "-----BEGIN CERTIFICATE-----\nMIIDljCCAn6gAwIBAgITEth/CdV4kCABggxrLpF5wyKqrTANBgkqhkiG9w0BAQsF\nADBbMQswCQYDVQQGEwJOTzENMAsGA1UECAwEb3NsbzENMAsGA1UEBwwEb3NsbzER\nMA8GA1UECgwIcGV4aXAgYXMxDDAKBgNVBAsMA2RldjENMAsGA1UEAwwEdGVzdDAe\nFw0yNjAyMjAxNDU0MjRaFw0yNzAyMjAxNDU0MjRaMFsxCzAJBgNVBAYTAk5PMQ0w\nCwYDVQQIDARvc2xvMQ0wCwYDVQQHDARvc2xvMREwDwYDVQQKDAhwZXhpcCBhczEM\nMAoGA1UECwwDZGV2MQ0wCwYDVQQDDAR0ZXN0MIIBIjANBgkqhkiG9w0BAQEFAAOC\nAQ8AMIIBCgKCAQEAtyO0f/mxvg08hH9twYX3ARuSQKnJS/2GnQW1akFtovYlmne8\niHxXC6OrTqYsjQhyJEhf7mZOmraV3X7Y8uSzdYtc0Mp6cZkqLHXHqA89alh8DP0r\n9HCQ3bk1bDz8IiT3N8F3sFCYeV2iK93tct6trkCVKEdKD3YtA4zNJ/nfWxXQQ0Kq\nYAn91Uq7gblcnN3KbAtNKO0sXf62vVE0d+3rumHzwe/QTl7y8bPjyIf6YffsFRk0\nlc21bq1oSiW2+u0wc4SC9LvV3cp6C7lui0+weYGgPSbk3TdpHCBfUJvxfkhZeswZ\nrHM76faf4N+pkbBSe+GIfclzFhPIazU5NmdKvQIDAQABo1MwUTAdBgNVHQ4EFgQU\nPo0hjDn7e6PH7vBDf26DxTGmyqEwHwYDVR0jBBgwFoAUPo0hjDn7e6PH7vBDf26D\nxTGmyqEwDwYDVR0TAQH/BAUwAwEB/zANBgkqhkiG9w0BAQsFAAOCAQEAk6DKXG+K\n3qAaAIFTqHx4KtE4MEBycHdNTsQb5dSqvkU/4TIHVaPSTZVzo4pP7a0GMLJOkxZ9\nuAwFrGoN8vWac6GZw3vvYCgAa0ncBuioSgf/PoKtXz9QoL0Qbzlpybghs+Rcje0i\nAufDIX1UAwyRjoJq9gSFx9IIjtBjtymTMFnsfjT9JDTEPdt4Pw2cywnqxwVF24nD\nC721z4cIXQsuyzXDav+adzKvuTqLJnHrtxE7SIF4jrDw8rBDDbJlUwWEqqmGHQf2\n6O1BuP0HM6b5OZnZeGgXclYuylh1Kht+4d4p9cCdyVbWwf18gPvfCfa3ZviKtkpU\nHU1n2Ch4Ufj8zQ==\n-----END CERTIFICATE-----"
  service_entity_id                   = "https://pexip.example.com/entity"
  service_public_key                  = "-----BEGIN CERTIFICATE-----\nMIIDljCCAn6gAwIBAgITEth/CdV4kCABggxrLpF5wyKqrTANBgkqhkiG9w0BAQsF\nADBbMQswCQYDVQQGEwJOTzENMAsGA1UECAwEb3NsbzENMAsGA1UEBwwEb3NsbzER\nMA8GA1UECgwIcGV4aXAgYXMxDDAKBgNVBAsMA2RldjENMAsGA1UEAwwEdGVzdDAe\nFw0yNjAyMjAxNDU0MjRaFw0yNzAyMjAxNDU0MjRaMFsxCzAJBgNVBAYTAk5PMQ0w\nCwYDVQQIDARvc2xvMQ0wCwYDVQQHDARvc2xvMREwDwYDVQQKDAhwZXhpcCBhczEM\nMAoGA1UECwwDZGV2MQ0wCwYDVQQDDAR0ZXN0MIIBIjANBgkqhkiG9w0BAQEFAAOC\nAQ8AMIIBCgKCAQEAtyO0f/mxvg08hH9twYX3ARuSQKnJS/2GnQW1akFtovYlmne8\niHxXC6OrTqYsjQhyJEhf7mZOmraV3X7Y8uSzdYtc0Mp6cZkqLHXHqA89alh8DP0r\n9HCQ3bk1bDz8IiT3N8F3sFCYeV2iK93tct6trkCVKEdKD3YtA4zNJ/nfWxXQQ0Kq\nYAn91Uq7gblcnN3KbAtNKO0sXf62vVE0d+3rumHzwe/QTl7y8bPjyIf6YffsFRk0\nlc21bq1oSiW2+u0wc4SC9LvV3cp6C7lui0+weYGgPSbk3TdpHCBfUJvxfkhZeswZ\nrHM76faf4N+pkbBSe+GIfclzFhPIazU5NmdKvQIDAQABo1MwUTAdBgNVHQ4EFgQU\nPo0hjDn7e6PH7vBDf26DxTGmyqEwHwYDVR0jBBgwFoAUPo0hjDn7e6PH7vBDf26D\nxTGmyqEwDwYDVR0TAQH/BAUwAwEB/zANBgkqhkiG9w0BAQsFAAOCAQEAk6DKXG+K\n3qAaAIFTqHx4KtE4MEBycHdNTsQb5dSqvkU/4TIHVaPSTZVzo4pP7a0GMLJOkxZ9\nuAwFrGoN8vWac6GZw3vvYCgAa0ncBuioSgf/PoKtXz9QoL0Qbzlpybghs+Rcje0i\nAufDIX1UAwyRjoJq9gSFx9IIjtBjtymTMFnsfjT9JDTEPdt4Pw2cywnqxwVF24nD\nC721z4cIXQsuyzXDav+adzKvuTqLJnHrtxE7SIF4jrDw8rBDDbJlUwWEqqmGHQf2\n6O1BuP0HM6b5OZnZeGgXclYuylh1Kht+4d4p9cCdyVbWwf18gPvfCfa3ZviKtkpU\nHU1n2Ch4Ufj8zQ==\n-----END CERTIFICATE-----"
  service_private_key                 = "-----BEGIN PRIVATE KEY-----\nMIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQC3I7R/+bG+DTyE\nf23BhfcBG5JAqclL/YadBbVqQW2i9iWad7yIfFcLo6tOpiyNCHIkSF/uZk6atpXd\nftjy5LN1i1zQynpxmSosdceoDz1qWHwM/Sv0cJDduTVsPPwiJPc3wXewUJh5XaIr\n3e1y3q2uQJUoR0oPdi0DjM0n+d9bFdBDQqpgCf3VSruBuVyc3cpsC00o7Sxd/ra9\nUTR37eu6YfPB79BOXvLxs+PIh/ph9+wVGTSVzbVurWhKJbb67TBzhIL0u9XdynoL\nuW6LT7B5gaA9JuTdN2kcIF9Qm/F+SFl6zBmsczvp9p/g36mRsFJ74Yh9yXMWE8hr\nNTk2Z0q9AgMBAAECggEAAJlqGMRxMD07qJT0WzPs/w/0lPdA4k1qk81jt8wH6tYA\nRGd8/qk+pFch1mV/50wMzUgFWOmfgTliibqzRP26Kgj+ZmayEq+EdFv6+N43hH42\nJ5xwwfnYCQjxCzc5HBXfySBHQIh4sk3srjBmgaFpPdgKNav788tts6knK5B8EICj\nIaQULnURk3ksXTWffHmOiULBM0YdZr8/WX5DEYW0H6uKytlYmMgiOBwPzLV+W56V\nx8vRNAR+8SDhvVZIFMMQNG3YVx5CtcfPFPPIyWJ/z9gmRjdHMVcEzRfMA27+yRBD\nPhVSaDa4zU8jLZJOYxmsoYtwN84sGLYt+e/liyiwgQKBgQDw2EAU1E6cOVVdPBv7\nboLZXeOZSaAMVhOwy8mLcOZlkehAIDDWzh9VOyxJTFAq86Pxe+ZI458Phmpvst81\nccfQgxKSITf9POZczYaU/YuTkdVFhc7D2M+/CYvjOJTzxX/Cr7bmTwHxUNvwjEZf\nSB3CYCNlr5NCk5OhH+m4ZTlufQKBgQDCqeKX+bZIQ2GJ7S4LB6ZGrqZe4kAgOWTs\nyWDKSIfgJuJ7mWUjspfzmV/HRLNBA4czpEfeZIWrN0mGCyDGCdcq8M/W1c977DAG\nKjT7AQS1omWyaukxIke4Q5rsXwQWH47lVmrghzYndZ7HYdWCee5KNlHuY1UlFJpV\nXpptOX7BQQKBgB6mJmGAMxDGawqWX6k7gwNOY2xaZjerrI3PITLRh0BdtQSUFCMB\n0BL0mMorH/iXUMGmVsPn05ISNPFJ/gW7GQZ7550ZBRsA6P/eV8YWDfEVmeJjCND7\nglR5Tof351yuQXfxDIF6hHDFLXgLIdl2P/NAcMC7+y15wku61+9TUl+xAoGAZYoK\nGS33KCCFm1Vtg/FciVgGjk3EF0r42w57/2fwADsoPkKYYBODcVyaei3attnpR0W0\n+0tB8jpnjpT1ZnexlcOBFlX24XQk5MJVWmyAkWIBXByqQKfZ80LIZ+10Czow5m26\nWB4PYGvZA7WDkoiZhHprKEcGHc5uZoNvV/P6q0ECgYEAgUnKnMQXzfmhSKghJ6t5\ntFl/o6nDfIuQdwV7lhYP9cY2UAU5eo3mz+bLr0NEstauD8NwqU/KJAxEEkle5xZp\nfuFQT6CAMfEwmDuX8IIFGW18eBZEu0IWNyTJB+Ueraz7pA73lGhhDDdmGccKYaPu\nVE1j4JF+iFBR/JrZgUEdbjM=\n-----END PRIVATE KEY-----"
  signature_algorithm                 = "http://www.w3.org/2001/04/xmldsig-more#rsa-sha384"
  digest_algorithm                    = "http://www.w3.org/2001/04/xmldsig-more#sha384"
  display_name_attribute_name         = "displayName"
  registration_alias_attribute_name   = "userPrincipalName"

  # Additional assertion consumer service URLs
  assertion_consumer_service_url2     = "https://test2.example.com/oidcconsumer/988d1247-7997-46e9-a89a-5a148b5c5f29"
  assertion_consumer_service_url3     = "https://test3.example.com/oidcconsumer/988d1247-7997-46e9-a89a-5a148b5c5f29"
  assertion_consumer_service_url4     = "https://test4.example.com/oidcconsumer/988d1247-7997-46e9-a89a-5a148b5c5f29"
  assertion_consumer_service_url5     = "https://test5.example.com/oidcconsumer/988d1247-7997-46e9-a89a-5a148b5c5f29"
  assertion_consumer_service_url6     = "https://test6.example.com/oidcconsumer/988d1247-7997-46e9-a89a-5a148b5c5f29"
  assertion_consumer_service_url7     = "https://test7.example.com/oidcconsumer/988d1247-7997-46e9-a89a-5a148b5c5f29"
  assertion_consumer_service_url8     = "https://test8.example.com/oidcconsumer/988d1247-7997-46e9-a89a-5a148b5c5f29"
  assertion_consumer_service_url9     = "https://test9.example.com/oidcconsumer/988d1247-7997-46e9-a89a-5a148b5c5f29"
  assertion_consumer_service_url10    = "https://test10.example.com/oidcconsumer/988d1247-7997-46e9-a89a-5a148b5c5f29"

  # Worker and popup settings
  worker_fqdn_acs_urls                = true
  disable_popup_flow                  = true

  # OIDC specific fields
  oidc_flow                           = "implicit"
  oidc_client_id                      = "test-client-id-12345"
  oidc_client_secret                  = "test-client-secret-67890"
  oidc_token_url                      = "https://idp.example.com/oauth2/token"
  oidc_user_info_url                  = "https://idp.example.com/oauth2/userinfo"
  oidc_jwks_url                       = "https://idp.example.com/.well-known/jwks.json"
  oidc_token_endpoint_auth_scheme     = "client_secret_basic"
  oidc_token_signature_scheme         = "hs256"
  oidc_display_name_claim_name        = "full_name"
  oidc_registration_alias_claim_name  = "preferred_username"
  oidc_additional_scopes              = "profile email phone address"
  oidc_france_connect_required_eidas_level = "eidas3"

  # Attributes
  attributes = [
    pexip_infinity_identity_provider_attribute.attr1.id,
    pexip_infinity_identity_provider_attribute.attr2.id,
  ]
}
