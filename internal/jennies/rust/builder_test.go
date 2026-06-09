package rust

import (
	"testing"

	"github.com/grafana/cog/internal/languages"
	"github.com/grafana/cog/internal/testutils"
)

// TestBuilder_Generate wires the golden-file harness for the Rust Builder
// jenny. The jenny itself is not implemented yet (Phase 4), so this test is
// skipped for now.
//
// Phase 4 un-skip checklist:
//   - Remove the t.Skip below.
//   - Construct the real Builder jenny (mirror python/builder_test.go: build the
//     Language via New(Config{}), initTemplates, apiRefCollector).
//   - Inside test.Run, unmarshal testutils.BuildersContextInputFile, run any
//     required language passes (e.g. languages.GenerateBuilderNilChecks), call
//     jenny.Generate, then tc.WriteFiles(files).
//   - Add golden files under each fixture dir at testdata/jennies/builders/<fixture>/RustBuilder/.
//     The harness fails if the jenny emits a file with no matching golden file,
//     and fails if a golden file exists that the jenny did not emit. Generate
//     the initial set with COG_UPDATE_GOLDEN=true.
//
// The Name "RustBuilder" matches the jenny's future JennyName(), mirroring how
// python uses "PythonBuilder". Golden files live under
// ${TestDataRoot}/<fixture>/RustBuilder.
func TestBuilder_Generate(t *testing.T) {
	t.Skip("Builder jenny not implemented until Phase 4")

	test := testutils.GoldenFilesTestSuite[languages.Context]{
		TestDataRoot: "../../../testdata/jennies/builders",
		Name:         "RustBuilder",
	}
	_ = test
}
