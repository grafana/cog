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

func (v *validators) scalarValidator(kind ast.ScalarKind, constraints []ast.TypeConstraint) string {
	if len(constraints) == 0 {
		return ""
	}
	if validator, ok := v.validatorDefinitions[kind]; ok {
		var buffer strings.Builder
		v.typeFormatter.packageMapper(fmt.Sprintf("github.com/hashicorp/terraform-plugin-framework-validators/%s", validator.importName))
		buffer.WriteString(fmt.Sprintf("[]validator.%s{\n", validator.name))
		buffer.WriteString(v.constraints(validator, constraints))
		buffer.WriteString("},\n")
		return buffer.String()
	}

	return ""
}

func (v *validators) constraints(validator scalarValidator, constraints []ast.TypeConstraint) string {
	var buffer strings.Builder
	for _, c := range constraints {
		args := make([]string, len(c.Args))
		for i, arg := range c.Args {
			args[i] = formatScalar(arg)
		}

		switch c.Op {
		case ast.MinLengthOp, ast.GreaterThanEqualOp, ast.GreaterThanOp:
			buffer.WriteString(fmt.Sprintf("%s.%s(%s),\n", validator.importName, validator.minFunc, v.calculateConstraint(c.Op, c.Args[0])))
		case ast.MaxLengthOp, ast.LessThanEqualOp, ast.LessThanOp:
			buffer.WriteString(fmt.Sprintf("%s.%s(%s),\n", validator.importName, validator.maxFunc, v.calculateConstraint(c.Op, c.Args[0])))
		case ast.NotEqualOp:
			buffer.WriteString(fmt.Sprintf("%s.%s(%+v),\n", validator.importName, validator.noneOfFunc, strings.Join(args, ", ")))
		case ast.EqualOp:
			buffer.WriteString(fmt.Sprintf("%s.%s(%+v),\n", validator.importName, validator.equalFunc, strings.Join(args, ", ")))
		default:
			fmt.Println("Unknown validator op", c.Op)
		}
	}

	return buffer.String()
}

func (v *validators) validateList(def ast.Type) string {
	var buffer strings.Builder
	switch def.Kind {
	case ast.KindRef:
		obj, ok := v.typeFormatter.context.LocateObject(def.AsRef().ReferredPkg, def.AsRef().ReferredType)
		if !ok {
			return "unknown validator"
		}

		return v.validateList(obj.Type)
	case ast.KindEnum:
		v.typeFormatter.packageMapper("github.com/hashicorp/terraform-plugin-framework-validators/listvalidator")
		buffer.WriteString("[]validator.List{\n")
		validatorType := "ValueStringsAre"
		kind := def.AsEnum().Values[0].Type.AsScalar().ScalarKind
		if kind == ast.KindInt64 {
			validatorType = "ValueInt64sAre"
		}

		constraints := formatEnumValuesAsConstraints(def.AsEnum().Values)
		buffer.WriteString(fmt.Sprintf("listvalidator.%s(%s),\n},\n", validatorType, v.constraints(v.validatorDefinitions[kind], constraints)))
	default:
		return ""
	}

	return buffer.String()
}

func (v *validators) calculateConstraint(op ast.Op, arg any) string {
	var value int64
	switch v := arg.(type) {
	case int64:
		value = v
	case int32:
		value = int64(v)
	case float32:
		value = int64(v)
	case float64:
		value = int64(v)
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
