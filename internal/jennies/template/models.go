package template

import "github.com/grafana/cog/internal/ast"

type Tmpl struct {
	Package     string
	BuilderName string
	ObjectName  string
	Imports     ImportMap
	ImportAlias string
	Options     []Option
	Constructor Constructor
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
	Path           string
	InitSafeguards []string
	Value          string
	IsBuilder      bool
	Method         ast.AssignmentMethod
	Constraints    []Constraint
}

type Constraint struct {
	Name     string
	Op       ast.Op
	Arg      any
	IsString bool
}
