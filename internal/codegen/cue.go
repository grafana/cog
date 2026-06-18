package codegen

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/load"
	"cuelang.org/go/cue/parser"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/httputil"
	"github.com/grafana/cog/internal/simplecue"
	"github.com/grafana/cog/internal/tools"
)

type genericCueLoader struct {
	*CueInput

	loader func(input CueInput) (ast.Schemas, error)
}

func (loader *genericCueLoader) LoadSchemas(_ context.Context) (ast.Schemas, error) {
	return loader.loader(*loader.CueInput)
}

type CueInput struct {
	InputBase `yaml:",inline"`

	// Entrypoint refers to a directory containing CUE files.
	Entrypoint string `yaml:"entrypoint"`

	// URL to a cue file
	URL string `yaml:"url"`

	// Subpath is used when the schema is not in the root of the schema
	Subpath string `yaml:"subpath"`

	// Value represents the CUE value to use as an input. If specified, it
	// supersedes the Entrypoint and URL options.
	Value *cue.Value `yaml:"-"`

	// ForcedEnvelope decorates the parsed cue Value with an envelope whose
	// name is given. This is useful for dataqueries for example, where the
	// schema doesn't define any suitable top-level object.
	ForcedEnvelope string `yaml:"forced_envelope"`

	// Package name to use for the input schema. If empty, it will be guessed.
	Package string `yaml:"package"`

	// CueImports allows importing additional libraries.
	// Format: [path]:[import]. Example: '../grafana/common-library:github.com/grafana/grafana/packages/grafana-schema/src/common
	CueImports []string `yaml:"cue_imports"`

	// NameFunc allows users to specify an alternative naming strategy for
	// objects and references. It is called with the value passed to the top
	// level method or function and the path to the entity being parsed.
	NameFunc simplecue.NameFunc `yaml:"-"`

	// InlineExternalReference instructs the parser to follow external
	// references (ie: references to objects outside the current schema)
	// and inline them.
	// By default, external references are parsed as actual `ast.Ref` to the
	// external objects.
	InlineExternalReference bool `yaml:"-"`
}

func (input *CueInput) packageName() string {
	if input.Package != "" {
		return input.Package
	}

	return filepath.Base(input.Entrypoint)
}

func (input *CueInput) schemaRootValue() (cue.Value, []simplecue.LibraryInclude, error) {
	if input.Value != nil {
		libraries, err := simplecue.ParseImports(input.CueImports)
		if err != nil {
			return cue.Value{}, nil, err
		}
		return *input.Value, libraries, nil
	}

	if input.Entrypoint == "" && input.URL == "" {
		return cue.Value{}, nil, fmt.Errorf("no entrypoint or url defined in cue input")
	}

	if input.Entrypoint != "" && input.URL != "" {
		return cue.Value{}, nil, fmt.Errorf("only one entrypoint or url defined in cue input")
	}

	libraries, err := simplecue.ParseImports(input.CueImports)
	if err != nil {
		return cue.Value{}, nil, err
	}

	value, err := input.parseCueEntrypoint(libraries)

	return value, libraries, err
}

func (input *CueInput) interpolateParameters(interpolator ParametersInterpolator) {
	input.InputBase.interpolateParameters(interpolator)

	input.Entrypoint = interpolator(input.Entrypoint)
	input.URL = interpolator(input.URL)
	input.CueImports = tools.Map(input.CueImports, interpolator)
}

func cueLoader(input CueInput) (ast.Schemas, error) {
	schemaRootValue, libraries, err := input.schemaRootValue()
	if err != nil {
		return nil, err
	}

	schema, err := simplecue.GenerateAST(schemaRootValue, simplecue.Config{
		Package:                 input.packageName(),
		ForceNamedEnvelope:      input.ForcedEnvelope,
		SchemaMetadata:          input.schemaMetadata(),
		Libraries:               libraries,
		NameFunc:                input.NameFunc,
		InlineExternalReference: input.InlineExternalReference,
	})
	if err != nil {
		return nil, err
	}

	return input.filterSchema(schema)
}

func (input *CueInput) parseCueEntrypoint(imports []simplecue.LibraryInclude) (cue.Value, error) {
	overlay := map[string]load.Source{
		"/cog/vfs/cue.mod/module.cue": load.FromBytes([]byte(
			"language: { version: \"v0.10.1\" }\nmodule: \"cog.vfs\"\n",
		)),
	}

	if err := addLibrariesToOverlay(overlay, imports); err != nil {
		return cue.Value{}, err
	}

	var instanceName string
	prefix := "/cog/vfs/cue.mod/pkg/"

	if input.URL != "" {
		var err error
		instanceName, err = addURLToOverlay(overlay, prefix, input.URL)
		if err != nil {
			return cue.Value{}, err
		}
	}

	if input.Entrypoint != "" {
		instanceName = "github.com/cog-vfs/" + filepath.Base(input.Entrypoint)

		if err := addDirToOverlay(overlay, prefix+instanceName, input.Entrypoint); err != nil {
			return cue.Value{}, err
		}
	}

	// Load Cue files into Cue build.Instances slice
	bis := load.Instances([]string{instanceName}, &load.Config{
		Overlay: overlay,
		// Point cue to a directory defined by the cueFsOverlay as base directory
		// for import path resolution instead of using the current working directory.
		// This ensures that only files/schemas defined in the vfs will be parsed.
		Dir: "/cog/vfs",
	})

	value := cuecontext.New().BuildInstance(bis[0])
	if err := value.Err(); err != nil {
		return cue.Value{}, fmt.Errorf("could not build cue instance: %w", err)
	}

	if input.Subpath != "" {
		value = value.LookupPath(cue.ParsePath(input.Subpath))
		if value.Err() != nil {
			return cue.Value{}, fmt.Errorf("could not read subpath '%s': %w", input.Subpath, value.Err())
		}
	}

	if input.Package == "" {
		input.Package = bis[0].PkgName
	}

	return value, nil
}

func addDirToOverlay(overlay map[string]load.Source, overlayPrefix string, directory string) error {
	absolutePath, err := filepath.Abs(directory)
	if err != nil {
		return err
	}

	return filepath.Walk(absolutePath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		overlay[filepath.Join(overlayPrefix, info.Name())] = load.FromBytes(content)

		return nil
	})
}

func addLibrariesToOverlay(overlay map[string]load.Source, imports []simplecue.LibraryInclude) error {
	for _, importDefinition := range imports {
		prefix := "/cog/vfs/cue.mod/pkg/" + importDefinition.ImportPath

		if err := addDirToOverlay(overlay, prefix, importDefinition.FSPath); err != nil {
			return err
		}
	}

	return nil
}

func addURLToOverlay(overlay map[string]load.Source, overlayPrefix string, entrypoint string) (string, error) {
	u, err := url.Parse(entrypoint)
	if err != nil {
		return "", err
	}

	if !strings.HasSuffix(u.Path, ".cue") {
		return "", fmt.Errorf("entrypoint %q must be a cue url", entrypoint)
	}

	ghInfo := parseRawGitHubURL(entrypoint)
	dirPath := repoDirPath(entrypoint)

	dirFiles, err := fetchGitHubDirectory(ghInfo.owner, ghInfo.repo, dirPath, ghInfo.ref)
	if err != nil {
		return "", fmt.Errorf("fetch directory: %w", err)
	}

	allFiles := make(map[string][]byte, len(dirFiles))
	fetched := make(map[string]bool)

	var instanceName string
	for name, content := range dirFiles {
		file, err := parser.ParseFile(name, content, parser.ParseComments)
		if err != nil {
			continue
		}

		if instanceName == "" {
			instanceName = "github.com/cog-vfs/" + file.PackageName()
		}

		overlay[overlayPrefix+"/"+instanceName+"/"+name] = load.FromBytes(content)
		allFiles[name] = content

		for _, imp := range file.Imports {
			importPath := strings.Trim(imp.Path.Value, `"`)
			if !strings.HasPrefix(importPath, "github.com") || fetched[importPath] {
				continue
			}
			fetched[importPath] = true
			ghInfo = parseRawGitHubURL(importPath)
			repoPath := strings.TrimPrefix(importPath, "github.com")
			pkgFiles, err := fetchGitHubDirectory(ghInfo.owner, ghInfo.repo, repoPath, ghInfo.ref)
			if err != nil {
				continue
			}
			for fname, fcontent := range pkgFiles {
				overlayPath := filepath.Join("/cog/vfs/cue.mod/pkg", importPath, fname)
				overlay[overlayPath] = load.FromBytes(fcontent)
			}
		}
	}

	return instanceName, nil
}

// repoDirPath extracts the repo-relative directory path from a raw GitHub URL.
// For "https://raw.githubusercontent.com/owner/repo/refs/heads/main/apps/iam/kinds/manifest.cue"
// it returns "apps/iam/kinds".
func repoDirPath(rawURL string) string {
	u, err := url.Parse(rawURL)
	if err != nil || u.Host != "raw.githubusercontent.com" {
		return ""
	}
	ghInfo := parseRawGitHubURL(rawURL)
	if ghInfo == nil {
		return ""
	}
	refParts := strings.Split(ghInfo.ref, "/")
	// allParts: [owner, repo, ref_part1, ..., ref_partN, path_part1, ..., filename]
	allParts := strings.Split(strings.TrimPrefix(u.Path, "/"), "/")
	skip := 2 + len(refParts) // skip owner + repo + ref segments
	if skip >= len(allParts) {
		return ""
	}
	pathParts := allParts[skip : len(allParts)-1] // exclude filename
	return strings.Join(pathParts, "/")
}

// rawGitHubURLInfo holds components extracted from a raw.githubusercontent.com URL.
type rawGitHubURLInfo struct {
	owner string
	repo  string
	ref   string
}

// parseRawGitHubURL extracts owner, repo, and ref from a raw GitHub URL.
// Handles both "refs/heads/main" and bare branch names.
func parseRawGitHubURL(rawURL string) *rawGitHubURLInfo {
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
	return &rawGitHubURLInfo{owner: owner, repo: repo, ref: ref}
}

// fetchGitHubDirectory fetches all *.cue files in a repository directory via
// the GitHub Contents API. Returns a map of filename → file content.
func fetchGitHubDirectory(owner string, repo string, dirPath string, ref string) (map[string][]byte, error) {
	// Note: ref may contain slashes (e.g. "refs/heads/main") which must NOT be
	// percent-encoded in the query string — the GitHub API expects them raw.
	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s?ref=%s", owner, repo, dirPath, ref)

	var entries []struct {
		Name        string `json:"name"`
		Type        string `json:"type"`
		DownloadURL string `json:"download_url"`
	}
	if err := httputil.LoadJSON(context.Background(), apiURL, &entries); err != nil {
		return nil, err
	}

	files := make(map[string][]byte)
	for _, entry := range entries {
		if entry.Type != "file" || !strings.HasSuffix(entry.Name, ".cue") {
			continue
		}

		body, err := httputil.LoadURL(context.Background(), entry.DownloadURL)
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
