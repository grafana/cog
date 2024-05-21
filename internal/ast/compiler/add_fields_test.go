package compiler

import (
	"testing"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/testutils"
	"github.com/stretchr/testify/require"
)

func TestAddFields(t *testing.T) {
	// Prepare test input
	schema := &ast.Schema{
		Package: "add_fields",
		Objects: testutils.ObjectsMap(
			ast.NewObject("add_fields", "AString", ast.String()),
			ast.NewObject("add_fields", "SomeObject", ast.NewStruct(
				ast.NewStructField("AString", ast.String()),
			)),
		),
	}
	expected := &ast.Schema{
		Package: "add_fields",
		Objects: testutils.ObjectsMap(
			ast.NewObject("add_fields", "AString", ast.String()),
			ast.NewObject("add_fields", "SomeObject", ast.NewStruct(
				ast.NewStructField("AString", ast.String()),
				ast.NewStructField("addedByPass", ast.Bool(), ast.PassesTrail("AddFields[created]")),
			)),
		),
	}

	pass := &AddFields{
		Object: ObjectReference{
			Package: schema.Package,
			Object:  "SomeObject",
		},
		Fields: []ast.StructField{ast.NewStructField("addedByPass", ast.Bool())},
	}

	// Run the compiler pass
	runPassOnSchema(t, pass, schema, expected)
}

func TestAddFields_withConflictingExistingField(t *testing.T) {
	// Prepare test input
	schema := &ast.Schema{
		Package: "add_fields",
		Objects: testutils.ObjectsMap(
			ast.NewObject("add_fields", "AString", ast.String()),
			ast.NewObject("add_fields", "SomeObject", ast.NewStruct(
				ast.NewStructField("foo", ast.Bool()),
			)),
		),
	}

	pass := &AddFields{
		Object: ObjectReference{
			Package: schema.Package,
			Object:  "SomeObject",
		},
		Fields: []ast.StructField{ast.NewStructField("foo", ast.String())},
	}

	// Run the compiler pass
	runPassOnSchema(t, pass, schema, schema)
}

func TestAddFields_withUnknownObjectRef(t *testing.T) {
	// Prepare test input
	schema := &ast.Schema{
		Package: "add_fields",
		Objects: testutils.ObjectsMap(
			ast.NewObject("add_fields", "AString", ast.String()),
			ast.NewObject("add_fields", "SomeObject", ast.NewStruct(
				ast.NewStructField("AString", ast.String()),
			)),
		),
	}

	pass := &AddFields{
		Object: ObjectReference{
			Package: schema.Package,
			Object:  "DoesNotExist",
		},
		Fields: []ast.StructField{ast.NewStructField("foo", ast.String())},
	}

	// Run the compiler pass
	runPassOnSchema(t, pass, schema, schema)
}

func TestAddFields_withNonStructObjectRef(t *testing.T) {
	// Prepare test input
	schema := &ast.Schema{
		Package: "add_fields",
		Objects: testutils.ObjectsMap(
			ast.NewObject("add_fields", "AString", ast.String()),
		),
	}

	pass := &AddFields{
		Object: ObjectReference{
			Package: schema.Package,
			Object:  "AString",
		},
		Fields: []ast.StructField{ast.NewStructField("foo", ast.String())},
	}

	// Run the compiler pass
	_, err := pass.Process(ast.Schemas{schema})
	require.Error(t, err)
}
