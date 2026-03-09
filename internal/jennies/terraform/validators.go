package terraform

import (
	"fmt"
	"strings"

	"github.com/grafana/cog/internal/ast"
)

type scalarValidator struct {
	importName string
	name       string
	minFunc    string
	maxFunc    string
	noneOfFunc string
	equalFunc  string
}

type validators struct {
	typeFormatter        *typeFormatter
	validatorDefinitions map[ast.ScalarKind]scalarValidator
}

func newValidators(typeFormatter *typeFormatter) *validators {
	return &validators{
		typeFormatter: typeFormatter,
		validatorDefinitions: map[ast.ScalarKind]scalarValidator{
			ast.KindString: {
				importName: "stringvalidator",
				name:       "String",
				minFunc:    "LengthAtLeast",
				maxFunc:    "LengthAtMost",
				noneOfFunc: "NoneOf",
				equalFunc:  "OneOf",
			},
			ast.KindInt64: {
				importName: "int64validator",
				name:       "Int64",
				minFunc:    "AtLeast",
				maxFunc:    "AtMost",
				noneOfFunc: "NoneOf",
				equalFunc:  "OneOf",
			},
			ast.KindUint64: {
				importName: "int64validator",
				name:       "Int64",
				minFunc:    "AtLeast",
				maxFunc:    "AtMost",
				noneOfFunc: "NoneOf",
				equalFunc:  "OneOf",
			},
			ast.KindInt32: {
				importName: "int32validator",
				name:       "Int32",
				minFunc:    "AtLeast",
				maxFunc:    "AtMost",
				noneOfFunc: "NoneOf",
				equalFunc:  "OneOf",
			},
			ast.KindUint32: {
				importName: "int32validator",
				name:       "Int32",
				minFunc:    "AtLeast",
				maxFunc:    "AtMost",
				noneOfFunc: "NoneOf",
				equalFunc:  "OneOf",
			},
			ast.KindFloat32: {
				importName: "float32validator",
				name:       "Float32",
				minFunc:    "AtLeast",
				maxFunc:    "AtMost",
				noneOfFunc: "NoneOf",
				equalFunc:  "OneOf",
			},
			ast.KindFloat64: {
				importName: "float64validator",
				name:       "Float64",
				minFunc:    "AtLeast",
				maxFunc:    "AtMost",
				noneOfFunc: "NoneOf",
				equalFunc:  "OneOf",
			},
			ast.KindBool: {
				importName: "boolvalidator",
				name:       "Bool",
				equalFunc:  "Equal",
			},
			ast.KindInt8: {
				importName: "numbervalidator",
				name:       "Number",
				minFunc:    "AtLeastOneOf",
				noneOfFunc: "NoneOf",
				equalFunc:  "OneOf",
			},
			ast.KindUint8: {
				importName: "numbervalidator",
				name:       "Number",
				minFunc:    "AtLeastOneOf",
				noneOfFunc: "NoneOf",
				equalFunc:  "OneOf",
			},
			ast.KindInt16: {
				importName: "numbervalidator",
				name:       "Number",
				minFunc:    "AtLeastOneOf",
				noneOfFunc: "NoneOf",
				equalFunc:  "OneOf",
			},
			ast.KindUint16: {
				importName: "numbervalidator",
				name:       "Number",
				minFunc:    "AtLeastOneOf",
				noneOfFunc: "NoneOf",
				equalFunc:  "OneOf",
			},
		},
	}
}

func (v *validators) validate(kind ast.ScalarKind, constraints []ast.TypeConstraint) string {
	if len(constraints) == 0 {
		fmt.Println("A")
		return ""
	}
	if validator, ok := v.validatorDefinitions[kind]; ok {
		var buffer strings.Builder
		buffer.WriteString(fmt.Sprintf("[]validator.%s{\n", validator.name))

		v.typeFormatter.packageMapper(fmt.Sprintf("github.com/hashicorp/terraform-plugin-framework-validators/%s", validator.importName))
		for _, c := range constraints {
			switch c.Op {
			case ast.MinLengthOp, ast.GreaterThanEqualOp, ast.GreaterThanOp:
				buffer.WriteString(fmt.Sprintf("%s.%s(%s),\n", validator.importName, validator.minFunc, v.calculateConstraint(c.Op, c.Args[0])))
			case ast.MaxLengthOp, ast.LessThanEqualOp, ast.LessThanOp:
				buffer.WriteString(fmt.Sprintf("%s.%s(%s),\n", validator.importName, validator.maxFunc, v.calculateConstraint(c.Op, c.Args[0])))
			case ast.NotEqualOp:
				buffer.WriteString(fmt.Sprintf("%s.%s(%s),\n", validator.importName, validator.noneOfFunc, formatScalar(c.Args[0])))
			case ast.EqualOp:
				buffer.WriteString(fmt.Sprintf("%s.%s(%s),\n", validator.importName, validator.equalFunc, formatScalar(c.Args[0])))
			default:
				fmt.Println("Unknown validator op", c.Op)
			}
		}

		buffer.WriteString("},\n")
		return buffer.String()
	}

	return ""
}

func (v *validators) calculateConstraint(op ast.Op, arg any) string {
	var value int64
	switch arg.(type) {
	case int64:
		value = arg.(int64)
	case int32:
		value = int64(arg.(int32))
	case float32:
		value = int64(arg.(float32))
	case float64:
		value = int64(arg.(float64))
	}
	switch op {
	case ast.GreaterThanOp:
		return fmt.Sprintf("%d", value+1)
	case ast.LessThanOp:
		return fmt.Sprintf("%d", value-1)
	default:
		return formatScalar(arg)
	}
}
