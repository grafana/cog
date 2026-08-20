package rust

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/grafana/codejen"
)

// rustEdition is the Rust edition rustfmt formats against. It matches the
// edition the generated crate targets.
const rustEdition = "2021"

// FormatRustFiles is a codejen postprocessor that runs rustfmt over the
// contents of every generated `.rs` file, letting rustfmt own all line
// wrapping and layout so the jennies can emit the simple single-line form
// everywhere.
//
// Non-`.rs` files pass through untouched. If rustfmt is not installed the file
// is returned unchanged rather than erroring, mirroring how the Go target's
// formatter is gated behind a flag: a missing formatter degrades to unformatted
// (but still valid) output rather than failing generation.
func FormatRustFiles(file codejen.File) (codejen.File, error) {
	if !strings.HasSuffix(file.RelativePath, ".rs") {
		return file, nil
	}

	formatted, err := formatRustSource(file.Data)
	if err != nil {
		return codejen.File{}, fmt.Errorf("rustfmt processing of generated file %q failed: %w", file.RelativePath, err)
	}

	return codejen.File{
		RelativePath: file.RelativePath,
		Data:         formatted,
		From:         file.From,
	}, nil
}

// formatRustSource pipes the given Rust source through `rustfmt` (reading from
// stdin, writing to stdout) and returns the formatted result. When rustfmt is
// not on the PATH the input is returned unchanged.
func formatRustSource(source []byte) ([]byte, error) {
	if _, err := exec.LookPath("rustfmt"); err != nil {
		return source, nil
	}

	cmd := exec.Command("rustfmt", "--edition", rustEdition, "--emit", "stdout")
	cmd.Stdin = bytes.NewReader(source)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("rustfmt failed: %w: %s", err, stderr.String())
	}

	return stdout.Bytes(), nil
}
