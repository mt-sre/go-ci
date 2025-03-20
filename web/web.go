// SPDX-FileCopyrightText: 2025 Red Hat, Inc. <sd-mt-sre@redhat.com>
//
// SPDX-License-Identifier: Apache-2.0

package web

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
)

// FailedRequestError is returned when HTTP requests fail
// with an unexpected status code.
type FailedRequestError int

func (e FailedRequestError) Error() string {
	return fmt.Sprintf("request failed with status %d", e)
}

// DownloadFile retrieves the data contained at 'url' and write it to
// the 'out' path provided. If the request fails or the file cannot be
// written at the given location an error is returned.
func DownloadFile(ctx context.Context, url, out string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("constructing request: %w", err)
	}

	var client http.Client

	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("downloading file: %w", err)
	}

	defer res.Body.Close()

	if status := res.StatusCode; status != http.StatusOK {
		return FailedRequestError(status)
	}

	f, err := os.Create(out)
	if err != nil {
		return fmt.Errorf("creating file %q: %w", out, err)
	}

	defer f.Close()

	if _, err := io.Copy(f, res.Body); err != nil {
		return fmt.Errorf("copying response: %w", err)
	}

	return nil
}
