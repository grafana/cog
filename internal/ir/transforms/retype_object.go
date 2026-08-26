package transforms

import (
	"fmt"

	"github.com/grafana/cog/internal/ir"
)

var _ Transform = (*RetypeObject)(nil)

type RetypeObject struct {
	Object   ObjectReference
	As       ir.Type
	Comments []string
}

func (pass *RetypeObject) Process(schemas ir.Schemas) (ir.Schemas, error) {
	visitor := &Visitor{
		OnObject: pass.processObject,
	}

	return visitor.VisitSchemas(schemas)
}

func (pass *RetypeObject) processObject(_ *Visitor, _ *ir.Schema, object ir.Object) (ir.Object, error) {
	if !pass.Object.Matches(object) {
		return object, nil
	}

	trailMessage := fmt.Sprintf("RetypeObject[%s → %s]", ir.TypeName(object.Type), ir.TypeName(pass.As))

	object.Type = pass.As
	object.AddToPassesTrail(trailMessage)

	if pass.Comments != nil {
		object.Comments = pass.Comments
	}

	return object, nil
}
