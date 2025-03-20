// SPDX-FileCopyrightText: 2025 Red Hat, Inc. <sd-mt-sre@redhat.com>
//
// SPDX-License-Identifier: Apache-2.0

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
type ErrorAssertionFunc func(t require.TestingT, err error, msgAndArgs ...interface{})

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
		Name      string
		Url       string
		Out       string
		Assertion ErrorAssertionFunc
	}{
		{
			Name:      "invalid URL",
			Url:       "invalid_url",
			Out:       "outfile1",
			Assertion: require.Error,
		},
		{
			Name:      "server error",
			Url:       "",
			Out:       "outfile2",
			Assertion: require.Error,
		},
		{
			Name:      "file creation error",
			Url:       srv.URL,
			Out:       "/nonexistent/dir/outfile3",
			Assertion: require.Error,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			out := filepath.Join(t.TempDir(), tc.Out)

			if tc.Name == "server error" {
				srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					http.Error(w, "server error", http.StatusInternalServerError)
				}))
				defer srv.Close()
				tc.Url = srv.URL
			}

			err := web.DownloadFile(context.Background(), tc.Url, out)
			tc.Assertion(t, err)
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
