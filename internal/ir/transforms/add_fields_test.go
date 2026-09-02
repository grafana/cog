package transforms

import (
	"testing"

	"github.com/grafana/cog/internal/ir"
	"github.com/grafana/cog/internal/testutils"
	"github.com/stretchr/testify/require"
)

func TestAddFields(t *testing.T) {
	// Prepare test input
	schema := &ir.Schema{
		Package: "add_fields",
		Objects: testutils.ObjectsMap(
			ir.NewObject("add_fields", "AString", ir.String()),
			ir.NewObject("add_fields", "SomeObject", ir.NewStruct(
				ir.NewStructField("AString", ir.String()),
			)),
		),
	}
	expected := &ir.Schema{
		Package: "add_fields",
		Objects: testutils.ObjectsMap(
			ir.NewObject("add_fields", "AString", ir.String()),
			ir.NewObject("add_fields", "SomeObject", ir.NewStruct(
				ir.NewStructField("AString", ir.String()),
				ir.NewStructField("addedByPass", ir.Bool(), ir.PassesTrail("AddFields[created]")),
			)),
		),
	}

	pass := &AddFields{
		Object: ObjectReference{
			Package: schema.Package,
			Object:  "SomeObject",
		},
		Fields: []ir.StructField{ir.NewStructField("addedByPass", ir.Bool())},
	}

	// Run the compiler pass
	runPassOnSchema(t, pass, schema, expected)
}

func TestAddFields_withConflictingExistingField(t *testing.T) {
	// Prepare test input
	schema := &ir.Schema{
		Package: "add_fields",
		Objects: testutils.ObjectsMap(
			ir.NewObject("add_fields", "AString", ir.String()),
			ir.NewObject("add_fields", "SomeObject", ir.NewStruct(
				ir.NewStructField("foo", ir.Bool()),
			)),
		),
	}

	pass := &AddFields{
		Object: ObjectReference{
			Package: schema.Package,
			Object:  "SomeObject",
		},
		Fields: []ir.StructField{ir.NewStructField("foo", ir.String())},
	}

	// Run the compiler pass
	runPassOnSchema(t, pass, schema, schema)
}

func TestAddFields_withUnknownObjectRef(t *testing.T) {
	// Prepare test input
	schema := &ir.Schema{
		Package: "add_fields",
		Objects: testutils.ObjectsMap(
			ir.NewObject("add_fields", "AString", ir.String()),
			ir.NewObject("add_fields", "SomeObject", ir.NewStruct(
				ir.NewStructField("AString", ir.String()),
			)),
		),
	}

	pass := &AddFields{
		Object: ObjectReference{
			Package: schema.Package,
			Object:  "DoesNotExist",
		},
		Fields: []ir.StructField{ir.NewStructField("foo", ir.String())},
	}

	// Run the compiler pass
	runPassOnSchema(t, pass, schema, schema)
}

func TestAddFields_withNonStructObjectRef(t *testing.T) {
	// Prepare test input
	schema := &ir.Schema{
		Package: "add_fields",
		Objects: testutils.ObjectsMap(
			ir.NewObject("add_fields", "AString", ir.String()),
		),
	}

	pass := &AddFields{
		Object: ObjectReference{
			Package: schema.Package,
			Object:  "AString",
		},
		Fields: []ir.StructField{ir.NewStructField("foo", ir.String())},
	}

	// Run the compiler pass
	_, err := pass.Process(ir.Schemas{schema})
	require.Error(t, err)
}
