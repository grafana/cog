package transforms

import (
	"testing"

	"github.com/grafana/cog/internal/ir"
	"github.com/grafana/cog/internal/testutils"
)

func TestNameAnonymousStruct(t *testing.T) {
	// Prepare test input
	schema := &ir.Schema{
		Package: "name_anonymous_struct",
		Objects: testutils.ObjectsMap(
			ir.NewObject("name_anonymous_struct", "SomeObject", ir.NewStruct(
				ir.NewStructField("inner", ir.NewStruct(
					ir.NewStructField("title", ir.String()),
				)),
			)),
		),
	}
	expected := &ir.Schema{
		Package: "name_anonymous_struct",
		Objects: testutils.ObjectsMap(
			ir.NewObject("name_anonymous_struct", "SomeObject", ir.NewStruct(
				ir.NewStructField("inner", ir.NewRef(schema.Package, "Inner")),
			)),
			ir.NewObject("name_anonymous_struct", "Inner", ir.NewStruct(
				ir.NewStructField("title", ir.String()),
			), "NameAnonymousStruct"),
		),
	}

	pass := &NameAnonymousStruct{
		Field: FieldReference{Package: schema.Package, Object: "SomeObject", Field: "inner"},
		As:    "Inner",
	}

	// Run the compiler pass
	runPassOnSchema(t, pass, schema, expected)
}

func TestNameAnonymousStruct_onNonStructObject(t *testing.T) {
	// Prepare test input
	schema := &ir.Schema{
		Package: "name_anonymous_struct",
		Objects: testutils.ObjectsMap(
			ir.NewObject("name_anonymous_struct", "SomeObject", ir.NewStruct(
				ir.NewStructField("inner", ir.NewStruct(
					ir.NewStructField("title", ir.String()),
				)),
			)),
		),
	}

	pass := &NameAnonymousStruct{
		// no-op since `doesNotExist` does not exist
		Field: FieldReference{Package: schema.Package, Object: "SomeObject", Field: "doesNotExist"},
		As:    "Inner",
	}

	// Run the compiler pass
	runPassOnSchema(t, pass, schema, schema)
}

func TestNameAnonymousStruct_onNonStructField(t *testing.T) {
	// Prepare test input
	schema := &ir.Schema{
		Package: "name_anonymous_struct",
		Objects: testutils.ObjectsMap(
			ir.NewObject("name_anonymous_struct", "SomeObject", ir.NewStruct(
				ir.NewStructField("inner", ir.Bool()),
			)),
		),
	}

	pass := &NameAnonymousStruct{
		// no-op since `AString` is not a struct
		Field: FieldReference{Package: schema.Package, Object: "SomeObject", Field: "inner"},
		As:    "Inner",
	}

	// Run the compiler pass
	runPassOnSchema(t, pass, schema, schema)
}

func TestNameAnonymousStruct_onUnknownField(t *testing.T) {
	// Prepare test input
	schema := &ir.Schema{
		Package: "name_anonymous_struct",
		Objects: testutils.ObjectsMap(
			ir.NewObject("name_anonymous_struct", "AString", ir.String()),
		),
	}

	pass := &NameAnonymousStruct{
		// no-op since `AString` is not a struct
		Field: FieldReference{Package: schema.Package, Object: "AString", Field: "doesNotExist"},
		As:    "Inner",
	}

	// Run the compiler pass
	runPassOnSchema(t, pass, schema, schema)
}
