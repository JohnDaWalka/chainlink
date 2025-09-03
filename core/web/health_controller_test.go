package web_test

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
	"text/template"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/config"
	solcfg "github.com/smartcontractkit/chainlink-solana/pkg/solana/config"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/mocks"
)

func TestHealthController_Readyz(t *testing.T) {
	var tt = []struct {
		name   string
		ready  bool
		status int
	}{
		{
			name:   "not ready",
			ready:  false,
			status: http.StatusServiceUnavailable,
		},
		{
			name:   "ready",
			ready:  true,
			status: http.StatusOK,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			app := cltest.NewApplicationWithKey(t)
			healthChecker := new(mocks.Checker)
			healthChecker.On("Start").Return(nil).Once()
			healthChecker.On("IsReady").Return(tc.ready, nil).Once()
			healthChecker.On("Close").Return(nil).Once()

			app.HealthChecker = healthChecker
			require.NoError(t, app.Start(testutils.Context(t)))

			client := app.NewHTTPClient(nil)
			resp, cleanup := client.Get("/readyz")
			t.Cleanup(cleanup)
			assert.Equal(t, tc.status, resp.StatusCode)
		})
	}
}

func TestHealthController_Health_status(t *testing.T) {
	var tt = []struct {
		name   string
		ready  bool
		status int
	}{
		{
			name:   "not ready",
			ready:  false,
			status: http.StatusMultiStatus,
		},
		{
			name:   "ready",
			ready:  true,
			status: http.StatusOK,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			app := cltest.NewApplicationWithKey(t)
			healthChecker := new(mocks.Checker)
			healthChecker.On("Start").Return(nil).Once()
			healthChecker.On("IsHealthy").Return(tc.ready, nil).Once()
			healthChecker.On("Close").Return(nil).Once()

			app.HealthChecker = healthChecker
			require.NoError(t, app.Start(testutils.Context(t)))

			client := app.NewHTTPClient(nil)
			resp, cleanup := client.Get("/health")
			t.Cleanup(cleanup)
			assert.Equal(t, tc.status, resp.StatusCode)
		})
	}
}

var (
	//go:embed testdata/body/health.json
	bodyJSON string
	//go:embed testdata/body/health.html
	bodyHTML string
	//go:embed testdata/body/health.txt
	bodyTXT string
	//go:embed testdata/body/health-failing.json
	bodyJSONFailing string
	//go:embed testdata/body/health-failing.html
	bodyHTMLFailing string
	//go:embed testdata/body/health-failing.txt
	bodyTXTFailing string
)

func TestHealthController_Health_body(t *testing.T) {
	templateData := map[string]interface{}{
		"chainID": testutils.FixtureChainID.String(),
	}

	bodyJSONTmplt, err := template.New("health.json").Parse(bodyJSON)
	require.NoError(t, err)
	bodyJSONRes := &bytes.Buffer{}
	bodyJSONTmplt.Execute(bodyJSONRes, templateData)

	bodyHTMLTmplt, err := template.New("health.html").Parse(bodyHTML)
	require.NoError(t, err)
	bodyHTMLRes := &bytes.Buffer{}
	bodyHTMLTmplt.Execute(bodyHTMLRes, templateData)

	bodyTXTmplt, err := template.New("health.txt").Parse(bodyTXT)
	require.NoError(t, err)
	bodyTXTRes := &bytes.Buffer{}
	bodyTXTmplt.Execute(bodyTXTRes, templateData)

	bodyJSONFailingTmplt, err := template.New("health.json").Parse(bodyJSONFailing)
	require.NoError(t, err)
	bodyJSONFailingRes := &bytes.Buffer{}
	bodyJSONFailingTmplt.Execute(bodyJSONFailingRes, templateData)

	bodyHTMLFailingTmplt, err := template.New("health.html").Parse(bodyHTMLFailing)
	require.NoError(t, err)
	bodyHTMLFailingRes := &bytes.Buffer{}
	bodyHTMLFailingTmplt.Execute(bodyHTMLFailingRes, templateData)

	bodyTXTFailingTmplt, err := template.New("health.txt").Parse(bodyTXTFailing)
	require.NoError(t, err)
	bodyTXTFailingRes := &bytes.Buffer{}
	bodyTXTFailingTmplt.Execute(bodyTXTFailingRes, templateData)

	for _, tc := range []struct {
		name    string
		path    string
		headers map[string]string
		expBody string
	}{
		{"default", "/health", nil, bodyJSONRes.String()},
		{"json", "/health", map[string]string{"Accept": gin.MIMEJSON}, bodyJSONRes.String()},
		{"html", "/health", map[string]string{"Accept": gin.MIMEHTML}, bodyHTMLRes.String()},
		{"text", "/health", map[string]string{"Accept": gin.MIMEPlain}, bodyTXTRes.String()},
		{".txt", "/health.txt", nil, bodyTXTRes.String()},

		{"default-failing", "/health?failing", nil, bodyJSONFailingRes.String()},
		{"json-failing", "/health?failing", map[string]string{"Accept": gin.MIMEJSON}, bodyJSONFailingRes.String()},
		{"html-failing", "/health?failing", map[string]string{"Accept": gin.MIMEHTML}, bodyHTMLFailingRes.String()},
		{"text-failing", "/health?failing", map[string]string{"Accept": gin.MIMEPlain}, bodyTXTFailingRes.String()},
		{".txt-failing", "/health.txt?failing", nil, bodyTXTFailingRes.String()},
	} {
		t.Run(tc.name, func(t *testing.T) {
			cfg := configtest.NewGeneralConfig(t, func(cfg *chainlink.Config, secrets *chainlink.Secrets) {
				cfg.Solana = append(cfg.Solana, &solcfg.TOMLConfig{
					ChainID: ptr("Bar"),
					Nodes: solcfg.Nodes{
						{Name: ptr("primary"), URL: config.MustParseURL("http://solana.web")},
					},
				})
				cfg.Solana[0].SetDefaults()
			})
			app := cltest.NewApplicationWithConfigAndKey(t, cfg)
			require.NoError(t, app.Start(testutils.Context(t)))

			client := app.NewHTTPClient(nil)
			resp, cleanup := client.Get(tc.path, tc.headers)
			t.Cleanup(cleanup)
			assert.Equal(t, http.StatusMultiStatus, resp.StatusCode)
			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			if tc.expBody == bodyJSONRes.String() {
				// pretty print for comparison
				var b bytes.Buffer
				require.NoError(t, json.Indent(&b, body, "", "  "))
				body = b.Bytes()
			}
			assert.Equal(t, strings.TrimSpace(tc.expBody), strings.TrimSpace(string(body)))
		})
	}
}
