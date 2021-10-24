package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/terdia/greenlight/config"
	"github.com/terdia/greenlight/infrastructures/logger"
	"github.com/terdia/greenlight/mock"
)

type testWriter struct{}

type TestResponse struct {
	StatusCode int
	Header     http.Header
	Body       []byte
}

func (tl *testWriter) Write(p []byte) (n int, err error) {
	return 0, nil
}

func newTestApplication(t *testing.T, totalMovies int) *application {

	cfg := new(config.Config)

	cfg.Limiter.Rps = 1
	cfg.Limiter.Burst = 2
	cfg.Limiter.Enabled = true
	cfg.Env = "test"
	cfg.Version = "v1"

	out := &testWriter{}
	logger := logger.New(out, logger.LevelInfo)

	wg := new(sync.WaitGroup)

	return &application{
		config:   cfg,
		logger:   logger,
		registry: mock.NewRegistry(logger, mock.NewMailerMock(), wg, totalMovies),
		wg:       wg,
	}
}

// Define a custom testServer type which anonymously embeds a httptest.Server instance.
type testServer struct {
	*httptest.Server
}

// Create a newTestServer helper which initalizes and returns a new instance
// of our custom testServer type.
func newTestServer(t *testing.T, h http.Handler) *testServer {
	ts := httptest.NewServer(h)

	// Disable redirect-following for the client. Essentially this function
	// is called after a 3xx response is received by the client, and returning
	// the http.ErrUseLastResponse error forces it to immediately return the
	// received response.
	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	return &testServer{ts}
}

func (ts *testServer) get(t *testing.T, urlPath string) TestResponse {
	rs, err := ts.Client().Get(ts.URL + urlPath)
	if err != nil {
		t.Fatal(err)
	}
	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	return TestResponse{
		StatusCode: rs.StatusCode,
		Header:     rs.Header,
		Body:       body,
	}
}
