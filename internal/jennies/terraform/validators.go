package terraform

import (
	"fmt"
	"strings"

	"github.com/grafana/cog/internal/ir"
	"github.com/grafana/cog/internal/languages"
)

type scalarValidator struct {
	importName string
	name       string
	minFunc    string
	maxFunc    string
	noneOfFunc string
	equalFunc  string
	regexFunc  string
}

type validators struct {
	context              languages.Context
	packageMapper        func(pkg string) string
	validatorDefinitions map[ir.ScalarKind]scalarValidator
}

func newValidators(context languages.Context, packageMapper func(pkg string) string) *validators {
	return &validators{
		context:       context,
		packageMapper: packageMapper,
		validatorDefinitions: map[ir.ScalarKind]scalarValidator{
			ir.KindString: {
				importName: "stringvalidator",
				name:       "String",
				minFunc:    "LengthAtLeast",
				maxFunc:    "LengthAtMost",
				noneOfFunc: "NoneOf",
				equalFunc:  "OneOf",
				regexFunc:  "RegexMatches",
			},
			ir.KindInt64: {
				importName: "int64validator",
				name:       "Int64",
				minFunc:    "AtLeast",
				maxFunc:    "AtMost",
				noneOfFunc: "NoneOf",
				equalFunc:  "OneOf",
			},
			ir.KindUint64: {
				importName: "int64validator",
				name:       "Int64",
				minFunc:    "AtLeast",
				maxFunc:    "AtMost",
				noneOfFunc: "NoneOf",
				equalFunc:  "OneOf",
			},
			ir.KindInt32: {
				importName: "int32validator",
				name:       "Int32",
				minFunc:    "AtLeast",
				maxFunc:    "AtMost",
				noneOfFunc: "NoneOf",
				equalFunc:  "OneOf",
			},
			ir.KindUint32: {
				importName: "int32validator",
				name:       "Int32",
				minFunc:    "AtLeast",
				maxFunc:    "AtMost",
				noneOfFunc: "NoneOf",
				equalFunc:  "OneOf",
			},
			ir.KindFloat32: {
				importName: "float32validator",
				name:       "Float32",
				minFunc:    "AtLeast",
				maxFunc:    "AtMost",
				noneOfFunc: "NoneOf",
				equalFunc:  "OneOf",
			},
			ir.KindFloat64: {
				importName: "float64validator",
				name:       "Float64",
				minFunc:    "AtLeast",
				maxFunc:    "AtMost",
				noneOfFunc: "NoneOf",
				equalFunc:  "OneOf",
			},
			ir.KindBool: {
				importName: "boolvalidator",
				name:       "Bool",
				equalFunc:  "Equal",
			},
			ir.KindInt8: {
				importName: "numbervalidator",
				name:       "Number",
				minFunc:    "AtLeastOneOf",
				noneOfFunc: "NoneOf",
				equalFunc:  "OneOf",
			},
			ir.KindUint8: {
				importName: "numbervalidator",
				name:       "Number",
				minFunc:    "AtLeastOneOf",
				noneOfFunc: "NoneOf",
				equalFunc:  "OneOf",
			},
			ir.KindInt16: {
				importName: "numbervalidator",
				name:       "Number",
				minFunc:    "AtLeastOneOf",
				noneOfFunc: "NoneOf",
				equalFunc:  "OneOf",
			},
			ir.KindUint16: {
				importName: "numbervalidator",
				name:       "Number",
				minFunc:    "AtLeastOneOf",
				noneOfFunc: "NoneOf",
				equalFunc:  "OneOf",
			},
		},
	}
}

func (v *validators) scalarValidator(kind ir.ScalarKind, constraints []ir.TypeConstraint) string {
	if len(constraints) == 0 {
		return ""
	}
	if validator, ok := v.validatorDefinitions[kind]; ok {
		constraintsStr := v.constraints(validator, constraints)
		if constraintsStr == "" {
			return ""
		}
		var buffer strings.Builder
		v.packageMapper("github.com/hashicorp/terraform-plugin-framework/schema/validator")
		v.packageMapper(fmt.Sprintf("github.com/hashicorp/terraform-plugin-framework-validators/%s", validator.importName))
		buffer.WriteString(fmt.Sprintf("[]validator.%s{\n", validator.name))
		buffer.WriteString(constraintsStr)
		buffer.WriteString("}")
		return buffer.String()
	}

	return ""
}

func (v *validators) constraints(validator scalarValidator, constraints []ir.TypeConstraint) string {
	var buffer strings.Builder
	for _, c := range constraints {
		args := make([]string, len(c.Args))
		for i, arg := range c.Args {
			args[i] = formatScalar(arg)
		}

		switch c.Op {
		case ir.MinLengthOp, ir.GreaterThanEqualOp, ir.GreaterThanOp:
			buffer.WriteString(fmt.Sprintf("%s.%s(%s),\n", validator.importName, validator.minFunc, v.calculateConstraint(c.Op, c.Args[0])))
		case ir.MaxLengthOp, ir.LessThanEqualOp, ir.LessThanOp:
			buffer.WriteString(fmt.Sprintf("%s.%s(%s),\n", validator.importName, validator.maxFunc, v.calculateConstraint(c.Op, c.Args[0])))
		case ir.NotEqualOp:
			buffer.WriteString(fmt.Sprintf("%s.%s(%+v),\n", validator.importName, validator.noneOfFunc, strings.Join(args, ", ")))
		case ir.EqualOp:
			buffer.WriteString(fmt.Sprintf("%s.%s(%+v),\n", validator.importName, validator.equalFunc, strings.Join(args, ", ")))
		case ir.RegexMatchOp:
			if validator.regexFunc != "" {
				v.packageMapper("regexp")
				buffer.WriteString(fmt.Sprintf("%s.%s(regexp.MustCompile(`%s`), \"\"),\n", validator.importName, validator.regexFunc, c.Args[0]))
			}
		}
	}

	return buffer.String()
}

func (v *validators) arrayConstraintValidator(constraints []ir.TypeConstraint) string {
	if len(constraints) == 0 {
		return ""
	}

	var buffer strings.Builder
	v.packageMapper("github.com/hashicorp/terraform-plugin-framework/schema/validator")
	v.packageMapper("github.com/hashicorp/terraform-plugin-framework-validators/listvalidator")
	buffer.WriteString("[]validator.List{\n")
	for _, c := range constraints {
		switch c.Op {
		case ir.MinItemsOp:
			buffer.WriteString(fmt.Sprintf("listvalidator.SizeAtLeast(%s),\n", formatScalar(c.Args[0])))
		case ir.MaxItemsOp:
			buffer.WriteString(fmt.Sprintf("listvalidator.SizeAtMost(%s),\n", formatScalar(c.Args[0])))
		case ir.UniqueItemsOp:
			buffer.WriteString("listvalidator.UniqueValues(),\n")
		}
	}
	buffer.WriteString("}")
	return buffer.String()
}

func (v *validators) validateList(def ir.Type) string {
	var buffer strings.Builder
	switch def.Kind {
	case ir.KindRef:
		obj, ok := v.context.LocateObject(def.AsRef().ReferredPkg, def.AsRef().ReferredType)
		if !ok {
			return "unknown validator"
		}

		return v.validateList(obj.Type)
	case ir.KindEnum:
		v.packageMapper("github.com/hashicorp/terraform-plugin-framework-validators/listvalidator")
		buffer.WriteString("[]validator.List{\n")
		validatorType := "ValueStringsAre"
		kind := def.AsEnum().Values[0].Type.AsScalar().ScalarKind
		if kind == ir.KindInt64 {
			validatorType = "ValueInt64sAre"
		}

		constraints := formatEnumValuesAsConstraints(def.AsEnum().Values)
		buffer.WriteString(fmt.Sprintf("listvalidator.%s(%s),\n}", validatorType, v.constraints(v.validatorDefinitions[kind], constraints)))
	default:
		return ""
	}

	return buffer.String()
}

func (v *validators) calculateConstraint(op ir.Op, arg any) string {
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
	case ir.GreaterThanOp:
		return fmt.Sprintf("%d", value+1)
	case ir.LessThanOp:
		return fmt.Sprintf("%d", value-1)
	default:
		return formatScalar(arg)
	}
}
