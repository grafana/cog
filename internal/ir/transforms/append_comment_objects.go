package transforms

import (
	"fmt"

	"github.com/grafana/cog/internal/ir"
)

var _ Transform = (*AppendCommentObjects)(nil)

// AppendCommentObjects appends the given comment to every object definition.
type AppendCommentObjects struct {
	Comment string
}

func (pass *AppendCommentObjects) Process(schemas ir.Schemas) (ir.Schemas, error) {
	visitor := &Visitor{
		OnObject: pass.processObject,
	}

	return visitor.VisitSchemas(schemas)
}

func (pass *AppendCommentObjects) processObject(_ *Visitor, _ *ir.Schema, object ir.Object) (ir.Object, error) {
	object.Comments = append(object.Comments, pass.Comment)
	object.AddToPassesTrail(fmt.Sprintf("AppendCommentObjects[%s]", pass.Comment))

	return object, nil
}
