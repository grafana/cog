package compiler

import (
	"fmt"
	"strings"

	"github.com/grafana/cog/internal/ast"
)

var _ Pass = (*HintObject)(nil)

type HintObject struct {
	Object ObjectReference
	Hints  ast.JenniesHints
}

func (pass *HintObject) Process(schemas []*ast.Schema) ([]*ast.Schema, error) {
	for i, schema := range schemas {
		schemas[i] = pass.processSchema(schema)
	}

	return schemas, nil
}

func (pass *HintObject) processSchema(schema *ast.Schema) *ast.Schema {
	schema.Objects = schema.Objects.Map(func(_ string, object ast.Object) ast.Object {
		return pass.processObject(object)
	})

	return schema
}

func (pass *HintObject) processObject(object ast.Object) ast.Object {
	if !pass.Object.Matches(object) {
		return object
	}

	hintsTrail := make([]string, 0, len(pass.Hints))
	for hint, val := range pass.Hints {
		object.Type.Hints[hint] = val
		hintsTrail = append(hintsTrail, fmt.Sprintf("%s=%v", hint, val))
	}

	object.AddToPassesTrail(fmt.Sprintf("HintObject[%s]", strings.Join(hintsTrail, ", ")))

	return object
}
