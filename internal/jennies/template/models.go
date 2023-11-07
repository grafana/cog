package template

import "github.com/grafana/cog/internal/ast"

type ValueType string

const (
	ValueTypeConstant  = "constant"
	ValueTypeAssigment = "assigment"
	ValueTypeEnvelope  = "envelope"
)

type Tmpl struct {
	Package        string
	BuilderName    string
	ObjectName     string
	Imports        ImportMap
	ImportAlias    string
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
	Path           string
	InitSafeguards []string
	Value          string
	IsBuilder      bool
	IsNullable     bool
	Method         ast.AssignmentMethod
	ValueType      ValueType
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
