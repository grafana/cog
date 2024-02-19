package common

import (
	"fmt"

	"github.com/grafana/cog/internal/orderedmap"
)

func NoopImportSanitizer(s string) string { return s }

type ImportMapConfig[M any] struct {
	Formatter           func(importMap M) string
	AliasSanitizer      func(string) string
	ImportPathSanitizer func(string) string
}

type ImportMapOption[M any] func(importMap *ImportMapConfig[M])

func WithAliasSanitizer[M any](sanitizer func(string) string) ImportMapOption[M] {
	return func(importMap *ImportMapConfig[M]) {
		importMap.AliasSanitizer = sanitizer
	}
}

func WithImportPathSanitizer[M any](sanitizer func(string) string) ImportMapOption[M] {
	return func(importMap *ImportMapConfig[M]) {
		importMap.ImportPathSanitizer = sanitizer
	}
}

func WithFormatter[M any](formatter func(M) string) ImportMapOption[M] {
	return func(importMap *ImportMapConfig[M]) {
		importMap.Formatter = formatter
	}
}

type DirectImportMap struct {
	Imports *orderedmap.Map[string, string] // alias → importPath
	config  ImportMapConfig[DirectImportMap]
}

func NewDirectImportMap(opts ...ImportMapOption[DirectImportMap]) *DirectImportMap {
	config := ImportMapConfig[DirectImportMap]{
		Formatter: func(importMap DirectImportMap) string {
			return fmt.Sprintf("%#v\n", importMap.Imports)
		},
		AliasSanitizer:      NoopImportSanitizer,
		ImportPathSanitizer: NoopImportSanitizer,
	}

	for _, opt := range opts {
		opt(&config)
	}

	return &DirectImportMap{
		Imports: orderedmap.New[string, string](),
		config:  config,
	}
}

func (im DirectImportMap) Add(alias string, importPath string) string {
	sanitizedAlias := im.config.AliasSanitizer(alias)
	im.Imports.Set(sanitizedAlias, im.config.ImportPathSanitizer(importPath))

	return sanitizedAlias
}

func (im DirectImportMap) String() string {
	return im.config.Formatter(im)
}
