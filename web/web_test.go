package web_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/mt-sre/go-ci/web"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDownloadFile(t *testing.T) {
	expectedContent := []byte("hello")

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write(expectedContent)
	}))

	defer srv.Close()

	out := filepath.Join(t.TempDir(), "outfile")

	err := web.DownloadFile(context.Background(), srv.URL, out)
	require.NoError(t, err)

	require.FileExists(t, out)

	data, err := os.ReadFile(out)
	require.NoError(t, err)

	assert.Equal(t, expectedContent, data)
}
