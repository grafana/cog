package httputil

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

func LoadURL(ctx context.Context, rawURL string) (io.ReadCloser, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}

	if u.Scheme != "https" {
		return nil, fmt.Errorf("unsupported scheme '%s'", u.Scheme)
	}

	allowedHosts := map[string]struct{}{
		"raw.githubusercontent.com": {},
	}

	if _, ok := allowedHosts[u.Hostname()]; !ok {
		return nil, fmt.Errorf("unsupported host '%s'", u.Hostname())
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	// Necessary for Github private repositories
	authToken := os.Getenv("GITHUB_AUTH_TOKEN")
	if authToken != "" && req.URL.Host == "raw.githubusercontent.com" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", authToken))
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("expecting 200 when loading '%s', got %d", u.String(), resp.StatusCode)
	}

	return resp.Body, nil
}
