package template

import (
	"fmt"

	"github.com/grafana/cog/internal/ast"
)

type Builder struct {
	Package              string
	BuilderSignatureType string
	BuilderName          string
	ObjectName           string
	Imports              fmt.Stringer
	ImportAlias          string // alias to the pkg in which the object being built lives.
	Comments             []string
	Constructor          Constructor
	Properties           []ast.StructField
	Options              []Option
	Defaults             []OptionCall
}

type Constructor struct {
	Args        []ast.Argument
	Assignments []Assignment
}

type Option struct {
	Name        string
	Comments    []string
	Args        []ast.Argument
	Assignments []Assignment
}

type Assignment struct {
	Path           ast.Path
	InitSafeguards []string
	Constraints    []ast.AssignmentConstraint
	Method         ast.AssignmentMethod
	Value          ast.AssignmentValue
}

type OptionCall struct {
	OptionName string
	Args       []string
}
