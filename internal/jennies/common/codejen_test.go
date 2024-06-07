package common

import (
	"fmt"
	"testing"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/languages"
	"github.com/stretchr/testify/require"
)

type fakeNamedJenny struct {
	Name string
}

func (jenny fakeNamedJenny) JennyName() string {
	return jenny.Name
}

func TestPrefixer(t *testing.T) {
	req := require.New(t)
	fileContent := []byte("with content")
	inputFile := codejen.NewFile("some.file", fileContent)

	resultFile, err := PathPrefixer("/the/prefix")(*inputFile)
	req.NoError(err)

	req.Equal("/the/prefix/some.file", resultFile.RelativePath)
	req.Equal(fileContent, resultFile.Data)
}

func TestSlashHeaderMapper(t *testing.T) {
	markdownContent := `# A document title

With some content`

	jsonContent := `{
  "status": "all green"
}`

	yamlContent := `status: "all green"`
	expectedYamlContent := fmt.Sprintf(`# Code generated - EDITING IS FUTILE. DO NOT EDIT.
#
# Using jennies:
#     SomeJenny

%s`, yamlContent)

	goContent := `type SomeType struct {
	Foo string
}`
	expectedGoContent := fmt.Sprintf(`// Code generated - EDITING IS FUTILE. DO NOT EDIT.
//
// Using jennies:
//     SomeJenny

%s`, goContent)

	tsContent := `export interface SomeType {
	foo: string;
}`
	expectedTSContent := fmt.Sprintf(`// Code generated - EDITING IS FUTILE. DO NOT EDIT.
//
// Using jennies:
//     SomeJenny

%s`, tsContent)
	expectedTSContentNoDebug := fmt.Sprintf(`// Code generated - EDITING IS FUTILE. DO NOT EDIT.

%s`, tsContent)

	namedJennies := []codejen.NamedJenny{
		fakeNamedJenny{Name: "SomeJenny"},
	}

	testCases := []struct {
		name            string
		inputFile       *codejen.File
		debug           bool
		expectedContent string
	}{
		{
			name:            "markdown file",
			inputFile:       codejen.NewFile("./dir/README.md", []byte(markdownContent), namedJennies...),
			expectedContent: markdownContent,
		},
		{
			name:            "json file",
			inputFile:       codejen.NewFile("./dir/data.json", []byte(jsonContent), namedJennies...),
			expectedContent: jsonContent,
		},
		{
			name:            "yaml file",
			debug:           true,
			inputFile:       codejen.NewFile("./dir/data.yaml", []byte(yamlContent), namedJennies...),
			expectedContent: expectedYamlContent,
		},
		{
			name:            "yml file",
			debug:           true,
			inputFile:       codejen.NewFile("./dir/data.yml", []byte(yamlContent), namedJennies...),
			expectedContent: expectedYamlContent,
		},
		{
			name:            "go file",
			debug:           true,
			inputFile:       codejen.NewFile("./dir/main.go", []byte(goContent), namedJennies...),
			expectedContent: expectedGoContent,
		},
		{
			name:            "ts type",
			debug:           true,
			inputFile:       codejen.NewFile("./dir/main.ts", []byte(tsContent), namedJennies...),
			expectedContent: expectedTSContent,
		},
		{
			name:            "ts type",
			debug:           false,
			inputFile:       codejen.NewFile("./dir/main.ts", []byte(tsContent), namedJennies...),
			expectedContent: expectedTSContentNoDebug,
		},
	}

	for _, testCase := range testCases {
		tc := testCase

		t.Run(tc.name, func(t *testing.T) {
			req := require.New(t)

			resultFile, err := GeneratedCommentHeader(languages.Config{Debug: tc.debug})(*tc.inputFile)
			req.NoError(err)

			req.Equal(tc.expectedContent, string(resultFile.Data))
		})
	}
}
