package loaders

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/ast/compiler"
	"github.com/grafana/cog/internal/tools"
)

func guessPackageFromFilename(filename string) string {
	pkg := filepath.Base(filepath.Dir(filename))
	if pkg != "." {
		return pkg
	}

	return strings.TrimSuffix(filepath.Base(filename), filepath.Ext(filename))
}

func dirExists(dir string) (bool, error) {
	stat, err := os.Stat(dir)
	//nolint:gocritic
	if os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, err
	} else if !stat.IsDir() {
		return false, fmt.Errorf("'%s' is not a directory", dir)
	}

	return true, nil
}

func filterSchema(schema *ast.Schema, allowedObjects []string) (ast.Schemas, error) {
	if len(allowedObjects) == 0 {
		return ast.Schemas{schema}, nil
	}

	filter := compiler.FilterSchemas{
		AllowedObjects: tools.Map(allowedObjects, func(objectName string) compiler.ObjectReference {
			return compiler.ObjectReference{Package: schema.Package, Object: objectName}
		}),
	}

	return filter.Process(ast.Schemas{schema})
}

func loadURL(ctx context.Context, url string) (io.ReadCloser, error) {
	client := http.DefaultClient

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
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
