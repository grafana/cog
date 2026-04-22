package compiler

import (
	"github.com/grafana/cog/internal/ast"
)

var _ Pass = (*DefaultAsTyped)(nil)

// DefaultAsTyped converts the raw `Default any` field on every ast.Type into a
// typed *ast.TypeDefault, populating TypedDefault alongside the existing Default.
// Jennies can then access TypedDefault instead of performing unsafe type assertions
// on Default.
type DefaultAsTyped struct{}

func (pass *DefaultAsTyped) Process(schemas []*ast.Schema) ([]*ast.Schema, error) {
	visitor := &Visitor{
		OnObject: pass.processObject,
	}

	return visitor.VisitSchemas(schemas)
}

func (pass *DefaultAsTyped) processObject(_ *Visitor, _ *ast.Schema, object ast.Object) (ast.Object, error) {
	object.Type = processTypeDefaults(object.Type)
	return object, nil
}

// processTypeDefaults recursively sets TypedDefault on any Type that has
// Default != nil, and recurses into child types.
func processTypeDefaults(t ast.Type) ast.Type {
	if t.Default != nil {
		t.TypedDefault = ast.AnyToTypedDefault(t.Default)
	}

	switch {
	case t.IsStruct():
		for i, field := range t.Struct.Fields {
			field.Type = processTypeDefaults(field.Type)
			t.Struct.Fields[i] = field
		}
	case t.IsIntersection():
		for i, branch := range t.Intersection.Branches {
			t.Intersection.Branches[i] = processTypeDefaults(branch)
		}
	case t.IsDisjunction():
		for i, branch := range t.Disjunction.Branches {
			t.Disjunction.Branches[i] = processTypeDefaults(branch)
		}
	case t.IsArray():
		t.Array.ValueType = processTypeDefaults(t.Array.ValueType)
	case t.IsMap():
		t.Map.ValueType = processTypeDefaults(t.Map.ValueType)
	}

	return t
}

