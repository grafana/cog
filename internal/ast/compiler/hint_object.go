package compiler

import (
	"fmt"
	"sort"
	"strings"

	"github.com/grafana/cog/internal/ast"
)

var _ Pass = (*HintObject)(nil)

type HintObject struct {
	Object      ObjectReference
	Hints       ast.JenniesHints
	objectFound bool
}

func (pass *HintObject) Process(schemas []*ast.Schema) ([]*ast.Schema, error) {
	pass.objectFound = false
	visitor := &Visitor{
		OnObject: pass.processObject,
	}

	return visitor.VisitSchemas(schemas)
}

func (pass *HintObject) processObject(_ *Visitor, _ *ast.Schema, object ast.Object) (ast.Object, error) {
	if !pass.Object.Matches(object) {
		return object, nil
	}

	pass.objectFound = true

	hintsTrail := make([]string, 0, len(pass.Hints))
	for hint, val := range pass.Hints {
		object.Type.Hints[hint] = val
		hintsTrail = append(hintsTrail, fmt.Sprintf("%s=%v", hint, val))
	}

	// to ensure a consistent trail
	sort.Strings(hintsTrail)

	object.AddToPassesTrail(fmt.Sprintf("HintObject[%s]", strings.Join(hintsTrail, ", ")))

	return object, nil
}

func (pass *HintObject) Diagnostics() []string {
	if pass.objectFound {
		return nil
	}

	return []string{
		fmt.Sprintf("object '%s' not found", pass.Object),
	}
}
