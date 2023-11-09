package template

import (
	"github.com/grafana/cog/internal/ast"
)

type Builder struct {
	Package     string
	BuilderName string
	ObjectName  string
	Imports     ImportMap
	ImportAlias string // alias to the pkg in which the object being built lives.
	Constructor Constructor
	Options     []Option
	Defaults    []OptionCall
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
	Constraints    []Constraint
	Method         ast.AssignmentMethod
	Value          ast.AssignmentValue
}

type Constraint struct {
	ArgName   string
	Op        ast.Op
	Parameter any
}

type OptionCall struct {
	OptionName string
	Args       []string
}
