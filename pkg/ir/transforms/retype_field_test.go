package transforms

import (
	"testing"

	"github.com/grafana/cog/internal/testutils"
	"github.com/grafana/cog/pkg/ir"
)

func TestRetypeField(t *testing.T) {
	// Prepare test input
	schema := &ir.Schema{
		Package: "retype_field",
		Objects: testutils.ObjectsMap(
			ir.NewObject("retype_field", "SomeObject", ir.NewStruct(
				ir.NewStructField("AString", ir.String()),
			)),
		),
	}
	expected := &ir.Schema{
		Package: "retype_field",
		Objects: testutils.ObjectsMap(
			ir.NewObject("retype_field", "SomeObject", ir.NewStruct(
				ir.NewStructField("AString", ir.Bool(), ir.PassesTrail("RetypeField[String → Bool]")),
			)),
		),
	}

	pass := &RetypeField{
		Field: FieldReference{Package: schema.Package, Object: "SomeObject", Field: "AString"},
		As:    ir.Bool(),
	}

	// Run the compiler pass
	runPassOnSchema(t, pass, schema, expected)
}

func TestRetypeField_notFoundFieldRef(t *testing.T) {
	// Prepare test input
	schema := &ir.Schema{
		Package: "retype_field",
		Objects: testutils.ObjectsMap(
			ir.NewObject("retype_field", "SomeObject", ir.NewStruct(
				ir.NewStructField("AString", ir.String()),
			)),
		),
	}

	pass := &RetypeField{
		// no-op since `SomeObject.NotFound` does not exist
		Field: FieldReference{Package: schema.Package, Object: "SomeObject", Field: "NotFound"},
		As:    ir.Bool(),
	}

	// Run the compiler pass
	runPassOnSchema(t, pass, schema, schema)
}

func TestRetypeField_onNonStruct(t *testing.T) {
	// Prepare test input
	schema := &ir.Schema{
		Package: "retype_field",
		Objects: testutils.ObjectsMap(
			ir.NewObject("retype_field", "AString", ir.String()),
		),
	}

	pass := &RetypeField{
		// no-op since `AString` is not a struct
		Field: FieldReference{Package: schema.Package, Object: "AString", Field: "NotAField"},
		As:    ir.Bool(),
	}

	// Run the compiler pass
	runPassOnSchema(t, pass, schema, schema)
}
