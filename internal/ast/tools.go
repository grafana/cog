package ast

import (
	"github.com/grafana/cog/internal/tools"
)

func TypeName(typeDef Type) string {
	if typeDef.Kind == KindRef {
		return tools.UpperCamelCase(typeDef.AsRef().ReferredType)
	}
	if typeDef.Kind == KindScalar {
		return tools.UpperCamelCase(string(typeDef.AsScalar().ScalarKind))
	}
	if typeDef.Kind == KindArray {
		return "ArrayOf" + TypeName(typeDef.AsArray().ValueType)
	}

	return tools.UpperCamelCase(string(typeDef.Kind))
}
