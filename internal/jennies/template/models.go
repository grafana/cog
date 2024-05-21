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
	Constructor          ast.Constructor
	Properties           []ast.StructField
	Options              []ast.Option
	Defaults             []OptionCall
}

type OptionCall struct {
	OptionName string
	Args       []string
}
