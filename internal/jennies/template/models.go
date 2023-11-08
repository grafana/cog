package template

import (
	"github.com/grafana/cog/internal/ast"
)

type Tmpl struct {
	Package        string
	BuilderName    string
	ObjectName     string
	Imports        ImportMap
	ImportAlias    string // alias to the pkg in which the object being built lives.
	Options        []Option
	Constructor    Constructor
	DefaultBuilder DefaultBuilder
}

type Constructor struct {
	Args        []Argument
	Assignments []Assignment
}

type Option struct {
	Name        string
	Comments    []string
	Args        []Argument
	Assignments []Assignment
}

type Argument struct {
	Name          string
	Type          string
	ReferredAlias string
	ReferredName  string
}

type Assignment struct {
	Path           ast.Path
	InitSafeguards []string
	Method         ast.AssignmentMethod
	Value          ast.AssignmentValue
	Constraints    []Constraint
}

type Constraint struct {
	Name     string
	Op       ast.Op
	Arg      any
	IsString bool
}

type DefaultBuilder struct {
	Name string
	Args []Argument
}
