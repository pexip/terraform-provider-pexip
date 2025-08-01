//go:build integration

package provider

import (
	"crypto/tls"
	"github.com/pexip/terraform-provider-pexip/internal/test"
	"github.com/stretchr/testify/require"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/pexip/go-infinity-sdk/v38"
)

func TestInfinityEventSinkIntegration(t *testing.T) {
	_ = os.Setenv("TF_ACC", "1")

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

	testInfinityEventSink(t, client)
}
