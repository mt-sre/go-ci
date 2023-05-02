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

// TestDownloadFile tests the behavior of the DownloadFile
// function in different scenarios
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

	tests := []struct {
		name        string
		url         string
		out         string
		expectError bool
	}{
		{
			name:        "invalid URL",
			url:         "invalid_url",
			out:         "outfile1",
			expectError: true,
		},
		{
			name:        "server error",
			url:         "",
			out:         "outfile2",
			expectError: true,
		},
		{
			name:        "file creation error",
			url:         srv.URL,
			out:         "/nonexistent/dir/outfile3",
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			out := filepath.Join(t.TempDir(), tc.out)

			if tc.name == "server error" {
				srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					http.Error(w, "server error", http.StatusInternalServerError)
				}))
				defer srv.Close()
				tc.url = srv.URL
			}

			err := web.DownloadFile(context.Background(), tc.url, out)
			if tc.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// TestFailedRequestError checks that the Error method returns the
// expected message.
func TestFailedRequestError(t *testing.T) {
	err := web.FailedRequestError(404)

	require.NotNil(t, err, "error should not be nil")
	assert.Equal(t, "request failed with status 404", err.Error(), "error message should match")
}
