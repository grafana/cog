package loaders

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
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
