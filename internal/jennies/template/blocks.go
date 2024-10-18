package template

import (
	"fmt"

	"github.com/grafana/cog/internal/ast"
)

func CustomObjectUnmarshalBlock(obj ast.Object) string {
	return fmt.Sprintf("object_%s_%s_custom_unmarshal", obj.SelfRef.ReferredPkg, obj.SelfRef.ReferredType)
}
