package transforms

import (
	"testing"

	"github.com/grafana/cog/internal/testutils"
	"github.com/grafana/cog/pkg/ir"
)

func TestRenameObject(t *testing.T) {
	// Prepare test input
	schema := &ir.Schema{
		Package: "rename_object",
		Objects: testutils.ObjectsMap(
			ir.NewObject("rename_object", "SomeObject", ir.NewStruct(
				ir.NewStructField("foo", ir.String(ir.Nullable())),
				ir.NewStructField("ref_to_nice_object", ir.NewRef("rename_object", "NotANiceName")),
			)),
			ir.NewObject("rename_object", "NotANiceName", ir.NewStruct(
				ir.NewStructField("AString", ir.String(ir.Nullable())),
			)),
		),
	}
	expected := &ir.Schema{
		Package: "rename_object",
		Objects: testutils.ObjectsMap(
			ir.NewObject("rename_object", "SomeObject", ir.NewStruct(
				ir.NewStructField("foo", ir.String(ir.Nullable())),
				ir.NewStructField("ref_to_nice_object", ir.NewRef("rename_object", "ReallyNiceName")),
			)),
			ir.NewObject("rename_object", "ReallyNiceName", ir.NewStruct(
				ir.NewStructField("AString", ir.String(ir.Nullable())),
			), "RenameObject[NotANiceName → ReallyNiceName]"),
		),
	}

	pass := &RenameObject{
		From: ObjectReference{Package: schema.Package, Object: "NotANiceName"},
		To:   "ReallyNiceName",
	}

	// Run the compiler pass
	runPassOnSchema(t, pass, schema, expected)
}
