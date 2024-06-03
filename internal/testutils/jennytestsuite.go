package testutils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/envvars"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type FileComparator func(t *testing.T, expectedContent []byte, gotContent []byte, filename string)

// TrimSpacesDiffComparator compares the content of two files after trimming leading and trailing spaces.
func TrimSpacesDiffComparator(t *testing.T, expectedContent []byte, gotContent []byte, filename string) {
	t.Helper()

	cleanedGoldenContent := strings.TrimSpace(string(expectedContent))
	cleanedGeneratedContent := strings.TrimSpace(string(gotContent))

	assert.Equal(t, cleanedGoldenContent, cleanedGeneratedContent, "result for %s differs", filename)
}

// A GoldenFilesTestSuite represents a test run that processes inputs from a
// given directory and compares results against "golden files". See the [Test]
// documentation for more details.
type GoldenFilesTestSuite[GoldenFileType any] struct {
	// Run the test suite on this directory.
	TestDataRoot string

	// Name is a unique name for this test. The golden files for this test are
	// expected to live under `${TestDataRoot}/[test]/${Name}`.
	Name string

	// Skip is a map of tests to skip; the key is the test name; the value is the
	// skip message.
	Skip map[string]string

	// FileComparator is called to compare a golden file with its test-generated equivalent.
	// TrimSpacesDiffComparator is used if no comparator is provided.
	FileComparator FileComparator
}

// A Test represents a single test.
//
// A Test embeds *[testing.T] and should be used to report errors.
//
// When the test function has returned, output written with [Test.Write], [Test.Writer],
// [Test.WriteFile] and friends is checked against the expected output files.
//
// If the output differs and $COG_UPDATE_GOLDEN is non-empty, the txtar file will be
// updated and written to disk with the actual output data replacing the
// out files.
type Test[GoldenFileType any] struct {
	// Allow Test to be used as a T.
	*testing.T

	fileComparator FileComparator

	outFiles map[string][]byte

	// RootDir is the path of the current test directory.
	RootDir string

	// OutputDir is the path where golden files live for this test.
	OutputDir string
}

// WriteFile writes a [codejen.File] to the main output.
func (t *Test[GoldenFileType]) WriteFile(f *codejen.File) {
	t.outFiles[f.RelativePath] = f.Data
}

// WriteJSON marshals and writes the given input to `filename`.
func (t *Test[GoldenFileType]) WriteJSON(filename string, input any) {
	t.Helper()

	marshaledIR, err := json.MarshalIndent(input, "", "  ")
	require.NoError(t, err)

	t.WriteFile(&codejen.File{
		RelativePath: filename,
		Data:         marshaledIR,
	})
}

// WriteFiles writes a list of [codejen.File] to the main output.
func (t *Test[GoldenFileType]) WriteFiles(files codejen.Files) {
	for i := range files {
		t.WriteFile(&files[i])
	}
}

// UnmarshalJSONInput reads and unmarshals the specified input file.
func (t *Test[GoldenFileType]) UnmarshalJSONInput(filename string) GoldenFileType {
	t.Helper()

	var parsed GoldenFileType
	if err := json.Unmarshal(t.ReadInput(filename), &parsed); err != nil {
		t.Fatalf("could not unmarshal input: %s", err)
	}

	return parsed
}

// OpenInput opens the specified input file and returns an [io.Reader] to it.
func (t *Test[GoldenFileType]) OpenInput(inputFile string) io.Reader {
	t.Helper()

	reader, err := os.Open(filepath.Join(t.RootDir, inputFile))
	if err != nil {
		t.Fatalf("could not open input file '%s': %s", inputFile, err)
	}

	return reader
}

// ReadInput reads the specified input file and its contents.
func (t *Test[GoldenFileType]) ReadInput(inputFile string) []byte {
	t.Helper()

	content, err := io.ReadAll(t.OpenInput(inputFile))
	if err != nil {
		t.Fatalf("could not read input file '%s': %s", inputFile, err)
	}

	return content
}

func (t *Test[GoldenFileType]) writeGoldenFiles() error {
	for filename, content := range t.outFiles {
		fullpath := filepath.Join(t.OutputDir, filename)

		if err := os.MkdirAll(filepath.Dir(fullpath), 0755); err != nil {
			return err
		}

		if err := os.WriteFile(fullpath, content, 0600); err != nil {
			return err
		}
	}

	return nil
}

// Run runs tests defined in `x.TestDataRoot`.
func (x *GoldenFilesTestSuite[GoldenFileType]) Run(t *testing.T, f func(tc *Test[GoldenFileType])) {
	t.Helper()

	entries, err := os.ReadDir(x.TestDataRoot)
	if err != nil {
		t.Fatal(err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		testRootDir := filepath.Join(x.TestDataRoot, entry.Name())
		outputDir := filepath.Join(testRootDir, x.Name)
		testName := filepath.Join(x.Name, entry.Name())

		t.Run(testName, func(t *testing.T) {
			tc := &Test[GoldenFileType]{
				T:              t,
				RootDir:        testRootDir,
				OutputDir:      outputDir,
				fileComparator: x.FileComparator,
				outFiles:       make(map[string][]byte),
			}

			if tc.fileComparator == nil {
				tc.fileComparator = TrimSpacesDiffComparator
			}

			if msg, ok := x.Skip[entry.Name()]; ok {
				t.Skip(msg)
			}

			f(tc)

			// check that the files generated by the test match the golden files
			for generatedFile, generatedContent := range tc.outFiles {
				goldenFilePath := filepath.Join(outputDir, generatedFile)

				_, err := os.Stat(goldenFilePath)
				if errors.Is(err, os.ErrNotExist) {
					if !envvars.UpdateGoldenFiles {
						assert.Fail(t, fmt.Sprintf("jenny generated file '%s', but it does not exist in the golden files", generatedFile))
					}
					continue
				} else if err != nil {
					t.Fatalf("error while checking golden file '%s': %s", goldenFilePath, err)
				}

				goldenContent, err := os.ReadFile(goldenFilePath)
				if err != nil {
					t.Fatalf("error while reading golden file '%s': %s", goldenFilePath, err)
				}

				tc.fileComparator(t, goldenContent, generatedContent, generatedFile)
			}

			// make sure that there is no golden file that wasn't covered by the test.
			err := filepath.WalkDir(outputDir, func(path string, d fs.DirEntry, err error) error {
				if path == outputDir {
					return nil
				}

				if d.IsDir() {
					return nil
				}

				relativeGoldenFilePath := strings.TrimLeft(strings.TrimPrefix(path, outputDir), "/")
				if _, found := tc.outFiles[relativeGoldenFilePath]; !found {
					if envvars.UpdateGoldenFiles {
						if err := os.Remove(path); err != nil {
							t.Fatalf("could not remove golden file: %s", err)
						}
					} else {
						assert.Fail(t, fmt.Sprintf("golden file '%s' exists but was not generated by jenny", relativeGoldenFilePath))
					}
				}

				return nil
			})
			if err != nil {
				t.Fatal(err.Error())
			}

			// If we need to: re-generate golden files
			if envvars.UpdateGoldenFiles {
				if err := tc.writeGoldenFiles(); err != nil {
					t.Fatalf("could not write golden files: %s", err)
				}
			}
		})
	}
}
