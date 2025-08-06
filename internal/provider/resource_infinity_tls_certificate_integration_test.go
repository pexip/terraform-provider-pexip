//go:build integration

package provider

import (
	"crypto/tls"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/pexip/terraform-provider-pexip/internal/test"

	"github.com/pexip/go-infinity-sdk/v38"
)

func TestInfinityTLSCertificateIntegration(t *testing.T) {
	_ = os.Setenv("TF_ACC", "1")

	tlsCertificate := test.LoadTestFile(t, "tls_certificate.pem")
	tlsPrivateKey := test.LoadTestFile(t, "tls_private_key.pem")

	client, err := infinity.New(
		infinity.WithBaseURL(test.INFINITY_BASE_URL),
		infinity.WithBasicAuth(test.INFINITY_USERNAME, test.INFINITY_PASSWORD),
		infinity.WithMaxRetries(2),
		infinity.WithTransport(&http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, // We need this because default certificate is not trusted
				MinVersion:         tls.VersionTLS12,
			},
			MaxIdleConns:        30,
			MaxIdleConnsPerHost: 5,
			IdleConnTimeout:     60 * time.Second,
		}),
	)
	require.NoError(t, err)

	testInfinityTLSCertificate(t, client, tlsPrivateKey, tlsCertificate)
}
