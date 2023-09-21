package builder

import (
	"github.com/grafana/cog/internal/ast"
)

type Selector func(builder ast.Builder) bool

func ByName(objectName string) Selector {
	return func(builder ast.Builder) bool {
		return builder.For.Name == objectName
	}
}

func StructGeneratedFromDisjunction() Selector {
	return func(builder ast.Builder) bool {
		if builder.For.Type.Kind != ast.KindStruct {
			return false
		}

		return builder.For.Type.AsStruct().Hint[ast.HintDisjunctionOfScalars] != nil ||
			builder.For.Type.AsStruct().Hint[ast.HintDiscriminatedDisjunctionOfRefs] != nil
	}
}

func EveryBuilder() Selector {
	return func(builder ast.Builder) bool {
		return true
	}
}
