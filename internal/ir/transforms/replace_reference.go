package transforms

import (
	"fmt"

	"github.com/grafana/cog/internal/ir"
)

// ReplaceReference replaces any usage of the `From` reference by the one given in `To`.
//
// Example:
//
//	```
//	From = { Package: common, Object: DataSourceRef }
//	To = { Package: common, Object: DataSourceDescriptor }
//	```
//
//	```
//	Panel: {
//		DataSource: common.DataSourceRef
//	}
//	```
//
// Will become:
//
//	```
//	Panel: {
//		DataSource: common.DataSourceDescriptor
//	}
//	```
type ReplaceReference struct {
	From     ObjectReference
	To       ObjectReference
	refFound bool
}

func (pass *ReplaceReference) Process(schemas ir.Schemas) (ir.Schemas, error) {
	pass.refFound = false

	visitor := Visitor{
		OnRef: pass.processRef,
	}

	return visitor.VisitSchemas(schemas)
}

func (pass *ReplaceReference) processRef(_ *Visitor, _ *ir.Schema, def ir.Type) (ir.Type, error) {
	if !pass.From.MatchesRef(def.AsRef()) {
		return def, nil
	}

	pass.refFound = true

	return ir.NewRef(pass.To.Package, pass.To.Object, ir.Trail(fmt.Sprintf("ReplaceReference[%s → %s]", def.Ref, pass.To))), nil
}

func (pass *ReplaceReference) Diagnostics() []string {
	if pass.refFound {
		return nil
	}

	return []string{
		fmt.Sprintf("reference '%s' not found", pass.From),
	}
}
