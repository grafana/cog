package csharp

import (
	"sort"
	"strings"
)

// importMap collects the set of `using` directives needed by a generated
// C# source file. Imports are de-duplicated by namespace and emitted in
// alphabetical order.
//
// `namespaceRoot` is prepended to schema-package imports so they are
// fully qualified (e.g. "Dashboard" -> "Grafana.Foundation.Dashboard").
// Built-in BCL namespaces (anything beginning with `System` or
// `System.*`) are passed through unchanged.
type importMap struct {
	namespaceRoot string
	currentNs     string // namespace of the file being generated; never imported
	namespaces    map[string]struct{}
}

func newImportMap(namespaceRoot string, currentNs string) *importMap {
	return &importMap{
		namespaceRoot: namespaceRoot,
		currentNs:     currentNs,
		namespaces:    make(map[string]struct{}),
	}
}

// addPackage records a `using` for a schema package (the namespace root
// is applied automatically). A no-op when the package corresponds to
// the file's own namespace.
func (im *importMap) addPackage(pkg string) {
	if pkg == "" {
		return
	}
	ns := qualifyNamespace(im.namespaceRoot, formatPackageName(pkg))
	if ns == im.currentNs {
		return
	}
	im.namespaces[ns] = struct{}{}
}

// addNamespace records a `using` for a fully-qualified namespace (e.g.
// "System.Text.Json.Serialization"). No prefixing is performed.
func (im *importMap) addNamespace(ns string) {
	if ns == "" || ns == im.currentNs {
		return
	}
	im.namespaces[ns] = struct{}{}
}

// String renders the collected imports as a sorted block of `using`
// statements terminated by a trailing newline. Returns the empty
// string when no imports were recorded.
func (im *importMap) String() string {
	if len(im.namespaces) == 0 {
		return ""
	}

	statements := make([]string, 0, len(im.namespaces))
	for ns := range im.namespaces {
		statements = append(statements, "using "+ns+";")
	}
	sort.Strings(statements)

	return strings.Join(statements, "\n") + "\n"
}

// qualifyNamespace prefixes a schema-package import with the configured
// namespace root, leaving BCL namespaces (System*) untouched.
func qualifyNamespace(namespaceRoot string, pkg string) string {
	if pkg == "" {
		return pkg
	}
	if pkg == "System" || strings.HasPrefix(pkg, "System.") {
		return pkg
	}
	if namespaceRoot == "" {
		return pkg
	}
	return namespaceRoot + "." + pkg
}
