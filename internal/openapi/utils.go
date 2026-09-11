package openapi

import (
	"errors"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/grafana/cog/pkg/ir"
)

func schemaComments(schema *openapi3.Schema) []string {
	lines := strings.Split(schema.Description, "\n")
	filtered := make([]string, 0, len(lines))

	for _, line := range lines {
		if line == "" {
			continue
		}

		filtered = append(filtered, line)
	}

	return filtered
}

func getEnumType(t string) (ir.Type, error) {
	switch t {
	case openapi3.TypeString:
		return ir.String(), nil
	case openapi3.TypeNumber:
		return ir.NewScalar(ir.KindInt32), nil
	case openapi3.TypeInteger:
		return ir.NewScalar(ir.KindInt64), nil
	default:
		return ir.Type{}, errors.New("only strings/numbers are supported")
	}
}

func getConstraints(schema *openapi3.Schema) []ir.TypeConstraint {
	constraints := make([]ir.TypeConstraint, 0)

	if schema.MinLength > 0 {
		constraints = append(constraints, ir.TypeConstraint{
			Op:   ir.MinLengthOp,
			Args: []any{schema.MinLength},
		})
	}
	if schema.MaxLength != nil {
		constraints = append(constraints, ir.TypeConstraint{
			Op:   ir.MaxLengthOp,
			Args: []any{*schema.MaxLength},
		})
	}

	if schema.MultipleOf != nil {
		constraints = append(constraints, ir.TypeConstraint{
			Op:   ir.MultipleOfOp,
			Args: getArgs(schema.MultipleOf, schema.Type.Slice()[0]),
		})
	}

	if schema.Min != nil || schema.ExclusiveMin.Value != nil {
		// In openapi3.0, schema.Min is used alongside a boolean exclusiveMinimum.
		// In 3.1, exclusiveMinimum is the actual value, rather than the boolean modifying minimum
		op := ir.GreaterThanEqualOp
		minVal := schema.Min
		if schema.ExclusiveMin.IsTrue() {
			op = ir.GreaterThanOp
		} else if schema.ExclusiveMin.Value != nil {
			// If the value is set, this is 3.1 style, and the value should be used as the min
			op = ir.GreaterThanOp
			minVal = schema.ExclusiveMin.Value
		}
		constraints = append(constraints, ir.TypeConstraint{
			Op:   op,
			Args: getArgs(minVal, schema.Type.Slice()[0]),
		})
	}

	if schema.Max != nil || schema.ExclusiveMax.Value != nil {
		// In openapi3.0, schema.Max is used alongside a boolean exclusiveMaximum.
		// In 3.1, exclusiveMaximum is the actual value, rather than the boolean modifying maximum
		op := ir.LessThanEqualOp
		maxVal := schema.Max
		if schema.ExclusiveMax.IsTrue() {
			op = ir.LessThanOp
		} else if schema.ExclusiveMax.Value != nil {
			// If the value is set, this is 3.1 style, and the value should be used as the max
			op = ir.LessThanOp
			maxVal = schema.ExclusiveMax.Value
		}
		constraints = append(constraints, ir.TypeConstraint{
			Op:   op,
			Args: getArgs(maxVal, schema.Type.Slice()[0]),
		})
	}

	return constraints
}

func getArgs(v *float64, t string) []any {
	args := []any{*v}
	if t == openapi3.TypeInteger {
		args = []any{int64(*v)}
	}
	return args
}

func isRef(ref string) bool {
	return ref != "" && strings.ContainsAny(ref, "#")
}
