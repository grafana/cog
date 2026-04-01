package httputil

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
)

func LoadURL(ctx context.Context, url string) (io.ReadCloser, error) {
	client := http.DefaultClient

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	// Necessary for Github private repositories
	authToken := os.Getenv("GITHUB_AUTH_TOKEN")
	if authToken != "" && req.URL.Host == "raw.githubusercontent.com" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", authToken))
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("expecting 200 when loading '%s', got %d", url, resp.StatusCode)
	}

	return resp.Body, nil
}
