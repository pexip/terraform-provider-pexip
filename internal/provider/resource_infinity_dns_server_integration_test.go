//go:build integration

package provider

import (
	"crypto/tls"
	"github.com/stretchr/testify/require"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/pexip/go-infinity-sdk/v38"
)

func TestInfinityDNSServerIntegration(t *testing.T) {
	_ = os.Setenv("TF_ACC", "1")

	client, err := infinity.New(
		infinity.WithBaseURL("https://dev-manager.dev.pexip.network"),
		infinity.WithBasicAuth("admin", "admin"),
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

	testInfinityDNSServer(t, client)
}
