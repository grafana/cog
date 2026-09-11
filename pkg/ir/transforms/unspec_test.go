package transforms

import (
	"testing"

	"github.com/grafana/cog/internal/testutils"
	"github.com/grafana/cog/pkg/ir"
)

func TestUnspec(t *testing.T) {
	// Prepare test input
	schemas := ir.Schemas{
		&ir.Schema{
			Package: "without_spec",
			Objects: testutils.ObjectsMap(
				ir.NewObject("without_spec", "NotAStruct", ir.String()),

				ir.NewObject("without_spec", "AStruct", ir.NewStruct(
					ir.NewStructField("AString", ir.String()),
				)),
			),
		},

		&ir.Schema{
			Package: "with_spec_no_meta_id",
			Objects: testutils.ObjectsMap(
				ir.NewObject("with_spec_no_meta_id", "Metadata", ir.NewStruct(
					ir.NewStructField("SomeMeta", ir.String()),
				)),
				ir.NewObject("with_spec_no_meta_id", "Spec", ir.NewStruct(
					ir.NewStructField("title", ir.String()),
				)),
			),
		},

		&ir.Schema{
			Package: "with_spec_and_meta_id",
			Metadata: ir.SchemaMeta{
				Identifier: "Dashboard",
			},
			Objects: testutils.ObjectsMap(
				ir.NewObject("with_spec_and_meta_id", "Metadata", ir.NewStruct(
					ir.NewStructField("SomeMeta", ir.String()),
				)),
				ir.NewObject("with_spec_and_meta_id", "Spec", ir.NewStruct(
					ir.NewStructField("title", ir.String()),
				)),
			),
		},
	}

	// Prepare expected output
	expected := ir.Schemas{
		// Unchanged
		schemas[0],

		// No identifier defined in schema metadata: the package is used as name instead of "Spec"
		&ir.Schema{
			Package: "with_spec_no_meta_id",
			Objects: testutils.ObjectsMap(
				ir.NewObject("with_spec_no_meta_id", "with_spec_no_meta_id", ir.NewStruct(
					ir.NewStructField("title", ir.String()),
				), "Unspec[Spec → with_spec_no_meta_id]"),
			),
		},

		// Identifier defined in the schema metadata: it's used as object name instead of "Spec"
		&ir.Schema{
			Package: "with_spec_and_meta_id",
			Metadata: ir.SchemaMeta{
				Identifier: "Dashboard",
			},
			Objects: testutils.ObjectsMap(
				ir.NewObject("with_spec_and_meta_id", "Dashboard", ir.NewStruct(
					ir.NewStructField("title", ir.String()),
				), "Unspec[Spec → Dashboard]"),
			),
		},
	}

	// Run the compiler pass
	runPassOnSchemas(t, &Unspec{}, schemas, expected)
}
