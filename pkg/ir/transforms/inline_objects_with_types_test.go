package transforms

import (
	"testing"

	"github.com/grafana/cog/internal/testutils"
	"github.com/grafana/cog/pkg/ir"
)

func TestInlineScalarAliases(t *testing.T) {
	// Prepare test input
	schema := &ir.Schema{
		Package: "inline_scalar_aliases",
		Objects: testutils.ObjectsMap(
			ir.NewObject("inline_scalar_aliases", "AliasToString", ir.String()),
			ir.NewObject("inline_scalar_aliases", "AliasToMap", ir.NewMap(ir.String(), ir.Any())),
			ir.NewObject("inline_scalar_aliases", "AliasToArray", ir.NewArray(ir.String())),
			ir.NewObject("inline_scalar_aliases", "Constant", ir.String(ir.Value("foo"))),
			ir.NewObject("inline_scalar_aliases", "SomeObject", ir.NewStruct(
				ir.NewStructField("aliasToString", ir.NewRef("inline_scalar_aliases", "AliasToString")),
				ir.NewStructField("aliasToMap", ir.NewRef("inline_scalar_aliases", "AliasToMap")),
				ir.NewStructField("aliasToArray", ir.NewRef("inline_scalar_aliases", "AliasToArray")),
			)),
		),
	}
	expected := &ir.Schema{
		Package: "inline_scalar_aliases",
		Objects: testutils.ObjectsMap(
			ir.NewObject("inline_scalar_aliases", "Constant", ir.String(ir.Value("foo"))),
			ir.NewObject("inline_scalar_aliases", "SomeObject", ir.NewStruct(
				ir.NewStructField("aliasToString", ir.String(ir.Trail("InlineObjectsWithTypes[original=inline_scalar_aliases.AliasToString]"))),
				ir.NewStructField("aliasToMap", ir.NewMap(ir.String(), ir.Any(), ir.Trail("InlineObjectsWithTypes[original=inline_scalar_aliases.AliasToMap]"))),
				ir.NewStructField("aliasToArray", ir.NewArray(ir.String(), ir.Trail("InlineObjectsWithTypes[original=inline_scalar_aliases.AliasToArray]"))),
			)),
		),
	}

	pass := &InlineObjectsWithTypes{
		InlineTypes: []ir.Kind{ir.KindScalar, ir.KindArray, ir.KindMap},
	}

	// Run the compiler pass
	runPassOnSchema(t, pass, schema, expected)
}
