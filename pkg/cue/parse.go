package cue

import (
	"context"
	"fmt"
	"io/fs"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"cuelang.org/go/cue"
	cueast "cuelang.org/go/cue/ast"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/load"
	"cuelang.org/go/cue/parser"
	"github.com/grafana/cog/internal/github"
)

type LibraryInclude struct {
	FSPath     string // path of the library on the filesystem
	ImportPath string // path used in CUE files to import that library
}

type Input struct {
	// Entrypoint refers to a directory containing CUE files.
	Entrypoint string

	// URL to a cue file
	URL string

	// Subpath is used when the schema is not in the root of the schema
	Subpath string

	// Value represents the CUE value to use as an input. If specified, it
	// supersedes the Entrypoint and URL options.
	Value *cue.Value

	// CueImports allows importing additional libraries.
	// Format: [path]:[import]. Example: '../grafana/common-library:github.com/grafana/grafana/packages/grafana-schema/src/common
	CueImports []string
}

type ParseResult struct {
	Package   string
	Value     cue.Value
	Libraries []LibraryInclude
}

// Parse the given input as a CUE value.
// Note: imports referring to modules hosted on GitHub will be resolved and
// downloaded on the fly.
func Parse(ctx context.Context, input Input) (ParseResult, error) {
	result := ParseResult{}

	libraries, err := parseImports(input.CueImports)
	if err != nil {
		return result, err
	}

	result.Libraries = libraries

	if input.Value != nil {
		result.Value = *input.Value
	} else {
		val, pkgName, err := parseCueEntrypoint(ctx, input, libraries)
		if err != nil {
			return result, fmt.Errorf("could not parse cue entrypoint: %w", err)
		}

		result.Value = val
		result.Package = pkgName
	}

	if input.Subpath != "" {
		result.Value = result.Value.LookupPath(cue.ParsePath(input.Subpath))
		if result.Value.Err() != nil {
			return result, fmt.Errorf("could not read subpath '%s': %w", input.Subpath, result.Value.Err())
		}
	}

	return result, nil
}

type cueOverlay struct {
	files     map[string]load.Source
	instances map[string]bool
}

func parseCueEntrypoint(ctx context.Context, input Input, imports []LibraryInclude) (cue.Value, string, error) {
	var err error
	var instanceName string

	if input.Entrypoint == "" && input.URL == "" {
		return cue.Value{}, "", fmt.Errorf("no entrypoint or url defined in cue input")
	}

	if input.Entrypoint != "" && input.URL != "" {
		return cue.Value{}, "", fmt.Errorf("only one entrypoint or url defined in cue input")
	}

	prefix := "/cog/vfs/cue.mod/pkg"
	overlay := cueOverlay{
		files: map[string]load.Source{
			"/cog/vfs/cue.mod/module.cue": load.FromBytes([]byte(
				"language: { version: \"v0.10.1\" }\nmodule: \"cog.vfs\"\n",
			)),
		},
		instances: map[string]bool{},
	}

	if err = addLibrariesToOverlay(ctx, overlay, prefix, imports); err != nil {
		return cue.Value{}, "", fmt.Errorf("could not add libraries to overlay: %w", err)
	}

	if input.URL != "" {
		instanceName, err = addURLToOverlay(ctx, overlay, prefix, input.URL)
		if err != nil {
			return cue.Value{}, "", fmt.Errorf("could not add URL '%s' to overlay: %w", input.URL, err)
		}
	}

	if input.Entrypoint != "" {
		instanceName, err = addDirToOverlay(ctx, overlay, prefix, input.Entrypoint)
		if err != nil {
			return cue.Value{}, "", fmt.Errorf("could not add entrypoint '%s' to overlay: %w", input.Entrypoint, err)
		}
	}

	// Load Cue files into Cue build.Instances slice
	bis := load.Instances([]string{instanceName}, &load.Config{
		Overlay: overlay.files,
		// Point cue to a directory defined by the cueFsOverlay as base directory
		// for import path resolution instead of using the current working directory.
		// This ensures that only files/schemas defined in the vfs will be parsed.
		Dir: "/cog/vfs",
	})

	value := cuecontext.New().BuildInstance(bis[0])
	if err := value.Err(); err != nil {
		return cue.Value{}, "", fmt.Errorf("could not build cue instance: %w", err)
	}

	return value, bis[0].PkgName, nil
}

// addDirToOverlay adds the contents of the given directory to the overlay, as an entrypoint.
func addDirToOverlay(ctx context.Context, overlay cueOverlay, prefix string, directory string) (string, error) {
	absolutePath, err := filepath.Abs(directory)
	if err != nil {
		return "", err
	}

	var instanceName string
	err = filepath.Walk(absolutePath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if !strings.HasSuffix(info.Name(), ".cue") {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		file, err := parser.ParseFile(path, content, parser.ParseComments)
		if err != nil {
			return fmt.Errorf("parse error in '%s': %w", path, err)
		}

		if instanceName == "" {
			instanceName = "github.com/cog-vfs/" + file.PackageName()
		}

		instance := prefix + "/" + instanceName

		overlay.files[instance+"/"+info.Name()] = load.FromBytes(content)
		overlay.instances[instance] = true

		err = addFileImportsToOverlay(ctx, overlay, prefix, file)
		if err != nil {
			return fmt.Errorf("could add file imports to overlay: %w", err)
		}

		return nil
	})
	if err != nil {
		return "", err
	}

	return instanceName, nil
}

// addLibrariesToOverlay adds libraries from the local filesystem to the overlay.
// The libraries included here will be available for the entrypoint to use,
// similar to how an "include path" would behave.
func addLibrariesToOverlay(ctx context.Context, overlay cueOverlay, prefix string, imports []LibraryInclude) error {
	for _, importDefinition := range imports {
		err := filepath.Walk(importDefinition.FSPath, func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			if !strings.HasSuffix(info.Name(), ".cue") {
				return nil
			}

			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			instance := prefix + "/" + importDefinition.ImportPath

			overlay.files[instance+"/"+info.Name()] = load.FromBytes(content)
			overlay.instances[instance] = true

			file, err := parser.ParseFile(path, content, parser.ParseComments)
			if err != nil {
				return fmt.Errorf("parse error in '%s': %w", path, err)
			}

			err = addFileImportsToOverlay(ctx, overlay, prefix, file)
			if err != nil {
				return fmt.Errorf("could add file imports to overlay: %w", err)
			}

			return nil
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// addDirToOverlay adds the contents of the given URL to the overlay, as an entrypoint.
func addURLToOverlay(ctx context.Context, overlay cueOverlay, overlayPrefix string, entrypoint string) (string, error) {
	u, err := url.Parse(entrypoint)
	if err != nil {
		return "", err
	}

	if !strings.HasSuffix(u.Path, ".cue") {
		return "", fmt.Errorf("entrypoint %s must be a cue url", entrypoint)
	}

	ghInfo := github.RawURLToRepoDescriptor(entrypoint)
	if ghInfo == nil {
		return "", fmt.Errorf("could not parse raw url as repo descriptor: %s", entrypoint)
	}

	dirFiles, err := github.FetchDirectory(ctx, *ghInfo, github.DirPath(entrypoint), ".cue")
	if err != nil {
		return "", fmt.Errorf("fetch directory: %w", err)
	}

	var instanceName string
	for name, content := range dirFiles {
		file, err := parser.ParseFile(name, content, parser.ParseComments)
		if err != nil {
			return "", fmt.Errorf("parse error in '%s': %w", name, err)
		}

		if instanceName == "" {
			instanceName = "github.com/cog-vfs/" + file.PackageName()
		}

		instance := overlayPrefix + "/" + instanceName

		overlay.files[instance+"/"+name] = load.FromBytes(content)
		overlay.instances[instance] = true

		err = addFileImportsToOverlay(ctx, overlay, overlayPrefix, file)
		if err != nil {
			return "", fmt.Errorf("could add file imports to overlay: %w", err)
		}
	}

	return instanceName, nil
}

func addFileImportsToOverlay(ctx context.Context, overlay cueOverlay, overlayPrefix string, file *cueast.File) error {
	for _, imp := range file.Imports {
		importPath := strings.Trim(imp.Path.Value, `"`)
		instance := overlayPrefix + "/" + importPath

		if !strings.HasPrefix(importPath, "github.com") || overlay.instances[instance] {
			// TODO: log? (could also just be an import if a standard pkg)
			continue
		}

		// TODO: make the ref configurable?
		ghInfo := github.URLToRepoDescriptor("https://"+importPath, "main")
		if ghInfo == nil {
			return fmt.Errorf("could not parse import path as repo descriptor: %s", importPath)
		}

		repoPath := strings.TrimPrefix(importPath, fmt.Sprintf("github.com/%s/%s/", ghInfo.Owner, ghInfo.Name))
		pkgFiles, err := github.FetchDirectory(ctx, *ghInfo, repoPath, ".cue")
		if err != nil {
			return fmt.Errorf("could not fetch import path %s: %w", importPath, err)
		}

		for fileName, fileContent := range pkgFiles {
			overlay.files[instance+"/"+fileName] = load.FromBytes(fileContent)
			overlay.instances[instance] = true
		}
	}

	return nil
}

func parseImports(cueImports []string) ([]LibraryInclude, error) {
	if len(cueImports) == 0 {
		return nil, nil
	}

	imports := make([]LibraryInclude, len(cueImports))
	for i, importDefinition := range cueImports {
		parts := strings.Split(importDefinition, ":")
		if len(parts) != 2 {
			return nil, fmt.Errorf("'%s' is not a valid import definition", importDefinition)
		}

		imports[i] = LibraryInclude{
			FSPath:     parts[0],
			ImportPath: parts[1],
		}
	}

	return imports, nil
}
