package jsonschema

import (
	"strings"

	"github.com/grafana/cog/internal/ast"
	schemaparser "github.com/santhosh-tekuri/jsonschema/v5"
)

func schemaComments(schema *schemaparser.Schema) []string {
	comment := schema.Description

	lines := strings.Split(comment, "\n")
	filtered := make([]string, 0, len(lines))

	for _, line := range lines {
		if line == "" {
			continue
		}

		filtered = append(filtered, line)
	}

	return filtered
}

func extractMetadata(schema *schemaparser.Schema) (ast.SchemaMeta, error) {
	if m, ok := schema.Extensions["metadata"]; ok {
		if err := m.Validate(schemaparser.ValidationContext{}, m); err != nil {
			return ast.SchemaMeta{}, err
		}
		metadata, ok := m.(metadataValidator)
		if !ok {
			return ast.SchemaMeta{}, nil
		}

		return ast.SchemaMeta{
			Kind:       ast.SchemaKind(metadata[MetadataKind]),
			Variant:    ast.SchemaVariant(metadata[MetadataVariant]),
			Identifier: metadata[MetadataIdentifier],
		}, nil
	}

	return ast.SchemaMeta{}, nil
}
