package rust

import (
	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/testutils"
)

// formatRustGoldenFiles runs the FormatRustFiles postprocessor (rustfmt) over
// each generated file so the golden-file tests compare against rustfmt-formatted
// output. The harness itself does not run postprocessors, so the jennies emit
// the simple single-line form and rely on this to reproduce the same formatting
// the real generation pipeline applies via AddPostprocessors(FormatRustFiles).
func formatRustGoldenFiles[T any](tc *testutils.Test[T], files codejen.Files) codejen.Files {
	tc.Helper()

	formatted := make(codejen.Files, 0, len(files))
	for i := range files {
		out, err := FormatRustFiles(files[i])
		if err != nil {
			tc.Fatalf("rustfmt formatting of %q failed: %s", files[i].RelativePath, err)
		}
		formatted = append(formatted, out)
	}

	return formatted
}
