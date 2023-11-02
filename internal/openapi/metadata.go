package openapi

import (
	"fmt"
	"strings"

	"github.com/grafana/cog/internal/ast"
)

const (
	MetadataMetadata   = "x-metadata"
	MetadataIdentifier = "identifier"
	MetadataKind       = "kind"
	MetadataVariant    = "variant"
)

const errMetadataValidationFailed = "metadata validation failed: %s"

type metadata map[string]string

func (m metadata) extractMetadata() (ast.SchemaMeta, error) {
	identifier, err := m.validateKeyword(MetadataIdentifier)
	if err != nil {
		return ast.SchemaMeta{}, fmt.Errorf(errMetadataValidationFailed, err)
	}

	kind, err := m.validateKeyword(MetadataKind)
	if err != nil {
		return ast.SchemaMeta{}, fmt.Errorf(errMetadataValidationFailed, err)
	}

	if err = m.validateValue(MetadataKind, string(ast.SchemaKindCore), string(ast.SchemaKindComposable)); err != nil {
		return ast.SchemaMeta{}, fmt.Errorf(errMetadataValidationFailed, err)
	}

	variant := ""
	if kind == string(ast.SchemaKindComposable) {
		variant, err = m.validateKeyword(MetadataVariant)
		if err != nil {
			return ast.SchemaMeta{}, fmt.Errorf(errMetadataValidationFailed, err)
		}
		if err = m.validateValue(MetadataVariant, string(ast.SchemaVariantPanel), string(ast.SchemaVariantDataQuery)); err != nil {
			return ast.SchemaMeta{}, fmt.Errorf(errMetadataValidationFailed, err)
		}
	}

	return ast.SchemaMeta{
		Identifier: identifier,
		Kind:       ast.SchemaKind(kind),
		Variant:    ast.SchemaVariant(variant),
	}, nil
}

func (m metadata) validateKeyword(key string) (string, error) {
	keyword, ok := m[key]
	if !ok {
		return "", fmt.Errorf("missing '%s' keyword", key)
	}

	if keyword == "" {
		return "", fmt.Errorf("'%s' cannot be empty", key)
	}

	return keyword, nil
}

func (m metadata) validateValue(key string, values ...string) error {
	hasValidValue := false
	for _, v := range values {
		if m[key] == v {
			hasValidValue = true
			break
		}
	}

	if hasValidValue {
		return nil
	}

	return fmt.Errorf("%s only accepts values: '%s'", key, strings.Join(values, ", "))
}
