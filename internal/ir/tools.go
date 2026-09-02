package ir

import (
	"github.com/grafana/cog/internal/tools"
)

// TypeName returns a human-readable name for the given type definition, based on its kind.
func TypeName(typeDef Type) string {
	if typeDef.IsRef() {
		return tools.UpperCamelCase(typeDef.AsRef().ReferredType)
	}
	if typeDef.IsScalar() {
		return tools.UpperCamelCase(string(typeDef.AsScalar().ScalarKind))
	}
	if typeDef.IsArray() {
		return "ArrayOf" + TypeName(typeDef.AsArray().ValueType)
	}

	return tools.UpperCamelCase(string(typeDef.Kind))
}
