package rust

import (
	"testing"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/testutils"
)

// TestRawTypes_Generate wires the golden-file harness for the Rust RawTypes
// jenny. The jenny itself is not implemented yet (Phase 3), so this test is
// skipped for now.
//
// Phase 3 un-skip checklist:
//   - Remove the t.Skip below.
//   - Construct the real RawTypes jenny (mirror python/rawtypes_test.go: build
//     Config, initTemplates, apiRefCollector) and run its compiler passes.
//   - Inside test.Run, unmarshal testutils.RawTypesIRInputFile, run the Rust
//     compiler passes, call jenny.Generate, then tc.WriteFiles(files).
//   - Add golden files under each fixture dir at testdata/jennies/rawtypes/<fixture>/RustRawTypes/.
//     The harness fails if the jenny emits a file with no matching golden file,
//     and fails if a golden file exists that the jenny did not emit. Generate
//     the initial set with COG_UPDATE_GOLDEN=true.
//
// The Name "RustRawTypes" matches the jenny's future JennyName(), mirroring how
// python uses "PythonRawTypes". Golden files live under
// ${TestDataRoot}/<fixture>/RustRawTypes.
func TestRawTypes_Generate(t *testing.T) {
	t.Skip("RawTypes jenny not implemented until Phase 3")

	test := testutils.GoldenFilesTestSuite[ast.Schema]{
		TestDataRoot: "../../../testdata/jennies/rawtypes",
		Name:         "RustRawTypes",
	}
	_ = test
}
