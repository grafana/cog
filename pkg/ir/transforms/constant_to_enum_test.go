package transforms

import (
	"testing"

	"github.com/grafana/cog/internal/testutils"
	"github.com/grafana/cog/pkg/ir"
)

func TestConstantToEnum(t *testing.T) {
	// Prepare test input
	strAlias := ir.NewObject("sandbox", "String", ir.String())
	strConstant := ir.NewObject("sandbox", "Mode", ir.String(ir.Value("auto")))
	notTargetedStrConstant := ir.NewObject("sandbox", "Alignment", ir.String(ir.Value("center")))
	intConstant := ir.NewObject("sandbox", "DefaultSize", ir.NewScalar(ir.KindInt32, ir.Value(42)))
	obj := ir.NewObject("sandbox", "Obj", ir.NewStruct(ir.NewStructField("foo", ir.String())))
	schema := &ir.Schema{
		Package: "sandbox",
		Objects: testutils.ObjectsMap(
			strAlias,
			strConstant,
			notTargetedStrConstant,
			intConstant,
			obj,
		),
	}

	newEnum := ir.NewObject("sandbox", "Mode", ir.NewEnum([]ir.EnumValue{
		{
			Type:  ir.String(),
			Name:  "auto",
			Value: "auto",
		},
	}))
	newEnum.AddToPassesTrail("ConstantToEnum")
	expected := &ir.Schema{
		Package: "sandbox",
		Objects: testutils.ObjectsMap(
			strAlias,
			newEnum,
			notTargetedStrConstant,
			intConstant,
			obj,
		),
	}

	// Run the compiler pass
	runPassOnSchema(t, &ConstantToEnum{
		Objects: []ObjectReference{
			{Package: "sandbox", Object: "String"},
			{Package: "sandbox", Object: "Mode"},
			{Package: "sandbox", Object: "DefaultSize"},
			{Package: "sandbox", Object: "Obj"},
		},
	}, schema, expected)
}
