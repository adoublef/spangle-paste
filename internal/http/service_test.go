package http_test

import (
	"io"
	"net/http/httptest"
	"testing"

	api "github.com/adoublef/spangle-paste/internal/http"
	"github.com/stretchr/testify/require"
)

func TestService(t *testing.T) {
	srv := newTestService(t)

	t.Cleanup(func() { srv.Close() })

	// TODO -- make http request
	t.Run("returns `Ciao, Foo!`", func(t *testing.T) {
		resp, err := srv.Client().Get(srv.URL+"?name=Foo")
		require.NoError(t, err, "failed to make get request")

		p, err := io.ReadAll(resp.Body)
		require.NoError(t, err, "failed to read body")

		require.Equal(t, "Ciao, Foo!", string(p))
	})
}

func newTestService(t *testing.T) *httptest.Server {
	h := api.NewService()
	return httptest.NewServer(h)
}
