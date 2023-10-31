package jsonschema

import (
	"fmt"
	"github.com/grafana/cog/internal/ast"
	schemaparser "github.com/santhosh-tekuri/jsonschema/v5"
	"strings"
)

const (
	MetadataMetadata   = "metadata"
	MetadataIdentifier = "identifier"
	MetadataKind       = "kind"
	MetadataVariant    = "variant"
)

func metadata() *schemaparser.Schema {
	return schemaparser.MustCompileString("metadata.json", `{
	"metadata": {
		"type": "object",
		"properties": {
			"identifier": {
				"type": "string"
			},
			"kind": {
				"type": "string"
			},
			"variant": {
				"type": "string"
			}
		}
	}
}`)
}

type metadataCompiler struct{}

func (metadataCompiler) Compile(ctx schemaparser.CompilerContext, m map[string]interface{}) (schemaparser.ExtSchema, error) {
	if pow, ok := m["metadata"]; ok {
		if v, ok := pow.(map[string]interface{}); ok {
			md := metadataValidator{}
			for key, value := range v {
				md[key] = fmt.Sprintf("%s", value)
			}
			return md, nil
		}
		return nil, nil
	}

	// nothing to compile, return nil
	return nil, nil
}

type metadataValidator map[string]string

func (m metadataValidator) Validate(ctx schemaparser.ValidationContext, v interface{}) error {
	if err := m.validateKeyword(MetadataIdentifier); err != nil {
		return err
	}

	if err := m.validateKeyword(MetadataKind); err != nil {
		return err
	}

	if err := m.validateValues(MetadataKind, string(ast.SchemaKindComposable), string(ast.SchemaKindCore)); err != nil {
		return err
	}

	if m[MetadataKind] == string(ast.SchemaKindComposable) {
		if err := m.validateKeyword(MetadataVariant); err != nil {
			return err
		}
		if err := m.validateValues(MetadataVariant, string(ast.SchemaVariantPanel), string(ast.SchemaVariantDataQuery)); err != nil {
			return err
		}
	}

	return nil
}

func (m metadataValidator) validateKeyword(key string) error {
	keyword, ok := m[key]
	if !ok {
		return &schemaparser.ValidationError{
			KeywordLocation: MetadataMetadata,
			Message:         fmt.Sprintf("missing '%s' keyword", key),
		}
	}

	if keyword == "" {
		return &schemaparser.ValidationError{
			KeywordLocation: MetadataMetadata,
			Message:         fmt.Sprintf("'%s' cannot be empty", key),
		}
	}

	return nil
}

func (m metadataValidator) validateValues(key string, values ...string) error {
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

	return &schemaparser.ValidationError{
		KeywordLocation: MetadataMetadata,
		Message:         fmt.Sprintf("%s only accepts values: '%s'", key, strings.Join(values, ", ")),
	}
}
