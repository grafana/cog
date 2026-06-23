package github

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"strings"

	"github.com/grafana/cog/internal/httputil"
)

// RepoDescriptor describes a repository.
type RepoDescriptor struct {
	Owner string
	Name  string
	Ref   string
}

// RawURLToRepoDescriptor extracts owner, repo, and ref from a raw GitHub URL.
// Handles both "refs/heads/main" and bare branch names.
func RawURLToRepoDescriptor(rawURL string) *RepoDescriptor {
	u, err := url.Parse(rawURL)
	if err != nil || u.Host != "raw.githubusercontent.com" {
		return nil
	}

	parts := strings.Split(strings.TrimPrefix(u.Path, "/"), "/")
	if len(parts) < 4 {
		return nil
	}

	owner, repo := parts[0], parts[1]
	var ref string
	if parts[2] == "refs" && len(parts) > 4 {
		ref = strings.Join(parts[2:5], "/")
	} else {
		ref = parts[2]
	}
	return &RepoDescriptor{Owner: owner, Name: repo, Ref: ref}
}

// FetchDirectory fetches all files with the given extension in a repository
// directory via the GitHub Contents API.
// Returns a map of filename → file content.
func FetchDirectory(ctx context.Context, repo RepoDescriptor, dirPath string, suffix string) (map[string][]byte, error) {
	// Note: ref may contain slashes (e.g. "refs/heads/main") which must NOT be
	// percent-encoded in the query string — the GitHub API expects them raw.
	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s?ref=%s", repo.Owner, repo.Name, dirPath, repo.Ref)

	var entries []struct {
		Name        string `json:"name"`
		Type        string `json:"type"`
		DownloadURL string `json:"download_url"`
	}
	if err := httputil.LoadJSON(ctx, apiURL, &entries); err != nil {
		return nil, err
	}

	files := make(map[string][]byte)
	for _, entry := range entries {
		if entry.Type != "file" || !strings.HasSuffix(entry.Name, suffix) {
			continue
		}

		body, err := httputil.LoadURL(ctx, entry.DownloadURL)
		if err != nil {
			return nil, err
		}
		defer func() { _ = body.Close() }()

		data, err := io.ReadAll(body)
		if err != nil {
			continue
		}

		files[entry.Name] = data
	}
	return files, nil
}

// DirPath extracts the repo-relative directory path from a raw GitHub URL.
// For "https://raw.githubusercontent.com/owner/repo/refs/heads/main/apps/iam/kinds/manifest.cue"
// it returns "apps/iam/kinds".
func DirPath(rawURL string) string {
	ghInfo := RawURLToRepoDescriptor(rawURL)
	if ghInfo == nil {
		return ""
	}

	u, err := url.Parse(rawURL)
	if err != nil || u.Host != "raw.githubusercontent.com" {
		return ""
	}

	refParts := strings.Split(ghInfo.Ref, "/")

	// allParts: [owner, repo, ref_part1, ..., ref_partN, path_part1, ..., filename]
	allParts := strings.Split(strings.TrimPrefix(u.Path, "/"), "/")

	skip := 2 + len(refParts) // skip owner + repo + ref segments
	if skip >= len(allParts) {
		return ""
	}

	pathParts := allParts[skip : len(allParts)-1] // exclude filename
	return strings.Join(pathParts, "/")
}
