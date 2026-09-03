package rust

import (
	"sort"
	"strings"
)

// importMap collects the `use` statements required by a generated Rust module,
// keeping them deduplicated and rendered in a stable, idiomatic order.
type importMap struct {
	paths map[string]struct{}
}

func newImportMap() *importMap {
	return &importMap{
		paths: make(map[string]struct{}),
	}
}

// Add records a single fully-qualified `use` path, e.g. "serde::Serialize".
func (im *importMap) Add(path string) {
	if path == "" {
		return
	}
	im.paths[path] = struct{}{}
}

// String renders the collected imports as `use` statements. Paths sharing the
// serde crate root are grouped into a single braced statement, matching what
// rustfmt produces, so the output is fmt-clean without post-processing.
func (im *importMap) String() string {
	if len(im.paths) == 0 {
		return ""
	}

	serdeItems := make([]string, 0)
	var others []string

	for path := range im.paths {
		if item, ok := strings.CutPrefix(path, "serde::"); ok && !strings.Contains(item, "::") {
			serdeItems = append(serdeItems, item)
			continue
		}
		others = append(others, "use "+path+";")
	}

	var statements []string
	if len(serdeItems) != 0 {
		sort.Strings(serdeItems)
		if len(serdeItems) == 1 {
			statements = append(statements, "use serde::"+serdeItems[0]+";")
		} else {
			statements = append(statements, "use serde::{"+strings.Join(serdeItems, ", ")+"};")
		}
	}

	statements = append(statements, others...)
	// `serde` items are grouped first (rustfmt style), then all other `use`
	// statements alphabetically. A single final sort keeps output deterministic.
	sort.Strings(statements)

	return strings.Join(statements, "\n")
}
